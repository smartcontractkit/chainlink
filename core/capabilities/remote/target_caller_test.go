package remote_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	commoncap "github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink-common/pkg/values"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote"
	remotetypes "github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/types"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/transmission"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"
)

const (
	executeValue1 = "triggerEvent1"
)

func Test_TargetCallerExecuteContextTimeout(t *testing.T) {
	lggr := logger.TestLogger(t)
	ctx := testutils.Context(t)

	p1 := p2ptypes.PeerID{}
	require.NoError(t, p1.UnmarshalText([]byte(PeerID1)))
	p2 := p2ptypes.PeerID{}
	require.NoError(t, p2.UnmarshalText([]byte(PeerID2)))
	capDonInfo := commoncap.DON{
		ID:      "capability-don",
		Members: []p2ptypes.PeerID{p1},
		F:       0,
	}

	capInfo := commoncap.CapabilityInfo{
		ID:             "cap_id",
		CapabilityType: commoncap.CapabilityTypeTarget,
		Description:    "Remote Target",
		Version:        "0.0.1",
	}

	workflowDonInfo := commoncap.DON{
		ID:      "workflow-don",
		Members: []p2ptypes.PeerID{p2},
		F:       0,
	}

	dispatcher := NewTestDispatcher()

	caller, err := remote.NewRemoteTargetCaller(lggr, capInfo, capDonInfo, workflowDonInfo, dispatcher)
	require.NoError(t, err)

	err = dispatcher.SetReceiver("cap_id", "workflow-don", caller)
	require.NoError(t, err)

	transmissionSchedule, err := values.NewMap(map[string]any{
		"schedule":   transmission.Schedule_AllAtOnce,
		"deltaStage": "100ms",
	})
	require.NoError(t, err)

	ctxWithTimeout, cancel := context.WithTimeout(ctx, 10*time.Millisecond)
	defer cancel()

	_, err = caller.Execute(ctxWithTimeout,
		commoncap.CapabilityRequest{
			Metadata: commoncap.RequestMetadata{
				WorkflowID:          "workflowID",
				WorkflowExecutionID: "workflowExecutionID",
			},
			Config: transmissionSchedule,
		})

	assert.NotNil(t, err)
}

func Test_TargetCallerExecute(t *testing.T) {

	lggr := logger.TestLogger(t)
	ctx := testutils.Context(t)

	p1 := p2ptypes.PeerID{}
	require.NoError(t, p1.UnmarshalText([]byte(PeerID1)))
	p2 := p2ptypes.PeerID{}
	require.NoError(t, p2.UnmarshalText([]byte(PeerID2)))
	capDonInfo := commoncap.DON{
		ID:      "capability-don",
		Members: []p2ptypes.PeerID{p1},
		F:       0,
	}

	capInfo := commoncap.CapabilityInfo{
		ID:             "cap_id",
		CapabilityType: commoncap.CapabilityTypeTarget,
		Description:    "Remote Target",
		Version:        "0.0.1",
	}

	workflowDonInfo := commoncap.DON{
		ID:      "workflow-don",
		Members: []p2ptypes.PeerID{p2},
		F:       0,
	}

	dispatcher := NewTestDispatcher()

	caller, err := remote.NewRemoteTargetCaller(lggr, capInfo, capDonInfo, workflowDonInfo, dispatcher)
	require.NoError(t, err)

	err = dispatcher.SetReceiver("cap_id", "workflow-don", caller)
	require.NoError(t, err)

	go func() {
		sentMessage := <-dispatcher.sentMessagesCh

		executeValue, err := values.Wrap(executeValue1)
		require.NoError(t, err)
		capResponse := commoncap.CapabilityResponse{
			Value: executeValue,
			Err:   nil,
		}
		marshaled, err := pb.MarshalCapabilityResponse(capResponse)
		require.NoError(t, err)
		executeResponse := &remotetypes.MessageBody{
			Sender:    p1[:],
			Method:    remotetypes.MethodExecute,
			Payload:   marshaled,
			MessageId: sentMessage.MessageId,
		}

		dispatcher.SendToReceiver(executeResponse)
	}()

	transmissionSchedule, err := values.NewMap(map[string]any{
		"schedule":   transmission.Schedule_AllAtOnce,
		"deltaStage": "100ms",
	})
	require.NoError(t, err)

	resultCh, err := caller.Execute(ctx,
		commoncap.CapabilityRequest{
			Metadata: commoncap.RequestMetadata{
				WorkflowID:          "workflowID",
				WorkflowExecutionID: "workflowExecutionID",
			},
			Config: transmissionSchedule,
		})

	require.NoError(t, err)

	response := <-resultCh

	responseValue, err := response.Value.Unwrap()
	assert.Equal(t, executeValue1, responseValue.(string))

}

