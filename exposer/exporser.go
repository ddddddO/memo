package exposer

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

func Run(dsn string) error {
	db, err := genDB(dsn)
	if err != nil {
		return errors.Wrap(err, "generate db connection error")
	}

	memos, err := fetchMemos(db)
	if err != nil {
		return errors.Wrap(err, "db error")
	}

	if err := genMDs(memos); err != nil {
		return errors.Wrap(err, "generate md file error")
	}

	if err := genSite(); err != nil {
		return errors.Wrap(err, "generate html error")
	}

	// if err := uploadSite(); err != nil {
	// 	return errors.Wrap(err, "upload site error")
	// }

	log.Println("succeeded")
	return nil
}

func genDB(dsn string) (*sql.DB, error) {
	if dsn == "" {
		return nil, errors.New("undefined dsn")
	}

	log.Println("using dsn:", dsn)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return db, nil
}

type Memo struct {
	subject string
	content string
}

func fetchMemos(db *sql.DB) ([]Memo, error) {
	const sql = `select subject, content from memos where id = 45`

	rows, err := db.Query(sql)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var memos []Memo
	for rows.Next() {
		var memo Memo
		if err := rows.Scan(&memo.subject, &memo.content); err != nil {
			return nil, errors.WithStack(err)
		}
		memos = append(memos, memo)
	}

	return memos, nil
}

func genMDs(memos []Memo) error {
	for _, memo := range memos {
		// TODO: ここを並列処理でいけないか
		if err := genMD(memo); err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}

func genMD(memo Memo) error {
	subject := cnvFileName(memo.subject)
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
	title := `title: "` + memo.subject + `"`
	_, err = f.WriteAt([]byte(title), 4)
	if err != nil {
		return errors.WithStack(err)
	}
	inf, err := f.Stat()
	if err != nil {
		return errors.WithStack(err)
	}
	// メモのcontentを追記するために、ファイルの最後尾から書き出す(inf.Size())
	_, err = f.WriteAt([]byte(memo.content), inf.Size())
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

func genSite() error {
	err := exec.Command("hugo", "-D").Run()
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func uploadSite() error {
	err := exec.Command("gsutil", "rsync", "-R", "public", "gs://www.dododo.site").Run()
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}
