package registration_test

import (
	"context"
	"errors"
	"time"

	//"errors"
	"testing"

	"github.com/stretchr/testify/require"

	commoncap "github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	caperrors "github.com/smartcontractkit/chainlink-common/pkg/capabilities/errors"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/trigger/registration"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/types"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	p2ptypes "github.com/smartcontractkit/libocr/ragep2p/types"
)

type testTriggerCapability struct {
	err error
}

var _ commoncap.TriggerCapability = (*testTriggerCapability)(nil)

func (t *testTriggerCapability) Info(ctx context.Context) (commoncap.CapabilityInfo, error) {
	return commoncap.CapabilityInfo{
		ID: "test-capability",
	}, nil
}

func (t *testTriggerCapability) RegisterTrigger(ctx context.Context, req commoncap.TriggerRegistrationRequest) (<-chan commoncap.TriggerResponse, error) {
	ch := make(chan commoncap.TriggerResponse, 1)
	return ch, t.err
}

func (t *testTriggerCapability) UnregisterTrigger(ctx context.Context, req commoncap.TriggerRegistrationRequest) error {
	return nil
}

type sendCall struct {
	peerID p2ptypes.PeerID
	msg    *types.MessageBody
}

func TestMultipleAddedRequests_SuccessfulTriggerRegistration(t *testing.T) {
	// Test that when multiple requests are added to a PublisherRegistration,
	// they all receive the same registration status update when the trigger is registered.

	// Setup
	lggr := logger.TestLogger(t)
	triggerID := "test-trigger"
	workflowID := "test-workflow"
	capabilityID := "test-capability"

	var sendCalls []sendCall
	pr := registration.NewPublisherRegistration(lggr, &testTriggerCapability{err: nil}, triggerID, workflowID, capabilityID, func(peerID p2ptypes.PeerID, msgBody *types.MessageBody) error {
		sendCalls = append(sendCalls, sendCall{peerID, msgBody})
		return nil
	})

	peers := make([]p2ptypes.PeerID, 8)
	require.NoError(t, peers[0].UnmarshalText([]byte(peerID1)))
	require.NoError(t, peers[1].UnmarshalText([]byte(peerID2)))
	require.NoError(t, peers[2].UnmarshalText([]byte(peerID3)))

	request := commoncap.TriggerRegistrationRequest{
		TriggerID: triggerID,
	}

	rawRequest, err := pb.MarshalTriggerRegistrationRequest(request)
	require.NoError(t, err)

	// Add multiple requests
	peerID1 := peers[0]
	peerID2 := peers[1]

	callerDon := commoncap.DON{
		ID: 1,
	}

	registrationExpiry := 10 * time.Minute

	pr.AddRegistrationRequest(context.Background(), peerID1, rawRequest, callerDon, registrationExpiry)
	pr.AddRegistrationRequest(context.Background(), peerID2, rawRequest, callerDon, registrationExpiry)

	// Confirm a trigger response message was sent to each peer
	require.Len(t, sendCalls, 2)

	foundPeer1 := false
	foundPeer2 := false
	for _, call := range sendCalls {
		require.Equal(t, types.TriggerRegistrationStatus, call.msg.Method)
		meta := call.msg.GetTriggerRegistrationMetadata()
		require.Equal(t, types.RegistrationStatus_REGISTERED, meta.Status)
		require.Equal(t, triggerID, meta.TriggerId)
		if call.peerID == peerID1 {
			foundPeer1 = true
		} else if call.peerID == peerID2 {
			foundPeer2 = true
		}
	}

	require.True(t, foundPeer1, "did not find message sent to peerID1")
	require.True(t, foundPeer2, "did not find message sent to peerID2")

	// Send an additional registration request after trigger registration has already completed, expect to get immediate response to the new peer
	peerID3 := peers[2]
	pr.AddRegistrationRequest(context.Background(), peerID3, rawRequest, callerDon, registrationExpiry)

	// Confirm a trigger response message was sent to the new peer
	require.Len(t, sendCalls, 3)

	foundPeer3 := false
	for _, call := range sendCalls {
		if call.peerID == peerID3 {
			foundPeer3 = true
			require.Equal(t, types.TriggerRegistrationStatus, call.msg.Method)
			meta := call.msg.GetTriggerRegistrationMetadata()
			require.Equal(t, types.RegistrationStatus_REGISTERED, meta.Status)
			require.Equal(t, triggerID, meta.TriggerId)
		}
	}

	require.True(t, foundPeer3, "did not find message sent to peerID3")
}

