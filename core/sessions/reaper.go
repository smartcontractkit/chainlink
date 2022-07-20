package sessions

import (
	"database/sql"
	"time"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"
)

type sessionReaper struct {
	db     *sql.DB
	config SessionReaperConfig
	lggr   logger.Logger
}

type SessionReaperConfig interface {
	SessionTimeout() models.Duration
	ReaperExpiration() models.Duration
}

// NewSessionReaper creates a reaper that cleans stale sessions from the store.
func NewSessionReaper(db *sql.DB, config SessionReaperConfig, lggr logger.Logger) utils.SleeperTask {
	return utils.NewSleeperTask(&sessionReaper{
		db,
		config,
		lggr.Named("SessionReaper"),
	})
}

func (sr *sessionReaper) Name() string {
	return "SessionReaper"
}

func (sr *sessionReaper) Work() {
	recordCreationStaleThreshold := sr.config.ReaperExpiration().Before(
		sr.config.SessionTimeout().Before(time.Now()))
	err := sr.deleteStaleSessions(recordCreationStaleThreshold)
	if err != nil {
		sr.lggr.Error("unable to reap stale sessions: ", err)
	}
}

// DeleteStaleSessions deletes all sessions before the passed time.
func (sr *sessionReaper) deleteStaleSessions(before time.Time) error {
	_, err := sr.db.Exec("DELETE FROM sessions WHERE last_used < $1", before)
	return err
}
