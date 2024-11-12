package memory

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient/simulated"
	"github.com/ethereum/go-ethereum/params"
	"github.com/stretchr/testify/require"

	chainsel "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
)

type EVMChain struct {
	Backend     *simulated.Backend
	DeployerKey *bind.TransactOpts
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

func GenerateChains(t *testing.T, numChains int) map[uint64]EVMChain {
	chains := make(map[uint64]EVMChain)
	for i := 0; i < numChains; i++ {
		chainID := chainsel.TEST_90000001.EvmChainID + uint64(i)
		key, err := crypto.GenerateKey()
		require.NoError(t, err)
		owner, err := bind.NewKeyedTransactorWithChainID(key, big.NewInt(1337))
		require.NoError(t, err)
		// there have to be enough initial funds on each chain to allocate for all the nodes that share the given chain in the test
		backend := simulated.NewBackend(types.GenesisAlloc{
			owner.From: {Balance: big.NewInt(0).Mul(big.NewInt(7000), big.NewInt(params.Ether))}},
			simulated.WithBlockGasLimit(50000000))
		backend.Commit() // ts will be now.
		chains[chainID] = EVMChain{
			Backend:     backend,
			DeployerKey: owner,
		}
	}
	return chains
}

func GenerateChainsWithIds(t *testing.T, chainIDs []uint64) map[uint64]EVMChain {
	chains := make(map[uint64]EVMChain)
	for _, chainID := range chainIDs {
		key, err := crypto.GenerateKey()
		require.NoError(t, err)
		owner, err := bind.NewKeyedTransactorWithChainID(key, big.NewInt(1337))
		require.NoError(t, err)
		backend := simulated.NewBackend(types.GenesisAlloc{
			owner.From: {Balance: big.NewInt(0).Mul(big.NewInt(700000), big.NewInt(params.Ether))}},
			simulated.WithBlockGasLimit(10000000))
		backend.Commit() // Note initializes block timestamp to now().
		chains[chainID] = EVMChain{
			Backend:     backend,
			DeployerKey: owner,
		}
	}
	return chains
}
