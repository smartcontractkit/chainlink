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

func TestHexDecodeTask(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		input  string
		result []byte
		error  string
	}{

		// success
		{"happy", "0x12345678", []byte{0x12, 0x34, 0x56, 0x78}, ""},
		{"happy zero", "0x00", []byte{0}, ""},

		// failure
		{"missing hex prefix", "12345678", nil, "hex string must have prefix 0x"},
		{"empty input", "", nil, "hex string must have prefix 0x"},
		{"wrong alphabet", "0xwq", nil, "failed to decode hex string"},
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
				task := pipeline.HexDecodeTask{BaseTask: pipeline.NewBaseTask(0, "task", nil, nil, 0)}
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
				task := pipeline.HexDecodeTask{
					BaseTask: pipeline.NewBaseTask(0, "task", nil, nil, 0),
					Input:    inputStr,
				}
				assertOK(task.Run(testutils.Context(t), logger.TestLogger(t), vars, []pipeline.Result{}))
			})
			t.Run("with vars", func(t *testing.T) {
				vars := pipeline.NewVarsFrom(map[string]interface{}{
					"foo": map[string]interface{}{"bar": test.input},
				})
				task := pipeline.HexDecodeTask{
					BaseTask: pipeline.NewBaseTask(0, "task", nil, nil, 0),
					Input:    "$(foo.bar)",
				}
				assertOK(task.Run(testutils.Context(t), logger.TestLogger(t), vars, []pipeline.Result{}))
			})
		})
	}
}
