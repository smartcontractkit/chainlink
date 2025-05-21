package memory

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient/simulated"
	"github.com/gagliardetto/solana-go"
	solRpc "github.com/gagliardetto/solana-go/rpc"

	"github.com/smartcontractkit/freeport"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"

	chainsel "github.com/smartcontractkit/chain-selectors"

	"golang.org/x/mod/modfile"

	"github.com/smartcontractkit/chainlink-evm/pkg/assets"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
)

type EVMChain struct {
	Backend     *simulated.Backend
	DeployerKey *bind.TransactOpts
	Users       []*bind.TransactOpts
}

type SolanaChain struct {
	Client      *solRpc.Client
	DeployerKey solana.PrivateKey
	URL         string
	WSURL       string
	KeypairPath string
}

func fundAddress(t *testing.T, from *bind.TransactOpts, to common.Address, amount *big.Int, backend *simulated.Backend) {
	ctx := t.Context()
	nonce, err := backend.Client().PendingNonceAt(ctx, from.From)
	require.NoError(t, err)
	gp, err := backend.Client().SuggestGasPrice(ctx)
	require.NoError(t, err)
	rawTx := types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		GasPrice: gp,
		Gas:      21000,
		To:       &to,
		Value:    amount,
	})
	signedTx, err := from.Signer(from.From, rawTx)
	require.NoError(t, err)
	err = backend.Client().SendTransaction(ctx, signedTx)
	require.NoError(t, err)
	backend.Commit()
}

