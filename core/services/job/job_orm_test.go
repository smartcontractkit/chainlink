package job_test

import (
	"context"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/core/bridges"

	"github.com/pelletier/go-toml"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/evmtest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/directrequest"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/keeper"
	"github.com/smartcontractkit/chainlink/core/services/offchainreporting"
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
		var OCROracleSpec job.OffchainReportingOracleSpec

		err = db.Get(&returnedSpec, "SELECT * FROM jobs WHERE jobs.id = $1", jb.ID)
		require.NoError(t, err)
		err = db.Get(&OCROracleSpec, "SELECT * FROM offchainreporting_oracle_specs WHERE offchainreporting_oracle_specs.id = $1", jb.OffchainreportingOracleSpecID)
		require.NoError(t, err)
		returnedSpec.OffchainreportingOracleSpec = &OCROracleSpec
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
		var OCROracleSpec job.OffchainReportingOracleSpec
		var jobSpecErrors []job.SpecError

		err = db.Get(&j2, "SELECT * FROM jobs WHERE jobs.id = $1", jobSpec.ID)
		require.NoError(t, err)
		err = db.Get(&OCROracleSpec, "SELECT * FROM offchainreporting_oracle_specs WHERE offchainreporting_oracle_specs.id = $1", j2.OffchainreportingOracleSpecID)
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

		expected1ID := specErrors[1].ID + 1
		expected2ID := specErrors[1].ID + 2
		specErr1, err := orm.FindSpecError(expected1ID)
		require.NoError(t, err)
		specErr2, err := orm.FindSpecError(expected2ID)
		require.NoError(t, err)

		assert.Equal(t, specErr1.ID, expected1ID)
		assert.Equal(t, specErr2.ID, expected2ID)
		assert.Equal(t, specErr1.Occurrences, uint(1))
		assert.Equal(t, specErr2.Occurrences, uint(1))
		assert.Equal(t, specErr1.Description, ocrSpecError1)
		assert.Equal(t, specErr2.Description, ocrSpecError2)
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
	korm := keeper.NewORM(db, logger.TestLogger(t), nil, nil, nil)

	t.Run("it deletes records for offchainreporting jobs", func(t *testing.T) {
		_, bridge := cltest.MustCreateBridge(t, db, cltest.BridgeOpts{}, config)
		_, bridge2 := cltest.MustCreateBridge(t, db, cltest.BridgeOpts{}, config)

		_, address := cltest.MustInsertRandomKey(t, keyStore.Eth())
		jb, err := offchainreporting.ValidatedOracleSpecToml(cc, testspecs.GenerateOCRSpec(testspecs.OCRSpecParams{
			TransmitterAddress: address.Hex(),
			DS1BridgeName:      bridge.Name.String(),
			DS2BridgeName:      bridge2.Name.String(),
		}).Toml())
		require.NoError(t, err)

		err = jobORM.CreateJob(&jb)
		require.NoError(t, err)

		cltest.AssertCount(t, db, "offchainreporting_oracle_specs", 1)
		cltest.AssertCount(t, db, "pipeline_specs", 1)

		err = jobORM.DeleteJob(jb.ID)
		require.NoError(t, err)
		cltest.AssertCount(t, db, "offchainreporting_oracle_specs", 0)
		cltest.AssertCount(t, db, "pipeline_specs", 0)
		cltest.AssertCount(t, db, "jobs", 0)
	})

	t.Run("it deletes records for keeper jobs", func(t *testing.T) {
		registry, keeperJob := cltest.MustInsertKeeperRegistry(t, db, korm, keyStore.Eth())
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
	keyStore.OCR().Add(cltest.DefaultOCRKey)

	pipelineORM := pipeline.NewORM(db, logger.TestLogger(t), config)
	cc := evmtest.NewChainSet(t, evmtest.TestChainOpts{DB: db, GeneralConfig: config})
	jobORM := job.NewTestORM(t, db, cc, pipelineORM, keyStore, config)

	jb, err := vrf.ValidatedVRFSpec(testspecs.GenerateVRFSpec(testspecs.VRFSpecParams{RequestedConfsDelay: 10}).Toml())
	require.NoError(t, err)

	err = jobORM.CreateJob(&jb)
	require.NoError(t, err)
	cltest.AssertCount(t, db, "vrf_specs", 1)
	cltest.AssertCount(t, db, "jobs", 1)
	var requestedConfsDelay int64
	require.NoError(t, db.Get(&requestedConfsDelay, `SELECT requested_confs_delay FROM vrf_specs LIMIT 1`))
	require.Equal(t, int64(10), requestedConfsDelay)
	jobORM.DeleteJob(jb.ID)
	cltest.AssertCount(t, db, "vrf_specs", 0)
	cltest.AssertCount(t, db, "jobs", 0)

	jb, err = vrf.ValidatedVRFSpec(testspecs.GenerateVRFSpec(testspecs.VRFSpecParams{}).Toml())
	require.NoError(t, err)
	err = jobORM.CreateJob(&jb)
	require.NoError(t, err)
	cltest.AssertCount(t, db, "vrf_specs", 1)
	cltest.AssertCount(t, db, "jobs", 1)
	require.NoError(t, db.Get(&requestedConfsDelay, `SELECT requested_confs_delay FROM vrf_specs LIMIT 1`))
	require.Equal(t, int64(0), requestedConfsDelay)
	jobORM.DeleteJob(jb.ID)
	cltest.AssertCount(t, db, "vrf_specs", 0)
	cltest.AssertCount(t, db, "jobs", 0)
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
	jb1, err := offchainreporting.ValidatedOracleSpecToml(cc,
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
	jb, err := offchainreporting.ValidatedOracleSpecToml(cc,
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

	t.Run("by id", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		jb, err = orm.FindJob(ctx, jb.ID)
		require.NoError(t, err)

		assert.Equal(t, jb.ID, jb.ID)
		assert.Equal(t, jb.Name, jb.Name)

		require.Greater(t, jb.PipelineSpecID, int32(0))
		require.NotNil(t, jb.PipelineSpec)
		require.NotNil(t, jb.OffchainreportingOracleSpecID)
		require.NotNil(t, jb.OffchainreportingOracleSpec)
	})

	t.Run("by external job id", func(t *testing.T) {
		jb, err := orm.FindJobByExternalJobID(externalJobID)
		require.NoError(t, err)

		assert.Equal(t, jb.ID, jb.ID)
		assert.Equal(t, jb.Name, jb.Name)

		require.Greater(t, jb.PipelineSpecID, int32(0))
		require.NotNil(t, jb.PipelineSpec)
		require.NotNil(t, jb.OffchainreportingOracleSpecID)
		require.NotNil(t, jb.OffchainreportingOracleSpec)
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
	jb, err := offchainreporting.ValidatedOracleSpecToml(cc,
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
	keyStore.OCR().Add(cltest.DefaultOCRKey)
	require.NoError(t, keyStore.P2P().Add(cltest.DefaultP2PKey))

	pipelineORM := pipeline.NewORM(db, logger.TestLogger(t), config)
	cc := evmtest.NewChainSet(t, evmtest.TestChainOpts{DB: db, GeneralConfig: config})
	orm := job.NewTestORM(t, db, cc, pipelineORM, keyStore, config)

	_, bridge := cltest.MustCreateBridge(t, db, cltest.BridgeOpts{}, config)
	_, bridge2 := cltest.MustCreateBridge(t, db, cltest.BridgeOpts{}, config)

	externalJobID := uuid.NewV4()
	_, address := cltest.MustInsertRandomKey(t, keyStore.Eth())
	jb, err := offchainreporting.ValidatedOracleSpecToml(cc,
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
	jb, err := offchainreporting.ValidatedOracleSpecToml(cc,
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
	jb, err := offchainreporting.ValidatedOracleSpecToml(cc,
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
	jb, err := offchainreporting.ValidatedOracleSpecToml(cc,
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
