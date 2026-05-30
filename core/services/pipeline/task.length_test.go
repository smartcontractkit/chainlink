package pipeline_test

import (
	"math/big"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
)

func TestLengthTask(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		input      any
		inputParam string
		want       decimal.Decimal
	}{
		{"normal bytes", []byte{0xaa, 0xbb, 0xcc, 0xdd}, string([]byte{0xaa, 0xbb, 0xcc, 0xdd}), decimal.NewFromInt(4)},
		{"empty bytes", []byte{}, "", decimal.NewFromInt(0)},
		{"string as bytes", []byte("stevetoshi sergeymoto"), "stevetoshi sergeymoto", decimal.NewFromInt(21)},
		{"string input gets converted to bytes", "stevetoshi sergeymoto", "stevetoshi sergeymoto", decimal.NewFromInt(21)},
		{"empty string", "", "", decimal.NewFromInt(0)},
		{"array input", []any{1, 2, 3}, "[1,2,3]", decimal.NewFromInt(3)},
		{"empty array input", []any{}, "[]", decimal.NewFromInt(0)},
		{"typed slice input", []*big.Int{big.NewInt(3), big.NewInt(11)}, "", decimal.NewFromInt(2)},
		{"JSON array string", "[1,2,3]", "[1,2,3]", decimal.NewFromInt(3)},
		{"empty JSON array string", "[]", "[]", decimal.NewFromInt(0)},
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
				if test.inputParam == "" {
					// empty input parameter is indistinguishable from not providing it at all
					// in that case the task will use an input defined by the job DAG
					return
				}
				vars := pipeline.NewVarsFrom(nil)
				task := pipeline.LengthTask{
					BaseTask: pipeline.NewBaseTask(0, "task", nil, nil, 0),
					Input:    test.inputParam,
				}
				assertOK(task.Run(testutils.Context(t), logger.TestLogger(t), vars, []pipeline.Result{}))
			})
			t.Run("with vars", func(t *testing.T) {
				vars := pipeline.NewVarsFrom(map[string]any{
					"foo": map[string]any{"bar": test.input},
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
