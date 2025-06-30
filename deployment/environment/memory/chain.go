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

	"github.com/gagliardetto/solana-go"
	solRpc "github.com/gagliardetto/solana-go/rpc"
	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/require"
	"golang.org/x/mod/modfile"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	cldf_solana_provider "github.com/smartcontractkit/chainlink-deployments-framework/chain/solana/provider"
)

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
		sig, err := solanaGoClient.RequestAirdrop(ctx, account, solAmount*solana.LAMPORTS_PER_SOL, solRpc.CommitmentConfirmed)
		if err != nil {
			return err
		}
		sigs = append(sigs, sig)
	}

	const timeout = 10 * time.Second
	const pollInterval = 50 * time.Millisecond

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

			unconfirmedTxCount := 0
			for _, res := range statusRes.Value {
				if res == nil || res.ConfirmationStatus == solRpc.ConfirmationStatusProcessed {
					unconfirmedTxCount++
				}
			}
			remaining = unconfirmedTxCount
		}
	}
	return nil
}

func generateChainsSol(t *testing.T, numChains int) []cldf_chain.BlockChain {
	t.Helper()

	once.Do(func() {
		err := DownloadSolanaCCIPProgramArtifacts(t.Context(), ProgramsPath, logger.Test(t), "")
		require.NoError(t, err)
	})

	testSolanaChainSelectors := getTestSolanaChainSelectors()
	if len(testSolanaChainSelectors) < numChains {
		t.Fatalf("not enough test solana chain selectors available")
	}

	chains := make([]cldf_chain.BlockChain, 0, numChains)
	for i := 0; i < numChains; i++ {
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
