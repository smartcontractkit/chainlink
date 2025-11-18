package memory

import (
	"path/filepath"
	"runtime"
	"sync"
	"testing"
	"time"

	"github.com/Masterminds/semver/v3"
	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	cldf_solana_provider "github.com/smartcontractkit/chainlink-deployments-framework/chain/solana/provider"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"

	"github.com/smartcontractkit/chainlink/deployment/utils/solutils"
)

var (
	// Instead of a relative path, use runtime.Caller or go-bindata
	ProgramsPath = getProgramsPath()

	once = &sync.Once{}
)

func getProgramsPath() string {
	// Get the directory of the current file (environment.go)
	_, currentFile, _, _ := runtime.Caller(0)
	// Go up to the root of the deployment package
	rootDir := filepath.Dir(filepath.Dir(filepath.Dir(currentFile)))
	// Construct the absolute path
	return filepath.Join(rootDir, "ccip/changeset/internal", "solana_contracts")
}

func getTestSolanaChainSelectors() []uint64 {
	result := []uint64{}
	for _, x := range chainsel.SolanaALL {
		if x.Name == x.ChainID {
			result = append(result, x.Selector)
		}
	}
	return result
}

func generateChainsSol(t *testing.T, numChains int, commitSha string) []cldf_chain.BlockChain {
	t.Helper()

	if numChains == 0 {
		// Avoid downloading Solana program artifacts
		return nil
	}

	once.Do(func() {
		// TODO PLEX-1718 use latest contracts sha for now. Derive commit sha from go.mod once contracts are in a separate go module
		err := solutils.DownloadChainlinkSolanaProgramArtifacts(t.Context(), ProgramsPath, "b0f7cd3fbdbb", logger.Test(t))
		require.NoError(t, err)
		err = solutils.DownloadChainlinkCCIPProgramArtifacts(t.Context(), ProgramsPath, commitSha, logger.Test(t))
		require.NoError(t, err)
	})

	testSolanaChainSelectors := getTestSolanaChainSelectors()
	if len(testSolanaChainSelectors) < numChains {
		t.Fatalf("not enough test solana chain selectors available")
	}

	chains := make([]cldf_chain.BlockChain, 0, numChains)
	for i := range numChains {
		selector := testSolanaChainSelectors[i]

		c, err := cldf_solana_provider.NewCTFChainProvider(t, selector,
			cldf_solana_provider.CTFChainProviderConfig{
				Once:                         once,
				DeployerKeyGen:               cldf_solana_provider.PrivateKeyRandom(),
				ProgramsPath:                 ProgramsPath,
				ProgramIDs:                   SolanaProgramIDs,
				WaitDelayAfterContainerStart: 15 * time.Second, // we have slot errors that force retries if the chain is not given enough time to boot
			},
		).Initialize(t.Context())
		require.NoError(t, err)

		chains = append(chains, c)
	}

	return chains
}

// chainlink-ccip has dynamic resolution which does not work across repos
var SolanaProgramIDs = map[string]string{
	"ccip_router":               "Ccip842gzYHhvdDkSyi2YVCoAWPbYJoApMFzSxQroE9C",
	"test_token_pool":           "JuCcZ4smxAYv9QHJ36jshA7pA3FuQ3vQeWLUeAtZduJ",
	"burnmint_token_pool":       "41FGToCmdaWa1dgZLKFAjvmx6e6AjVTX7SVRibvsMGVB",
	"lockrelease_token_pool":    "8eqh8wppT9c5rw4ERqNCffvU6cNFJWff9WmkcYtmGiqC",
	"fee_quoter":                "FeeQPGkKDeRV1MgoYfMH6L8o3KeuYjwUZrgn4LRKfjHi",
	"test_ccip_receiver":        "EvhgrPhTDt4LcSPS2kfJgH6T6XWZ6wT3X9ncDGLT1vui",
	"ccip_offramp":              "offqSMQWgQud6WJz694LRzkeN5kMYpCHTpXQr3Rkcjm",
	"mcm":                       "5vNJx78mz7KVMjhuipyr9jKBKcMrKYGdjGkgE4LUmjKk",
	"timelock":                  "DoajfR5tK24xVw51fWcawUZWhAXD8yrBJVacc13neVQA",
	"access_controller":         "6KsN58MTnRQ8FfPaXHiFPPFGDRioikj9CdPvPxZJdCjb",
	"external_program_cpi_stub": "2zZwzyptLqwFJFEFxjPvrdhiGpH9pJ3MfrrmZX6NTKxm",
	"rmn_remote":                "RmnXLft1mSEwDgMKu2okYuHkiazxntFFcZFrrcXxYg7",
	"cctp_token_pool":           "CCiTPESGEevd7TBU8EGBKrcxuRq7jx3YtW6tPidnscaZ",
	"keystone_forwarder":        "whV7Q5pi17hPPyaPksToDw1nMx6Lh8qmNWKFaLRQ4wz",
	"data_feeds_cache":          "3kX63udXtYcsdj2737Wi2KGd2PhqiKPgAFAxstrjtRUa",
}

// Not deployed as part of the other solana programs, as it has its unique
// repository.
var SolanaNonCcipProgramIDs = map[string]string{
	"ccip_signer_registry": "S1GN4jus9XzKVVnoHqfkjo1GN8bX46gjXZQwsdGBPHE",
}

// Populates datastore with the predeployed program addresses
// pass map [programName]:ContractType of contracts to populate datastore with
func PopulateDatastore(ds *datastore.MemoryAddressRefStore, contracts map[string]datastore.ContractType, version *semver.Version, qualifier string, chainSel uint64) error {
	for programName, programID := range SolanaProgramIDs {
		ct, ok := contracts[programName]
		if !ok {
			continue
		}

		err := ds.Add(datastore.AddressRef{
			Address:       programID,
			ChainSelector: chainSel,
			Qualifier:     qualifier,
			Type:          ct,
			Version:       version,
		})

		if err != nil {
			return err
		}
	}

	return nil
}
