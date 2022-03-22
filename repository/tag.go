package repository

import (
	"github.com/ddddddO/memo/models"
)

type TagRepository interface {
	FetchList(userID int) ([]*models.Tag, error)
	FetchListByMemoID(memoID int) ([]*models.Tag, error)
	Fetch(tagID int) (*models.Tag, error)
	Update(tag *models.Tag) error
	Delete(tagID int) error
	Create(tag *models.Tag) error
}
