package handlers

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/sessions"
	_ "github.com/lib/pq"
)

func AuthHandler(store sessions.Store) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//return func(w http.ResponseWriter, r *http.Request) {
		// TODO: name -> email へ変更したい
		name := r.PostFormValue("name")
		passwd := r.PostFormValue("passwd")

		// name, ok := c.GetPostForm("name")
		// if !ok {
		// 	// c.JSON(http.StatusBadRequest, gin.H{
		// 	// 	"message": "empty key 'name'",
		// 	// })
		// 	return
		// }
		// passwd, ok := c.GetPostForm("passwd")
		// if !ok {
		// 	// c.JSON(http.StatusBadRequest, gin.H{
		// 	// 	"message": "empty key 'passwd'",
		// 	// })
		// 	return
		// }

		// TODO: 共通化
		DBDSN := os.Getenv("DBDSN")
		if len(DBDSN) == 0 {
			log.Println("set default DSN")
			DBDSN = "host=localhost dbname=tag-mng user=postgres password=postgres sslmode=disable"
		}

		conn, err := sql.Open("postgres", DBDSN)
		if err != nil {
			// c.JSON(http.StatusInternalServerError, gin.H{
			// 	"message": "failed to connect db 1",
			// })
			return
		}

		const query = "SELECT id FROM users WHERE name=$1 AND passwd=$2"
		rows, err := conn.Query(query, name, genSecuredPasswd(passwd, name))
		if err != nil {
			// c.JSON(http.StatusInternalServerError, gin.H{
			// 	"message": "failed to connect db 2",
			// })
			return
		}

		// TODO: ユーザー登録をしてもらう or 正しいname/passwdを指定してもらう
		if !rows.Next() {
			// c.JSON(http.StatusUnauthorized, gin.H{
			// 	"message": "faild to authenticate",
			// })
			return
		}
		var userId int
		if err := rows.Scan(&userId); err != nil {
			log.Println(err)
			// c.JSON(http.StatusInternalServerError, gin.H{
			// 	"message": "failed to connect db 3",
			// })
			return
		}

		session, _ := store.New(r, "STORE")
		session.Values["authed"] = true
		if err := session.Save(r, w); err != nil {
			return
		}

		type response struct {
			UserID int `json:"user_id"`
		}
		res := response{
			UserID: userId,
		}

		resJson, err := json.Marshal(res)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(resJson))
	})
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
