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

	// The upper bound is loose on purpose: the call is expected to return
	// shortly after the 50ms request timeout, but under t.Parallel() on a
	// loaded CI runner the goroutine may not be scheduled to resume for
	// several hundred milliseconds after the timer fires. 1s is still tight
	// enough to catch a real stall (e.g. the timer never firing, which would
	// block for the store's 10m expiry).
	start := time.Now()
	_, err := provider.GetDONTime()
	elapsed := time.Since(start)

	require.NoError(t, err)
	require.GreaterOrEqual(t, elapsed, 50*time.Millisecond)
	require.Less(t, elapsed, time.Second)

	// A timed-out request must not block subsequent sequence numbers.
	start = time.Now()
	_, err = provider.GetDONTime()
	elapsed = time.Since(start)
	require.NoError(t, err)
	require.GreaterOrEqual(t, elapsed, 50*time.Millisecond)
	require.Less(t, elapsed, time.Second)
}
