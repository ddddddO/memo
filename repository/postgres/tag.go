package postgres

import (
	"context"
	"database/sql"

	sq "github.com/Masterminds/squirrel"
	_ "github.com/lib/pq"

	"github.com/ddddddO/memo/adapter"
	"github.com/ddddddO/memo/models"
)

type tagRepository struct {
	db *sql.DB
}

func NewTagRepository(db *sql.DB) *tagRepository {
	return &tagRepository{
		db: db,
	}
}

func (pg *tagRepository) FetchList(userID int) ([]*models.Tag, error) {
	ctx := context.Background()
	usersID := sql.NullInt64{
		Int64: int64(userID),
		Valid: true,
	}
	tags, err := models.TagsByUsersID(ctx, pg.db, usersID)
	if err != nil {
		return nil, err
	}

	return tags, nil
}

// FIXME:
func (pg *tagRepository) FetchListByMemoID(memoID int) ([]adapter.Tag, error) {
	var (
		rows *sql.Rows
		tags []adapter.Tag
		err  error
	)
	query := "SELECT id, name FROM tags WHERE id IN (SELECT tags_id FROM memo_tag WHERE memos_id=$1) ORDER BY id"
	rows, err = pg.db.Query(query, memoID)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var tag adapter.Tag
		if err := rows.Scan(&tag.ID, &tag.Name); err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}
	return tags, nil
}

func (pg *tagRepository) Fetch(tagID int) (*models.Tag, error) {
	ctx := context.Background()
	tag, err := models.TagByID(ctx, pg.db, tagID)
	if err != nil {
		return nil, err
	}
	return tag, nil
}

func (pg *tagRepository) Update(tag *models.Tag) error {
	tx, err := pg.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	ctx := context.Background()
	if err := tag.Update(ctx, tx); err != nil {
		return err
	}
	return tx.Commit()
}

func (pg *tagRepository) Delete(tagID int) error {
	tx, err := pg.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	ctx := context.Background()
	tag, err := models.TagByID(ctx, tx, tagID)
	if err != nil {
		return err
	}

	if err := tag.Delete(ctx, tx); err != nil {
		return err
	}

	// FIXME: using sq
	const deleteMemoTagQuery = `
	DELETE FROM memo_tag WHERE tags_id = $1
	`
	_, err = tx.Exec(deleteMemoTagQuery, tag.ID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (pg *tagRepository) Create(tag *models.Tag) error {
	tx, err := pg.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := sq.Insert("tags").
		Columns("name", "users_id").
		Values(tag.Name, tag.UsersID).
		Suffix("RETURNING \"id\"").
		RunWith(tx).
		PlaceholderFormat(sq.Dollar)

	if err := query.QueryRow().Scan(&tag.ID); err != nil {
		return err
	}

	// const createTagQuery = `
	// INSERT INTO tags(name, users_id) VALUES($1, $2) RETURNING id
	// `
	// result, err := pg.db.Exec(createTagQuery,
	// 	tag.Name, tag.UserID,
	// )
	// if err != nil {
	// 	return err
	// }
	// n, err := result.RowsAffected()
	// if err != nil {
	// 	return err
	// }
	// if n != 1 {
	// 	return errors.New("unexpected")
	// }
	return tx.Commit()
}
