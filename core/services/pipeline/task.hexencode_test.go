package pipeline_test

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
)

func TestHexEncodeTask(t *testing.T) {
	t.Parallel()
	bigTwo, bigThree := big.NewInt(2), big.NewInt(3)

	tests := []struct {
		name   string
		input  interface{}
		result string
		error  string
	}{

		// success integers
		{"zero", 0, "0x0", ""},
		{"small int", 1, "0x1", ""},
		{"two-byte integer", 256, "0x100", ""},
		{"uint8", uint8(10), "0xa", ""},
		{"small int64", int64(456), "0x1c8", ""},
		{"large int64", int64(999000000000), "0xe8990a4600", ""},
		{"bigint 1", bigTwo.Exp(bigTwo, big.NewInt(100), nil), "0x10000000000000000000000000", ""},
		{"bigint 2", bigThree.Exp(bigThree, big.NewInt(100), nil), "0x5a4653ca673768565b41f775d6947d55cf3813d1", ""},
		{"decimal type but integer value", 1.0, "0x1", ""},
		{"decimal type but integer value zero", 0.0, "0x0", ""},
		{"decimal.Decimal type but integer value", mustDecimal(t, "256"), "0x100", ""},

		// success strings/bytes
		{"string ascii bytes", "xyz", "0x78797a", ""},
		{"string with whitespace", "1 x *", "0x312078202a", ""},
		{"string shouldn't convert to int", "456", "0x343536", ""},
		{"don't detect hex in string", "0xff", "0x30786666", ""},
		// NOTE: for byte arrays, output is padded to full bytes (i.e. a potential leading zero)
		{"bytes remain bytes", []byte{0xa, 0x0, 0xff, 0x1}, "0x0a00ff01", ""},

		// success empty results
		{"empty string", "", "", ""},
		{"empty byte array", []byte{}, "", ""},

		// failure
		{"negative int", -1, "", "negative integer"},
		{"negative float", -1.0, "", "negative integer"},
		{"negative int64", int64(-10), "", "negative integer"},
		{"negative bigint", big.NewInt(-100), "", "negative integer"},
		{"input of type bool", true, "", "bad input for task"},
		{"input of type decimal", 1.44, "", "decimal input"},
		{"input of type decimal and negative", -0.44, "", "decimal input"},
		{"input of decimal.Decimal type but not integer", mustDecimal(t, "3.14"), "", "decimal input"},
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
			t.Run("without vars", func(t *testing.T) {
				vars := pipeline.NewVarsFrom(nil)
				task := pipeline.HexEncodeTask{BaseTask: pipeline.NewBaseTask(0, "task", nil, nil, 0)}
				assertOK(task.Run(testutils.Context(t), logger.TestLogger(t), vars, []pipeline.Result{{Value: test.input}}))
			})
			t.Run("with vars", func(t *testing.T) {
				vars := pipeline.NewVarsFrom(map[string]interface{}{
					"foo": map[string]interface{}{"bar": test.input},
				})
				task := pipeline.HexEncodeTask{
					BaseTask: pipeline.NewBaseTask(0, "task", nil, nil, 0),
					Input:    "$(foo.bar)",
				}
				assertOK(task.Run(testutils.Context(t), logger.TestLogger(t), vars, []pipeline.Result{}))
			})
		})
	}
}

func TestHexEncodeTaskInputParamLiteral(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		input  interface{}
		result string
	}{
		// Only strings can be passed via input param literals (other types will get converted to strings anyway)
		{"string ascii bytes", "xyz", "0x78797a"},
		{"string with whitespace", "1 x *", "0x312078202a"},
		{"string shouldn't convert to int", "456", "0x343536"},
		{"don't detect hex in string", "0xff", "0x30786666"},
		{"int gets converted to string", 256, "0x323536"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			vars := pipeline.NewVarsFrom(nil)
			task := pipeline.HexEncodeTask{
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
