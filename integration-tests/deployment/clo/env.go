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

//type DonEnvironment deployment.Environment

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
				if chain.Ocr1Config.IsBootstrap {
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

func NewDonEnvWithMemoryChains(t *testing.T, cfg DonEnvConfig) *deployment.Environment {
	e := NewDonEnv(t, cfg)
	// overwrite the chains with memory chains
	chains := make(map[uint64]struct{})
	for _, nop := range cfg.Nops {
		for _, node := range nop.Nodes {
			for _, chain := range node.ChainConfigs {
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

func NewMultiDonEnvironment(lggr logger.Logger, dons map[string]*deployment.Environment) deployment.MultiDonEnvironment {
	out := deployment.MultiDonEnvironment{
		Logger:   lggr,
		DonToEnv: dons,
	}
	return out
}
