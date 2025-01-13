package memory

import (
	"encoding/json"
	"math/big"
	"os"
	"path"
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

	solTestUtil "github.com/smartcontractkit/chainlink-ccip/chains/solana/contracts/tests/testutils"

	chainsel "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"

	chainselectors "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/assets"
)

type EVMChain struct {
	Backend     *simulated.Backend
	DeployerKey *bind.TransactOpts
	Users       []*bind.TransactOpts
}

type SolanaChain struct {
	Client      *solRpc.Client
	URL         string
	WSURL       string
	DeployerKey solana.PrivateKey
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

func GenerateChainsSol(t *testing.T, numChains int) map[uint64]SolanaChain {
	testSolanaChainSelectors := getTestSolanaChainSelectors()
	if len(testSolanaChainSelectors) < numChains {
		t.Fatalf("not enough test solana chain selectors available")
	}
	chains := make(map[uint64]SolanaChain)
	for i := 0; i < numChains; i++ {
		chainID := testSolanaChainSelectors[i]
		solChain := solChain(t)
		admin := solChain.DeployerKey
		solTestUtil.FundTestAccounts(t, []solana.PublicKey{admin.PublicKey()}, solChain.URL)
		chains[chainID] = SolanaChain{
			Client:      solChain.Client,
			DeployerKey: solChain.DeployerKey,
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

var once = &sync.Once{}

func solChain(t *testing.T) SolanaChain {
	t.Helper()

	// initialize the docker network used by CTF
	err := framework.DefaultNetwork(once)
	require.NoError(t, err)

	deployerKey, err := solana.NewRandomPrivateKey()
	require.NoError(t, err)

	t.TempDir()
	// store the generated keypair somewhere
	bytes, err := json.Marshal([]byte(deployerKey))
	require.NoError(t, err)
	keypairPath := path.Join(t.TempDir(), "solana-keypair.json")
	err = os.WriteFile(keypairPath, bytes, 0600)
	require.NoError(t, err)

	port := freeport.GetOne(t)

	bcInput := &blockchain.Input{
		Type:      "solana",
		ChainID:   chainselectors.SOLANA_DEVNET.ChainID,
		PublicKey: deployerKey.PublicKey().String(),
		Port:      strconv.Itoa(port),
		// TODO: ContractsDir & SolanaPrograms via env vars
	}
	output, err := blockchain.NewBlockchainNetwork(bcInput)
	require.NoError(t, err)
	testcontainers.CleanupContainer(t, output.Container)

	url := output.Nodes[0].HostHTTPUrl
	wsURL := output.Nodes[0].HostWSUrl

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

	return SolanaChain{
		Client:      client,
		URL:         url,
		WSURL:       wsURL,
		DeployerKey: deployerKey,
	}
}
