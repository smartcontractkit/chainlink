package migrationsv2

import (
	"github.com/jinzhu/gorm"
	gormigrate "gopkg.in/gormigrate.v1"
)

var Migrations []*gormigrate.Migration

func Migrate(db *gorm.DB) error {
	return MigrateUp(db, "")
}

func MigrateUp(db *gorm.DB, to string) error {
	g := gormigrate.New(db, &gormigrate.Options{
		UseTransaction:            true,
		ValidateUnknownMigrations: false,
	}, Migrations)

	if to == "" {
		to = Migrations[len(Migrations)-1].ID
	}
	if err := g.MigrateTo(to); err != nil {
		return err
	}
	return nil
}

func MigrateDown(db *gorm.DB) error {
	g := gormigrate.New(db, &gormigrate.Options{
		UseTransaction:            true,
		ValidateUnknownMigrations: false,
	}, Migrations)
	return g.RollbackLast()
}
