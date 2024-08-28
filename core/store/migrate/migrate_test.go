package migrate_test

import (
	"math/big"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink-common/pkg/types"

	evmcfg "github.com/smartcontractkit/chainlink/v2/core/chains/evm/config/toml"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	ubig "github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils/big"
	"github.com/smartcontractkit/chainlink/v2/core/config/env"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/configtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/v2/core/store/migrate"
	"github.com/smartcontractkit/chainlink/v2/core/store/models"
	"github.com/smartcontractkit/chainlink/v2/core/utils/testutils/heavyweight"
)

type OffchainReporting2OracleSpec100 struct {
	ID                                int32           `toml:"-"`
	ContractID                        string          `toml:"contractID"`
	Relay                             string          `toml:"relay"` // RelayID.Network
	RelayConfig                       job.JSONConfig  `toml:"relayConfig"`
	P2PBootstrapPeers                 pq.StringArray  `toml:"p2pBootstrapPeers"`
	OCRKeyBundleID                    null.String     `toml:"ocrKeyBundleID"`
	MonitoringEndpoint                null.String     `toml:"monitoringEndpoint"`
	TransmitterID                     null.String     `toml:"transmitterID"`
	BlockchainTimeout                 models.Interval `toml:"blockchainTimeout"`
	ContractConfigTrackerPollInterval models.Interval `toml:"contractConfigTrackerPollInterval"`
	ContractConfigConfirmations       uint16          `toml:"contractConfigConfirmations"`
	JuelsPerFeeCoinPipeline           string          `toml:"juelsPerFeeCoinSource"`
	CreatedAt                         time.Time       `toml:"-"`
	UpdatedAt                         time.Time       `toml:"-"`
}

func getOCR2Spec100() OffchainReporting2OracleSpec100 {
	return OffchainReporting2OracleSpec100{
		ID:                                100,
		ContractID:                        "terra_187246hr3781h9fd198fh391g8f924",
		Relay:                             "terra",
		RelayConfig:                       map[string]interface{}{"chainID": float64(1337)},
		P2PBootstrapPeers:                 pq.StringArray{""},
		OCRKeyBundleID:                    null.String{},
		MonitoringEndpoint:                null.StringFrom("endpoint:chainlink.monitor"),
		TransmitterID:                     null.String{},
		BlockchainTimeout:                 1337,
		ContractConfigTrackerPollInterval: 16,
		ContractConfigConfirmations:       32,
		JuelsPerFeeCoinPipeline: `ds1          [type=bridge name=voter_turnout];
	ds1_parse    [type=jsonparse path="one,two"];
	ds1_multiply [type=multiply times=1.23];
	ds1 -> ds1_parse -> ds1_multiply -> answer1;
	answer1      [type=median index=0];`,
	}
}

