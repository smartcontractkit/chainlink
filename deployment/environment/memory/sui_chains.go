package memory

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	chainsel "github.com/smartcontractkit/chain-selectors"

	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	cldf_sui_provider "github.com/smartcontractkit/chainlink-deployments-framework/chain/sui/provider"
)

func getTestSuiChainSelectors() []uint64 {
	// TODO: CTF to support different chain ids, need to investigate if it's possible (thru node config.yaml?)
	return []uint64{chainsel.SUI_LOCALNET.Selector}
}

func randomSeed() []byte {
	seed := make([]byte, ed25519.SeedSize)
	_, err := rand.Read(seed)
	if err != nil {
		panic(fmt.Sprintf("failed to generate random seed: %+v", err))
	}

	return seed
}

func GenerateChainsSui(t *testing.T, numChains int) []cldf_chain.BlockChain {
	testSuiChainSelectors := getTestSuiChainSelectors()
	if len(testSuiChainSelectors) < numChains {
		t.Fatalf("not enough test sui chain selectors available")
	}
	chains := make([]cldf_chain.BlockChain, 0, numChains)
	for i := range numChains {
		selector := testSuiChainSelectors[i]

		seeded := ed25519.NewKeyFromSeed(randomSeed()) // 64 bytes: seed||pub
		seed := seeded[:32]                            // or: seeded.Seed() if available
		hexKey := hex.EncodeToString(seed)             // 64 hex chars

		platform := "linux/amd64"
		img := "mysten/sui-tools:devnet"
		// generate adhoc sui privKey
		c, err := cldf_sui_provider.NewCTFChainProvider(t, selector,
			cldf_sui_provider.CTFChainProviderConfig{
				Once:              once,
				DeployerSignerGen: cldf_sui_provider.AccountGenPrivateKey(hexKey),
				Image:             &img,
				Platform:          &platform,
			},
		).Initialize(t.Context())
		require.NoError(t, err)

		chains = append(chains, c)
	}

	t.Logf("Created %d Sui chains: %+v", len(chains), chains)
	return chains
}
