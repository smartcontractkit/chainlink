package testdb

import (
	"context"
	"database/sql"

	"github.com/peterldowns/pgtestdb"
	"github.com/smartcontractkit/chainlink/v2/core/store/migrate"
)

const (
	// PristineDBName is a clean copy of test DB with migrations.
	PristineDBName = "chainlink_test_pristine"
	// TestDBNamePrefix is a common prefix that will be auto-removed by the dangling DB cleanup process.
	TestDBNamePrefix = "chainlink_test_"
)

type migrator struct {
	withTemplate bool
}

func (m *migrator) Hash() (string, error) {
	if m.withTemplate {
		return "withTemplate", nil
	}
	return "empty", nil
}

func (m *migrator) Migrate(ctx context.Context, db *sql.DB, config pgtestdb.Config) error {
	if !m.withTemplate {
		return nil
	}
	return migrate.Migrate(ctx, db)
}

func Migrator(withTemplate bool) pgtestdb.Migrator {
	return &migrator{withTemplate: withTemplate}
}
