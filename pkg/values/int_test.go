package values

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_IntUnwrapTo(t *testing.T) {
	expected := int64(100)
	v := NewInt64(expected)

	var got int
	err := v.UnwrapTo(&got)
	require.NoError(t, err)

	assert.Equal(t, expected, int64(got))

	var gotInt64 int64
	err = v.UnwrapTo(&gotInt64)
	require.NoError(t, err)

	assert.Equal(t, expected, gotInt64)
}
