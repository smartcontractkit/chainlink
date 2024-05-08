package values

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_BoolUnwrapTo(t *testing.T) {
	tr := NewBool(true)

	var got bool
	err := tr.UnwrapTo(&got)
	require.NoError(t, err)

	assert.True(t, got)

	fa := NewBool(false)

	got = true
	err = fa.UnwrapTo(&got)
	require.NoError(t, err)

	assert.False(t, got)

	var s string
	err = fa.UnwrapTo(&s)
	require.Error(t, err)

	nilb := (*bool)(nil)
	err = fa.UnwrapTo(nilb)
	assert.ErrorContains(t, err, "cannot unwrap to nil pointer")

	var varAny any
	err = fa.UnwrapTo(&varAny)
	require.NoError(t, err)
	assert.False(t, varAny.(bool))
}
