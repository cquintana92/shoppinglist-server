package storage

import (
	"database/sql"
	"errors"
)

func findById(tx *sql.Tx, id int) (*ItemDB, error) {
	stmt, err := prepareStmt(tx, "SELECT * FROM items WHERE id = ?")
	if err != nil {
		return nil, err
	}
	res, err := stmt.Query(id)
	if err != nil {
		return nil, err
	}

	defer res.Close()

	if res.Next() {
		item, err := scanItem(res)
		if err != nil {
			return nil, err
		} else {
			return item, nil
		}
	} else {
		return nil, errors.New("Query did not return any row")
	}
}

func getAtPosition(tx *sql.Tx, position int) (*ItemDB, error) {
	stmt, err := prepareStmt(tx, "SELECT * FROM items WHERE listOrder = ?")
	if err != nil {
		return nil, err
	}
	res, err := stmt.Query(position)
	defer res.Close()

	if err != nil {
		return nil, err
	}
	if res.Next() {
		item, err := scanItem(res)
		if err != nil {
			return nil, err
		} else {
			return item, nil
		}
	} else {
		return nil, errors.New("Query did not return any row")
	}
}

func countAll(tx *sql.Tx) (int, error) {
	res, err := tx.Query("SELECT COUNT(*) FROM items")
	if err != nil {
		return 0, err
	}

	defer res.Close()

	if res.Next() {
		var count sql.NullInt64
		if err = res.Scan(&count); err != nil {
			return 0, err
		}

		if count.Valid {
			return int(count.Int64), nil
		} else {
			return 0, nil
		}
	} else {
		return 0, errors.New("Query did not return any row")
	}
}

func isItemAddedAndUnchecked(tx *sql.Tx, name string) (bool, error) {
	stmt, err := prepareStmt(tx, "SELECT COUNT(*) FROM items WHERE name = ? AND checked = 0")
	if err != nil {
		return false, err
	}

	res, err := stmt.Query(name)
	if err != nil {
		return false, err
	}

	defer res.Close()

	if res.Next() {
		var count int
		if err = res.Scan(&count); err != nil {
			return false, err
		} else {
			return count > 0, nil
		}
	} else {
		return false, errors.New("Query did not return any row")
	}
}

func scanItem(rows *sql.Rows) (*ItemDB, error) {
	var id int
	var name string
	var checked int
	var listOrder int
	var createdAt string

	if err := rows.Scan(&id, &name, &checked, &listOrder, &createdAt); err != nil {
		return nil, err
	}
	return &ItemDB{id, name, checked, listOrder, createdAt}, nil
}
