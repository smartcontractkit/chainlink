package capabilities

import (
	"context"
	"crypto/rand"
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
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
	"github.com/smartcontractkit/chainlink-common/pkg/values"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote"
	remotetypes "github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/types"
	remoteMocks "github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/types/mocks"
	kcr "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/capabilities_registry_1_1_0"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"
	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"
	"github.com/smartcontractkit/chainlink/v2/core/services/p2p/types/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/services/registrysyncer"
)

var _ capabilities.TriggerCapability = (*mockTrigger)(nil)

type mockDonNotifier struct {
}

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

var _ capabilities.TargetCapability = (*mockCapability)(nil)

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

func randomWord() [32]byte {
	word := make([]byte, 32)
	_, err := rand.Read(word)
	if err != nil {
		panic(err)
	}
	return [32]byte(word)
}

func TestLauncher(t *testing.T) {
	t.Run("OK-wires_up_external_capabilities", func(t *testing.T) {
		ctx := tests.Context(t)
		lggr := logger.TestLogger(t)
		registry := NewRegistry(lggr)
		dispatcher := remoteMocks.NewDispatcher(t)

		var pid ragetypes.PeerID
		err := pid.UnmarshalText([]byte("12D3KooWBCF1XT5Wi8FzfgNCqRL76Swv8TRU3TiD4QiJm8NMNX7N"))
		require.NoError(t, err)
		peer := mocks.NewPeer(t)
		peer.On("UpdateConnections", mock.Anything).Return(nil)
		peer.On("ID").Return(pid)
		wrapper := mocks.NewPeerWrapper(t)
		wrapper.On("GetPeer").Return(peer)

		nodes := []ragetypes.PeerID{
			pid,
			randomWord(),
			randomWord(),
			randomWord(),
		}

		fullTriggerCapID := "streams-trigger@1.0.0"
		mt := newMockTrigger(capabilities.MustNewCapabilityInfo(
			fullTriggerCapID,
			capabilities.CapabilityTypeTrigger,
			"streams trigger",
		))
		require.NoError(t, registry.Add(ctx, mt))

		fullTargetID := "write-chain_evm_1@1.0.0"
		mtarg := &mockCapability{
			CapabilityInfo: capabilities.MustNewCapabilityInfo(
				fullTargetID,
				capabilities.CapabilityTypeTarget,
				"write chain",
			),
		}
		require.NoError(t, registry.Add(ctx, mtarg))

		triggerCapID := randomWord()
		targetCapID := randomWord()
		// one capability from onchain registry is not set up locally
		fullMissingTargetID := "super-duper-target@6.6.6"
		missingTargetCapID := randomWord()
		dID := uint32(1)
		// The below state describes a Workflow DON (AcceptsWorkflows = true),
		// which exposes the streams-trigger and write_chain capabilities.
		// We expect a publisher to be wired up with this configuration, and
		// no entries should be added to the registry.
		state := &registrysyncer.LocalRegistry{
			IDsToDONs: map[registrysyncer.DonID]registrysyncer.DON{
				registrysyncer.DonID(dID): {
					DON: capabilities.DON{
						ID:               dID,
						ConfigVersion:    uint32(0),
						F:                uint8(1),
						IsPublic:         true,
						AcceptsWorkflows: true,
						Members:          nodes,
					},
					CapabilityConfigurations: map[string]registrysyncer.CapabilityConfiguration{
						fullTriggerCapID:    {},
						fullTargetID:        {},
						fullMissingTargetID: {},
					},
				},
			},
			IDsToCapabilities: map[string]registrysyncer.Capability{
				fullTriggerCapID: {
					ID:             "streams-trigger@1.0.0",
					CapabilityType: capabilities.CapabilityTypeTrigger,
				},
				fullTargetID: {
					ID:             "write-chain_evm_1@1.0.0",
					CapabilityType: capabilities.CapabilityTypeTarget,
				},
				fullMissingTargetID: {
					ID:             fullMissingTargetID,
					CapabilityType: capabilities.CapabilityTypeTarget,
				},
			},
			IDsToNodes: map[p2ptypes.PeerID]kcr.INodeInfoProviderNodeInfo{
				nodes[0]: {
					NodeOperatorId:      1,
					Signer:              randomWord(),
					P2pId:               nodes[0],
					EncryptionPublicKey: randomWord(),
					HashedCapabilityIds: [][32]byte{triggerCapID, targetCapID, missingTargetCapID},
				},
				nodes[1]: {
					NodeOperatorId:      1,
					Signer:              randomWord(),
					P2pId:               nodes[1],
					EncryptionPublicKey: randomWord(),
					HashedCapabilityIds: [][32]byte{triggerCapID, targetCapID, missingTargetCapID},
				},
				nodes[2]: {
					NodeOperatorId:      1,
					Signer:              randomWord(),
					P2pId:               nodes[2],
					EncryptionPublicKey: randomWord(),
					HashedCapabilityIds: [][32]byte{triggerCapID, targetCapID, missingTargetCapID},
				},
				nodes[3]: {
					NodeOperatorId:      1,
					Signer:              randomWord(),
					P2pId:               nodes[3],
					EncryptionPublicKey: randomWord(),
					HashedCapabilityIds: [][32]byte{triggerCapID, targetCapID, missingTargetCapID},
				},
			},
		}

		launcher := NewLauncher(
			lggr,
			wrapper,
			dispatcher,
			registry,
			&mockDonNotifier{},
		)

		dispatcher.On("SetReceiver", fullTriggerCapID, dID, mock.AnythingOfType("*remote.triggerPublisher")).Return(nil)
		dispatcher.On("SetReceiver", fullTargetID, dID, mock.AnythingOfType("*executable.server")).Return(nil)

		err = launcher.Launch(ctx, state)
		require.NoError(t, err)
		defer launcher.Close()
	})

	t.Run("NOK-invalid_trigger_capability", func(t *testing.T) {
		ctx := tests.Context(t)
		lggr, observedLogs := logger.TestLoggerObserved(t, zapcore.DebugLevel)
		registry := NewRegistry(lggr)
		dispatcher := remoteMocks.NewDispatcher(t)

		var pid ragetypes.PeerID
		err := pid.UnmarshalText([]byte("12D3KooWBCF1XT5Wi8FzfgNCqRL76Swv8TRU3TiD4QiJm8NMNX7N"))
		require.NoError(t, err)
		peer := mocks.NewPeer(t)
		peer.On("UpdateConnections", mock.Anything).Return(nil)
		peer.On("ID").Return(pid)
		wrapper := mocks.NewPeerWrapper(t)
		wrapper.On("GetPeer").Return(peer)

		nodes := []ragetypes.PeerID{
			pid,
			randomWord(),
			randomWord(),
			randomWord(),
		}

		// We intentionally create a Trigger capability with a Target type
		fullTriggerCapID := "streams-trigger@1.0.0"
		mtarg := &mockCapability{
			CapabilityInfo: capabilities.MustNewCapabilityInfo(
				fullTriggerCapID,
				capabilities.CapabilityTypeTarget,
				"wrong type capability",
			),
		}
		require.NoError(t, registry.Add(ctx, mtarg))

		triggerCapID := randomWord()

		dID := uint32(1)
		// The below state describes a Workflow DON (AcceptsWorkflows = true),
		// which exposes the streams-trigger and write_chain capabilities.
		// We expect a publisher to be wired up with this configuration, and
		// no entries should be added to the registry.
		state := &registrysyncer.LocalRegistry{
			IDsToDONs: map[registrysyncer.DonID]registrysyncer.DON{
				registrysyncer.DonID(dID): {
					DON: capabilities.DON{
						ID:               dID,
						ConfigVersion:    uint32(0),
						F:                uint8(1),
						IsPublic:         true,
						AcceptsWorkflows: true,
						Members:          nodes,
					},
					CapabilityConfigurations: map[string]registrysyncer.CapabilityConfiguration{
						fullTriggerCapID: {},
					},
				},
			},
			IDsToCapabilities: map[string]registrysyncer.Capability{
				fullTriggerCapID: {
					ID:             "streams-trigger@1.0.0",
					CapabilityType: capabilities.CapabilityTypeTrigger,
				},
			},
			IDsToNodes: map[p2ptypes.PeerID]kcr.INodeInfoProviderNodeInfo{
				nodes[0]: {
					NodeOperatorId:      1,
					Signer:              randomWord(),
					P2pId:               nodes[0],
					EncryptionPublicKey: randomWord(),
					HashedCapabilityIds: [][32]byte{triggerCapID},
				},
				nodes[1]: {
					NodeOperatorId:      1,
					Signer:              randomWord(),
					P2pId:               nodes[1],
					EncryptionPublicKey: randomWord(),
					HashedCapabilityIds: [][32]byte{triggerCapID},
				},
				nodes[2]: {
					NodeOperatorId:      1,
					Signer:              randomWord(),
					P2pId:               nodes[2],
					EncryptionPublicKey: randomWord(),
					HashedCapabilityIds: [][32]byte{triggerCapID},
				},
				nodes[3]: {
					NodeOperatorId:      1,
					Signer:              randomWord(),
					P2pId:               nodes[3],
					EncryptionPublicKey: randomWord(),
					HashedCapabilityIds: [][32]byte{triggerCapID},
				},
			},
		}

		launcher := NewLauncher(
			lggr,
			wrapper,
			dispatcher,
			registry,
			&mockDonNotifier{},
		)

		err = launcher.Launch(ctx, state)
		require.NoError(t, err)

		assert.Equal(t, 1, observedLogs.FilterMessage("failed to add server-side receiver for a trigger capability - it won't be exposed remotely").Len())
		defer launcher.Close()
	})

	t.Run("NOK-invalid_target_capability", func(t *testing.T) {
		ctx := tests.Context(t)
		lggr, observedLogs := logger.TestLoggerObserved(t, zapcore.DebugLevel)
		registry := NewRegistry(lggr)
		dispatcher := remoteMocks.NewDispatcher(t)

		var pid ragetypes.PeerID
		err := pid.UnmarshalText([]byte("12D3KooWBCF1XT5Wi8FzfgNCqRL76Swv8TRU3TiD4QiJm8NMNX7N"))
		require.NoError(t, err)
		peer := mocks.NewPeer(t)
		peer.On("UpdateConnections", mock.Anything).Return(nil)
		peer.On("ID").Return(pid)
		wrapper := mocks.NewPeerWrapper(t)
		wrapper.On("GetPeer").Return(peer)

		nodes := []ragetypes.PeerID{
			pid,
			randomWord(),
			randomWord(),
			randomWord(),
		}

		fullTargetID := "write-chain_evm_1@1.0.0"
		mt := newMockTrigger(capabilities.MustNewCapabilityInfo(
			fullTargetID,
			capabilities.CapabilityTypeTrigger,
			"streams trigger",
		))
		require.NoError(t, registry.Add(ctx, mt))

		targetCapID := randomWord()
		dID := uint32(1)
		// The below state describes a Workflow DON (AcceptsWorkflows = true),
		// which exposes the streams-trigger and write_chain capabilities.
		// We expect a publisher to be wired up with this configuration, and
		// no entries should be added to the registry.
		state := &registrysyncer.LocalRegistry{
			IDsToDONs: map[registrysyncer.DonID]registrysyncer.DON{
				registrysyncer.DonID(dID): {
					DON: capabilities.DON{
						ID:               dID,
						ConfigVersion:    uint32(0),
						F:                uint8(1),
						IsPublic:         true,
						AcceptsWorkflows: true,
						Members:          nodes,
					},
					CapabilityConfigurations: map[string]registrysyncer.CapabilityConfiguration{
						fullTargetID: {},
					},
				},
			},
			IDsToCapabilities: map[string]registrysyncer.Capability{
				fullTargetID: {
					ID:             "write-chain_evm_1@1.0.0",
					CapabilityType: capabilities.CapabilityTypeTarget,
				},
			},
			IDsToNodes: map[p2ptypes.PeerID]kcr.INodeInfoProviderNodeInfo{
				nodes[0]: {
					NodeOperatorId:      1,
					Signer:              randomWord(),
					P2pId:               nodes[0],
					EncryptionPublicKey: randomWord(),
					HashedCapabilityIds: [][32]byte{targetCapID},
				},
				nodes[1]: {
					NodeOperatorId:      1,
					Signer:              randomWord(),
					P2pId:               nodes[1],
					EncryptionPublicKey: randomWord(),
					HashedCapabilityIds: [][32]byte{targetCapID},
				},
				nodes[2]: {
					NodeOperatorId:      1,
					Signer:              randomWord(),
					P2pId:               nodes[2],
					EncryptionPublicKey: randomWord(),
					HashedCapabilityIds: [][32]byte{targetCapID},
				},
				nodes[3]: {
					NodeOperatorId:      1,
					Signer:              randomWord(),
					P2pId:               nodes[3],
					EncryptionPublicKey: randomWord(),
					HashedCapabilityIds: [][32]byte{targetCapID},
				},
			},
		}

		launcher := NewLauncher(
			lggr,
			wrapper,
			dispatcher,
			registry,
			&mockDonNotifier{},
		)

		err = launcher.Launch(ctx, state)
		require.NoError(t, err)

		assert.Equal(t, 1, observedLogs.FilterMessage("failed to add server-side receiver for a target capability - it won't be exposed remotely").Len())
		defer launcher.Close()
	})
}

