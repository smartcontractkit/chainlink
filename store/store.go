package store

import (
	"errors"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/asdine/storm"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/smartcontractkit/chainlink/logger"
	"github.com/smartcontractkit/chainlink/store/models"
)

// Store contains fields for the database, Config, KeyStore, and TxManager
// for keeping the application state in sync with the database.
type Store struct {
	*models.ORM
	Config      Config
	Clock       AfterNower
	Exiter      func(int)
	KeyStore    *KeyStore
	TxManager   *TxManager
	HeadTracker *HeadTracker
	sigs        chan os.Signal
}

// NewStore will create a new database file at the config's RootDir if
// it is not already present, otherwise it will use the existing db.bolt
// file.
func NewStore(config Config) *Store {
	err := os.MkdirAll(config.RootDir, os.FileMode(0700))
	if err != nil {
		logger.Fatal(err)
	}
	orm := models.NewORM(config.RootDir)
	ethrpc, err := rpc.Dial(config.EthereumURL)
	if err != nil {
		logger.Fatal(err)
	}
	keyStore := NewKeyStore(config.KeysDir())

	ht, err := NewHeadTracker(orm)
	if err != nil {
		logger.Fatal(err)
	}

	store := &Store{
		ORM:         orm,
		Config:      config,
		KeyStore:    keyStore,
		Exiter:      os.Exit,
		Clock:       Clock{},
		HeadTracker: ht,
		TxManager: &TxManager{
			Config:    config,
			EthClient: &EthClient{ethrpc},
			KeyStore:  keyStore,
			ORM:       orm,
		},
	}
	return store
}

// Start listens for interrupt signals from the operating system so
// that the database can be properly closed before the application
// exits.
func (s *Store) Start() {
	s.sigs = make(chan os.Signal, 1)
	signal.Notify(s.sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-s.sigs
		s.Close()
		s.Exiter(1)
	}()
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

// Holds and stores the latest block number experienced by this particular node
// in a thread safe manner. Reconstitutes the last block number from the data
// store on reboot.
type HeadTracker struct {
	orm    *models.ORM
	number *models.IndexableBlockNumber
	mutex  sync.RWMutex
}

// Updates the latest block number from the header#Number.
func (ht *HeadTracker) SaveFromHeader(header models.BlockHeader) error {
	return ht.Save(header.IndexableBlockNumber())
}

// Updates the latest block number, if indeed the latest, and persists
// this number in case of reboot. Thread safe.
func (ht *HeadTracker) Save(n *models.IndexableBlockNumber) error {
	if n == nil {
		return errors.New("Cannot save a nil block header")
	}

	ht.mutex.Lock()
	if ht.number == nil || ht.number.ToInt().Cmp(n.ToInt()) < 0 {
		copy := *n
		ht.number = &copy
	}
	ht.mutex.Unlock()
	return ht.orm.Save(n)
}

// Returns the latest block header being tracked, or nil.
func (ht *HeadTracker) Get() *models.IndexableBlockNumber {
	ht.mutex.RLock()
	defer ht.mutex.RUnlock()
	return ht.number
}

// Instantiates a new HeadTracker using the orm to persist
// new block numbers
func NewHeadTracker(orm *models.ORM) (*HeadTracker, error) {
	ht := &HeadTracker{orm: orm}
	numbers := []models.IndexableBlockNumber{}
	err := orm.Select().OrderBy("Digits", "Number").Limit(1).Reverse().Find(&numbers)
	if err != nil && err != storm.ErrNotFound {
		return nil, err
	}
	if len(numbers) > 0 {
		ht.number = &numbers[0]
		logger.Info("Tracking logs from the last received block header ", ht.number.String())
	}
	return ht, nil
}
