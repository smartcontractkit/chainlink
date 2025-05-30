package memory

import (
	"testing"

	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/require"

	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	cldf_aptos "github.com/smartcontractkit/chainlink-deployments-framework/chain/aptos"
	cldf_aptos_provider "github.com/smartcontractkit/chainlink-deployments-framework/chain/aptos/provider"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
)

func getTestAptosChainSelectors() []uint64 {
	// TODO: CTF to support different chain ids, need to investigate if it's possible (thru node config.yaml?)
	return []uint64{chainsel.APTOS_LOCALNET.Selector}
}

func generateChainsAptos(t *testing.T, numChains int) []cldf_chain.BlockChain {
	t.Helper()

	testAptosChainSelectors := getTestAptosChainSelectors()
	if len(testAptosChainSelectors) < numChains {
		t.Fatalf("not enough test aptos chain selectors available")
	}

	chains := make([]cldf_chain.BlockChain, 0, numChains)
	for i := range numChains {
		selector := testAptosChainSelectors[i]

		c, err := cldf_aptos_provider.NewCTFChainProvider(t, selector,
			cldf_aptos_provider.CTFChainProviderConfig{
				Once:              once,
				DeployerSignerGen: cldf_aptos_provider.AccountGenCTFDefault(),
			},
		).Initialize()
		require.NoError(t, err)

		chains = append(chains, c)
	}

	t.Logf("Created %d Aptos chains", len(chains))

	return chains
}

func createAptosChainConfig(chainID string, chain cldf_aptos.Chain) chainlink.RawConfig {
	chainConfig := chainlink.RawConfig{}

	chainConfig["Enabled"] = true
	chainConfig["ChainID"] = chainID
	chainConfig["NetworkName"] = "aptos-local"
	chainConfig["NetworkNameFull"] = "aptos-local"
	chainConfig["Nodes"] = []any{
		map[string]any{
			"Name": "primary",
			"URL":  chain.URL,
		},
	}

	return chainConfig
}
