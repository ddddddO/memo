package domain

type Memo struct {
	ID          int      `json:"id"`
	Subject     string   `json:"subject"`
	Content     string   `json:"content"`
	IsExposed   bool     `json:"is_exposed"`
	TagIds      []int    `json:"tag_ids"`
	TagNames    []string `json:"tag_names"`
	NotifiedCnt int      `json:"notified_cnt"`
	RowVariant  string   `json:"_rowVariant"` // for vue
}
