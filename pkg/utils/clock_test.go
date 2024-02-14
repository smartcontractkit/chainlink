package utils_test

import (
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/utils"

	"github.com/stretchr/testify/assert"
)

const (
	TestInterval = 100 * time.Millisecond
)

func TestNewRealClock(t *testing.T) {
	t.Parallel()

	clock := utils.NewRealClock()
	now := clock.Now()
	time.Sleep(TestInterval)
	interval := time.Since(now)
	assert.GreaterOrEqual(t, interval, TestInterval)
}

func TestNewFixedClock(t *testing.T) {
	t.Parallel()

	now := time.Now()
	clock := utils.NewFixedClock(now)
	time.Sleep(TestInterval)
	assert.Equal(t, now, clock.Now())
}
