package capabilities

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	ragetypes "github.com/smartcontractkit/libocr/ragep2p/types"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	commonMocks "github.com/smartcontractkit/chainlink-common/pkg/types/mocks"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
	remoteMocks "github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/types/mocks"
	kcr "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/keystone_capability_registry"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"
	"github.com/smartcontractkit/chainlink/v2/core/services/p2p/types/mocks"
)

func TestSyncer_CleanStartClose(t *testing.T) {
	lggr := logger.TestLogger(t)
	ctx := testutils.Context(t)
	var pid ragetypes.PeerID
	err := pid.UnmarshalText([]byte("12D3KooWBCF1XT5Wi8FzfgNCqRL76Swv8TRU3TiD4QiJm8NMNX7N"))
	require.NoError(t, err)
	peer := mocks.NewPeer(t)
	peer.On("UpdateConnections", mock.Anything).Return(nil)
	peer.On("ID").Return(pid)
	wrapper := mocks.NewPeerWrapper(t)
	wrapper.On("GetPeer").Return(peer)
	registry := commonMocks.NewCapabilitiesRegistry(t)
	dispatcher := remoteMocks.NewDispatcher(t)

	networkSetup, err := NewHardcodedDonNetworkSetup()
	require.NoError(t, err)
	mr := &mockReader{}
	syncer := newRegistrySyncer(make(services.StopChan), wrapper, registry, dispatcher, lggr, networkSetup, mr)
	require.NoError(t, err)
	require.NoError(t, syncer.Start(ctx))
	require.NoError(t, syncer.Close())
}

type mockReader struct {
	s   state
	err error
}

func (m mockReader) state(ctx context.Context) (state, error) {
	return m.s, m.err
}

func (m mockReader) LocalNode(ctx context.Context) (capabilities.Node, error) {
	return capabilities.Node{}, nil
}

type mockTrigger struct {
	capabilities.CapabilityInfo
}

func (m *mockTrigger) RegisterTrigger(ctx context.Context, request capabilities.CapabilityRequest) (<-chan capabilities.CapabilityResponse, error) {
	return nil, nil
}

func (m *mockTrigger) UnregisterTrigger(ctx context.Context, request capabilities.CapabilityRequest) error {
	return nil
}

func newMockTrigger(info capabilities.CapabilityInfo) *mockTrigger {
	return &mockTrigger{CapabilityInfo: info}
}

type mockCapability struct {
	capabilities.CapabilityInfo
}

func (m *mockCapability) Execute(ctx context.Context, req capabilities.CapabilityRequest) (<-chan capabilities.CapabilityResponse, error) {
	return nil, nil
}

func (m *mockCapability) RegisterToWorkflow(ctx context.Context, request capabilities.RegisterToWorkflowRequest) error {
	return nil
}

func (m *mockCapability) UnregisterFromWorkflow(ctx context.Context, request capabilities.UnregisterFromWorkflowRequest) error {
	return nil
}