func TestMigrate_0100_BootstrapConfigs(t *testing.T) {
	cfg, db := heavyweight.FullTestDBEmptyV2(t, nil)
	lggr := logger.TestLogger(t)
	p, err := migrate.NewProvider(testutils.Context(t), db.DB)
	require.NoError(t, err)
	results, err := p.UpTo(testutils.Context(t), 99)
	require.NoError(t, err)
	assert.Len(t, results, 99)

	pipelineORM := pipeline.NewORM(db, lggr, cfg.JobPipeline().MaxSuccessfulRuns())
	ctx := testutils.Context(t)
	pipelineID, err := pipelineORM.CreateSpec(ctx, pipeline.Pipeline{}, 0)
	require.NoError(t, err)
	pipelineID2, err := pipelineORM.CreateSpec(ctx, pipeline.Pipeline{}, 0)
	require.NoError(t, err)
	nonBootstrapPipelineID, err := pipelineORM.CreateSpec(ctx, pipeline.Pipeline{}, 0)
	require.NoError(t, err)
	newFormatBoostrapPipelineID2, err := pipelineORM.CreateSpec(ctx, pipeline.Pipeline{}, 0)
	require.NoError(t, err)

	// OCR2 struct at migration v0099
	type OffchainReporting2OracleSpec struct {
		OffchainReporting2OracleSpec100
		IsBootstrapPeer bool
	}

	// Job struct at migration v0099
	type Job struct {
		job.Job
		OffchainreportingOracleSpecID  *int32
		Offchainreporting2OracleSpecID *int32
		Offchainreporting2OracleSpec   *OffchainReporting2OracleSpec
	}

	spec := OffchainReporting2OracleSpec{
		OffchainReporting2OracleSpec100: getOCR2Spec100(),
		IsBootstrapPeer:                 true,
	}
	spec2 := OffchainReporting2OracleSpec{
		OffchainReporting2OracleSpec100: OffchainReporting2OracleSpec100{
			ID:                                200,
			ContractID:                        "sol_187246hr3781h9fd198fh391g8f924",
			Relay:                             "sol",
			RelayConfig:                       job.JSONConfig{},
			P2PBootstrapPeers:                 pq.StringArray{""},
			OCRKeyBundleID:                    null.String{},
			MonitoringEndpoint:                null.StringFrom("endpoint:chain.link.monitor"),
			TransmitterID:                     null.String{},
			BlockchainTimeout:                 1338,
			ContractConfigTrackerPollInterval: 17,
			ContractConfigConfirmations:       33,
			JuelsPerFeeCoinPipeline:           "",
		},
		IsBootstrapPeer: true,
	}

	jb := Job{
		Job: job.Job{
			ID:             10,
			ExternalJobID:  uuid.New(),
			Type:           job.OffchainReporting2,
			SchemaVersion:  1,
			PipelineSpecID: pipelineID,
		},
		Offchainreporting2OracleSpec:   &spec,
		Offchainreporting2OracleSpecID: &spec.ID,
	}

	jb2 := Job{
		Job: job.Job{
			ID:             20,
			ExternalJobID:  uuid.New(),
			Type:           job.OffchainReporting2,
			SchemaVersion:  1,
			PipelineSpecID: pipelineID2,
		},
		Offchainreporting2OracleSpec:   &spec2,
		Offchainreporting2OracleSpecID: &spec2.ID,
	}

	nonBootstrapSpec := OffchainReporting2OracleSpec{
		OffchainReporting2OracleSpec100: OffchainReporting2OracleSpec100{
			ID:                101,
			P2PBootstrapPeers: pq.StringArray{""},
			ContractID:        "empty",
		},
		IsBootstrapPeer: false,
	}
	nonBootstrapJob := Job{
		Job: job.Job{
			ID:             11,
			ExternalJobID:  uuid.New(),
			Type:           job.OffchainReporting2,
			SchemaVersion:  1,
			PipelineSpecID: nonBootstrapPipelineID,
		},
		Offchainreporting2OracleSpec:   &nonBootstrapSpec,
		Offchainreporting2OracleSpecID: &nonBootstrapSpec.ID,
	}

	newFormatBoostrapSpec := job.BootstrapSpec{
		ID:                                1,
		ContractID:                        "evm_187246hr3781h9fd198fh391g8f924",
		Relay:                             "evm",
		RelayConfig:                       job.JSONConfig{},
		MonitoringEndpoint:                null.StringFrom("new:chain.link.monitor"),
		BlockchainTimeout:                 2448,
		ContractConfigTrackerPollInterval: 18,
		ContractConfigConfirmations:       34,
	}

	newFormatBootstrapJob := Job{
		Job: job.Job{
			ID:              30,
			ExternalJobID:   uuid.New(),
			Type:            job.Bootstrap,
			SchemaVersion:   1,
			PipelineSpecID:  newFormatBoostrapPipelineID2,
			BootstrapSpecID: &newFormatBoostrapSpec.ID,
			BootstrapSpec:   &newFormatBoostrapSpec,
		},
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

	sql = `INSERT INTO bootstrap_specs (contract_id, relay, relay_config, monitoring_endpoint,
					blockchain_timeout, contract_config_tracker_poll_interval, 
					contract_config_confirmations, created_at, updated_at)
			VALUES ( :contract_id, :relay, :relay_config, :monitoring_endpoint, 
					:blockchain_timeout, :contract_config_tracker_poll_interval, 
					:contract_config_confirmations, NOW(), NOW())
			RETURNING id;`

	_, err = db.NamedExec(sql, newFormatBootstrapJob.BootstrapSpec)
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
	_, err = db.NamedExec(sql, newFormatBootstrapJob)
	require.NoError(t, err)

	// Migrate up
	_, err = p.UpByOne(ctx)
	require.NoError(t, err)

	var bootstrapSpecs []job.BootstrapSpec
	sql = `SELECT * FROM bootstrap_specs;`
	err = db.Select(&bootstrapSpecs, sql)
	require.NoError(t, err)
	require.Len(t, bootstrapSpecs, 3)
	t.Logf("bootstrap count %d\n", len(bootstrapSpecs))
	for _, bootstrapSpec := range bootstrapSpecs {
		t.Logf("bootstrap id: %d\n", bootstrapSpec.ID)
	}

	var jobs []Job
	sql = `SELECT * FROM jobs ORDER BY created_at DESC, id DESC;`
	err = db.Select(&jobs, sql)

	require.NoError(t, err)
	require.Len(t, jobs, 4)
	t.Logf("jobs count %d\n", len(jobs))
	for _, jb := range jobs {
		t.Logf("job id: %d with BootstrapSpecID: %d\n", jb.ID, jb.BootstrapSpecID)
	}
	require.Nil(t, jobs[2].BootstrapSpecID)

	migratedJob := jobs[3]
	require.Nil(t, migratedJob.Offchainreporting2OracleSpecID)
	require.NotNil(t, migratedJob.BootstrapSpecID)

	var resultingBootstrapSpec job.BootstrapSpec
	err = db.Get(&resultingBootstrapSpec, `SELECT * FROM bootstrap_specs WHERE id = $1`, *migratedJob.BootstrapSpecID)
	migratedJob.BootstrapSpec = &resultingBootstrapSpec
	require.NoError(t, err)

	require.Equal(t, &job.BootstrapSpec{
		ID:                                2,
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
	var count int
	err = db.Get(&count, sql)
	require.NoError(t, err)
	require.Equal(t, 1, count)

	// Migrate down
	_, err = p.Down(ctx)
	require.NoError(t, err)

	var oldJobs []Job
	sql = `SELECT * FROM jobs;`
	err = db.Select(&oldJobs, sql)
	require.NoError(t, err)
	require.Len(t, oldJobs, 4)

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
	require.Len(t, oldOCR2Spec, 4)
	bootSpec := oldOCR2Spec[2]

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

	type jobIdAndContractId struct {
		ID         int32
		ContractID string
	}

	var jobsAndContracts []jobIdAndContractId
	sql = `SELECT jobs.id, ocr2.contract_id
FROM jobs 
INNER JOIN offchainreporting2_oracle_specs as ocr2 
ON jobs.offchainreporting2_oracle_spec_id = ocr2.id`
	err = db.Select(&jobsAndContracts, sql)
	require.NoError(t, err)

	require.Len(t, jobsAndContracts, 4)
	require.Equal(t, jobIdAndContractId{ID: 11, ContractID: "empty"}, jobsAndContracts[0])
	require.Equal(t, jobIdAndContractId{ID: 30, ContractID: "evm_187246hr3781h9fd198fh391g8f924"}, jobsAndContracts[1])
	require.Equal(t, jobIdAndContractId{ID: 10, ContractID: "terra_187246hr3781h9fd198fh391g8f924"}, jobsAndContracts[2])
	require.Equal(t, jobIdAndContractId{ID: 20, ContractID: "sol_187246hr3781h9fd198fh391g8f924"}, jobsAndContracts[3])
}

func TestMigrate_101_GenericOCR2(t *testing.T) {
	_, db := heavyweight.FullTestDBEmptyV2(t, nil)
	ctx := testutils.Context(t)
	p, err := migrate.NewProvider(ctx, db.DB)
	require.NoError(t, err)
	results, err := p.UpTo(ctx, 100)
	require.NoError(t, err)
	assert.Len(t, results, 100)

	sql := `INSERT INTO offchainreporting2_oracle_specs (id, contract_id, relay, relay_config, p2p_bootstrap_peers, ocr_key_bundle_id, transmitter_id,
					blockchain_timeout, contract_config_tracker_poll_interval, contract_config_confirmations, juels_per_fee_coin_pipeline,
					monitoring_endpoint, created_at, updated_at)
			VALUES (:id, :contract_id, :relay, :relay_config, :p2p_bootstrap_peers, :ocr_key_bundle_id, :transmitter_id,
					 :blockchain_timeout, :contract_config_tracker_poll_interval, :contract_config_confirmations, :juels_per_fee_coin_pipeline,
					:monitoring_endpoint, NOW(), NOW())
			RETURNING id;`

	spec := getOCR2Spec100()

	_, err = db.NamedExec(sql, spec)
	require.NoError(t, err)

	_, err = p.UpByOne(ctx)
	require.NoError(t, err)

	type PluginValues struct {
		PluginType   types.OCR2PluginType
		PluginConfig job.JSONConfig
	}

	var pluginValues PluginValues

	sql = `SELECT plugin_type, plugin_config FROM ocr2_oracle_specs`
	err = db.Get(&pluginValues, sql)
	require.NoError(t, err)

	require.Equal(t, types.Median, pluginValues.PluginType)
	require.Equal(t, job.JSONConfig{"juelsPerFeeCoinSource": spec.JuelsPerFeeCoinPipeline}, pluginValues.PluginConfig)

	_, err = p.Down(ctx)
	require.NoError(t, err)

	sql = `SELECT plugin_type, plugin_config FROM offchainreporting2_oracle_specs`
	err = db.Get(&pluginValues, sql)
	require.Error(t, err)

	var juels string
	sql = `SELECT juels_per_fee_coin_pipeline FROM offchainreporting2_oracle_specs`
	err = db.Get(&juels, sql)
	require.NoError(t, err)
	require.Equal(t, spec.JuelsPerFeeCoinPipeline, juels)
}

func TestMigrate(t *testing.T) {
	ctx := testutils.Context(t)
	_, db := heavyweight.FullTestDBEmptyV2(t, nil)

	p, err := migrate.NewProvider(ctx, db.DB)
	require.NoError(t, err)
	results, err := p.UpTo(ctx, 100)
	require.NoError(t, err)
	assert.Len(t, results, 100)

	err = migrate.Status(ctx, db.DB)
	require.NoError(t, err)

	ver, err := migrate.Current(ctx, db.DB)
	require.NoError(t, err)
	require.Equal(t, int64(100), ver)

	err = migrate.Migrate(ctx, db.DB)
	require.NoError(t, err)

	err = migrate.Rollback(ctx, db.DB, null.IntFrom(99))
	require.NoError(t, err)

	ver, err = migrate.Current(ctx, db.DB)
	require.NoError(t, err)
	require.Equal(t, int64(99), ver)
}

func TestSetMigrationENVVars(t *testing.T) {
	t.Run("ValidEVMConfig", func(t *testing.T) {
		chainID := ubig.New(big.NewInt(1337))
		testConfig := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
			evmEnabled := true
			c.EVM = evmcfg.EVMConfigs{&evmcfg.EVMConfig{
				ChainID: chainID,
				Enabled: &evmEnabled,
			}}
		})

		require.NoError(t, migrate.SetMigrationENVVars(testConfig))

		actualChainID := os.Getenv(env.EVMChainIDNotNullMigration0195)
		require.Equal(t, actualChainID, chainID.String())
	})

	t.Run("EVMConfigMissing", func(t *testing.T) {
		chainID := ubig.New(big.NewInt(1337))
		testConfig := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) { c.EVM = nil })

		require.NoError(t, migrate.SetMigrationENVVars(testConfig))

		actualChainID := os.Getenv(env.EVMChainIDNotNullMigration0195)
		require.Equal(t, actualChainID, chainID.String())
	})
}

