package migrations

import (
	"sort"
	"sync"

	"github.com/smartcontractkit/chainlink/logger"
	"github.com/smartcontractkit/chainlink/store/migrations/migration0"
	"github.com/smartcontractkit/chainlink/store/migrations/migration1536696950"
	"github.com/smartcontractkit/chainlink/store/migrations/migration1536764911"
	"github.com/smartcontractkit/chainlink/store/migrations/migration1537223654"
	"github.com/smartcontractkit/chainlink/store/orm"
)

func init() {
	registerMigration(migration0.Migration{})
	registerMigration(migration1536696950.Migration{})
	registerMigration(migration1536764911.Migration{})
	registerMigration(migration1537223654.Migration{})
}

type migration interface {
	Migrate(orm *orm.ORM) error
	Timestamp() string
}

// MigrationTimestamp tracks already run and available migrations.
type MigrationTimestamp struct {
	Timestamp string `json:"timestamp" storm:"id"`
}

// Migrate iterates through available migrations, running and tracking
// migrations that have not been run.
func Migrate(orm *orm.ORM) error {
	err := orm.InitBucket(&MigrationTimestamp{})
	if err != nil {
		return err
	}

	var migrationTimestamps []MigrationTimestamp
	err = orm.AllByIndex("Timestamp", &migrationTimestamps)
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
				return err
			}
			err = orm.DB.Save(&MigrationTimestamp{ts})
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
