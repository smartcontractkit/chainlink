package v2

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/workflows/dontime"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestDonTimeProvider_GetDONTime_requestTimeout(t *testing.T) {
	t.Parallel()

	store := dontime.NewStore(dontime.DefaultRequestTimeout)
	provider := NewDonTimeProvider(store, "exec-id", 50*time.Millisecond, logger.TestLogger(t), nil)

	start := time.Now()
	_, err := provider.GetDONTime()
	elapsed := time.Since(start)

	require.NoError(t, err)
	require.GreaterOrEqual(t, elapsed, 50*time.Millisecond)
	require.Less(t, elapsed, 200*time.Millisecond)

	// A timed-out request must not block subsequent sequence numbers.
	start = time.Now()
	_, err = provider.GetDONTime()
	elapsed = time.Since(start)
	require.NoError(t, err)
	require.GreaterOrEqual(t, elapsed, 50*time.Millisecond)
	require.Less(t, elapsed, 200*time.Millisecond)
}
