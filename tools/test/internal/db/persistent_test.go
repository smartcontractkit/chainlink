package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPersistentTestDBLabelFilter(t *testing.T) {
	t.Parallel()
	got := PersistentTestDBLabelFilter()
	assert.Equal(t, "label=chainlink.tools.test.persistent-testdb=1", got)
}

func TestShellExportLine(t *testing.T) {
	t.Parallel()
	url := `postgres://u:p@h:1/db?sslmode=disable&x='y`
	assert.Equal(t,
		"export CL_DATABASE_URL='postgres://u:p@h:1/db?sslmode=disable&x='\\''y'\n",
		ShellExportLine(url, "bash"),
	)
	assert.Equal(t,
		`set -gx CL_DATABASE_URL "postgres://u:p@h:1/db?sslmode=disable&x='y"`+"\n",
		ShellExportLine(url, "fish"),
	)
}

func TestFishDoubleQuoteEscapes(t *testing.T) {
	t.Parallel()
	s := `a"b\c`
	assert.Equal(t, `"a\"b\\c"`, fishDoubleQuote(s))
}
