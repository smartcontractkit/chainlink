package ccipdeployment

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
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
	homeChainSel, feedChainSel := allocateCCIPChainSelectors(e.Chains)
	_ = DeployTestContracts(t, lggr, e.ExistingAddresses, homeChainSel, feedChainSel, e.Chains, MockLinkPrice, MockWethPrice)

	nodes, err := deployment.NodeInfo(e.NodeIDs, e.Offchain)
	require.NoError(t, err)

	_, err = DeployHomeChain(lggr, e, e.ExistingAddresses, e.Chains[homeChainSel],
		NewTestRMNStaticConfig(),
		NewTestRMNDynamicConfig(),
		NewTestNodeOperator(e.Chains[homeChainSel].DeployerKey.From),
		map[string][][32]byte{
			"NodeOperator": nodes.NonBootstraps().PeerIDs(),
		},
	)
	require.NoError(t, err)
	// Load the state after deploying the cap reg and feeds.
	s, err := LoadOnchainState(e)
	require.NoError(t, err)
	require.NotNil(t, s.Chains[homeChainSel].CapabilityRegistry)
	require.NotNil(t, s.Chains[homeChainSel].CCIPHome)
	require.NotNil(t, s.Chains[feedChainSel].USDFeeds)

	newAddresses := deployment.NewMemoryAddressBook()
	err = DeployPrerequisiteChainContracts(e, newAddresses, e.AllChainSelectors())
	require.NoError(t, err)
	require.NoError(t, e.ExistingAddresses.Merge(newAddresses))

	newAddresses = deployment.NewMemoryAddressBook()
	err = DeployCCIPContracts(e, newAddresses, DeployCCIPContractConfig{
		HomeChainSel:   homeChainSel,
		FeedChainSel:   feedChainSel,
		ChainsToDeploy: e.AllChainSelectors(),
		TokenConfig:    NewTokenConfig(),
		MCMSConfig:     NewTestMCMSConfig(t, e),
		OCRSecrets:     deployment.XXXGenerateTestOCRSecrets(),
	})
	require.NoError(t, err)
	require.NoError(t, e.ExistingAddresses.Merge(newAddresses))
	state, err := LoadOnchainState(e)
	require.NoError(t, err)
	snap, err := state.View(e.AllChainSelectors())
	require.NoError(t, err)

	// Assert expect every deployed address to be in the address book.
	// TODO (CCIP-3047): Add the rest of CCIPv2 representation
	b, err := json.MarshalIndent(snap, "", "	")
	require.NoError(t, err)
	fmt.Println(string(b))
}
