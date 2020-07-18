package exposer

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

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
	// TODO: where句
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
