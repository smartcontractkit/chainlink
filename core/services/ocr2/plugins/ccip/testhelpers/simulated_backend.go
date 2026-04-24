package testhelpers

import (
	"context"
	"errors"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/eth/ethconfig"
	"github.com/ethereum/go-ethereum/ethclient/simulated"
	"github.com/stretchr/testify/require"
)

// FirstBlockAge is used to compute first block's timestamp in SimulatedBackend (time.Now() - FirstBlockAge)
const FirstBlockAge = 24 * time.Hour

func SetupChain(t *testing.T) (*simulated.Backend, *bind.TransactOpts) {
	key, err := crypto.GenerateKey()
	require.NoError(t, err)
	user, err := bind.NewKeyedTransactorWithChainID(key, big.NewInt(1337))
	require.NoError(t, err)
	chain := simulated.NewBackend(ethtypes.GenesisAlloc{
		user.From: {
			Balance: new(big.Int).Mul(big.NewInt(1000), big.NewInt(1e18)),
		},
		common.Address{}: {
			Balance: new(big.Int).Mul(big.NewInt(1000), big.NewInt(1e18)),
		},
	}, simulated.WithBlockGasLimit(ethconfig.Defaults.Miner.GasCeil))
	require.NoError(t, err)
	chain.Commit()
	return chain, user
}

func ConfirmTxs(t *testing.T, txs []*ethtypes.Transaction, chain *Backend) {
	chain.Commit()
	ctx := t.Context()
	for _, tx := range txs {
		rec, err := bind.WaitMined(ctx, chain.Client(), tx)
		require.NoError(t, err)
		require.Equal(t, uint64(1), rec.Status)
		if rec.Status == uint64(1) {
			r, err := getFailureReason(chain.Client(), common.Address{}, tx, rec.BlockNumber)
			t.Log("Reverted", r, err)
		}
	}
}

func createCallMsgFromTransaction(from common.Address, tx *ethtypes.Transaction) ethereum.CallMsg {
	return ethereum.CallMsg{
		From:     from,
		To:       tx.To(),
		Gas:      tx.Gas(),
		GasPrice: tx.GasPrice(),
		Value:    tx.Value(),
		Data:     tx.Data(),
	}
}
func getFailureReason(client simulated.Client, from common.Address, tx *ethtypes.Transaction, blockNumber *big.Int) (string, error) {
	code, err := client.CallContract(context.Background(), createCallMsgFromTransaction(from, tx), blockNumber)
	if err != nil {
		return "", err
	}
	if len(code) == 0 {
		return "", errors.New("no error message or out of gas")
	}
	return string(code), nil
}
