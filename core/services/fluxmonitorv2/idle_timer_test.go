package fluxmonitorv2_test

import (
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/core/services/fluxmonitorv2"
	"github.com/stretchr/testify/assert"
)

func TestIdleTimer_Getters(t *testing.T) {
	t.Parallel()

	ticker := fluxmonitorv2.NewIdleTimer(time.Second, false)

	t.Run("Period", func(t *testing.T) {
		assert.Equal(t, time.Second, ticker.Period())
	})

	t.Run("IsEnabled", func(t *testing.T) {
		assert.Equal(t, true, ticker.IsEnabled())
	})

	t.Run("IsDisabled", func(t *testing.T) {
		assert.Equal(t, false, ticker.IsDisabled())
	})
}

// TODO - Test the ticker functions
func TestIdleTimer_Ticker(t *testing.T) {
	t.Skip()
}
