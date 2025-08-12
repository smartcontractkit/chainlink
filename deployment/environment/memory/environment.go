package memory

import (
	"context"
	"slices"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/smartcontractkit/freeport"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	cldf_evm "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink/deployment"
)

const (
	Memory = "memory"
)

type MemoryEnvironmentConfig struct {
	Chains             int
	SolChains          int
	AptosChains        int
	ZkChains           int
	TonChains          int
	TronChains         int
	NumOfUsersPerChain int
	Nodes              int
	Bootstraps         int
	RegistryConfig     deployment.CapabilityRegistryConfig
	CustomDBSetup      []string // SQL queries to run after DB creation

	// Solana Handle different contract versions
	CCIPSolanaContractVersion CCIPSolanaContractVersion
}

// TODO: This shouldn't be duplicated from solana_changesets_V0_1_1/utils.go
// This is a temporary solution to avoid circular dependencies.
// We should refactor the code to avoid this duplication.
type CCIPSolanaContractVersion string

const (
	SolanaContractV0_1_0 CCIPSolanaContractVersion = "v0.1.0"
	SolanaContractV0_1_1 CCIPSolanaContractVersion = "v0.1.1"
)

var ContractVersionShortSha = map[CCIPSolanaContractVersion]string{
	SolanaContractV0_1_0: "0ee732e80586",
	SolanaContractV0_1_1: "7f8a0f403c3a",
}

type NewNodesConfig struct {
	LogLevel zapcore.Level
	// BlockChains to be configured
	BlockChains    cldf_chain.BlockChains
	NumNodes       int
	NumBootstraps  int
	RegistryConfig deployment.CapabilityRegistryConfig
	// SQL queries to run after DB creation, typically used for setting up testing state. Optional.
	CustomDBSetup []string
}

// For placeholders like aptos
func NewMemoryChain(t *testing.T, selector uint64) cldf_evm.Chain {
	return cldf_evm.Chain{
		Selector:    selector,
		Client:      nil,
		DeployerKey: &bind.TransactOpts{},
		Confirm: func(tx *types.Transaction) (uint64, error) {
			return 0, nil
		},
	}
}

func NewMemoryChainsEVM(t *testing.T, numChains int, numUsers int) []cldf_chain.BlockChain {
	t.Helper()

	return generateChainsEVM(t, numChains, numUsers)
}

func NewMemoryChainsEVMWithChainIDs(
	t *testing.T, chainIDs []uint64, numUsers int,
) []cldf_chain.BlockChain {
	t.Helper()

	return generateChainsEVMWithIDs(t, chainIDs, numUsers)
}

func NewMemoryChainsSol(t *testing.T, numChains int, commitSha string) []cldf_chain.BlockChain {
	return generateChainsSol(t, numChains, commitSha)
}

func NewMemoryChainsAptos(t *testing.T, numChains int) []cldf_chain.BlockChain {
	return generateChainsAptos(t, numChains)
}

func NewMemoryChainsZk(t *testing.T, numChains int) []cldf_chain.BlockChain {
	return GenerateChainsZk(t, numChains)
}

func NewMemoryChainsTon(t *testing.T, numChains int) []cldf_chain.BlockChain {
	return generateChainsTon(t, numChains)
}

func NewMemoryChainsTron(t *testing.T, numChains int) []cldf_chain.BlockChain {
	return generateChainsTron(t, numChains)
}

