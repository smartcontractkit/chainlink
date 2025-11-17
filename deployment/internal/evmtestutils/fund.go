package evmtestutils

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient/simulated"
	"github.com/stretchr/testify/require"
)

// FundAddress funds to an EVM address using a given transaction options.
func FundAddress(t *testing.T, from *bind.TransactOpts, to common.Address, amount *big.Int, backend *simulated.Backend) {
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
