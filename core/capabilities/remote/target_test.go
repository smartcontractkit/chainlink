package remote_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	commoncap "github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote"
	remoteMocks "github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/types/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"
)

func TestTarget_Placeholder(t *testing.T) {
	lggr := logger.TestLogger(t)
	ctx := testutils.Context(t)
	donInfo := &capabilities.DON{
		Members: []p2ptypes.PeerID{{}},
	}
	dispatcher := remoteMocks.NewDispatcher(t)
	dispatcher.On("Send", mock.Anything, mock.Anything).Return(nil)
	target := remote.NewRemoteTargetCaller(commoncap.CapabilityInfo{}, donInfo, dispatcher, lggr)

	_, err := target.Execute(ctx, commoncap.CapabilityRequest{})
	assert.NoError(t, err)
}
