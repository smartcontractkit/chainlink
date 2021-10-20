package migrations

// NOTE: This package is copied from gormigrate, with code cleanup and applying
// some fixes around transactional migrations
// Source: https://github.com/go-gormigrate/gormigrate/blob/dacf763a39d8b491bd13f34bbc583ecb1640094f/go
//
// One thing we may wish to do in future is add the InitSchema function which
// skips migrations and loads the schema from a dump on a virgin database.

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/services/postgres"
	"go.uber.org/multierr"
	"gorm.io/gorm"
)

const (
	// MaxIndividualMigrationTime is the maximum amount of time a single
	// migration is allowed to take before declaring it a failure
	MaxIndividualMigrationTime = 5 * time.Minute
	initSchemaMigrationID      = "SCHEMA_INIT"
)

// MigrateFunc is the func signature for migrating.
type MigrateFunc func(*gorm.DB) error

// RollbackFunc is the func signature for rollbacking.
type RollbackFunc func(*gorm.DB) error

// InitSchemaFunc is the func signature for initializing the schema.
type InitSchemaFunc func(*gorm.DB) error

// Options define options for all migrations.
type Options struct {
	// TableName is the migration table.
	TableName string
	// IDColumnName is the name of column where the migration id will be stored.
	IDColumnName string
	// IDColumnSize is the length of the migration id column
	IDColumnSize int
	// ValidateUnknownMigrations will cause migrate to fail if there's unknown migration
	// IDs in the database
	ValidateUnknownMigrations bool
}

// Migration represents a database migration (a modification to be made on the database).
type Migration struct {
	// ID is the migration identifier. Usually a timestamp like "201601021504".
	ID string
	// Migrate is a function that will br executed while running this migration.
	Migrate MigrateFunc
	// Rollback will be executed on rollback. Can be nil.
	Rollback RollbackFunc
	// DisableTransaction disables wrapping the migration in a transaction. Useful in
	// very rare cases, be careful because it can leave the database in an
	// inconsistent state
	DisableTransaction bool
}

// Gormigrate represents a collection of all migrations of a database schema.
type Gormigrate struct {
	db         *gorm.DB
	options    *Options
	migrations []*Migration
}

// ReservedIDError is returned when a migration is using a reserved ID
type ReservedIDError struct {
	ID string
}

func (e *ReservedIDError) Error() string {
	return fmt.Sprintf(`gormigrate: Reserved migration ID: "%s"`, e.ID)
}

// DuplicatedIDError is returned when more than one migration have the same ID
type DuplicatedIDError struct {
	ID string
}

func (e *DuplicatedIDError) Error() string {
	return fmt.Sprintf(`gormigrate: Duplicated migration ID: "%s"`, e.ID)
}

var (
	// DefaultOptions can be used if you don't want to think about options.
	DefaultOptions = &Options{
		TableName:                 "migrations",
		IDColumnName:              "id",
		IDColumnSize:              255,
		ValidateUnknownMigrations: false,
	}

	// ErrRollbackImpossible is returned when trying to rollback a migration
	// that has no rollback function.
	ErrRollbackImpossible = errors.New("gormigrate: It's impossible to rollback this migration")

	// ErrNoMigrationDefined is returned when no migration is defined.
	ErrNoMigrationDefined = errors.New("gormigrate: No migration defined")

	// ErrMissingID is returned when the ID od migration is equal to ""
	ErrMissingID = errors.New("gormigrate: Missing ID in migration")

	// ErrNoRunMigration is returned when any run migration was found while
	// running RollbackLast
	ErrNoRunMigration = errors.New("gormigrate: Could not find last run migration")

	// ErrMigrationIDDoesNotExist is returned when migrating or rolling back to a migration ID that
	// does not exist in the list of migrations
	ErrMigrationIDDoesNotExist = errors.New("gormigrate: Tried to migrate to an ID that doesn't exist")

	// ErrUnknownPastMigration is returned if a migration exists in the DB that doesn't exist in the code
	ErrUnknownPastMigration = errors.New("gormigrate: Found migration in DB that does not exist in code")
)

// New returns a new Gormigrate.
func New(db *gorm.DB, options *Options, migrations []*Migration) *Gormigrate {
	if options.TableName == "" {
		options.TableName = DefaultOptions.TableName
	}
	if options.IDColumnName == "" {
		options.IDColumnName = DefaultOptions.IDColumnName
	}
	if options.IDColumnSize == 0 {
		options.IDColumnSize = DefaultOptions.IDColumnSize
	}
	return &Gormigrate{
		db:         db,
		options:    options,
		migrations: migrations,
	}
}

