package capabilities

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/durationpb"

	ragetypes "github.com/smartcontractkit/libocr/ragep2p/types"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"

	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink-protos/cre/go/values"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote"
	remotetypes "github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/types"
	remoteMocks "github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/types/mocks"
	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"
	"github.com/smartcontractkit/chainlink/v2/core/services/p2p/types/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/services/registrysyncer"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

var _ capabilities.TriggerCapability = (*mockTrigger)(nil)

type mockDonNotifier struct{}

func (m *mockDonNotifier) NotifyDonSet(don capabilities.DON) {
}

type mockTrigger struct {
	capabilities.CapabilityInfo
}

func (m *mockTrigger) RegisterTrigger(ctx context.Context, request capabilities.TriggerRegistrationRequest) (<-chan capabilities.TriggerResponse, error) {
	return nil, nil
}

func (m *mockTrigger) UnregisterTrigger(ctx context.Context, request capabilities.TriggerRegistrationRequest) error {
	return nil
}

func newMockTrigger(info capabilities.CapabilityInfo) *mockTrigger {
	return &mockTrigger{CapabilityInfo: info}
}

var _ capabilities.ExecutableCapability = (*mockCapability)(nil)

type mockCapability struct {
	capabilities.CapabilityInfo
}

func (m *mockCapability) Execute(ctx context.Context, req capabilities.CapabilityRequest) (capabilities.CapabilityResponse, error) {
	return capabilities.CapabilityResponse{}, nil
}

func (m *mockCapability) RegisterToWorkflow(ctx context.Context, request capabilities.RegisterToWorkflowRequest) error {
	return nil
}

func (m *mockCapability) UnregisterFromWorkflow(ctx context.Context, request capabilities.UnregisterFromWorkflowRequest) error {
	return nil
}

func TestLauncher(t *testing.T) {
	t.Run("OK-wires_up_external_capabilities", func(t *testing.T) {
		lggr := logger.Test(t)
		registry := NewRegistry(lggr)
		dispatcher := remoteMocks.NewDispatcher(t)

		nodes := newNodes(4)
		capabilityDonNodes := newNodes(4)
		peer := mocks.NewPeer(t)
		peer.On("UpdateConnections", mock.Anything).Return(nil)
		peer.On("ID").Return(capabilityDonNodes[0])
		peer.On("IsBootstrap").Return(false)
		wrapper := mocks.NewPeerWrapper(t)
		wrapper.On("GetPeer").Return(peer)

		fullTriggerCapID := "streams-trigger@1.0.0"
		mt := newMockTrigger(capabilities.MustNewCapabilityInfo(
			fullTriggerCapID,
			capabilities.CapabilityTypeTrigger,
			"streams trigger",
		))
		require.NoError(t, registry.Add(t.Context(), mt))

		fullTargetID := "write-chain_evm_1@1.0.0"
		mtarg := &mockCapability{
			CapabilityInfo: capabilities.MustNewCapabilityInfo(
				fullTargetID,
				capabilities.CapabilityTypeTarget,
				"write chain",
			),
		}
		require.NoError(t, registry.Add(t.Context(), mtarg))

		triggerCapID := RandomUTF8BytesWord()
		targetCapID := RandomUTF8BytesWord()
		// one capability from onchain registry is not set up locally
		fullMissingTargetID := "super-duper-target@6.6.6"
		missingTargetCapID := RandomUTF8BytesWord()
		dID := uint32(1)
		capDonID := uint32(2)

		localRegistry := buildLocalRegistry()
		addDON(localRegistry, dID, uint32(0), uint8(1), true, true, nodes, []string{"zone-a"}, 1, [][32]byte{triggerCapID, targetCapID, missingTargetCapID})
		addDON(localRegistry, capDonID, uint32(0), uint8(1), true, false, capabilityDonNodes, []string{"zone-a"}, 1, [][32]byte{triggerCapID, targetCapID})
		addCapabilityToDON(localRegistry, capDonID, fullTriggerCapID, capabilities.CapabilityTypeTrigger, nil)
		addCapabilityToDON(localRegistry, capDonID, fullTargetID, capabilities.CapabilityTypeTarget, nil)
		addCapabilityToDON(localRegistry, capDonID, fullMissingTargetID, capabilities.CapabilityTypeTarget, nil)

		launcher := NewLauncher(
			lggr,
			wrapper,
			nil,
			nil,
			dispatcher,
			registry,
			&mockDonNotifier{},
		)
		require.NoError(t, launcher.Start(t.Context()))
		defer launcher.Close()

		dispatcher.On("SetReceiver", fullTriggerCapID, capDonID, mock.AnythingOfType("*remote.triggerPublisher")).Return(nil)
		dispatcher.On("SetReceiver", fullTargetID, capDonID, mock.AnythingOfType("*executable.server")).Return(nil)

		require.NoError(t, launcher.OnNewRegistry(t.Context(), localRegistry))
	})

	t.Run("NOK-invalid_trigger_capability", func(t *testing.T) {
		lggr, observedLogs := logger.TestObserved(t, zapcore.DebugLevel)
		registry := NewRegistry(lggr)
		dispatcher := remoteMocks.NewDispatcher(t)

		nodes := newNodes(4)
		peer := mocks.NewPeer(t)
		peer.On("UpdateConnections", mock.Anything).Return(nil)
		peer.On("ID").Return(nodes[0])
		peer.On("IsBootstrap").Return(false)
		wrapper := mocks.NewPeerWrapper(t)
		wrapper.On("GetPeer").Return(peer)

		// We intentionally create a Trigger capability with a Target type
		fullTriggerCapID := "streams-trigger@1.0.0"
		mtarg := &mockCapability{
			CapabilityInfo: capabilities.MustNewCapabilityInfo(
				fullTriggerCapID,
				capabilities.CapabilityTypeTarget, // intentionally wrong type
				"wrong type capability",
			),
		}
		require.NoError(t, registry.Add(t.Context(), mtarg))
		triggerCapID := RandomUTF8BytesWord()

		dID := uint32(1)
		localRegistry := buildLocalRegistry()
		addDON(localRegistry, dID, uint32(0), uint8(1), true, true, nodes, []string{"zone-a"}, 1, [][32]byte{triggerCapID})
		addCapabilityToDON(localRegistry, dID, fullTriggerCapID, capabilities.CapabilityTypeTrigger, nil)

		launcher := NewLauncher(
			lggr,
			wrapper,
			nil,
			nil,
			dispatcher,
			registry,
			&mockDonNotifier{},
		)
		require.NoError(t, launcher.Start(t.Context()))
		defer launcher.Close()

		require.NoError(t, launcher.OnNewRegistry(t.Context(), localRegistry))
		assert.Equal(t, 1, observedLogs.FilterMessage("failed to serve capability").Len())
	})

	t.Run("NOK-invalid_target_capability", func(t *testing.T) {
		lggr, observedLogs := logger.TestObserved(t, zapcore.DebugLevel)
		registry := NewRegistry(lggr)
		dispatcher := remoteMocks.NewDispatcher(t)

		nodes := newNodes(4)
		peer := mocks.NewPeer(t)
		peer.On("UpdateConnections", mock.Anything).Return(nil)
		peer.On("ID").Return(nodes[0])
		peer.On("IsBootstrap").Return(false)
		wrapper := mocks.NewPeerWrapper(t)
		wrapper.On("GetPeer").Return(peer)

		fullTargetID := "write-chain_evm_1@1.0.0"
		mt := newMockTrigger(capabilities.MustNewCapabilityInfo(
			fullTargetID,
			capabilities.CapabilityTypeTrigger, // intentionally wrong type
			"streams trigger",
		))
		require.NoError(t, registry.Add(t.Context(), mt))

		targetCapID := RandomUTF8BytesWord()
		dID := uint32(1)
		localRegistry := buildLocalRegistry()
		addDON(localRegistry, dID, uint32(0), uint8(1), true, true, nodes, []string{"zone-a"}, 1, [][32]byte{targetCapID})
		addCapabilityToDON(localRegistry, dID, fullTargetID, capabilities.CapabilityTypeTarget, nil)

		launcher := NewLauncher(
			lggr,
			wrapper,
			nil,
			nil,
			dispatcher,
			registry,
			&mockDonNotifier{},
		)
		require.NoError(t, launcher.Start(t.Context()))
		defer launcher.Close()

		require.NoError(t, launcher.OnNewRegistry(t.Context(), localRegistry))
		assert.Equal(t, 1, observedLogs.FilterMessage("failed to serve capability").Len())
	})

	t.Run("start and close with nil peer wrapper", func(t *testing.T) {
		lggr := logger.Test(t)
		registry := NewRegistry(lggr)
		dispatcher := remoteMocks.NewDispatcher(t)
		sharedPeer := mocks.NewSharedPeer(t)
		sharedPeer.On("ID").Return(ragetypes.PeerID(RandomUTF8BytesWord()))
		launcher := NewLauncher(
			lggr,
			nil,
			sharedPeer,
			nil,
			dispatcher,
			registry,
			&mockDonNotifier{},
		)
		require.NoError(t, launcher.Start(t.Context()))
		require.NoError(t, launcher.Close())
	})
}

