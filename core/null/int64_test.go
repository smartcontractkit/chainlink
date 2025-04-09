package null_test

import (
	"encoding/json"
	"math"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/null"
)

func TestInt64From(t *testing.T) {
	tests := []struct {
		input int64
	}{
		{12345},
		{0},
	}

	for _, test := range tests {
		t.Run(strconv.FormatInt(test.input, 10), func(t *testing.T) {
			i := null.Int64From(test.input)
			assert.True(t, i.Valid)
			assert.Equal(t, test.input, i.Int64)
		})
	}
}

func TestUnmarshalInt64_Valid(t *testing.T) {
	tests := []struct {
		name, input string
	}{
		{"int json", `12345`},
		{"int string json", `"12345"`},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var i null.Int64
			err := json.Unmarshal([]byte(test.input), &i)
			require.NoError(t, err)
			assert.True(t, i.Valid)
			assert.Equal(t, int64(12345), i.Int64)
		})
	}
}

func TestUnmarshalInt64_Invalid(t *testing.T) {
	tests := []struct {
		name, input string
	}{
		{"blank json string", `""`},
		{"null json", `null`},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var i null.Int64
			err := json.Unmarshal([]byte(test.input), &i)
			require.NoError(t, err)
			assert.False(t, i.Valid)
		})
	}
}

func TestUnmarshalInt64_Error(t *testing.T) {
	tests := []struct {
		name, input string
	}{
		{"wrong type json", `true`},
		{"invalid json", `:)`},
		{"float", `1.2345`},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var i null.Int64
			err := json.Unmarshal([]byte(test.input), &i)
			require.Error(t, err)
			assert.False(t, i.Valid)
		})
	}
}

func TestUnmarshalUint64Overflow(t *testing.T) {
	// Max int64 should decode successfully
	var i null.Int64
	err := json.Unmarshal([]byte(strconv.FormatInt(math.MaxInt64, 10)), &i)
	require.NoError(t, err)

	// Attempt to overflow
	err = json.Unmarshal([]byte(strconv.FormatUint(math.MaxUint64, 10)), &i)
	require.Error(t, err)
}

func TestTextUnmarshalInt64_Valid(t *testing.T) {
	var i null.Int64
	err := i.UnmarshalText([]byte("12345"))
	require.NoError(t, err)
	assert.True(t, i.Valid)
	assert.Equal(t, int64(12345), i.Int64)
}

func TestTextUnmarshalInt64_Invalid(t *testing.T) {
	tests := []struct {
		name, input string
	}{
		{"empty", ""},
		{"null", "null"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var i null.Int64
			err := i.UnmarshalText([]byte(test.input))
			require.NoError(t, err)
			assert.False(t, i.Valid)
		})
	}
}

func TestMarshalInt64(t *testing.T) {
	i := null.Int64From(12345)
	data, err := json.Marshal(i)
	require.NoError(t, err)
	assertJSONEquals(t, data, "12345", "non-empty json marshal")

	// invalid values should be encoded as null
	null := null.NewInt64(0, false)
	data, err = json.Marshal(null)
	require.NoError(t, err)
	assertJSONEquals(t, data, "null", "null json marshal")
}

func TestMarshalInt64Text(t *testing.T) {
	i := null.Int64From(12345)
	data, err := i.MarshalText()
	require.NoError(t, err)
	assertJSONEquals(t, data, "12345", "non-empty text marshal")

	// invalid values should be encoded as null
	null := null.NewInt64(0, false)
	data, err = null.MarshalText()
	require.NoError(t, err)
	assertJSONEquals(t, data, "", "null text marshal")
}

func TestInt64SetValid(t *testing.T) {
	change := null.NewInt64(0, false)
	change.SetValid(12345)
	assert.True(t, change.Valid)
	assert.Equal(t, int64(12345), change.Int64)
}

func TestInt64Scan(t *testing.T) {
	var i null.Int64
	err := i.Scan(int(12345))
	require.NoError(t, err)
	assert.True(t, i.Valid)
	assert.Equal(t, int64(12345), i.Int64)

	err = i.Scan(int32(12345))
	require.NoError(t, err)
	assert.True(t, i.Valid)
	assert.Equal(t, int64(12345), i.Int64)

	err = i.Scan(int64(12345))
	require.NoError(t, err)
	assert.True(t, i.Valid)
	assert.Equal(t, int64(12345), i.Int64)

	err = i.Scan(math.MaxInt64)
	require.NoError(t, err)
	assert.True(t, i.Valid)
	assert.Equal(t, int64(math.MaxInt64), i.Int64)

	err = i.Scan(uint(12345))
	require.NoError(t, err)
	assert.True(t, i.Valid)
	assert.Equal(t, int64(12345), i.Int64)

	err = i.Scan(uint64(12345))
	require.NoError(t, err)
	assert.True(t, i.Valid)
	assert.Equal(t, int64(12345), i.Int64)

	// uint64 overflows int64
	overflowingUint64 := uint64(math.MaxInt64) + 1
	err = i.Scan(overflowingUint64)
	require.Error(t, err)

	// uint overflows int64
	overflowingUint := uint(math.MaxInt64) + 1
	err = i.Scan(overflowingUint)
	require.Error(t, err)

	err = i.Scan(nil)
	require.NoError(t, err)
	assert.False(t, i.Valid)
}