func newTriggerEventMsg(t *testing.T,
	senderPeerID types.PeerID,
	workflowID string,
	triggerEvent map[string]any,
	triggerEventID string) (*remotetypes.MessageBody, *values.Map) {
	triggerEventValue, err := values.NewMap(triggerEvent)
	require.NoError(t, err)
	capResponse := capabilities.TriggerResponse{
		Event: capabilities.TriggerEvent{
			Outputs: triggerEventValue,
			ID:      triggerEventID,
		},
		Err: nil,
	}
	marshaled, err := pb.MarshalTriggerResponse(capResponse)
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
	ctx := tests.Context(t)
	lggr := logger.TestLogger(t)
	registry := NewRegistry(lggr)
	dispatcher := remoteMocks.NewDispatcher(t)

	var pid ragetypes.PeerID
	err := pid.UnmarshalText([]byte("12D3KooWBCF1XT5Wi8FzfgNCqRL76Swv8TRU3TiD4QiJm8NMNX7N"))
	require.NoError(t, err)
	peer := mocks.NewPeer(t)
	peer.On("UpdateConnections", mock.Anything).Return(nil)
	peer.On("ID").Return(pid)
	wrapper := mocks.NewPeerWrapper(t)
	wrapper.On("GetPeer").Return(peer)

	workflowDonNodes := []ragetypes.PeerID{
		pid,
		randomWord(),
		randomWord(),
		randomWord(),
	}

	capabilityDonNodes := []ragetypes.PeerID{
		randomWord(),
		randomWord(),
		randomWord(),
		randomWord(),
	}

	fullTriggerCapID := "log-event-trigger-evm-43113@1.0.0"
	fullTargetID := "write-chain_evm_1@1.0.0"
	triggerCapID := randomWord()
	targetCapID := randomWord()
	dID := uint32(1)
	capDonID := uint32(2)
	// The below state describes a Workflow DON (AcceptsWorkflows = true),
	// which exposes the log-event-trigger and write_chain capabilities.
	// We expect receivers to be wired up and both capabilities to be added to the registry.
	rtc := &capabilities.RemoteTriggerConfig{}
	rtc.ApplyDefaults()

	cfg, err := proto.Marshal(&capabilitiespb.CapabilityConfig{
		RemoteConfig: &capabilitiespb.CapabilityConfig_RemoteTriggerConfig{
			RemoteTriggerConfig: &capabilitiespb.RemoteTriggerConfig{
				RegistrationRefresh:     durationpb.New(1 * time.Second),
				MinResponsesToAggregate: 3,
			},
		},
	})
	require.NoError(t, err)

	state := &registrysyncer.LocalRegistry{
		IDsToDONs: map[registrysyncer.DonID]registrysyncer.DON{
			registrysyncer.DonID(dID): {
				DON: capabilities.DON{
					ID:               dID,
					ConfigVersion:    uint32(0),
					F:                uint8(1),
					IsPublic:         true,
					AcceptsWorkflows: true,
					Members:          workflowDonNodes,
				},
			},
			registrysyncer.DonID(capDonID): {
				DON: capabilities.DON{
					ID:               capDonID,
					ConfigVersion:    uint32(0),
					F:                uint8(1),
					IsPublic:         true,
					AcceptsWorkflows: false,
					Members:          capabilityDonNodes,
				},
				CapabilityConfigurations: map[string]registrysyncer.CapabilityConfiguration{
					fullTriggerCapID: {
						Config: cfg,
					},
					fullTargetID: {
						Config: cfg,
					},
				},
			},
		},
		IDsToCapabilities: map[string]registrysyncer.Capability{
			fullTriggerCapID: {
				ID:             fullTriggerCapID,
				CapabilityType: capabilities.CapabilityTypeTrigger,
			},
			fullTargetID: {
				ID:             fullTargetID,
				CapabilityType: capabilities.CapabilityTypeTarget,
			},
		},
		IDsToNodes: map[p2ptypes.PeerID]kcr.INodeInfoProviderNodeInfo{
			capabilityDonNodes[0]: {
				NodeOperatorId:      1,
				Signer:              randomWord(),
				P2pId:               capabilityDonNodes[0],
				EncryptionPublicKey: randomWord(),
				HashedCapabilityIds: [][32]byte{triggerCapID, targetCapID},
			},
			capabilityDonNodes[1]: {
				NodeOperatorId:      1,
				Signer:              randomWord(),
				P2pId:               capabilityDonNodes[1],
				EncryptionPublicKey: randomWord(),
				HashedCapabilityIds: [][32]byte{triggerCapID, targetCapID},
			},
			capabilityDonNodes[2]: {
				NodeOperatorId:      1,
				Signer:              randomWord(),
				P2pId:               capabilityDonNodes[2],
				EncryptionPublicKey: randomWord(),
				HashedCapabilityIds: [][32]byte{triggerCapID, targetCapID},
			},
			capabilityDonNodes[3]: {
				NodeOperatorId:      1,
				Signer:              randomWord(),
				P2pId:               capabilityDonNodes[3],
				EncryptionPublicKey: randomWord(),
				HashedCapabilityIds: [][32]byte{triggerCapID, targetCapID},
			},
			workflowDonNodes[0]: {
				NodeOperatorId:      1,
				Signer:              randomWord(),
				P2pId:               workflowDonNodes[0],
				EncryptionPublicKey: randomWord(),
			},
			workflowDonNodes[1]: {
				NodeOperatorId:      1,
				Signer:              randomWord(),
				P2pId:               workflowDonNodes[1],
				EncryptionPublicKey: randomWord(),
			},
			workflowDonNodes[2]: {
				NodeOperatorId:      1,
				Signer:              randomWord(),
				P2pId:               workflowDonNodes[2],
				EncryptionPublicKey: randomWord(),
			},
			workflowDonNodes[3]: {
				NodeOperatorId:      1,
				Signer:              randomWord(),
				P2pId:               workflowDonNodes[3],
				EncryptionPublicKey: randomWord(),
			},
		},
	}

	launcher := NewLauncher(
		lggr,
		wrapper,
		dispatcher,
		registry,
		&mockDonNotifier{},
	)

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

	err = launcher.Launch(ctx, state)
	require.NoError(t, err)
	defer launcher.Close()

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
	ctx := tests.Context(t)
	lggr := logger.TestLogger(t)
	registry := NewRegistry(lggr)
	dispatcher := remoteMocks.NewDispatcher(t)

	var pid ragetypes.PeerID
	err := pid.UnmarshalText([]byte("12D3KooWBCF1XT5Wi8FzfgNCqRL76Swv8TRU3TiD4QiJm8NMNX7N"))
	require.NoError(t, err)
	peer := mocks.NewPeer(t)
	peer.On("UpdateConnections", mock.Anything).Return(nil)
	peer.On("ID").Return(pid)
	wrapper := mocks.NewPeerWrapper(t)
	wrapper.On("GetPeer").Return(peer)

	nodes := []ragetypes.PeerID{
		pid,
		randomWord(),
		randomWord(),
		randomWord(),
	}

	dID := uint32(1)
	triggerID := "streams-trigger@1.0.0"
	hashedTriggerID := randomWord()
	targetID := "write-chain_evm_1@1.0.0"
	hashedTargetID := randomWord()

	// The below state describes a Workflow DON (AcceptsWorkflows = true),
	// which isn't public (IsPublic = false), but hosts the
	// the streams-trigger and write_chain capabilities.
	// We expect no action to be taken by the syncer.
	state := &registrysyncer.LocalRegistry{
		IDsToDONs: map[registrysyncer.DonID]registrysyncer.DON{
			registrysyncer.DonID(dID): {
				DON: capabilities.DON{
					ID:               dID,
					ConfigVersion:    uint32(0),
					F:                uint8(1),
					IsPublic:         false,
					AcceptsWorkflows: true,
					Members:          nodes,
				},
				CapabilityConfigurations: map[string]registrysyncer.CapabilityConfiguration{
					triggerID: {},
					targetID:  {},
				},
			},
		},
		IDsToCapabilities: map[string]registrysyncer.Capability{
			triggerID: {
				ID:             triggerID,
				CapabilityType: capabilities.CapabilityTypeTrigger,
			},
			targetID: {
				ID:             targetID,
				CapabilityType: capabilities.CapabilityTypeTarget,
			},
		},
		IDsToNodes: map[p2ptypes.PeerID]kcr.INodeInfoProviderNodeInfo{
			nodes[0]: {
				NodeOperatorId:      1,
				Signer:              randomWord(),
				P2pId:               nodes[0],
				EncryptionPublicKey: randomWord(),
				HashedCapabilityIds: [][32]byte{hashedTriggerID, hashedTargetID},
			},
			nodes[1]: {
				NodeOperatorId:      1,
				Signer:              randomWord(),
				P2pId:               nodes[1],
				EncryptionPublicKey: randomWord(),
				HashedCapabilityIds: [][32]byte{hashedTriggerID, hashedTargetID},
			},
			nodes[2]: {
				NodeOperatorId:      1,
				Signer:              randomWord(),
				P2pId:               nodes[2],
				EncryptionPublicKey: randomWord(),
				HashedCapabilityIds: [][32]byte{hashedTriggerID, hashedTargetID},
			},
			nodes[3]: {
				NodeOperatorId:      1,
				Signer:              randomWord(),
				P2pId:               nodes[3],
				EncryptionPublicKey: randomWord(),
				HashedCapabilityIds: [][32]byte{hashedTriggerID, hashedTargetID},
			},
		},
	}

	launcher := NewLauncher(
		lggr,
		wrapper,
		dispatcher,
		registry,
		&mockDonNotifier{},
	)

	// If the DON were public, this would fail with two errors:
	// - error fetching the capabilities from the registry since they haven't been added
	// - erroneous calls to dispatcher.SetReceiver, since the call hasn't been registered.
	err = launcher.Launch(ctx, state)
	require.NoError(t, err)
	defer launcher.Close()

	// Finally, assert that no services were added.
	assert.Len(t, launcher.subServices, 0)
}

