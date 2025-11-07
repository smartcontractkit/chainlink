package executable_test

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/anypb"

	commoncap "github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/v2/chain-capabilities/evm"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	sdkpb "github.com/smartcontractkit/chainlink-protos/cre/go/sdk"
	"github.com/smartcontractkit/chainlink-protos/cre/go/values"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/executable"
	remotetypes "github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/types"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"
)

func Test_Server_Execute_SlowCapabilityExecutionDoesNotImpactSubsequentCall(t *testing.T) {
	ctx := testutils.Context(t)

	numCapabilityPeers := 4

	workflowIDToPause := map[string]time.Duration{}
	workflowIDToPause[workflowID1] = 1 * time.Minute
	workflowIDToPause[workflowID2] = 1 * time.Second

	callers, srvcs := testRemoteExecutableCapabilityServer(ctx, t, &commoncap.RemoteExecutableConfig{}, &TestSlowExecutionCapability{workflowIDToPause: workflowIDToPause}, 10, 9, numCapabilityPeers, 3, 10*time.Minute, nil)

	for _, caller := range callers {
		_, err := caller.Execute(context.Background(),
			commoncap.CapabilityRequest{
				Metadata: commoncap.RequestMetadata{
					WorkflowID:          workflowID1,
					WorkflowExecutionID: workflowExecutionID1,
				},
			})
		require.NoError(t, err)
	}

	for _, caller := range callers {
		_, err := caller.Execute(context.Background(),
			commoncap.CapabilityRequest{
				Metadata: commoncap.RequestMetadata{
					WorkflowID:          workflowID2,
					WorkflowExecutionID: workflowExecutionID2,
				},
			})
		require.NoError(t, err)
	}

	for _, caller := range callers {
		for range numCapabilityPeers {
			msg := <-caller.receivedMessages
			assert.Equal(t, remotetypes.Error_OK, msg.Error)

			capabilityResponse, err := pb.UnmarshalCapabilityResponse(msg.Payload)
			require.NoError(t, err)
			val := capabilityResponse.Value.Underlying["response"]

			var valAsStr string
			err = val.UnwrapTo(&valAsStr)
			require.NoError(t, err)

			assert.Equal(t, "1s", valAsStr)
		}
	}

	closeServices(t, srvcs)
}

func Test_Server_DefaultExcludedAttributes(t *testing.T) {
	ctx := testutils.Context(t)

	numCapabilityPeers := 4

	callers, srvcs := testRemoteExecutableCapabilityServer(ctx, t, &commoncap.RemoteExecutableConfig{},
		&TestCapability{}, 10, 9, numCapabilityPeers, 3, 10*time.Minute, nil)

	for idx, caller := range callers {
		rawInputs := map[string]any{
			"StepDependency": strconv.Itoa(idx),
		}

		inputs, err := values.NewMap(rawInputs)
		require.NoError(t, err)

		_, err = caller.Execute(context.Background(),
			commoncap.CapabilityRequest{
				Metadata: commoncap.RequestMetadata{
					WorkflowID:          workflowID1,
					WorkflowExecutionID: workflowExecutionID1,
				},
				Inputs: inputs,
			})
		require.NoError(t, err)
	}

	for _, caller := range callers {
		for range numCapabilityPeers {
			msg := <-caller.receivedMessages
			assert.Equal(t, remotetypes.Error_OK, msg.Error)
		}
	}
	closeServices(t, srvcs)
}

func Test_Server_ExcludesNonDeterministicInputAttributes(t *testing.T) {
	ctx := testutils.Context(t)

	numCapabilityPeers := 4

	callers, srvcs := testRemoteExecutableCapabilityServer(ctx, t, &commoncap.RemoteExecutableConfig{RequestHashExcludedAttributes: []string{"signed_report.Signatures"}},
		&TestCapability{}, 10, 9, numCapabilityPeers, 3, 10*time.Minute, nil)

	for idx, caller := range callers {
		rawInputs := map[string]any{
			"signed_report": map[string]any{"Signatures": "sig" + strconv.Itoa(idx), "Price": 20},
		}

		inputs, err := values.NewMap(rawInputs)
		require.NoError(t, err)

		_, err = caller.Execute(context.Background(),
			commoncap.CapabilityRequest{
				Metadata: commoncap.RequestMetadata{
					WorkflowID:          workflowID1,
					WorkflowExecutionID: workflowExecutionID1,
				},
				Inputs: inputs,
			})
		require.NoError(t, err)
	}

	for _, caller := range callers {
		for range numCapabilityPeers {
			msg := <-caller.receivedMessages
			assert.Equal(t, remotetypes.Error_OK, msg.Error)
		}
	}
	closeServices(t, srvcs)
}

