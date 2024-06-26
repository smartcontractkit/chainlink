package utils

import (
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeferableWriteCloser_Close(t *testing.T) {
	d := t.TempDir()
	f, err := os.Create(filepath.Join(d, "test-file"))
	require.NoError(t, err)

	wc := NewDeferableWriteCloser(f)
	wantStr := "wanted"
	_, err = io.WriteString(wc, wantStr)
	assert.NoError(t, err)
	defer func() {
		assert.NoError(t, wc.Close())
	}()

	assert.NoError(t, wc.Close())
	assert.True(t, wc.closed)
	// safe to close multiple times
	assert.NoError(t, wc.Close())

	_, err = io.WriteString(f, "after close")
	assert.ErrorIs(t, err, os.ErrClosed)

	_, err = io.WriteString(f, "write to wc after close")
	assert.ErrorIs(t, err, os.ErrClosed)

	r, err := os.ReadFile(f.Name())
	assert.NoError(t, err)
	assert.Equal(t, wantStr, string(r))
}