func Test_TargetCallerExecuteWithError(t *testing.T) {

	lggr := logger.TestLogger(t)
	ctx := testutils.Context(t)

	p1 := p2ptypes.PeerID{}
	require.NoError(t, p1.UnmarshalText([]byte(PeerID1)))
	p2 := p2ptypes.PeerID{}
	require.NoError(t, p2.UnmarshalText([]byte(PeerID2)))
	capDonInfo := commoncap.DON{
		ID:      "capability-don",
		Members: []p2ptypes.PeerID{p1},
		F:       0,
	}

	capInfo := commoncap.CapabilityInfo{
		ID:             "cap_id",
		CapabilityType: commoncap.CapabilityTypeTarget,
		Description:    "Remote Target",
		Version:        "0.0.1",
	}

	workflowDonInfo := commoncap.DON{
		ID:      "workflow-don",
		Members: []p2ptypes.PeerID{p2},
		F:       0,
	}

	dispatcher := NewTestDispatcher()

	caller, err := remote.NewRemoteTargetCaller(lggr, capInfo, capDonInfo, workflowDonInfo, dispatcher)
	require.NoError(t, err)

	err = dispatcher.SetReceiver("cap_id", "workflow-don", caller)
	require.NoError(t, err)

	go func() {
		sentMessage := <-dispatcher.sentMessagesCh

		require.NoError(t, err)
		executeResponse := &remotetypes.MessageBody{
			Sender:    p1[:],
			Method:    remotetypes.MethodExecute,
			MessageId: sentMessage.MessageId,
			Error:     remotetypes.Error_CAPABILITY_NOT_FOUND,
		}

		dispatcher.SendToReceiver(executeResponse)
	}()

	transmissionSchedule, err := values.NewMap(map[string]any{
		"schedule":   transmission.Schedule_AllAtOnce,
		"deltaStage": "100ms",
	})
	require.NoError(t, err)

	_, err = caller.Execute(ctx,
		commoncap.CapabilityRequest{
			Metadata: commoncap.RequestMetadata{
				WorkflowID:          "workflowID",
				WorkflowExecutionID: "workflowExecutionID",
			},
			Config: transmissionSchedule,
		})

	require.NotNil(t, err)
}

type TestDispatcher struct {
	sentMessagesCh chan *remotetypes.MessageBody
	receiver       remotetypes.Receiver
}

func NewTestDispatcher() *TestDispatcher {
	return &TestDispatcher{
		sentMessagesCh: make(chan *remotetypes.MessageBody, 1),
	}
}

func (t *TestDispatcher) SendToReceiver(msgBody *remotetypes.MessageBody) {
	t.receiver.Receive(msgBody)
}

func (t *TestDispatcher) SetReceiver(capabilityId string, donId string, receiver remotetypes.Receiver) error {
	t.receiver = receiver
	return nil
}

func (t *TestDispatcher) RemoveReceiver(capabilityId string, donId string) {}

func (t *TestDispatcher) Send(peerID p2ptypes.PeerID, msgBody *remotetypes.MessageBody) error {
	t.sentMessagesCh <- msgBody
	return nil
}
