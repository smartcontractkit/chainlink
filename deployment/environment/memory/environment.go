package memory

import (
	"context"
	"path/filepath"
	"runtime"
	"slices"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/smartcontractkit/freeport"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	cldf_aptos "github.com/smartcontractkit/chainlink-deployments-framework/chain/aptos"
	cldf_evm "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
	cldf_solana "github.com/smartcontractkit/chainlink-deployments-framework/chain/solana"
	cldf_ton "github.com/smartcontractkit/chainlink-deployments-framework/chain/ton"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink/deployment"
)

const (
	Memory = "memory"
)

var (
	// Instead of a relative path, use runtime.Caller or go-bindata
	ProgramsPath = GetProgramsPath()
)

func GetProgramsPath() string {
	// Get the directory of the current file (environment.go)
	_, currentFile, _, _ := runtime.Caller(0)
	// Go up to the root of the deployment package
	rootDir := filepath.Dir(filepath.Dir(filepath.Dir(currentFile)))
	// Construct the absolute path
	return filepath.Join(rootDir, "ccip/changeset/internal", "solana_contracts")
}

type MemoryEnvironmentConfig struct {
	Chains             int
	SolChains          int
	AptosChains        int
	ZkChains           int
	TonChains          int
	NumOfUsersPerChain int
	Nodes              int
	Bootstraps         int
	RegistryConfig     deployment.CapabilityRegistryConfig
	CustomDBSetup      []string // SQL queries to run after DB creation
}

type NewNodesConfig struct {
	LogLevel zapcore.Level
	// EVM chains to be configured. Optional.
	Chains map[uint64]cldf_evm.Chain
	// Solana chains to be configured. Optional.
	SolChains map[uint64]cldf_solana.Chain
	// Aptos chains to be configured. Optional.
	AptosChains map[uint64]cldf_aptos.Chain
	// TON chains to be configured. Optional.
	TonChains      map[uint64]cldf_ton.Chain
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

func NewMemoryChainsSol(t *testing.T, numChains int) []cldf_chain.BlockChain {
	return generateChainsSol(t, numChains)
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
			Chains:         cfg.Chains,
			Solchains:      cfg.SolChains,
			Aptoschains:    cfg.AptosChains,
			Tonchains:      cfg.TonChains,
			LogLevel:       cfg.LogLevel,
			Bootstrap:      true,
			RegistryConfig: cfg.RegistryConfig,
			CustomDBSetup:  cfg.CustomDBSetup,
		}
		node := NewNode(t, c, configOpts...)
		nodesByPeerID[node.Keys.PeerID.String()] = *node
		// Note in real env, this ID is allocated by JD.
	}
	for i := 0; i < cfg.NumNodes; i++ {
		c := NewNodeConfig{
			Port:           ports[cfg.NumBootstraps+i],
			Chains:         cfg.Chains,
			Solchains:      cfg.SolChains,
			Aptoschains:    cfg.AptosChains,
			Tonchains:      cfg.TonChains,
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
	evmChains map[uint64]cldf_evm.Chain,
	solChains map[uint64]cldf_solana.Chain,
	aptosChains map[uint64]cldf_aptos.Chain,
	tonChains map[uint64]cldf_ton.Chain,
	nodes map[string]Node,
) cldf.Environment {
	var nodeIDs []string
	for id := range nodes {
		nodeIDs = append(nodeIDs, id)

	}

	blockChains := map[uint64]cldf_chain.BlockChain{}
	for _, c := range evmChains {
		blockChains[c.Selector] = c
	}
	for _, c := range solChains {
		blockChains[c.Selector] = c
	}
	for _, c := range aptosChains {
		blockChains[c.Selector] = c
	}
	for _, c := range tonChains {
		blockChains[c.Selector] = c
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
		cldf_chain.NewBlockChains(blockChains),
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
	solChains := NewMemoryChainsSol(t, config.SolChains)
	aptosChains := NewMemoryChainsAptos(t, config.AptosChains)
	zkChains := NewMemoryChainsZk(t, config.ZkChains)
	tonChains := NewMemoryChainsTon(t, config.TonChains)

	chains := cldf_chain.NewBlockChainsFromSlice(
		slices.Concat(evmChains, solChains, aptosChains, zkChains, tonChains),
	)

	c := NewNodesConfig{
		LogLevel:       logLevel,
		Chains:         chains.EVMChains(),
		SolChains:      chains.SolanaChains(),
		AptosChains:    chains.AptosChains(),
		TonChains:      chains.TonChains(),
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
