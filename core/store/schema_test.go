package store

import (
	"bytes"
	_ "embed"
	"fmt"
	"os/exec"
	"strings"
	"testing"

	"github.com/kylelemons/godebug/diff"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/internal/testutils/configtest"
)

//go:embed schema.sql
var schema string

func TestSchema(t *testing.T) {
	cfg := configtest.NewTestGeneralConfig(t)
	dbURL := cfg.DatabaseURL()
	dumped, err := dbSchemaDump(dbURL.String())
	require.NoError(t, err)

	if dumped != schema {
		t.Errorf("schema.sql doesn't match database:\n\n%s", shortDiff(schema, dumped))
		t.Error("run `tools/bin/db_schema_dump` to regenerate schema.sql")
	}
}

func dbSchemaDump(dbURL string) (string, error) {
	cmd := exec.Command("../../tools/bin/db_schema_dump")
	cmd.Env = []string{"DATABASE_URL=" + dbURL}

	schema, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to dump schema: %v", err)
	}
	return string(schema), nil
}

// shortDiff is copied from diff.Diff, but skips printing equal lines.
func shortDiff(A, B string) string {
	aLines := strings.Split(A, "\n")
	bLines := strings.Split(B, "\n")

	chunks := diff.DiffChunks(aLines, bLines)

	buf := new(bytes.Buffer)
	for _, c := range chunks {
		for _, line := range c.Added {
			fmt.Fprintf(buf, "+%s\n", line)
		}
		for _, line := range c.Deleted {
			fmt.Fprintf(buf, "-%s\n", line)
		}
		// omit
		//for _, line := range c.Equal {
		//	fmt.Fprintf(buf, " %s\n", line)
		//}
	}
	return strings.TrimRight(buf.String(), "\n")
}