func TestSyncer_WiresUpExternalCapabilities(t *testing.T) {
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

	nodes := [][32]byte{
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
	mr := &mockReader{
		s: state{
			IDsToDONs: map[donID]kcr.CapabilitiesRegistryDONInfo{
				donID(dID): {
					Id:               dID,
					ConfigCount:      uint32(0),
					F:                uint8(1),
					IsPublic:         true,
					AcceptsWorkflows: true,
					NodeP2PIds:       nodes,
					CapabilityConfigurations: []kcr.CapabilitiesRegistryCapabilityConfiguration{
						{
							CapabilityId: triggerCapID,
							Config:       []byte(""),
						},
						{
							CapabilityId: targetCapID,
							Config:       []byte(""),
						},
					},
				},
			},
			IDsToCapabilities: map[hashedCapabilityID]kcr.CapabilitiesRegistryCapability{
				triggerCapID: {
					LabelledName:   "streams-trigger",
					Version:        "1.0.0",
					CapabilityType: 0,
				},
				targetCapID: {
					LabelledName:   "write-chain_evm_1",
					Version:        "1.0.0",
					CapabilityType: 3,
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
		},
	}
	syncer := newRegistrySyncer(make(services.StopChan), wrapper, registry, dispatcher, lggr, HardcodedDonNetworkSetup{}, mr)
	require.NoError(t, err)

	dispatcher.On("SetReceiver", fullTriggerCapID, fmt.Sprint(dID), mock.AnythingOfType("*remote.triggerPublisher")).Return(nil)
	dispatcher.On("SetReceiver", fullTargetID, fmt.Sprint(dID), mock.AnythingOfType("*target.server")).Return(nil)

	err = syncer.sync(ctx)
	require.NoError(t, err)
	defer syncer.Close()
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

	nodes := [][32]byte{
		pid,
		randomWord(),
		randomWord(),
		randomWord(),
	}

	triggerCapID := randomWord()
	targetCapID := randomWord()
	dID := uint32(1)
	// The below state describes a Workflow DON (AcceptsWorkflows = true),
	// which isn't public (IsPublic = false), but hosts the
	// the streams-trigger and write_chain capabilities.
	// We expect no action to be taken by the syncer.
	mr := &mockReader{
		s: state{
			IDsToDONs: map[donID]kcr.CapabilitiesRegistryDONInfo{
				donID(dID): {
					Id:               dID,
					ConfigCount:      uint32(0),
					F:                uint8(1),
					IsPublic:         false,
					AcceptsWorkflows: true,
					NodeP2PIds:       nodes,
					CapabilityConfigurations: []kcr.CapabilitiesRegistryCapabilityConfiguration{
						{
							CapabilityId: triggerCapID,
							Config:       []byte(""),
						},
						{
							CapabilityId: targetCapID,
							Config:       []byte(""),
						},
					},
				},
			},
			IDsToCapabilities: map[hashedCapabilityID]kcr.CapabilitiesRegistryCapability{
				triggerCapID: {
					LabelledName:   "streams-trigger",
					Version:        "1.0.0",
					CapabilityType: 0,
				},
				targetCapID: {
					LabelledName:   "write-chain_evm_1",
					Version:        "1.0.0",
					CapabilityType: 3,
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
		},
	}
	syncer := newRegistrySyncer(make(services.StopChan), wrapper, registry, dispatcher, lggr, HardcodedDonNetworkSetup{}, mr)
	require.NoError(t, err)

	// If the DON were public, this would fail with two errors:
	// - error fetching the capabilities from the registry since they haven't been added
	// - erroneous calls to dispatcher.SetReceiver, since the call hasn't been registered.
	err = syncer.sync(ctx)
	require.NoError(t, err)
	defer syncer.Close()

	// Finally, assert that no services were added.
	assert.Len(t, syncer.subServices, 0)
}

func TestSyncer_WiresUpClientsForPublicWorkflowDON(t *testing.T) {
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

	workflowDonNodes := [][32]byte{
		pid,
		randomWord(),
		randomWord(),
		randomWord(),
	}

	capabilityDonNodes := [][32]byte{
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
	mr := &mockReader{
		s: state{
			IDsToDONs: map[donID]kcr.CapabilitiesRegistryDONInfo{
				donID(dID): {
					Id:               dID,
					ConfigCount:      uint32(0),
					F:                uint8(1),
					IsPublic:         true,
					AcceptsWorkflows: true,
					NodeP2PIds:       workflowDonNodes,
				},
				donID(capDonID): {
					Id:               capDonID,
					ConfigCount:      uint32(0),
					F:                uint8(1),
					IsPublic:         true,
					AcceptsWorkflows: false,
					NodeP2PIds:       capabilityDonNodes,
					CapabilityConfigurations: []kcr.CapabilitiesRegistryCapabilityConfiguration{
						{
							CapabilityId: triggerCapID,
							Config:       []byte(""),
						},
						{
							CapabilityId: targetCapID,
							Config:       []byte(""),
						},
					},
				},
			},
			IDsToCapabilities: map[hashedCapabilityID]kcr.CapabilitiesRegistryCapability{
				triggerCapID: {
					LabelledName:   "streams-trigger",
					Version:        "1.0.0",
					CapabilityType: 0,
				},
				targetCapID: {
					LabelledName:   "write-chain_evm_1",
					Version:        "1.0.0",
					CapabilityType: 3,
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
		},
	}
	syncer := newRegistrySyncer(make(services.StopChan), wrapper, registry, dispatcher, lggr, HardcodedDonNetworkSetup{}, mr)
	require.NoError(t, err)

	dispatcher.On("SetReceiver", fullTriggerCapID, fmt.Sprint(capDonID), mock.AnythingOfType("*remote.triggerSubscriber")).Return(nil)
	dispatcher.On("SetReceiver", fullTargetID, fmt.Sprint(capDonID), mock.AnythingOfType("*target.client")).Return(nil)

	err = syncer.sync(ctx)
	require.NoError(t, err)
	defer syncer.Close()

	_, err = registry.Get(ctx, fullTriggerCapID)
	require.NoError(t, err)

	_, err = registry.Get(ctx, fullTargetID)
	require.NoError(t, err)
}

func TestSyncer_WiresUpClientsForPublicWorkflowDONButIgnoresPrivateCapabilities(t *testing.T) {
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

	workflowDonNodes := [][32]byte{
		pid,
		randomWord(),
		randomWord(),
		randomWord(),
	}

	capabilityDonNodes := [][32]byte{
		randomWord(),
		randomWord(),
		randomWord(),
		randomWord(),
	}

	fullTriggerCapID := "streams-trigger@1.0.0"
	triggerCapID := randomWord()
	targetCapID := randomWord()
	dID := uint32(1)
	triggerCapDonID := uint32(2)
	targetCapDonID := uint32(3)
	// The below state describes a Workflow DON (AcceptsWorkflows = true),
	// which exposes the streams-trigger and write_chain capabilities.
	// We expect receivers to be wired up and both capabilities to be added to the registry.
	mr := &mockReader{
		s: state{
			IDsToDONs: map[donID]kcr.CapabilitiesRegistryDONInfo{
				donID(dID): {
					Id:               dID,
					ConfigCount:      uint32(0),
					F:                uint8(1),
					IsPublic:         true,
					AcceptsWorkflows: true,
					NodeP2PIds:       workflowDonNodes,
				},
				donID(triggerCapDonID): {
					Id:               triggerCapDonID,
					ConfigCount:      uint32(0),
					F:                uint8(1),
					IsPublic:         true,
					AcceptsWorkflows: false,
					NodeP2PIds:       capabilityDonNodes,
					CapabilityConfigurations: []kcr.CapabilitiesRegistryCapabilityConfiguration{
						{
							CapabilityId: triggerCapID,
							Config:       []byte(""),
						},
					},
				},
				donID(targetCapDonID): {
					Id:               targetCapDonID,
					ConfigCount:      uint32(0),
					F:                uint8(1),
					IsPublic:         false,
					AcceptsWorkflows: false,
					NodeP2PIds:       capabilityDonNodes,
					CapabilityConfigurations: []kcr.CapabilitiesRegistryCapabilityConfiguration{
						{
							CapabilityId: targetCapID,
							Config:       []byte(""),
						},
					},
				},
			},
			IDsToCapabilities: map[hashedCapabilityID]kcr.CapabilitiesRegistryCapability{
				triggerCapID: {
					LabelledName:   "streams-trigger",
					Version:        "1.0.0",
					CapabilityType: 0,
				},
				targetCapID: {
					LabelledName:   "write-chain_evm_1",
					Version:        "1.0.0",
					CapabilityType: 3,
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
		},
	}
	syncer := newRegistrySyncer(make(services.StopChan), wrapper, registry, dispatcher, lggr, HardcodedDonNetworkSetup{}, mr)
	require.NoError(t, err)

	dispatcher.On("SetReceiver", fullTriggerCapID, fmt.Sprint(triggerCapDonID), mock.AnythingOfType("*remote.triggerSubscriber")).Return(nil)

	err = syncer.sync(ctx)
	require.NoError(t, err)
	defer syncer.Close()

	_, err = registry.Get(ctx, fullTriggerCapID)
	require.NoError(t, err)
}
