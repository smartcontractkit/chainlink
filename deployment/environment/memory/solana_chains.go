package memory

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"runtime"
	"sync"
	"testing"
	"time"

	"github.com/Masterminds/semver/v3"
	"github.com/gagliardetto/solana-go"
	solRpc "github.com/gagliardetto/solana-go/rpc"
	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	cldf_solana "github.com/smartcontractkit/chainlink-deployments-framework/chain/solana"
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
