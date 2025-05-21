package job_test

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/pelletier/go-toml/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/jsonserializable"
	pkgworkflows "github.com/smartcontractkit/chainlink-common/pkg/workflows"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/artifacts"

	"github.com/smartcontractkit/chainlink-evm/pkg/assets"
	configtoml "github.com/smartcontractkit/chainlink-evm/pkg/config/toml"
	"github.com/smartcontractkit/chainlink-evm/pkg/keys"
	evmtypes "github.com/smartcontractkit/chainlink-evm/pkg/types"
	"github.com/smartcontractkit/chainlink-evm/pkg/utils/big"

	"github.com/smartcontractkit/chainlink/v2/core/bridges"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/configtest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/evmtest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/blockhashstore"
	"github.com/smartcontractkit/chainlink/v2/core/services/blockheaderfeeder"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/services/directrequest"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/keeper"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr"
	ocr2validate "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/validate"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocrbootstrap"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay"
	"github.com/smartcontractkit/chainlink/v2/core/services/standardcapabilities"
	"github.com/smartcontractkit/chainlink/v2/core/services/streams"
	"github.com/smartcontractkit/chainlink/v2/core/services/vrf/vrfcommon"
	"github.com/smartcontractkit/chainlink/v2/core/services/webhook"
	"github.com/smartcontractkit/chainlink/v2/core/testdata/testspecs"
	"github.com/smartcontractkit/chainlink/v2/core/utils/testutils/heavyweight"
)

const mercuryOracleTOML = `name = 'LINK / ETH | 0x0000000000000000000000000000000000000000000000000000000000000001 | verifier_proxy 0x0000000000000000000000000000000000000001'
type = 'offchainreporting2'
schemaVersion = 1
externalJobID = '00000000-0000-0000-0000-000000000001'
contractID = '0x0000000000000000000000000000000000000006'
transmitterID = '%s'
feedID = '%s'
relay = 'evm'
pluginType = 'mercury'
observationSource = """
	ds          [type=http method=GET url="https://chain.link/ETH-USD"];
	ds_parse    [type=jsonparse path="data.price" separator="."];
	ds_multiply [type=multiply times=100];
	ds -> ds_parse -> ds_multiply;
"""

[relayConfig]
chainID = 1
fromBlock = 1000

[onchainSigningStrategy]
strategyName = 'single-chain'
[onchainSigningStrategy.config]
publicKey = '8fa807463ad73f9ee855cfd60ba406dcf98a2855b3dd8af613107b0f6890a707'

[pluginConfig]
serverURL = 'wss://localhost:8080'
serverPubKey = '8fa807463ad73f9ee855cfd60ba406dcf98a2855b3dd8af613107b0f6890a707'
`

func TestORM(t *testing.T) {
	t.Parallel()

	config := configtest.NewTestGeneralConfig(t)
	db := pgtest.NewSqlxDB(t)
	keyStore := cltest.NewKeyStore(t, db)
	ethKeyStore := keyStore.Eth()

	func() {
		ctx := testutils.Context(t)
		require.NoError(t, keyStore.OCR().Add(ctx, cltest.DefaultOCRKey))
		require.NoError(t, keyStore.P2P().Add(ctx, cltest.DefaultP2PKey))
	}()

	pipelineORM := pipeline.NewORM(db, logger.TestLogger(t), config.JobPipeline().MaxSuccessfulRuns())
	borm := bridges.NewORM(db)
	orm := NewTestORM(t, db, pipelineORM, borm, keyStore)

	_, bridge := cltest.MustCreateBridge(t, db, cltest.BridgeOpts{})
	_, bridge2 := cltest.MustCreateBridge(t, db, cltest.BridgeOpts{})
	_, address := cltest.MustInsertRandomKey(t, ethKeyStore)
	jb := makeOCRJobSpec(t, address, bridge.Name.String(), bridge2.Name.String())

	t.Run("it creates job specs", func(t *testing.T) {
		err := orm.CreateJob(testutils.Context(t), jb)
		require.NoError(t, err)

		var returnedSpec job.Job
		var OCROracleSpec job.OCROracleSpec

		err = db.Get(&returnedSpec, "SELECT * FROM jobs WHERE jobs.id = $1", jb.ID)
		require.NoError(t, err)
		err = db.Get(&OCROracleSpec, "SELECT * FROM ocr_oracle_specs WHERE ocr_oracle_specs.id = $1", jb.OCROracleSpecID)
		require.NoError(t, err)
		returnedSpec.OCROracleSpec = &OCROracleSpec
		compareOCRJobSpecs(t, *jb, returnedSpec)
	})

	t.Run("it correctly mark job_pipeline_specs as primary when creating a job", func(t *testing.T) {
		ctx := testutils.Context(t)
		jb2 := makeOCRJobSpec(t, address, bridge.Name.String(), bridge2.Name.String())
		err := orm.CreateJob(ctx, jb2)
		require.NoError(t, err)

		var pipelineSpec pipeline.Spec
		err = db.Get(&pipelineSpec, "SELECT pipeline_specs.* FROM pipeline_specs JOIN job_pipeline_specs ON (pipeline_specs.id = job_pipeline_specs.pipeline_spec_id) WHERE job_pipeline_specs.job_id = $1", jb2.ID)
		require.NoError(t, err)
		var jobPipelineSpec job.PipelineSpec
		err = db.Get(&jobPipelineSpec, "SELECT * FROM job_pipeline_specs WHERE job_id = $1 AND pipeline_spec_id = $2", jb2.ID, pipelineSpec.ID)
		require.NoError(t, err)

		// `jb2.PipelineSpecID` gets loaded when calling `orm.CreateJob()` so we can compare it directly
		assert.Equal(t, jb2.PipelineSpecID, pipelineSpec.ID)
		assert.True(t, jobPipelineSpec.IsPrimary)
	})

	t.Run("autogenerates external job ID if missing", func(t *testing.T) {
		jb2 := makeOCRJobSpec(t, address, bridge.Name.String(), bridge2.Name.String())
		jb2.ExternalJobID = uuid.UUID{}
		err := orm.CreateJob(testutils.Context(t), jb2)
		require.NoError(t, err)

		var returnedSpec job.Job
		err = db.Get(&returnedSpec, "SELECT * FROM jobs WHERE jobs.id = $1", jb.ID)
		require.NoError(t, err)

		assert.NotEqual(t, uuid.UUID{}, returnedSpec.ExternalJobID)
	})

	t.Run("it deletes jobs from the DB", func(t *testing.T) {
		var dbSpecs []job.Job

		err := db.Select(&dbSpecs, "SELECT * FROM jobs")
		require.NoError(t, err)
		require.Len(t, dbSpecs, 3)

		err = orm.DeleteJob(testutils.Context(t), jb.ID, jb.Type)
		require.NoError(t, err)

		dbSpecs = []job.Job{}
		err = db.Select(&dbSpecs, "SELECT * FROM jobs")
		require.NoError(t, err)
		require.Len(t, dbSpecs, 2)
	})

	t.Run("increase job spec error occurrence", func(t *testing.T) {
		ctx := testutils.Context(t)
		jb3 := makeOCRJobSpec(t, address, bridge.Name.String(), bridge2.Name.String())
		err := orm.CreateJob(ctx, jb3)
		require.NoError(t, err)
		var jobSpec job.Job
		err = db.Get(&jobSpec, "SELECT * FROM jobs")
		require.NoError(t, err)

		ocrSpecError1 := "ocr spec 1 errored"
		ocrSpecError2 := "ocr spec 2 errored"
		require.NoError(t, orm.RecordError(ctx, jobSpec.ID, ocrSpecError1))
		require.NoError(t, orm.RecordError(ctx, jobSpec.ID, ocrSpecError1))
		require.NoError(t, orm.RecordError(ctx, jobSpec.ID, ocrSpecError2))

		var specErrors []job.SpecError
		err = db.Select(&specErrors, "SELECT * FROM job_spec_errors")
		require.NoError(t, err)
		require.Len(t, specErrors, 2)

		assert.Equal(t, uint(2), specErrors[0].Occurrences)
		assert.Equal(t, uint(1), specErrors[1].Occurrences)
		assert.True(t, specErrors[0].CreatedAt.Before(specErrors[0].UpdatedAt), "expected created_at (%s) to be before updated_at (%s)", specErrors[0].CreatedAt, specErrors[0].UpdatedAt)
		assert.Equal(t, specErrors[0].Description, ocrSpecError1)
		assert.Equal(t, specErrors[1].Description, ocrSpecError2)
		assert.True(t, specErrors[1].CreatedAt.After(specErrors[0].UpdatedAt))
		var j2 job.Job
		var OCROracleSpec job.OCROracleSpec
		var jobSpecErrors []job.SpecError

		err = db.Get(&j2, "SELECT * FROM jobs WHERE jobs.id = $1", jobSpec.ID)
		require.NoError(t, err)
		err = db.Get(&OCROracleSpec, "SELECT * FROM ocr_oracle_specs WHERE ocr_oracle_specs.id = $1", j2.OCROracleSpecID)
		require.NoError(t, err)
		err = db.Select(&jobSpecErrors, "SELECT * FROM job_spec_errors WHERE job_spec_errors.job_id = $1", j2.ID)
		require.NoError(t, err)
		require.Len(t, jobSpecErrors, 2)
	})

	t.Run("finds job spec error by ID", func(t *testing.T) {
		ctx := testutils.Context(t)
		jb3 := makeOCRJobSpec(t, address, bridge.Name.String(), bridge2.Name.String())
		err := orm.CreateJob(ctx, jb3)
		require.NoError(t, err)
		var jobSpec job.Job
		err = db.Get(&jobSpec, "SELECT * FROM jobs")
		require.NoError(t, err)

		var specErrors []job.SpecError
		err = db.Select(&specErrors, "SELECT * FROM job_spec_errors")
		require.NoError(t, err)
		require.Len(t, specErrors, 2)

		ocrSpecError1 := "ocr spec 3 errored"
		ocrSpecError2 := "ocr spec 4 errored"
		require.NoError(t, orm.RecordError(ctx, jobSpec.ID, ocrSpecError1))
		require.NoError(t, orm.RecordError(ctx, jobSpec.ID, ocrSpecError2))

		var updatedSpecError []job.SpecError

		err = db.Select(&updatedSpecError, "SELECT * FROM job_spec_errors ORDER BY id ASC")
		require.NoError(t, err)
		require.Len(t, updatedSpecError, 4)

		assert.Equal(t, uint(1), updatedSpecError[2].Occurrences)
		assert.Equal(t, uint(1), updatedSpecError[3].Occurrences)
		assert.Equal(t, ocrSpecError1, updatedSpecError[2].Description)
		assert.Equal(t, ocrSpecError2, updatedSpecError[3].Description)

		dbSpecErr1, err := orm.FindSpecError(ctx, updatedSpecError[2].ID)
		require.NoError(t, err)
		dbSpecErr2, err := orm.FindSpecError(ctx, updatedSpecError[3].ID)
		require.NoError(t, err)

		assert.Equal(t, uint(1), dbSpecErr1.Occurrences)
		assert.Equal(t, uint(1), dbSpecErr2.Occurrences)
		assert.Equal(t, ocrSpecError1, dbSpecErr1.Description)
		assert.Equal(t, ocrSpecError2, dbSpecErr2.Description)
	})

	t.Run("creates a job with a direct request spec", func(t *testing.T) {
		drSpec := fmt.Sprintf(`
		type                = "directrequest"
		schemaVersion       = 1
		evmChainID          = "0"
		name                = "example eth request event spec"
		contractAddress     = "0x613a38AC1659769640aaE063C651F48E0250454C"
		externalJobID       = "%s"
		observationSource   = """
		    ds1          [type=http method=GET url="http://example.com" allowunrestrictednetworkaccess="true"];
		    ds1_merge    [type=merge left="{}"]
		    ds1_parse    [type=jsonparse path="USD"];
		    ds1_multiply [type=multiply times=100];
		    ds1 -> ds1_parse -> ds1_multiply;
		"""
		`, uuid.New())

		drJob, err := directrequest.ValidatedDirectRequestSpec(drSpec)
		require.NoError(t, err)
		err = orm.CreateJob(testutils.Context(t), &drJob)
		require.NoError(t, err)
	})

	t.Run("creates webhook specs along with external_initiator_webhook_specs", func(t *testing.T) {
		ctx := testutils.Context(t)
		eiFoo := cltest.MustInsertExternalInitiator(t, borm)
		eiBar := cltest.MustInsertExternalInitiator(t, borm)

		eiWS := []webhook.TOMLWebhookSpecExternalInitiator{
			{Name: eiFoo.Name, Spec: cltest.JSONFromString(t, `{}`)},
			{Name: eiBar.Name, Spec: cltest.JSONFromString(t, `{"bar": 1}`)},
		}
		eim := webhook.NewExternalInitiatorManager(db, nil)
		jb, err := webhook.ValidatedWebhookSpec(ctx, testspecs.GenerateWebhookSpec(testspecs.WebhookSpecParams{ExternalInitiators: eiWS}).Toml(), eim)
		require.NoError(t, err)

		err = orm.CreateJob(testutils.Context(t), &jb)
		require.NoError(t, err)

		cltest.AssertCount(t, db, "external_initiator_webhook_specs", 2)
	})

	t.Run("it creates and deletes records for blockhash store jobs", func(t *testing.T) {
		ctx := testutils.Context(t)
		bhsJob, err := blockhashstore.ValidatedSpec(
			testspecs.GenerateBlockhashStoreSpec(testspecs.BlockhashStoreSpecParams{}).Toml())
		require.NoError(t, err)

		err = orm.CreateJob(ctx, &bhsJob)
		require.NoError(t, err)
		savedJob, err := orm.FindJob(testutils.Context(t), bhsJob.ID)
		require.NoError(t, err)
		require.Equal(t, bhsJob.ID, savedJob.ID)
		require.Equal(t, bhsJob.Type, savedJob.Type)
		require.Equal(t, bhsJob.BlockhashStoreSpec.ID, savedJob.BlockhashStoreSpec.ID)
		require.Equal(t, bhsJob.BlockhashStoreSpec.CoordinatorV1Address, savedJob.BlockhashStoreSpec.CoordinatorV1Address)
		require.Equal(t, bhsJob.BlockhashStoreSpec.CoordinatorV2Address, savedJob.BlockhashStoreSpec.CoordinatorV2Address)
		require.Equal(t, bhsJob.BlockhashStoreSpec.CoordinatorV2PlusAddress, savedJob.BlockhashStoreSpec.CoordinatorV2PlusAddress)
		require.Equal(t, bhsJob.BlockhashStoreSpec.WaitBlocks, savedJob.BlockhashStoreSpec.WaitBlocks)
		require.Equal(t, bhsJob.BlockhashStoreSpec.LookbackBlocks, savedJob.BlockhashStoreSpec.LookbackBlocks)
		require.Equal(t, bhsJob.BlockhashStoreSpec.HeartbeatPeriod, savedJob.BlockhashStoreSpec.HeartbeatPeriod)
		require.Equal(t, bhsJob.BlockhashStoreSpec.BlockhashStoreAddress, savedJob.BlockhashStoreSpec.BlockhashStoreAddress)
		require.Equal(t, bhsJob.BlockhashStoreSpec.TrustedBlockhashStoreAddress, savedJob.BlockhashStoreSpec.TrustedBlockhashStoreAddress)
		require.Equal(t, bhsJob.BlockhashStoreSpec.TrustedBlockhashStoreBatchSize, savedJob.BlockhashStoreSpec.TrustedBlockhashStoreBatchSize)
		require.Equal(t, bhsJob.BlockhashStoreSpec.PollPeriod, savedJob.BlockhashStoreSpec.PollPeriod)
		require.Equal(t, bhsJob.BlockhashStoreSpec.RunTimeout, savedJob.BlockhashStoreSpec.RunTimeout)
		require.Equal(t, bhsJob.BlockhashStoreSpec.EVMChainID, savedJob.BlockhashStoreSpec.EVMChainID)
		require.Equal(t, bhsJob.BlockhashStoreSpec.FromAddresses, savedJob.BlockhashStoreSpec.FromAddresses)
		err = orm.DeleteJob(ctx, bhsJob.ID, bhsJob.Type)
		require.NoError(t, err)
		_, err = orm.FindJob(testutils.Context(t), bhsJob.ID)
		require.Error(t, err)
	})

	t.Run("it creates and deletes records for blockheaderfeeder jobs", func(t *testing.T) {
		ctx := testutils.Context(t)
		bhsJob, err := blockheaderfeeder.ValidatedSpec(
			testspecs.GenerateBlockHeaderFeederSpec(testspecs.BlockHeaderFeederSpecParams{}).Toml())
		require.NoError(t, err)

		err = orm.CreateJob(ctx, &bhsJob)
		require.NoError(t, err)
		savedJob, err := orm.FindJob(testutils.Context(t), bhsJob.ID)
		require.NoError(t, err)
		require.Equal(t, bhsJob.ID, savedJob.ID)
		require.Equal(t, bhsJob.Type, savedJob.Type)
		require.Equal(t, bhsJob.BlockHeaderFeederSpec.ID, savedJob.BlockHeaderFeederSpec.ID)
		require.Equal(t, bhsJob.BlockHeaderFeederSpec.CoordinatorV1Address, savedJob.BlockHeaderFeederSpec.CoordinatorV1Address)
		require.Equal(t, bhsJob.BlockHeaderFeederSpec.CoordinatorV2Address, savedJob.BlockHeaderFeederSpec.CoordinatorV2Address)
		require.Equal(t, bhsJob.BlockHeaderFeederSpec.CoordinatorV2PlusAddress, savedJob.BlockHeaderFeederSpec.CoordinatorV2PlusAddress)
		require.Equal(t, bhsJob.BlockHeaderFeederSpec.WaitBlocks, savedJob.BlockHeaderFeederSpec.WaitBlocks)
		require.Equal(t, bhsJob.BlockHeaderFeederSpec.LookbackBlocks, savedJob.BlockHeaderFeederSpec.LookbackBlocks)
		require.Equal(t, bhsJob.BlockHeaderFeederSpec.BlockhashStoreAddress, savedJob.BlockHeaderFeederSpec.BlockhashStoreAddress)
		require.Equal(t, bhsJob.BlockHeaderFeederSpec.BatchBlockhashStoreAddress, savedJob.BlockHeaderFeederSpec.BatchBlockhashStoreAddress)
		require.Equal(t, bhsJob.BlockHeaderFeederSpec.PollPeriod, savedJob.BlockHeaderFeederSpec.PollPeriod)
		require.Equal(t, bhsJob.BlockHeaderFeederSpec.RunTimeout, savedJob.BlockHeaderFeederSpec.RunTimeout)
		require.Equal(t, bhsJob.BlockHeaderFeederSpec.EVMChainID, savedJob.BlockHeaderFeederSpec.EVMChainID)
		require.Equal(t, bhsJob.BlockHeaderFeederSpec.FromAddresses, savedJob.BlockHeaderFeederSpec.FromAddresses)
		require.Equal(t, bhsJob.BlockHeaderFeederSpec.GetBlockhashesBatchSize, savedJob.BlockHeaderFeederSpec.GetBlockhashesBatchSize)
		require.Equal(t, bhsJob.BlockHeaderFeederSpec.StoreBlockhashesBatchSize, savedJob.BlockHeaderFeederSpec.StoreBlockhashesBatchSize)
		err = orm.DeleteJob(ctx, bhsJob.ID, bhsJob.Type)
		require.NoError(t, err)
		_, err = orm.FindJob(testutils.Context(t), bhsJob.ID)
		require.Error(t, err)
	})
}

