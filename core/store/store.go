package store

import (
	"database/sql"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"sync"

	"github.com/coreos/go-semver/semver"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/periodicbackup"
	"github.com/smartcontractkit/chainlink/core/services/versioning"
	"github.com/smartcontractkit/chainlink/core/static"

	"github.com/smartcontractkit/chainlink/core/gracefulpanic"
	"github.com/smartcontractkit/chainlink/core/services/postgres"
	"github.com/smartcontractkit/chainlink/core/store/config"
	"github.com/smartcontractkit/chainlink/core/store/migrate"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/store/orm"
	"github.com/smartcontractkit/chainlink/core/utils"

	"github.com/pkg/errors"
	"go.uber.org/multierr"
	"gorm.io/gorm"
)

const (
	// AutoMigrate is a flag that automatically migrates the DB when passed to initializeORM
	AutoMigrate = "auto_migrate"
)

// Store contains fields for the database, Config
// for keeping the application state in sync with the database.
type Store struct {
	*orm.ORM
	Config         config.GeneralConfig
	Clock          utils.AfterNower
	AdvisoryLocker postgres.AdvisoryLocker
	closeOnce      *sync.Once
}

// NewStore will create a new store
// func NewStore(config config.GeneralConfig, ethClient eth.Client, advisoryLock postgres.AdvisoryLocker, shutdownSignal gracefulpanic.Signal, keyStoreGenerator KeyStoreGenerator) (*Store, error) {
func NewStore(config config.GeneralConfig, advisoryLock postgres.AdvisoryLocker, shutdownSignal gracefulpanic.Signal) (*Store, error) {
	return newStore(config, advisoryLock, shutdownSignal)
}

// NewInsecureStore creates a new store with the given config using an insecure keystore.
// NOTE: Should only be used for testing!
func NewInsecureStore(config config.GeneralConfig, advisoryLocker postgres.AdvisoryLocker, shutdownSignal gracefulpanic.Signal) (*Store, error) {
	return newStore(config, advisoryLocker, shutdownSignal)
}

func newStore(
	config config.GeneralConfig,
	advisoryLocker postgres.AdvisoryLocker,
	shutdownSignal gracefulpanic.Signal,
) (*Store, error) {
	if err := utils.EnsureDirAndMaxPerms(config.RootDir(), os.FileMode(0700)); err != nil {
		return nil, errors.Wrap(err, "error while creating project root dir")
	}

	orm, err := initializeORM(config, shutdownSignal)
	if err != nil {
		return nil, errors.Wrap(err, "failed to initialize ORM")
	}

	store := &Store{
		Clock:          utils.Clock{},
		AdvisoryLocker: advisoryLocker,
		Config:         config,
		ORM:            orm,
		closeOnce:      &sync.Once{},
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
		err = multierr.Append(err, s.AdvisoryLocker.Close())
	})
	return err
}

func (s *Store) Ready() error {
	return nil
}

func (s *Store) Healthy() error {
	return nil
}

// Unscoped returns a shallow copy of the store, with an unscoped ORM allowing
// one to work with soft deleted records.
func (s *Store) Unscoped() *Store {
	cpy := *s
	cpy.ORM = s.ORM.Unscoped()
	return &cpy
}

// AuthorizedUserWithSession will return the one API user if the Session ID exists
// and hasn't expired, and update session's LastUsed field.
func (s *Store) AuthorizedUserWithSession(sessionID string) (models.User, error) {
	return s.ORM.AuthorizedUserWithSession(
		sessionID, s.Config.SessionTimeout().Duration())
}

func CheckSquashUpgrade(db *gorm.DB) error {
	// Ensure that we don't try to run a node version later than the
	// squashed database versions without first migrating up to just before the squash.
	// If we don't see the latest migration and node version >= S, error and recommend
	// first running version S - 1 (S = version in which migrations are squashed).
	if static.Version == "unset" {
		return nil
	}
	squashVersionMinus1 := semver.New("0.9.10")
	currentVersion, err := semver.NewVersion(static.Version)
	if err != nil {
		return errors.Wrapf(err, "expected VERSION to be valid semver (for example 1.42.3). Got: %s", static.Version)
	}
	lastV1Migration := "1611847145"
	if squashVersionMinus1.LessThan(*currentVersion) {
		// Completely empty database is fine to run squashed migrations on
		if !db.Migrator().HasTable("migrations") {
			return nil
		}
		// Running code later than S - 1. Ensure that we see
		// the last v1 migration.
		q := db.Exec("SELECT * FROM migrations WHERE id = ?", lastV1Migration)
		if q.Error != nil {
			return q.Error
		}
		if q.RowsAffected == 0 {
			// Do not have the S-1 migration.
			return errors.Errorf("Need to upgrade to chainlink version %v first before upgrading to version %v in order to run migrations", squashVersionMinus1, currentVersion)
		}
	}
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
				logger.Default.Debugf("Failed to find any node version in the DB: %w", err)
			} else if strings.Contains(err.Error(), "relation \"node_versions\" does not exist") {
				logger.Default.Debugf("Failed to find any node version in the DB, the node_versions table does not exist yet: %w", err)
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
	if err = CheckSquashUpgrade(dbOrm.DB); err != nil {
		panic(err)
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
