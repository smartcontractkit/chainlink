package registration_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/trigger/registration"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/types"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	p2ptypes "github.com/smartcontractkit/libocr/ragep2p/types"
)

const (
	peerID1 = "12D3KooWF3dVeJ6YoT5HFnYhmwQWWMoEwVFzJQ5kKCMX3ZityxMC"
	peerID2 = "12D3KooWQsmok6aD8PZqt3RnJhQRrNzKHLficq7zYFRp7kZ1hHP8"
	peerID3 = "12D3KooWPumsXxg6mJ4hmRizjBD7oFMtN9vm3kTwN8BLEinyDPJS"
	peerID4 = "12D3KooWNuumb38Jpw6DoRbgwejcZYwfsWbbzPU4fWy5imrW3dyD"
)

func TestSubscriberRegistration_SuccessfulRegistration(t *testing.T) {
	lggr := logger.TestLogger(t)
	ctx := t.Context()

	peers := make([]p2ptypes.PeerID, 8)
	require.NoError(t, peers[0].UnmarshalText([]byte(peerID1)))
	require.NoError(t, peers[1].UnmarshalText([]byte(peerID2)))
	require.NoError(t, peers[2].UnmarshalText([]byte(peerID3)))
	require.NoError(t, peers[3].UnmarshalText([]byte(peerID4)))

	reg := registration.NewSubscriberRegistration(lggr, []byte("rawRequest"))

	registrationResponseMessage := &types.MessageBody{
		CapabilityId:    "evm1",
		CapabilityDonId: 2,
		CallerDonId:     3,
		Method:          types.TriggerRegistrationStatus,
		Metadata: &types.MessageBody_TriggerRegistrationMetadata{
			TriggerRegistrationMetadata: &types.TriggerRegistrationMetadata{
				TriggerId:  "w1t1",
				WorkflowId: "w1",
				Status:     types.RegistrationStatus_REGISTERED,
			},
		},
	}

	go func() {
		capDonF := uint8(1)
		minResponsesToAggregate := uint32(capDonF + 1)
		numberOfCapabilityNodes := int(capDonF*2 + 1)
		messageExpiry := 60000 * time.Millisecond
		reg.HandleTriggerRegistrationStatusUpdate(peers[0], registrationResponseMessage, minResponsesToAggregate, messageExpiry, numberOfCapabilityNodes)
		reg.HandleTriggerRegistrationStatusUpdate(peers[1], registrationResponseMessage, minResponsesToAggregate, messageExpiry, numberOfCapabilityNodes)
		reg.HandleTriggerRegistrationStatusUpdate(peers[2], registrationResponseMessage, minResponsesToAggregate, messageExpiry, numberOfCapabilityNodes)
	}()

	err := reg.AwaitRegistration(ctx, 10*time.Minute)
	require.NoError(t, err)

	// Ensure that the second wait for response returns immediately
	err = reg.AwaitRegistration(ctx, 10*time.Minute)
	require.NoError(t, err)

	require.NotNil(t, reg.GetTriggerResponseChannel())
}

func TestSubscriberRegistration_ZeroTimeout(t *testing.T) {
	lggr := logger.TestLogger(t)
	ctx := t.Context()

	// Zero timeout indicates the caller does not want to wait for a registration status update
	reg := registration.NewSubscriberRegistration(lggr, []byte("rawRequest"))

	err := reg.AwaitRegistration(ctx, 0)

	require.ErrorIs(t, err, registration.TerminatedBeforeRegistrationUpdate)
}

