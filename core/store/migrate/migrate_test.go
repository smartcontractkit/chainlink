package migrate_test

import (
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

func TestMigrate_0099_BootstrapConfigs(t *testing.T) {
	_, db := heavyweight.FullTestDB(t, "migrations", false, false)
	lggr := logger.TestLogger(t)
	cfg := configtest.NewTestGeneralConfig(t)
	err := goose.UpTo(db.DB, "migrations", 98)
	require.NoError(t, err)

	pipelineORM := pipeline.NewORM(db, lggr, cfg)
	pipelineID, err := pipelineORM.CreateSpec(pipeline.Pipeline{}, 0)
	require.NoError(t, err)
	jobORM := job.NewORM(db, nil, pipelineORM, nil, lggr, cfg)

	// OCR2 struct at migration v0098
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

	// Job struct at migration v0098
	type Job struct {
		ID                             int32     `toml:"-"`
		ExternalJobID                  uuid.UUID `toml:"externalJobID"`
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
		Pipeline                       pipeline.Pipeline `toml:"observationSource"`
		CreatedAt                      time.Time
	}

	spec := OffchainReporting2OracleSpec{
		ID:                                100,
		ContractID:                        "terra_187246hr3781h9fd198fh391g8f924",
		Relay:                             "evm",
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

	jb := Job{
		ID:                             10,
		ExternalJobID:                  uuid.NewV4(),
		Type:                           job.OffchainReporting2,
		SchemaVersion:                  1,
		PipelineSpecID:                 pipelineID,
		Offchainreporting2OracleSpecID: &spec.ID,
		Offchainreporting2OracleSpec:   &spec,
		BootstrapSpecID:                nil,
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

	jobInsert := `INSERT INTO jobs (id, pipeline_spec_id, external_job_id, schema_version, type, offchainreporting2_oracle_spec_id, bootstrap_spec_id, created_at)
		VALUES (:id, :pipeline_spec_id, :external_job_id, :schema_version, :type, :offchainreporting2_oracle_spec_id, :bootstrap_spec_id, NOW())
		RETURNING *;`

	_, err = db.NamedExec(jobInsert, jb)
	require.NoError(t, err)

	// Migrate up
	err = goose.UpByOne(db.DB, "migrations")
	require.NoError(t, err)

	jobs, count, err := jobORM.FindJobs(0, 1000)
	require.NoError(t, err)
	require.Equal(t, 1, count)

	migratedJob := jobs[0]
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
	require.Equal(t, 0, count)

	// Migrate down
	err = goose.Down(db.DB, "migrations")
	require.NoError(t, err)

	var oldJobs []Job
	sql = `SELECT * FROM jobs;`
	err = db.Select(&oldJobs, sql)
	require.NoError(t, err)

	require.Len(t, oldJobs, 1)
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
	require.Len(t, oldOCR2Spec, 1)

	require.Equal(t, spec.Relay, oldOCR2Spec[0].Relay)
	require.Equal(t, spec.ContractID, oldOCR2Spec[0].ContractID)
	require.Equal(t, spec.RelayConfig, oldOCR2Spec[0].RelayConfig)
	require.Equal(t, spec.ContractConfigConfirmations, oldOCR2Spec[0].ContractConfigConfirmations)
	require.Equal(t, spec.ContractConfigTrackerPollInterval, oldOCR2Spec[0].ContractConfigTrackerPollInterval)
	require.Equal(t, spec.BlockchainTimeout, oldOCR2Spec[0].BlockchainTimeout)
	require.True(t, oldOCR2Spec[0].IsBootstrapPeer)

	count = -1
	sql = `SELECT COUNT(*) FROM bootstrap_specs;`
	err = db.Get(&count, sql)
	require.NoError(t, err)
	require.Equal(t, 0, count)
}
