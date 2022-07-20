package null_test

import (
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"testing"

	"github.com/smartcontractkit/chainlink/core/null"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUint32From(t *testing.T) {
	tests := []struct {
		input uint32
	}{
		{12345},
		{0},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%d", test.input), func(t *testing.T) {
			i := null.Uint32From(test.input)
			assert.True(t, i.Valid)
			assert.Equal(t, test.input, i.Uint32)
		})
	}
}

func TestUnmarshalUint32_Valid(t *testing.T) {
	tests := []struct {
		name, input string
	}{
		{"int json", `12345`},
		{"int string json", `"12345"`},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var i null.Uint32
			err := json.Unmarshal([]byte(test.input), &i)
			require.NoError(t, err)
			assert.True(t, i.Valid)
			assert.Equal(t, uint32(12345), i.Uint32)
		})
	}
}

func TestUnmarshalUint32_Invalid(t *testing.T) {
	tests := []struct {
		name, input string
	}{
		{"blank json string", `""`},
		{"null json", `null`},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var i null.Uint32
			err := json.Unmarshal([]byte(test.input), &i)
			require.NoError(t, err)
			assert.False(t, i.Valid)
		})
	}
}

func TestUnmarshalUint32_Error(t *testing.T) {
	tests := []struct {
		name, input string
	}{
		{"wrong type json", `true`},
		{"invalid json", `:)`},
		{"float", `1.2345`},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var i null.Uint32
			err := json.Unmarshal([]byte(test.input), &i)
			require.Error(t, err)
			assert.False(t, i.Valid)
		})
	}
}

func TestUnmarshalUint32Overflow(t *testing.T) {
	maxUint32 := uint64(math.MaxUint32)

	// Max uint32 should decode successfully
	var i null.Uint32
	err := json.Unmarshal([]byte(strconv.FormatUint(maxUint32, 10)), &i)
	require.NoError(t, err)

	// Attempt to overflow
	err = json.Unmarshal([]byte(strconv.FormatUint(maxUint32+1, 10)), &i)
	require.Error(t, err)
}

func TestTextUnmarshalInt_Valid(t *testing.T) {
	var i null.Uint32
	err := i.UnmarshalText([]byte("12345"))
	require.NoError(t, err)
	assert.True(t, i.Valid)
	assert.Equal(t, uint32(12345), i.Uint32)
}

func TestTextUnmarshalInt_Invalid(t *testing.T) {
	tests := []struct {
		name, input string
	}{
		{"empty", ""},
		{"null", "null"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var i null.Uint32
			err := i.UnmarshalText([]byte(test.input))
			require.NoError(t, err)
			assert.False(t, i.Valid)
		})
	}
}

func TestMarshalInt(t *testing.T) {
	i := null.Uint32From(12345)
	data, err := json.Marshal(i)
	require.NoError(t, err)
	assertJSONEquals(t, data, "12345", "non-empty json marshal")

	// invalid values should be encoded as null
	null := null.NewUint32(0, false)
	data, err = json.Marshal(null)
	require.NoError(t, err)
	assertJSONEquals(t, data, "null", "null json marshal")
}

func TestMarshalIntText(t *testing.T) {
	i := null.Uint32From(12345)
	data, err := i.MarshalText()
	require.NoError(t, err)
	assertJSONEquals(t, data, "12345", "non-empty text marshal")

	// invalid values should be encoded as null
	null := null.NewUint32(0, false)
	data, err = null.MarshalText()
	require.NoError(t, err)
	assertJSONEquals(t, data, "", "null text marshal")
}

func TestUint32SetValid(t *testing.T) {
	change := null.NewUint32(0, false)
	change.SetValid(12345)
	assert.True(t, change.Valid)
	assert.Equal(t, uint32(12345), change.Uint32)
}

func TestUint32Scan(t *testing.T) {
	var i null.Uint32
	err := i.Scan(12345)
	require.NoError(t, err)
	assert.True(t, i.Valid)
	assert.Equal(t, uint32(12345), i.Uint32)

	err = i.Scan(int64(12345))
	require.NoError(t, err)
	assert.True(t, i.Valid)
	assert.Equal(t, uint32(12345), i.Uint32)

	// int64 overflows uint32
	err = i.Scan(int64(math.MaxInt64))
	require.Error(t, err)

	err = i.Scan(uint(12345))
	require.NoError(t, err)
	assert.True(t, i.Valid)
	assert.Equal(t, uint32(12345), i.Uint32)

	// uint overflows uint32
	err = i.Scan(uint(math.MaxUint64))
	require.Error(t, err)

	err = i.Scan(uint32(12345))
	require.NoError(t, err)
	assert.True(t, i.Valid)
	assert.Equal(t, uint32(12345), i.Uint32)

	err = i.Scan(nil)
	require.NoError(t, err)
	assert.False(t, i.Valid)
}

func assertJSONEquals(t *testing.T, data []byte, cmp string, from string) {
	if string(data) != cmp {
		t.Errorf("bad %s data: %s ≠ %s\n", from, data, cmp)
	}
}
