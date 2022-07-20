package pipeline_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
)

func TestTask_Conditional_Success(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		input     interface{}
		expectErr bool
	}{
		{"true string", "true", false},
		{"false string", "false", true},
		{"empty string", "", true},
		{"0 string", "0", true},
		{"1 string", "1", false},
		{"abc string", "abc", true},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			vars := pipeline.NewVarsFrom(nil)
			task := pipeline.ConditionalTask{
				BaseTask: pipeline.NewBaseTask(0, "task", nil, nil, 0),
				Data:     test.input.(string)}
			result, runInfo := task.Run(context.Background(), logger.TestLogger(t), vars, []pipeline.Result{{Value: test.input}})

			assert.False(t, runInfo.IsPending)
			assert.False(t, runInfo.IsRetryable)
			if test.expectErr {
				require.Error(t, result.Error)
				require.Equal(t, nil, result.Value)
			} else {
				require.NoError(t, result.Error)
				require.Equal(t, true, result.Value.(bool))
			}
		})
	}

	for _, test := range tests {
		test := test
		t.Run(test.name+" (with pipeline.Vars)", func(t *testing.T) {
			vars := pipeline.NewVarsFrom(map[string]interface{}{
				"foo": map[string]interface{}{"bar": test.input},
			})
			task := pipeline.ConditionalTask{
				BaseTask: pipeline.NewBaseTask(0, "task", nil, nil, 0),
				Data:     "$(foo.bar)",
			}
			result, runInfo := task.Run(context.Background(), logger.TestLogger(t), vars, []pipeline.Result{})

			assert.False(t, runInfo.IsPending)
			assert.False(t, runInfo.IsRetryable)
			if test.expectErr {
				require.Error(t, result.Error)
				require.Equal(t, nil, result.Value)
			} else {
				require.NoError(t, result.Error)
				require.Equal(t, true, result.Value.(bool))
			}
		})
	}
}
