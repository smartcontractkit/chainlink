package migrationsv2_test

import (
	"testing"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/store/migrationsv2"

	"github.com/stretchr/testify/assert"

	_ "github.com/jinzhu/gorm/dialects/postgres" // http://doc.gorm.io/database.html#connecting-to-a-database
	"github.com/stretchr/testify/require"
)

func TestMigrate_Migrations(t *testing.T) {
	_, orm, cleanup := cltest.BootstrapThrowawayORM(t, "migrationsv2", false)
	defer cleanup()

	err := migrationsv2.MigrateUp(orm.DB, "")
	require.NoError(t, err)
	tables := []string{
		"bridge_types",
		"configurations",
		"direct_request_specs",
		"encrypted_ocr_key_bundles",
		"encrypted_p2p_keys",
		"encrypted_vrf_keys",
		"encumbrances",
		"eth_receipts",
		"eth_task_run_txes",
		"eth_tx_attempts",
		"eth_txes",
		"external_initiators",
		"flux_monitor_round_stats",
		"flux_monitor_specs",
		"heads",
		"initiators",
		"job_runs",
		"job_spec_errors",
		"job_spec_errors_v2",
		"job_specs",
		"jobs",
		"keys",
		"log_consumptions",
		"offchainreporting_contract_configs",
		"offchainreporting_oracle_specs",
		"offchainreporting_pending_transmissions",
		"offchainreporting_persistent_states",
		"p2p_peers",
		"pipeline_runs",
		"pipeline_specs",
		"pipeline_task_runs",
		"pipeline_task_specs",
		"run_requests",
		"run_results",
		"service_agreements",
		"sessions",
		"sync_events",
		"task_runs",
		"task_specs",
		"eth_tx_attempts",
		"eth_txes",
		"users",
	}
	for _, table := range tables {
		assert.True(t, orm.DB.HasTable(table), "table %v not found", table)
	}
	err = migrationsv2.MigrateDown(orm.DB)
	require.NoError(t, err)

	for _, table := range tables {
		assert.False(t, orm.DB.HasTable(table), "table %v found", table)
	}
}