// Migrate executes all migrations that did not run yet.
func (g *Gormigrate) Migrate() error {
	if !g.hasMigrations() {
		return ErrNoMigrationDefined
	}
	var targetMigrationID string
	if len(g.migrations) > 0 {
		targetMigrationID = g.migrations[len(g.migrations)-1].ID
	}
	return g.migrate(targetMigrationID)
}

// MigrateTo executes all migrations that did not run yet up to the migration that matches `migrationID`.
func (g *Gormigrate) MigrateTo(migrationID string) error {
	if err := g.checkIDExist(migrationID); err != nil {
		return err
	}
	return g.migrate(migrationID)
}

func (g *Gormigrate) migrate(migrationID string) error {
	if !g.hasMigrations() {
		return ErrNoMigrationDefined
	}

	if err := g.checkReservedID(); err != nil {
		return err
	}

	if err := g.checkDuplicatedID(); err != nil {
		return err
	}

	if err := g.createMigrationTableIfNotExists(); err != nil {
		return err
	}

	if g.options.ValidateUnknownMigrations {
		unknownMigrations, err := g.unknownMigrationsHaveHappened()
		if err != nil {
			return err
		}
		if unknownMigrations {
			return ErrUnknownPastMigration
		}
	}

	for _, migration := range g.migrations {
		if err := g.runMigration(migration); err != nil {
			return err
		}
		if migrationID != "" && migration.ID == migrationID {
			break
		}
	}
	return nil
}

// There are migrations to apply if either there's a defined
// initSchema function or if the list of migrations is not empty.
func (g *Gormigrate) hasMigrations() bool {
	return len(g.migrations) > 0
}

// Check whether any migration is using a reserved ID.
// For now there's only have one reserved ID, but there may be more in the future.
func (g *Gormigrate) checkReservedID() error {
	for _, m := range g.migrations {
		if m.ID == initSchemaMigrationID {
			return &ReservedIDError{ID: m.ID}
		}
	}
	return nil
}

func (g *Gormigrate) checkDuplicatedID() error {
	lookup := make(map[string]struct{}, len(g.migrations))
	for _, m := range g.migrations {
		if _, ok := lookup[m.ID]; ok {
			return &DuplicatedIDError{ID: m.ID}
		}
		lookup[m.ID] = struct{}{}
	}
	return nil
}

func (g *Gormigrate) checkIDExist(migrationID string) error {
	for _, migrate := range g.migrations {
		if migrate.ID == migrationID {
			return nil
		}
	}
	return ErrMigrationIDDoesNotExist
}

// RollbackLast undo the last migration
func (g *Gormigrate) RollbackLast() error {
	if len(g.migrations) == 0 {
		return ErrNoMigrationDefined
	}

	lastRunMigration, err := g.getLastRunMigration()
	if err != nil {
		return err
	}

	return g.RollbackMigration(lastRunMigration)
}

// RollbackTo undoes migrations up to the given migration that matches the `migrationID`.
// Migration with the matching `migrationID` is not rolled back.
func (g *Gormigrate) RollbackTo(migrationID string) error {
	if len(g.migrations) == 0 {
		return ErrNoMigrationDefined
	}

	if err := g.checkIDExist(migrationID); err != nil {
		return err
	}

	for i := len(g.migrations) - 1; i >= 0; i-- {
		migration := g.migrations[i]
		if migration.ID == migrationID {
			break
		}
		if err := g.RollbackMigration(migration); err != nil {
			return err
		}
	}
	return nil
}

func (g *Gormigrate) getLastRunMigration() (*Migration, error) {
	for i := len(g.migrations) - 1; i >= 0; i-- {
		migration := g.migrations[i]

		migrationRan, err := migrationRan(g.db, migration, g.options)
		if err != nil {
			return nil, err
		}

		if migrationRan {
			return migration, nil
		}
	}
	return nil, ErrNoRunMigration
}