func TestLauncher_WiresUpClientsForPublicWorkflowDON(t *testing.T) {
	ctx := tests.Context(t)
	lggr := logger.TestLogger(t)
	registry := NewRegistry(lggr)
	dispatcher := remoteMocks.NewDispatcher(t)

	var pid ragetypes.PeerID
	err := pid.UnmarshalText([]byte("12D3KooWBCF1XT5Wi8FzfgNCqRL76Swv8TRU3TiD4QiJm8NMNX7N"))
	require.NoError(t, err)
	peer := mocks.NewPeer(t)
	peer.On("UpdateConnections", mock.Anything).Return(nil)
	peer.On("ID").Return(pid)
	wrapper := mocks.NewPeerWrapper(t)
	wrapper.On("GetPeer").Return(peer)

	workflowDonNodes := []ragetypes.PeerID{
		pid,
		randomWord(),
		randomWord(),
		randomWord(),
	}

	capabilityDonNodes := []ragetypes.PeerID{
		randomWord(),
		randomWord(),
		randomWord(),
		randomWord(),
	}

	fullTriggerCapID := "streams-trigger@1.0.0"
	fullTargetID := "write-chain_evm_1@1.0.0"
	triggerCapID := randomWord()
	targetCapID := randomWord()
	dID := uint32(1)
	capDonID := uint32(2)
	// The below state describes a Workflow DON (AcceptsWorkflows = true),
	// which exposes the streams-trigger and write_chain capabilities.
	// We expect receivers to be wired up and both capabilities to be added to the registry.
	rtc := &capabilities.RemoteTriggerConfig{}
	rtc.ApplyDefaults()

	cfg, err := proto.Marshal(&capabilitiespb.CapabilityConfig{
		RemoteConfig: &capabilitiespb.CapabilityConfig_RemoteTriggerConfig{
			RemoteTriggerConfig: &capabilitiespb.RemoteTriggerConfig{
				RegistrationRefresh: durationpb.New(1 * time.Second),
			},
		},
	})
	require.NoError(t, err)

	state := &registrysyncer.LocalRegistry{
		IDsToDONs: map[registrysyncer.DonID]registrysyncer.DON{
			registrysyncer.DonID(dID): {
				DON: capabilities.DON{
					ID:               dID,
					ConfigVersion:    uint32(0),
					F:                uint8(1),
					IsPublic:         true,
					AcceptsWorkflows: true,
					Members:          workflowDonNodes,
				},
			},
			registrysyncer.DonID(capDonID): {
				DON: capabilities.DON{
					ID:               capDonID,
					ConfigVersion:    uint32(0),
					F:                uint8(1),
					IsPublic:         true,
					AcceptsWorkflows: false,
					Members:          capabilityDonNodes,
				},
				CapabilityConfigurations: map[string]registrysyncer.CapabilityConfiguration{
					fullTriggerCapID: {
						Config: cfg,
					},
					fullTargetID: {
						Config: cfg,
					},
				},
			},
		},
		IDsToCapabilities: map[string]registrysyncer.Capability{
			fullTriggerCapID: {
				ID:             fullTriggerCapID,
				CapabilityType: capabilities.CapabilityTypeTrigger,
			},
			fullTargetID: {
				ID:             fullTargetID,
				CapabilityType: capabilities.CapabilityTypeTarget,
			},
		},
		IDsToNodes: map[p2ptypes.PeerID]kcr.INodeInfoProviderNodeInfo{
			capabilityDonNodes[0]: {
				NodeOperatorId:      1,
				Signer:              randomWord(),
				P2pId:               capabilityDonNodes[0],
				EncryptionPublicKey: randomWord(),
				HashedCapabilityIds: [][32]byte{triggerCapID, targetCapID},
			},
			capabilityDonNodes[1]: {
				NodeOperatorId:      1,
				Signer:              randomWord(),
				P2pId:               capabilityDonNodes[1],
				EncryptionPublicKey: randomWord(),
				HashedCapabilityIds: [][32]byte{triggerCapID, targetCapID},
			},
			capabilityDonNodes[2]: {
				NodeOperatorId:      1,
				Signer:              randomWord(),
				P2pId:               capabilityDonNodes[2],
				EncryptionPublicKey: randomWord(),
				HashedCapabilityIds: [][32]byte{triggerCapID, targetCapID},
			},
			capabilityDonNodes[3]: {
				NodeOperatorId:      1,
				Signer:              randomWord(),
				P2pId:               capabilityDonNodes[3],
				EncryptionPublicKey: randomWord(),
				HashedCapabilityIds: [][32]byte{triggerCapID, targetCapID},
			},
			workflowDonNodes[0]: {
				NodeOperatorId:      1,
				Signer:              randomWord(),
				P2pId:               workflowDonNodes[0],
				EncryptionPublicKey: randomWord(),
			},
			workflowDonNodes[1]: {
				NodeOperatorId:      1,
				Signer:              randomWord(),
				P2pId:               workflowDonNodes[1],
				EncryptionPublicKey: randomWord(),
			},
			workflowDonNodes[2]: {
				NodeOperatorId:      1,
				Signer:              randomWord(),
				P2pId:               workflowDonNodes[2],
				EncryptionPublicKey: randomWord(),
			},
			workflowDonNodes[3]: {
				NodeOperatorId:      1,
				Signer:              randomWord(),
				P2pId:               workflowDonNodes[3],
				EncryptionPublicKey: randomWord(),
			},
		},
	}

	launcher := NewLauncher(
		lggr,
		wrapper,
		dispatcher,
		registry,
		&mockDonNotifier{},
	)

	dispatcher.On("SetReceiver", fullTriggerCapID, capDonID, mock.AnythingOfType("*remote.triggerSubscriber")).Return(nil)
	dispatcher.On("SetReceiver", fullTargetID, capDonID, mock.AnythingOfType("*executable.client")).Return(nil)

	err = launcher.Launch(ctx, state)
	require.NoError(t, err)
	defer launcher.Close()

	_, err = registry.Get(ctx, fullTriggerCapID)
	require.NoError(t, err)

	_, err = registry.Get(ctx, fullTargetID)
	require.NoError(t, err)
}

