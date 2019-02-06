package migrations

import (
	"fmt"
	"sort"
	"strings"
	"sync"

	"github.com/smartcontractkit/chainlink/logger"
	"github.com/smartcontractkit/chainlink/store/migrations/migration0"
	"github.com/smartcontractkit/chainlink/store/orm"
	"go.uber.org/multierr"
)

func init() {
	registerMigration(migration0.Migration{})
}

type migration interface {
	Migrate(orm *orm.ORM) error
	Timestamp() string
}

// MigrationTimestamp tracks already run and available migrations.
type MigrationTimestamp struct {
	Timestamp string `json:"timestamp" gorm:"primary_key"`
}

// Migrate iterates through available migrations, running and tracking
// migrations that have not been run.
func Migrate(orm *orm.ORM) error {
	db := orm.DB
	if err := db.AutoMigrate(&MigrationTimestamp{}).Error; err != nil {
		return err
	}

	if err := runAlways(orm); err != nil {
		return err
	}

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
	return multierr.Append(
		setTimezone(orm),
		setForeignKeysOn(orm),
	)
}

func setTimezone(orm *orm.ORM) error {
	if orm.DB.Dialect().GetName() == "postgres" {
		return orm.DB.Exec(`SET TIME ZONE 'UTC';`).Error
	}
	return nil
}

func setForeignKeysOn(orm *orm.ORM) error {
	if strings.HasPrefix(orm.DB.Dialect().GetName(), "sqlite") {
		return orm.DB.Exec("PRAGMA foreign_keys = ON").Error
	}
	return nil
}
