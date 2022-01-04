package pipeline

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stretchr/testify/assert"
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

func Test_convertToETHABIType_Errors(t *testing.T) {
	for _, tt := range []struct {
		val    interface{}
		errStr string
	}{
		{"0x1234", "expected 20, got 2"},
		{"0xasdfasdfasdfasdfasdfsadfasdfasdfasdfasdf", "invalid hex"},
	} {
		tt := tt
		t.Run(fmt.Sprintf("%T,%s", tt.val, tt.errStr), func(t *testing.T) {
			_, err := convertToETHABIType(tt.val, mustABIType(t, "bytes20"))
			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.errStr)
		})
	}
}

func Test_convertToETHABIBytes_Errors(t *testing.T) {
	for _, tt := range []struct {
		val    interface{}
		errStr string
	}{
		{"test", "expected 20, got 4"},
		{"12345", "expected 20, got 5"},
		{"0x1234", "expected 20, got 2"},
		{"0xzZ", "expected 20, got 1"},
		{"0xasdfasdfasdfasdfasdfsadfasdfasdfasdfasdf", "invalid hex"},
	} {
		tt := tt
		t.Run(fmt.Sprintf("%T,%s", tt.val, tt.errStr), func(t *testing.T) {
			a := reflect.TypeOf([20]byte{})
			b := reflect.ValueOf(tt.val)
			_, err := convertToETHABIBytes(a, b, 20)
			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.errStr)
		})
	}
}
