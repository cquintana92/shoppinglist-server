package storage

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_prepareStmtPostgres(t *testing.T) {
	res := prepareStmtPostgres("SELECT * FROM items WHERE a = ? AND b = ?")
	assert.Equal(t, "SELECT * FROM items WHERE a = $1 AND b = $2", res)
}
