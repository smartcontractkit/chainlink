package cre

import (
	"context"
	"database/sql"
	"fmt"
	"slices"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"

	crontypes "github.com/smartcontractkit/chainlink/core/scripts/cre/environment/examples/workflows/v2/cron/types"

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

// countPendingDurableEvents returns the current number of rows in
// cre.chip_durable_events (events that haven't been delivered yet).
func countPendingDurableEvents(ctx context.Context, db *sql.DB) (int64, error) {
	var count int64
	err := db.QueryRowContext(ctx,
		`SELECT count(*) FROM cre.chip_durable_events`,
	).Scan(&count)
	return count, err
}

// ExecuteDurableEmitterTest verifies the DurableEmitter is active and
// functioning by deploying a cron workflow that emits events, then checking
// that chip_durable_events sees sustained insert+delete activity over time.
func ExecuteDurableEmitterTest(t *testing.T, testEnv *ttypes.TestEnvironment) {
	lggr := framework.L
	workflowFileLocation := "../../../../core/scripts/cre/environment/examples/workflows/v2/cron/main.go"

	db := connectWorkflowDONDB(t, testEnv.Config.NodeSets)

	_, err := countPendingDurableEvents(t.Context(), db)
	require.NoError(t, err, "cre.chip_durable_events table should exist — check migration 0295")

	baseline, err := snapshotDurableEventStats(t.Context(), db)
	require.NoError(t, err)
	t.Logf("baseline chip_durable_events stats: inserts=%d deletes=%d", baseline.inserts, baseline.deletes)

	// Deploy a cron workflow that fires every 5 seconds.
	lggr.Info().Msg("Deploying cron workflow for durable emitter test...")
	workflowConfig := crontypes.WorkflowConfig{
		Schedule: "*/5 * * * * *",
	}
	_ = t_helpers.CompileAndDeployWorkflow(t, testEnv, lggr, "durable-emitter-test", &workflowConfig, workflowFileLocation)

	// Wait for a meaningful volume of events to flow through the pipeline.
	// Each cron execution emits ~3-5 beholder events across the DON.
	// At every-5s with 4 nodes, expect ~50+ events per minute.
	const minExpectedEvents int64 = 30

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
	}, 4*time.Minute, 10*time.Second, "expected at least %d insert+delete events", minExpectedEvents)

	pending, err := countPendingDurableEvents(t.Context(), db)
	require.NoError(t, err)
	t.Logf("pending durable events at end of test: %d", pending)
	assert.LessOrEqual(t, pending, int64(10),
		"durable event queue should be near-empty when chip ingress is healthy")

	final, err := snapshotDurableEventStats(t.Context(), db)
	require.NoError(t, err)
	t.Logf("final chip_durable_events stats: inserts=%d (+%d) deletes=%d (+%d)",
		final.inserts, final.inserts-baseline.inserts,
		final.deletes, final.deletes-baseline.deletes)

	lggr.Info().Msg("Durable emitter test completed successfully")
}

