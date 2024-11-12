package memory

import (
	"context"
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/hashicorp/consul/sdk/freeport"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	chainsel "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink/deployment"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
)

const (
	Memory = "memory"
)

type MemoryEnvironmentConfig struct {
	Chains         int
	Nodes          int
	Bootstraps     int
	RegistryConfig deployment.CapabilityRegistryConfig
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
func NewMemoryChains(t *testing.T, numChains int) map[uint64]deployment.Chain {
	mchains := GenerateChains(t, numChains)
	return generateMemoryChain(t, mchains)
}

func NewMemoryChainsWithChainIDs(t *testing.T, chainIDs []uint64) map[uint64]deployment.Chain {
	mchains := GenerateChainsWithIds(t, chainIDs)
	return generateMemoryChain(t, mchains)
}

func generateMemoryChain(t *testing.T, inputs map[uint64]EVMChain) map[uint64]deployment.Chain {
	chains := make(map[uint64]deployment.Chain)
	for cid, chain := range inputs {
		chain := chain
		sel, err := chainsel.SelectorFromChainId(cid)
		require.NoError(t, err)
		backend := NewBackend(chain.Backend)
		chains[sel] = deployment.Chain{
			Selector:    sel,
			Client:      backend,
			DeployerKey: chain.DeployerKey,
			Confirm: func(tx *types.Transaction) (uint64, error) {
				if tx == nil {
					return 0, fmt.Errorf("tx was nil, nothing to confirm")
				}
				for {
					backend.Commit()
					receipt, err := backend.TransactionReceipt(context.Background(), tx.Hash())
					if err != nil {
						t.Log("failed to get receipt", err)
						continue
					}
					if receipt.Status == 0 {
						errReason, err := deployment.GetErrorReasonFromTx(chain.Backend.Client(), chain.DeployerKey.From, tx, receipt)
						if err == nil && errReason != "" {
							return 0, fmt.Errorf("tx %s reverted,error reason: %s", tx.Hash().Hex(), errReason)
						}
						return 0, fmt.Errorf("tx %s reverted, could not decode error reason", tx.Hash().Hex())
					}
					return receipt.BlockNumber.Uint64(), nil
				}
			},
		}
	}
	return chains
}

func NewNodes(t *testing.T, logLevel zapcore.Level, chains map[uint64]deployment.Chain, numNodes, numBootstraps int, registryConfig deployment.CapabilityRegistryConfig) map[string]Node {
	nodesByPeerID := make(map[string]Node)
	ports := freeport.GetN(t, numBootstraps+numNodes)
	// bootstrap nodes must be separate nodes from plugin nodes,
	// since we won't run a bootstrapper and a plugin oracle on the same
	// chainlink node in production.
	for i := 0; i < numBootstraps; i++ {
		node := NewNode(t, ports[i], chains, logLevel, true /* bootstrap */, registryConfig)
		nodesByPeerID[node.Keys.PeerID.String()] = *node
		// Note in real env, this ID is allocated by JD.
	}
	for i := 0; i < numNodes; i++ {
		// grab port offset by numBootstraps, since above loop also takes some ports.
		node := NewNode(t, ports[numBootstraps+i], chains, logLevel, false /* bootstrap */, registryConfig)
		nodesByPeerID[node.Keys.PeerID.String()] = *node
		// Note in real env, this ID is allocated by JD.
	}
	return nodesByPeerID
}

func NewMemoryEnvironmentFromChainsNodes(t *testing.T,
	lggr logger.Logger,
	chains map[uint64]deployment.Chain,
	nodes map[string]Node) deployment.Environment {
	var nodeIDs []string
	for id := range nodes {
		nodeIDs = append(nodeIDs, id)
	}
	return *deployment.NewEnvironment(
		Memory,
		lggr,
		deployment.NewMemoryAddressBook(),
		chains,
		nodeIDs, // Note these have the p2p_ prefix.
		NewMemoryJobClient(nodes),
	)
}

// To be used by tests and any kind of deployment logic.
func NewMemoryEnvironment(t *testing.T, lggr logger.Logger, logLevel zapcore.Level, config MemoryEnvironmentConfig) deployment.Environment {
	chains := NewMemoryChains(t, config.Chains)
	nodes := NewNodes(t, logLevel, chains, config.Nodes, config.Bootstraps, config.RegistryConfig)
	var nodeIDs []string
	for id := range nodes {
		nodeIDs = append(nodeIDs, id)
	}
	return *deployment.NewEnvironment(
		Memory,
		lggr,
		deployment.NewMemoryAddressBook(),
		chains,
		nodeIDs,
		NewMemoryJobClient(nodes),
	)
}
