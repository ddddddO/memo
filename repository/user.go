package repository

import (
	"github.com/ddddddO/tag-mng/domain"
)

type UserRepository interface {
	Fetch(name string, password string) (*domain.User, error)
}
