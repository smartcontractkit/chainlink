package bulletprooftxmanager_test

import (
	"errors"
	"testing"

	"github.com/smartcontractkit/chainlink/core/services/bulletprooftxmanager"
	"github.com/stretchr/testify/assert"
)

func TestBulletproofTxManager_Errors(t *testing.T) {
	t.Parallel()
	randomError := bulletprooftxmanager.NewSendError("some old bollocks")

	// IsNonceTooLowError
	assert.False(t, randomError.IsNonceTooLowError())

	// Geth
	err := bulletprooftxmanager.NewSendError("nonce too low")
	assert.True(t, err.IsNonceTooLowError())
	// Parity
	err = bulletprooftxmanager.NewSendError("Transaction nonce is too low. Try incrementing the nonce.")
	assert.True(t, err.IsNonceTooLowError())

	// IsReplacementUnderpriced

	// Geth
	err = bulletprooftxmanager.NewSendError("replacement transaction underpriced")
	assert.True(t, err.IsReplacementUnderpriced())
	// Parity
	s := "Transaction gas price 100wei is too low. There is another transaction with same nonce in the queue with gas price 150wei. Try increasing the gas price or incrementing the nonce."
	err = bulletprooftxmanager.NewSendError(s)
	assert.True(t, err.IsReplacementUnderpriced())
	s = "There are too many transactions in the queue. Your transaction was dropped due to limit. Try increasing the fee."
	err = bulletprooftxmanager.NewSendError(s)
	assert.False(t, err.IsReplacementUnderpriced())

	// IsTransactionAlreadyInMempool
	assert.False(t, randomError.IsTransactionAlreadyInMempool())

	// Geth
	// I have seen this in log output
	err = bulletprooftxmanager.NewSendError("known transaction: 0x7f657507aee0511e36d2d1972a6b22e917cc89f92b6c12c4dbd57eaabb236960")
	assert.True(t, err.IsTransactionAlreadyInMempool())
	// This comes from the geth source - https://github.com/ethereum/go-ethereum/blob/eb9d7d15ecf08cd5104e01a8af64489f01f700b0/core/tx_pool.go#L57
	err = bulletprooftxmanager.NewSendError("already known")
	assert.True(t, err.IsTransactionAlreadyInMempool())
	// This one is present in the light client (?!)
	err = bulletprooftxmanager.NewSendError("Known transaction (7f65)")
	assert.True(t, err.IsTransactionAlreadyInMempool())
	// Parity
	s = "Transaction with the same hash was already imported."
	err = bulletprooftxmanager.NewSendError(s)
	assert.True(t, err.IsTransactionAlreadyInMempool())

	// IsTerminallyUnderpriced
	assert.False(t, randomError.IsTerminallyUnderpriced())

	// Geth
	err = bulletprooftxmanager.NewSendError("transaction underpriced")
	assert.True(t, err.IsTerminallyUnderpriced())
	// Parity
	err = bulletprooftxmanager.NewSendError("There are too many transactions in the queue. Your transaction was dropped due to limit. Try increasing the fee.")
	assert.False(t, err.IsTerminallyUnderpriced())
	err = bulletprooftxmanager.NewSendError("Transaction gas price is too low. It does not satisfy your node's minimal gas price (minimal: 100 got: 50). Try increasing the gas price.")
	assert.True(t, err.IsTerminallyUnderpriced())

	// IsTemporarilyUnderpriced
	// Parity
	err = bulletprooftxmanager.NewSendError("There are too many transactions in the queue. Your transaction was dropped due to limit. Try increasing the fee.")
	assert.True(t, err.IsTemporarilyUnderpriced())
	err = bulletprooftxmanager.NewSendError("Transaction gas price is too low. It does not satisfy your node's minimal gas price (minimal: 100 got: 50). Try increasing the gas price.")
	assert.False(t, err.IsTemporarilyUnderpriced())

	// IsInsufficientEth
	// Geth
	err = bulletprooftxmanager.NewSendError("insufficient funds for transfer")
	assert.True(t, err.IsInsufficientEth())
	err = bulletprooftxmanager.NewSendError("insufficient funds for gas * price + value")
	assert.True(t, err.IsInsufficientEth())
	err = bulletprooftxmanager.NewSendError("insufficient balance for transfer")
	assert.True(t, err.IsInsufficientEth())
	// Parity
	err = bulletprooftxmanager.NewSendError("Insufficient balance for transaction. Balance=100.25, Cost=200.50")
	assert.True(t, err.IsInsufficientEth())
	err = bulletprooftxmanager.NewSendError("Insufficient funds. The account you tried to send transaction from does not have enough funds. Required 200.50 and got: 100.25.")
	assert.True(t, err.IsInsufficientEth())
	// Nil
	err = bulletprooftxmanager.SendError(nil)
	assert.False(t, err.IsInsufficientEth())
}

func TestBulletproofTxManager_Errors_Fatal(t *testing.T) {
	t.Parallel()

	tests := []struct {
		errStr      string
		expectFatal bool
	}{
		{"some old bollocks", false},

		// Geth
		{"insufficient funds for transfer", false},

		{"exceeds block gas limit", true},
		{"invalid sender", true},
		{"negative value", true},
		{"oversized data", true},
		{"gas uint64 overflow", true},
		{"intrinsic gas too low", true},
		{"nonce too high", true},

		// Parity
		{"Insufficient funds. The account you tried to send transaction from does not have enough funds. Required 100 and got: 50.", false},

		{"Supplied gas is beyond limit.", true},
		{"Sender is banned in local queue.", true},
		{"Recipient is banned in local queue.", true},
		{"Code is banned in local queue.", true},
		{"Transaction is not permitted.", true},
		{"Transaction is too big, see chain specification for the limit.", true},
		{"Transaction gas is too low. There is not enough gas to cover minimal cost of the transaction (minimal: 100 got: 50) Try increasing supplied gas.", true},
		{"Transaction cost exceeds current gas limit. Limit: 50, got: 100. Try decreasing supplied gas.", true},
		{"Invalid signature: some old bollocks", true},
		{"Invalid RLP data: some old bollocks", true},
	}

	for _, test := range tests {
		t.Run(test.errStr, func(t *testing.T) {
			err := bulletprooftxmanager.SendError(errors.New(test.errStr))
			assert.Equal(t, test.expectFatal, err.Fatal())
		})
	}
}
