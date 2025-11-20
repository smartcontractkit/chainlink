package memory

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/require"
	"github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/tlb"
	"golang.org/x/mod/modfile"

	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	cldf_ton "github.com/smartcontractkit/chainlink-deployments-framework/chain/ton"
	cldf_ton_provider "github.com/smartcontractkit/chainlink-deployments-framework/chain/ton/provider"
	"github.com/smartcontractkit/chainlink-ton/deployment/utils"
)

var (
	deployerFundAmount = tlb.MustFromTON("1000")
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
		utils.FundWallets(t, tonChain.Client, []*address.Address{tonChain.WalletAddress}, []tlb.Coins{deployerFundAmount})
	}

	return chains
}

func getModFilePath() (string, error) {
	_, currentFile, _, _ := runtime.Caller(0)
	// Get the root directory by walking up from current file until we find go.mod
	rootDir := filepath.Dir(currentFile)
	for {
		if _, err := os.Stat(filepath.Join(rootDir, "go.mod")); err == nil {
			break
		}
		parent := filepath.Dir(rootDir)
		if parent == rootDir {
			return "", errors.New("could not find project root directory containing go.mod")
		}
		rootDir = parent
	}
	return filepath.Join(rootDir, "go.mod"), nil
}
