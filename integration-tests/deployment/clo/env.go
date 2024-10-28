package clo

import (
	"strconv"
	"testing"

	"github.com/test-go/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink/integration-tests/deployment"
	"github.com/smartcontractkit/chainlink/integration-tests/deployment/clo/models"
	"github.com/smartcontractkit/chainlink/integration-tests/deployment/memory"
)

type DonEnvConfig struct {
	DonName string
	Chains  map[uint64]deployment.Chain
	Logger  logger.Logger
	Nops    []*models.NodeOperator
}

func NewDonEnv(t *testing.T, cfg DonEnvConfig) *deployment.Environment {
	// no bootstraps in the don as far as capabilities registry is concerned
	for _, nop := range cfg.Nops {
		for _, node := range nop.Nodes {
			for _, chain := range node.ChainConfigs {
				if chain.Ocr2Config.IsBootstrap {
					t.Fatalf("Don nodes should not be bootstraps nop %s node %s chain %s", nop.ID, node.ID, chain.Network.ChainID)
				}
			}
		}
	}
	out := deployment.Environment{
		Name:     cfg.DonName,
		Offchain: NewJobClient(cfg.Logger, cfg.Nops),
		NodeIDs:  make([]string, 0),
		Chains:   cfg.Chains,
		Logger:   cfg.Logger,
	}
	// assume that all the nodes in the provided input nops are part of the don
	for _, nop := range cfg.Nops {
		for _, node := range nop.Nodes {
			out.NodeIDs = append(out.NodeIDs, node.ID)
		}
	}

	return &out
}

func NewDonEnvWithMemoryChains(t *testing.T, cfg DonEnvConfig, ignore func(*models.NodeChainConfig) bool) *deployment.Environment {
	e := NewDonEnv(t, cfg)
	// overwrite the chains with memory chains
	chains := make(map[uint64]struct{})
	for _, nop := range cfg.Nops {
		for _, node := range nop.Nodes {
			for _, chain := range node.ChainConfigs {
				if ignore(chain) {
					continue
				}
				id, err := strconv.ParseUint(chain.Network.ChainID, 10, 64)
				require.NoError(t, err, "failed to parse chain id to uint64")
				chains[id] = struct{}{}
			}
		}
	}
	var cs []uint64
	for c := range chains {
		cs = append(cs, c)
	}
	memoryChains := memory.NewMemoryChainsWithChainIDs(t, cs)
	e.Chains = memoryChains
	return e
}

// MultiDonEnvironment is a single logical deployment environment (like dev, testnet, prod,...).
// It represents the idea that different nodesets host different capabilities.
// Each element in the DonEnv is a logical set of nodes that host the same capabilities.
// This model allows us to reuse the existing Environment abstraction while supporting multiple nodesets at
// expense of slightly abusing the original abstraction. Specifically, the abuse is that
// each Environment in the DonToEnv map is a subset of the target deployment environment.
// One element cannot represent dev and other testnet for example.
type MultiDonEnvironment struct {
	donToEnv map[string]*deployment.Environment
	Logger   logger.Logger
	// hacky but temporary to transition to Environment abstraction. set by New
	Chains map[uint64]deployment.Chain
}

func (mde MultiDonEnvironment) Flatten(name string) *deployment.Environment {
	return &deployment.Environment{
		Name:   name,
		Chains: mde.Chains,
		Logger: mde.Logger,

		// TODO: KS-460 integrate with the clo offchain client impl
		// may need to extend the Environment abstraction use maps rather than slices for Nodes
		// somehow we need to capture the fact that each nodes belong to nodesets which have different capabilities
		// purposely nil to catch misuse until we do that work
		Offchain: nil,
		NodeIDs:  nil,
	}
}

func newMultiDonEnvironment(logger logger.Logger, donToEnv map[string]*deployment.Environment) *MultiDonEnvironment {
	chains := make(map[uint64]deployment.Chain)
	for _, env := range donToEnv {
		for sel, chain := range env.Chains {
			if _, exists := chains[sel]; !exists {
				chains[sel] = chain
			}
		}
	}
	return &MultiDonEnvironment{
		donToEnv: donToEnv,
		Logger:   logger,
		Chains:   chains,
	}
}

func NewTestEnv(t *testing.T, lggr logger.Logger, dons map[string]*deployment.Environment) *MultiDonEnvironment {
	for _, don := range dons {
		//don := don
		seen := make(map[uint64]deployment.Chain)
		// ensure that generated chains are the same for all environments. this ensures that he in memory representation
		// points to a common object for all dons given the same selector.
		for sel, chain := range don.Chains {
			c, exists := seen[sel]
			if exists {
				don.Chains[sel] = c
			} else {
				seen[sel] = chain
			}
		}
	}
	return newMultiDonEnvironment(lggr, dons)
}
