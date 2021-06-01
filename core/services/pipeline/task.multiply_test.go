package pipeline_test

import (
	"context"
	"testing"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/services/pipeline"
)

func mustDecimal(t *testing.T, arg string) *decimal.Decimal {
	ret, err := decimal.NewFromString(arg)
	require.NoError(t, err)
	return &ret
}

func TestMultiplyTask_Happy(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input interface{}
		times string
		want  decimal.Decimal
	}{
		{"string, by 100", "1.23", "100", *mustDecimal(t, "123")},
		{"string, negative", "1.23", "-5", *mustDecimal(t, "-6.15")},
		{"string, no times parameter", "1.23", "1", *mustDecimal(t, "1.23")},
		{"string, zero", "1.23", "0", *mustDecimal(t, "0")},
		{"string, large value", "1.23", "1000000000000000000", *mustDecimal(t, "1230000000000000000")},

		{"int, by 100", int(2), "100", *mustDecimal(t, "200")},
		{"int, negative", int(2), "-5", *mustDecimal(t, "-10")},
		{"int, no times parameter", int(2), "1", *mustDecimal(t, "2")},
		{"int, zero", int(2), "0", *mustDecimal(t, "0")},
		{"int, large value", int(2), "1000000000000000000", *mustDecimal(t, "2000000000000000000")},

		{"int8, by 100", int8(2), "100", *mustDecimal(t, "200")},
		{"int8, negative", int8(2), "-5", *mustDecimal(t, "-10")},
		{"int8, no times parameter", int8(2), "1", *mustDecimal(t, "2")},
		{"int8, zero", int8(2), "0", *mustDecimal(t, "0")},
		{"int8, large value", int8(2), "1000000000000000000", *mustDecimal(t, "2000000000000000000")},

		{"int16, by 100", int16(2), "100", *mustDecimal(t, "200")},
		{"int16, negative", int16(2), "-5", *mustDecimal(t, "-10")},
		{"int16, no times parameter", int16(2), "1", *mustDecimal(t, "2")},
		{"int16, zero", int16(2), "0", *mustDecimal(t, "0")},
		{"int16, large value", int16(2), "1000000000000000000", *mustDecimal(t, "2000000000000000000")},

		{"int32, by 100", int32(2), "100", *mustDecimal(t, "200")},
		{"int32, negative", int32(2), "-5", *mustDecimal(t, "-10")},
		{"int32, no times parameter", int32(2), "1", *mustDecimal(t, "2")},
		{"int32, zero", int32(2), "0", *mustDecimal(t, "0")},
		{"int32, large value", int32(2), "1000000000000000000", *mustDecimal(t, "2000000000000000000")},

		{"int64, by 100", int64(2), "100", *mustDecimal(t, "200")},
		{"int64, negative", int64(2), "-5", *mustDecimal(t, "-10")},
		{"int64, no times parameter", int64(2), "1", *mustDecimal(t, "2")},
		{"int64, zero", int64(2), "0", *mustDecimal(t, "0")},
		{"int64, large value", int64(2), "1000000000000000000", *mustDecimal(t, "2000000000000000000")},

		{"uint, by 100", uint(2), "100", *mustDecimal(t, "200")},
		{"uint, negative", uint(2), "-5", *mustDecimal(t, "-10")},
		{"uint, no times parameter", uint(2), "1", *mustDecimal(t, "2")},
		{"uint, zero", uint(2), "0", *mustDecimal(t, "0")},
		{"uint, large value", uint(2), "1000000000000000000", *mustDecimal(t, "2000000000000000000")},

		{"uint8, by 100", uint8(2), "100", *mustDecimal(t, "200")},
		{"uint8, negative", uint8(2), "-5", *mustDecimal(t, "-10")},
		{"uint8, no times parameter", uint8(2), "1", *mustDecimal(t, "2")},
		{"uint8, zero", uint8(2), "0", *mustDecimal(t, "0")},
		{"uint8, large value", uint8(2), "1000000000000000000", *mustDecimal(t, "2000000000000000000")},

		{"uint16, by 100", uint16(2), "100", *mustDecimal(t, "200")},
		{"uint16, negative", uint16(2), "-5", *mustDecimal(t, "-10")},
		{"uint16, no times parameter", uint16(2), "1", *mustDecimal(t, "2")},
		{"uint16, zero", uint16(2), "0", *mustDecimal(t, "0")},
		{"uint16, large value", uint16(2), "1000000000000000000", *mustDecimal(t, "2000000000000000000")},

		{"uint32, by 100", uint32(2), "100", *mustDecimal(t, "200")},
		{"uint32, negative", uint32(2), "-5", *mustDecimal(t, "-10")},
		{"uint32, no times parameter", uint32(2), "1", *mustDecimal(t, "2")},
		{"uint32, zero", uint32(2), "0", *mustDecimal(t, "0")},
		{"uint32, large value", uint32(2), "1000000000000000000", *mustDecimal(t, "2000000000000000000")},

		{"uint64, by 100", uint64(2), "100", *mustDecimal(t, "200")},
		{"uint64, negative", uint64(2), "-5", *mustDecimal(t, "-10")},
		{"uint64, no times parameter", uint64(2), "1", *mustDecimal(t, "2")},
		{"uint64, zero", uint64(2), "0", *mustDecimal(t, "0")},
		{"uint64, large value", uint64(2), "1000000000000000000", *mustDecimal(t, "2000000000000000000")},

		{"float32, by 100", float32(1.23), "10", *mustDecimal(t, "12.3")},
		{"float32, negative", float32(1.23), "-5", *mustDecimal(t, "-6.15")},
		{"float32, no times parameter", float32(1.23), "1", *mustDecimal(t, "1.23")},
		{"float32, zero", float32(1.23), "0", *mustDecimal(t, "0")},
		{"float32, large value", float32(1.23), "1000000000000000000", *mustDecimal(t, "1230000000000000000")},

		{"float64, by 100", float64(1.23), "10", *mustDecimal(t, "12.3")},
		{"float64, negative", float64(1.23), "-5", *mustDecimal(t, "-6.15")},
		{"float64, no times parameter", float64(1.23), "1", *mustDecimal(t, "1.23")},
		{"float64, zero", float64(1.23), "0", *mustDecimal(t, "0")},
		{"float64, large value", float64(1.23), "1000000000000000000", *mustDecimal(t, "1230000000000000000")},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			vars := pipeline.NewVarsFrom(nil)
			task := pipeline.MultiplyTask{BaseTask: pipeline.NewBaseTask(0, "task", nil, 0), Times: test.times}
			result := task.Run(context.Background(), vars, pipeline.JSONSerializable{}, []pipeline.Result{{Value: test.input}})
			require.NoError(t, result.Error)
			require.Equal(t, test.want.String(), result.Value.(decimal.Decimal).String())
		})
	}

	for _, test := range tests {
		test := test
		t.Run(test.name+" (with pipeline.Vars)", func(t *testing.T) {
			vars := pipeline.NewVarsFrom(map[string]interface{}{
				"foo":   map[string]interface{}{"bar": test.input},
				"chain": map[string]interface{}{"link": test.times},
			})
			task := pipeline.MultiplyTask{
				BaseTask: pipeline.NewBaseTask(0, "task", nil, 0),
				Input:    "$(foo.bar)",
				Times:    "$(chain.link)",
			}
			result := task.Run(context.Background(), vars, pipeline.JSONSerializable{}, []pipeline.Result{})
			require.NoError(t, result.Error)
			require.Equal(t, test.want.String(), result.Value.(decimal.Decimal).String())
		})
	}
}

func TestMultiplyTask_Unhappy(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name              string
		times             string
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
		{"times as missing var", "$(foo)", "", []pipeline.Result{{Value: "123"}}, pipeline.NewVarsFrom(nil), pipeline.ErrKeypathNotFound, "times"},
	}

	for _, tt := range tests {
		test := tt
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			task := pipeline.MultiplyTask{
				BaseTask: pipeline.NewBaseTask(0, "task", nil, 0),
				Input:    test.input,
				Times:    test.times,
			}
			result := task.Run(context.Background(), test.vars, pipeline.JSONSerializable{}, test.inputs)
			require.Equal(t, test.wantErrorCause, errors.Cause(result.Error))
			if test.wantErrorContains != "" {
				require.Contains(t, result.Error.Error(), test.wantErrorContains)
			}
		})
	}
}
