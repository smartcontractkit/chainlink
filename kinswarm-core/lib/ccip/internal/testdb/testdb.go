package testdb

import (
	"database/sql"
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/smartcontractkit/chainlink/v2/core/store/dialects"
)

const (
	// PristineDBName is a clean copy of test DB with migrations.
	PristineDBName = "chainlink_test_pristine"
	// TestDBNamePrefix is a common prefix that will be auto-removed by the dangling DB cleanup process.
	TestDBNamePrefix = "chainlink_test_"
)

// CreateOrReplace creates a database named with a common prefix and the given suffix, and returns the URL.
// If the database already exists, it will be dropped and re-created.
// If withTemplate is true, the pristine DB will be used as a template.
func CreateOrReplace(parsed url.URL, suffix string, withTemplate bool) (string, error) {
	if parsed.Path == "" {
		return "", errors.New("path missing from database URL")
	}

	// Match the naming schema that our dangling DB cleanup methods expect
	dbname := TestDBNamePrefix + suffix
	if l := len(dbname); l > 63 {
		return "", fmt.Errorf("dbname %v too long (%d), max is 63 bytes. Try a shorter suffix", dbname, l)
	}
	// Cannot drop test database if we are connected to it, so we must connect
	// to a different one. 'postgres' should be present on all postgres installations
	parsed.Path = "/postgres"
	db, err := sql.Open(string(dialects.Postgres), parsed.String())
	if err != nil {
		return "", fmt.Errorf("in order to drop the test database, we need to connect to a separate database"+
			" called 'postgres'. But we are unable to open 'postgres' database: %+v\n", err)
	}
	defer db.Close()

	_, err = db.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS %s", dbname))
	if err != nil {
		return "", fmt.Errorf("unable to drop postgres migrations test database: %v", err)
	}
	if withTemplate {
		_, err = db.Exec(fmt.Sprintf("CREATE DATABASE %s WITH TEMPLATE %s", dbname, PristineDBName))
	} else {
		_, err = db.Exec(fmt.Sprintf("CREATE DATABASE %s", dbname))
	}
	if err != nil {
		return "", fmt.Errorf("unable to create postgres test database with name '%s': %v", dbname, err)
	}
	parsed.Path = fmt.Sprintf("/%s", dbname)
	return parsed.String(), nil
}

// Drop drops the database at the given URL.
func Drop(dbURL url.URL) error {
	if dbURL.Path == "" {
		return errors.New("path missing from database URL")
	}
	dbname := strings.TrimPrefix(dbURL.Path, "/")

	// Cannot drop test database if we are connected to it, so we must connect
	// to a different one. 'postgres' should be present on all postgres installations
	dbURL.Path = "/postgres"
	db, err := sql.Open(string(dialects.Postgres), dbURL.String())
	if err != nil {
		return fmt.Errorf("in order to drop the test database, we need to connect to a separate database"+
			" called 'postgres'. But we are unable to open 'postgres' database: %+v\n", err)
	}
	defer db.Close()

	_, err = db.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS %s", dbname))
	if err != nil {
		return fmt.Errorf("unable to drop postgres migrations test database: %v", err)
	}
	return nil
}
