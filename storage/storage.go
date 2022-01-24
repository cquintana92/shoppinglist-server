package storage

import (
	"database/sql"
	"strings"
	"sync"
)

var (
	globalStorage *Storage
	mutex         sync.Mutex
)

type Storage struct {
	db *sql.DB
}

func InitStorage(dbUrl string) error {
	if strings.Contains(dbUrl, "sqlite") {
		return initSqliteStorage(dbUrl)
	} else {
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