func TestORM_DeleteJob_DeletesAssociatedRecords(t *testing.T) {
	t.Parallel()
	ctx := testutils.Context(t)
	config := configtest.NewGeneralConfig(t, nil)

	db := pgtest.NewSqlxDB(t)
	keyStore := cltest.NewKeyStore(t, db)
	require.NoError(t, keyStore.OCR().Add(ctx, cltest.DefaultOCRKey))
	require.NoError(t, keyStore.P2P().Add(ctx, cltest.DefaultP2PKey))

	lggr := logger.TestLogger(t)
	pipelineORM := pipeline.NewORM(db, lggr, config.JobPipeline().MaxSuccessfulRuns())
	bridgesORM := bridges.NewORM(db)
	jobORM := NewTestORM(t, db, pipelineORM, bridgesORM, keyStore)
	korm := keeper.NewORM(db, logger.TestLogger(t))

	t.Run("it deletes records for offchainreporting jobs", func(t *testing.T) {
		ctx := testutils.Context(t)
		_, bridge := cltest.MustCreateBridge(t, db, cltest.BridgeOpts{})
		_, bridge2 := cltest.MustCreateBridge(t, db, cltest.BridgeOpts{})

		_, address := cltest.MustInsertRandomKey(t, keyStore.Eth())
		legacyChains := evmtest.NewLegacyChains(t, evmtest.TestChainOpts{
			ChainConfigs:   config.EVMConfigs(),
			DatabaseConfig: config.Database(),
			FeatureConfig:  config.Feature(),
			ListenerConfig: config.Database().Listener(),
			DB:             db,
			KeyStore:       keyStore.Eth(),
		})
		jb, err := ocr.ValidatedOracleSpecToml(config, legacyChains, testspecs.GenerateOCRSpec(testspecs.OCRSpecParams{
			TransmitterAddress: address.Hex(),
			DS1BridgeName:      bridge.Name.String(),
			DS2BridgeName:      bridge2.Name.String(),
		}).Toml())
		require.NoError(t, err)

		err = jobORM.CreateJob(ctx, &jb)
		require.NoError(t, err)

		cltest.AssertCount(t, db, "ocr_oracle_specs", 1)
		cltest.AssertCount(t, db, "pipeline_specs", 1)

		err = jobORM.DeleteJob(ctx, jb.ID, jb.Type)
		require.NoError(t, err)
		cltest.AssertCount(t, db, "ocr_oracle_specs", 0)
		cltest.AssertCount(t, db, "pipeline_specs", 0)
		cltest.AssertCount(t, db, "jobs", 0)
	})

	t.Run("it deletes records for keeper jobs", func(t *testing.T) {
		ctx := testutils.Context(t)
		registry, keeperJob := cltest.MustInsertKeeperRegistry(t, db, korm, keyStore.Eth(), 0, 1, 20)
		cltest.MustInsertUpkeepForRegistry(t, db, registry)

		cltest.AssertCount(t, db, "keeper_specs", 1)
		cltest.AssertCount(t, db, "keeper_registries", 1)
		cltest.AssertCount(t, db, "upkeep_registrations", 1)

		err := jobORM.DeleteJob(ctx, keeperJob.ID, keeperJob.Type)
		require.NoError(t, err)
		cltest.AssertCount(t, db, "keeper_specs", 0)
		cltest.AssertCount(t, db, "keeper_registries", 0)
		cltest.AssertCount(t, db, "upkeep_registrations", 0)
		cltest.AssertCount(t, db, "jobs", 0)
	})

	t.Run("it creates and deletes records for vrf jobs", func(t *testing.T) {
		ctx := testutils.Context(t)
		key, err := keyStore.VRF().Create(ctx)
		require.NoError(t, err)
		pk := key.PublicKey
		jb, err := vrfcommon.ValidatedVRFSpec(testspecs.GenerateVRFSpec(testspecs.VRFSpecParams{PublicKey: pk.String()}).Toml())
		require.NoError(t, err)

		err = jobORM.CreateJob(ctx, &jb)
		require.NoError(t, err)
		cltest.AssertCount(t, db, "vrf_specs", 1)
		cltest.AssertCount(t, db, "jobs", 1)
		err = jobORM.DeleteJob(ctx, jb.ID, jb.Type)
		require.NoError(t, err)
		cltest.AssertCount(t, db, "vrf_specs", 0)
		cltest.AssertCount(t, db, "jobs", 0)
	})

	t.Run("it deletes records for webhook jobs", func(t *testing.T) {
		ctx := testutils.Context(t)
		ei := cltest.MustInsertExternalInitiator(t, bridges.NewORM(db))
		jb, webhookSpec := cltest.MustInsertWebhookSpec(t, db)
		_, err := db.Exec(`INSERT INTO external_initiator_webhook_specs (external_initiator_id, webhook_spec_id, spec) VALUES ($1,$2,$3)`, ei.ID, webhookSpec.ID, `{"ei": "foo", "name": "webhookSpecTwoEIs"}`)
		require.NoError(t, err)

		err = jobORM.DeleteJob(ctx, jb.ID, jb.Type)
		require.NoError(t, err)
		cltest.AssertCount(t, db, "webhook_specs", 0)
		cltest.AssertCount(t, db, "external_initiator_webhook_specs", 0)
		cltest.AssertCount(t, db, "jobs", 0)
	})

	t.Run("it creates and deletes records for stream jobs", func(t *testing.T) {
		ctx := testutils.Context(t)
		jb, err := streams.ValidatedStreamSpec(testspecs.GenerateStreamSpec(testspecs.StreamSpecParams{Name: "Test-stream", StreamID: 1}).Toml())
		require.NoError(t, err)
		err = jobORM.CreateJob(ctx, &jb)
		require.NoError(t, err)
		cltest.AssertCount(t, db, "jobs", 1)
		err = jobORM.DeleteJob(ctx, jb.ID, jb.Type)
		require.NoError(t, err)
		cltest.AssertCount(t, db, "jobs", 0)
	})

	t.Run("does not allow to delete external initiators if they have referencing external_initiator_webhook_specs", func(t *testing.T) {
		// create new db because this will rollback transaction and poison it
		db := pgtest.NewSqlxDB(t)
		ei := cltest.MustInsertExternalInitiator(t, bridges.NewORM(db))
		_, webhookSpec := cltest.MustInsertWebhookSpec(t, db)
		_, err := db.Exec(`INSERT INTO external_initiator_webhook_specs (external_initiator_id, webhook_spec_id, spec) VALUES ($1,$2,$3)`, ei.ID, webhookSpec.ID, `{"ei": "foo", "name": "webhookSpecTwoEIs"}`)
		require.NoError(t, err)

		_, err = db.Exec(`DELETE FROM external_initiators`)
		require.EqualError(t, err, "ERROR: update or delete on table \"external_initiators\" violates foreign key constraint \"external_initiator_webhook_specs_external_initiator_id_fkey\" on table \"external_initiator_webhook_specs\" (SQLSTATE 23503)")
	})
}

func TestORM_CreateJob_VRFV2(t *testing.T) {
	ctx := testutils.Context(t)
	config := configtest.NewTestGeneralConfig(t)
	db := pgtest.NewSqlxDB(t)
	keyStore := cltest.NewKeyStore(t, db)
	require.NoError(t, keyStore.OCR().Add(ctx, cltest.DefaultOCRKey))

	lggr := logger.TestLogger(t)
	pipelineORM := pipeline.NewORM(db, lggr, config.JobPipeline().MaxSuccessfulRuns())
	bridgesORM := bridges.NewORM(db)

	jobORM := NewTestORM(t, db, pipelineORM, bridgesORM, keyStore)

	fromAddresses := []string{cltest.NewEIP55Address().String(), cltest.NewEIP55Address().String()}
	jb, err := vrfcommon.ValidatedVRFSpec(testspecs.GenerateVRFSpec(
		testspecs.VRFSpecParams{
			RequestedConfsDelay: 10,
			FromAddresses:       fromAddresses,
			ChunkSize:           25,
			BackoffInitialDelay: time.Minute,
			BackoffMaxDelay:     time.Hour,
			GasLanePrice:        assets.GWei(100),
			VRFOwnerAddress:     "0x32891BD79647DC9136Fc0a59AAB48c7825eb624c",
		}).
		Toml())
	require.NoError(t, err)

	require.NoError(t, jobORM.CreateJob(ctx, &jb))
	cltest.AssertCount(t, db, "vrf_specs", 1)
	cltest.AssertCount(t, db, "jobs", 1)
	var requestedConfsDelay int64
	require.NoError(t, db.Get(&requestedConfsDelay, `SELECT requested_confs_delay FROM vrf_specs LIMIT 1`))
	require.Equal(t, int64(10), requestedConfsDelay)
	var batchFulfillmentEnabled bool
	require.NoError(t, db.Get(&batchFulfillmentEnabled, `SELECT batch_fulfillment_enabled FROM vrf_specs LIMIT 1`))
	require.False(t, batchFulfillmentEnabled)
	var customRevertsPipelineEnabled bool
	require.NoError(t, db.Get(&customRevertsPipelineEnabled, `SELECT custom_reverts_pipeline_enabled FROM vrf_specs LIMIT 1`))
	require.False(t, customRevertsPipelineEnabled)
	var batchFulfillmentGasMultiplier float64
	require.NoError(t, db.Get(&batchFulfillmentGasMultiplier, `SELECT batch_fulfillment_gas_multiplier FROM vrf_specs LIMIT 1`))
	require.Equal(t, float64(1.0), batchFulfillmentGasMultiplier)
	var requestTimeout time.Duration
	require.NoError(t, db.Get(&requestTimeout, `SELECT request_timeout FROM vrf_specs LIMIT 1`))
	require.Equal(t, 24*time.Hour, requestTimeout)
	var backoffInitialDelay time.Duration
	require.NoError(t, db.Get(&backoffInitialDelay, `SELECT backoff_initial_delay FROM vrf_specs LIMIT 1`))
	require.Equal(t, time.Minute, backoffInitialDelay)
	var backoffMaxDelay time.Duration
	require.NoError(t, db.Get(&backoffMaxDelay, `SELECT backoff_max_delay FROM vrf_specs LIMIT 1`))
	require.Equal(t, time.Hour, backoffMaxDelay)
	var chunkSize int
	require.NoError(t, db.Get(&chunkSize, `SELECT chunk_size FROM vrf_specs LIMIT 1`))
	require.Equal(t, 25, chunkSize)
	var gasLanePrice assets.Wei
	require.NoError(t, db.Get(&gasLanePrice, `SELECT gas_lane_price FROM vrf_specs LIMIT 1`))
	require.Equal(t, jb.VRFSpec.GasLanePrice, &gasLanePrice)
	var fa pq.ByteaArray
	require.NoError(t, db.Get(&fa, `SELECT from_addresses FROM vrf_specs LIMIT 1`))
	var actual []string
	for _, b := range fa {
		actual = append(actual, common.BytesToAddress(b).String())
	}
	require.ElementsMatch(t, fromAddresses, actual)
	var vrfOwnerAddress evmtypes.EIP55Address
	require.NoError(t, db.Get(&vrfOwnerAddress, `SELECT vrf_owner_address FROM vrf_specs LIMIT 1`))
	require.Equal(t, "0x32891BD79647DC9136Fc0a59AAB48c7825eb624c", vrfOwnerAddress.Address().String())
	require.NoError(t, jobORM.DeleteJob(ctx, jb.ID, jb.Type))
	cltest.AssertCount(t, db, "vrf_specs", 0)
	cltest.AssertCount(t, db, "jobs", 0)

	jb, err = vrfcommon.ValidatedVRFSpec(testspecs.GenerateVRFSpec(testspecs.VRFSpecParams{RequestTimeout: 1 * time.Hour}).Toml())
	require.NoError(t, err)
	require.NoError(t, jobORM.CreateJob(ctx, &jb))
	cltest.AssertCount(t, db, "vrf_specs", 1)
	cltest.AssertCount(t, db, "jobs", 1)
	require.NoError(t, db.Get(&requestedConfsDelay, `SELECT requested_confs_delay FROM vrf_specs LIMIT 1`))
	require.Equal(t, int64(0), requestedConfsDelay)
	require.NoError(t, db.Get(&requestTimeout, `SELECT request_timeout FROM vrf_specs LIMIT 1`))
	require.Equal(t, 1*time.Hour, requestTimeout)
	require.NoError(t, jobORM.DeleteJob(ctx, jb.ID, jb.Type))
	cltest.AssertCount(t, db, "vrf_specs", 0)
	cltest.AssertCount(t, db, "jobs", 0)
}

