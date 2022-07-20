package pipeline_test

import (
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink/core/services/pipeline"
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
				AllErrors:   pipeline.RunErrors{},
				FatalErrors: pipeline.RunErrors{},
				Outputs:     pipeline.JSONSerializable{},
				FinishedAt:  null.Time{},
			},
			want: pipeline.RunStatusRunning,
		},
		{
			name: "Completed",
			run: &pipeline.Run{
				AllErrors:   pipeline.RunErrors{},
				FatalErrors: pipeline.RunErrors{},
				Outputs:     pipeline.JSONSerializable{Val: []interface{}{10, 10}, Valid: true},
				FinishedAt:  now,
			},
			want: pipeline.RunStatusCompleted,
		},
		{
			name: "Error",
			run: &pipeline.Run{
				AllErrors:   pipeline.RunErrors{null.StringFrom(errors.New("fail").Error())},
				FatalErrors: pipeline.RunErrors{null.StringFrom(errors.New("fail").Error())},
				Outputs:     pipeline.JSONSerializable{},
				FinishedAt:  null.Time{},
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

func TestRunErrors_ToError(t *testing.T) {
	runErrors := pipeline.RunErrors{}
	runErrors = append(runErrors, null.NewString("bad thing happened", true))
	runErrors = append(runErrors, null.NewString("pretty bad thing happened", true))
	runErrors = append(runErrors, null.NewString("", false))
	expected := errors.New("bad thing happened; pretty bad thing happened")
	require.Equal(t, expected.Error(), runErrors.ToError().Error())
}
