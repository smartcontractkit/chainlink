package pipeline

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

var testsABIDecode = []struct {
	name                  string
	abi                   string
	data                  string
	vars                  Vars
	inputs                []Result
	expected              map[string]interface{}
	expectedErrorCause    error
	expectedErrorContains string
}{
	{
		"uint256",
		"uint256 data",
		"$(data)",
		NewVarsFrom(map[string]interface{}{
			"data": "0x000000000000000000000000000000000000000000000000105ba6a589b23a81",
		}),
		nil,
		map[string]interface{}{
			"data": big.NewInt(1178718957397490305),
		},
		nil,
		"",
	},
	{
		"uint256, bool, int256, string",
		"uint256 u, bool b, int256 i, string s",
		"$(foo)",
		NewVarsFrom(map[string]interface{}{
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
		"bytes memory performData, uint256 maxLinkPayment, uint256 gasLimit, uint256 adjustedGasWei, uint256 linkEth",
		"bytes memory performData, uint256 maxLinkPayment, uint256 gasLimit, uint256 adjustedGasWei, uint256 linkEth",
		"$(foo)",
		NewVarsFrom(map[string]interface{}{
			"foo": "0x00000000000000000000000000000000000000000000000000000000000000a00000000000000000000000000000000000000000000000002cc18069c8a2800000000000000000000000000000000000000000000000000000000000002625a000000000000000000000000000000000000000000000000000000000000000c8000000000000000000000000000000000000000000000000000000000bebc20000000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000000",
		}),
		nil,
		map[string]interface{}{
			"performData":    []uint8{0x0},
			"maxLinkPayment": big.NewInt(3225000000000000000),
			"gasLimit":       big.NewInt(2500000),
			"adjustedGasWei": big.NewInt(200),
			"linkEth":        big.NewInt(200000000),
		},
		nil,
		"",
	},
	{
		"weird spaces / address, uint80[3][], bytes, bytes32",
		"address  a , uint80[3][] u , bytes b, bytes32 b32  ",
		"$(foo)",
		NewVarsFrom(map[string]interface{}{
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
		NewVarsFrom(map[string]interface{}{
			"foo": "0x000000000000000000000000deadbeefdeadbeefdeadbeefdeadbeefdeadbeef000000000000000000000000000000000000000000000000000000000000008000000000000000000000000000000000000000000000000000000000000001607374657665746f7368692073657267616d6f746f0000000000000000000000000000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000005c000000000000000000000000000000000000000000000000000000000000003d000000000000000000000000000000000000000000000000000000000000001e000000000000000000000000000000000000000000000000000000000000002100000000000000000000000000000000000000000000000000000000000000420000000000000000000000000000000000000000000000000000000000000063000000000000000000000000000000000000000000000000000000000000000c666f6f206261722062617a0a0000000000000000000000000000000000000000",
		}),
		nil,
		nil,
		ErrBadInput,
		"",
	},
	{
		"errored task inputs",
		"uint256 u, bool b, int256 i, string s",
		"$(foo)",
		NewVarsFrom(map[string]interface{}{
			"foo": "0x000000000000000000000000000000000000000000000000000000000000007b0000000000000000000000000000000000000000000000000000000000000001fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffebf0000000000000000000000000000000000000000000000000000000000000080000000000000000000000000000000000000000000000000000000000000000b666f6f206261722062617a000000000000000000000000000000000000000000",
		}),
		[]Result{{Error: errors.New("uh oh")}},
		nil,
		ErrTooManyErrors,
		"task inputs",
	},
}

func TestETHABIDecodeTask(t *testing.T) {
	for _, test := range testsABIDecode {
		test := test

		t.Run(test.name, func(t *testing.T) {
			task := ETHABIDecodeTask{
				BaseTask: NewBaseTask(0, "decode", nil, nil, 0),
				ABI:      test.abi,
				Data:     test.data,
			}

			result, runInfo := task.Run(testutils.Context(t), logger.TestLogger(t), test.vars, test.inputs)
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
