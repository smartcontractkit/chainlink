package remote

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	commoncap "github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	remotetypes "github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/types"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/types/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"
)

type stubTrigger struct{}

func (stubTrigger) Info(context.Context) (commoncap.CapabilityInfo, error) {
	return commoncap.CapabilityInfo{ID: "stub@1", CapabilityType: commoncap.CapabilityTypeTrigger, Description: "stub"}, nil
}

func (stubTrigger) RegisterTrigger(context.Context, commoncap.TriggerRegistrationRequest) (<-chan commoncap.TriggerResponse, error) {
	return make(chan commoncap.TriggerResponse), nil
}

func (stubTrigger) UnregisterTrigger(context.Context, commoncap.TriggerRegistrationRequest) error {
	return nil
}

func (stubTrigger) AckEvent(context.Context, string, string, string) error {
	return nil
}

func regCachePeer(i int) p2ptypes.PeerID {
	var p p2ptypes.PeerID
	p[30] = byte(i >> 8)
	p[31] = byte(i)
	return p
}

// A workflow-DON member that can never reach the 2F+1 registration quorum must not be
// able to grow the registration staging cache without bound. Pre-quorum registrations
// never become active registrations, so the unregister path cannot remove them;
// cacheCleanupLoop must reap stale staging entries (older than RegistrationExpiry).
func TestTriggerPublisher_RegistrationStagingCacheIsBounded(t *testing.T) {
	ctx := testutils.Context(t)
	lggr := logger.Test(t)

	const capabilityDONID = uint32(1)
	const workflowDONID = uint32(2)

	capDON := commoncap.DON{ID: capabilityDONID, Members: []p2ptypes.PeerID{regCachePeer(1000)}, F: 0}
	wfMembers := []p2ptypes.PeerID{regCachePeer(1), regCachePeer(2), regCachePeer(3)}
	workflowDON := commoncap.DON{ID: workflowDONID, Members: wfMembers, F: 1} // quorum 2F+1 = 3
	attacker := wfMembers[0]

	dispatcher := mocks.NewDispatcher(t)
	dispatcher.On("Send", mock.Anything, mock.Anything).Return(nil).Maybe()

	cfg := &commoncap.RemoteTriggerConfig{
		RegistrationRefresh:     100 * time.Millisecond,
		RegistrationExpiry:      500 * time.Millisecond,
		MinResponsesToAggregate: 3,
		MessageExpiry:           200 * time.Millisecond, // cacheCleanupLoop tick interval
		MaxBatchSize:            1,
		BatchCollectionPeriod:   time.Second,
	}

	publisher := NewTriggerPublisher("cap_id@1", "", dispatcher, lggr)
	require.NoError(t, publisher.SetConfig(cfg, stubTrigger{}, capDON, map[uint32]commoncap.DON{workflowDON.ID: workflowDON}))
	require.NoError(t, publisher.Start(ctx))
	t.Cleanup(func() { _ = publisher.Close() })

	const n = 1000
	for i := 0; i < n; i++ {
		wfID := fmt.Sprintf("%064x", i) // distinct, valid 32-byte-hex workflow id => distinct staging key
		req := commoncap.TriggerRegistrationRequest{Metadata: commoncap.RequestMetadata{WorkflowID: wfID}}
		marshaled, err := pb.MarshalTriggerRegistrationRequest(req)
		require.NoError(t, err)
		publisher.Receive(ctx, &remotetypes.MessageBody{
			Sender:      attacker[:],
			Method:      remotetypes.MethodRegisterTrigger,
			CallerDonId: workflowDONID,
			Payload:     marshaled,
		})
	}

	cacheLen := func() int {
		publisher.mu.Lock()
		defer publisher.mu.Unlock()
		return publisher.messageCache.Len()
	}

	require.GreaterOrEqual(t, cacheLen(), n*9/10, "below-quorum registrations should be staged in messageCache")

	require.Eventually(t, func() bool { return cacheLen() == 0 }, 4*time.Second, 100*time.Millisecond,
		"stale below-quorum registration staging entries must be reaped by cacheCleanupLoop")
}
