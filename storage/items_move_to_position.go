package storage

import (
	"database/sql"
	"errors"
)

func MoveToPosition(id int, newPosition int) error {
	return WithStorage(func(storage *Storage) error {
		tx, err := storage.db.Begin()
		if err != nil {
			return err
		}

		if err = moveToPosition(tx, id, newPosition); err != nil {
			tx.Rollback()
			return err
		}

		return tx.Commit()
	})
}

func moveToPosition(tx *sql.Tx, id int, newPosition int) error {
	item, err := findById(tx, id)
	if err != nil {
		return err
	}

	if item.ListOrder == newPosition {
		// No movement needed
		return nil
	}

	isValid, err := isMovementValid(tx, item, newPosition)
	if err != nil {
		return err
	}

	oldPosition := item.ListOrder
	if isValid {
		if err = moveAllOthers(tx, oldPosition, newPosition); err != nil {
			return err
		}
		if err = setNewPosition(tx, item, newPosition); err != nil {
			return err
		}
	} else {
		return errors.New("Movement not valid.")
	}

	return nil
}

func isMovementValid(tx *sql.Tx, item *ItemDB, newPosition int) (bool, error) {

	if newPosition < 0 {
		return false, errors.New("Cannot move to a negative position")
	}

	if max, err := countAll(tx); err != nil {
		return false, err
	} else {
		if newPosition >= max {
			return false, errors.New("Cannot move to a position greater than the number of elements")
		}
	}

	itemAtPosition, err := getAtPosition(tx, newPosition)
	if err != nil {
		return false, err
	}

	return itemAtPosition.Checked == item.Checked, nil
}

func moveAllOthers(tx *sql.Tx, oldPosition int, newPosition int) error {
	if oldPosition > newPosition {
		// Moving up
		return moveItemsDown(tx, newPosition, oldPosition)
	} else {
		// Moving down
		return moveItemsUp(tx, oldPosition, newPosition)
	}
}

func moveItemsDown(tx *sql.Tx, newPosition int, oldPosition int) error {
	stmt, err := prepareStmt(tx, "UPDATE ITEMS SET listOrder = listOrder + 1 WHERE listOrder >= ? AND listOrder <= ?")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(newPosition, oldPosition)
	return err
}

func moveItemsUp(tx *sql.Tx, oldPosition int, newPosition int) error {
	stmt, err := prepareStmt(tx, "UPDATE ITEMS SET listOrder = listOrder - 1 WHERE listOrder >= ? AND listOrder <= ?")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(oldPosition, newPosition)
	return err
}

func setNewPosition(tx *sql.Tx, item *ItemDB, newPosition int) error {
	stmt, err := prepareStmt(tx, "UPDATE items SET listOrder = ? WHERE id = ?")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(newPosition, item.Id)
	return err
}
