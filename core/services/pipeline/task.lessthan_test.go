package pipeline_test

import (
	"fmt"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
)

func TestLessThanTask_Happy(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input interface{}
		limit string
		want  bool
	}{
		{"string, lt 100", "1.23", "100", true},
		{"string, lt negative", "1.23", "-5", false},
		{"string, lt zero", "1.23", "0", false},
		{"string, lt large value", "1.23", "1000000000000000000", true},
		{"large string, lt large value", "10000000000000000001", "1000000000000000000", false},

		{"int, true", int(2), "100", true},
		{"int, false", int(2), "-5", false},

		{"int8, true", int8(2), "100", true},
		{"int8, false", int8(2), "-5", false},

		{"int16, true", int16(2), "100", true},
		{"int16, false", int16(2), "-5", false},

		{"int32,true", int32(2), "100", true},
		{"int32, false", int32(2), "-5", false},

		{"int64, true", int64(2), "100", true},
		{"int64, false", int64(2), "-5", false},

		{"uint, true", uint(2), "100", true},
		{"uint, false", uint(2), "-5", false},

		{"uint8, true", uint8(2), "100", true},
		{"uint8, false", uint8(2), "-5", false},

		{"uint16, true", uint16(2), "100", true},
		{"uint16, false", uint16(2), "-5", false},

		{"uint32, true", uint32(2), "100", true},
		{"uint32, false", uint32(2), "-5", false},

		{"uint64, true", uint64(2), "100", true},
		{"uint64, false", uint64(2), "-5", false},

		{"float32, true", float32(1.23), "10", true},
		{"float32, false", float32(1.23), "-5", false},

		{"float64, true", float64(1.23), "10", true},
		{"float64, false", float64(1.23), "-5", false},
	}

	for _, test := range tests {
		assertOK := func(result pipeline.Result, runInfo pipeline.RunInfo) {
			assert.False(t, runInfo.IsPending)
			assert.False(t, runInfo.IsRetryable)
			require.NoError(t, result.Error)
			require.Equal(t, test.want, result.Value.(bool))
		}
		t.Run(test.name, func(t *testing.T) {
			t.Run("without vars through job DAG", func(t *testing.T) {
				vars := pipeline.NewVarsFrom(nil)
				task := pipeline.LessThan{BaseTask: pipeline.NewBaseTask(0, "task", nil, nil, 0), Limit: test.limit}
				assertOK(task.Run(testutils.Context(t), logger.TestLogger(t), vars, []pipeline.Result{{Value: test.input}}))
			})
			t.Run("without vars through input param", func(t *testing.T) {
				vars := pipeline.NewVarsFrom(nil)
				task := pipeline.LessThan{
					BaseTask: pipeline.NewBaseTask(0, "task", nil, nil, 0),
					Input:    fmt.Sprintf("%v", test.input),
					Limit:    test.limit,
				}
				assertOK(task.Run(testutils.Context(t), logger.TestLogger(t), vars, []pipeline.Result{}))
			})
			t.Run("with vars", func(t *testing.T) {
				vars := pipeline.NewVarsFrom(map[string]interface{}{
					"foo":   map[string]interface{}{"bar": test.input},
					"chain": map[string]interface{}{"link": test.limit},
				})
				task := pipeline.LessThan{
					BaseTask: pipeline.NewBaseTask(0, "task", nil, nil, 0),
					Input:    "$(foo.bar)",
					Limit:    "$(chain.link)",
				}
				assertOK(task.Run(testutils.Context(t), logger.TestLogger(t), vars, []pipeline.Result{}))
			})
		})
	}
}

func TestLessThanTask_Unhappy(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name              string
		limit             string
		input             string
		inputs            []pipeline.Result
		vars              pipeline.Vars
		wantErrorCause    error
		wantErrorContains string
	}{
		{"map as input from inputs", "100", "", []pipeline.Result{{Value: map[string]interface{}{"chain": "link"}}}, pipeline.NewVarsFrom(nil), pipeline.ErrBadInput, "input"},
		{"map as input from var", "100", "$(foo)", nil, pipeline.NewVarsFrom(map[string]interface{}{"foo": map[string]interface{}{"chain": "link"}}), pipeline.ErrBadInput, "input"},
		{"slice as input from inputs", "100", "", []pipeline.Result{{Value: []interface{}{"chain", "link"}}}, pipeline.NewVarsFrom(nil), pipeline.ErrBadInput, "input"},
		{"slice as input from var", "100", "$(foo)", nil, pipeline.NewVarsFrom(map[string]interface{}{"foo": []interface{}{"chain", "link"}}), pipeline.ErrBadInput, "input"},
		{"input as missing var", "100", "$(foo)", nil, pipeline.NewVarsFrom(nil), pipeline.ErrKeypathNotFound, "input"},
		{"limit as missing var", "$(foo)", "", []pipeline.Result{{Value: "123"}}, pipeline.NewVarsFrom(nil), pipeline.ErrKeypathNotFound, "limit"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			task := pipeline.LessThan{
				BaseTask: pipeline.NewBaseTask(0, "task", nil, nil, 0),
				Input:    test.input,
				Limit:    test.limit,
			}
			result, runInfo := task.Run(testutils.Context(t), logger.TestLogger(t), test.vars, test.inputs)
			assert.False(t, runInfo.IsPending)
			assert.False(t, runInfo.IsRetryable)
			require.Equal(t, test.wantErrorCause, errors.Cause(result.Error))
			if test.wantErrorContains != "" {
				require.Contains(t, result.Error.Error(), test.wantErrorContains)
			}
		})
	}
}
