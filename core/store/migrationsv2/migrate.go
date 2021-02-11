package migrationsv2

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
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

	for i := len(Migrations) - 1; i >= 0; i-- {
		err := g.RollbackMigration(Migrations[i])
		if err != nil {
			return err
		}
	}
	return nil
}

func Rollback(db *gorm.DB, m *gormigrate.Migration) error {
	g := gormigrate.New(db, &gormigrate.Options{
		UseTransaction:            true,
		ValidateUnknownMigrations: false,
	}, Migrations)

	return g.RollbackMigration(m)
}