func TestSubscriberRegistration_UnsuccessfulRegistration_SameError(t *testing.T) {
	lggr := logger.TestLogger(t)
	ctx := t.Context()

	peers := make([]p2ptypes.PeerID, 8)
	require.NoError(t, peers[0].UnmarshalText([]byte(peerID1)))
	require.NoError(t, peers[1].UnmarshalText([]byte(peerID2)))
	require.NoError(t, peers[2].UnmarshalText([]byte(peerID3)))
	require.NoError(t, peers[3].UnmarshalText([]byte(peerID4)))

	reg := registration.NewSubscriberRegistration(lggr, []byte("rawRequest"))

	errMsg := "its broken"
	registrationResponseMessage := &types.MessageBody{
		CapabilityId:    "evm1",
		CapabilityDonId: 2,
		CallerDonId:     3,
		Method:          types.TriggerRegistrationStatus,
		Metadata: &types.MessageBody_TriggerRegistrationMetadata{
			TriggerRegistrationMetadata: &types.TriggerRegistrationMetadata{
				TriggerId:         "w1t1",
				WorkflowId:        "w1",
				Status:            types.RegistrationStatus_REGISTRATION_ERROR,
				RegistrationError: errMsg,
			},
		},
	}

	go func() {
		capDonF := uint8(1)
		minResponsesToAggregate := uint32(capDonF + 1)
		numberOfCapabilityNodes := int(capDonF*2 + 1)
		messageExpiry := 60000 * time.Millisecond
		reg.HandleTriggerRegistrationStatusUpdate(peers[0], registrationResponseMessage, minResponsesToAggregate, messageExpiry, numberOfCapabilityNodes)
		reg.HandleTriggerRegistrationStatusUpdate(peers[1], registrationResponseMessage, minResponsesToAggregate, messageExpiry, numberOfCapabilityNodes)
		reg.HandleTriggerRegistrationStatusUpdate(peers[2], registrationResponseMessage, minResponsesToAggregate, messageExpiry, numberOfCapabilityNodes)
	}()

	err := reg.AwaitRegistration(ctx, 10*time.Minute)
	require.Error(t, err)
	require.Equal(t, err.Error(), "[2]Unknown: its broken")

	// Ensure that the second wait for response returns immediately
	err = reg.AwaitRegistration(ctx, 10*time.Minute)
	require.Error(t, err)
	require.Equal(t, err.Error(), "[2]Unknown: its broken")
}

func TestSubscriberRegistration_UnsuccessfulRegistration_MixedErrors(t *testing.T) {
	lggr := logger.TestLogger(t)
	ctx := t.Context()

	peers := make([]p2ptypes.PeerID, 8)
	require.NoError(t, peers[0].UnmarshalText([]byte(peerID1)))
	require.NoError(t, peers[1].UnmarshalText([]byte(peerID2)))
	require.NoError(t, peers[2].UnmarshalText([]byte(peerID3)))
	require.NoError(t, peers[3].UnmarshalText([]byte(peerID4)))

	reg := registration.NewSubscriberRegistration(lggr, []byte("rawRequest"))

	errMsg1 := "its broken1"
	errMsg2 := "its broken2"
	errMsg3 := "its broken3"

	go func() {
		capDonF := uint8(1)
		minResponsesToAggregate := uint32(capDonF + 1)
		numberOfCapabilityNodes := int(capDonF*2 + 1)
		messageExpiry := 60000 * time.Millisecond
		reg.HandleTriggerRegistrationStatusUpdate(peers[0], createRegisterResponseMessageWithError(errMsg1), minResponsesToAggregate, messageExpiry, numberOfCapabilityNodes)
		reg.HandleTriggerRegistrationStatusUpdate(peers[1], createRegisterResponseMessageWithError(errMsg2), minResponsesToAggregate, messageExpiry, numberOfCapabilityNodes)
		reg.HandleTriggerRegistrationStatusUpdate(peers[2], createRegisterResponseMessageWithError(errMsg3), minResponsesToAggregate, messageExpiry, numberOfCapabilityNodes)
	}()

	err := reg.AwaitRegistration(ctx, 10*time.Minute)
	require.Error(t, err)
	require.NotNil(t, reg.GetTriggerResponseChannel())
	require.Contains(t, err.Error(), "[2]Unknown: received 3 errors, last error OK : its broken3")

	// Ensure that the second wait for response returns immediately
	err = reg.AwaitRegistration(ctx, 10*time.Minute)
	require.NotNil(t, err)
	require.Contains(t, err.Error(), "[2]Unknown: received 3 errors, last error OK : its broken3")
}

