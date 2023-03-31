package pipeline_test

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
)

func TestModeTask(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		inputs          []pipeline.Result
		allowedFaults   string
		wantResults     []interface{}
		wantOccurrences uint64
		wantErrorCause  error
	}{
		{
			"happy (one winner)",
			[]pipeline.Result{{Value: "foo"}, {Value: "foo"}, {Value: "bar"}, {Value: true}},
			"1",
			[]interface{}{"foo"}, 2, nil,
		},
		{
			"happy (multiple winners)",
			[]pipeline.Result{{Value: "foo"}, {Value: "foo"}, {Value: "bar"}, {Value: "bar"}},
			"1",
			[]interface{}{"foo", "bar"}, 2, nil,
		},
		{
			"happy (one winner expressed as different types)",
			[]pipeline.Result{{Value: mustDecimal(t, "1.234")}, {Value: float64(1.234)}, {Value: float32(1.234)}, {Value: "1.234"}},
			"1",
			[]interface{}{"1.234"}, 4, nil,
		},
		{
			"one input",
			[]pipeline.Result{{Value: common.Address{1}}},
			"0",
			[]interface{}{common.Address{1}}, 1, nil,
		},
		{
			"zero inputs",
			[]pipeline.Result{},
			"0",
			nil, 0, pipeline.ErrWrongInputCardinality,
		},
		{
			"fewer errors than threshold",
			[]pipeline.Result{{Error: errors.New("")}, {Value: mustDecimal(t, "2")}, {Value: mustDecimal(t, "2")}, {Value: []byte("foo bar")}},
			"2",
			[]interface{}{mustDecimal(t, "2")}, 2, nil,
		},
		{
			"exactly threshold of errors",
			[]pipeline.Result{{Error: errors.New("")}, {Error: errors.New("")}, {Value: []interface{}{1, 2, 3}}, {Value: []interface{}{1, 2, 3}}},
			"2",
			[]interface{}{[]interface{}{1, 2, 3}}, 2, nil,
		},
		{
			"more errors than threshold",
			[]pipeline.Result{{Error: errors.New("")}, {Error: errors.New("")}, {Error: errors.New("")}, {Value: mustDecimal(t, "4")}},
			"2",
			nil, 0, pipeline.ErrTooManyErrors,
		},
		{
			"(unspecified AllowedFaults) fewer errors than threshold",
			[]pipeline.Result{{Error: errors.New("")}, {Error: errors.New("")}, {Value: big.NewInt(123)}, {Value: big.NewInt(123)}},
			"",
			[]interface{}{big.NewInt(123)}, 2, nil,
		},
		{
			"(unspecified AllowedFaults) exactly threshold of errors",
			[]pipeline.Result{{Error: errors.New("")}, {Error: errors.New("")}, {Error: errors.New("")}, {Value: 123}},
			"",
			[]interface{}{123}, 1, nil,
		},
		{
			"(unspecified AllowedFaults) more errors than threshold",
			[]pipeline.Result{{Error: errors.New("")}, {Error: errors.New("")}, {Error: errors.New("")}},
			"",
			nil, 0, pipeline.ErrTooManyErrors,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Run("without vars", func(t *testing.T) {
				task := pipeline.ModeTask{
					BaseTask:      pipeline.NewBaseTask(0, "mode", nil, nil, 0),
					AllowedFaults: test.allowedFaults,
				}
				output, runInfo := task.Run(testutils.Context(t), logger.TestLogger(t), pipeline.NewVarsFrom(nil), test.inputs)
				assert.False(t, runInfo.IsPending)
				assert.False(t, runInfo.IsRetryable)
				if output.Error != nil {
					require.Equal(t, test.wantErrorCause, errors.Cause(output.Error))
					require.Nil(t, output.Value)
				} else {
					require.Equal(t, map[string]interface{}{
						"results":     test.wantResults,
						"occurrences": test.wantOccurrences,
					}, output.Value)
					require.NoError(t, output.Error)
				}
			})
			t.Run("with vars", func(t *testing.T) {
				var inputs []interface{}
				for _, input := range test.inputs {
					if input.Error != nil {
						inputs = append(inputs, input.Error)
					} else {
						inputs = append(inputs, input.Value)
					}
				}
				vars := pipeline.NewVarsFrom(map[string]interface{}{
					"foo": map[string]interface{}{"bar": inputs},
				})
				task := pipeline.ModeTask{
					BaseTask:      pipeline.NewBaseTask(0, "task", nil, nil, 0),
					Values:        "$(foo.bar)",
					AllowedFaults: test.allowedFaults,
				}
				output, runInfo := task.Run(testutils.Context(t), logger.TestLogger(t), vars, nil)
				assert.False(t, runInfo.IsPending)
				assert.False(t, runInfo.IsRetryable)
				if output.Error != nil {
					require.Equal(t, test.wantErrorCause, errors.Cause(output.Error))
					require.Nil(t, output.Value)
				} else {
					require.Equal(t, map[string]interface{}{
						"results":     test.wantResults,
						"occurrences": test.wantOccurrences,
					}, output.Value)
					require.NoError(t, output.Error)
				}
			})
			t.Run("with json vars", func(t *testing.T) {
				var inputs []interface{}
				for _, input := range test.inputs {
					if input.Error != nil {
						inputs = append(inputs, input.Error)
					} else {
						inputs = append(inputs, input.Value)
					}
				}
				var valuesParam string
				var vars pipeline.Vars
				switch len(inputs) {
				case 0:
					valuesParam = "[]"
					vars = pipeline.NewVarsFrom(nil)
				case 1:
					valuesParam = "[ $(foo) ]"
					vars = pipeline.NewVarsFrom(map[string]interface{}{"foo": inputs[0]})
				case 3:
					valuesParam = "[ $(foo), $(bar), $(chain) ]"
					vars = pipeline.NewVarsFrom(map[string]interface{}{"foo": inputs[0], "bar": inputs[1], "chain": inputs[2]})
				case 4:
					valuesParam = "[ $(foo), $(bar), $(chain), $(link) ]"
					vars = pipeline.NewVarsFrom(map[string]interface{}{"foo": inputs[0], "bar": inputs[1], "chain": inputs[2], "link": inputs[3]})
				}

				task := pipeline.ModeTask{
					BaseTask:      pipeline.NewBaseTask(0, "task", nil, nil, 0),
					Values:        valuesParam,
					AllowedFaults: test.allowedFaults,
				}
				output, runInfo := task.Run(testutils.Context(t), logger.TestLogger(t), vars, nil)
				assert.False(t, runInfo.IsPending)
				assert.False(t, runInfo.IsRetryable)
				if output.Error != nil {
					require.Equal(t, test.wantErrorCause, errors.Cause(output.Error))
					require.Nil(t, output.Value)
				} else {
					require.Equal(t, map[string]interface{}{
						"results":     test.wantResults,
						"occurrences": test.wantOccurrences,
					}, output.Value)
					require.NoError(t, output.Error)
				}
			})
		})
	}
}
