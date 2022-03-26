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

func (*gce) filterRemoveMarkdowns(existingMarkdowns, newFileNames []string) []string {
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
		if err := g.removeFile(path); err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}

func (*gce) removeFile(path string) error {
	err := exec.Command("rm", path).Run()
	return errors.WithStack(err)
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
		err := g.removeFile(path)
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

	md, err := newMarkdown(fullFilePath, subject, body, createdAt, updatedAt)
	if err != nil {
		return errors.WithStack(err)
	}
	defer md.close()

	err = md.write()
	return errors.WithStack(err)
}

func (*gce) generateSites() error {
	err := exec.Command("hugo", "-D", "--cleanDestinationDir").Run()
	return errors.WithStack(err)
}

const (
	gcs = "gs://www.dododo.site"
)

func (*gce) uploadSites() error {
	err := exec.Command("gsutil", "-h", "Cache-Control:public, max-age=180", "rsync", "-d", "-r", "public", gcs).Run()
	return errors.WithStack(err)
}
