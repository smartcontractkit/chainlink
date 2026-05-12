package utils

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestBillingService(t *testing.T) {
	t.Parallel()

	t.Run("ephemeral port", func(t *testing.T) {
		t.Parallel()

		lggr := logger.TestLogger(t)
		svc := NewBillingServiceOnPort(lggr, 0)
		require.NoError(t, svc.Start(t.Context()))
		t.Cleanup(func() { _ = svc.Close() })

		addr := svc.Addr()
		assert.NotEmpty(t, addr)
		assert.NotEqual(t, fmt.Sprintf("localhost:%d", defaultBillingPort), addr)
	})

	t.Run("specific port", func(t *testing.T) {
		t.Parallel()

		lggr := logger.TestLogger(t)
		svc := NewBillingService(lggr)
		require.NoError(t, svc.Start(t.Context()))
		t.Cleanup(func() { _ = svc.Close() })

		assert.Equal(t, fmt.Sprintf("127.0.0.1:%d", defaultBillingPort), svc.Addr())
	})

	t.Run("parallel instances do not conflict", func(t *testing.T) {
		t.Parallel()

		lggr := logger.TestLogger(t)
		svc1 := NewBillingServiceOnPort(lggr, 0)
		svc2 := NewBillingServiceOnPort(lggr, 0)

		require.NoError(t, svc1.Start(t.Context()))
		t.Cleanup(func() { _ = svc1.Close() })

		require.NoError(t, svc2.Start(t.Context()))
		t.Cleanup(func() { _ = svc2.Close() })

		assert.NotEqual(t, svc1.Addr(), svc2.Addr())
	})
}
