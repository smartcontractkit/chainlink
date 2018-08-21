package services

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/asdine/storm"
	"github.com/smartcontractkit/chainlink/logger"
	"github.com/smartcontractkit/chainlink/store"
	"github.com/smartcontractkit/chainlink/store/models"
)

// Reaper interface is a gateway to an instance that can reap stale objects such as sessions.
type Reaper interface {
	Start() error
	Stop() error
	ReapSessions()
}

type storeReaper struct {
	store     *store.Store
	config    store.Config
	bootMutex sync.Mutex
	started   bool
	listener  chan struct{}
	semaphore singleSemaphore
}

// NewStoreReaper creates a reaper that cleans stale objects from the store.
func NewStoreReaper(store *store.Store) Reaper {
	return &storeReaper{
		store:    store,
		config:   store.Config,
		listener: make(chan struct{}, 1),
	}
}

// Start starts the reaper instance so that it can listen for cleanup asynchronously.
func (sr *storeReaper) Start() error {
	sr.bootMutex.Lock()
	defer sr.bootMutex.Unlock()

	sr.listener = make(chan struct{}, 1)
	go sr.listenForReaps()
	sr.started = true
	return nil
}

// Stop stops the reaper from listening to clean up messages asynchronously.
func (sr *storeReaper) Stop() error {
	sr.bootMutex.Lock()
	defer sr.bootMutex.Unlock()

	if sr.started {
		close(sr.listener)
		sr.started = false
	}

	return nil
}

// ReapSessions signals the reaper to clean up sessions asynchronously.
func (sr *storeReaper) ReapSessions() {
	if sr.semaphore.CanRun() {
		sr.listener <- struct{}{}
	}
}

func (sr *storeReaper) listenForReaps() {
	for {
		select {
		case _, ok := <-sr.listener:
			defer sr.semaphore.Done()
			if !ok {
				return
			}
			sr.deleteStaleSessions()
		}
	}
}

func (sr *storeReaper) deleteStaleSessions() {
	var sessions []models.Session
	offset := time.Now().Add(-sr.config.ReaperExpiration.Duration).Add(-sr.config.SessionTimeout.Duration)
	stale := models.Time{offset}
	err := sr.store.Range("LastUsed", models.Time{}, stale, &sessions)
	if err != nil && err != storm.ErrNotFound {
		logger.Error("unable to reap stale sessions: ", err)
		return
	}

	for _, s := range sessions {
		err := sr.store.DeleteStruct(&s)
		if err != nil {
			logger.Error("unable to delete stale session: ", err)
		}
	}
}

// singleSemaphore constrains consumption to a single consumer.
type singleSemaphore struct {
	semaphore int32
}

func (l *singleSemaphore) CanRun() bool {
	return atomic.CompareAndSwapInt32(&l.semaphore, 0, 1)
}
func (l *singleSemaphore) Done() {
	atomic.CompareAndSwapInt32(&l.semaphore, 1, 0)
}
