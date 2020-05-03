package infra

import (
	"database/sql"
	"errors"
	"log"
	"os"

	_ "github.com/lib/pq"

	"github.com/ddddddO/tag-mng/internal/api/domain"
	"github.com/ddddddO/tag-mng/internal/api/domain/model"
)

type user struct{}

func NewUser() domain.User {
	return user{}
}

func (u user) FetchUser(name, passwd string) (*model.User, error) {
	// TODO: 共通化
	DBDSN := os.Getenv("DBDSN")
	if len(DBDSN) == 0 {
		log.Println("set default DSN")
		DBDSN = "host=localhost dbname=tag-mng user=postgres password=postgres sslmode=disable"
	}

	conn, err := sql.Open("postgres", DBDSN)
	if err != nil {
		return nil, err
	}

	const query = "SELECT id, name, passwd FROM users WHERE name=$1 AND passwd=$2"
	rows, err := conn.Query(query, name, passwd)
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

