package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/sessions"
	_ "github.com/lib/pq"

	"github.com/ddddddO/memo/api/adapter"
)

type userUsecase interface {
	Fetch(name string, password string) (*adapter.User, error)
}

type authHandler struct {
	userUsecase userUsecase
	store       sessions.Store
}

func NewAuthHandler(uc userUsecase, store sessions.Store) *authHandler {
	return &authHandler{
		userUsecase: uc,
		store:       store,
	}
}

func (h *authHandler) Auth(w http.ResponseWriter, r *http.Request) {
	// TODO: name -> email へ変更したい
	name := r.PostFormValue("name")
	if len(name) == 0 {
		errResponse(w, http.StatusBadRequest, "empty key 'name'")
		return
	}

	password := r.PostFormValue("passwd")
	if len(password) == 0 {
		errResponse(w, http.StatusBadRequest, "empty key 'passwd'")
		return
	}

	user, err := h.userUsecase.Fetch(name, password)
	if err != nil {
		// TODO: error handling(dbに無かった or db接続エラー)
		errResponse(w, http.StatusUnauthorized, "failed")
		return
	}

	session, _ := h.store.New(r, "STORE")
	session.Values["authed"] = true
	if err := session.Save(r, w); err != nil {
		errResponse(w, http.StatusInternalServerError, "failed")
		return
	}

	res := struct {
		UserID int `json:"user_id"`
	}{
		UserID: user.ID,
	}
	if err := json.NewEncoder(w).Encode(res); err != nil {
		errResponse(w, http.StatusInternalServerError, "failed")
		return
	}
}
