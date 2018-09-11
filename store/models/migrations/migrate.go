package migrations

import (
	"sort"
	"sync"

	"github.com/smartcontractkit/chainlink/store/models/migrations/migration1536521223"
	"github.com/smartcontractkit/chainlink/store/models/migrations/migration1536696950"
	"github.com/smartcontractkit/chainlink/store/models/orm"
)

func init() {
	registerMigration(migration1536521223.Migration{})
	registerMigration(migration1536696950.Migration{})
}

type migration interface {
	Migrate(orm *orm.ORM) error
	Timestamp() string
}

type MigrationTimestamp struct {
	Timestamp string `json:"timestamp" storm:"id"`
}

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
			err = availableMigrations[ts].Migrate(orm)
			if err != nil {
				return err
			}
			err = orm.Save(&MigrationTimestamp{ts})
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