func TestORM_CreateJob_VRFV2Plus(t *testing.T) {
	ctx := testutils.Context(t)
	config := configtest.NewTestGeneralConfig(t)
	db := pgtest.NewSqlxDB(t)
	keyStore := cltest.NewKeyStore(t, db)
	require.NoError(t, keyStore.OCR().Add(ctx, cltest.DefaultOCRKey))

	lggr := logger.TestLogger(t)
	pipelineORM := pipeline.NewORM(db, lggr, config.JobPipeline().MaxSuccessfulRuns())
	bridgesORM := bridges.NewORM(db)
	jobORM := NewTestORM(t, db, pipelineORM, bridgesORM, keyStore)

	fromAddresses := []string{cltest.NewEIP55Address().String(), cltest.NewEIP55Address().String()}
	jb, err := vrfcommon.ValidatedVRFSpec(testspecs.GenerateVRFSpec(
		testspecs.VRFSpecParams{
			VRFVersion:                   vrfcommon.V2Plus,
			RequestedConfsDelay:          10,
			FromAddresses:                fromAddresses,
			ChunkSize:                    25,
			BackoffInitialDelay:          time.Minute,
			BackoffMaxDelay:              time.Hour,
			GasLanePrice:                 assets.GWei(100),
			CustomRevertsPipelineEnabled: true,
		}).
		Toml())
	require.NoError(t, err)

	require.NoError(t, jobORM.CreateJob(ctx, &jb))
	cltest.AssertCount(t, db, "vrf_specs", 1)
	cltest.AssertCount(t, db, "jobs", 1)
	var requestedConfsDelay int64
	require.NoError(t, db.Get(&requestedConfsDelay, `SELECT requested_confs_delay FROM vrf_specs LIMIT 1`))
	require.Equal(t, int64(10), requestedConfsDelay)
	var batchFulfillmentEnabled bool
	require.NoError(t, db.Get(&batchFulfillmentEnabled, `SELECT batch_fulfillment_enabled FROM vrf_specs LIMIT 1`))
	require.False(t, batchFulfillmentEnabled)
	var customRevertsPipelineEnabled bool
	require.NoError(t, db.Get(&customRevertsPipelineEnabled, `SELECT custom_reverts_pipeline_enabled FROM vrf_specs LIMIT 1`))
	require.True(t, customRevertsPipelineEnabled)
	var batchFulfillmentGasMultiplier float64
	require.NoError(t, db.Get(&batchFulfillmentGasMultiplier, `SELECT batch_fulfillment_gas_multiplier FROM vrf_specs LIMIT 1`))
	require.Equal(t, float64(1.0), batchFulfillmentGasMultiplier)
	var requestTimeout time.Duration
	require.NoError(t, db.Get(&requestTimeout, `SELECT request_timeout FROM vrf_specs LIMIT 1`))
	require.Equal(t, 24*time.Hour, requestTimeout)
	var backoffInitialDelay time.Duration
	require.NoError(t, db.Get(&backoffInitialDelay, `SELECT backoff_initial_delay FROM vrf_specs LIMIT 1`))
	require.Equal(t, time.Minute, backoffInitialDelay)
	var backoffMaxDelay time.Duration
	require.NoError(t, db.Get(&backoffMaxDelay, `SELECT backoff_max_delay FROM vrf_specs LIMIT 1`))
	require.Equal(t, time.Hour, backoffMaxDelay)
	var chunkSize int
	require.NoError(t, db.Get(&chunkSize, `SELECT chunk_size FROM vrf_specs LIMIT 1`))
	require.Equal(t, 25, chunkSize)
	var gasLanePrice assets.Wei
	require.NoError(t, db.Get(&gasLanePrice, `SELECT gas_lane_price FROM vrf_specs LIMIT 1`))
	require.Equal(t, jb.VRFSpec.GasLanePrice, &gasLanePrice)
	var fa pq.ByteaArray
	require.NoError(t, db.Get(&fa, `SELECT from_addresses FROM vrf_specs LIMIT 1`))
	var actual []string
	for _, b := range fa {
		actual = append(actual, common.BytesToAddress(b).String())
	}
	require.ElementsMatch(t, fromAddresses, actual)
	var vrfOwnerAddress evmtypes.EIP55Address
	require.Error(t, db.Get(&vrfOwnerAddress, `SELECT vrf_owner_address FROM vrf_specs LIMIT 1`))
	require.NoError(t, jobORM.DeleteJob(ctx, jb.ID, jb.Type))
	cltest.AssertCount(t, db, "vrf_specs", 0)
	cltest.AssertCount(t, db, "jobs", 0)

	jb, err = vrfcommon.ValidatedVRFSpec(testspecs.GenerateVRFSpec(testspecs.VRFSpecParams{
		VRFVersion:     vrfcommon.V2Plus,
		RequestTimeout: 1 * time.Hour,
		FromAddresses:  fromAddresses,
	}).Toml())
	require.NoError(t, err)
	require.NoError(t, jobORM.CreateJob(ctx, &jb))
	cltest.AssertCount(t, db, "vrf_specs", 1)
	cltest.AssertCount(t, db, "jobs", 1)
	require.NoError(t, db.Get(&requestedConfsDelay, `SELECT requested_confs_delay FROM vrf_specs LIMIT 1`))
	require.Equal(t, int64(0), requestedConfsDelay)
	require.NoError(t, db.Get(&requestTimeout, `SELECT request_timeout FROM vrf_specs LIMIT 1`))
	require.Equal(t, 1*time.Hour, requestTimeout)
	require.NoError(t, jobORM.DeleteJob(ctx, jb.ID, jb.Type))
	cltest.AssertCount(t, db, "vrf_specs", 0)
	cltest.AssertCount(t, db, "jobs", 0)
}

func TestORM_CreateJob_OCRBootstrap(t *testing.T) {
	ctx := testutils.Context(t)
	config := configtest.NewTestGeneralConfig(t)
	db := pgtest.NewSqlxDB(t)
	keyStore := cltest.NewKeyStore(t, db)
	require.NoError(t, keyStore.OCR().Add(ctx, cltest.DefaultOCRKey))

	lggr := logger.TestLogger(t)
	pipelineORM := pipeline.NewORM(db, lggr, config.JobPipeline().MaxSuccessfulRuns())
	bridgesORM := bridges.NewORM(db)
	jobORM := NewTestORM(t, db, pipelineORM, bridgesORM, keyStore)

	jb, err := ocrbootstrap.ValidatedBootstrapSpecToml(testspecs.GetOCRBootstrapSpec())
	require.NoError(t, err)

	err = jobORM.CreateJob(ctx, &jb)
	require.NoError(t, err)
	cltest.AssertCount(t, db, "bootstrap_specs", 1)
	cltest.AssertCount(t, db, "jobs", 1)
	var relay string
	require.NoError(t, db.Get(&relay, `SELECT relay FROM bootstrap_specs LIMIT 1`))
	require.Equal(t, "evm", relay)

	require.NoError(t, jobORM.DeleteJob(ctx, jb.ID, jb.Type))
	cltest.AssertCount(t, db, "bootstrap_specs", 0)
	cltest.AssertCount(t, db, "jobs", 0)
}

func TestORM_CreateJob_EVMChainID_Validation(t *testing.T) {
	config := configtest.NewGeneralConfig(t, nil)
	db := pgtest.NewSqlxDB(t)
	keyStore := cltest.NewKeyStore(t, db)

	lggr := logger.TestLogger(t)
	pipelineORM := pipeline.NewORM(db, lggr, config.JobPipeline().MaxSuccessfulRuns())
	bridgesORM := bridges.NewORM(db)

	jobORM := NewTestORM(t, db, pipelineORM, bridgesORM, keyStore)

	t.Run("evm chain id validation for ocr works", func(t *testing.T) {
		jb := job.Job{
			Type:          job.OffchainReporting,
			OCROracleSpec: &job.OCROracleSpec{},
		}
		assert.Equal(t, "CreateJobFailed: evm chain id must be defined", jobORM.CreateJob(testutils.Context(t), &jb).Error())
	})

	t.Run("evm chain id validation for direct request works", func(t *testing.T) {
		jb := job.Job{
			Type:              job.DirectRequest,
			DirectRequestSpec: &job.DirectRequestSpec{},
		}
		assert.Equal(t, "CreateJobFailed: evm chain id must be defined", jobORM.CreateJob(testutils.Context(t), &jb).Error())
	})

	t.Run("evm chain id validation for flux monitor works", func(t *testing.T) {
		jb := job.Job{
			Type:            job.FluxMonitor,
			FluxMonitorSpec: &job.FluxMonitorSpec{},
		}
		assert.Equal(t, "CreateJobFailed: evm chain id must be defined", jobORM.CreateJob(testutils.Context(t), &jb).Error())
	})

	t.Run("evm chain id validation for keepers works", func(t *testing.T) {
		jb := job.Job{
			Type:       job.Keeper,
			KeeperSpec: &job.KeeperSpec{},
		}
		assert.Equal(t, "CreateJobFailed: evm chain id must be defined", jobORM.CreateJob(testutils.Context(t), &jb).Error())
	})

	t.Run("evm chain id validation for vrf works", func(t *testing.T) {
		jb := job.Job{
			Type:    job.VRF,
			VRFSpec: &job.VRFSpec{},
		}
		assert.Equal(t, "CreateJobFailed: evm chain id must be defined", jobORM.CreateJob(testutils.Context(t), &jb).Error())
	})

	t.Run("evm chain id validation for block hash store works", func(t *testing.T) {
		jb := job.Job{
			Type:               job.BlockhashStore,
			BlockhashStoreSpec: &job.BlockhashStoreSpec{},
		}
		assert.Equal(t, "CreateJobFailed: evm chain id must be defined", jobORM.CreateJob(testutils.Context(t), &jb).Error())
	})

	t.Run("evm chain id validation for block header feeder works", func(t *testing.T) {
		jb := job.Job{
			Type:                  job.BlockHeaderFeeder,
			BlockHeaderFeederSpec: &job.BlockHeaderFeederSpec{},
		}
		assert.Equal(t, "CreateJobFailed: evm chain id must be defined", jobORM.CreateJob(testutils.Context(t), &jb).Error())
	})

	t.Run("evm chain id validation for legacy gas station server spec works", func(t *testing.T) {
		jb := job.Job{
			Type:                       job.LegacyGasStationServer,
			LegacyGasStationServerSpec: &job.LegacyGasStationServerSpec{},
		}
		assert.Equal(t, "CreateJobFailed: evm chain id must be defined", jobORM.CreateJob(testutils.Context(t), &jb).Error())
	})

	t.Run("evm chain id validation for legacy gas station sidecar spec works", func(t *testing.T) {
		jb := job.Job{
			Type:                        job.LegacyGasStationSidecar,
			LegacyGasStationSidecarSpec: &job.LegacyGasStationSidecarSpec{},
		}
		assert.Equal(t, "CreateJobFailed: evm chain id must be defined", jobORM.CreateJob(testutils.Context(t), &jb).Error())
	})
}