func createRegisterResponseMessageWithError(errMsg1 string) *types.MessageBody {
	return &types.MessageBody{
		CapabilityId:    "evm1",
		CapabilityDonId: 2,
		CallerDonId:     3,
		Method:          types.TriggerRegistrationStatus,
		Metadata: &types.MessageBody_TriggerRegistrationMetadata{
			TriggerRegistrationMetadata: &types.TriggerRegistrationMetadata{
				TriggerId:         "w1t1",
				WorkflowId:        "w1",
				Status:            types.RegistrationStatus_REGISTRATION_ERROR,
				RegistrationError: errMsg1,
			},
		},
	}
}

func TestSubscriberRegistration_Timesout(t *testing.T) {
	lggr := logger.TestLogger(t)
	ctx := t.Context()

	peers := make([]p2ptypes.PeerID, 8)
	require.NoError(t, peers[0].UnmarshalText([]byte(peerID1)))
	require.NoError(t, peers[1].UnmarshalText([]byte(peerID2)))
	require.NoError(t, peers[2].UnmarshalText([]byte(peerID3)))
	require.NoError(t, peers[3].UnmarshalText([]byte(peerID4)))

	reg := registration.NewSubscriberRegistration(lggr, []byte("rawRequest"))

	registrationResponseMessage := &types.MessageBody{
		CapabilityId:    "evm1",
		CapabilityDonId: 2,
		CallerDonId:     3,
		Method:          types.TriggerRegistrationStatus,
		Metadata: &types.MessageBody_TriggerRegistrationMetadata{
			TriggerRegistrationMetadata: &types.TriggerRegistrationMetadata{
				TriggerId:  "w1t1",
				WorkflowId: "w1",
				Status:     types.RegistrationStatus_REGISTERED,
			},
		},
	}

	go func() {
		capDonF := uint8(1)
		minResponsesToAggregate := uint32(capDonF + 1)
		numberOfCapabilityNodes := int(capDonF*2 + 1)
		messageExpiry := 60000 * time.Millisecond
		reg.HandleTriggerRegistrationStatusUpdate(peers[0], registrationResponseMessage, minResponsesToAggregate, messageExpiry, numberOfCapabilityNodes)
	}()

	err := reg.AwaitRegistration(ctx, 1*time.Millisecond)
	require.Error(t, err)
	require.Equal(t, registration.TerminatedBeforeRegistrationUpdate, err)
	require.NotNil(t, reg.GetTriggerResponseChannel())
}

func TestSubscriberRegistration_InsufficientResponses(t *testing.T) {
	lggr := logger.TestLogger(t)
	ctx := t.Context()

	peers := make([]p2ptypes.PeerID, 8)
	require.NoError(t, peers[0].UnmarshalText([]byte(peerID1)))
	require.NoError(t, peers[1].UnmarshalText([]byte(peerID2)))
	require.NoError(t, peers[2].UnmarshalText([]byte(peerID3)))
	require.NoError(t, peers[3].UnmarshalText([]byte(peerID4)))

	reg := registration.NewSubscriberRegistration(lggr, []byte("rawRequest"))

	registrationResponseMessage := &types.MessageBody{
		CapabilityId:    "evm1",
		CapabilityDonId: 2,
		CallerDonId:     3,
		Method:          types.TriggerRegistrationStatus,
		Metadata: &types.MessageBody_TriggerRegistrationMetadata{
			TriggerRegistrationMetadata: &types.TriggerRegistrationMetadata{
				TriggerId:  "w1t1",
				WorkflowId: "w1",
				Status:     types.RegistrationStatus_REGISTERED,
			},
		},
	}

	ctxWithCancel, cancel := context.WithCancel(ctx)

	go func() {
		capDonF := uint8(1)
		minResponsesToAggregate := uint32(capDonF + 1)
		numberOfCapabilityNodes := int(capDonF*2 + 1)
		messageExpiry := 60000 * time.Millisecond
		reg.HandleTriggerRegistrationStatusUpdate(peers[0], registrationResponseMessage, minResponsesToAggregate, messageExpiry, numberOfCapabilityNodes)

		// Simulate insufficient responses by not sending the third response
		// Cancel the context to simulate timeout or cancellation
		cancel()
	}()

	err := reg.AwaitRegistration(ctxWithCancel, 10*time.Minute)
	require.Equal(t, registration.TerminatedBeforeRegistrationUpdate, err)

	require.Error(t, err)
	require.NotNil(t, reg.GetTriggerResponseChannel())
}
