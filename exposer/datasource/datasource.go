package datasource

import (
	"github.com/ddddddO/tag-mng/domain"
)

type DataSource interface {
	FetchAllExposedMemoSubjects() ([]string, error)
	FetchMemos() ([]domain.Memo, error)
	UpdateMemosExposedAt(memos []domain.Memo) error
}
