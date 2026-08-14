package wasmtest

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetTestBinary(t *testing.T) {
	t.Parallel()

	binary := GetTestBinary(t, "core/services/job/testdata/wasm", false)
	require.NotEmpty(t, binary)
	assert.Equal(t, []byte("\x00asm"), binary[:4])

	compressedBinary := GetTestBinary(t, "core/services/job/testdata/wasm", true)
	require.NotEmpty(t, compressedBinary)
}
