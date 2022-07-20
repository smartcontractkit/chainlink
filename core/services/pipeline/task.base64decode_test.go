package pipeline_test

import (
	"context"
	"testing"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTask_Base64Decode(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		input  string
		result []byte
		error  string
	}{

		// succsess
		{"happy", "SGVsbG8sIHBsYXlncm91bmQ=", []byte("Hello, playground"), ""},
		{"empty input", "", []byte{}, ""},

		// failure
		{"missing characters", "S.G_VsbG8sIHBsYXlncm91bmQ=", nil, "failed to decode base64 string"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			vars := pipeline.NewVarsFrom(nil)
			task := pipeline.Base64DecodeTask{BaseTask: pipeline.NewBaseTask(0, "task", nil, nil, 0), Input: test.input}
			result, runInfo := task.Run(context.Background(), logger.TestLogger(t), vars, []pipeline.Result{{Value: test.input}})

			assert.False(t, runInfo.IsPending)
			assert.False(t, runInfo.IsRetryable)

			if test.error == "" {
				require.NoError(t, result.Error)
				require.Equal(t, test.result, result.Value)
			} else {
				require.ErrorContains(t, result.Error, test.error)
			}
		})
	}

	for _, test := range tests {
		t.Run(test.name+" (with pipeline.Vars)", func(t *testing.T) {
			vars := pipeline.NewVarsFrom(map[string]interface{}{
				"foo": map[string]interface{}{"bar": test.input},
			})
			task := pipeline.Base64DecodeTask{
				BaseTask: pipeline.NewBaseTask(0, "task", nil, nil, 0),
				Input:    "$(foo.bar)",
			}
			result, runInfo := task.Run(context.Background(), logger.TestLogger(t), vars, []pipeline.Result{})

			assert.False(t, runInfo.IsPending)
			assert.False(t, runInfo.IsRetryable)

			if test.error == "" {
				require.NoError(t, result.Error)
				require.Equal(t, test.result, result.Value)
			} else {
				require.ErrorContains(t, result.Error, test.error)
			}
		})
	}
}
