package capabilities

import (
	"context"
	"crypto/rand"
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/durationpb"

	"github.com/smartcontractkit/libocr/ragep2p"
	ragetypes "github.com/smartcontractkit/libocr/ragep2p/types"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services/servicetest"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote"
	remoteMocks "github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/types/mocks"
	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"
	"github.com/smartcontractkit/chainlink/v2/core/services/p2p/types/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/services/registrysyncer"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

// RandomUTF8BytesWord generates a [32]byte array containing random UTF-8 encoded characters.
func RandomUTF8BytesWord() [32]byte {
	var result [32]byte
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	for i := range 32 {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			panic(err)
		}
		result[i] = letters[num.Int64()]
	}
	return result
}

var _ capabilities.TriggerCapability = (*mockTrigger)(nil)

type mockDonNotifier struct{}

func (m *mockDonNotifier) NotifyDonSet(don capabilities.DON) {
}

type mockTrigger struct {
	capabilities.CapabilityInfo
}

func (m *mockTrigger) AckEvent(ctx context.Context, triggerID string, eventID string, method string) error {
	return nil
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
	t.Parallel()
	t.Run("start and close", func(t *testing.T) {
		t.Parallel()
		lggr := logger.Test(t)
		registry := NewRegistry(lggr)
		dispatcher := remoteMocks.NewDispatcher(t)
		sharedPeer := mocks.NewSharedPeer(t)
		sharedPeer.On("ID").Return(ragetypes.PeerID(RandomUTF8BytesWord()))
		launcher, err := NewLauncher(
			lggr,
			sharedPeer,
			nil,
			dispatcher,
			registry,
			&mockDonNotifier{}, limits.Factory{},
		)
		require.NoError(t, err)
		require.NoError(t, launcher.Start(t.Context()))
		require.NoError(t, launcher.Close())
	})
}

