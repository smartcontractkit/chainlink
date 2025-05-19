package testdb

import (
	"net/url"
	"testing"

	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil/sqltest"
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
func CreateOrReplace(t testing.TB, parsed url.URL, suffix string, withTemplate bool) url.URL {
	// Match the naming schema that our dangling DB cleanup methods expect
	dbname := TestDBNamePrefix + suffix
	var template string
	if withTemplate {
		template = PristineDBName
	}
	return sqltest.CreateOrReplace(t, parsed, dbname, template)
}
