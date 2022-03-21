package repository

// TODO: adapter -> modelsに置き換える
import (
	"github.com/ddddddO/memo/adapter"
	"github.com/ddddddO/memo/models"
)

type MemoRepository interface {
	FetchList(userID int, tagID int) ([]*models.Memo, error)
	Fetch(memoID int) (*models.Memo, error)
	Update(memo *models.Memo, tagIDs []int) error
	Create(memo adapter.Memo) error
	Delete(memo adapter.Memo) error
}
