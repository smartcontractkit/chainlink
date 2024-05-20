package remote_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	commoncap "github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink-common/pkg/values"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote"
	remotetypes "github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/types"
	remoteMocks "github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/types/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"

	"github.com/stretchr/testify/mock"
)

const (
	executeValue1 = "triggerEvent1"
)

func Test_TargetCallerExecute(t *testing.T) {

	lggr := logger.TestLogger(t)
	ctx := testutils.Context(t)
	capInfo := commoncap.CapabilityInfo{
		ID:             "cap_id",
		CapabilityType: commoncap.CapabilityTypeTarget,
		Description:    "Remote Target",
		Version:        "0.0.1",
	}
	p1 := p2ptypes.PeerID{}
	require.NoError(t, p1.UnmarshalText([]byte(PeerID1)))
	p2 := p2ptypes.PeerID{}
	require.NoError(t, p2.UnmarshalText([]byte(PeerID2)))
	capDonInfo := commoncap.DON{
		ID:      "capability-don",
		Members: []p2ptypes.PeerID{p1},
		F:       0,
	}

	/*
		workflowDonInfo := commoncap.DON{
			ID:      "workflow-don",
			Members: []p2ptypes.PeerID{p2},
			F:       0,
		}*/

	dispatcher := remoteMocks.NewDispatcher(t)

	awaitExecuteCh := make(chan struct{})
	dispatcher.On("Send", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		select {
		case awaitExecuteCh <- struct{}{}:
		default:
		}
	})

	caller := remote.NewRemoteTargetCaller(lggr, capInfo, capDonInfo, dispatcher)

	go func() {
		<-awaitExecuteCh

		executeValue, err := values.Wrap(executeValue1)
		require.NoError(t, err)
		capResponse := commoncap.CapabilityResponse{
			Value: executeValue,
			Err:   nil,
		}
		marshaled, err := pb.MarshalCapabilityResponse(capResponse)
		require.NoError(t, err)
		executeResponse := &remotetypes.MessageBody{
			Sender:  p1[:],
			Method:  remotetypes.MethodExecute,
			Payload: marshaled,
		}

		caller.Receive(executeResponse)
	}()

	resultCh, err := caller.Execute(ctx,
		commoncap.CapabilityRequest{
			Metadata: commoncap.RequestMetadata{
				WorkflowID: WorkflowID1,
			},
		})

	require.NoError(t, err)

	response := <-resultCh

	responseValue, err := response.Value.Unwrap()
	assert.Equal(t, executeValue1, responseValue.(string))

}
