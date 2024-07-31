package registrysyncer

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/eth/ethconfig"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/durationpb"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
	"github.com/smartcontractkit/chainlink-common/pkg/values"

	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	evmclient "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/headtracker"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	kcr "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/capabilities_registry"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"
	"github.com/smartcontractkit/chainlink/v2/core/services/p2p/types/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm"
	evmrelaytypes "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"
)

var writeChainCapability = kcr.CapabilitiesRegistryCapability{
	LabelledName:   "write-chain",
	Version:        "1.0.1",
	CapabilityType: uint8(3),
}

func startNewChainWithRegistry(t *testing.T) (*kcr.CapabilitiesRegistry, common.Address, *bind.TransactOpts, *backends.SimulatedBackend) {
	owner := testutils.MustNewSimTransactor(t)

	i := &big.Int{}
	oneEth, _ := i.SetString("100000000000000000000", 10)
	gasLimit := ethconfig.Defaults.Miner.GasCeil * 2 // 60 M blocks

	simulatedBackend := backends.NewSimulatedBackend(core.GenesisAlloc{owner.From: {
		Balance: oneEth,
	}}, gasLimit)
	simulatedBackend.Commit()

	CapabilitiesRegistryAddress, _, CapabilitiesRegistry, err := kcr.DeployCapabilitiesRegistry(owner, simulatedBackend)
	require.NoError(t, err, "DeployCapabilitiesRegistry failed")

	fmt.Println("Deployed CapabilitiesRegistry at", CapabilitiesRegistryAddress.Hex())
	simulatedBackend.Commit()

	return CapabilitiesRegistry, CapabilitiesRegistryAddress, owner, simulatedBackend
}

type crFactory struct {
	lggr      logger.Logger
	ht        logpoller.HeadTracker
	logPoller logpoller.LogPoller
	client    evmclient.Client
}

func (c *crFactory) NewContractReader(ctx context.Context, cfg []byte) (types.ContractReader, error) {
	crCfg := &evmrelaytypes.ChainReaderConfig{}
	if err := json.Unmarshal(cfg, crCfg); err != nil {
		return nil, err
	}

	svc, err := evm.NewChainReaderService(ctx, c.lggr, c.logPoller, c.ht, c.client, *crCfg)
	if err != nil {
		return nil, err
	}

	return svc, svc.Start(ctx)
}

func newContractReaderFactory(t *testing.T, simulatedBackend *backends.SimulatedBackend) *crFactory {
	lggr := logger.TestLogger(t)
	client := evmclient.NewSimulatedBackendClient(
		t,
		simulatedBackend,
		testutils.SimulatedChainID,
	)
	db := pgtest.NewSqlxDB(t)
	const finalityDepth = 2
	ht := headtracker.NewSimulatedHeadTracker(client, false, finalityDepth)
	lp := logpoller.NewLogPoller(
		logpoller.NewORM(testutils.SimulatedChainID, db, lggr),
		client,
		lggr,
		ht,
		logpoller.Opts{
			PollPeriod:               100 * time.Millisecond,
			FinalityDepth:            finalityDepth,
			BackfillBatchSize:        3,
			RpcBatchSize:             2,
			KeepFinalizedBlocksDepth: 1000,
		},
	)
	return &crFactory{
		lggr:      lggr,
		client:    client,
		ht:        ht,
		logPoller: lp,
	}
}

func randomWord() [32]byte {
	word := make([]byte, 32)
	_, err := rand.Read(word)
	if err != nil {
		panic(err)
	}
	return [32]byte(word)
}

type launcher struct {
	localRegistry *LocalRegistry
}

func (l *launcher) Launch(ctx context.Context, localRegistry *LocalRegistry) error {
	l.localRegistry = localRegistry
	return nil
}

func toPeerIDs(ids [][32]byte) []p2ptypes.PeerID {
	pids := []p2ptypes.PeerID{}
	for _, id := range ids {
		pids = append(pids, id)
	}
	return pids
}

