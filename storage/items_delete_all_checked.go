package storage

import "database/sql"

func DeleteAllChecked() error {
	return WithStorage(func(s *Storage) error {
		tx, err := s.db.Begin()
		if err != nil {
			return err
		}
		if err = deleteAllChecked(tx); err != nil {
			tx.Rollback()
			return err
		}
		return tx.Commit()
	})
}

func deleteAllChecked(tx *sql.Tx) error {
	_, err := tx.Exec("DELETE FROM items WHERE checked = 1")
	return err
}
