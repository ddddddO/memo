package repository

import (
	"github.com/ddddddO/memo/domain"
)

type UserRepository interface {
	Fetch(name string, password string) (*domain.User, error)
}
