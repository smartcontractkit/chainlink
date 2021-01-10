package store

import (
	"context"
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

	"github.com/ethereum/go-ethereum/common"
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

type KeyStoreGenerator func(*orm.Config) *KeyStore

func StandardKeyStoreGen(config *orm.Config) *KeyStore {
	scryptParams := utils.GetScryptParams(config)
	return NewKeyStore(config.KeysDir(), scryptParams)
}

func InsecureKeyStoreGen(config *orm.Config) *KeyStore {
	return NewInsecureKeyStore(config.KeysDir())
}

// NewStore will create a new store
func NewStore(config *orm.Config, ethClient eth.Client, advisoryLock postgres.AdvisoryLocker, shutdownSignal gracefulpanic.Signal, keyStoreGenerator KeyStoreGenerator) *Store {
	return newStoreWithKeyStore(config, ethClient, advisoryLock, keyStoreGenerator, shutdownSignal)
}

// NewInsecureStore creates a new store with the given config using an insecure keystore.
// NOTE: Should only be used for testing!
func NewInsecureStore(config *orm.Config, ethClient eth.Client, advisoryLocker postgres.AdvisoryLocker, shutdownSignal gracefulpanic.Signal) *Store {
	return newStoreWithKeyStore(config, ethClient, advisoryLocker, InsecureKeyStoreGen, shutdownSignal)
}

// TODO(sam): Remove ethClient from here completely after legacy tx manager is gone
// See: https://www.pivotaltracker.com/story/show/175493792
func newStoreWithKeyStore(
	config *orm.Config,
	ethClient eth.Client,
	advisoryLocker postgres.AdvisoryLocker,
	keyStoreGenerator KeyStoreGenerator,
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

	keyStore := keyStoreGenerator(config)
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

// DeleteKey hard-deletes a key whose address matches the supplied address.
func (s *Store) DeleteKey(address common.Address) error {
	return postgres.GormTransaction(context.Background(), s.ORM.DB, func(tx *gorm.DB) error {
		err := tx.Where("address = ?", address).Delete(models.Key{}).Error
		if err != nil {
			return errors.Wrap(err, "while deleting ETH key from DB")
		}
		return s.KeyStore.Delete(address)
	})
}

// ArchiveKey soft-deletes a key whose address matches the supplied address.
func (s *Store) ArchiveKey(address common.Address) error {
	err := s.ORM.DB.Where("address = ?", address).Delete(models.Key{}).Error
	if err != nil {
		return err
	}

	acct, err := s.KeyStore.GetAccountByAddress(address)
	if err != nil {
		return err
	}

	archivedKeysDir := filepath.Join(s.Config.RootDir(), "archivedkeys")
	err = utils.EnsureDirAndMaxPerms(archivedKeysDir, os.FileMode(0700))
	if err != nil {
		return errors.Wrap(err, "could not create "+archivedKeysDir)
	}

	basename := filepath.Base(acct.URL.Path)
	dst := filepath.Join(archivedKeysDir, basename)
	err = utils.CopyFileWithMaxPerms(acct.URL.Path, dst, os.FileMode(0700))
	if err != nil {
		return errors.Wrap(err, "could not copy "+acct.URL.Path+" to "+dst)
	}

	return s.KeyStore.Delete(address)
}

func (s *Store) ImportKey(keyJSON []byte, oldPassword string) error {
	return postgres.GormTransaction(context.Background(), s.ORM.DB, func(tx *gorm.DB) error {
		_, err := s.KeyStore.Import(keyJSON, oldPassword)
		if err != nil {
			return err
		}
		return s.SyncDiskKeyStoreToDB()
	})
}

func initializeORM(config *orm.Config, shutdownSignal gracefulpanic.Signal) (*orm.ORM, error) {
	orm, err := orm.NewORM(config.DatabaseURL(), config.DatabaseTimeout(), shutdownSignal, config.GetDatabaseDialectConfiguredOrDefault(), config.GetAdvisoryLockIDConfiguredOrDefault(), config.GlobalLockRetryInterval().Duration())
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
