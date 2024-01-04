package localauth

import (
	"database/sql"
	"time"

	commonconfig "github.com/smartcontractkit/chainlink-common/pkg/config"
	"github.com/smartcontractkit/chainlink-common/pkg/utils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

type sessionReaper struct {
	db     *sql.DB
	config SessionReaperConfig
	lggr   logger.Logger
}

type SessionReaperConfig interface {
	SessionTimeout() commonconfig.Duration
	SessionReaperExpiration() commonconfig.Duration
}

// NewSessionReaper creates a reaper that cleans stale sessions from the store.
func NewSessionReaper(db *sql.DB, config SessionReaperConfig, lggr logger.Logger) *utils.SleeperTask {
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
	recordCreationStaleThreshold := sr.config.SessionReaperExpiration().Before(
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
