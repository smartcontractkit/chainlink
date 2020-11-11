package job_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/mocks"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/core/services/postgres"
	"github.com/smartcontractkit/chainlink/core/store/models"
	ormpkg "github.com/smartcontractkit/chainlink/core/store/orm"
)

func TestORM(t *testing.T) {
	t.Parallel()
	config, oldORM, cleanupDB := cltest.BootstrapThrowawayORM(t, "services_job_orm", true, true)
	defer cleanupDB()
	db := oldORM.DB

	pipelineORM, eventBroadcaster, cleanupORM := cltest.NewPipelineORM(t, config, db)
	defer cleanupORM()

	orm := job.NewORM(db, config, pipelineORM, eventBroadcaster, &postgres.NullAdvisoryLocker{})
	defer orm.Close()

	ocrSpec, dbSpec := makeOCRJobSpec(t, db)

	t.Run("it creates job specs", func(t *testing.T) {
		err := orm.CreateJob(context.Background(), dbSpec, ocrSpec.TaskDAG())
		require.NoError(t, err)

		var returnedSpec models.JobSpecV2
		err = db.
			Preload("OffchainreportingOracleSpec").
			Where("id = ?", dbSpec.ID).First(&returnedSpec).Error
		require.NoError(t, err)
		compareOCRJobSpecs(t, *dbSpec, returnedSpec)
	})

	db2, err := gorm.Open(string(ormpkg.DialectPostgres), config.DatabaseURL())
	require.NoError(t, err)
	defer db2.Close()

	orm2 := job.NewORM(db2, config, pipeline.NewORM(db2, config, eventBroadcaster), eventBroadcaster, &postgres.NullAdvisoryLocker{})
	defer orm2.Close()

	t.Run("it correctly returns the unclaimed jobs in the DB", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		unclaimed, err := orm.ClaimUnclaimedJobs(ctx)
		require.NoError(t, err)
		require.Len(t, unclaimed, 1)
		compareOCRJobSpecs(t, *dbSpec, unclaimed[0])
		require.Equal(t, int32(1), unclaimed[0].ID)
		require.Equal(t, int32(1), unclaimed[0].OffchainreportingOracleSpecID)
		require.Equal(t, int32(1), unclaimed[0].PipelineSpecID)
		require.Equal(t, int32(1), unclaimed[0].OffchainreportingOracleSpec.ID)

		ocrSpec2, dbSpec2 := makeOCRJobSpec(t, db)
		err = orm.CreateJob(context.Background(), dbSpec2, ocrSpec2.TaskDAG())
		require.NoError(t, err)

		ctx2, cancel2 := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel2()

		unclaimed, err = orm2.ClaimUnclaimedJobs(ctx2)
		require.NoError(t, err)
		require.Len(t, unclaimed, 1)
		compareOCRJobSpecs(t, *dbSpec2, unclaimed[0])
		require.Equal(t, int32(2), unclaimed[0].ID)
		require.Equal(t, int32(2), unclaimed[0].OffchainreportingOracleSpecID)
		require.Equal(t, int32(2), unclaimed[0].PipelineSpecID)
		require.Equal(t, int32(2), unclaimed[0].OffchainreportingOracleSpec.ID)
	})

	t.Run("it can delete jobs claimed by other nodes", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := orm2.DeleteJob(ctx, dbSpec.ID)
		require.NoError(t, err)

		var dbSpecs []models.JobSpecV2
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

		var dbSpecs []models.JobSpecV2
		err = db.Find(&dbSpecs).Error
		require.NoError(t, err)
		require.Len(t, dbSpecs, 1)

		var oracleSpecs []models.OffchainReportingOracleSpec
		err = db.Find(&oracleSpecs).Error
		require.NoError(t, err)
		require.Len(t, oracleSpecs, 1)

		var pipelineSpecs []pipeline.Spec
		err = db.Find(&pipelineSpecs).Error
		require.NoError(t, err)
		require.Len(t, pipelineSpecs, 1)

		var pipelineTaskSpecs []pipeline.TaskSpec
		err = db.Find(&pipelineTaskSpecs).Error
		require.NoError(t, err)
		require.Len(t, pipelineTaskSpecs, 9) // 8 explicitly-defined tasks + 1 automatically added ResultTask
	})

	t.Run("increase job spec error occurrence", func(t *testing.T) {
		ocrSpec3, dbSpec3 := makeOCRJobSpec(t, db)
		err := orm.CreateJob(context.Background(), dbSpec3, ocrSpec3.TaskDAG())
		require.NoError(t, err)
		var jobSpec models.JobSpecV2
		err = db.
			First(&jobSpec).
			Error
		require.NoError(t, err)

		ocrSpecError1 := "ocr spec 1 errored"
		ocrSpecError2 := "ocr spec 2 errored"
		orm.RecordError(context.Background(), jobSpec.ID, ocrSpecError1)
		orm.RecordError(context.Background(), jobSpec.ID, ocrSpecError1)
		orm.RecordError(context.Background(), jobSpec.ID, ocrSpecError2)

		var specErrors []models.JobSpecErrorV2
		err = db.Find(&specErrors).Error
		require.NoError(t, err)
		require.Len(t, specErrors, 2)

		assert.Equal(t, specErrors[0].Occurrences, uint(2))
		assert.Equal(t, specErrors[1].Occurrences, uint(1))
		assert.True(t, specErrors[0].CreatedAt.Before(specErrors[0].UpdatedAt))
		assert.Equal(t, specErrors[0].Description, ocrSpecError1)
		assert.Equal(t, specErrors[1].Description, ocrSpecError2)
		assert.True(t, specErrors[1].CreatedAt.After(specErrors[0].UpdatedAt))
		var j2 models.JobSpecV2
		err = db.
			Preload("OffchainreportingOracleSpec").
			Preload("JobSpecErrors").
			First(&j2, "jobs.id = ?", jobSpec.ID).
			Error
		require.NoError(t, err)
	})
}

func TestORM_CheckForDeletedJobs(t *testing.T) {
	t.Parallel()

	config, cleanup := cltest.NewConfig(t)
	defer cleanup()
	store, cleanup := cltest.NewStoreWithConfig(config)
	defer cleanup()
	db := store.DB

	pipelineORM, eventBroadcaster, cleanupORM := cltest.NewPipelineORM(t, config, db)
	defer cleanupORM()

	orm := job.NewORM(db, config, pipelineORM, eventBroadcaster, &postgres.NullAdvisoryLocker{})
	defer orm.Close()

	claimedJobs := make([]models.JobSpecV2, 3)
	for i := range claimedJobs {
		ocrSpec, dbSpec := makeOCRJobSpec(t, db)
		require.NoError(t, orm.CreateJob(context.Background(), dbSpec, ocrSpec.TaskDAG()))
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
	store, cleanup := cltest.NewStoreWithConfig(config)
	defer cleanup()
	db := store.DB

	pipelineORM, eventBroadcaster, cleanupORM := cltest.NewPipelineORM(t, config, db)
	defer cleanupORM()

	advisoryLocker := new(mocks.AdvisoryLocker)
	orm := job.NewORM(db, config, pipelineORM, eventBroadcaster, advisoryLocker)
	defer orm.Close()

	require.NoError(t, orm.UnclaimJob(context.Background(), 42))

	claimedJobs := make([]models.JobSpecV2, 3)
	for i := range claimedJobs {
		_, dbSpec := makeOCRJobSpec(t, db)
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
