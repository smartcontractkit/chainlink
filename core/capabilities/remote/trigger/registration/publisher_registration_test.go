package registration_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	p2ptypes "github.com/smartcontractkit/libocr/ragep2p/types"

	"github.com/smartcontractkit/chainlink-common/pkg/services"

	commoncap "github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	caperrors "github.com/smartcontractkit/chainlink-common/pkg/capabilities/errors"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"

	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/trigger/registration"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/types"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

type testTriggerCapability struct {
	t   *testing.T
	err error
}

func (t *testTriggerCapability) AckEvent(ctx context.Context, triggerID string, eventID string, method string) error {
	return nil
}

var _ commoncap.TriggerCapability = (*testTriggerCapability)(nil)

func newTestTriggerCapability(t *testing.T, err error) *testTriggerCapability {
	return &testTriggerCapability{t: t, err: err}
}

func (t *testTriggerCapability) Info(ctx context.Context) (commoncap.CapabilityInfo, error) {
	return commoncap.CapabilityInfo{
		ID: "test-capability",
	}, nil
}

func (t *testTriggerCapability) RegisterTrigger(ctx context.Context, req commoncap.TriggerRegistrationRequest) (<-chan commoncap.TriggerResponse, error) {
	require.NotNil(t.t, req)
	ch := make(chan commoncap.TriggerResponse, 1)
	return ch, t.err
}

