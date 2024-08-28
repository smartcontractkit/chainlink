package remote_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	commoncap "github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote"
	remotetypes "github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/types"
	remoteMocks "github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/types/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"
)

func TestTriggerPublisher_Register(t *testing.T) {
	lggr := logger.TestLogger(t)
	ctx := testutils.Context(t)
	capInfo := commoncap.CapabilityInfo{
		ID:             "cap_id@1",
		CapabilityType: commoncap.CapabilityTypeTrigger,
		Description:    "Remote Trigger",
	}
	p1 := p2ptypes.PeerID{}
	require.NoError(t, p1.UnmarshalText([]byte(peerID1)))
	p2 := p2ptypes.PeerID{}
	require.NoError(t, p2.UnmarshalText([]byte(peerID2)))
	capDonInfo := commoncap.DON{
		ID:      1,
		Members: []p2ptypes.PeerID{p1},
		F:       0,
	}
	workflowDonInfo := commoncap.DON{
		ID:      2,
		Members: []p2ptypes.PeerID{p2},
		F:       0,
	}

	dispatcher := remoteMocks.NewDispatcher(t)
	config := &commoncap.RemoteTriggerConfig{
		RegistrationRefresh:     100 * time.Millisecond,
		RegistrationExpiry:      100 * time.Second,
		MinResponsesToAggregate: 1,
		MessageExpiry:           100 * time.Second,
	}
	workflowDONs := map[uint32]commoncap.DON{
		workflowDonInfo.ID: workflowDonInfo,
	}
	underlying := &testTrigger{
		info:            capInfo,
		registrationsCh: make(chan commoncap.TriggerRegistrationRequest, 2),
	}
	publisher := remote.NewTriggerPublisher(config, underlying, capInfo, capDonInfo, workflowDONs, dispatcher, lggr)
	require.NoError(t, publisher.Start(ctx))

	// trigger registration event
	triggerRequest := commoncap.TriggerRegistrationRequest{
		Metadata: commoncap.RequestMetadata{
			WorkflowID: workflowID1,
		},
	}
	marshaled, err := pb.MarshalTriggerRegistrationRequest(triggerRequest)
	require.NoError(t, err)
	regEvent := &remotetypes.MessageBody{
		Sender:      p1[:],
		Method:      remotetypes.MethodRegisterTrigger,
		CallerDonId: workflowDonInfo.ID,
		Payload:     marshaled,
	}
	publisher.Receive(ctx, regEvent)
	// node p1 is not a member of the workflow DON so registration shoudn't happen
	require.Empty(t, underlying.registrationsCh)

	regEvent.Sender = p2[:]
	publisher.Receive(ctx, regEvent)
	require.NotEmpty(t, underlying.registrationsCh)
	forwarded := <-underlying.registrationsCh
	require.Equal(t, triggerRequest.Metadata.WorkflowID, forwarded.Metadata.WorkflowID)

	require.NoError(t, publisher.Close())
}

type testTrigger struct {
	info            commoncap.CapabilityInfo
	registrationsCh chan commoncap.TriggerRegistrationRequest
}

func (t *testTrigger) Info(_ context.Context) (commoncap.CapabilityInfo, error) {
	return t.info, nil
}

func (t *testTrigger) RegisterTrigger(_ context.Context, request commoncap.TriggerRegistrationRequest) (<-chan commoncap.TriggerResponse, error) {
	t.registrationsCh <- request
	return nil, nil
}

func (t *testTrigger) UnregisterTrigger(_ context.Context, request commoncap.TriggerRegistrationRequest) error {
	return nil
}
