package store

import (
	"context"
	"database/sql"
	"fmt"
	"regexp"

	"github.com/lib/pq"
)

const (
	dropDatabaseSQLFmt   = "DROP DATABASE IF EXISTS %s WITH (FORCE)"
	createDatabaseSQLFmt = "CREATE DATABASE %s"
)

var postgresDBNamePattern = regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*$`)

func quotePostgresDBName(name string) (string, error) {
	if name == "" || len(name) > 63 || !postgresDBNamePattern.MatchString(name) {
		return "", fmt.Errorf("invalid postgres database name: %q", name)
	}
	return pq.QuoteIdentifier(name), nil
}

func execDropDatabase(ctx context.Context, db *sql.DB, name string) error {
	quoted, err := quotePostgresDBName(name)
	if err != nil {
		return err
	}
	// PostgreSQL does not support bound parameters for database identifiers.
	//nolint:gosec // G701 -- name validated by quotePostgresDBName; identifier escaped with pq.QuoteIdentifier
	_, err = db.ExecContext(ctx, fmt.Sprintf(dropDatabaseSQLFmt, quoted))
	return err
}

func execCreateDatabase(ctx context.Context, db *sql.DB, name string) error {
	quoted, err := quotePostgresDBName(name)
	if err != nil {
		return err
	}
	// PostgreSQL does not support bound parameters for database identifiers.
	//nolint:gosec // G701 -- name validated by quotePostgresDBName; identifier escaped with pq.QuoteIdentifier
	_, err = db.ExecContext(ctx, fmt.Sprintf(createDatabaseSQLFmt, quoted))
	return err
}
