package storage

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"shoppinglistserver/log"
	"sync"
)

var (
	globalStorage *Storage
	mutex         sync.Mutex
)

type Storage struct {
	db *sql.DB
}

func InitStorage(dbPath string) error {

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return err
	}

	if _, err := os.Stat(dbPath); err != nil {
		if os.IsNotExist(err) {
			log.Logger.Warnf("DB in [%s] did not exist. Creating it", dbPath)
			if err = performSetup(db); err != nil {
				return err
			}
		} else {
			return err
		}
	}

	mutex.Lock()
	defer mutex.Unlock()
	globalStorage = &Storage{db: db}
	return nil
}

func performSetup(db *sql.DB) error {
	sqlStmt := `
	CREATE TABLE items (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		checked TINYINT NOT NULL,
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

func (s *Storage) Close() error {
	return s.db.Close()
}

func WithStorage(f func(*Storage) error) error {
	mutex.Lock()
	defer mutex.Unlock()
	return f(globalStorage)
}
