package pipeline_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/services/pipeline"
)

func TestSerializable_Empty(t *testing.T) {
	t.Parallel()

	s := pipeline.CBORSerializable{}
	require.True(t, s.Empty())

	s.Valid = true
	require.False(t, s.Empty())

	s1 := pipeline.CBORSerializable{
		Val:   "test",
		Valid: false,
	}
	require.True(t, s1.Empty())

	s1.Valid = true
	require.False(t, s1.Empty())
}

func TestSerializable_ScanValue(t *testing.T) {
	t.Parallel()

	s := &pipeline.CBORSerializable{
		Val: map[string]interface{}{
			"nested": map[string]interface{}{
				"nested": map[string]interface{}{
					"nested": map[string]interface{}{
						"bin": []byte{0x11, 0x22},
					},
				},
			},
		},
		Valid: true,
	}

	v, err := s.Value()
	require.NoError(t, err)
	require.NotEmpty(t, interface{}(v))

	s2 := &pipeline.CBORSerializable{}
	err = s2.Scan(v)
	require.NoError(t, err)
	require.False(t, s2.Empty())
	require.Equal(t, s, s2)
}
