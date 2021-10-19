package job_test

import (
	"context"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/cltest/heavyweight"
	"github.com/smartcontractkit/chainlink/core/internal/mocks"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/core/services/directrequest"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/keeper"
	"github.com/smartcontractkit/chainlink/core/services/offchainreporting"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/core/services/postgres"
	"github.com/smartcontractkit/chainlink/core/services/vrf"
	"github.com/smartcontractkit/chainlink/core/services/webhook"
	"github.com/smartcontractkit/chainlink/core/testdata/testspecs"

	"github.com/pelletier/go-toml"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	gormpostgres "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestORM(t *testing.T) {
	t.Parallel()
	config, oldORM, cleanupDB := heavyweight.FullTestORM(t, "services_job_orm", true, true)
	defer cleanupDB()
	db := oldORM.DB
	keyStore := cltest.NewKeyStore(t, db)
	ethKeyStore := keyStore.Eth()

	keyStore.OCR().Add(cltest.DefaultOCRKey)
	keyStore.P2P().Add(cltest.DefaultP2PKey)

	pipelineORM, eventBroadcaster, cleanupORM := cltest.NewPipelineORM(t, config, db)
	defer cleanupORM()

	orm := job.NewORM(db, config, pipelineORM, eventBroadcaster, &postgres.NullAdvisoryLocker{}, keyStore)
	defer orm.Close()

	_, bridge := cltest.NewBridgeType(t, "voter_turnout", "http://blah.com")
	require.NoError(t, db.Create(bridge).Error)
	_, bridge2 := cltest.NewBridgeType(t, "election_winner", "http://blah.com")
	require.NoError(t, db.Create(bridge2).Error)
	_, address := cltest.MustInsertRandomKey(t, ethKeyStore)
	dbSpec := makeOCRJobSpec(t, address)

	t.Run("it creates job specs", func(t *testing.T) {
		jb, err := orm.CreateJob(context.Background(), dbSpec, dbSpec.Pipeline)
		require.NoError(t, err)

		var returnedSpec job.Job
		err = db.
			Preload("OffchainreportingOracleSpec").
			Where("id = ?", dbSpec.ID).First(&returnedSpec).Error
		require.NoError(t, err)
		compareOCRJobSpecs(t, jb, returnedSpec)
	})

	t.Run("autogenerates external job ID if missing", func(t *testing.T) {
		job2 := makeOCRJobSpec(t, address)
		job2.ExternalJobID = uuid.UUID{}
		_, err := orm.CreateJob(context.Background(), job2, job2.Pipeline)
		require.NoError(t, err)

		var returnedSpec job.Job
		err = db.Where("id = ?", job2.ID).First(&returnedSpec).Error
		require.NoError(t, err)

		assert.NotEqual(t, uuid.UUID{}, returnedSpec.ExternalJobID)
	})

	dbURL := config.DatabaseURL()
	db2, err := gorm.Open(gormpostgres.New(gormpostgres.Config{
		DSN: dbURL.String(),
	}), &gorm.Config{})
	require.NoError(t, err)
	d, err := db2.DB()
	require.NoError(t, err)
	defer d.Close()

	orm2 := job.NewORM(db2, config, pipeline.NewORM(db2), eventBroadcaster, &postgres.NullAdvisoryLocker{}, keyStore)
	defer orm2.Close()

	t.Run("it correctly returns the unclaimed jobs in the DB", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		unclaimed, err := orm.ClaimUnclaimedJobs(ctx)
		require.NoError(t, err)
		require.Len(t, unclaimed, 2)
		compareOCRJobSpecs(t, *dbSpec, unclaimed[0])
		require.Equal(t, int32(1), unclaimed[0].ID)
		require.Equal(t, int32(1), *unclaimed[0].OffchainreportingOracleSpecID)
		require.Equal(t, int32(1), unclaimed[0].PipelineSpecID)
		require.Equal(t, int32(1), unclaimed[0].OffchainreportingOracleSpec.ID)

		dbSpec2 := makeOCRJobSpec(t, address)
		_, err = orm.CreateJob(context.Background(), dbSpec2, dbSpec2.Pipeline)
		require.NoError(t, err)

		ctx2, cancel2 := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel2()

		unclaimed, err = orm2.ClaimUnclaimedJobs(ctx2)
		require.NoError(t, err)
		require.Len(t, unclaimed, 1)
		compareOCRJobSpecs(t, *dbSpec2, unclaimed[0])
		require.Equal(t, int32(3), unclaimed[0].ID)
		require.Equal(t, int32(3), *unclaimed[0].OffchainreportingOracleSpecID)
		require.Equal(t, int32(3), unclaimed[0].PipelineSpecID)
		require.Equal(t, int32(3), unclaimed[0].OffchainreportingOracleSpec.ID)
	})

	t.Run("it can delete jobs claimed by other nodes", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := orm2.DeleteJob(ctx, dbSpec.ID)
		require.NoError(t, err)

		var dbSpecs []job.Job
		err = db.Find(&dbSpecs).Error
		require.NoError(t, err)
		require.Len(t, dbSpecs, 2)
	})

	t.Run("it deletes its own claimed jobs from the DB", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Check that it is claimed
		claimedJobIDs := job.GetORMClaimedJobIDs(orm)
		require.Contains(t, claimedJobIDs, dbSpec.ID)

		err := orm.DeleteJob(ctx, dbSpec.ID)
		require.NoError(t, err)

		// Check that it is no longer claimed
		claimedJobIDs = job.GetORMClaimedJobIDs(orm)
		assert.NotContains(t, claimedJobIDs, dbSpec.ID)

		var dbSpecs []job.Job
		err = db.Find(&dbSpecs).Error
		require.NoError(t, err)
		require.Len(t, dbSpecs, 2)

		var oracleSpecs []job.OffchainReportingOracleSpec
		err = db.Find(&oracleSpecs).Error
		require.NoError(t, err)
		require.Len(t, oracleSpecs, 2)

		var pipelineSpecs []pipeline.Spec
		err = db.Find(&pipelineSpecs).Error
		require.NoError(t, err)
		require.Len(t, pipelineSpecs, 2)
	})

	t.Run("increase job spec error occurrence", func(t *testing.T) {
		dbSpec3 := makeOCRJobSpec(t, address)
		_, err := orm.CreateJob(context.Background(), dbSpec3, dbSpec3.Pipeline)
		require.NoError(t, err)
		var jobSpec job.Job
		err = db.
			First(&jobSpec).
			Error
		require.NoError(t, err)

		ocrSpecError1 := "ocr spec 1 errored"
		ocrSpecError2 := "ocr spec 2 errored"
		orm.RecordError(context.Background(), jobSpec.ID, ocrSpecError1)
		orm.RecordError(context.Background(), jobSpec.ID, ocrSpecError1)
		orm.RecordError(context.Background(), jobSpec.ID, ocrSpecError2)

		var specErrors []job.SpecError
		err = db.Find(&specErrors).Error
		require.NoError(t, err)
		require.Len(t, specErrors, 2)

		assert.Equal(t, specErrors[0].Occurrences, uint(2))
		assert.Equal(t, specErrors[1].Occurrences, uint(1))
		assert.True(t, specErrors[0].CreatedAt.Before(specErrors[0].UpdatedAt))
		assert.Equal(t, specErrors[0].Description, ocrSpecError1)
		assert.Equal(t, specErrors[1].Description, ocrSpecError2)
		assert.True(t, specErrors[1].CreatedAt.After(specErrors[0].UpdatedAt))
		var j2 job.Job
		err = db.
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
		_, err = orm.CreateJob(context.Background(), &jb, jb.Pipeline)
		require.NoError(t, err)
	})

	t.Run("creates webhook specs along with external_initiator_webhook_specs", func(t *testing.T) {
		eiFoo := cltest.MustInsertExternalInitiator(t, db)
		eiBar := cltest.MustInsertExternalInitiator(t, db)

		eiWS := []webhook.TOMLWebhookSpecExternalInitiator{
			{Name: eiFoo.Name, Spec: cltest.JSONFromString(t, `{}`)},
			{Name: eiBar.Name, Spec: cltest.JSONFromString(t, `{"bar": 1}`)},
		}
		eim := webhook.NewExternalInitiatorManager(db, nil)
		jb, err := webhook.ValidatedWebhookSpec(testspecs.GenerateWebhookSpec(testspecs.WebhookSpecParams{ExternalInitiators: eiWS}).Toml(), eim)
		require.NoError(t, err)

		_, err = orm.CreateJob(context.Background(), &jb, jb.Pipeline)
		require.NoError(t, err)

		cltest.AssertCount(t, db, job.ExternalInitiatorWebhookSpec{}, 2)
	})
}

