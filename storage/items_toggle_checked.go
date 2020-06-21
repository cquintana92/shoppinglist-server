package storage

import (
	"database/sql"
	"errors"
)

func ToggleChecked(id int) error {
	return WithStorage(func(s *Storage) error {
		tx, err := s.db.Begin()
		if err != nil {
			return err
		}
		err = toggleChecked(tx, id)
		if err != nil {
			tx.Rollback()
			return err
		}
		return tx.Commit()
	})
}

func toggleChecked(tx *sql.Tx, id int) error {
	item, err := findById(tx, id)
	if err != nil {
		return err
	}
	if item.Checked == 1 {
		return setAsUnchecked(tx, item)
	} else {
		return setAsChecked(tx, item)
	}
}

func setAsChecked(tx *sql.Tx, item *ItemDB) error {
	previousPosition := item.ListOrder
	newPosition, err := getMaxChecked(tx)
	if err != nil {
		return err
	}
	if err := moveAllNextOneUp(tx, previousPosition); err != nil {
		return err
	}
	err = setCheckedAndPosition(tx, item.Id, newPosition, 1)
	return err
}

func moveAllNextOneUp(tx *sql.Tx, previousPosition int) error {
	stmt, err := tx.Prepare("UPDATE items SET listOrder = listOrder - 1 WHERE listOrder > ?")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(previousPosition)
	return err
}

func setCheckedAndPosition(tx *sql.Tx, id int, newPosition int, checked int) error {
	stmt, err := tx.Prepare("UPDATE items SET listOrder = ?, checked = ? WHERE id = ?")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(newPosition, checked, id)
	return err
}

func setAsUnchecked(tx *sql.Tx, item *ItemDB) error {
	previousPosition := item.ListOrder
	if previousPosition != 0 {
		newPosition, err := getMaxUnchecked(tx)
		if err != nil {
			return err
		}
		if err := moveAllNextOneDown(tx, newPosition, item.Id); err != nil {
			return err
		}
		err = setCheckedAndPosition(tx, item.Id, newPosition+1, 0)
		return err
	} else {
		err := setCheckedAndPosition(tx, item.Id, previousPosition, 0)
		return err
	}
}

func moveAllNextOneDown(tx *sql.Tx, previousPosition int, currentId int) error {
	stmt, err := tx.Prepare("UPDATE items SET listOrder = listOrder + 1 WHERE listOrder > ? AND id <> ?")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(previousPosition, currentId)
	return err
}

func getMaxChecked(tx *sql.Tx) (int, error) {
	res, err := tx.Query("SELECT MAX(listOrder) FROM items WHERE checked = 1")
	if err != nil {
		return 0, err
	}
	if res.Next() {
		var max sql.NullInt64
		err = res.Scan(&max)
		res.Close()
		// Check if returned number or null
		if max.Valid {
			return int(max.Int64), nil
		} else {
			// Returned null
			// No checked items
			// Use the max id
			res, err = tx.Query("SELECT MAX(listOrder) FROM items")
			if err != nil {
				return 0, err
			}
			if res.Next() {
				var max sql.NullInt64
				err = res.Scan(&max)
				res.Close()
				if max.Valid {
					return int(max.Int64), nil
				} else {
					return 0, nil
				}
			} else {
				return 0, errors.New("Query did not return any row")
			}
		}
	} else {
		return 0, errors.New("Query did not return any row")
	}
}

func getMaxUnchecked(tx *sql.Tx) (int, error) {
	res, err := tx.Query("SELECT MAX(listOrder) FROM items WHERE checked = 0")
	if err != nil {
		return 0, err
	}
	if res.Next() {
		var max sql.NullInt64
		err = res.Scan(&max)
		res.Close()
		if max.Valid {
			return int(max.Int64), nil
		} else {
			return -1, nil
		}
	} else {
		return 0, errors.New("Query did not return any row")
	}
}
