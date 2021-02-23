package api

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/sessions"
	_ "github.com/lib/pq"

	"github.com/ddddddO/tag-mng/domain"
)

func fetchUserID(DB *sql.DB, name, passwd string) (int, error) {
	user, err := fetchUser(DB, name, genSecuredPasswd(passwd, name))
	if err != nil {
		return 0, err
	}

	return user.ID, nil
}

func fetchUser(DB *sql.DB, name, passwd string) (*domain.User, error) {
	const query = "SELECT id, name, passwd FROM users WHERE name=$1 AND passwd=$2"
	rows, err := DB.Query(query, name, passwd)
	if err != nil {
		return nil, err
	}

	// TODO: ユーザー登録をしてもらう or 正しいname/passwdを指定してもらう
	if !rows.Next() {
		return nil, errors.New("error !")
	}

	user := domain.User{}
	if err := rows.Scan(&user.ID, &user.Name, &user.Password); err != nil {
		log.Println(err)
		return nil, err
	}

	return &user, nil
}

func genSecuredPasswd(name, passwd string) string {
	secStrPass := name + passwd
	secPass := sha256.Sum256([]byte(secStrPass))
	for i := 0; i < 99999; i++ {
		secStrPass = hex.EncodeToString(secPass[:])
		secPass = sha256.Sum256([]byte(secStrPass))
	}
	return strings.ToLower(hex.EncodeToString(secPass[:]))
}

func NewAuthHandler(DB *sql.DB, store sessions.Store) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: name -> email へ変更したい
		name := r.PostFormValue("name")
		if len(name) == 0 {
			errResponse(w, http.StatusBadRequest, "empty key 'name'", nil)
			return
		}

		passwd := r.PostFormValue("passwd")
		if len(passwd) == 0 {
			errResponse(w, http.StatusBadRequest, "empty key 'passwd'", nil)
			return
		}

		userID, err := fetchUserID(DB, name, passwd)
		if err != nil {
			errResponse(w, http.StatusUnauthorized, "failed", err)
			return
		}

		session, _ := store.New(r, "STORE")
		session.Values["authed"] = true
		if err := session.Save(r, w); err != nil {
			errResponse(w, http.StatusInternalServerError, "failed", err)
			return
		}

		type response struct {
			UserID int `json:"user_id"`
		}
		res := response{
			UserID: userID,
		}

		resJson, err := json.Marshal(res)
		if err != nil {
			errResponse(w, http.StatusInternalServerError, "failed", err)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(resJson))
	})
}