func TestORM_CheckForDeletedJobs(t *testing.T) {
	t.Parallel()

	config := cltest.NewTestEVMConfig(t)
	store, cleanup := cltest.NewStoreWithConfig(t, config)
	defer cleanup()
	db := store.DB
	keyStore := cltest.NewKeyStore(t, db)
	ethKeyStore := keyStore.Eth()

	_, address := cltest.MustInsertRandomKey(t, ethKeyStore)
	keyStore.OCR().Add(cltest.DefaultOCRKey)
	keyStore.P2P().Add(cltest.DefaultP2PKey)

	_, bridge := cltest.NewBridgeType(t, "voter_turnout", "http://blah.com")
	require.NoError(t, db.Create(bridge).Error)
	_, bridge2 := cltest.NewBridgeType(t, "election_winner", "http://blah.com")
	require.NoError(t, db.Create(bridge2).Error)

	pipelineORM, eventBroadcaster, cleanupORM := cltest.NewPipelineORM(t, config, db)
	defer cleanupORM()

	orm := job.NewORM(db, config, pipelineORM, eventBroadcaster, &postgres.NullAdvisoryLocker{}, keyStore)
	defer orm.Close()

	claimedJobs := make([]job.Job, 3)
	for i := range claimedJobs {
		dbSpec := makeOCRJobSpec(t, address)
		_, err := orm.CreateJob(context.Background(), dbSpec, dbSpec.Pipeline)
		require.NoError(t, err)
		claimedJobs[i] = *dbSpec
	}
	job.SetORMClaimedJobs(orm, claimedJobs)

	deletedID := claimedJobs[0].ID
	require.NoError(t, db.Exec(`DELETE FROM jobs WHERE id = ?`, deletedID).Error)

	deletedJobIDs, err := orm.CheckForDeletedJobs(context.Background())
	require.NoError(t, err)

	require.Len(t, deletedJobIDs, 1)
	require.Equal(t, deletedID, deletedJobIDs[0])

}