func (t *testTriggerCapability) UnregisterTrigger(ctx context.Context, req commoncap.TriggerRegistrationRequest) error {
	require.NotNil(t.t, req)
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

	dispatcher := &testDispatcher{}

	pr := registration.NewPublisherRegistration(lggr,
		registration.NewID(1, workflowID, triggerID), func(registrationID registration.ID, response commoncap.TriggerResponse) {}, newTestTriggerCapability(t, nil),
		capabilityID, 1, "triggerCap", dispatcher,
	)

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
	require.Len(t, dispatcher.sendCalls, 2)

	foundPeer1 := false
	foundPeer2 := false
	for _, call := range dispatcher.sendCalls {
		require.Equal(t, types.TriggerRegistrationStatus, call.msg.Method)
		meta := call.msg.GetTriggerRegistrationMetadata()
		require.Equal(t, types.RegistrationStatus_REGISTERED, meta.Status)
		require.Equal(t, triggerID, meta.TriggerId)
		switch call.peerID {
		case peerID1:
			foundPeer1 = true
		case peerID2:
			foundPeer2 = true
		}
	}

	require.True(t, foundPeer1, "did not find message sent to peerID1")
	require.True(t, foundPeer2, "did not find message sent to peerID2")

	// Send an additional registration request after trigger registration has already completed, expect to get immediate response to the new peer
	peerID3 := peers[2]
	pr.AddRegistrationRequest(context.Background(), peerID3, rawRequest, callerDon, registrationExpiry)

	// Confirm a trigger response message was sent to the new peer
	require.Len(t, dispatcher.sendCalls, 3)

	foundPeer3 := false
	for _, call := range dispatcher.sendCalls {
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

	dispatcher := &testDispatcher{}
	pr := registration.NewPublisherRegistration(lggr,
		registration.NewID(1, workflowID, triggerID), func(registrationID registration.ID, response commoncap.TriggerResponse) {},
		newTestTriggerCapability(t, errors.New("its borken")),
		capabilityID, 1, "triggerCap", dispatcher)

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
	require.Len(t, dispatcher.sendCalls, 2)

	foundPeer1 := false
	foundPeer2 := false
	for _, call := range dispatcher.sendCalls {
		require.Equal(t, types.TriggerRegistrationStatus, call.msg.Method)
		meta := call.msg.GetTriggerRegistrationMetadata()
		require.Equal(t, types.RegistrationStatus_REGISTRATION_ERROR, meta.Status)
		require.Equal(t, triggerID, meta.TriggerId)
		caperror := caperrors.DeserializeErrorFromString(meta.RegistrationError)
		require.Equal(t, caperrors.Unknown, caperror.Code())
		require.Equal(t, caperrors.OriginSystem, caperror.Origin())
		require.Equal(t, caperrors.VisibilityPublic, caperror.Visibility())
		require.Equal(t, "[2]Unknown: its borken", caperror.Error())
		switch call.peerID {
		case peerID1:
			foundPeer1 = true
		case peerID2:
			foundPeer2 = true
		}
	}

	require.True(t, foundPeer1, "did not find message sent to peerID1")
	require.True(t, foundPeer2, "did not find message sent to peerID2")

	// Send an additional registration request after trigger registration has already completed with error, expect to get immediate error response to the new peer
	peerID3 := peers[2]

	pr.AddRegistrationRequest(context.Background(), peerID3, rawRequest, callerDon, registrationExpiry)

	// Confirm a trigger response message was sent to the new peer
	require.Len(t, dispatcher.sendCalls, 3)

	foundPeer3 := false
	for _, call := range dispatcher.sendCalls {
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

	dispatcher := &testDispatcher{}
	pr := registration.NewPublisherRegistration(lggr,
		registration.NewID(1, workflowID, triggerID), func(registrationID registration.ID, response commoncap.TriggerResponse) {},
		newTestTriggerCapability(t, caperrors.NewPublicUserError(errors.New("its borken"), caperrors.InvalidArgument)),
		capabilityID, 1, "triggerCap",
		dispatcher)

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
	require.Len(t, dispatcher.sendCalls, 2)

	foundPeer1 := false
	foundPeer2 := false
	for _, call := range dispatcher.sendCalls {
		require.Equal(t, types.TriggerRegistrationStatus, call.msg.Method)
		meta := call.msg.GetTriggerRegistrationMetadata()
		require.Equal(t, types.RegistrationStatus_REGISTRATION_ERROR, meta.Status)
		require.Equal(t, triggerID, meta.TriggerId)
		caperror := caperrors.DeserializeErrorFromString(meta.RegistrationError)
		require.Equal(t, caperrors.InvalidArgument, caperror.Code())
		require.Equal(t, caperrors.OriginUser, caperror.Origin())
		require.Equal(t, caperrors.VisibilityPublic, caperror.Visibility())
		require.Equal(t, "[3]InvalidArgument: its borken", caperror.Error())
		switch call.peerID {
		case peerID1:
			foundPeer1 = true
		case peerID2:
			foundPeer2 = true
		}
	}

	require.True(t, foundPeer1, "did not find message sent to peerID1")
	require.True(t, foundPeer2, "did not find message sent to peerID2")

	// Send an additional registration request after trigger registration has already completed with error, expect to get immediate error response to the new peer
	peerID3 := peers[2]

	pr.AddRegistrationRequest(context.Background(), peerID3, rawRequest, callerDon, registrationExpiry)

	// Confirm a trigger response message was sent to the new peer
	require.Len(t, dispatcher.sendCalls, 3)

	foundPeer3 := false
	for _, call := range dispatcher.sendCalls {
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

	dispatcher := &testDispatcher{}
	testCapability := newTestTriggerCapability(t, errors.New("its borken"))
	pr := registration.NewPublisherRegistration(lggr,
		registration.NewID(1, workflowID, triggerID), func(registrationID registration.ID, response commoncap.TriggerResponse) {},
		testCapability, capabilityID, 1,
		"triggerCap", dispatcher)

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
	require.Len(t, dispatcher.sendCalls, 2)

	foundPeer1 := false
	foundPeer2 := false
	for _, call := range dispatcher.sendCalls {
		require.Equal(t, types.TriggerRegistrationStatus, call.msg.Method)
		meta := call.msg.GetTriggerRegistrationMetadata()
		require.Equal(t, types.RegistrationStatus_REGISTRATION_ERROR, meta.Status)
		require.Equal(t, triggerID, meta.TriggerId)
		caperror := caperrors.DeserializeErrorFromString(meta.RegistrationError)
		require.Equal(t, caperrors.Unknown, caperror.Code())
		require.Equal(t, caperrors.OriginSystem, caperror.Origin())
		require.Equal(t, caperrors.VisibilityPublic, caperror.Visibility())
		require.Equal(t, "[2]Unknown: its borken", caperror.Error())
		switch call.peerID {
		case peerID1:
			foundPeer1 = true
		case peerID2:
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
	require.Len(t, dispatcher.sendCalls, 3)

	foundPeer3 := false
	for _, call := range dispatcher.sendCalls {
		if call.peerID == peerID3 {
			foundPeer3 = true
			require.Equal(t, types.TriggerRegistrationStatus, call.msg.Method)
			meta := call.msg.GetTriggerRegistrationMetadata()
			require.Equal(t, types.RegistrationStatus_REGISTERED, meta.Status)
			require.Empty(t, meta.RegistrationError)
			require.Equal(t, triggerID, meta.TriggerId)
		}
	}

	require.True(t, foundPeer3, "did not find message sent to peerID3")

	// Resend a request from peerID1 to confirm it now also gets success response
	pr.AddRegistrationRequest(context.Background(), peerID1, rawRequest, callerDon, registrationExpiry)

	// Confirm a trigger response message was sent to peerID1 again
	require.Len(t, dispatcher.sendCalls, 4)

	call := dispatcher.sendCalls[3]
	if call.peerID == peerID1 {
		require.Equal(t, types.TriggerRegistrationStatus, call.msg.Method)
		meta := call.msg.GetTriggerRegistrationMetadata()
		require.Equal(t, types.RegistrationStatus_REGISTERED, meta.Status)
		require.Empty(t, meta.RegistrationError)
		require.Equal(t, triggerID, meta.TriggerId)
	}
}

type testDispatcher struct {
	services.StateMachine
	sendCalls []sendCall
}

func (d *testDispatcher) Send(peerID p2ptypes.PeerID, msgBody *types.MessageBody) error {
	d.sendCalls = append(d.sendCalls, sendCall{peerID, msgBody})
	return nil
}
func (d *testDispatcher) SetReceiver(capabilityID string, donID uint32, receiver types.Receiver) error {
	return nil
}

func (d *testDispatcher) RemoveReceiver(capabilityID string, donID uint32) {}

func (d *testDispatcher) SetReceiverForMethod(capabilityID string, donID uint32, method string, receiver types.Receiver) error {
	return nil
}

func (d *testDispatcher) RemoveReceiverForMethod(capabilityID string, donID uint32, method string) {}

func (d *testDispatcher) Start(ctx context.Context) error {
	return nil
}

func (d *testDispatcher) Close() error {
	return nil
}

func (d *testDispatcher) Ready() error {
	return nil
}

func (d *testDispatcher) HealthReport() map[string]error {
	return nil
}

func (d *testDispatcher) Name() string {
	return ""
}
