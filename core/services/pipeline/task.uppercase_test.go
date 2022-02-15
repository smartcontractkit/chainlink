package pipeline_test

import (
	"context"
	"testing"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTask_Uppercase_Success(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input interface{}
		want  string
	}{
		{"uppercase string", "UPPERCASE", "UPPERCASE"},
		{"camelCase string", "camelCase", "CAMELCASE"},
		{"PascalCase string", "PascalCase", "PASCALCASE"},
		{"mixed string", "mIxEd", "MIXED"},
		{"lowercase string", "lowercase", "LOWERCASE"},
		{"empty string", "", ""},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			vars := pipeline.NewVarsFrom(nil)
			task := pipeline.UppercaseTask{BaseTask: pipeline.NewBaseTask(0, "task", nil, nil, 0), Input: test.input.(string)}
			result, runInfo := task.Run(context.Background(), logger.TestLogger(t), vars, []pipeline.Result{{Value: test.input}})

			assert.False(t, runInfo.IsPending)
			assert.False(t, runInfo.IsRetryable)
			require.NoError(t, result.Error)
			require.Equal(t, test.want, result.Value.(string))
		})
	}

	for _, test := range tests {
		test := test
		t.Run(test.name+" (with pipeline.Vars)", func(t *testing.T) {
			vars := pipeline.NewVarsFrom(map[string]interface{}{
				"foo": map[string]interface{}{"bar": test.input},
			})
			task := pipeline.UppercaseTask{
				BaseTask: pipeline.NewBaseTask(0, "task", nil, nil, 0),
				Input:    "$(foo.bar)",
			}
			result, runInfo := task.Run(context.Background(), logger.TestLogger(t), vars, []pipeline.Result{})
			assert.False(t, runInfo.IsPending)
			assert.False(t, runInfo.IsRetryable)
			require.NoError(t, result.Error)
			require.Equal(t, test.want, result.Value.(string))
		})
	}
}
