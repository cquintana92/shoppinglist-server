package storage

import (
	"database/sql"
	"errors"
	"fmt"
)

func DeleteOne(id int) error {
	return WithStorage(func(storage *Storage) error {
		tx, err := storage.db.Begin()
		if err != nil {
			return err
		}

		err = deleteById(tx, id)
		if err != nil {
			tx.Rollback()
			return err
		}

		return tx.Commit()
	})
}

func deleteById(tx *sql.Tx, id int) error {
	stmt, err := prepareStmt(tx, "DELETE FROM items WHERE id = ?")
	if err != nil {
		return err
	}
	res, err := stmt.Exec(id)
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
