package services

import (
	"sync"
	"time"

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
	semaphore chan struct{}
}

// NewStoreReaper creates a reaper that cleans stale objects from the store.
func NewStoreReaper(store *store.Store) Reaper {
	return &storeReaper{
		store:  store,
		config: store.Config,
	}
}

// Start starts the reaper instance so that it can listen for cleanup asynchronously.
func (sr *storeReaper) Start() error {
	sr.bootMutex.Lock()
	defer sr.bootMutex.Unlock()
	sr.semaphore = make(chan struct{}, 1)
	sr.semaphore <- struct{}{}
	return nil
}

// Stop stops the reaper from listening to clean up messages asynchronously.
func (sr *storeReaper) Stop() error {
	sr.bootMutex.Lock()
	defer sr.bootMutex.Unlock()

	if sr.semaphore != nil {
		close(sr.semaphore)
		sr.semaphore = nil
	}
	return nil
}

// ReapSessions signals the reaper to clean up sessions asynchronously.
func (sr *storeReaper) ReapSessions() {
	go sr.reapOrSkip()
}

func (sr *storeReaper) reapOrSkip() {
	select {
	case _, ok := <-sr.semaphore:
		if ok {
			sr.deleteStaleSessions()
			sr.semaphore <- struct{}{}
		}
		return
	default: // skip
	}
}

func (sr *storeReaper) deleteStaleSessions() {
	var sessions []models.Session
	offset := time.Now().Add(-sr.config.ReaperExpiration.Duration).Add(-sr.config.SessionTimeout.Duration)
	stale := models.Time{offset}
	err := sr.store.Range("LastUsed", models.Time{}, stale, &sessions)
	if err != nil {
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