func newTriggerEventMsg(t *testing.T,
	senderPeerID p2ptypes.PeerID,
	workflowID string,
	triggerEvent map[string]any,
	triggerEventID string,
) (*remotetypes.MessageBody, *values.Map) {
	triggerEventValue, err := values.NewMap(triggerEvent)
	require.NoError(t, err)
	capResponse := capabilities.TriggerResponse{
		Event: capabilities.TriggerEvent{
			Outputs: triggerEventValue,
			ID:      triggerEventID,
		},
		Err: nil,
	}
	marshaled, err := capabilitiespb.MarshalTriggerResponse(capResponse)
	require.NoError(t, err)
	return &remotetypes.MessageBody{
		Sender: senderPeerID[:],
		Method: remotetypes.MethodTriggerEvent,
		Metadata: &remotetypes.MessageBody_TriggerEventMetadata{
			TriggerEventMetadata: &remotetypes.TriggerEventMetadata{
				WorkflowIds: []string{workflowID},
			},
		},
		Payload: marshaled,
	}, triggerEventValue
}

func TestLauncher_RemoteTriggerModeAggregatorShim(t *testing.T) {
	ctx := t.Context()
	lggr := logger.Test(t)
	registry := NewRegistry(lggr)
	dispatcher := remoteMocks.NewDispatcher(t)

	workflowDonNodes, capabilityDonNodes := newNodes(4), newNodes(4)
	peer := mocks.NewPeer(t)
	peer.On("UpdateConnections", mock.Anything).Return(nil)
	peer.On("ID").Return(workflowDonNodes[0])
	peer.On("IsBootstrap").Return(false)
	wrapper := mocks.NewPeerWrapper(t)
	wrapper.On("GetPeer").Return(peer)

	fullTriggerCapID := "log-event-trigger-evm-43113@1.0.0"
	fullTargetID := "write-chain_evm_1@1.0.0"
	triggerCapID := RandomUTF8BytesWord()
	targetCapID := RandomUTF8BytesWord()
	dID := uint32(1)
	capDonID := uint32(2)

	cfg, err := proto.Marshal(&capabilitiespb.CapabilityConfig{
		RemoteConfig: &capabilitiespb.CapabilityConfig_RemoteTriggerConfig{
			RemoteTriggerConfig: &capabilitiespb.RemoteTriggerConfig{
				RegistrationRefresh:     durationpb.New(1 * time.Second),
				MinResponsesToAggregate: 3,
			},
		},
	})
	require.NoError(t, err)

	localRegistry := buildLocalRegistry()
	addDON(localRegistry, dID, uint32(0), uint8(1), true, true, workflowDonNodes, []string{"zone-a"}, 1, nil)
	addDON(localRegistry, capDonID, uint32(0), uint8(1), true, false, capabilityDonNodes, []string{"zone-a"}, 1, [][32]byte{triggerCapID, targetCapID})
	addCapabilityToDON(localRegistry, capDonID, fullTriggerCapID, capabilities.CapabilityTypeTrigger, cfg)
	addCapabilityToDON(localRegistry, capDonID, fullTargetID, capabilities.CapabilityTypeTarget, cfg)

	launcher := NewLauncher(
		lggr,
		wrapper,
		nil,
		nil,
		dispatcher,
		registry,
		&mockDonNotifier{},
	)
	require.NoError(t, launcher.Start(t.Context()))
	defer launcher.Close()

	dispatcher.On("SetReceiver", fullTriggerCapID, capDonID, mock.AnythingOfType("*remote.triggerSubscriber")).Return(nil)
	dispatcher.On("SetReceiver", fullTargetID, capDonID, mock.AnythingOfType("*executable.client")).Return(nil)
	dispatcher.On("Ready").Return(nil).Maybe()
	awaitRegistrationMessageCh := make(chan struct{})
	dispatcher.On("Send", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		select {
		case awaitRegistrationMessageCh <- struct{}{}:
		default:
		}
	})

	err = launcher.OnNewRegistry(ctx, localRegistry)
	require.NoError(t, err)

	baseCapability, err := registry.Get(ctx, fullTriggerCapID)
	require.NoError(t, err)

	remoteTriggerSubscriber, ok := baseCapability.(remote.TriggerSubscriber)
	require.True(t, ok, "remote trigger capability")

	// Register trigger
	workflowID1 := "15c631d295ef5e32deb99a10ee6804bc4af13855687559d7ff6552ac6dbb2ce0"
	workflowExecutionID1 := "95ef5e32deb99a10ee6804bc4af13855687559d7ff6552ac6dbb2ce0abbadeed"
	req := capabilities.TriggerRegistrationRequest{
		TriggerID: "logeventtrigger_log1",
		Metadata: capabilities.RequestMetadata{
			ReferenceID:         "logeventtrigger",
			WorkflowID:          workflowID1,
			WorkflowExecutionID: workflowExecutionID1,
		},
	}
	triggerEventCallbackCh, err := remoteTriggerSubscriber.RegisterTrigger(ctx, req)
	require.NoError(t, err)
	<-awaitRegistrationMessageCh

	// Receive trigger event
	triggerEvent1 := map[string]any{"event": "triggerEvent1"}
	triggerEvent2 := map[string]any{"event": "triggerEvent2"}
	triggerEventMsg1, triggerEventValue := newTriggerEventMsg(t, capabilityDonNodes[0], workflowID1, triggerEvent1, "TriggerEventID1")
	triggerEventMsg2, _ := newTriggerEventMsg(t, capabilityDonNodes[1], workflowID1, triggerEvent1, "TriggerEventID1")
	// One Faulty Node (F = 1) sending bad event data for the same TriggerEventID1
	triggerEventMsg3, _ := newTriggerEventMsg(t, capabilityDonNodes[2], workflowID1, triggerEvent2, "TriggerEventID1")
	remoteTriggerSubscriber.Receive(ctx, triggerEventMsg1)
	remoteTriggerSubscriber.Receive(ctx, triggerEventMsg2)
	remoteTriggerSubscriber.Receive(ctx, triggerEventMsg3)

	// After MinResponsesToAggregate, we should get a response
	response := <-triggerEventCallbackCh

	// Checks if response is same as minIdenticalResponses = F + 1, F = 1
	require.Equal(t, response.Event.Outputs, triggerEventValue)
}

