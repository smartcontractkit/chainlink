package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_diskSpaceAvailable(t *testing.T) {
	t.Parallel()

	size, err := diskSpaceAvailable(".")
	require.NoError(t, err)
	assert.NotZero(t, size)

	_, err = diskSpaceAvailable("")
	assert.Error(t, err)
}
