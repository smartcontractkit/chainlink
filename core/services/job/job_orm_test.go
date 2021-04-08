package job_test

import (
	"context"
	"testing"
	"time"

	gormpostgres "gorm.io/driver/postgres"

	"github.com/pelletier/go-toml"
	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/mocks"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/keeper"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/core/services/postgres"
)

func TestORM(t *testing.T) {
	t.Parallel()
	config, oldORM, cleanupDB := cltest.BootstrapThrowawayORM(t, "services_job_orm", true, true)
	defer cleanupDB()
	db := oldORM.DB

	pipelineORM, eventBroadcaster, cleanupORM := cltest.NewPipelineORM(t, config, db)
	defer cleanupORM()

	orm := job.NewORM(db, config.Config, pipelineORM, eventBroadcaster, &postgres.NullAdvisoryLocker{})
	defer orm.Close()

	_, bridge := cltest.NewBridgeType(t, "voter_turnout", "http://blah.com")
	require.NoError(t, db.Create(bridge).Error)
	_, bridge2 := cltest.NewBridgeType(t, "election_winner", "http://blah.com")
	require.NoError(t, db.Create(bridge2).Error)
	key := cltest.MustInsertRandomKey(t, db)
	address := key.Address.Address()
	dbSpec := makeOCRJobSpec(t, address)

	t.Run("it creates job specs", func(t *testing.T) {
		err := orm.CreateJob(context.Background(), dbSpec, dbSpec.Pipeline)
		require.NoError(t, err)

		var returnedSpec job.Job
		err = db.
			Preload("OffchainreportingOracleSpec").
			Where("id = ?", dbSpec.ID).First(&returnedSpec).Error
		require.NoError(t, err)
		compareOCRJobSpecs(t, *dbSpec, returnedSpec)
	})

	dbURL := config.DatabaseURL()
	db2, err := gorm.Open(gormpostgres.New(gormpostgres.Config{
		DSN: dbURL.String(),
	}), &gorm.Config{})
	require.NoError(t, err)
	d, err := db2.DB()
	require.NoError(t, err)
	defer d.Close()

	orm2 := job.NewORM(db2, config.Config, pipeline.NewORM(db2, config, eventBroadcaster), eventBroadcaster, &postgres.NullAdvisoryLocker{})
	defer orm2.Close()

	t.Run("it correctly returns the unclaimed jobs in the DB", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		unclaimed, err := orm.ClaimUnclaimedJobs(ctx)
		require.NoError(t, err)
		require.Len(t, unclaimed, 1)
		compareOCRJobSpecs(t, *dbSpec, unclaimed[0])
		require.Equal(t, int32(1), unclaimed[0].ID)
		require.Equal(t, int32(1), *unclaimed[0].OffchainreportingOracleSpecID)
		require.Equal(t, int32(1), unclaimed[0].PipelineSpecID)
		require.Equal(t, int32(1), unclaimed[0].OffchainreportingOracleSpec.ID)

		dbSpec2 := makeOCRJobSpec(t, address)
		err = orm.CreateJob(context.Background(), dbSpec2, dbSpec2.Pipeline)
		require.NoError(t, err)

		ctx2, cancel2 := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel2()

		unclaimed, err = orm2.ClaimUnclaimedJobs(ctx2)
		require.NoError(t, err)
		require.Len(t, unclaimed, 1)
		compareOCRJobSpecs(t, *dbSpec2, unclaimed[0])
		require.Equal(t, int32(2), unclaimed[0].ID)
		require.Equal(t, int32(2), *unclaimed[0].OffchainreportingOracleSpecID)
		require.Equal(t, int32(2), unclaimed[0].PipelineSpecID)
		require.Equal(t, int32(2), unclaimed[0].OffchainreportingOracleSpec.ID)
	})

	t.Run("it can delete jobs claimed by other nodes", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := orm2.DeleteJob(ctx, dbSpec.ID)
		require.NoError(t, err)

		var dbSpecs []job.Job
		err = db.Find(&dbSpecs).Error
		require.NoError(t, err)
		require.Len(t, dbSpecs, 1)
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
		require.Len(t, dbSpecs, 1)

		var oracleSpecs []job.OffchainReportingOracleSpec
		err = db.Find(&oracleSpecs).Error
		require.NoError(t, err)
		require.Len(t, oracleSpecs, 1)

		var pipelineSpecs []pipeline.Spec
		err = db.Find(&pipelineSpecs).Error
		require.NoError(t, err)
		require.Len(t, pipelineSpecs, 1)
	})

	t.Run("increase job spec error occurrence", func(t *testing.T) {
		dbSpec3 := makeOCRJobSpec(t, address)
		err := orm.CreateJob(context.Background(), dbSpec3, dbSpec3.Pipeline)
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
		spec := job.Job{}
		tree, err := toml.LoadFile("../../cmd/testdata/direct-request-spec.toml")
		require.NoError(t, err)
		err = tree.Unmarshal(&spec)
		require.NoError(t, err)

		var drSpec job.DirectRequestSpec
		err = tree.Unmarshal(&drSpec)
		require.NoError(t, err)
		spec.DirectRequestSpec = &drSpec

		err = orm.CreateJob(context.Background(), &spec, spec.Pipeline)
		require.NoError(t, err)
	})
}

