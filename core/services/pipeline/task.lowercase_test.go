package pipeline_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
)

func TestLowercaseTask(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input interface{}
		want  string
	}{
		{"uppercase string", "UPPERCASE", "uppercase"},
		{"camelCase string", "camelCase", "camelcase"},
		{"PascalCase string", "PascalCase", "pascalcase"},
		{"mixed string", "mIxEd", "mixed"},
		{"lowercase string", "lowercase", "lowercase"},
		{"empty string", "", ""},
	}

	for _, test := range tests {
		assertOK := func(result pipeline.Result, runInfo pipeline.RunInfo) {
			assert.False(t, runInfo.IsPending)
			assert.False(t, runInfo.IsRetryable)
			require.NoError(t, result.Error)
			require.Equal(t, test.want, result.Value.(string))
		}
		t.Run(test.name, func(t *testing.T) {
			t.Run("without vars through job DAG", func(t *testing.T) {
				vars := pipeline.NewVarsFrom(nil)
				task := pipeline.LowercaseTask{BaseTask: pipeline.NewBaseTask(0, "task", nil, nil, 0)}
				assertOK(task.Run(testutils.Context(t), logger.TestLogger(t), vars, []pipeline.Result{{Value: test.input}}))
			})
			t.Run("without vars through input param", func(t *testing.T) {
				inputStr := fmt.Sprintf("%v", test.input)
				if inputStr == "" {
					// empty input parameter is indistinguishable from not providing it at all
					// in that case the task will use an input defined by the job DAG
					return
				}
				vars := pipeline.NewVarsFrom(nil)
				task := pipeline.LowercaseTask{
					BaseTask: pipeline.NewBaseTask(0, "task", nil, nil, 0),
					Input:    inputStr,
				}
				assertOK(task.Run(testutils.Context(t), logger.TestLogger(t), vars, []pipeline.Result{}))
			})
			t.Run("with vars", func(t *testing.T) {
				vars := pipeline.NewVarsFrom(map[string]interface{}{
					"foo": map[string]interface{}{"bar": test.input},
				})
				task := pipeline.LowercaseTask{
					BaseTask: pipeline.NewBaseTask(0, "task", nil, nil, 0),
					Input:    "$(foo.bar)",
				}
				assertOK(task.Run(testutils.Context(t), logger.TestLogger(t), vars, []pipeline.Result{}))
			})
		})
	}
}