func GenerateChains(t *testing.T, numChains int, numUsers int) map[uint64]EVMChain {
	chains := make(map[uint64]EVMChain)
	for i := 0; i < numChains; i++ {
		chainID := chainsel.TEST_90000001.EvmChainID + uint64(i)
		chains[chainID] = evmChain(t, numUsers)
	}
	return chains
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

func generateSolanaKeypair(t testing.TB) (solana.PrivateKey, string, error) {
	// Create a temporary directory that will be cleaned up after the test
	tmpDir := t.TempDir()

	privateKey, err := solana.NewRandomPrivateKey()
	if err != nil {
		return solana.PrivateKey{}, "", fmt.Errorf("failed to generate private key: %w", err)
	}

	// Convert private key bytes to JSON array
	privateKeyBytes := []byte(privateKey)

	// Convert bytes to array of integers for JSON
	intArray := make([]int, len(privateKeyBytes))
	for i, b := range privateKeyBytes {
		intArray[i] = int(b)
	}

	keypairJSON, err := json.Marshal(intArray)
	if err != nil {
		return solana.PrivateKey{}, "", fmt.Errorf("failed to marshal keypair: %w", err)
	}

	// Create the keypair file in the temporary directory
	keypairPath := filepath.Join(tmpDir, "solana-keypair.json")
	if err := os.WriteFile(keypairPath, keypairJSON, 0600); err != nil {
		return solana.PrivateKey{}, "", fmt.Errorf("failed to write keypair to file: %w", err)
	}

	return privateKey, keypairPath, nil
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

func GenerateChainsSol(t *testing.T, numChains int) map[uint64]SolanaChain {
	testSolanaChainSelectors := getTestSolanaChainSelectors()
	if len(testSolanaChainSelectors) < numChains {
		t.Fatalf("not enough test solana chain selectors available")
	}
	chains := make(map[uint64]SolanaChain)
	for i := 0; i < numChains; i++ {
		chainID := testSolanaChainSelectors[i]
		admin, keypairPath, err := generateSolanaKeypair(t)
		require.NoError(t, err)
		url, wsURL, err := solChain(t, chainID, &admin)
		require.NoError(t, err)
		client := solRpc.New(url)
		balance, err := client.GetBalance(context.Background(), admin.PublicKey(), solRpc.CommitmentConfirmed)
		require.NoError(t, err)
		require.NotEqual(t, 0, balance.Value) // auto funded 500000000.000000000 SOL
		chains[chainID] = SolanaChain{
			Client:      client,
			DeployerKey: admin,
			URL:         url,
			WSURL:       wsURL,
			KeypairPath: keypairPath,
		}
	}
	return chains
}

func GenerateChainsWithIds(t *testing.T, chainIDs []uint64, numUsers int) map[uint64]EVMChain {
	chains := make(map[uint64]EVMChain)
	for _, chainID := range chainIDs {
		chains[chainID] = evmChain(t, numUsers)
	}
	return chains
}

func evmChain(t *testing.T, numUsers int) EVMChain {
	key, err := crypto.GenerateKey()
	require.NoError(t, err)
	owner, err := bind.NewKeyedTransactorWithChainID(key, big.NewInt(1337))
	require.NoError(t, err)
	genesis := types.GenesisAlloc{
		owner.From: {Balance: assets.Ether(1_000_000).ToInt()}}
	// create a set of user keys
	var users []*bind.TransactOpts
	for j := 0; j < numUsers; j++ {
		key, err := crypto.GenerateKey()
		require.NoError(t, err)
		user, err := bind.NewKeyedTransactorWithChainID(key, big.NewInt(1337))
		require.NoError(t, err)
		users = append(users, user)
		genesis[user.From] = types.Account{Balance: assets.Ether(1_000_000).ToInt()}
	}
	// there have to be enough initial funds on each chain to allocate for all the nodes that share the given chain in the test
	backend := simulated.NewBackend(genesis, simulated.WithBlockGasLimit(50000000))
	backend.Commit() // ts will be now.
	return EVMChain{
		Backend:     backend,
		DeployerKey: owner,
		Users:       users,
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
}

var once = &sync.Once{}

func solChain(t *testing.T, chainID uint64, adminKey *solana.PrivateKey) (string, string, error) {
	t.Helper()

	once.Do(func() {
		err := DownloadSolanaCCIPProgramArtifacts(t.Context(), ProgramsPath, logger.Test(t), "")
		require.NoError(t, err)
	})

	// initialize the docker network used by CTF
	err := framework.DefaultNetwork(once)
	require.NoError(t, err)

	maxRetries := 10
	var url, wsURL string
	for i := 0; i < maxRetries; i++ {
		// solana requires 2 ports, one for http and one for ws, but only allows one to be specified
		// the other is +1 of the first one
		// must reserve 2 to avoid port conflicts in the freeport library with other tests
		// https://github.com/smartcontractkit/chainlink-testing-framework/blob/e109695d311e6ed42ca3194907571ce6454fae8d/framework/components/blockchain/blockchain.go#L39
		ports := freeport.GetN(t, 2)

		image := ""
		if runtime.GOOS == "linux" {
			image = "solanalabs/solana:v1.18.26" // TODO: workaround on linux
		}

		bcInput := &blockchain.Input{
			Image:          image,
			Type:           "solana",
			ChainID:        strconv.FormatUint(chainID, 10),
			PublicKey:      adminKey.PublicKey().String(),
			Port:           strconv.Itoa(ports[0]),
			ContractsDir:   ProgramsPath,
			SolanaPrograms: SolanaProgramIDs,
		}
		output, err := blockchain.NewBlockchainNetwork(bcInput)
		if err != nil {
			t.Logf("Error creating solana network: %v", err)
			freeport.Return(ports)
			time.Sleep(time.Second)
			maxRetries -= 1
			continue
		}
		require.NoError(t, err)
		testcontainers.CleanupContainer(t, output.Container)
		url = output.Nodes[0].ExternalHTTPUrl
		wsURL = output.Nodes[0].ExternalWSUrl
		break
	}
	require.NoError(t, err)

	// Wait for api server to boot
	client := solRpc.New(url)
	var ready bool
	for i := 0; i < 30; i++ {
		time.Sleep(time.Second)
		out, err := client.GetHealth(t.Context())
		if err != nil || out != solRpc.HealthOk {
			t.Logf("API server not ready yet (attempt %d)\n", i+1)
			continue
		}
		ready = true
		break
	}
	if !ready {
		t.Logf("solana-test-validator is not ready after 30 attempts")
	}
	require.True(t, ready)
	t.Logf("solana-test-validator is ready at %s", url)
	time.Sleep(15 * time.Second) // we have slot errors that force retries if the chain is not given enough time to boot

	return url, wsURL, nil
}

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
