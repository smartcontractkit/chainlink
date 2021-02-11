package dbutil

import (
	"strings"

	gormv1 "github.com/jinzhu/gorm"
	"gorm.io/gorm"
)

func IsPostgresURL(url string) bool {
	return strings.HasPrefix(strings.ToLower(url), "postgres")
}

// IsPostgres returns true if the underlying database is postgres.
func IsPostgres(db *gormv1.DB) bool {
	return db.Dialect().GetName() == "postgres"
}

// SetTimezone sets the time zone to UTC
func SetTimezone(db *gorm.DB) error {
	return db.Exec(`SET TIME ZONE 'UTC'`).Error
}
