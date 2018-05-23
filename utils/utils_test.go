package utils_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
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

func TestUtils_HexToString(t *testing.T) {
	tests := []struct {
		hex     string
		utf8    string
		errored bool
	}{
		{"0x616263", "abc", false},
		{"616263", "abc", false},
		{"0x4869204d6f6d21", "Hi Mom!", false},
		{"0x", "", false},
		{"uh oh", "", true},
	}

	for _, tt := range tests {
		test := tt
		t.Run(test.hex, func(t *testing.T) {
			t.Parallel()
			actualUtf8, err := utils.HexToString(test.hex)
			assert.Equal(t, test.errored, err != nil)
			assert.Equal(t, test.utf8, actualUtf8)
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
