package repository

import (
	"github.com/ddddddO/memo/domain"
)

type TagRepository interface {
	FetchList(userID int) ([]domain.Tag, error)
	Fetch(tagID int) (domain.Tag, error)
	Update(tag domain.Tag) error
	Delete(tag domain.Tag) error
	Create(tag domain.Tag) error
}
