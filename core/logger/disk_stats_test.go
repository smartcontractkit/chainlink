package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_diskSpaceAvailable(t *testing.T) {
	t.Parallel()

	size, err := diskSpaceAvailable(".")
	assert.NoError(t, err)
	assert.NotZero(t, size)

	_, err = diskSpaceAvailable("")
	assert.Error(t, err)
}
