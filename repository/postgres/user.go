package postgres

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"log"
	"strings"

	_ "github.com/lib/pq"
	"github.com/pkg/errors"

	"github.com/ddddddO/tag-mng/domain"
)

type userPGRepository struct {
	db *sql.DB
}

func NewUserPGRepository(db *sql.DB) *userPGRepository {
	return &userPGRepository{
		db: db,
	}
}

func (pg *userPGRepository) Fetch(name string, password string) (*domain.User, error) {
	user, err := pg.fetchUser(name, genSecuredPassword(password, name))
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (pg *userPGRepository) fetchUser(name, password string) (*domain.User, error) {
	const query = "SELECT id, name, passwd FROM users WHERE name=$1 AND passwd=$2"
	rows, err := pg.db.Query(query, name, password)
	if err != nil {
		return nil, err
	}

	// TODO: ユーザー登録をしてもらう or 正しいname/passwordを指定してもらう
	if !rows.Next() {
		return nil, errors.New("error !")
	}

	user := domain.User{}
	if err := rows.Scan(&user.ID, &user.Name, &user.Password); err != nil {
		log.Println(err)
		return nil, err
	}

	return &user, nil
}

func genSecuredPassword(name, password string) string {
	secStrPass := name + password
	secPass := sha256.Sum256([]byte(secStrPass))
	for i := 0; i < 99999; i++ {
		secStrPass = hex.EncodeToString(secPass[:])
		secPass = sha256.Sum256([]byte(secStrPass))
	}
	return strings.ToLower(hex.EncodeToString(secPass[:]))
}