func TestMultipleAddedRequests_TriggerRegistrationError(t *testing.T) {
	// Test that when multiple requests are added to a PublisherRegistration,
	// they all receive the same error response when the trigger registration fails.

	// Setup
	lggr := logger.TestLogger(t)
	triggerID := "test-trigger"
	workflowID := "test-workflow"
	capabilityID := "test-capability"

	var sendCalls []sendCall

	pr := registration.NewPublisherRegistration(lggr, &testTriggerCapability{err: errors.New("its borken")}, triggerID, workflowID, capabilityID, func(peerID p2ptypes.PeerID, msgBody *types.MessageBody) error {
		sendCalls = append(sendCalls, sendCall{peerID, msgBody})
		return nil
	})

	peers := make([]p2ptypes.PeerID, 4)
	require.NoError(t, peers[0].UnmarshalText([]byte(peerID1)))
	require.NoError(t, peers[1].UnmarshalText([]byte(peerID2)))

	request := commoncap.TriggerRegistrationRequest{
		TriggerID: triggerID,
	}

	rawRequest, err := pb.MarshalTriggerRegistrationRequest(request)
	require.NoError(t, err)

	// Add multiple requests
	peerID1 := peers[0]
	peerID2 := peers[1]

	callerDon := commoncap.DON{
		ID: 1,
	}

	registrationExpiry := 10 * time.Minute

	pr.AddRegistrationRequest(context.Background(), peerID1, rawRequest, callerDon, registrationExpiry)
	pr.AddRegistrationRequest(context.Background(), peerID2, rawRequest, callerDon, registrationExpiry)

	// Confirm a trigger response message was sent to each peer
	require.Len(t, sendCalls, 2)

	foundPeer1 := false
	foundPeer2 := false
	for _, call := range sendCalls {
		require.Equal(t, types.TriggerRegistrationStatus, call.msg.Method)
		meta := call.msg.GetTriggerRegistrationMetadata()
		require.Equal(t, types.RegistrationStatus_REGISTRATION_ERROR, meta.Status)
		require.Equal(t, triggerID, meta.TriggerId)
		caperror := caperrors.DeserializeErrorFromString(meta.RegistrationError)
		require.Equal(t, caperrors.Unknown, caperror.Code())
		require.Equal(t, caperrors.OriginSystem, caperror.Origin())
		require.Equal(t, caperrors.VisibilityPublic, caperror.Visibility())
		require.Equal(t, "[2]Unknown: its borken", caperror.Error())
		if call.peerID == peerID1 {
			foundPeer1 = true
		} else if call.peerID == peerID2 {
			foundPeer2 = true
		}
	}

	require.True(t, foundPeer1, "did not find message sent to peerID1")
	require.True(t, foundPeer2, "did not find message sent to peerID2")

	// Send an additional registration request after trigger registration has already completed with error, expect to get immediate error response to the new peer
	peerID3 := peers[2]

	pr.AddRegistrationRequest(context.Background(), peerID3, rawRequest, callerDon, registrationExpiry)

	// Confirm a trigger response message was sent to the new peer
	require.Len(t, sendCalls, 3)

	foundPeer3 := false
	for _, call := range sendCalls {
		if call.peerID == peerID3 {
			foundPeer3 = true
			require.Equal(t, types.TriggerRegistrationStatus, call.msg.Method)
			meta := call.msg.GetTriggerRegistrationMetadata()
			require.Equal(t, types.RegistrationStatus_REGISTRATION_ERROR, meta.Status)
			require.Equal(t, triggerID, meta.TriggerId)
			caperror := caperrors.DeserializeErrorFromString(meta.RegistrationError)
			require.Equal(t, caperrors.Unknown, caperror.Code())
			require.Equal(t, caperrors.OriginSystem, caperror.Origin())
			require.Equal(t, caperrors.VisibilityPublic, caperror.Visibility())
			require.Equal(t, "[2]Unknown: its borken", caperror.Error())

			require.Equal(t, triggerID, meta.TriggerId)
		}
	}

	require.True(t, foundPeer3, "did not find message sent to peerID3")
}

