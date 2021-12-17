package pipeline

import (
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stretchr/testify/require"
)

func mustABIType(t *testing.T, ty string) abi.Type {
	typ, err := abi.NewType(ty, "", nil)
	require.NoError(t, err)
	return typ
}

func Test_convertToETHABIType(t *testing.T) {
	emptyHash := common.Hash{}
	emptyAddr := common.Address{}
	emptyFunc := [24]byte{}
	for _, tt := range []struct {
		val     interface{}
		abiType string
	}{
		{"test", "string"},

		{emptyHash, "bytes"},
		{emptyHash[:], "bytes"},
		{emptyHash.Hex(), "bytes"},
		{emptyHash, "bytes32"},
		{emptyHash[:], "bytes32"},
		{emptyHash.Hex(), "bytes32"},

		{emptyAddr, "bytes"},
		{emptyAddr, "bytes20"},
		{emptyAddr[:], "bytes20"},
		{emptyAddr.Hex(), "bytes20"},
		{emptyAddr, "address"},
		{emptyAddr[:], "address"},
		{emptyAddr.Hex(), "address"},

		{emptyFunc, "bytes"},
		{emptyFunc, "bytes24"},
		{emptyFunc[:], "bytes24"},
		{hexutil.Encode(emptyFunc[:]), "bytes24"},
	} {
		tt := tt
		t.Run(fmt.Sprintf("%T,%s", tt.val, tt.abiType), func(t *testing.T) {
			got, err := convertToETHABIType(tt.val, mustABIType(t, tt.abiType))
			require.NoError(t, err)
			t.Logf("got: 0x%x\n", got)
			require.NotNil(t, got)
			//TODO more validation
		})
	}
}