func TestSyncer_IgnoresCapabilitiesForPrivateDON(t *testing.T) {
	lggr := logger.Test(t)
	registry := NewRegistry(lggr)
	dispatcher := remoteMocks.NewDispatcher(t)

	nodes := newNodes(4)
	peer := mocks.NewPeer(t)
	peer.On("UpdateConnections", mock.Anything).Return(nil)
	peer.On("ID").Return(nodes[0])
	peer.On("IsBootstrap").Return(false)
	wrapper := mocks.NewPeerWrapper(t)
	wrapper.On("GetPeer").Return(peer)

	dID := uint32(1)
	triggerID := "streams-trigger@1.0.0"
	hashedTriggerID := RandomUTF8BytesWord()
	targetID := "write-chain_evm_1@1.0.0"
	hashedTargetID := RandomUTF8BytesWord()

	localRegistry := buildLocalRegistry()
	addDON(localRegistry, dID, uint32(0), uint8(1), false, true, nodes, []string{"zone-a"}, 1, [][32]byte{hashedTriggerID, hashedTargetID})
	addCapabilityToDON(localRegistry, dID, triggerID, capabilities.CapabilityTypeTrigger, nil)
	addCapabilityToDON(localRegistry, dID, targetID, capabilities.CapabilityTypeTarget, nil)

	launcher := NewLauncher(
		lggr,
		wrapper,
		nil,
		nil,
		dispatcher,
		registry,
		&mockDonNotifier{},
	)
	require.NoError(t, launcher.Start(t.Context()))
	defer launcher.Close()

	// If the DON were public, this would fail with two errors:
	// - error fetching the capabilities from the registry since they haven't been added
	// - erroneous calls to dispatcher.SetReceiver, since the call hasn't been registered.
	err := launcher.OnNewRegistry(t.Context(), localRegistry)
	require.NoError(t, err)

	// Finally, assert that no services were added.
	assert.Empty(t, launcher.subServices)
}

