package eth_test

import (
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/services/eth"
	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newSendErrorWrapped(s string) *eth.SendError {
	return eth.NewSendError(errors.Wrap(errors.New(s), "wrapped with some old bollocks"))
}

func Test_Eth_Errors(t *testing.T) {
	t.Parallel()
	var err *eth.SendError
	randomError := eth.NewSendErrorS("some old bollocks")

	t.Run("IsNonceTooLowError", func(t *testing.T) {
		assert.False(t, randomError.IsNonceTooLowError())

		// Geth
		err = eth.NewSendErrorS("nonce too low")
		assert.True(t, err.IsNonceTooLowError())
		err = newSendErrorWrapped("nonce too low")
		assert.True(t, err.IsNonceTooLowError())
		// Parity
		err = eth.NewSendErrorS("Transaction nonce is too low. Try incrementing the nonce.")
		assert.True(t, err.IsNonceTooLowError())
		err = newSendErrorWrapped("Transaction nonce is too low. Try incrementing the nonce.")
		assert.True(t, err.IsNonceTooLowError())
		// Arbitrum
		err = eth.NewSendErrorS("transaction rejected: nonce too low")
		assert.True(t, err.IsNonceTooLowError())
		err = newSendErrorWrapped("transaction rejected: nonce too low")
		assert.True(t, err.IsNonceTooLowError())
		// Optimism
		err = eth.NewSendErrorS("invalid transaction: nonce too low")
		assert.True(t, err.IsNonceTooLowError())
		err = newSendErrorWrapped("invalid transaction: nonce too low")
		assert.True(t, err.IsNonceTooLowError())
	})

	t.Run("IsReplacementUnderpriced", func(t *testing.T) {
		// Geth
		err = eth.NewSendErrorS("replacement transaction underpriced")
		assert.True(t, err.IsReplacementUnderpriced())
		err = newSendErrorWrapped("replacement transaction underpriced")
		assert.True(t, err.IsReplacementUnderpriced())
		// Parity
		s := "Transaction gas price 100wei is too low. There is another transaction with same nonce in the queue with gas price 150wei. Try increasing the gas price or incrementing the nonce."
		err = eth.NewSendErrorS(s)
		assert.True(t, err.IsReplacementUnderpriced())
		err = newSendErrorWrapped(s)
		assert.True(t, err.IsReplacementUnderpriced())

		s = "There are too many transactions in the queue. Your transaction was dropped due to limit. Try increasing the fee."
		err = eth.NewSendErrorS(s)
		assert.False(t, err.IsReplacementUnderpriced())
	})

	t.Run("IsTransactionAlreadyInMempool", func(t *testing.T) {
		assert.False(t, randomError.IsTransactionAlreadyInMempool())

		// Geth
		// I have seen this in log output
		err = eth.NewSendErrorS("known transaction: 0x7f657507aee0511e36d2d1972a6b22e917cc89f92b6c12c4dbd57eaabb236960")
		assert.True(t, err.IsTransactionAlreadyInMempool())
		err = newSendErrorWrapped("known transaction: 0x7f657507aee0511e36d2d1972a6b22e917cc89f92b6c12c4dbd57eaabb236960")
		assert.True(t, err.IsTransactionAlreadyInMempool())
		// This comes from the geth source - https://github.com/ethereum/go-ethereum/blob/eb9d7d15ecf08cd5104e01a8af64489f01f700b0/core/tx_pool.go#L57
		err = eth.NewSendErrorS("already known")
		assert.True(t, err.IsTransactionAlreadyInMempool())
		// This one is present in the light client (?!)
		err = eth.NewSendErrorS("Known transaction (7f65)")
		assert.True(t, err.IsTransactionAlreadyInMempool())
		// Parity
		s := "Transaction with the same hash was already imported."
		err = eth.NewSendErrorS(s)
		assert.True(t, err.IsTransactionAlreadyInMempool())
	})

	t.Run("IsTerminallyUnderpriced", func(t *testing.T) {
		assert.False(t, randomError.IsTerminallyUnderpriced())

		// Geth
		err = eth.NewSendErrorS("transaction underpriced")
		assert.True(t, err.IsTerminallyUnderpriced())
		err = newSendErrorWrapped("transaction underpriced")
		assert.True(t, err.IsTerminallyUnderpriced())

		err = eth.NewSendErrorS("replacement transaction underpriced")
		assert.False(t, err.IsTerminallyUnderpriced())
		// Parity
		err = eth.NewSendErrorS("There are too many transactions in the queue. Your transaction was dropped due to limit. Try increasing the fee.")
		assert.False(t, err.IsTerminallyUnderpriced())
		err = eth.NewSendErrorS("Transaction gas price is too low. It does not satisfy your node's minimal gas price (minimal: 100 got: 50). Try increasing the gas price.")
		assert.True(t, err.IsTerminallyUnderpriced())
	})

	t.Run("IsTemporarilyUnderpriced", func(t *testing.T) {
		// Parity
		err = eth.NewSendErrorS("There are too many transactions in the queue. Your transaction was dropped due to limit. Try increasing the fee.")
		assert.True(t, err.IsTemporarilyUnderpriced())
		err = newSendErrorWrapped("There are too many transactions in the queue. Your transaction was dropped due to limit. Try increasing the fee.")
		assert.True(t, err.IsTemporarilyUnderpriced())
		err = eth.NewSendErrorS("Transaction gas price is too low. It does not satisfy your node's minimal gas price (minimal: 100 got: 50). Try increasing the gas price.")
		assert.False(t, err.IsTemporarilyUnderpriced())
	})

	t.Run("IsInsufficientEth", func(t *testing.T) {
		// Geth
		err = eth.NewSendErrorS("insufficient funds for transfer")
		assert.True(t, err.IsInsufficientEth())
		err = newSendErrorWrapped("insufficient funds for transfer")
		assert.True(t, err.IsInsufficientEth())
		err = eth.NewSendErrorS("insufficient funds for gas * price + value")
		assert.True(t, err.IsInsufficientEth())
		err = eth.NewSendErrorS("insufficient balance for transfer")
		assert.True(t, err.IsInsufficientEth())
		// Parity
		err = eth.NewSendErrorS("Insufficient balance for transaction. Balance=100.25, Cost=200.50")
		assert.True(t, err.IsInsufficientEth())
		err = eth.NewSendErrorS("Insufficient funds. The account you tried to send transaction from does not have enough funds. Required 200.50 and got: 100.25.")
		assert.True(t, err.IsInsufficientEth())
		// Arbitrum
		err = eth.NewSendErrorS("transaction rejected: insufficient funds for gas * price + value")
		assert.True(t, err.IsInsufficientEth())
		// Optimism
		err = eth.NewSendErrorS("invalid transaction: insufficient funds for gas * price + value")
		assert.True(t, err.IsInsufficientEth())
		// Nil
		err = eth.NewSendError(nil)
		assert.False(t, err.IsInsufficientEth())
	})

	t.Run("IsTooExpensive", func(t *testing.T) {
		// Geth
		err = eth.NewSendErrorS("tx fee (1.10 ether) exceeds the configured cap (1.00 ether)")
		assert.True(t, err.IsTooExpensive())
		err = newSendErrorWrapped("tx fee (1.10 ether) exceeds the configured cap (1.00 ether)")
		assert.True(t, err.IsTooExpensive())

		assert.False(t, randomError.IsTooExpensive())
		// Nil
		err = eth.NewSendError(nil)
		assert.False(t, err.IsTooExpensive())
	})
}

func Test_Eth_Errors_Fatal(t *testing.T) {
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
			err := eth.NewSendError(errors.New(test.errStr))
			assert.Equal(t, test.expectFatal, err.Fatal())
		})
	}
}

