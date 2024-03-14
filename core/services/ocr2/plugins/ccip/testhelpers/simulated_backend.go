package testhelpers

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/eth/ethconfig"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
)

// FirstBlockAge is used to compute first block's timestamp in SimulatedBackend (time.Now() - FirstBlockAge)
const FirstBlockAge = 24 * time.Hour

func SetupChain(t *testing.T) (*backends.SimulatedBackend, *bind.TransactOpts) {
	key, err := crypto.GenerateKey()
	require.NoError(t, err)
	user, err := bind.NewKeyedTransactorWithChainID(key, big.NewInt(1337))
	require.NoError(t, err)
	chain := backends.NewSimulatedBackend(core.GenesisAlloc{
		user.From: {Balance: new(big.Int).Mul(big.NewInt(1000), big.NewInt(1e18))}},
		ethconfig.Defaults.Miner.GasCeil)
	// CCIP relies on block timestamps, but SimulatedBackend uses by default clock starting from 1970-01-01
	// This trick is used to move the clock closer to the current time. We set first block to be X hours ago.
	// Tests create plenty of transactions so this number can't be too low, every new block mined will tick the clock,
	// if you mine more than "X hours" transactions, SimulatedBackend will panic because generated timestamps will be in the future.
	// IMPORTANT: Any adjustments to FirstBlockAge will automatically update PermissionLessExecutionThresholdSeconds in tests
	blockTime := time.UnixMilli(int64(chain.Blockchain().CurrentHeader().Time))
	err = chain.AdjustTime(time.Since(blockTime) - FirstBlockAge)
	require.NoError(t, err)
	chain.Commit()
	return chain, user
}

type EthKeyStoreSim struct {
	ETHKS keystore.Eth
	CSAKS keystore.CSA
}

func (ks EthKeyStoreSim) CSA() keystore.CSA {
	return ks.CSAKS
}

func (ks EthKeyStoreSim) Eth() keystore.Eth {
	return ks.ETHKS
}

func (ks EthKeyStoreSim) SignTx(address common.Address, tx *ethtypes.Transaction, chainID *big.Int) (*ethtypes.Transaction, error) {
	if chainID.String() == "1000" {
		// A terrible hack, just for the multichain test. All simulation clients run on chainID 1337.
		// We let the DestChainSelector actually use 1337 to make sure the offchainConfig digests are properly generated.
		return ks.ETHKS.SignTx(context.Background(), address, tx, big.NewInt(1337))
	}
	return ks.ETHKS.SignTx(context.Background(), address, tx, chainID)
}

var _ keystore.Eth = EthKeyStoreSim{}.ETHKS

func ConfirmTxs(t *testing.T, txs []*ethtypes.Transaction, chain *backends.SimulatedBackend) {
	chain.Commit()
	for _, tx := range txs {
		rec, err := bind.WaitMined(context.Background(), chain, tx)
		require.NoError(t, err)
		require.Equal(t, uint64(1), rec.Status)
	}
}
