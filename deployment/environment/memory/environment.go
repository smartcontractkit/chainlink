package memory

import (
	"context"
	"fmt"
	"math/rand"
	"path/filepath"
	"runtime"
	"strconv"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/gagliardetto/solana-go"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
	"golang.org/x/exp/maps"

	cldf_aptos "github.com/smartcontractkit/chainlink-deployments-framework/chain/aptos"

	"github.com/smartcontractkit/freeport"

	chainsel "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-deployments-framework/chain"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment"

	solRpc "github.com/gagliardetto/solana-go/rpc"

	solCommonUtil "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/common"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
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
	Chains map[uint64]cldf.Chain
	// Solana chains to be configured. Optional.
	SolChains map[uint64]cldf.SolChain
	// Aptos chains to be configured. Optional.
	AptosChains    map[uint64]cldf_aptos.Chain
	NumNodes       int
	NumBootstraps  int
	RegistryConfig deployment.CapabilityRegistryConfig
	// SQL queries to run after DB creation, typically used for setting up testing state. Optional.
	CustomDBSetup []string
	// ChainTopology is the chain-node topology of the role DON.
	ChainTopology ChainTopology
	// HomeChainSel is the chain selector of the home chain.
	HomeChainSel uint64
}

// For placeholders like aptos
func NewMemoryChain(t *testing.T, selector uint64) cldf.Chain {
	return cldf.Chain{
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
func NewMemoryChains(t *testing.T, numChains int, numUsers int) (map[uint64]cldf.Chain, map[uint64][]*bind.TransactOpts) {
	mchains := GenerateChains(t, numChains, numUsers)
	users := make(map[uint64][]*bind.TransactOpts)
	for id, chain := range mchains {
		sel, err := chainsel.SelectorFromChainId(id)
		require.NoError(t, err)
		users[sel] = chain.Users
	}
	return generateMemoryChain(t, mchains), users
}

func NewMemoryChainsSol(t *testing.T, numChains int) map[uint64]cldf.SolChain {
	mchains := GenerateChainsSol(t, numChains)
	return generateMemoryChainSol(mchains)
}

func NewMemoryChainsAptos(t *testing.T, numChains int) map[uint64]cldf_aptos.Chain {
	return GenerateChainsAptos(t, numChains)
}

func NewMemoryChainsZk(t *testing.T, numChains int) map[uint64]cldf.Chain {
	return GenerateChainsZk(t, numChains)
}

func NewMemoryChainsWithChainIDs(t *testing.T, chainIDs []uint64, numUsers int) (map[uint64]cldf.Chain, map[uint64][]*bind.TransactOpts) {
	mchains := GenerateChainsWithIds(t, chainIDs, numUsers)
	users := make(map[uint64][]*bind.TransactOpts)
	for id, chain := range mchains {
		sel, err := chainsel.SelectorFromChainId(id)
		require.NoError(t, err)
		users[sel] = chain.Users
	}
	return generateMemoryChain(t, mchains), users
}

func generateMemoryChain(t *testing.T, inputs map[uint64]EVMChain) map[uint64]cldf.Chain {
	chains := make(map[uint64]cldf.Chain)
	for cid, chain := range inputs {
		chain := chain
		chainInfo, err := chainsel.GetChainDetailsByChainIDAndFamily(strconv.FormatUint(cid, 10), chainsel.FamilyEVM)
		require.NoError(t, err)
		backend := NewBackend(chain.Backend)
		chains[chainInfo.ChainSelector] = cldf.Chain{
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

func generateMemoryChainSol(inputs map[uint64]SolanaChain) map[uint64]cldf.SolChain {
	chains := make(map[uint64]cldf.SolChain)
	for cid, chain := range inputs {
		chain := chain
		chains[cid] = cldf.SolChain{
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
	for i := range cfg.NumBootstraps {
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

	var newNodeConfigs []NewNodeConfig
	if cfg.ChainTopology.FChainToNumChains != nil {
		newNodeConfigs = createNewNodeConfigsWithChainTopology(t, cfg, ports)
	}
	for i := range cfg.NumNodes {
		c := NewNodeConfig{
			// grab port offset by NumBootstraps, since above loop also takes some ports.
			Port:           ports[cfg.NumBootstraps+i],
			Chains:         cfg.Chains,
			Solchains:      cfg.SolChains,
			Aptoschains:    cfg.AptosChains,
			LogLevel:       cfg.LogLevel,
			Bootstrap:      false,
			RegistryConfig: cfg.RegistryConfig,
			CustomDBSetup:  cfg.CustomDBSetup,
		}
		// if chain topology is set, use the new node config from the chain topology,
		// which may include a different set of chains to support.
		if cfg.ChainTopology.FChainToNumChains != nil {
			c = newNodeConfigs[i]
		}
		node := NewNode(t, c, configOpts...)
		nodesByPeerID[node.Keys.PeerID.String()] = *node
		// Note in real env, this ID is allocated by JD.
	}
	return nodesByPeerID
}

func createNewNodeConfigsWithChainTopology(t *testing.T, cfg NewNodesConfig, ports []int) []NewNodeConfig {
	homeChain, ok := cfg.Chains[cfg.HomeChainSel]
	require.Truef(t, ok, "home chain %d not found in chains, %+v", cfg.HomeChainSel, cfg.Chains)

	allEVMChains := maps.Values(cfg.Chains)
	allSolChains := maps.Values(cfg.SolChains)
	allAptosChains := maps.Values(cfg.AptosChains)

	// combine all the chains into a single slice.
	allChains := make([]any, 0, len(allEVMChains)+len(allSolChains)+len(allAptosChains))
	for _, chain := range allEVMChains {
		// Home chain is always EVM, so we can only check here.
		// We don't include it in allChains because it must be supported by all nodes,
		// so its a special case.
		if chain.ChainSelector() == cfg.HomeChainSel {
			continue
		}
		allChains = append(allChains, chain)
	}
	for _, chain := range allSolChains {
		allChains = append(allChains, chain)
	}
	for _, chain := range allAptosChains {
		allChains = append(allChains, chain)
	}

	// Validate that the chain topology is valid.
	// This should've been done already but we do it again for safety.
	totalChains := 0
	for _, numChains := range cfg.ChainTopology.FChainToNumChains {
		totalChains += numChains
	}
	require.Equalf(
		t,
		totalChains,
		len(allChains),
		"chain topology is invalid, totalChains: %d, expected: %d (total num chains minus home chain, which must be supported by all nodes)",
		totalChains,
		len(allChains),
	)

	var chainIdxToFChain = make(map[int]int) // index into allChains -> fChain value
	var chainIdx int
	for fChain, numChains := range cfg.ChainTopology.FChainToNumChains {
		for range numChains {
			chainIdxToFChain[chainIdx] = fChain
			chainIdx++
		}
	}

	// pseudocode:
	// for each chain:
	//   get the fChain value
	//   use the fChain value to "draw" the nodes that will support this chain.
	//   assign the chains to the nodes.
	nodeIdxToChainIdxs := make(map[int][]int)
	for chainIdx, fChain := range chainIdxToFChain {
		// "draw" the nodes that will support this chain.
		nodeIdxs := drawNodesForChain(t, fChain, cfg.NumNodes, cfg.ChainTopology.Seed)
		// assign the chains to the nodes.
		for _, nodeIdx := range nodeIdxs {
			nodeIdxToChainIdxs[nodeIdx] = append(nodeIdxToChainIdxs[nodeIdx], chainIdx)
		}
	}

	t.Logf("nodeIdxToChainIdxs: %v", nodeIdxToChainIdxs)

	// create the new node configs.
	var cfgs []NewNodeConfig
	for i := range cfg.NumNodes {
		nodeCfg := NewNodeConfig{
			Port: ports[cfg.NumBootstraps+i],
			// Every node must support the home chain, so include it in the starting config.
			Chains:         map[uint64]cldf.Chain{homeChain.ChainSelector(): homeChain},
			LogLevel:       cfg.LogLevel,
			Bootstrap:      false,
			RegistryConfig: cfg.RegistryConfig,
			CustomDBSetup:  cfg.CustomDBSetup,
		}
		chainsSupported, ok := nodeIdxToChainIdxs[i]
		require.Truef(t, ok, "node %d is not assigned to any chains", i)
		for _, chainIdx := range chainsSupported {
			chain := allChains[chainIdx]
			switch theChain := chain.(type) {
			case cldf.Chain:
				nodeCfg.Chains[theChain.ChainSelector()] = theChain
			case cldf.SolChain:
				nodeCfg.Solchains[theChain.ChainSelector()] = theChain
			case cldf_aptos.Chain:
				nodeCfg.Aptoschains[theChain.ChainSelector()] = theChain
			default:
				require.Failf(
					t,
					"unsupported chain type",
					"unsupported chain type: %T, forgot to add it to the switch statement?",
					theChain,
				)
			}
		}
		cfgs = append(cfgs, nodeCfg)
	}

	return cfgs
}

// drawNodesForChain draws a set of nodes that will support a given fChain value.
// Due to the randomness involved in the setup, a node might end up supporting multiple chains.
// However, this setup ensures that at most 3 * fChain + 1 nodes will support a given chain.
// TODO: we might want to make this more sophisticated in the future, e.g a node has a max # of chains it can support.
func drawNodesForChain(t *testing.T, fChain int, numNodes int, seed int64) (nodeIdxs []int) {
	require.Greaterf(t, fChain, 0, "fChain must be greater than 0, got %d", fChain)
	require.GreaterOrEqualf(t, numNodes, 3*fChain+1, "numNodes must be at least 3*fChain+1, got %d", numNodes)

	// Create a generator with a seed for reproducible setups.
	gen := rand.New(rand.NewSource(seed))

	numNodesToDraw := 3*fChain + 1
	nodeIdxs = make([]int, 0, numNodesToDraw)
	alreadyDrawn := make(map[int]bool)

	for len(nodeIdxs) < numNodesToDraw {
		drawnNodeIndex := gen.Intn(numNodes)

		if !alreadyDrawn[drawnNodeIndex] {
			nodeIdxs = append(nodeIdxs, drawnNodeIndex)
			alreadyDrawn[drawnNodeIndex] = true
		}
	}

	return nodeIdxs
}

func NewMemoryEnvironmentFromChainsNodes(
	ctx func() context.Context,
	lggr logger.Logger,
	chains map[uint64]cldf.Chain,
	solChains map[uint64]cldf.SolChain,
	aptosChains map[uint64]cldf_aptos.Chain,
	nodes map[string]Node,
) cldf.Environment {
	var nodeIDs []string
	for id := range nodes {
		nodeIDs = append(nodeIDs, id)
	}

	blockChains := map[uint64]chain.BlockChain{}
	for _, c := range chains {
		blockChains[c.Selector] = c
	}
	for _, c := range solChains {
		blockChains[c.Selector] = c
	}
	for _, c := range aptosChains {
		blockChains[c.Selector] = c
	}

	return *cldf.NewCLDFEnvironment(
		Memory,
		lggr,
		cldf.NewMemoryAddressBook(),
		datastore.NewMemoryDataStore[
			datastore.DefaultMetadata,
			datastore.DefaultMetadata,
		]().Seal(),
		chains,
		solChains,
		aptosChains,
		nodeIDs, // Note these have the p2p_ prefix.
		NewMemoryJobClient(nodes),
		ctx,
		cldf.XXXGenerateTestOCRSecrets(),
		chain.NewBlockChains(blockChains),
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
	c := NewNodesConfig{
		LogLevel:       logLevel,
		Chains:         chains,
		SolChains:      solChains,
		AptosChains:    aptosChains,
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

	blockChains := map[uint64]chain.BlockChain{}
	for _, c := range chains {
		blockChains[c.Selector] = c
	}
	for _, c := range solChains {
		blockChains[c.Selector] = c
	}
	for _, c := range aptosChains {
		blockChains[c.Selector] = c
	}
	return *cldf.NewCLDFEnvironment(
		Memory,
		lggr,
		cldf.NewMemoryAddressBook(),
		datastore.NewMemoryDataStore[
			datastore.DefaultMetadata,
			datastore.DefaultMetadata,
		]().Seal(),
		chains,
		solChains,
		nil, // this field will be deleted in future since env.BlockChains will now contain all the chains.
		nodeIDs,
		NewMemoryJobClient(nodes),
		t.Context,
		cldf.XXXGenerateTestOCRSecrets(),
		chain.NewBlockChains(blockChains),
	)
}
