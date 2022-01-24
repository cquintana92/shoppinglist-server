package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"shoppinglistserver/log"
	"strings"
	"sync"
)

type dbMode int

const (
	dbModeSqlite dbMode = iota
	dbModePostgres
)

var (
	globalStorage *Storage
	databaseMode  dbMode
	mutex         sync.Mutex
)

type Storage struct {
	db *sql.DB
}

func InitStorage(dbUrl string) error {
	db, err := initDb(dbUrl)
	if err != nil {
		return errors.New(fmt.Sprintf("Error initializing db connection: %+v", err))
	}
	mutex.Lock()
	defer mutex.Unlock()
	globalStorage = &Storage{db: db}
	return nil
}

func initDb(dbUrl string) (*sql.DB, error) {
	if strings.Contains(dbUrl, "sqlite") {
		log.Logger.Infof("Using SQLITE backend")
		databaseMode = dbModeSqlite
		return initSqliteStorage(dbUrl)
	} else {
		databaseMode = dbModePostgres
		log.Logger.Infof("Using POSTGRESQL backend")
		return initPostgresqlStorage(dbUrl)
	}
}

func (s *Storage) Close() error {
	return s.db.Close()
}

func WithStorage(f func(*Storage) error) error {
	mutex.Lock()
	defer mutex.Unlock()
	return f(globalStorage)
}