func TestLauncher_WiresUpClientsForPublicWorkflowDON(t *testing.T) {
	lggr := logger.Test(t)
	registry := NewRegistry(lggr)
	dispatcher := remoteMocks.NewDispatcher(t)

	workflowDonNodes, capabilityDonNodes := newNodes(4), newNodes(4)
	peer := mocks.NewPeer(t)
	peer.On("UpdateConnections", mock.Anything).Return(nil)
	peer.On("ID").Return(workflowDonNodes[0])
	peer.On("IsBootstrap").Return(false)
	wrapper := mocks.NewPeerWrapper(t)
	wrapper.On("GetPeer").Return(peer)

	fullTriggerCapID := "streams-trigger@1.0.0"
	fullTargetID := "write-chain_evm_1@1.0.0"
	triggerCapID := RandomUTF8BytesWord()
	targetCapID := RandomUTF8BytesWord()
	dID := uint32(1)
	capDonID := uint32(2)

	cfg, err := proto.Marshal(&capabilitiespb.CapabilityConfig{
		RemoteConfig: &capabilitiespb.CapabilityConfig_RemoteTriggerConfig{
			RemoteTriggerConfig: &capabilitiespb.RemoteTriggerConfig{
				RegistrationRefresh: durationpb.New(1 * time.Second),
			},
		},
	})
	require.NoError(t, err)

	localRegistry := buildLocalRegistry()
	addDON(localRegistry, dID, uint32(0), uint8(1), true, true, workflowDonNodes, []string{"zone-a"}, 1, nil)
	addDON(localRegistry, capDonID, uint32(0), uint8(1), true, false, capabilityDonNodes, []string{"zone-a"}, 1, [][32]byte{triggerCapID, targetCapID})
	addCapabilityToDON(localRegistry, capDonID, fullTriggerCapID, capabilities.CapabilityTypeTrigger, cfg)
	addCapabilityToDON(localRegistry, capDonID, fullTargetID, capabilities.CapabilityTypeTarget, cfg)

	launcher := NewLauncher(
		lggr,
		wrapper,
		nil,
		nil,
		dispatcher,
		registry,
		&mockDonNotifier{},
	)
	require.NoError(t, launcher.Start(t.Context()))
	defer launcher.Close()

	dispatcher.On("SetReceiver", fullTriggerCapID, capDonID, mock.AnythingOfType("*remote.triggerSubscriber")).Return(nil)
	dispatcher.On("SetReceiver", fullTargetID, capDonID, mock.AnythingOfType("*executable.client")).Return(nil)

	err = launcher.OnNewRegistry(t.Context(), localRegistry)
	require.NoError(t, err)

	_, err = registry.Get(t.Context(), fullTriggerCapID)
	require.NoError(t, err)

	_, err = registry.Get(t.Context(), fullTargetID)
	require.NoError(t, err)
}

func TestLauncher_WiresUpClientsForPublicWorkflowDONButIgnoresPrivateCapabilities(t *testing.T) {
	lggr := logger.Test(t)
	registry := NewRegistry(lggr)
	dispatcher := remoteMocks.NewDispatcher(t)

	workflowDonNodes, capabilityDonNodes := newNodes(4), newNodes(4)
	peer := mocks.NewPeer(t)
	peer.On("UpdateConnections", mock.Anything).Return(nil)
	peer.On("ID").Return(workflowDonNodes[0])
	peer.On("IsBootstrap").Return(false)
	wrapper := mocks.NewPeerWrapper(t)
	wrapper.On("GetPeer").Return(peer)

	fullTriggerCapID := "streams-trigger@1.0.0"
	fullTargetID := "write-chain_evm_1@1.0.0"
	triggerCapID := RandomUTF8BytesWord()
	targetCapID := RandomUTF8BytesWord()
	dID := uint32(1)
	triggerCapDonID := uint32(2)
	targetCapDonID := uint32(3)

	cfg, err := proto.Marshal(&capabilitiespb.CapabilityConfig{
		RemoteConfig: &capabilitiespb.CapabilityConfig_RemoteTriggerConfig{
			RemoteTriggerConfig: &capabilitiespb.RemoteTriggerConfig{},
		},
	})
	require.NoError(t, err)

	localRegistry := buildLocalRegistry()
	addDON(localRegistry, dID, uint32(0), uint8(1), true, true, workflowDonNodes, []string{"zone-a"}, 1, nil)
	addDON(localRegistry, triggerCapDonID, uint32(0), uint8(1), true, false, capabilityDonNodes, []string{"zone-a"}, 1, [][32]byte{triggerCapID, targetCapID})
	addCapabilityToDON(localRegistry, triggerCapDonID, fullTriggerCapID, capabilities.CapabilityTypeTrigger, cfg)
	addDON(localRegistry, targetCapDonID, uint32(0), uint8(1), false, false, capabilityDonNodes, []string{"zone-a"}, 1, [][32]byte{triggerCapID, targetCapID})
	addCapabilityToDON(localRegistry, targetCapDonID, fullTargetID, capabilities.CapabilityTypeTarget, cfg)

	launcher := NewLauncher(
		lggr,
		wrapper,
		nil,
		nil,
		dispatcher,
		registry,
		&mockDonNotifier{},
	)
	require.NoError(t, launcher.Start(t.Context()))
	defer launcher.Close()
	dispatcher.On("SetReceiver", fullTriggerCapID, triggerCapDonID, mock.AnythingOfType("*remote.triggerSubscriber")).Return(nil)

	require.NoError(t, launcher.OnNewRegistry(t.Context(), localRegistry))

	_, err = registry.Get(t.Context(), fullTriggerCapID)
	require.NoError(t, err)
}

