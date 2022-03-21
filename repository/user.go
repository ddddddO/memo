package repository

import (
	"github.com/ddddddO/memo/adapter"
)

// TODO: adapter -> modelsに置き換える
type UserRepository interface {
	Fetch(name string, password string) (*adapter.User, error)
}
