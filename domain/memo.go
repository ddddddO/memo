package domain

// Memo is ...
type Memo struct {
	ID          int    `json:"id"`
	Subject     string `json:"subject"`
	Content     string `json:"content"`
	IsExposed   bool   `json:"is_exposed"`
	UserID      int    `json:"user_id"`
	Tags        []Tag  `json:"tags"`
	NotifiedCnt int    `json:"notified_cnt"`
	RowVariant  string `json:"_rowVariant"` // for vue
}
