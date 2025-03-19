package pipeline_test

import (
	"math/big"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink-common/pkg/utils/jsonserializable"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
)

func TestRun_Status(t *testing.T) {
	t.Parallel()

	assert.False(t, pipeline.RunStatusUnknown.Finished())
	assert.False(t, pipeline.RunStatusRunning.Finished())
	assert.True(t, pipeline.RunStatusCompleted.Finished())
	assert.True(t, pipeline.RunStatusErrored.Finished())

	assert.False(t, pipeline.RunStatusUnknown.Errored())
	assert.False(t, pipeline.RunStatusRunning.Errored())
	assert.False(t, pipeline.RunStatusCompleted.Errored())
	assert.True(t, pipeline.RunStatusErrored.Errored())

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
				Outputs:     jsonserializable.JSONSerializable{},
				FinishedAt:  null.Time{},
			},
			want: pipeline.RunStatusRunning,
		},
		{
			name: "Completed",
			run: &pipeline.Run{
				AllErrors:   pipeline.RunErrors{},
				FatalErrors: pipeline.RunErrors{},
				Outputs:     jsonserializable.JSONSerializable{Val: []interface{}{10, 10}, Valid: true},
				FinishedAt:  now,
			},
			want: pipeline.RunStatusCompleted,
		},
		{
			name: "Error",
			run: &pipeline.Run{
				AllErrors:   pipeline.RunErrors{null.StringFrom(errors.New("fail").Error())},
				FatalErrors: pipeline.RunErrors{null.StringFrom(errors.New("fail").Error())},
				Outputs:     jsonserializable.JSONSerializable{},
				FinishedAt:  null.Time{},
			},
			want: pipeline.RunStatusErrored,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
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

func TestRun_StringOutputs(t *testing.T) {
	t.Parallel()

	t.Run("invalid outputs", func(t *testing.T) {
		run := &pipeline.Run{
			Outputs: jsonserializable.JSONSerializable{
				Valid: false,
			},
		}
		outputs, err := run.StringOutputs()
		assert.NoError(t, err)
		assert.Empty(t, outputs)
	})

	big := big.NewInt(123)
	dec := mustDecimal(t, "123")

	testCases := []struct {
		name string
		val  interface{}
		want string
	}{
		{"int64", int64(123), "123"},
		{"uint64", uint64(123), "123"},
		{"float64", float64(123.456), "123.456"},
		{"large float64", float64(9007199254740991231), "9007199254740991000"},
		{"big.Int", *big, "123"},
		{"*big.Int", big, "123"},
		{"decimal", *dec, "123"},
		{"*decimal", dec, "123"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			run := &pipeline.Run{
				Outputs: jsonserializable.JSONSerializable{
					Valid: true,
					Val:   []interface{}{tc.val},
				},
			}
			t.Log(tc.val)
			outputs, err := run.StringOutputs()
			assert.NoError(t, err)
			assert.NotNil(t, outputs)
			assert.Len(t, outputs, 1)
			assert.Equal(t, tc.want, *outputs[0])
		})
	}
}