func TestORM_CheckForDeletedJobs(t *testing.T) {
	t.Parallel()

	config, cleanup := cltest.NewConfig(t)
	defer cleanup()
	store, cleanup := cltest.NewStoreWithConfig(t, config)
	defer cleanup()
	db := store.DB

	key := cltest.MustInsertRandomKey(t, db)
	address := key.Address.Address()

	_, bridge := cltest.NewBridgeType(t, "voter_turnout", "http://blah.com")
	require.NoError(t, db.Create(bridge).Error)
	_, bridge2 := cltest.NewBridgeType(t, "election_winner", "http://blah.com")
	require.NoError(t, db.Create(bridge2).Error)

	pipelineORM, eventBroadcaster, cleanupORM := cltest.NewPipelineORM(t, config, db)
	defer cleanupORM()

	orm := job.NewORM(db, config.Config, pipelineORM, eventBroadcaster, &postgres.NullAdvisoryLocker{})
	defer orm.Close()

	claimedJobs := make([]job.Job, 3)
	for i := range claimedJobs {
		dbSpec := makeOCRJobSpec(t, address)
		require.NoError(t, orm.CreateJob(context.Background(), dbSpec, dbSpec.Pipeline))
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

	config, cleanup := cltest.NewConfig(t)
	defer cleanup()
	store, cleanup := cltest.NewStoreWithConfig(t, config)
	defer cleanup()
	db := store.DB

	key := cltest.MustInsertRandomKey(t, db)
	address := key.Address.Address()

	pipelineORM, eventBroadcaster, cleanupORM := cltest.NewPipelineORM(t, config, db)
	defer cleanupORM()

	advisoryLocker := new(mocks.AdvisoryLocker)
	orm := job.NewORM(db, config.Config, pipelineORM, eventBroadcaster, advisoryLocker)
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
	config, cleanup := cltest.NewConfig(t)
	defer cleanup()
	store, cleanup := cltest.NewStoreWithConfig(t, config)
	defer cleanup()
	db := store.DB

	pipelineORM, eventBroadcaster, cleanupORM := cltest.NewPipelineORM(t, config, db)
	defer cleanupORM()
	orm := job.NewORM(db, config.Config, pipelineORM, eventBroadcaster, &postgres.NullAdvisoryLocker{})
	defer orm.Close()

	t.Run("it deletes records for offchainreporting jobs", func(t *testing.T) {
		_, bridge := cltest.NewBridgeType(t, "voter_turnout", "http://blah.com")
		require.NoError(t, db.Create(bridge).Error)
		_, bridge2 := cltest.NewBridgeType(t, "election_winner", "http://blah.com")
		require.NoError(t, db.Create(bridge2).Error)

		key := cltest.MustInsertRandomKey(t, store.DB)
		address := key.Address.Address()
		dbSpec := makeOCRJobSpec(t, address)

		err := orm.CreateJob(context.Background(), dbSpec, dbSpec.Pipeline)
		require.NoError(t, err)

		var ocrJob job.Job
		err = store.DB.First(&ocrJob).Error
		require.NoError(t, err)

		cltest.AssertCount(t, store, job.OffchainReportingOracleSpec{}, 1)
		cltest.AssertCount(t, store, pipeline.Spec{}, 1)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		err = orm.DeleteJob(ctx, ocrJob.ID)
		require.NoError(t, err)
		cltest.AssertCount(t, store, job.OffchainReportingOracleSpec{}, 0)
		cltest.AssertCount(t, store, pipeline.Spec{}, 0)
		cltest.AssertCount(t, store, job.Job{}, 0)
	})

	t.Run("it deletes records for keeper jobs", func(t *testing.T) {
		registry, keeperJob := cltest.MustInsertKeeperRegistry(t, store)
		cltest.MustInsertUpkeepForRegistry(t, store, registry)

		cltest.AssertCount(t, store, job.KeeperSpec{}, 1)
		cltest.AssertCount(t, store, keeper.Registry{}, 1)
		cltest.AssertCount(t, store, keeper.UpkeepRegistration{}, 1)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		err := orm.DeleteJob(ctx, keeperJob.ID)
		require.NoError(t, err)
		cltest.AssertCount(t, store, job.KeeperSpec{}, 0)
		cltest.AssertCount(t, store, keeper.Registry{}, 0)
		cltest.AssertCount(t, store, keeper.UpkeepRegistration{}, 0)
		cltest.AssertCount(t, store, job.Job{}, 0)
	})
}
