package repository

import (
	"github.com/ddddddO/memo/models"
)

type MemoRepository interface {
	FetchList(userID int) ([]*models.Memo, error)
	FetchListByTagID(userID, tagID int) ([]*models.Memo, error)
	Fetch(memoID int) (*models.Memo, error)
	Update(memo *models.Memo, tagIDs []int) error
	Create(memo *models.Memo, tagIDs []int) error
	Delete(memoID int) error
}
