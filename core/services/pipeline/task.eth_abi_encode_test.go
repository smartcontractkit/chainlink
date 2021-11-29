package pipeline_test

import (
	"context"
	"math"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
)

func TestETHABIEncodeTask(t *testing.T) {
	var bytes32 [32]byte
	copy(bytes32[:], []byte("chainlink chainlink chainlink"))

	tests := []struct {
		name                  string
		abi                   string
		data                  string
		vars                  pipeline.Vars
		inputs                []pipeline.Result
		expected              string
		expectedErrorCause    error
		expectedErrorContains string
	}{
		{
			"unusual characters in method name / uint256, bool, int256, string",
			"foo_Bar__3928 ( uint256 u, bool b, int256 i, string s )",
			`{ "u": $(foo), "b": $(bar), "i": $(baz), "s": $(quux) }`,
			pipeline.NewVarsFrom(map[string]interface{}{
				"foo":  big.NewInt(123),
				"bar":  true,
				"baz":  big.NewInt(-321),
				"quux": "foo bar baz",
			}),
			nil,
			"0xae506917000000000000000000000000000000000000000000000000000000000000007b0000000000000000000000000000000000000000000000000000000000000001fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffebf0000000000000000000000000000000000000000000000000000000000000080000000000000000000000000000000000000000000000000000000000000000b666f6f206261722062617a000000000000000000000000000000000000000000",
			nil,
			"",
		},
		{
			"bytes32, bytes, address",
			"asdf(bytes32 b, bytes bs, address a)",
			`{ "b": $(foo), "bs": $(bar), "a": $(baz) }`,
			pipeline.NewVarsFrom(map[string]interface{}{
				"foo": bytes32,
				"bar": []byte("stevetoshi sergeymoto"),
				"baz": common.HexToAddress("0xdeadbeefdeadbeefdeadbeefdeadbeefdeadbeef"),
			}),
			nil,
			"0x4f5e7a89636861696e6c696e6b20636861696e6c696e6b20636861696e6c696e6b0000000000000000000000000000000000000000000000000000000000000000000060000000000000000000000000deadbeefdeadbeefdeadbeefdeadbeefdeadbeef00000000000000000000000000000000000000000000000000000000000000157374657665746f736869207365726765796d6f746f0000000000000000000000",
			nil,
			"",
		},
		{
			"address[] calldata, uint80, uint32[2]",
			"chainLink(address[] calldata a, uint80 x, uint32[2] s)",
			`{ "a": $(foo), "x": $(bar), "s": $(baz) }`,
			pipeline.NewVarsFrom(map[string]interface{}{
				"foo": []common.Address{
					common.HexToAddress("0x6c91b062a774cbe8b9bf52f224c37badf98fc40b"),
					common.HexToAddress("0xc4f27ead9083c756cc2c02aaa39b223fe8d0a0e5"),
					common.HexToAddress("0x749e4598819b2b0e915a02120696c7b8fe16c09c"),
				},
				"bar": big.NewInt(8293),
				"baz": []*big.Int{big.NewInt(192), big.NewInt(4182)},
			}),
			nil,
			"0xa3a122020000000000000000000000000000000000000000000000000000000000000080000000000000000000000000000000000000000000000000000000000000206500000000000000000000000000000000000000000000000000000000000000c0000000000000000000000000000000000000000000000000000000000000105600000000000000000000000000000000000000000000000000000000000000030000000000000000000000006c91b062a774cbe8b9bf52f224c37badf98fc40b000000000000000000000000c4f27ead9083c756cc2c02aaa39b223fe8d0a0e5000000000000000000000000749e4598819b2b0e915a02120696c7b8fe16c09c",
			nil,
			"",
		},
		{
			"bool[2][] calldata, uint96[2][] calldata",
			"arrayOfArrays(bool[2][] calldata bools, uint96[2][] calldata uints)",
			`{ "bools": $(foo), "uints": $(bar) }`,
			pipeline.NewVarsFrom(map[string]interface{}{
				"foo": [][]bool{{true, false}, {false, true}, {false, false}, {true, true}},
				"bar": [][]*big.Int{{big.NewInt(123), big.NewInt(456)}, {big.NewInt(22), big.NewInt(19842)}},
			}),
			nil,
			"0xb04bee77000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000001600000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000007b00000000000000000000000000000000000000000000000000000000000001c800000000000000000000000000000000000000000000000000000000000000160000000000000000000000000000000000000000000000000000000000004d82",
			nil,
			"",
		},
		{
			"no args",
			"noArgs()",
			``,
			pipeline.NewVarsFrom(nil),
			nil,
			"0x83c962bb",
			nil,
			"",
		},
		{
			"number too large for uint32",
			"willFail(uint32 s)",
			`{ "s": $(foo) }`,
			pipeline.NewVarsFrom(map[string]interface{}{
				"foo": big.NewInt(math.MaxInt64),
			}),
			nil,
			"",
			pipeline.ErrBadInput,
			"overflow",
		},
		{
			"string too large for address",
			"willFail(address a)",
			`{ "a": $(foo) }`,
			pipeline.NewVarsFrom(map[string]interface{}{
				"foo": "0xdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeef",
			}),
			nil,
			"",
			pipeline.ErrBadInput,
			"incorrect length",
		},
		{
			"too many array elements",
			"willFail(uint32[2] a)",
			`{ "a": $(foo) }`,
			pipeline.NewVarsFrom(map[string]interface{}{
				"foo": []interface{}{123, 456, 789},
			}),
			nil,
			"",
			pipeline.ErrBadInput,
			"incorrect length",
		},
		{
			"too many array elements (nested)",
			"willFail(uint32[2][] a)",
			`{ "a": $(foo) }`,
			pipeline.NewVarsFrom(map[string]interface{}{
				"foo": []interface{}{
					[]interface{}{123, 456, 789},
				},
			}),
			nil,
			"",
			pipeline.ErrBadInput,
			"incorrect length",
		},
		{
			"no argument names",
			"willFail(address, uint256[])",
			``,
			pipeline.NewVarsFrom(nil),
			nil,
			"",
			pipeline.ErrBadInput,
			"missing argument name",
		},
		{
			"no argument names (calldata)",
			"willFail(uint256[] calldata)",
			``,
			pipeline.NewVarsFrom(nil),
			nil,
			"",
			pipeline.ErrBadInput,
			"missing argument name",
		},
		{
			"errored task inputs",
			"asdf(bytes32 b, bytes bs, address a)",
			`{ "b": $(foo), "bs": $(bar), "a": $(baz) }`,
			pipeline.NewVarsFrom(map[string]interface{}{
				"foo": bytes32,
				"bar": []byte("stevetoshi sergeymoto"),
				"baz": common.HexToAddress("0xdeadbeefdeadbeefdeadbeefdeadbeefdeadbeef"),
			}),
			[]pipeline.Result{{Error: errors.New("uh oh")}},
			"",
			pipeline.ErrTooManyErrors,
			"task inputs",
		},
		{
			"hex string to fixed size byte array (note used by fulfillOracleRequest(..., bytes32 data))",
			"asdf(bytes32 b)",
			`{ "b": $(foo)}`,
			pipeline.NewVarsFrom(map[string]interface{}{
				"foo": "0x0000000000000000000000000000000000000000000000000000000000000001",
			}),
			nil,
			"0x628507ac0000000000000000000000000000000000000000000000000000000000000001",
			nil,
			"",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			task := pipeline.ETHABIEncodeTask{
				BaseTask: pipeline.NewBaseTask(0, "encode", nil, nil, 0),
				ABI:      test.abi,
				Data:     test.data,
			}

			result, runInfo := task.Run(context.Background(), logger.TestLogger(t), test.vars, test.inputs)
			assert.False(t, runInfo.IsPending)
			assert.False(t, runInfo.IsRetryable)

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
