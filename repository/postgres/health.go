package postgres

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type HealthPGRepository struct {
	db *sql.DB
}

func NewHealthPGRepository(db *sql.DB) *HealthPGRepository {
	return &HealthPGRepository{
		db: db,
	}
}

func (pg *HealthPGRepository) Check() error {
	_, err := pg.db.Query("SELECT 1")
	if err != nil {
		return err
	}
	return nil
}
