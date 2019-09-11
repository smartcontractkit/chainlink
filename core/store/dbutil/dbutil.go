package dbutil

import (
	"strings"

	"github.com/jinzhu/gorm"
)

func IsPostgresURL(url string) bool {
	return strings.HasPrefix(strings.ToLower(url), "postgres")
}

// IsPostgres returns true if the underlying database is postgres.
func IsPostgres(db *gorm.DB) bool {
	return db.Dialect().GetName() == "postgres"
}

// IsSqlite returns true if the underlying database is sqlite.
func IsSqlite(db *gorm.DB) bool {
	return strings.HasPrefix(db.Dialect().GetName(), "sqlite")
}

// SetTimezone sets the time zone to UTC
func SetTimezone(db *gorm.DB) error {
	if IsPostgres(db) {
		return db.Exec(`SET TIME ZONE 'UTC'`).Error
	}
	return nil
}

// SetSqlitePragmas sets some optimization params for SQLite
func SetSqlitePragmas(db *gorm.DB) error {
	if IsSqlite(db) {
		return db.Exec(`
			PRAGMA foreign_keys = ON;
			PRAGMA journal_mode = WAL;
		`).Error
	}
	return nil
}

// LimitSqliteOpenConnections deliberately limits Sqlites concurrency
// to reduce contention, reduce errors, and improve performance:
// https://stackoverflow.com/a/35805826/639773
func LimitSqliteOpenConnections(db *gorm.DB) error {
	if IsSqlite(db) {
		db.DB().SetMaxOpenConns(1)
	}
	return nil
}
