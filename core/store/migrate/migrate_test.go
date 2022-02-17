package migrate_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/lib/pq"
	"github.com/pressly/goose/v3"
	uuid "github.com/satori/go.uuid"
	"github.com/smartcontractkit/chainlink/core/internal/cltest/heavyweight"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/configtest"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	relaytypes "github.com/smartcontractkit/chainlink/core/services/relay/types"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v4"
)

func TestMigrate_0100_BootstrapConfigs(t *testing.T) {
	_, db := heavyweight.FullTestDB(t, "migrations", false, false)
	lggr := logger.TestLogger(t)
	cfg := configtest.NewTestGeneralConfig(t)
	err := goose.UpTo(db.DB, "migrations", 99)
	require.NoError(t, err)

	pipelineORM := pipeline.NewORM(db, lggr, cfg)
	pipelineID, err := pipelineORM.CreateSpec(pipeline.Pipeline{}, 0)
	require.NoError(t, err)
	pipelineID2, err := pipelineORM.CreateSpec(pipeline.Pipeline{}, 0)
	require.NoError(t, err)
	nonBootstrapPipelineID, err := pipelineORM.CreateSpec(pipeline.Pipeline{}, 0)
	require.NoError(t, err)
	jobORM := job.NewORM(db, nil, pipelineORM, nil, lggr, cfg)

	// OCR2 struct at migration v0099
	type OffchainReporting2OracleSpec struct {
		ID                                int32              `toml:"-"`
		ContractID                        string             `toml:"contractID"`
		Relay                             relaytypes.Network `toml:"relay"`
		RelayConfig                       job.RelayConfig    `toml:"relayConfig"`
		P2PBootstrapPeers                 pq.StringArray     `toml:"p2pBootstrapPeers"`
		OCRKeyBundleID                    null.String        `toml:"ocrKeyBundleID"`
		MonitoringEndpoint                null.String        `toml:"monitoringEndpoint"`
		TransmitterID                     null.String        `toml:"transmitterID"`
		BlockchainTimeout                 models.Interval    `toml:"blockchainTimeout"`
		ContractConfigTrackerPollInterval models.Interval    `toml:"contractConfigTrackerPollInterval"`
		ContractConfigConfirmations       uint16             `toml:"contractConfigConfirmations"`
		JuelsPerFeeCoinPipeline           string             `toml:"juelsPerFeeCoinSource"`
		IsBootstrapPeer                   bool
		CreatedAt                         time.Time `toml:"-"`
		UpdatedAt                         time.Time `toml:"-"`
	}

	// Job struct at migration v0099
	type Job struct {
		ID                             int32
		ExternalJobID                  uuid.UUID
		OffchainreportingOracleSpecID  *int32
		OffchainreportingOracleSpec    *job.OffchainReportingOracleSpec
		Offchainreporting2OracleSpecID *int32
		Offchainreporting2OracleSpec   *OffchainReporting2OracleSpec
		CronSpecID                     *int32
		CronSpec                       *job.CronSpec
		DirectRequestSpecID            *int32
		DirectRequestSpec              *job.DirectRequestSpec
		FluxMonitorSpecID              *int32
		FluxMonitorSpec                *job.FluxMonitorSpec
		KeeperSpecID                   *int32
		KeeperSpec                     *job.KeeperSpec
		VRFSpecID                      *int32
		VRFSpec                        *job.VRFSpec
		WebhookSpecID                  *int32
		WebhookSpec                    *job.WebhookSpec
		BlockhashStoreSpecID           *int32
		BlockhashStoreSpec             *job.BlockhashStoreSpec
		BootstrapSpec                  *job.BootstrapSpec
		BootstrapSpecID                *int32
		PipelineSpecID                 int32
		PipelineSpec                   *pipeline.Spec
		JobSpecErrors                  []job.SpecError
		Type                           job.Type
		SchemaVersion                  uint32
		Name                           null.String
		MaxTaskDuration                models.Interval
		Pipeline                       pipeline.Pipeline
		CreatedAt                      time.Time
	}

	spec := OffchainReporting2OracleSpec{
		ID:                                100,
		ContractID:                        "terra_187246hr3781h9fd198fh391g8f924",
		Relay:                             "terra",
		RelayConfig:                       job.RelayConfig{},
		P2PBootstrapPeers:                 pq.StringArray{""},
		OCRKeyBundleID:                    null.String{},
		MonitoringEndpoint:                null.StringFrom("endpoint:chainlink.monitor"),
		TransmitterID:                     null.String{},
		BlockchainTimeout:                 1337,
		ContractConfigTrackerPollInterval: 16,
		ContractConfigConfirmations:       32,
		JuelsPerFeeCoinPipeline:           "",
		IsBootstrapPeer:                   true,
	}
	spec2 := OffchainReporting2OracleSpec{
		ID:                                200,
		ContractID:                        "sol_187246hr3781h9fd198fh391g8f924",
		Relay:                             "sol",
		RelayConfig:                       job.RelayConfig{},
		P2PBootstrapPeers:                 pq.StringArray{""},
		OCRKeyBundleID:                    null.String{},
		MonitoringEndpoint:                null.StringFrom("endpoint:chain.link.monitor"),
		TransmitterID:                     null.String{},
		BlockchainTimeout:                 1338,
		ContractConfigTrackerPollInterval: 17,
		ContractConfigConfirmations:       33,
		JuelsPerFeeCoinPipeline:           "",
		IsBootstrapPeer:                   true,
	}

	jb := Job{
		ID:                             10,
		ExternalJobID:                  uuid.NewV4(),
		Type:                           job.OffchainReporting2,
		SchemaVersion:                  1,
		PipelineSpecID:                 pipelineID,
		Offchainreporting2OracleSpecID: &spec.ID,
		Offchainreporting2OracleSpec:   &spec,
	}

	jb2 := Job{
		ID:                             20,
		ExternalJobID:                  uuid.NewV4(),
		Type:                           job.OffchainReporting2,
		SchemaVersion:                  1,
		PipelineSpecID:                 pipelineID2,
		Offchainreporting2OracleSpecID: &spec2.ID,
		Offchainreporting2OracleSpec:   &spec2,
	}

	nonBootstrapSpec := OffchainReporting2OracleSpec{
		ID:                101,
		IsBootstrapPeer:   false,
		P2PBootstrapPeers: pq.StringArray{""},
		ContractID:        "empty",
	}
	nonBootstrapJob := Job{
		ID:                             11,
		ExternalJobID:                  uuid.NewV4(),
		Type:                           job.OffchainReporting2,
		SchemaVersion:                  1,
		PipelineSpecID:                 nonBootstrapPipelineID,
		Offchainreporting2OracleSpecID: &nonBootstrapSpec.ID,
		Offchainreporting2OracleSpec:   &nonBootstrapSpec,
	}

	sql := `INSERT INTO offchainreporting2_oracle_specs (id, contract_id, relay, relay_config, p2p_bootstrap_peers, ocr_key_bundle_id, transmitter_id,
					blockchain_timeout, contract_config_tracker_poll_interval, contract_config_confirmations, juels_per_fee_coin_pipeline, is_bootstrap_peer,
					monitoring_endpoint, created_at, updated_at)
			VALUES (:id, :contract_id, :relay, :relay_config, :p2p_bootstrap_peers, :ocr_key_bundle_id, :transmitter_id,
					 :blockchain_timeout, :contract_config_tracker_poll_interval, :contract_config_confirmations, :juels_per_fee_coin_pipeline, :is_bootstrap_peer,
					:monitoring_endpoint, NOW(), NOW())
			RETURNING id;`
	_, err = db.NamedExec(sql, jb.Offchainreporting2OracleSpec)
	require.NoError(t, err)
	_, err = db.NamedExec(sql, nonBootstrapJob.Offchainreporting2OracleSpec)
	require.NoError(t, err)
	_, err = db.NamedExec(sql, jb2.Offchainreporting2OracleSpec)
	require.NoError(t, err)

	sql = `INSERT INTO jobs (id, pipeline_spec_id, external_job_id, schema_version, type, offchainreporting2_oracle_spec_id, bootstrap_spec_id, created_at)
		VALUES (:id, :pipeline_spec_id, :external_job_id, :schema_version, :type, :offchainreporting2_oracle_spec_id, :bootstrap_spec_id, NOW())
		RETURNING *;`
	_, err = db.NamedExec(sql, jb)
	require.NoError(t, err)
	_, err = db.NamedExec(sql, nonBootstrapJob)
	require.NoError(t, err)
	_, err = db.NamedExec(sql, jb2)
	require.NoError(t, err)

	// Migrate up
	err = goose.UpByOne(db.DB, "migrations")
	require.NoError(t, err)

	var bootstrapSpecs []job.BootstrapSpec
	sql = `SELECT * FROM bootstrap_specs;`
	err = db.Select(&bootstrapSpecs, sql)
	require.NoError(t, err)
	require.Len(t, bootstrapSpecs, 2)
	fmt.Printf("bootstrap count %d\n", len(bootstrapSpecs))
	for _, bootstrapSpec := range bootstrapSpecs {
		fmt.Printf("bootstrap id: %d\n", bootstrapSpec.ID)
	}

	var jobs []job.Job
	jobs, count, err := jobORM.FindJobs(0, 1000)
	require.NoError(t, err)
	require.Equal(t, 3, count)
	fmt.Printf("jobs count %d\n", len(jobs))
	for _, jb := range jobs {
		fmt.Printf("job id: %d with BootstrapSpecID: %d\n", jb.ID, jb.BootstrapSpecID)
	}
	require.Nil(t, jobs[1].BootstrapSpecID)

	migratedJob := jobs[2]
	require.Nil(t, migratedJob.Offchainreporting2OracleSpecID)
	require.NotNil(t, migratedJob.BootstrapSpecID)
	require.Equal(t, &job.BootstrapSpec{
		ID:                                1,
		ContractID:                        spec.ContractID,
		Relay:                             spec.Relay,
		RelayConfig:                       spec.RelayConfig,
		MonitoringEndpoint:                spec.MonitoringEndpoint,
		BlockchainTimeout:                 spec.BlockchainTimeout,
		ContractConfigTrackerPollInterval: spec.ContractConfigTrackerPollInterval,
		ContractConfigConfirmations:       spec.ContractConfigConfirmations,
		CreatedAt:                         migratedJob.BootstrapSpec.CreatedAt,
		UpdatedAt:                         migratedJob.BootstrapSpec.UpdatedAt,
	}, migratedJob.BootstrapSpec)
	require.Equal(t, job.Bootstrap, migratedJob.Type)

	sql = `SELECT COUNT(*) FROM offchainreporting2_oracle_specs;`
	err = db.Get(&count, sql)
	require.NoError(t, err)
	require.Equal(t, 1, count)

	// Migrate down
	err = goose.Down(db.DB, "migrations")
	require.NoError(t, err)

	var oldJobs []Job
	sql = `SELECT * FROM jobs;`
	err = db.Select(&oldJobs, sql)
	require.NoError(t, err)
	require.Len(t, oldJobs, 3)

	revertedJob := oldJobs[0]
	require.NotNil(t, revertedJob.Offchainreporting2OracleSpecID)
	require.Nil(t, revertedJob.BootstrapSpecID)

	var oldOCR2Spec []OffchainReporting2OracleSpec
	sql = `SELECT contract_id, relay, relay_config, p2p_bootstrap_peers, ocr_key_bundle_id, transmitter_id,
					blockchain_timeout, contract_config_tracker_poll_interval, contract_config_confirmations, juels_per_fee_coin_pipeline, is_bootstrap_peer,
					monitoring_endpoint, created_at, updated_at
		FROM offchainreporting2_oracle_specs;`
	err = db.Select(&oldOCR2Spec, sql)
	require.NoError(t, err)
	require.Len(t, oldOCR2Spec, 3)
	bootSpec := oldOCR2Spec[1]

	require.Equal(t, spec.Relay, bootSpec.Relay)
	require.Equal(t, spec.ContractID, bootSpec.ContractID)
	require.Equal(t, spec.RelayConfig, bootSpec.RelayConfig)
	require.Equal(t, spec.ContractConfigConfirmations, bootSpec.ContractConfigConfirmations)
	require.Equal(t, spec.ContractConfigTrackerPollInterval, bootSpec.ContractConfigTrackerPollInterval)
	require.Equal(t, spec.BlockchainTimeout, bootSpec.BlockchainTimeout)
	require.True(t, bootSpec.IsBootstrapPeer)

	sql = `SELECT COUNT(*) FROM bootstrap_specs;`
	err = db.Get(&count, sql)
	require.NoError(t, err)
	require.Equal(t, 0, count)
}
