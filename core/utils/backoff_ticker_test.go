package utils_test

import (
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/stretchr/testify/assert"
)

// assertTick asserts for a tick during the waitDuration
func assertTick(t *testing.T, ticker *utils.BackoffTicker, doesTick bool, waitDuration time.Duration) {
	durationDeviation := 20 * time.Millisecond
	ticked := false

	waitCh := time.After(waitDuration + durationDeviation)

	select {
	case <-ticker.Ticks():
		ticked = true
	case <-waitCh:
	}

	assert.Equal(t, doesTick, ticked)
}

func TestBackoffTicker(t *testing.T) {
	ticker := utils.NewBackoffTicker(50*time.Millisecond, 200*time.Millisecond)

	ticker.Start()

	// Increases tick time by a factor of 2
	assertTick(t, &ticker, true, 50*time.Millisecond)
	assertTick(t, &ticker, true, 100*time.Millisecond)
	assertTick(t, &ticker, true, 200*time.Millisecond)
	// It continues to tick at the max duration
	assertTick(t, &ticker, true, 200*time.Millisecond)

	ticker.Stop()
	assertTick(t, &ticker, false, 200*time.Millisecond)

	ticker.Start()
	assertTick(t, &ticker, true, 50*time.Millisecond)
	// Does not tick before the next backoff period
	assertTick(t, &ticker, false, 50*time.Millisecond)
}
