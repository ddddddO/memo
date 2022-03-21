package postgres

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/pkg/errors"

	"github.com/ddddddO/memo/adapter"
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

func (pg *memoRepository) FetchList(userID int, tagID int) ([]*models.Memo, error) {
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

	ms := []*models.Memo{}
	// NOTE: tagIdが設定されている場合
	if tagID != -1 {
		for _, memo := range memos {
			memoTag := &models.MemoTag{
				MemosID: sql.NullInt64{
					Int64: int64(memo.ID),
					Valid: true,
				},
				TagsID: sql.NullInt64{
					Int64: int64(tagID),
					Valid: true,
				},
			}

			m, err := memoTag.Memo(ctx, pg.db)
			if err != nil {
				return nil, err
			}
			if m == nil {
				continue
			}

			t, err := memoTag.Tag(ctx, pg.db)
			if err != nil {
				return nil, err
			}
			if t == nil {
				continue
			}
			ms = append(ms, memo)
		}
	}
	if len(ms) != 0 {
		memos = ms
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

func (pg *memoRepository) Create(memo adapter.Memo) error {
	var createMemoQuery = `
	WITH inserted AS (INSERT INTO memos(subject, content, users_id, is_exposed) VALUES($1, $2, $3, $4) RETURNING id)
	INSERT INTO memo_tag(memos_id, tags_id) VALUES
		((SELECT id FROM inserted), 1)
		%s;
	`

	var values string
	for _, tag := range memo.Tags {
		values += fmt.Sprintf(",((SELECT id FROM inserted), %d)", tag.ID)
	}
	createMemoQuery = fmt.Sprintf(createMemoQuery, values)

	_, err := pg.db.Exec(createMemoQuery,
		memo.Subject, memo.Content, memo.UserID, memo.IsExposed,
	)
	if err != nil {
		return err
	}

	return nil
}

func (pg *memoRepository) Delete(memo adapter.Memo) error {
	const deleteMemoQuery = `
	DELETE FROM memos WHERE users_id = $1 AND id = $2;
	`

	tx, err := pg.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	result, err := tx.Exec(deleteMemoQuery,
		memo.UserID, memo.ID,
	)
	if err != nil {
		return err
	}
	n, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if n != 1 {
		return errors.New("unexpected")
	}

	const deleteMemoTagQuery = `
	DELETE FROM memo_tag WHERE memos_id = $1
	`

	_, err = tx.Exec(deleteMemoTagQuery, memo.ID)
	if err != nil {
		return err
	}

	tx.Commit()

	return nil
}