func TestORM_UnclaimJob(t *testing.T) {
	t.Parallel()

	config := cltest.NewTestEVMConfig(t)
	store, cleanup := cltest.NewStoreWithConfig(t, config)
	defer cleanup()
	db := store.DB
	keyStore := cltest.NewKeyStore(t, db)
	ethKeyStore := keyStore.Eth()

	_, address := cltest.MustInsertRandomKey(t, ethKeyStore)

	pipelineORM, eventBroadcaster, cleanupORM := cltest.NewPipelineORM(t, config, db)
	defer cleanupORM()

	advisoryLocker := new(mocks.AdvisoryLocker)
	orm := job.NewORM(db, config, pipelineORM, eventBroadcaster, advisoryLocker, keyStore)
	defer orm.Close()

	require.NoError(t, orm.UnclaimJob(context.Background(), 42))

	claimedJobs := make([]job.Job, 3)
	for i := range claimedJobs {
		dbSpec := makeOCRJobSpec(t, address)
		dbSpec.ID = int32(i)
		claimedJobs[i] = *dbSpec
	}

	job.SetORMClaimedJobs(orm, claimedJobs)

	jobID := claimedJobs[0].ID
	advisoryLocker.On("Unlock", mock.Anything, job.GetORMAdvisoryLockClassID(orm), jobID).Once().Return(nil)

	require.NoError(t, orm.UnclaimJob(context.Background(), jobID))

	claimedJobs = job.GetORMClaimedJobs(orm)
	require.Len(t, claimedJobs, 2)
	require.NotContains(t, claimedJobs, jobID)

	advisoryLocker.AssertExpectations(t)
}

