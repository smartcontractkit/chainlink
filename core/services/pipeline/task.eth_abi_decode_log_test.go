package pipeline_test

import (
	"context"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/core/utils"
)

func TestETHABIDecodeLogTask(t *testing.T) {
	tests := []struct {
		name                  string
		abi                   string
		data                  string
		topics                string
		vars                  pipeline.Vars
		inputs                []pipeline.Result
		expected              map[string]interface{}
		expectedErrorCause    error
		expectedErrorContains string
	}{
		{
			"AggregatorV2V3#NewRound",
			"NewRound(uint256 indexed roundId, address indexed startedBy, uint256 startedAt)",
			`$(foo.data)`,
			`$(foo.topics)`,
			pipeline.NewVarsFrom(map[string]interface{}{
				"foo": map[string]interface{}{
					"address": common.HexToAddress("0x2fCeA879fDC9FE5e90394faf0CA644a1749d0ad6"),
					"data":    hexutil.MustDecode("0x000000000000000000000000000000000000000000000000000000000000000f"),
					"topics": []common.Hash{
						common.HexToHash("0x0109fc6f55cf40689f02fbaad7af7fe7bbac8a3d2186600afc7d3e10cac60271"),
						common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000009"),
						common.HexToHash("0x000000000000000000000000f17f52151ebef6c7334fad080c5704d77216b732"),
					},
				},
			}),
			nil,
			map[string]interface{}{
				"roundId":   big.NewInt(9),
				"startedBy": common.HexToAddress("0xf17f52151ebef6c7334fad080c5704d77216b732"),
				"startedAt": big.NewInt(15),
			},
			nil,
			"",
		},
		{
			"Operator#OracleRequest",
			"OracleRequest(bytes32 indexed specId, address requester, bytes32 requestId, uint256 payment, address callbackAddr, bytes4 callbackFunctionId, uint256 cancelExpiration, uint256 dataVersion, bytes data)",
			`$(foo.data)`,
			`$(foo.topics)`,
			pipeline.NewVarsFrom(map[string]interface{}{
				"foo": map[string]interface{}{
					"address": common.HexToAddress("0x2fCeA879fDC9FE5e90394faf0CA644a1749d0ad6"),
					"data":    hexutil.MustDecode("0x000000000000000000000000deadbeefdeadbeefdeadbeefdeadbeefdeadbeef74686520726571756573742069640000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000020a6e000000000000000000000000cafebabecafebabecafebabecafebabecafebabe61736466000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000003039000000000000000000000000000000000000000000000000000000000000d431000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000000147374657665746f7368692073657267616d6f746f000000000000000000000000"),
					"topics": []common.Hash{
						common.HexToHash("0xd8d7ecc4800d25fa53ce0372f13a416d98907a7ef3d8d3bdd79cf4fe75529c65"),
						common.HexToHash("0x746865206a6f6220696400000000000000000000000000000000000000000000"),
					},
				},
			}),
			nil,
			map[string]interface{}{
				"specId":             utils.Bytes32FromString("the job id"),
				"requester":          common.HexToAddress("0xDeaDbeefdEAdbeefdEadbEEFdeadbeEFdEaDbeeF"),
				"requestId":          utils.Bytes32FromString("the request id"),
				"payment":            big.NewInt(133742),
				"callbackAddr":       common.HexToAddress("0xCafEBAbECAFEbAbEcaFEbabECAfebAbEcAFEBaBe"),
				"callbackFunctionId": utils.Bytes4FromString("asdf"),
				"cancelExpiration":   big.NewInt(12345),
				"dataVersion":        big.NewInt(54321),
				"data":               []byte("stevetoshi sergamoto"),
			},
			nil,
			"",
		},
		{
			"Operator#AuthorizedSendersChanged",
			"AuthorizedSendersChanged(address[] senders)",
			`$(foo.data)`,
			`$(foo.topics)`,
			pipeline.NewVarsFrom(map[string]interface{}{
				"foo": map[string]interface{}{
					"address": common.HexToAddress("0x2fCeA879fDC9FE5e90394faf0CA644a1749d0ad6"),
					"data":    hexutil.MustDecode("0x00000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000002000000000000000000000000deadbeefdeadbeefdeadbeefdeadbeefdeadbeef000000000000000000000000cafebabecafebabecafebabecafebabecafebabe"),
					"topics": []common.Hash{
						common.HexToHash("0xe720bc96024900ba647b8faa27766eb59f72cadf3c7ec34a7365c999f78320db"),
					},
				},
			}),
			nil,
			map[string]interface{}{
				"senders": []common.Address{
					common.HexToAddress("0xDeaDbeefdEAdbeefdEadbEEFdeadbeEFdEaDbeeF"),
					common.HexToAddress("0xCafEBAbECAFEbAbEcaFEbabECAfebAbEcAFEBaBe"),
				},
			},
			nil,
			"",
		},

		{
			"missing arg name",
			"SomeEvent(bytes32)",
			`$(foo.data)`,
			`$(foo.topics)`,
			pipeline.NewVarsFrom(map[string]interface{}{
				"foo": map[string]interface{}{
					"address": common.HexToAddress("0x2fCeA879fDC9FE5e90394faf0CA644a1749d0ad6"),
					"data":    hexutil.MustDecode("0x000000000000000000000000deadbeefdeadbeefdeadbeefdeadbeefdeadbeef74686520726571756573742069640000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000020a6e000000000000000000000000cafebabecafebabecafebabecafebabecafebabe61736466000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000003039000000000000000000000000000000000000000000000000000000000000d431000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000000147374657665746f7368692073657267616d6f746f000000000000000000000000"),
					"topics": []common.Hash{
						common.HexToHash("0xd8d7ecc4800d25fa53ce0372f13a416d98907a7ef3d8d3bdd79cf4fe75529c65"),
					},
				},
			}),
			nil,
			nil,
			pipeline.ErrBadInput,
			"bad ABI specification",
		},
		{
			"missing arg name (with 'indexed' modifier)",
			"SomeEvent(bytes32 indexed)",
			`$(foo.data)`,
			`$(foo.topics)`,
			pipeline.NewVarsFrom(map[string]interface{}{
				"foo": map[string]interface{}{
					"address": common.HexToAddress("0x2fCeA879fDC9FE5e90394faf0CA644a1749d0ad6"),
					"data":    hexutil.MustDecode("0x000000000000000000000000deadbeefdeadbeefdeadbeefdeadbeefdeadbeef74686520726571756573742069640000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000020a6e000000000000000000000000cafebabecafebabecafebabecafebabecafebabe61736466000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000003039000000000000000000000000000000000000000000000000000000000000d431000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000000147374657665746f7368692073657267616d6f746f000000000000000000000000"),
					"topics": []common.Hash{
						common.HexToHash("0xd8d7ecc4800d25fa53ce0372f13a416d98907a7ef3d8d3bdd79cf4fe75529c65"),
					},
				},
			}),
			nil,
			nil,
			pipeline.ErrBadInput,
			"bad ABI specification",
		},
		{
			"missing topic data",
			"OracleRequest(bytes32 indexed specId, address requester, bytes32 requestId, uint256 payment, address callbackAddr, bytes4 callbackFunctionId, uint256 cancelExpiration, uint256 dataVersion, bytes data)",
			`$(foo.data)`,
			`$(foo.topics)`,
			pipeline.NewVarsFrom(map[string]interface{}{
				"foo": map[string]interface{}{
					"address": common.HexToAddress("0x2fCeA879fDC9FE5e90394faf0CA644a1749d0ad6"),
					"data":    hexutil.MustDecode("0x000000000000000000000000deadbeefdeadbeefdeadbeefdeadbeefdeadbeef74686520726571756573742069640000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000020a6e000000000000000000000000cafebabecafebabecafebabecafebabecafebabe61736466000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000003039000000000000000000000000000000000000000000000000000000000000d431000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000000147374657665746f7368692073657267616d6f746f000000000000000000000000"),
					"topics": []common.Hash{
						common.HexToHash("0xd8d7ecc4800d25fa53ce0372f13a416d98907a7ef3d8d3bdd79cf4fe75529c65"),
					},
				},
			}),
			nil,
			nil,
			pipeline.ErrBadInput,
			"topic/field count mismatch",
		},
		{
			"not enough data: len(data) % 32 != 0",
			"OracleRequest(bytes32 indexed specId, address requester, bytes32 requestId, uint256 payment, address callbackAddr, bytes4 callbackFunctionId, uint256 cancelExpiration, uint256 dataVersion, bytes data)",
			`$(foo.data)`,
			`$(foo.topics)`,
			pipeline.NewVarsFrom(map[string]interface{}{
				"foo": map[string]interface{}{
					"address": common.HexToAddress("0x2fCeA879fDC9FE5e90394faf0CA644a1749d0ad6"),
					"data":    hexutil.MustDecode("0x000000000000000000000000deadbeefdeadbeefdeadbeefdeadbeefdeadbeef74686520726571756573742069640000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000020a6e000000000000000000000000cafebabecafebabecafebabecafebabecafebabe61736466000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000003039000000000000000000000000000000000000000000000000000000000000d4310000000000000000000000000000000000000000000000000000"),
					"topics": []common.Hash{
						common.HexToHash("0xd8d7ecc4800d25fa53ce0372f13a416d98907a7ef3d8d3bdd79cf4fe75529c65"),
						common.HexToHash("0x746865206a6f6220696400000000000000000000000000000000000000000000"),
					},
				},
			}),
			nil,
			nil,
			pipeline.ErrBadInput,
			"length insufficient 250 require 256",
		},
		{
			"not enough data: len(data) < len(non-indexed args) * 32",
			"OracleRequest(bytes32 indexed specId, address requester, bytes32 requestId, uint256 foobar)",
			`$(foo.data)`,
			`$(foo.topics)`,
			pipeline.NewVarsFrom(map[string]interface{}{
				"foo": map[string]interface{}{
					"address": common.HexToAddress("0x2fCeA879fDC9FE5e90394faf0CA644a1749d0ad6"),
					"data":    hexutil.MustDecode("0x000000000000000000000000deadbeefdeadbeefdeadbeefdeadbeefdeadbeef7468652072657175657374206964000000000000000000000000000000000000"),
					"topics": []common.Hash{
						common.HexToHash("0xd8d7ecc4800d25fa53ce0372f13a416d98907a7ef3d8d3bdd79cf4fe75529c65"),
						common.HexToHash("0x746865206a6f6220696400000000000000000000000000000000000000000000"),
					},
				},
			}),
			nil,
			nil,
			pipeline.ErrBadInput,
			"length insufficient 64 require 96",
		},
		{
			"errored task inputs",
			"NewRound(uint256 indexed roundId, address indexed startedBy, uint256 startedAt)",
			`$(foo.data)`,
			`$(foo.topics)`,
			pipeline.NewVarsFrom(map[string]interface{}{
				"foo": map[string]interface{}{
					"address": common.HexToAddress("0x2fCeA879fDC9FE5e90394faf0CA644a1749d0ad6"),
					"data":    hexutil.MustDecode("0x000000000000000000000000000000000000000000000000000000000000000f"),
					"topics": []common.Hash{
						common.HexToHash("0x0109fc6f55cf40689f02fbaad7af7fe7bbac8a3d2186600afc7d3e10cac60271"),
						common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000009"),
						common.HexToHash("0x000000000000000000000000f17f52151ebef6c7334fad080c5704d77216b732"),
					},
				},
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
			task := pipeline.ETHABIDecodeLogTask{
				BaseTask: pipeline.NewBaseTask(0, "decodelog", nil, nil, 0),
				ABI:      test.abi,
				Data:     test.data,
				Topics:   test.topics,
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
