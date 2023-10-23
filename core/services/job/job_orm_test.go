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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink-relay/pkg/types"

	"github.com/smartcontractkit/chainlink/v2/core/services/relay"

	"github.com/smartcontractkit/chainlink/v2/core/assets"
	"github.com/smartcontractkit/chainlink/v2/core/bridges"
	evmcfg "github.com/smartcontractkit/chainlink/v2/core/chains/evm/config/toml"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	configtest "github.com/smartcontractkit/chainlink/v2/core/internal/testutils/configtest/v2"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/evmtest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/blockhashstore"
	"github.com/smartcontractkit/chainlink/v2/core/services/blockheaderfeeder"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/services/directrequest"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/keeper"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr"
	ocr2validate "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/validate"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocrbootstrap"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
	evmrelay "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm"
	"github.com/smartcontractkit/chainlink/v2/core/services/vrf/vrfcommon"
	"github.com/smartcontractkit/chainlink/v2/core/services/webhook"
	"github.com/smartcontractkit/chainlink/v2/core/testdata/testspecs"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
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

[pluginConfig]
serverURL = 'wss://localhost:8080'
serverPubKey = '8fa807463ad73f9ee855cfd60ba406dcf98a2855b3dd8af613107b0f6890a707'
`

func TestORM(t *testing.T) {
	t.Parallel()
	config := configtest.NewTestGeneralConfig(t)
	db := pgtest.NewSqlxDB(t)
	keyStore := cltest.NewKeyStore(t, db, config.Database())
	ethKeyStore := keyStore.Eth()

	require.NoError(t, keyStore.OCR().Add(cltest.DefaultOCRKey))
	require.NoError(t, keyStore.P2P().Add(cltest.DefaultP2PKey))

	pipelineORM := pipeline.NewORM(db, logger.TestLogger(t), config.Database(), config.JobPipeline().MaxSuccessfulRuns())
	bridgesORM := bridges.NewORM(db, logger.TestLogger(t), config.Database())

	relayExtenders := evmtest.NewChainRelayExtenders(t, evmtest.TestChainOpts{DB: db, GeneralConfig: config, KeyStore: ethKeyStore})
	legacyChains := evmrelay.NewLegacyChainsFromRelayerExtenders(relayExtenders)
	orm := NewTestORM(t, db, legacyChains, pipelineORM, bridgesORM, keyStore, config.Database())
	borm := bridges.NewORM(db, logger.TestLogger(t), config.Database())

	_, bridge := cltest.MustCreateBridge(t, db, cltest.BridgeOpts{}, config.Database())
	_, bridge2 := cltest.MustCreateBridge(t, db, cltest.BridgeOpts{}, config.Database())
	_, address := cltest.MustInsertRandomKey(t, ethKeyStore)
	jb := makeOCRJobSpec(t, address, bridge.Name.String(), bridge2.Name.String())

	t.Run("it creates job specs", func(t *testing.T) {
		err := orm.CreateJob(jb)
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

	t.Run("autogenerates external job ID if missing", func(t *testing.T) {
		jb2 := makeOCRJobSpec(t, address, bridge.Name.String(), bridge2.Name.String())
		jb2.ExternalJobID = uuid.UUID{}
		err := orm.CreateJob(jb2)
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
		require.Len(t, dbSpecs, 2)

		err = orm.DeleteJob(jb.ID)
		require.NoError(t, err)

		dbSpecs = []job.Job{}
		err = db.Select(&dbSpecs, "SELECT * FROM jobs")
		require.NoError(t, err)
		require.Len(t, dbSpecs, 1)
	})

	t.Run("increase job spec error occurrence", func(t *testing.T) {
		jb3 := makeOCRJobSpec(t, address, bridge.Name.String(), bridge2.Name.String())
		err := orm.CreateJob(jb3)
		require.NoError(t, err)
		var jobSpec job.Job
		err = db.Get(&jobSpec, "SELECT * FROM jobs")
		require.NoError(t, err)

		ocrSpecError1 := "ocr spec 1 errored"
		ocrSpecError2 := "ocr spec 2 errored"
		require.NoError(t, orm.RecordError(jobSpec.ID, ocrSpecError1))
		require.NoError(t, orm.RecordError(jobSpec.ID, ocrSpecError1))
		require.NoError(t, orm.RecordError(jobSpec.ID, ocrSpecError2))

		var specErrors []job.SpecError
		err = db.Select(&specErrors, "SELECT * FROM job_spec_errors")
		require.NoError(t, err)
		require.Len(t, specErrors, 2)

		assert.Equal(t, specErrors[0].Occurrences, uint(2))
		assert.Equal(t, specErrors[1].Occurrences, uint(1))
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
		jb3 := makeOCRJobSpec(t, address, bridge.Name.String(), bridge2.Name.String())
		err := orm.CreateJob(jb3)
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
		require.NoError(t, orm.RecordError(jobSpec.ID, ocrSpecError1))
		require.NoError(t, orm.RecordError(jobSpec.ID, ocrSpecError2))

		var updatedSpecError []job.SpecError

		err = db.Select(&updatedSpecError, "SELECT * FROM job_spec_errors ORDER BY id ASC")
		require.NoError(t, err)
		require.Len(t, updatedSpecError, 4)

		assert.Equal(t, uint(1), updatedSpecError[2].Occurrences)
		assert.Equal(t, uint(1), updatedSpecError[3].Occurrences)
		assert.Equal(t, ocrSpecError1, updatedSpecError[2].Description)
		assert.Equal(t, ocrSpecError2, updatedSpecError[3].Description)

		dbSpecErr1, err := orm.FindSpecError(updatedSpecError[2].ID)
		require.NoError(t, err)
		dbSpecErr2, err := orm.FindSpecError(updatedSpecError[3].ID)
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
		err = orm.CreateJob(&drJob)
		require.NoError(t, err)
	})

	t.Run("creates webhook specs along with external_initiator_webhook_specs", func(t *testing.T) {
		eiFoo := cltest.MustInsertExternalInitiator(t, borm)
		eiBar := cltest.MustInsertExternalInitiator(t, borm)

		eiWS := []webhook.TOMLWebhookSpecExternalInitiator{
			{Name: eiFoo.Name, Spec: cltest.JSONFromString(t, `{}`)},
			{Name: eiBar.Name, Spec: cltest.JSONFromString(t, `{"bar": 1}`)},
		}
		eim := webhook.NewExternalInitiatorManager(db, nil, logger.TestLogger(t), config.Database())
		jb, err := webhook.ValidatedWebhookSpec(testspecs.GenerateWebhookSpec(testspecs.WebhookSpecParams{ExternalInitiators: eiWS}).Toml(), eim)
		require.NoError(t, err)

		err = orm.CreateJob(&jb)
		require.NoError(t, err)

		cltest.AssertCount(t, db, "external_initiator_webhook_specs", 2)
	})

	t.Run("it creates and deletes records for blockhash store jobs", func(t *testing.T) {
		bhsJob, err := blockhashstore.ValidatedSpec(
			testspecs.GenerateBlockhashStoreSpec(testspecs.BlockhashStoreSpecParams{}).Toml())
		require.NoError(t, err)

		err = orm.CreateJob(&bhsJob)
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
		err = orm.DeleteJob(bhsJob.ID)
		require.NoError(t, err)
		_, err = orm.FindJob(testutils.Context(t), bhsJob.ID)
		require.Error(t, err)
	})

	t.Run("it creates and deletes records for blockheaderfeeder jobs", func(t *testing.T) {
		bhsJob, err := blockheaderfeeder.ValidatedSpec(
			testspecs.GenerateBlockHeaderFeederSpec(testspecs.BlockHeaderFeederSpecParams{}).Toml())
		require.NoError(t, err)

		err = orm.CreateJob(&bhsJob)
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
		err = orm.DeleteJob(bhsJob.ID)
		require.NoError(t, err)
		_, err = orm.FindJob(testutils.Context(t), bhsJob.ID)
		require.Error(t, err)
	})
}

func TestORM_DeleteJob_DeletesAssociatedRecords(t *testing.T) {
	t.Parallel()
	config := configtest.NewGeneralConfig(t, nil)

	db := pgtest.NewSqlxDB(t)
	keyStore := cltest.NewKeyStore(t, db, config.Database())
	require.NoError(t, keyStore.OCR().Add(cltest.DefaultOCRKey))
	require.NoError(t, keyStore.P2P().Add(cltest.DefaultP2PKey))

	lggr := logger.TestLogger(t)
	pipelineORM := pipeline.NewORM(db, lggr, config.Database(), config.JobPipeline().MaxSuccessfulRuns())
	bridgesORM := bridges.NewORM(db, lggr, config.Database())
	relayExtenders := evmtest.NewChainRelayExtenders(t, evmtest.TestChainOpts{DB: db, GeneralConfig: config, KeyStore: keyStore.Eth()})
	legacyChains := evmrelay.NewLegacyChainsFromRelayerExtenders(relayExtenders)
	jobORM := NewTestORM(t, db, legacyChains, pipelineORM, bridgesORM, keyStore, config.Database())
	scopedConfig := evmtest.NewChainScopedConfig(t, config)
	korm := keeper.NewORM(db, logger.TestLogger(t), scopedConfig.Database())

	t.Run("it deletes records for offchainreporting jobs", func(t *testing.T) {
		_, bridge := cltest.MustCreateBridge(t, db, cltest.BridgeOpts{}, config.Database())
		_, bridge2 := cltest.MustCreateBridge(t, db, cltest.BridgeOpts{}, config.Database())

		_, address := cltest.MustInsertRandomKey(t, keyStore.Eth())
		jb, err := ocr.ValidatedOracleSpecToml(legacyChains, testspecs.GenerateOCRSpec(testspecs.OCRSpecParams{
			TransmitterAddress: address.Hex(),
			DS1BridgeName:      bridge.Name.String(),
			DS2BridgeName:      bridge2.Name.String(),
		}).Toml())
		require.NoError(t, err)

		err = jobORM.CreateJob(&jb)
		require.NoError(t, err)

		cltest.AssertCount(t, db, "ocr_oracle_specs", 1)
		cltest.AssertCount(t, db, "pipeline_specs", 1)

		err = jobORM.DeleteJob(jb.ID)
		require.NoError(t, err)
		cltest.AssertCount(t, db, "ocr_oracle_specs", 0)
		cltest.AssertCount(t, db, "pipeline_specs", 0)
		cltest.AssertCount(t, db, "jobs", 0)
	})

	t.Run("it deletes records for keeper jobs", func(t *testing.T) {
		registry, keeperJob := cltest.MustInsertKeeperRegistry(t, db, korm, keyStore.Eth(), 0, 1, 20)
		scoped := evmtest.NewChainScopedConfig(t, config)
		cltest.MustInsertUpkeepForRegistry(t, db, scoped.Database(), registry)

		cltest.AssertCount(t, db, "keeper_specs", 1)
		cltest.AssertCount(t, db, "keeper_registries", 1)
		cltest.AssertCount(t, db, "upkeep_registrations", 1)

		err := jobORM.DeleteJob(keeperJob.ID)
		require.NoError(t, err)
		cltest.AssertCount(t, db, "keeper_specs", 0)
		cltest.AssertCount(t, db, "keeper_registries", 0)
		cltest.AssertCount(t, db, "upkeep_registrations", 0)
		cltest.AssertCount(t, db, "jobs", 0)
	})

	t.Run("it creates and deletes records for vrf jobs", func(t *testing.T) {
		key, err := keyStore.VRF().Create()
		require.NoError(t, err)
		pk := key.PublicKey
		jb, err := vrfcommon.ValidatedVRFSpec(testspecs.GenerateVRFSpec(testspecs.VRFSpecParams{PublicKey: pk.String()}).Toml())
		require.NoError(t, err)

		err = jobORM.CreateJob(&jb)
		require.NoError(t, err)
		cltest.AssertCount(t, db, "vrf_specs", 1)
		cltest.AssertCount(t, db, "jobs", 1)
		err = jobORM.DeleteJob(jb.ID)
		require.NoError(t, err)
		cltest.AssertCount(t, db, "vrf_specs", 0)
		cltest.AssertCount(t, db, "jobs", 0)
	})

	t.Run("it deletes records for webhook jobs", func(t *testing.T) {
		ei := cltest.MustInsertExternalInitiator(t, bridges.NewORM(db, logger.TestLogger(t), config.Database()))
		jb, webhookSpec := cltest.MustInsertWebhookSpec(t, db)
		_, err := db.Exec(`INSERT INTO external_initiator_webhook_specs (external_initiator_id, webhook_spec_id, spec) VALUES ($1,$2,$3)`, ei.ID, webhookSpec.ID, `{"ei": "foo", "name": "webhookSpecTwoEIs"}`)
		require.NoError(t, err)

		err = jobORM.DeleteJob(jb.ID)
		require.NoError(t, err)
		cltest.AssertCount(t, db, "webhook_specs", 0)
		cltest.AssertCount(t, db, "external_initiator_webhook_specs", 0)
		cltest.AssertCount(t, db, "jobs", 0)
	})

	t.Run("does not allow to delete external initiators if they have referencing external_initiator_webhook_specs", func(t *testing.T) {
		// create new db because this will rollback transaction and poison it
		db := pgtest.NewSqlxDB(t)
		ei := cltest.MustInsertExternalInitiator(t, bridges.NewORM(db, logger.TestLogger(t), config.Database()))
		_, webhookSpec := cltest.MustInsertWebhookSpec(t, db)
		_, err := db.Exec(`INSERT INTO external_initiator_webhook_specs (external_initiator_id, webhook_spec_id, spec) VALUES ($1,$2,$3)`, ei.ID, webhookSpec.ID, `{"ei": "foo", "name": "webhookSpecTwoEIs"}`)
		require.NoError(t, err)

		_, err = db.Exec(`DELETE FROM external_initiators`)
		require.EqualError(t, err, "ERROR: update or delete on table \"external_initiators\" violates foreign key constraint \"external_initiator_webhook_specs_external_initiator_id_fkey\" on table \"external_initiator_webhook_specs\" (SQLSTATE 23503)")
	})
}

func TestORM_CreateJob_VRFV2(t *testing.T) {
	config := configtest.NewTestGeneralConfig(t)
	db := pgtest.NewSqlxDB(t)
	keyStore := cltest.NewKeyStore(t, db, config.Database())
	require.NoError(t, keyStore.OCR().Add(cltest.DefaultOCRKey))

	lggr := logger.TestLogger(t)
	pipelineORM := pipeline.NewORM(db, lggr, config.Database(), config.JobPipeline().MaxSuccessfulRuns())
	bridgesORM := bridges.NewORM(db, lggr, config.Database())
	relayExtenders := evmtest.NewChainRelayExtenders(t, evmtest.TestChainOpts{DB: db, GeneralConfig: config, KeyStore: keyStore.Eth()})
	legacyChains := evmrelay.NewLegacyChainsFromRelayerExtenders(relayExtenders)

	jobORM := NewTestORM(t, db, legacyChains, pipelineORM, bridgesORM, keyStore, config.Database())

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

	require.NoError(t, jobORM.CreateJob(&jb))
	cltest.AssertCount(t, db, "vrf_specs", 1)
	cltest.AssertCount(t, db, "jobs", 1)
	var requestedConfsDelay int64
	require.NoError(t, db.Get(&requestedConfsDelay, `SELECT requested_confs_delay FROM vrf_specs LIMIT 1`))
	require.Equal(t, int64(10), requestedConfsDelay)
	var batchFulfillmentEnabled bool
	require.NoError(t, db.Get(&batchFulfillmentEnabled, `SELECT batch_fulfillment_enabled FROM vrf_specs LIMIT 1`))
	require.False(t, batchFulfillmentEnabled)
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
	var vrfOwnerAddress ethkey.EIP55Address
	require.NoError(t, db.Get(&vrfOwnerAddress, `SELECT vrf_owner_address FROM vrf_specs LIMIT 1`))
	require.Equal(t, "0x32891BD79647DC9136Fc0a59AAB48c7825eb624c", vrfOwnerAddress.Address().String())
	require.NoError(t, jobORM.DeleteJob(jb.ID))
	cltest.AssertCount(t, db, "vrf_specs", 0)
	cltest.AssertCount(t, db, "jobs", 0)

	jb, err = vrfcommon.ValidatedVRFSpec(testspecs.GenerateVRFSpec(testspecs.VRFSpecParams{RequestTimeout: 1 * time.Hour}).Toml())
	require.NoError(t, err)
	require.NoError(t, jobORM.CreateJob(&jb))
	cltest.AssertCount(t, db, "vrf_specs", 1)
	cltest.AssertCount(t, db, "jobs", 1)
	require.NoError(t, db.Get(&requestedConfsDelay, `SELECT requested_confs_delay FROM vrf_specs LIMIT 1`))
	require.Equal(t, int64(0), requestedConfsDelay)
	require.NoError(t, db.Get(&requestTimeout, `SELECT request_timeout FROM vrf_specs LIMIT 1`))
	require.Equal(t, 1*time.Hour, requestTimeout)
	require.NoError(t, jobORM.DeleteJob(jb.ID))
	cltest.AssertCount(t, db, "vrf_specs", 0)
	cltest.AssertCount(t, db, "jobs", 0)
}

func TestORM_CreateJob_VRFV2Plus(t *testing.T) {
	config := configtest.NewTestGeneralConfig(t)
	db := pgtest.NewSqlxDB(t)
	keyStore := cltest.NewKeyStore(t, db, config.Database())
	require.NoError(t, keyStore.OCR().Add(cltest.DefaultOCRKey))

	lggr := logger.TestLogger(t)
	pipelineORM := pipeline.NewORM(db, lggr, config.Database(), config.JobPipeline().MaxSuccessfulRuns())
	bridgesORM := bridges.NewORM(db, lggr, config.Database())
	cc := evmtest.NewChainRelayExtenders(t, evmtest.TestChainOpts{DB: db, GeneralConfig: config, KeyStore: keyStore.Eth()})
	legacyChains := evmrelay.NewLegacyChainsFromRelayerExtenders(cc)
	jobORM := NewTestORM(t, db, legacyChains, pipelineORM, bridgesORM, keyStore, config.Database())

	fromAddresses := []string{cltest.NewEIP55Address().String(), cltest.NewEIP55Address().String()}
	jb, err := vrfcommon.ValidatedVRFSpec(testspecs.GenerateVRFSpec(
		testspecs.VRFSpecParams{
			VRFVersion:          vrfcommon.V2Plus,
			RequestedConfsDelay: 10,
			FromAddresses:       fromAddresses,
			ChunkSize:           25,
			BackoffInitialDelay: time.Minute,
			BackoffMaxDelay:     time.Hour,
			GasLanePrice:        assets.GWei(100),
		}).
		Toml())
	require.NoError(t, err)

	require.NoError(t, jobORM.CreateJob(&jb))
	cltest.AssertCount(t, db, "vrf_specs", 1)
	cltest.AssertCount(t, db, "jobs", 1)
	var requestedConfsDelay int64
	require.NoError(t, db.Get(&requestedConfsDelay, `SELECT requested_confs_delay FROM vrf_specs LIMIT 1`))
	require.Equal(t, int64(10), requestedConfsDelay)
	var batchFulfillmentEnabled bool
	require.NoError(t, db.Get(&batchFulfillmentEnabled, `SELECT batch_fulfillment_enabled FROM vrf_specs LIMIT 1`))
	require.False(t, batchFulfillmentEnabled)
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
	var vrfOwnerAddress ethkey.EIP55Address
	require.Error(t, db.Get(&vrfOwnerAddress, `SELECT vrf_owner_address FROM vrf_specs LIMIT 1`))
	require.NoError(t, jobORM.DeleteJob(jb.ID))
	cltest.AssertCount(t, db, "vrf_specs", 0)
	cltest.AssertCount(t, db, "jobs", 0)

	jb, err = vrfcommon.ValidatedVRFSpec(testspecs.GenerateVRFSpec(testspecs.VRFSpecParams{
		VRFVersion:     vrfcommon.V2Plus,
		RequestTimeout: 1 * time.Hour,
		FromAddresses:  fromAddresses,
	}).Toml())
	require.NoError(t, err)
	require.NoError(t, jobORM.CreateJob(&jb))
	cltest.AssertCount(t, db, "vrf_specs", 1)
	cltest.AssertCount(t, db, "jobs", 1)
	require.NoError(t, db.Get(&requestedConfsDelay, `SELECT requested_confs_delay FROM vrf_specs LIMIT 1`))
	require.Equal(t, int64(0), requestedConfsDelay)
	require.NoError(t, db.Get(&requestTimeout, `SELECT request_timeout FROM vrf_specs LIMIT 1`))
	require.Equal(t, 1*time.Hour, requestTimeout)
	require.NoError(t, jobORM.DeleteJob(jb.ID))
	cltest.AssertCount(t, db, "vrf_specs", 0)
	cltest.AssertCount(t, db, "jobs", 0)
}

func TestORM_CreateJob_OCRBootstrap(t *testing.T) {
	config := configtest.NewTestGeneralConfig(t)
	db := pgtest.NewSqlxDB(t)
	keyStore := cltest.NewKeyStore(t, db, config.Database())
	require.NoError(t, keyStore.OCR().Add(cltest.DefaultOCRKey))

	lggr := logger.TestLogger(t)
	pipelineORM := pipeline.NewORM(db, lggr, config.Database(), config.JobPipeline().MaxSuccessfulRuns())
	bridgesORM := bridges.NewORM(db, lggr, config.Database())
	relayExtenders := evmtest.NewChainRelayExtenders(t, evmtest.TestChainOpts{DB: db, GeneralConfig: config, KeyStore: keyStore.Eth()})
	legacyChains := evmrelay.NewLegacyChainsFromRelayerExtenders(relayExtenders)
	jobORM := NewTestORM(t, db, legacyChains, pipelineORM, bridgesORM, keyStore, config.Database())

	jb, err := ocrbootstrap.ValidatedBootstrapSpecToml(testspecs.GetOCRBootstrapSpec())
	require.NoError(t, err)

	err = jobORM.CreateJob(&jb)
	require.NoError(t, err)
	cltest.AssertCount(t, db, "bootstrap_specs", 1)
	cltest.AssertCount(t, db, "jobs", 1)
	var relay string
	require.NoError(t, db.Get(&relay, `SELECT relay FROM bootstrap_specs LIMIT 1`))
	require.Equal(t, "evm", relay)

	require.NoError(t, jobORM.DeleteJob(jb.ID))
	cltest.AssertCount(t, db, "bootstrap_specs", 0)
	cltest.AssertCount(t, db, "jobs", 0)
}

func TestORM_CreateJob_EVMChainID_Validation(t *testing.T) {
	config := configtest.NewGeneralConfig(t, nil)
	db := pgtest.NewSqlxDB(t)
	keyStore := cltest.NewKeyStore(t, db, config.Database())

	lggr := logger.TestLogger(t)
	pipelineORM := pipeline.NewORM(db, lggr, config.Database(), config.JobPipeline().MaxSuccessfulRuns())
	bridgesORM := bridges.NewORM(db, lggr, config.Database())

	relayExtenders := evmtest.NewChainRelayExtenders(t, evmtest.TestChainOpts{DB: db, GeneralConfig: config, KeyStore: keyStore.Eth()})
	legacyChains := evmrelay.NewLegacyChainsFromRelayerExtenders(relayExtenders)
	jobORM := NewTestORM(t, db, legacyChains, pipelineORM, bridgesORM, keyStore, config.Database())

	t.Run("evm chain id validation for ocr works", func(t *testing.T) {
		jb := job.Job{
			Type:          job.OffchainReporting,
			OCROracleSpec: &job.OCROracleSpec{},
		}
		assert.Equal(t, "CreateJobFailed: evm chain id must be defined", jobORM.CreateJob(&jb).Error())
	})

	t.Run("evm chain id validation for direct request works", func(t *testing.T) {
		jb := job.Job{
			Type:              job.DirectRequest,
			DirectRequestSpec: &job.DirectRequestSpec{},
		}
		assert.Equal(t, "CreateJobFailed: evm chain id must be defined", jobORM.CreateJob(&jb).Error())
	})

	t.Run("evm chain id validation for flux monitor works", func(t *testing.T) {
		jb := job.Job{
			Type:            job.FluxMonitor,
			FluxMonitorSpec: &job.FluxMonitorSpec{},
		}
		assert.Equal(t, "CreateJobFailed: evm chain id must be defined", jobORM.CreateJob(&jb).Error())
	})

	t.Run("evm chain id validation for keepers works", func(t *testing.T) {
		jb := job.Job{
			Type:       job.Keeper,
			KeeperSpec: &job.KeeperSpec{},
		}
		assert.Equal(t, "CreateJobFailed: evm chain id must be defined", jobORM.CreateJob(&jb).Error())
	})

	t.Run("evm chain id validation for vrf works", func(t *testing.T) {
		jb := job.Job{
			Type:    job.VRF,
			VRFSpec: &job.VRFSpec{},
		}
		assert.Equal(t, "CreateJobFailed: evm chain id must be defined", jobORM.CreateJob(&jb).Error())
	})

	t.Run("evm chain id validation for block hash store works", func(t *testing.T) {
		jb := job.Job{
			Type:               job.BlockhashStore,
			BlockhashStoreSpec: &job.BlockhashStoreSpec{},
		}
		assert.Equal(t, "CreateJobFailed: evm chain id must be defined", jobORM.CreateJob(&jb).Error())
	})

	t.Run("evm chain id validation for block header feeder works", func(t *testing.T) {
		jb := job.Job{
			Type:                  job.BlockHeaderFeeder,
			BlockHeaderFeederSpec: &job.BlockHeaderFeederSpec{},
		}
		assert.Equal(t, "CreateJobFailed: evm chain id must be defined", jobORM.CreateJob(&jb).Error())
	})

	t.Run("evm chain id validation for legacy gas station server spec works", func(t *testing.T) {
		jb := job.Job{
			Type:                       job.LegacyGasStationServer,
			LegacyGasStationServerSpec: &job.LegacyGasStationServerSpec{},
		}
		assert.Equal(t, "CreateJobFailed: evm chain id must be defined", jobORM.CreateJob(&jb).Error())
	})

	t.Run("evm chain id validation for legacy gas station sidecar spec works", func(t *testing.T) {
		jb := job.Job{
			Type:                        job.LegacyGasStationSidecar,
			LegacyGasStationSidecarSpec: &job.LegacyGasStationSidecarSpec{},
		}
		assert.Equal(t, "CreateJobFailed: evm chain id must be defined", jobORM.CreateJob(&jb).Error())
	})
}

func TestORM_CreateJob_OCR_DuplicatedContractAddress(t *testing.T) {
	customChainID := utils.NewBig(testutils.NewRandomEVMChainID())

	config := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		enabled := true
		c.EVM = append(c.EVM, &evmcfg.EVMConfig{
			ChainID: customChainID,
			Chain:   evmcfg.Defaults(customChainID),
			Enabled: &enabled,
			Nodes:   evmcfg.EVMNodes{{}},
		})
	})
	db := pgtest.NewSqlxDB(t)
	keyStore := cltest.NewKeyStore(t, db, config.Database())
	require.NoError(t, keyStore.OCR().Add(cltest.DefaultOCRKey))

	lggr := logger.TestLogger(t)
	pipelineORM := pipeline.NewORM(db, lggr, config.Database(), config.JobPipeline().MaxSuccessfulRuns())
	bridgesORM := bridges.NewORM(db, lggr, config.Database())

	relayExtenders := evmtest.NewChainRelayExtenders(t, evmtest.TestChainOpts{DB: db, GeneralConfig: config, KeyStore: keyStore.Eth()})
	legacyChains := evmrelay.NewLegacyChainsFromRelayerExtenders(relayExtenders)
	jobORM := NewTestORM(t, db, legacyChains, pipelineORM, bridgesORM, keyStore, config.Database())

	// defaultChainID is deprecated
	defaultChainID := customChainID
	_, address := cltest.MustInsertRandomKey(t, keyStore.Eth())
	_, bridge := cltest.MustCreateBridge(t, db, cltest.BridgeOpts{}, config.Database())
	_, bridge2 := cltest.MustCreateBridge(t, db, cltest.BridgeOpts{}, config.Database())

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
	jb, err := ocr.ValidatedOracleSpecToml(legacyChains, spec.Toml())
	require.NoError(t, err)

	t.Run("with a set chain id", func(t *testing.T) {
		err = jobORM.CreateJob(&jb) // Add job with custom chain id
		require.NoError(t, err)

		cltest.AssertCount(t, db, "ocr_oracle_specs", 1)
		cltest.AssertCount(t, db, "jobs", 1)

		externalJobID = uuid.NullUUID{UUID: uuid.New(), Valid: true}
		spec.JobID = externalJobID.UUID.String()
		jba, err := ocr.ValidatedOracleSpecToml(legacyChains, spec.Toml())
		require.NoError(t, err)
		err = jobORM.CreateJob(&jba) // Try to add duplicate job with default id
		require.Error(t, err)
		assert.Equal(t, fmt.Sprintf("CreateJobFailed: a job with contract address %s already exists for chain ID %s", jb.OCROracleSpec.ContractAddress, defaultChainID.String()), err.Error())

		externalJobID = uuid.NullUUID{UUID: uuid.New(), Valid: true}
		spec.JobID = externalJobID.UUID.String()
		jb2, err := ocr.ValidatedOracleSpecToml(legacyChains, spec.Toml())
		require.NoError(t, err)

		err = jobORM.CreateJob(&jb2) // Try to add duplicate job with custom id
		require.Error(t, err)
		assert.Equal(t, fmt.Sprintf("CreateJobFailed: a job with contract address %s already exists for chain ID %s", jb2.OCROracleSpec.ContractAddress, customChainID), err.Error())
	})
}

func TestORM_CreateJob_OCR2_DuplicatedContractAddress(t *testing.T) {
	customChainID := utils.NewBig(testutils.NewRandomEVMChainID())

	config := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		enabled := true
		c.EVM = append(c.EVM, &evmcfg.EVMConfig{
			ChainID: customChainID,
			Chain:   evmcfg.Defaults(customChainID),
			Enabled: &enabled,
			Nodes:   evmcfg.EVMNodes{{}},
		})
	})
	db := pgtest.NewSqlxDB(t)
	keyStore := cltest.NewKeyStore(t, db, config.Database())
	require.NoError(t, keyStore.OCR2().Add(cltest.DefaultOCR2Key))

	lggr := logger.TestLogger(t)
	pipelineORM := pipeline.NewORM(db, lggr, config.Database(), config.JobPipeline().MaxSuccessfulRuns())
	bridgesORM := bridges.NewORM(db, lggr, config.Database())

	relayExtenders := evmtest.NewChainRelayExtenders(t, evmtest.TestChainOpts{DB: db, GeneralConfig: config, KeyStore: keyStore.Eth()})
	legacyChains := evmrelay.NewLegacyChainsFromRelayerExtenders(relayExtenders)
	jobORM := NewTestORM(t, db, legacyChains, pipelineORM, bridgesORM, keyStore, config.Database())

	_, address := cltest.MustInsertRandomKey(t, keyStore.Eth())

	jb, err := ocr2validate.ValidatedOracleSpecToml(config.OCR2(), config.Insecure(), testspecs.GetOCR2EVMSpecMinimal())
	require.NoError(t, err)

	const juelsPerFeeCoinSource = `
	ds          [type=http method=GET url="https://chain.link/ETH-USD"];
	ds_parse    [type=jsonparse path="data.price" separator="."];
	ds_multiply [type=multiply times=100];
	ds -> ds_parse -> ds_multiply;`

	jb.Name = null.StringFrom("Job 1")
	jb.OCR2OracleSpec.TransmitterID = null.StringFrom(address.String())
	jb.OCR2OracleSpec.PluginConfig["juelsPerFeeCoinSource"] = juelsPerFeeCoinSource

	err = jobORM.CreateJob(&jb)
	require.NoError(t, err)

	jb2, err := ocr2validate.ValidatedOracleSpecToml(config.OCR2(), config.Insecure(), testspecs.GetOCR2EVMSpecMinimal())
	require.NoError(t, err)

	jb2.Name = null.StringFrom("Job with same chain id & contract address")
	jb2.OCR2OracleSpec.TransmitterID = null.StringFrom(address.String())
	jb.OCR2OracleSpec.PluginConfig["juelsPerFeeCoinSource"] = juelsPerFeeCoinSource

	err = jobORM.CreateJob(&jb2)
	require.Error(t, err)

	jb3, err := ocr2validate.ValidatedOracleSpecToml(config.OCR2(), config.Insecure(), testspecs.GetOCR2EVMSpecMinimal())
	require.NoError(t, err)
	jb3.Name = null.StringFrom("Job with different chain id & same contract address")
	jb3.OCR2OracleSpec.TransmitterID = null.StringFrom(address.String())
	jb3.OCR2OracleSpec.RelayConfig["chainID"] = customChainID.Int64()
	jb.OCR2OracleSpec.PluginConfig["juelsPerFeeCoinSource"] = juelsPerFeeCoinSource

	err = jobORM.CreateJob(&jb3)
	require.Error(t, err)
}

func TestORM_CreateJob_OCR2_Sending_Keys_Transmitter_Keys_Validations(t *testing.T) {
	customChainID := utils.NewBig(testutils.NewRandomEVMChainID())

	config := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		enabled := true
		c.EVM = append(c.EVM, &evmcfg.EVMConfig{
			ChainID: customChainID,
			Chain:   evmcfg.Defaults(customChainID),
			Enabled: &enabled,
			Nodes:   evmcfg.EVMNodes{{}},
		})
	})
	db := pgtest.NewSqlxDB(t)
	keyStore := cltest.NewKeyStore(t, db, config.Database())
	require.NoError(t, keyStore.OCR2().Add(cltest.DefaultOCR2Key))

	lggr := logger.TestLogger(t)
	pipelineORM := pipeline.NewORM(db, lggr, config.Database(), config.JobPipeline().MaxSuccessfulRuns())
	bridgesORM := bridges.NewORM(db, lggr, config.Database())

	relayExtenders := evmtest.NewChainRelayExtenders(t, evmtest.TestChainOpts{DB: db, GeneralConfig: config, KeyStore: keyStore.Eth()})
	require.True(t, relayExtenders.Len() > 0)
	legacyChains := evmrelay.NewLegacyChainsFromRelayerExtenders(relayExtenders)
	jobORM := NewTestORM(t, db, legacyChains, pipelineORM, bridgesORM, keyStore, config.Database())

	jb, err := ocr2validate.ValidatedOracleSpecToml(config.OCR2(), config.Insecure(), testspecs.GetOCR2EVMSpecMinimal())
	require.NoError(t, err)

	t.Run("sending keys or transmitterID must be defined", func(t *testing.T) {
		jb.OCR2OracleSpec.TransmitterID = null.String{}
		assert.Equal(t, "CreateJobFailed: neither sending keys nor transmitter ID is defined", jobORM.CreateJob(&jb).Error())
	})

	_, address := cltest.MustInsertRandomKey(t, keyStore.Eth())
	t.Run("sending keys validation works properly", func(t *testing.T) {
		jb.OCR2OracleSpec.TransmitterID = null.String{}
		_, address2 := cltest.MustInsertRandomKey(t, keyStore.Eth())
		jb.OCR2OracleSpec.RelayConfig["sendingKeys"] = interface{}([]any{address.String(), address2.String(), common.HexToAddress("0X0").String()})
		assert.Equal(t, "CreateJobFailed: no EVM key matching: \"0x0000000000000000000000000000000000000000\": no such sending key exists", jobORM.CreateJob(&jb).Error())

		jb.OCR2OracleSpec.RelayConfig["sendingKeys"] = interface{}([]any{1, 2, 3})
		assert.Equal(t, "CreateJobFailed: sending keys are of wrong type", jobORM.CreateJob(&jb).Error())
	})

	t.Run("sending keys and transmitter ID can't both be defined", func(t *testing.T) {
		jb.OCR2OracleSpec.TransmitterID = null.StringFrom(address.String())
		jb.OCR2OracleSpec.RelayConfig["sendingKeys"] = interface{}([]any{address.String()})
		assert.Equal(t, "CreateJobFailed: sending keys and transmitter ID can't both be defined", jobORM.CreateJob(&jb).Error())
	})

	t.Run("transmitter validation works", func(t *testing.T) {
		jb.OCR2OracleSpec.TransmitterID = null.StringFrom("transmitterID that doesn't have a match in key store")
		jb.OCR2OracleSpec.RelayConfig["sendingKeys"] = nil
		assert.Equal(t, "CreateJobFailed: no EVM key matching: \"transmitterID that doesn't have a match in key store\": no such transmitter key exists", jobORM.CreateJob(&jb).Error())
	})
}

func TestORM_ValidateKeyStoreMatch(t *testing.T) {
	config := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {})

	keyStore := cltest.NewKeyStore(t, pgtest.NewSqlxDB(t), config.Database())
	require.NoError(t, keyStore.OCR2().Add(cltest.DefaultOCR2Key))

	var jb job.Job
	{
		var err error
		jb, err = ocr2validate.ValidatedOracleSpecToml(config.OCR2(), config.Insecure(), testspecs.GetOCR2EVMSpecMinimal())
		require.NoError(t, err)
	}

	t.Run("test ETH key validation", func(t *testing.T) {
		jb.OCR2OracleSpec.Relay = relay.EVM
		err := job.ValidateKeyStoreMatch(jb.OCR2OracleSpec, keyStore, "bad key")
		require.EqualError(t, err, "no EVM key matching: \"bad key\"")

		_, evmKey := cltest.MustInsertRandomKey(t, keyStore.Eth())
		err = job.ValidateKeyStoreMatch(jb.OCR2OracleSpec, keyStore, evmKey.String())
		require.NoError(t, err)
	})

	t.Run("test Cosmos key validation", func(t *testing.T) {
		jb.OCR2OracleSpec.Relay = relay.Cosmos
		err := job.ValidateKeyStoreMatch(jb.OCR2OracleSpec, keyStore, "bad key")
		require.EqualError(t, err, "no Cosmos key matching: \"bad key\"")

		cosmosKey, err := keyStore.Cosmos().Create()
		require.NoError(t, err)
		err = job.ValidateKeyStoreMatch(jb.OCR2OracleSpec, keyStore, cosmosKey.ID())
		require.NoError(t, err)
	})

	t.Run("test Solana key validation", func(t *testing.T) {
		jb.OCR2OracleSpec.Relay = relay.Solana

		err := job.ValidateKeyStoreMatch(jb.OCR2OracleSpec, keyStore, "bad key")
		require.EqualError(t, err, "no Solana key matching: \"bad key\"")

		solanaKey, err := keyStore.Solana().Create()
		require.NoError(t, err)
		err = job.ValidateKeyStoreMatch(jb.OCR2OracleSpec, keyStore, solanaKey.ID())
		require.NoError(t, err)
	})

	t.Run("test Starknet key validation", func(t *testing.T) {
		jb.OCR2OracleSpec.Relay = relay.StarkNet
		err := job.ValidateKeyStoreMatch(jb.OCR2OracleSpec, keyStore, "bad key")
		require.EqualError(t, err, "no Starknet key matching: \"bad key\"")

		starkNetKey, err := keyStore.StarkNet().Create()
		require.NoError(t, err)
		err = job.ValidateKeyStoreMatch(jb.OCR2OracleSpec, keyStore, starkNetKey.ID())
		require.NoError(t, err)
	})

	t.Run("test Mercury ETH key validation", func(t *testing.T) {
		jb.OCR2OracleSpec.PluginType = types.Mercury
		err := job.ValidateKeyStoreMatch(jb.OCR2OracleSpec, keyStore, "bad key")
		require.EqualError(t, err, "no CSA key matching: \"bad key\"")

		csaKey, err := keyStore.CSA().Create()
		require.NoError(t, err)
		err = job.ValidateKeyStoreMatch(jb.OCR2OracleSpec, keyStore, csaKey.ID())
		require.NoError(t, err)
	})
}

func Test_FindJobs(t *testing.T) {
	t.Parallel()

	config := configtest.NewTestGeneralConfig(t)
	db := pgtest.NewSqlxDB(t)
	keyStore := cltest.NewKeyStore(t, db, config.Database())
	require.NoError(t, keyStore.OCR().Add(cltest.DefaultOCRKey))
	require.NoError(t, keyStore.P2P().Add(cltest.DefaultP2PKey))

	pipelineORM := pipeline.NewORM(db, logger.TestLogger(t), config.Database(), config.JobPipeline().MaxSuccessfulRuns())
	bridgesORM := bridges.NewORM(db, logger.TestLogger(t), config.Database())
	relayExtenders := evmtest.NewChainRelayExtenders(t, evmtest.TestChainOpts{DB: db, GeneralConfig: config, KeyStore: keyStore.Eth()})
	legacyChains := evmrelay.NewLegacyChainsFromRelayerExtenders(relayExtenders)
	orm := NewTestORM(t, db, legacyChains, pipelineORM, bridgesORM, keyStore, config.Database())

	_, bridge := cltest.MustCreateBridge(t, db, cltest.BridgeOpts{}, config.Database())
	_, bridge2 := cltest.MustCreateBridge(t, db, cltest.BridgeOpts{}, config.Database())

	_, address := cltest.MustInsertRandomKey(t, keyStore.Eth())
	jb1, err := ocr.ValidatedOracleSpecToml(legacyChains,
		testspecs.GenerateOCRSpec(testspecs.OCRSpecParams{
			JobID:              uuid.New().String(),
			TransmitterAddress: address.Hex(),
			DS1BridgeName:      bridge.Name.String(),
			DS2BridgeName:      bridge2.Name.String(),
		}).Toml(),
	)
	require.NoError(t, err)

	err = orm.CreateJob(&jb1)
	require.NoError(t, err)

	jb2, err := directrequest.ValidatedDirectRequestSpec(
		testspecs.GetDirectRequestSpec(),
	)
	require.NoError(t, err)

	err = orm.CreateJob(&jb2)
	require.NoError(t, err)

	t.Run("jobs are ordered by latest first", func(t *testing.T) {
		jobs, count, err2 := orm.FindJobs(0, 2)
		require.NoError(t, err2)
		require.Len(t, jobs, 2)
		assert.Equal(t, count, 2)

		expectedJobs := []job.Job{jb2, jb1}

		for i, exp := range expectedJobs {
			assert.Equal(t, exp.ID, jobs[i].ID)
		}
	})

	t.Run("jobs respect pagination", func(t *testing.T) {
		jobs, count, err2 := orm.FindJobs(0, 1)
		require.NoError(t, err2)
		require.Len(t, jobs, 1)
		assert.Equal(t, count, 2)

		expectedJobs := []job.Job{jb2}

		for i, exp := range expectedJobs {
			assert.Equal(t, exp.ID, jobs[i].ID)
		}
	})
}

func Test_FindJob(t *testing.T) {
	t.Parallel()

	// Create a config with multiple EVM chains. The test fixtures already load 1337
	// Additional chains will need additional fixture statements to add a chain to evm_chains.
	config := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		chainID := utils.NewBigI(1337)
		enabled := true
		c.EVM = append(c.EVM, &evmcfg.EVMConfig{
			ChainID: chainID,
			Chain:   evmcfg.Defaults(chainID),
			Enabled: &enabled,
			Nodes:   evmcfg.EVMNodes{{}},
		})
	})

	db := pgtest.NewSqlxDB(t)
	keyStore := cltest.NewKeyStore(t, db, config.Database())
	require.NoError(t, keyStore.OCR().Add(cltest.DefaultOCRKey))
	require.NoError(t, keyStore.P2P().Add(cltest.DefaultP2PKey))
	require.NoError(t, keyStore.CSA().Add(cltest.DefaultCSAKey))

	pipelineORM := pipeline.NewORM(db, logger.TestLogger(t), config.Database(), config.JobPipeline().MaxSuccessfulRuns())
	bridgesORM := bridges.NewORM(db, logger.TestLogger(t), config.Database())
	relayExtenders := evmtest.NewChainRelayExtenders(t, evmtest.TestChainOpts{DB: db, GeneralConfig: config, KeyStore: keyStore.Eth()})
	legacyChains := evmrelay.NewLegacyChainsFromRelayerExtenders(relayExtenders)
	orm := NewTestORM(t, db, legacyChains, pipelineORM, bridgesORM, keyStore, config.Database())

	_, bridge := cltest.MustCreateBridge(t, db, cltest.BridgeOpts{}, config.Database())
	_, bridge2 := cltest.MustCreateBridge(t, db, cltest.BridgeOpts{}, config.Database())

	// Create two jobs.  Each job has the same Transmitter Address but on a different chain.
	// Must uniquely name the OCR Specs to properly insert a new job in the job table.
	externalJobID := uuid.New()
	_, address := cltest.MustInsertRandomKey(t, keyStore.Eth())
	job, err := ocr.ValidatedOracleSpecToml(legacyChains,
		testspecs.GenerateOCRSpec(testspecs.OCRSpecParams{
			JobID:              externalJobID.String(),
			Name:               "orig ocr spec",
			TransmitterAddress: address.Hex(),
			DS1BridgeName:      bridge.Name.String(),
			DS2BridgeName:      bridge2.Name.String(),
		}).Toml(),
	)
	require.NoError(t, err)

	jobSameAddress, err := ocr.ValidatedOracleSpecToml(legacyChains,
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

	jobOCR2, err := ocr2validate.ValidatedOracleSpecToml(config.OCR2(), config.Insecure(), testspecs.GetOCR2EVMSpecMinimal())
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
		config.OCR2(),
		config.Insecure(),
		fmt.Sprintf(mercuryOracleTOML, cltest.DefaultCSAKey.PublicKeyString(), ocr2WithFeedID1),
	)
	require.NoError(t, err)

	jobOCR2WithFeedID2, err := ocr2validate.ValidatedOracleSpecToml(
		config.OCR2(),
		config.Insecure(),
		fmt.Sprintf(mercuryOracleTOML, cltest.DefaultCSAKey.PublicKeyString(), ocr2WithFeedID2),
	)
	jobOCR2WithFeedID2.ExternalJobID = uuid.New()
	jobOCR2WithFeedID2.Name = null.StringFrom("new name")
	require.NoError(t, err)

	err = orm.CreateJob(&job)
	require.NoError(t, err)

	err = orm.CreateJob(&jobSameAddress)
	require.NoError(t, err)

	err = orm.CreateJob(&jobOCR2)
	require.NoError(t, err)

	err = orm.CreateJob(&jobOCR2WithFeedID1)
	require.NoError(t, err)

	// second ocr2 job with same contract id but different feed id
	err = orm.CreateJob(&jobOCR2WithFeedID2)
	require.NoError(t, err)

	t.Run("by id", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(testutils.Context(t), 5*time.Second)
		defer cancel()
		jb, err2 := orm.FindJob(ctx, job.ID)
		require.NoError(t, err2)

		assert.Equal(t, jb.ID, job.ID)
		assert.Equal(t, jb.Name, job.Name)

		require.Greater(t, jb.PipelineSpecID, int32(0))
		require.NotNil(t, jb.PipelineSpec)
		require.NotNil(t, jb.OCROracleSpecID)
		require.NotNil(t, jb.OCROracleSpec)
	})

	t.Run("by external job id", func(t *testing.T) {
		jb, err2 := orm.FindJobByExternalJobID(externalJobID)
		require.NoError(t, err2)

		assert.Equal(t, jb.ID, job.ID)
		assert.Equal(t, jb.Name, job.Name)

		require.Greater(t, jb.PipelineSpecID, int32(0))
		require.NotNil(t, jb.PipelineSpec)
		require.NotNil(t, jb.OCROracleSpecID)
		require.NotNil(t, jb.OCROracleSpec)
	})

	t.Run("by address", func(t *testing.T) {
		jbID, err2 := orm.FindJobIDByAddress(job.OCROracleSpec.ContractAddress, job.OCROracleSpec.EVMChainID)
		require.NoError(t, err2)

		assert.Equal(t, job.ID, jbID)

		_, err2 = orm.FindJobIDByAddress("not-existing", utils.NewBigI(0))
		require.Error(t, err2)
		require.ErrorIs(t, err2, sql.ErrNoRows)
	})

	t.Run("by address yet chain scoped", func(t *testing.T) {
		commonAddr := jobSameAddress.OCROracleSpec.ContractAddress

		// Find job ID for job on chain 1337 with common address.
		jbID, err2 := orm.FindJobIDByAddress(commonAddr, jobSameAddress.OCROracleSpec.EVMChainID)
		require.NoError(t, err2)

		assert.Equal(t, jobSameAddress.ID, jbID)

		// Find job ID for job on default evm chain with common address.
		jbID, err2 = orm.FindJobIDByAddress(commonAddr, job.OCROracleSpec.EVMChainID)
		require.NoError(t, err2)

		assert.Equal(t, job.ID, jbID)
	})

	t.Run("by contract id without feed id", func(t *testing.T) {
		contractID := "0x613a38AC1659769640aaE063C651F48E0250454C"

		// Find job ID for ocr2 job without feedID.
		jbID, err2 := orm.FindOCR2JobIDByAddress(contractID, nil)
		require.NoError(t, err2)

		assert.Equal(t, jobOCR2.ID, jbID)
	})

	t.Run("by contract id with valid feed id", func(t *testing.T) {
		contractID := "0x0000000000000000000000000000000000000006"
		feedID := common.HexToHash(ocr2WithFeedID1)

		// Find job ID for ocr2 job with feed ID
		jbID, err2 := orm.FindOCR2JobIDByAddress(contractID, &feedID)
		require.NoError(t, err2)

		assert.Equal(t, jobOCR2WithFeedID1.ID, jbID)
	})

	t.Run("with duplicate contract id but different feed id", func(t *testing.T) {
		contractID := "0x0000000000000000000000000000000000000006"
		feedID := common.HexToHash(ocr2WithFeedID2)

		// Find job ID for ocr2 job with feed ID
		jbID, err2 := orm.FindOCR2JobIDByAddress(contractID, &feedID)
		require.NoError(t, err2)

		assert.Equal(t, jobOCR2WithFeedID2.ID, jbID)
	})
}

func Test_FindJobsByPipelineSpecIDs(t *testing.T) {
	t.Parallel()

	config := configtest.NewTestGeneralConfig(t)
	db := pgtest.NewSqlxDB(t)
	keyStore := cltest.NewKeyStore(t, db, config.Database())
	require.NoError(t, keyStore.OCR().Add(cltest.DefaultOCRKey))

	pipelineORM := pipeline.NewORM(db, logger.TestLogger(t), config.Database(), config.JobPipeline().MaxSuccessfulRuns())
	bridgesORM := bridges.NewORM(db, logger.TestLogger(t), config.Database())
	relayExtenders := evmtest.NewChainRelayExtenders(t, evmtest.TestChainOpts{DB: db, GeneralConfig: config, KeyStore: keyStore.Eth()})
	legacyChains := evmrelay.NewLegacyChainsFromRelayerExtenders(relayExtenders)
	orm := NewTestORM(t, db, legacyChains, pipelineORM, bridgesORM, keyStore, config.Database())

	jb, err := directrequest.ValidatedDirectRequestSpec(testspecs.GetDirectRequestSpec())
	require.NoError(t, err)
	jb.DirectRequestSpec.EVMChainID = utils.NewBigI(0)

	err = orm.CreateJob(&jb)
	require.NoError(t, err)

	t.Run("with jobs", func(t *testing.T) {
		jbs, err2 := orm.FindJobsByPipelineSpecIDs([]int32{jb.PipelineSpecID})
		require.NoError(t, err2)
		assert.Len(t, jbs, 1)

		assert.Equal(t, jb.ID, jbs[0].ID)
		assert.Equal(t, jb.Name, jbs[0].Name)

		require.Greater(t, jbs[0].PipelineSpecID, int32(0))
		require.Equal(t, jb.PipelineSpecID, jbs[0].PipelineSpecID)
		require.NotNil(t, jbs[0].PipelineSpec)
	})

	t.Run("without jobs", func(t *testing.T) {
		jbs, err2 := orm.FindJobsByPipelineSpecIDs([]int32{-1})
		require.NoError(t, err2)
		assert.Len(t, jbs, 0)
	})

	t.Run("with chainID disabled", func(t *testing.T) {
		newCfg := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
			c.EVM[0] = &evmcfg.EVMConfig{
				ChainID: utils.NewBigI(0),
				Enabled: ptr(false),
			}
			c.EVM = append(c.EVM, &evmcfg.EVMConfig{
				ChainID: utils.NewBigI(123123123),
				Enabled: ptr(true),
			})
		})
		relayExtenders2 := evmtest.NewChainRelayExtenders(t, evmtest.TestChainOpts{DB: db, GeneralConfig: newCfg, KeyStore: keyStore.Eth()})
		legacyChains2 := evmrelay.NewLegacyChainsFromRelayerExtenders(relayExtenders2)

		orm2 := NewTestORM(t, db, legacyChains2, pipelineORM, bridgesORM, keyStore, config.Database())

		jbs, err2 := orm2.FindJobsByPipelineSpecIDs([]int32{jb.PipelineSpecID})
		require.NoError(t, err2)
		assert.Len(t, jbs, 1)
	})
}

func Test_FindPipelineRuns(t *testing.T) {
	t.Parallel()

	config := configtest.NewTestGeneralConfig(t)
	db := pgtest.NewSqlxDB(t)
	keyStore := cltest.NewKeyStore(t, db, config.Database())
	require.NoError(t, keyStore.OCR().Add(cltest.DefaultOCRKey))
	require.NoError(t, keyStore.P2P().Add(cltest.DefaultP2PKey))

	pipelineORM := pipeline.NewORM(db, logger.TestLogger(t), config.Database(), config.JobPipeline().MaxSuccessfulRuns())
	bridgesORM := bridges.NewORM(db, logger.TestLogger(t), config.Database())
	relayExtenders := evmtest.NewChainRelayExtenders(t, evmtest.TestChainOpts{DB: db, GeneralConfig: config, KeyStore: keyStore.Eth()})
	legacyChains := evmrelay.NewLegacyChainsFromRelayerExtenders(relayExtenders)
	orm := NewTestORM(t, db, legacyChains, pipelineORM, bridgesORM, keyStore, config.Database())

	_, bridge := cltest.MustCreateBridge(t, db, cltest.BridgeOpts{}, config.Database())
	_, bridge2 := cltest.MustCreateBridge(t, db, cltest.BridgeOpts{}, config.Database())

	externalJobID := uuid.New()
	_, address := cltest.MustInsertRandomKey(t, keyStore.Eth())
	jb, err := ocr.ValidatedOracleSpecToml(legacyChains,
		testspecs.GenerateOCRSpec(testspecs.OCRSpecParams{
			JobID:              externalJobID.String(),
			TransmitterAddress: address.Hex(),
			DS1BridgeName:      bridge.Name.String(),
			DS2BridgeName:      bridge2.Name.String(),
		}).Toml(),
	)
	require.NoError(t, err)

	err = orm.CreateJob(&jb)
	require.NoError(t, err)

	t.Run("with no pipeline runs", func(t *testing.T) {
		runs, count, err2 := orm.PipelineRuns(nil, 0, 10)
		require.NoError(t, err2)
		assert.Equal(t, count, 0)
		assert.Empty(t, runs)
	})

	t.Run("with a pipeline run", func(t *testing.T) {
		run := mustInsertPipelineRun(t, pipelineORM, jb)

		runs, count, err2 := orm.PipelineRuns(nil, 0, 10)
		require.NoError(t, err2)

		assert.Equal(t, count, 1)
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

	config := configtest.NewTestGeneralConfig(t)
	db := pgtest.NewSqlxDB(t)

	keyStore := cltest.NewKeyStore(t, db, config.Database())
	require.NoError(t, keyStore.OCR().Add(cltest.DefaultOCRKey))
	require.NoError(t, keyStore.P2P().Add(cltest.DefaultP2PKey))

	pipelineORM := pipeline.NewORM(db, logger.TestLogger(t), config.Database(), config.JobPipeline().MaxSuccessfulRuns())
	bridgesORM := bridges.NewORM(db, logger.TestLogger(t), config.Database())
	relayExtenders := evmtest.NewChainRelayExtenders(t, evmtest.TestChainOpts{DB: db, GeneralConfig: config, KeyStore: keyStore.Eth()})
	legacyChains := evmrelay.NewLegacyChainsFromRelayerExtenders(relayExtenders)
	orm := NewTestORM(t, db, legacyChains, pipelineORM, bridgesORM, keyStore, config.Database())

	_, bridge := cltest.MustCreateBridge(t, db, cltest.BridgeOpts{}, config.Database())
	_, bridge2 := cltest.MustCreateBridge(t, db, cltest.BridgeOpts{}, config.Database())

	externalJobID := uuid.New()
	_, address := cltest.MustInsertRandomKey(t, keyStore.Eth())
	jb, err := ocr.ValidatedOracleSpecToml(legacyChains,
		testspecs.GenerateOCRSpec(testspecs.OCRSpecParams{
			JobID:              externalJobID.String(),
			TransmitterAddress: address.Hex(),
			DS1BridgeName:      bridge.Name.String(),
			DS2BridgeName:      bridge2.Name.String(),
		}).Toml(),
	)
	require.NoError(t, err)

	err = orm.CreateJob(&jb)
	require.NoError(t, err)

	t.Run("with no pipeline runs", func(t *testing.T) {
		runs, count, err2 := orm.PipelineRuns(&jb.ID, 0, 10)
		require.NoError(t, err2)
		assert.Equal(t, count, 0)
		assert.Empty(t, runs)
	})

	t.Run("with a pipeline run", func(t *testing.T) {
		run := mustInsertPipelineRun(t, pipelineORM, jb)

		runs, count, err2 := orm.PipelineRuns(&jb.ID, 0, 10)
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
	var jb job.Job

	config := configtest.NewTestGeneralConfig(t)
	db := pgtest.NewSqlxDB(t)

	keyStore := cltest.NewKeyStore(t, db, config.Database())
	require.NoError(t, keyStore.OCR().Add(cltest.DefaultOCRKey))
	require.NoError(t, keyStore.P2P().Add(cltest.DefaultP2PKey))

	lggr := logger.TestLogger(t)
	pipelineORM := pipeline.NewORM(db, logger.TestLogger(t), config.Database(), config.JobPipeline().MaxSuccessfulRuns())
	bridgesORM := bridges.NewORM(db, lggr, config.Database())
	relayExtenders := evmtest.NewChainRelayExtenders(t, evmtest.TestChainOpts{DB: db, GeneralConfig: config, KeyStore: keyStore.Eth()})
	legacyChains := evmrelay.NewLegacyChainsFromRelayerExtenders(relayExtenders)
	orm := NewTestORM(t, db, legacyChains, pipelineORM, bridgesORM, keyStore, config.Database())

	_, address := cltest.MustInsertRandomKey(t, keyStore.Eth())

	jobs := make([]job.Job, 11)
	for j := 0; j < len(jobs); j++ {
		_, bridge := cltest.MustCreateBridge(t, db, cltest.BridgeOpts{}, config.Database())
		_, bridge2 := cltest.MustCreateBridge(t, db, cltest.BridgeOpts{}, config.Database())
		jobID := uuid.New().String()
		key, err := ethkey.NewV2()

		require.NoError(t, err)
		jb, err = ocr.ValidatedOracleSpecToml(legacyChains,
			testspecs.GenerateOCRSpec(testspecs.OCRSpecParams{
				JobID:              jobID,
				Name:               fmt.Sprintf("Job #%v", jobID),
				DS1BridgeName:      bridge.Name.String(),
				DS2BridgeName:      bridge2.Name.String(),
				TransmitterAddress: address.Hex(),
				ContractAddress:    key.Address.String(),
			}).Toml())

		require.NoError(t, err)

		err = orm.CreateJob(&jb)
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
		runIDs, err := orm.FindPipelineRunIDsByJobID(jb.ID, 0, 10)
		require.NoError(t, err)
		assert.Empty(t, runIDs)
	})

	t.Run("with a pipeline run", func(t *testing.T) {
		run := mustInsertPipelineRun(t, pipelineORM, jb)

		runIDs, err := orm.FindPipelineRunIDsByJobID(jb.ID, 0, 10)
		require.NoError(t, err)
		require.Len(t, runIDs, 1)

		assert.Equal(t, run.ID, runIDs[0])
	})

	// Internally these queries are batched by 1000, this tests case requiring concatenation
	//  of more than 1 batch
	t.Run("with batch concatenation limit 10", func(t *testing.T) {
		runIDs, err := orm.FindPipelineRunIDsByJobID(jobs[3].ID, 95, 10)
		require.NoError(t, err)
		require.Len(t, runIDs, 10)
		assert.Equal(t, int64(4*(len(jobs)-1)), runIDs[3]-runIDs[7])
	})

	// Internally these queries are batched by 1000, this tests case requiring concatenation
	//  of more than 1 batch
	t.Run("with batch concatenation limit 100", func(t *testing.T) {
		runIDs, err := orm.FindPipelineRunIDsByJobID(jobs[3].ID, 95, 100)
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
		runIDs, err := orm.FindPipelineRunIDsByJobID(jobs[3].ID, 0, 25)
		require.NoError(t, err)
		require.Len(t, runIDs, 25)
		assert.Equal(t, int64(16*(len(jobs)-1)), runIDs[7]-runIDs[23])
	})

	// Same as previous, but where there are fewer matching jobs than the limit
	t.Run("with first batch empty, under limit", func(t *testing.T) {
		runIDs, err := orm.FindPipelineRunIDsByJobID(jobs[3].ID, 143, 190)
		require.NoError(t, err)
		require.Len(t, runIDs, 107)
		assert.Equal(t, int64(16*(len(jobs)-1)), runIDs[7]-runIDs[23])
	})
}

func Test_FindPipelineRunsByIDs(t *testing.T) {
	t.Parallel()

	config := configtest.NewTestGeneralConfig(t)
	db := pgtest.NewSqlxDB(t)

	keyStore := cltest.NewKeyStore(t, db, config.Database())
	require.NoError(t, keyStore.OCR().Add(cltest.DefaultOCRKey))
	require.NoError(t, keyStore.P2P().Add(cltest.DefaultP2PKey))

	pipelineORM := pipeline.NewORM(db, logger.TestLogger(t), config.Database(), config.JobPipeline().MaxSuccessfulRuns())
	bridgesORM := bridges.NewORM(db, logger.TestLogger(t), config.Database())
	relayExtenders := evmtest.NewChainRelayExtenders(t, evmtest.TestChainOpts{DB: db, GeneralConfig: config, KeyStore: keyStore.Eth()})
	legacyChains := evmrelay.NewLegacyChainsFromRelayerExtenders(relayExtenders)
	orm := NewTestORM(t, db, legacyChains, pipelineORM, bridgesORM, keyStore, config.Database())

	_, bridge := cltest.MustCreateBridge(t, db, cltest.BridgeOpts{}, config.Database())
	_, bridge2 := cltest.MustCreateBridge(t, db, cltest.BridgeOpts{}, config.Database())

	externalJobID := uuid.New()
	_, address := cltest.MustInsertRandomKey(t, keyStore.Eth())
	jb, err := ocr.ValidatedOracleSpecToml(legacyChains,
		testspecs.GenerateOCRSpec(testspecs.OCRSpecParams{
			JobID:              externalJobID.String(),
			TransmitterAddress: address.Hex(),
			DS1BridgeName:      bridge.Name.String(),
			DS2BridgeName:      bridge2.Name.String(),
		}).Toml(),
	)
	require.NoError(t, err)

	err = orm.CreateJob(&jb)
	require.NoError(t, err)

	t.Run("with no pipeline runs", func(t *testing.T) {
		runs, err2 := orm.FindPipelineRunsByIDs([]int64{-1})
		require.NoError(t, err2)
		assert.Empty(t, runs)
	})

	t.Run("with a pipeline run", func(t *testing.T) {
		run := mustInsertPipelineRun(t, pipelineORM, jb)

		actual, err2 := orm.FindPipelineRunsByIDs([]int64{run.ID})
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

	config := configtest.NewTestGeneralConfig(t)
	db := pgtest.NewSqlxDB(t)

	keyStore := cltest.NewKeyStore(t, db, config.Database())
	err := keyStore.OCR().Add(cltest.DefaultOCRKey)
	require.NoError(t, err)

	pipelineORM := pipeline.NewORM(db, logger.TestLogger(t), config.Database(), config.JobPipeline().MaxSuccessfulRuns())
	bridgesORM := bridges.NewORM(db, logger.TestLogger(t), config.Database())
	relayExtenders := evmtest.NewChainRelayExtenders(t, evmtest.TestChainOpts{DB: db, GeneralConfig: config, KeyStore: keyStore.Eth()})
	legacyChains := evmrelay.NewLegacyChainsFromRelayerExtenders(relayExtenders)
	orm := NewTestORM(t, db, legacyChains, pipelineORM, bridgesORM, keyStore, config.Database())

	jb, err := directrequest.ValidatedDirectRequestSpec(testspecs.GetDirectRequestSpec())
	require.NoError(t, err)

	err = orm.CreateJob(&jb)
	require.NoError(t, err)

	t.Run("with no pipeline run", func(t *testing.T) {
		run, err2 := orm.FindPipelineRunByID(-1)
		assert.Equal(t, run, pipeline.Run{})
		require.ErrorIs(t, err2, sql.ErrNoRows)
	})

	t.Run("with a pipeline run", func(t *testing.T) {
		run := mustInsertPipelineRun(t, pipelineORM, jb)

		actual, err2 := orm.FindPipelineRunByID(run.ID)
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

	config := configtest.NewTestGeneralConfig(t)
	db := pgtest.NewSqlxDB(t)

	keyStore := cltest.NewKeyStore(t, db, config.Database())
	err := keyStore.OCR().Add(cltest.DefaultOCRKey)
	require.NoError(t, err)

	pipelineORM := pipeline.NewORM(db, logger.TestLogger(t), config.Database(), config.JobPipeline().MaxSuccessfulRuns())
	bridgesORM := bridges.NewORM(db, logger.TestLogger(t), config.Database())
	relayExtenders := evmtest.NewChainRelayExtenders(t, evmtest.TestChainOpts{DB: db, GeneralConfig: config, KeyStore: keyStore.Eth()})
	legacyChains := evmrelay.NewLegacyChainsFromRelayerExtenders(relayExtenders)
	orm := NewTestORM(t, db, legacyChains, pipelineORM, bridgesORM, keyStore, config.Database())

	jb, err := directrequest.ValidatedDirectRequestSpec(testspecs.GetDirectRequestSpec())
	require.NoError(t, err)

	err = orm.CreateJob(&jb)
	require.NoError(t, err)
	var jobSpec job.Job
	err = db.Get(&jobSpec, "SELECT * FROM jobs")
	require.NoError(t, err)

	ocrSpecError1 := "ocr spec 1 errored"
	ocrSpecError2 := "ocr spec 2 errored"
	require.NoError(t, orm.RecordError(jobSpec.ID, ocrSpecError1))
	require.NoError(t, orm.RecordError(jobSpec.ID, ocrSpecError2))

	jb, err = orm.FindJobWithoutSpecErrors(jobSpec.ID)
	require.NoError(t, err)
	jbWithErrors, err := orm.FindJobTx(jobSpec.ID)
	require.NoError(t, err)

	assert.Equal(t, len(jb.JobSpecErrors), 0)
	assert.Equal(t, len(jbWithErrors.JobSpecErrors), 2)
}

func Test_FindSpecErrorsByJobIDs(t *testing.T) {
	t.Parallel()

	config := configtest.NewTestGeneralConfig(t)
	db := pgtest.NewSqlxDB(t)

	keyStore := cltest.NewKeyStore(t, db, config.Database())
	err := keyStore.OCR().Add(cltest.DefaultOCRKey)
	require.NoError(t, err)

	pipelineORM := pipeline.NewORM(db, logger.TestLogger(t), config.Database(), config.JobPipeline().MaxSuccessfulRuns())
	bridgesORM := bridges.NewORM(db, logger.TestLogger(t), config.Database())
	relayExtenders := evmtest.NewChainRelayExtenders(t, evmtest.TestChainOpts{DB: db, GeneralConfig: config, KeyStore: keyStore.Eth()})
	legacyChains := evmrelay.NewLegacyChainsFromRelayerExtenders(relayExtenders)
	orm := NewTestORM(t, db, legacyChains, pipelineORM, bridgesORM, keyStore, config.Database())

	jb, err := directrequest.ValidatedDirectRequestSpec(testspecs.GetDirectRequestSpec())
	require.NoError(t, err)

	err = orm.CreateJob(&jb)
	require.NoError(t, err)
	var jobSpec job.Job
	err = db.Get(&jobSpec, "SELECT * FROM jobs")
	require.NoError(t, err)

	ocrSpecError1 := "ocr spec 1 errored"
	ocrSpecError2 := "ocr spec 2 errored"
	require.NoError(t, orm.RecordError(jobSpec.ID, ocrSpecError1))
	require.NoError(t, orm.RecordError(jobSpec.ID, ocrSpecError2))

	specErrs, err := orm.FindSpecErrorsByJobIDs([]int32{jobSpec.ID})
	require.NoError(t, err)

	assert.Equal(t, len(specErrs), 2)
}

func Test_CountPipelineRunsByJobID(t *testing.T) {
	t.Parallel()

	config := configtest.NewTestGeneralConfig(t)
	db := pgtest.NewSqlxDB(t)

	keyStore := cltest.NewKeyStore(t, db, config.Database())
	require.NoError(t, keyStore.OCR().Add(cltest.DefaultOCRKey))
	require.NoError(t, keyStore.P2P().Add(cltest.DefaultP2PKey))

	pipelineORM := pipeline.NewORM(db, logger.TestLogger(t), config.Database(), config.JobPipeline().MaxSuccessfulRuns())
	bridgesORM := bridges.NewORM(db, logger.TestLogger(t), config.Database())
	relayExtenders := evmtest.NewChainRelayExtenders(t, evmtest.TestChainOpts{DB: db, GeneralConfig: config, KeyStore: keyStore.Eth()})
	legacyChains := evmrelay.NewLegacyChainsFromRelayerExtenders(relayExtenders)
	orm := NewTestORM(t, db, legacyChains, pipelineORM, bridgesORM, keyStore, config.Database())

	_, bridge := cltest.MustCreateBridge(t, db, cltest.BridgeOpts{}, config.Database())
	_, bridge2 := cltest.MustCreateBridge(t, db, cltest.BridgeOpts{}, config.Database())

	externalJobID := uuid.New()
	_, address := cltest.MustInsertRandomKey(t, keyStore.Eth())
	jb, err := ocr.ValidatedOracleSpecToml(legacyChains,
		testspecs.GenerateOCRSpec(testspecs.OCRSpecParams{
			JobID:              externalJobID.String(),
			TransmitterAddress: address.Hex(),
			DS1BridgeName:      bridge.Name.String(),
			DS2BridgeName:      bridge2.Name.String(),
		}).Toml(),
	)
	require.NoError(t, err)

	err = orm.CreateJob(&jb)
	require.NoError(t, err)

	t.Run("with no pipeline runs", func(t *testing.T) {
		count, err2 := orm.CountPipelineRunsByJobID(jb.ID)
		require.NoError(t, err2)
		assert.Equal(t, int32(0), count)
	})

	t.Run("with a pipeline run", func(t *testing.T) {
		mustInsertPipelineRun(t, pipelineORM, jb)

		count, err2 := orm.CountPipelineRunsByJobID(jb.ID)
		require.NoError(t, err2)
		require.Equal(t, int32(1), count)
	})
}

func mustInsertPipelineRun(t *testing.T, orm pipeline.ORM, j job.Job) pipeline.Run {
	t.Helper()

	run := pipeline.Run{
		PipelineSpecID: j.PipelineSpecID,
		State:          pipeline.RunStatusRunning,
		Outputs:        pipeline.JSONSerializable{Valid: false},
		AllErrors:      pipeline.RunErrors{},
		CreatedAt:      time.Now(),
		FinishedAt:     null.Time{},
	}
	err := orm.CreateRun(&run)
	require.NoError(t, err)
	return run
}
