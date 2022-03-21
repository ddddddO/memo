package repository

import (
	"github.com/ddddddO/memo/adapter"
)

type UserRepository interface {
	Fetch(name string, password string) (*adapter.User, error)
}
