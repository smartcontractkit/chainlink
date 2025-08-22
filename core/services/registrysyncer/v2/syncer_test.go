package v2_test

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"errors"
	"math/big"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	gethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/eth/ethconfig"
	"github.com/ethereum/go-ethereum/ethclient/simulated"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/durationpb"

	p2ptypes "github.com/smartcontractkit/libocr/ragep2p/types"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query/primitives"
	capabilities_registry_v2 "github.com/smartcontractkit/chainlink-evm/gethwrappers/workflow/generated/capabilities_registry_wrapper_v2"
	evmclient "github.com/smartcontractkit/chainlink-evm/pkg/client"
	"github.com/smartcontractkit/chainlink-evm/pkg/heads/headstest"
	"github.com/smartcontractkit/chainlink-evm/pkg/logpoller"
	evmtestutils "github.com/smartcontractkit/chainlink-evm/pkg/testutils"
	"github.com/smartcontractkit/chainlink-protos/cre/go/values"
	"github.com/smartcontractkit/chainlink/v2/core/logger"

	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/services/registrysyncer"
	syncerMocks "github.com/smartcontractkit/chainlink/v2/core/services/registrysyncer/mocks"
	registrysyncer_v2 "github.com/smartcontractkit/chainlink/v2/core/services/registrysyncer/v2"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm"
	evmrelaytypes "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"
)

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

	return svc, nil
}

