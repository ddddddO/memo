package repository

import (
	"github.com/ddddddO/memo/models"
)

type UserRepository interface {
	Fetch(name string, password string) (*models.User, error)
}
