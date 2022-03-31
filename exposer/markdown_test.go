package exposer

import (
	"testing"
	"time"
)

func TestHeader(t *testing.T) {
	want := `
| 新規作成 | 最終更新 |
| -- | -- |
| 2022-1-18 | 2022-3-31 |
`

	in := struct {
		createdAt time.Time
		updatedAt time.Time
	}{
		createdAt: time.Date(2022, 1, 18, 0, 0, 0, 0, time.UTC),
		updatedAt: time.Date(2022, 3, 31, 0, 0, 0, 0, time.UTC),
	}

	got := header(in.createdAt, in.updatedAt)
	if got != want {
		t.Errorf("want: %s, got: %s", want, got)
	}
}
