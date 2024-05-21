package pipeline

import (
	"fmt"
	"reflect"
	"strings"
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
	t.Parallel()

	emptyHash := common.Hash{}
	emptyAddr := common.Address{}
	emptyFunc := [24]byte{}

	fullHash := common.HexToHash(strings.Repeat("FF", 32))
	fullAddr := common.HexToAddress(strings.Repeat("FF", 20))
	fullFunc := [24]byte{}

	oneHash := common.Hash{31: 0x1}
	oneAddr := common.Address{19: 0x1}
	oneFunc := [24]byte{23: 0x1}
	type testCase struct {
		abiType string
		exp     interface{}
	}
	for _, tt := range []struct {
		vals  []interface{}
		cases []testCase
	}{
		{[]interface{}{emptyHash, emptyHash[:], emptyHash.Hex(), string(emptyHash[:])}, []testCase{
			{"bytes", make([]byte, 32)},
			{"bytes32", [32]byte{}},
		}},
		{[]interface{}{emptyAddr, emptyAddr[:], emptyAddr.Hex(), string(emptyAddr[:])}, []testCase{
			{"bytes", make([]byte, 20)},
			{"bytes20", [20]byte{}},
			{"address", common.Address{}},
		}},
		{[]interface{}{emptyFunc, emptyFunc[:], hexutil.Encode(emptyFunc[:]), string(emptyFunc[:])}, []testCase{
			{"bytes", make([]byte, 24)},
			{"bytes24", [24]byte{}},
		}},

		{[]interface{}{fullHash, fullHash[:], fullHash.Hex()}, []testCase{
			{"bytes", fullHash[:]},
			{"bytes32", [32]byte(fullHash)},
		}},
		{[]interface{}{fullAddr, fullAddr[:], fullAddr.Hex()}, []testCase{
			{"bytes", fullAddr[:]},
			{"bytes20", [20]byte(fullAddr)},
			{"address", fullAddr},
		}},
		{[]interface{}{fullFunc, fullFunc[:], hexutil.Encode(fullFunc[:])}, []testCase{
			{"bytes", fullFunc[:]},
			{"bytes24", fullFunc},
		}},

		{[]interface{}{oneHash, oneHash[:], oneHash.Hex()}, []testCase{
			{"bytes", oneHash[:]},
			{"bytes32", [32]byte(oneHash)},
		}},
		{[]interface{}{oneAddr, oneAddr[:], oneAddr.Hex()}, []testCase{
			{"bytes", oneAddr[:]},
			{"bytes20", [20]byte{19: 0x1}},
			{"address", common.Address{19: 0x1}},
		}},
		{[]interface{}{oneFunc, oneFunc[:], hexutil.Encode(oneFunc[:])}, []testCase{
			{"bytes", oneFunc[:]},
			{"bytes24", [24]byte{23: 0x1}},
		}},

		{[]interface{}{"test", []byte("test")}, []testCase{
			{"string", "test"},
		}},

		{[]interface{}{true, "true", "1"}, []testCase{
			{"bool", true},
		}},
	} {
		tt := tt
		for _, tc := range tt.cases {
			tc := tc
			abiType := mustABIType(t, tc.abiType)
			t.Run(fmt.Sprintf("%s:%T", tc.abiType, tc.exp), func(t *testing.T) {
				for _, val := range tt.vals {
					val := val
					t.Run(fmt.Sprintf("%T", val), func(t *testing.T) {
						got, err := convertToETHABIType(val, abiType)
						require.NoError(t, err)
						require.NotNil(t, got)
						require.Equal(t, tc.exp, got)
					})
				}
			})
		}
	}
}

func Test_convertToETHABIType_Errors(t *testing.T) {
	t.Parallel()

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
	t.Parallel()

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
