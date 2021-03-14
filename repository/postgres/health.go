package postgres

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type healthRepository struct {
	db *sql.DB
}

func NewHealthRepository(db *sql.DB) *healthRepository {
	return &healthRepository{
		db: db,
	}
}

func (pg *healthRepository) Check() error {
	_, err := pg.db.Query("SELECT 1")
	if err != nil {
		return err
	}
	return nil
}
