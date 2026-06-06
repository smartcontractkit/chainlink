package store

import (
	"context"
	"database/sql"
	"fmt"
	"regexp"

	"github.com/lib/pq"
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
	_, err = db.ExecContext(ctx, "DROP DATABASE IF EXISTS "+quoted+" WITH (FORCE)")
	return err
}

func execCreateDatabase(ctx context.Context, db *sql.DB, name string) error {
	quoted, err := quotePostgresDBName(name)
	if err != nil {
		return err
	}
	_, err = db.ExecContext(ctx, "CREATE DATABASE "+quoted)
	return err
}
