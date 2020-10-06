package job_test

import (
	"context"
	"testing"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/core/store/models"
	ormpkg "github.com/smartcontractkit/chainlink/core/store/orm"
)

func TestORM(t *testing.T) {
	config, oldORM, cleanupDB := cltest.BootstrapThrowawayORM(t, "chainlink_test_temp_job_orm", true)
	defer cleanupDB()
	db := oldORM.DB

	orm := job.NewORM(db, config.DatabaseURL(), pipeline.NewORM(db, config.DatabaseURL()))
	defer orm.Close()

	ocrSpec, dbSpec := makeOCRJobSpec(t, db)

	t.Run("it creates job specs", func(t *testing.T) {
		err := orm.CreateJob(dbSpec, ocrSpec.TaskDAG())
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

	orm2 := job.NewORM(db2, config.DatabaseURL(), pipeline.NewORM(db2, config.DatabaseURL()))
	defer orm2.Close()

	t.Run("it correctly returns the unclaimed jobs in the DB", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		unclaimed, err := orm.ClaimUnclaimedJobs(ctx)
		require.NoError(t, err)
		require.Len(t, unclaimed, 1)
		compareOCRJobSpecs(t, *dbSpec, unclaimed[0])

		ocrSpec2, dbSpec2 := makeOCRJobSpec(t, db)
		err = orm.CreateJob(dbSpec2, ocrSpec2.TaskDAG())
		require.NoError(t, err)

		ctx2, cancel2 := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel2()

		unclaimed, err = orm2.ClaimUnclaimedJobs(ctx2)
		require.NoError(t, err)
		require.Len(t, unclaimed, 1)
		compareOCRJobSpecs(t, *dbSpec, unclaimed[0])
	})

	t.Run("it cannot delete jobs claimed by other nodes", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := orm2.DeleteJob(ctx, dbSpec.ID)
		require.Error(t, err)
	})

	t.Run("it deletes its own claimed jobs from the DB", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := orm.DeleteJob(ctx, dbSpec.ID)
		require.NoError(t, err)

		var dbSpecs []models.JobSpecV2
		err = db.Find(&dbSpecs).Error
		require.NoError(t, err)
		require.Len(t, dbSpecs, 0)

		var oracleSpecs []models.OffchainReportingOracleSpec
		err = db.Find(&oracleSpecs).Error
		require.NoError(t, err)
		require.Len(t, oracleSpecs, 0)

		var pipelineSpecs []pipeline.Spec
		err = db.Find(&pipelineSpecs).Error
		require.NoError(t, err)
		require.Len(t, pipelineSpecs, 0)

		var pipelineTaskSpecs []pipeline.TaskSpec
		err = db.Find(&pipelineTaskSpecs).Error
		require.NoError(t, err)
		require.Len(t, pipelineTaskSpecs, 0)
	})
}