func TestORM_CreateJob_OCR_DuplicatedContractAddress(t *testing.T) {
	ctx := testutils.Context(t)
	customChainID := big.New(testutils.NewRandomEVMChainID())

	config := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		enabled := true
		c.EVM = append(c.EVM, &configtoml.EVMConfig{
			ChainID: customChainID,
			Chain:   configtoml.Defaults(customChainID),
			Enabled: &enabled,
			Nodes:   configtoml.EVMNodes{{}},
		})
	})
	db := pgtest.NewSqlxDB(t)
	keyStore := cltest.NewKeyStore(t, db)
	require.NoError(t, keyStore.OCR().Add(ctx, cltest.DefaultOCRKey))

	lggr := logger.TestLogger(t)
	pipelineORM := pipeline.NewORM(db, lggr, config.JobPipeline().MaxSuccessfulRuns())
	bridgesORM := bridges.NewORM(db)

	jobORM := NewTestORM(t, db, pipelineORM, bridgesORM, keyStore)

	// defaultChainID is deprecated
	defaultChainID := customChainID
	_, address := cltest.MustInsertRandomKey(t, keyStore.Eth())
	_, bridge := cltest.MustCreateBridge(t, db, cltest.BridgeOpts{})
	_, bridge2 := cltest.MustCreateBridge(t, db, cltest.BridgeOpts{})

	// Custom Chain Job
	externalJobID := uuid.NullUUID{UUID: uuid.New(), Valid: true}
	spec := testspecs.GenerateOCRSpec(testspecs.OCRSpecParams{
		Name:               "job3",
		EVMChainID:         customChainID.String(),
		DS1BridgeName:      bridge.Name.String(),
		DS2BridgeName:      bridge2.Name.String(),
		TransmitterAddress: address.Hex(),
		JobID:              externalJobID.UUID.String(),
	})
	legacyChains := evmtest.NewLegacyChains(t, evmtest.TestChainOpts{
		ChainConfigs:   config.EVMConfigs(),
		DatabaseConfig: config.Database(),
		FeatureConfig:  config.Feature(),
		ListenerConfig: config.Database().Listener(),
		DB:             db,
		KeyStore:       keyStore.Eth(),
	})
	jb, err := ocr.ValidatedOracleSpecToml(config, legacyChains, spec.Toml())
	require.NoError(t, err)

	t.Run("with a set chain id", func(t *testing.T) {
		ctx := testutils.Context(t)
		err = jobORM.CreateJob(ctx, &jb) // Add job with custom chain id
		require.NoError(t, err)

		cltest.AssertCount(t, db, "ocr_oracle_specs", 1)
		cltest.AssertCount(t, db, "jobs", 1)

		externalJobID = uuid.NullUUID{UUID: uuid.New(), Valid: true}
		spec.JobID = externalJobID.UUID.String()
		jba, err := ocr.ValidatedOracleSpecToml(config, legacyChains, spec.Toml())
		require.NoError(t, err)
		err = jobORM.CreateJob(ctx, &jba) // Try to add duplicate job with default id
		require.Error(t, err)
		assert.Equal(t, fmt.Sprintf("CreateJobFailed: a job with contract address %s already exists for chain ID %s", jb.OCROracleSpec.ContractAddress, defaultChainID.String()), err.Error())

		externalJobID = uuid.NullUUID{UUID: uuid.New(), Valid: true}
		spec.JobID = externalJobID.UUID.String()
		jb2, err := ocr.ValidatedOracleSpecToml(config, legacyChains, spec.Toml())
		require.NoError(t, err)

		err = jobORM.CreateJob(ctx, &jb2) // Try to add duplicate job with custom id
		require.Error(t, err)
		assert.Equal(t, fmt.Sprintf("CreateJobFailed: a job with contract address %s already exists for chain ID %s", jb2.OCROracleSpec.ContractAddress, customChainID), err.Error())
	})
}

func TestORM_CreateJob_OCR2_DuplicatedContractAddress(t *testing.T) {
	ctx := testutils.Context(t)
	customChainID := big.New(testutils.NewRandomEVMChainID())

	config := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		enabled := true
		c.EVM = append(c.EVM, &configtoml.EVMConfig{
			ChainID: customChainID,
			Chain:   configtoml.Defaults(customChainID),
			Enabled: &enabled,
			Nodes:   configtoml.EVMNodes{{}},
		})
	})
	db := pgtest.NewSqlxDB(t)
	keyStore := cltest.NewKeyStore(t, db)
	require.NoError(t, keyStore.OCR2().Add(ctx, cltest.DefaultOCR2Key))

	lggr := logger.TestLogger(t)
	pipelineORM := pipeline.NewORM(db, lggr, config.JobPipeline().MaxSuccessfulRuns())
	bridgesORM := bridges.NewORM(db)

	jobORM := NewTestORM(t, db, pipelineORM, bridgesORM, keyStore)

	_, address := cltest.MustInsertRandomKey(t, keyStore.Eth())

	jb, err := ocr2validate.ValidatedOracleSpecToml(testutils.Context(t), config.OCR2(), config.Insecure(), testspecs.GetOCR2EVMSpecMinimal(), nil)
	require.NoError(t, err)

	const juelsPerFeeCoinSource = `
	ds          [type=http method=GET url="https://chain.link/ETH-USD"];
	ds_parse    [type=jsonparse path="data.price" separator="."];
	ds_multiply [type=multiply times=100];
	ds -> ds_parse -> ds_multiply;`

	jb.Name = null.StringFrom("Job 1")
	jb.OCR2OracleSpec.TransmitterID = null.StringFrom(address.String())
	jb.OCR2OracleSpec.PluginConfig["juelsPerFeeCoinSource"] = juelsPerFeeCoinSource

	err = jobORM.CreateJob(ctx, &jb)
	require.NoError(t, err)

	jb2, err := ocr2validate.ValidatedOracleSpecToml(testutils.Context(t), config.OCR2(), config.Insecure(), testspecs.GetOCR2EVMSpecMinimal(), nil)
	require.NoError(t, err)

	jb2.Name = null.StringFrom("Job with same chain id & contract address")
	jb2.OCR2OracleSpec.TransmitterID = null.StringFrom(address.String())
	jb.OCR2OracleSpec.PluginConfig["juelsPerFeeCoinSource"] = juelsPerFeeCoinSource

	err = jobORM.CreateJob(ctx, &jb2)
	require.Error(t, err)

	jb3, err := ocr2validate.ValidatedOracleSpecToml(testutils.Context(t), config.OCR2(), config.Insecure(), testspecs.GetOCR2EVMSpecMinimal(), nil)
	require.NoError(t, err)
	jb3.Name = null.StringFrom("Job with different chain id & same contract address")
	jb3.OCR2OracleSpec.TransmitterID = null.StringFrom(address.String())
	jb3.OCR2OracleSpec.RelayConfig["chainID"] = customChainID.Int64()
	jb.OCR2OracleSpec.PluginConfig["juelsPerFeeCoinSource"] = juelsPerFeeCoinSource

	err = jobORM.CreateJob(ctx, &jb3)
	require.Error(t, err)
}

func TestORM_CreateJob_OCR2_Sending_Keys_Transmitter_Keys_Validations(t *testing.T) {
	ctx := testutils.Context(t)
	customChainID := big.New(testutils.NewRandomEVMChainID())

	config := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		enabled := true
		c.EVM = append(c.EVM, &configtoml.EVMConfig{
			ChainID: customChainID,
			Chain:   configtoml.Defaults(customChainID),
			Enabled: &enabled,
			Nodes:   configtoml.EVMNodes{{}},
		})
	})
	db := pgtest.NewSqlxDB(t)
	keyStore := cltest.NewKeyStore(t, db)
	require.NoError(t, keyStore.OCR2().Add(ctx, cltest.DefaultOCR2Key))

	lggr := logger.TestLogger(t)
	pipelineORM := pipeline.NewORM(db, lggr, config.JobPipeline().MaxSuccessfulRuns())
	bridgesORM := bridges.NewORM(db)

	jobORM := NewTestORM(t, db, pipelineORM, bridgesORM, keyStore)

	jb, err := ocr2validate.ValidatedOracleSpecToml(testutils.Context(t), config.OCR2(), config.Insecure(), testspecs.GetOCR2EVMSpecMinimal(), nil)
	require.NoError(t, err)

	t.Run("sending keys or transmitterID must be defined", func(t *testing.T) {
		ctx := testutils.Context(t)
		jb.OCR2OracleSpec.TransmitterID = null.String{}
		assert.Equal(t, "CreateJobFailed: neither sending keys nor transmitter ID is defined", jobORM.CreateJob(ctx, &jb).Error())
	})

	_, address := cltest.MustInsertRandomKey(t, keyStore.Eth())
	t.Run("sending keys validation works properly", func(t *testing.T) {
		ctx := testutils.Context(t)
		jb.OCR2OracleSpec.TransmitterID = null.String{}
		_, address2 := cltest.MustInsertRandomKey(t, keyStore.Eth())
		jb.OCR2OracleSpec.RelayConfig["sendingKeys"] = interface{}([]any{address.String(), address2.String(), common.HexToAddress("0X0").String()})
		assert.Equal(t, "CreateJobFailed: no EVM key matching: \"0x0000000000000000000000000000000000000000\": no such sending key exists", jobORM.CreateJob(ctx, &jb).Error())

		jb.OCR2OracleSpec.RelayConfig["sendingKeys"] = interface{}([]any{1, 2, 3})
		assert.Equal(t, "CreateJobFailed: sending keys are of wrong type", jobORM.CreateJob(ctx, &jb).Error())
	})

	t.Run("sending keys and transmitter ID can't both be defined", func(t *testing.T) {
		ctx := testutils.Context(t)
		jb.OCR2OracleSpec.TransmitterID = null.StringFrom(address.String())
		jb.OCR2OracleSpec.RelayConfig["sendingKeys"] = interface{}([]any{address.String()})
		assert.Equal(t, "CreateJobFailed: sending keys and transmitter ID can't both be defined", jobORM.CreateJob(ctx, &jb).Error())
	})

	t.Run("transmitter validation works", func(t *testing.T) {
		ctx := testutils.Context(t)
		jb.OCR2OracleSpec.TransmitterID = null.StringFrom("transmitterID that doesn't have a match in key store")
		jb.OCR2OracleSpec.RelayConfig["sendingKeys"] = nil
		assert.Equal(t, "CreateJobFailed: no EVM key matching: \"transmitterID that doesn't have a match in key store\": no such transmitter key exists", jobORM.CreateJob(ctx, &jb).Error())
	})
}

func TestORM_ValidateKeyStoreMatch(t *testing.T) {
	ctx := testutils.Context(t)
	config := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {})

	keyStore := cltest.NewKeyStore(t, pgtest.NewSqlxDB(t))
	require.NoError(t, keyStore.OCR2().Add(ctx, cltest.DefaultOCR2Key))

	var jb job.Job
	{
		var err error
		jb, err = ocr2validate.ValidatedOracleSpecToml(testutils.Context(t), config.OCR2(), config.Insecure(), testspecs.GetOCR2EVMSpecMinimal(), nil)
		require.NoError(t, err)
	}

	t.Run("test ETH key validation", func(t *testing.T) {
		ctx := testutils.Context(t)
		jb.OCR2OracleSpec.Relay = relay.NetworkEVM
		err := job.ValidateKeyStoreMatch(ctx, jb.OCR2OracleSpec, keyStore, "bad key")
		require.EqualError(t, err, "no EVM key matching: \"bad key\"")

		_, evmKey := cltest.MustInsertRandomKey(t, keyStore.Eth())
		err = job.ValidateKeyStoreMatch(ctx, jb.OCR2OracleSpec, keyStore, evmKey.String())
		require.NoError(t, err)
	})

	t.Run("test Cosmos key validation", func(t *testing.T) {
		ctx := testutils.Context(t)
		jb.OCR2OracleSpec.Relay = relay.NetworkCosmos
		err := job.ValidateKeyStoreMatch(ctx, jb.OCR2OracleSpec, keyStore, "bad key")
		require.EqualError(t, err, "no Cosmos key matching: \"bad key\"")

		cosmosKey, err := keyStore.Cosmos().Create(ctx)
		require.NoError(t, err)
		err = job.ValidateKeyStoreMatch(ctx, jb.OCR2OracleSpec, keyStore, cosmosKey.ID())
		require.NoError(t, err)
	})

	t.Run("test Solana key validation", func(t *testing.T) {
		ctx := testutils.Context(t)
		jb.OCR2OracleSpec.Relay = relay.NetworkSolana

		err := job.ValidateKeyStoreMatch(ctx, jb.OCR2OracleSpec, keyStore, "bad key")
		require.EqualError(t, err, "no Solana key matching: \"bad key\"")

		solanaKey, err := keyStore.Solana().Create(ctx)
		require.NoError(t, err)
		err = job.ValidateKeyStoreMatch(ctx, jb.OCR2OracleSpec, keyStore, solanaKey.ID())
		require.NoError(t, err)
	})

	t.Run("test Starknet key validation", func(t *testing.T) {
		ctx := testutils.Context(t)
		jb.OCR2OracleSpec.Relay = relay.NetworkStarkNet
		err := job.ValidateKeyStoreMatch(ctx, jb.OCR2OracleSpec, keyStore, "bad key")
		require.EqualError(t, err, "no Starknet key matching: \"bad key\"")

		starkNetKey, err := keyStore.StarkNet().Create(ctx)
		require.NoError(t, err)
		err = job.ValidateKeyStoreMatch(ctx, jb.OCR2OracleSpec, keyStore, starkNetKey.ID())
		require.NoError(t, err)
	})

	t.Run(("test Aptos key validation"), func(t *testing.T) {
		ctx := testutils.Context(t)
		jb.OCR2OracleSpec.Relay = relay.NetworkAptos
		err := job.ValidateKeyStoreMatch(ctx, jb.OCR2OracleSpec, keyStore, "bad key")
		require.EqualError(t, err, "no Aptos key matching: \"bad key\"")

		aptosKey, err := keyStore.Aptos().Create(ctx)
		require.NoError(t, err)
		err = job.ValidateKeyStoreMatch(ctx, jb.OCR2OracleSpec, keyStore, aptosKey.ID())
		require.NoError(t, err)
	})

	t.Run(("test Tron key validation"), func(t *testing.T) {
		ctx := testutils.Context(t)
		jb.OCR2OracleSpec.Relay = relay.NetworkTron
		err := job.ValidateKeyStoreMatch(ctx, jb.OCR2OracleSpec, keyStore, "bad key")
		require.EqualError(t, err, "no Tron key matching: \"bad key\"")

		tronKey, err := keyStore.Tron().Create(ctx)
		require.NoError(t, err)
		err = job.ValidateKeyStoreMatch(ctx, jb.OCR2OracleSpec, keyStore, tronKey.ID())
		require.NoError(t, err)
	})

	t.Run("test Mercury ETH key validation", func(t *testing.T) {
		ctx := testutils.Context(t)
		jb.OCR2OracleSpec.PluginType = types.Mercury
		err := job.ValidateKeyStoreMatch(ctx, jb.OCR2OracleSpec, keyStore, "bad key")
		require.EqualError(t, err, "no CSA key matching: \"bad key\"")

		csaKey, err := keyStore.CSA().Create(ctx)
		require.NoError(t, err)
		err = job.ValidateKeyStoreMatch(ctx, jb.OCR2OracleSpec, keyStore, csaKey.ID())
		require.NoError(t, err)
	})
}

