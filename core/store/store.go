package store

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/smartcontractkit/chainlink/core/gracefulpanic"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/eth"
	"github.com/smartcontractkit/chainlink/core/store/migrations"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/store/orm"
	"github.com/smartcontractkit/chainlink/core/utils"

	gethClient "github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"github.com/tevino/abool"
	"go.uber.org/multierr"
	"golang.org/x/time/rate"
)

const (
	// AutoMigrate is a flag that automatically migrates the DB when passed to initializeORM
	AutoMigrate = "auto_migrate"
)

type GethClientWrapper interface {
	GethClient(func(gethClient eth.GethClient) error) error
	RPCClient(func(rpcClient eth.RPCClient) error) error
}

// NotifyNewEthTx allows to notify the ethBroadcaster of a new transaction
//go:generate mockery --name NotifyNewEthTx --output ../internal/mocks/ --case=underscore
type NotifyNewEthTx interface {
	Trigger()
}

// Store contains fields for the database, Config, KeyStore, and TxManager
// for keeping the application state in sync with the database.
type Store struct {
	*orm.ORM
	Config            *orm.Config
	Clock             utils.AfterNower
	KeyStore          KeyStoreInterface
	VRFKeyStore       *VRFKeyStore
	TxManager         TxManager
	GethClientWrapper GethClientWrapper
	NotifyNewEthTx    NotifyNewEthTx
	closeOnce         *sync.Once
}

type lazyRPCWrapper struct {
	client      *rpc.Client
	url         *url.URL
	mutex       *sync.Mutex
	initialized *abool.AtomicBool
	limiter     *rate.Limiter
}

func newLazyRPCWrapper(urlString string, limiter *rate.Limiter) (eth.CallerSubscriber, error) {
	parsed, err := url.ParseRequestURI(urlString)
	if err != nil {
		return nil, err
	}
	if parsed.Scheme != "ws" && parsed.Scheme != "wss" {
		return nil, fmt.Errorf("ethereum url scheme must be websocket: %s", parsed.String())
	}
	return &lazyRPCWrapper{
		url:         parsed,
		mutex:       &sync.Mutex{},
		initialized: abool.New(),
		limiter:     limiter,
	}, nil
}

// lazyDialInitializer initializes the Dial instance used to interact with
// an ethereum node using the Double-checked locking optimization:
// https://en.wikipedia.org/wiki/Double-checked_locking
func (wrapper *lazyRPCWrapper) lazyDialInitializer() error {
	if wrapper.initialized.IsSet() {
		return nil
	}

	wrapper.mutex.Lock()
	defer wrapper.mutex.Unlock()

	if wrapper.client == nil {
		client, err := rpc.Dial(wrapper.url.String())
		if err != nil {
			return err
		}
		wrapper.client = client
		wrapper.initialized.Set()
	}
	return nil
}

// GethClient allows callers to access go-ethereum's ethclient through the
// wrapper's rate limiting
func (wrapper *lazyRPCWrapper) GethClient(callback func(gethClient eth.GethClient) error) error {
	err := wrapper.lazyDialInitializer()
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	logger.ErrorIf(wrapper.limiter.Wait(ctx))

	client := gethClient.NewClient(wrapper.client)

	return callback(client)
}

// RPCClient allows callers to access go-ethereum's rpcclient through the
// wrapper's rate limiting
func (wrapper *lazyRPCWrapper) RPCClient(callback func(rpcClient eth.RPCClient) error) error {
	err := wrapper.lazyDialInitializer()
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	logger.ErrorIf(wrapper.limiter.Wait(ctx))

	return callback(wrapper.client)
}

func (wrapper *lazyRPCWrapper) Call(result interface{}, method string, args ...interface{}) error {
	err := wrapper.lazyDialInitializer()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	err = wrapper.limiter.Wait(ctx)
	if err != nil {
		return err
	}

	return wrapper.client.Call(result, method, args...)
}

func (wrapper *lazyRPCWrapper) Subscribe(ctx context.Context, channel interface{}, args ...interface{}) (eth.Subscription, error) {
	err := wrapper.lazyDialInitializer()
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	err = wrapper.limiter.Wait(ctx)
	if err != nil {
		return nil, err
	}

	return wrapper.client.EthSubscribe(ctx, channel, args...)
}

// Dialer implements Dial which is a function that creates a client for that url
type Dialer interface {
	Dial(string) (eth.CallerSubscriber, error)
}

// EthDialer is Dialer which accesses rpc urls
type EthDialer struct {
	limiter *rate.Limiter
}

// NewEthDialer returns an eth dialer with the specified rate limit
func NewEthDialer(rateLimit uint64) *EthDialer {
	return &EthDialer{
		limiter: rate.NewLimiter(rate.Limit(rateLimit), 1),
	}
}

// Dial will dial the given url and return a CallerSubscriber
func (ed *EthDialer) Dial(urlString string) (eth.CallerSubscriber, error) {
	return newLazyRPCWrapper(urlString, ed.limiter)
}

// NewStore will create a new store using the Eth dialer
func NewStore(config *orm.Config, shutdownSignal gracefulpanic.Signal) *Store {
	keyStore := func() *KeyStore { return NewKeyStore(config.KeysDir()) }
	dialer := NewEthDialer(config.MaxRPCCallsPerSecond())
	return newStoreWithDialerAndKeyStore(config, dialer, keyStore, shutdownSignal)
}

// NewInsecureStore creates a new store with the given config and
// dialer, using an insecure keystore.
// NOTE: Should only be used for testing!
func NewInsecureStore(config *orm.Config, shutdownSignal gracefulpanic.Signal) *Store {
	dialer := NewEthDialer(config.MaxRPCCallsPerSecond())
	keyStore := func() *KeyStore { return NewInsecureKeyStore(config.KeysDir()) }
	return newStoreWithDialerAndKeyStore(config, dialer, keyStore, shutdownSignal)
}

func newStoreWithDialerAndKeyStore(
	config *orm.Config,
	dialer Dialer,
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
	ethrpc, err := dialer.Dial(config.EthereumURL())
	if err != nil {
		logger.Fatal(fmt.Sprintf("Unable to dial ETH RPC port: %+v", err))
	}
	if e := orm.ClobberDiskKeyStoreWithDBKeys(config.KeysDir()); e != nil {
		logger.Fatal(fmt.Sprintf("Unable to migrate key store to disk: %+v", e))
	}

	keyStore := keyStoreGenerator()
	callerSubscriberClient := &eth.CallerSubscriberClient{CallerSubscriber: ethrpc}
	txManager := NewEthTxManager(callerSubscriberClient, config, keyStore, orm)

	if err != nil {
		logger.Fatalf("Unable to dial ETH client: %+v", err)
	}

	store := &Store{
		Clock:             utils.Clock{},
		Config:            config,
		KeyStore:          keyStore,
		ORM:               orm,
		TxManager:         txManager,
		GethClientWrapper: ethrpc,
		closeOnce:         &sync.Once{},
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

		err = s.CreateKeyIfNotExists(key)
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
