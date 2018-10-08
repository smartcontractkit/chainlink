package store

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/rpc"
	"github.com/smartcontractkit/chainlink/logger"
	"github.com/smartcontractkit/chainlink/store/migrations"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/smartcontractkit/chainlink/store/orm"
	"go.uber.org/multierr"
)

// Store contains fields for the database, Config, KeyStore, and TxManager
// for keeping the application state in sync with the database.
type Store struct {
	*orm.ORM
	Config     Config
	Clock      AfterNower
	KeyStore   *KeyStore
	RunChannel RunChannel
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
		logger.Fatal(fmt.Sprintf("Unable to create project root dir: %+v", err))
	}
	orm, err := initializeORM(config)
	if err != nil {
		logger.Fatal(fmt.Sprintf("Unable to initialize ORM: %+v", err))
	}
	ethrpc, err := dialer.Dial(config.EthereumURL)
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

	err = s.TxManager.ActivateAccount(acc)
	if err != nil {
		return err
	}

	return s.cleanUpAbruptShutdown()
}

func (s *Store) cleanUpAbruptShutdown() error {
	runs, err := s.recoverInProgress()
	if err != nil {
		return err
	}

	return s.resumeInProgress(runs)
}

func (s *Store) recoverInProgress() ([]models.JobRun, error) {
	runs, err := s.JobRunsWithStatus(models.RunStatusInProgress, models.RunStatusUnstarted)
	if err != nil {
		return runs, err
	}

	var merr error
	for _, jr := range runs {
		jr.Status = models.RunStatusUnstarted
		multierr.Append(merr, s.Save(&jr))
	}
	return runs, merr
}

func (s *Store) resumeInProgress(runs []models.JobRun) error {
	for _, run := range runs {
		if err := s.RunChannel.Send(run.ID, nil); err != nil {
			return err
		}
	}
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
	return s.ORM.AuthorizedUserWithSession(sessionID, s.Config.SessionTimeout.Duration)
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
	path := path.Join(config.RootDir, "db.bolt")
	duration := config.DatabaseTimeout.Duration
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
	ID          string
	BlockNumber *models.IndexableBlockNumber
}

// RunChannel manages and dispatches incoming runs.
type RunChannel interface {
	Send(jobRunID string, ibn *models.IndexableBlockNumber) error
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
func (rq *QueuedRunChannel) Send(jobRunID string, ibn *models.IndexableBlockNumber) error {
	rq.mutex.Lock()
	defer rq.mutex.Unlock()

	if rq.closed {
		return errors.New("QueuedRunChannel.Add: cannot add to a closed QueuedRunChannel")
	}

	if jobRunID == "" {
		return errors.New("QueuedRunChannel.Add: cannot add an empty jobRunID")
	}

	rq.queue <- RunRequest{
		ID:          jobRunID,
		BlockNumber: ibn,
	}
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
