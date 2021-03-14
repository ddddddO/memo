package postgres

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type healthPGRepository struct {
	db *sql.DB
}

func NewHealthPGRepository(db *sql.DB) *healthPGRepository {
	return &healthPGRepository{
		db: db,
	}
}

func (pg *healthPGRepository) Check() error {
	_, err := pg.db.Query("SELECT 1")
	if err != nil {
		return err
	}
	return nil
}
