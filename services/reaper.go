package services

import (
	"time"

	"github.com/asdine/storm"
	"github.com/smartcontractkit/chainlink/logger"
	"github.com/smartcontractkit/chainlink/store"
	"github.com/smartcontractkit/chainlink/store/models"
)

type storeReaper struct {
	store  *store.Store
	config store.Config
}

// NewStoreReaper creates a reaper that cleans stale objects from the store.
func NewStoreReaper(store *store.Store) SleeperTask {
	return NewSleeperTask(&storeReaper{
		store:  store,
		config: store.Config,
	})
}

func (sr *storeReaper) Work() {
	var sessions []models.Session
	offset := time.Now().Add(-sr.config.ReaperExpiration()).Add(-sr.config.SessionTimeout())
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
