package memory

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"strconv"
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
	"github.com/hashicorp/consul/sdk/freeport"

	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"

	chainsel "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"

	"github.com/smartcontractkit/chainlink-integrations/evm/assets"
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
	ctx := tests.Context(t)
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
	ctx context.Context, t *testing.T, accounts []solana.PublicKey, solAmount uint64, solanaGoClient *solRpc.Client,
) {
	t.Helper()

	var sigs = make([]solana.Signature, 0, len(accounts))
	for _, account := range accounts {
		sig, err := solanaGoClient.RequestAirdrop(ctx, account, solAmount*solana.LAMPORTS_PER_SOL, solRpc.CommitmentConfirmed)
		require.NoError(t, err)
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
			require.NoError(t, errors.New("unable to find transaction within timeout"))
		case <-ticker.C:
			statusRes, sigErr := solanaGoClient.GetSignatureStatuses(ctx, true, sigs...)
			require.NoError(t, sigErr)
			require.NotNil(t, statusRes)
			require.NotNil(t, statusRes.Value)

			unconfirmedTxCount := 0
			for _, res := range statusRes.Value {
				if res == nil || res.ConfirmationStatus == solRpc.ConfirmationStatusProcessed {
					unconfirmedTxCount++
				}
			}
			remaining = unconfirmedTxCount
		}
	}
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
	"ccip_router":                    "Ccip842gzYHhvdDkSyi2YVCoAWPbYJoApMFzSxQroE9C",
	"test_token_pool":                "JuCcZ4smxAYv9QHJ36jshA7pA3FuQ3vQeWLUeAtZduJ",
	"example_burnmint_token_pool":    "41FGToCmdaWa1dgZLKFAjvmx6e6AjVTX7SVRibvsMGVB",
	"example_lockrelease_token_pool": "8eqh8wppT9c5rw4ERqNCffvU6cNFJWff9WmkcYtmGiqC",
	"fee_quoter":                     "FeeQPGkKDeRV1MgoYfMH6L8o3KeuYjwUZrgn4LRKfjHi",
	"test_ccip_receiver":             "EvhgrPhTDt4LcSPS2kfJgH6T6XWZ6wT3X9ncDGLT1vui",
	"ccip_offramp":                   "offqSMQWgQud6WJz694LRzkeN5kMYpCHTpXQr3Rkcjm",
	"mcm":                            "5vNJx78mz7KVMjhuipyr9jKBKcMrKYGdjGkgE4LUmjKk",
	"timelock":                       "DoajfR5tK24xVw51fWcawUZWhAXD8yrBJVacc13neVQA",
	"access_controller":              "6KsN58MTnRQ8FfPaXHiFPPFGDRioikj9CdPvPxZJdCjb",
	"external_program_cpi_stub":      "2zZwzyptLqwFJFEFxjPvrdhiGpH9pJ3MfrrmZX6NTKxm",
	"rmn_remote":                     "RmnXLft1mSEwDgMKu2okYuHkiazxntFFcZFrrcXxYg7",
}

var once = &sync.Once{}

func solChain(t *testing.T, chainID uint64, adminKey *solana.PrivateKey) (string, string, error) {
	t.Helper()

	// initialize the docker network used by CTF
	err := framework.DefaultNetwork(once)
	require.NoError(t, err)

	maxRetries := 10
	var url, wsURL string
	for i := 0; i < maxRetries; i++ {
		port := freeport.GetOne(t)

		bcInput := &blockchain.Input{
			Type:           "solana",
			ChainID:        strconv.FormatUint(chainID, 10),
			PublicKey:      adminKey.PublicKey().String(),
			Port:           strconv.Itoa(port),
			ContractsDir:   ProgramsPath,
			SolanaPrograms: SolanaProgramIDs,
		}
		output, err := blockchain.NewBlockchainNetwork(bcInput)
		if err != nil {
			t.Logf("Error creating solana network: %v", err)
			time.Sleep(time.Second)
			maxRetries -= 1
			continue
		}
		require.NoError(t, err)
		testcontainers.CleanupContainer(t, output.Container)
		url = output.Nodes[0].HostHTTPUrl
		wsURL = output.Nodes[0].HostWSUrl
		break
	}
	require.NoError(t, err)

	// Wait for api server to boot
	client := solRpc.New(url)
	var ready bool
	for i := 0; i < 30; i++ {
		time.Sleep(time.Second)
		out, err := client.GetHealth(tests.Context(t))
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
