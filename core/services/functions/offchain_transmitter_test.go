package functions_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/services/functions"
)

func TestOffchainTransmitter(t *testing.T) {
	t.Parallel()

	transmitter := functions.NewOffchainTransmitter(1)
	ch := transmitter.ReportChannel()
	report := &functions.OffchainResponse{RequestId: []byte("testID")}
	ctx := testutils.Context(t)

	require.NoError(t, transmitter.TransmitReport(ctx, report))
	require.Equal(t, report, <-ch)

	require.NoError(t, transmitter.TransmitReport(ctx, report))

	ctxTimeout, cancel := context.WithTimeout(ctx, time.Millisecond*20)
	defer cancel()
	// should not freeze
	require.Error(t, transmitter.TransmitReport(ctxTimeout, report))
}
