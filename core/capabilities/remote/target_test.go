package remote_test

import (
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	commoncap "github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/types"
	remoteMocks "github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/types/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"
)

func TestTarget_Placeholder(t *testing.T) {
	lggr := logger.TestLogger(t)
	ctx := testutils.Context(t)
	donInfo := &types.DON{
		Members: []p2ptypes.PeerID{{}},
	}
	dispatcher := remoteMocks.NewDispatcher(t)
	dispatcher.On("Send", mock.Anything, mock.Anything).Return(nil)
	target := remote.NewRemoteTargetCaller(commoncap.CapabilityInfo{}, donInfo, dispatcher, lggr)
	require.NoError(t, target.Execute(ctx, nil, commoncap.CapabilityRequest{}))
}
