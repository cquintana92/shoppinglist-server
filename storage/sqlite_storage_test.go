package storage

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_extractDbPath(t *testing.T) {
	cases := []struct {
		Case           string
		ExpectedString string
		ExpectedError  bool
	}{
		{Case: "", ExpectedString: "", ExpectedError: true},
		{Case: "sqlite3://something", ExpectedString: "something", ExpectedError: false},
		{Case: "sqlite3://./something.db", ExpectedString: "./something.db", ExpectedError: false},
	}

	for _, c := range cases {
		resString, resErr := extractDbPath(c.Case)
		assert.Equal(t, c.ExpectedString, resString)
		if c.ExpectedError {
			assert.Error(t, resErr)
		} else {
			assert.NoError(t, resErr)
		}
	}
}
