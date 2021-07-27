package store

import (
	"fmt"
	"math/rand"
	"os"
	"sync"

	"github.com/coreos/go-semver/semver"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/periodicbackup"
	"github.com/smartcontractkit/chainlink/core/static"

	"github.com/smartcontractkit/chainlink/core/gracefulpanic"
	"github.com/smartcontractkit/chainlink/core/services/eth"
	"github.com/smartcontractkit/chainlink/core/services/postgres"
	"github.com/smartcontractkit/chainlink/core/store/config"
	"github.com/smartcontractkit/chainlink/core/store/migrations"
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
	Config         *config.Config
	Clock          utils.AfterNower
	AdvisoryLocker postgres.AdvisoryLocker
	closeOnce      *sync.Once
}

// NewStore will create a new store
// func NewStore(config *config.Config, ethClient eth.Client, advisoryLock postgres.AdvisoryLocker, shutdownSignal gracefulpanic.Signal, keyStoreGenerator KeyStoreGenerator) (*Store, error) {
func NewStore(config *config.Config, ethClient eth.Client, advisoryLock postgres.AdvisoryLocker, shutdownSignal gracefulpanic.Signal) (*Store, error) {
	// return newStore(config, ethClient, advisoryLock, keyStoreGenerator, shutdownSignal)
	return newStore(config, ethClient, advisoryLock, shutdownSignal)
}

// NewInsecureStore creates a new store with the given config using an insecure keystore.
// NOTE: Should only be used for testing!
func NewInsecureStore(config *config.Config, ethClient eth.Client, advisoryLocker postgres.AdvisoryLocker, shutdownSignal gracefulpanic.Signal) (*Store, error) {
	// return newStore(config, ethClient, advisoryLocker, InsecureKeyStoreGen, shutdownSignal)
	return newStore(config, ethClient, advisoryLocker, shutdownSignal)
}

// TODO(sam): Remove ethClient from here completely after legacy tx manager is gone
// See: https://www.pivotaltracker.com/story/show/175493792
func newStore(
	config *config.Config,
	ethClient eth.Client,
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

// Start initiates all of Store's dependencies
func (s *Store) Start() error {
	return checkV1JobSpecs(s.DB)
}

func checkV1JobSpecs(db *gorm.DB) error {
	var count int
	if err := db.Raw(`SELECT count(*) FROM job_specs`).Scan(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		logger.Warnf(`Found %d legacy job_specs. The JSON style of job spec is now deprecated and support for jobs using this format will be REMOVED in an upcoming release. You should migrate all these jobs to V2 (TOML) format. For help doing this, please refer to the docs (https://docs.chain.link/docs/jobs/). To test your node to see how it would behave after support for these jobs is removed, you may set ENABLE_LEGACY_JOB_PIPELINE=false`, count)
	}
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

func initializeORM(cfg *config.Config, shutdownSignal gracefulpanic.Signal) (*orm.ORM, error) {
	dbURL := cfg.DatabaseURL()
	dbOrm, err := orm.NewORM(dbURL.String(), cfg.DatabaseTimeout(), shutdownSignal, cfg.GetDatabaseDialectConfiguredOrDefault(), cfg.GetAdvisoryLockIDConfiguredOrDefault(), cfg.GlobalLockRetryInterval().Duration(), cfg.ORMMaxOpenConns(), cfg.ORMMaxIdleConns())
	if err != nil {
		return nil, errors.Wrap(err, "initializeORM#NewORM")
	}
	if cfg.DatabaseBackupMode() != config.DatabaseBackupModeNone {

		version, err2 := dbOrm.FindLatestNodeVersion()
		if err2 != nil {
			return nil, errors.Wrap(err2, "initializeORM#FindLatestNodeVersion")
		}
		var versionString string
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
			return migrations.Migrate(db)
		})
		if err != nil {
			return nil, errors.Wrap(err, "initializeORM#Migrate")
		}
	}

	nodeVersion := static.Version
	if nodeVersion == "unset" {
		nodeVersion = fmt.Sprintf("random_%d", rand.Uint32())
	}
	version := models.NewNodeVersion(nodeVersion)
	err = dbOrm.UpsertNodeVersion(version)
	if err != nil {
		return nil, errors.Wrap(err, "initializeORM#UpsertNodeVersion")
	}
	dbOrm.SetLogging(cfg.LogSQLStatements())
	return dbOrm, nil
}