func NewNodes(
	t *testing.T,
	cfg NewNodesConfig,
	configOpts ...ConfigOpt,
) map[string]Node {
	nodesByPeerID := make(map[string]Node)
	if cfg.NumNodes+cfg.NumBootstraps == 0 {
		return nodesByPeerID
	}
	ports := freeport.GetN(t, cfg.NumNodes+cfg.NumBootstraps)
	// bootstrap nodes must be separate nodes from plugin nodes,
	// since we won't run a bootstrapper and a plugin oracle on the same
	// chainlink node in production.
	for i := 0; i < cfg.NumBootstraps; i++ {
		// TODO: bootstrap nodes don't have to support anything other than the home chain.
		// We should remove all non-home chains from the config below and make sure things
		// run smoothly.
		c := NewNodeConfig{
			Port:           ports[i],
			BlockChains:    cfg.BlockChains,
			LogLevel:       cfg.LogLevel,
			Bootstrap:      true,
			RegistryConfig: cfg.RegistryConfig,
			CustomDBSetup:  cfg.CustomDBSetup,
		}
		node := NewNode(t, c, configOpts...)
		nodesByPeerID[node.Keys.PeerID.String()] = *node
		// Note in real env, this ID is allocated by JD.
	}
	for i := range cfg.NumNodes {
		c := NewNodeConfig{
			Port:           ports[cfg.NumBootstraps+i],
			BlockChains:    cfg.BlockChains,
			LogLevel:       cfg.LogLevel,
			Bootstrap:      false,
			RegistryConfig: cfg.RegistryConfig,
			CustomDBSetup:  cfg.CustomDBSetup,
		}
		// grab port offset by numBootstraps, since above loop also takes some ports.
		node := NewNode(t, c, configOpts...)
		nodesByPeerID[node.Keys.PeerID.String()] = *node
		// Note in real env, this ID is allocated by JD.
	}
	return nodesByPeerID
}

func NewMemoryEnvironmentFromChainsNodes(
	ctx func() context.Context,
	lggr logger.Logger,
	blockchains cldf_chain.BlockChains,
	nodes map[string]Node,
) cldf.Environment {
	var nodeIDs []string
	for id := range nodes {
		nodeIDs = append(nodeIDs, id)

	}

	return *cldf.NewEnvironment(
		Memory,
		lggr,
		cldf.NewMemoryAddressBook(),
		datastore.NewMemoryDataStore().Seal(),
		nodeIDs, // Note these have the p2p_ prefix.
		NewMemoryJobClient(nodes),
		ctx,
		cldf.XXXGenerateTestOCRSecrets(),
		blockchains,
	)
}

// To be used by tests and any kind of deployment logic.
func NewMemoryEnvironment(
	t *testing.T,
	lggr logger.Logger,
	logLevel zapcore.Level,
	config MemoryEnvironmentConfig,
) cldf.Environment {
	evmChains := NewMemoryChainsEVM(t, config.Chains, config.NumOfUsersPerChain)

	var solanaCommitSha string
	ccipContractVersion := config.CCIPSolanaContractVersion
	if ccipContractVersion == SolanaContractV0_1_1 {
		solanaCommitSha = ContractVersionShortSha[ccipContractVersion]
	} else {
		solanaCommitSha = ""
	}
	solChains := NewMemoryChainsSol(t, config.SolChains, solanaCommitSha)
	aptosChains := NewMemoryChainsAptos(t, config.AptosChains)
	zkChains := NewMemoryChainsZk(t, config.ZkChains)
	tonChains := NewMemoryChainsTon(t, config.TonChains)
	tronChains := NewMemoryChainsTron(t, config.TronChains)

	chains := cldf_chain.NewBlockChainsFromSlice(
		slices.Concat(evmChains, solChains, aptosChains, zkChains, tonChains, tronChains),
	)

	c := NewNodesConfig{
		LogLevel:       logLevel,
		BlockChains:    chains,
		NumNodes:       config.Nodes,
		NumBootstraps:  config.Bootstraps,
		RegistryConfig: config.RegistryConfig,
		CustomDBSetup:  config.CustomDBSetup,
	}
	nodes := NewNodes(t, c)
	var nodeIDs []string
	for id, node := range nodes {
		require.NoError(t, node.App.Start(t.Context()))
		t.Cleanup(func() {
			require.NoError(t, node.App.Stop())
		})
		nodeIDs = append(nodeIDs, id)
	}

	return *cldf.NewEnvironment(
		Memory,
		lggr,
		cldf.NewMemoryAddressBook(),
		datastore.NewMemoryDataStore().Seal(),
		nodeIDs,
		NewMemoryJobClient(nodes),
		t.Context,
		cldf.XXXGenerateTestOCRSecrets(),
		chains,
	)
}
