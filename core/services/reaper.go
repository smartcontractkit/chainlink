package services

import (
	"time"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/store"
	"github.com/smartcontractkit/chainlink/core/store/orm"
	"github.com/smartcontractkit/chainlink/core/utils"
)

type sessionReaper struct {
	store  *store.Store
	config orm.ConfigReader
}

// NewSessionReaper creates a reaper that cleans stale sessions from the store.
func NewSessionReaper(store *store.Store) utils.SleeperTask {
	return utils.NewSleeperTask(&sessionReaper{
		store:  store,
		config: store.Config,
	})
}

func (sr *sessionReaper) Work() {
	recordCreationStaleThreshold := sr.config.ReaperExpiration().Before(
		sr.config.SessionTimeout().Before(time.Now()))
	err := sr.store.DeleteStaleSessions(recordCreationStaleThreshold)
	if err != nil {
		logger.Error("unable to reap stale sessions: ", err)
	}
}