func TestLauncher_SucceedsEvenIfDispatcherAlreadyHasReceiver(t *testing.T) {
	lggr := logger.Test(t)
	registry := NewRegistry(lggr)
	dispatcher := remoteMocks.NewDispatcher(t)

	var pid ragetypes.PeerID
	err := pid.UnmarshalText([]byte("12D3KooWBCF1XT5Wi8FzfgNCqRL76Swv8TRU3TiD4QiJm8NMNX7N"))
	require.NoError(t, err)

	peer := mocks.NewPeer(t)
	peer.On("UpdateConnections", mock.Anything).Return(nil)
	peer.On("ID").Return(pid)
	peer.On("IsBootstrap").Return(false)

	wrapper := mocks.NewPeerWrapper(t)
	wrapper.On("GetPeer").Return(peer)

	fullTriggerCapID := "streams-trigger@1.0.0"
	mt := newMockTrigger(capabilities.MustNewCapabilityInfo(
		fullTriggerCapID,
		capabilities.CapabilityTypeTrigger,
		"streams trigger",
	))
	require.NoError(t, registry.Add(t.Context(), mt))

	tt := NewTestTopology(pid, 4, 4)

	triggerCapID := RandomUTF8BytesWord()
	workflowDONID := uint32(1)
	capabilitiesDONID := uint32(2)
	workflowNCapabilitiesDONID := uint32(3)

	// The below state describes a Capability DON (AcceptsWorkflows = true),
	// which exposes the streams-trigger and write_chain capabilities.
	// We expect receivers to be wired up.
	localRegistry := tt.MakeLocalRegistry(
		workflowDONID,
		capabilitiesDONID,
		workflowNCapabilitiesDONID,
		triggerCapID,
		fullTriggerCapID,
	)

	dispatcher.On(
		"SetReceiver",
		fullTriggerCapID,
		capabilitiesDONID,
		mock.AnythingOfType("*remote.triggerPublisher"),
	).Return(remote.ErrReceiverExists)

	launcher := NewLauncher(
		lggr,
		wrapper,
		nil,
		nil,
		dispatcher,
		registry,
		&mockDonNotifier{},
	)
	require.NoError(t, launcher.Start(t.Context()))
	defer launcher.Close()
	err = launcher.OnNewRegistry(t.Context(), localRegistry)
	require.NoError(t, err)
}

func TestLauncher_SuccessfullyFilterDon2Don(t *testing.T) {
	lggr := logger.Test(t)
	registry := NewRegistry(lggr)
	dispatcher := remoteMocks.NewDispatcher(t)

	var pid ragetypes.PeerID
	err := pid.UnmarshalText([]byte("12D3KooWBCF1XT5Wi8FzfgNCqRL76Swv8TRU3TiD4QiJm8NMNX7N"))
	require.NoError(t, err)

	peer := mocks.NewPeer(t)
	peer.On("UpdateConnections", mock.Anything).Return(nil)
	peer.On("ID").Return(pid)
	peer.On("IsBootstrap").Return(false)

	wrapper := mocks.NewPeerWrapper(t)
	wrapper.On("GetPeer").Return(peer)

	fullTriggerCapID := "streams-trigger@1.0.0"
	mt := newMockTrigger(capabilities.MustNewCapabilityInfo(
		fullTriggerCapID,
		capabilities.CapabilityTypeTrigger,
		"streams trigger",
	))
	require.NoError(t, registry.Add(t.Context(), mt))

	tt := NewTestTopology(pid, 4, 4)

	triggerCapID := RandomUTF8BytesWord()
	workflowDONID := uint32(1)
	capabilitiesDONID := uint32(2)
	workflowNCapabilitiesDONID := uint32(3)

	localRegistry := tt.MakeLocalRegistry(
		workflowDONID,
		capabilitiesDONID,
		workflowNCapabilitiesDONID,
		triggerCapID,
		fullTriggerCapID,
	)

	dispatcher.On(
		"SetReceiver",
		fullTriggerCapID,
		capabilitiesDONID,
		mock.AnythingOfType("*remote.triggerPublisher"),
	).Return(remote.ErrReceiverExists)

	launcher := NewLauncher(
		lggr,
		wrapper,
		nil,
		nil,
		dispatcher,
		registry,
		&mockDonNotifier{},
	)
	require.NoError(t, launcher.Start(t.Context()))
	defer launcher.Close()

	inputs := [][]bool{
		// { belongsToACapabilityDON, belongsToAWorkflowDON, isBootstrap }
		{true, true, false},
		{true, false, false},
		{false, true, false},
		{false, false, false},
		{false, false, true},
		{true, true, true}, // invalid
	}

	expectedPeerCount := []int{
		8, // we expect all DONs members
		5, // we expect all capability DONs members (4+1)
		4, // we expect all workflow DONs members
		0, // the node does nothing, we expect no peers
		8, // bootstrap node always adds all peers
		8, // bootstrap node always adds all peers
	}

	for i := range inputs {
		allPeers := launcher.peers(
			inputs[i][0],
			inputs[i][1],
			inputs[i][2],
			localRegistry,
		)
		require.Len(t, allPeers, expectedPeerCount[i])
	}

	err = launcher.OnNewRegistry(t.Context(), localRegistry)
	require.NoError(t, err)
}

