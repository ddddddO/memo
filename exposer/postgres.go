package exposer

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"

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
	id      int
	subject string
	content string
}

func fetchMemos(db *sql.DB) ([]Memo, error) {
	const sql = `
	select id, subject, content from memos
	where (is_exposed = true and exposed_at is null)
	or (is_exposed = true and (exposed_at < updated_at))
	`

	rows, err := db.Query(sql)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var memos []Memo
	for rows.Next() {
		var memo Memo
		if err := rows.Scan(&memo.id, &memo.subject, &memo.content); err != nil {
			return nil, errors.WithStack(err)
		}
		memos = append(memos, memo)
	}
	return memos, nil
}

func updateMemosExposedAt(db *sql.DB, memos []Memo) error {
	var sql = `update memos set exposed_at = now() where id in (%s)`

	tmp := ""
	for _, m := range memos {
		id := strconv.Itoa(m.id)
		tmp += id + ","
	}
	tmp = tmp[:len(tmp)-1]

	_, err := db.Exec(fmt.Sprintf(sql, tmp))
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}
