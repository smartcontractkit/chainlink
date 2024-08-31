package memory

import (
	"context"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/hashicorp/consul/sdk/freeport"
	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

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
	RegistryConfig RegistryConfig
}

// Needed for environment variables on the node which point to prexisitng addresses.
// i.e. CapReg.
func NewMemoryChains(t *testing.T, numChains int) (map[uint64]deployment.Chain, map[uint64]EVMChain) {
	mchains := GenerateChains(t, numChains)
	chains := make(map[uint64]deployment.Chain)
	for cid, chain := range mchains {
		sel, err := chainsel.SelectorFromChainId(cid)
		require.NoError(t, err)
		chains[sel] = deployment.Chain{
			Selector:    sel,
			Client:      chain.Backend.Client(),
			DeployerKey: chain.DeployerKey,
			Confirm: func(tx common.Hash) error {
				for {
					chain.Backend.Commit()
					receipt, err := chain.Backend.Client().TransactionReceipt(context.Background(), tx)
					if err != nil {
						t.Log("failed to get receipt", err)
						continue
					}
					if receipt.Status == 0 {
						t.Logf("Status (reverted) %d for txhash %s\n", receipt.Status, tx.String())
					}
					return nil
				}
			},
		}
	}
	return chains, mchains
}

func NewNodes(t *testing.T, logLevel zapcore.Level, mchains map[uint64]EVMChain, numNodes, numBootstraps int, registryConfig RegistryConfig) map[string]Node {
	nodesByPeerID := make(map[string]Node)
	ports := freeport.GetN(t, numNodes)
	var existingNumBootstraps int
	for i := 0; i < numNodes; i++ {
		bootstrap := false
		if existingNumBootstraps < numBootstraps {
			bootstrap = true
			existingNumBootstraps++
		}
		node := NewNode(t, ports[i], mchains, logLevel, bootstrap, registryConfig)
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
//	nodes := NewNodes(t, chains, config.Nodes, config.Bootstraps, config.RegistryConfig)
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
	chains, mchains := NewMemoryChains(t, config.Chains)
	nodes := NewNodes(t, logLevel, mchains, config.Nodes, config.Bootstraps, config.RegistryConfig)
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
