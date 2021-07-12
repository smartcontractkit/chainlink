package pipeline_test

import (
	"context"
	"testing"
	"time"

	"github.com/bmizerany/assert"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func Test_PipelineORM_CreateSpec(t *testing.T) {
	db, orm := setupORM(t)

	var (
		source          = ""
		maxTaskDuration = models.Interval(1 * time.Minute)
	)

	p := pipeline.Pipeline{
		Source: source,
	}

	id, err := orm.CreateSpec(context.Background(), db, p, maxTaskDuration)
	require.NoError(t, err)

	actual := pipeline.Spec{}
	err = db.Find(&actual, id).Error
	require.NoError(t, err)
	assert.Equal(t, source, actual.DotDagSource)
	assert.Equal(t, maxTaskDuration, actual.MaxTaskDuration)
}

func Test_PipelineORM_FindRun(t *testing.T) {
	db, orm := setupORM(t)

	require.NoError(t, db.Exec(`SET CONSTRAINTS pipeline_runs_pipeline_spec_id_fkey DEFERRED`).Error)
	expected := mustInsertPipelineRun(t, db)

	run, err := orm.FindRun(expected.ID)
	require.NoError(t, err)

	require.Equal(t, expected.ID, run.ID)
}

func mustInsertPipelineRun(t *testing.T, db *gorm.DB) pipeline.Run {
	t.Helper()

	run := pipeline.Run{
		Outputs:    pipeline.JSONSerializable{Null: true},
		Errors:     pipeline.RunErrors{},
		FinishedAt: nil,
	}
	require.NoError(t, db.Create(&run).Error)
	return run
}

func setupORM(t *testing.T) (*gorm.DB, pipeline.ORM) {
	t.Helper()

	db := pgtest.NewGormDB(t)
	orm := pipeline.NewORM(db)

	return db, orm
}
