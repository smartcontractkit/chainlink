package pipeline_test

import (
	"context"
	"sort"
	"testing"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/core/services/postgres"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v4"
)

func Test_PipelineORM_CreateRun(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()
	db := store.DB

	job := cltest.MustInsertSampleDirectRequestJob(t, db)

	eventBroadcaster := postgres.NewEventBroadcaster(store.Config.DatabaseURL(), 0, 0)

	orm := pipeline.NewORM(db, store.Config, eventBroadcaster)

	meta := make(map[string]interface{})

	runID, err := orm.CreateRun(context.Background(), job.ID, meta)
	require.NoError(t, err)

	// Check that JobRun, TaskRuns and QueueItems were created

	var prs []pipeline.Run
	var trs []pipeline.TaskRun
	var qis []pipeline.QueueItem

	require.NoError(t, db.Find(&prs).Error)
	require.NoError(t, db.Find(&trs).Error)
	require.NoError(t, db.Find(&qis).Error)

	require.Len(t, prs, 1)
	require.Equal(t, runID, prs[0].ID)
	require.Len(t, trs, 4)
	require.Len(t, qis, 2)
}

func Test_PipelineORM_TaskRunsToQueueItems(t *testing.T) {
	taskRuns := []pipeline.TaskRun{
		pipeline.TaskRun{
			ID: 10,
			PipelineTaskSpec: pipeline.TaskSpec{
				ID:          0,
				SuccessorID: null.IntFrom(2),
			},
		},
		pipeline.TaskRun{
			ID: 11,
			PipelineTaskSpec: pipeline.TaskSpec{
				ID:          1,
				SuccessorID: null.IntFrom(2),
			},
		},
		pipeline.TaskRun{
			ID: 12,
			PipelineTaskSpec: pipeline.TaskSpec{
				ID:          2,
				SuccessorID: null.IntFrom(4),
			},
		},
		pipeline.TaskRun{
			ID: 13,
			PipelineTaskSpec: pipeline.TaskSpec{
				ID:          3,
				SuccessorID: null.Int{},
			},
		},
		pipeline.TaskRun{
			ID: 14,
			PipelineTaskSpec: pipeline.TaskSpec{
				ID:          4,
				SuccessorID: null.IntFrom(5),
			},
		},
		pipeline.TaskRun{
			ID: 15,
			PipelineTaskSpec: pipeline.TaskSpec{
				ID:          5,
				SuccessorID: null.Int{},
			},
		},
		pipeline.TaskRun{
			ID: 16,
			PipelineTaskSpec: pipeline.TaskSpec{
				ID:          6,
				SuccessorID: null.Int{},
			},
		},
		pipeline.TaskRun{
			ID: 17,
			PipelineTaskSpec: pipeline.TaskSpec{
				ID:          7,
				SuccessorID: null.IntFrom(6),
			},
		},
		pipeline.TaskRun{
			ID: 18,
			PipelineTaskSpec: pipeline.TaskSpec{
				ID:          8,
				SuccessorID: null.IntFrom(5),
			},
		},
	}

	qis := pipeline.TaskRunsToQueueItems(taskRuns)

	require.Len(t, qis, 4)

	// QueueItems are generated in a non-deterministic ordering so simply sort
	// by first element in PipelineTaskRunIDs in order to make assertions
	sort.SliceStable(qis, func(i, j int) bool { return qis[i].PipelineTaskRunIDs[0] < qis[j].PipelineTaskRunIDs[0] })

	require.Equal(t, []int64{11, 12, 14, 15}, []int64(qis[0].PipelineTaskRunIDs))
	require.Equal(t, []int64{13}, []int64(qis[1].PipelineTaskRunIDs))
	require.Equal(t, []int64{17, 16}, []int64(qis[2].PipelineTaskRunIDs))
	require.Equal(t, []int64{18, 15}, []int64(qis[3].PipelineTaskRunIDs))
}
