package job_test

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/lib/pq"
	"github.com/pelletier/go-toml"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink/core/bridges"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/evmtest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/blockhashstore"
	"github.com/smartcontractkit/chainlink/core/services/directrequest"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/keeper"
	"github.com/smartcontractkit/chainlink/core/services/ocr"
	"github.com/smartcontractkit/chainlink/core/services/ocrbootstrap"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/core/services/vrf"
	"github.com/smartcontractkit/chainlink/core/services/webhook"
	"github.com/smartcontractkit/chainlink/core/testdata/testspecs"
)

func TestORM(t *testing.T) {
	t.Parallel()
	config := cltest.NewTestGeneralConfig(t)
	db := pgtest.NewSqlxDB(t)
	keyStore := cltest.NewKeyStore(t, db, config)
	ethKeyStore := keyStore.Eth()

	require.NoError(t, keyStore.OCR().Add(cltest.DefaultOCRKey))
	require.NoError(t, keyStore.P2P().Add(cltest.DefaultP2PKey))

	pipelineORM := pipeline.NewORM(db, logger.TestLogger(t), config)

	cc := evmtest.NewChainSet(t, evmtest.TestChainOpts{DB: db, GeneralConfig: config})
	orm := job.NewTestORM(t, db, cc, pipelineORM, keyStore, config)
	borm := bridges.NewORM(db, logger.TestLogger(t), config)

	_, bridge := cltest.MustCreateBridge(t, db, cltest.BridgeOpts{}, config)
	_, bridge2 := cltest.MustCreateBridge(t, db, cltest.BridgeOpts{}, config)
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
		tree, err := toml.LoadFile("../../testdata/tomlspecs/direct-request-spec.toml")
		require.NoError(t, err)
		jb, err := directrequest.ValidatedDirectRequestSpec(tree.String())
		require.NoError(t, err)
		err = orm.CreateJob(&jb)
		require.NoError(t, err)
	})

	t.Run("creates webhook specs along with external_initiator_webhook_specs", func(t *testing.T) {
		eiFoo := cltest.MustInsertExternalInitiator(t, borm)
		eiBar := cltest.MustInsertExternalInitiator(t, borm)

		eiWS := []webhook.TOMLWebhookSpecExternalInitiator{
			{Name: eiFoo.Name, Spec: cltest.JSONFromString(t, `{}`)},
			{Name: eiBar.Name, Spec: cltest.JSONFromString(t, `{"bar": 1}`)},
		}
		eim := webhook.NewExternalInitiatorManager(db, nil, logger.TestLogger(t), config)
		jb, err := webhook.ValidatedWebhookSpec(testspecs.GenerateWebhookSpec(testspecs.WebhookSpecParams{ExternalInitiators: eiWS}).Toml(), eim)
		require.NoError(t, err)

		err = orm.CreateJob(&jb)
		require.NoError(t, err)

		cltest.AssertCount(t, db, "external_initiator_webhook_specs", 2)
	})

	t.Run("it creates and deletes records for blockhash store jobs", func(t *testing.T) {
		jb, err := blockhashstore.ValidatedSpec(
			testspecs.GenerateBlockhashStoreSpec(testspecs.BlockhashStoreSpecParams{}).Toml())
		require.NoError(t, err)

		err = orm.CreateJob(&jb)
		require.NoError(t, err)
		savedJob, err := orm.FindJob(context.Background(), jb.ID)
		require.NoError(t, err)
		require.Equal(t, jb.ID, savedJob.ID)
		require.Equal(t, jb.Type, savedJob.Type)
		require.Equal(t, jb.BlockhashStoreSpec.ID, savedJob.BlockhashStoreSpec.ID)
		require.Equal(t, jb.BlockhashStoreSpec.CoordinatorV1Address, savedJob.BlockhashStoreSpec.CoordinatorV1Address)
		require.Equal(t, jb.BlockhashStoreSpec.CoordinatorV2Address, savedJob.BlockhashStoreSpec.CoordinatorV2Address)
		require.Equal(t, jb.BlockhashStoreSpec.WaitBlocks, savedJob.BlockhashStoreSpec.WaitBlocks)
		require.Equal(t, jb.BlockhashStoreSpec.LookbackBlocks, savedJob.BlockhashStoreSpec.LookbackBlocks)
		require.Equal(t, jb.BlockhashStoreSpec.BlockhashStoreAddress, savedJob.BlockhashStoreSpec.BlockhashStoreAddress)
		require.Equal(t, jb.BlockhashStoreSpec.PollPeriod, savedJob.BlockhashStoreSpec.PollPeriod)
		require.Equal(t, jb.BlockhashStoreSpec.RunTimeout, savedJob.BlockhashStoreSpec.RunTimeout)
		require.Equal(t, jb.BlockhashStoreSpec.EVMChainID, savedJob.BlockhashStoreSpec.EVMChainID)
		require.Equal(t, jb.BlockhashStoreSpec.FromAddress, savedJob.BlockhashStoreSpec.FromAddress)
		err = orm.DeleteJob(jb.ID)
		require.NoError(t, err)
		_, err = orm.FindJob(context.Background(), jb.ID)
		require.Error(t, err)
	})
}