func TestMultipleAddedRequests_TriggerRegistrationCapabilityUserError(t *testing.T) {
	// Test that when multiple requests are added to a PublisherRegistration,
	// they all receive the same error response when the trigger registration fails.

	// Setup
	lggr := logger.TestLogger(t)
	triggerID := "test-trigger"
	workflowID := "test-workflow"
	capabilityID := "test-capability"

	var sendCalls []sendCall

	pr := registration.NewPublisherRegistration(lggr, &testTriggerCapability{err: caperrors.NewPublicUserError(errors.New("its borken"), caperrors.InvalidArgument)}, triggerID, workflowID, capabilityID, func(peerID p2ptypes.PeerID, msgBody *types.MessageBody) error {
		sendCalls = append(sendCalls, sendCall{peerID, msgBody})
		return nil
	})

	peers := make([]p2ptypes.PeerID, 4)
	require.NoError(t, peers[0].UnmarshalText([]byte(peerID1)))
	require.NoError(t, peers[1].UnmarshalText([]byte(peerID2)))

	request := commoncap.TriggerRegistrationRequest{
		TriggerID: triggerID,
	}

	rawRequest, err := pb.MarshalTriggerRegistrationRequest(request)
	require.NoError(t, err)

	// Add multiple requests
	peerID1 := peers[0]
	peerID2 := peers[1]

	callerDon := commoncap.DON{
		ID: 1,
	}

	registrationExpiry := 10 * time.Minute

	pr.AddRegistrationRequest(context.Background(), peerID1, rawRequest, callerDon, registrationExpiry)
	pr.AddRegistrationRequest(context.Background(), peerID2, rawRequest, callerDon, registrationExpiry)

	// Confirm a trigger response message was sent to each peer
	require.Len(t, sendCalls, 2)

	foundPeer1 := false
	foundPeer2 := false
	for _, call := range sendCalls {
		require.Equal(t, types.TriggerRegistrationStatus, call.msg.Method)
		meta := call.msg.GetTriggerRegistrationMetadata()
		require.Equal(t, types.RegistrationStatus_REGISTRATION_ERROR, meta.Status)
		require.Equal(t, triggerID, meta.TriggerId)
		caperror := caperrors.DeserializeErrorFromString(meta.RegistrationError)
		require.Equal(t, caperrors.InvalidArgument, caperror.Code())
		require.Equal(t, caperrors.OriginUser, caperror.Origin())
		require.Equal(t, caperrors.VisibilityPublic, caperror.Visibility())
		require.Equal(t, "[3]InvalidArgument: its borken", caperror.Error())
		if call.peerID == peerID1 {
			foundPeer1 = true
		} else if call.peerID == peerID2 {
			foundPeer2 = true
		}
	}

	require.True(t, foundPeer1, "did not find message sent to peerID1")
	require.True(t, foundPeer2, "did not find message sent to peerID2")

	// Send an additional registration request after trigger registration has already completed with error, expect to get immediate error response to the new peer
	peerID3 := peers[2]

	pr.AddRegistrationRequest(context.Background(), peerID3, rawRequest, callerDon, registrationExpiry)

	// Confirm a trigger response message was sent to the new peer
	require.Len(t, sendCalls, 3)

	foundPeer3 := false
	for _, call := range sendCalls {
		if call.peerID == peerID3 {
			foundPeer3 = true
			require.Equal(t, types.TriggerRegistrationStatus, call.msg.Method)
			meta := call.msg.GetTriggerRegistrationMetadata()
			require.Equal(t, types.RegistrationStatus_REGISTRATION_ERROR, meta.Status)
			require.Equal(t, triggerID, meta.TriggerId)
			caperror := caperrors.DeserializeErrorFromString(meta.RegistrationError)
			require.Equal(t, caperrors.InvalidArgument, caperror.Code())
			require.Equal(t, caperrors.OriginUser, caperror.Origin())
			require.Equal(t, caperrors.VisibilityPublic, caperror.Visibility())
			require.Equal(t, "[3]InvalidArgument: its borken", caperror.Error())

			require.Equal(t, triggerID, meta.TriggerId)
		}
	}

	require.True(t, foundPeer3, "did not find message sent to peerID3")
}

