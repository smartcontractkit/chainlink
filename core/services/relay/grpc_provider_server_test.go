package relay

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-relay/pkg/types"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestProviderServer(t *testing.T) {
	r := &mockRelayer{}
	sa := NewServerAdapter(r, mockRelayerExt{})
	mp, _ := sa.NewPluginProvider(context.Background(), types.RelayArgs{ProviderType: string(types.Median)}, types.PluginArgs{})

	lggr := logger.TestLogger(t)
	_, err := NewProviderServer(mp, "unsupported-type", lggr)
	require.Error(t, err)

	ps, err := NewProviderServer(staticMedianProvider{}, types.Median, lggr)
	require.NoError(t, err)

	_, err = ps.GetConn()
	require.NoError(t, err)
}
