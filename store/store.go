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

// Holds and stores the latest block header experienced by this particular node
// in a thread safe manner. Reconstitutes the last block header from the data
// store on reboot.
type HeadTracker struct {
	orm         *models.ORM
	blockHeader *models.BlockHeader
	mutex       sync.RWMutex
}

// Updates the latest block header, if indeed the latest, and persists
// this block header in case of reboot. Thread safe.
func (ht *HeadTracker) Save(bh *models.BlockHeader) error {
	if bh == nil {
		return errors.New("Cannot save a nil block header")
	}

	ht.mutex.Lock()
	if ht.blockHeader == nil || ht.blockHeader.Number.ToInt().Cmp(bh.Number.ToInt()) < 0 {
		copy := *bh
		ht.blockHeader = &copy
	}
	ht.mutex.Unlock()
	return ht.orm.Save(bh)
}

// Returns the latest block header being tracked, or nil.
func (ht *HeadTracker) Get() *models.BlockHeader {
	ht.mutex.RLock()
	defer ht.mutex.RUnlock()
	return ht.blockHeader
}

// Instantiates a new HeadTracker using the orm to persist
// new BlockHeaders
func NewHeadTracker(orm *models.ORM) (*HeadTracker, error) {
	ht := &HeadTracker{orm: orm}
	blockHeaders := []models.BlockHeader{}
	err := orm.AllByIndex("Number", &blockHeaders, storm.Limit(1), storm.Reverse())
	if err != nil {
		return nil, err
	}
	if len(blockHeaders) > 0 {
		ht.blockHeader = &blockHeaders[0]
	}
	return ht, nil
}