func Test_FindJobs(t *testing.T) {
	t.Parallel()
	ctx := testutils.Context(t)

	config := configtest.NewTestGeneralConfig(t)
	db := pgtest.NewSqlxDB(t)
	keyStore := cltest.NewKeyStore(t, db)
	require.NoError(t, keyStore.OCR().Add(ctx, cltest.DefaultOCRKey))
	require.NoError(t, keyStore.P2P().Add(ctx, cltest.DefaultP2PKey))

	pipelineORM := pipeline.NewORM(db, logger.TestLogger(t), config.JobPipeline().MaxSuccessfulRuns())
	bridgesORM := bridges.NewORM(db)

	orm := NewTestORM(t, db, pipelineORM, bridgesORM, keyStore)

	_, bridge := cltest.MustCreateBridge(t, db, cltest.BridgeOpts{})
	_, bridge2 := cltest.MustCreateBridge(t, db, cltest.BridgeOpts{})

	_, address := cltest.MustInsertRandomKey(t, keyStore.Eth())
	legacyChains := evmtest.NewLegacyChains(t, evmtest.TestChainOpts{
		ChainConfigs:   config.EVMConfigs(),
		DatabaseConfig: config.Database(),
		FeatureConfig:  config.Feature(),
		ListenerConfig: config.Database().Listener(),
		DB:             db,
		KeyStore:       keyStore.Eth(),
	})
	jb1, err := ocr.ValidatedOracleSpecToml(config, legacyChains,
		testspecs.GenerateOCRSpec(testspecs.OCRSpecParams{
			JobID:              uuid.New().String(),
			TransmitterAddress: address.Hex(),
			DS1BridgeName:      bridge.Name.String(),
			DS2BridgeName:      bridge2.Name.String(),
		}).Toml(),
	)
	require.NoError(t, err)

	err = orm.CreateJob(ctx, &jb1)
	require.NoError(t, err)

	jb2, err := directrequest.ValidatedDirectRequestSpec(
		testspecs.GetDirectRequestSpec(),
	)
	require.NoError(t, err)

	err = orm.CreateJob(ctx, &jb2)
	require.NoError(t, err)

	t.Run("jobs are ordered by latest first", func(t *testing.T) {
		jobs, count, err2 := orm.FindJobs(testutils.Context(t), 0, 2)
		require.NoError(t, err2)
		require.Len(t, jobs, 2)
		assert.Equal(t, 2, count)

		expectedJobs := []job.Job{jb2, jb1}

		for i, exp := range expectedJobs {
			assert.Equal(t, exp.ID, jobs[i].ID)
		}
	})

	t.Run("jobs respect pagination", func(t *testing.T) {
		jobs, count, err2 := orm.FindJobs(testutils.Context(t), 0, 1)
		require.NoError(t, err2)
		require.Len(t, jobs, 1)
		assert.Equal(t, 2, count)

		expectedJobs := []job.Job{jb2}

		for i, exp := range expectedJobs {
			assert.Equal(t, exp.ID, jobs[i].ID)
		}
	})
}

func Test_FindJob(t *testing.T) {
	t.Parallel()
	ctx := testutils.Context(t)

	// Create a config with multiple EVM chains. The test fixtures already load 1337
	// Additional chains will need additional fixture statements to add a chain to evm_chains.
	config := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		chainID := big.NewI(1337)
		enabled := true
		c.EVM = append(c.EVM, &configtoml.EVMConfig{
			ChainID: chainID,
			Chain:   configtoml.Defaults(chainID),
			Enabled: &enabled,
			Nodes:   configtoml.EVMNodes{{}},
		})
	})

	db := pgtest.NewSqlxDB(t)
	keyStore := cltest.NewKeyStore(t, db)
	require.NoError(t, keyStore.OCR().Add(ctx, cltest.DefaultOCRKey))
	require.NoError(t, keyStore.P2P().Add(ctx, cltest.DefaultP2PKey))
	require.NoError(t, keyStore.CSA().Add(ctx, cltest.DefaultCSAKey))

	pipelineORM := pipeline.NewORM(db, logger.TestLogger(t), config.JobPipeline().MaxSuccessfulRuns())
	bridgesORM := bridges.NewORM(db)

	orm := NewTestORM(t, db, pipelineORM, bridgesORM, keyStore)

	_, bridge := cltest.MustCreateBridge(t, db, cltest.BridgeOpts{})
	_, bridge2 := cltest.MustCreateBridge(t, db, cltest.BridgeOpts{})

	// Create two jobs.  Each job has the same Transmitter Address but on a different chain.
	// Must uniquely name the OCR Specs to properly insert a new job in the job table.
	externalJobID := uuid.New()
	_, address := cltest.MustInsertRandomKey(t, keyStore.Eth())
	legacyChains := evmtest.NewLegacyChains(t, evmtest.TestChainOpts{
		ChainConfigs:   config.EVMConfigs(),
		DatabaseConfig: config.Database(),
		FeatureConfig:  config.Feature(),
		ListenerConfig: config.Database().Listener(),
		DB:             db,
		KeyStore:       keyStore.Eth(),
	})
	job, err := ocr.ValidatedOracleSpecToml(config, legacyChains,
		testspecs.GenerateOCRSpec(testspecs.OCRSpecParams{
			JobID:              externalJobID.String(),
			Name:               "orig ocr spec",
			TransmitterAddress: address.Hex(),
			DS1BridgeName:      bridge.Name.String(),
			DS2BridgeName:      bridge2.Name.String(),
		}).Toml(),
	)
	require.NoError(t, err)

	jobSameAddress, err := ocr.ValidatedOracleSpecToml(config, legacyChains,
		testspecs.GenerateOCRSpec(testspecs.OCRSpecParams{
			JobID:              uuid.New().String(),
			TransmitterAddress: address.Hex(),
			Name:               "ocr spec dup addr",
			EVMChainID:         "1337",
			DS1BridgeName:      bridge.Name.String(),
			DS2BridgeName:      bridge2.Name.String(),
		}).Toml(),
	)
	require.NoError(t, err)

	jobOCR2, err := ocr2validate.ValidatedOracleSpecToml(testutils.Context(t), config.OCR2(), config.Insecure(), testspecs.GetOCR2EVMSpecMinimal(), nil)
	require.NoError(t, err)
	jobOCR2.OCR2OracleSpec.TransmitterID = null.StringFrom(address.String())

	const juelsPerFeeCoinSource = `
	ds          [type=http method=GET url="https://chain.link/ETH-USD"];
	ds_parse    [type=jsonparse path="data.price" separator="."];
	ds_multiply [type=multiply times=100];
	ds -> ds_parse -> ds_multiply;`

	jobOCR2.OCR2OracleSpec.PluginConfig["juelsPerFeeCoinSource"] = juelsPerFeeCoinSource

	ocr2WithFeedID1 := "0x0001000000000000000000000000000000000000000000000000000000000001"
	ocr2WithFeedID2 := "0x0001000000000000000000000000000000000000000000000000000000000002"
	jobOCR2WithFeedID1, err := ocr2validate.ValidatedOracleSpecToml(
		testutils.Context(t),
		config.OCR2(),
		config.Insecure(),
		fmt.Sprintf(mercuryOracleTOML, cltest.DefaultCSAKey.PublicKeyString(), ocr2WithFeedID1),
		nil,
	)
	require.NoError(t, err)

	jobOCR2WithFeedID2, err := ocr2validate.ValidatedOracleSpecToml(
		testutils.Context(t),
		config.OCR2(),
		config.Insecure(),
		fmt.Sprintf(mercuryOracleTOML, cltest.DefaultCSAKey.PublicKeyString(), ocr2WithFeedID2),
		nil,
	)
	jobOCR2WithFeedID2.ExternalJobID = uuid.New()
	jobOCR2WithFeedID2.Name = null.StringFrom("new name")
	require.NoError(t, err)

	err = orm.CreateJob(ctx, &job)
	require.NoError(t, err)

	err = orm.CreateJob(ctx, &jobSameAddress)
	require.NoError(t, err)

	err = orm.CreateJob(ctx, &jobOCR2)
	require.NoError(t, err)

	err = orm.CreateJob(ctx, &jobOCR2WithFeedID1)
	require.NoError(t, err)

	// second ocr2 job with same contract id but different feed id
	err = orm.CreateJob(ctx, &jobOCR2WithFeedID2)
	require.NoError(t, err)

	t.Run("by id", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(testutils.Context(t), 5*time.Second)
		defer cancel()
		jb, err2 := orm.FindJob(ctx, job.ID)
		require.NoError(t, err2)

		assert.Equal(t, jb.ID, job.ID)
		assert.Equal(t, jb.Name, job.Name)

		require.Positive(t, jb.PipelineSpecID)
		require.NotNil(t, jb.PipelineSpec)
		require.NotNil(t, jb.OCROracleSpecID)
		require.NotNil(t, jb.OCROracleSpec)
	})

	t.Run("by external job id", func(t *testing.T) {
		ctx := testutils.Context(t)
		jb, err2 := orm.FindJobByExternalJobID(ctx, externalJobID)
		require.NoError(t, err2)

		assert.Equal(t, jb.ID, job.ID)
		assert.Equal(t, jb.Name, job.Name)

		require.Positive(t, jb.PipelineSpecID)
		require.NotNil(t, jb.PipelineSpec)
		require.NotNil(t, jb.OCROracleSpecID)
		require.NotNil(t, jb.OCROracleSpec)
	})

	t.Run("by address", func(t *testing.T) {
		ctx := testutils.Context(t)
		jbID, err2 := orm.FindJobIDByAddress(ctx, job.OCROracleSpec.ContractAddress, job.OCROracleSpec.EVMChainID)
		require.NoError(t, err2)

		assert.Equal(t, job.ID, jbID)

		_, err2 = orm.FindJobIDByAddress(ctx, "not-existing", big.NewI(0))
		require.Error(t, err2)
		require.ErrorIs(t, err2, sql.ErrNoRows)
	})

	t.Run("by address yet chain scoped", func(t *testing.T) {
		ctx := testutils.Context(t)
		commonAddr := jobSameAddress.OCROracleSpec.ContractAddress

		// Find job ID for job on chain 1337 with common address.
		jbID, err2 := orm.FindJobIDByAddress(ctx, commonAddr, jobSameAddress.OCROracleSpec.EVMChainID)
		require.NoError(t, err2)

		assert.Equal(t, jobSameAddress.ID, jbID)

		// Find job ID for job on default evm chain with common address.
		jbID, err2 = orm.FindJobIDByAddress(ctx, commonAddr, job.OCROracleSpec.EVMChainID)
		require.NoError(t, err2)

		assert.Equal(t, job.ID, jbID)
	})

	t.Run("by contract id without feed id", func(t *testing.T) {
		ctx := testutils.Context(t)
		contractID := "0x613a38AC1659769640aaE063C651F48E0250454C"

		// Find job ID for ocr2 job without feedID.
		jbID, err2 := orm.FindOCR2JobIDByAddress(ctx, contractID, nil)
		require.NoError(t, err2)

		assert.Equal(t, jobOCR2.ID, jbID)
	})

	t.Run("by contract id with valid feed id", func(t *testing.T) {
		ctx := testutils.Context(t)
		contractID := "0x0000000000000000000000000000000000000006"
		feedID := common.HexToHash(ocr2WithFeedID1)

		// Find job ID for ocr2 job with feed ID
		jbID, err2 := orm.FindOCR2JobIDByAddress(ctx, contractID, &feedID)
		require.NoError(t, err2)

		assert.Equal(t, jobOCR2WithFeedID1.ID, jbID)
	})

	t.Run("with duplicate contract id but different feed id", func(t *testing.T) {
		ctx := testutils.Context(t)
		contractID := "0x0000000000000000000000000000000000000006"
		feedID := common.HexToHash(ocr2WithFeedID2)

		// Find job ID for ocr2 job with feed ID
		jbID, err2 := orm.FindOCR2JobIDByAddress(ctx, contractID, &feedID)
		require.NoError(t, err2)

		assert.Equal(t, jobOCR2WithFeedID2.ID, jbID)
	})
}

func Test_FindJobsByPipelineSpecIDs(t *testing.T) {
	t.Parallel()
	ctx := testutils.Context(t)

	config := configtest.NewTestGeneralConfig(t)
	db := pgtest.NewSqlxDB(t)
	keyStore := cltest.NewKeyStore(t, db)
	require.NoError(t, keyStore.OCR().Add(ctx, cltest.DefaultOCRKey))

	pipelineORM := pipeline.NewORM(db, logger.TestLogger(t), config.JobPipeline().MaxSuccessfulRuns())
	bridgesORM := bridges.NewORM(db)
	orm := NewTestORM(t, db, pipelineORM, bridgesORM, keyStore)

	jb, err := directrequest.ValidatedDirectRequestSpec(testspecs.GetDirectRequestSpec())
	require.NoError(t, err)
	jb.DirectRequestSpec.EVMChainID = big.NewI(0)

	err = orm.CreateJob(testutils.Context(t), &jb)
	require.NoError(t, err)

	t.Run("with jobs", func(t *testing.T) {
		ctx := testutils.Context(t)
		jbs, err2 := orm.FindJobsByPipelineSpecIDs(ctx, []int32{jb.PipelineSpecID})
		require.NoError(t, err2)
		assert.Len(t, jbs, 1)

		assert.Equal(t, jb.ID, jbs[0].ID)
		assert.Equal(t, jb.Name, jbs[0].Name)

		require.Positive(t, jbs[0].PipelineSpecID)
		require.Equal(t, jb.PipelineSpecID, jbs[0].PipelineSpecID)
		require.NotNil(t, jbs[0].PipelineSpec)
	})

	t.Run("without jobs", func(t *testing.T) {
		ctx := testutils.Context(t)
		jbs, err2 := orm.FindJobsByPipelineSpecIDs(ctx, []int32{-1})
		require.NoError(t, err2)
		assert.Empty(t, jbs)
	})

	t.Run("with chainID disabled", func(t *testing.T) {
		ctx := testutils.Context(t)
		orm2 := NewTestORM(t, db, pipelineORM, bridgesORM, keyStore)

		jbs, err2 := orm2.FindJobsByPipelineSpecIDs(ctx, []int32{jb.PipelineSpecID})
		require.NoError(t, err2)
		assert.Len(t, jbs, 1)
	})
}

