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

func TestBase64DecodeTask(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		input  string
		result []byte
		error  string
	}{

		// success
		{"happy", "SGVsbG8sIHBsYXlncm91bmQ=", []byte("Hello, playground"), ""},
		{"empty input", "", []byte{}, ""},

		// failure
		{"wrong characters", "S.G_VsbG8sIHBsYXlncm91bmQ=", nil, "failed to decode base64 string"},
	}

	for _, test := range tests {
		assertOK := func(result pipeline.Result, runInfo pipeline.RunInfo) {
			assert.False(t, runInfo.IsPending)
			assert.False(t, runInfo.IsRetryable)
			if test.error == "" {
				require.NoError(t, result.Error)
				require.Equal(t, test.result, result.Value)
			} else {
				require.ErrorContains(t, result.Error, test.error)
			}
		}
		t.Run(test.name, func(t *testing.T) {
			t.Run("without vars through job DAG", func(t *testing.T) {
				vars := pipeline.NewVarsFrom(nil)
				task := pipeline.Base64DecodeTask{BaseTask: pipeline.NewBaseTask(0, "task", nil, nil, 0)}
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
				task := pipeline.Base64DecodeTask{
					BaseTask: pipeline.NewBaseTask(0, "task", nil, nil, 0),
					Input:    inputStr,
				}
				assertOK(task.Run(testutils.Context(t), logger.TestLogger(t), vars, []pipeline.Result{}))
			})
			t.Run("with vars", func(t *testing.T) {
				vars := pipeline.NewVarsFrom(map[string]interface{}{
					"foo": map[string]interface{}{"bar": test.input},
				})
				task := pipeline.Base64DecodeTask{
					BaseTask: pipeline.NewBaseTask(0, "task", nil, nil, 0),
					Input:    "$(foo.bar)",
				}
				assertOK(task.Run(testutils.Context(t), logger.TestLogger(t), vars, []pipeline.Result{}))
			})
		})
	}
}
