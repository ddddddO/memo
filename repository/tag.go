package repository

import (
	"github.com/ddddddO/memo/adapter"
)

// TODO: adapter -> modelsに置き換える
type TagRepository interface {
	FetchList(userID int) ([]adapter.Tag, error)
	FetchListByMemoID(memoID int) ([]adapter.Tag, error)
	Fetch(tagID int) (adapter.Tag, error)
	Update(tag adapter.Tag) error
	Delete(tag adapter.Tag) error
	Create(tag adapter.Tag) error
}
