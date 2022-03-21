package repository

import (
	"github.com/ddddddO/memo/adapter"
)

type MemoRepository interface {
	FetchList(userID int, tagID int) ([]adapter.Memo, error)
	Fetch(userID int, memoID int) (adapter.Memo, error)
	Update(memo adapter.Memo) error
	Create(memo adapter.Memo) error
	Delete(memo adapter.Memo) error
}