func Test_FindPipelineRuns(t *testing.T) {
	t.Parallel()
	ctx := testutils.Context(t)

	config := configtest.NewTestGeneralConfig(t)
	db := pgtest.NewSqlxDB(t)
	keyStore := cltest.NewKeyStore(t, db)
	require.NoError(t, keyStore.OCR().Add(ctx, cltest.DefaultOCRKey))
	require.NoError(t, keyStore.P2P().Add(ctx, cltest.DefaultP2PKey))

	pipelineORM := pipeline.NewORM(db, logger.TestLogger(t), config.JobPipeline().MaxSuccessfulRuns())
	bridgesORM := bridges.NewORM(db)
	legacyChains := evmtest.NewLegacyChains(t, evmtest.TestChainOpts{
		ChainConfigs:   config.EVMConfigs(),
		DatabaseConfig: config.Database(),
		FeatureConfig:  config.Feature(),
		ListenerConfig: config.Database().Listener(),
		KeyStore:       keyStore.Eth(),
		DB:             db,
	})
	orm := NewTestORM(t, db, pipelineORM, bridgesORM, keyStore)

	_, bridge := cltest.MustCreateBridge(t, db, cltest.BridgeOpts{})
	_, bridge2 := cltest.MustCreateBridge(t, db, cltest.BridgeOpts{})

	externalJobID := uuid.New()
	_, address := cltest.MustInsertRandomKey(t, keyStore.Eth())
	jb, err := ocr.ValidatedOracleSpecToml(config, legacyChains,
		testspecs.GenerateOCRSpec(testspecs.OCRSpecParams{
			JobID:              externalJobID.String(),
			TransmitterAddress: address.Hex(),
			DS1BridgeName:      bridge.Name.String(),
			DS2BridgeName:      bridge2.Name.String(),
		}).Toml(),
	)
	require.NoError(t, err)

	err = orm.CreateJob(testutils.Context(t), &jb)
	require.NoError(t, err)

	t.Run("with no pipeline runs", func(t *testing.T) {
		ctx := testutils.Context(t)
		runs, count, err2 := orm.PipelineRuns(ctx, nil, 0, 10)
		require.NoError(t, err2)
		assert.Equal(t, 0, count)
		assert.Empty(t, runs)
	})

	t.Run("with a pipeline run", func(t *testing.T) {
		ctx := testutils.Context(t)
		run := mustInsertPipelineRun(t, pipelineORM, jb)

		runs, count, err2 := orm.PipelineRuns(ctx, nil, 0, 10)
		require.NoError(t, err2)

		assert.Equal(t, 1, count)
		actual := runs[0]

		// Test pipeline run fields
		assert.Equal(t, run.State, actual.State)
		assert.Equal(t, run.PipelineSpecID, actual.PipelineSpecID)

		// Test preloaded pipeline spec
		require.NotNil(t, jb.PipelineSpec)
		assert.Equal(t, jb.PipelineSpec.ID, actual.PipelineSpec.ID)
		assert.Equal(t, jb.ID, actual.PipelineSpec.JobID)
	})
}

func Test_PipelineRunsByJobID(t *testing.T) {
	t.Parallel()
	ctx := testutils.Context(t)

	config := configtest.NewTestGeneralConfig(t)
	db := pgtest.NewSqlxDB(t)

	keyStore := cltest.NewKeyStore(t, db)
	require.NoError(t, keyStore.OCR().Add(ctx, cltest.DefaultOCRKey))
	require.NoError(t, keyStore.P2P().Add(ctx, cltest.DefaultP2PKey))

	pipelineORM := pipeline.NewORM(db, logger.TestLogger(t), config.JobPipeline().MaxSuccessfulRuns())
	bridgesORM := bridges.NewORM(db)
	legacyChains := evmtest.NewLegacyChains(t, evmtest.TestChainOpts{
		ChainConfigs:   config.EVMConfigs(),
		DatabaseConfig: config.Database(),
		FeatureConfig:  config.Feature(),
		ListenerConfig: config.Database().Listener(),
		DB:             db,
		KeyStore:       keyStore.Eth(),
	})
	orm := NewTestORM(t, db, pipelineORM, bridgesORM, keyStore)

	_, bridge := cltest.MustCreateBridge(t, db, cltest.BridgeOpts{})
	_, bridge2 := cltest.MustCreateBridge(t, db, cltest.BridgeOpts{})

	externalJobID := uuid.New()
	_, address := cltest.MustInsertRandomKey(t, keyStore.Eth())
	jb, err := ocr.ValidatedOracleSpecToml(config, legacyChains,
		testspecs.GenerateOCRSpec(testspecs.OCRSpecParams{
			JobID:              externalJobID.String(),
			TransmitterAddress: address.Hex(),
			DS1BridgeName:      bridge.Name.String(),
			DS2BridgeName:      bridge2.Name.String(),
		}).Toml(),
	)
	require.NoError(t, err)

	err = orm.CreateJob(testutils.Context(t), &jb)
	require.NoError(t, err)

	t.Run("with no pipeline runs", func(t *testing.T) {
		ctx := testutils.Context(t)
		runs, count, err2 := orm.PipelineRuns(ctx, &jb.ID, 0, 10)
		require.NoError(t, err2)
		assert.Equal(t, 0, count)
		assert.Empty(t, runs)
	})

	t.Run("with a pipeline run", func(t *testing.T) {
		ctx := testutils.Context(t)
		run := mustInsertPipelineRun(t, pipelineORM, jb)

		runs, count, err2 := orm.PipelineRuns(ctx, &jb.ID, 0, 10)
		require.NoError(t, err2)

		assert.Equal(t, 1, count)
		actual := runs[0]

		// Test pipeline run fields
		assert.Equal(t, run.State, actual.State)
		assert.Equal(t, run.PipelineSpecID, actual.PipelineSpecID)

		// Test preloaded pipeline spec
		assert.Equal(t, jb.PipelineSpec.ID, actual.PipelineSpec.ID)
		assert.Equal(t, jb.ID, actual.PipelineSpec.JobID)
	})
}

func Test_FindPipelineRunIDsByJobID(t *testing.T) {
	ctx := testutils.Context(t)
	var jb job.Job

	config := configtest.NewTestGeneralConfig(t)
	_, db := heavyweight.FullTestDBV2(t, nil)

	keyStore := cltest.NewKeyStore(t, db)
	require.NoError(t, keyStore.OCR().Add(ctx, cltest.DefaultOCRKey))
	require.NoError(t, keyStore.P2P().Add(ctx, cltest.DefaultP2PKey))

	pipelineORM := pipeline.NewORM(db, logger.TestLogger(t), config.JobPipeline().MaxSuccessfulRuns())
	bridgesORM := bridges.NewORM(db)
	legacyChains := evmtest.NewLegacyChains(t, evmtest.TestChainOpts{
		ChainConfigs:   config.EVMConfigs(),
		DatabaseConfig: config.Database(),
		FeatureConfig:  config.Feature(),
		ListenerConfig: config.Database().Listener(),
		DB:             db,
		KeyStore:       keyStore.Eth(),
	})
	orm := NewTestORM(t, db, pipelineORM, bridgesORM, keyStore)

	_, address := cltest.MustInsertRandomKey(t, keyStore.Eth())

	jobs := make([]job.Job, 11)
	for j := 0; j < len(jobs); j++ {
		_, bridge := cltest.MustCreateBridge(t, db, cltest.BridgeOpts{})
		_, bridge2 := cltest.MustCreateBridge(t, db, cltest.BridgeOpts{})
		jobID := uuid.New().String()
		key, err := ethkey.NewV2()

		require.NoError(t, err)
		jb, err = ocr.ValidatedOracleSpecToml(config, legacyChains,
			testspecs.GenerateOCRSpec(testspecs.OCRSpecParams{
				JobID:              jobID,
				Name:               fmt.Sprintf("Job #%v", jobID),
				DS1BridgeName:      bridge.Name.String(),
				DS2BridgeName:      bridge2.Name.String(),
				TransmitterAddress: address.Hex(),
				ContractAddress:    key.Address.String(),
			}).Toml())

		require.NoError(t, err)

		err = orm.CreateJob(testutils.Context(t), &jb)
		require.NoError(t, err)
		jobs[j] = jb
	}

	for i, j := 0, 0; i < 2500; i++ {
		mustInsertPipelineRun(t, pipelineORM, jobs[j])
		j++
		if j == len(jobs)-1 {
			j = 0
		}
	}

	t.Run("with no pipeline runs", func(t *testing.T) {
		ctx := testutils.Context(t)
		runIDs, err := orm.FindPipelineRunIDsByJobID(ctx, jb.ID, 0, 10)
		require.NoError(t, err)
		assert.Empty(t, runIDs)
	})

	t.Run("with a pipeline run", func(t *testing.T) {
		ctx := testutils.Context(t)
		run := mustInsertPipelineRun(t, pipelineORM, jb)

		runIDs, err := orm.FindPipelineRunIDsByJobID(ctx, jb.ID, 0, 10)
		require.NoError(t, err)
		require.Len(t, runIDs, 1)

		assert.Equal(t, run.ID, runIDs[0])
	})

	// Internally these queries are batched by 1000, this tests case requiring concatenation
	//  of more than 1 batch
	t.Run("with batch concatenation limit 10", func(t *testing.T) {
		ctx := testutils.Context(t)
		runIDs, err := orm.FindPipelineRunIDsByJobID(ctx, jobs[3].ID, 95, 10)
		require.NoError(t, err)
		require.Len(t, runIDs, 10)
		assert.Equal(t, int64(4*(len(jobs)-1)), runIDs[3]-runIDs[7])
	})

	// Internally these queries are batched by 1000, this tests case requiring concatenation
	//  of more than 1 batch
	t.Run("with batch concatenation limit 100", func(t *testing.T) {
		ctx := testutils.Context(t)
		runIDs, err := orm.FindPipelineRunIDsByJobID(ctx, jobs[3].ID, 95, 100)
		require.NoError(t, err)
		require.Len(t, runIDs, 100)
		assert.Equal(t, int64(67*(len(jobs)-1)), runIDs[12]-runIDs[79])
	})

	for i := 0; i < 2100; i++ {
		mustInsertPipelineRun(t, pipelineORM, jb)
	}

	// There is a COUNT query which doesn't run unless the query for the most recent 1000 rows
	//  returns empty.  This can happen if the job id being requested hasn't run in a while,
	//  but many other jobs have run since.
	t.Run("with first batch empty, over limit", func(t *testing.T) {
		ctx := testutils.Context(t)
		runIDs, err := orm.FindPipelineRunIDsByJobID(ctx, jobs[3].ID, 0, 25)
		require.NoError(t, err)
		require.Len(t, runIDs, 25)
		assert.Equal(t, int64(16*(len(jobs)-1)), runIDs[7]-runIDs[23])
	})

	// Same as previous, but where there are fewer matching jobs than the limit
	t.Run("with first batch empty, under limit", func(t *testing.T) {
		ctx := testutils.Context(t)
		runIDs, err := orm.FindPipelineRunIDsByJobID(ctx, jobs[3].ID, 143, 190)
		require.NoError(t, err)
		require.Len(t, runIDs, 107)
		assert.Equal(t, int64(16*(len(jobs)-1)), runIDs[7]-runIDs[23])
	})
}

func Test_FindPipelineRunsByIDs(t *testing.T) {
	t.Parallel()
	ctx := testutils.Context(t)

	config := configtest.NewTestGeneralConfig(t)
	db := pgtest.NewSqlxDB(t)

	keyStore := cltest.NewKeyStore(t, db)
	require.NoError(t, keyStore.OCR().Add(ctx, cltest.DefaultOCRKey))
	require.NoError(t, keyStore.P2P().Add(ctx, cltest.DefaultP2PKey))

	pipelineORM := pipeline.NewORM(db, logger.TestLogger(t), config.JobPipeline().MaxSuccessfulRuns())
	bridgesORM := bridges.NewORM(db)
	legacyChains := evmtest.NewLegacyChains(t, evmtest.TestChainOpts{
		ChainConfigs:   config.EVMConfigs(),
		DatabaseConfig: config.Database(),
		FeatureConfig:  config.Feature(),
		ListenerConfig: config.Database().Listener(),
		DB:             db,
		KeyStore:       keyStore.Eth(),
	})
	orm := NewTestORM(t, db, pipelineORM, bridgesORM, keyStore)

	_, bridge := cltest.MustCreateBridge(t, db, cltest.BridgeOpts{})
	_, bridge2 := cltest.MustCreateBridge(t, db, cltest.BridgeOpts{})

	externalJobID := uuid.New()
	_, address := cltest.MustInsertRandomKey(t, keyStore.Eth())
	jb, err := ocr.ValidatedOracleSpecToml(config, legacyChains,
		testspecs.GenerateOCRSpec(testspecs.OCRSpecParams{
			JobID:              externalJobID.String(),
			TransmitterAddress: address.Hex(),
			DS1BridgeName:      bridge.Name.String(),
			DS2BridgeName:      bridge2.Name.String(),
		}).Toml(),
	)
	require.NoError(t, err)

	err = orm.CreateJob(testutils.Context(t), &jb)
	require.NoError(t, err)

	t.Run("with no pipeline runs", func(t *testing.T) {
		ctx := testutils.Context(t)
		runs, err2 := orm.FindPipelineRunsByIDs(ctx, []int64{-1})
		require.NoError(t, err2)
		assert.Empty(t, runs)
	})

	t.Run("with a pipeline run", func(t *testing.T) {
		ctx := testutils.Context(t)
		run := mustInsertPipelineRun(t, pipelineORM, jb)

		actual, err2 := orm.FindPipelineRunsByIDs(ctx, []int64{run.ID})
		require.NoError(t, err2)
		require.Len(t, actual, 1)

		actualRun := actual[0]
		// Test pipeline run fields
		assert.Equal(t, run.State, actualRun.State)
		assert.Equal(t, run.PipelineSpecID, actualRun.PipelineSpecID)

		// Test preloaded pipeline spec
		assert.Equal(t, jb.PipelineSpec.ID, actualRun.PipelineSpec.ID)
		assert.Equal(t, jb.ID, actualRun.PipelineSpec.JobID)
	})
}

func Test_FindPipelineRunByID(t *testing.T) {
	t.Parallel()
	ctx := testutils.Context(t)

	config := configtest.NewTestGeneralConfig(t)
	db := pgtest.NewSqlxDB(t)

	keyStore := cltest.NewKeyStore(t, db)
	err := keyStore.OCR().Add(ctx, cltest.DefaultOCRKey)
	require.NoError(t, err)

	pipelineORM := pipeline.NewORM(db, logger.TestLogger(t), config.JobPipeline().MaxSuccessfulRuns())
	bridgesORM := bridges.NewORM(db)
	orm := NewTestORM(t, db, pipelineORM, bridgesORM, keyStore)

	jb, err := directrequest.ValidatedDirectRequestSpec(testspecs.GetDirectRequestSpec())
	require.NoError(t, err)

	err = orm.CreateJob(testutils.Context(t), &jb)
	require.NoError(t, err)

	t.Run("with no pipeline run", func(t *testing.T) {
		ctx := testutils.Context(t)
		run, err2 := orm.FindPipelineRunByID(ctx, -1)
		assert.Equal(t, pipeline.Run{}, run)
		require.ErrorIs(t, err2, sql.ErrNoRows)
	})

	t.Run("with a pipeline run", func(t *testing.T) {
		ctx := testutils.Context(t)
		run := mustInsertPipelineRun(t, pipelineORM, jb)

		actual, err2 := orm.FindPipelineRunByID(ctx, run.ID)
		require.NoError(t, err2)

		actualRun := actual
		// Test pipeline run fields
		assert.Equal(t, run.State, actualRun.State)
		assert.Equal(t, run.PipelineSpecID, actualRun.PipelineSpecID)

		// Test preloaded pipeline spec
		assert.Equal(t, jb.PipelineSpec.ID, actualRun.PipelineSpec.ID)
		assert.Equal(t, jb.ID, actualRun.PipelineSpec.JobID)
	})
}

