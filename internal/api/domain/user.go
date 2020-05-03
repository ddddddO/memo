package domain

import (
	"github.com/ddddddO/tag-mng/internal/api/domain/model"
)

type User interface {
	FetchUser(name, passwd string) (*model.User, error)
}
