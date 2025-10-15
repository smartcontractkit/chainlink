package memory

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/Masterminds/semver/v3"
	"github.com/gagliardetto/solana-go"
	solRpc "github.com/gagliardetto/solana-go/rpc"
	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/require"
	"golang.org/x/mod/modfile"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	cldf_solana "github.com/smartcontractkit/chainlink-deployments-framework/chain/solana"
	cldf_solana_provider "github.com/smartcontractkit/chainlink-deployments-framework/chain/solana/provider"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
)

var (
	// Instead of a relative path, use runtime.Caller or go-bindata
	ProgramsPath = getProgramsPath()
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

func FundSolanaAccounts(
	ctx context.Context, accounts []solana.PublicKey, solAmount uint64, solanaGoClient *solRpc.Client,
) error {
	var sigs = make([]solana.Signature, 0, len(accounts))
	for _, account := range accounts {
		sig, err := solanaGoClient.RequestAirdrop(ctx, account, solAmount*solana.LAMPORTS_PER_SOL,
			solRpc.CommitmentFinalized)
		if err != nil {
			return err
		}
		sigs = append(sigs, sig)
	}

	const timeout = 100 * time.Second
	const pollInterval = 500 * time.Millisecond

	timeoutCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()

	remaining := len(sigs)
	for remaining > 0 {
		select {
		case <-timeoutCtx.Done():
			return errors.New("unable to find transaction within timeout")
		case <-ticker.C:
			statusRes, sigErr := solanaGoClient.GetSignatureStatuses(ctx, true, sigs...)
			if sigErr != nil {
				return sigErr
			}
			if statusRes == nil {
				return errors.New("Status response is nil")
			}
			if statusRes.Value == nil {
				return errors.New("Status response value is nil")
			}

			unfinalizedCount := 0
			for _, res := range statusRes.Value {
				if res == nil || res.ConfirmationStatus == solRpc.ConfirmationStatusFinalized {
					unfinalizedCount++
				}
			}
			remaining = unfinalizedCount
		}
	}
	return nil
}

// FundSolanaAccountsWithLogging requests airdrops for the provided accounts and waits for confirmation.
// It waits until all transactions reach at least "Confirmed" commitment level with enhanced logging and timeouts.
// Solana commitment levels: Processed < Confirmed < Finalized
// - Processed: Transaction processed by a validator but may be rolled back
// - Confirmed: Transaction confirmed by supermajority of cluster stake
// - Finalized: Transaction finalized and cannot be rolled back
func FundSolanaAccountsWithLogging(
	ctx context.Context, accounts []solana.PublicKey, solAmount uint64, solanaGoClient *solRpc.Client,
	lggr logger.Logger,
) error {
	if len(accounts) == 0 {
		return nil
	}

	var sigs = make([]solana.Signature, 0, len(accounts))
	var successfulAccounts = make([]solana.PublicKey, 0, len(accounts))

	lggr.Infow("Starting Solana airdrop requests", "accountCount", len(accounts), "amountSOL", solAmount)

	// Request airdrops with better error tracking
	// Note: Using CommitmentConfirmed here means the RequestAirdrop call itself waits for confirmed status
	for i, account := range accounts {
		sig, err := solanaGoClient.RequestAirdrop(ctx, account, solAmount*solana.LAMPORTS_PER_SOL, solRpc.CommitmentFinalized)
		if err != nil {
			// Return partial success information
			if len(sigs) > 0 {
				return fmt.Errorf("airdrop request failed for account %d (%s): %w (note: %d previous requests may have succeeded)",
					i, account.String(), err, len(sigs))
			}
			return fmt.Errorf("airdrop request failed for account %d (%s): %w", i, account.String(), err)
		}
		sigs = append(sigs, sig)
		successfulAccounts = append(successfulAccounts, account)

		lggr.Debugw("Airdrop request completed",
			"progress", fmt.Sprintf("%d/%d", i+1, len(accounts)),
			"account", account.String(),
			"signature", sig.String())

		// small delay to avoid rate limiting issues
		time.Sleep(100 * time.Millisecond)
	}

	// Adaptive timeout based on batch size - each airdrop can take several seconds
	// Base timeout of 30s + 5s per account for larger batches
	baseTimeout := 60 * time.Second
	if len(accounts) > 5 {
		baseTimeout += time.Duration(len(accounts)) * 5 * time.Second
	}
	timeout := baseTimeout
	const pollInterval = 500 * time.Millisecond

	lggr.Infow("Starting confirmation polling", "timeout", timeout, "accounts", len(accounts))

	timeoutCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()

	remaining := len(sigs)
	pollCount := 0
	for remaining > 0 {
		select {
		case <-timeoutCtx.Done():
			// Log which transactions are still unconfirmed for debugging
			unfinalizedSigs := []string{}
			statusRes, _ := solanaGoClient.GetSignatureStatuses(ctx, true, sigs...)
			if statusRes != nil && statusRes.Value != nil {
				for i, res := range statusRes.Value {
					if res == nil || res.ConfirmationStatus != solRpc.ConfirmationStatusFinalized {
						unfinalizedSigs = append(unfinalizedSigs, fmt.Sprintf("%s (account: %s)",
							sigs[i].String(), successfulAccounts[i].String()))
					}
				}
			}
			lggr.Errorw("Timeout waiting for transaction confirmations",
				"remaining", remaining,
				"total", len(sigs),
				"timeout", timeout,
				"unfinalizedSigs", unfinalizedSigs)

			return fmt.Errorf("timeout waiting for transaction confirmations,"+
				"remaining: %d, total: %d, timeout: %s"+
				"unfinalizedSigs: %v",
				remaining, len(sigs), timeout, unfinalizedSigs)
		case <-ticker.C:
			pollCount++
			statusRes, sigErr := solanaGoClient.GetSignatureStatuses(timeoutCtx, true, sigs...)
			if sigErr != nil {
				return fmt.Errorf("failed to get signature statuses: %w", sigErr)
			}
			if statusRes == nil {
				return errors.New("signature status response is nil")
			}
			if statusRes.Value == nil {
				return errors.New("signature status response value is nil")
			}

			unfinalizedTxCount := 0
			for i, res := range statusRes.Value {
				if res == nil {
					// Transaction status not yet available
					unfinalizedTxCount++
					continue
				}

				if res.Err != nil {
					// Transaction failed
					lggr.Errorw("Transaction failed",
						"account", successfulAccounts[i].String(),
						"signature", sigs[i].String(),
						"error", res.Err)
					return fmt.Errorf("transaction failed for account %s (sig: %s): %v",
						successfulAccounts[i].String(), sigs[i].String(), res.Err)
				}

				// Check confirmation status - we want at least "Confirmed" level
				// Solana confirmation levels: Processed < Confirmed < Finalized
				switch res.ConfirmationStatus {
				case solRpc.ConfirmationStatusProcessed, solRpc.ConfirmationStatusConfirmed:
					// Still only processed, not yet confirmed
					unfinalizedTxCount++
				case solRpc.ConfirmationStatusFinalized:
					// Transaction is finalized - we're good
					// Don't increment unfinalizedTxCount
				default:
					// Unknown status, treat as unconfirmed
					unfinalizedTxCount++
				}
			}
			remaining = unfinalizedTxCount

			// Log progress every 10 polls (5 seconds) for large batches
			if pollCount%10 == 0 {
				finalized := len(sigs) - remaining
				lggr.Infow("Confirmation progress",
					"finalized", finalized,
					"total", len(sigs),
					"pollCount", pollCount)
			}
		}
	}

	// Log successful completion
	lggr.Infow("Successfully funded all accounts",
		"accountCount", len(accounts),
		"amountSOL", solAmount)
	return nil
}

// DownloadSolanaProgramArtifactsForTest downloads the Solana program artifacts for the test environment.
//
// This is a temporary function which will be replaced by a more comprehensive package which can
// handle a more customizable download of artifacts.
func DownloadSolanaProgramArtifactsForTest(t *testing.T) {
	once.Do(func() {
		err := DownloadSolanaProgramArtifacts(t.Context(), ProgramsPath, logger.Test(t), "b0f7cd3fbdbb")
		require.NoError(t, err)
		err = DownloadSolanaCCIPProgramArtifacts(t.Context(), ProgramsPath, logger.Test(t), "")
		require.NoError(t, err)
	})
}

func generateChainsSol(t *testing.T, numChains int, commitSha string) []cldf_chain.BlockChain {
	t.Helper()

	if numChains == 0 {
		// Avoid downloading Solana program artifacts
		return nil
	}

	once.Do(func() {
		// TODO PLEX-1718 use latest contracts sha for now. Derive commit sha from go.mod once contracts are in a separate go module
		err := DownloadSolanaProgramArtifacts(t.Context(), ProgramsPath, logger.Test(t), "b0f7cd3fbdbb")
		require.NoError(t, err)
		err = DownloadSolanaCCIPProgramArtifacts(t.Context(), ProgramsPath, logger.Test(t), commitSha)
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

func fundNodesSol(t *testing.T, solChain cldf_solana.Chain, nodes []*Node) {
	for _, node := range nodes {
		solkeys, err := node.App.GetKeyStore().Solana().GetAll()
		require.NoError(t, err)
		require.Len(t, solkeys, 1)
		transmitter := solkeys[0]
		_, err = solChain.Client.RequestAirdrop(t.Context(), transmitter.PublicKey(), 1000*solana.LAMPORTS_PER_SOL, solRpc.CommitmentConfirmed)
		require.NoError(t, err)
		// we don't wait for confirmation so we don't block the tests, it'll take a while before nodes start transmitting
	}
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

var once = &sync.Once{}

// TODO: these functions should be moved to a better location

func withGetRequest[T any](ctx context.Context, url string, cb func(res *http.Response) (T, error)) (T, error) {
	var empty T

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return empty, err
	}

	res, err := (&http.Client{}).Do(req)
	if err != nil {
		return empty, err
	}
	defer res.Body.Close()

	return cb(res)
}

func DownloadTarGzReleaseAssetFromGithub(
	ctx context.Context,
	owner string,
	repo string,
	name string,
	tag string,
	cb func(r *tar.Reader, h *tar.Header) error,
) error {
	url := fmt.Sprintf(
		"https://github.com/%s/%s/releases/download/%s/%s",
		owner,
		repo,
		tag,
		name,
	)

	_, err := withGetRequest(ctx, url, func(res *http.Response) (any, error) {
		if res.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("request failed with status %d - could not download tar.gz release artifact from Github (url = '%s')", res.StatusCode, url)
		}

		gzipReader, err := gzip.NewReader(res.Body)
		if err != nil {
			return nil, err
		}
		defer gzipReader.Close()

		tarReader := tar.NewReader(gzipReader)
		for {
			header, err := tarReader.Next()
			if err == io.EOF {
				break
			}
			if err != nil {
				return nil, err
			}
			if err := cb(tarReader, header); err != nil {
				return nil, err
			}
		}

		return nil, nil
	})

	return err
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

func getSolanaCcipDependencyVersion(gomodPath string) (string, error) {
	const dependency = "github.com/smartcontractkit/chainlink-ccip/chains/solana"

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

func GetSha() (version string, err error) {
	modFilePath, err := getModFilePath()
	if err != nil {
		return "", err
	}
	go_mod_version, err := getSolanaCcipDependencyVersion(modFilePath)
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

func DownloadSolanaProgramArtifacts(ctx context.Context, dir string, lggr logger.Logger, sha string) error {
	const ownr = "smartcontractkit"
	const repo = "chainlink-solana"
	const name = "artifacts.tar.gz"

	tag := "solana-artifacts-localtest-" + sha

	if lggr != nil {
		lggr.Infof("Downloading Solana chainlink-solana program artifacts (tag = %s)", tag)
	}

	return DownloadTarGzReleaseAssetFromGithub(ctx, ownr, repo, name, tag, func(r *tar.Reader, h *tar.Header) error {
		if h.Typeflag != tar.TypeReg {
			return nil
		}

		outPath := filepath.Join(dir, filepath.Base(h.Name))
		if err := os.MkdirAll(filepath.Dir(outPath), os.ModePerm); err != nil {
			return err
		}

		outFile, err := os.Create(outPath)
		if err != nil {
			return err
		}
		defer outFile.Close()

		if _, err := io.Copy(outFile, r); err != nil {
			return err
		}

		if lggr != nil {
			lggr.Infof("Extracted Solana chainlink-solana artifact: %s", outPath)
		}

		return nil
	})
}

func DownloadSolanaCCIPProgramArtifacts(ctx context.Context, dir string, lggr logger.Logger, sha string) error {
	const ownr = "smartcontractkit"
	const repo = "chainlink-ccip"
	const name = "artifacts.tar.gz"

	if sha == "" {
		version, err := GetSha()
		if err != nil {
			return err
		}
		sha = version
	}
	tag := "solana-artifacts-localtest-" + sha

	if lggr != nil {
		lggr.Infof("Downloading Solana CCIP program artifacts (tag = %s)", tag)
	}

	return DownloadTarGzReleaseAssetFromGithub(ctx, ownr, repo, name, tag, func(r *tar.Reader, h *tar.Header) error {
		if h.Typeflag != tar.TypeReg {
			return nil
		}

		outPath := filepath.Join(dir, filepath.Base(h.Name))
		if err := os.MkdirAll(filepath.Dir(outPath), os.ModePerm); err != nil {
			return err
		}

		outFile, err := os.Create(outPath)
		if err != nil {
			return err
		}
		defer outFile.Close()

		if _, err := io.Copy(outFile, r); err != nil {
			return err
		}

		if lggr != nil {
			lggr.Infof("Extracted Solana CCIP artifact: %s", outPath)
		}

		return nil
	})
}
