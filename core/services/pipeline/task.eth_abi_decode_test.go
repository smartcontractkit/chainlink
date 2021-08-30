package pipeline_test

import (
	"context"
	"math/big"
	"testing"

	"github.com/pkg/errors"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/core/utils"
)

func TestETHABIDecodeTask(t *testing.T) {
	tests := []struct {
		name                  string
		abi                   string
		data                  string
		vars                  pipeline.Vars
		inputs                []pipeline.Result
		expected              map[string]interface{}
		expectedErrorCause    error
		expectedErrorContains string
	}{
		{
			"uint256, bool, int256, string",
			"uint256 u, bool b, int256 i, string s",
			"$(foo)",
			pipeline.NewVarsFrom(map[string]interface{}{
				"foo": "0x000000000000000000000000000000000000000000000000000000000000007b0000000000000000000000000000000000000000000000000000000000000001fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffebf0000000000000000000000000000000000000000000000000000000000000080000000000000000000000000000000000000000000000000000000000000000b666f6f206261722062617a000000000000000000000000000000000000000000",
			}),
			nil,
			map[string]interface{}{
				"u": big.NewInt(123),
				"b": true,
				"i": big.NewInt(-321),
				"s": "foo bar baz",
			},
			nil,
			"",
		},
		{
			"weird spaces / address, uint80[3][], bytes, bytes32",
			"address  a , uint80[3][] u , bytes b, bytes32 b32  ",
			"$(foo)",
			pipeline.NewVarsFrom(map[string]interface{}{
				"foo": "0x000000000000000000000000deadbeefdeadbeefdeadbeefdeadbeefdeadbeef000000000000000000000000000000000000000000000000000000000000008000000000000000000000000000000000000000000000000000000000000001607374657665746f7368692073657267616d6f746f0000000000000000000000000000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000005c000000000000000000000000000000000000000000000000000000000000003d000000000000000000000000000000000000000000000000000000000000001e000000000000000000000000000000000000000000000000000000000000002100000000000000000000000000000000000000000000000000000000000000420000000000000000000000000000000000000000000000000000000000000063000000000000000000000000000000000000000000000000000000000000000c666f6f206261722062617a0a0000000000000000000000000000000000000000",
			}),
			nil,
			map[string]interface{}{
				"a": common.HexToAddress("0xdeadbeefdeadbeefdeadbeefdeadbeefdeadbeef"),
				"u": [][3]*big.Int{
					{big.NewInt(92), big.NewInt(61), big.NewInt(30)},
					{big.NewInt(33), big.NewInt(66), big.NewInt(99)},
				},
				"b":   hexutil.MustDecode("0x666f6f206261722062617a0a"),
				"b32": utils.Bytes32FromString("stevetoshi sergamoto"),
			},
			nil,
			"",
		},
		{
			"no attribute names",
			"address, bytes32",
			"$(foo)",
			pipeline.NewVarsFrom(map[string]interface{}{
				"foo": "0x000000000000000000000000deadbeefdeadbeefdeadbeefdeadbeefdeadbeef000000000000000000000000000000000000000000000000000000000000008000000000000000000000000000000000000000000000000000000000000001607374657665746f7368692073657267616d6f746f0000000000000000000000000000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000005c000000000000000000000000000000000000000000000000000000000000003d000000000000000000000000000000000000000000000000000000000000001e000000000000000000000000000000000000000000000000000000000000002100000000000000000000000000000000000000000000000000000000000000420000000000000000000000000000000000000000000000000000000000000063000000000000000000000000000000000000000000000000000000000000000c666f6f206261722062617a0a0000000000000000000000000000000000000000",
			}),
			nil,
			nil,
			pipeline.ErrBadInput,
			"",
		},
		{
			"errored task inputs",
			"uint256 u, bool b, int256 i, string s",
			"$(foo)",
			pipeline.NewVarsFrom(map[string]interface{}{
				"foo": "0x000000000000000000000000000000000000000000000000000000000000007b0000000000000000000000000000000000000000000000000000000000000001fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffebf0000000000000000000000000000000000000000000000000000000000000080000000000000000000000000000000000000000000000000000000000000000b666f6f206261722062617a000000000000000000000000000000000000000000",
			}),
			[]pipeline.Result{{Error: errors.New("uh oh")}},
			nil,
			pipeline.ErrTooManyErrors,
			"task inputs",
		},
	}

	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {
			task := pipeline.ETHABIDecodeTask{
				BaseTask: pipeline.NewBaseTask(0, "decode", nil, nil, 0),
				ABI:      test.abi,
				Data:     test.data,
			}

			result := task.Run(context.Background(), test.vars, test.inputs)

			if test.expectedErrorCause != nil {
				require.Equal(t, test.expectedErrorCause, errors.Cause(result.Error))
				require.Nil(t, result.Value)
				if test.expectedErrorContains != "" {
					require.Contains(t, result.Error.Error(), test.expectedErrorContains)
				}
			} else {
				require.NoError(t, result.Error)
				require.Equal(t, test.expected, result.Value)
			}
		})
	}
}
