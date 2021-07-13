package pipeline_test

import (
	"errors"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/stretchr/testify/assert"
	"gopkg.in/guregu/null.v4"
)

func TestRunStatus(t *testing.T) {
	t.Parallel()

	assert.Equal(t, pipeline.RunStatusUnknown.Finished(), false)
	assert.Equal(t, pipeline.RunStatusRunning.Finished(), false)
	assert.Equal(t, pipeline.RunStatusCompleted.Finished(), true)
	assert.Equal(t, pipeline.RunStatusErrored.Finished(), true)

	assert.Equal(t, pipeline.RunStatusUnknown.Errored(), false)
	assert.Equal(t, pipeline.RunStatusRunning.Errored(), false)
	assert.Equal(t, pipeline.RunStatusCompleted.Errored(), false)
	assert.Equal(t, pipeline.RunStatusErrored.Errored(), true)
}

func TestRun_Status(t *testing.T) {
	now := null.TimeFrom(time.Now())

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
				FinishedAt: null.Time{},
			},
			want: pipeline.RunStatusRunning,
		},
		{
			name: "Completed",
			run: &pipeline.Run{
				Errors:     pipeline.RunErrors{},
				Outputs:    pipeline.JSONSerializable{Val: []interface{}{10, 10}, Null: false},
				FinishedAt: now,
			},
			want: pipeline.RunStatusCompleted,
		},
		{
			name: "Error",
			run: &pipeline.Run{
				Outputs:    pipeline.JSONSerializable{},
				Errors:     pipeline.RunErrors{null.StringFrom(errors.New("fail").Error())},
				FinishedAt: null.Time{},
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
