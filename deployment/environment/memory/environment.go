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
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/freeport"

	chainsel "github.com/smartcontractkit/chain-selectors"

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
	NumOfUsersPerChain int
	Nodes              int
	Bootstraps         int
	RegistryConfig     deployment.CapabilityRegistryConfig
	CustomDBSetup      []string // SQL queries to run after DB creation
}

// For placeholders like aptos
func NewMemoryChain(t *testing.T, selector uint64) deployment.Chain {
	return deployment.Chain{
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
func NewMemoryChains(t *testing.T, numChains int, numUsers int) (map[uint64]deployment.Chain, map[uint64][]*bind.TransactOpts) {
	mchains := GenerateChains(t, numChains, numUsers)
	users := make(map[uint64][]*bind.TransactOpts)
	for id, chain := range mchains {
		sel, err := chainsel.SelectorFromChainId(id)
		require.NoError(t, err)
		users[sel] = chain.Users
	}
	return generateMemoryChain(t, mchains), users
}

func NewMemoryChainsSol(t *testing.T, numChains int) map[uint64]deployment.SolChain {
	mchains := GenerateChainsSol(t, numChains)
	return generateMemoryChainSol(mchains)
}

func NewMemoryChainsAptos(t *testing.T, numChains int) map[uint64]deployment.AptosChain {
	return GenerateChainsAptos(t, numChains)
}

func NewMemoryChainsWithChainIDs(t *testing.T, chainIDs []uint64, numUsers int) (map[uint64]deployment.Chain, map[uint64][]*bind.TransactOpts) {
	mchains := GenerateChainsWithIds(t, chainIDs, numUsers)
	users := make(map[uint64][]*bind.TransactOpts)
	for id, chain := range mchains {
		sel, err := chainsel.SelectorFromChainId(id)
		require.NoError(t, err)
		users[sel] = chain.Users
	}
	return generateMemoryChain(t, mchains), users
}

func generateMemoryChain(t *testing.T, inputs map[uint64]EVMChain) map[uint64]deployment.Chain {
	chains := make(map[uint64]deployment.Chain)
	for cid, chain := range inputs {
		chain := chain
		chainInfo, err := chainsel.GetChainDetailsByChainIDAndFamily(strconv.FormatUint(cid, 10), chainsel.FamilyEVM)
		require.NoError(t, err)
		backend := NewBackend(chain.Backend)
		chains[chainInfo.ChainSelector] = deployment.Chain{
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

func generateMemoryChainSol(inputs map[uint64]SolanaChain) map[uint64]deployment.SolChain {
	chains := make(map[uint64]deployment.SolChain)
	for cid, chain := range inputs {
		chain := chain
		chains[cid] = deployment.SolChain{
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
	logLevel zapcore.Level,
	chains map[uint64]deployment.Chain,
	solChains map[uint64]deployment.SolChain,
	aptosChains map[uint64]deployment.AptosChain,
	numNodes,
	numBootstraps int,
	registryConfig deployment.CapabilityRegistryConfig,
	customDBSetup []string, // SQL queries to run after DB creation
	configOpts ...ConfigOpt,
) map[string]Node {
	nodesByPeerID := make(map[string]Node)
	if numNodes+numBootstraps == 0 {
		return nodesByPeerID
	}
	ports := freeport.GetN(t, numBootstraps+numNodes)
	// bootstrap nodes must be separate nodes from plugin nodes,
	// since we won't run a bootstrapper and a plugin oracle on the same
	// chainlink node in production.
	for i := 0; i < numBootstraps; i++ {
		node := NewNode(t, ports[i], chains, solChains, aptosChains, logLevel, true /* bootstrap */, registryConfig, customDBSetup, configOpts...)
		nodesByPeerID[node.Keys.PeerID.String()] = *node
		// Note in real env, this ID is allocated by JD.
	}
	for i := 0; i < numNodes; i++ {
		// grab port offset by numBootstraps, since above loop also takes some ports.
		node := NewNode(t, ports[numBootstraps+i], chains, solChains, aptosChains, logLevel, false /* bootstrap */, registryConfig, customDBSetup, configOpts...)
		nodesByPeerID[node.Keys.PeerID.String()] = *node
		// Note in real env, this ID is allocated by JD.
	}
	return nodesByPeerID
}

func NewMemoryEnvironmentFromChainsNodes(
	ctx func() context.Context,
	lggr logger.Logger,
	chains map[uint64]deployment.Chain,
	solChains map[uint64]deployment.SolChain,
	aptosChains map[uint64]deployment.AptosChain,
	nodes map[string]Node,
) deployment.Environment {
	var nodeIDs []string
	for id := range nodes {
		nodeIDs = append(nodeIDs, id)
	}
	return *deployment.NewEnvironment(
		Memory,
		lggr,
		deployment.NewMemoryAddressBook(),
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
	)
}

// To be used by tests and any kind of deployment logic.
func NewMemoryEnvironment(t *testing.T, lggr logger.Logger, logLevel zapcore.Level, config MemoryEnvironmentConfig) deployment.Environment {
	chains, _ := NewMemoryChains(t, config.Chains, config.NumOfUsersPerChain)
	solChains := NewMemoryChainsSol(t, config.SolChains)
	aptosChains := NewMemoryChainsAptos(t, config.AptosChains)
	nodes := NewNodes(t, logLevel, chains, solChains, aptosChains, config.Nodes, config.Bootstraps, config.RegistryConfig, config.CustomDBSetup)
	var nodeIDs []string
	for id, node := range nodes {
		require.NoError(t, node.App.Start(t.Context()))
		t.Cleanup(func() {
			require.NoError(t, node.App.Stop())
		})
		nodeIDs = append(nodeIDs, id)
	}
	return *deployment.NewEnvironment(
		Memory,
		lggr,
		deployment.NewMemoryAddressBook(),
		datastore.NewMemoryDataStore[
			datastore.DefaultMetadata,
			datastore.DefaultMetadata,
		]().Seal(),
		chains,
		solChains,
		aptosChains,
		nodeIDs,
		NewMemoryJobClient(nodes),
		t.Context,
		cldf.XXXGenerateTestOCRSecrets(),
	)
}
