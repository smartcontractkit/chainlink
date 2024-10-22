package ccipdeployment

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink/integration-tests/deployment"
	"github.com/smartcontractkit/chainlink/integration-tests/deployment/memory"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestDeployCCIPContracts(t *testing.T) {
	lggr := logger.TestLogger(t)
	e := memory.NewMemoryEnvironment(t, lggr, zapcore.InfoLevel, memory.MemoryEnvironmentConfig{
		Bootstraps: 1,
		Chains:     2,
		Nodes:      4,
	})
	// Deploy all the CCIP contracts.
	ab := deployment.NewMemoryAddressBook()
	homeChainSel, feedChainSel := allocateCCIPChainSelectors(e.Chains)
	feeTokenContracts, _ := DeployTestContracts(t, lggr, ab, homeChainSel, feedChainSel, e.Chains)

	// Load the state after deploying the cap reg and feeds.
	s, err := LoadOnchainState(e, ab)
	require.NoError(t, err)
	require.NotNil(t, s.Chains[homeChainSel].CapabilityRegistry)
	require.NotNil(t, s.Chains[homeChainSel].CCIPHome)
	require.NotNil(t, s.Chains[feedChainSel].USDFeeds)

	err = DeployCCIPContracts(e, ab, DeployCCIPContractConfig{
		HomeChainSel:       homeChainSel,
		FeedChainSel:       feedChainSel,
		ChainsToDeploy:     e.AllChainSelectors(),
		TokenConfig:        NewTokenConfig(),
		CapabilityRegistry: s.Chains[homeChainSel].CapabilityRegistry.Address(),
		FeeTokenContracts:  feeTokenContracts,
		MCMSConfig:         NewTestMCMSConfig(t, e),
		OCRSecrets:         deployment.XXXGenerateTestOCRSecrets(),
	})
	require.NoError(t, err)
	state, err := LoadOnchainState(e, ab)
	require.NoError(t, err)
	snap, err := state.View(e.AllChainSelectors())
	require.NoError(t, err)

	// Assert expect every deployed address to be in the address book.
	// TODO (CCIP-3047): Add the rest of CCIPv2 representation
	b, err := json.MarshalIndent(snap, "", "	")
	require.NoError(t, err)
	fmt.Println(string(b))
}

func TestJobSpecGeneration(t *testing.T) {
	lggr := logger.TestLogger(t)
	e := memory.NewMemoryEnvironment(t, lggr, zapcore.InfoLevel, memory.MemoryEnvironmentConfig{
		Chains: 1,
		Nodes:  1,
	})
	js, err := NewCCIPJobSpecs(e.NodeIDs, e.Offchain)
	require.NoError(t, err)
	for node, jb := range js {
		fmt.Println(node, jb)
	}
	// TODO: Add job assertions
}
