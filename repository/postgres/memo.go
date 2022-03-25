package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	_ "github.com/lib/pq"

	"github.com/ddddddO/memo/models"
)

type memoRepository struct {
	db *sql.DB
}

func NewMemoRepository(db *sql.DB) *memoRepository {
	return &memoRepository{
		db: db,
	}
}

func (pg *memoRepository) FetchList(userID int) ([]*models.Memo, error) {
	var (
		memos []*models.Memo
		err   error
		ctx   = context.Background()
	)

	usersID := sql.NullInt64{
		Int64: int64(userID),
		Valid: true,
	}
	memos, err = models.MemosByUsersID(ctx, pg.db, usersID)
	if err != nil {
		return nil, err
	}

	return memos, nil
}

func (pg *memoRepository) FetchListByTagID(userID, tagID int) ([]*models.Memo, error) {
	// FIXME: using sq
	query := `select id, subject, content, users_id, created_at, updated_at, notified_cnt, is_exposed, exposed_at from memos where id in (select memos_id from memo_tag where tags_id = $1)`
	rows, err := pg.db.Query(query, tagID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var memos []*models.Memo
	for rows.Next() {
		m := models.Memo{}
		if err := rows.Scan(&m.ID, &m.Subject, &m.Content, &m.UsersID, &m.CreatedAt, &m.UpdatedAt, &m.NotifiedCnt, &m.IsExposed, &m.ExposedAt); err != nil {
			return nil, err
		}
		if m.UsersID.Int64 == int64(userID) {
			memos = append(memos, &m)
		}
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return memos, nil
}

func (pg *memoRepository) Fetch(memoID int) (*models.Memo, error) {
	ctx := context.Background()
	memo, err := models.MemoByID(ctx, pg.db, memoID)
	if err != nil {
		return nil, err
	}
	return memo, nil
}

// FIXME: memoRepositoryでmemo_tagの操作やめる
func (pg *memoRepository) Update(memo *models.Memo, tagIDs []int) error {
	tx, err := pg.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	ctx := context.Background()
	if err := memo.Update(ctx, tx); err != nil {
		return err
	}

	const deleteMemoTagQuery = `
	DELETE FROM memo_tag WHERE memos_id=$1 AND tags_id <> 1
	`

	_, err = tx.Exec(deleteMemoTagQuery, memo.ID)
	if err != nil {
		return err
	}

	var insertMemoTagQuery = `
	INSERT INTO memo_tag(memos_id, tags_id) VALUES
	%s
	`

	var values string
	for _, tid := range tagIDs {
		values += fmt.Sprintf("(%v, %d),", memo.ID, tid)
	}
	insertMemoTagQuery = fmt.Sprintf(insertMemoTagQuery, values[:len(values)-1])

	_, err = tx.Exec(insertMemoTagQuery)
	if err != nil {
		return err
	}

	tx.Commit()

	return nil
}

func (pg *memoRepository) UpdateExposedAt(memo *models.Memo) error {
	tx, err := pg.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	ctx := context.Background()
	memo.ExposedAt = sql.NullTime{Valid: true, Time: time.Now()}
	if err := memo.Update(ctx, tx); err != nil {
		return err
	}

	return tx.Commit()
}

// FIXME: memo_tag handling
func (pg *memoRepository) Create(memo *models.Memo, tagIDs []int) error {
	tx, err := pg.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := sq.Insert("memos").
		Columns("subject", "content", "users_id", "is_exposed").
		Values(memo.Subject, memo.Content, memo.UsersID, memo.IsExposed).
		Suffix("RETURNING \"id\"").
		RunWith(tx).
		PlaceholderFormat(sq.Dollar)

	if err := query.QueryRow().Scan(&memo.ID); err != nil {
		return err
	}

	// FIXME: using sq
	insertMemoTagQuery := "INSERT INTO memo_tag(memos_id, tags_id) VALUES($1, $2)"
	for _, tid := range tagIDs {
		_, err := tx.Exec(insertMemoTagQuery, memo.ID, tid)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (pg *memoRepository) Delete(memoID int) error {
	tx, err := pg.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	ctx := context.Background()
	memo, err := models.MemoByID(ctx, tx, memoID)
	if err != nil {
		return err
	}

	// FIXME: using sq
	const deleteMemoTagQuery = `
	DELETE FROM memo_tag WHERE memos_id = $1
	`
	_, err = tx.Exec(deleteMemoTagQuery, memo.ID)
	if err != nil {
		return err
	}

	if err := memo.Delete(ctx, tx); err != nil {
		return err
	}

	return tx.Commit()
}