// ExecuteDurableEmitterLoadTest deploys multiple high-frequency cron workflows
// to stress the durable emitter pipeline. This measures the maximum sustained
// throughput of the persist → publish → delete cycle against real Postgres.
func ExecuteDurableEmitterLoadTest(t *testing.T, testEnv *ttypes.TestEnvironment) {
	lggr := framework.L
	workflowFileLocation := "../../../../core/scripts/cre/environment/examples/workflows/v2/cron/main.go"

	db := connectWorkflowDONDB(t, testEnv.Config.NodeSets)

	_, err := countPendingDurableEvents(t.Context(), db)
	require.NoError(t, err, "cre.chip_durable_events table should exist")

	baseline, err := snapshotDurableEventStats(t.Context(), db)
	require.NoError(t, err)
	t.Logf("baseline: inserts=%d deletes=%d", baseline.inserts, baseline.deletes)

	// Deploy multiple cron workflows, each firing every second.
	// Each execution emits ~3-5 events per node. With 4 nodes and N workflows,
	// we expect roughly N * 4 * 4 = 16N events/sec across the DON.
	const numWorkflows = 5
	cronConfig := crontypes.WorkflowConfig{
		Schedule: "*/1 * * * * *", // every second
	}

	lggr.Info().Msgf("Deploying %d high-frequency cron workflows...", numWorkflows)
	for i := 0; i < numWorkflows; i++ {
		name := fmt.Sprintf("durable-load-%d", i)
		_ = t_helpers.CompileAndDeployWorkflow(t, testEnv, lggr, name, &cronConfig, workflowFileLocation)
	}

	// Let the load run for a fixed observation window.
	const observationPeriod = 3 * time.Minute
	lggr.Info().Msgf("Load running for %s — monitoring durable event stats...", observationPeriod)

	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()
	deadline := time.After(observationPeriod)

	var maxPending int64
	var lastStats durableEventStats

	for {
		select {
		case <-deadline:
			goto done
		case <-ticker.C:
			stats, statsErr := snapshotDurableEventStats(t.Context(), db)
			if statsErr != nil {
				t.Logf("stats error: %v", statsErr)
				continue
			}
			pending, _ := countPendingDurableEvents(t.Context(), db)

			newInserts := stats.inserts - baseline.inserts
			newDeletes := stats.deletes - baseline.deletes

			if pending > maxPending {
				maxPending = pending
			}

			// Calculate rates over the last interval.
			var insertRate, deleteRate float64
			if lastStats.inserts > 0 {
				insertRate = float64(stats.inserts-lastStats.inserts) / 15.0
				deleteRate = float64(stats.deletes-lastStats.deletes) / 15.0
			}
			lastStats = stats

			t.Logf("durable events: +%d ins, +%d del | pending: %d (max %d) | rate: %.1f ins/s, %.1f del/s",
				newInserts, newDeletes, pending, maxPending, insertRate, deleteRate)
		}
	}

done:
	final, err := snapshotDurableEventStats(t.Context(), db)
	require.NoError(t, err)
	pending, err := countPendingDurableEvents(t.Context(), db)
	require.NoError(t, err)

	totalInserts := final.inserts - baseline.inserts
	totalDeletes := final.deletes - baseline.deletes
	avgInsertRate := float64(totalInserts) / observationPeriod.Seconds()
	avgDeleteRate := float64(totalDeletes) / observationPeriod.Seconds()

	t.Logf("╔════════════════════════════════════════════════╗")
	t.Logf("║        DURABLE EMITTER LOAD TEST RESULTS      ║")
	t.Logf("╠════════════════════════════════════════════════╣")
	t.Logf("║ Workflows deployed:   %d                       ║", numWorkflows)
	t.Logf("║ Observation period:   %s                   ║", observationPeriod)
	t.Logf("║ Total inserts:        %-6d                   ║", totalInserts)
	t.Logf("║ Total deletes:        %-6d                   ║", totalDeletes)
	t.Logf("║ Avg insert rate:      %-6.1f events/sec        ║", avgInsertRate)
	t.Logf("║ Avg delete rate:      %-6.1f events/sec        ║", avgDeleteRate)
	t.Logf("║ Max queue depth:      %-6d                   ║", maxPending)
	t.Logf("║ Final pending:        %-6d                   ║", pending)
	t.Logf("╚════════════════════════════════════════════════╝")

	// Sanity checks.
	assert.Greater(t, totalInserts, int64(100),
		"expected significant event volume from %d workflows", numWorkflows)
	assert.Greater(t, totalDeletes, int64(0),
		"deletes must occur — chip delivery is required")
	assert.LessOrEqual(t, pending, int64(50),
		"queue should not grow unboundedly with healthy chip ingress")

	lggr.Info().Msg("Durable emitter load test completed")
}