func Test_Server_Execute_RespondsAfterSufficientRequests(t *testing.T) {
	ctx := testutils.Context(t)

	numCapabilityPeers := 4

	callers, srvcs := testRemoteExecutableCapabilityServer(ctx, t, &commoncap.RemoteExecutableConfig{}, &TestCapability{}, 10, 9, numCapabilityPeers, 3, 10*time.Minute, nil)

	for _, caller := range callers {
		_, err := caller.Execute(context.Background(),
			commoncap.CapabilityRequest{
				Metadata: commoncap.RequestMetadata{
					WorkflowID:          workflowID1,
					WorkflowExecutionID: workflowExecutionID1,
				},
			})
		require.NoError(t, err)
	}

	for _, caller := range callers {
		for range numCapabilityPeers {
			msg := <-caller.receivedMessages
			assert.Equal(t, remotetypes.Error_OK, msg.Error)
		}
	}
	closeServices(t, srvcs)
}

func Test_Server_InsufficientCallers(t *testing.T) {
	ctx := testutils.Context(t)

	numCapabilityPeers := 4

	callers, srvcs := testRemoteExecutableCapabilityServer(ctx, t, &commoncap.RemoteExecutableConfig{}, &TestCapability{}, 10, 10, numCapabilityPeers, 3, 100*time.Millisecond, nil)

	for _, caller := range callers {
		_, err := caller.Execute(context.Background(),
			commoncap.CapabilityRequest{
				Metadata: commoncap.RequestMetadata{
					WorkflowID:          workflowID1,
					WorkflowExecutionID: workflowExecutionID1,
				},
			})
		require.NoError(t, err)
	}

	for _, caller := range callers {
		for range numCapabilityPeers {
			msg := <-caller.receivedMessages
			assert.Equal(t, remotetypes.Error_TIMEOUT, msg.Error)
		}
	}
	closeServices(t, srvcs)
}

func Test_Server_CapabilityError(t *testing.T) {
	ctx := testutils.Context(t)

	numCapabilityPeers := 4

	callers, srvcs := testRemoteExecutableCapabilityServer(ctx, t, &commoncap.RemoteExecutableConfig{}, &TestErrorCapability{}, 10, 9, numCapabilityPeers, 3, 100*time.Millisecond, nil)

	for _, caller := range callers {
		_, err := caller.Execute(context.Background(),
			commoncap.CapabilityRequest{
				Metadata: commoncap.RequestMetadata{
					WorkflowID:          workflowID1,
					WorkflowExecutionID: workflowExecutionID1,
				},
			})
		require.NoError(t, err)
	}

	for _, caller := range callers {
		for range numCapabilityPeers {
			msg := <-caller.receivedMessages
			assert.Equal(t, remotetypes.Error_INTERNAL_ERROR, msg.Error)
		}
	}
	closeServices(t, srvcs)
}

func Test_Server_V2Request_ExcludesNonDeterministicInputAttributes(t *testing.T) {
	ctx := testutils.Context(t)

	numCapabilityPeers := 4

	callers, srvcs := testRemoteExecutableCapabilityServer(ctx, t, &commoncap.RemoteExecutableConfig{RequestHashExcludedAttributes: []string{"signed_report.Signatures"}},
		&TestCapability{}, 10, 9, numCapabilityPeers, 3, 10*time.Minute, &v2WriteChainMessageHasher{})

	report := []byte("report01234")
	for idx, caller := range callers {
		if idx < 0 || idx > 4294967295 { // Check bounds for uint32
			require.Fail(t, "idx out of range for uint32")
		}
		payload := &evm.WriteReportRequest{
			Receiver: []byte("abcdef"),
			Report: &sdkpb.ReportResponse{
				RawReport: report,
				Sigs: []*sdkpb.AttributedSignature{ // non-deterministic set of sigs that we want to ignore when hashing
					{
						SignerId:  uint32(idx), // Now safe after bounds check
						Signature: []byte("sig" + strconv.Itoa(idx)),
					},
				},
			},
		}
		anyPayload, err := anypb.New(payload)
		require.NoError(t, err)

		_, err = caller.Execute(context.Background(),
			commoncap.CapabilityRequest{
				Metadata: commoncap.RequestMetadata{
					WorkflowID:          workflowID1,
					WorkflowExecutionID: workflowExecutionID1,
				},
				Payload: anyPayload,
			})
		require.NoError(t, err)
	}

	for _, caller := range callers {
		for range numCapabilityPeers {
			msg := <-caller.receivedMessages
			assert.Equal(t, remotetypes.Error_OK, msg.Error)
		}
	}
	closeServices(t, srvcs)
}

