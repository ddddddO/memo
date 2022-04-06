package usecase

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"

	"github.com/pkg/errors"

	"github.com/ddddddO/memo/api/adapter"
	"github.com/ddddddO/memo/models"
)

type userRepository interface {
	Fetch(name string, password string) (*models.User, error)
}

type userUsecase struct {
	repo userRepository
}

func NewUserUsecase(repo userRepository) *userUsecase {
	return &userUsecase{
		repo: repo,
	}
}

func (u userUsecase) Fetch(name string, password string) (*adapter.User, error) {
	user, err := u.repo.Fetch(name, genSecuredPassword(password, name))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	au := &adapter.User{
		ID:   user.ID,
		Name: user.Name,
	}
	return au, nil
}

// TODO: rubyで旧memoアプリ作ってた時の名残。ライブラリを使うようにする。
func genSecuredPassword(name, password string) string {
	secStrPass := name + password
	secPass := sha256.Sum256([]byte(secStrPass))
	for i := 0; i < 99999; i++ {
		secStrPass = hex.EncodeToString(secPass[:])
		secPass = sha256.Sum256([]byte(secStrPass))
	}
	return strings.ToLower(hex.EncodeToString(secPass[:]))
}
