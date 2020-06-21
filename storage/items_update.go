package storage

import (
	"database/sql"
	"errors"
	"fmt"
)

func Update(name string, id int) error {
	return WithStorage(func(storage *Storage) error {
		tx, err := storage.db.Begin()
		if err != nil {
			return err
		}

		err = setName(tx, name, id)
		if err != nil {
			tx.Rollback()
			return err
		}

		return tx.Commit()
	})
}

func setName(tx *sql.Tx, name string, id int) error {
	stmt, err := tx.Prepare("UPDATE items SET name = ? WHERE id = ?")
	if err != nil {
		return err
	}
	res, err := stmt.Exec(name, id)
	if err != nil {
		return err
	}
	if rows, err := res.RowsAffected(); err != nil {
		return err
	} else {
		if rows != 1 {
			return errors.New(fmt.Sprintf("Should have only affected 1 row, but affected %d", rows))
		} else {
			return nil
		}
	}
}
