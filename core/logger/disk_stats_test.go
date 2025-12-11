package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
)

func Test_diskSpaceAvailable(t *testing.T) {
	t.Parallel()
	tests.BelongsToCISuite(t, "unit")

	size, err := diskSpaceAvailable(".")
	assert.NoError(t, err)
	assert.NotZero(t, size)

	_, err = diskSpaceAvailable("")
	assert.Error(t, err)
}
