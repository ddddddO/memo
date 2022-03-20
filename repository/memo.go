package repository

import (
	"github.com/ddddddO/memo/domain"
)

type MemoRepository interface {
	FetchList(userID int, tagID int) ([]domain.Memo, error)
	Fetch(userID int, memoID int) (domain.Memo, error)
	Update(memo domain.Memo) error
	Create(memo domain.Memo) error
	Delete(memo domain.Memo) error
}
