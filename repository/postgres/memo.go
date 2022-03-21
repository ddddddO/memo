package postgres

import (
	"context"
	"database/sql"
	"fmt"

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
