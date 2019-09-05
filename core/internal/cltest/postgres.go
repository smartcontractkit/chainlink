package cltest

import (
	"database/sql"
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/store/orm"
)

func CreatePostgresDatabase(t testing.TB, config *TestConfig) func() {
	t.Helper()

	originalURL := config.DatabaseURL()
	if strings.HasPrefix(strings.ToLower(originalURL), string(orm.DialectPostgres)) {
		db, err := gorm.Open(string(orm.DialectPostgres), originalURL)
		if err != nil {
			t.Fatalf("unable to open postgres database for creating test db: %+v", err)
		}
		defer db.Close()

		dbname := fmt.Sprintf("chainlink_test_%s", models.NewID().String())

		//`CREATE DATABASE $1` does not seem to work w CREATE DATABASE
		err = db.Exec(
			fmt.Sprintf("CREATE DATABASE %s", dbname),
		).Error
		if err != nil {
			t.Fatalf("unable to create postgres test database: %+v", err)
		}

		config.Set("DATABASE_URL", swapNewNameIntoDatabase(originalURL, dbname))

		return func() {
			reapChainlinkTestDB(t, originalURL, dbname)
			config.Set("DATABASE_URL", originalURL)
		}
	}

	return func() {}
}

func reapChainlinkTestDB(t testing.TB, dbURL, testdb string) {
	db, err := sql.Open(string(orm.DialectPostgres), dbURL)
	if err != nil {
		t.Fatalf("Unable to connect to parent CL db to clean up test database: %v", err)
	}
	defer db.Close()
	dbsSQL := "DROP DATABASE IF EXISTS " + testdb
	_, err = db.Exec(dbsSQL)
	if err != nil {
		t.Fatalf("Unable to clean up previous test database: %v", err)
	}
}

var chainlinkTestRe = regexp.MustCompile(`(/chainlink_test[_a-zA-Z0-9]*)`)

// swapNewNameIntoDatabase uses regex to swap the databasename from a postgres URL. eg:
// postgres://localhost:5432/chainlink_test?sslmode=disable becomes
// postgres://localhost:5432/chainlink_test_4d63a0af83c34e348292189c0648a2af?sslmode=disable
func swapNewNameIntoDatabase(dburl, newdb string) string {
	return chainlinkTestRe.ReplaceAllString(dburl, "/"+newdb)
}
