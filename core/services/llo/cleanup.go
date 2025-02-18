package llo

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
)

func Cleanup(ctx context.Context, lp LogPoller, addr common.Address, donID uint32, ds sqlutil.DataSource, chainSelector uint64) error {
	if (addr != common.Address{} && donID > 0) {
		if err := lp.UnregisterFilter(ctx, filterName(addr, donID)); err != nil {
			return fmt.Errorf("failed to unregister filter: %w", err)
		}
		orm := NewChainScopedORM(ds, chainSelector)
		if err := orm.CleanupChannelDefinitions(ctx, addr, donID); err != nil {
			return fmt.Errorf("failed to cleanup channel definitions: %w", err)
		}
	}
	// Don't bother deleting transmission records since it can be really slow
	// to do that if you have a job that's been erroring for a long time. Let
	// the reaper handle it async instead.
	return nil
}

const (
	// TransmissionReaperBatchSize is the number of transmissions to delete in a
	// single batch.
	TransmissionReaperBatchSize = 10_000
	// OvertimeDeleteTimeout is the maximum time we will spend trying to reap
	// after exit signal before giving up and logging an error.
	OvertimeDeleteTimeout = 2 * time.Second
)

type transmissionReaper struct {
	services.Service
	eng      *services.Engine
	ds       sqlutil.DataSource
	lggr     logger.Logger
	reapFreq time.Duration
	maxAge   time.Duration
}

// NewTransmissionReaper returns a new transmission reaper service
//
// In theory, if everything is working properly, there will never be stale
// transmissions. In practice there can be bugs, jobs that get deleted without
// proper cleanup etc. This acts as a sanity check to evict obviously stale
// entries from the llo_mercury_transmit_queue table.
func NewTransmissionReaper(ds sqlutil.DataSource, lggr logger.Logger, freq, maxAge time.Duration) services.Service {
	t := &transmissionReaper{ds: ds, lggr: lggr, reapFreq: freq, maxAge: maxAge}
	t.Service, t.eng = services.Config{
		Name:  "LLOTransmissionReaper",
		Start: t.start,
	}.NewServiceEngine(lggr)
	return t
}

func (t *transmissionReaper) start(context.Context) error {
	if t.reapFreq == 0 || t.maxAge == 0 {
		t.eng.Debugw("Transmission reaper disabled", "reapFreq", t.reapFreq, "maxAge", t.maxAge)
		return nil
	}
	t.eng.Go(t.runLoop)
	return nil
}

func (t *transmissionReaper) runLoop(ctx context.Context) {
	t.eng.Debugw("Transmission reaper running", "reapFreq", t.reapFreq, "maxAge", t.maxAge)
	ticker := services.TickerConfig{
		// Don't reap right away, wait some time for the application to settle
		// down first
		Initial:   services.DefaultJitter.Apply(t.reapFreq),
		JitterPct: services.DefaultJitter,
	}.NewTicker(t.reapFreq)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			// make a final effort to clear the database that goes into
			// overtime
			overtimeCtx, cancel := context.WithTimeout(context.Background(), OvertimeDeleteTimeout)
			if n, err := t.reapStale(overtimeCtx, TransmissionReaperBatchSize); err != nil {
				t.lggr.Errorw("Failed to reap stale transmissions on exit", "err", err)
			} else if n > 0 {
				t.lggr.Infow("Reaped stale transmissions on exit", "nDeleted", n)
			}
			cancel()
			return
		case <-ticker.C:
			// TODO: Could also automatically reap orphaned transmissions
			// that don't have a job with a matching DON ID (from job
			// deletion)
			//
			// https://smartcontract-it.atlassian.net/browse/MERC-6807
			// TODO: Should also reap other LLO garbage that can be left
			// behind e.g. channel definitions etc
			n, err := t.reapStale(ctx, TransmissionReaperBatchSize)
			if err != nil {
				t.lggr.Errorw("Failed to reap", "err", err)
				continue
			}
			if n > 0 {
				t.lggr.Infow("Reaped stale transmissions", "nDeleted", n)
			}
		}
	}
}

func (t *transmissionReaper) reapStale(ctx context.Context, batchSize int) (rowsDeleted int64, err error) {
	for {
		var res sql.Result
		res, err = t.ds.ExecContext(ctx, `
DELETE FROM llo_mercury_transmit_queue AS q
USING (
    SELECT transmission_hash 
    FROM llo_mercury_transmit_queue
    WHERE inserted_at < $1
    ORDER BY inserted_at ASC
    LIMIT $2
) AS to_delete
WHERE q.transmission_hash = to_delete.transmission_hash;
		`, time.Now().Add(-t.maxAge), batchSize)
		if err != nil {
			return rowsDeleted, fmt.Errorf("transmissionReaper: failed to delete stale transmissions: %w", err)
		}
		var rowsAffected int64
		rowsAffected, err = res.RowsAffected()
		if err != nil {
			return rowsDeleted, fmt.Errorf("transmissionReaper: failed to get rows affected: %w", err)
		}
		if rowsAffected == 0 {
			break
		}
		rowsDeleted += rowsAffected
	}
	return rowsDeleted, nil
}
