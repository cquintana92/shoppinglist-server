package utils

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_Replacements(t *testing.T) {
	source := "This is an example SeNTeNcE"
	replacements := []replacement{
		{from: "example", into: "replaced"},
		{from: "sentence", into: "phrase"},
	}

	res := performSanitize(source, replacements)
	require.Equal(t, "This is an replaced phrase", res)
}
