package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"shoppinglistserver/log"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

func initSqliteStorage(dbUrl string) (*sql.DB, error) {
	dbPath, err := extractDbPath(dbUrl)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error extracting dbPath: %+v", err))
	}
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error opening sqlite database: %+v", err))
	}

	if _, err := os.Stat(dbPath); err != nil {
		if os.IsNotExist(err) {
			log.Logger.Warnf("DB in [%s] did not exist. Creating it", dbPath)
			if err = performSqliteSetup(db); err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	} else {
		log.Logger.Info("DB already existed")
	}

	return db, nil
}

func extractDbPath(dbUrl string) (string, error) {
	parts := strings.Split(dbUrl, "://")
	if len(parts) != 2 {
		return "", errors.New(fmt.Sprintf("Invalid sqlite url: %s", dbUrl))
	}

	return parts[1], nil
}

func performSqliteSetup(db *sql.DB) error {
	exists, err := checkIfTableExists(db)
	if err != nil {
		return errors.New(fmt.Sprintf("Error in startup: %+v", err))
	}
	if exists {
		return nil
	}
	sqlStmt := `
	CREATE TABLE items (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		checked TINYINT NOT NULL,
        listOrder INT NOT NULL,
        createdAt TEXT NOT NULL
    );
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		return err
	}
	return nil
}

func checkIfTableExists(db *sql.DB) (bool, error) {
	q := "SELECT name FROM sqlite_master WHERE type='table' AND name='items'"
	res, err := db.Query(q)
	if err != nil {
		return false, errors.New(fmt.Sprintf("Error checking if items table exists: %+v", err))
	}
	defer res.Close()
	if res.Next() {
		return true, nil
	}
	return false, nil

}
