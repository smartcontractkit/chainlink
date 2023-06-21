package utils_test

import (
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/utils"

	"github.com/stretchr/testify/assert"
)

func TestNewRealClock(t *testing.T) {
	t.Parallel()

	clock := utils.NewRealClock()
	now := clock.Now()
	time.Sleep(testutils.TestInterval)
	interval := time.Since(now)
	assert.GreaterOrEqual(t, interval, testutils.TestInterval)
}

func TestNewFixedClock(t *testing.T) {
	t.Parallel()

	now := time.Now()
	clock := utils.NewFixedClock(now)
	time.Sleep(testutils.TestInterval)
	assert.Equal(t, now, clock.Now())
}
