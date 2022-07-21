package pipeline_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
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
		t.Run(test.name, func(t *testing.T) {
			t.Run("without vars", func(t *testing.T) {
				vars := pipeline.NewVarsFrom(nil)
				task := pipeline.HexDecodeTask{BaseTask: pipeline.NewBaseTask(0, "task", nil, nil, 0), Input: test.input}
				result, runInfo := task.Run(testutils.Context(t), logger.TestLogger(t), vars, []pipeline.Result{{Value: test.input}})

				assert.False(t, runInfo.IsPending)
				assert.False(t, runInfo.IsRetryable)

				if test.error == "" {
					require.NoError(t, result.Error)
					require.Equal(t, test.result, result.Value)
				} else {
					require.ErrorContains(t, result.Error, test.error)
				}
			})
			t.Run("with vars", func(t *testing.T) {
				vars := pipeline.NewVarsFrom(map[string]interface{}{
					"foo": map[string]interface{}{"bar": test.input},
				})
				task := pipeline.HexDecodeTask{
					BaseTask: pipeline.NewBaseTask(0, "task", nil, nil, 0),
					Input:    "$(foo.bar)",
				}
				result, runInfo := task.Run(testutils.Context(t), logger.TestLogger(t), vars, []pipeline.Result{})

				assert.False(t, runInfo.IsPending)
				assert.False(t, runInfo.IsRetryable)

				if test.error == "" {
					require.NoError(t, result.Error)
					require.Equal(t, test.result, result.Value)
				} else {
					require.ErrorContains(t, result.Error, test.error)
				}
			})
		})
	}
}
