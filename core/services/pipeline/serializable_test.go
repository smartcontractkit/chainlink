package pipeline_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/services/pipeline"
)

func TestSerializable_Empty(t *testing.T) {
	t.Parallel()

	s := pipeline.Serializable{}
	require.True(t, s.Empty())

	s.Valid = true
	require.False(t, s.Empty())

	s1 := pipeline.Serializable{
		Val:   "test",
		Valid: false,
	}
	require.True(t, s1.Empty())

	s1.Valid = true
	require.False(t, s1.Empty())
}

func TestSerializable_ScanValue(t *testing.T) {
	t.Parallel()

	s := pipeline.NewValidSerializable(map[string]interface{}{
		"nested": map[string]interface{}{
			"nested": map[string]interface{}{
				"nested": map[string]interface{}{
					"bin": []byte{0x11, 0x22},
				},
			},
		},
	})

	v, err := s.Value()
	require.NoError(t, err)
	require.NotEmpty(t, interface{}(v))

	s2 := &pipeline.Serializable{}
	err = s2.Scan(v)
	require.NoError(t, err)
	require.False(t, s2.Empty())
	require.Equal(t, s, s2)
}

func TestSerializable_String(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		val      interface{}
		valid    bool
		expected string
	}{
		{
			name:     "nil value",
			val:      nil,
			valid:    true,
			expected: `"nil"`,
		},
		{
			name:     "invalid Serializable",
			val:      "foo",
			valid:    false,
			expected: "invalid",
		},
		{
			name:     "[]byte value",
			val:      []byte{0x11, 0x22, 0x33},
			valid:    true,
			expected: `"0x112233"`,
		},
		{
			name:     "array of []byte values",
			val:      [][]byte{{0x11, 0x22}, {0x33, 0x44}},
			valid:    true,
			expected: `["0x1122","0x3344"]`,
		},
		{
			name: "map[string]interface{} value with bytes",
			val: map[string]interface{}{
				"foo": "bar",
				"bin": []byte{0x11, 0x22},
			},
			valid:    true,
			expected: `{"bin":"0x1122","foo":"bar"}`,
		},
		{
			name: "map[string]interface{} value with array of bytes",
			val: map[string]interface{}{
				"foo":  "bar",
				"bins": [][]byte{{0x11, 0x22}, {0x33, 0x44}},
			},
			valid:    true,
			expected: `{"bins":["0x1122","0x3344"],"foo":"bar"}`,
		},
		{
			name: "nested map[string]interface{} values",
			val: map[string]interface{}{
				"nested": map[string]interface{}{
					"nested": map[string]interface{}{
						"nested": map[string]interface{}{
							"bin": []byte{0x11, 0x22},
						},
					},
				},
			},
			valid:    true,
			expected: `{"nested":{"nested":{"nested":{"bin":"0x1122"}}}}`,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			s := pipeline.Serializable{
				Val:   tc.val,
				Valid: tc.valid,
			}
			require.Equal(t, tc.expected, s.String())
		})
	}
}
