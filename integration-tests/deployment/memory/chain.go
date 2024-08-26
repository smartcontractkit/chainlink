package memory

import (
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	gethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/params"
	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/require"
)

type EVMChain struct {
	Backend     *backends.SimulatedBackend
	DeployerKey *bind.TransactOpts
}

// CCIP relies on block timestamps, but SimulatedBackend uses by default clock starting from 1970-01-01
// This trick is used to move the clock closer to the current time. We set first block to be X hours ago.
// Tests create plenty of transactions so this number can't be too low, every new block mined will tick the clock,
// if you mine more than "X hours" transactions, SimulatedBackend will panic because generated timestamps will be in the future.
func tweakChainTimestamp(t *testing.T, backend *backends.SimulatedBackend, tweak time.Duration) {
	blockTime := time.Unix(int64(backend.Blockchain().CurrentHeader().Time), 0)
	sinceBlockTime := time.Since(blockTime)
	diff := sinceBlockTime - tweak
	err := backend.AdjustTime(diff)
	require.NoError(t, err, "unable to adjust time on simulated chain")
	backend.Commit()
	backend.Commit()
}

func fundAddress(t *testing.T, from *bind.TransactOpts, to common.Address, amount *big.Int, backend *backends.SimulatedBackend) {
	nonce, err := backend.PendingNonceAt(Context(t), from.From)
	require.NoError(t, err)
	gp, err := backend.SuggestGasPrice(Context(t))
	require.NoError(t, err)
	rawTx := gethtypes.NewTx(&gethtypes.LegacyTx{
		Nonce:    nonce,
		GasPrice: gp,
		Gas:      21000,
		To:       &to,
		Value:    amount,
	})
	signedTx, err := from.Signer(from.From, rawTx)
	require.NoError(t, err)
	err = backend.SendTransaction(Context(t), signedTx)
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
		backend := backends.NewSimulatedBackend(core.GenesisAlloc{
			owner.From: {Balance: big.NewInt(0).Mul(big.NewInt(100), big.NewInt(params.Ether))}}, 10000000)
		tweakChainTimestamp(t, backend, time.Hour*8)
		chains[chainID] = EVMChain{
			Backend:     backend,
			DeployerKey: owner,
		}
	}
	return chains
}