type v2WriteChainMessageHasher struct{}

func (r *v2WriteChainMessageHasher) Hash(msg *remotetypes.MessageBody) ([32]byte, error) {
	req, err := pb.UnmarshalCapabilityRequest(msg.Payload)
	if err != nil {
		return [32]byte{}, fmt.Errorf("failed to unmarshal capability request: %w", err)
	}
	if req.Payload == nil {
		return [32]byte{}, errors.New("request payload is nil")
	}
	var writeReportRequest evm.WriteReportRequest
	if err = req.Payload.UnmarshalTo(&writeReportRequest); err != nil {
		return [32]byte{}, fmt.Errorf("failed to unmarshal payload to WriteReportRequest: %w", err)
	}
	writeReportRequest.Report.Sigs = nil // exclude signatures from the hash

	req.Payload, err = anypb.New(&writeReportRequest)
	if err != nil {
		return [32]byte{}, fmt.Errorf("failed to marshal WriteReportRequest to anypb: %w", err)
	}
	reqBytes, err := pb.MarshalCapabilityRequest(req)
	if err != nil {
		return [32]byte{}, fmt.Errorf("failed to marshal capability request: %w", err)
	}
	hash := sha256.Sum256(reqBytes)
	return hash, nil
}

func testRemoteExecutableCapabilityServer(ctx context.Context, t *testing.T,
	config *commoncap.RemoteExecutableConfig,
	underlying commoncap.ExecutableCapability,
	numWorkflowPeers int, workflowDonF uint8,
	numCapabilityPeers int, capabilityDonF uint8, capabilityNodeResponseTimeout time.Duration,
	messageHasher remotetypes.MessageHasher) ([]*serverTestClient, []services.Service) {
	lggr := logger.Test(t)
	if config.RequestTimeout == 0 {
		config.RequestTimeout = capabilityNodeResponseTimeout
	}
	if config.ServerMaxParallelRequests == 0 {
		config.ServerMaxParallelRequests = 10
	}

	capabilityPeers := make([]p2ptypes.PeerID, numCapabilityPeers)
	for i := range numCapabilityPeers {
		capabilityPeerID := NewP2PPeerID(t)
		capabilityPeers[i] = capabilityPeerID
	}

	capDonInfo := commoncap.DON{
		ID:      1,
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
	for i := range numWorkflowPeers {
		workflowPeers[i] = NewP2PPeerID(t)
	}

	workflowDonInfo := commoncap.DON{
		Members: workflowPeers,
		ID:      2,
		F:       workflowDonF,
	}

	var srvcs []services.Service
	broker := newTestAsyncMessageBroker(t, 1000)
	err := broker.Start(context.Background())
	require.NoError(t, err)
	srvcs = append(srvcs, broker)

	workflowDONs := map[uint32]commoncap.DON{
		workflowDonInfo.ID: workflowDonInfo,
	}

	capabilityNodes := make([]remotetypes.Receiver, numCapabilityPeers)

	for i := range numCapabilityPeers {
		capabilityPeer := capabilityPeers[i]
		capabilityDispatcher := broker.NewDispatcherForNode(capabilityPeer)
		capabilityNode := executable.NewServer(capInfo.ID, "", capabilityPeer, capabilityDispatcher, lggr)
		require.NoError(t, capabilityNode.SetConfig(config, underlying, capInfo, capDonInfo, workflowDONs, messageHasher))
		require.NoError(t, capabilityNode.Start(ctx))
		broker.RegisterReceiverNode(capabilityPeer, capabilityNode)
		capabilityNodes[i] = capabilityNode
		srvcs = append(srvcs, capabilityNode)
	}

	workflowNodes := make([]*serverTestClient, numWorkflowPeers)
	for i := range numWorkflowPeers {
		workflowPeerDispatcher := broker.NewDispatcherForNode(workflowPeers[i])
		workflowNode := newServerTestClient(workflowPeers[i], capDonInfo, workflowPeerDispatcher)
		broker.RegisterReceiverNode(workflowPeers[i], workflowNode)
		workflowNodes[i] = workflowNode
	}

	return workflowNodes, srvcs
}

func closeServices(t *testing.T, srvcs []services.Service) {
	for _, srv := range srvcs {
		require.NoError(t, srv.Close())
	}
}

type serverTestClient struct {
	peerID            p2ptypes.PeerID
	dispatcher        remotetypes.Dispatcher
	capabilityDonInfo commoncap.DON
	receivedMessages  chan *remotetypes.MessageBody
	callerDonID       string
}

func (r *serverTestClient) Receive(_ context.Context, msg *remotetypes.MessageBody) {
	r.receivedMessages <- msg
}

func newServerTestClient(peerID p2ptypes.PeerID, capabilityDonInfo commoncap.DON,
	dispatcher remotetypes.Dispatcher) *serverTestClient {
	return &serverTestClient{peerID: peerID, dispatcher: dispatcher, capabilityDonInfo: capabilityDonInfo,
		receivedMessages: make(chan *remotetypes.MessageBody, 100), callerDonID: "workflow-don"}
}

func (r *serverTestClient) Info(ctx context.Context) (commoncap.CapabilityInfo, error) {
	panic("not implemented")
}

func (r *serverTestClient) Execute(ctx context.Context, req commoncap.CapabilityRequest) (<-chan commoncap.CapabilityResponse, error) {
	rawRequest, err := pb.MarshalCapabilityRequest(req)
	if err != nil {
		return nil, err
	}

	messageID := remotetypes.MethodExecute + ":" + req.Metadata.WorkflowExecutionID

	for _, node := range r.capabilityDonInfo.Members {
		message := &remotetypes.MessageBody{
			CapabilityId:    "capability-id",
			CapabilityDonId: 1,
			CallerDonId:     2,
			Method:          remotetypes.MethodExecute,
			Payload:         rawRequest,
			MessageId:       []byte(messageID),
			Sender:          r.peerID[:],
			Receiver:        node[:],
		}

		if err = r.dispatcher.Send(node, message); err != nil {
			return nil, err
		}
	}

	return nil, nil
}

func Test_Server_SetConfig(t *testing.T) {
	lggr := logger.Test(t)
	peerID := NewP2PPeerID(t)

	// Create broker and dispatcher
	broker := newTestAsyncMessageBroker(t, 100)
	dispatcher := broker.NewDispatcherForNode(peerID)

	// Create server instance
	server := executable.NewServer("test-capability-id", "test-method", peerID, dispatcher, lggr)

	// Create test data
	capInfo := commoncap.CapabilityInfo{
		ID:             "test-capability-id",
		CapabilityType: commoncap.CapabilityTypeTarget,
		Description:    "Test capability",
	}

	localDonInfo := commoncap.DON{
		ID:      1,
		Members: []p2ptypes.PeerID{peerID},
		F:       0,
	}

	workflowDONs := map[uint32]commoncap.DON{
		2: {
			ID:      2,
			Members: []p2ptypes.PeerID{NewP2PPeerID(t)},
			F:       0,
		},
	}

	underlying := &TestCapability{}
	requestTimeout := 10 * time.Second
	maxParallelRequests := uint32(5)

	t.Run("valid config should succeed", func(t *testing.T) {
		config := &commoncap.RemoteExecutableConfig{
			RequestHashExcludedAttributes: []string{"test"},
			RequestTimeout:                requestTimeout,
			ServerMaxParallelRequests:     maxParallelRequests,
		}

		err := server.SetConfig(config, underlying, capInfo, localDonInfo, workflowDONs, nil)
		require.NoError(t, err)
	})

	t.Run("mismatched capability ID should return error", func(t *testing.T) {
		invalidCapInfo := commoncap.CapabilityInfo{
			ID:             "different-capability-id",
			CapabilityType: commoncap.CapabilityTypeTarget,
		}

		err := server.SetConfig(&commoncap.RemoteExecutableConfig{}, underlying, invalidCapInfo,
			localDonInfo, workflowDONs, nil)
		require.Error(t, err)
		require.Contains(t, err.Error(), "capability info provided does not match")
	})

	t.Run("nil underlying capability should return error", func(t *testing.T) {
		err := server.SetConfig(&commoncap.RemoteExecutableConfig{}, nil, capInfo,
			localDonInfo, workflowDONs, nil)
		require.Error(t, err)
		require.Contains(t, err.Error(), "underlying capability cannot be nil")
	})

	t.Run("empty local DON members should fail", func(t *testing.T) {
		server := executable.NewServer("test-capability-id", "test-method", peerID, dispatcher, lggr)
		emptyLocalDon := commoncap.DON{
			ID:      1,
			Members: []p2ptypes.PeerID{},
			F:       0,
		}
		config := &commoncap.RemoteExecutableConfig{
			RequestTimeout:            10 * time.Second,
			ServerMaxParallelRequests: 5,
		}
		err := server.SetConfig(config, underlying, capInfo, emptyLocalDon, workflowDONs, nil)
		require.Error(t, err)
		require.Contains(t, err.Error(), "empty localDonInfo provided")
	})

	t.Run("nil message hasher should use default", func(t *testing.T) {
		server := executable.NewServer("test-capability-id", "test-method", peerID, dispatcher, lggr)
		config := &commoncap.RemoteExecutableConfig{
			RequestTimeout:            10 * time.Second,
			ServerMaxParallelRequests: 5,
		}
		err := server.SetConfig(config, underlying, capInfo, localDonInfo, workflowDONs, nil)
		require.NoError(t, err)
	})

	t.Run("zero timeout should fail", func(t *testing.T) {
		server := executable.NewServer("test-capability-id", "test-method", peerID, dispatcher, lggr)
		config := &commoncap.RemoteExecutableConfig{
			RequestTimeout:            0,
			ServerMaxParallelRequests: 5,
		}
		err := server.SetConfig(config, underlying, capInfo, localDonInfo, workflowDONs, nil)
		require.Error(t, err)
		require.Contains(t, err.Error(), "RequestTimeout must be positive")
	})

	t.Run("zero max parallel requests should fail", func(t *testing.T) {
		server := executable.NewServer("test-capability-id", "test-method", peerID, dispatcher, lggr)
		config := &commoncap.RemoteExecutableConfig{
			RequestTimeout:            10 * time.Second,
			ServerMaxParallelRequests: 0,
		}
		err := server.SetConfig(config, underlying, capInfo, localDonInfo, workflowDONs, nil)
		require.Error(t, err)
		require.Contains(t, err.Error(), "ServerMaxParallelRequests must be positive")
	})

	t.Run("empty workflow DONs should fail", func(t *testing.T) {
		server := executable.NewServer("test-capability-id", "test-method", peerID, dispatcher, lggr)
		emptyWorkflowDONs := map[uint32]commoncap.DON{}
		config := &commoncap.RemoteExecutableConfig{
			RequestTimeout:            10 * time.Second,
			ServerMaxParallelRequests: 5,
		}
		err := server.SetConfig(config, underlying, capInfo, localDonInfo, emptyWorkflowDONs, nil)
		require.Error(t, err)
		require.Contains(t, err.Error(), "empty workflowDONs provided")
	})
}

func Test_Server_SetConfig_ConfigReplacement(t *testing.T) {
	lggr := logger.Test(t)
	peerID := NewP2PPeerID(t)
	broker := newTestAsyncMessageBroker(t, 100)
	dispatcher := broker.NewDispatcherForNode(peerID)
	server := executable.NewServer("test-capability-id", "test-method", peerID, dispatcher, lggr)

	capInfo := commoncap.CapabilityInfo{
		ID:             "test-capability-id",
		CapabilityType: commoncap.CapabilityTypeTarget,
		Description:    "Test capability",
	}

	localDonInfo := commoncap.DON{
		ID:      1,
		Members: []p2ptypes.PeerID{peerID},
		F:       0,
	}

	workflowDONs := map[uint32]commoncap.DON{
		2: {
			ID:      2,
			Members: []p2ptypes.PeerID{NewP2PPeerID(t)},
			F:       0,
		},
	}

	underlying := &TestCapability{}

	// Set initial config
	config1 := &commoncap.RemoteExecutableConfig{
		RequestHashExcludedAttributes: []string{"attr1"},
		RequestTimeout:                5 * time.Second,
		ServerMaxParallelRequests:     3,
	}
	err := server.SetConfig(config1, underlying, capInfo, localDonInfo, workflowDONs, nil)
	require.NoError(t, err)

	// Verify server can start with valid config
	ctx := testutils.Context(t)
	err = server.Start(ctx)
	require.NoError(t, err)

	// Replace with new config
	config2 := &commoncap.RemoteExecutableConfig{
		RequestHashExcludedAttributes: []string{"attr2", "attr3"},
		RequestTimeout:                10 * time.Second,
		ServerMaxParallelRequests:     5,
	}
	err = server.SetConfig(config2, underlying, capInfo, localDonInfo, workflowDONs, nil)
	require.NoError(t, err)

	// Clean up
	err = server.Close()
	require.NoError(t, err)
}

func Test_Server_SetConfig_StartValidation(t *testing.T) {
	ctx := testutils.Context(t)

	t.Run("Start without SetConfig should fail", func(t *testing.T) {
		lggr := logger.Test(t)
		peerID := NewP2PPeerID(t)
		broker := newTestAsyncMessageBroker(t, 100)
		dispatcher := broker.NewDispatcherForNode(peerID)
		server := executable.NewServer("test-capability-id", "test-method", peerID, dispatcher, lggr)

		err := server.Start(ctx)
		require.Error(t, err)
		require.Contains(t, err.Error(), "config not set - call SetConfig() before Start()")
	})

	t.Run("Start with valid config should succeed", func(t *testing.T) {
		lggr := logger.Test(t)
		peerID := NewP2PPeerID(t)
		broker := newTestAsyncMessageBroker(t, 100)
		dispatcher := broker.NewDispatcherForNode(peerID)
		server := executable.NewServer("test-capability-id", "test-method", peerID, dispatcher, lggr)

		// Set valid config
		capInfo := commoncap.CapabilityInfo{
			ID:             "test-capability-id",
			CapabilityType: commoncap.CapabilityTypeTarget,
			Description:    "Test capability",
		}

		localDonInfo := commoncap.DON{
			ID:      1,
			Members: []p2ptypes.PeerID{peerID},
			F:       0,
		}

		workflowDONs := map[uint32]commoncap.DON{
			2: {
				ID:      2,
				Members: []p2ptypes.PeerID{NewP2PPeerID(t)},
				F:       0,
			},
		}

		underlying := &TestCapability{}
		cfg := &commoncap.RemoteExecutableConfig{
			RequestTimeout:            10 * time.Second,
			ServerMaxParallelRequests: 5,
		}
		err := server.SetConfig(cfg, underlying, capInfo,
			localDonInfo, workflowDONs, nil)
		require.NoError(t, err)

		err = server.Start(ctx)
		require.NoError(t, err)

		// Clean up
		err = server.Close()
		require.NoError(t, err)
	})
}

func Test_Server_SetConfig_DONMembershipChange(t *testing.T) {
	ctx := testutils.Context(t)
	lggr := logger.Test(t)
	peerID := NewP2PPeerID(t)
	broker := newTestAsyncMessageBroker(t, 100)
	dispatcher := broker.NewDispatcherForNode(peerID)
	server := executable.NewServer("test-capability-id", "test-method", peerID, dispatcher, lggr)

	capInfo := commoncap.CapabilityInfo{
		ID:             "test-capability-id",
		CapabilityType: commoncap.CapabilityTypeTarget,
		Description:    "Test capability",
	}

	localDonInfo := commoncap.DON{
		ID:      1,
		Members: []p2ptypes.PeerID{peerID},
		F:       0,
	}

	workflowPeer1 := NewP2PPeerID(t)
	workflowPeer2 := NewP2PPeerID(t)
	workflowDONs := map[uint32]commoncap.DON{
		2: {
			ID:      2,
			Members: []p2ptypes.PeerID{workflowPeer1},
			F:       0,
		},
	}

	underlying := &TestSlowExecutionCapability{
		workflowIDToPause: map[string]time.Duration{
			workflowID1: 1 * time.Second,
		},
	}

	config := &commoncap.RemoteExecutableConfig{
		RequestTimeout:            10 * time.Second,
		ServerMaxParallelRequests: 5,
	}
	err := server.SetConfig(config, underlying, capInfo, localDonInfo, workflowDONs, nil)
	require.NoError(t, err)

	// Set up workflow node before starting servers
	workflowDispatcher := broker.NewDispatcherForNode(workflowPeer1)
	workflowNode := newServerTestClient(workflowPeer1, localDonInfo, workflowDispatcher)
	broker.RegisterReceiverNode(workflowPeer1, workflowNode)
	broker.RegisterReceiverNode(peerID, server)

	err = server.Start(ctx)
	require.NoError(t, err)
	err = broker.Start(ctx)
	require.NoError(t, err)

	// Start a request
	_, err = workflowNode.Execute(context.Background(), commoncap.CapabilityRequest{
		Metadata: commoncap.RequestMetadata{
			WorkflowID:          workflowID1,
			WorkflowExecutionID: workflowExecutionID1,
		},
	})
	require.NoError(t, err)

	// Change DON membership while request is in flight
	time.Sleep(100 * time.Millisecond)
	newWorkflowDONs := map[uint32]commoncap.DON{
		2: {
			ID:      2,
			Members: []p2ptypes.PeerID{workflowPeer1, workflowPeer2},
			F:       0,
		},
	}
	err = server.SetConfig(config, underlying, capInfo, localDonInfo, newWorkflowDONs, nil)
	require.NoError(t, err)

	// Original request should still complete
	select {
	case msg := <-workflowNode.receivedMessages:
		assert.NotNil(t, msg)
	case <-time.After(5 * time.Second):
		t.Fatal("request did not complete after DON change")
	}

	// Clean up
	require.NoError(t, server.Close())
	require.NoError(t, broker.Close())
}

func Test_Server_SetConfig_ShutdownRaces(t *testing.T) {
	ctx := testutils.Context(t)
	lggr := logger.Test(t)
	peerID := NewP2PPeerID(t)
	broker := newTestAsyncMessageBroker(t, 100)
	dispatcher := broker.NewDispatcherForNode(peerID)
	server := executable.NewServer("test-capability-id", "test-method", peerID, dispatcher, lggr)

	capInfo := commoncap.CapabilityInfo{
		ID:             "test-capability-id",
		CapabilityType: commoncap.CapabilityTypeTarget,
		Description:    "Test capability",
	}

	localDonInfo := commoncap.DON{
		ID:      1,
		Members: []p2ptypes.PeerID{peerID},
		F:       0,
	}

	workflowDONs := map[uint32]commoncap.DON{
		2: {
			ID:      2,
			Members: []p2ptypes.PeerID{NewP2PPeerID(t)},
			F:       0,
		},
	}

	underlying := &TestCapability{}
	config := &commoncap.RemoteExecutableConfig{
		RequestTimeout:            10 * time.Second,
		ServerMaxParallelRequests: 5,
	}

	err := server.SetConfig(config, underlying, capInfo, localDonInfo, workflowDONs, nil)
	require.NoError(t, err)
	err = server.Start(ctx)
	require.NoError(t, err)

	// Concurrently call SetConfig and Close
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := range 50 {
			newConfig := &commoncap.RemoteExecutableConfig{
				RequestTimeout:            time.Duration(5+i) * time.Millisecond,
				ServerMaxParallelRequests: 5,
			}
			_ = server.SetConfig(newConfig, underlying, capInfo, localDonInfo, workflowDONs, nil)
			time.Sleep(1 * time.Millisecond)
		}
	}()

	go func() {
		defer wg.Done()
		time.Sleep(25 * time.Millisecond)
		_ = server.Close()
	}()

	wg.Wait()
	// Test passes if no panic occurs
}

