package utils_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/smartcontractkit/chainlink/core/utils"
)

func TestDiskStatsProvider_AvailableSpace(t *testing.T) {
	t.Parallel()

	provider := utils.NewDiskStatsProvider()

	size, err := provider.AvailableSpace(".")
	assert.NoError(t, err)
	assert.NotZero(t, size)

	_, err = provider.AvailableSpace("")
	assert.Error(t, err)
}
