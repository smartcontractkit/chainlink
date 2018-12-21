package store

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"os"
	"path"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/rpc"
	"github.com/smartcontractkit/chainlink/logger"
	"github.com/smartcontractkit/chainlink/store/migrations"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/smartcontractkit/chainlink/store/orm"
	"github.com/tevino/abool"
)

// Store contains fields for the database, Config, KeyStore, and TxManager
// for keeping the application state in sync with the database.
type Store struct {
	*orm.ORM
	Config     Config
	Clock      AfterNower
	KeyStore   *KeyStore
	RunChannel RunChannel
	TxManager  TxManager
	closed     bool
}

type lazyRPCWrapper struct {
	client      *rpc.Client
	url         *url.URL
	mutex       *sync.Mutex
	initialized *abool.AtomicBool
}

func newLazyRPCWrapper(urlString string) (CallerSubscriber, error) {
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
	url models.WebURL
}

// Dial will dial the given url and return a CallerSubscriber
func (ed *EthDialer) Dial(urlString string) (CallerSubscriber, error) {
	return newLazyRPCWrapper(urlString)
}

// NewStore will create a new database file at the config's RootDir if
// it is not already present, otherwise it will use the existing db.bolt
// file.
func NewStore(config Config) *Store {
	return NewStoreWithDialer(config, &EthDialer{})
}

// NewStoreWithDialer creates a new store with the given config and dialer
func NewStoreWithDialer(config Config, dialer Dialer) *Store {
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
	keyStore := NewKeyStore(config.KeysDir())

	store := &Store{
		Clock:      Clock{},
		Config:     config,
		KeyStore:   keyStore,
		ORM:        orm,
		RunChannel: NewQueuedRunChannel(),
		TxManager:  NewEthTxManager(&EthClient{ethrpc}, config, keyStore, orm),
	}
	return store
}

// Start initiates all of Store's dependencies including the TxManager.
func (s *Store) Start() error {
	s.TxManager.Register(s.KeyStore.Accounts())
	return nil
}

// Close shuts down all of the working parts of the store.
func (s *Store) Close() error {
	s.RunChannel.Close()
	return s.ORM.Close()
}

// AuthorizedUserWithSession will return the one API user if the Session ID exists
// and hasn't expired, and update session's LastUsed field.
func (s *Store) AuthorizedUserWithSession(sessionID string) (models.User, error) {
	return s.ORM.AuthorizedUserWithSession(sessionID, s.Config.SessionTimeout())
}

// AfterNower is an interface that fulfills the `After()` and `Now()`
// methods.
type AfterNower interface {
	After(d time.Duration) <-chan time.Time
	Now() time.Time
}

// Clock is a basic type for scheduling events in the application.
type Clock struct{}

// Now returns the current time.
func (Clock) Now() time.Time {
	return time.Now()
}

// After returns the current time if the given duration has elapsed.
func (Clock) After(d time.Duration) <-chan time.Time {
	return time.After(d)
}

func initializeORM(config Config) (*orm.ORM, error) {
	path := path.Join(config.RootDir(), "db.bolt")
	duration := config.DatabaseTimeout()
	logger.Infof("Waiting %s for lock on db file %s", friendlyDuration(duration), path)
	orm, err := orm.NewORM(path, duration)
	if err != nil {
		return nil, err
	}
	return orm, migrations.Migrate(orm)
}

const zeroDuration = time.Duration(0)

func friendlyDuration(duration time.Duration) string {
	if duration == zeroDuration {
		return "indefinitely"
	}
	return fmt.Sprintf("%v", duration)
}

// RunRequest is the type that the RunChannel uses to package all the necessary
// pieces to execute a Job Run.
type RunRequest struct {
	ID string
}

// RunChannel manages and dispatches incoming runs.
type RunChannel interface {
	Send(jobRunID string) error
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
func (rq *QueuedRunChannel) Send(jobRunID string) error {
	rq.mutex.Lock()
	defer rq.mutex.Unlock()

	if rq.closed {
		return errors.New("QueuedRunChannel.Add: cannot add to a closed QueuedRunChannel")
	}

	if jobRunID == "" {
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