func Test_Server_Execute_WithConcurrentSetConfig(t *testing.T) {
	ctx := testutils.Context(t)
	lggr := logger.Test(t)
	numWorkflowPeers := 4

	peerID := NewP2PPeerID(t)
	capabilityPeers := []p2ptypes.PeerID{peerID}

	capDonInfo := commoncap.DON{
		ID:      1,
		Members: capabilityPeers,
		F:       0,
	}

	capInfo := commoncap.CapabilityInfo{
		ID:             "cap_id@1.0.0",
		CapabilityType: commoncap.CapabilityTypeTarget,
		Description:    "Remote Target",
		DON:            &capDonInfo,
	}

	workflowPeers := make([]p2ptypes.PeerID, numWorkflowPeers)
	for i := range numWorkflowPeers {
		workflowPeers[i] = NewP2PPeerID(t)
	}

	workflowDonInfo := commoncap.DON{
		Members: workflowPeers,
		ID:      2,
		F:       1,
	}

	broker := newTestAsyncMessageBroker(t, 1000)
	err := broker.Start(context.Background())
	require.NoError(t, err)
	defer broker.Close()

	workflowDONs := map[uint32]commoncap.DON{
		workflowDonInfo.ID: workflowDonInfo,
	}

	// Create and set up server
	dispatcher := broker.NewDispatcherForNode(peerID)
	server := executable.NewServer(capInfo.ID, "", peerID, dispatcher, lggr)

	underlying := &TestSlowExecutionCapability{
		workflowIDToPause: map[string]time.Duration{
			workflowID1: 50 * time.Millisecond,
		},
	}

	initialConfig := &commoncap.RemoteExecutableConfig{
		RequestTimeout:            10 * time.Second,
		ServerMaxParallelRequests: 10,
	}
	err = server.SetConfig(initialConfig, underlying, capInfo, capDonInfo, workflowDONs, nil)
	require.NoError(t, err)

	err = server.Start(ctx)
	require.NoError(t, err)
	defer server.Close()

	broker.RegisterReceiverNode(peerID, server)

	// Create workflow nodes (callers)
	workflowNodes := make([]*serverTestClient, numWorkflowPeers)
	for i := range numWorkflowPeers {
		workflowPeerDispatcher := broker.NewDispatcherForNode(workflowPeers[i])
		workflowNode := newServerTestClient(workflowPeers[i], capDonInfo, workflowPeerDispatcher)
		broker.RegisterReceiverNode(workflowPeers[i], workflowNode)
		workflowNodes[i] = workflowNode
	}

	var wg sync.WaitGroup
	numExecuteCalls := 20
	numSetConfigCalls := 10

	// Track successful responses
	responseCount := sync.Map{}

	// Start goroutine for concurrent SetConfig calls with randomized delays
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := range numSetConfigCalls {
			// Random delay between 5-50ms
			delay := time.Duration(5+i*2) * time.Millisecond
			time.Sleep(delay)

			newConfig := &commoncap.RemoteExecutableConfig{
				RequestTimeout:            time.Duration(10+i) * time.Second,
				ServerMaxParallelRequests: uint32(5),
			}
			assert.NoError(t, server.SetConfig(newConfig, underlying, capInfo, capDonInfo, workflowDONs, nil))
		}
	}()

	// Start multiple goroutines for concurrent Execute calls with randomized delays
	for callerIdx, caller := range workflowNodes {
		for execIdx := range numExecuteCalls {
			wg.Add(1)
			go func(callerID int, execID int, node *serverTestClient) {
				defer wg.Done()

				// Random delay between 0-100ms
				delay := time.Duration(execID*5) * time.Millisecond
				time.Sleep(delay)

				workflowExecutionID := fmt.Sprintf("exec-%d", execID)
				_, err := node.Execute(context.Background(),
					commoncap.CapabilityRequest{
						Metadata: commoncap.RequestMetadata{
							WorkflowID:          workflowID1,
							WorkflowExecutionID: workflowExecutionID,
						},
					})
				if err != nil {
					t.Logf("Execute error for caller %d exec %d: %v", callerID, execID, err)
				}
			}(callerIdx, execIdx, caller)
		}
	}

	// Collect responses
	wg.Add(1)
	go func() {
		defer wg.Done()
		expectedResponses := numWorkflowPeers * numExecuteCalls

		for i := range expectedResponses {
			// Try to receive from all callers
			for _, caller := range workflowNodes {
				select {
				case msg := <-caller.receivedMessages:
					if msg.Error == remotetypes.Error_OK {
						count, _ := responseCount.LoadOrStore("success", 0)
						responseCount.Store("success", count.(int)+1)
					} else {
						count, _ := responseCount.LoadOrStore("error", 0)
						responseCount.Store("error", count.(int)+1)
					}
				case <-time.After(15 * time.Second):
					t.Logf("Timeout waiting for response %d/%d", i+1, expectedResponses)
					return
				}
			}
		}
	}()

	wg.Wait()

	// Verify we received responses (most should succeed)
	successCount := 0
	if val, ok := responseCount.Load("success"); ok {
		successCount = val.(int)
	}
	expectedResponses := numWorkflowPeers * numExecuteCalls
	require.Equal(t, expectedResponses, successCount)
}
