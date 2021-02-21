package exposer

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	_ "github.com/lib/pq"
	"github.com/pkg/errors"

	"github.com/ddddddO/tag-mng/domain"
)

func removeMarkdwonsNotIncluded(subjects []string) ([]string, error) {
	if len(subjects) == 0 {
		return nil, nil
	}

	var convSubjects []string
	for _, subject := range subjects {
		convSubjects = append(convSubjects, cnvFileName(subject))
	}

	dir, err := os.Getwd()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	absDirPath := fmt.Sprintf("%s/content/posts/", dir)
	files, err := ioutil.ReadDir(absDirPath)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var fileNames []string
	for _, file := range files {
		name := file.Name()
		if strings.HasSuffix(name, ".md") {
			fileNames = append(fileNames, strings.TrimSuffix(name, ".md"))
		}
	}

	var removeMarkdowns []string
	for _, fileName := range fileNames {
		isRemoving := false
		for _, subject := range convSubjects {
			if fileName == subject {
				isRemoving = false
				break
			}
			isRemoving = true
		}
		if isRemoving {
			removeMarkdowns = append(removeMarkdowns, fmt.Sprintf("%s.md", fileName))
		}
	}

	if len(removeMarkdowns) == 0 {
		return nil, nil
	}

	for _, fileName := range removeMarkdowns {
		absFilePath := fmt.Sprintf("%s/content/posts/%s", dir, fileName)
		if err := exec.Command("rm", absFilePath).Run(); err != nil {
			return nil, errors.WithStack(err)
		}
	}

	return removeMarkdowns, nil
}

func generateMarkdowns(memos []domain.Memo) error {
	if len(memos) == 0 {
		return nil
	}

	for _, memo := range memos {
		// TODO: ここを並列処理でいけないか
		if err := generateMarkdown(memo); err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}

func generateMarkdown(memo domain.Memo) error {
	subject := cnvFileName(memo.Subject)
	fileName := fmt.Sprintf("%s.md", subject)

	dir, err := os.Getwd()
	if err != nil {
		return errors.WithStack(err)
	}
	absFilePath := fmt.Sprintf("%s/content/posts/%s", dir, fileName)

	// 既に同名のmdファイルが存在していた場合、hugo new fuga.mdは失敗する。なので、削除する。
	if exists(absFilePath) {
		err := exec.Command("rm", absFilePath).Run()
		if err != nil {
			return errors.WithStack(err)
		}
	}

	// hugo new site hogehoge で生成したhogehogeディレクトリ内でhugo new fuga.md　しないと失敗する。
	err = exec.Command("hugo", "new", fmt.Sprintf("posts/%s", fileName)).Run()
	if err != nil {
		return errors.WithStack(err)
	}

	f, err := os.OpenFile(absFilePath, os.O_RDWR, 0644)
	if err != nil {
		return errors.WithStack(err)
	}
	defer f.Close()
	// HUGOで生成したmdファイルに、titleへメモのsubjectを書き出すため(4バイト目から)
	title := `title: "` + memo.Subject + `"`
	_, err = f.WriteAt([]byte(title), 4)
	if err != nil {
		return errors.WithStack(err)
	}
	inf, err := f.Stat()
	if err != nil {
		return errors.WithStack(err)
	}
	// メモのcontentを追記するために、ファイルの最後尾から書き出す(inf.Size())
	_, err = f.WriteAt([]byte(memo.Content), inf.Size())
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func exists(filePath string) bool {
	_, err := os.Stat(filePath)
	return err == nil
}

const (
	invalidChars = "/"
	cnvChar      = "_"
)

func cnvFileName(fileName string) string {
	if !strings.ContainsAny(fileName, invalidChars) {
		return fileName
	}

	cnvFileName := fileName
	i := strings.IndexAny(cnvFileName, invalidChars)
	for i != -1 {
		cnvFileName = cnvFileName[:i] + cnvChar + cnvFileName[i+1:]
		i = strings.IndexAny(cnvFileName, invalidChars)
	}
	return cnvFileName
}

func generateSites() error {
	err := exec.Command("hugo", "-D", "--cleanDestinationDir").Run()
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func uploadSites() error {
	err := exec.Command("gsutil", "-h", "Cache-Control:public, max-age=180", "rsync", "-d", "-r", "public", "gs://www.dododo.site").Run()
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}
