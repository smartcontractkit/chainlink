package utils

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeferableWriterCloser_Close(t *testing.T) {

	d := t.TempDir()
	f, err := os.Create(filepath.Join(d, "test-file"))
	require.NoError(t, err)

	wc := NewDeferableWriterCloser(f)
	wantStr := "wanted"
	_, err = wc.Write([]byte(wantStr))
	assert.NoError(t, err)
	defer func() {
		assert.NoError(t, wc.Close())
	}()

	assert.NoError(t, wc.Close())
	assert.Nil(t, wc.WriteCloser)
	// safe to close multiple times
	assert.NoError(t, wc.Close())

	r, err := os.ReadFile(f.Name())
	assert.NoError(t, err)
	assert.Equal(t, wantStr, string(r))
}
