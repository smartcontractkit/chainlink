package store

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/smartcontractkit/chainlink/core/gracefulpanic"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/eth"
	"github.com/smartcontractkit/chainlink/core/services/offchainreporting"
	"github.com/smartcontractkit/chainlink/core/services/postgres"
	"github.com/smartcontractkit/chainlink/core/store/migrations"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/store/orm"
	"github.com/smartcontractkit/chainlink/core/utils"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"go.uber.org/multierr"
)

const (
	// AutoMigrate is a flag that automatically migrates the DB when passed to initializeORM
	AutoMigrate = "auto_migrate"
)

// NotifyNewEthTx allows to notify the ethBroadcaster of a new transaction
//go:generate mockery --name NotifyNewEthTx --output ../internal/mocks/ --case=underscore
type NotifyNewEthTx interface {
	Trigger()
}

// Store contains fields for the database, Config, and KeyStore
// for keeping the application state in sync with the database.
type Store struct {
	*orm.ORM
	Config         *orm.Config
	Clock          utils.AfterNower
	KeyStore       KeyStoreInterface
	VRFKeyStore    *VRFKeyStore
	OCRKeyStore    *offchainreporting.KeyStore
	EthClient      eth.Client
	NotifyNewEthTx NotifyNewEthTx
	AdvisoryLocker postgres.AdvisoryLocker
	closeOnce      *sync.Once
}

// NewStore will create a new store
func NewStore(config *orm.Config, ethClient eth.Client, advisoryLock postgres.AdvisoryLocker, shutdownSignal gracefulpanic.Signal) *Store {
	keyStore := func() *KeyStore {
		scryptParams := utils.GetScryptParams(config)
		return NewKeyStore(config.KeysDir(), scryptParams)
	}
	return newStoreWithKeyStore(config, ethClient, advisoryLock, keyStore, shutdownSignal)
}

// NewInsecureStore creates a new store with the given config using an insecure keystore.
// NOTE: Should only be used for testing!
func NewInsecureStore(config *orm.Config, ethClient eth.Client, advisoryLocker postgres.AdvisoryLocker, shutdownSignal gracefulpanic.Signal) *Store {
	keyStore := func() *KeyStore { return NewInsecureKeyStore(config.KeysDir()) }
	return newStoreWithKeyStore(config, ethClient, advisoryLocker, keyStore, shutdownSignal)
}

// TODO(sam): Remove ethClient from here completely after legacy tx manager is gone
func newStoreWithKeyStore(
	config *orm.Config,
	ethClient eth.Client,
	advisoryLocker postgres.AdvisoryLocker,
	keyStoreGenerator func() *KeyStore,
	shutdownSignal gracefulpanic.Signal,
) *Store {
	if err := utils.EnsureDirAndMaxPerms(config.RootDir(), os.FileMode(0700)); err != nil {
		logger.Fatal(fmt.Sprintf("Unable to create project root dir: %+v", err))
	}
	orm, err := initializeORM(config, shutdownSignal)
	if err != nil {
		logger.Fatal(fmt.Sprintf("Unable to initialize ORM: %+v", err))
	}
	if e := orm.ClobberDiskKeyStoreWithDBKeys(config.KeysDir()); e != nil {
		logger.Fatal(fmt.Sprintf("Unable to migrate key store to disk: %+v", e))
	}

	keyStore := keyStoreGenerator()
	scryptParams := utils.GetScryptParams(config)

	store := &Store{
		Clock:          utils.Clock{},
		AdvisoryLocker: advisoryLocker,
		Config:         config,
		KeyStore:       keyStore,
		OCRKeyStore:    offchainreporting.NewKeyStore(orm.DB, scryptParams),
		ORM:            orm,
		EthClient:      ethClient,
		closeOnce:      &sync.Once{},
	}
	store.VRFKeyStore = NewVRFKeyStore(store)
	return store
}

// Start initiates all of Store's dependencies
func (s *Store) Start() error {
	return s.SyncDiskKeyStoreToDB()
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

// SyncDiskKeyStoreToDB writes all keys in the keys directory to the underlying
// orm.
func (s *Store) SyncDiskKeyStoreToDB() error {
	files, err := utils.FilesInDir(s.Config.KeysDir())
	if err != nil {
		return multierr.Append(errors.New("unable to sync disk keystore to db"), err)
	}

	var merr error
	for _, f := range files {
		key, err := models.NewKeyFromFile(filepath.Join(s.Config.KeysDir(), f))
		if err != nil {
			merr = multierr.Append(err, merr)
			continue
		}

		err = s.CreateKeyIfNotExists(key)
		if err != nil {
			merr = multierr.Append(err, merr)
		}
	}
	return merr
}

func initializeORM(config *orm.Config, shutdownSignal gracefulpanic.Signal) (*orm.ORM, error) {
	orm, err := orm.NewORM(config.DatabaseURL(), config.DatabaseTimeout(), shutdownSignal, config.GetDatabaseDialectConfiguredOrDefault(), config.GetAdvisoryLockIDConfiguredOrDefault())
	if err != nil {
		return nil, errors.Wrap(err, "initializeORM#NewORM")
	}
	if config.MigrateDatabase() {
		orm.SetLogging(config.LogSQLStatements() || config.LogSQLMigrations())

		err = orm.RawDB(func(db *gorm.DB) error {
			return migrations.Migrate(db)
		})
		if err != nil {
			return nil, errors.Wrap(err, "initializeORM#Migrate")
		}
	}
	orm.SetLogging(config.LogSQLStatements())
	return orm, nil
}
