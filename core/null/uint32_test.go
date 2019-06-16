package null

import (
	"encoding/json"
	"math"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	intJSON        = []byte(`12345`)
	intStringJSON  = []byte(`"12345"`)
	floatBlankJSON = []byte(`""`)
	nullJSON       = []byte(`null`)
	invalidJSON    = []byte(`:)`)
	boolJSON       = []byte(`true`)
)

func TestUint32From(t *testing.T) {
	i := Uint32From(12345)
	assertUint32(t, i, "Uint32From()")

	zero := Uint32From(0)
	if !zero.Valid {
		t.Error("Uint32From(0)", "is invalid, but should be valid")
	}
}

// TODO: Make table driven
func TestUnmarshalUint32(t *testing.T) {
	var i Uint32
	err := json.Unmarshal(intJSON, &i)
	require.NoError(t, err)
	assertUint32(t, i, "int json")

	var si Uint32
	err = json.Unmarshal(intStringJSON, &si)
	require.NoError(t, err)
	assertUint32(t, si, "int string json")

	var bi Uint32
	err = json.Unmarshal(floatBlankJSON, &bi)
	require.NoError(t, err)
	assertNullUint32(t, bi, "blank json string")

	var null Uint32
	err = json.Unmarshal(nullJSON, &null)
	require.NoError(t, err)
	assertNullUint32(t, null, "null json")

	var badType Uint32
	require.Error(t, json.Unmarshal(boolJSON, &badType))
	assertNullUint32(t, badType, "wrong type json")

	var invalid Uint32
	err = invalid.UnmarshalJSON(invalidJSON)
	if _, ok := err.(*json.SyntaxError); !ok {
		t.Errorf("expected json.SyntaxError, not %T", err)
	}
	assertNullUint32(t, invalid, "invalid json")
}

func TestUnmarshalNonIntegerNumber(t *testing.T) {
	var i Uint32
	floatJSON := []byte(`1.2345`)
	require.Error(t, json.Unmarshal(floatJSON, &i))
}

func TestUnmarshalInt64Overflow(t *testing.T) {
	int32Overflow := uint64(math.MaxInt32)

	// Max int32 should decode successfully
	var i Uint32
	err := json.Unmarshal([]byte(strconv.FormatUint(int32Overflow, 10)), &i)
	require.NoError(t, err)

	// Attempt to overflow
	int32Overflow = math.MaxUint64
	err = json.Unmarshal([]byte(strconv.FormatUint(int32Overflow, 10)), &i)
	require.Error(t, err)
}

func TestTextUnmarshalInt(t *testing.T) {
	var i Uint32
	err := i.UnmarshalText([]byte("12345"))
	require.NoError(t, err)
	assertUint32(t, i, "UnmarshalText() int")

	var blank Uint32
	err = blank.UnmarshalText([]byte(""))
	require.NoError(t, err)
	assertNullUint32(t, blank, "UnmarshalText() empty int")

	var null Uint32
	err = null.UnmarshalText([]byte("null"))
	require.NoError(t, err)
	assertNullUint32(t, null, `UnmarshalText() "null"`)
}

func TestMarshalInt(t *testing.T) {
	i := Uint32From(12345)
	data, err := json.Marshal(i)
	require.NoError(t, err)
	assertJSONEquals(t, data, "12345", "non-empty json marshal")

	// invalid values should be encoded as null
	null := NewUint32(0, false)
	data, err = json.Marshal(null)
	require.NoError(t, err)
	assertJSONEquals(t, data, "null", "null json marshal")
}

func TestMarshalIntText(t *testing.T) {
	i := Uint32From(12345)
	data, err := i.MarshalText()
	require.NoError(t, err)
	assertJSONEquals(t, data, "12345", "non-empty text marshal")

	// invalid values should be encoded as null
	null := NewUint32(0, false)
	data, err = null.MarshalText()
	require.NoError(t, err)
	assertJSONEquals(t, data, "", "null text marshal")
}

func TestUint32SetValid(t *testing.T) {
	change := NewUint32(0, false)
	assertNullUint32(t, change, "SetValid()")
	change.SetValid(12345)
	assertUint32(t, change, "SetValid()")
}

func TestUint32Scan(t *testing.T) {
	var i Uint32
	err := i.Scan(12345)
	require.NoError(t, err)
	assertUint32(t, i, "scanned int")

	err = i.Scan(int64(12345))
	require.NoError(t, err)
	assertUint32(t, i, "scanned int")

	var null Uint32
	err = null.Scan(nil)
	require.NoError(t, err)
	assertNullUint32(t, null, "scanned null")
}

func assertUint32(t *testing.T, i Uint32, from string) {
	if i.Uint32 != 12345 {
		t.Errorf("bad %s int: %d ≠ %d\n", from, i.Uint32, 12345)
	}
	if !i.Valid {
		t.Error(from, "is invalid, but should be valid")
	}
}

func assertNullUint32(t *testing.T, i Uint32, from string) {
	if i.Valid {
		t.Error(from, "is valid, but should be invalid")
	}
}

func assertJSONEquals(t *testing.T, data []byte, cmp string, from string) {
	if string(data) != cmp {
		t.Errorf("bad %s data: %s ≠ %s\n", from, data, cmp)
	}
}
