package cltest

import (
	"database/sql"
	"fmt"
	"net/url"
	"strings"
	"testing"

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
	db, err := sql.Open(string(orm.DialectPostgres), originalURL)
	if err != nil {
		t.Fatalf("unable to open postgres database for creating test db: %+v", err)
	}
	defer db.Close()

	originalDB := extractDB(t, originalURL)
	dbname := fmt.Sprintf("%s_%s", originalDB, models.NewID().String())

	//`CREATE DATABASE $1` does not seem to work w CREATE DATABASE
	_, err = db.Exec(
		fmt.Sprintf("CREATE DATABASE %s", dbname),
	)
	if err != nil {
		t.Fatalf("unable to create postgres test database: %+v", err)
	}

	tc.Set("DATABASE_URL", swapDB(originalDB, originalURL, dbname))

	return func() {
		reapPostgresChildDB(t, originalURL, dbname)
		tc.Set("DATABASE_URL", originalURL)
	}
}

func reapPostgresChildDB(t testing.TB, dbURL, testdb string) {
	db, err := sql.Open(string(orm.DialectPostgres), dbURL)
	if err != nil {
		t.Fatalf("Unable to connect to parent CL db to clean up test database: %v", err)
	}
	defer db.Close()
	dbsSQL := "DROP DATABASE " + testdb
	_, err = db.Exec(dbsSQL)
	if err != nil {
		t.Fatalf("Unable to clean up previous test database: %v", err)
	}
}

func extractDB(t testing.TB, originalURL string) string {
	parsed, err := url.Parse(originalURL)
	if err != nil {
		t.Fatalf("unable to extract database from %v: %v", originalURL, err)
	}
	return parsed.Path[1:]
}

// swapDB uses replaces the DB part of the URL:
// postgres://localhost:5432/chainlink_test?sslmode=disable becomes
// postgres://localhost:5432/chainlink_test_4d63a0af83c34e348292189c0648a2af?sslmode=disable
func swapDB(originalDB, dburl, newdb string) string {
	return strings.Replace(dburl, originalDB, newdb, 1)
}
