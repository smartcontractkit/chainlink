package migrations_test

import (
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/core/services/job"

	uuid "github.com/satori/go.uuid"

	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/cltest/heavyweight"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/core/store/migrations"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"gopkg.in/guregu/null.v4"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

func TestMigrate_Initial(t *testing.T) {
	_, orm, cleanup := heavyweight.FullTestORM(t, "migrations", false)
	defer cleanup()

	err := migrations.MigrateUp(orm.DB, "1611847145")
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
	err = migrations.Rollback(orm.DB, migrations.Migrations[0])
	require.NoError(t, err)

	for _, table := range tables {
		r := orm.DB.Exec("SELECT * from information_schema.tables where table_name = ?", table)
		require.NoError(t, r.Error)
		assert.False(t, r.RowsAffected > 0, "table %v found", table)
	}
}

// V2 pipeline TaskSpec
// DEPRECATED
type TaskSpec struct {
	ID             int32                     `json:"-" gorm:"primary_key"`
	DotID          string                    `json:"dotId"`
	PipelineSpecID int32                     `json:"-"`
	PipelineSpec   pipeline.Spec             `json:"-"`
	Type           pipeline.TaskType         `json:"-"`
	JSON           pipeline.JSONSerializable `json:"-" gorm:"type:jsonb"`
	Index          int32                     `json:"-"`
	SuccessorID    null.Int                  `json:"-"`
	CreatedAt      time.Time                 `json:"-"`
	BridgeName     *string                   `json:"-"`
	Bridge         models.BridgeType         `json:"-" gorm:"foreignKey:BridgeName;->"`
}

func (TaskSpec) TableName() string {
	return "pipeline_task_specs"
}

func TestMigrate_BridgeFK(t *testing.T) {
	_, orm, cleanup := heavyweight.FullTestORM(t, "migrations_bridgefk", false)
	defer cleanup()

	require.NoError(t, migrations.MigrateUp(orm.DB, "0009_add_min_payment_to_flux_monitor_spec"))
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
	require.NoError(t, migrations.MigrateUp(orm.DB, "0010_bridge_fk"))

	var p TaskSpec
	require.NoError(t, orm.DB.Find(&p, "id = ?", bts.ID).Error)

	assert.Equal(t, *p.BridgeName, string(bt.Name))

	// Run the down migration
	require.NoError(t, migrations.MigrateDownFrom(orm.DB, "0010_bridge_fk"))
}

func TestMigrate_ChangeJobsToNumeric(t *testing.T) {
	_, orm, cleanup := heavyweight.FullTestORM(t, "migrations_change_jobs_to_numeric", false)
	defer cleanup()

	require.NoError(t, migrations.MigrateUp(orm.DB, "0010_bridge_fk"))

	jobSpec := cltest.NewJob()
	jobSpec.MinPayment = assets.NewLink(100)
	orm.DB.Create(&jobSpec)

	// These structs are copied from the `job` package because otherwise, changes
	// to them in subsequent migrations will cause this test to fail.
	type FluxMonitorSpec struct {
		ID                int32               `toml:"-" gorm:"primary_key"`
		ContractAddress   ethkey.EIP55Address `toml:"contractAddress"`
		Threshold         float32             `toml:"threshold,float"`
		AbsoluteThreshold float32             `toml:"absoluteThreshold,float" gorm:"type:float;not null"`
		PollTimerPeriod   time.Duration       `gorm:"type:jsonb"`
		PollTimerDisabled bool                `gorm:"type:jsonb"`
		IdleTimerPeriod   time.Duration       `gorm:"type:jsonb"`
		IdleTimerDisabled bool                `gorm:"type:jsonb"`
		MinPayment        *assets.Link
		CreatedAt         time.Time `toml:"-"`
		UpdatedAt         time.Time `toml:"-"`
	}

	type Job struct {
		ID                int32     `toml:"-" gorm:"primary_key"`
		ExternalJobID     uuid.UUID `toml:"externalJobID"`
		FluxMonitorSpecID *int32
		FluxMonitorSpec   *FluxMonitorSpec
		PipelineSpecID    int32
		PipelineSpec      *pipeline.Spec
	}

	fmSpec := FluxMonitorSpec{
		MinPayment:        assets.NewLink(100),
		ContractAddress:   cltest.NewEIP55Address(),
		PollTimerDisabled: true,
		IdleTimerDisabled: true,
	}
	orm.DB.Create(&fmSpec)

	require.NoError(t, migrations.MigrateUp(orm.DB, "0012_change_jobs_to_numeric"))

	var js models.JobSpec
	require.NoError(t, orm.DB.Find(&js, "id = ?", jobSpec.ID).Error)
	require.Equal(t, assets.NewLink(100), js.MinPayment)

	var fms FluxMonitorSpec
	require.NoError(t, orm.DB.Find(&fms, "id = ?", fmSpec.ID).Error)
	require.Equal(t, assets.NewLink(100), fms.MinPayment)

	require.NoError(t, migrations.MigrateDownFrom(orm.DB, "0012_change_jobs_to_numeric"))
}

func TestMigrate_PipelineTaskRunDotID(t *testing.T) {
	_, orm, cleanup := heavyweight.FullTestORM(t, "migrations_task_run_dot_id", false)
	defer cleanup()

	require.NoError(t, migrations.MigrateUp(orm.DB, "0015_simplify_log_broadcaster"))
	// Add some task specs
	ps := pipeline.Spec{
		DotDagSource: "blah",
	}
	require.NoError(t, orm.DB.Create(&ps).Error)
	result := TaskSpec{
		DotID:          "__result__",
		PipelineSpecID: ps.ID,
		Type:           "result",
		JSON:           pipeline.JSONSerializable{},
		SuccessorID:    null.Int{},
	}
	require.NoError(t, orm.DB.Create(&result).Error)
	ds := TaskSpec{
		DotID:          "ds1",
		PipelineSpecID: ps.ID,
		Type:           "http",
		JSON:           pipeline.JSONSerializable{},
		SuccessorID:    null.NewInt(int64(result.ID), true),
	}
	require.NoError(t, orm.DB.Create(&ds).Error)
	// Add a pipeline run
	type PipelineRun struct {
		ID             int64                     `json:"-" gorm:"primary_key"`
		PipelineSpecID int32                     `json:"-"`
		PipelineSpec   pipeline.Spec             `json:"pipelineSpec"`
		Meta           pipeline.JSONSerializable `json:"meta"`
		// The errors are only ever strings
		// DB example: [null, null, "my error"]
		Errors pipeline.RunErrors `json:"errors" gorm:"type:jsonb"`
		// The outputs can be anything.
		// DB example: [1234, {"a": 10}, null]
		Outputs          pipeline.JSONSerializable `json:"outputs" gorm:"type:jsonb"`
		CreatedAt        time.Time                 `json:"createdAt"`
		FinishedAt       *time.Time                `json:"finishedAt"`
		PipelineTaskRuns []pipeline.TaskRun        `json:"taskRuns" gorm:"foreignkey:PipelineRunID;->"`
	}
	pr := PipelineRun{
		PipelineSpecID: ps.ID,
		Meta:           pipeline.JSONSerializable{},
		Errors:         pipeline.RunErrors{},
		Outputs:        pipeline.JSONSerializable{Null: true},
	}
	require.NoError(t, orm.DB.Create(&pr).Error)

	// Add some task runs
	type PipelineTaskRun struct {
		ID                 int64                      `json:"-" gorm:"primary_key"`
		Type               pipeline.TaskType          `json:"type"`
		PipelineRun        PipelineRun                `json:"-"`
		PipelineRunID      int64                      `json:"-"`
		Output             *pipeline.JSONSerializable `json:"output" gorm:"type:jsonb"`
		Error              null.String                `json:"error"`
		CreatedAt          time.Time                  `json:"createdAt"`
		FinishedAt         *time.Time                 `json:"finishedAt"`
		Index              int32
		PipelineTaskSpecID int32 `json:"-"`
	}
	tr1 := PipelineTaskRun{
		Type:               pipeline.TaskTypeAny,
		PipelineRunID:      pr.ID,
		PipelineTaskSpecID: result.ID,
		Output:             &pipeline.JSONSerializable{Null: true},
		Error:              null.String{},
	}
	require.NoError(t, orm.DB.Create(&tr1).Error)
	tr2 := PipelineTaskRun{
		Type:               pipeline.TaskTypeHTTP,
		PipelineTaskSpecID: ds.ID,
		PipelineRunID:      pr.ID,
		Output:             &pipeline.JSONSerializable{Null: true},
		Error:              null.String{},
	}
	require.NoError(t, orm.DB.Create(&tr2).Error)

	require.NoError(t, migrations.MigrateUp(orm.DB, "0016_pipeline_task_run_dot_id"))
	type NewPipelineTaskRun struct {
		ID            int64                      `json:"-" gorm:"primary_key"`
		Type          pipeline.TaskType          `json:"type"`
		PipelineRun   PipelineRun                `json:"-"`
		PipelineRunID int64                      `json:"-"`
		Output        *pipeline.JSONSerializable `json:"output" gorm:"type:jsonb"`
		Error         null.String                `json:"error"`
		CreatedAt     time.Time                  `json:"createdAt"`
		FinishedAt    *time.Time                 `json:"finishedAt"`
		Index         int32
		DotID         string `json:"dotId"`
	}
	var ptrs []NewPipelineTaskRun
	require.NoError(t, orm.DB.Table("pipeline_task_runs").Find(&ptrs).Error)
	assert.Equal(t, "__result__", ptrs[0].DotID)
	assert.Equal(t, "ds1", ptrs[1].DotID)

	require.NoError(t, migrations.MigrateDownFrom(orm.DB, "0016_pipeline_task_run_dot_id"))

}

func TestMigrate_RemoveResultTask(t *testing.T) {
	_, orm, cleanup := heavyweight.FullTestORM(t, "migrations_result_task", false)
	defer cleanup()

	require.NoError(t, migrations.MigrateUp(orm.DB, "0019_last_run_height_column_to_keeper_table"))
	// Add some task specs
	ps := pipeline.Spec{
		DotDagSource: "blah",
	}
	require.NoError(t, orm.DB.Create(&ps).Error)
	// Add a pipeline run
	type PipelineRun struct {
		ID             int64                     `json:"-" gorm:"primary_key"`
		PipelineSpecID int32                     `json:"-"`
		PipelineSpec   pipeline.Spec             `json:"pipelineSpec"`
		Meta           pipeline.JSONSerializable `json:"meta"`
		// The errors are only ever strings
		// DB example: [null, null, "my error"]
		Errors pipeline.RunErrors `json:"errors" gorm:"type:jsonb"`
		// The outputs can be anything.
		// DB example: [1234, {"a": 10}, null]
		Outputs          pipeline.JSONSerializable `json:"outputs" gorm:"type:jsonb"`
		CreatedAt        time.Time                 `json:"createdAt"`
		FinishedAt       *time.Time                `json:"finishedAt"`
		PipelineTaskRuns []pipeline.TaskRun        `json:"taskRuns" gorm:"foreignkey:PipelineRunID;->"`
	}
	pr := PipelineRun{
		PipelineSpecID: ps.ID,
		Meta:           pipeline.JSONSerializable{},
		Errors:         pipeline.RunErrors{},
		Outputs:        pipeline.JSONSerializable{Null: true},
	}
	type PipelineTaskRun struct {
		ID            int64                      `json:"-" gorm:"primary_key"`
		Type          pipeline.TaskType          `json:"type"`
		PipelineRun   PipelineRun                `json:"-"`
		PipelineRunID int64                      `json:"-"`
		Output        *pipeline.JSONSerializable `json:"output" gorm:"type:jsonb"`
		Error         null.String                `json:"error"`
		CreatedAt     time.Time                  `json:"createdAt"`
		FinishedAt    *time.Time                 `json:"finishedAt"`
		DotID         string                     `json:"dotId"`
		Index         int32
	}
	require.NoError(t, orm.DB.Create(&pr).Error)
	tr1 := PipelineTaskRun{
		Type:          pipeline.TaskTypeAny,
		DotID:         "any",
		PipelineRunID: pr.ID,
		Output:        &pipeline.JSONSerializable{Null: true},
		Error:         null.String{},
	}
	require.NoError(t, orm.DB.Create(&tr1).Error)
	f := time.Now()
	tr2 := PipelineTaskRun{
		Type:          "result",
		DotID:         "result",
		PipelineRunID: pr.ID,
		Output:        &pipeline.JSONSerializable{Val: "10"},
		Error:         null.StringFrom("[null]"),
		FinishedAt:    &f,
	}
	require.NoError(t, orm.DB.Create(&tr2).Error)

	require.NoError(t, migrations.MigrateUp(orm.DB, "0020_remove_result_task"))
	var ptrs []PipelineTaskRun
	require.NoError(t, orm.DB.Find(&ptrs).Error)
	assert.Equal(t, 1, len(ptrs))

	require.NoError(t, migrations.MigrateDownFrom(orm.DB, "0020_remove_result_task"))
}

func TestMigrate_LogConfigTables(t *testing.T) {
	_, orm, cleanup := heavyweight.FullTestORM(t, "migrations_create_log_config_tables", false)
	defer cleanup()

	require.NoError(t, migrations.MigrateUp(orm.DB, "0025_create_log_config_table"))

	sql, err := orm.DB.DB()
	require.NoError(t, err)
	res, err := sql.Query("select 1 from pg_type where typname = 'log_level'")
	require.NoError(t, err)
	assert.True(t, res.Next())

	lgCfg := logger.LogConfig{
		ServiceName: "head_tracker",
		LogLevel:    "warn",
		CreatedAt:   time.Time{},
		UpdatedAt:   time.Time{},
	}
	require.NoError(t, orm.DB.Create(&lgCfg).Error)
	require.NoError(t, orm.DB.Find(&lgCfg).Error)
	require.NoError(t, migrations.MigrateDownFrom(orm.DB, "0025_create_log_config_table"))

	err = orm.DB.Create(&lgCfg).Error
	require.Error(t, err)
	assert.Contains(t, err.Error(), "relation \"log_configs\" does not exist")

	res, err = sql.Query("select 1 from pg_type where typname = 'log_level'")
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.False(t, res.Next())
}

func TestMigrate_CreateCronTables(t *testing.T) {
	_, orm, cleanup := heavyweight.FullTestORM(t, "migrations_create_cron_tables", false)
	defer cleanup()

	require.NoError(t, migrations.MigrateUp(orm.DB, "0024_add_cron_spec_tables"))

	type CronSpec struct {
		ID           int32     `toml:"-" gorm:"primary_key"`
		CronSchedule string    `toml:"schedule"`
		CreatedAt    time.Time `toml:"-"`
		UpdatedAt    time.Time `toml:"-"`
	}

	cs := CronSpec{
		ID:           int32(1),
		CronSchedule: "0 0 0 1 1 *",
	}
	require.NoError(t, orm.DB.Create(&cs).Error)
	require.NoError(t, orm.DB.Find(&cs).Error)
	require.NoError(t, migrations.MigrateDownFrom(orm.DB, "0024_add_cron_spec_tables"))
}

type WebhookSpec struct {
	ID               int32     `toml:"-" gorm:"primary_key"`
	CreatedAt        time.Time `json:"createdAt" toml:"-"`
	UpdatedAt        time.Time `json:"updatedAt" toml:"-"`
	OnChainJobSpecID uuid.UUID
}

func TestMigrate_CreateWebhookTables(t *testing.T) {
	_, orm, cleanup := heavyweight.FullTestORM(t, "migrations_create_webhook_tables", false)
	defer cleanup()

	require.NoError(t, migrations.MigrateUp(orm.DB, "0029_add_webhook_spec_tables"))

	cs := WebhookSpec{
		ID:               int32(1),
		OnChainJobSpecID: uuid.FromStringOrNil("0x79BE667EF9DCBBAC55A06295CE870B07029BFCDB2DCE28D959F2815B16F8179800"),
	}
	require.NoError(t, orm.DB.Create(&cs).Error)
	require.NoError(t, orm.DB.Find(&cs).Error)
	require.NoError(t, migrations.MigrateDownFrom(orm.DB, "0029_add_webhook_spec_tables"))
}

func TestMigrate_ExternalJobID(t *testing.T) {
	_, orm, cleanup := heavyweight.FullTestORM(t, "migrations_external_jobid", false)
	defer cleanup()
	require.NoError(t, migrations.MigrateUp(orm.DB, "0035_create_feeds_managers"))
	cs := WebhookSpec{
		ID:               int32(1),
		OnChainJobSpecID: uuid.FromStringOrNil("0x79BE667EF9DCBBAC55A06295CE870B07029BFCDB2DCE28D959F2815B16F8179800"),
	}
	require.NoError(t, orm.DB.Create(&cs).Error)
	type Job struct {
		ID                            int32 `toml:"-" gorm:"primary_key"`
		OffchainreportingOracleSpecID *int32
		CronSpecID                    *int32
		DirectRequestSpecID           *int32
		FluxMonitorSpecID             *int32
		KeeperSpecID                  *int32
		VRFSpecID                     *int32
		WebhookSpecId                 *int32
		PipelineSpecID                int32
		PipelineSpec                  *pipeline.Spec
		JobSpecErrors                 []job.SpecError `gorm:"foreignKey:JobID"`
		Type                          job.Type
		SchemaVersion                 uint32
		Name                          null.String
		MaxTaskDuration               models.Interval
		Pipeline                      pipeline.Pipeline `toml:"observationSource" gorm:"-"`
	}
	var ps []pipeline.Spec
	for i := 0; i < 10; i++ {
		ps = append(ps, pipeline.Spec{
			DotDagSource: "",
		})
	}
	require.NoError(t, orm.DB.Create(&ps).Error)
	var jbs []Job
	for i := 0; i < 10; i++ {
		jbs = append(jbs, Job{
			SchemaVersion:  1,
			WebhookSpecId:  &cs.ID,
			PipelineSpecID: ps[i].ID,
			Type:           job.Webhook,
		})
	}
	require.NoError(t, orm.DB.Create(&jbs).Error)
	require.NoError(t, migrations.MigrateUp(orm.DB, "0036_external_job_id"))
	var jb2 []migrations.Job36
	require.NoError(t, orm.DB.Find(&jb2).Error)
	var seen = make(map[uuid.UUID]struct{})
	for i := range jb2 {
		if _, ok := seen[jb2[i].ExternalJobID]; ok {
			t.Error("all uuid's should be unique")
		}
		assert.NotEqual(t, uuid.UUID{}.String(), jb2[i].ExternalJobID.String())
		t.Log(jb2[i].ExternalJobID.String())
	}
	require.NoError(t, migrations.MigrateDownFrom(orm.DB, "0036_external_job_id"))
}

func TestMigrate_CascadeDeletes(t *testing.T) {
	_, orm, cleanup := heavyweight.FullTestORM(t, "migrations_cascade_deletes", false)
	t.Cleanup(cleanup)
	require.NoError(t, migrations.MigrateUp(orm.DB, "0037_cascade_deletes"))
	require.NoError(t, migrations.MigrateDownFrom(orm.DB, "0037_cascade_deletes"))
}
