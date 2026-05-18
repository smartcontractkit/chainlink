package capabilities

import (
	"context"
	"crypto/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/durationpb"

	ragetypes "github.com/smartcontractkit/libocr/ragep2p/types"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote"
	remoteMocks "github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/types/mocks"
	kcr "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/capabilities_registry"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"
	"github.com/smartcontractkit/chainlink/v2/core/services/p2p/types/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/services/registrysyncer"
)

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

func TestLauncher_WiresUpExternalCapabilities(t *testing.T) {
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
					fullTargetID:     {},
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
		},
		IDsToNodes: map[p2ptypes.PeerID]kcr.CapabilitiesRegistryNodeInfo{
			nodes[0]: {
				NodeOperatorId:      1,
				Signer:              randomWord(),
				P2pId:               nodes[0],
				HashedCapabilityIds: [][32]byte{triggerCapID, targetCapID},
			},
			nodes[1]: {
				NodeOperatorId:      1,
				Signer:              randomWord(),
				P2pId:               nodes[1],
				HashedCapabilityIds: [][32]byte{triggerCapID, targetCapID},
			},
			nodes[2]: {
				NodeOperatorId:      1,
				Signer:              randomWord(),
				P2pId:               nodes[2],
				HashedCapabilityIds: [][32]byte{triggerCapID, targetCapID},
			},
			nodes[3]: {
				NodeOperatorId:      1,
				Signer:              randomWord(),
				P2pId:               nodes[3],
				HashedCapabilityIds: [][32]byte{triggerCapID, targetCapID},
			},
		},
	}

	launcher := NewLauncher(
		lggr,
		wrapper,
		dispatcher,
		registry,
	)

	dispatcher.On("SetReceiver", fullTriggerCapID, dID, mock.AnythingOfType("*remote.triggerPublisher")).Return(nil)
	dispatcher.On("SetReceiver", fullTargetID, dID, mock.AnythingOfType("*target.server")).Return(nil)

	err = launcher.Launch(ctx, state)
	require.NoError(t, err)
	defer launcher.Close()
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
		IDsToNodes: map[p2ptypes.PeerID]kcr.CapabilitiesRegistryNodeInfo{
			nodes[0]: {
				NodeOperatorId:      1,
				Signer:              randomWord(),
				P2pId:               nodes[0],
				HashedCapabilityIds: [][32]byte{hashedTriggerID, hashedTargetID},
			},
			nodes[1]: {
				NodeOperatorId:      1,
				Signer:              randomWord(),
				P2pId:               nodes[1],
				HashedCapabilityIds: [][32]byte{hashedTriggerID, hashedTargetID},
			},
			nodes[2]: {
				NodeOperatorId:      1,
				Signer:              randomWord(),
				P2pId:               nodes[2],
				HashedCapabilityIds: [][32]byte{hashedTriggerID, hashedTargetID},
			},
			nodes[3]: {
				NodeOperatorId:      1,
				Signer:              randomWord(),
				P2pId:               nodes[3],
				HashedCapabilityIds: [][32]byte{hashedTriggerID, hashedTargetID},
			},
		},
	}

	launcher := NewLauncher(
		lggr,
		wrapper,
		dispatcher,
		registry,
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
		IDsToNodes: map[p2ptypes.PeerID]kcr.CapabilitiesRegistryNodeInfo{
			capabilityDonNodes[0]: {
				NodeOperatorId:      1,
				Signer:              randomWord(),
				P2pId:               capabilityDonNodes[0],
				HashedCapabilityIds: [][32]byte{triggerCapID, targetCapID},
			},
			capabilityDonNodes[1]: {
				NodeOperatorId:      1,
				Signer:              randomWord(),
				P2pId:               capabilityDonNodes[1],
				HashedCapabilityIds: [][32]byte{triggerCapID, targetCapID},
			},
			capabilityDonNodes[2]: {
				NodeOperatorId:      1,
				Signer:              randomWord(),
				P2pId:               capabilityDonNodes[2],
				HashedCapabilityIds: [][32]byte{triggerCapID, targetCapID},
			},
			capabilityDonNodes[3]: {
				NodeOperatorId:      1,
				Signer:              randomWord(),
				P2pId:               capabilityDonNodes[3],
				HashedCapabilityIds: [][32]byte{triggerCapID, targetCapID},
			},
			workflowDonNodes[0]: {
				NodeOperatorId: 1,
				Signer:         randomWord(),
				P2pId:          workflowDonNodes[0],
			},
			workflowDonNodes[1]: {
				NodeOperatorId: 1,
				Signer:         randomWord(),
				P2pId:          workflowDonNodes[1],
			},
			workflowDonNodes[2]: {
				NodeOperatorId: 1,
				Signer:         randomWord(),
				P2pId:          workflowDonNodes[2],
			},
			workflowDonNodes[3]: {
				NodeOperatorId: 1,
				Signer:         randomWord(),
				P2pId:          workflowDonNodes[3],
			},
		},
	}

	launcher := NewLauncher(
		lggr,
		wrapper,
		dispatcher,
		registry,
	)

	dispatcher.On("SetReceiver", fullTriggerCapID, capDonID, mock.AnythingOfType("*remote.triggerSubscriber")).Return(nil)
	dispatcher.On("SetReceiver", fullTargetID, capDonID, mock.AnythingOfType("*target.client")).Return(nil)

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
		IDsToNodes: map[p2ptypes.PeerID]kcr.CapabilitiesRegistryNodeInfo{
			capabilityDonNodes[0]: {
				NodeOperatorId:      1,
				Signer:              randomWord(),
				P2pId:               capabilityDonNodes[0],
				HashedCapabilityIds: [][32]byte{triggerCapID, targetCapID},
			},
			capabilityDonNodes[1]: {
				NodeOperatorId:      1,
				Signer:              randomWord(),
				P2pId:               capabilityDonNodes[1],
				HashedCapabilityIds: [][32]byte{triggerCapID, targetCapID},
			},
			capabilityDonNodes[2]: {
				NodeOperatorId:      1,
				Signer:              randomWord(),
				P2pId:               capabilityDonNodes[2],
				HashedCapabilityIds: [][32]byte{triggerCapID, targetCapID},
			},
			capabilityDonNodes[3]: {
				NodeOperatorId:      1,
				Signer:              randomWord(),
				P2pId:               capabilityDonNodes[3],
				HashedCapabilityIds: [][32]byte{triggerCapID, targetCapID},
			},
			workflowDonNodes[0]: {
				NodeOperatorId: 1,
				Signer:         randomWord(),
				P2pId:          workflowDonNodes[0],
			},
			workflowDonNodes[1]: {
				NodeOperatorId: 1,
				Signer:         randomWord(),
				P2pId:          workflowDonNodes[1],
			},
			workflowDonNodes[2]: {
				NodeOperatorId: 1,
				Signer:         randomWord(),
				P2pId:          workflowDonNodes[2],
			},
			workflowDonNodes[3]: {
				NodeOperatorId: 1,
				Signer:         randomWord(),
				P2pId:          workflowDonNodes[3],
			},
		},
	}

	launcher := NewLauncher(
		lggr,
		wrapper,
		dispatcher,
		registry,
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
		IDsToNodes: map[p2ptypes.PeerID]kcr.CapabilitiesRegistryNodeInfo{
			capabilityDonNodes[0]: {
				NodeOperatorId:      1,
				Signer:              randomWord(),
				P2pId:               capabilityDonNodes[0],
				HashedCapabilityIds: [][32]byte{triggerCapID},
			},
			capabilityDonNodes[1]: {
				NodeOperatorId:      1,
				Signer:              randomWord(),
				P2pId:               capabilityDonNodes[1],
				HashedCapabilityIds: [][32]byte{triggerCapID},
			},
			capabilityDonNodes[2]: {
				NodeOperatorId:      1,
				Signer:              randomWord(),
				P2pId:               capabilityDonNodes[2],
				HashedCapabilityIds: [][32]byte{triggerCapID},
			},
			capabilityDonNodes[3]: {
				NodeOperatorId:      1,
				Signer:              randomWord(),
				P2pId:               capabilityDonNodes[3],
				HashedCapabilityIds: [][32]byte{triggerCapID},
			},
			workflowDonNodes[0]: {
				NodeOperatorId: 1,
				Signer:         randomWord(),
				P2pId:          workflowDonNodes[0],
			},
			workflowDonNodes[1]: {
				NodeOperatorId: 1,
				Signer:         randomWord(),
				P2pId:          workflowDonNodes[1],
			},
			workflowDonNodes[2]: {
				NodeOperatorId: 1,
				Signer:         randomWord(),
				P2pId:          workflowDonNodes[2],
			},
			workflowDonNodes[3]: {
				NodeOperatorId: 1,
				Signer:         randomWord(),
				P2pId:          workflowDonNodes[3],
			},
		},
	}

	dispatcher.On("SetReceiver", fullTriggerCapID, capDonID, mock.AnythingOfType("*remote.triggerPublisher")).Return(remote.ErrReceiverExists)

	launcher := NewLauncher(
		lggr,
		wrapper,
		dispatcher,
		registry,
	)

	err = launcher.Launch(ctx, state)
	require.NoError(t, err)
	defer launcher.Close()
}