func newContractReaderFactory(t *testing.T, simulatedBackend *simulated.Backend) *crFactory {
	lggr := logger.TestLogger(t)
	client := evmclient.NewSimulatedBackendClient(
		t,
		simulatedBackend,
		testutils.SimulatedChainID,
	)
	db := pgtest.NewSqlxDB(t)
	const finalityDepth = 2
	ht := headstest.NewSimulatedHeadTracker(client, false, finalityDepth)
	lp := logpoller.NewLogPoller(
		logpoller.NewORM(testutils.SimulatedChainID, db, lggr),
		client,
		lggr,
		ht,
		logpoller.Opts{
			PollPeriod:               100 * time.Millisecond,
			FinalityDepth:            finalityDepth,
			BackfillBatchSize:        3,
			RPCBatchSize:             2,
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
	localRegistry *registrysyncer.LocalRegistry
	mu            sync.RWMutex
}

func (l *launcher) OnNewRegistry(_ context.Context, localRegistry *registrysyncer.LocalRegistry) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.localRegistry = localRegistry
	return nil
}

type orm struct {
	ormMock               *syncerMocks.ORM
	mu                    sync.RWMutex
	latestLocalRegistryCh chan struct{}
	addLocalRegistryCh    chan struct{}
}

func newORM(t *testing.T) *orm {
	t.Helper()

	return &orm{
		ormMock:               syncerMocks.NewORM(t),
		latestLocalRegistryCh: make(chan struct{}, 1),
		addLocalRegistryCh:    make(chan struct{}, 1),
	}
}

func (o *orm) Cleanup() {
	o.mu.Lock()
	defer o.mu.Unlock()
	close(o.latestLocalRegistryCh)
	close(o.addLocalRegistryCh)
}

func (o *orm) AddLocalRegistry(ctx context.Context, localRegistry registrysyncer.LocalRegistry) error {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.addLocalRegistryCh <- struct{}{}
	err := o.ormMock.AddLocalRegistry(ctx, localRegistry)
	return err
}

func (o *orm) LatestLocalRegistry(ctx context.Context) (*registrysyncer.LocalRegistry, error) {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.latestLocalRegistryCh <- struct{}{}
	return o.ormMock.LatestLocalRegistry(ctx)
}

func toPeerIDs(ids [][32]byte) []p2ptypes.PeerID {
	pids := make([]p2ptypes.PeerID, len(ids))
	for i, id := range ids {
		pids[i] = id
	}
	return pids
}

func TestReader_Integration(t *testing.T) {
	ctx := testutils.Context(t)
	lggr := logger.TestLogger(t)

	// Create a simulated backend similar to V1 tests
	owner := evmtestutils.MustNewSimTransactor(t)
	i := &big.Int{}
	oneEth, _ := i.SetString("100000000000000000000", 10)
	gasLimit := ethconfig.Defaults.Miner.GasCeil * 2 // 60 M blocks

	simulatedBackend := simulated.NewBackend(gethtypes.GenesisAlloc{owner.From: {
		Balance: oneEth,
	}}, simulated.WithBlockGasLimit(gasLimit))
	simulatedBackend.Commit()

	// Deploy a V2 capabilities registry
	regAddress, _, reg, err := capabilities_registry_v2.DeployCapabilitiesRegistry(owner, simulatedBackend.Client(), capabilities_registry_v2.CapabilitiesRegistryConstructorParams{})
	require.NoError(t, err, "DeployCapabilitiesRegistry failed")
	simulatedBackend.Commit()

	// Add a V2 capability with string ID and metadata
	writeChainCapabilityV2 := capabilities_registry_v2.CapabilitiesRegistryCapability{
		CapabilityId:          "write-chain@1.0.1",
		ConfigurationContract: common.Address{},
		Metadata:              []byte(`{"capabilityType": 3, "responseType": 1}`), // 3 = target capability
	}

	// Add capability
	_, err = reg.AddCapabilities(owner, []capabilities_registry_v2.CapabilitiesRegistryCapability{writeChainCapabilityV2})
	require.NoError(t, err, "AddCapability failed for %s", writeChainCapabilityV2.CapabilityId)
	simulatedBackend.Commit()

	// V2 uses string capability IDs directly
	cid := writeChainCapabilityV2.CapabilityId

	// Add node operator
	_, err = reg.AddNodeOperators(owner, []capabilities_registry_v2.CapabilitiesRegistryNodeOperator{
		{
			Admin: owner.From,
			Name:  "TEST_NOP_V2",
		},
	})
	require.NoError(t, err, "Failed to add node operator")
	simulatedBackend.Commit()

	// Create test nodes
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

	encPubKey1 := randomWord()
	encPubKey2 := randomWord()
	encPubKey3 := randomWord()

	csaKey1 := randomWord()
	csaKey2 := randomWord()
	csaKey3 := randomWord()

	// V2 nodes use string capability IDs and require CsaKey
	nodes := []capabilities_registry_v2.CapabilitiesRegistryNodeParams{
		{
			NodeOperatorId:      uint32(1),
			Signer:              signersSet[0],
			P2pId:               nodeSet[0],
			EncryptionPublicKey: encPubKey1,
			CsaKey:              csaKey1,
			CapabilityIds:       []string{cid},
		},
		{
			NodeOperatorId:      uint32(1),
			Signer:              signersSet[1],
			P2pId:               nodeSet[1],
			EncryptionPublicKey: encPubKey2,
			CsaKey:              csaKey2,
			CapabilityIds:       []string{cid},
		},
		{
			NodeOperatorId:      uint32(1),
			Signer:              signersSet[2],
			P2pId:               nodeSet[2],
			EncryptionPublicKey: encPubKey3,
			CsaKey:              csaKey3,
			CapabilityIds:       []string{cid},
		},
	}
	_, err = reg.AddNodes(owner, nodes)
	require.NoError(t, err, "Failed to add nodes")
	simulatedBackend.Commit()

	// Create capability configuration
	config := &capabilitiespb.CapabilityConfig{
		DefaultConfig: values.Proto(values.EmptyMap()).GetMapValue(),
		RemoteConfig: &capabilitiespb.CapabilityConfig_RemoteTriggerConfig{
			RemoteTriggerConfig: &capabilitiespb.RemoteTriggerConfig{
				RegistrationRefresh:     durationpb.New(20 * time.Second),
				RegistrationExpiry:      durationpb.New(60 * time.Second),
				MinResponsesToAggregate: uint32(1) + 1,
				MessageExpiry:           durationpb.New(120 * time.Second),
			},
		},
	}
	configb, err := proto.Marshal(config)
	require.NoError(t, err)

	// V2 DON configuration uses string capability IDs
	cfgs := []capabilities_registry_v2.CapabilitiesRegistryCapabilityConfiguration{
		{
			CapabilityId: cid,
			Config:       configb,
		},
	}

	// Add DON using AddDONs with DON family (V2 feature)
	newDONs := []capabilities_registry_v2.CapabilitiesRegistryNewDONParams{
		{
			Name:                     "test-don-v2",
			DonFamilies:              []string{"workflow-don-family"},
			Config:                   []byte("test-don-v2-config"),
			CapabilityConfigurations: cfgs,
			Nodes:                    nodeSet,
			F:                        1,
			IsPublic:                 true,
			AcceptsWorkflows:         true,
		},
	}
	_, err = reg.AddDONs(owner, newDONs)
	require.NoError(t, err)
	simulatedBackend.Commit()

	db := pgtest.NewSqlxDB(t)
	factory := newContractReaderFactory(t, simulatedBackend)
	syncerORM := registrysyncer.NewORM(db, lggr)
	syncer, err := registrysyncer_v2.New(lggr, func() (p2ptypes.PeerID, error) { return p2ptypes.PeerID{}, nil }, factory, regAddress.Hex(), syncerORM)
	require.NoError(t, err)

	l := &launcher{}
	syncer.AddListener(l)

	err = syncer.Sync(ctx, false)
	require.NoError(t, err)

	s := l.localRegistry
	require.NotNil(t, s)

	// Test V2 capabilities with string IDs
	assert.Len(t, s.IDsToCapabilities, 1)
	gotCap := s.IDsToCapabilities[cid]
	assert.Equal(t, registrysyncer.Capability{
		CapabilityType: capabilities.CapabilityTypeTarget,
		ID:             "write-chain@1.0.1",
	}, gotCap)

	// Test V2 DON with family
	assert.Len(t, s.IDsToDONs, 1)
	expectedDON := capabilities.DON{
		Name:             "test-don-v2",
		ID:               1,
		Families:         []string{"workflow-don-family"},
		ConfigVersion:    1,
		IsPublic:         true,
		AcceptsWorkflows: true,
		F:                1,
		Members:          toPeerIDs(nodeSet),
		Config:           []byte("test-don-v2-config"),
	}
	gotDon := s.IDsToDONs[1]
	assert.Equal(t, expectedDON, gotDon.DON)
	assert.Equal(t, configb, gotDon.CapabilityConfigurations[cid].Config)

	hashedID, err := registrysyncer_v2.HashCapabilityID(cid)
	require.NoError(t, err, "Failed to hash capability ID")

	// Test V2 node info with string capability IDs
	expectedNodesInfo := []registrysyncer.NodeInfo{
		{
			NodeOperatorID:      uint32(1),
			ConfigCount:         1,
			WorkflowDONId:       1,
			Signer:              signersSet[0],
			P2pID:               p2ptypes.PeerID(nodeSet[0]),
			EncryptionPublicKey: encPubKey1,
			CapabilityIDs:       []string{cid}, // V2 uses string IDs
			CapabilitiesDONIds:  []*big.Int{},
			HashedCapabilityIDs: [][32]byte{hashedID},
			CsaKey:              csaKey1,
		},
		{
			NodeOperatorID:      uint32(1),
			ConfigCount:         1,
			WorkflowDONId:       1,
			Signer:              signersSet[1],
			P2pID:               p2ptypes.PeerID(nodeSet[1]),
			EncryptionPublicKey: encPubKey2,
			CapabilityIDs:       []string{cid}, // V2 uses string IDs
			CapabilitiesDONIds:  []*big.Int{},
			HashedCapabilityIDs: [][32]byte{hashedID},
			CsaKey:              csaKey2,
		},
		{
			NodeOperatorID:      uint32(1),
			ConfigCount:         1,
			WorkflowDONId:       1,
			Signer:              signersSet[2],
			P2pID:               p2ptypes.PeerID(nodeSet[2]),
			EncryptionPublicKey: encPubKey3,
			CapabilityIDs:       []string{cid}, // V2 uses string IDs
			CapabilitiesDONIds:  []*big.Int{},
			HashedCapabilityIDs: [][32]byte{hashedID},
			CsaKey:              csaKey3,
		},
	}

	assert.Len(t, s.IDsToNodes, 3)
	assert.Equal(t, map[p2ptypes.PeerID]registrysyncer.NodeInfo{
		nodeSet[0]: expectedNodesInfo[0],
		nodeSet[1]: expectedNodesInfo[1],
		nodeSet[2]: expectedNodesInfo[2],
	}, s.IDsToNodes)
}

func TestSyncer_V2_DBIntegration(t *testing.T) {
	ctx := testutils.Context(t)
	lggr := logger.TestLogger(t)

	// Create a simulated backend similar to V1 tests
	owner := evmtestutils.MustNewSimTransactor(t)
	i := &big.Int{}
	oneEth, _ := i.SetString("100000000000000000000", 10)
	gasLimit := ethconfig.Defaults.Miner.GasCeil * 2 // 60 M blocks

	simulatedBackend := simulated.NewBackend(gethtypes.GenesisAlloc{owner.From: {
		Balance: oneEth,
	}}, simulated.WithBlockGasLimit(gasLimit))
	simulatedBackend.Commit()

	// Deploy a V2 capabilities registry
	regAddress, _, reg, err := capabilities_registry_v2.DeployCapabilitiesRegistry(owner, simulatedBackend.Client(), capabilities_registry_v2.CapabilitiesRegistryConstructorParams{})
	require.NoError(t, err, "DeployCapabilitiesRegistry failed")
	simulatedBackend.Commit()

	// Add a V2 capability
	writeChainCapabilityV2 := capabilities_registry_v2.CapabilitiesRegistryCapability{
		CapabilityId:          "write-chain@1.0.1",
		ConfigurationContract: common.Address{},
		Metadata:              []byte(`{"capabilityType": 3, "responseType": 1}`), // 3 = target capability
	}

	_, err = reg.AddCapabilities(owner, []capabilities_registry_v2.CapabilitiesRegistryCapability{writeChainCapabilityV2})
	require.NoError(t, err, "AddCapability failed for %s", writeChainCapabilityV2.CapabilityId)
	simulatedBackend.Commit()

	cid := writeChainCapabilityV2.CapabilityId

	// Add node operator
	_, err = reg.AddNodeOperators(owner, []capabilities_registry_v2.CapabilitiesRegistryNodeOperator{
		{
			Admin: owner.From,
			Name:  "TEST_NOP_V2",
		},
	})
	require.NoError(t, err)
	simulatedBackend.Commit()

	// Create test nodes
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

	nodes := []capabilities_registry_v2.CapabilitiesRegistryNodeParams{
		{
			NodeOperatorId:      uint32(1),
			Signer:              signersSet[0],
			P2pId:               nodeSet[0],
			EncryptionPublicKey: randomWord(),
			CsaKey:              randomWord(),
			CapabilityIds:       []string{cid},
		},
		{
			NodeOperatorId:      uint32(1),
			Signer:              signersSet[1],
			P2pId:               nodeSet[1],
			EncryptionPublicKey: randomWord(),
			CsaKey:              randomWord(),
			CapabilityIds:       []string{cid},
		},
		{
			NodeOperatorId:      uint32(1),
			Signer:              signersSet[2],
			P2pId:               nodeSet[2],
			EncryptionPublicKey: randomWord(),
			CsaKey:              randomWord(),
			CapabilityIds:       []string{cid},
		},
	}
	_, err = reg.AddNodes(owner, nodes)
	require.NoError(t, err)
	simulatedBackend.Commit()

	// Create capability configuration
	config := &capabilitiespb.CapabilityConfig{
		DefaultConfig: values.Proto(values.EmptyMap()).GetMapValue(),
		RemoteConfig: &capabilitiespb.CapabilityConfig_RemoteTriggerConfig{
			RemoteTriggerConfig: &capabilitiespb.RemoteTriggerConfig{
				RegistrationRefresh:     durationpb.New(20 * time.Second),
				RegistrationExpiry:      durationpb.New(60 * time.Second),
				MinResponsesToAggregate: uint32(1) + 1,
			},
		},
	}
	configb, err := proto.Marshal(config)
	require.NoError(t, err)

	cfgs := []capabilities_registry_v2.CapabilitiesRegistryCapabilityConfiguration{
		{
			CapabilityId: cid,
			Config:       configb,
		},
	}

	// Add DON using AddDONs with DON family (V2 feature)
	newDONs := []capabilities_registry_v2.CapabilitiesRegistryNewDONParams{
		{
			Name:                     "test-don-v2-db",
			DonFamilies:              []string{"workflow-don-family-v2"},
			Config:                   []byte("test-don-v2-db-config"),
			CapabilityConfigurations: cfgs,
			Nodes:                    nodeSet,
			F:                        1,
			IsPublic:                 true,
			AcceptsWorkflows:         true,
		},
	}
	_, err = reg.AddDONs(owner, newDONs)
	require.NoError(t, err)
	simulatedBackend.Commit()

	// Test database integration
	syncerORM := newORM(t)
	syncerORM.ormMock.On("LatestLocalRegistry", mock.Anything).Return(nil, errors.New("no state found"))
	syncerORM.ormMock.On("AddLocalRegistry", mock.Anything, mock.Anything).Return(nil)

	factory := newContractReaderFactory(t, simulatedBackend)

	syncer, err := registrysyncer_v2.New(
		lggr,
		func() (p2ptypes.PeerID, error) { return p2ptypes.PeerID{}, nil },
		factory,
		regAddress.Hex(),
		syncerORM,
	)
	require.NoError(t, err)
	require.NoError(t, syncer.Start(ctx))

	t.Cleanup(func() {
		syncerORM.Cleanup()
		require.NoError(t, syncer.Close())
	})

	l := &launcher{}
	syncer.AddListener(l)

	// Test that the syncer calls the ORM methods
	var latestLocalRegistryCalled, addLocalRegistryCalled bool
	timeout := time.After(testutils.WaitTimeout(t))

	for !latestLocalRegistryCalled || !addLocalRegistryCalled {
		select {
		case val := <-syncerORM.latestLocalRegistryCh:
			assert.Equal(t, struct{}{}, val)
			latestLocalRegistryCalled = true
		case val := <-syncerORM.addLocalRegistryCh:
			assert.Equal(t, struct{}{}, val)
			addLocalRegistryCalled = true
		case <-timeout:
			t.Fatal("test timed out; channels did not received data")
		}
	}
}

func TestSyncer_V2_LocalNode(t *testing.T) {
	ctx := testutils.Context(t)
	lggr := logger.TestLogger(t)

	var pid p2ptypes.PeerID
	err := pid.UnmarshalText([]byte("12D3KooWBCF1XT5Wi8FzfgNCqRL76Swv8TRU3TiD4QiJm8NMNX7N"))
	require.NoError(t, err)

	workflowDonNodes := []p2ptypes.PeerID{
		pid,
		randomWord(),
		randomWord(),
		randomWord(),
	}

	dID := uint32(1)
	dName := "test-don-v2-db"
	dFamilies := []string{"workflow-don-family-v2"}
	dConfig := []byte("test-don-v2-db-config")
	// Test local registry with string capability IDs
	localRegistry := registrysyncer.NewLocalRegistry(
		lggr,
		func() (p2ptypes.PeerID, error) { return pid, nil },
		map[registrysyncer.DonID]registrysyncer.DON{
			registrysyncer.DonID(dID): {
				DON: capabilities.DON{
					Name:             dName,
					Families:         dFamilies,
					Config:           dConfig,
					ID:               dID,
					ConfigVersion:    uint32(2),
					F:                uint8(1),
					IsPublic:         true,
					AcceptsWorkflows: true,
					Members:          workflowDonNodes,
				},
			},
		},
		map[p2ptypes.PeerID]registrysyncer.NodeInfo{
			workflowDonNodes[0]: {
				NodeOperatorID:      1,
				Signer:              randomWord(),
				P2pID:               workflowDonNodes[0],
				EncryptionPublicKey: randomWord(),
				CapabilityIDs:       []string{"write-chain@1.0.1", "trigger@1.0.0"}, // V2 uses string IDs
			},
			workflowDonNodes[1]: {
				NodeOperatorID:      1,
				Signer:              randomWord(),
				P2pID:               workflowDonNodes[1],
				EncryptionPublicKey: randomWord(),
				CapabilityIDs:       []string{"write-chain@1.0.1", "trigger@1.0.0"}, // V2 uses string IDs
			},
			workflowDonNodes[2]: {
				NodeOperatorID:      1,
				Signer:              randomWord(),
				P2pID:               workflowDonNodes[2],
				EncryptionPublicKey: randomWord(),
				CapabilityIDs:       []string{"write-chain@1.0.1"}, // V2 uses string IDs
			},
			workflowDonNodes[3]: {
				NodeOperatorID:      1,
				Signer:              randomWord(),
				P2pID:               workflowDonNodes[3],
				EncryptionPublicKey: randomWord(),
				CapabilityIDs:       []string{"write-chain@1.0.1"}, // V2 uses string IDs
			},
		},
		map[string]registrysyncer.Capability{
			"write-chain@1.0.1": {
				CapabilityType: capabilities.CapabilityTypeTarget,
				ID:             "write-chain@1.0.1",
			},
			"trigger@1.0.0": {
				CapabilityType: capabilities.CapabilityTypeTrigger,
				ID:             "trigger@1.0.0",
			},
		},
	)

	node, err := localRegistry.LocalNode(ctx)
	require.NoError(t, err)

	don := capabilities.DON{
		ID:               dID,
		Name:             dName,
		Families:         dFamilies,
		Config:           dConfig,
		ConfigVersion:    2,
		Members:          workflowDonNodes,
		F:                1,
		IsPublic:         true,
		AcceptsWorkflows: true,
	}
	expectedNode := capabilities.Node{
		PeerID:              &pid,
		NodeOperatorID:      1,
		Signer:              localRegistry.IDsToNodes[pid].Signer,
		EncryptionPublicKey: localRegistry.IDsToNodes[pid].EncryptionPublicKey,
		WorkflowDON:         don,
		CapabilityDONs:      []capabilities.DON{don},
	}
	assert.Equal(t, expectedNode, node)

	// Test that V2 capabilities are properly handled
	assert.Len(t, localRegistry.IDsToCapabilities, 2)
	assert.Contains(t, localRegistry.IDsToCapabilities, "write-chain@1.0.1")
	assert.Contains(t, localRegistry.IDsToCapabilities, "trigger@1.0.0")

	// Test that V2 node info has string capability IDs
	nodeInfo := localRegistry.IDsToNodes[pid]
	assert.NotNil(t, nodeInfo.CapabilityIDs)
	assert.Equal(t, []string{"write-chain@1.0.1", "trigger@1.0.0"}, nodeInfo.CapabilityIDs)
}

func TestReader_V2_FamilyOperations(t *testing.T) {
	ctx := testutils.Context(t)
	lggr := logger.TestLogger(t)

	// Create a simulated backend
	owner := evmtestutils.MustNewSimTransactor(t)
	i := &big.Int{}
	oneEth, _ := i.SetString("100000000000000000000", 10)
	gasLimit := ethconfig.Defaults.Miner.GasCeil * 2

	simulatedBackend := simulated.NewBackend(gethtypes.GenesisAlloc{owner.From: {
		Balance: oneEth,
	}}, simulated.WithBlockGasLimit(gasLimit))
	simulatedBackend.Commit()

	// Deploy a V2 capabilities registry
	regAddress, _, reg, err := capabilities_registry_v2.DeployCapabilitiesRegistry(owner, simulatedBackend.Client(), capabilities_registry_v2.CapabilitiesRegistryConstructorParams{})
	require.NoError(t, err, "DeployCapabilitiesRegistry failed")
	simulatedBackend.Commit()

	// Add V2 capabilities
	writeChainCapabilityV2 := capabilities_registry_v2.CapabilitiesRegistryCapability{
		CapabilityId:          "write-chain@1.0.1",
		ConfigurationContract: common.Address{},
		Metadata:              []byte(`{"capabilityType": 3, "responseType": 1}`),
	}
	triggerCapabilityV2 := capabilities_registry_v2.CapabilitiesRegistryCapability{
		CapabilityId:          "trigger@1.0.0",
		ConfigurationContract: common.Address{},
		Metadata:              []byte(`{"capabilityType": 1, "responseType": 1}`),
	}

	_, err = reg.AddCapabilities(owner, []capabilities_registry_v2.CapabilitiesRegistryCapability{
		writeChainCapabilityV2,
		triggerCapabilityV2,
	})
	require.NoError(t, err, "AddCapabilities failed")
	simulatedBackend.Commit()

	// Add node operator
	_, err = reg.AddNodeOperators(owner, []capabilities_registry_v2.CapabilitiesRegistryNodeOperator{
		{
			Admin: owner.From,
			Name:  "TEST_NOP_V2_FAMILY",
		},
	})
	require.NoError(t, err, "Failed to add node operator")
	simulatedBackend.Commit()

	// Create different node sets for different DONs
	nodeSetA := [][32]byte{randomWord(), randomWord(), randomWord()}
	nodeSetB := [][32]byte{randomWord(), randomWord(), randomWord()}
	nodeSetC := [][32]byte{randomWord(), randomWord(), randomWord()}
	nodeSetD := [][32]byte{randomWord(), randomWord(), randomWord()}

	// Create all nodes with both capabilities
	allNodes := []capabilities_registry_v2.CapabilitiesRegistryNodeParams{}

	// Add nodes for DON A (workflow-family-a)
	for _, nodeID := range nodeSetA {
		allNodes = append(allNodes, capabilities_registry_v2.CapabilitiesRegistryNodeParams{
			NodeOperatorId:      uint32(1),
			Signer:              randomWord(),
			P2pId:               nodeID,
			EncryptionPublicKey: randomWord(),
			CsaKey:              randomWord(),
			CapabilityIds:       []string{"write-chain@1.0.1", "trigger@1.0.0"},
		})
	}

	// Add nodes for DON B (workflow-family-b)
	for _, nodeID := range nodeSetB {
		allNodes = append(allNodes, capabilities_registry_v2.CapabilitiesRegistryNodeParams{
			NodeOperatorId:      uint32(1),
			Signer:              randomWord(),
			P2pId:               nodeID,
			EncryptionPublicKey: randomWord(),
			CsaKey:              randomWord(),
			CapabilityIds:       []string{"write-chain@1.0.1", "trigger@1.0.0"},
		})
	}

	// Add nodes for DON C (multiple families)
	for _, nodeID := range nodeSetC {
		allNodes = append(allNodes, capabilities_registry_v2.CapabilitiesRegistryNodeParams{
			NodeOperatorId:      uint32(1),
			Signer:              randomWord(),
			P2pId:               nodeID,
			EncryptionPublicKey: randomWord(),
			CsaKey:              randomWord(),
			CapabilityIds:       []string{"write-chain@1.0.1", "trigger@1.0.0"},
		})
	}

	// Add nodes for DON D (no family)
	for _, nodeID := range nodeSetD {
		allNodes = append(allNodes, capabilities_registry_v2.CapabilitiesRegistryNodeParams{
			NodeOperatorId:      uint32(1),
			Signer:              randomWord(),
			P2pId:               nodeID,
			EncryptionPublicKey: randomWord(),
			CsaKey:              randomWord(),
			CapabilityIds:       []string{"write-chain@1.0.1"},
		})
	}

	_, err = reg.AddNodes(owner, allNodes)
	require.NoError(t, err, "Failed to add nodes")
	simulatedBackend.Commit()

	// Create capability configurations
	config := &capabilitiespb.CapabilityConfig{
		DefaultConfig: values.Proto(values.EmptyMap()).GetMapValue(),
		RemoteConfig: &capabilitiespb.CapabilityConfig_RemoteTriggerConfig{
			RemoteTriggerConfig: &capabilitiespb.RemoteTriggerConfig{
				RegistrationRefresh:     durationpb.New(20 * time.Second),
				RegistrationExpiry:      durationpb.New(60 * time.Second),
				MinResponsesToAggregate: uint32(1) + 1,
				MessageExpiry:           durationpb.New(120 * time.Second),
			},
		},
	}
	configb, err := proto.Marshal(config)
	require.NoError(t, err)

	cfgs := []capabilities_registry_v2.CapabilitiesRegistryCapabilityConfiguration{
		{
			CapabilityId: "write-chain@1.0.1",
			Config:       configb,
		},
		{
			CapabilityId: "trigger@1.0.0",
			Config:       configb,
		},
	}

	// Create multiple DONs with different family configurations
	newDONs := []capabilities_registry_v2.CapabilitiesRegistryNewDONParams{
		{
			Name:                     "don-family-a",
			DonFamilies:              []string{"workflow-family-a"},
			Config:                   []byte("config-family-a"),
			CapabilityConfigurations: cfgs,
			Nodes:                    nodeSetA,
			F:                        1,
			IsPublic:                 true,
			AcceptsWorkflows:         true,
		},
		{
			Name:                     "don-family-b",
			DonFamilies:              []string{"workflow-family-b"},
			Config:                   []byte("config-family-b"),
			CapabilityConfigurations: cfgs,
			Nodes:                    nodeSetB,
			F:                        1,
			IsPublic:                 true,
			AcceptsWorkflows:         true,
		},
		{
			Name:                     "don-multi-family",
			DonFamilies:              []string{"workflow-family-a", "workflow-family-c"},
			Config:                   []byte("config-multi-family"),
			CapabilityConfigurations: cfgs,
			Nodes:                    nodeSetC,
			F:                        1,
			IsPublic:                 true,
			AcceptsWorkflows:         true,
		},
		{
			Name:        "don-no-family",
			DonFamilies: []string{}, // empty families array
			Config:      []byte("config-no-family"),
			CapabilityConfigurations: []capabilities_registry_v2.CapabilitiesRegistryCapabilityConfiguration{
				{
					CapabilityId: "write-chain@1.0.1",
					Config:       configb,
				},
			},
			Nodes:            nodeSetD,
			F:                1,
			IsPublic:         true,
			AcceptsWorkflows: true,
		},
	}
	_, err = reg.AddDONs(owner, newDONs)
	require.NoError(t, err)
	simulatedBackend.Commit()

	// Set up syncer and reader
	db := pgtest.NewSqlxDB(t)
	factory := newContractReaderFactory(t, simulatedBackend)
	syncerORM := registrysyncer.NewORM(db, lggr)
	syncer, err := registrysyncer_v2.New(lggr, func() (p2ptypes.PeerID, error) { return p2ptypes.PeerID{}, nil }, factory, regAddress.Hex(), syncerORM)
	require.NoError(t, err)

	l := &launcher{}
	syncer.AddListener(l)

	err = syncer.Sync(ctx, false)
	require.NoError(t, err)

	s := l.localRegistry
	require.NotNil(t, s)

	// Verify we have 4 DONs
	assert.Len(t, s.IDsToDONs, 4)

	// Create a V2 reader to test family operations directly
	// First, we need to create a properly configured contract reader
	contractReaderConfig := evmrelaytypes.ChainReaderConfig{
		Contracts: map[string]evmrelaytypes.ChainContractReader{
			"CapabilitiesRegistry": {
				ContractABI: capabilities_registry_v2.CapabilitiesRegistryABI,
				Configs: map[string]*evmrelaytypes.ChainReaderDefinition{
					"getDONs": {
						ChainSpecificName: "getDONs",
					},
					"getCapabilities": {
						ChainSpecificName: "getCapabilities",
					},
					"getNodes": {
						ChainSpecificName: "getNodes",
					},
					"getDONsInFamily": {
						ChainSpecificName: "getDONsInFamily",
					},
					"getHistoricalDONInfo": {
						ChainSpecificName: "getHistoricalDONInfo",
					},
					"getNode": {
						ChainSpecificName: "getNode",
					},
					"getNodeOperator": {
						ChainSpecificName: "getNodeOperator",
					},
					"getNodeOperators": {
						ChainSpecificName: "getNodeOperators",
					},
					"getNodesByP2PIds": {
						ChainSpecificName: "getNodesByP2PIds",
					},
					"isCapabilityDeprecated": {
						ChainSpecificName: "isCapabilityDeprecated",
					},
				},
			},
		},
	}
	contractReaderConfigEncoded, err := json.Marshal(contractReaderConfig)
	require.NoError(t, err)

	contractReader, err := factory.NewContractReader(ctx, contractReaderConfigEncoded)
	require.NoError(t, err)

	// Bind the contract
	err = contractReader.Bind(ctx, []types.BoundContract{
		{
			Address: regAddress.Hex(),
			Name:    "CapabilitiesRegistry",
		},
	})
	require.NoError(t, err)

	capabilitiesRegistryReader := CapabilitiesRegistryReader{
		boundedContract: types.BoundContract{Address: regAddress.Hex(), Name: "CapabilitiesRegistry"},
		contractReader:  contractReader,
	}

	// Test GetDONsInFamily functionality
	t.Run("GetDONsInFamily_SingleFamily", func(t *testing.T) {
		// Query "workflow-family-a" -> should return DON IDs [1, 3] (don-family-a and don-multi-family)
		familyADONs, err := capabilitiesRegistryReader.GetDONsInFamily(ctx, "workflow-family-a")
		require.NoError(t, err)
		require.Len(t, familyADONs, 2, "Should find 2 DONs in workflow-family-a")
		assert.Contains(t, familyADONs, *big.NewInt(1), "DON 1 should be in workflow-family-a")
		assert.Contains(t, familyADONs, *big.NewInt(3), "DON 3 should be in workflow-family-a")

		// Query "workflow-family-b" -> should return DON ID [2]
		familyBDONs, err := capabilitiesRegistryReader.GetDONsInFamily(ctx, "workflow-family-b")
		require.NoError(t, err)
		require.Len(t, familyBDONs, 1, "Should find 1 DON in workflow-family-b")
		assert.Contains(t, familyBDONs, *big.NewInt(2), "DON 2 should be in workflow-family-b")

		// Query "workflow-family-c" -> should return DON ID [3]
		familyCDONs, err := capabilitiesRegistryReader.GetDONsInFamily(ctx, "workflow-family-c")
		require.NoError(t, err)
		require.Len(t, familyCDONs, 1, "Should find 1 DON in workflow-family-c")
		assert.Contains(t, familyCDONs, *big.NewInt(3), "DON 3 should be in workflow-family-c")
	})

	t.Run("GetDONsInFamily_NonExistentFamily", func(t *testing.T) {
		// Query "non-existent-family" -> should return empty
		nonExistentDONs, err := capabilitiesRegistryReader.GetDONsInFamily(ctx, "non-existent-family")
		require.NoError(t, err)
		assert.Empty(t, nonExistentDONs, "Non-existent family should return empty")
	})

	t.Run("GetHistoricalDONInfo_FamilyData", func(t *testing.T) {
		// Test GetHistoricalDONInfo with configCount=1 for each DON
		for donID := uint32(1); donID <= 4; donID++ {
			historicalDON, err := capabilitiesRegistryReader.GetHistoricalDONInfo(ctx, donID, 1)
			require.NoError(t, err, "GetHistoricalDONInfo should work for DON %d", donID)
			require.NotNil(t, historicalDON, "Historical DON info should not be nil for DON %d", donID)

			// Verify family data is preserved in historical records
			switch donID {
			case 1:
				assert.Equal(t, "don-family-a", historicalDON.Name)
				assert.Equal(t, []string{"workflow-family-a"}, historicalDON.DonFamilies)
			case 2:
				assert.Equal(t, "don-family-b", historicalDON.Name)
				assert.Equal(t, []string{"workflow-family-b"}, historicalDON.DonFamilies)
			case 3:
				assert.Equal(t, "don-multi-family", historicalDON.Name)
				assert.Equal(t, []string{"workflow-family-a", "workflow-family-c"}, historicalDON.DonFamilies)
			case 4:
				assert.Equal(t, "don-no-family", historicalDON.Name)
				if historicalDON.DonFamilies != nil {
					assert.Empty(t, historicalDON.DonFamilies, "DON 4 should have empty families")
				}
			}
		}
	})

	t.Run("GetHistoricalDONInfo_InvalidConfigCount", func(t *testing.T) {
		// Test with invalid configCount -> should handle gracefully
		_, err := capabilitiesRegistryReader.GetHistoricalDONInfo(ctx, 1, 999)
		require.Error(t, err, "Should return error for invalid configCount")
		assert.ErrorContains(t, err, "invalid type: contract call: execution reverted")
	})

	// Verify synced data integrity in local registry
	t.Run("LocalRegistry_FamilyIntegrity", func(t *testing.T) {
		// Verify each DON has correct family data
		don1 := s.IDsToDONs[1].DON
		assert.Equal(t, "don-family-a", don1.Name)
		assert.Equal(t, []string{"workflow-family-a"}, don1.Families)
		assert.Equal(t, []byte("config-family-a"), don1.Config)

		don2 := s.IDsToDONs[2].DON
		assert.Equal(t, "don-family-b", don2.Name)
		assert.Equal(t, []string{"workflow-family-b"}, don2.Families)
		assert.Equal(t, []byte("config-family-b"), don2.Config)

		don3 := s.IDsToDONs[3].DON
		assert.Equal(t, "don-multi-family", don3.Name)
		assert.Equal(t, []string{"workflow-family-a", "workflow-family-c"}, don3.Families)
		assert.Equal(t, []byte("config-multi-family"), don3.Config)

		don4 := s.IDsToDONs[4].DON
		assert.Equal(t, "don-no-family", don4.Name)
		assert.Empty(t, don4.Families)
		assert.Equal(t, []byte("config-no-family"), don4.Config)
	})

	// Verify capabilities are correctly synced
	assert.Len(t, s.IDsToCapabilities, 2)
	assert.Contains(t, s.IDsToCapabilities, "write-chain@1.0.1")
	assert.Contains(t, s.IDsToCapabilities, "trigger@1.0.0")

	// Verify all nodes are correctly synced
	assert.Len(t, s.IDsToNodes, 12) // 3 nodes * 4 DONs = 12 nodes
}

type CapabilitiesRegistryReader struct {
	contractReader  types.ContractReader
	boundedContract types.BoundContract
}

func (r *CapabilitiesRegistryReader) GetDONsInFamily(ctx context.Context, family string) ([]big.Int, error) {
	var familyADONs []big.Int
	err := r.contractReader.GetLatestValue(
		ctx,
		r.boundedContract.ReadIdentifier("getDONsInFamily"),
		primitives.Unconfirmed,
		map[string]any{
			"donFamily": family,
		},
		&familyADONs,
	)
	return familyADONs, err
}

func (r *CapabilitiesRegistryReader) GetHistoricalDONInfo(ctx context.Context, donID uint32, configCount uint32) (*capabilities_registry_v2.CapabilitiesRegistryDONInfo, error) {
	var historicalDON capabilities_registry_v2.CapabilitiesRegistryDONInfo
	err := r.contractReader.GetLatestValue(
		ctx,
		r.boundedContract.ReadIdentifier("getHistoricalDONInfo"),
		primitives.Unconfirmed,
		map[string]any{
			"donId":       donID,
			"configCount": configCount,
		},
		&historicalDON,
	)
	return &historicalDON, err
}
