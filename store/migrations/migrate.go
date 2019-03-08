package migrations

import (
	"fmt"
	"sort"
	"strings"
	"sync"

	"github.com/smartcontractkit/chainlink/logger"
	"github.com/smartcontractkit/chainlink/store/migrations/migration0"
	"github.com/smartcontractkit/chainlink/store/migrations/migration1549496047"
	"github.com/smartcontractkit/chainlink/store/migrations/migration1551816486"
	"github.com/smartcontractkit/chainlink/store/orm"
	"go.uber.org/multierr"
)

func init() {
	registerMigration(migration0.Migration{})
	registerMigration(migration1549496047.Migration{})
	registerMigration(migration1551816486.Migration{})
}

type migration interface {
	Migrate(orm *orm.ORM) error
	Timestamp() string
}

// MigrationTimestamp tracks already run and available migrations.
type MigrationTimestamp struct {
	Timestamp string `json:"timestamp" gorm:"primary_key;type:varchar(12)"`
}

// Migrate iterates through available migrations, running and tracking
// migrations that have not been run.
func Migrate(orm *orm.ORM) error {
	if err := runAlways(orm); err != nil {
		return err
	}

	db := orm.DB
	var migrationTimestamps []MigrationTimestamp
	err := db.Order("timestamp asc").Find(&migrationTimestamps).Error
	if err != nil {
		return err
	}

	alreadyMigratedSet := make(map[string]bool)
	for _, mt := range migrationTimestamps {
		alreadyMigratedSet[mt.Timestamp] = true
	}

	sortedTimestamps := availableMigrationTimestamps()
	for _, ts := range sortedTimestamps {
		_, already := alreadyMigratedSet[ts]
		if !already {
			logger.Debug("Migrating ", ts)
			err = availableMigrations[ts].Migrate(orm)
			if err != nil {
				return multierr.Append(fmt.Errorf("Failed migration %v", ts), err)
			}
			err = db.Create(&MigrationTimestamp{Timestamp: ts}).Error
			if err != nil {
				return err
			}
		}
	}

	return nil
}

var migrationMutex sync.RWMutex
var availableMigrations = make(map[string]migration)

func registerMigration(migration migration) {
	migrationMutex.Lock()
	availableMigrations[migration.Timestamp()] = migration
	migrationMutex.Unlock()
}

func availableMigrationTimestamps() []string {
	migrationMutex.RLock()
	defer migrationMutex.RUnlock()

	var sortedTimestamps []string
	for k := range availableMigrations {
		sortedTimestamps = append(sortedTimestamps, k)
	}
	sort.Strings(sortedTimestamps)
	return sortedTimestamps
}

func runAlways(orm *orm.ORM) error {
	return multierr.Combine(
		setTimezone(orm),
		setForeignKeysOn(orm),
		limitSqliteOpenConnections(orm),
		automigrateMigrationsTable(orm),
	)
}

func automigrateMigrationsTable(orm *orm.ORM) error {
	return orm.DB.AutoMigrate(&MigrationTimestamp{}).Error
}

func setTimezone(orm *orm.ORM) error {
	if orm.DB.Dialect().GetName() == "postgres" {
		return orm.DB.Exec(`SET TIME ZONE 'UTC'`).Error
	}
	return nil
}

func setForeignKeysOn(orm *orm.ORM) error {
	if strings.HasPrefix(orm.DB.Dialect().GetName(), "sqlite") {
		return orm.DB.Exec(`
			PRAGMA foreign_keys = ON;
			PRAGMA journal_mode = WAL;
		`).Error
	}
	return nil
}

// limitSqliteOpenConnections deliberately limits Sqlites concurrency
// to reduce contention, reduce errors, and improve performance:
// https://stackoverflow.com/a/35805826/639773
func limitSqliteOpenConnections(orm *orm.ORM) error {
	if strings.HasPrefix(orm.DB.Dialect().GetName(), "sqlite") {
		orm.DB.DB().SetMaxOpenConns(1)
	}
	return nil
}
