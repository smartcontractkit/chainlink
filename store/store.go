package store

import (
	"context"
	"errors"
	"os"
	"path"
	"sync"
	"time"

	"github.com/coreos/bbolt"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/smartcontractkit/chainlink/logger"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/smartcontractkit/chainlink/utils"
)

// Store contains fields for the database, Config, KeyStore, and TxManager
// for keeping the application state in sync with the database.
type Store struct {
	*models.ORM
	Config     Config
	Clock      AfterNower
	KeyStore   *KeyStore
	RunChannel *RunChannel
	TxManager  *TxManager
	closed     bool
}

type rpcSubscriptionWrapper struct {
	*rpc.Client
}

func (wrapper rpcSubscriptionWrapper) EthSubscribe(ctx context.Context, channel interface{}, args ...interface{}) (models.EthSubscription, error) {
	return wrapper.Client.EthSubscribe(ctx, channel, args...)
}

// Dialer implements Dial which is a function that creates a client for that url
type Dialer interface {
	Dial(string) (CallerSubscriber, error)
}

// EthDialer is Dialer which accesses rpc urls
type EthDialer struct{}

// Dial will dial the given url and return a CallerSubscriber
func (EthDialer) Dial(url string) (CallerSubscriber, error) {
	dialed, err := rpc.Dial(url)
	if err != nil {
		return nil, err
	}
	return rpcSubscriptionWrapper{dialed}, nil
}

// NewStore will create a new database file at the config's RootDir if
// it is not already present, otherwise it will use the existing db.bolt
// file.
func NewStore(config Config) *Store {
	return NewStoreWithDialer(config, EthDialer{})
}

// NewStoreWithDialer creates a new store with the given config and dialer
func NewStoreWithDialer(config Config, dialer Dialer) *Store {
	err := os.MkdirAll(config.RootDir, os.FileMode(0700))
	if err != nil {
		logger.Fatal("Unable to create project root dir: %+v", err)
	}
	orm := initializeORM(config)
	ethrpc, err := dialer.Dial(config.EthereumURL)
	if err != nil {
		logger.Fatal("Unable to dial ETH RPC port: %+v", err)
	}
	keyStore := NewKeyStore(config.KeysDir())

	store := &Store{
		Clock:      Clock{},
		Config:     config,
		KeyStore:   keyStore,
		ORM:        orm,
		RunChannel: NewRunChannel(),
		TxManager: &TxManager{
			EthClient: &EthClient{ethrpc},
			config:    config,
			keyStore:  keyStore,
			orm:       orm,
		},
	}
	return store
}

// Start initiates all of Store's dependencies including the TxManager.
func (s *Store) Start() error {
	acc, err := s.KeyStore.GetAccount()
	if err != nil {
		return err
	}
	return s.TxManager.ActivateAccount(acc)
}

// Stop shuts down all of the working parts of the store.
func (s *Store) Stop() error {
	s.RunChannel.Close()
	return s.Close()
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

func initializeORM(config Config) *models.ORM {
	var orm *models.ORM
	var err error
	sleeper := utils.NewConstantSleeper(config.DatabasePollInterval.Duration)
	for {
		orm, err = models.NewORM(path.Join(config.RootDir, "db.bolt"))
		if err != nil && err == bolt.ErrTimeout {
			logger.Info("BoltDB is locked, sleeping", "sleepDuration", sleeper.Duration())
			sleeper.Sleep()
		} else {
			break
		}
	}
	return orm
}

// RunRequest is the type that the RunChannel uses to package all the necessary
// pieces to execute a Job Run.
type RunRequest struct {
	Input       models.RunResult
	BlockNumber *models.IndexableBlockNumber
}

// RunChannel manages accepting a queue of incoming runs.
type RunChannel struct {
	queue  chan RunRequest
	closed bool
	mutex  sync.Mutex
	Waiter sync.WaitGroup
}

// NewRunChannel initializes a RunChannel.
func NewRunChannel() *RunChannel {
	return &RunChannel{
		queue: make(chan RunRequest, 1000),
	}
}

// Send adds another entry to the queue of runs.
func (rq *RunChannel) Send(rr models.RunResult, ibn *models.IndexableBlockNumber) error {
	rq.mutex.Lock()
	defer rq.mutex.Unlock()

	if rq.closed {
		return errors.New("RunChannel.Add: cannot add to a closed RunChannel")
	}

	rq.Waiter.Add(1)
	rq.queue <- RunRequest{
		Input:       rr,
		BlockNumber: ibn,
	}
	return nil
}

// Receive returs a channel for listening to sent runs.
func (rq *RunChannel) Receive() <-chan RunRequest {
	return rq.queue
}

// Close closes the RunChannel so that no runs can be added to it without
// throwing an error.
func (rq *RunChannel) Close() {
	rq.mutex.Lock()
	defer rq.mutex.Unlock()

	rq.Waiter.Wait()

	if !rq.closed {
		rq.closed = true
		close(rq.queue)
	}
}
