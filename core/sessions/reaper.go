package sessions

import (
	"database/sql"
	"time"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/store/models"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

type sessionReaper struct {
	db     *sql.DB
	config SessionReaperConfig
	lggr   logger.Logger

	// Receive from this for testing via sr.RunSignal()
	// to be notified after each reaper run.
	runSignal chan struct{}
}

type SessionReaperConfig interface {
	SessionTimeout() models.Duration
	SessionReaperExpiration() models.Duration
}

// NewSessionReaper creates a reaper that cleans stale sessions from the store.
func NewSessionReaper(db *sql.DB, config SessionReaperConfig, lggr logger.Logger) utils.SleeperTask {
	return utils.NewSleeperTask(NewSessionReaperWorker(db, config, lggr))
}

func NewSessionReaperWorker(db *sql.DB, config SessionReaperConfig, lggr logger.Logger) *sessionReaper {
	return &sessionReaper{
		db,
		config,
		lggr.Named("SessionReaper"),

		// For testing only.
		make(chan struct{}, 10),
	}
}

func (sr *sessionReaper) Name() string {
	return "SessionReaper"
}

func (sr *sessionReaper) Work() {
	recordCreationStaleThreshold := sr.config.SessionReaperExpiration().Before(
		sr.config.SessionTimeout().Before(time.Now()))
	err := sr.deleteStaleSessions(recordCreationStaleThreshold)
	if err != nil {
		sr.lggr.Error("unable to reap stale sessions: ", err)
	}

	select {
	case sr.runSignal <- struct{}{}:
	default:
	}
}

// DeleteStaleSessions deletes all sessions before the passed time.
func (sr *sessionReaper) deleteStaleSessions(before time.Time) error {
	_, err := sr.db.Exec("DELETE FROM sessions WHERE last_used < $1", before)
	return err
}
