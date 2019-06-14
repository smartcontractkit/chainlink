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

func TestUSmallFrom(t *testing.T) {
	i := USmallFrom(12345)
	assertUSmall(t, i, "USmallFrom()")

	zero := USmallFrom(0)
	if !zero.Valid {
		t.Error("USmallFrom(0)", "is invalid, but should be valid")
	}
}

// TODO: Make table driven
func TestUnmarshalUSmall(t *testing.T) {
	var i USmall
	err := json.Unmarshal(intJSON, &i)
	require.NoError(t, err)
	assertUSmall(t, i, "int json")

	var si USmall
	err = json.Unmarshal(intStringJSON, &si)
	require.NoError(t, err)
	assertUSmall(t, si, "int string json")

	var bi USmall
	err = json.Unmarshal(floatBlankJSON, &bi)
	require.NoError(t, err)
	assertNullUSmall(t, bi, "blank json string")

	var null USmall
	err = json.Unmarshal(nullJSON, &null)
	require.NoError(t, err)
	assertNullUSmall(t, null, "null json")

	var badType USmall
	require.Error(t, json.Unmarshal(boolJSON, &badType))
	assertNullUSmall(t, badType, "wrong type json")

	var invalid USmall
	err = invalid.UnmarshalJSON(invalidJSON)
	if _, ok := err.(*json.SyntaxError); !ok {
		t.Errorf("expected json.SyntaxError, not %T", err)
	}
	assertNullUSmall(t, invalid, "invalid json")
}

func TestUnmarshalNonIntegerNumber(t *testing.T) {
	var i USmall
	floatJSON := []byte(`1.2345`)
	require.Error(t, json.Unmarshal(floatJSON, &i))
}

func TestUnmarshalInt64Overflow(t *testing.T) {
	int32Overflow := uint64(math.MaxInt32)

	// Max int32 should decode successfully
	var i USmall
	err := json.Unmarshal([]byte(strconv.FormatUint(int32Overflow, 10)), &i)
	require.NoError(t, err)

	// Attempt to overflow
	int32Overflow = math.MaxUint64
	err = json.Unmarshal([]byte(strconv.FormatUint(int32Overflow, 10)), &i)
	require.Error(t, err)
}

func TestTextUnmarshalInt(t *testing.T) {
	var i USmall
	err := i.UnmarshalText([]byte("12345"))
	require.NoError(t, err)
	assertUSmall(t, i, "UnmarshalText() int")

	var blank USmall
	err = blank.UnmarshalText([]byte(""))
	require.NoError(t, err)
	assertNullUSmall(t, blank, "UnmarshalText() empty int")

	var null USmall
	err = null.UnmarshalText([]byte("null"))
	require.NoError(t, err)
	assertNullUSmall(t, null, `UnmarshalText() "null"`)
}

func TestMarshalInt(t *testing.T) {
	i := USmallFrom(12345)
	data, err := json.Marshal(i)
	require.NoError(t, err)
	assertJSONEquals(t, data, "12345", "non-empty json marshal")

	// invalid values should be encoded as null
	null := NewUSmall(0, false)
	data, err = json.Marshal(null)
	require.NoError(t, err)
	assertJSONEquals(t, data, "null", "null json marshal")
}

func TestMarshalIntText(t *testing.T) {
	i := USmallFrom(12345)
	data, err := i.MarshalText()
	require.NoError(t, err)
	assertJSONEquals(t, data, "12345", "non-empty text marshal")

	// invalid values should be encoded as null
	null := NewUSmall(0, false)
	data, err = null.MarshalText()
	require.NoError(t, err)
	assertJSONEquals(t, data, "", "null text marshal")
}

func TestUSmallSetValid(t *testing.T) {
	change := NewUSmall(0, false)
	assertNullUSmall(t, change, "SetValid()")
	change.SetValid(12345)
	assertUSmall(t, change, "SetValid()")
}

func TestUSmallScan(t *testing.T) {
	var i USmall
	err := i.Scan(12345)
	require.NoError(t, err)
	assertUSmall(t, i, "scanned int")

	err = i.Scan(int64(12345))
	require.NoError(t, err)
	assertUSmall(t, i, "scanned int")

	var null USmall
	err = null.Scan(nil)
	require.NoError(t, err)
	assertNullUSmall(t, null, "scanned null")
}

func assertUSmall(t *testing.T, i USmall, from string) {
	if i.Uint32 != 12345 {
		t.Errorf("bad %s int: %d ≠ %d\n", from, i.Uint32, 12345)
	}
	if !i.Valid {
		t.Error(from, "is invalid, but should be valid")
	}
}

func assertNullUSmall(t *testing.T, i USmall, from string) {
	if i.Valid {
		t.Error(from, "is valid, but should be invalid")
	}
}

func assertJSONEquals(t *testing.T, data []byte, cmp string, from string) {
	if string(data) != cmp {
		t.Errorf("bad %s data: %s ≠ %s\n", from, data, cmp)
	}
}
