package usecase

import (
	"net/http"

	"github.com/ddddddO/tag-mng/domain"
	"github.com/ddddddO/tag-mng/repository"
	"github.com/gorilla/sessions"

	"github.com/pkg/errors"
)

type authUsecase struct {
	userRepo repository.UserRepository
	store    sessions.Store
}

func NewAuth(userRepo repository.UserRepository, store sessions.Store) *authUsecase {
	return &authUsecase{
		userRepo: userRepo,
		store:    store,
	}
}

func (u *authUsecase) Login(name string, password string, w http.ResponseWriter, r *http.Request) (*domain.User, error) {
	user, err := u.userRepo.Fetch(name, password)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	session, _ := u.store.New(r, "STORE")
	session.Values["authed"] = true
	if err := session.Save(r, w); err != nil {
		return nil, errors.WithStack(err)
	}

	return user, nil
}