func Test_FindJobWithoutSpecErrors(t *testing.T) {
	t.Parallel()
	ctx := testutils.Context(t)

	config := configtest.NewTestGeneralConfig(t)
	db := pgtest.NewSqlxDB(t)

	keyStore := cltest.NewKeyStore(t, db)
	err := keyStore.OCR().Add(ctx, cltest.DefaultOCRKey)
	require.NoError(t, err)

	pipelineORM := pipeline.NewORM(db, logger.TestLogger(t), config.JobPipeline().MaxSuccessfulRuns())
	bridgesORM := bridges.NewORM(db)
	orm := NewTestORM(t, db, pipelineORM, bridgesORM, keyStore)

	jb, err := directrequest.ValidatedDirectRequestSpec(testspecs.GetDirectRequestSpec())
	require.NoError(t, err)

	err = orm.CreateJob(ctx, &jb)
	require.NoError(t, err)
	var jobSpec job.Job
	err = db.Get(&jobSpec, "SELECT * FROM jobs")
	require.NoError(t, err)

	ocrSpecError1 := "ocr spec 1 errored"
	ocrSpecError2 := "ocr spec 2 errored"
	require.NoError(t, orm.RecordError(ctx, jobSpec.ID, ocrSpecError1))
	require.NoError(t, orm.RecordError(ctx, jobSpec.ID, ocrSpecError2))

	jb, err = orm.FindJobWithoutSpecErrors(ctx, jobSpec.ID)
	require.NoError(t, err)
	jbWithErrors, err := orm.FindJob(testutils.Context(t), jobSpec.ID)
	require.NoError(t, err)

	assert.Empty(t, jb.JobSpecErrors)
	assert.Len(t, jbWithErrors.JobSpecErrors, 2)
}

func Test_FindSpecErrorsByJobIDs(t *testing.T) {
	t.Parallel()
	ctx := testutils.Context(t)

	config := configtest.NewTestGeneralConfig(t)
	db := pgtest.NewSqlxDB(t)

	keyStore := cltest.NewKeyStore(t, db)
	err := keyStore.OCR().Add(ctx, cltest.DefaultOCRKey)
	require.NoError(t, err)

	pipelineORM := pipeline.NewORM(db, logger.TestLogger(t), config.JobPipeline().MaxSuccessfulRuns())
	bridgesORM := bridges.NewORM(db)
	orm := NewTestORM(t, db, pipelineORM, bridgesORM, keyStore)

	jb, err := directrequest.ValidatedDirectRequestSpec(testspecs.GetDirectRequestSpec())
	require.NoError(t, err)

	err = orm.CreateJob(ctx, &jb)
	require.NoError(t, err)
	var jobSpec job.Job
	err = db.Get(&jobSpec, "SELECT * FROM jobs")
	require.NoError(t, err)

	ocrSpecError1 := "ocr spec 1 errored"
	ocrSpecError2 := "ocr spec 2 errored"
	require.NoError(t, orm.RecordError(ctx, jobSpec.ID, ocrSpecError1))
	require.NoError(t, orm.RecordError(ctx, jobSpec.ID, ocrSpecError2))

	specErrs, err := orm.FindSpecErrorsByJobIDs(ctx, []int32{jobSpec.ID})
	require.NoError(t, err)

	assert.Len(t, specErrs, 2)
}

func Test_CountPipelineRunsByJobID(t *testing.T) {
	t.Parallel()
	ctx := testutils.Context(t)

	config := configtest.NewTestGeneralConfig(t)
	db := pgtest.NewSqlxDB(t)

	keyStore := cltest.NewKeyStore(t, db)
	require.NoError(t, keyStore.OCR().Add(ctx, cltest.DefaultOCRKey))
	require.NoError(t, keyStore.P2P().Add(ctx, cltest.DefaultP2PKey))

	pipelineORM := pipeline.NewORM(db, logger.TestLogger(t), config.JobPipeline().MaxSuccessfulRuns())
	bridgesORM := bridges.NewORM(db)
	legacyChains := evmtest.NewLegacyChains(t, evmtest.TestChainOpts{
		ChainConfigs:   config.EVMConfigs(),
		DatabaseConfig: config.Database(),
		FeatureConfig:  config.Feature(),
		ListenerConfig: config.Database().Listener(),
		DB:             db,
		KeyStore:       keyStore.Eth(),
	})
	orm := NewTestORM(t, db, pipelineORM, bridgesORM, keyStore)

	_, bridge := cltest.MustCreateBridge(t, db, cltest.BridgeOpts{})
	_, bridge2 := cltest.MustCreateBridge(t, db, cltest.BridgeOpts{})

	externalJobID := uuid.New()
	_, address := cltest.MustInsertRandomKey(t, keyStore.Eth())
	jb, err := ocr.ValidatedOracleSpecToml(config, legacyChains,
		testspecs.GenerateOCRSpec(testspecs.OCRSpecParams{
			JobID:              externalJobID.String(),
			TransmitterAddress: address.Hex(),
			DS1BridgeName:      bridge.Name.String(),
			DS2BridgeName:      bridge2.Name.String(),
		}).Toml(),
	)
	require.NoError(t, err)

	err = orm.CreateJob(testutils.Context(t), &jb)
	require.NoError(t, err)

	t.Run("with no pipeline runs", func(t *testing.T) {
		ctx := testutils.Context(t)
		count, err2 := orm.CountPipelineRunsByJobID(ctx, jb.ID)
		require.NoError(t, err2)
		assert.Equal(t, int32(0), count)
	})

	t.Run("with a pipeline run", func(t *testing.T) {
		ctx := testutils.Context(t)
		mustInsertPipelineRun(t, pipelineORM, jb)

		count, err2 := orm.CountPipelineRunsByJobID(ctx, jb.ID)
		require.NoError(t, err2)
		require.Equal(t, int32(1), count)
	})
}

func Test_ORM_FindJobByWorkflow(t *testing.T) {
	var addr1 = "0x0123456789012345678901234567890123456789"
	var addr2 = "0xabcdefabcdefabcdefabcdefabcdefabcdefabcd"
	t.Parallel()
	type fields struct {
		ds sqlutil.DataSource
	}
	type args struct {
		spec   *job.WorkflowSpec
		before func(t *testing.T, o job.ORM, s *job.WorkflowSpec) int32
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{

		{
			name: "wf not job found",
			fields: fields{
				ds: pgtest.NewSqlxDB(t),
			},
			args: args{
				// before is nil, so no job is inserted
				spec: &job.WorkflowSpec{
					ID:       1,
					Workflow: pkgworkflows.WFYamlSpec(t, "workflow00", addr1),
				},
			},
			wantErr: true,
		},

		{
			name: "wf job found",
			fields: fields{
				ds: pgtest.NewSqlxDB(t),
			},
			args: args{
				spec: &job.WorkflowSpec{
					ID:       1,
					Workflow: pkgworkflows.WFYamlSpec(t, "workflow01", addr1),
					SpecType: job.YamlSpec,
				},
				before: mustInsertWFJob,
			},
			wantErr: false,
		},

		{
			name: "wf wrong name",
			fields: fields{
				ds: pgtest.NewSqlxDB(t),
			},
			args: args{
				spec: &job.WorkflowSpec{
					ID:       1,
					Workflow: pkgworkflows.WFYamlSpec(t, "workflow02", addr1),
				},
				before: func(t *testing.T, o job.ORM, s *job.WorkflowSpec) int32 {
					var c job.WorkflowSpec
					c.ID = s.ID
					c.Workflow = pkgworkflows.WFYamlSpec(t, "workflow99", addr1) // insert with mismatched name
					c.SpecType = job.YamlSpec
					c.SecretsID = s.SecretsID
					return mustInsertWFJob(t, o, &c)
				},
			},
			wantErr: true,
		},
		{
			name: "wf wrong owner",
			fields: fields{
				ds: pgtest.NewSqlxDB(t),
			},
			args: args{
				spec: &job.WorkflowSpec{
					ID:       1,
					Workflow: pkgworkflows.WFYamlSpec(t, "workflow03", addr1),
				},
				before: func(t *testing.T, o job.ORM, s *job.WorkflowSpec) int32 {
					var c job.WorkflowSpec
					c.ID = s.ID
					c.Workflow = pkgworkflows.WFYamlSpec(t, "workflow03", addr2) // insert with mismatched owner
					c.SecretsID = s.SecretsID
					return mustInsertWFJob(t, o, &c)
				},
			},
			wantErr: true,
		},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := testutils.Context(t)
			ks := cltest.NewKeyStore(t, tt.fields.ds)

			secretsORM := artifacts.NewWorkflowRegistryDS(tt.fields.ds, logger.TestLogger(t))

			sid, err := secretsORM.Create(ctx, "some-url.com", fmt.Sprintf("some-hash-%d", i), "some-contentz")
			require.NoError(t, err)
			tt.args.spec.SecretsID = sql.NullInt64{Int64: sid, Valid: true}

			pipelineORM := pipeline.NewORM(tt.fields.ds, logger.TestLogger(t), configtest.NewTestGeneralConfig(t).JobPipeline().MaxSuccessfulRuns())
			bridgesORM := bridges.NewORM(tt.fields.ds)
			o := NewTestORM(t, tt.fields.ds, pipelineORM, bridgesORM, ks)

			var wantJobID int32
			if tt.args.before != nil {
				wantJobID = tt.args.before(t, o, tt.args.spec)
			}

			gotJ, err := o.FindJobIDByWorkflow(ctx, *tt.args.spec)
			if (err != nil) != tt.wantErr {
				t.Errorf("orm.FindJobByWorkflow() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err == nil {
				assert.Equal(t, wantJobID, gotJ, "mismatch job id")
			}
		})
	}
}

func Test_ORM_FindJobByWorkflow_Multiple(t *testing.T) {
	var addr1 = "0x012345678901234567890123456789012345ffff"
	var addr2 = "0xabcdefabcdefabcdefabcdefabcdefabcdef0000"
	t.Parallel()
	t.Run("multiple jobs", func(t *testing.T) {
		db := pgtest.NewSqlxDB(t)
		o := NewTestORM(t,
			db,
			pipeline.NewORM(db,
				logger.TestLogger(t),
				configtest.NewTestGeneralConfig(t).JobPipeline().MaxSuccessfulRuns()),
			bridges.NewORM(db),
			cltest.NewKeyStore(t, db))
		ctx := testutils.Context(t)
		secretsORM := artifacts.NewWorkflowRegistryDS(db, logger.TestLogger(t))

		var sids []int64
		for i := 0; i < 3; i++ {
			sid, err := secretsORM.Create(ctx, "some-url.com", fmt.Sprintf("some-hash-%d", i), "some-contentz")
			require.NoError(t, err)
			sids = append(sids, sid)
		}

		wfYaml1 := pkgworkflows.WFYamlSpec(t, "workflow00", addr1)
		s1 := job.WorkflowSpec{
			Workflow:  wfYaml1,
			SpecType:  job.YamlSpec,
			SecretsID: sql.NullInt64{Int64: sids[0], Valid: true},
		}
		wantJobID1 := mustInsertWFJob(t, o, &s1)

		wfYaml2 := pkgworkflows.WFYamlSpec(t, "workflow01", addr1)
		s2 := job.WorkflowSpec{
			Workflow:  wfYaml2,
			SpecType:  job.YamlSpec,
			SecretsID: sql.NullInt64{Int64: sids[1], Valid: true},
		}
		wantJobID2 := mustInsertWFJob(t, o, &s2)

		wfYaml3 := pkgworkflows.WFYamlSpec(t, "workflow00", addr2)
		s3 := job.WorkflowSpec{
			Workflow:  wfYaml3,
			SpecType:  job.YamlSpec,
			SecretsID: sql.NullInt64{Int64: sids[2], Valid: true},
		}
		wantJobID3 := mustInsertWFJob(t, o, &s3)

		expectedIDs := []int32{wantJobID1, wantJobID2, wantJobID3}
		for i, s := range []job.WorkflowSpec{s1, s2, s3} {
			gotJ, err := o.FindJobIDByWorkflow(ctx, s)
			require.NoError(t, err)
			assert.Equal(t, expectedIDs[i], gotJ, "mismatch job id case %d, spec %v", i, s)
			j, err := o.FindJob(ctx, expectedIDs[i])
			require.NoError(t, err)
			assert.NotNil(t, j)
			t.Logf("found job %v", j)
			assert.EqualValues(t, j.WorkflowSpec.Workflow, s.Workflow)
			assert.EqualValues(t, j.WorkflowSpec.WorkflowID, s.WorkflowID)
			assert.EqualValues(t, j.WorkflowSpec.WorkflowOwner, s.WorkflowOwner)
			assert.EqualValues(t, j.WorkflowSpec.WorkflowName, s.WorkflowName)
			assert.Equal(t, job.YamlSpec, j.WorkflowSpec.SpecType)
		}
	})
}

func mustInsertWFJob(t *testing.T, orm job.ORM, s *job.WorkflowSpec) int32 {
	t.Helper()
	err := s.Validate(testutils.Context(t))
	require.NoError(t, err, "failed to validate spec %v", s)
	ctx := testutils.Context(t)
	_, err = toml.Marshal(s.Workflow)
	require.NoError(t, err, "failed to TOML marshal workflow %v", s.Workflow)
	j := job.Job{
		Type:          job.Workflow,
		WorkflowSpec:  s,
		ExternalJobID: uuid.New(),
		Name:          null.StringFrom(s.WorkflowOwner + "_" + s.WorkflowName),
		SchemaVersion: 1,
	}

	err = orm.CreateJob(ctx, &j)
	require.NoError(t, err, "failed to insert job with wf spec %+v %s", s, err)
	return j.ID
}

func mustInsertPipelineRun(t *testing.T, orm pipeline.ORM, j job.Job) pipeline.Run {
	t.Helper()
	ctx := testutils.Context(t)

	run := pipeline.Run{
		PipelineSpecID: j.PipelineSpecID,
		PruningKey:     j.ID,
		State:          pipeline.RunStatusRunning,
		Outputs:        jsonserializable.JSONSerializable{Valid: false},
		AllErrors:      pipeline.RunErrors{},
		CreatedAt:      time.Now(),
		FinishedAt:     null.Time{},
	}
	err := orm.CreateRun(ctx, &run)
	require.NoError(t, err)
	return run
}

