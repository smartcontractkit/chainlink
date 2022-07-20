package pipeline_test

import (
	"math/big"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/services/pipeline"
)

func TestObjectParam_UnmarshalPipelineParamValid(t *testing.T) {
	t.Parallel()

	decimalValue := decimal.New(173, -1)

	tests := []struct {
		name   string
		input  interface{}
		output string
	}{
		{"identity", pipeline.ObjectParam{Type: pipeline.BoolType, BoolValue: true}, "true"},

		{"nil", nil, "null"},

		{"bool", true, "true"},
		{"bool false", false, "false"},

		{"uint8", uint8(17), `"17"`},
		{"uint16", uint16(17), `"17"`},
		{"uint32", uint32(17), `"17"`},
		{"uint64", uint64(17), `"17"`},
		{"uint", 17, `"17"`},

		{"int8", int8(17), `"17"`},
		{"int16", int16(17), `"17"`},
		{"int32", int32(17), `"17"`},
		{"int64", int64(17), `"17"`},
		{"integer", 17, `"17"`},

		{"negative integer", -19, `"-19"`},
		{"float32", float32(17.3), `"17.3"`},
		{"float", 17.3, `"17.3"`},
		{"negative float", -17.3, `"-17.3"`},

		{"bigintp", big.NewInt(-17), `"-17"`},
		{"bigint", *big.NewInt(29), `"29"`},

		{"decimalp", &decimalValue, `"17.3"`},
		{"decimal", decimalValue, `"17.3"`},

		{"string", "hello world", `"hello world"`},

		{"array", []int{17, 19}, "[17,19]"},
		{"empty array", []interface{}{}, "[]"},
		{"interface array", []interface{}{17, 19}, "[17,19]"},
		{"string array", []string{"hello", "world"}, `["hello","world"]`},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var value pipeline.ObjectParam
			err := value.UnmarshalPipelineParam(test.input)
			require.NoError(t, err)
			marshalledValue, err := value.Marshal()
			require.NoError(t, err)
			assert.Equal(t, test.output, marshalledValue)
		})
	}
}

func TestObjectParam_Marshal(t *testing.T) {
	tests := []struct {
		name   string
		input  *pipeline.ObjectParam
		output string
	}{
		{"nil", mustNewObjectParam(t, nil), "null"},
		{"bool", mustNewObjectParam(t, true), "true"},
		{"integer", mustNewObjectParam(t, 17), `"17"`},
		{"string", mustNewObjectParam(t, "hello world"), `"hello world"`},
		{"array", mustNewObjectParam(t, []int{17, 19}), "[17,19]"},
		{"map", mustNewObjectParam(t, map[string]interface{}{"key": 19}), `{"key":19}`},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			marshalledValue, err := test.input.Marshal()
			require.NoError(t, err)
			assert.Equal(t, test.output, marshalledValue)
		})
	}
}
