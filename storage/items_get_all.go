package storage

import "database/sql"

func GetAll() ([]*ItemDB, error) {
	var allItems []*ItemDB
	err := WithStorage(func(s *Storage) error {
		tx, err := s.db.Begin()
		if err != nil {
			return err
		}

		items, err := retrieveAll(tx)
		if err != nil {
			tx.Rollback()
			return err
		}
		allItems = items
		return tx.Commit()
	})
	if err != nil {
		return nil, err
	}
	return allItems, nil
}

func retrieveAll(tx *sql.Tx) ([]*ItemDB, error) {
	res, err := tx.Query("SELECT * FROM items ORDER BY listOrder")
	if err != nil {
		return nil, err
	}

	defer res.Close()

	list := make([]*ItemDB, 0)
	for res.Next() {
		item, err := scanItem(res)
		if err != nil {
			return nil, err
		}
		list = append(list, item)
	}

	return list, nil
}
