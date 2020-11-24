package adapters

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	Uint256 abi.Type
	Int256  abi.Type
	Bool    abi.Type
)

func init() {
	Uint256, _ = abi.NewType("uint256", "", nil)
	Int256, _ = abi.NewType("int256", "", nil)
	Bool, _ = abi.NewType("bool", "", nil)
}

func TestGetTxData(t *testing.T) {
	var tt = []struct {
		name        string
		abiEncoding []string
		argTypes    abi.Arguments
		args        []interface{}
		err         error
		assertion   func(t *testing.T, vals []interface{})
	}{
		{
			name:        "uint256",
			abiEncoding: []string{"uint256"},
			argTypes:    abi.Arguments{{Type: Uint256}},
			args:        []interface{}{1234},
			assertion: func(t *testing.T, vals []interface{}) {
				require.Len(t, vals, 1)
				assert.Equal(t, big.NewInt(1234), vals[0])
			},
		},
		{
			name:        "int256",
			abiEncoding: []string{"int256"},
			argTypes:    abi.Arguments{{Type: Int256}},
			args:        []interface{}{-1234},
			assertion: func(t *testing.T, vals []interface{}) {
				require.Len(t, vals, 1)
				assert.Equal(t, big.NewInt(-1234), vals[0])
			},
		},
		{
			name:        "multiple int256",
			abiEncoding: []string{"int256", "int256"},
			argTypes:    abi.Arguments{{Type: Int256}, {Type: Int256}},
			args:        []interface{}{-1234, 10923810298},
			assertion: func(t *testing.T, vals []interface{}) {
				require.Len(t, vals, 2)
				assert.Equal(t, big.NewInt(-1234), vals[0])
				assert.Equal(t, big.NewInt(10923810298), vals[1])
			},
		},
		{
			name:        "bool",
			abiEncoding: []string{"bool"},
			argTypes:    abi.Arguments{{Type: Bool}},
			args:        []interface{}{true},
			assertion: func(t *testing.T, vals []interface{}) {
				require.Len(t, vals, 1)
				assert.Equal(t, true, vals[0])
			},
		},
		{
			name:        "bytes32",
			abiEncoding: []string{"bytes32"},
			argTypes:    abi.Arguments{{Type: Bool}},
			args:        []interface{}{true},
			assertion: func(t *testing.T, vals []interface{}) {
				require.Len(t, vals, 1)
				assert.Equal(t, true, vals[0])
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			j := models.JSON{}
			d, err := j.Add("__chainlink_result_collection__", tc.args)
			require.NoError(t, err)
			b, err := getTxData2(&EthTx{ABIEncoding: tc.abiEncoding,
				FunctionSelector: models.HexToFunctionSelector("0x70a08231"),
			}, d)
			if tc.err != nil {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			// TODO We should be able to decode and get back the same args
			//t.Log(hex.EncodeToString(b))
			vals, err := tc.argTypes.UnpackValues(b[4:])
			require.NoError(t, err)
			tc.assertion(t, vals)
		})
	}
}