func Test_ExtractRevertReasonFromRPCError(t *testing.T) {
	message := "important revert reason"
	messageHex := utils.RemoveHexPrefix(hexutil.Encode([]byte(message)))
	sigHash := "12345678"
	var jsonErr error = &eth.JsonError{
		Code:    1,
		Data:    fmt.Sprintf("0x%s%s", sigHash, messageHex),
		Message: "something different",
	}

	t.Run("it extracts revert reasons when present", func(tt *testing.T) {
		revertReason, err := eth.ExtractRevertReasonFromRPCError(jsonErr)
		require.NoError(t, err)
		require.Equal(t, message, revertReason)
	})

	t.Run("it unwraps wrapped errors", func(tt *testing.T) {
		wrappedErr := errors.Wrap(jsonErr, "wrapped message")
		revertReason, err := eth.ExtractRevertReasonFromRPCError(wrappedErr)
		require.NoError(t, err)
		require.Equal(t, message, revertReason)
	})

	t.Run("it unwraps multi-wrapped errors", func(tt *testing.T) {
		wrappedErr := errors.Wrap(jsonErr, "wrapped message")
		wrappedErr = errors.Wrap(wrappedErr, "wrapped again!!")
		revertReason, err := eth.ExtractRevertReasonFromRPCError(wrappedErr)
		require.NoError(t, err)
		require.Equal(t, message, revertReason)
	})

	t.Run("it gracefully errors when no data present", func(tt *testing.T) {
		var jsonErr error = &eth.JsonError{
			Code:    1,
			Message: "something different",
		}
		_, err := eth.ExtractRevertReasonFromRPCError(jsonErr)
		require.Error(t, err)
	})

	t.Run("gracefully errors when given a normal error", func(tt *testing.T) {
		_, err := eth.ExtractRevertReasonFromRPCError(errors.New("normal error"))
		require.Error(tt, err)
	})

	t.Run("gracefully errors when given no error", func(tt *testing.T) {
		_, err := eth.ExtractRevertReasonFromRPCError(nil)
		require.Error(tt, err)
	})
}
