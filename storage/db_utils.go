package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

func prepareStmt(tx *sql.Tx, stmt string) (*sql.Stmt, error) {
	switch databaseMode {
	case dbModeSqlite:
		return tx.Prepare(stmt)
	case dbModePostgres:
		query := prepareStmtPostgres(stmt)
		return tx.Prepare(query)
	default:
		return nil, errors.New(fmt.Sprintf("invalid databaseMode: %d", databaseMode))
	}
}

func prepareStmtPostgres(stmt string) string {
	counter := 0
	newStmt := ""
	tokens := strings.Split(stmt, " ")
	for idx, token := range tokens {
		if idx > 0 {
			newStmt += " "
		}

		if strings.Contains(token, "?") {
			counter += 1
			newPart := fmt.Sprintf("$%d", counter)
			replaced := strings.Replace(token, "?", newPart, 1)
			newStmt += replaced
		} else {
			newStmt += token
		}
	}

	return newStmt
}