func TestORM_DeleteJob_DeletesAssociatedRecords(t *testing.T) {
	t.Parallel()
	config := cltest.NewTestEVMConfig(t)
	store, cleanup := cltest.NewStoreWithConfig(t, config)
	defer cleanup()
	db := store.DB
	keyStore := cltest.NewKeyStore(t, store.DB)
	keyStore.OCR().Add(cltest.DefaultOCRKey)
	keyStore.P2P().Add(cltest.DefaultP2PKey)

	pipelineORM, eventBroadcaster, cleanupORM := cltest.NewPipelineORM(t, config, db)
	defer cleanupORM()
	orm := job.NewORM(db, config, pipelineORM, eventBroadcaster, &postgres.NullAdvisoryLocker{}, keyStore)
	defer orm.Close()

	t.Run("it deletes records for offchainreporting jobs", func(t *testing.T) {
		_, bridge := cltest.NewBridgeType(t, "voter_turnout", "http://blah.com")
		require.NoError(t, db.Create(bridge).Error)
		_, bridge2 := cltest.NewBridgeType(t, "election_winner", "http://blah.com")
		require.NoError(t, db.Create(bridge2).Error)

		_, address := cltest.MustInsertRandomKey(t, keyStore.Eth())
		jb, err := offchainreporting.ValidatedOracleSpecToml(config, testspecs.GenerateOCRSpec(testspecs.OCRSpecParams{TransmitterAddress: address.Hex()}).Toml())
		require.NoError(t, err)

		ocrJob, err := orm.CreateJob(context.Background(), &jb, jb.Pipeline)
		require.NoError(t, err)

		cltest.AssertCount(t, db, job.OffchainReportingOracleSpec{}, 1)
		cltest.AssertCount(t, db, pipeline.Spec{}, 1)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		err = orm.DeleteJob(ctx, ocrJob.ID)
		require.NoError(t, err)
		cltest.AssertCount(t, db, job.OffchainReportingOracleSpec{}, 0)
		cltest.AssertCount(t, db, pipeline.Spec{}, 0)
		cltest.AssertCount(t, db, job.Job{}, 0)
	})

	t.Run("it deletes records for keeper jobs", func(t *testing.T) {
		registry, keeperJob := cltest.MustInsertKeeperRegistry(t, store, keyStore.Eth())
		cltest.MustInsertUpkeepForRegistry(t, store, registry)

		cltest.AssertCount(t, db, job.KeeperSpec{}, 1)
		cltest.AssertCount(t, db, keeper.Registry{}, 1)
		cltest.AssertCount(t, db, keeper.UpkeepRegistration{}, 1)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		err := orm.DeleteJob(ctx, keeperJob.ID)
		require.NoError(t, err)
		cltest.AssertCount(t, db, job.KeeperSpec{}, 0)
		cltest.AssertCount(t, db, keeper.Registry{}, 0)
		cltest.AssertCount(t, db, keeper.UpkeepRegistration{}, 0)
		cltest.AssertCount(t, db, job.Job{}, 0)
	})

	t.Run("it deletes records for vrf jobs", func(t *testing.T) {
		key, err := keyStore.VRF().Create()
		require.NoError(t, err)
		pk := key.PublicKey
		jb, err := vrf.ValidatedVRFSpec(testspecs.GenerateVRFSpec(testspecs.VRFSpecParams{PublicKey: pk.String()}).Toml())
		require.NoError(t, err)

		_, err = orm.CreateJob(context.Background(), &jb, jb.Pipeline)
		require.NoError(t, err)
		ctx, cancel := postgres.DefaultQueryCtx()
		defer cancel()
		err = orm.DeleteJob(ctx, jb.ID)
		require.NoError(t, err)
		cltest.AssertCount(t, db, job.VRFSpec{}, 0)
		cltest.AssertCount(t, db, job.Job{}, 0)
	})

	t.Run("it deletes records for webhook jobs", func(t *testing.T) {
		ei := cltest.MustInsertExternalInitiator(t, db)
		jb, webhookSpec := cltest.MustInsertWebhookSpec(t, db)
		err := db.Exec(`INSERT INTO external_initiator_webhook_specs (external_initiator_id, webhook_spec_id, spec) VALUES (?,?,?)`, ei.ID, webhookSpec.ID, `{"ei": "foo", "name": "webhookSpecTwoEIs"}`).Error
		require.NoError(t, err)

		ctx, cancel := postgres.DefaultQueryCtx()
		defer cancel()
		err = orm.DeleteJob(ctx, jb.ID)
		require.NoError(t, err)
		cltest.AssertCount(t, db, job.WebhookSpec{}, 0)
		cltest.AssertCount(t, db, job.ExternalInitiatorWebhookSpec{}, 0)
		cltest.AssertCount(t, db, job.Job{}, 0)
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
