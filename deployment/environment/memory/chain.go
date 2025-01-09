package memory

import (
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient/simulated"
	"github.com/gagliardetto/solana-go"
	solRpc "github.com/gagliardetto/solana-go/rpc"
	"github.com/mr-tron/base58"

	"github.com/stretchr/testify/require"

	solTestUtil "github.com/smartcontractkit/chainlink-ccip/chains/solana/contracts/tests/testutils"

	chainsel "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/assets"
)

type EVMChain struct {
	Backend     *simulated.Backend
	DeployerKey *bind.TransactOpts
	Users       []*bind.TransactOpts
}

type SolanaChain struct {
	Client      *solRpc.Client
	DeployerKey *solana.PrivateKey
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

func generateAndStoreKeypair() (solana.PrivateKey, string, error) {
	// Generate a random private key
	privateKey, err := solana.NewRandomPrivateKey()
	if err != nil {
		return solana.PrivateKey{}, "", fmt.Errorf("failed to generate private key: %w", err)
	}

	privateKeyBytes, err := base58.Decode(privateKey.String())
	if err != nil {
		return solana.PrivateKey{}, "", fmt.Errorf("failed to decode Base58 private key: %w", err)
	}

	intArray := make([]int, len(privateKeyBytes))
	for i, b := range privateKeyBytes {
		intArray[i] = int(b)
	}

	// Marshal the integer array to JSON
	keypairJSON, err := json.Marshal(intArray)
	if err != nil {
		return solana.PrivateKey{}, "", fmt.Errorf("failed to marshal keypair to JSON: %w", err)
	}

	// Create a temporary file
	tempFile, err := os.CreateTemp("", "solana-keypair-*.json")
	if err != nil {
		return solana.PrivateKey{}, "", fmt.Errorf("failed to create temporary file: %w", err)
	}
	defer tempFile.Close()

	// Write the keypair data to the file
	if err := os.WriteFile(tempFile.Name(), keypairJSON, 0600); err != nil {
		return solana.PrivateKey{}, "", fmt.Errorf("failed to write keypair to file: %w", err)
	}

	// Return the path to the temporary file
	return privateKey, tempFile.Name(), nil
}

func GenerateChainsSol(t *testing.T, numChains int) map[uint64]SolanaChain {
	testSolanaChainSelectors := getTestSolanaChainSelectors()
	if len(testSolanaChainSelectors) < numChains {
		t.Fatalf("not enough test solana chain selectors available")
	}
	chains := make(map[uint64]SolanaChain)
	for i := 0; i < numChains; i++ {
		chainID := testSolanaChainSelectors[i]
		url, wsurl := solTestUtil.SetupLocalSolNodeWithFlags(t)
		admin, keypairPath, gerr := generateAndStoreKeypair()
		// byteSlice, err := base58.Decode(admin)
		t.Log("keypairPath", keypairPath)
		t.Log("admin private key", admin)
		key, err := solana.PrivateKeyFromSolanaKeygenFile(keypairPath)
		require.NoError(t, err)
		t.Log("keypair key", key)
		require.NoError(t, gerr)
		solTestUtil.FundTestAccounts(t, []solana.PublicKey{admin.PublicKey()}, url)
		require.NoError(t, gerr)
		chains[chainID] = SolanaChain{
			Client:      solRpc.New(url),
			DeployerKey: &admin,
			URL:         url,
			WSURL:       wsurl,
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
