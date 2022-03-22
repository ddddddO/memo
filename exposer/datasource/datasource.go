package datasource

import (
	"github.com/ddddddO/memo/domain"
)

type DataSource interface {
	FetchAllExposedMemoSubjects() ([]string, error)
	FetchMemos() ([]domain.Memo, error)
	UpdateMemosExposedAt(memos []domain.Memo) error
}
