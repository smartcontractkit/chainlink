package memory

import (
	"context"
	"fmt"
	"path/filepath"
	"runtime"
	"strconv"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/gagliardetto/solana-go"
	solRpc "github.com/gagliardetto/solana-go/rpc"
	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/freeport"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	solCommonUtil "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/common"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	cldf_aptos "github.com/smartcontractkit/chainlink-deployments-framework/chain/aptos"
	cldf_evm "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
	cldf_solana "github.com/smartcontractkit/chainlink-deployments-framework/chain/solana"
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
	AptosChains    map[uint64]cldf_aptos.Chain
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

// Needed for environment variables on the node which point to prexisitng addresses.
// i.e. CapReg.
func NewMemoryChains(t *testing.T, numChains int, numUsers int) (map[uint64]cldf_evm.Chain, map[uint64][]*bind.TransactOpts) {
	mchains := GenerateChains(t, numChains, numUsers)
	users := make(map[uint64][]*bind.TransactOpts)
	for id, chain := range mchains {
		sel, err := chainsel.SelectorFromChainId(id)
		require.NoError(t, err)
		users[sel] = chain.Users
	}
	return generateMemoryChain(t, mchains), users
}

func NewMemoryChainsSol(t *testing.T, numChains int) map[uint64]cldf_solana.Chain {
	mchains := GenerateChainsSol(t, numChains)
	return generateMemoryChainSol(mchains)
}

func NewMemoryChainsAptos(t *testing.T, numChains int) []cldf_chain.BlockChain {
	return generateChainsAptos(t, numChains)
}

func NewMemoryChainsZk(t *testing.T, numChains int) map[uint64]cldf_evm.Chain {
	return GenerateChainsZk(t, numChains)
}

func NewMemoryChainsWithChainIDs(t *testing.T, chainIDs []uint64, numUsers int) (map[uint64]cldf_evm.Chain, map[uint64][]*bind.TransactOpts) {
	mchains := GenerateChainsWithIds(t, chainIDs, numUsers)
	users := make(map[uint64][]*bind.TransactOpts)
	for id, chain := range mchains {
		sel, err := chainsel.SelectorFromChainId(id)
		require.NoError(t, err)
		users[sel] = chain.Users
	}
	return generateMemoryChain(t, mchains), users
}

func generateMemoryChain(t *testing.T, inputs map[uint64]EVMChain) map[uint64]cldf_evm.Chain {
	chains := make(map[uint64]cldf_evm.Chain)
	for cid, chain := range inputs {
		chain := chain
		chainInfo, err := chainsel.GetChainDetailsByChainIDAndFamily(strconv.FormatUint(cid, 10), chainsel.FamilyEVM)
		require.NoError(t, err)
		backend := NewBackend(chain.Backend)
		chains[chainInfo.ChainSelector] = cldf_evm.Chain{
			Selector:    chainInfo.ChainSelector,
			Client:      backend,
			DeployerKey: chain.DeployerKey,
			Confirm: func(tx *types.Transaction) (uint64, error) {
				if tx == nil {
					return 0, fmt.Errorf("tx was nil, nothing to confirm, chain %s", chainInfo.ChainName)
				}
				for {
					backend.Commit()
					receipt, err := func() (*types.Receipt, error) {
						ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
						defer cancel()
						return bind.WaitMined(ctx, backend, tx)
					}()
					if err != nil {
						return 0, fmt.Errorf("tx %s failed to confirm: %w, chain %d", tx.Hash().Hex(), err, chainInfo.ChainSelector)
					}
					if receipt.Status == 0 {
						errReason, err := deployment.GetErrorReasonFromTx(chain.Backend.Client(), chain.DeployerKey.From, tx, receipt)
						if err == nil && errReason != "" {
							return 0, fmt.Errorf("tx %s reverted,error reason: %s chain %s", tx.Hash().Hex(), errReason, chainInfo.ChainName)
						}
						return 0, fmt.Errorf("tx %s reverted, could not decode error reason chain %s", tx.Hash().Hex(), chainInfo.ChainName)
					}
					return receipt.BlockNumber.Uint64(), nil
				}
			},
			Users: chain.Users,
		}
	}
	return chains
}

func generateMemoryChainSol(inputs map[uint64]SolanaChain) map[uint64]cldf_solana.Chain {
	chains := make(map[uint64]cldf_solana.Chain)
	for cid, chain := range inputs {
		chain := chain
		chains[cid] = cldf_solana.Chain{
			Selector:     cid,
			Client:       chain.Client,
			DeployerKey:  &chain.DeployerKey,
			URL:          chain.URL,
			WSURL:        chain.WSURL,
			KeypairPath:  chain.KeypairPath,
			ProgramsPath: ProgramsPath,
			Confirm: func(instructions []solana.Instruction, opts ...solCommonUtil.TxModifier) error {
				_, err := solCommonUtil.SendAndConfirm(
					context.Background(), chain.Client, instructions, chain.DeployerKey, solRpc.CommitmentConfirmed, opts...,
				)
				return err
			},
		}
	}
	return chains
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
	chains map[uint64]cldf_evm.Chain,
	solChains map[uint64]cldf_solana.Chain,
	aptosChains map[uint64]cldf_aptos.Chain,
	nodes map[string]Node,
) cldf.Environment {
	var nodeIDs []string
	for id := range nodes {
		nodeIDs = append(nodeIDs, id)
	}

	blockChains := map[uint64]cldf_chain.BlockChain{}
	for _, c := range chains {
		blockChains[c.Selector] = c
	}
	for _, c := range solChains {
		blockChains[c.Selector] = c
	}
	for _, c := range aptosChains {
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
func NewMemoryEnvironment(t *testing.T, lggr logger.Logger, logLevel zapcore.Level, config MemoryEnvironmentConfig) cldf.Environment {
	chains, _ := NewMemoryChains(t, config.Chains, config.NumOfUsersPerChain)
	solChains := NewMemoryChainsSol(t, config.SolChains)
	aptosChains := NewMemoryChainsAptos(t, config.AptosChains)
	zkChains := NewMemoryChainsZk(t, config.ZkChains)
	for chainSel, chain := range zkChains {
		chains[chainSel] = chain
	}

	// Convert aptos chains to concrete type to pass to the NewNodesConfig.
	// This is a temporary workaround until we have a better way to handle bringing up configuring
	// the memory nodes.
	concreteAptosChains := cldf_chain.NewBlockChainsFromSlice(aptosChains).AptosChains()

	c := NewNodesConfig{
		LogLevel:       logLevel,
		Chains:         chains,
		SolChains:      solChains,
		AptosChains:    concreteAptosChains,
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

	blockChains := map[uint64]cldf_chain.BlockChain{}
	for _, c := range chains {
		blockChains[c.Selector] = c
	}
	for _, c := range solChains {
		blockChains[c.Selector] = c
	}
	for _, c := range aptosChains {
		blockChains[c.ChainSelector()] = c
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
		cldf_chain.NewBlockChains(blockChains),
	)
}
