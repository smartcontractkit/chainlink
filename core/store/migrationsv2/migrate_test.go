package migrationsv2_test

import (
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/core/store/migrationsv2"
	"gopkg.in/guregu/null.v4"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

func TestMigrate_Initial(t *testing.T) {
	_, orm, cleanup := cltest.BootstrapThrowawayORM(t, "migrationsv2", false)
	defer cleanup()

	err := migrationsv2.MigrateUp(orm.DB, "1611847145")
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
		r := orm.DB.Exec("SELECT * from information_schema.tables where table_name = ?", table)
		require.NoError(t, r.Error)
		assert.True(t, r.RowsAffected > 0, "table %v not found", table)
	}
	err = migrationsv2.Rollback(orm.DB, migrationsv2.Migrations[0])
	require.NoError(t, err)

	for _, table := range tables {
		r := orm.DB.Exec("SELECT * from information_schema.tables where table_name = ?", table)
		require.NoError(t, r.Error)
		assert.False(t, r.RowsAffected > 0, "table %v found", table)
	}
}

func TestMigrate_BridgeFK(t *testing.T) {
	_, orm, cleanup := cltest.BootstrapThrowawayORM(t, "migrationsv2_bridgefk", false)
	defer cleanup()

	require.NoError(t, migrationsv2.MigrateUp(orm.DB, "0009_add_min_payment_to_flux_monitor_spec"))
	_, bt := cltest.NewBridgeType(t)
	require.NoError(t, orm.DB.Create(bt).Error)
	ps := &pipeline.Spec{}
	require.NoError(t, orm.DB.Create(&ps).Error)

	// V1 pipeline.TaskSpec
	type PipelineTaskSpec struct {
		ID             int32                     `json:"-" gorm:"primary_key"`
		DotID          string                    `json:"dotId"`
		PipelineSpecID int32                     `json:"-"`
		PipelineSpec   pipeline.Spec             `json:"-"`
		Type           pipeline.TaskType         `json:"-"`
		JSON           pipeline.JSONSerializable `json:"-" gorm:"type:jsonb"`
		Index          int32                     `json:"-"`
		SuccessorID    null.Int                  `json:"-"`
		CreatedAt      time.Time                 `json:"-"`
	}
	bts := &PipelineTaskSpec{
		PipelineSpecID: ps.ID,
		Type:           pipeline.TaskTypeBridge,
		JSON: pipeline.JSONSerializable{
			Val: pipeline.BridgeTask{
				Name: string(bt.Name),
			}, Null: false},
	}
	require.NoError(t, orm.DB.Create(&bts).Error)
	require.NoError(t, orm.DB.Create(&PipelineTaskSpec{
		PipelineSpecID: ps.ID,
		Type:           pipeline.TaskTypeAny,
		SuccessorID:    null.NewInt(int64(bts.ID), true),
	}).Error)

	// Migrating up should populate the bridge field
	require.NoError(t, migrationsv2.MigrateUp(orm.DB, "0010_bridge_fk"))

	// V2 pipeline TaskSpec
	var p pipeline.TaskSpec
	require.NoError(t, orm.DB.Find(&p, "id = ?", bts.ID).Error)

	assert.Equal(t, *p.BridgeName, string(bt.Name))
}
