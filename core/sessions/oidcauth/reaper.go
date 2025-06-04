package oidcauth

import (
	"context"
	"time"

	commonconfig "github.com/smartcontractkit/chainlink-common/pkg/config"
	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	"github.com/smartcontractkit/chainlink-common/pkg/utils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

type sessionReaper struct {
	ds     sqlutil.DataSource
	config SessionReaperConfig
	lggr   logger.Logger
}

type SessionReaperConfig interface {
	SessionTimeout() commonconfig.Duration
	SessionReaperExpiration() commonconfig.Duration
}

// NewSessionReaper creates a reaper that cleans stale sessions from the store.
func NewSessionReaper(ds sqlutil.DataSource, config SessionReaperConfig, lggr logger.Logger) *utils.SleeperTask {
	return utils.NewSleeperTaskCtx(&sessionReaper{
		ds,
		config,
		lggr.Named("OIDCSessionReaper"),
	})
}

func (sr *sessionReaper) Name() string {
	return "OIDCSessionReaper"
}

func (sr *sessionReaper) Work(ctx context.Context) {
	recordCreationStaleThreshold := sr.config.SessionReaperExpiration().Before(
		sr.config.SessionTimeout().Before(time.Now()))
	err := sr.deleteStaleSessions(ctx, recordCreationStaleThreshold)
	if err != nil {
		sr.lggr.Error("unable to reap stale sessions: ", err)
	}
}

// DeleteStaleSessions deletes all sessions before the passed time.
func (sr *sessionReaper) deleteStaleSessions(ctx context.Context, before time.Time) error {
	_, err := sr.ds.ExecContext(ctx, "DELETE FROM oidc_sessions WHERE created_at < $1", before)
	return err
}
