package store

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/smartcontractkit/chainlink/core/gracefulpanic"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/eth"
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

// Store contains fields for the database, Config, KeyStore, and TxManager
// for keeping the application state in sync with the database.
type Store struct {
	*orm.ORM
	Config         *orm.Config
	Clock          utils.AfterNower
	KeyStore       KeyStoreInterface
	VRFKeyStore    *VRFKeyStore
	TxManager      TxManager
	EthClient      eth.Client
	NotifyNewEthTx NotifyNewEthTx
	closeOnce      *sync.Once
}

// NewStore will create a new store
func NewStore(config *orm.Config, shutdownSignal gracefulpanic.Signal) *Store {
	keyStore := func() *KeyStore { return NewKeyStore(config.KeysDir()) }
	return newStoreWithKeyStore(config, keyStore, shutdownSignal)
}

// NewInsecureStore creates a new store with the given config using an insecure keystore.
// NOTE: Should only be used for testing!
func NewInsecureStore(config *orm.Config, shutdownSignal gracefulpanic.Signal) *Store {
	keyStore := func() *KeyStore { return NewInsecureKeyStore(config.KeysDir()) }
	return newStoreWithKeyStore(config, keyStore, shutdownSignal)
}

func newStoreWithKeyStore(
	config *orm.Config,
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

	ethClient, err := eth.NewClient(config.EthereumURL())
	if err != nil {
		logger.Fatal(fmt.Sprintf("Unable to create ETH client: %+v", err))
	}
	keyStore := keyStoreGenerator()
	txManager := NewEthTxManager(ethClient, config, keyStore, orm)

	store := &Store{
		Clock:     utils.Clock{},
		Config:    config,
		KeyStore:  keyStore,
		ORM:       orm,
		TxManager: txManager,
		EthClient: ethClient,
		closeOnce: &sync.Once{},
	}
	store.VRFKeyStore = NewVRFKeyStore(store)
	return store
}

// Start initiates all of Store's dependencies including the TxManager.
func (s *Store) Start() error {
	if s.Config.EnableBulletproofTxManager() {
		if err := setNonceFromLegacyTxManager(s.DB); err != nil {
			return err
		}
	} else {
		s.TxManager.Register(s.KeyStore.Accounts())
	}

	return s.SyncDiskKeyStoreToDB()
}

func setNonceFromLegacyTxManager(db *gorm.DB) error {
	return db.Exec(`
	UPDATE keys
	SET next_nonce = (SELECT max(nonce) FROM txes WHERE txes.from = keys.address)+1
	WHERE next_nonce < (SELECT max(nonce) FROM txes WHERE txes.from = keys.address)+1;
	`).Error
}

// Close shuts down all of the working parts of the store.
func (s *Store) Close() error {
	var err error
	s.closeOnce.Do(func() {
		err = s.ORM.Close()
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

		err = s.UpsertKey(key)
		if err != nil {
			fmt.Println("Balls", err)
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
