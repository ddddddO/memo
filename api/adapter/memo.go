package adapter

import (
	"time"
)

// NOTE: adapterの役割りは、repositoryから取得したmodelsを、apiやbatchで使いやすくするために変換する用
//       api/batch用にそれぞれadapterを用意した方がよさそう

// Memo is ...
type Memo struct {
	ID          int        `json:"id"`
	Subject     string     `json:"subject"`
	Content     string     `json:"content"`
	IsExposed   bool       `json:"is_exposed"`
	UserID      int        `json:"user_id"`
	Tags        []Tag      `json:"tags"`
	NotifiedCnt int        `json:"notified_cnt"`
	CreatedAt   *time.Time `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
	ExposedAt   *time.Time `json:"exposed_at"`
	RowVariant  string     `json:"_rowVariant"` // for vue
}
