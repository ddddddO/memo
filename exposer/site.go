package exposer

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/ddddddO/memo/models"
)

// FIXME: gceだと具体すぎる気がする。
type gce struct {
	currentDir string
}

func newGCE(current string) *gce {
	return &gce{
		currentDir: current,
	}
}

func (g *gce) removeMarkdwonsNotIncludedInDB(subjects []string) ([]string, error) {
	if len(subjects) == 0 {
		return nil, nil
	}

	var newFileNames []string
	for _, s := range subjects {
		newFileNames = append(newFileNames, newFileName(s))
	}

	existingMarkdowns, err := g.searchExistingMarkdowns()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	removeMarkdowns := g.filterRemoveMarkdowns(existingMarkdowns, newFileNames)
	if len(removeMarkdowns) == 0 {
		return nil, nil
	}

	if err := g.removeContentFiles(removeMarkdowns); err != nil {
		return nil, errors.WithStack(err)
	}

	return removeMarkdowns, nil
}

const (
	invalidChars = "/"
	sep          = "_"
)

func newFileName(old string) string {
	if !strings.ContainsAny(old, invalidChars) {
		return old
	}

	new := old
	i := strings.IndexAny(new, invalidChars)
	for i != -1 {
		new = new[:i] + sep + new[i+1:]
		i = strings.IndexAny(new, invalidChars)
	}
	return new
}

func (g *gce) searchExistingMarkdowns() ([]string, error) {
	path := filepath.Join(g.currentDir, "content", "posts")
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var existingMarkdowns []string
	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".md") {
			existingMarkdowns = append(existingMarkdowns, strings.TrimSuffix(f.Name(), ".md"))
		}
	}
	return existingMarkdowns, nil
}

func (g *gce) filterRemoveMarkdowns(existingMarkdowns, newFileNames []string) []string {
	var removeMarkdowns []string
	for _, existing := range existingMarkdowns {
		isRemoving := false
		for _, new := range newFileNames {
			if existing == new {
				isRemoving = false
				break
			}
			isRemoving = true
		}
		if isRemoving {
			removeMarkdowns = append(removeMarkdowns, fmt.Sprintf("%s.md", existing))
		}
	}
	return removeMarkdowns
}

func (g *gce) removeContentFiles(fileNames []string) error {
	for _, f := range fileNames {
		path := filepath.Join(g.currentDir, "content", "posts", f)
		if err := exec.Command("rm", path).Run(); err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}

func (g *gce) removeExistingFiles(memos []*models.Memo) error {
	for _, m := range memos {
		fileName := fmt.Sprintf("%s.md", newFileName(m.Subject))
		fullFilePath := filepath.Join(g.currentDir, "content", "posts", fileName)

		// 既に同名のmdファイルが存在していた場合、hugo new fuga.mdは失敗する。なので、削除する。
		if err := g.removeExistingFile(fullFilePath); err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}

func (g *gce) removeExistingFile(path string) error {
	if exists(path) {
		err := exec.Command("rm", path).Run()
		if err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}

func exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func (g *gce) generateMarkdowns(memos []*models.Memo) error {
	for _, m := range memos {
		fileName := fmt.Sprintf("%s.md", newFileName(m.Subject))
		fullFilePath := filepath.Join(g.currentDir, "content", "posts", fileName)

		if err := g.generateMarkdown(fileName, fullFilePath, m.Subject, m.Content, m.CreatedAt.Time, m.UpdatedAt.Time); err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}

func (g *gce) generateMarkdown(fileName, fullFilePath, subject, body string, createdAt, updatedAt time.Time) error {
	// hugo new site hogehoge で生成したhogehogeディレクトリ内でhugo new fuga.md　しないと失敗する。
	err := exec.Command("hugo", "new", filepath.Join("posts", fileName)).Run()
	if err != nil {
		return errors.WithStack(err)
	}

	title := fmt.Sprintf("title: \"%s\"", subject)
	content := content(body, createdAt, updatedAt)
	md, err := newMarkdown(fullFilePath, title, content)
	if err != nil {
		return errors.WithStack(err)
	}
	defer md.close()

	err = md.write()
	return errors.WithStack(err)
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

type markdown struct {
	f       *os.File
	title   string
	content string
}

func newMarkdown(path, title, content string) (*markdown, error) {
	f, err := os.OpenFile(path, os.O_RDWR, 0644)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return &markdown{
		f:       f,
		title:   title,
		content: content,
	}, nil
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

func (g *gce) generateSites() error {
	err := exec.Command("hugo", "-D", "--cleanDestinationDir").Run()
	return errors.WithStack(err)
}

const (
	gcs = "gs://www.dododo.site"
)

func (g *gce) uploadSites() error {
	err := exec.Command("gsutil", "-h", "Cache-Control:public, max-age=180", "rsync", "-d", "-r", "public", gcs).Run()
	return errors.WithStack(err)
}
