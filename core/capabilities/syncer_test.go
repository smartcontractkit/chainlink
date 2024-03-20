package capabilities_test

import (
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	commonMocks "github.com/smartcontractkit/chainlink-common/pkg/types/mocks"
	coreCapabilities "github.com/smartcontractkit/chainlink/v2/core/capabilities"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/p2p/types/mocks"
)

func TestSyncer_CleanStartClose(t *testing.T) {
	lggr := logger.TestLogger(t)
	ctx := testutils.Context(t)
	peer := mocks.NewPeer(t)
	peer.On("UpdateConnections", mock.Anything).Return(nil)
	wrapper := mocks.NewPeerWrapper(t)
	wrapper.On("GetPeer").Return(peer)
	registry := commonMocks.NewCapabilitiesRegistry(t)

	syncer := coreCapabilities.NewRegistrySyncer(wrapper, registry, lggr)
	require.NoError(t, syncer.Start(ctx))
	require.NoError(t, syncer.Close())
}
