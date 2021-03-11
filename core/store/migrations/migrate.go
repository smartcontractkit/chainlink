package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

var Migrations []*gormigrate.Migration

func Migrate(db *gorm.DB) error {
	return MigrateUp(db, "")
}

func MigrateUp(db *gorm.DB, to string) error {
	// We don't want to wrap all the migrations in a tx.
	// Gorm v2 uses a transaction by default.
	g := gormigrate.New(db, &gormigrate.Options{
		UseTransaction:            false,
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
		UseTransaction:            false,
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

func MigrateDownFrom(db *gorm.DB, name string) error {
	var from *gormigrate.Migration
	for _, m := range Migrations {
		if m.ID == name {
			from = m
		}
	}
	g := gormigrate.New(db, &gormigrate.Options{
		UseTransaction:            false,
		ValidateUnknownMigrations: false,
	}, Migrations)

	return g.RollbackMigration(from)
}

func Rollback(db *gorm.DB, m *gormigrate.Migration) error {
	g := gormigrate.New(db, &gormigrate.Options{
		UseTransaction:            false,
		ValidateUnknownMigrations: false,
	}, Migrations)

	return g.RollbackMigration(m)
}
