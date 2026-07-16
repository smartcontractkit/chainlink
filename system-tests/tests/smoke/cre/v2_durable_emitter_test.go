package cre

import (
	"context"
	"database/sql"
	"fmt"
	"slices"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"

	crontypes "github.com/smartcontractkit/chainlink/core/scripts/cre/environment/examples/workflows/cron/types"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	t_helpers "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers"
	ttypes "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers/configuration"
)

// connectWorkflowDONDB connects to the Postgres database of the first workflow
// DON NodeSet. This is where cre.chip_durable_events lives.
func connectWorkflowDONDB(t *testing.T, nodeSets []*cre.NodeSet) *sql.DB {
	t.Helper()

	var port int
	var label string
	for _, ns := range nodeSets {
		if slices.Contains(ns.DONTypes, cre.WorkflowDON) {
			port = ns.DbInput.Port
			label = ns.Name
			break
		}
	}
	require.NotZerof(t, port, "no workflow DON NodeSet found")

	dsn := fmt.Sprintf(
		"host=localhost port=%d user=chainlink password=thispasswordislongenough dbname=db_0 sslmode=disable",
		port,
	)
	db, err := sql.Open("postgres", dsn)
	require.NoError(t, err)
	require.NoError(t, db.PingContext(t.Context()))
	t.Logf("connected to %s workflow DON DB (port %d) for durable emitter tracking", label, port)
	t.Cleanup(func() { _ = db.Close() })
	return db
}

type durableEventStats struct {
	inserts int64
	deletes int64
}

// snapshotDurableEventStats returns cumulative insert/delete counts for
// chip_durable_events from pg_stat_user_tables.
func snapshotDurableEventStats(ctx context.Context, db *sql.DB) (durableEventStats, error) {
	var s durableEventStats
	err := db.QueryRowContext(ctx,
		`SELECT COALESCE(n_tup_ins,0), COALESCE(n_tup_del,0)
		   FROM pg_stat_user_tables
		  WHERE relname = 'chip_durable_events'`,
	).Scan(&s.inserts, &s.deletes)
	if err == sql.ErrNoRows {
		return durableEventStats{}, nil
	}
	return s, err
}

// countPendingDurableEvents returns rows still awaiting delivery (delivered_at IS NULL).
func countPendingDurableEvents(ctx context.Context, db *sql.DB) (int64, error) {
	var count int64
	err := db.QueryRowContext(ctx,
		`SELECT count(*) FROM cre.chip_durable_events WHERE delivered_at IS NULL`,
	).Scan(&count)
	return count, err
}

// resetDurableEventQueue removes all pending durable events so queue depth and pending
// counts don't carry over from other tests or earlier suite steps on the same DB.
func resetDurableEventQueue(ctx context.Context, t *testing.T, db *sql.DB) {
	t.Helper()
	res, err := db.ExecContext(ctx, `DELETE FROM cre.chip_durable_events`)
	require.NoError(t, err)
	n, err := res.RowsAffected()
	require.NoError(t, err)
	t.Logf("cleared cre.chip_durable_events (%d rows removed before test)", n)
}

// ExecuteDurableEmitterTest verifies the DurableEmitter is active and
// functioning by deploying a cron workflow that emits events, then checking
// that chip_durable_events sees sustained insert+delete activity over time.
func ExecuteDurableEmitterTest(t *testing.T, testEnv *ttypes.TestEnvironment) {
	lggr := framework.L
	workflowFileLocation := "../../../../core/scripts/cre/environment/examples/workflows/cron/main.go"

	db := connectWorkflowDONDB(t, testEnv.Config.NodeSets)

	_, err := countPendingDurableEvents(t.Context(), db)
	require.NoError(t, err)

	resetDurableEventQueue(t.Context(), t, db)

	baseline, err := snapshotDurableEventStats(t.Context(), db)
	require.NoError(t, err)
	t.Logf("baseline chip_durable_events stats: inserts=%d deletes=%d", baseline.inserts, baseline.deletes)

	// Deploy a cron workflow that fires every 5 seconds.
	lggr.Info().Msg("Deploying cron workflow for durable emitter test...")
	workflowConfig := crontypes.WorkflowConfig{
		Schedule: "*/5 * * * * *",
	}
	_ = t_helpers.CompileAndDeployWorkflow(t, testEnv, lggr, "durable-emitter-test", &workflowConfig, workflowFileLocation)

	const minExpectedEvents int64 = 4
	lggr.Info().Msg("Waiting for sustained durable event activity...")

	require.Eventually(t, func() bool {
		stats, statsErr := snapshotDurableEventStats(t.Context(), db)
		if statsErr != nil {
			t.Logf("failed to snapshot stats: %v", statsErr)
			return false
		}

		newInserts := stats.inserts - baseline.inserts
		newDeletes := stats.deletes - baseline.deletes

		pending, _ := countPendingDurableEvents(t.Context(), db)
		t.Logf("chip_durable_events: +%d inserts, +%d deletes, %d pending", newInserts, newDeletes, pending)

		return newInserts >= minExpectedEvents && newDeletes >= minExpectedEvents
	}, 2*time.Minute, 5*time.Second, "expected at least %d insert+delete events", minExpectedEvents)
}
