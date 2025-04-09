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

	"github.com/smartcontractkit/chainlink/v2/core/services/llo/types"
)

type LogPoller interface {
	UnregisterFilter(ctx context.Context, filterName string) error
}

func Cleanup(ctx context.Context, lp LogPoller, addr common.Address, donID uint32, ds sqlutil.DataSource, chainSelector uint64) error {
	if (addr != common.Address{} && donID > 0) {
		if err := lp.UnregisterFilter(ctx, types.ChannelDefinitionCacheFilterName(addr, donID)); err != nil {
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
			if n, err := t.reap(overtimeCtx, TransmissionReaperBatchSize, "stale"); err != nil {
				t.lggr.Errorw("Failed to reap stale transmissions on exit", "err", err)
			} else if n > 0 {
				t.lggr.Infow("Reaped stale transmissions on exit", "nDeleted", n)
			}
			cancel()
			return
		case <-ticker.C:
			// TODO: Should also reap other LLO garbage that can be left
			// behind e.g. channel definitions etc
			t.reapAndLog(ctx, TransmissionReaperBatchSize, "stale")
			t.reapAndLog(ctx, TransmissionReaperBatchSize, "orphaned")
		}
	}
}

func (t *transmissionReaper) reapAndLog(ctx context.Context, batchSize int, reapType string) {
	n, err := t.reap(ctx, batchSize, reapType)
	if err != nil {
		t.lggr.Errorw("Failed to reap", "type", reapType, "err", err)
		return
	}
	if n > 0 {
		t.lggr.Infow("Reaped transmissions", "type", reapType, "nDeleted", n)
	}
}

func (t *transmissionReaper) reap(ctx context.Context, batchSize int, reapType string) (rowsDeleted int64, err error) {
	for {
		var res sql.Result
		switch reapType {
		case "stale":
			res, err = t.reapStale(ctx, batchSize)
		case "orphaned":
			res, err = t.reapOrphaned(ctx, batchSize)
		default:
			return 0, fmt.Errorf("transmissionReaper: unknown reap type: %s", reapType)
		}

		if err != nil {
			return rowsDeleted, fmt.Errorf("transmissionReaper: failed to delete %s transmissions: %w", reapType, err)
		}

		var rowsAffected int64
		rowsAffected, err = res.RowsAffected()
		if err != nil {
			return rowsDeleted, fmt.Errorf("transmissionReaper: failed to get %s rows affected: %w", reapType, err)
		}
		if rowsAffected == 0 {
			break
		}
		rowsDeleted += rowsAffected
	}
	return rowsDeleted, nil
}

func (t *transmissionReaper) reapStale(ctx context.Context, batchSize int) (sql.Result, error) {
	return t.ds.ExecContext(ctx, `
DELETE FROM llo_mercury_transmit_queue AS q
USING (
    SELECT transmission_hash 
    FROM llo_mercury_transmit_queue
    WHERE inserted_at < NOW() - ($1 * INTERVAL '1 MICROSECOND')
    ORDER BY inserted_at ASC
    LIMIT $2
) AS to_delete
WHERE q.transmission_hash = to_delete.transmission_hash;
`, t.maxAge.Microseconds(), batchSize)
}

func (t *transmissionReaper) reapOrphaned(ctx context.Context, batchSize int) (sql.Result, error) {
	return t.ds.ExecContext(ctx, `
WITH activeDonIds AS (
    SELECT DISTINCT cast(relay_config->>'lloDonID' as bigint) as don_id
	FROM ocr2_oracle_specs
	WHERE 
		relay_config->>'lloDonID' IS NOT NULL
		AND relay_config->>'lloDonID' <> ''
)
DELETE FROM llo_mercury_transmit_queue as q
USING (
    SELECT transmission_hash 
    FROM llo_mercury_transmit_queue
    WHERE don_id NOT IN (SELECT don_id FROM activeDonIds)
    ORDER BY inserted_at ASC
    LIMIT $1
) AS to_delete
WHERE q.transmission_hash = to_delete.transmission_hash;
`, batchSize)
}
