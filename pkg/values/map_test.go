package values

import (
	"math/big"
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

	List      []string
	ListValue *List

	Map      map[string]any
	MapValue *Map
}

func TestMap_UnwrapTo(t *testing.T) {
	im := map[string]any{
		"foo": "bar",
	}
	mv, err := NewMap(im)
	require.NoError(t, err)

	l := []string{"a", "b", "c"}
	lv, err := Wrap(l)
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

		List:      l,
		ListValue: lv.(*List),

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

		"list":      []string{"a", "b", "c"},
		"listValue": lv,

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

func TestMap_UnwrapTo_WorksOnMaps(t *testing.T) {
	im := map[string]any{
		"foo": "bar",
	}
	mv, err := NewMap(im)
	require.NoError(t, err)

	expected := map[string]any{
		"string":  "a",
		"bool":    true,
		"byte":    []byte("byte"),
		"int64":   int64(123),
		"decimal": decimal.NewFromFloat32(1.00),
		"map":     im,
	}
	mv, err = NewMap(expected)
	require.NoError(t, err)

	got := map[string]any{}
	err = mv.UnwrapTo(&got)
	require.NoError(t, err)

	assert.Equal(t, expected, got)
}

func TestMap_UnwrapTo_OtherMapTypes(t *testing.T) {
	testCases := []struct {
		name     string
		expected any
		got      any
	}{
		{
			name: "map[string]string",
			expected: map[string]string{
				"string":        "a",
				"anotherString": "b",
			},
			got: map[string]string{},
		},
		{
			name: "map[string]int",
			expected: map[string]int{
				"string":        1,
				"anotherString": 2,
			},
			got: map[string]int{},
		},
		{
			name: "map[string]int64",
			expected: map[string]int64{
				"string":        int64(1),
				"anotherString": int64(2),
			},
			got: map[string]int64{},
		},
		{
			name: "map[string]decimal.Decimal",
			expected: map[string]decimal.Decimal{
				"string":        decimal.NewFromFloat(1.00),
				"anotherString": decimal.NewFromFloat(1.32),
			},
			got: map[string]decimal.Decimal{},
		},
		{
			name: "map[string]big.Int",
			expected: map[string]*big.Int{
				"string":        big.NewInt(1),
				"anotherString": big.NewInt(2),
			},
			got: map[string]*big.Int{},
		},
		{
			name: "map[string][]byte",
			expected: map[string][]byte{
				"string":        []byte("hello"),
				"anotherString": []byte("world"),
			},
			got: map[string][]byte{},
		},
		{
			name: "map[string]bool",
			expected: map[string]bool{
				"string":        true,
				"anotherString": false,
			},
			got: map[string]bool{},
		},
		{
			name: "map[string]any",
			expected: map[string]any{
				"nested": map[string]any{
					"inner": "value",
				},
			},
			got: map[string]any{},
		},
		{
			name: "map[string]any nested list",
			expected: map[string]any{
				"nested": []any{
					map[string]any{
						"inner": "value",
					},
					map[string]any{
						"inner2": "value2",
					},
				},
			},
			got: map[string]any{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mv, err := Wrap(tc.expected)
			require.NoError(t, err)

			err = mv.UnwrapTo(&tc.got) // #nosec G601
			require.NoError(t, err)

			assert.Equal(t, tc.expected, tc.got)
		})
	}
}
