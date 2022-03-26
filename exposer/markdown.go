package exposer

import (
	"fmt"
	"os"
	"time"

	"github.com/pkg/errors"
)

type markdown struct {
	f       *os.File
	title   string
	content string
}

func newMarkdown(path, subject, body string, createdAt, updatedAt time.Time) (*markdown, error) {
	f, err := os.OpenFile(path, os.O_RDWR, 0644)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return &markdown{
		f:       f,
		title:   title(subject),
		content: content(body, createdAt, updatedAt),
	}, nil
}

func title(subject string) string {
	return fmt.Sprintf("title: \"%s\"", subject)
}

func content(body string, createdAt, updatedAt time.Time) string {
	return header(createdAt, updatedAt) +
		"\n\n" +
		"---" +
		"\n\n" +
		body
}

const (
	layout         = "2006-1-2"
	headerTemplate = `
| 新規作成 | 最終更新 |
| -- | -- |
| %s | %s |
`
)

var jst = time.FixedZone("Asia/Tokyo", 9*60*60)

func header(createdAt, updatedAt time.Time) string {
	return fmt.Sprintf(
		headerTemplate,
		createdAt.In(jst).Format(layout),
		updatedAt.In(jst).Format(layout),
	)
}

func (m *markdown) write() error {
	// HUGOで生成したmdファイルに、titleへメモのsubjectを書き出すため(4バイト目から)
	_, err := m.f.WriteAt([]byte(m.title), 4)
	if err != nil {
		return errors.WithStack(err)
	}
	inf, err := m.f.Stat()
	if err != nil {
		return errors.WithStack(err)
	}

	// メモのcontentを追記するために、ファイルの最後尾から書き出す(inf.Size())
	_, err = m.f.WriteAt([]byte(m.content), inf.Size())
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (m *markdown) close() error {
	return m.f.Close()
}
