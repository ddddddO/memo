package postgres

import (
	"database/sql"

	_ "github.com/lib/pq"
	"github.com/pkg/errors"

	"github.com/ddddddO/tag-mng/domain"
)

type tagRepository struct {
	db *sql.DB
}

func NewTagRepository(db *sql.DB) *tagRepository {
	return &tagRepository{
		db: db,
	}
}

func (pg *tagRepository) FetchList(userID int) ([]domain.Tag, error) {
	var (
		rows *sql.Rows
		tags []domain.Tag
		err  error
	)
	query := "SELECT id, name FROM tags WHERE users_id = $1 ORDER BY id"
	rows, err = pg.db.Query(query, userID)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var tag domain.Tag
		if err := rows.Scan(&tag.ID, &tag.Name); err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}
	return tags, nil
}

func (pg *tagRepository) Fetch(tagID int) (domain.Tag, error) {
	var tag domain.Tag
	query := "SELECT id, name FROM tags WHERE id = $1"
	if err := pg.db.QueryRow(query, tagID).Scan(&tag.ID, &tag.Name); err != nil {
		return domain.Tag{}, err
	}
	return tag, nil
}

func (pg *tagRepository) Update(tag domain.Tag) error {
	const updateTagQuery = `
	UPDATE tags SET name = $1 WHERE id = $2
	`
	result, err := pg.db.Exec(updateTagQuery,
		tag.Name, tag.ID,
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
	return nil
}

func (pg *tagRepository) Delete(tag domain.Tag) error {
	const deleteTagQuery = `
	DELETE FROM tags WHERE id = $1
	`

	tx, err := pg.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	result, err := tx.Exec(deleteTagQuery, tag.ID)
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
	DELETE FROM memo_tag WHERE tags_id = $1
	`

	_, err = tx.Exec(deleteMemoTagQuery, tag.ID)
	if err != nil {
		return err
	}

	tx.Commit()

	return nil
}

func (pg *tagRepository) Create(tag domain.Tag) error {
	const createTagQuery = `
	INSERT INTO tags(name, users_id) VALUES($1, $2) RETURNING id
	`
	result, err := pg.db.Exec(createTagQuery,
		tag.Name, tag.UserID,
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
	return nil
}
