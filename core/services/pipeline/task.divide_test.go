package pipeline_test

import (
	"context"
	"testing"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
)

func TestDivideTask_Happy(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name                  string
		input                 interface{}
		divisor               string
		precision             string
		expected              *decimal.Decimal
		expectedErrorCause    error
		expectedErrorContains string
	}{
		{"string", "12345.67", "100", "", mustDecimal(t, "123.4567"), nil, ""},
		{"string, negative", "12345.67", "-5", "", mustDecimal(t, "-2469.134"), nil, ""},
		{"string, large value", "12345.67", "1000000000000000000", "", mustDecimal(t, "0.0000000000000123"), nil, ""},

		{"int", int(200), "16", "", mustDecimal(t, "12.5"), nil, ""},
		{"int, negative", int(200), "-5", "", mustDecimal(t, "-40"), nil, ""},
		{"int, large value", int(200), "1000000000000000000", "", mustDecimal(t, "0.0000000000000002"), nil, ""},

		{"int8", int8(20), "16", "", mustDecimal(t, "1.25"), nil, ""},
		{"int8, negative", int8(20), "-5", "", mustDecimal(t, "-4"), nil, ""},
		{"int8, large value", int8(20), "10000000000000000", "", mustDecimal(t, "0.000000000000002"), nil, ""},

		{"int16", int16(200), "16", "", mustDecimal(t, "12.5"), nil, ""},
		{"int16, negative", int16(200), "-5", "", mustDecimal(t, "-40"), nil, ""},
		{"int16, large value", int16(200), "1000000000000000000", "", mustDecimal(t, "0.0000000000000002"), nil, ""},

		{"int32", int32(200), "16", "", mustDecimal(t, "12.5"), nil, ""},
		{"int32, negative", int32(200), "-5", "", mustDecimal(t, "-40"), nil, ""},
		{"int32, large value", int32(200), "1000000000000000000", "", mustDecimal(t, "0.0000000000000002"), nil, ""},

		{"int64", int64(200), "16", "", mustDecimal(t, "12.5"), nil, ""},
		{"int64, negative", int64(200), "-5", "", mustDecimal(t, "-40"), nil, ""},
		{"int64, large value", int64(200), "1000000000000000000", "", mustDecimal(t, "0.0000000000000002"), nil, ""},

		{"uint", uint(200), "16", "", mustDecimal(t, "12.5"), nil, ""},
		{"uint, negative", uint(200), "-5", "", mustDecimal(t, "-40"), nil, ""},
		{"uint, large value", uint(200), "1000000000000000000", "", mustDecimal(t, "0.0000000000000002"), nil, ""},

		{"uint8", uint8(200), "16", "", mustDecimal(t, "12.5"), nil, ""},
		{"uint8, negative", uint8(200), "-5", "", mustDecimal(t, "-40"), nil, ""},
		{"uint8, large value", uint8(200), "1000000000000000000", "", mustDecimal(t, "0.0000000000000002"), nil, ""},

		{"uint16", uint16(200), "16", "", mustDecimal(t, "12.5"), nil, ""},
		{"uint16, negative", uint16(200), "-5", "", mustDecimal(t, "-40"), nil, ""},
		{"uint16, large value", uint16(200), "1000000000000000000", "", mustDecimal(t, "0.0000000000000002"), nil, ""},

		{"uint32", uint32(200), "16", "", mustDecimal(t, "12.5"), nil, ""},
		{"uint32, negative", uint32(200), "-5", "", mustDecimal(t, "-40"), nil, ""},
		{"uint32, large value", uint32(200), "1000000000000000000", "", mustDecimal(t, "0.0000000000000002"), nil, ""},

		{"uint64", uint64(200), "16", "", mustDecimal(t, "12.5"), nil, ""},
		{"uint64, negative", uint64(200), "-5", "", mustDecimal(t, "-40"), nil, ""},
		{"uint64, large value", uint64(200), "1000000000000000000", "", mustDecimal(t, "0.0000000000000002"), nil, ""},

		{"float32", float32(12345.67), "1000", "", mustDecimal(t, "12.34567"), nil, ""},
		{"float32, negative", float32(12345.67), "-5", "", mustDecimal(t, "-2469.134"), nil, ""},
		{"float32, large value", float32(12345.67), "1000000000000000000", "", mustDecimal(t, "0.0000000000000123"), nil, ""},

		{"float64", float64(12345.67), "1000", "", mustDecimal(t, "12.34567"), nil, ""},
		{"float64, negative", float64(12345.67), "-5", "", mustDecimal(t, "-2469.134"), nil, ""},
		{"float64, large value", float64(12345.67), "1000000000000000000", "", mustDecimal(t, "0.0000000000000123"), nil, ""},

		{"precision", float64(12345.67), "1000", "2", mustDecimal(t, "12.35"), nil, ""},
		{"precision (> 16)", float64(200), "6", "18", mustDecimal(t, "33.333333333333333333"), nil, ""},
		{"precision (negative)", float64(12345.67), "1000", "-1", mustDecimal(t, "10"), nil, ""},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			vars := pipeline.NewVarsFrom(nil)
			task := pipeline.DivideTask{
				BaseTask:  pipeline.NewBaseTask(0, "task", nil, nil, 0),
				Divisor:   test.divisor,
				Precision: test.precision,
			}
			result, runInfo := task.Run(context.Background(), logger.TestLogger(t), vars, []pipeline.Result{{Value: test.input}})
			assert.False(t, runInfo.IsPending)
			assert.False(t, runInfo.IsRetryable)
			require.NoError(t, result.Error)
			require.Equal(t, test.expected.String(), result.Value.(decimal.Decimal).String())
		})
	}

	for _, test := range tests {
		test := test
		t.Run(test.name+" (with pipeline.Vars)", func(t *testing.T) {
			vars := pipeline.NewVarsFrom(map[string]interface{}{
				"foo":    map[string]interface{}{"bar": test.input},
				"chain":  map[string]interface{}{"link": test.divisor},
				"sergey": map[string]interface{}{"steve": test.precision},
			})
			task := pipeline.DivideTask{
				BaseTask:  pipeline.NewBaseTask(0, "task", nil, nil, 0),
				Input:     "$(foo.bar)",
				Divisor:   "$(chain.link)",
				Precision: "$(sergey.steve)",
			}
			result, runInfo := task.Run(context.Background(), logger.TestLogger(t), vars, []pipeline.Result{})
			assert.False(t, runInfo.IsPending)
			assert.False(t, runInfo.IsRetryable)
			require.NoError(t, result.Error)
			require.Equal(t, test.expected.String(), result.Value.(decimal.Decimal).String())
		})
	}
}

func TestDivideTask_Unhappy(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name              string
		divisor           string
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
		{"divisor as missing var", "$(foo)", "", []pipeline.Result{{Value: "123"}}, pipeline.NewVarsFrom(nil), pipeline.ErrKeypathNotFound, "divisor"},
		{"errored inputs", "1000", "", []pipeline.Result{{Error: errors.New("uh oh")}}, pipeline.NewVarsFrom(nil), pipeline.ErrTooManyErrors, "task inputs"},
	}

	for _, tt := range tests {
		test := tt
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			task := pipeline.DivideTask{
				BaseTask: pipeline.NewBaseTask(0, "task", nil, nil, 0),
				Input:    test.input,
				Divisor:  test.divisor,
			}
			result, runInfo := task.Run(context.Background(), logger.TestLogger(t), test.vars, test.inputs)
			assert.False(t, runInfo.IsPending)
			assert.False(t, runInfo.IsRetryable)
			require.Equal(t, test.wantErrorCause, errors.Cause(result.Error))
			if test.wantErrorContains != "" {
				require.Contains(t, result.Error.Error(), test.wantErrorContains)
			}
		})
	}
}