func TestDatabaseBackFillWithMigration202(t *testing.T) {
	_, db := heavyweight.FullTestDBEmptyV2(t, nil)
	ctx := testutils.Context(t)

	p, err := migrate.NewProvider(ctx, db.DB)
	require.NoError(t, err)
	results, err := p.UpTo(ctx, 201)
	require.NoError(t, err)
	assert.Len(t, results, 201)

	simulatedOrm := logpoller.NewORM(testutils.SimulatedChainID, db, logger.TestLogger(t))
	require.NoError(t, simulatedOrm.InsertBlock(ctx, testutils.Random32Byte(), 10, time.Now(), 0), err)
	require.NoError(t, simulatedOrm.InsertBlock(ctx, testutils.Random32Byte(), 51, time.Now(), 0), err)
	require.NoError(t, simulatedOrm.InsertBlock(ctx, testutils.Random32Byte(), 90, time.Now(), 0), err)
	require.NoError(t, simulatedOrm.InsertBlock(ctx, testutils.Random32Byte(), 120, time.Now(), 23), err)

	baseOrm := logpoller.NewORM(big.NewInt(int64(84531)), db, logger.TestLogger(t))
	require.NoError(t, baseOrm.InsertBlock(ctx, testutils.Random32Byte(), 400, time.Now(), 0), err)

	klaytnOrm := logpoller.NewORM(big.NewInt(int64(1001)), db, logger.TestLogger(t))
	require.NoError(t, klaytnOrm.InsertBlock(ctx, testutils.Random32Byte(), 100, time.Now(), 0), err)

	_, err = p.UpTo(ctx, 202)
	require.NoError(t, err)

	tests := []struct {
		name                   string
		blockNumber            int64
		expectedFinalizedBlock int64
		orm                    logpoller.ORM
	}{
		{
			name:                   "last finalized block not changed if finality is too deep",
			blockNumber:            10,
			expectedFinalizedBlock: 0,
			orm:                    simulatedOrm,
		},
		{
			name:                   "last finalized block is updated for first block",
			blockNumber:            51,
			expectedFinalizedBlock: 1,
			orm:                    simulatedOrm,
		},
		{
			name:                   "last finalized block is updated",
			blockNumber:            90,
			expectedFinalizedBlock: 40,
			orm:                    simulatedOrm,
		},
		{
			name:                   "last finalized block is not changed when finality is set",
			blockNumber:            120,
			expectedFinalizedBlock: 23,
			orm:                    simulatedOrm,
		},
		{
			name:                   "use non default finality depth for chain 84531",
			blockNumber:            400,
			expectedFinalizedBlock: 200,
			orm:                    baseOrm,
		},
		{
			name:                   "use default finality depth for chain 1001",
			blockNumber:            100,
			expectedFinalizedBlock: 99,
			orm:                    klaytnOrm,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			block, err := tt.orm.SelectBlockByNumber(ctx, tt.blockNumber)
			require.NoError(t, err)
			require.Equal(t, tt.expectedFinalizedBlock, block.FinalizedBlockNumber)
		})
	}
}