func TestMultipleAddedRequests_TriggerRegistrationInitiallyErrorsThenRecovers(t *testing.T) {
	// Test that when multiple requests are added to a PublisherRegistration,
	// they all receive the same error response when the trigger registration fails.

	// Setup
	lggr := logger.TestLogger(t)
	triggerID := "test-trigger"
	workflowID := "test-workflow"
	capabilityID := "test-capability"

	var sendCalls []sendCall

	testCapability := &testTriggerCapability{err: errors.New("its borken")}
	pr := registration.NewPublisherRegistration(lggr, testCapability, triggerID, workflowID, capabilityID, func(peerID p2ptypes.PeerID, msgBody *types.MessageBody) error {
		sendCalls = append(sendCalls, sendCall{peerID, msgBody})
		return nil
	})

	peers := make([]p2ptypes.PeerID, 4)
	require.NoError(t, peers[0].UnmarshalText([]byte(peerID1)))
	require.NoError(t, peers[1].UnmarshalText([]byte(peerID2)))

	request := commoncap.TriggerRegistrationRequest{
		TriggerID: triggerID,
	}

	rawRequest, err := pb.MarshalTriggerRegistrationRequest(request)
	require.NoError(t, err)

	// Add multiple requests
	peerID1 := peers[0]
	peerID2 := peers[1]

	callerDon := commoncap.DON{
		ID: 1,
	}

	registrationExpiry := 10 * time.Minute

	pr.AddRegistrationRequest(context.Background(), peerID1, rawRequest, callerDon, registrationExpiry)
	pr.AddRegistrationRequest(context.Background(), peerID2, rawRequest, callerDon, registrationExpiry)

	// Confirm a trigger response message was sent to each peer
	require.Len(t, sendCalls, 2)

	foundPeer1 := false
	foundPeer2 := false
	for _, call := range sendCalls {
		require.Equal(t, types.TriggerRegistrationStatus, call.msg.Method)
		meta := call.msg.GetTriggerRegistrationMetadata()
		require.Equal(t, types.RegistrationStatus_REGISTRATION_ERROR, meta.Status)
		require.Equal(t, triggerID, meta.TriggerId)
		caperror := caperrors.DeserializeErrorFromString(meta.RegistrationError)
		require.Equal(t, caperrors.Unknown, caperror.Code())
		require.Equal(t, caperrors.OriginSystem, caperror.Origin())
		require.Equal(t, caperrors.VisibilityPublic, caperror.Visibility())
		require.Equal(t, "[2]Unknown: its borken", caperror.Error())
		if call.peerID == peerID1 {
			foundPeer1 = true
		} else if call.peerID == peerID2 {
			foundPeer2 = true
		}
	}

	require.True(t, foundPeer1, "did not find message sent to peerID1")
	require.True(t, foundPeer2, "did not find message sent to peerID2")

	// Now update the test capability to succeed
	testCapability.err = nil
	// Send another registration request to trigger re-registration, expect it to return success
	peerID3 := peers[2]

	pr.AddRegistrationRequest(context.Background(), peerID3, rawRequest, callerDon, registrationExpiry)

	// Confirm a trigger response message was sent to the new peer
	require.Len(t, sendCalls, 3)

	foundPeer3 := false
	for _, call := range sendCalls {
		if call.peerID == peerID3 {
			foundPeer3 = true
			require.Equal(t, types.TriggerRegistrationStatus, call.msg.Method)
			meta := call.msg.GetTriggerRegistrationMetadata()
			require.Equal(t, types.RegistrationStatus_REGISTERED, meta.Status)
			require.Equal(t, "", meta.RegistrationError)
			require.Equal(t, triggerID, meta.TriggerId)
		}
	}

	require.True(t, foundPeer3, "did not find message sent to peerID3")

	// Resend a request from peerID1 to confirm it now also gets success response
	pr.AddRegistrationRequest(context.Background(), peerID1, rawRequest, callerDon, registrationExpiry)

	// Confirm a trigger response message was sent to peerID1 again
	require.Len(t, sendCalls, 4)

	call := sendCalls[3]
	if call.peerID == peerID1 {
		require.Equal(t, types.TriggerRegistrationStatus, call.msg.Method)
		meta := call.msg.GetTriggerRegistrationMetadata()
		require.Equal(t, types.RegistrationStatus_REGISTERED, meta.Status)
		require.Equal(t, "", meta.RegistrationError)
		require.Equal(t, triggerID, meta.TriggerId)
	}
}
