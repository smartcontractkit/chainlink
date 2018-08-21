package services

import (
	"sync"
	"time"

	"github.com/asdine/storm"
	"github.com/smartcontractkit/chainlink/logger"
	"github.com/smartcontractkit/chainlink/store"
	"github.com/smartcontractkit/chainlink/store/models"
)

// Reaper interface defines the methods used to reap stale objects such as sessions.
type Reaper interface {
	Start() error
	Stop() error
	ReapSessions()
}

type storeReaper struct {
	store   *store.Store
	config  store.Config
	cond    *sync.Cond
	started bool
}

// NewStoreReaper creates a reaper that cleans stale objects from the store.
func NewStoreReaper(store *store.Store) Reaper {
	var m sync.Mutex
	return &storeReaper{
		store:  store,
		config: store.Config,
		cond:   sync.NewCond(&m),
	}
}

// Start starts the reaper instance so that it can listen for cleanup asynchronously.
func (sr *storeReaper) Start() error {
	sr.cond.L.Lock()
	sr.started = true
	sr.cond.L.Unlock()
	go sr.listenForReaps()
	return nil
}

// Stop stops the reaper from listening to clean up messages asynchronously.
func (sr *storeReaper) Stop() error {
	sr.cond.L.Lock()
	sr.started = false
	sr.cond.Signal()
	sr.cond.L.Unlock()
	return nil
}

// ReapSessions signals the reaper to clean up sessions asynchronously.
func (sr *storeReaper) ReapSessions() {
	sr.cond.Signal()
}

func (sr *storeReaper) listenForReaps() {
	for {
		sr.cond.L.Lock()
		sr.cond.Wait()
		if sr.started == false {
			sr.cond.L.Unlock()
			return
		}
		sr.deleteStaleSessions()
		sr.cond.L.Unlock()
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
