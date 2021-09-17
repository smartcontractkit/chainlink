package store

import (
	"database/sql"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"sync"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/periodicbackup"
	"github.com/smartcontractkit/chainlink/core/services/versioning"
	"github.com/smartcontractkit/chainlink/core/static"

	"github.com/smartcontractkit/chainlink/core/gracefulpanic"
	"github.com/smartcontractkit/chainlink/core/services/postgres"
	"github.com/smartcontractkit/chainlink/core/store/config"
	"github.com/smartcontractkit/chainlink/core/store/migrate"
	"github.com/smartcontractkit/chainlink/core/store/orm"
	"github.com/smartcontractkit/chainlink/core/utils"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

// Store contains fields for the database, Config
// for keeping the application state in sync with the database.
type Store struct {
	*orm.ORM
	Config    config.GeneralConfig
	closeOnce *sync.Once
}

// NewStore will create a new store
func NewStore(config config.GeneralConfig, advisoryLock postgres.AdvisoryLocker, shutdownSignal gracefulpanic.Signal) (*Store, error) {
	if err := utils.EnsureDirAndMaxPerms(config.RootDir(), os.FileMode(0700)); err != nil {
		return nil, errors.Wrap(err, "error while creating project root dir")
	}

	orm, err := initializeORM(config, shutdownSignal)
	if err != nil {
		return nil, errors.Wrap(err, "failed to initialize ORM")
	}

	store := &Store{
		Config:    config,
		ORM:       orm,
		closeOnce: &sync.Once{},
	}
	return store, nil
}

func (s *Store) Start() error {
	return nil
}

// Close shuts down all of the working parts of the store.
func (s *Store) Close() error {
	var err error
	s.closeOnce.Do(func() {
		err = s.ORM.Close()
	})
	return err
}

func (s *Store) Ready() error {
	return nil
}

func (s *Store) Healthy() error {
	return nil
}

func initializeORM(cfg config.GeneralConfig, shutdownSignal gracefulpanic.Signal) (*orm.ORM, error) {
	dbURL := cfg.DatabaseURL()
	dbOrm, err := orm.NewORM(dbURL.String(), cfg.DatabaseTimeout(), shutdownSignal, cfg.GetDatabaseDialectConfiguredOrDefault(), cfg.GetAdvisoryLockIDConfiguredOrDefault(), cfg.GlobalLockRetryInterval().Duration(), cfg.ORMMaxOpenConns(), cfg.ORMMaxIdleConns())
	if err != nil {
		return nil, errors.Wrap(err, "initializeORM#NewORM")
	}

	// Set up the versioning ORM
	verORM := versioning.NewORM(postgres.WrapDbWithSqlx(
		postgres.MustSQLDB(dbOrm.DB)),
	)

	if cfg.DatabaseBackupMode() != config.DatabaseBackupModeNone {
		var version *versioning.NodeVersion
		var versionString string

		version, err = verORM.FindLatestNodeVersion()
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				logger.Default.Debugf("Failed to find any node version in the DB: %v", err)
			} else if strings.Contains(err.Error(), "relation \"node_versions\" does not exist") {
				logger.Default.Debugf("Failed to find any node version in the DB, the node_versions table does not exist yet: %v", err)
			} else {
				return nil, errors.Wrap(err, "initializeORM#FindLatestNodeVersion")
			}
		}

		if version != nil {
			versionString = version.Version
		}

		databaseBackup := periodicbackup.NewDatabaseBackup(cfg, logger.Default)
		databaseBackup.RunBackupGracefully(versionString)
	}
	if cfg.MigrateDatabase() {
		dbOrm.SetLogging(cfg.LogSQLStatements() || cfg.LogSQLMigrations())

		err = dbOrm.RawDBWithAdvisoryLock(func(db *gorm.DB) error {
			return migrate.Migrate(postgres.UnwrapGormDB(db).DB)
		})
		if err != nil {
			return nil, errors.Wrap(err, "initializeORM#Migrate")
		}
	}

	nodeVersion := static.Version
	if nodeVersion == "unset" {
		nodeVersion = fmt.Sprintf("random_%d", rand.Uint32())
	}
	version := versioning.NewNodeVersion(nodeVersion)
	err = verORM.UpsertNodeVersion(version)
	if err != nil {
		return nil, errors.Wrap(err, "initializeORM#UpsertNodeVersion")
	}

	dbOrm.SetLogging(cfg.LogSQLStatements())
	return dbOrm, nil
}
