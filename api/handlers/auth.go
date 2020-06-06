package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/sessions"
	_ "github.com/lib/pq"

	"github.com/ddddddO/tag-mng/api/usecase"
)

type AuthHandler interface {
	Login(store sessions.Store) http.Handler
}

type authHandler struct {
	userUseCase usecase.UserUseCase
}

func NewAuthHandler(uu usecase.UserUseCase) AuthHandler {
	return authHandler{
		userUseCase: uu,
	}
}

func (ah authHandler) Login(store sessions.Store) http.Handler {
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

		userID, err := ah.userUseCase.FetchUserID(name, passwd)
		if err != nil {
			errResponse(w, http.StatusInternalServerError, "failed", err)
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