func TestORM_DeleteJob_DeletesAssociatedRecords(t *testing.T) {
	t.Parallel()
	config := evmtest.NewChainScopedConfig(t, cltest.NewTestGeneralConfig(t))
	db := pgtest.NewSqlxDB(t)
	keyStore := cltest.NewKeyStore(t, db, config)
	require.NoError(t, keyStore.OCR().Add(cltest.DefaultOCRKey))
	require.NoError(t, keyStore.P2P().Add(cltest.DefaultP2PKey))

	pipelineORM := pipeline.NewORM(db, logger.TestLogger(t), config)
	cc := evmtest.NewChainSet(t, evmtest.TestChainOpts{DB: db, GeneralConfig: config})
	jobORM := job.NewTestORM(t, db, cc, pipelineORM, keyStore, config)
	korm := keeper.NewORM(db, logger.TestLogger(t), nil, nil)

	t.Run("it deletes records for offchainreporting jobs", func(t *testing.T) {
		_, bridge := cltest.MustCreateBridge(t, db, cltest.BridgeOpts{}, config)
		_, bridge2 := cltest.MustCreateBridge(t, db, cltest.BridgeOpts{}, config)

		_, address := cltest.MustInsertRandomKey(t, keyStore.Eth())
		jb, err := ocr.ValidatedOracleSpecToml(cc, testspecs.GenerateOCRSpec(testspecs.OCRSpecParams{
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
		cltest.MustInsertUpkeepForRegistry(t, db, config, registry)

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
		jb, err := vrf.ValidatedVRFSpec(testspecs.GenerateVRFSpec(testspecs.VRFSpecParams{PublicKey: pk.String()}).Toml())
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
		ei := cltest.MustInsertExternalInitiator(t, bridges.NewORM(db, logger.TestLogger(t), config))
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
		ei := cltest.MustInsertExternalInitiator(t, bridges.NewORM(db, logger.TestLogger(t), config))
		_, webhookSpec := cltest.MustInsertWebhookSpec(t, db)
		_, err := db.Exec(`INSERT INTO external_initiator_webhook_specs (external_initiator_id, webhook_spec_id, spec) VALUES ($1,$2,$3)`, ei.ID, webhookSpec.ID, `{"ei": "foo", "name": "webhookSpecTwoEIs"}`)
		require.NoError(t, err)

		_, err = db.Exec(`DELETE FROM external_initiators`)
		require.EqualError(t, err, "ERROR: update or delete on table \"external_initiators\" violates foreign key constraint \"external_initiator_webhook_specs_external_initiator_id_fkey\" on table \"external_initiator_webhook_specs\" (SQLSTATE 23503)")
	})
}

func TestORM_CreateJob_VRFV2(t *testing.T) {
	config := evmtest.NewChainScopedConfig(t, cltest.NewTestGeneralConfig(t))
	db := pgtest.NewSqlxDB(t)
	keyStore := cltest.NewKeyStore(t, db, config)
	require.NoError(t, keyStore.OCR().Add(cltest.DefaultOCRKey))

	pipelineORM := pipeline.NewORM(db, logger.TestLogger(t), config)
	cc := evmtest.NewChainSet(t, evmtest.TestChainOpts{DB: db, GeneralConfig: config})
	jobORM := job.NewTestORM(t, db, cc, pipelineORM, keyStore, config)

	fromAddresses := []string{cltest.NewEIP55Address().String(), cltest.NewEIP55Address().String()}
	jb, err := vrf.ValidatedVRFSpec(testspecs.GenerateVRFSpec(
		testspecs.VRFSpecParams{
			RequestedConfsDelay: 10,
			FromAddresses:       fromAddresses,
			ChunkSize:           25,
			BackoffInitialDelay: time.Minute,
			BackoffMaxDelay:     time.Hour}).
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
	var fa pq.ByteaArray
	require.NoError(t, db.Get(&fa, `SELECT from_addresses FROM vrf_specs LIMIT 1`))
	var actual []string
	for _, b := range fa {
		actual = append(actual, common.BytesToAddress(b).String())
	}
	require.ElementsMatch(t, fromAddresses, actual)
	require.NoError(t, jobORM.DeleteJob(jb.ID))
	cltest.AssertCount(t, db, "vrf_specs", 0)
	cltest.AssertCount(t, db, "jobs", 0)

	jb, err = vrf.ValidatedVRFSpec(testspecs.GenerateVRFSpec(testspecs.VRFSpecParams{RequestTimeout: 1 * time.Hour}).Toml())
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
	config := evmtest.NewChainScopedConfig(t, cltest.NewTestGeneralConfig(t))
	db := pgtest.NewSqlxDB(t)
	keyStore := cltest.NewKeyStore(t, db, config)
	keyStore.OCR().Add(cltest.DefaultOCRKey)

	pipelineORM := pipeline.NewORM(db, logger.TestLogger(t), config)
	cc := evmtest.NewChainSet(t, evmtest.TestChainOpts{DB: db, GeneralConfig: config})
	jobORM := job.NewTestORM(t, db, cc, pipelineORM, keyStore, config)

	jb, err := ocrbootstrap.ValidatedBootstrapSpecToml(testspecs.OCRBootstrapSpec)
	require.NoError(t, err)

	err = jobORM.CreateJob(&jb)
	require.NoError(t, err)
	cltest.AssertCount(t, db, "bootstrap_specs", 1)
	cltest.AssertCount(t, db, "jobs", 1)
	var relay string
	require.NoError(t, db.Get(&relay, `SELECT relay FROM bootstrap_specs LIMIT 1`))
	require.Equal(t, "evm", relay)

	jobORM.DeleteJob(jb.ID)
	cltest.AssertCount(t, db, "bootstrap_specs", 0)
	cltest.AssertCount(t, db, "jobs", 0)
}

func TestORM_CreateJob_OCR_DuplicatedContractAddress(t *testing.T) {
	config := evmtest.NewChainScopedConfig(t, cltest.NewTestGeneralConfig(t))
	db := pgtest.NewSqlxDB(t)
	keyStore := cltest.NewKeyStore(t, db, config)
	keyStore.OCR().Add(cltest.DefaultOCRKey)

	pipelineORM := pipeline.NewORM(db, logger.TestLogger(t), config)
	cc := evmtest.NewChainSet(t, evmtest.TestChainOpts{DB: db, GeneralConfig: config})
	jobORM := job.NewTestORM(t, db, cc, pipelineORM, keyStore, config)

	chain, err := cc.Default()
	require.NoError(t, err)

	_, address := cltest.MustInsertRandomKey(t, keyStore.Eth())
	_, bridge := cltest.MustCreateBridge(t, db, cltest.BridgeOpts{}, config)
	_, bridge2 := cltest.MustCreateBridge(t, db, cltest.BridgeOpts{}, config)

	t.Run("with the default chain id", func(t *testing.T) {
		spec := testspecs.GenerateOCRSpec(testspecs.OCRSpecParams{
			DS1BridgeName:      bridge.Name.String(),
			DS2BridgeName:      bridge2.Name.String(),
			TransmitterAddress: address.Hex(),
		})
		jb, err := ocr.ValidatedOracleSpecToml(cc, spec.Toml())
		require.NoError(t, err)

		err = jobORM.CreateJob(&jb)
		require.NoError(t, err)
		cltest.AssertCount(t, db, "ocr_oracle_specs", 1)
		cltest.AssertCount(t, db, "jobs", 1)

		spec2 := testspecs.GenerateOCRSpec(testspecs.OCRSpecParams{
			EVMChainID:         chain.ID().String(),
			DS1BridgeName:      bridge.Name.String(),
			DS2BridgeName:      bridge2.Name.String(),
			TransmitterAddress: address.Hex(),
		})
		jb2, err := ocr.ValidatedOracleSpecToml(cc, spec2.Toml())
		require.NoError(t, err)

		err = jobORM.CreateJob(&jb2)
		require.Error(t, err)
		assert.Equal(t, err.Error(), fmt.Sprintf("CreateJobFailed: a job with contract address %s already exists for chain ID %d", jb2.OCROracleSpec.ContractAddress, jb2.OCROracleSpec.EVMChainID.ToInt()))
	})

	t.Run("with a set chain id", func(t *testing.T) {
		externalJobID := uuid.NullUUID{UUID: uuid.NewV4(), Valid: true}
		_, contractAddress := cltest.MustInsertRandomKey(t, keyStore.Eth())

		spec := testspecs.GenerateOCRSpec(testspecs.OCRSpecParams{
			EVMChainID:         chain.ID().String(),
			DS1BridgeName:      bridge.Name.String(),
			DS2BridgeName:      bridge2.Name.String(),
			TransmitterAddress: address.Hex(),
			ContractAddress:    contractAddress.Hex(),
			JobID:              externalJobID.UUID.String(),
			Name:               "with a chain id",
		})

		jb, err := ocr.ValidatedOracleSpecToml(cc, spec.Toml())
		require.NoError(t, err)

		err = jobORM.CreateJob(&jb)
		require.NoError(t, err)

		spec2 := testspecs.GenerateOCRSpec(testspecs.OCRSpecParams{
			EVMChainID:         chain.ID().String(),
			DS1BridgeName:      bridge.Name.String(),
			DS2BridgeName:      bridge2.Name.String(),
			TransmitterAddress: address.Hex(),
			ContractAddress:    contractAddress.Hex(),
			JobID:              externalJobID.UUID.String(),
			Name:               "with a chain id 2",
		})
		jb2, err := ocr.ValidatedOracleSpecToml(cc, spec2.Toml())
		require.NoError(t, err)

		err = jobORM.CreateJob(&jb2)
		require.Error(t, err)
		assert.Equal(t, err.Error(), fmt.Sprintf("CreateJobFailed: a job with contract address %s already exists for chain ID %d", jb2.OCROracleSpec.ContractAddress, chain.ID()))
	})
}

func Test_FindJobs(t *testing.T) {
	t.Parallel()

	config := cltest.NewTestGeneralConfig(t)
	db := pgtest.NewSqlxDB(t)
	keyStore := cltest.NewKeyStore(t, db, config)
	require.NoError(t, keyStore.OCR().Add(cltest.DefaultOCRKey))
	require.NoError(t, keyStore.P2P().Add(cltest.DefaultP2PKey))

	pipelineORM := pipeline.NewORM(db, logger.TestLogger(t), config)
	cc := evmtest.NewChainSet(t, evmtest.TestChainOpts{DB: db, GeneralConfig: config})
	orm := job.NewTestORM(t, db, cc, pipelineORM, keyStore, config)

	_, bridge := cltest.MustCreateBridge(t, db, cltest.BridgeOpts{}, config)
	_, bridge2 := cltest.MustCreateBridge(t, db, cltest.BridgeOpts{}, config)

	_, address := cltest.MustInsertRandomKey(t, keyStore.Eth())
	jb1, err := ocr.ValidatedOracleSpecToml(cc,
		testspecs.GenerateOCRSpec(testspecs.OCRSpecParams{
			JobID:              uuid.NewV4().String(),
			TransmitterAddress: address.Hex(),
			DS1BridgeName:      bridge.Name.String(),
			DS2BridgeName:      bridge2.Name.String(),
		}).Toml(),
	)
	require.NoError(t, err)

	err = orm.CreateJob(&jb1)
	require.NoError(t, err)

	jb2, err := directrequest.ValidatedDirectRequestSpec(
		testspecs.DirectRequestSpec,
	)
	require.NoError(t, err)

	err = orm.CreateJob(&jb2)
	require.NoError(t, err)

	t.Run("jobs are ordered by latest first", func(t *testing.T) {
		jobs, count, err := orm.FindJobs(0, 2)
		require.NoError(t, err)
		require.Len(t, jobs, 2)
		assert.Equal(t, count, 2)

		expectedJobs := []job.Job{jb2, jb1}

		for i, exp := range expectedJobs {
			assert.Equal(t, exp.ID, jobs[i].ID)
		}
	})

	t.Run("jobs respect pagination", func(t *testing.T) {
		jobs, count, err := orm.FindJobs(0, 1)
		require.NoError(t, err)
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

	config := cltest.NewTestGeneralConfig(t)
	db := pgtest.NewSqlxDB(t)
	keyStore := cltest.NewKeyStore(t, db, config)
	require.NoError(t, keyStore.OCR().Add(cltest.DefaultOCRKey))
	require.NoError(t, keyStore.P2P().Add(cltest.DefaultP2PKey))

	pipelineORM := pipeline.NewORM(db, logger.TestLogger(t), config)
	cc := evmtest.NewChainSet(t, evmtest.TestChainOpts{DB: db, GeneralConfig: config})
	orm := job.NewTestORM(t, db, cc, pipelineORM, keyStore, config)

	_, bridge := cltest.MustCreateBridge(t, db, cltest.BridgeOpts{}, config)
	_, bridge2 := cltest.MustCreateBridge(t, db, cltest.BridgeOpts{}, config)

	externalJobID := uuid.NewV4()
	_, address := cltest.MustInsertRandomKey(t, keyStore.Eth())
	job, err := ocr.ValidatedOracleSpecToml(cc,
		testspecs.GenerateOCRSpec(testspecs.OCRSpecParams{
			JobID:              externalJobID.String(),
			TransmitterAddress: address.Hex(),
			DS1BridgeName:      bridge.Name.String(),
			DS2BridgeName:      bridge2.Name.String(),
		}).Toml(),
	)
	require.NoError(t, err)

	err = orm.CreateJob(&job)
	require.NoError(t, err)

	t.Run("by id", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		jb, err := orm.FindJob(ctx, job.ID)
		require.NoError(t, err)

		assert.Equal(t, jb.ID, job.ID)
		assert.Equal(t, jb.Name, job.Name)

		require.Greater(t, jb.PipelineSpecID, int32(0))
		require.NotNil(t, jb.PipelineSpec)
		require.NotNil(t, jb.OCROracleSpecID)
		require.NotNil(t, jb.OCROracleSpec)
	})

	t.Run("by external job id", func(t *testing.T) {
		jb, err := orm.FindJobByExternalJobID(externalJobID)
		require.NoError(t, err)

		assert.Equal(t, jb.ID, job.ID)
		assert.Equal(t, jb.Name, job.Name)

		require.Greater(t, jb.PipelineSpecID, int32(0))
		require.NotNil(t, jb.PipelineSpec)
		require.NotNil(t, jb.OCROracleSpecID)
		require.NotNil(t, jb.OCROracleSpec)
	})

	t.Run("by address", func(t *testing.T) {
		jbID, err := orm.FindJobIDByAddress(job.OCROracleSpec.ContractAddress)
		require.NoError(t, err)

		assert.Equal(t, job.ID, jbID)

		_, err = orm.FindJobIDByAddress("not-existing")
		require.Error(t, err)
		require.ErrorIs(t, err, sql.ErrNoRows)
	})
}

func Test_FindJobsByPipelineSpecIDs(t *testing.T) {
	t.Parallel()

	config := cltest.NewTestGeneralConfig(t)
	db := pgtest.NewSqlxDB(t)
	keyStore := cltest.NewKeyStore(t, db, config)
	require.NoError(t, keyStore.OCR().Add(cltest.DefaultOCRKey))

	pipelineORM := pipeline.NewORM(db, logger.TestLogger(t), config)
	cc := evmtest.NewChainSet(t, evmtest.TestChainOpts{DB: db, GeneralConfig: config})
	orm := job.NewTestORM(t, db, cc, pipelineORM, keyStore, config)

	jb, err := directrequest.ValidatedDirectRequestSpec(testspecs.DirectRequestSpec)
	require.NoError(t, err)

	err = orm.CreateJob(&jb)
	require.NoError(t, err)

	t.Run("with jobs", func(t *testing.T) {
		jbs, err := orm.FindJobsByPipelineSpecIDs([]int32{jb.PipelineSpecID})
		require.NoError(t, err)
		assert.Len(t, jbs, 1)

		assert.Equal(t, jb.ID, jbs[0].ID)
		assert.Equal(t, jb.Name, jbs[0].Name)

		require.Greater(t, jbs[0].PipelineSpecID, int32(0))
		require.Equal(t, jb.PipelineSpecID, jbs[0].PipelineSpecID)
		require.NotNil(t, jbs[0].PipelineSpec)
	})

	t.Run("without jobs", func(t *testing.T) {
		jbs, err := orm.FindJobsByPipelineSpecIDs([]int32{-1})
		require.NoError(t, err)
		assert.Len(t, jbs, 0)
	})
}

func Test_FindPipelineRuns(t *testing.T) {
	t.Parallel()

	config := cltest.NewTestGeneralConfig(t)
	db := pgtest.NewSqlxDB(t)
	keyStore := cltest.NewKeyStore(t, db, config)
	require.NoError(t, keyStore.OCR().Add(cltest.DefaultOCRKey))
	require.NoError(t, keyStore.P2P().Add(cltest.DefaultP2PKey))

	pipelineORM := pipeline.NewORM(db, logger.TestLogger(t), config)
	cc := evmtest.NewChainSet(t, evmtest.TestChainOpts{DB: db, GeneralConfig: config})
	orm := job.NewTestORM(t, db, cc, pipelineORM, keyStore, config)

	_, bridge := cltest.MustCreateBridge(t, db, cltest.BridgeOpts{}, config)
	_, bridge2 := cltest.MustCreateBridge(t, db, cltest.BridgeOpts{}, config)

	externalJobID := uuid.NewV4()
	_, address := cltest.MustInsertRandomKey(t, keyStore.Eth())
	jb, err := ocr.ValidatedOracleSpecToml(cc,
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
		runs, count, err := orm.PipelineRuns(nil, 0, 10)
		require.NoError(t, err)
		assert.Equal(t, count, 0)
		assert.Empty(t, runs)
	})

	t.Run("with a pipeline run", func(t *testing.T) {
		run := mustInsertPipelineRun(t, pipelineORM, jb)

		runs, count, err := orm.PipelineRuns(nil, 0, 10)
		require.NoError(t, err)

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

	config := cltest.NewTestGeneralConfig(t)
	db := pgtest.NewSqlxDB(t)

	keyStore := cltest.NewKeyStore(t, db, config)
	require.NoError(t, keyStore.OCR().Add(cltest.DefaultOCRKey))
	require.NoError(t, keyStore.P2P().Add(cltest.DefaultP2PKey))

	pipelineORM := pipeline.NewORM(db, logger.TestLogger(t), config)
	cc := evmtest.NewChainSet(t, evmtest.TestChainOpts{DB: db, GeneralConfig: config})
	orm := job.NewTestORM(t, db, cc, pipelineORM, keyStore, config)

	_, bridge := cltest.MustCreateBridge(t, db, cltest.BridgeOpts{}, config)
	_, bridge2 := cltest.MustCreateBridge(t, db, cltest.BridgeOpts{}, config)

	externalJobID := uuid.NewV4()
	_, address := cltest.MustInsertRandomKey(t, keyStore.Eth())
	jb, err := ocr.ValidatedOracleSpecToml(cc,
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
		runs, count, err := orm.PipelineRuns(&jb.ID, 0, 10)
		require.NoError(t, err)
		assert.Equal(t, count, 0)
		assert.Empty(t, runs)
	})

	t.Run("with a pipeline run", func(t *testing.T) {
		run := mustInsertPipelineRun(t, pipelineORM, jb)

		runs, count, err := orm.PipelineRuns(&jb.ID, 0, 10)
		require.NoError(t, err)

		assert.Equal(t, count, 1)
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
	t.Parallel()

	config := cltest.NewTestGeneralConfig(t)
	db := pgtest.NewSqlxDB(t)

	keyStore := cltest.NewKeyStore(t, db, config)
	require.NoError(t, keyStore.OCR().Add(cltest.DefaultOCRKey))
	require.NoError(t, keyStore.P2P().Add(cltest.DefaultP2PKey))

	pipelineORM := pipeline.NewORM(db, logger.TestLogger(t), config)
	cc := evmtest.NewChainSet(t, evmtest.TestChainOpts{DB: db, GeneralConfig: config})
	orm := job.NewTestORM(t, db, cc, pipelineORM, keyStore, config)

	_, bridge := cltest.MustCreateBridge(t, db, cltest.BridgeOpts{}, config)
	_, bridge2 := cltest.MustCreateBridge(t, db, cltest.BridgeOpts{}, config)

	externalJobID := uuid.NewV4()
	_, address := cltest.MustInsertRandomKey(t, keyStore.Eth())
	jb, err := ocr.ValidatedOracleSpecToml(cc,
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
}

func Test_FindPipelineRunsByIDs(t *testing.T) {
	t.Parallel()

	config := cltest.NewTestGeneralConfig(t)
	db := pgtest.NewSqlxDB(t)

	keyStore := cltest.NewKeyStore(t, db, config)
	require.NoError(t, keyStore.OCR().Add(cltest.DefaultOCRKey))
	require.NoError(t, keyStore.P2P().Add(cltest.DefaultP2PKey))

	pipelineORM := pipeline.NewORM(db, logger.TestLogger(t), config)
	cc := evmtest.NewChainSet(t, evmtest.TestChainOpts{DB: db, GeneralConfig: config})
	orm := job.NewTestORM(t, db, cc, pipelineORM, keyStore, config)

	_, bridge := cltest.MustCreateBridge(t, db, cltest.BridgeOpts{}, config)
	_, bridge2 := cltest.MustCreateBridge(t, db, cltest.BridgeOpts{}, config)

	externalJobID := uuid.NewV4()
	_, address := cltest.MustInsertRandomKey(t, keyStore.Eth())
	jb, err := ocr.ValidatedOracleSpecToml(cc,
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
		runs, err := orm.FindPipelineRunsByIDs([]int64{-1})
		require.NoError(t, err)
		assert.Empty(t, runs)
	})

	t.Run("with a pipeline run", func(t *testing.T) {
		run := mustInsertPipelineRun(t, pipelineORM, jb)

		actual, err := orm.FindPipelineRunsByIDs([]int64{run.ID})
		require.NoError(t, err)
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

	config := cltest.NewTestGeneralConfig(t)
	db := pgtest.NewSqlxDB(t)

	keyStore := cltest.NewKeyStore(t, db, config)
	err := keyStore.OCR().Add(cltest.DefaultOCRKey)
	require.NoError(t, err)

	pipelineORM := pipeline.NewORM(db, logger.TestLogger(t), config)
	cc := evmtest.NewChainSet(t, evmtest.TestChainOpts{DB: db, GeneralConfig: config})
	orm := job.NewTestORM(t, db, cc, pipelineORM, keyStore, config)

	jb, err := directrequest.ValidatedDirectRequestSpec(testspecs.DirectRequestSpec)
	require.NoError(t, err)

	err = orm.CreateJob(&jb)
	require.NoError(t, err)

	t.Run("with no pipeline run", func(t *testing.T) {
		run, err := orm.FindPipelineRunByID(-1)
		assert.Equal(t, run, pipeline.Run{})
		require.ErrorIs(t, err, sql.ErrNoRows)
	})

	t.Run("with a pipeline run", func(t *testing.T) {
		run := mustInsertPipelineRun(t, pipelineORM, jb)

		actual, err := orm.FindPipelineRunByID(run.ID)
		require.NoError(t, err)

		actualRun := actual
		// Test pipeline run fields
		assert.Equal(t, run.State, actualRun.State)
		assert.Equal(t, run.PipelineSpecID, actualRun.PipelineSpecID)

		// Test preloaded pipeline spec
		assert.Equal(t, jb.PipelineSpec.ID, actualRun.PipelineSpec.ID)
		assert.Equal(t, jb.ID, actualRun.PipelineSpec.JobID)
	})
}

func Test_CountPipelineRunsByJobID(t *testing.T) {
	t.Parallel()

	config := cltest.NewTestGeneralConfig(t)
	db := pgtest.NewSqlxDB(t)

	keyStore := cltest.NewKeyStore(t, db, config)
	require.NoError(t, keyStore.OCR().Add(cltest.DefaultOCRKey))
	require.NoError(t, keyStore.P2P().Add(cltest.DefaultP2PKey))

	pipelineORM := pipeline.NewORM(db, logger.TestLogger(t), config)
	cc := evmtest.NewChainSet(t, evmtest.TestChainOpts{DB: db, GeneralConfig: config})
	orm := job.NewTestORM(t, db, cc, pipelineORM, keyStore, config)

	_, bridge := cltest.MustCreateBridge(t, db, cltest.BridgeOpts{}, config)
	_, bridge2 := cltest.MustCreateBridge(t, db, cltest.BridgeOpts{}, config)

	externalJobID := uuid.NewV4()
	_, address := cltest.MustInsertRandomKey(t, keyStore.Eth())
	jb, err := ocr.ValidatedOracleSpecToml(cc,
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
		count, err := orm.CountPipelineRunsByJobID(jb.ID)
		require.NoError(t, err)
		assert.Equal(t, int32(0), count)
	})

	t.Run("with a pipeline run", func(t *testing.T) {
		mustInsertPipelineRun(t, pipelineORM, jb)

		count, err := orm.CountPipelineRunsByJobID(jb.ID)
		require.NoError(t, err)
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