func TestORM_CreateJob_OCR2_With_DualTransmission(t *testing.T) {
	ctx := testutils.Context(t)
	customChainID := big.New(testutils.NewRandomEVMChainID())

	config := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		enabled := true
		c.EVM = append(c.EVM, &configtoml.EVMConfig{
			ChainID: customChainID,
			Chain:   configtoml.Defaults(customChainID),
			Enabled: &enabled,
			Nodes:   configtoml.EVMNodes{{}},
		})
	})
	db := pgtest.NewSqlxDB(t)
	keyStore := cltest.NewKeyStore(t, db)
	require.NoError(t, keyStore.OCR2().Add(ctx, cltest.DefaultOCR2Key))
	_, transmitterID := cltest.MustInsertRandomKey(t, keyStore.Eth())

	baseJobSpec := fmt.Sprintf(testspecs.OCR2EVMDualTransmissionSpecMinimalTemplate, transmitterID.String())

	lggr := logger.TestLogger(t)
	pipelineORM := pipeline.NewORM(db, lggr, config.JobPipeline().MaxSuccessfulRuns())
	bridgesORM := bridges.NewORM(db)

	jobORM := NewTestORM(t, db, pipelineORM, bridgesORM, keyStore)

	// Enabled but no config set
	enabledDualTransmissionSpec := `
		enableDualTransmission=true`

	jb, err := ocr2validate.ValidatedOracleSpecToml(testutils.Context(t), config.OCR2(), config.Insecure(), baseJobSpec+enabledDualTransmissionSpec, nil)
	require.NoError(t, err)
	require.ErrorContains(t, jobORM.CreateJob(ctx, &jb), "dual transmission is enabled but no dual transmission config present")

	// ContractAddress not set
	emptyContractAddress := `
		enableDualTransmission=true
		[relayConfig.dualTransmission]
		contractAddress=""
	`
	jb, err = ocr2validate.ValidatedOracleSpecToml(testutils.Context(t), config.OCR2(), config.Insecure(), baseJobSpec+emptyContractAddress, nil)
	require.NoError(t, err)
	require.ErrorContains(t, jobORM.CreateJob(ctx, &jb), "invalid contract address in dual transmission config")

	// Transmitter address not set
	emptyTransmitterAddress := `
		enableDualTransmission=true
		[relayConfig.dualTransmission]
		contractAddress = '0x613a38AC1659769640aaE063C651F48E0250454C'
		transmitterAddress = ''
	`
	jb, err = ocr2validate.ValidatedOracleSpecToml(testutils.Context(t), config.OCR2(), config.Insecure(), baseJobSpec+emptyTransmitterAddress, nil)
	require.NoError(t, err)
	require.ErrorContains(t, jobORM.CreateJob(ctx, &jb), "invalid transmitter address in dual transmission config")

	dtTransmitterAddress := cltest.MustGenerateRandomKey(t)

	completeDualTransmissionSpec := fmt.Sprintf(`
		enableDualTransmission=true
		[relayConfig.dualTransmission]
		contractAddress = '0x613a38AC1659769640aaE063C651F48E0250454C'
		transmitterAddress = '%s'
		[relayConfig.dualTransmission.meta]
		key1 = ['val1']
		key2 = ['val2','val3']
		`,
		dtTransmitterAddress.Address.String())

	jb, err = ocr2validate.ValidatedOracleSpecToml(testutils.Context(t), config.OCR2(), config.Insecure(), baseJobSpec+completeDualTransmissionSpec, nil)
	require.NoError(t, err)

	jb.OCR2OracleSpec.TransmitterID = null.StringFrom(transmitterID.String())

	// Unknown transmitter address
	require.ErrorContains(t, jobORM.CreateJob(ctx, &jb), "unknown dual transmission transmitterAddress: no EVM key matching:")

	// Should not error
	keyStore.Eth().XXXTestingOnlyAdd(ctx, dtTransmitterAddress)
	require.NoError(t, jobORM.CreateJob(ctx, &jb))
}

func TestORM_CreateJob_KeyLocking(t *testing.T) {
	ctx := testutils.Context(t)
	customChainID := big.New(testutils.NewRandomEVMChainID())

	config := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		enabled := true
		c.EVM = append(c.EVM, &configtoml.EVMConfig{
			ChainID: customChainID,
			Chain:   configtoml.Defaults(customChainID),
			Enabled: &enabled,
			Nodes:   configtoml.EVMNodes{{}},
		})
	})
	db := pgtest.NewSqlxDB(t)
	ks := cltest.NewKeyStore(t, db)
	require.NoError(t, ks.OCR2().Add(ctx, cltest.DefaultOCR2Key))
	_, transmitterID := cltest.MustInsertRandomKey(t, ks.Eth())
	dtTransmitterAddress := cltest.MustGenerateRandomKey(t)
	ks.Eth().XXXTestingOnlyAdd(ctx, dtTransmitterAddress)

	baseJobSpec := fmt.Sprintf(testspecs.OCR2EVMDualTransmissionSpecMinimalTemplate, transmitterID.String())

	lggr := logger.TestLogger(t)
	pipelineORM := pipeline.NewORM(db, lggr, config.JobPipeline().MaxSuccessfulRuns())
	bridgesORM := bridges.NewORM(db)

	jobORM := NewTestORM(t, db, pipelineORM, bridgesORM, ks)

	t.Run("keys not locked", func(t *testing.T) {
		completeDualTransmissionSpec := fmt.Sprintf(`
		enableDualTransmission=true
		[relayConfig.dualTransmission]
		contractAddress = '0x613a38AC1659769640aaE063C651F48E0250454C'
		transmitterAddress = '%s'
		[relayConfig.dualTransmission.meta]
		key1 = ['val1']
		key2 = ['val2','val3']
		`,
			dtTransmitterAddress.Address.String())

		jb, err := ocr2validate.ValidatedOracleSpecToml(testutils.Context(t), config.OCR2(), config.Insecure(), baseJobSpec+completeDualTransmissionSpec, nil)
		require.NoError(t, err)

		jb.OCR2OracleSpec.TransmitterID = null.StringFrom(transmitterID.String())
		jb.Name = null.StringFrom(uuid.NewString())

		require.NoError(t, jobORM.CreateJob(ctx, &jb))
	})

	t.Run("keys locked", func(t *testing.T) {
		completeDualTransmissionSpec := fmt.Sprintf(`
		enableDualTransmission=true
		[relayConfig.dualTransmission]
		contractAddress = '0x613a38AC1659769640aaE063C651F48E0250454C'
		transmitterAddress = '%s'
		[relayConfig.dualTransmission.meta]
		key1 = ['val1']
		key2 = ['val2','val3']
		`,
			dtTransmitterAddress.Address.String())

		rm := ks.Eth().GetResourceMutex(ctx, transmitterID)
		require.NoError(t, rm.TryLock(keys.TXMv1))
		rm = ks.Eth().GetResourceMutex(ctx, dtTransmitterAddress.Address)
		require.NoError(t, rm.TryLock(keys.TXMv2))
		jb, err := ocr2validate.ValidatedOracleSpecToml(testutils.Context(t), config.OCR2(), config.Insecure(), baseJobSpec+completeDualTransmissionSpec, nil)
		require.NoError(t, err)

		jb.OCR2OracleSpec.TransmitterID = null.StringFrom(transmitterID.String())
		jb.Name = null.StringFrom(uuid.NewString())

		require.NoError(t, jobORM.CreateJob(ctx, &jb))
	})

	t.Run("keys locked but job spec misconfigured", func(t *testing.T) {
		rm := ks.Eth().GetResourceMutex(ctx, transmitterID)
		require.NoError(t, rm.TryLock(keys.TXMv1))
		rm = ks.Eth().GetResourceMutex(ctx, dtTransmitterAddress.Address)
		require.NoError(t, rm.TryLock(keys.TXMv2))

		completeDualTransmissionSpec := fmt.Sprintf(`
		enableDualTransmission=true
		[relayConfig.dualTransmission]
		contractAddress = '0x613a38AC1659769640aaE063C651F48E0250454C'
		transmitterAddress = '%s'
		[relayConfig.dualTransmission.meta]
		key1 = ['val1']
		key2 = ['val2','val3']
		`,
			transmitterID.String())

		jb, err := ocr2validate.ValidatedOracleSpecToml(testutils.Context(t), config.OCR2(), config.Insecure(), baseJobSpec+completeDualTransmissionSpec, nil)
		require.NoError(t, err)

		jb.OCR2OracleSpec.TransmitterID = null.StringFrom(dtTransmitterAddress.Address.String())
		jb.Name = null.StringFrom(uuid.NewString())

		require.ErrorContains(t, jobORM.CreateJob(ctx, &jb), "cannot be a secondary transmitter address because it's used a primary transmitter in another job")
	})
}

func Test_FindGatewayJobID(t *testing.T) {
	t.Parallel()
	ctx := testutils.Context(t)

	config := configtest.NewTestGeneralConfig(t)
	db := pgtest.NewSqlxDB(t)

	keyStore := cltest.NewKeyStore(t, db)
	err := keyStore.OCR().Add(ctx, cltest.DefaultOCRKey)
	require.NoError(t, err, "failed to add OCR key")

	pipelineORM := pipeline.NewORM(db, logger.TestLogger(t), config.JobPipeline().MaxSuccessfulRuns())
	bridgesORM := bridges.NewORM(db)
	orm := NewTestORM(t, db, pipelineORM, bridgesORM, keyStore)

	gatewayJob, err := gateway.ValidatedGatewaySpec(testspecs.GetGatewaySpec())
	require.NoError(t, err, "failed to validate gateway spec")

	err = orm.CreateJob(ctx, &gatewayJob)
	require.NoError(t, err, "failed to create gateway job")
	var jobSpec job.Job
	err = db.Get(&jobSpec, "SELECT * FROM jobs")
	require.NoError(t, err, "failed to get gateway job from db")

	// find only by auth gateway id
	gatewayJobSpec := job.GatewaySpec{
		GatewayConfig: job.JSONConfig{
			"ConnectionManagerConfig": map[string]interface{}{
				"AuthGatewayId": "gateway",
			},
		},
	}
	id, err := orm.FindGatewayJobID(ctx, gatewayJobSpec)
	require.NoError(t, err, "failed to find gateway job by auth gateway id")
	require.Equal(t, jobSpec.ID, id, "mismatch job id")
}

func Test_FindGatewayJobID_NoMatch(t *testing.T) {
	t.Parallel()
	ctx := testutils.Context(t)

	config := configtest.NewTestGeneralConfig(t)
	db := pgtest.NewSqlxDB(t)

	keyStore := cltest.NewKeyStore(t, db)
	err := keyStore.OCR().Add(ctx, cltest.DefaultOCRKey)
	require.NoError(t, err, "failed to add OCR key")

	pipelineORM := pipeline.NewORM(db, logger.TestLogger(t), config.JobPipeline().MaxSuccessfulRuns())
	bridgesORM := bridges.NewORM(db)
	orm := NewTestORM(t, db, pipelineORM, bridgesORM, keyStore)

	gatewayJob, err := gateway.ValidatedGatewaySpec(testspecs.GetGatewaySpec())
	require.NoError(t, err, "failed to validate gateway spec")

	err = orm.CreateJob(ctx, &gatewayJob)
	require.NoError(t, err, "failed to create gateway job")
	var jobSpec job.Job
	err = db.Get(&jobSpec, "SELECT * FROM jobs")
	require.NoError(t, err, "failed to get gateway job from db")

	// different auth gateway id
	gatewayJobSpec := job.GatewaySpec{
		GatewayConfig: job.JSONConfig{
			"ConnectionManagerConfig": map[string]interface{}{
				"AuthGatewayId": "another_gateway",
			},
		},
	}
	id, err := orm.FindGatewayJobID(ctx, gatewayJobSpec)
	require.Error(t, err, "found gateway job by auth gateway id")
	require.Equal(t, int32(0), id, "found non-zero job id")
}

func Test_FindStandardCapabilityJobID(t *testing.T) {
	t.Parallel()
	ctx := testutils.Context(t)

	config := configtest.NewTestGeneralConfig(t)
	db := pgtest.NewSqlxDB(t)

	keyStore := cltest.NewKeyStore(t, db)
	err := keyStore.OCR().Add(ctx, cltest.DefaultOCRKey)
	require.NoError(t, err, "failed to add OCR key")

	pipelineORM := pipeline.NewORM(db, logger.TestLogger(t), config.JobPipeline().MaxSuccessfulRuns())
	bridgesORM := bridges.NewORM(db)
	orm := NewTestORM(t, db, pipelineORM, bridgesORM, keyStore)

	stdJob, err := standardcapabilities.ValidatedStandardCapabilitiesSpec(testspecs.GetStandardCapabilitySpec())
	require.NoError(t, err, "failed to validate standard capabilities spec")

	err = orm.CreateJob(ctx, &stdJob)
	require.NoError(t, err, "failed to create standard capabilities job")
	var jobSpec job.Job
	err = db.Get(&jobSpec, "SELECT * FROM jobs")
	require.NoError(t, err, "failed to get standard capabilities job from db")

	stdCapJobSpec := job.StandardCapabilitiesSpec{
		Command: "/home/capabilities/some_capability_linux_amd64",
	}

	id, err := orm.FindStandardCapabilityJobID(ctx, stdCapJobSpec)
	require.NoError(t, err, "failed to find standard capabilities by command")
	require.Equal(t, jobSpec.ID, id, "mismatch job id")
}

func Test_FindStandardCapabilityJobID_NoMatch(t *testing.T) {
	t.Parallel()
	ctx := testutils.Context(t)

	config := configtest.NewTestGeneralConfig(t)
	db := pgtest.NewSqlxDB(t)

	keyStore := cltest.NewKeyStore(t, db)
	err := keyStore.OCR().Add(ctx, cltest.DefaultOCRKey)
	require.NoError(t, err, "failed to add OCR key")

	pipelineORM := pipeline.NewORM(db, logger.TestLogger(t), config.JobPipeline().MaxSuccessfulRuns())
	bridgesORM := bridges.NewORM(db)
	orm := NewTestORM(t, db, pipelineORM, bridgesORM, keyStore)

	stdJob, err := standardcapabilities.ValidatedStandardCapabilitiesSpec(testspecs.GetStandardCapabilitySpec())
	require.NoError(t, err, "failed to validate standard capabilities spec")

	err = orm.CreateJob(ctx, &stdJob)
	require.NoError(t, err, "failed to create standard capabilities job")
	var jobSpec job.Job
	err = db.Get(&jobSpec, "SELECT * FROM jobs")
	require.NoError(t, err, "failed to get standard capabilities job from db")

	stdCapJobSpec := job.StandardCapabilitiesSpec{
		Command: "/home/capabilities/some_other_capability_linux_amd64",
	}

	id, err := orm.FindStandardCapabilityJobID(ctx, stdCapJobSpec)
	require.Error(t, err, "found standard capabilities with different command")
	require.Equal(t, int32(0), id, "found non-zero job id")
}
