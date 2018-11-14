package utils_test

import (
	"math"
	"math/big"
	"reflect"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/smartcontractkit/chainlink/internal/cltest"
	"github.com/smartcontractkit/chainlink/utils"
	"github.com/stretchr/testify/assert"
)

func TestUtils_NewBytes32ID(t *testing.T) {
	t.Parallel()
	id := utils.NewBytes32ID()
	assert.NotContains(t, id, "-")
}

func TestUtils_IsEmptyAddress(t *testing.T) {
	tests := []struct {
		name string
		addr common.Address
		want bool
	}{
		{"zero address", common.Address{}, true},
		{"non-zero address", cltest.NewAddress(), false},
	}

	for _, tt := range tests {
		test := tt
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			actual := utils.IsEmptyAddress(test.addr)
			assert.Equal(t, test.want, actual)
		})
	}
}

func TestUtils_StringToHex(t *testing.T) {
	tests := []struct {
		utf8 string
		hex  string
	}{
		{"abc", "0x616263"},
		{"Hi Mom!", "0x4869204d6f6d21"},
		{"", "0x"},
	}

	for _, tt := range tests {
		test := tt
		t.Run(test.utf8, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, test.hex, utils.StringToHex(test.utf8))
		})
	}
}

func TestUtils_BackoffSleeper(t *testing.T) {
	bs := utils.NewBackoffSleeper()
	d := 1 * time.Nanosecond
	bs.Min = d
	bs.Factor = 2
	assert.Equal(t, d, bs.Duration())
	bs.Sleep()
	d2 := 2 * time.Nanosecond
	assert.Equal(t, d2, bs.Duration())
}

func TestCoerceInterfaceMapToStringMap(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		input     interface{}
		want      interface{}
		wantError bool
	}{
		{"empty map", map[interface{}]interface{}{}, map[string]interface{}{}, false},
		{"simple map", map[interface{}]interface{}{"key": "value"}, map[string]interface{}{"key": "value"}, false},
		{"int map", map[int]interface{}{1: "value"}, map[int]interface{}{1: "value"}, false},
		{"error map", map[interface{}]interface{}{1: "value"}, map[int]interface{}{}, true},
		{
			"nested string map map",
			map[string]interface{}{"key": map[interface{}]interface{}{"nk": "nv"}},
			map[string]interface{}{"key": map[string]interface{}{"nk": "nv"}},
			false,
		},
		{
			"nested map map",
			map[interface{}]interface{}{"key": map[interface{}]interface{}{"nk": "nv"}},
			map[string]interface{}{"key": map[string]interface{}{"nk": "nv"}},
			false,
		},
		{
			"nested map array",
			map[interface{}]interface{}{"key": []interface{}{1, "value"}},
			map[string]interface{}{"key": []interface{}{1, "value"}},
			false,
		},
		{"empty array", []interface{}{}, []interface{}{}, false},
		{"simple array", []interface{}{1, "value"}, []interface{}{1, "value"}, false},
		{
			"error array",
			[]interface{}{map[interface{}]interface{}{1: "value"}},
			[]interface{}{},
			true,
		},
		{
			"nested array map",
			[]interface{}{map[interface{}]interface{}{"key": map[interface{}]interface{}{"nk": "nv"}}},
			[]interface{}{map[string]interface{}{"key": map[string]interface{}{"nk": "nv"}}},
			false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			decoded, err := utils.CoerceInterfaceMapToStringMap(test.input)
			if test.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.True(t, reflect.DeepEqual(test.want, decoded))
			}
		})
	}
}

func TestParseUintHex(t *testing.T) {
	t.Parallel()

	evmUint256Max, ok := (&big.Int{}).SetString("115792089237316195423570985008687907853269984665640564039457584007913129639935", 0)
	assert.True(t, ok)
	biggerThanUint256Max, ok := (&big.Int{}).SetString("121416805764108066932466369176469931665150427440758720078238275608681517825325531135", 0)
	assert.True(t, ok)
	tests := []struct {
		name      string
		input     string
		want      *big.Int
		wantError bool
	}{
		{"basic", "0x09", big.NewInt(9), false},
		{"large number", "0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff",
			evmUint256Max, false},
		{"bigger than EVM word", "0xfffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff",
			biggerThanUint256Max, false},
		{"negative", "-0xffffff", big.NewInt(-16777215), false},
		{"error", "!!!!", nil, true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := utils.ParseUintHex(test.input)
			if test.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.want, result)
			}
		})
	}
}

func TestKeccak256(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"basic", "0xf00b", "0x2433bb36d5f9b14e4fea87c2d32d79abfe34e56808b891e471f4400fca2a336c"},
		{"long input", "0xf00b2433bb36d5f9b14e4fea87c2d32d79abfe34e56808b891e471f4400fca2a336c", "0x6b917c56ad7bea7d09132b9e1e29bb5d9aa7d32d067c638dfa886bbbf6874cdf"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			input, err := hexutil.Decode(test.input)
			assert.NoError(t, err)
			result, err := utils.Keccak256(input)
			assert.NoError(t, err)

			assert.Equal(t, test.want, hexutil.Encode(result))
		})
	}
}

func TestEVMWordUint64(t *testing.T) {
	assert.Equal(t,
		[]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		utils.EVMWordUint64(1))
	assert.Equal(t,
		[]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0},
		utils.EVMWordUint64(256))
	assert.Equal(t,
		[]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
		utils.EVMWordUint64(math.MaxUint64))
}

func TestEVMWordSignedBigInt(t *testing.T) {
	val, err := utils.EVMWordSignedBigInt(&big.Int{})
	assert.NoError(t, err)
	assert.Equal(t,
		[]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		val)

	val, err = utils.EVMWordSignedBigInt(new(big.Int).SetInt64(1))
	assert.NoError(t, err)
	assert.Equal(t,
		[]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		val)

	val, err = utils.EVMWordSignedBigInt(new(big.Int).SetInt64(256))
	assert.NoError(t, err)
	assert.Equal(t,
		[]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0},
		val)

	val, err = utils.EVMWordSignedBigInt(new(big.Int).SetInt64(-1))
	assert.NoError(t, err)
	assert.Equal(t,
		[]byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
		val)

	val, err = utils.EVMWordSignedBigInt(utils.MaxInt256)
	assert.NoError(t, err)
	assert.Equal(t,
		[]byte{0x7f, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
		val)

	val, err = utils.EVMWordSignedBigInt(new(big.Int).Add(utils.MaxInt256, big.NewInt(1)))
	assert.Error(t, err)
}

func TestEVMWordBigInt(t *testing.T) {
	val, err := utils.EVMWordBigInt(&big.Int{})
	assert.NoError(t, err)
	assert.Equal(t,
		[]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		val)

	val, err = utils.EVMWordBigInt(new(big.Int).SetInt64(1))
	assert.NoError(t, err)
	assert.Equal(t,
		[]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		val)

	val, err = utils.EVMWordBigInt(new(big.Int).SetInt64(256))
	assert.NoError(t, err)
	assert.Equal(t,
		[]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0},
		val)

	val, err = utils.EVMWordBigInt(new(big.Int).SetInt64(-1))
	assert.Error(t, err)

	val, err = utils.EVMWordBigInt(utils.MaxInt256)
	assert.NoError(t, err)
	assert.Equal(t,
		[]byte{0x7f, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
		val)

	val, err = utils.EVMWordBigInt(new(big.Int).Add(utils.MaxUint256, big.NewInt(1)))
	assert.Error(t, err)
}
