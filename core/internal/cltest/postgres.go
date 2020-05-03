package cltest

import (
	"database/sql"
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"testing"

	"github.com/smartcontractkit/chainlink/core/logger"

	"github.com/smartcontractkit/chainlink/core/store/dbutil"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/store/orm"
)

// PrepareTestDB prepares the database to run tests, functionality varies
// on the underlying database.
// SQLite: No-op.
// Postgres: Creates a second database, and returns a cleanup callback
// that drops said DB.
func PrepareTestDB(tc *TestConfig) func() {
	t := tc.t
	t.Helper()

	originalURL := tc.DatabaseURL()
	if dbutil.IsPostgresURL(originalURL) {
		return createPostgresChildDB(tc, originalURL)
	}

	return func() {}
}

func createPostgresChildDB(tc *TestConfig, originalURL string) func() {
	t := tc.t

	parsed, err := url.Parse(originalURL)
	if err != nil {
		t.Fatalf("unable to extract database from %v: %v", originalURL, err)
	}

	testdb := createTestDB(t, parsed)
	tc.Set("DATABASE_URL", testdb.String())

	return func() {
		reapPostgresChildDB(t, parsed, testdb)
		tc.Set("DATABASE_URL", originalURL)
	}
}

func createTestDB(t testing.TB, parsed *url.URL) *url.URL {
	dbname := fmt.Sprintf("%s_%s_%s", parsed.Path[0:], models.NewID().String()[:8], t.Name())
	dbname = strings.Replace(dbname, "/", "_", -1)
	dbname = strings.Replace(dbname, "Test", "", 1)
	dbname = strings.ToLower(dbname)
	var safeCharsOnly = regexp.MustCompile(`[^0-9a-z_]+`)
	dbname = safeCharsOnly.ReplaceAllString(dbname, "")
	// NOTE: PostgreSQL's Max Identifier Length Is 63 Bytes
	if len(dbname) > 63 {
		dbname = dbname[:63]
	}
	db, err := sql.Open(string(orm.DialectPostgres), parsed.String())
	if err != nil {
		t.Fatalf("unable to open postgres database for creating test db: %+v", err)
	}
	defer logger.ErrorIfCalling(db.Close)

	//`CREATE DATABASE $1` does not seem to work w CREATE DATABASE
	_, err = db.Exec(fmt.Sprintf("CREATE DATABASE %s", dbname))
	if err != nil {
		t.Fatalf("unable to create postgres test database: %+v", err)
	}

	newURL := *parsed
	newURL.Path = "/" + dbname
	return &newURL
}

func reapPostgresChildDB(t testing.TB, parentURL, testURL *url.URL) {
	db, err := sql.Open(string(orm.DialectPostgres), parentURL.String())
	if err != nil {
		t.Fatalf("Unable to connect to parent CL db to clean up test database: %v", err)
	}
	defer logger.ErrorIfCalling(db.Close)

	testdb := testURL.Path[1:]
	dbsSQL := "DROP DATABASE " + testdb
	_, err = db.Exec(dbsSQL)
	if err != nil {
		t.Fatalf("Unable to clean up previous test database: %v", err)
	}
}
