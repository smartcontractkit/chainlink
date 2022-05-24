package pipeline_test

import (
	"context"
	"fmt"
	"math"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/core/utils"
)

func TestETHABIEncodeTask(t *testing.T) {
	var bytes32 [32]byte
	copy(bytes32[:], []byte("chainlink chainlink chainlink"))

	bytes32hex := utils.StringToHex(string(bytes32[:]))

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
			"bytes32 (hex), bytes, address",
			"asdf(bytes32 b, bytes bs, address a)",
			`{ "b": $(foo), "bs": $(bar), "a": $(baz) }`,
			pipeline.NewVarsFrom(map[string]interface{}{
				"foo": bytes32hex,
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

func TestETHABIEncode_EncodeIntegers(t *testing.T) {
	testCases := []struct {
		name                  string
		abi                   string
		data                  string
		vars                  pipeline.Vars
		inputs                []pipeline.Result
		expected              string
		expectedErrorCause    error
		expectedErrorContains string
	}{
		// no overflow cases
		// 8 bit ints.
		{
			"encode 1 to int8",
			"asdf(int8 i)",
			`{ "i": $(foo) }`,
			pipeline.NewVarsFrom(map[string]interface{}{
				"foo": int8(1),
			}),
			nil,
			"0xa8d7f3cd0000000000000000000000000000000000000000000000000000000000000001",
			nil,
			"",
		},
		{
			"encode 1 to uint8",
			"asdf(uint8 i)",
			`{ "i": $(foo) }`,
			pipeline.NewVarsFrom(map[string]interface{}{
				"foo": uint8(1),
			}),
			nil,
			"0x6b377be20000000000000000000000000000000000000000000000000000000000000001",
			nil,
			"",
		},
		// 16 bit ints.
		{
			"encode 1 to int16",
			"asdf(int16 i)",
			`{ "i": $(foo) }`,
			pipeline.NewVarsFrom(map[string]interface{}{
				"foo": int16(1),
			}),
			nil,
			"0xabd195460000000000000000000000000000000000000000000000000000000000000001",
			nil,
			"",
		},
		{
			"encode 1 to uint16",
			"asdf(uint16 i)",
			`{ "i": $(foo) }`,
			pipeline.NewVarsFrom(map[string]interface{}{
				"foo": uint16(1),
			}),
			nil,
			"0x8f3294d20000000000000000000000000000000000000000000000000000000000000001",
			nil,
			"",
		},
		// 24 bit ints.
		{
			"encode 1 to int24",
			"asdf(int24 i)",
			`{ "i": $(foo) }`,
			pipeline.NewVarsFrom(map[string]interface{}{
				"foo": int32(1),
			}),
			nil,
			"0xfdc8ca190000000000000000000000000000000000000000000000000000000000000001",
			nil,
			"",
		},
		{
			"encode 1 to uint24",
			"asdf(uint24 i)",
			`{ "i": $(foo) }`,
			pipeline.NewVarsFrom(map[string]interface{}{
				"foo": uint32(1),
			}),
			nil,
			"0xd3f78f380000000000000000000000000000000000000000000000000000000000000001",
			nil,
			"",
		},
		// 32 bit ints.
		{
			"encode 1 to int32",
			"asdf(int32 i)",
			`{ "i": $(foo) }`,
			pipeline.NewVarsFrom(map[string]interface{}{
				"foo": int32(1),
			}),
			nil,
			"0x5124903a0000000000000000000000000000000000000000000000000000000000000001",
			nil,
			"",
		},
		{
			"encode 1 to uint32",
			"asdf(uint32 i)",
			`{ "i": $(foo) }`,
			pipeline.NewVarsFrom(map[string]interface{}{
				"foo": uint32(1),
			}),
			nil,
			"0xeea24d600000000000000000000000000000000000000000000000000000000000000001",
			nil,
			"",
		},
		// 40 bit ints.
		{
			"encode 1 to int40",
			"asdf(int40 i)",
			`{ "i": $(foo) }`,
			pipeline.NewVarsFrom(map[string]interface{}{
				"foo": int64(1),
			}),
			nil,
			"0x8fdcab050000000000000000000000000000000000000000000000000000000000000001",
			nil,
			"",
		},
		{
			"encode 1 to uint40",
			"asdf(uint40 i)",
			`{ "i": $(foo) }`,
			pipeline.NewVarsFrom(map[string]interface{}{
				"foo": uint64(1),
			}),
			nil,
			"0xcb53df3b0000000000000000000000000000000000000000000000000000000000000001",
			nil,
			"",
		},
		// 48 bit ints.
		{
			"encode 1 to int48",
			"asdf(int48 i)",
			`{ "i": $(foo) }`,
			pipeline.NewVarsFrom(map[string]interface{}{
				"foo": int64(1),
			}),
			nil,
			"0xeeab50db0000000000000000000000000000000000000000000000000000000000000001",
			nil,
			"",
		},
		{
			"encode 1 to uint48",
			"asdf(uint48 i)",
			`{ "i": $(foo) }`,
			pipeline.NewVarsFrom(map[string]interface{}{
				"foo": uint64(1),
			}),
			nil,
			"0x2d4a67fd0000000000000000000000000000000000000000000000000000000000000001",
			nil,
			"",
		},
		// 56 bit ints.
		{
			"encode 1 to int56",
			"asdf(int56 i)",
			`{ "i": $(foo) }`,
			pipeline.NewVarsFrom(map[string]interface{}{
				"foo": int64(1),
			}),
			nil,
			"0x5f4d36420000000000000000000000000000000000000000000000000000000000000001",
			nil,
			"",
		},
		{
			"encode 1 to uint56",
			"asdf(uint56 i)",
			`{ "i": $(foo) }`,
			pipeline.NewVarsFrom(map[string]interface{}{
				"foo": uint64(1),
			}),
			nil,
			"0xfe0d590c0000000000000000000000000000000000000000000000000000000000000001",
			nil,
			"",
		},
		// 64 bit ints.
		{
			"encode 1 to int64",
			"asdf(int64 i)",
			`{ "i": $(foo) }`,
			pipeline.NewVarsFrom(map[string]interface{}{
				"foo": int64(1),
			}),
			nil,
			"0x9089b4180000000000000000000000000000000000000000000000000000000000000001",
			nil,
			"",
		},
		{
			"encode 1 to uint64",
			"asdf(uint64 i)",
			`{ "i": $(foo) }`,
			pipeline.NewVarsFrom(map[string]interface{}{
				"foo": uint64(1),
			}),
			nil,
			"0x237643700000000000000000000000000000000000000000000000000000000000000001",
			nil,
			"",
		},
		// Integer sizes strictly larger than 64 bits should resolve in convertToETHABIType rather than
		// in convertToETHABIInteger, since geth uses big.Int to represent integers larger than 64 bits.
		{
			"encode 1 to int96",
			"asdf(int96 i)",
			`{ "i": $(foo) }`,
			pipeline.NewVarsFrom(map[string]interface{}{
				"foo": big.NewInt(1),
			}),
			nil,
			"0x7d14efc00000000000000000000000000000000000000000000000000000000000000001",
			nil,
			"",
		},
		{
			"encode 1 to uint96",
			"asdf(uint96 i)",
			`{ "i": $(foo) }`,
			pipeline.NewVarsFrom(map[string]interface{}{
				"foo": big.NewInt(1),
			}),
			nil,
			"0x605171600000000000000000000000000000000000000000000000000000000000000001",
			nil,
			"",
		},
		{
			"encode 1 to int128",
			"asdf(int128 i)",
			`{ "i": $(foo) }`,
			pipeline.NewVarsFrom(map[string]interface{}{
				"foo": big.NewInt(1),
			}),
			nil,
			"0x633a67090000000000000000000000000000000000000000000000000000000000000001",
			nil,
			"",
		},
		{
			"encode 1 to uint128",
			"asdf(uint128 i)",
			`{ "i": $(foo) }`,
			pipeline.NewVarsFrom(map[string]interface{}{
				"foo": big.NewInt(1),
			}),
			nil,
			"0x8209afa10000000000000000000000000000000000000000000000000000000000000001",
			nil,
			"",
		},
	}

	for _, test := range testCases {
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
				assert.Equal(t, test.expected, result.Value, fmt.Sprintf("test: %s", test.name))
			}
		})
	}
}

func TestETHABIEncode_EncodeIntegers_Overflow(t *testing.T) {
	testCases := []struct {
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
			"encode 1 to int8",
			"asdf(int8 i)",
			`{ "i": $(foo) }`,
			pipeline.NewVarsFrom(map[string]interface{}{
				"foo": int16(129),
			}),
			nil,
			"",
			pipeline.ErrBadInput,
			pipeline.ErrOverflow.Error(),
		},
		{
			"encode 1 to uint8",
			"asdf(uint8 i)",
			`{ "i": $(foo) }`,
			pipeline.NewVarsFrom(map[string]interface{}{
				"foo": uint16(257),
			}),
			nil,
			"",
			pipeline.ErrBadInput,
			pipeline.ErrOverflow.Error(),
		},
	}

	for _, test := range testCases {
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
				assert.Equal(t, test.expected, result.Value, fmt.Sprintf("test: %s", test.name))
			}
		})
	}
}
