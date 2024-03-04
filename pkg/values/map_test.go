package values

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testStruct struct {
	String      string
	StringValue *String

	Bool      bool
	BoolValue *Bool

	Byte      []byte
	ByteValue *Bytes

	Int64      int64
	Int64Value *Int64

	Int int

	Decimal      decimal.Decimal
	DecimalValue *Decimal

	Map      map[string]any
	MapValue *Map
}

func TestMap_UnwrapTo(t *testing.T) {
	im := map[string]any{
		"foo": "bar",
	}
	mv, err := NewMap(im)
	require.NoError(t, err)

	expected := &testStruct{
		String:      "a",
		StringValue: NewString("b"),

		Bool:      true,
		BoolValue: NewBool(false),

		Byte:      []byte("byte"),
		ByteValue: NewBytes([]byte("byte")),

		Int64:      int64(123),
		Int64Value: NewInt64(123),

		Int: 456,

		Decimal:      decimal.NewFromFloat(1.00),
		DecimalValue: NewDecimal(decimal.NewFromFloat(1.00)),

		Map:      im,
		MapValue: mv,
	}

	m := map[string]any{
		"string":      "a",
		"stringValue": "b",

		"bool":      true,
		"boolValue": false,

		"byte":      []byte("byte"),
		"byteValue": []byte("byte"),

		"int64":      int64(123),
		"int64Value": int64(123),

		"int": 456,

		"decimal":      decimal.NewFromFloat32(1.00),
		"decimalValue": NewDecimal(decimal.NewFromFloat(1.00)),

		"map":      im,
		"mapValue": mv,
	}
	mv, err = NewMap(m)
	require.NoError(t, err)

	got := &testStruct{}
	err = mv.UnwrapTo(got)
	require.NoError(t, err)

	assert.Equal(t, expected, got)
}
