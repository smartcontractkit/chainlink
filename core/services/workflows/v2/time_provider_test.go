package v2

import (
	"testing"
	"time"

	"github.com/jonboulle/clockwork"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/workflows/dontime"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestDonTimeProvider_GetDONTime_requestTimeout(t *testing.T) {
	t.Parallel()

	fakeClock := clockwork.NewFakeClock()
	store := dontime.NewStore(dontime.DefaultRequestTimeout)
	provider := NewDonTimeProvider(store, "exec-id", 50*time.Millisecond, logger.TestLogger(t), nil, fakeClock)

	callGetDONTime := func() {
		errCh := make(chan error, 1)
		go func() {
			_, err := provider.GetDONTime()
			errCh <- err
		}()

		err := fakeClock.BlockUntilContext(t.Context(), 1)
		require.NoError(t, err)
		fakeClock.Advance(50 * time.Millisecond)

		require.NoError(t, <-errCh)
	}

	callGetDONTime()

	// A timed-out request must not block subsequent sequence numbers.
	callGetDONTime()
}
