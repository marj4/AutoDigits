package repository

import (
	"database/sql"

	_ "github.com/lib/pq"
)

func InitDB(connString string) (*sql.DB, error) {
	db, err := sql.Open("postgres", connString)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func Ping(db *sql.DB) error {
	if err := db.Ping(); err != nil {
		return err
	}
	return nil
}
