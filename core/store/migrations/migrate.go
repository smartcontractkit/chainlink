package migrations

import (
	"gorm.io/gorm"
)

var Migrations []*Migration

func Migrate(db *gorm.DB) error {
	return MigrateUp(db, "")
}

func MigrateUp(db *gorm.DB, to string) error {
	g := New(db, &Options{
		ValidateUnknownMigrations: false,
	}, Migrations)

	if to == "" {
		to = Migrations[len(Migrations)-1].ID
	}
	return g.MigrateTo(to)
}

func MigrateDown(db *gorm.DB) error {
	g := New(db, &Options{
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
	var from *Migration
	for _, m := range Migrations {
		if m.ID == name {
			from = m
		}
	}
	g := New(db, &Options{
		ValidateUnknownMigrations: false,
	}, Migrations)

	return g.RollbackMigration(from)
}

func Rollback(db *gorm.DB, m *Migration) error {
	g := New(db, &Options{
		ValidateUnknownMigrations: false,
	}, Migrations)

	return g.RollbackMigration(m)
}

func Current(db *gorm.DB) (*Migration, error) {
	g := New(db, &Options{
		ValidateUnknownMigrations: false,
	}, Migrations)

	if err := g.createMigrationTableIfNotExists(); err != nil {
		return nil, err
	}

	migration, err := g.getLastRunMigration()
	if err != nil {
		return nil, err
	}

	return migration, nil
}