func TestLauncher_DonPairsToUpdate(t *testing.T) {
	registry := NewRegistry(logger.Test(t))
	dispatcher := remoteMocks.NewDispatcher(t)

	var pid, other ragetypes.PeerID
	require.NoError(t, pid.UnmarshalText([]byte(utils.MustNewPeerID())))
	require.NoError(t, other.UnmarshalText([]byte(utils.MustNewPeerID())))
	sharedPeer := mocks.NewSharedPeer(t)

	fullTriggerCapID := "streams-trigger@1.0.0"
	mt := newMockTrigger(capabilities.MustNewCapabilityInfo(
		fullTriggerCapID,
		capabilities.CapabilityTypeTrigger,
		"streams trigger",
	))
	require.NoError(t, registry.Add(t.Context(), mt))

	tt := NewTestTopology(pid, 4, 4)
	wfDONID, capDONID, mixedDONID := registrysyncer.DonID(7), registrysyncer.DonID(12), registrysyncer.DonID(33)
	localRegistry := tt.MakeLocalRegistry(uint32(wfDONID), uint32(capDONID), uint32(mixedDONID), RandomUTF8BytesWord(), fullTriggerCapID)
	launcher := NewLauncher(logger.Test(t), nil, sharedPeer, nil, dispatcher, registry, &mockDonNotifier{})

	sharedPeer.On("IsBootstrap").Return(false).Times(3)
	// capability DON connects to DONs: workflow and mixed
	res := launcher.donPairsToUpdate(tt.capabilityDonNodes[0], localRegistry)
	require.Len(t, res, 2)
	require.Equal(t, p2ptypes.DonPair{localRegistry.IDsToDONs[wfDONID].DON, localRegistry.IDsToDONs[capDONID].DON}, res[0])
	require.Equal(t, p2ptypes.DonPair{localRegistry.IDsToDONs[capDONID].DON, localRegistry.IDsToDONs[mixedDONID].DON}, res[1])

	// workflow DON connects to DONs: capability and mixed
	res = launcher.donPairsToUpdate(tt.workflowDonNodes[0], localRegistry)
	require.Len(t, res, 2)
	require.Equal(t, p2ptypes.DonPair{localRegistry.IDsToDONs[wfDONID].DON, localRegistry.IDsToDONs[capDONID].DON}, res[0])
	require.Equal(t, p2ptypes.DonPair{localRegistry.IDsToDONs[wfDONID].DON, localRegistry.IDsToDONs[mixedDONID].DON}, res[1])

	// peer (not bootstrap) that doesn't belong to any DON connects to nobody
	require.Empty(t, launcher.donPairsToUpdate(other, localRegistry))

	// bootstrap node adds all 3 DON pairs
	sharedPeer.On("IsBootstrap").Return(true).Once()
	res = launcher.donPairsToUpdate(pid, localRegistry)
	require.Len(t, res, 3)
	require.Equal(t, p2ptypes.DonPair{localRegistry.IDsToDONs[wfDONID].DON, localRegistry.IDsToDONs[capDONID].DON}, res[0])
	require.Equal(t, p2ptypes.DonPair{localRegistry.IDsToDONs[wfDONID].DON, localRegistry.IDsToDONs[mixedDONID].DON}, res[1])
	require.Equal(t, p2ptypes.DonPair{localRegistry.IDsToDONs[capDONID].DON, localRegistry.IDsToDONs[mixedDONID].DON}, res[2])

	// bootstrap node adds only allowed DON pairs
	mixedDON := localRegistry.IDsToDONs[mixedDONID]
	mixedDON.AcceptsWorkflows = false
	localRegistry.IDsToDONs[mixedDONID] = mixedDON
	sharedPeer.On("IsBootstrap").Return(true).Once()
	res = launcher.donPairsToUpdate(pid, localRegistry)
	require.Len(t, res, 2)
	require.Equal(t, p2ptypes.DonPair{localRegistry.IDsToDONs[wfDONID].DON, localRegistry.IDsToDONs[capDONID].DON}, res[0])
	require.Equal(t, p2ptypes.DonPair{localRegistry.IDsToDONs[wfDONID].DON, localRegistry.IDsToDONs[mixedDONID].DON}, res[1])
}

func TestLauncher_DonPairsToUpdate_SkipsDifferentFamilies(t *testing.T) {
	registry := NewRegistry(logger.Test(t))
	dispatcher := remoteMocks.NewDispatcher(t)

	var pid ragetypes.PeerID
	require.NoError(t, pid.UnmarshalText([]byte(utils.MustNewPeerID())))
	sharedPeer := mocks.NewSharedPeer(t)

	fullTriggerCapID := "streams-trigger@1.0.0"
	mt := newMockTrigger(capabilities.MustNewCapabilityInfo(
		fullTriggerCapID,
		capabilities.CapabilityTypeTrigger,
		"streams trigger",
	))
	require.NoError(t, registry.Add(t.Context(), mt))

	// Create DONs with different families
	workflowDonNodes := newNodes(4)
	capabilityDonNodesZoneA := newNodes(4)
	capabilityDonNodesZoneB := newNodes(4)
	workflowDonNodes[0] = pid // node belongs to workflow DON

	wfDONID := uint32(1)
	capDONZoneAID := uint32(2)
	capDONZoneBID := uint32(3)

	triggerCapID := RandomUTF8BytesWord()
	localRegistry := buildLocalRegistry()

	// Workflow DON in zone-a
	addDON(localRegistry, wfDONID, uint32(0), uint8(1), true, true, workflowDonNodes, []string{"zone-a"}, 1, nil)
	// Capability DON in zone-a (should be included in pairs)
	addDON(localRegistry, capDONZoneAID, uint32(0), uint8(1), true, false, capabilityDonNodesZoneA, []string{"zone-a"}, 1, [][32]byte{triggerCapID})
	addCapabilityToDON(localRegistry, capDONZoneAID, fullTriggerCapID, capabilities.CapabilityTypeTrigger, nil)
	// Capability DON in zone-b (should be filtered out due to family mismatch)
	addDON(localRegistry, capDONZoneBID, uint32(0), uint8(1), true, false, capabilityDonNodesZoneB, []string{"zone-b"}, 1, [][32]byte{triggerCapID})
	addCapabilityToDON(localRegistry, capDONZoneBID, fullTriggerCapID, capabilities.CapabilityTypeTrigger, nil)

	launcher := NewLauncher(logger.Test(t), nil, sharedPeer, nil, dispatcher, registry, &mockDonNotifier{})

	sharedPeer.On("IsBootstrap").Return(false).Once()
	// Node belongs to workflow DON, should only connect to capability DON in same family (zone-a)
	res := launcher.donPairsToUpdate(pid, localRegistry)
	require.Len(t, res, 1, "expected only one DON pair (zone-a workflow to zone-a capability)")
	require.Equal(t, p2ptypes.DonPair{localRegistry.IDsToDONs[registrysyncer.DonID(wfDONID)].DON, localRegistry.IDsToDONs[registrysyncer.DonID(capDONZoneAID)].DON}, res[0])

	// Bootstrap node should still respect family boundaries
	sharedPeer.On("IsBootstrap").Return(true).Once()
	res = launcher.donPairsToUpdate(pid, localRegistry)
	require.Len(t, res, 1, "bootstrap should also filter based on families")
	require.Equal(t, p2ptypes.DonPair{localRegistry.IDsToDONs[registrysyncer.DonID(wfDONID)].DON, localRegistry.IDsToDONs[registrysyncer.DonID(capDONZoneAID)].DON}, res[0])
}

