package postgres

import (
	"context"
	"database/sql"

	sq "github.com/Masterminds/squirrel"
	_ "github.com/lib/pq"

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
	return models.TagsByUsersID(ctx, pg.db, usersID)
}

func (pg *tagRepository) FetchListByMemoID(memoID int) ([]*models.Tag, error) {
	query, args, err := sq.Select("id, name").Distinct().
		From("tags t").
		InnerJoin("memo_tag mt on mt.tags_id = t.id").
		Where(sq.Eq{"mt.memos_id": memoID}).
		OrderBy("id").
		PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := pg.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tags := []*models.Tag{}
	for rows.Next() {
		var tag models.Tag
		if err := rows.Scan(&tag.ID, &tag.Name); err != nil {
			return nil, err
		}
		tags = append(tags, &tag)
	}
	return tags, nil
}

func (pg *tagRepository) Fetch(tagID int) (*models.Tag, error) {
	ctx := context.Background()
	return models.TagByID(ctx, pg.db, tagID)
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

	deleteQuery, args, err := sq.Delete("memo_tag").Where(sq.Eq{"tags_id": tag.ID}).PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return err
	}

	_, err = tx.Exec(deleteQuery, args...)
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
	return tx.Commit()
}
