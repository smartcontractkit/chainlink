package services

import (
	"time"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/store"
	"github.com/smartcontractkit/chainlink/core/store/orm"
)

type storeReaper struct {
	store  *store.Store
	config orm.Depot
}

// NewStoreReaper creates a reaper that cleans stale objects from the store.
func NewStoreReaper(store *store.Store) SleeperTask {
	return NewSleeperTask(&storeReaper{
		store:  store,
		config: store.Config,
	})
}

func (sr *storeReaper) Work() {
	offset := time.Now().Add(-sr.config.ReaperExpiration()).Add(-sr.config.SessionTimeout())
	err := sr.store.DeleteStaleSessions(offset)
	if err != nil {
		logger.Error("unable to reap stale sessions: ", err)
	}
}
