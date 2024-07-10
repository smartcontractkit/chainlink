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

	"github.com/smartcontractkit/chainlink-common/pkg/types"

	evmclient "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/headtracker"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	kcr "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/capabilities_registry"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm"
	evmrelaytypes "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"
)

var writeChainCapability = kcr.CapabilitiesRegistryCapability{
	LabelledName: "write-chain",
	Version:      "1.0.1",
	ResponseType: uint8(1),
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
	logPoller logpoller.LogPoller
	client    evmclient.Client
}

func (c *crFactory) NewContractReader(ctx context.Context, cfg []byte) (types.ContractReader, error) {
	crCfg := &evmrelaytypes.ChainReaderConfig{}
	if err := json.Unmarshal(cfg, crCfg); err != nil {
		return nil, err
	}
	svc, err := evm.NewChainReaderService(ctx, c.lggr, c.logPoller, c.client, *crCfg)
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

func TestReader_Integration(t *testing.T) {
	ctx := testutils.Context(t)
	reg, regAddress, owner, sim := startNewChainWithRegistry(t)

	_, err := reg.AddCapabilities(owner, []kcr.CapabilitiesRegistryCapability{writeChainCapability})
	require.NoError(t, err, "AddCapability failed for %s", writeChainCapability.LabelledName)
	sim.Commit()

	cid, err := reg.GetHashedCapabilityId(&bind.CallOpts{}, writeChainCapability.LabelledName, writeChainCapability.Version)
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
			HashedCapabilityIds: [][32]byte{cid},
		},
		{
			// The first NodeOperatorId has id 1 since the id is auto-incrementing.
			NodeOperatorId:      uint32(1),
			Signer:              signersSet[1],
			P2pId:               nodeSet[1],
			HashedCapabilityIds: [][32]byte{cid},
		},
		{
			// The first NodeOperatorId has id 1 since the id is auto-incrementing.
			NodeOperatorId:      uint32(1),
			Signer:              signersSet[2],
			P2pId:               nodeSet[2],
			HashedCapabilityIds: [][32]byte{cid},
		},
	}
	_, err = reg.AddNodes(owner, nodes)
	require.NoError(t, err)

	cfgs := []kcr.CapabilitiesRegistryCapabilityConfiguration{
		{
			CapabilityId: cid,
			Config:       []byte(`{"hello": "world"}`),
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

	factory := newContractReaderFactory(t, sim)
	reader, err := New(logger.TestLogger(t), factory, regAddress.Hex())
	require.NoError(t, err)

	s, err := reader.state(ctx)
	require.NoError(t, err)
	assert.Len(t, s.IDsToCapabilities, 1)

	gotCap := s.IDsToCapabilities[cid]
	assert.Equal(t, kcr.CapabilitiesRegistryCapabilityInfo{
		HashedId:       cid,
		LabelledName:   "write-chain",
		Version:        "1.0.1",
		ResponseType:   uint8(1),
		CapabilityType: uint8(0),
		IsDeprecated:   false,
	}, gotCap)

	assert.Len(t, s.IDsToDONs, 1)
	assert.Equal(t, kcr.CapabilitiesRegistryDONInfo{
		Id:                       1, // initial Id
		ConfigCount:              1, // initial Count
		IsPublic:                 true,
		AcceptsWorkflows:         true,
		F:                        1,
		NodeP2PIds:               nodeSet,
		CapabilityConfigurations: cfgs,
	}, s.IDsToDONs[1])

	nodesInfo := []kcr.CapabilitiesRegistryNodeInfo{
		{
			// The first NodeOperatorId has id 1 since the id is auto-incrementing.
			NodeOperatorId:      uint32(1),
			ConfigCount:         1,
			WorkflowDONId:       1,
			Signer:              signersSet[0],
			P2pId:               nodeSet[0],
			HashedCapabilityIds: [][32]byte{cid},
			CapabilitiesDONIds:  []*big.Int{},
		},
		{
			// The first NodeOperatorId has id 1 since the id is auto-incrementing.
			NodeOperatorId:      uint32(1),
			ConfigCount:         1,
			WorkflowDONId:       1,
			Signer:              signersSet[1],
			P2pId:               nodeSet[1],
			HashedCapabilityIds: [][32]byte{cid},
			CapabilitiesDONIds:  []*big.Int{},
		},
		{
			// The first NodeOperatorId has id 1 since the id is auto-incrementing.
			NodeOperatorId:      uint32(1),
			ConfigCount:         1,
			WorkflowDONId:       1,
			Signer:              signersSet[2],
			P2pId:               nodeSet[2],
			HashedCapabilityIds: [][32]byte{cid},
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
