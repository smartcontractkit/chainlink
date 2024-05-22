package values

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_StringUnwrapTo(t *testing.T) {
	expected := "hello"
	v := NewString(expected)

	var got string
	err := v.UnwrapTo(&got)
	require.NoError(t, err)

	assert.Equal(t, expected, got)

	gotn := (*string)(nil)
	err = v.UnwrapTo(gotn)
	assert.ErrorContains(t, err, "cannot unwrap to nil pointer")

	var varAny any
	err = v.UnwrapTo(&varAny)
	require.NoError(t, err)
	assert.Equal(t, expected, varAny)
}
