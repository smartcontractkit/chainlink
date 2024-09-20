package memory

import (
	"context"
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/hashicorp/consul/sdk/freeport"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	chainsel "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink/integration-tests/deployment"

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

// Needed for environment variables on the node which point to prexisitng addresses.
// i.e. CapReg.
func NewMemoryChains(t *testing.T, numChains int) map[uint64]deployment.Chain {
	mchains := GenerateChains(t, numChains)
	chains := make(map[uint64]deployment.Chain)
	for cid, chain := range mchains {
		sel, err := chainsel.SelectorFromChainId(cid)
		require.NoError(t, err)
		chains[sel] = deployment.Chain{
			Selector:    sel,
			Client:      chain.Backend,
			DeployerKey: chain.DeployerKey,
			Confirm: func(tx *types.Transaction) (uint64, error) {
				if tx == nil {
					return 0, fmt.Errorf("tx was nil, nothing to confirm")
				}
				for {
					chain.Backend.Commit()
					receipt, err := chain.Backend.TransactionReceipt(context.Background(), tx.Hash())
					if err != nil {
						t.Log("failed to get receipt", err)
						continue
					}
					if receipt.Status == 0 {
						t.Logf("Status (reverted) %d for txhash %s\n", receipt.Status, tx.Hash().Hex())
					}
					return receipt.BlockNumber.Uint64(), nil
				}
			},
		}
	}
	return chains
}

func NewNodes(t *testing.T, logLevel zapcore.Level, chains map[uint64]deployment.Chain, numNodes, numBootstraps int, registryConfig deployment.CapabilityRegistryConfig) map[string]Node {
	mchains := make(map[uint64]EVMChain)
	for _, chain := range chains {
		evmChainID, err := chainsel.ChainIdFromSelector(chain.Selector)
		if err != nil {
			t.Fatal(err)
		}
		mchains[evmChainID] = EVMChain{
			Backend:     chain.Client.(*backends.SimulatedBackend),
			DeployerKey: chain.DeployerKey,
		}
	}
	nodesByPeerID := make(map[string]Node)
	ports := freeport.GetN(t, numBootstraps+numNodes)
	// bootstrap nodes must be separate nodes from plugin nodes,
	// since we won't run a bootstrapper and a plugin oracle on the same
	// chainlink node in production.
	for i := 0; i < numBootstraps; i++ {
		node := NewNode(t, ports[i], mchains, logLevel, true /* bootstrap */, registryConfig)
		nodesByPeerID[node.Keys.PeerID.String()] = *node
		// Note in real env, this ID is allocated by JD.
	}
	for i := 0; i < numNodes; i++ {
		// grab port offset by numBootstraps, since above loop also takes some ports.
		node := NewNode(t, ports[numBootstraps+i], mchains, logLevel, false /* bootstrap */, registryConfig)
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
	return deployment.Environment{
		Name:     Memory,
		Offchain: NewMemoryJobClient(nodes),
		// Note these have the p2p_ prefix.
		NodeIDs: nodeIDs,
		Chains:  chains,
		Logger:  lggr,
	}
}

//func NewMemoryEnvironmentExistingChains(t *testing.T, lggr logger.Logger,
//	chains map[uint64]deployment.Chain, config MemoryEnvironmentConfig) deployment.Environment {
//	nodes := NewNodes(t, chains, config.Nodes, config.Bootstraps, config.CapabilityRegistryConfig)
//	var nodeIDs []string
//	for id := range nodes {
//		nodeIDs = append(nodeIDs, id)
//	}
//	return deployment.Environment{
//		Name:     Memory,
//		Offchain: NewMemoryJobClient(nodes),
//		// Note these have the p2p_ prefix.
//		NodeIDs: nodeIDs,
//		Chains:  chains,
//		Logger:  lggr,
//	}
//}

// To be used by tests and any kind of deployment logic.
func NewMemoryEnvironment(t *testing.T, lggr logger.Logger, logLevel zapcore.Level, config MemoryEnvironmentConfig) deployment.Environment {
	chains := NewMemoryChains(t, config.Chains)
	nodes := NewNodes(t, logLevel, chains, config.Nodes, config.Bootstraps, config.RegistryConfig)
	var nodeIDs []string
	for id := range nodes {
		nodeIDs = append(nodeIDs, id)
	}
	return deployment.Environment{
		Name:     Memory,
		Offchain: NewMemoryJobClient(nodes),
		NodeIDs:  nodeIDs,
		Chains:   chains,
		Logger:   lggr,
	}
}
