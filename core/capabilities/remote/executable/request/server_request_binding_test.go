package request_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	commoncap "github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
	"github.com/smartcontractkit/chainlink-protos/cre/go/values"

	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/executable/request"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/types"
	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"
)

// Test_ServerRequest_WorkflowDONBinding verifies that, when the binding gate is
// open, the server rejects requests whose Metadata.WorkflowDonID does not match
// the authenticated calling DON (callingDon.ID), and otherwise executes.
func Test_ServerRequest_WorkflowDONBinding(t *testing.T) {
	t.Parallel()

	lggr := logger.Test(t)
	capability := TestCapability{}
	capabilityPeerID := NewP2PPeerID(t)
	workflowPeers := []p2ptypes.PeerID{NewP2PPeerID(t), NewP2PPeerID(t)}

	const callingDonID = uint32(1)
	callingDon := commoncap.DON{Members: workflowPeers, ID: callingDonID, F: 1}

	executeInputs, err := values.NewMap(map[string]any{"executeValue1": "aValue1"})
	require.NoError(t, err)

	buildPayload := func(workflowDonID uint32) []byte {
		raw, merr := pb.MarshalCapabilityRequest(commoncap.CapabilityRequest{
			Metadata: commoncap.RequestMetadata{
				WorkflowID:          "workflowID",
				WorkflowExecutionID: "workflowExecutionID",
				WorkflowDonID:       workflowDonID,
			},
			Inputs: executeInputs,
		})
		require.NoError(t, merr)
		return raw
	}

	send := func(req serverRequest, peer p2ptypes.PeerID, payload []byte) error {
		return req.OnMessage(context.Background(), &types.MessageBody{
			Sender:          peer[:],
			Receiver:        capabilityPeerID[:],
			MessageId:       []byte("workflowID" + "workflowExecutionID"),
			CapabilityId:    "capabilityID",
			CapabilityDonId: 2,
			CallerDonId:     callingDonID,
			Method:          types.MethodExecute,
			Payload:         payload,
		})
	}

	// Drive an F+1 (=2) quorum of identical requests from distinct DON peers,
	// which is what triggers execution on the server.
	newReqAndSendQuorum := func(t *testing.T, gate limits.GateLimiter, workflowDonID uint32) *testDispatcher {
		t.Helper()
		d := &testDispatcher{}
		req, rerr := request.NewServerRequest(capability, types.MethodExecute, "capabilityID", 2,
			capabilityPeerID, callingDon, "requestMessageID", d, 10*time.Minute, "", gate, lggr)
		require.NoError(t, rerr)
		payload := buildPayload(workflowDonID)
		require.NoError(t, send(req, workflowPeers[0], payload))
		require.NoError(t, send(req, workflowPeers[1], payload))
		require.Len(t, d.msgs, 2)
		return d
	}

	t.Run("gate enabled, WorkflowDonID mismatch is rejected", func(t *testing.T) {
		t.Parallel()
		d := newReqAndSendQuorum(t, limits.NewGateLimiter(true), callingDonID+99)
		require.Equal(t, types.Error_INTERNAL_ERROR, d.msgs[0].Error)
		require.Contains(t, d.msgs[0].ErrorMsg, "does not match calling DON")
	})

	t.Run("gate enabled, WorkflowDonID match executes", func(t *testing.T) {
		t.Parallel()
		d := newReqAndSendQuorum(t, limits.NewGateLimiter(true), callingDonID)
		require.Equal(t, types.Error_OK, d.msgs[0].Error)
	})

	t.Run("gate disabled, WorkflowDonID mismatch still executes", func(t *testing.T) {
		t.Parallel()
		d := newReqAndSendQuorum(t, limits.NewGateLimiter(false), callingDonID+99)
		require.Equal(t, types.Error_OK, d.msgs[0].Error)
	})

	t.Run("nil gate is rejected at construction", func(t *testing.T) {
		t.Parallel()
		_, rerr := request.NewServerRequest(capability, types.MethodExecute, "capabilityID", 2,
			capabilityPeerID, callingDon, "requestMessageID", &testDispatcher{}, 10*time.Minute, "", nil, lggr)
		require.ErrorContains(t, rerr, "workflowDONBindingGate must not be nil")
	})
}
