package pipeline_test

import (
	"context"
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/services/pipeline"
)

func marshallToBigInt(v interface{}) *big.Int {
	res := new(big.Int)
	bytes := v.([]byte)
	res.SetBytes(bytes)
	return res
}

func TestUInt256Task_Succeed(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string

		input interface{}
		want  big.Int
	}{
		{"convert 1", 0, *big.NewInt(0)},
		{"convert 1", 1, *big.NewInt(1)},
		{"convert 1000000000000000000", 1000000000000000000, *big.NewInt(1000000000000000000)},

		{"convert int(0)", int(0), *big.NewInt(0)},
		{"convert int(1)", int(1), *big.NewInt(1)},
		{"convert int(1000000000000000000)", int(1000000000000000000), *big.NewInt(1000000000000000000)},

		{"convert int8(0)", int8(0), *big.NewInt(0)},
		{"convert int8(1)", int8(1), *big.NewInt(1)},
		{"convert int8(100)", int8(100), *big.NewInt(100)},

		{"convert int16(0)", int16(0), *big.NewInt(0)},
		{"convert int16(1)", int16(1), *big.NewInt(1)},
		{"convert int16(10000)", int16(10000), *big.NewInt(10000)},

		{"convert int32(0)", int32(0), *big.NewInt(0)},
		{"convert int32(1)", int32(1), *big.NewInt(1)},
		{"convert int32(1000000000)", int32(1000000000), *big.NewInt(1000000000)},

		{"convert int64(0)", int64(0), *big.NewInt(0)},
		{"convert int64(1)", int64(1), *big.NewInt(1)},
		{"convert int64(1000000000000000000)", int64(1000000000000000000), *big.NewInt(1000000000000000000)},

		{"convert float32(0)", float32(0), *big.NewInt(0)},
		{"convert float32(1)", float32(1), *big.NewInt(1)},
		{"convert float32(1000000000000000000)", float32(1000000000000000000), *big.NewInt(1000000000000000000)},

		{"convert float64(0)", float64(0), *big.NewInt(0)},
		{"convert float64(1)", float64(1), *big.NewInt(1)},
		{"convert float64(1000000000000000000)", float64(1000000000000000000), *big.NewInt(1000000000000000000)},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			task := pipeline.UInt256Task{}
			result := task.Run(context.Background(), pipeline.JSONSerializable{}, []pipeline.Result{{Value: test.input}})
			require.NoError(t, result.Error)
			require.Equal(t, test.want.Cmp(marshallToBigInt(result.Value)), 0)
		})
	}
}

func TestUInt256Task_Fail(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input interface{}
	}{
		{"map", map[string]interface{}{"chain": "link"}},
		{"slice", []interface{}{"chain", "link"}},
	}

	for _, tt := range tests {
		test := tt
		t.Run(test.name, func(t *testing.T) {
			task := pipeline.UInt256Task{}
			result := task.Run(context.Background(), pipeline.JSONSerializable{}, []pipeline.Result{{Value: test.input}})
			require.Error(t, result.Error)
		})
	}
}
