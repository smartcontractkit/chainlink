package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestBillingService(t *testing.T) {
	t.Parallel()

	t.Run("parallel instances do not conflict", func(t *testing.T) {
		t.Parallel()

		lggr := logger.TestLogger(t)
		svc1 := NewBillingServiceOnPort(lggr, 0)
		svc2 := NewBillingServiceOnPort(lggr, 0)

		require.NoError(t, svc1.Start(t.Context()))
		t.Cleanup(func() { require.NoError(t, svc1.Close()) })

		require.NoError(t, svc2.Start(t.Context()))
		t.Cleanup(func() { require.NoError(t, svc2.Close()) })

		assert.NotEqual(t, svc1.Addr(), svc2.Addr())
	})
}
