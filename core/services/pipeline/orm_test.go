package pipeline_test

import (
	"sort"
	"testing"

	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v4"
)

func Test_Pipeline_TaskRunsToQueueItems(t *testing.T) {
	taskRuns := []pipeline.TaskRun{
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
	}

	qis := pipeline.TaskRunsToQueueItems(taskRuns)

	require.Len(t, qis, 3)

	// QueueItems are generated in a non-deterministic ordering so simply sort
	// by first element in PipelineTaskRunIDs in order to make assertions
	sort.SliceStable(qis, func(i, j int) bool { return qis[i].PipelineTaskRunIDs[0] < qis[j].PipelineTaskRunIDs[0] })

	require.Equal(t, []int64{11, 12, 14, 15}, qis[0].PipelineTaskRunIDs)
	require.Equal(t, []int64{13}, qis[1].PipelineTaskRunIDs)
	require.Equal(t, []int64{17, 16}, qis[2].PipelineTaskRunIDs)
}