func TestLauncher_V2CapabilitiesAddViaCombinedClient(t *testing.T) {
	lggr := logger.Test(t)
	registry := NewRegistry(lggr)
	dispatcher := remoteMocks.NewDispatcher(t)

	workflowDonNodes, capabilityDonNodes, zoneBDonNodes := newNodes(4), newNodes(4), newNodes(4)
	fullTriggerCapID := "streams-trigger@1.0.0"
	fullExecutableCapID := "evm@1.0.0"
	triggerCapID := RandomUTF8BytesWord()
	executableCapID := RandomUTF8BytesWord()
	wfDonID := uint32(1)
	capDonID := uint32(2)
	zoneBDonID := uint32(4)

	triggerCfg, err := proto.Marshal(&capabilitiespb.CapabilityConfig{
		MethodConfigs: map[string]*capabilitiespb.CapabilityMethodConfig{
			"StreamsTrigger": {
				RemoteConfig: &capabilitiespb.CapabilityMethodConfig_RemoteTriggerConfig{
					RemoteTriggerConfig: &capabilitiespb.RemoteTriggerConfig{
						RegistrationRefresh:     durationpb.New(1 * time.Second),
						MinResponsesToAggregate: 3,
					},
				},
			},
		},
	})
	require.NoError(t, err)

	execCfg, err := proto.Marshal(&capabilitiespb.CapabilityConfig{
		MethodConfigs: map[string]*capabilitiespb.CapabilityMethodConfig{
			"Write": {
				RemoteConfig: &capabilitiespb.CapabilityMethodConfig_RemoteExecutableConfig{
					RemoteExecutableConfig: &capabilitiespb.RemoteExecutableConfig{
						RequestTimeout: durationpb.New(30 * time.Second),
						DeltaStage:     durationpb.New(1 * time.Second),
					},
				},
			},
		},
	})
	require.NoError(t, err)

	localRegistry := buildLocalRegistry()
	addDON(localRegistry, wfDonID, 0, 1, true, true, workflowDonNodes, []string{"zone-a"}, 1, nil)
	addDON(localRegistry, capDonID, 0, 1, true, false, capabilityDonNodes, []string{"zone-a"}, 1, [][32]byte{triggerCapID, executableCapID})
	addDON(localRegistry, zoneBDonID, 0, 1, true, false, zoneBDonNodes, []string{"zone-b"}, 1, [][32]byte{triggerCapID, executableCapID})
	addCapabilityToDON(localRegistry, capDonID, fullTriggerCapID, capabilities.CapabilityTypeTrigger, triggerCfg)
	addCapabilityToDON(localRegistry, capDonID, fullExecutableCapID, capabilities.CapabilityTypeTarget, execCfg)
	addCapabilityToDON(localRegistry, zoneBDonID, fullExecutableCapID, capabilities.CapabilityTypeTarget, execCfg)

	sharedPeer := mocks.NewSharedPeer(t)
	sharedPeer.On("ID").Return(workflowDonNodes[0])
	sharedPeer.On("IsBootstrap").Return(false)
	sharedPeer.On("UpdateConnectionsByDONs", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	launcher := NewLauncher(
		lggr,
		nil,
		sharedPeer,
		nil,
		dispatcher,
		registry,
		&mockDonNotifier{},
	)
	require.NoError(t, launcher.Start(t.Context()))
	defer launcher.Close()

	dispatcher.On("SetReceiverForMethod", fullTriggerCapID, capDonID, "StreamsTrigger", mock.AnythingOfType("*remote.triggerSubscriber")).Return(nil)
	dispatcher.On("SetReceiverForMethod", fullExecutableCapID, capDonID, "Write", mock.AnythingOfType("*executable.client")).Return(nil)

	// first test the initial CombinedClient creation
	err = launcher.OnNewRegistry(t.Context(), localRegistry)
	require.NoError(t, err)

	trigCap, err := registry.Get(t.Context(), fullTriggerCapID)
	require.NoError(t, err)
	trigCC, ok := trigCap.(remote.CombinedClient)
	assert.True(t, ok, "expected CombinedClient object")
	subscriber := trigCC.GetTriggerSubscriber("StreamsTrigger")
	capInfo, err := subscriber.Info(t.Context())
	require.NoError(t, err)
	assert.Equal(t, fullTriggerCapID, capInfo.ID)
	assert.Len(t, capInfo.DON.Members, 4)

	execCap, err := registry.Get(t.Context(), fullExecutableCapID)
	require.NoError(t, err)
	execCC, ok := execCap.(remote.CombinedClient)
	assert.True(t, ok, "expected CombinedClient object")
	require.NotNil(t, execCC.GetExecutableClient("Write"))

	// Now update config for one capability and verify it's propagated correctly (DON size)
	capDon := localRegistry.IDsToDONs[registrysyncer.DonID(capDonID)]
	capDon.Members = append(capDon.Members, ragetypes.PeerID(RandomUTF8BytesWord()))
	localRegistry.IDsToDONs[registrysyncer.DonID(capDonID)] = capDon
	err = launcher.OnNewRegistry(t.Context(), localRegistry)
	require.NoError(t, err)

	capInfo, err = subscriber.Info(t.Context())
	require.NoError(t, err)
	assert.Equal(t, fullTriggerCapID, capInfo.ID)
	assert.Len(t, capInfo.DON.Members, 5)
}

func TestLauncher_V2CapabilitiesExposeRemotely(t *testing.T) {
	lggr := logger.Test(t)
	registry := NewRegistry(lggr)
	fullTriggerCapID := "streams-trigger@1.0.0"
	mt := newMockTrigger(capabilities.MustNewCapabilityInfo(
		fullTriggerCapID,
		capabilities.CapabilityTypeTrigger,
		"streams trigger",
	))
	require.NoError(t, registry.Add(t.Context(), mt))

	fullExecutableCapID := "evm@1.0.0"
	mtarg := &mockCapability{
		CapabilityInfo: capabilities.MustNewCapabilityInfo(
			fullExecutableCapID,
			capabilities.CapabilityTypeTarget,
			"evm",
		),
	}
	require.NoError(t, registry.Add(t.Context(), mtarg))

	dispatcher := remoteMocks.NewDispatcher(t)

	workflowDonNodes, capabilityDonNodes := newNodes(4), newNodes(4)
	triggerCapID := RandomUTF8BytesWord()
	executableCapID := RandomUTF8BytesWord()
	wfDonID := uint32(1)
	capDonID := uint32(2)

	triggerCfg, err := proto.Marshal(&capabilitiespb.CapabilityConfig{
		MethodConfigs: map[string]*capabilitiespb.CapabilityMethodConfig{
			"StreamsTrigger": {
				RemoteConfig: &capabilitiespb.CapabilityMethodConfig_RemoteTriggerConfig{
					RemoteTriggerConfig: &capabilitiespb.RemoteTriggerConfig{
						RegistrationRefresh:     durationpb.New(1 * time.Second),
						MinResponsesToAggregate: 3,
					},
				},
			},
		},
	})
	require.NoError(t, err)

	execCfg, err := proto.Marshal(&capabilitiespb.CapabilityConfig{
		MethodConfigs: map[string]*capabilitiespb.CapabilityMethodConfig{
			"Write": {
				RemoteConfig: &capabilitiespb.CapabilityMethodConfig_RemoteExecutableConfig{
					RemoteExecutableConfig: &capabilitiespb.RemoteExecutableConfig{
						RequestTimeout:            durationpb.New(30 * time.Second),
						ServerMaxParallelRequests: 10,
						DeltaStage:                durationpb.New(1 * time.Second),
					},
				},
			},
		},
	})
	require.NoError(t, err)

	localRegistry := buildLocalRegistry()
	addDON(localRegistry, wfDonID, 0, 1, true, true, workflowDonNodes, []string{"zone-a"}, 1, nil)
	addDON(localRegistry, capDonID, 0, 1, true, false, capabilityDonNodes, []string{"zone-a"}, 1, [][32]byte{triggerCapID, executableCapID})
	addCapabilityToDON(localRegistry, capDonID, fullTriggerCapID, capabilities.CapabilityTypeTrigger, triggerCfg)
	addCapabilityToDON(localRegistry, capDonID, fullExecutableCapID, capabilities.CapabilityTypeTarget, execCfg)

	sharedPeer := mocks.NewSharedPeer(t)
	sharedPeer.On("ID").Return(capabilityDonNodes[0])
	sharedPeer.On("IsBootstrap").Return(false)
	sharedPeer.On("UpdateConnectionsByDONs", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	launcher := NewLauncher(
		lggr,
		nil,
		sharedPeer,
		nil,
		dispatcher,
		registry,
		&mockDonNotifier{},
	)
	require.NoError(t, launcher.Start(t.Context()))
	defer launcher.Close()

	dispatcher.On("SetReceiverForMethod", fullTriggerCapID, capDonID, "StreamsTrigger", mock.AnythingOfType("*remote.triggerPublisher")).Return(nil)
	dispatcher.On("SetReceiverForMethod", fullExecutableCapID, capDonID, "Write", mock.AnythingOfType("*executable.server")).Return(nil)

	err = launcher.OnNewRegistry(t.Context(), localRegistry)
	require.NoError(t, err)
}

// Helper functions for building LocalRegistry
func newNodes(count int) []ragetypes.PeerID {
	nodes := make([]ragetypes.PeerID, count)
	for i := range count {
		nodes[i] = RandomUTF8BytesWord()
	}
	return nodes
}

func buildLocalRegistry() *registrysyncer.LocalRegistry {
	return &registrysyncer.LocalRegistry{
		IDsToDONs:         make(map[registrysyncer.DonID]registrysyncer.DON),
		IDsToCapabilities: make(map[string]registrysyncer.Capability),
		IDsToNodes:        make(map[ragetypes.PeerID]registrysyncer.NodeInfo),
	}
}

func addDON(registry *registrysyncer.LocalRegistry, donID uint32, configVersion uint32, f uint8, isPublic bool, acceptsWorkflows bool, members []ragetypes.PeerID, families []string, operatorID uint32, hashedCapabilityIDs [][32]byte) {
	registry.IDsToDONs[registrysyncer.DonID(donID)] = registrysyncer.DON{
		DON: capabilities.DON{
			ID:               donID,
			ConfigVersion:    configVersion,
			F:                f,
			IsPublic:         isPublic,
			AcceptsWorkflows: acceptsWorkflows,
			Members:          members,
			Families:         families,
		},
		CapabilityConfigurations: make(map[string]registrysyncer.CapabilityConfiguration),
	}

	// Add each member node to the registry
	for _, peerID := range members {
		registry.IDsToNodes[peerID] = registrysyncer.NodeInfo{
			NodeOperatorID:      operatorID,
			Signer:              RandomUTF8BytesWord(),
			P2pID:               peerID,
			EncryptionPublicKey: RandomUTF8BytesWord(),
			HashedCapabilityIDs: hashedCapabilityIDs,
		}
	}
}

func addCapabilityToDON(registry *registrysyncer.LocalRegistry, donID uint32, capabilityID string, capabilityType capabilities.CapabilityType, config []byte) {
	don := registry.IDsToDONs[registrysyncer.DonID(donID)]
	don.CapabilityConfigurations[capabilityID] = registrysyncer.CapabilityConfiguration{
		Config: config,
	}
	registry.IDsToDONs[registrysyncer.DonID(donID)] = don

	registry.IDsToCapabilities[capabilityID] = registrysyncer.Capability{
		ID:             capabilityID,
		CapabilityType: capabilityType,
	}
}
