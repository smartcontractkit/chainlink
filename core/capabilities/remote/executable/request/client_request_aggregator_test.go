package request_test

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	commoncap "github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-protos/cre/go/values"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/executable/request"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/types"
	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"
)

type countingAggregator struct {
	want    int
	oks     []p2ptypes.PeerID
	errMsgs []string
}

func (a *countingAggregator) decide() (*commoncap.CapabilityResponse, error) {
	if len(a.oks)+len(a.errMsgs) < a.want {
		return nil, nil
	}
	if len(a.errMsgs) > 0 {
		return nil, errors.New(a.errMsgs[0])
	}
	m, err := values.NewMap(map[string]any{"responses": len(a.oks)})
	if err != nil {
		return nil, err
	}
	return &commoncap.CapabilityResponse{Value: m}, nil
}

func (a *countingAggregator) OnResponse(peer p2ptypes.PeerID, _ commoncap.CapabilityResponse) (*commoncap.CapabilityResponse, error) {
	a.oks = append(a.oks, peer)
	return a.decide()
}

func (a *countingAggregator) OnError(_ p2ptypes.PeerID, errMsg string) (*commoncap.CapabilityResponse, error) {
	a.errMsgs = append(a.errMsgs, errMsg)
	return a.decide()
}

func Test_ClientRequest_WithAggregator(t *testing.T) {
	t.Parallel()
	workflowDonInfo := commoncap.DON{Members: []p2ptypes.PeerID{NewP2PPeerID(t)}, ID: 2}
	capabilityRequest := commoncap.CapabilityRequest{
		Metadata: commoncap.RequestMetadata{
			WorkflowID:          workflowID1,
			WorkflowExecutionID: workflowExecutionID1,
			ReferenceID:         stepRef1,
		},
	}

	// Every peer returns a different payload; the default quorum would never agree.
	message := func(t *testing.T, sender p2ptypes.PeerID, i int, errMsg string) *types.MessageBody {
		m, err := values.NewMap(map[string]any{"node": i})
		require.NoError(t, err)
		raw, err := pb.MarshalCapabilityResponse(commoncap.CapabilityResponse{Value: m})
		require.NoError(t, err)
		msg := &types.MessageBody{Method: types.MethodExecute, Payload: raw, MessageId: []byte("messageID"), Sender: sender[:]}
		if errMsg != "" {
			msg.Error = types.Error_INTERNAL_ERROR
			msg.ErrorMsg = errMsg
			msg.Payload = nil
		}
		return msg
	}

	t.Run("aggregator decides the response", func(t *testing.T) {
		t.Parallel()
		capabilityPeers, _, capInfo := capabilityDon(t, 4, 1)
		agg := &countingAggregator{want: 3}
		req, err := request.NewClientExecuteRequestWithAggregator(t.Context(), logger.Test(t), capabilityRequest, capInfo,
			workflowDonInfo, newClientRequestTestDispatcher(), 10*time.Minute, "", agg)
		require.NoError(t, err)
		defer req.Cancel(errors.New("test end"))

		require.NoError(t, req.OnMessage(t.Context(), message(t, capabilityPeers[0], 0, "")))
		require.NoError(t, req.OnMessage(t.Context(), message(t, capabilityPeers[1], 1, "")))
		require.Error(t, req.OnMessage(t.Context(), message(t, capabilityPeers[1], 1, "")), "duplicate senders are still rejected")
		require.NoError(t, req.OnMessage(t.Context(), message(t, capabilityPeers[2], 2, "")))

		resp := <-req.ResponseChan()
		require.NoError(t, resp.Err)
		got, err := pb.UnmarshalCapabilityResponse(resp.Result)
		require.NoError(t, err)
		var v map[string]int
		require.NoError(t, got.Value.UnwrapTo(&v))
		assert.Equal(t, 3, v["responses"])
		assert.Len(t, agg.oks, 3)
	})

	t.Run("aggregator error is returned", func(t *testing.T) {
		t.Parallel()
		capabilityPeers, _, capInfo := capabilityDon(t, 4, 1)
		agg := &countingAggregator{want: 2}
		req, err := request.NewClientExecuteRequestWithAggregator(t.Context(), logger.Test(t), capabilityRequest, capInfo,
			workflowDonInfo, newClientRequestTestDispatcher(), 10*time.Minute, "", agg)
		require.NoError(t, err)
		defer req.Cancel(errors.New("test end"))

		require.NoError(t, req.OnMessage(t.Context(), message(t, capabilityPeers[0], 0, "boom")))
		require.NoError(t, req.OnMessage(t.Context(), message(t, capabilityPeers[1], 1, "")))

		resp := <-req.ResponseChan()
		require.EqualError(t, resp.Err, "boom")
	})

	t.Run("nil aggregator is rejected", func(t *testing.T) {
		t.Parallel()
		_, _, capInfo := capabilityDon(t, 4, 1)
		_, err := request.NewClientExecuteRequestWithAggregator(t.Context(), logger.Test(t), capabilityRequest, capInfo,
			workflowDonInfo, newClientRequestTestDispatcher(), 10*time.Minute, "", nil)
		require.Error(t, err)
	})
}
