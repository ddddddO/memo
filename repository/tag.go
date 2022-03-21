package repository

import (
	"github.com/ddddddO/memo/adapter"
	"github.com/ddddddO/memo/models"
)

// TODO: adapter -> modelsに置き換える
type TagRepository interface {
	FetchList(userID int) ([]*models.Tag, error)
	FetchListByMemoID(memoID int) ([]adapter.Tag, error)
	Fetch(tagID int) (*models.Tag, error)
	Update(tag *models.Tag) error
	Delete(tag adapter.Tag) error
	Create(tag adapter.Tag) error
}
