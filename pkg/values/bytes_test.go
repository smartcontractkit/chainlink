package values

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_BytesUnwrapTo(t *testing.T) {
	hs := []byte("hello")
	tr := NewBytes(hs)

	got := []byte{}
	err := tr.UnwrapTo(&got)
	require.NoError(t, err)

	assert.Equal(t, hs, got)

	var s string
	err = tr.UnwrapTo(&s)
	require.Error(t, err)

	gotB := (*[]byte)(nil)
	err = tr.UnwrapTo(gotB)
	assert.ErrorContains(t, err, "cannot unwrap to nil pointer")
}
