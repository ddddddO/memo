package datasource

// TODO: apiと共通で使えるモデルとして移動する
type Memo struct {
	ID      int
	Subject string
	Content string
}

type DataSource interface {
	FetchAllExposedMemoSubjects() ([]string, error)
	FetchMemos() ([]Memo, error)
	UpdateMemosExposedAt(memos []Memo) error
}