func TestSyncer_IgnoresCapabilitiesForPrivateDON(t *testing.T) {
	t.Parallel()
	lggr := logger.Test(t)
	registry := NewRegistry(lggr)
	dispatcher := remoteMocks.NewDispatcher(t)

	nodes := newNodes(4)
	sharedPeer := mocks.NewSharedPeer(t)
	sharedPeer.On("ID").Return(nodes[0])
	sharedPeer.On("IsBootstrap").Return(false)
	sharedPeer.On("UpdateConnectionsByDONs", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	dID := uint32(1)
	triggerID := "streams-trigger@1.0.0"
	hashedTriggerID := RandomUTF8BytesWord()
	targetID := "write-chain_evm_1@1.0.0"
	hashedTargetID := RandomUTF8BytesWord()

	localRegistry := buildLocalRegistry()
	addDON(localRegistry, dID, uint32(0), uint8(1), false, true, nodes, []string{"zone-a"}, 1, [][32]byte{hashedTriggerID, hashedTargetID})
	addCapabilityToDON(localRegistry, dID, triggerID, capabilities.CapabilityTypeTrigger, nil)
	addCapabilityToDON(localRegistry, dID, targetID, capabilities.CapabilityTypeTarget, nil)

	launcher, err := NewLauncher(
		lggr,
		sharedPeer,
		nil,
		dispatcher,
		registry,
		&mockDonNotifier{}, limits.Factory{},
	)
	require.NoError(t, err)
	require.NoError(t, launcher.Start(t.Context()))
	defer launcher.Close()

	// If the DON were public, this would fail with two errors:
	// - error fetching the capabilities from the registry since they haven't been added
	// - erroneous calls to dispatcher.SetReceiver, since the call hasn't been registered.
	err = launcher.OnNewRegistry(t.Context(), localRegistry)
	require.NoError(t, err)

	// Finally, assert that no services were added.
	assert.Empty(t, launcher.subServices)
}

func TestLauncher_DonPairsToUpdate(t *testing.T) {
	t.Parallel()
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

	workflowDonNodes := newNodes(4)
	capabilityDonNodes := newNodes(4)
	capabilityDonNodes[0] = pid // node belongs to the capability DON
	triggerCapID := RandomUTF8BytesWord()

	wfDONID, capDONID, mixedDONID := registrysyncer.DonID(7), registrysyncer.DonID(12), registrysyncer.DonID(33)
	localRegistry := buildLocalRegistry()
	addDON(localRegistry, uint32(wfDONID), uint32(0), uint8(1), true, true, workflowDonNodes, nil, 1, nil)
	addDON(localRegistry, uint32(capDONID), uint32(0), uint8(1), true, false, capabilityDonNodes, nil, 1, [][32]byte{triggerCapID})
	addCapabilityToDON(localRegistry, uint32(capDONID), fullTriggerCapID, capabilities.CapabilityTypeTrigger, nil)
	addDON(localRegistry, uint32(mixedDONID), uint32(0), uint8(1), true, true, capabilityDonNodes[2:3], nil, 1, [][32]byte{triggerCapID})
	addCapabilityToDON(localRegistry, uint32(mixedDONID), fullTriggerCapID, capabilities.CapabilityTypeTrigger, nil)
	launcher, err := NewLauncher(logger.Test(t), sharedPeer, nil, dispatcher, registry, &mockDonNotifier{}, limits.Factory{})
	require.NoError(t, err)

	sharedPeer.On("IsBootstrap").Return(false).Times(3)
	// capability DON connects to DONs: workflow and mixed
	res := launcher.donPairsToUpdate(capabilityDonNodes[0], localRegistry)
	require.Len(t, res, 2)
	require.Equal(t, p2ptypes.DonPair{localRegistry.IDsToDONs[wfDONID].DON, localRegistry.IDsToDONs[capDONID].DON}, res[0])
	require.Equal(t, p2ptypes.DonPair{localRegistry.IDsToDONs[capDONID].DON, localRegistry.IDsToDONs[mixedDONID].DON}, res[1])

	// workflow DON connects to DONs: capability and mixed
	res = launcher.donPairsToUpdate(workflowDonNodes[0], localRegistry)
	require.Len(t, res, 2)
	require.Equal(t, p2ptypes.DonPair{localRegistry.IDsToDONs[wfDONID].DON, localRegistry.IDsToDONs[capDONID].DON}, res[0])
	require.Equal(t, p2ptypes.DonPair{localRegistry.IDsToDONs[wfDONID].DON, localRegistry.IDsToDONs[mixedDONID].DON}, res[1])

	// peer (not bootstrap) that doesn't belong to any DON connects to nobody
	require.Empty(t, launcher.donPairsToUpdate(other, localRegistry))

	// bootstrap node adds all 3 DON pairs + 3 self-pairs
	sharedPeer.On("IsBootstrap").Return(true).Once()
	res = launcher.donPairsToUpdate(pid, localRegistry)
	require.Len(t, res, 6)
	require.Equal(t, p2ptypes.DonPair{localRegistry.IDsToDONs[wfDONID].DON, localRegistry.IDsToDONs[capDONID].DON}, res[0])
	require.Equal(t, p2ptypes.DonPair{localRegistry.IDsToDONs[wfDONID].DON, localRegistry.IDsToDONs[mixedDONID].DON}, res[1])
	require.Equal(t, p2ptypes.DonPair{localRegistry.IDsToDONs[capDONID].DON, localRegistry.IDsToDONs[mixedDONID].DON}, res[2])
	require.Equal(t, p2ptypes.DonPair{localRegistry.IDsToDONs[wfDONID].DON, localRegistry.IDsToDONs[wfDONID].DON}, res[3])
	require.Equal(t, p2ptypes.DonPair{localRegistry.IDsToDONs[capDONID].DON, localRegistry.IDsToDONs[capDONID].DON}, res[4])
	require.Equal(t, p2ptypes.DonPair{localRegistry.IDsToDONs[mixedDONID].DON, localRegistry.IDsToDONs[mixedDONID].DON}, res[5])

	// bootstrap node adds only allowed DON pairs + 3 self-pairs
	mixedDON := localRegistry.IDsToDONs[mixedDONID]
	mixedDON.AcceptsWorkflows = false
	localRegistry.IDsToDONs[mixedDONID] = mixedDON
	sharedPeer.On("IsBootstrap").Return(true).Once()
	res = launcher.donPairsToUpdate(pid, localRegistry)
	require.Len(t, res, 5)
	require.Equal(t, p2ptypes.DonPair{localRegistry.IDsToDONs[wfDONID].DON, localRegistry.IDsToDONs[capDONID].DON}, res[0])
	require.Equal(t, p2ptypes.DonPair{localRegistry.IDsToDONs[wfDONID].DON, localRegistry.IDsToDONs[mixedDONID].DON}, res[1])
	require.Equal(t, p2ptypes.DonPair{localRegistry.IDsToDONs[wfDONID].DON, localRegistry.IDsToDONs[wfDONID].DON}, res[2])
	require.Equal(t, p2ptypes.DonPair{localRegistry.IDsToDONs[capDONID].DON, localRegistry.IDsToDONs[capDONID].DON}, res[3])
	require.Equal(t, p2ptypes.DonPair{localRegistry.IDsToDONs[mixedDONID].DON, localRegistry.IDsToDONs[mixedDONID].DON}, res[4])
}

func TestLauncher_DonPairsToUpdate_SkipsDifferentFamilies(t *testing.T) {
	t.Parallel()
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

	launcher, err := NewLauncher(logger.Test(t), sharedPeer, nil, dispatcher, registry, &mockDonNotifier{}, limits.Factory{})
	require.NoError(t, err)

	sharedPeer.On("IsBootstrap").Return(false).Once()
	// Node belongs to workflow DON, should only connect to capability DON in same family (zone-a)
	res := launcher.donPairsToUpdate(pid, localRegistry)
	require.Len(t, res, 1, "expected only one DON pair (zone-a workflow to zone-a capability)")
	require.Equal(t, p2ptypes.DonPair{localRegistry.IDsToDONs[registrysyncer.DonID(wfDONID)].DON, localRegistry.IDsToDONs[registrysyncer.DonID(capDONZoneAID)].DON}, res[0])

	// Bootstrap node should still respect family boundaries, plus add self-pairs for all DONs
	sharedPeer.On("IsBootstrap").Return(true).Once()
	res = launcher.donPairsToUpdate(pid, localRegistry)
	require.Len(t, res, 4, "bootstrap should also filter based on families but add self-pairs")
	require.Equal(t, p2ptypes.DonPair{localRegistry.IDsToDONs[registrysyncer.DonID(wfDONID)].DON, localRegistry.IDsToDONs[registrysyncer.DonID(capDONZoneAID)].DON}, res[0])
	require.Equal(t, p2ptypes.DonPair{localRegistry.IDsToDONs[registrysyncer.DonID(wfDONID)].DON, localRegistry.IDsToDONs[registrysyncer.DonID(wfDONID)].DON}, res[1])
	require.Equal(t, p2ptypes.DonPair{localRegistry.IDsToDONs[registrysyncer.DonID(capDONZoneAID)].DON, localRegistry.IDsToDONs[registrysyncer.DonID(capDONZoneAID)].DON}, res[2])
	require.Equal(t, p2ptypes.DonPair{localRegistry.IDsToDONs[registrysyncer.DonID(capDONZoneBID)].DON, localRegistry.IDsToDONs[registrysyncer.DonID(capDONZoneBID)].DON}, res[3])
}

func TestLauncher_ShardedCapabilityRoutingByFamily(t *testing.T) {
	t.Parallel()
	lggr := logger.Test(t)
	registry := NewRegistry(lggr)
	dispatcher := remoteMocks.NewDispatcher(t)

	workflowDonNodes := newNodes(4)
	capShard0Nodes, capShard1Nodes, capShard2Nodes, sharedCapNodes := newNodes(4), newNodes(4), newNodes(4), newNodes(4)
	sharedPeer := mocks.NewSharedPeer(t)
	sharedPeer.On("ID").Return(workflowDonNodes[0])
	sharedPeer.On("IsBootstrap").Return(false)
	sharedPeer.On("UpdateConnectionsByDONs", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	fullShardedTargetID := "write-chain_evm_1@1.0.0"
	fullSharedTargetID := "read-chain_evm_1@1.0.0"
	shardedCapIDHash := RandomUTF8BytesWord()
	sharedCapIDHash := RandomUTF8BytesWord()

	wfDONID := uint32(1)
	capShard0ID := uint32(2)
	capShard1ID := uint32(3)
	capShard2ID := uint32(4)
	sharedCapDONID := uint32(5)

	cfg, err := proto.Marshal(&capabilitiespb.CapabilityConfig{
		MethodConfigs: map[string]*capabilitiespb.CapabilityMethodConfig{
			"Write": {
				RemoteConfig: &capabilitiespb.CapabilityMethodConfig_RemoteExecutableConfig{
					RemoteExecutableConfig: &capabilitiespb.RemoteExecutableConfig{
						RequestTimeout: durationpb.New(30 * time.Second),
					},
				},
			},
		},
	})
	require.NoError(t, err)

	localRegistry := buildLocalRegistry()
	addDON(localRegistry, wfDONID, uint32(0), uint8(1), true, true, workflowDonNodes, []string{"zone-a", "zone-a_shard-0"}, 1, nil)
	addDON(localRegistry, capShard0ID, uint32(0), uint8(1), true, false, capShard0Nodes, []string{"zone-a_shard-0"}, 1, [][32]byte{shardedCapIDHash})
	addCapabilityToDON(localRegistry, capShard0ID, fullShardedTargetID, capabilities.CapabilityTypeTarget, cfg)
	addDON(localRegistry, capShard1ID, uint32(0), uint8(1), true, false, capShard1Nodes, []string{"zone-a_shard-1"}, 1, [][32]byte{shardedCapIDHash})
	addCapabilityToDON(localRegistry, capShard1ID, fullShardedTargetID, capabilities.CapabilityTypeTarget, cfg)
	addDON(localRegistry, capShard2ID, uint32(0), uint8(1), true, false, capShard2Nodes, []string{"zone-a_shard-2"}, 1, [][32]byte{shardedCapIDHash})
	addCapabilityToDON(localRegistry, capShard2ID, fullShardedTargetID, capabilities.CapabilityTypeTarget, cfg)
	addDON(localRegistry, sharedCapDONID, uint32(0), uint8(1), true, false, sharedCapNodes, []string{"zone-a"}, 1, [][32]byte{sharedCapIDHash})
	addCapabilityToDON(localRegistry, sharedCapDONID, fullSharedTargetID, capabilities.CapabilityTypeTarget, cfg)

	launcher, err := NewLauncher(
		lggr,
		sharedPeer,
		nil,
		dispatcher,
		registry,
		&mockDonNotifier{}, limits.Factory{},
	)
	require.NoError(t, err)
	require.NoError(t, launcher.Start(t.Context()))
	defer launcher.Close()

	dispatcher.On("SetReceiverForMethod", fullShardedTargetID, capShard0ID, "Write", mock.AnythingOfType("*executable.client")).Return(nil)
	dispatcher.On("SetReceiverForMethod", fullSharedTargetID, sharedCapDONID, "Write", mock.AnythingOfType("*executable.client")).Return(nil)

	require.NoError(t, launcher.OnNewRegistry(t.Context(), localRegistry))

	shardedCap, err := registry.Get(t.Context(), fullShardedTargetID)
	require.NoError(t, err)
	shardedInfo, err := shardedCap.Info(t.Context())
	require.NoError(t, err)
	require.NotNil(t, shardedInfo.DON)
	require.Equal(t, capShard0ID, shardedInfo.DON.ID)

	sharedCap, err := registry.Get(t.Context(), fullSharedTargetID)
	require.NoError(t, err)
	sharedInfo, err := sharedCap.Info(t.Context())
	require.NoError(t, err)
	require.NotNil(t, sharedInfo.DON)
	require.Equal(t, sharedCapDONID, sharedInfo.DON.ID)
}

func TestLauncher_DonPairsToUpdate_ShardedFamilies(t *testing.T) {
	t.Parallel()
	registry := NewRegistry(logger.Test(t))
	dispatcher := remoteMocks.NewDispatcher(t)

	var pid ragetypes.PeerID
	require.NoError(t, pid.UnmarshalText([]byte(utils.MustNewPeerID())))
	sharedPeer := mocks.NewSharedPeer(t)

	fullTargetID := "write-chain_evm_1@1.0.0"
	capIDHash := RandomUTF8BytesWord()

	workflowDonNodes := newNodes(4)
	capShard0Nodes, capShard1Nodes, sharedCapNodes := newNodes(4), newNodes(4), newNodes(4)
	workflowDonNodes[0] = pid // node belongs to the workflow shard

	wfDONID := uint32(1)
	capShard0ID := uint32(2)
	capShard1ID := uint32(3)
	sharedCapDONID := uint32(4)

	localRegistry := buildLocalRegistry()
	addDON(localRegistry, wfDONID, uint32(0), uint8(1), true, true, workflowDonNodes, []string{"zone-a", "zone-a_shard-0"}, 1, nil)
	addDON(localRegistry, capShard0ID, uint32(0), uint8(1), true, false, capShard0Nodes, []string{"zone-a_shard-0"}, 1, [][32]byte{capIDHash})
	addCapabilityToDON(localRegistry, capShard0ID, fullTargetID, capabilities.CapabilityTypeTarget, nil)
	addDON(localRegistry, capShard1ID, uint32(0), uint8(1), true, false, capShard1Nodes, []string{"zone-a_shard-1"}, 1, [][32]byte{capIDHash})
	addCapabilityToDON(localRegistry, capShard1ID, fullTargetID, capabilities.CapabilityTypeTarget, nil)
	addDON(localRegistry, sharedCapDONID, uint32(0), uint8(1), true, false, sharedCapNodes, []string{"zone-a"}, 1, [][32]byte{capIDHash})
	addCapabilityToDON(localRegistry, sharedCapDONID, fullTargetID, capabilities.CapabilityTypeTarget, nil)

	launcher, err := NewLauncher(logger.Test(t), sharedPeer, nil, dispatcher, registry, &mockDonNotifier{}, limits.Factory{})
	require.NoError(t, err)

	sharedPeer.On("IsBootstrap").Return(false).Once()
	res := launcher.donPairsToUpdate(pid, localRegistry)
	require.Len(t, res, 2, "workflow shard connects only to its shard-family capability DON and the shared-family capability DON")
	require.Equal(t, p2ptypes.DonPair{localRegistry.IDsToDONs[registrysyncer.DonID(wfDONID)].DON, localRegistry.IDsToDONs[registrysyncer.DonID(capShard0ID)].DON}, res[0])
	require.Equal(t, p2ptypes.DonPair{localRegistry.IDsToDONs[registrysyncer.DonID(wfDONID)].DON, localRegistry.IDsToDONs[registrysyncer.DonID(sharedCapDONID)].DON}, res[1])
}

func TestLauncher_DonPairsToUpdate_CapShardPairsOnlyWithWorkflowShard(t *testing.T) {
	t.Parallel()
	registry := NewRegistry(logger.Test(t))
	dispatcher := remoteMocks.NewDispatcher(t)

	var capNodePID ragetypes.PeerID
	require.NoError(t, capNodePID.UnmarshalText([]byte(utils.MustNewPeerID())))
	sharedPeer := mocks.NewSharedPeer(t)

	wfShard0Nodes, wfShard1Nodes, capShard0Nodes := newNodes(4), newNodes(4), newNodes(4)
	capShard0Nodes[0] = capNodePID

	wfShard0ID := uint32(1)
	wfShard1ID := uint32(2)
	capShard0ID := uint32(3)

	localRegistry := buildLocalRegistry()
	addDON(localRegistry, wfShard0ID, uint32(0), uint8(1), true, true, wfShard0Nodes, []string{"zone-a", "zone-a_shard-0"}, 1, nil)
	addDON(localRegistry, wfShard1ID, uint32(0), uint8(1), true, true, wfShard1Nodes, []string{"zone-a", "zone-a_shard-1"}, 1, nil)
	addDON(localRegistry, capShard0ID, uint32(0), uint8(1), true, false, capShard0Nodes, []string{"zone-a_shard-0"}, 1, [][32]byte{RandomUTF8BytesWord()})
	addCapabilityToDON(localRegistry, capShard0ID, "write-chain_evm_1@1.0.0", capabilities.CapabilityTypeTarget, nil)

	launcher, err := NewLauncher(logger.Test(t), sharedPeer, nil, dispatcher, registry, &mockDonNotifier{}, limits.Factory{})
	require.NoError(t, err)

	sharedPeer.On("IsBootstrap").Return(false).Once()
	res := launcher.donPairsToUpdate(capNodePID, localRegistry)
	require.Len(t, res, 1)
	require.Equal(t, p2ptypes.DonPair{
		localRegistry.IDsToDONs[registrysyncer.DonID(wfShard0ID)].DON,
		localRegistry.IDsToDONs[registrysyncer.DonID(capShard0ID)].DON,
	}, res[0])
}

// A bootstrap node that belongs to no shard family still wires connectivity for an isolated
// cap shard: the self-pair lets the shard discover its own members and the cross-pair connects
// it to its paired workflow shard. Bootstrap family membership is irrelevant.
func TestLauncher_DonPairsToUpdate_BootstrapConnectsIsolatedCapShard(t *testing.T) {
	t.Parallel()
	registry := NewRegistry(logger.Test(t))
	dispatcher := remoteMocks.NewDispatcher(t)

	var bootstrapPID ragetypes.PeerID
	require.NoError(t, bootstrapPID.UnmarshalText([]byte(utils.MustNewPeerID())))
	sharedPeer := mocks.NewSharedPeer(t)

	wfShard0Nodes, wfShard1Nodes, capShard0Nodes, capShard1Nodes := newNodes(4), newNodes(4), newNodes(4), newNodes(4)

	wfShard0ID := uint32(1)
	wfShard1ID := uint32(2)
	capShard0ID := uint32(3)
	capShard1ID := uint32(4)

	localRegistry := buildLocalRegistry()
	addDON(localRegistry, wfShard0ID, uint32(0), uint8(1), true, true, wfShard0Nodes, []string{"zone-a", "zone-a_shard-0"}, 1, nil)
	addDON(localRegistry, wfShard1ID, uint32(0), uint8(1), true, true, wfShard1Nodes, []string{"zone-a", "zone-a_shard-1"}, 1, nil)
	addDON(localRegistry, capShard0ID, uint32(0), uint8(1), true, false, capShard0Nodes, []string{"zone-a_shard-0"}, 1, [][32]byte{RandomUTF8BytesWord()})
	addCapabilityToDON(localRegistry, capShard0ID, "write-chain_evm_1@1.0.0", capabilities.CapabilityTypeTarget, nil)
	addDON(localRegistry, capShard1ID, uint32(0), uint8(1), true, false, capShard1Nodes, []string{"zone-a_shard-1"}, 1, [][32]byte{RandomUTF8BytesWord()})
	addCapabilityToDON(localRegistry, capShard1ID, "write-chain_evm_1@1.0.0", capabilities.CapabilityTypeTarget, nil)

	launcher, err := NewLauncher(logger.Test(t), sharedPeer, nil, dispatcher, registry, &mockDonNotifier{}, limits.Factory{})
	require.NoError(t, err)

	// bootstrapPID is a member of no DON and therefore of no shard family.
	sharedPeer.On("IsBootstrap").Return(true).Once()
	res := launcher.donPairsToUpdate(bootstrapPID, localRegistry)

	capShard0Self := p2ptypes.DonPair{
		localRegistry.IDsToDONs[registrysyncer.DonID(capShard0ID)].DON,
		localRegistry.IDsToDONs[registrysyncer.DonID(capShard0ID)].DON,
	}
	shard0Cross := p2ptypes.DonPair{
		localRegistry.IDsToDONs[registrysyncer.DonID(wfShard0ID)].DON,
		localRegistry.IDsToDONs[registrysyncer.DonID(capShard0ID)].DON,
	}
	crossShard := p2ptypes.DonPair{
		localRegistry.IDsToDONs[registrysyncer.DonID(wfShard1ID)].DON,
		localRegistry.IDsToDONs[registrysyncer.DonID(capShard0ID)].DON,
	}

	require.Contains(t, res, capShard0Self, "isolated cap shard must get a self-pair for member discovery via bootstrap")
	require.Contains(t, res, shard0Cross, "cap shard must connect to its paired workflow shard")
	require.NotContains(t, res, crossShard, "cap shard must not connect to a workflow shard in a different shard family")

	// 2 cross-pairs (shard-0, shard-1) + 4 self-pairs (one per DON).
	require.Len(t, res, 6)
}

func TestLauncher_V2CapabilitiesAddViaCombinedClient(t *testing.T) {
	t.Parallel()
	lggr := logger.Test(t)
	registry := NewRegistry(lggr)
	dispatcher := remoteMocks.NewDispatcher(t)

	workflowDonNodes, capabilityDonNodes, zoneBDonNodes := newNodes(4), newNodes(4), newNodes(4)
	fullTriggerCapID := "streams-trigger@1.0.0"
	fullExecutableCapID := "evm@1.0.0"
	fullLocalCapID := "cron-trigger@1.0.0"
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

	localCfg, err := proto.Marshal(&capabilitiespb.CapabilityConfig{
		LocalOnly: true,
	})
	require.NoError(t, err)

	localRegistry := buildLocalRegistry()
	addDON(localRegistry, wfDonID, 0, 1, true, true, workflowDonNodes, []string{"zone-a"}, 1, nil)
	addDON(localRegistry, capDonID, 0, 1, true, false, capabilityDonNodes, []string{"zone-a"}, 1, [][32]byte{triggerCapID, executableCapID})
	addDON(localRegistry, zoneBDonID, 0, 1, true, false, zoneBDonNodes, []string{"zone-b"}, 1, [][32]byte{triggerCapID, executableCapID})
	addCapabilityToDON(localRegistry, capDonID, fullTriggerCapID, capabilities.CapabilityTypeTrigger, triggerCfg)
	addCapabilityToDON(localRegistry, capDonID, fullExecutableCapID, capabilities.CapabilityTypeTarget, execCfg)
	addCapabilityToDON(localRegistry, zoneBDonID, fullExecutableCapID, capabilities.CapabilityTypeTarget, execCfg)
	addCapabilityToDON(localRegistry, capDonID, fullLocalCapID, capabilities.CapabilityTypeAction, localCfg) // should be skipped

	customStreamConfig := p2ptypes.StreamConfig{
		IncomingMessageBufferSize: 999,
		OutgoingMessageBufferSize: 888,
		MaxMessageLenBytes:        777777,
		MessageRateLimiter: ragep2p.TokenBucketParams{
			Rate:     50.0,
			Capacity: 250,
		},
		BytesRateLimiter: ragep2p.TokenBucketParams{
			Rate:     2500000.0,
			Capacity: 5000000,
		},
	}

	sharedPeer := mocks.NewSharedPeer(t)
	sharedPeer.On("ID").Return(workflowDonNodes[0])
	sharedPeer.On("IsBootstrap").Return(false)
	sharedPeer.On("UpdateConnectionsByDONs", mock.Anything, mock.Anything, customStreamConfig).Return(nil)

	launcher, err := NewLauncher(
		lggr,
		sharedPeer,
		nil,
		dispatcher,
		registry,
		&mockDonNotifier{}, limits.Factory{},
	)
	require.NoError(t, err)
	launcher.p2pStreamConfig = customStreamConfig
	servicetest.Run(t, launcher)

	dispatcher.On("SetReceiverForMethod", fullTriggerCapID, capDonID, "StreamsTrigger", mock.AnythingOfType("*remote.triggerSubscriber")).Return(nil)
	dispatcher.On("SetReceiverForMethod", fullExecutableCapID, capDonID, "Write", mock.AnythingOfType("*executable.client")).Return(nil)

	// first test the initial CombinedClient creation
	err = launcher.OnNewRegistry(t.Context(), localRegistry)
	require.NoError(t, err)

	trigCap, err := registry.Get(t.Context(), fullTriggerCapID)
	require.NoError(t, err)
	atomCC, ok := trigCap.(interface {
		Load() *capabilities.ExecutableAndTriggerCapability
	})
	require.True(t, ok, "expected CombinedClient object but got: %T", atomCC)
	loaded := atomCC.Load()
	require.NotNil(t, loaded)
	trigCC, ok := (*loaded).(remote.CombinedClient)
	require.True(t, ok, "expected CombinedClient object")
	subscriber := trigCC.GetTriggerSubscriber("StreamsTrigger")
	capInfo, err := subscriber.Info(t.Context())
	require.NoError(t, err)
	assert.Equal(t, fullTriggerCapID, capInfo.ID)
	assert.Len(t, capInfo.DON.Members, 4)

	execCap, err := registry.Get(t.Context(), fullExecutableCapID)
	require.NoError(t, err)
	atomCC, ok = execCap.(interface {
		Load() *capabilities.ExecutableAndTriggerCapability
	})
	require.True(t, ok, "expected CombinedClient object but got: %T", atomCC)
	loaded = atomCC.Load()
	require.NotNil(t, loaded)
	execCC, ok := (*loaded).(remote.CombinedClient)
	require.True(t, ok, "expected CombinedClient object")
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
	t.Parallel()
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

	fullLocalCapID := "cron-trigger@1.0.0"
	mlocal := newMockTrigger(capabilities.MustNewCapabilityInfo(
		fullLocalCapID,
		capabilities.CapabilityTypeTrigger,
		"cron",
	))
	require.NoError(t, registry.Add(t.Context(), mlocal))

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

	localCfg, err := proto.Marshal(&capabilitiespb.CapabilityConfig{
		MethodConfigs: map[string]*capabilitiespb.CapabilityMethodConfig{
			"CronTrigger": {
				RemoteConfig: &capabilitiespb.CapabilityMethodConfig_RemoteTriggerConfig{
					RemoteTriggerConfig: &capabilitiespb.RemoteTriggerConfig{
						RegistrationRefresh:     durationpb.New(1 * time.Second),
						MinResponsesToAggregate: 3,
					},
				},
			},
		},
		LocalOnly: true,
	})
	require.NoError(t, err)

	localRegistry := buildLocalRegistry()
	addDON(localRegistry, wfDonID, 0, 1, true, true, workflowDonNodes, []string{"zone-a"}, 1, nil)
	addDON(localRegistry, capDonID, 0, 1, true, false, capabilityDonNodes, []string{"zone-a"}, 1, [][32]byte{triggerCapID, executableCapID})
	addCapabilityToDON(localRegistry, capDonID, fullTriggerCapID, capabilities.CapabilityTypeTrigger, triggerCfg)
	addCapabilityToDON(localRegistry, capDonID, fullExecutableCapID, capabilities.CapabilityTypeTarget, execCfg)
	addCapabilityToDON(localRegistry, capDonID, fullLocalCapID, capabilities.CapabilityTypeAction, localCfg) // should be skipped

	sharedPeer := mocks.NewSharedPeer(t)
	sharedPeer.On("ID").Return(capabilityDonNodes[0])
	sharedPeer.On("IsBootstrap").Return(false)
	sharedPeer.On("UpdateConnectionsByDONs", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	launcher, err := NewLauncher(
		lggr,
		sharedPeer,
		nil,
		dispatcher,
		registry,
		&mockDonNotifier{}, limits.Factory{},
	)
	require.NoError(t, err)
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

func TestLauncher_OnNewRegistry_CallsLocalCapabilityManagerReconcile(t *testing.T) {
	t.Parallel()
	lggr := logger.Test(t)
	registry := NewRegistry(lggr)
	dispatcher := remoteMocks.NewDispatcher(t)

	capabilityDonNodes := newNodes(4)
	sharedPeer := mocks.NewSharedPeer(t)
	sharedPeer.On("ID").Return(capabilityDonNodes[0])
	sharedPeer.On("IsBootstrap").Return(false)
	sharedPeer.On("UpdateConnectionsByDONs", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	fullTriggerCapID := "streams-trigger@1.0.0"
	mt := newMockTrigger(capabilities.MustNewCapabilityInfo(
		fullTriggerCapID,
		capabilities.CapabilityTypeTrigger,
		"streams trigger",
	))
	require.NoError(t, registry.Add(t.Context(), mt))

	triggerCapIDHash := RandomUTF8BytesWord()
	capDonID := uint32(1)

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

	localRegistry := buildLocalRegistry()
	addDON(localRegistry, capDonID, uint32(0), uint8(1), true, false, capabilityDonNodes, []string{"zone-a"}, 1, [][32]byte{triggerCapIDHash})
	addCapabilityToDON(localRegistry, capDonID, fullTriggerCapID, capabilities.CapabilityTypeTrigger, triggerCfg)

	reconcileCalled := make(chan struct{}, 1)
	mockLCM := &mockLocalCapabilityManager{
		reconcileFn: func(ctx context.Context, dons []registrysyncer.DON) error {
			assert.Len(t, dons, 1, "should pass all DONs")
			assert.Equal(t, capDonID, dons[0].ID)
			reconcileCalled <- struct{}{}
			return nil
		},
	}

	dispatcher.On("SetReceiverForMethod", fullTriggerCapID, capDonID, "StreamsTrigger", mock.AnythingOfType("*remote.triggerPublisher")).Return(nil)

	launcher, err := NewLauncher(
		lggr,
		sharedPeer,
		nil,
		dispatcher,
		registry,
		&mockDonNotifier{}, limits.Factory{},
	)
	require.NoError(t, err)
	launcher.SetLocalCapabilityManager(mockLCM)
	require.NoError(t, launcher.Start(t.Context()))
	defer launcher.Close()

	err = launcher.OnNewRegistry(t.Context(), localRegistry)
	require.NoError(t, err)

	select {
	case <-reconcileCalled:
		// success
	default:
		t.Fatal("Reconcile was not called on LocalCapabilityManager")
	}
}

func TestLauncher_OnNewRegistry_NilLocalCapabilityManager(t *testing.T) {
	t.Parallel()
	lggr := logger.Test(t)
	registry := NewRegistry(lggr)
	dispatcher := remoteMocks.NewDispatcher(t)

	nodes := newNodes(4)
	sharedPeer := mocks.NewSharedPeer(t)
	sharedPeer.On("ID").Return(nodes[0])
	sharedPeer.On("IsBootstrap").Return(false)
	sharedPeer.On("UpdateConnectionsByDONs", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	localRegistry := buildLocalRegistry()
	dID := uint32(1)
	addDON(localRegistry, dID, uint32(0), uint8(1), true, true, nodes, []string{"zone-a"}, 1, nil)

	// Don't set localCapMgr - should not panic.
	launcher, err := NewLauncher(
		lggr,
		sharedPeer,
		nil,
		dispatcher,
		registry,
		&mockDonNotifier{}, limits.Factory{},
	)
	require.NoError(t, err)
	require.NoError(t, launcher.Start(t.Context()))
	defer launcher.Close()

	err = launcher.OnNewRegistry(t.Context(), localRegistry)
	require.NoError(t, err)
}

// mockLocalCapabilityManager is a test mock that records calls to Reconcile.
type mockLocalCapabilityManager struct {
	reconcileFn func(ctx context.Context, dons []registrysyncer.DON) error
}

func (m *mockLocalCapabilityManager) Start(context.Context) error    { return nil }
func (m *mockLocalCapabilityManager) Close() error                   { return nil }
func (m *mockLocalCapabilityManager) Ready() error                   { return nil }
func (m *mockLocalCapabilityManager) HealthReport() map[string]error { return nil }
func (m *mockLocalCapabilityManager) Name() string                   { return "mockLocalCapMgr" }
func (m *mockLocalCapabilityManager) Reconcile(ctx context.Context, dons []registrysyncer.DON) error {
	if m.reconcileFn != nil {
		return m.reconcileFn(ctx, dons)
	}
	return nil
}
