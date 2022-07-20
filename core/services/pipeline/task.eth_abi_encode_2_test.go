package pipeline_test

import (
	"context"
	"math"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
)

func TestETHABIEncodeTask2(t *testing.T) {
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
			`{
				"name": "foo_Bar__3928",
				"inputs": [
					{
						"indexed":      false,
						"name":         "u",
						"type":         "uint256"
					},
					{
						"indexed":      false,
						"name":         "b",
						"type":         "bool"
					},
					{
						"indexed":      false,
						"name":         "i",
						"type":         "int256"
					},
					{
						"indexed":      false,
						"name":         "s",
						"type":         "string"
					}
				],
				"stateMutability": "view",
				"type":            "function",
				"outputs":         []
			}`,
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
			`{
				"name": "asdf",
				"inputs": [
					{
						"name":         "b",
						"type":         "bytes32"
					},
					{
						"name":         "bs",
						"type":         "bytes"
					},
					{
						"name":         "a",
						"type":         "address"
					}
				]
			}`,
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
			`{
				"name": "chainLink",
				"inputs": [
					{
						"name":         "a",
						"type":         "address[]"
					},
					{
						"name":         "x",
						"type":         "uint80"
					},
					{
						"name":         "s",
						"type":         "uint32[2]"
					}
				]
			}`,
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
			`{
				"name": "arrayOfArrays",
				"inputs": [
					{
						"name":         "bools",
						"type":         "bool[2][]"
					},
					{
						"name":         "uints",
						"type":         "uint96[2][]"
					}
				]
			}`,
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
			`{"name": "noArgs"}`,
			``,
			pipeline.NewVarsFrom(nil),
			nil,
			"0x83c962bb",
			nil,
			"",
		},
		{
			"number too large for uint32",
			`{
				"name": "willFail",
				"inputs": [
					{
						"name":         "s",
						"type":         "uint32"
					}
				]
			}`,
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
			`{
				"name": "willFail",
				"inputs": [
					{
						"name":         "a",
						"type":         "address"
					}
				]
			}`,
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
			"no argument names",
			`{
				"name": "willFail",
				"inputs": [
					{
						"type":         "address"
					},
					{
						"type":         "uint256[]"
					}
				]
			}`,
			``,
			pipeline.NewVarsFrom(nil),
			nil,
			"",
			pipeline.ErrBadInput,
			"missing argument name",
		},
		{
			"errored task inputs",
			`{
				"name": "asdf",
				"inputs": [
					{
						"name":         "b",
						"type":         "bytes32"
					},
					{
						"name":         "bs",
						"type":         "bytes"
					},
					{
						"name":         "a",
						"type":         "address"
					}
				]
			}`,
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
			`{
				"name": "asdf",
				"inputs": [
					{
						"name":         "b",
						"type":         "bytes32"
					}
				]
			}`,
			`{ "b": $(foo)}`,
			pipeline.NewVarsFrom(map[string]interface{}{
				"foo": "0x0000000000000000000000000000000000000000000000000000000000000001",
			}),
			nil,
			"0x628507ac0000000000000000000000000000000000000000000000000000000000000001",
			nil,
			"",
		},
		{
			"tuple with map",
			`{
				"name": "call",
				"inputs": [
					{
						"name": "value",
						"type": "tuple",
						"components": [
							{
								"name": "first",
								"type": "bytes32"
							},
							{
								"name": "last",
								"type": "bool"
							}
						]
					}
				]
			}`,
			`{ "value": $(value) }`,
			pipeline.NewVarsFrom(map[string]interface{}{
				"value": map[string]interface{}{
					"first": "0x0000000000000000000000000000000000000000000000000000000000000001",
					"last":  true,
				},
			}),
			nil,
			"0xb06b167500000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000001",
			nil,
			"",
		},
		{
			"tuple with array",
			`{
				"name": "call",
				"inputs": [
					{
						"name": "value",
						"type": "tuple",
						"components": [
							{
								"name": "first",
								"type": "bytes32"
							},
							{
								"name": "last",
								"type": "bool"
							}
						]
					}
				]
			}`,
			`{ "value": $(value) }`,
			pipeline.NewVarsFrom(map[string]interface{}{
				"value": []interface{}{
					"0x0000000000000000000000000000000000000000000000000000000000000001",
					true,
				},
			}),
			nil,
			"0xb06b167500000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000001",
			nil,
			"",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			task := pipeline.ETHABIEncodeTask2{
				BaseTask: pipeline.NewBaseTask(0, "encode", nil, nil, 0),
				ABI:      test.abi,
				Data:     test.data,
			}

			result, _ := task.Run(context.Background(), logger.TestLogger(t), test.vars, test.inputs)

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
