package usecase

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"

	"github.com/ddddddO/tag-mng/internal/api/domain"
)

type UserUseCase interface {
	FetchUserID(name, passwd string) (int, error)
}

type userUseCase struct {
	user domain.User
}

func NewUserUseCase(u domain.User) UserUseCase {
	return userUseCase{
		user: u,
	}
}

func (uu userUseCase) FetchUserID(name, passwd string) (int, error) {
	user, err := uu.user.FetchUser(name, genSecuredPasswd(passwd, name))
	if err != nil {
		return 0, err
	}

	return user.ID, nil
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
