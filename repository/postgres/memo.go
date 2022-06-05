package postgres

import (
	"context"
	"database/sql"
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
	usersID := sql.NullInt64{
		Int64: int64(userID),
		Valid: true,
	}
	ctx := context.Background()
	return models.MemosByUsersID(ctx, pg.db, usersID)
}

func (pg *memoRepository) FetchListByTagID(userID, tagID int) ([]*models.Memo, error) {
	query, args, err := sq.Select("id, subject, content, users_id, created_at, updated_at, notified_cnt, is_exposed, exposed_at").Distinct().
		From("memos m").
		InnerJoin("memo_tag mt on mt.memos_id = m.id").
		Where(sq.Eq{"mt.tags_id": tagID}).PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := pg.db.Query(query, args...)
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
	return models.MemoByID(ctx, pg.db, memoID)
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

	const tagAll = 1
	deleteQuery, args, err := sq.Delete("memo_tag").Where(sq.And{sq.Eq{"memos_id": memo.ID}, sq.NotEq{"tags_id": tagAll}}).PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return err
	}

	_, err = tx.Exec(deleteQuery, args...)
	if err != nil {
		return err
	}

	insert := sq.Insert("memo_tag").Columns("memos_id", "tags_id")
	for _, tagID := range tagIDs {
		insert = insert.Values(
			memo.ID,
			tagID,
		)
	}
	insertQuery, args, err := insert.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return err
	}

	_, err = tx.Exec(insertQuery, args...)
	if err != nil {
		return err
	}

	return tx.Commit()
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

	insert := sq.Insert("memo_tag").Columns("memos_id", "tags_id")
	for _, tagID := range tagIDs {
		insert = insert.Values(
			memo.ID,
			tagID,
		)
	}
	insertQuery, args, err := insert.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return err
	}

	_, err = tx.Exec(insertQuery, args...)
	if err != nil {
		return err
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

	deleteQuery, args, err := sq.Delete("memo_tag").Where(sq.Eq{"memos_id": memo.ID}).PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return err
	}
	_, err = tx.Exec(deleteQuery, args...)
	if err != nil {
		return err
	}

	if err := memo.Delete(ctx, tx); err != nil {
		return err
	}

	return tx.Commit()
}
