package pipeline_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
)

func TestConditionalTask(t *testing.T) {
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
		t.Run(test.name, func(t *testing.T) {
			t.Run("without vars", func(t *testing.T) {
				vars := pipeline.NewVarsFrom(nil)
				task := pipeline.ConditionalTask{
					BaseTask: pipeline.NewBaseTask(0, "task", nil, nil, 0),
					Data:     test.input.(string)}
				result, runInfo := task.Run(testutils.Context(t), logger.TestLogger(t), vars, []pipeline.Result{{Value: test.input}})

				assert.False(t, runInfo.IsPending)
				assert.False(t, runInfo.IsRetryable)
				if test.expectErr {
					require.Error(t, result.Error)
					require.Nil(t, result.Value)
				} else {
					require.NoError(t, result.Error)
					require.True(t, result.Value.(bool))
				}
			})
			t.Run("with vars", func(t *testing.T) {
				vars := pipeline.NewVarsFrom(map[string]interface{}{
					"foo": map[string]interface{}{"bar": test.input},
				})
				task := pipeline.ConditionalTask{
					BaseTask: pipeline.NewBaseTask(0, "task", nil, nil, 0),
					Data:     "$(foo.bar)",
				}
				result, runInfo := task.Run(testutils.Context(t), logger.TestLogger(t), vars, []pipeline.Result{})

				assert.False(t, runInfo.IsPending)
				assert.False(t, runInfo.IsRetryable)
				if test.expectErr {
					require.Error(t, result.Error)
					require.Nil(t, result.Value)
				} else {
					require.NoError(t, result.Error)
					require.True(t, result.Value.(bool))
				}
			})
		})
	}
}
