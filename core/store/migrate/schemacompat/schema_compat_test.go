package schemacompat_test

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Rollback safety: a node may be rolled back to an older binary while the database keeps
// the columns added by newer migrations. Readers that SELECT * / RETURNING * into structs
// fail on unknown columns (sqlx strict mode), so new queries must list columns explicitly.
//
// legacyWildcardQueries is a ratchet of pre-existing offenders (non-test files under core/,
// counted per string literal). It may only shrink: new code must not add wildcard queries.
var legacyWildcardQueries = map[string]int{
	"bridges/orm.go":            9,
	"services/feeds/orm.go":     6,
	"services/job/orm.go":       22,
	"services/keystore/eth.go":  3,
	"services/keystore/orm.go":  2,
	"services/pipeline/orm.go":  10,
	"sessions/ldapauth/ldap.go": 7,
	"sessions/localauth/orm.go": 11,
	"sessions/oidcauth/oidc.go": 5,
}

var wildcardQuery = regexp.MustCompile(`(?is)(select\s+(\w+\.)?\*|returning\s+\*)`)

func countWildcardQueries(t *testing.T, root string) map[string]int {
	t.Helper()
	counts := map[string]int{}
	fset := token.NewFileSet()
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if d.Name() == "testdata" || path == filepath.Join(root, "internal") {
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") || d.Name() == "platform_arch_guard.go" {
			return nil
		}
		f, perr := parser.ParseFile(fset, path, nil, parser.SkipObjectResolution)
		if perr != nil {
			return perr
		}
		ast.Inspect(f, func(n ast.Node) bool {
			lit, ok := n.(*ast.BasicLit)
			if !ok || lit.Kind != token.STRING {
				return true
			}
			s, uerr := strconv.Unquote(lit.Value)
			if uerr == nil && wildcardQuery.MatchString(s) {
				rel, _ := filepath.Rel(root, path)
				counts[filepath.ToSlash(rel)]++
			}
			return true
		})
		return nil
	})
	require.NoError(t, err)
	return counts
}

func TestNoNewWildcardQueries(t *testing.T) {
	t.Parallel()
	actual := countWildcardQueries(t, filepath.Join("..", "..", ".."))

	for file, n := range actual {
		assert.LessOrEqualf(t, n, legacyWildcardQueries[file],
			"%s has %d SELECT */RETURNING * queries (allowed %d): list columns explicitly so older binaries tolerate newer schemas",
			file, n, legacyWildcardQueries[file])
	}
	for file, allowed := range legacyWildcardQueries {
		assert.LessOrEqualf(t, allowed, actual[file],
			"%s now has %d wildcard queries: lower its legacyWildcardQueries entry from %d", file, actual[file], allowed)
	}
}

// bridgeGuards blocks ADD COLUMN on tables whose readers in already-released binaries still
// use SELECT *. Raise the floor only once every release that reads the table with a wildcard
// has left the rollback window (check the release tags, not just the oldest supported one).
var bridgeGuards = map[string]int{
	// Every release through v2.72.0-rc.0 reads workflow_specs_v2 with SELECT *; 309 is their last migration.
	"workflow_specs_v2": 309,
}

const migrationsDir = "../migrations"

var migrationFile = regexp.MustCompile(`^(\d+)_.*\.sql$`)

func TestNoColumnAddsOnBridgeGuardedTables(t *testing.T) {
	t.Parallel()
	entries, err := os.ReadDir(migrationsDir)
	require.NoError(t, err)

	for _, e := range entries {
		m := migrationFile.FindStringSubmatch(e.Name())
		if m == nil {
			continue
		}
		num, _ := strconv.Atoi(m[1])
		body, rerr := fs.ReadFile(os.DirFS(migrationsDir), e.Name())
		require.NoError(t, rerr)
		up, _, _ := strings.Cut(string(body), "-- +goose Down")

		for table, floor := range bridgeGuards {
			if num <= floor {
				continue
			}
			addColumn := regexp.MustCompile(`(?is)alter\s+table\s+(if\s+exists\s+)?(only\s+)?("?\w+"?\.)?"?` + table + `"?\s[^;]*?add\s+column\b`)
			assert.Falsef(t, addColumn.MatchString(up),
				"%s adds a column to %s, which released binaries still read with SELECT *; a rollback would break. See bridgeGuards in schema_compat_test.go",
				e.Name(), table)
		}
	}
}
