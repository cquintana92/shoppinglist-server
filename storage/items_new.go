package storage

import (
	"database/sql"
	"errors"
	"shoppinglistserver/log"
	"shoppinglistserver/utils"
)

var (
	ItemAlreadyExistsError = errors.New("ItemAlreadyExists")
)

func New(name string) error {
	return WithStorage(func(s *Storage) error {
		tx, err := s.db.Begin()
		if err != nil {
			return err
		}

		if err = create(tx, name); err != nil {
			tx.Rollback()
			return err
		}
		return tx.Commit()
	})
}

func create(tx *sql.Tx, name string) error {
	alreadyAdded, err := isItemAddedAndUnchecked(tx, name)
	if err != nil {
		return err
	}
	if alreadyAdded {
		return ItemAlreadyExistsError
	}
	numUnchecked, err := getUncheckedCount(tx)
	if err != nil {
		return err
	}

	if err = createNew(tx, name, numUnchecked); err != nil {
		return err
	}
	if err = updateRest(tx); err != nil {
		return err
	}
	return nil

}

func createNew(tx *sql.Tx, name string, numUnchecked int) error {
	stmt, err := tx.Prepare("INSERT INTO items (name, checked, listOrder, createdAt) VALUES (?, ?, ?, ?)")
	if err != nil {
		return err
	}

	res, err := stmt.Exec(name, 0, numUnchecked, utils.Now())
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	log.Logger.Infof("Created new with ID=%d", id)
	return nil
}

func updateRest(tx *sql.Tx) error {
	_, err := tx.Exec("UPDATE items SET listOrder = listOrder + 1 WHERE checked = 1")
	return err
}

func getUncheckedCount(tx *sql.Tx) (int, error) {
	r, err := tx.Query("SELECT COUNT(*) AS total FROM items WHERE checked = 0")
	if err != nil {
		return 0, err
	}

	if r.Next() {
		var count int
		err = r.Scan(&count)
		if err != nil {
			return 0, err
		}
		r.Close()
		return count, nil
	} else {
		return 0, errors.New("Query returned 0 rows")
	}
}
