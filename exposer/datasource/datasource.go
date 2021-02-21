package datasource

import (
	"github.com/ddddddO/tag-mng/domain"
)

type DataSource interface {
	FetchAllExposedMemoSubjects() ([]string, error)
	FetchMemos() ([]domain.MemoDetail, error)
	UpdateMemosExposedAt(memos []domain.MemoDetail) error
}