func TestReader_Integration(t *testing.T) {
	ctx := testutils.Context(t)
	reg, regAddress, owner, sim := startNewChainWithRegistry(t)

	_, err := reg.AddCapabilities(owner, []kcr.CapabilitiesRegistryCapability{writeChainCapability})
	require.NoError(t, err, "AddCapability failed for %s", writeChainCapability.LabelledName)
	sim.Commit()

	cid := fmt.Sprintf("%s@%s", writeChainCapability.LabelledName, writeChainCapability.Version)

	hid, err := reg.GetHashedCapabilityId(&bind.CallOpts{}, writeChainCapability.LabelledName, writeChainCapability.Version)
	require.NoError(t, err)

	_, err = reg.AddNodeOperators(owner, []kcr.CapabilitiesRegistryNodeOperator{
		{
			Admin: owner.From,
			Name:  "TEST_NOP",
		},
	})
	require.NoError(t, err)

	nodeSet := [][32]byte{
		randomWord(),
		randomWord(),
		randomWord(),
	}

	signersSet := [][32]byte{
		randomWord(),
		randomWord(),
		randomWord(),
	}

	nodes := []kcr.CapabilitiesRegistryNodeParams{
		{
			// The first NodeOperatorId has id 1 since the id is auto-incrementing.
			NodeOperatorId:      uint32(1),
			Signer:              signersSet[0],
			P2pId:               nodeSet[0],
			HashedCapabilityIds: [][32]byte{hid},
		},
		{
			// The first NodeOperatorId has id 1 since the id is auto-incrementing.
			NodeOperatorId:      uint32(1),
			Signer:              signersSet[1],
			P2pId:               nodeSet[1],
			HashedCapabilityIds: [][32]byte{hid},
		},
		{
			// The first NodeOperatorId has id 1 since the id is auto-incrementing.
			NodeOperatorId:      uint32(1),
			Signer:              signersSet[2],
			P2pId:               nodeSet[2],
			HashedCapabilityIds: [][32]byte{hid},
		},
	}
	_, err = reg.AddNodes(owner, nodes)
	require.NoError(t, err)

	config := &capabilitiespb.CapabilityConfig{
		DefaultConfig: values.Proto(values.EmptyMap()).GetMapValue(),
		RemoteConfig: &capabilitiespb.CapabilityConfig_RemoteTriggerConfig{
			RemoteTriggerConfig: &capabilitiespb.RemoteTriggerConfig{
				RegistrationRefresh: durationpb.New(20 * time.Second),
				RegistrationExpiry:  durationpb.New(60 * time.Second),
				// F + 1
				MinResponsesToAggregate: uint32(1) + 1,
			},
		},
	}
	configb, err := proto.Marshal(config)
	if err != nil {
		panic(err)
	}

	cfgs := []kcr.CapabilitiesRegistryCapabilityConfiguration{
		{
			CapabilityId: hid,
			Config:       configb,
		},
	}
	_, err = reg.AddDON(
		owner,
		nodeSet,
		cfgs,
		true,
		true,
		1,
	)
	sim.Commit()

	require.NoError(t, err)

	wrapper := mocks.NewPeerWrapper(t)
	factory := newContractReaderFactory(t, sim)
	syncer, err := New(logger.TestLogger(t), wrapper, factory, regAddress.Hex())
	require.NoError(t, err)

	l := &launcher{}
	syncer.AddLauncher(l)

	err = syncer.sync(ctx)
	s := l.localRegistry
	require.NoError(t, err)
	assert.Len(t, s.IDsToCapabilities, 1)

	gotCap := s.IDsToCapabilities[cid]
	assert.Equal(t, Capability{
		CapabilityType: capabilities.CapabilityTypeTarget,
		ID:             "write-chain@1.0.1",
	}, gotCap)

	assert.Len(t, s.IDsToDONs, 1)
	rtc := capabilities.RemoteTriggerConfig{
		RegistrationRefresh:     20 * time.Second,
		MinResponsesToAggregate: 2,
		RegistrationExpiry:      60 * time.Second,
		MessageExpiry:           120 * time.Second,
	}
	expectedDON := DON{
		DON: capabilities.DON{
			ID:               1,
			ConfigVersion:    1,
			IsPublic:         true,
			AcceptsWorkflows: true,
			F:                1,
			Members:          toPeerIDs(nodeSet),
		},
		CapabilityConfigurations: map[string]capabilities.CapabilityConfiguration{
			cid: {
				DefaultConfig:       values.EmptyMap(),
				RemoteTriggerConfig: rtc,
			},
		},
	}
	assert.Equal(t, expectedDON, s.IDsToDONs[1])

	nodesInfo := []kcr.CapabilitiesRegistryNodeInfo{
		{
			// The first NodeOperatorId has id 1 since the id is auto-incrementing.
			NodeOperatorId:      uint32(1),
			ConfigCount:         1,
			WorkflowDONId:       1,
			Signer:              signersSet[0],
			P2pId:               nodeSet[0],
			HashedCapabilityIds: [][32]byte{hid},
			CapabilitiesDONIds:  []*big.Int{},
		},
		{
			// The first NodeOperatorId has id 1 since the id is auto-incrementing.
			NodeOperatorId:      uint32(1),
			ConfigCount:         1,
			WorkflowDONId:       1,
			Signer:              signersSet[1],
			P2pId:               nodeSet[1],
			HashedCapabilityIds: [][32]byte{hid},
			CapabilitiesDONIds:  []*big.Int{},
		},
		{
			// The first NodeOperatorId has id 1 since the id is auto-incrementing.
			NodeOperatorId:      uint32(1),
			ConfigCount:         1,
			WorkflowDONId:       1,
			Signer:              signersSet[2],
			P2pId:               nodeSet[2],
			HashedCapabilityIds: [][32]byte{hid},
			CapabilitiesDONIds:  []*big.Int{},
		},
	}

	assert.Len(t, s.IDsToNodes, 3)
	assert.Equal(t, map[p2ptypes.PeerID]kcr.CapabilitiesRegistryNodeInfo{
		nodeSet[0]: nodesInfo[0],
		nodeSet[1]: nodesInfo[1],
		nodeSet[2]: nodesInfo[2],
	}, s.IDsToNodes)
}

