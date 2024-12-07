package changeset

import (
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"

	"github.com/smartcontractkit/chainlink/deployment"
)

type TestConfigs struct {
	Chains                   int
	NumOfUsersPerChain       int
	Nodes                    int
	Bootstraps               int
	IsUSDC                   bool
	IsUSDCAttestationMissing bool
	IsMultiCall3             bool
	OCRConfigOverride        func(CCIPOCRParams) CCIPOCRParams
	RMNEnabled               bool
}

func (t *TestConfigs) Validate() error {
	if t.Chains < 2 {
		return fmt.Errorf("chains must be at least 2")
	}
	if t.Nodes < 4 {
		return fmt.Errorf("nodes must be at least 4")
	}
	if t.Bootstraps < 1 {
		return fmt.Errorf("bootstraps must be at least 1")
	}
	return nil
}

func defaultTestConfigs() *TestConfigs {
	return &TestConfigs{
		Chains:             2,
		NumOfUsersPerChain: 1,
		Nodes:              4,
		Bootstraps:         1,
	}
}

type TestOps func(testCfg *TestConfigs)

func WithMultiCall3() TestOps {
	return func(testCfg *TestConfigs) {
		testCfg.IsMultiCall3 = true
	}
}

func WithRMNEnabled() TestOps {
	return func(testCfg *TestConfigs) {
		testCfg.RMNEnabled = true
	}
}

func WithOCRConfigOverride(override func(CCIPOCRParams) CCIPOCRParams) TestOps {
	return func(testCfg *TestConfigs) {
		testCfg.OCRConfigOverride = override
	}
}

func WithUSDCAttestationMissing() TestOps {
	return func(testCfg *TestConfigs) {
		testCfg.IsUSDCAttestationMissing = true
	}
}

func WithUSDC() TestOps {
	return func(testCfg *TestConfigs) {
		testCfg.IsUSDC = true
	}
}

func WithChains(numChains int) TestOps {
	return func(testCfg *TestConfigs) {
		testCfg.Chains = numChains
	}
}

func WithUsersPerChain(numUsers int) TestOps {
	return func(testCfg *TestConfigs) {
		testCfg.NumOfUsersPerChain = numUsers
	}
}

func WithNodes(numNodes int) TestOps {
	return func(testCfg *TestConfigs) {
		testCfg.Nodes = numNodes
	}
}

func WithBootstraps(numBootstraps int) TestOps {
	return func(testCfg *TestConfigs) {
		testCfg.Bootstraps = numBootstraps
	}
}

type TestEnvironment interface {
	SetupJobs(t *testing.T)
	HomeChainSel() uint64
	FeedChainSel() uint64
	ReplayLogs(t *testing.T, oc deployment.OffchainClient, replayBlocks map[uint64]uint64)
	Users() map[uint64][]*bind.TransactOpts
	New(t *testing.T, tc *TestConfigs) deployment.Environment
	StartEnvironmentWithJobsAndContracts(t *testing.T, tc *TestConfigs) deployment.Environment
}
