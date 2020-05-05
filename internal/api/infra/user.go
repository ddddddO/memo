package infra

import (
	"database/sql"
	"errors"
	"log"

	_ "github.com/lib/pq"

	"github.com/ddddddO/tag-mng/internal/api/domain"
	"github.com/ddddddO/tag-mng/internal/api/domain/model"
)

type user struct {
	DB *sql.DB
}

func NewUser(db *sql.DB) domain.User {
	return user{
		DB: db,
	}
}

func (u user) FetchUser(name, passwd string) (*model.User, error) {
	const query = "SELECT id, name, passwd FROM users WHERE name=$1 AND passwd=$2"
	rows, err := u.DB.Query(query, name, passwd)
	if err != nil {
		return nil, err
	}

	// TODO: ユーザー登録をしてもらう or 正しいname/passwdを指定してもらう
	if !rows.Next() {
		return nil, errors.New("error !")
	}

	us := model.User{}
	if err := rows.Scan(&us.ID, &us.Name, &us.Passwd); err != nil {
		log.Println(err)
		return nil, err
	}

	return &us, nil
}