func TestLauncher_WiresUpClientsForPublicWorkflowDONButIgnoresPrivateCapabilities(t *testing.T) {
	ctx := tests.Context(t)
	lggr := logger.TestLogger(t)
	registry := NewRegistry(lggr)
	dispatcher := remoteMocks.NewDispatcher(t)

	var pid ragetypes.PeerID
	err := pid.UnmarshalText([]byte("12D3KooWBCF1XT5Wi8FzfgNCqRL76Swv8TRU3TiD4QiJm8NMNX7N"))
	require.NoError(t, err)
	peer := mocks.NewPeer(t)
	peer.On("UpdateConnections", mock.Anything).Return(nil)
	peer.On("ID").Return(pid)
	wrapper := mocks.NewPeerWrapper(t)
	wrapper.On("GetPeer").Return(peer)

	workflowDonNodes := []ragetypes.PeerID{
		pid,
		randomWord(),
		randomWord(),
		randomWord(),
	}

	capabilityDonNodes := []ragetypes.PeerID{
		randomWord(),
		randomWord(),
		randomWord(),
		randomWord(),
	}

	fullTriggerCapID := "streams-trigger@1.0.0"
	fullTargetID := "write-chain_evm_1@1.0.0"
	triggerCapID := randomWord()
	targetCapID := randomWord()
	dID := uint32(1)
	triggerCapDonID := uint32(2)
	targetCapDonID := uint32(3)
	// The below state describes a Workflow DON (AcceptsWorkflows = true),
	// which exposes the streams-trigger and write_chain capabilities.
	// We expect receivers to be wired up and both capabilities to be added to the registry.
	state := &registrysyncer.LocalRegistry{
		IDsToDONs: map[registrysyncer.DonID]registrysyncer.DON{
			registrysyncer.DonID(dID): {
				DON: capabilities.DON{
					ID:               dID,
					ConfigVersion:    uint32(0),
					F:                uint8(1),
					IsPublic:         true,
					AcceptsWorkflows: true,
					Members:          workflowDonNodes,
				},
			},
			registrysyncer.DonID(triggerCapDonID): {
				DON: capabilities.DON{
					ID:               triggerCapDonID,
					ConfigVersion:    uint32(0),
					F:                uint8(1),
					IsPublic:         true,
					AcceptsWorkflows: false,
					Members:          capabilityDonNodes,
				},
				CapabilityConfigurations: map[string]registrysyncer.CapabilityConfiguration{
					fullTriggerCapID: {},
				},
			},
			registrysyncer.DonID(targetCapDonID): {
				DON: capabilities.DON{
					ID:               targetCapDonID,
					ConfigVersion:    uint32(0),
					F:                uint8(1),
					IsPublic:         false,
					AcceptsWorkflows: false,
					Members:          capabilityDonNodes,
				},
				CapabilityConfigurations: map[string]registrysyncer.CapabilityConfiguration{
					fullTargetID: {},
				},
			},
		},
		IDsToCapabilities: map[string]registrysyncer.Capability{
			fullTriggerCapID: {
				ID:             fullTriggerCapID,
				CapabilityType: capabilities.CapabilityTypeTrigger,
			},
			fullTargetID: {
				ID:             fullTargetID,
				CapabilityType: capabilities.CapabilityTypeTarget,
			},
		},
		IDsToNodes: map[p2ptypes.PeerID]kcr.INodeInfoProviderNodeInfo{
			capabilityDonNodes[0]: {
				NodeOperatorId:      1,
				Signer:              randomWord(),
				P2pId:               capabilityDonNodes[0],
				EncryptionPublicKey: randomWord(),
				HashedCapabilityIds: [][32]byte{triggerCapID, targetCapID},
			},
			capabilityDonNodes[1]: {
				NodeOperatorId:      1,
				Signer:              randomWord(),
				P2pId:               capabilityDonNodes[1],
				EncryptionPublicKey: randomWord(),
				HashedCapabilityIds: [][32]byte{triggerCapID, targetCapID},
			},
			capabilityDonNodes[2]: {
				NodeOperatorId:      1,
				Signer:              randomWord(),
				P2pId:               capabilityDonNodes[2],
				EncryptionPublicKey: randomWord(),
				HashedCapabilityIds: [][32]byte{triggerCapID, targetCapID},
			},
			capabilityDonNodes[3]: {
				NodeOperatorId:      1,
				Signer:              randomWord(),
				P2pId:               capabilityDonNodes[3],
				EncryptionPublicKey: randomWord(),
				HashedCapabilityIds: [][32]byte{triggerCapID, targetCapID},
			},
			workflowDonNodes[0]: {
				NodeOperatorId:      1,
				Signer:              randomWord(),
				P2pId:               workflowDonNodes[0],
				EncryptionPublicKey: randomWord(),
			},
			workflowDonNodes[1]: {
				NodeOperatorId:      1,
				Signer:              randomWord(),
				P2pId:               workflowDonNodes[1],
				EncryptionPublicKey: randomWord(),
			},
			workflowDonNodes[2]: {
				NodeOperatorId:      1,
				Signer:              randomWord(),
				P2pId:               workflowDonNodes[2],
				EncryptionPublicKey: randomWord(),
			},
			workflowDonNodes[3]: {
				NodeOperatorId:      1,
				Signer:              randomWord(),
				P2pId:               workflowDonNodes[3],
				EncryptionPublicKey: randomWord(),
			},
		},
	}

	launcher := NewLauncher(
		lggr,
		wrapper,
		dispatcher,
		registry,
		&mockDonNotifier{},
	)

	dispatcher.On("SetReceiver", fullTriggerCapID, triggerCapDonID, mock.AnythingOfType("*remote.triggerSubscriber")).Return(nil)

	err = launcher.Launch(ctx, state)
	require.NoError(t, err)
	defer launcher.Close()

	_, err = registry.Get(ctx, fullTriggerCapID)
	require.NoError(t, err)
}

