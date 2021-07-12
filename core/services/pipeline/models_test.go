package pipeline_test

import (
	"errors"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/stretchr/testify/assert"
)

func TestRunStatus(t *testing.T) {
	t.Parallel()

	assert.Equal(t, pipeline.RunStatusUnknown.Finished(), false)
	assert.Equal(t, pipeline.RunStatusInProgress.Finished(), false)
	assert.Equal(t, pipeline.RunStatusCompleted.Finished(), true)
	assert.Equal(t, pipeline.RunStatusErrored.Finished(), true)

	assert.Equal(t, pipeline.RunStatusUnknown.Errored(), false)
	assert.Equal(t, pipeline.RunStatusInProgress.Errored(), false)
	assert.Equal(t, pipeline.RunStatusCompleted.Errored(), false)
	assert.Equal(t, pipeline.RunStatusErrored.Errored(), true)
}

func TestRun_Status(t *testing.T) {
	now := time.Now()
	var success = pipeline.TaskRunResults{
		{
			Task: &pipeline.HTTPTask{},
			Result: pipeline.Result{
				Value: 10,
				Error: nil,
			},
			FinishedAt: time.Now(),
		},
		{
			Task: &pipeline.HTTPTask{},
			Result: pipeline.Result{
				Value: 10,
				Error: nil,
			},
			FinishedAt: time.Now(),
		},
	}
	var fail = pipeline.TaskRunResults{
		{
			Task: &pipeline.HTTPTask{},
			Result: pipeline.Result{
				Value: nil,
				Error: errors.New("fail"),
			},
			FinishedAt: time.Now(),
		},
		{
			Task: &pipeline.HTTPTask{},
			Result: pipeline.Result{
				Value: nil,
				Error: errors.New("fail"),
			},
			FinishedAt: time.Now(),
		},
	}

	testCases := []struct {
		name string
		run  *pipeline.Run
		want pipeline.RunStatus
	}{
		{
			name: "In Progress",
			run: &pipeline.Run{
				Errors:     pipeline.RunErrors{},
				Outputs:    pipeline.JSONSerializable{},
				FinishedAt: nil,
			},
			want: pipeline.RunStatusInProgress,
		},
		{
			name: "Completed",
			run: &pipeline.Run{
				Errors:     success.FinalResult().ErrorsDB(),
				Outputs:    success.FinalResult().OutputsDB(),
				FinishedAt: &now,
			},
			want: pipeline.RunStatusCompleted,
		},
		{
			name: "Error",
			run: &pipeline.Run{
				Outputs:    fail.FinalResult().OutputsDB(),
				Errors:     fail.FinalResult().ErrorsDB(),
				FinishedAt: nil,
			},
			want: pipeline.RunStatusErrored,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tc.want, tc.run.Status())
		})
	}
}
