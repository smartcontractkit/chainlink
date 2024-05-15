package values

import (
	"math"
	"math/big"
	"testing"

	"github.com/mitchellh/mapstructure"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type TestValueEvent struct {
	TriggerType string       `json:"triggerType"`
	ID          string       `json:"id"`
	Timestamp   string       `json:"timestamp"`
	Payload     []TestReport `json:"payload"`
}

type TestReport struct {
	FeedID     int64  `json:"feedId"`
	FullReport string `json:"fullreport"`
}

func Test_Value(t *testing.T) {
	testCases := []struct {
		name     string
		newValue func() (any, Value, error)
		equal    func(t *testing.T, expected any, unwrapped any)
	}{
		{
			name: "map",
			newValue: func() (any, Value, error) {
				m := map[string]any{
					"hello": "world",
				}
				mv, err := NewMap(m)
				return m, mv, err
			},
		},
		{
			name: "list",
			newValue: func() (any, Value, error) {
				l := []any{
					1,
					"2",
					decimal.NewFromFloat(1.0),
				}
				lv, err := NewList(l)
				return l, lv, err
			},
			equal: func(t *testing.T, expected any, unwrapped any) {
				e, u := expected.([]any), unwrapped.([]any)
				assert.Equal(t, int64(e[0].(int)), u[0])
				assert.Equal(t, e[1], u[1])
				assert.Equal(t, e[2].(decimal.Decimal).String(), u[2].(decimal.Decimal).String())
			},
		},
		{
			name: "decimal",
			newValue: func() (any, Value, error) {
				dec, err := decimal.NewFromString("1.03")
				if err != nil {
					return nil, nil, err
				}
				decv := NewDecimal(dec)
				return dec, decv, err
			},
		},
		{
			name: "string",
			newValue: func() (any, Value, error) {
				s := "hello"
				sv := NewString(s)
				return s, sv, nil
			},
		},
		{
			name: "bytes",
			newValue: func() (any, Value, error) {
				b := []byte("hello")
				bv := NewBytes(b)
				return b, bv, nil
			},
		},
		{
			name: "bool",
			newValue: func() (any, Value, error) {
				b := true
				bv := NewBool(b)
				return b, bv, nil
			},
		},
		{
			name: "bigInt",
			newValue: func() (any, Value, error) {
				b := big.NewInt(math.MaxInt64)
				bv := NewBigInt(b)
				return b, bv, nil
			},
		},
		{
			name: "recursive map",
			newValue: func() (any, Value, error) {
				m := map[string]any{
					"hello": map[string]any{
						"world": "foo",
					},
					"baz": []any{
						int64(1), int64(2), int64(3),
					},
				}
				mv, err := NewMap(m)
				return m, mv, err
			},
		},
		{
			name: "struct",
			newValue: func() (any, Value, error) {
				var v TestReport
				m := map[string]any{
					"FeedID":     int64(2),
					"FullReport": "hello",
				}
				err := mapstructure.Decode(m, &v)
				if err != nil {
					return nil, nil, err
				}
				vv, err := Wrap(v)
				return m, vv, err
			},
		},
		{
			name: "structPointer",
			newValue: func() (any, Value, error) {
				var v TestReport
				m := map[string]any{
					"FeedID":     int64(3),
					"FullReport": "world",
				}
				err := mapstructure.Decode(m, &v)
				if err != nil {
					return nil, nil, err
				}
				vv, err := Wrap(&v)
				return m, vv, err
			},
		},
		{
			name: "nestedStruct",
			newValue: func() (any, Value, error) {
				var v TestValueEvent
				m := map[string]any{
					"TriggerType": "mercury",
					"ID":          "123",
					"Timestamp":   "123",
					"Payload": []any{
						map[string]any{
							"FeedID":     int64(4),
							"FullReport": "hello",
						},
						map[string]any{
							"FeedID":     int64(5),
							"FullReport": "world",
						},
					},
				}
				err := mapstructure.Decode(m, &v)
				if err != nil {
					return nil, nil, err
				}
				vv, err := Wrap(v)
				return m, vv, err
			},
		},
		{
			name: "map of values",
			newValue: func() (any, Value, error) {
				bar := "bar"
				str := &String{Underlying: bar}
				l, err := NewList([]any{1, 2, 3})
				if err != nil {
					return nil, nil, err
				}
				m := map[string]any{
					"hello": map[string]any{
						"string": str,
						"nil":    nil,
						"list":   l,
					},
				}
				mv, err := NewMap(m)

				list := []any{int64(1), int64(2), int64(3)}
				expectedUnwrapped := map[string]any{
					"hello": map[string]any{
						"string": bar,
						"nil":    nil,
						"list":   list,
					},
				}

				return expectedUnwrapped, mv, err
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(st *testing.T) {
			originalValue, wrapped, err := tc.newValue()
			require.NoError(st, err)

			pb := Proto(wrapped)

			rehydratedValue := FromProto(pb)
			assert.Equal(st, wrapped, rehydratedValue)

			unwrapped, err := Unwrap(rehydratedValue)
			require.NoError(st, err)
			if tc.equal != nil {
				tc.equal(st, originalValue, unwrapped)
			} else {
				assert.Equal(st, originalValue, unwrapped)
			}
		})
	}
}
