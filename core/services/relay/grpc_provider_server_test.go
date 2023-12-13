package relay

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/types"

	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestProviderServer(t *testing.T) {
	r := &mockRelayer{}
	sa := NewServerAdapter(r, mockRelayerExt{})
	mp, _ := sa.NewPluginProvider(testutils.Context(t), types.RelayArgs{ProviderType: string(types.Median)}, types.PluginArgs{})

	lggr := logger.TestLogger(t)
	_, err := NewProviderServer(mp, "unsupported-type", lggr)
	require.ErrorContains(t, err, "unsupported-type")

	ps, err := NewProviderServer(staticMedianProvider{}, types.Median, lggr)
	require.NoError(t, err)

	_, err = ps.GetConn()
	require.NoError(t, err)
}
