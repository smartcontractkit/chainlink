package store

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/rpc"
	"github.com/pkg/errors"
	"chainlink/core/logger"
	"chainlink/core/services/synchronization"
	"chainlink/core/store/migrations"
	"chainlink/core/store/models"
	"chainlink/core/store/orm"
	"chainlink/core/utils"
	"github.com/tevino/abool"
	"go.uber.org/multierr"
	"golang.org/x/time/rate"
)

// Store contains fields for the database, Config, KeyStore, and TxManager
// for keeping the application state in sync with the database.
type Store struct {
	*orm.ORM
	Config      *orm.Config
	Clock       utils.AfterNower
	KeyStore    *KeyStore
	RunChannel  RunChannel
	TxManager   TxManager
	StatsPusher *synchronization.StatsPusher
}

type lazyRPCWrapper struct {
	client      *rpc.Client
	url         *url.URL
	mutex       *sync.Mutex
	initialized *abool.AtomicBool
	limiter     *rate.Limiter
}

func newLazyRPCWrapper(urlString string, limiter *rate.Limiter) (CallerSubscriber, error) {
	parsed, err := url.ParseRequestURI(urlString)
	if err != nil {
		return nil, err
	}
	if parsed.Scheme != "ws" && parsed.Scheme != "wss" {
		return nil, fmt.Errorf("Ethereum url scheme must be websocket: %s", parsed.String())
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

func (wrapper *lazyRPCWrapper) Call(result interface{}, method string, args ...interface{}) error {
	err := wrapper.lazyDialInitializer()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	wrapper.limiter.Wait(ctx)

	return wrapper.client.Call(result, method, args...)
}

func (wrapper *lazyRPCWrapper) EthSubscribe(ctx context.Context, channel interface{}, args ...interface{}) (models.EthSubscription, error) {
	err := wrapper.lazyDialInitializer()
	if err != nil {
		return nil, err
	}
	return wrapper.client.EthSubscribe(ctx, channel, args...)
}

// Dialer implements Dial which is a function that creates a client for that url
type Dialer interface {
	Dial(string) (CallerSubscriber, error)
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
func (ed *EthDialer) Dial(urlString string) (CallerSubscriber, error) {
	return newLazyRPCWrapper(urlString, ed.limiter)
}

// NewStore will create a new database file at the config's RootDir if
// it is not already present, otherwise it will use the existing db.sqlite3
// file.
func NewStore(config *orm.Config) *Store {
	return NewStoreWithDialer(config, NewEthDialer(config.MaxRPCCallsPerSecond()))
}

// NewStoreWithDialer creates a new store with the given config and dialer
func NewStoreWithDialer(config *orm.Config, dialer Dialer) *Store {
	err := os.MkdirAll(config.RootDir(), os.FileMode(0700))
	if err != nil {
		logger.Fatal(fmt.Sprintf("Unable to create project root dir: %+v", err))
	}
	orm, err := initializeORM(config)
	if err != nil {
		logger.Fatal(fmt.Sprintf("Unable to initialize ORM: %+v", err))
	}
	ethrpc, err := dialer.Dial(config.EthereumURL())
	if err != nil {
		logger.Fatal(fmt.Sprintf("Unable to dial ETH RPC port: %+v", err))
	}
	if err := orm.ClobberDiskKeyStoreWithDBKeys(config.KeysDir()); err != nil {
		logger.Fatal(fmt.Sprintf("Unable to migrate key store to disk: %+v", err))
	}
	keyStore := NewKeyStore(config.KeysDir())

	store := &Store{
		Clock:       utils.Clock{},
		Config:      config,
		KeyStore:    keyStore,
		ORM:         orm,
		RunChannel:  NewQueuedRunChannel(),
		TxManager:   NewEthTxManager(&EthCallerSubscriber{ethrpc}, config, keyStore, orm),
		StatsPusher: synchronization.NewStatsPusher(orm, config.ExplorerURL(), config.ExplorerAccessKey(), config.ExplorerSecret()),
	}
	return store
}

// Start initiates all of Store's dependencies including the TxManager.
func (s *Store) Start() error {
	s.TxManager.Register(s.KeyStore.Accounts())
	return multierr.Combine(
		s.SyncDiskKeyStoreToDB(),
		s.StatsPusher.Start(),
	)
}

// Close shuts down all of the working parts of the store.
func (s *Store) Close() error {
	s.RunChannel.Close()
	return multierr.Combine(
		s.ORM.Close(),
		s.StatsPusher.Close(),
	)
}

// Unscoped returns a shallow copy of the store, with an unscoped ORM allowing
// one to work with soft deleted records.
func (s *Store) Unscoped() *Store {
	cpy := *s
	cpy.ORM = cpy.ORM.Unscoped()
	return &cpy
}

// AuthorizedUserWithSession will return the one API user if the Session ID exists
// and hasn't expired, and update session's LastUsed field.
func (s *Store) AuthorizedUserWithSession(sessionID string) (models.User, error) {
	return s.ORM.AuthorizedUserWithSession(sessionID, s.Config.SessionTimeout())
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

		err = s.FirstOrCreateKey(key)
		if err != nil {
			merr = multierr.Append(err, merr)
		}
	}
	return merr
}

func initializeORM(config *orm.Config) (*orm.ORM, error) {
	orm, err := orm.NewORM(orm.NormalizedDatabaseURL(config), config.DatabaseTimeout())
	if err != nil {
		return nil, errors.Wrap(err, "initializeORM#NewORM")
	}
	orm.SetLogging(config.LogSQLStatements() || config.LogSQLMigrations())
	if err = migrations.Migrate(orm.DB); err != nil {
		return nil, errors.Wrap(err, "initializeORM#Migrate")
	}
	orm.SetLogging(config.LogSQLStatements())
	return orm, nil
}

// RunRequest is the type that the RunChannel uses to package all the necessary
// pieces to execute a Job Run.
type RunRequest struct {
	ID *models.ID
}

// RunChannel manages and dispatches incoming runs.
type RunChannel interface {
	Send(jobRunID *models.ID) error
	Receive() <-chan RunRequest
	Close()
}

// QueuedRunChannel manages incoming results and blocks by enqueuing them
// in a queue per run.
type QueuedRunChannel struct {
	queue  chan RunRequest
	closed bool
	mutex  sync.Mutex
}

// NewQueuedRunChannel initializes a QueuedRunChannel.
func NewQueuedRunChannel() RunChannel {
	return &QueuedRunChannel{
		queue: make(chan RunRequest, 1000),
	}
}

// Send adds another entry to the queue of runs.
func (rq *QueuedRunChannel) Send(jobRunID *models.ID) error {
	rq.mutex.Lock()
	defer rq.mutex.Unlock()

	if rq.closed {
		return errors.New("QueuedRunChannel.Add: cannot add to a closed QueuedRunChannel")
	}

	if jobRunID == nil {
		return errors.New("QueuedRunChannel.Add: cannot add an empty jobRunID")
	}

	rq.queue <- RunRequest{ID: jobRunID}
	return nil
}

// Receive returns a channel for listening to sent runs.
func (rq *QueuedRunChannel) Receive() <-chan RunRequest {
	return rq.queue
}

// Close closes the QueuedRunChannel so that no runs can be added to it without
// throwing an error.
func (rq *QueuedRunChannel) Close() {
	rq.mutex.Lock()
	defer rq.mutex.Unlock()

	if !rq.closed {
		rq.closed = true
		close(rq.queue)
	}
}
