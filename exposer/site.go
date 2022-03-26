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
	for _, subject := range subjects {
		newFileNames = append(newFileNames, newFileName(subject))
	}

	dir, err := os.Getwd()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	fullDirPath := filepath.Join(dir, "content", "posts")
	files, err := ioutil.ReadDir(fullDirPath)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var existingFileNames []string
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".md") {
			existingFileNames = append(existingFileNames, strings.TrimSuffix(file.Name(), ".md"))
		}
	}

	var removeMarkdowns []string
	for _, existing := range existingFileNames {
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

	if len(removeMarkdowns) == 0 {
		return nil, nil
	}

	for _, fileName := range removeMarkdowns {
		fullFilePath := filepath.Join(dir, "content", "posts", fileName)
		if err := exec.Command("rm", fullFilePath).Run(); err != nil {
			return nil, errors.WithStack(err)
		}
	}

	return removeMarkdowns, nil
}

func generateMarkdowns(memos []*models.Memo) error {
	if len(memos) == 0 {
		return nil
	}

	for _, memo := range memos {
		if err := generateMarkdown(memo); err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}

func generateMarkdown(memo *models.Memo) error {
	subject := newFileName(memo.Subject)
	fileName := fmt.Sprintf("%s.md", subject)

	dir, err := os.Getwd()
	if err != nil {
		return errors.WithStack(err)
	}
	fullFilePath := filepath.Join(dir, "content", "posts", fileName)

	// 既に同名のmdファイルが存在していた場合、hugo new fuga.mdは失敗する。なので、削除する。
	if exists(fullFilePath) {
		err := exec.Command("rm", fullFilePath).Run()
		if err != nil {
			return errors.WithStack(err)
		}
	}

	// hugo new site hogehoge で生成したhogehogeディレクトリ内でhugo new fuga.md　しないと失敗する。
	err = exec.Command("hugo", "new", filepath.Join("posts", fileName)).Run()
	if err != nil {
		return errors.WithStack(err)
	}

	f, err := os.OpenFile(fullFilePath, os.O_RDWR, 0644)
	if err != nil {
		return errors.WithStack(err)
	}
	defer f.Close()
	// HUGOで生成したmdファイルに、titleへメモのsubjectを書き出すため(4バイト目から)
	title := fmt.Sprintf("title: \"%s\"", memo.Subject)
	_, err = f.WriteAt([]byte(title), 4)
	if err != nil {
		return errors.WithStack(err)
	}
	inf, err := f.Stat()
	if err != nil {
		return errors.WithStack(err)
	}

	// メモのcontentを追記するために、ファイルの最後尾から書き出す(inf.Size())
	_, err = f.WriteAt([]byte(content(memo)), inf.Size())
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func content(memo *models.Memo) string {
	return header(memo.CreatedAt.Time, memo.UpdatedAt.Time) +
		"\n\n" +
		"---" +
		"\n\n" +
		memo.Content
}

const (
	layout         = "2006-1-2"
	headerTemplate = `
| 新規作成 | 最終更新 |
| -- | -- |
| %s | %s |
`
)

func header(createdAt, updatedAt time.Time) string {
	jst := time.FixedZone("Asia/Tokyo", 9*60*60)
	return fmt.Sprintf(
		headerTemplate,
		createdAt.In(jst).Format(layout),
		updatedAt.In(jst).Format(layout),
	)
}

func exists(filePath string) bool {
	_, err := os.Stat(filePath)
	return err == nil
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
