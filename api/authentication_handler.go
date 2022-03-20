package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/sessions"
	_ "github.com/lib/pq"

	"github.com/ddddddO/memo/repository"
)

type authHandler struct {
	repo  repository.UserRepository
	store sessions.Store
}

func NewAuthHandler(repo repository.UserRepository, store sessions.Store) *authHandler {
	return &authHandler{
		repo:  repo,
		store: store,
	}
}

func (h *authHandler) Auth(w http.ResponseWriter, r *http.Request) {
	// TODO: name -> email へ変更したい
	name := r.PostFormValue("name")
	if len(name) == 0 {
		errResponse(w, http.StatusBadRequest, "empty key 'name'", nil)
		return
	}

	password := r.PostFormValue("passwd")
	if len(password) == 0 {
		errResponse(w, http.StatusBadRequest, "empty key 'passwd'", nil)
		return
	}

	user, err := h.repo.Fetch(name, password)
	if err != nil {
		errResponse(w, http.StatusUnauthorized, "failed", err)
		return
	}

	session, _ := h.store.New(r, "STORE")
	session.Values["authed"] = true
	if err := session.Save(r, w); err != nil {
		errResponse(w, http.StatusInternalServerError, "failed", err)
		return
	}

	res := struct {
		UserID int `json:"user_id"`
	}{
		UserID: user.ID,
	}
	if err := json.NewEncoder(w).Encode(res); err != nil {
		errResponse(w, http.StatusInternalServerError, "failed", err)
		return
	}
}
