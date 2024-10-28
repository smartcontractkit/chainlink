package utils_test

import (
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

func TestFiniteTicker(t *testing.T) {
	t.Parallel()

	var counter atomic.Int32

	onTick := func() {
		counter.Add(1)
	}

	now := time.Now()
	stop := utils.FiniteTicker(testutils.TestInterval, onTick)

	require.Eventually(t, func() bool {
		return counter.Load() >= 10
	}, testutils.WaitTimeout(t), testutils.TestInterval)

	assert.Greater(t, time.Now().Add(10*testutils.TestInterval), now)

	stop()
	last := counter.Load()
	time.Sleep(2 * testutils.TestInterval)
	assert.Equal(t, last, counter.Load())
}
