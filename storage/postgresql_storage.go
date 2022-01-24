package storage

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
)

func initPostgresqlStorage(dbUrl string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error establishing connection to PostgreSQL: %+v", err))
	}
	err = performPostgresqlSetup(db)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error performing PostgreSQL setup: %+v", err))
	}

	return db, nil
}

func performPostgresqlSetup(db *sql.DB) error {
	sqlStmt := `
	CREATE TABLE IF NOT EXISTS items (
		id SERIAL PRIMARY KEY,
		name TEXT NOT NULL,
		checked SMALLINT NOT NULL,
        listOrder INT NOT NULL,
        createdAt TEXT NOT NULL
    );
	`
	_, err := db.Exec(sqlStmt)
	if err != nil {
		return err
	}
	return nil
}