func TestLauncher_SucceedsEvenIfDispatcherAlreadyHasReceiver(t *testing.T) {
	ctx := tests.Context(t)
	lggr := logger.TestLogger(t)
	registry := NewRegistry(lggr)
	dispatcher := remoteMocks.NewDispatcher(t)

	var pid ragetypes.PeerID
	err := pid.UnmarshalText([]byte("12D3KooWBCF1XT5Wi8FzfgNCqRL76Swv8TRU3TiD4QiJm8NMNX7N"))
	require.NoError(t, err)
	peer := mocks.NewPeer(t)
	peer.On("UpdateConnections", mock.Anything).Return(nil)
	peer.On("ID").Return(pid)
	wrapper := mocks.NewPeerWrapper(t)
	wrapper.On("GetPeer").Return(peer)

	fullTriggerCapID := "streams-trigger@1.0.0"
	mt := newMockTrigger(capabilities.MustNewCapabilityInfo(
		fullTriggerCapID,
		capabilities.CapabilityTypeTrigger,
		"streams trigger",
	))
	require.NoError(t, registry.Add(ctx, mt))

	workflowDonNodes := []p2ptypes.PeerID{
		randomWord(),
		randomWord(),
		randomWord(),
		randomWord(),
	}

	capabilityDonNodes := []p2ptypes.PeerID{
		pid,
		randomWord(),
		randomWord(),
		randomWord(),
	}

	triggerCapID := randomWord()
	dID := uint32(1)
	capDonID := uint32(2)
	// The below state describes a Capability DON (AcceptsWorkflows = true),
	// which exposes the streams-trigger and write_chain capabilities.
	// We expect receivers to be wired up.
	state := &registrysyncer.LocalRegistry{
		IDsToDONs: map[registrysyncer.DonID]registrysyncer.DON{
			registrysyncer.DonID(dID): {
				DON: capabilities.DON{
					ID:               dID,
					ConfigVersion:    uint32(0),
					F:                uint8(1),
					IsPublic:         true,
					AcceptsWorkflows: true,
					Members:          workflowDonNodes,
				},
			},
			registrysyncer.DonID(capDonID): {
				DON: capabilities.DON{
					ID:               capDonID,
					ConfigVersion:    uint32(0),
					F:                uint8(1),
					IsPublic:         true,
					AcceptsWorkflows: false,
					Members:          capabilityDonNodes,
				},
				CapabilityConfigurations: map[string]registrysyncer.CapabilityConfiguration{
					fullTriggerCapID: {},
				},
			},
		},
		IDsToCapabilities: map[string]registrysyncer.Capability{
			fullTriggerCapID: {
				ID:             fullTriggerCapID,
				CapabilityType: capabilities.CapabilityTypeTrigger,
			},
		},
		IDsToNodes: map[p2ptypes.PeerID]kcr.INodeInfoProviderNodeInfo{
			capabilityDonNodes[0]: {
				NodeOperatorId:      1,
				Signer:              randomWord(),
				P2pId:               capabilityDonNodes[0],
				EncryptionPublicKey: randomWord(),
				HashedCapabilityIds: [][32]byte{triggerCapID},
			},
			capabilityDonNodes[1]: {
				NodeOperatorId:      1,
				Signer:              randomWord(),
				P2pId:               capabilityDonNodes[1],
				EncryptionPublicKey: randomWord(),
				HashedCapabilityIds: [][32]byte{triggerCapID},
			},
			capabilityDonNodes[2]: {
				NodeOperatorId:      1,
				Signer:              randomWord(),
				P2pId:               capabilityDonNodes[2],
				EncryptionPublicKey: randomWord(),
				HashedCapabilityIds: [][32]byte{triggerCapID},
			},
			capabilityDonNodes[3]: {
				NodeOperatorId:      1,
				Signer:              randomWord(),
				P2pId:               capabilityDonNodes[3],
				EncryptionPublicKey: randomWord(),
				HashedCapabilityIds: [][32]byte{triggerCapID},
			},
			workflowDonNodes[0]: {
				NodeOperatorId:      1,
				Signer:              randomWord(),
				P2pId:               workflowDonNodes[0],
				EncryptionPublicKey: randomWord(),
			},
			workflowDonNodes[1]: {
				NodeOperatorId:      1,
				Signer:              randomWord(),
				P2pId:               workflowDonNodes[1],
				EncryptionPublicKey: randomWord(),
			},
			workflowDonNodes[2]: {
				NodeOperatorId:      1,
				Signer:              randomWord(),
				P2pId:               workflowDonNodes[2],
				EncryptionPublicKey: randomWord(),
			},
			workflowDonNodes[3]: {
				NodeOperatorId:      1,
				Signer:              randomWord(),
				P2pId:               workflowDonNodes[3],
				EncryptionPublicKey: randomWord(),
			},
		},
	}

	dispatcher.On("SetReceiver", fullTriggerCapID, capDonID, mock.AnythingOfType("*remote.triggerPublisher")).Return(remote.ErrReceiverExists)

	launcher := NewLauncher(
		lggr,
		wrapper,
		dispatcher,
		registry,
		&mockDonNotifier{},
	)

	err = launcher.Launch(ctx, state)
	require.NoError(t, err)
	defer launcher.Close()
}
