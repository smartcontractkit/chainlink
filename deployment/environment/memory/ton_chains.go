package memory

import (
	"fmt"
	"os"
	"strings"
	"testing"

	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/require"
	"github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/tlb"
	"golang.org/x/mod/modfile"

	"github.com/xssnick/tonutils-go/ton/wallet"

	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	cldf_ton "github.com/smartcontractkit/chainlink-deployments-framework/chain/ton"
	cldf_ton_provider "github.com/smartcontractkit/chainlink-deployments-framework/chain/ton/provider"
	"github.com/smartcontractkit/chainlink-ton/deployment/utils"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
)

func getTestTonChainSelectors() []uint64 {
	return []uint64{chainsel.TON_LOCALNET.Selector}
}

func GetTONSha() (version string, err error) {
	modFilePath, err := getModFilePath()
	if err != nil {
		return "", err
	}
	go_mod_version, err := getTONCcipDependencyVersion(modFilePath)
	if err != nil {
		return "", err
	}
	tokens := strings.Split(go_mod_version, "-")
	if len(tokens) == 3 {
		version := tokens[len(tokens)-1]
		return version, nil
	} else {
		return "", fmt.Errorf("invalid go.mod version: %s", go_mod_version)
	}
}

func getTONCcipDependencyVersion(gomodPath string) (string, error) {
	const dependency = "github.com/smartcontractkit/chainlink-ton"

	gomod, err := os.ReadFile(gomodPath)
	if err != nil {
		return "", err
	}

	modFile, err := modfile.ParseLax("go.mod", gomod, nil)
	if err != nil {
		return "", err
	}

	for _, dep := range modFile.Require {
		if dep.Mod.Path == dependency {
			return dep.Mod.Version, nil
		}
	}

	return "", fmt.Errorf("dependency %s not found", dependency)
}

func generateChainsTon(t *testing.T, numChains int) []cldf_chain.BlockChain {
	testTonChainSelectors := getTestTonChainSelectors()
	if numChains > 1 {
		t.Fatalf("only one ton chain is supported for now, got %d", numChains)
	}
	if len(testTonChainSelectors) < numChains {
		t.Fatalf("not enough test ton chain selectors available")
	}

	chains := make([]cldf_chain.BlockChain, 0, numChains)
	for i := range numChains {
		selector := testTonChainSelectors[i]

		c, err := cldf_ton_provider.NewCTFChainProvider(t, selector,
			cldf_ton_provider.CTFChainProviderConfig{
				Once: once,
			},
		).Initialize(t.Context())
		require.NoError(t, err)

		chains = append(chains, c)
		tonChain, ok := c.(cldf_ton.Chain)
		if !ok {
			t.Fatalf("expected cldf_ton.Chain, got %T", c)
		}

		// memory environment doesn't block on funding so changesets can execute before the env is fully ready, manually call fund so we block here
		utils.FundWallets(t, tonChain.Client, []*address.Address{tonChain.WalletAddress}, []tlb.Coins{tlb.MustFromTON("1000")})
	}

	return chains
}

func createTonChainConfig(chainID string, chain cldf_ton.Chain) chainlink.RawConfig {
	chainConfig := chainlink.RawConfig{}

	chainConfig["Enabled"] = true
	chainConfig["ChainID"] = chainID
	chainConfig["NetworkName"] = "ton-local"
	chainConfig["NetworkNameFull"] = "ton-local"
	chainConfig["Nodes"] = []any{
		map[string]any{
			"Name": "primary",
			"URL":  chain.URL,
		},
	}

	return chainConfig
}

func fundNodesTon(t *testing.T, tonChain cldf_ton.Chain, nodes []*Node) {
	messages := make([]*wallet.Message, 0, len(nodes))
	for _, node := range nodes {
		tonkeys, err := node.App.GetKeyStore().TON().GetAll()
		require.NoError(t, err)
		require.Len(t, tonkeys, 1)
		transmitter := tonkeys[0].PubkeyToAddress()
		msg, err := tonChain.Wallet.BuildTransfer(transmitter, tlb.MustFromTON("1000"), false, "")
		require.NoError(t, err)
		messages = append(messages, msg)
	}
	_, _, err := tonChain.Wallet.SendManyWaitTransaction(t.Context(), messages)
	require.NoError(t, err)
}
