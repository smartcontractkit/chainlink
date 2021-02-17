package dbutil

import (
	gormv1 "github.com/jinzhu/gorm"
	"gorm.io/gorm"
)

// IsPostgres returns true if the underlying database is postgres.
func IsPostgres(db *gormv1.DB) bool {
	return db.Dialect().GetName() == "postgres"
}

// SetTimezone sets the time zone to UTC
func SetTimezone(db *gorm.DB) error {
	return db.Exec(`SET TIME ZONE 'UTC'`).Error
}
