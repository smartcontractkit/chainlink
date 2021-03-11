package fluxmonitorv2_test

import (
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/core/services/fluxmonitorv2"
	"github.com/stretchr/testify/assert"
)

func TestPollTicker_Getters(t *testing.T) {
	t.Parallel()

	ticker := fluxmonitorv2.NewPollTicker(time.Second, false)

	t.Run("Interval", func(t *testing.T) {
		assert.Equal(t, time.Second, ticker.Interval())
	})

	t.Run("IsEnabled", func(t *testing.T) {
		assert.Equal(t, true, ticker.IsEnabled())
	})

	t.Run("IsDisabled", func(t *testing.T) {
		assert.Equal(t, false, ticker.IsDisabled())
	})
}

// TODO - Test the ticker functions
func TestPollTimer_Ticker(t *testing.T) {
}
