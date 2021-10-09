package handler

import (
	"encoding/json"
	"net/http"

	_ "github.com/lib/pq"

	"github.com/ddddddO/tag-mng/domain"
)

type authUsecase interface {
	Login(name string, password string, w http.ResponseWriter, r *http.Request) (*domain.User, error)
}

type authHandler struct {
	usecase authUsecase
}

func NewAuth(usecase authUsecase) *authHandler {
	return &authHandler{
		usecase: usecase,
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

	// user, err := h.userRepo.Fetch(name, password)
	// if err != nil {
	// 	errResponse(w, http.StatusUnauthorized, "failed", err)
	// 	return
	// }

	// session, _ := h.store.New(r, "STORE")
	// session.Values["authed"] = true
	// if err := session.Save(r, w); err != nil {
	// 	errResponse(w, http.StatusInternalServerError, "failed", err)
	// 	return
	// }

	user, err := h.usecase.Login(name, password, w, r)
	if err != nil {
		// TODO: 認証エラーは4xx系エラーを返すようにする
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