func TestSyncer_LocalNode(t *testing.T) {
	ctx := tests.Context(t)
	lggr := logger.TestLogger(t)

	var pid p2ptypes.PeerID
	err := pid.UnmarshalText([]byte("12D3KooWBCF1XT5Wi8FzfgNCqRL76Swv8TRU3TiD4QiJm8NMNX7N"))
	require.NoError(t, err)
	peer := mocks.NewPeer(t)
	peer.On("ID").Return(pid)
	wrapper := mocks.NewPeerWrapper(t)
	wrapper.On("GetPeer").Return(peer)

	workflowDonNodes := []p2ptypes.PeerID{
		pid,
		randomWord(),
		randomWord(),
		randomWord(),
	}

	dID := uint32(1)
	// The below state describes a Workflow DON (AcceptsWorkflows = true),
	// which exposes the streams-trigger and write_chain capabilities.
	// We expect receivers to be wired up and both capabilities to be added to the registry.
	localRegistry := LocalRegistry{
		lggr:        lggr,
		peerWrapper: wrapper,
		IDsToDONs: map[DonID]DON{
			DonID(dID): {
				DON: capabilities.DON{
					ID:               dID,
					ConfigVersion:    uint32(2),
					F:                uint8(1),
					IsPublic:         true,
					AcceptsWorkflows: true,
					Members:          workflowDonNodes,
				},
			},
		},
		IDsToNodes: map[p2ptypes.PeerID]kcr.CapabilitiesRegistryNodeInfo{
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

	node, err := localRegistry.LocalNode(ctx)
	require.NoError(t, err)

	don := capabilities.DON{
		ID:               dID,
		ConfigVersion:    2,
		Members:          workflowDonNodes,
		F:                1,
		IsPublic:         true,
		AcceptsWorkflows: true,
	}
	expectedNode := capabilities.Node{
		PeerID:         &pid,
		WorkflowDON:    don,
		CapabilityDONs: []capabilities.DON{don},
	}
	assert.Equal(t, expectedNode, node)
}
