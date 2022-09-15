package pipeline_test

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
)

func TestLengthTask(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input interface{}
		want  decimal.Decimal
	}{
		{"normal bytes", []byte{0xaa, 0xbb, 0xcc, 0xdd}, decimal.NewFromInt(4)},
		{"empty bytes", []byte{}, decimal.NewFromInt(0)},
		{"string as bytes", []byte("stevetoshi sergeymoto"), decimal.NewFromInt(21)},
		{"string input gets converted to bytes", "stevetoshi sergeymoto", decimal.NewFromInt(21)},
		{"empty string", "", decimal.NewFromInt(0)},
	}

	for _, test := range tests {
		assertOK := func(result pipeline.Result, runInfo pipeline.RunInfo) {
			assert.False(t, runInfo.IsPending)
			assert.False(t, runInfo.IsRetryable)
			require.NoError(t, result.Error)
			require.Equal(t, test.want.String(), result.Value.(decimal.Decimal).String())
		}
		t.Run(test.name, func(t *testing.T) {
			t.Run("without vars through job DAG", func(t *testing.T) {
				vars := pipeline.NewVarsFrom(nil)
				task := pipeline.LengthTask{BaseTask: pipeline.NewBaseTask(0, "task", nil, nil, 0)}
				assertOK(task.Run(testutils.Context(t), logger.TestLogger(t), vars, []pipeline.Result{{Value: test.input}}))
			})
			t.Run("without vars through input param", func(t *testing.T) {
				var inputStr string
				if _, ok := test.input.([]byte); ok {
					inputStr = string(test.input.([]byte))
				} else {
					inputStr = test.input.(string)
				}
				if inputStr == "" {
					// empty input parameter is indistinguishable from not providing it at all
					// in that case the task will use an input defined by the job DAG
					return
				}
				vars := pipeline.NewVarsFrom(nil)
				task := pipeline.LengthTask{
					BaseTask: pipeline.NewBaseTask(0, "task", nil, nil, 0),
					Input:    inputStr,
				}
				assertOK(task.Run(testutils.Context(t), logger.TestLogger(t), vars, []pipeline.Result{}))
			})
			t.Run("with vars", func(t *testing.T) {
				vars := pipeline.NewVarsFrom(map[string]interface{}{
					"foo": map[string]interface{}{"bar": test.input},
				})
				task := pipeline.LengthTask{
					BaseTask: pipeline.NewBaseTask(0, "task", nil, nil, 0),
					Input:    "$(foo.bar)",
				}
				assertOK(task.Run(testutils.Context(t), logger.TestLogger(t), vars, []pipeline.Result{}))
			})
		})
	}
}
