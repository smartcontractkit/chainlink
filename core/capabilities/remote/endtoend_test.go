package remote_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	commoncap "github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink-common/pkg/services/servicetest"
	"github.com/smartcontractkit/chainlink-common/pkg/values"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote"
	remoteutils "github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/testutils"
	remotetypes "github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/types"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"
)

func testRemoteTrigger(ctx context.Context,
	t *testing.T,
	underlying commoncap.TriggerCapability,
	numWorkflowPeers int,
	workflowDonF uint8,
	workflowNodeTimeout time.Duration,
	numCapabilityPeers int,
	capabilityDonF uint8,
	capabilityNodeResponseTimeout time.Duration,
	transmissionSchedule *values.Map,
	responseTest func(t *testing.T, response commoncap.CapabilityResponse, responseError error)) {
	lggr := logger.TestLogger(t)

	capabilityPeers := make([]p2ptypes.PeerID, numCapabilityPeers)
	for i := 0; i < numCapabilityPeers; i++ {
		capabilityPeerID := p2ptypes.PeerID{}
		require.NoError(t, capabilityPeerID.UnmarshalText([]byte(remoteutils.NewPeerID())))
		capabilityPeers[i] = capabilityPeerID
	}

	capabilityPeerID := p2ptypes.PeerID{}
	require.NoError(t, capabilityPeerID.UnmarshalText([]byte(remoteutils.NewPeerID())))

	capDonInfo := commoncap.DON{
		ID:      2,
		Members: capabilityPeers,
		F:       capabilityDonF,
	}

	capInfo := commoncap.CapabilityInfo{
		ID:             "cap_id@1.0.0",
		CapabilityType: commoncap.CapabilityTypeTarget,
		Description:    "Remote Target",
		DON:            &capDonInfo,
	}

	workflowPeers := make([]p2ptypes.PeerID, numWorkflowPeers)
	for i := 0; i < numWorkflowPeers; i++ {
		workflowPeerID := p2ptypes.PeerID{}
		require.NoError(t, workflowPeerID.UnmarshalText([]byte(remoteutils.NewPeerID())))
		workflowPeers[i] = workflowPeerID
	}

	workflowDonInfo := commoncap.DON{
		Members: workflowPeers,
		ID:      1,
		F:       workflowDonF,
	}

	broker := remoteutils.NewTestAsyncMessageBroker(t, 1000)

	workflowDONs := map[uint32]commoncap.DON{
		workflowDonInfo.ID: workflowDonInfo,
	}

	capabilityNodes := make([]remotetypes.Receiver, numCapabilityPeers)
	for i := 0; i < numCapabilityPeers; i++ {
		capabilityPeer := capabilityPeers[i]
		capabilityDispatcher := broker.NewDispatcherForNode(capabilityPeer)
		config := &commoncap.RemoteTriggerConfig{
			RegistrationRefresh:     100 * time.Millisecond,
			RegistrationExpiry:      100 * time.Second,
			MinResponsesToAggregate: 1,
			MessageExpiry:           100 * time.Second,
			MaxBatchSize:            1,
			BatchCollectionPeriod:   time.Second,
		}
		capabilityNode := remote.NewTriggerPublisher(config, underlying, capInfo, capDonInfo, workflowDONs, capabilityDispatcher, lggr)
		servicetest.Run(t, capabilityNode)
		broker.RegisterReceiverNode(capabilityPeer, capabilityNode)
		capabilityNodes[i] = capabilityNode
	}

	workflowNodes := make([]commoncap.TriggerCapability, numWorkflowPeers)
	for i := 0; i < numWorkflowPeers; i++ {
		workflowPeerDispatcher := broker.NewDispatcherForNode(workflowPeers[i])
		// register trigger
		config := &commoncap.RemoteTriggerConfig{
			RegistrationRefresh:     100 * time.Millisecond,
			RegistrationExpiry:      100 * time.Second,
			MinResponsesToAggregate: 1,
			MessageExpiry:           100 * time.Second,
		}
		workflowNode := remote.NewTriggerSubscriber(config, capInfo, capDonInfo, workflowDonInfo, workflowPeerDispatcher, nil, lggr)
		// workflowNode := target.NewClient(capInfo, workflowDonInfo, workflowPeerDispatcher, workflowNodeTimeout, lggr)
		servicetest.Run(t, workflowNode)
		<-broker.SendCh
		broker.RegisterReceiverNode(workflowPeers[i], workflowNode)
		workflowNodes[i] = workflowNode
	}

	servicetest.Run(t, broker)

	executeInputs, err := values.NewMap(
		map[string]any{
			"executeValue1": "aValue1",
		},
	)

	require.NoError(t, err)

	wg := &sync.WaitGroup{}
	wg.Add(len(workflowNodes))

	for _, caller := range workflowNodes {
		go func(caller commoncap.TriggerCapability) {
			defer wg.Done()
			req := commoncap.TriggerRegistrationRequest{
				TriggerID: "logeventtrigger_log1",
				Metadata: commoncap.RequestMetadata{
					ReferenceID:         "logeventtrigger",
					WorkflowID:          remoteutils.WorkflowID1,
					WorkflowExecutionID: remoteutils.WorkflowExecutionID1,
				},
			}
			response, err := caller.RegisterTrigger(ctx, req)

			// receive trigger event
			triggerEventValue, err := values.NewMap(triggerEvent1)
			require.NoError(t, err)
			capResponse := commoncap.TriggerResponse{
				Event: commoncap.TriggerEvent{
					Outputs: triggerEventValue,
				},
				Err: nil,
			}
			marshaled, err := pb.MarshalTriggerResponse(capResponse)
			require.NoError(t, err)
			triggerEvent := &remotetypes.MessageBody{
				Sender: p1[:],
				Method: remotetypes.MethodTriggerEvent,
				Metadata: &remotetypes.MessageBody_TriggerEventMetadata{
					TriggerEventMetadata: &remotetypes.TriggerEventMetadata{
						WorkflowIds: []string{workflowID1},
					},
				},
				Payload: marshaled,
			}
			subscriber.Receive(ctx, triggerEvent)

			// response, err := caller.Execute(ctx,
			// 	commoncap.CapabilityRequest{
			// 		Metadata: commoncap.RequestMetadata{
			// 			WorkflowID:          remoteutils.WorkflowID1,
			// 			WorkflowExecutionID: remoteutils.WorkflowExecutionID1,
			// 		},
			// 		Config: transmissionSchedule,
			// 		Inputs: executeInputs,
			// 	})

			responseTest(t, response, err)
		}(caller)
	}

	wg.Wait()
}
