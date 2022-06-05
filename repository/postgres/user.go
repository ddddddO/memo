package postgres

import (
	"database/sql"

	sq "github.com/Masterminds/squirrel"
	_ "github.com/lib/pq"

	"github.com/ddddddO/memo/models"
)

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *userRepository {
	return &userRepository{
		db: db,
	}
}

// TODO: nameとpasswdで複合index貼ってxoを再実行。生成された関数でUserを取得したい。
func (pg *userRepository) Fetch(name string, password string) (*models.User, error) {
	query, args, err := sq.Select("id, name, passwd").
		From("users").
		Where(sq.And{sq.Eq{"name": name}, sq.Eq{"passwd": password}}).PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, err
	}

	user := models.User{}
	err = pg.db.QueryRow(query, args...).Scan(&user.ID, &user.Name, &user.Passwd)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
