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

func removeMarkdwonsNotIncludedInDB(subjects []string) ([]string, error) {
	if len(subjects) == 0 {
		return nil, nil
	}

	var newFileNames []string
	for _, s := range subjects {
		newFileNames = append(newFileNames, newFileName(s))
	}

	dir, err := os.Getwd()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	existingMarkdowns, err := seachExistingMarkdowns(dir)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	removeMarkdowns := filterRemoveMarkdowns(existingMarkdowns, newFileNames)
	if len(removeMarkdowns) == 0 {
		return nil, nil
	}

	if err := removeContentFiles(dir, removeMarkdowns); err != nil {
		return nil, errors.WithStack(err)
	}

	return removeMarkdowns, nil
}

func seachExistingMarkdowns(dir string) ([]string, error) {
	path := filepath.Join(dir, "content", "posts")
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

func filterRemoveMarkdowns(existingMarkdowns, newFileNames []string) []string {
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

func removeContentFiles(dir string, fileNames []string) error {
	for _, f := range fileNames {
		path := filepath.Join(dir, "content", "posts", f)
		if err := exec.Command("rm", path).Run(); err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}

func removeExistingFiles(memos []*models.Memo) error {
	dir, err := os.Getwd()
	if err != nil {
		return errors.WithStack(err)
	}

	for _, m := range memos {
		fileName := fmt.Sprintf("%s.md", newFileName(m.Subject))
		fullFilePath := filepath.Join(dir, "content", "posts", fileName)

		// 既に同名のmdファイルが存在していた場合、hugo new fuga.mdは失敗する。なので、削除する。
		if err := removeExistingFile(fullFilePath); err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}

func removeExistingFile(path string) error {
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

func generateMarkdowns(memos []*models.Memo) error {
	dir, err := os.Getwd()
	if err != nil {
		return errors.WithStack(err)
	}

	for _, m := range memos {
		fileName := fmt.Sprintf("%s.md", newFileName(m.Subject))
		fullFilePath := filepath.Join(dir, "content", "posts", fileName)

		if err := generateMarkdown(fileName, fullFilePath, m.Subject, m.Content, m.CreatedAt.Time, m.UpdatedAt.Time); err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}

func generateMarkdown(fileName, fullFilePath, subject, body string, createdAt, updatedAt time.Time) error {
	// hugo new site hogehoge で生成したhogehogeディレクトリ内でhugo new fuga.md　しないと失敗する。
	err := exec.Command("hugo", "new", filepath.Join("posts", fileName)).Run()
	if err != nil {
		return errors.WithStack(err)
	}

	f, err := os.OpenFile(fullFilePath, os.O_RDWR, 0644)
	if err != nil {
		return errors.WithStack(err)
	}
	defer f.Close()

	title := fmt.Sprintf("title: \"%s\"", subject)
	// HUGOで生成したmdファイルに、titleへメモのsubjectを書き出すため(4バイト目から)
	_, err = f.WriteAt([]byte(title), 4)
	if err != nil {
		return errors.WithStack(err)
	}
	inf, err := f.Stat()
	if err != nil {
		return errors.WithStack(err)
	}

	// メモのcontentを追記するために、ファイルの最後尾から書き出す(inf.Size())
	_, err = f.WriteAt([]byte(content(body, createdAt, updatedAt)), inf.Size())
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

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

const (
	invalidChars = "/"
	sep          = "_"
)

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

func generateSites() error {
	err := exec.Command("hugo", "-D", "--cleanDestinationDir").Run()
	return errors.WithStack(err)
}

const (
	gcs = "gs://www.dododo.site"
)

func uploadSites() error {
	err := exec.Command("gsutil", "-h", "Cache-Control:public, max-age=180", "rsync", "-d", "-r", "public", gcs).Run()
	return errors.WithStack(err)
}
