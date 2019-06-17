package utils_test

import (
	"fmt"
	"math/big"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/multierr"
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
	assert.Equal(t, time.Duration(0), bs.Duration(), "should initially return immediately")
	bs.Sleep()

	d := 1 * time.Nanosecond
	bs.Min = d
	bs.Factor = 2
	assert.Equal(t, d, bs.Duration())
	bs.Sleep()

	d2 := 2 * time.Nanosecond
	assert.Equal(t, d2, bs.Duration())

	bs.Reset()
	assert.Equal(t, time.Duration(0), bs.Duration(), "should initially return immediately")
}

func TestUtils_DurationFromNow(t *testing.T) {
	t.Parallel()
	future := time.Now().Add(time.Second)
	duration := utils.DurationFromNow(future)
	assert.True(t, 0 < duration)
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

// From https://github.com/ethereum/EIPs/blob/master/EIPS/eip-55.md#test-cases
var testAddresses = []string{
	"0x52908400098527886E0F7030069857D2E4169EE7",
	"0x8617E340B3D01FA5F11F306F4090FD50E238070D",
	"0xde709f2102306220921060314715629080e2fb77",
	"0x27b1fdb04752bbc536007a920d24acb045561c26",
	"0x5aAeb6053F3E94C9b9A09f33669435E7Ef1BeAed",
	"0xfB6916095ca1df60bB79Ce92cE3Ea74c37c5d359",
	"0xdbF03B407c01E7cD3CBea99509d93f8DDDC8C6FB",
	"0xD1220A0cf47c7B9Be7A2E6BA89F429762e7b9aDb",
}

func TestClient_EIP55CapitalizedAddress(t *testing.T) {
	valid := utils.EIP55CapitalizedAddress
	for _, address := range testAddresses {
		assert.True(t, valid(address))
		assert.False(t, valid(strings.ToLower(address)) &&
			valid(strings.ToUpper(address)))
	}
}

func TestClient_ParseEthereumAddress(t *testing.T) {
	parse := utils.ParseEthereumAddress
	for _, address := range testAddresses {
		a1, err := parse(address)
		assert.NoError(t, err)
		no0xPrefix := address[2:]
		a2, err := parse(no0xPrefix)
		assert.NoError(t, err)
		assert.True(t, a1 == a2)
		_, lowerErr := parse(strings.ToLower(address))
		_, upperErr := parse(strings.ToUpper(address))
		shouldBeError := multierr.Combine(lowerErr, upperErr)
		assert.Error(t, shouldBeError)
		assert.True(t, strings.Contains(shouldBeError.Error(), no0xPrefix))
	}
	_, notHexErr := parse("0xCeci n'est pas une chaîne hexadécimale")
	assert.Error(t, notHexErr)
	_, tooLongErr := parse("0x0123456789abcdef0123456789abcdef0123456789abcdef")
	assert.Error(t, tooLongErr)
}

func TestMinBigs(t *testing.T) {
	tests := []struct {
		min, max string
	}{
		{"0", "0"},
		{"-1", "0"},
		{"99", "100"},
		{"0", "1"},
		{"4294967295", "4294967296"},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%s < %s", test.min, test.max), func(t *testing.T) {
			left, ok := big.NewInt(0).SetString(test.min, 10)
			require.True(t, ok)
			right, ok := big.NewInt(0).SetString(test.max, 10)
			require.True(t, ok)

			min := utils.MinBigs(left, right)
			assert.Equal(t, left, min)
		})
	}
}