// RollbackMigration rolls back a migration.
func (g *Gormigrate) RollbackMigration(migration *Migration) error {
	if migration.Rollback == nil {
		return ErrRollbackImpossible
	}

	ctx, cancel := context.WithTimeout(context.Background(), MaxIndividualMigrationTime)
	defer cancel()

	var err error
	if migration.DisableTransaction {
		db := g.db.WithContext(ctx)
		err = errors.Wrap(rollbackMigrationNoDDL(db, migration, g.options), "WARNING: DDL was disabled, your database may be in an inconsistent state")
	} else {
		err = postgres.GormTransaction(ctx, g.db, func(dbtx *gorm.DB) error {
			return rollbackMigrationNoDDL(dbtx, migration, g.options)
		})
	}

	return errors.Wrapf(err, "failed to rollback migration %s", migration.ID)
}

func rollbackMigrationNoDDL(db *gorm.DB, migration *Migration, options *Options) error {
	migrationRan, err := migrationRan(db, migration, options)
	if err != nil {
		return err
	}
	if migrationRan {
		if err := migration.Rollback(db); err != nil {
			return err
		}

		/* #nosec G201 */
		sql := fmt.Sprintf("DELETE FROM %s WHERE %s = ?", options.TableName, options.IDColumnName)
		return db.Exec(sql, migration.ID).Error
	}

	return nil
}

func (g *Gormigrate) runMigration(migration *Migration) error {
	if len(migration.ID) == 0 {
		return ErrMissingID
	}

	ctx, cancel := context.WithTimeout(context.Background(), MaxIndividualMigrationTime)
	defer cancel()

	var err error
	if migration.DisableTransaction {
		db := g.db.WithContext(ctx)
		err = errors.Wrap(runMigrationNoDDL(db, migration, g.options), "WARNING: DDL was disabled, your database may be in an inconsistent state")
	} else {
		err = postgres.GormTransaction(ctx, g.db, func(dbtx *gorm.DB) error {
			return runMigrationNoDDL(dbtx, migration, g.options)
		})
	}

	return errors.Wrapf(err, "failed to run migration %s", migration.ID)

}

func runMigrationNoDDL(db *gorm.DB, migration *Migration, options *Options) error {
	migrationRan, err := migrationRan(db, migration, options)
	if err != nil {
		return err
	}
	if !migrationRan {
		if err := migration.Migrate(db); err != nil {
			return err
		}

		return insertMigration(db, migration.ID, options)
	}
	return nil
}

func (g *Gormigrate) createMigrationTableIfNotExists() error {
	if g.db.Migrator().HasTable("goose_migrations") {
		return errors.New("a newer version of chainlink has already migrated this database; it is not safe to run this release")
	}
	if g.db.Migrator().HasTable(g.options.TableName) {
		return nil
	}

	/* #nosec G201 */
	sql := fmt.Sprintf("CREATE TABLE %s (%s VARCHAR(%d) PRIMARY KEY)", g.options.TableName, g.options.IDColumnName, g.options.IDColumnSize)
	return g.db.Exec(sql).Error
}

func migrationRan(db *gorm.DB, m *Migration, options *Options) (bool, error) {
	var count int64
	err := db.
		Table(options.TableName).
		/* #nosec G201 */
		Where(fmt.Sprintf("%s = ?", options.IDColumnName), m.ID).
		Count(&count).
		Error
	return count > 0, err
}

func (g *Gormigrate) unknownMigrationsHaveHappened() (unknown bool, merr error) {
	/* #nosec G201 */
	sql := fmt.Sprintf("SELECT %s FROM %s", g.options.IDColumnName, g.options.TableName)
	rows, err := g.db.Raw(sql).Rows()
	if err != nil {
		merr = err
		return
	}
	defer func() {
		merr = multierr.Combine(merr, rows.Close())
	}()

	validIDSet := make(map[string]struct{}, len(g.migrations)+1)
	validIDSet[initSchemaMigrationID] = struct{}{}
	for _, migration := range g.migrations {
		validIDSet[migration.ID] = struct{}{}
	}

	for rows.Next() {
		var pastMigrationID string
		if err := rows.Scan(&pastMigrationID); err != nil {
			merr = err
			return
		}
		if _, ok := validIDSet[pastMigrationID]; !ok {
			unknown = true
			return
		}
	}

	return
}

func insertMigration(db *gorm.DB, id string, options *Options) error {
	/* #nosec G201 */
	sql := fmt.Sprintf("INSERT INTO %s (%s) VALUES (?)", options.TableName, options.IDColumnName)
	return db.Exec(sql, id).Error
}
