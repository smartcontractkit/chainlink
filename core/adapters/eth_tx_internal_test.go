package adapters

import (
	"encoding/hex"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	Uint256 abi.Type
	Int256  abi.Type
	Bool    abi.Type
	Bytes32 abi.Type
	Bytes4  abi.Type
	Address abi.Type
	Bytes   abi.Type
)

func init() {
	// Static types
	Uint256, _ = abi.NewType("uint256", "", nil)
	Int256, _ = abi.NewType("int256", "", nil)
	Bool, _ = abi.NewType("bool", "", nil)
	Bytes32, _ = abi.NewType("bytes32", "", nil)
	Bytes4, _ = abi.NewType("bytes4", "", nil)
	Address, _ = abi.NewType("address", "", nil)

	// Dynamic types
	Bytes, _ = abi.NewType("bytes", "", nil)
}

func TestGetTxData(t *testing.T) {
	var tt = []struct {
		name        string
		abiEncoding []string
		argTypes    abi.Arguments // Helpers to assert the unpacking works.
		args        []interface{}
		errLike     string
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
			name:        "ints as strings",
			abiEncoding: []string{"int256", "uint256"},
			argTypes:    abi.Arguments{{Type: Int256}, {Type: Uint256}},
			args:        []interface{}{"-1234", "1234"},
			assertion: func(t *testing.T, vals []interface{}) {
				require.Len(t, vals, 2)
				assert.Equal(t, big.NewInt(-1234), vals[0])
				assert.Equal(t, big.NewInt(1234), vals[1])
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
			argTypes:    abi.Arguments{{Type: Bytes32}},
			args:        []interface{}{"0x0000000000000000000000000000000000000000000000000000000000000001"},
			assertion: func(t *testing.T, vals []interface{}) {
				require.Len(t, vals, 1)
				b, _ := hex.DecodeString("0000000000000000000000000000000000000000000000000000000000000001")
				var expected [32]byte
				copy(expected[:], b[:])
				assert.Equal(t, expected, vals[0])
			},
		},
		{
			name:        "address",
			abiEncoding: []string{"address"},
			argTypes:    abi.Arguments{{Type: Address}},
			args:        []interface{}{"0x32Be343B94f860124dC4fEe278FDCBD38C102D88"},
			assertion: func(t *testing.T, vals []interface{}) {
				require.Len(t, vals, 1)
				addr := common.HexToAddress("32Be343B94f860124dC4fEe278FDCBD38C102D88")
				assert.Equal(t, addr, vals[0])
			},
		},
		{
			name:        "invalid address",
			abiEncoding: []string{"address"},
			argTypes:    abi.Arguments{{Type: Address}},
			args:        []interface{}{"0x32Be343B94f860124dC4fEe278FDCBD38C102D"},
			errLike:     "invalid address",
		},
		{
			name:        "bytes",
			abiEncoding: []string{"bytes"},
			argTypes:    abi.Arguments{{Type: Bytes}},
			args:        []interface{}{"0x00000000000000000000000000000000000000000000000000000000000000010101"},
			assertion: func(t *testing.T, vals []interface{}) {
				require.Len(t, vals, 1)
				b, _ := hex.DecodeString("00000000000000000000000000000000000000000000000000000000000000010101")
				assert.Equal(t, b, vals[0])
			},
		},
		{
			name:        "bytes4",
			abiEncoding: []string{"bytes4"},
			argTypes:    abi.Arguments{{Type: Bytes4}},
			args:        []interface{}{"0x01010101"},
			assertion: func(t *testing.T, vals []interface{}) {
				require.Len(t, vals, 1)
				b, _ := hex.DecodeString("01010101")
				var expected [4]byte
				copy(expected[:], b[:])
				assert.Equal(t, expected, vals[0])
			},
		},
		{
			name:        "multiple bytes",
			abiEncoding: []string{"bytes", "bytes"},
			argTypes:    abi.Arguments{{Type: Bytes}, {Type: Bytes}},
			args:        []interface{}{"0x00000000000000000000000000000000000000000000000000000000000000010101", "0x0000000000000000000000000000000000000000000000000000000000000001"},
			assertion: func(t *testing.T, vals []interface{}) {
				require.Len(t, vals, 2)
				b1, _ := hex.DecodeString("00000000000000000000000000000000000000000000000000000000000000010101")
				b2, _ := hex.DecodeString("0000000000000000000000000000000000000000000000000000000000000001")
				assert.Equal(t, b1, vals[0])
				assert.Equal(t, b2, vals[1])
			},
		},
		{
			name:        "type mismatch",
			abiEncoding: []string{"uint256"},
			args:        []interface{}{"0x0123"},
			errLike:     "can't convert 0x0123 (String) to uint256",
		},
		{
			name:        "invalid bytes32",
			abiEncoding: []string{"bytes32"},
			args:        []interface{}{"0x0123"},
			errLike:     "can't convert 0x0123 (String) to bytes32", // Could consider relaxing this to just <= 32?
		},
		{
			name:        "unsupported type",
			abiEncoding: []string{"uint8"},
			args:        []interface{}{18},
			errLike:     "uint8 is unsupported", // Could consider relaxing this to just <= 32?
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			j := models.JSON{}
			d, err := j.Add(models.ResultCollectionKey, tc.args)
			require.NoError(t, err)
			b, err := getTxDataUsingABIEncoding(tc.abiEncoding,
				d.Get(models.ResultCollectionKey).Array())
			if tc.errLike != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.errLike)
				return
			}
			require.NoError(t, err)
			// We should be able to decode and get back the same args we specified.
			vals, err := tc.argTypes.UnpackValues(b)
			require.NoError(t, err)
			tc.assertion(t, vals)
		})
	}
}
