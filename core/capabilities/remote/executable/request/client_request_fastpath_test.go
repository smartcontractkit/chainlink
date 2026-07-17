package request_test

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/anypb"

	commoncap "github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	vaultcommon "github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-protos/cre/go/values"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/executable/request"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/types"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/vaultshare"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/transmission"
	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"
)

// Test_ClientRequest_FastPathAggregator_OnMessage_MergesFPlusOneShares verifies that the
// ClientRequest routes incoming responses through the provided aggregator when the response count
// reaches F+1. It is the requester-side counterpart to the end-to-end integration test in
// executable/fastpath_integration_test.go.
func Test_ClientRequest_FastPathAggregator_OnMessage_MergesFPlusOneShares(t *testing.T) {
	ctx := t.Context()

	capabilityPeers, capDonInfo, capInfo := capabilityDon(t, 4, 1)
	workflowPeers := make([]p2ptypes.PeerID, 1)
	workflowPeers[0] = NewP2PPeerID(t)
	workflowDonInfo := commoncap.DON{Members: workflowPeers, ID: 2}

	executeInputs, err := values.NewMap(map[string]any{"executeValue1": "aValue1"})
	require.NoError(t, err)

	capabilityRequest := commoncap.CapabilityRequest{
		Metadata: commoncap.RequestMetadata{
			WorkflowID:          workflowID1,
			WorkflowExecutionID: workflowExecutionID1,
			ReferenceID:         stepRef1,
		},
		Inputs: executeInputs,
		Config: mustFastPathTransmissionConfig(t),
	}

	dispatcher := newClientRequestTestDispatcher()
	agg := vaultshare.NewAggregator(2, 3) // F+1=2, max=2F+1=3
	req, err := request.NewClientExecuteRequestWithAggregator(ctx, logger.Test(t), capabilityRequest, capInfo,
		workflowDonInfo, dispatcher, 10*time.Minute, nil, "", nil, agg)
	require.NoError(t, err)
	defer req.Cancel(errors.New("test end"))

	drainInitialPeerSends(t, dispatcher, len(capabilityPeers))

	encKey := "enc-key"
	for i, peer := range capabilityPeers[:2] {
		msg := &types.MessageBody{
			CapabilityId:    capInfo.ID,
			CapabilityDonId: capDonInfo.ID,
			CallerDonId:     workflowDonInfo.ID,
			Method:          types.MethodExecute,
			Payload:         makeVaultCapabilityResponse(t, "owner", encKey, []byte{byte(i + 1)}),
			MessageId:       []byte(req.ID()),
			Sender:          peer[:],
		}
		require.NoError(t, req.OnMessage(ctx, msg))
	}

	select {
	case resp := <-req.ResponseChan():
		require.NoError(t, resp.Err)
		capabilityResp, err := pb.UnmarshalCapabilityResponse(resp.Result)
		require.NoError(t, err)
		merged := &vaultcommon.GetSecretsResponse{}
		require.NoError(t, capabilityResp.Payload.UnmarshalTo(merged))
		require.Len(t, merged.Responses, 1)
		shares := merged.Responses[0].GetData().EncryptedDecryptionKeyShares
		require.Len(t, shares, 1)
		require.Len(t, shares[0].BinaryShares, 2)
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for aggregated fast-path response")
	}
}

// Test_ClientRequest_FastPathAggregator_Cancel_InvokesOnTimeout verifies that cancelling a
// ClientRequest that has an aggregator calls Aggregator.OnTimeout() and returns the resulting error
// on the response channel. This exercises the path taken when fewer than F+1 peers respond to a
// fast-path GetSecrets request.
func Test_ClientRequest_FastPathAggregator_Cancel_InvokesOnTimeout(t *testing.T) {
	ctx := t.Context()

	capabilityPeers, capDonInfo, capInfo := capabilityDon(t, 4, 1)
	workflowPeers := make([]p2ptypes.PeerID, 1)
	workflowPeers[0] = NewP2PPeerID(t)
	workflowDonInfo := commoncap.DON{Members: workflowPeers, ID: 2}

	executeInputs, err := values.NewMap(map[string]any{"executeValue1": "aValue1"})
	require.NoError(t, err)

	capabilityRequest := commoncap.CapabilityRequest{
		Metadata: commoncap.RequestMetadata{
			WorkflowID:          workflowID1,
			WorkflowExecutionID: workflowExecutionID1,
			ReferenceID:         stepRef1,
		},
		Inputs: executeInputs,
		Config: mustFastPathTransmissionConfig(t),
	}

	dispatcher := newClientRequestTestDispatcher()
	agg := vaultshare.NewAggregator(2, 3)
	req, err := request.NewClientExecuteRequestWithAggregator(ctx, logger.Test(t), capabilityRequest, capInfo,
		workflowDonInfo, dispatcher, 10*time.Minute, nil, "", nil, agg)
	require.NoError(t, err)
	defer req.Cancel(errors.New("test end"))

	drainInitialPeerSends(t, dispatcher, len(capabilityPeers))

	// Only one peer responds, below the F+1 threshold. Then we cancel the request.
	msg := &types.MessageBody{
		CapabilityId:    capInfo.ID,
		CapabilityDonId: capDonInfo.ID,
		CallerDonId:     workflowDonInfo.ID,
		Method:          types.MethodExecute,
		Payload:         makeVaultCapabilityResponse(t, "owner", "enc-key", []byte{1}),
		MessageId:       []byte(req.ID()),
		Sender:          capabilityPeers[0][:],
	}
	require.NoError(t, req.OnMessage(ctx, msg))

	req.Cancel(errors.New("request timed out"))

	select {
	case resp := <-req.ResponseChan():
		require.ErrorIs(t, resp.Err, vaultshare.ErrFastPathReconstruction)
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for fast-path timeout response")
	}
}

func makeVaultCapabilityResponse(t *testing.T, owner, encKey string, share []byte) []byte {
	t.Helper()
	vaultResp := &vaultcommon.GetSecretsResponse{
		Responses: []*vaultcommon.SecretResponse{{
			Id: &vaultcommon.SecretIdentifier{Owner: owner, Namespace: "main", Key: "secret"},
			Result: &vaultcommon.SecretResponse_Data{
				Data: &vaultcommon.SecretData{
					EncryptedValue: "abc123",
					EncryptedDecryptionKeyShares: []*vaultcommon.EncryptedShares{{
						EncryptionKey: encKey,
						BinaryShares:  [][]byte{share},
					}},
				},
			},
		}},
	}
	anyPayload, err := anypb.New(vaultResp)
	require.NoError(t, err)
	payload, err := pb.MarshalCapabilityResponse(commoncap.CapabilityResponse{Payload: anyPayload})
	require.NoError(t, err)
	return payload
}

func mustFastPathTransmissionConfig(t *testing.T) *values.Map {
	t.Helper()
	transmissionSchedule, err := values.NewMap(map[string]any{
		"schedule":   transmission.Schedule_AllAtOnce,
		"deltaStage": "10ms",
	})
	require.NoError(t, err)
	return transmissionSchedule
}
