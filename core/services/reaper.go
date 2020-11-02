package services

import (
	"time"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/store"
	"github.com/smartcontractkit/chainlink/core/store/orm"
	"github.com/smartcontractkit/chainlink/core/utils"
)

type storeReaper struct {
	store  *store.Store
	config orm.ConfigReader
}

// NewStoreReaper creates a reaper that cleans stale objects from the store.
func NewStoreReaper(store *store.Store) utils.SleeperTask {
	return utils.NewSleeperTask(&storeReaper{
		store:  store,
		config: store.Config,
	})
}

func (sr *storeReaper) Work() {
	recordCreationStaleThreshold := sr.config.ReaperExpiration().Before(
		sr.config.SessionTimeout().Before(time.Now()))
	err := sr.store.DeleteStaleSessions(recordCreationStaleThreshold)
	if err != nil {
		logger.Error("unable to reap stale sessions: ", err)
	}
}
