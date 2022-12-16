package pipeline_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
)

func TestBase64EncodeTask(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		input  interface{}
		result string
		error  string
	}{

		// success
		{"string input 1", "Hello, playground", "SGVsbG8sIHBsYXlncm91bmQ=", ""},
		{"string input 2", "=test=test=", "PXRlc3Q9dGVzdD0=", ""},
		{"empty string", "", "", ""},
		{"bytes input 1", []byte{0xaa, 0xbb, 0xcc, 0xdd}, "qrvM3Q==", ""},
		{"empty bytes", []byte{}, "", ""},

		// failure (unsupported types)
		{"int", 234, "", "bad input for task"},
		{"bool", false, "", "bad input for task"},
		{"float", 3.14, "", "bad input for task"},
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
				task := pipeline.Base64EncodeTask{BaseTask: pipeline.NewBaseTask(0, "task", nil, nil, 0)}
				assertOK(task.Run(testutils.Context(t), logger.TestLogger(t), vars, []pipeline.Result{{Value: test.input}}))
			})
			t.Run("with vars", func(t *testing.T) {
				vars := pipeline.NewVarsFrom(map[string]interface{}{
					"foo": map[string]interface{}{"bar": test.input},
				})
				task := pipeline.Base64EncodeTask{
					BaseTask: pipeline.NewBaseTask(0, "task", nil, nil, 0),
					Input:    "$(foo.bar)",
				}
				assertOK(task.Run(testutils.Context(t), logger.TestLogger(t), vars, []pipeline.Result{}))
			})
		})
	}
}

func TestBase64EncodeTaskInputParamLiteral(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		input  interface{}
		result string
	}{
		// Only strings can be passed via input param literals (other types will get converted to strings anyway)
		{"string input 1", "Hello, playground", "SGVsbG8sIHBsYXlncm91bmQ="},
		{"string input 2", "=test=test=", "PXRlc3Q9dGVzdD0="},
		{"int gets converted to a string", 234, "MjM0"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			vars := pipeline.NewVarsFrom(nil)
			task := pipeline.Base64EncodeTask{
				BaseTask: pipeline.NewBaseTask(0, "task", nil, nil, 0),
				Input:    fmt.Sprintf("%v", test.input),
			}
			result, runInfo := task.Run(testutils.Context(t), logger.TestLogger(t), vars, []pipeline.Result{})
			assert.False(t, runInfo.IsPending)
			assert.False(t, runInfo.IsRetryable)
			require.NoError(t, result.Error)
			require.Equal(t, test.result, result.Value)
		})
	}
}
