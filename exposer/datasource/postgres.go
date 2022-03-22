package datasource

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"

	_ "github.com/lib/pq"
	"github.com/pkg/errors"

	"github.com/ddddddO/memo/domain"
)

type Postgres struct {
	db *sql.DB
}

func NewPostgres(dsn string) (*Postgres, error) {
	if dsn == "" {
		return nil, errors.New("undefined dsn")
	}

	log.Println("using dsn:", dsn)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &Postgres{
		db: db,
	}, nil
}

func (p *Postgres) FetchAllExposedMemoSubjects() ([]string, error) {
	const sql = `
	select subject from memos
	where is_exposed = true
	`

	rows, err := p.db.Query(sql)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var subjects []string
	for rows.Next() {
		var subject string
		if err := rows.Scan(&subject); err != nil {
			return nil, errors.WithStack(err)
		}
		subjects = append(subjects, subject)
	}
	return subjects, nil
}

func (p *Postgres) FetchMemos() ([]domain.Memo, error) {
	const sql = `
	select id, subject, content, created_at, updated_at from memos
	where (is_exposed = true and exposed_at is null)
	or (is_exposed = true and (exposed_at < updated_at))
	`

	rows, err := p.db.Query(sql)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var memos []domain.Memo
	for rows.Next() {
		var memo domain.Memo
		if err := rows.Scan(&memo.ID, &memo.Subject, &memo.Content, &memo.CreatedAt, &memo.UpdatedAt); err != nil {
			return nil, errors.WithStack(err)
		}
		memos = append(memos, memo)
	}
	return memos, nil
}

func (p *Postgres) UpdateMemosExposedAt(memos []domain.Memo) error {
	var sql = `update memos set exposed_at = now() where id in (%s)`

	tmp := ""
	for _, m := range memos {
		id := strconv.Itoa(m.ID)
		tmp += id + ","
	}
	tmp = tmp[:len(tmp)-1]

	_, err := p.db.Exec(fmt.Sprintf(sql, tmp))
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}