func TestNoTriggers(t *testing.T) {
	_, db := heavyweight.FullTestDBEmptyV2(t, nil)

	assert_num_triggers := func(expected int) {
		row := db.DB.QueryRow("select count(*) from information_schema.triggers")
		var count int
		err := row.Scan(&count)

		require.NoError(t, err)
		require.Equal(t, expected, count)
	}

	// if you find yourself here and are tempted to add a trigger, something has gone wrong
	// and you should talk to the foundations team before proceeding
	assert_num_triggers(0)

	// version prior to removal of all triggers
	v := int64(217)
	p, err := migrate.NewProvider(testutils.Context(t), db.DB)
	require.NoError(t, err)
	_, err = p.UpTo(testutils.Context(t), v)
	require.NoError(t, err)
	assert_num_triggers(1)
}

func BenchmarkBackfillingRecordsWithMigration202(b *testing.B) {
	ctx := testutils.Context(b)
	previousMigration := int64(201)
	backfillMigration := int64(202)
	chainCount := 2
	// By default, log poller keeps up to 100_000 blocks in the database, this is the pessimistic case
	maxLogsSize := 100_000
	// Disable Goose logging for benchmarking
	goose.SetLogger(goose.NopLogger())
	_, db := heavyweight.FullTestDBEmptyV2(b, nil)

	p, err := migrate.NewProvider(ctx, db.DB)
	require.NoError(b, err)
	results, err := p.UpTo(ctx, previousMigration)
	require.NoError(b, err)
	assert.Len(b, results, int(previousMigration))

	for j := 0; j < chainCount; j++ {
		// Insert 100_000 block to database, can't do all at once, so batching by 10k
		var blocks []logpoller.LogPollerBlock
		for i := 0; i < maxLogsSize; i++ {
			blocks = append(blocks, logpoller.LogPollerBlock{
				EvmChainId:           ubig.NewI(int64(j + 1)),
				BlockHash:            testutils.Random32Byte(),
				BlockNumber:          int64(i + 1000),
				FinalizedBlockNumber: 0,
			})
		}
		batchInsertSize := 10_000
		for i := 0; i < maxLogsSize; i += batchInsertSize {
			start, end := i, i+batchInsertSize
			if end > maxLogsSize {
				end = maxLogsSize
			}

			_, err = db.NamedExecContext(ctx, `
			INSERT INTO evm.log_poller_blocks
				(evm_chain_id, block_hash, block_number, finalized_block_number, block_timestamp, created_at)
			VALUES 
				(:evm_chain_id, :block_hash, :block_number, :finalized_block_number, NOW(), NOW())
			ON CONFLICT DO NOTHING`, blocks[start:end])
			require.NoError(b, err)
		}
	}

	b.ResetTimer()

	// 1. Measure time of migration 200
	// 2. Goose down to 199
	// 3. Reset last_finalized_block_number to 0
	// Repeat 1-3
	for i := 0; i < b.N; i++ {
		b.StartTimer()
		_, err = p.UpTo(ctx, backfillMigration)
		require.NoError(b, err)
		b.StopTimer()

		// Cleanup
		_, err = p.DownTo(ctx, previousMigration)
		require.NoError(b, err)

		_, err = db.ExecContext(ctx, `
			UPDATE evm.log_poller_blocks
			SET finalized_block_number = 0`)
		require.NoError(b, err)
	}
}

func TestRollback_247_TxStateEnumUpdate(t *testing.T) {
	ctx := testutils.Context(t)
	_, db := heavyweight.FullTestDBV2(t, nil)
	p, err := migrate.NewProvider(ctx, db.DB)
	require.NoError(t, err)
	_, err = p.DownTo(ctx, 54)
	require.NoError(t, err)
	_, err = p.UpTo(ctx, 247)
	require.NoError(t, err)
}
