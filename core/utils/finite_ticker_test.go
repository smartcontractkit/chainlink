package utils_test

import (
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/stretchr/testify/assert"

	"go.uber.org/atomic"
)

func TestFiniteTicker(t *testing.T) {
	t.Parallel()

	var counter atomic.Int32

	onTick := func() {
		counter.Inc()
	}

	now := time.Now()
	stop := utils.FiniteTicker(testutils.TestInterval, onTick)

	assert.Eventually(t, func() bool {
		return counter.Load() >= 10
	}, testutils.WaitTimeout(t), testutils.TestInterval)

	assert.Greater(t, time.Now().Add(10*testutils.TestInterval), now)

	stop()
	last := counter.Load()
	time.Sleep(2 * testutils.TestInterval)
	assert.Equal(t, last, counter.Load())
}
