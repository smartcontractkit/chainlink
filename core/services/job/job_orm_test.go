package job_test

import (
	"context"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/evmtest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/directrequest"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/keeper"
	"github.com/smartcontractkit/chainlink/core/services/offchainreporting"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/core/services/postgres"
	"github.com/smartcontractkit/chainlink/core/services/vrf"
	"github.com/smartcontractkit/chainlink/core/services/webhook"
	"github.com/smartcontractkit/chainlink/core/testdata/testspecs"
	"gopkg.in/guregu/null.v4"

	"github.com/pelletier/go-toml"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestORM(t *testing.T) {
	t.Parallel()
	config := cltest.NewTestGeneralConfig(t)
	db := pgtest.NewSqlxDB(t)
	gdb := pgtest.GormDBFromSql(t, db.DB)
	keyStore := cltest.NewKeyStore(t, db)
	ethKeyStore := keyStore.Eth()

	keyStore.OCR().Add(cltest.DefaultOCRKey)

	pipelineORM := pipeline.NewORM(db, logger.TestLogger(t))

	cc := evmtest.NewChainSet(t, evmtest.TestChainOpts{DB: gdb, GeneralConfig: config})
	orm := job.NewTestORM(t, db, cc, pipelineORM, keyStore)

	_, bridge := cltest.MustCreateBridge(t, db, cltest.BridgeOpts{})
	_, bridge2 := cltest.MustCreateBridge(t, db, cltest.BridgeOpts{})
	_, address := cltest.MustInsertRandomKey(t, ethKeyStore)
	jb := makeOCRJobSpec(t, address, bridge.Name.String(), bridge2.Name.String())

	t.Run("it creates job specs", func(t *testing.T) {
		err := orm.CreateJob(jb)
		require.NoError(t, err)

		var returnedSpec job.Job
		err = gdb.
			Preload("OffchainreportingOracleSpec").
			Where("id = ?", jb.ID).First(&returnedSpec).Error
		require.NoError(t, err)
		compareOCRJobSpecs(t, *jb, returnedSpec)
	})

	t.Run("autogenerates external job ID if missing", func(t *testing.T) {
		jb2 := makeOCRJobSpec(t, address, bridge.Name.String(), bridge2.Name.String())
		jb2.ExternalJobID = uuid.UUID{}
		err := orm.CreateJob(jb2)
		require.NoError(t, err)

		var returnedSpec job.Job
		err = gdb.Where("id = ?", jb.ID).First(&returnedSpec).Error
		require.NoError(t, err)

		assert.NotEqual(t, uuid.UUID{}, returnedSpec.ExternalJobID)
	})

	t.Run("it deletes jobs from the DB", func(t *testing.T) {
		var dbSpecs []job.Job
		err := gdb.Find(&dbSpecs).Error
		require.NoError(t, err)
		require.Len(t, dbSpecs, 2)

		err = orm.DeleteJob(jb.ID)
		require.NoError(t, err)

		err = gdb.Find(&dbSpecs).Error
		require.NoError(t, err)
		require.Len(t, dbSpecs, 1)
	})

	t.Run("increase job spec error occurrence", func(t *testing.T) {
		jb3 := makeOCRJobSpec(t, address, bridge.Name.String(), bridge2.Name.String())
		err := orm.CreateJob(jb3)
		require.NoError(t, err)
		var jobSpec job.Job
		err = gdb.
			First(&jobSpec).
			Error
		require.NoError(t, err)

		ocrSpecError1 := "ocr spec 1 errored"
		ocrSpecError2 := "ocr spec 2 errored"
		orm.RecordError(context.Background(), jobSpec.ID, ocrSpecError1)
		orm.RecordError(context.Background(), jobSpec.ID, ocrSpecError1)
		orm.RecordError(context.Background(), jobSpec.ID, ocrSpecError2)

		var specErrors []job.SpecError
		err = gdb.Find(&specErrors).Error
		require.NoError(t, err)
		require.Len(t, specErrors, 2)

		assert.Equal(t, specErrors[0].Occurrences, uint(2))
		assert.Equal(t, specErrors[1].Occurrences, uint(1))
		assert.True(t, specErrors[0].CreatedAt.Before(specErrors[0].UpdatedAt), "expected created_at (%s) to be before updated_at (%s)", specErrors[0].CreatedAt, specErrors[0].UpdatedAt)
		assert.Equal(t, specErrors[0].Description, ocrSpecError1)
		assert.Equal(t, specErrors[1].Description, ocrSpecError2)
		assert.True(t, specErrors[1].CreatedAt.After(specErrors[0].UpdatedAt))
		var j2 job.Job
		err = gdb.
			Preload("OffchainreportingOracleSpec").
			Preload("JobSpecErrors").
			First(&j2, "jobs.id = ?", jobSpec.ID).
			Error
		require.NoError(t, err)
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
		eiFoo := cltest.MustInsertExternalInitiator(t, gdb)
		eiBar := cltest.MustInsertExternalInitiator(t, gdb)

		eiWS := []webhook.TOMLWebhookSpecExternalInitiator{
			{Name: eiFoo.Name, Spec: cltest.JSONFromString(t, `{}`)},
			{Name: eiBar.Name, Spec: cltest.JSONFromString(t, `{"bar": 1}`)},
		}
		eim := webhook.NewExternalInitiatorManager(gdb, nil)
		jb, err := webhook.ValidatedWebhookSpec(testspecs.GenerateWebhookSpec(testspecs.WebhookSpecParams{ExternalInitiators: eiWS}).Toml(), eim)
		require.NoError(t, err)

		err = orm.CreateJob(&jb)
		require.NoError(t, err)

		cltest.AssertCount(t, gdb, job.ExternalInitiatorWebhookSpec{}, 2)
	})
}

func TestORM_DeleteJob_DeletesAssociatedRecords(t *testing.T) {
	t.Parallel()
	config := evmtest.NewChainScopedConfig(t, cltest.NewTestGeneralConfig(t))
	gdb := pgtest.NewGormDB(t)
	db := postgres.UnwrapGormDB(gdb)
	keyStore := cltest.NewKeyStore(t, db)
	keyStore.OCR().Add(cltest.DefaultOCRKey)

	pipelineORM := pipeline.NewORM(db, logger.TestLogger(t))
	cc := evmtest.NewChainSet(t, evmtest.TestChainOpts{DB: gdb, GeneralConfig: config})
	orm := job.NewTestORM(t, db, cc, pipelineORM, keyStore)

	t.Run("it deletes records for offchainreporting jobs", func(t *testing.T) {
		_, bridge := cltest.MustCreateBridge(t, db, cltest.BridgeOpts{})
		_, bridge2 := cltest.MustCreateBridge(t, db, cltest.BridgeOpts{})

		_, address := cltest.MustInsertRandomKey(t, keyStore.Eth())
		jb, err := offchainreporting.ValidatedOracleSpecToml(cc, testspecs.GenerateOCRSpec(testspecs.OCRSpecParams{
			TransmitterAddress: address.Hex(),
			DS1BridgeName:      bridge.Name.String(),
			DS2BridgeName:      bridge2.Name.String(),
		}).Toml())
		require.NoError(t, err)

		err = orm.CreateJob(&jb)
		require.NoError(t, err)

		cltest.AssertCount(t, gdb, job.OffchainReportingOracleSpec{}, 1)
		cltest.AssertCount(t, gdb, pipeline.Spec{}, 1)

		err = orm.DeleteJob(jb.ID)
		require.NoError(t, err)
		cltest.AssertCount(t, gdb, job.OffchainReportingOracleSpec{}, 0)
		cltest.AssertCount(t, gdb, pipeline.Spec{}, 0)
		cltest.AssertCount(t, gdb, job.Job{}, 0)
	})

	t.Run("it deletes records for keeper jobs", func(t *testing.T) {
		registry, keeperJob := cltest.MustInsertKeeperRegistry(t, gdb, keyStore.Eth())
		cltest.MustInsertUpkeepForRegistry(t, gdb, config, registry)

		cltest.AssertCount(t, gdb, job.KeeperSpec{}, 1)
		cltest.AssertCount(t, gdb, keeper.Registry{}, 1)
		cltest.AssertCount(t, gdb, keeper.UpkeepRegistration{}, 1)

		err := orm.DeleteJob(keeperJob.ID)
		require.NoError(t, err)
		cltest.AssertCount(t, gdb, job.KeeperSpec{}, 0)
		cltest.AssertCount(t, gdb, keeper.Registry{}, 0)
		cltest.AssertCount(t, gdb, keeper.UpkeepRegistration{}, 0)
		cltest.AssertCount(t, gdb, job.Job{}, 0)
	})

	t.Run("it deletes records for vrf jobs", func(t *testing.T) {
		key, err := keyStore.VRF().Create()
		require.NoError(t, err)
		pk := key.PublicKey
		jb, err := vrf.ValidatedVRFSpec(testspecs.GenerateVRFSpec(testspecs.VRFSpecParams{PublicKey: pk.String()}).Toml())
		require.NoError(t, err)

		err = orm.CreateJob(&jb)
		require.NoError(t, err)
		err = orm.DeleteJob(jb.ID)
		require.NoError(t, err)
		cltest.AssertCount(t, gdb, job.VRFSpec{}, 0)
		cltest.AssertCount(t, gdb, job.Job{}, 0)
	})

	t.Run("it deletes records for webhook jobs", func(t *testing.T) {
		ei := cltest.MustInsertExternalInitiator(t, gdb)
		jb, webhookSpec := cltest.MustInsertWebhookSpec(t, gdb)
		err := gdb.Exec(`INSERT INTO external_initiator_webhook_specs (external_initiator_id, webhook_spec_id, spec) VALUES (?,?,?)`, ei.ID, webhookSpec.ID, `{"ei": "foo", "name": "webhookSpecTwoEIs"}`).Error
		require.NoError(t, err)

		err = orm.DeleteJob(jb.ID)
		require.NoError(t, err)
		cltest.AssertCount(t, gdb, job.WebhookSpec{}, 0)
		cltest.AssertCount(t, gdb, job.ExternalInitiatorWebhookSpec{}, 0)
		cltest.AssertCount(t, gdb, job.Job{}, 0)
	})

	t.Run("does not allow to delete external initiators if they have referencing external_initiator_webhook_specs", func(t *testing.T) {
		// create new db because this will rollback transaction and poison it
		db2 := pgtest.NewGormDB(t)
		ei := cltest.MustInsertExternalInitiator(t, db2)
		_, webhookSpec := cltest.MustInsertWebhookSpec(t, db2)
		err := db2.Exec(`INSERT INTO external_initiator_webhook_specs (external_initiator_id, webhook_spec_id, spec) VALUES (?,?,?)`, ei.ID, webhookSpec.ID, `{"ei": "foo", "name": "webhookSpecTwoEIs"}`).Error
		require.NoError(t, err)

		err = db2.Exec(`DELETE FROM external_initiators`).Error
		require.EqualError(t, err, "ERROR: update or delete on table \"external_initiators\" violates foreign key constraint \"external_initiator_webhook_specs_external_initiator_id_fkey\" on table \"external_initiator_webhook_specs\" (SQLSTATE 23503)")
	})
}

func Test_FindJob(t *testing.T) {
	t.Parallel()

	config := cltest.NewTestGeneralConfig(t)
	db := pgtest.NewSqlxDB(t)
	gdb := pgtest.GormDBFromSql(t, db.DB)
	keyStore := cltest.NewKeyStore(t, db)
	keyStore.OCR().Add(cltest.DefaultOCRKey)

	pipelineORM := pipeline.NewORM(db, logger.TestLogger(t))
	cc := evmtest.NewChainSet(t, evmtest.TestChainOpts{DB: gdb, GeneralConfig: config})
	orm := job.NewTestORM(t, db, cc, pipelineORM, keyStore)

	_, bridge := cltest.MustCreateBridge(t, db, cltest.BridgeOpts{})
	_, bridge2 := cltest.MustCreateBridge(t, db, cltest.BridgeOpts{})

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
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		jb, err := orm.FindJobByExternalJobID(ctx, externalJobID)
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
	gdb := pgtest.NewGormDB(t)
	db := postgres.UnwrapGormDB(gdb)
	keyStore := cltest.NewKeyStore(t, db)
	keyStore.OCR().Add(cltest.DefaultOCRKey)

	pipelineORM := pipeline.NewORM(db, logger.TestLogger(t))
	cc := evmtest.NewChainSet(t, evmtest.TestChainOpts{DB: gdb, GeneralConfig: config})
	orm := job.NewTestORM(t, db, cc, pipelineORM, keyStore)

	_, bridge := cltest.MustCreateBridge(t, db, cltest.BridgeOpts{})
	_, bridge2 := cltest.MustCreateBridge(t, db, cltest.BridgeOpts{})

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
		run := mustInsertPipelineRun(t, gdb, jb)

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
	gdb := pgtest.NewGormDB(t)
	db := postgres.UnwrapGormDB(gdb)

	keyStore := cltest.NewKeyStore(t, db)
	keyStore.OCR().Add(cltest.DefaultOCRKey)

	pipelineORM := pipeline.NewORM(db, logger.TestLogger(t))
	cc := evmtest.NewChainSet(t, evmtest.TestChainOpts{DB: gdb, GeneralConfig: config})
	orm := job.NewTestORM(t, db, cc, pipelineORM, keyStore)

	_, bridge := cltest.MustCreateBridge(t, db, cltest.BridgeOpts{})
	_, bridge2 := cltest.MustCreateBridge(t, db, cltest.BridgeOpts{})

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
		run := mustInsertPipelineRun(t, gdb, jb)

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

func mustInsertPipelineRun(t *testing.T, db *gorm.DB, j job.Job) pipeline.Run {
	t.Helper()

	run := pipeline.Run{
		PipelineSpecID: j.PipelineSpecID,
		State:          pipeline.RunStatusRunning,
		Outputs:        pipeline.JSONSerializable{Valid: false},
		AllErrors:      pipeline.RunErrors{},
		FinishedAt:     null.Time{},
	}
	require.NoError(t, db.Create(&run).Error)
	return run
}
