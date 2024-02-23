package solana_test

import (
	"testing"

	"github.com/test-go/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/solana"
)

func TestSolanaChainReaderService_ServiceCtx(t *testing.T) {
	t.Parallel()

	ctx := testutils.Context(t)
	svc, err := solana.NewChainReaderService(logger.TestLogger(t))

	require.NoError(t, err)
	require.NotNil(t, svc)

	require.Error(t, svc.Ready())
	require.Len(t, svc.HealthReport(), 1)
	require.Contains(t, svc.HealthReport(), solana.ServiceName)
	require.Error(t, svc.HealthReport()[solana.ServiceName])

	require.NoError(t, svc.Start(ctx))
	require.NoError(t, svc.Ready())
	require.Equal(t, map[string]error{solana.ServiceName: nil}, svc.HealthReport())

	require.Error(t, svc.Start(ctx))

	require.NoError(t, svc.Close())
	require.Error(t, svc.Ready())
	require.Error(t, svc.Close())
}
