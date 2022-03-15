package client_test

import (
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	evmclient "github.com/smartcontractkit/chainlink/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/core/utils"
)

func newSendErrorWrapped(s string) *evmclient.SendError {
	return evmclient.NewSendError(errors.Wrap(errors.New(s), "wrapped with some old bollocks"))
}

type errorCase struct {
	message string
	expect  bool
	network string
}

func Test_Eth_Errors(t *testing.T) {
	t.Parallel()
	var err *evmclient.SendError
	randomError := evmclient.NewSendErrorS("some old bollocks")

	t.Run("IsNonceTooLowError", func(t *testing.T) {
		assert.False(t, randomError.IsNonceTooLowError())

		tests := []errorCase{
			{"nonce too low", true, "Geth"},
			{"Transaction nonce is too low. Try incrementing the nonce.", true, "Parity"},
			{"transaction rejected: nonce too low", true, "Arbitrum"},
			{"invalid transaction nonce", true, "Arbitrum"},
			{"invalid transaction: nonce too low", true, "Optimism"},
			{"call failed: nonce too low: address 0x0499BEA33347cb62D79A9C0b1EDA01d8d329894c current nonce (5833) > tx nonce (5511)", true, "Avalanche"},
			{"call failed: OldNonce", true, "Nethermind"},
		}

		for _, test := range tests {
			t.Run(test.network, func(t *testing.T) {
				err = evmclient.NewSendErrorS(test.message)
				assert.Equal(t, err.IsNonceTooLowError(), test.expect)
				err = newSendErrorWrapped(test.message)
				assert.Equal(t, err.IsNonceTooLowError(), test.expect)
			})
		}
	})

	t.Run("IsTransactionAlreadyMined", func(t *testing.T) {
		assert.False(t, randomError.IsTransactionAlreadyMined())

		tests := []errorCase{
			{"transaction already finalized", true, "Harmony"},
		}

		for _, test := range tests {
			t.Run(test.network, func(t *testing.T) {
				err = evmclient.NewSendErrorS(test.message)
				assert.Equal(t, err.IsTransactionAlreadyMined(), test.expect)
				err = newSendErrorWrapped(test.message)
				assert.Equal(t, err.IsTransactionAlreadyMined(), test.expect)
			})
		}
	})

	t.Run("IsReplacementUnderpriced", func(t *testing.T) {

		tests := []errorCase{
			{"replacement transaction underpriced", true, "geth"},
			{"Transaction gas price 100wei is too low. There is another transaction with same nonce in the queue with gas price 150wei. Try increasing the gas price or incrementing the nonce.", true, "Parity"},
			{"There are too many transactions in the queue. Your transaction was dropped due to limit. Try increasing the fee.", false, "Parity"},
			{"gas price too low", false, "Arbitrum"},
		}

		for _, test := range tests {
			err = evmclient.NewSendErrorS(test.message)
			assert.Equal(t, err.IsReplacementUnderpriced(), test.expect)
			err = newSendErrorWrapped(test.message)
			assert.Equal(t, err.IsReplacementUnderpriced(), test.expect)
		}
	})

	t.Run("IsTransactionAlreadyInMempool", func(t *testing.T) {
		assert.False(t, randomError.IsTransactionAlreadyInMempool())

		tests := []errorCase{
			// I have seen this in log output
			{"known transaction: 0x7f657507aee0511e36d2d1972a6b22e917cc89f92b6c12c4dbd57eaabb236960", true, "Geth"},
			// This comes from the geth source - https://github.com/ethereum/go-ethereum/blob/eb9d7d15ecf08cd5104e01a8af64489f01f700b0/core/tx_pool.go#L57
			{"already known", true, "Geth"},
			// This one is present in the light client (?!)
			{"Known transaction (7f65)", true, "Geth"},
			{"Transaction with the same hash was already imported.", true, "Parity"},
			{"call failed: AlreadyKnown", true, "Nethermind"},
			{"call failed: OwnNonceAlreadyUsed", true, "Nethermind"},
		}
		for _, test := range tests {
			err = evmclient.NewSendErrorS(test.message)
			assert.Equal(t, err.IsTransactionAlreadyInMempool(), test.expect)
			err = newSendErrorWrapped(test.message)
			assert.Equal(t, err.IsTransactionAlreadyInMempool(), test.expect)
		}
	})

	t.Run("IsTerminallyUnderpriced", func(t *testing.T) {
		assert.False(t, randomError.IsTerminallyUnderpriced())

		tests := []errorCase{
			{"transaction underpriced", true, "geth"},
			{"replacement transaction underpriced", false, "geth"},
			{"There are too many transactions in the queue. Your transaction was dropped due to limit. Try increasing the fee.", false, "Parity"},
			{"Transaction gas price is too low. It does not satisfy your node's minimal gas price (minimal: 100 got: 50). Try increasing the gas price.", true, "Parity"},
			{"gas price too low", true, "Arbitrum"},
		}

		for _, test := range tests {
			err = evmclient.NewSendErrorS(test.message)
			assert.Equal(t, err.IsTerminallyUnderpriced(), test.expect)
			err = newSendErrorWrapped(test.message)
			assert.Equal(t, err.IsTerminallyUnderpriced(), test.expect)
		}
	})

	t.Run("IsTemporarilyUnderpriced", func(t *testing.T) {
		tests := []errorCase{
			{"There are too many transactions in the queue. Your transaction was dropped due to limit. Try increasing the fee.", true, "Parity"},
			{"There are too many transactions in the queue. Your transaction was dropped due to limit. Try increasing the fee.", true, "Parity"},
			{"Transaction gas price is too low. It does not satisfy your node's minimal gas price (minimal: 100 got: 50). Try increasing the gas price.", false, "Parity"},
		}
		for _, test := range tests {
			err = evmclient.NewSendErrorS(test.message)
			assert.Equal(t, err.IsTemporarilyUnderpriced(), test.expect)
			err = newSendErrorWrapped(test.message)
			assert.Equal(t, err.IsTemporarilyUnderpriced(), test.expect)
		}
	})

	t.Run("IsInsufficientEth", func(t *testing.T) {
		tests := []errorCase{
			{"insufficient funds for transfer", true, "Geth"},
			{"insufficient funds for gas * price + value", true, "Geth"},
			{"insufficient balance for transfer", true, "Geth"},
			{"Insufficient balance for transaction. Balance=100.25, Cost=200.50", true, "Parity"},
			{"Insufficient funds. The account you tried to send transaction from does not have enough funds. Required 200.50 and got: 100.25.", true, "Parity"},
			{"transaction rejected: insufficient funds for gas * price + value", true, "Arbitrum"},
			{"not enough funds for gas", true, "Arbitrum"},
			{"invalid transaction: insufficient funds for gas * price + value", true, "Optimism"},
			{"call failed: InsufficientFunds", true, "Nethermind"},
		}
		for _, test := range tests {
			err = evmclient.NewSendErrorS(test.message)
			assert.Equal(t, err.IsInsufficientEth(), test.expect)
			err = newSendErrorWrapped(test.message)
			assert.Equal(t, err.IsInsufficientEth(), test.expect)
		}
	})

	t.Run("IsTooExpensive", func(t *testing.T) {
		tests := []errorCase{
			{"tx fee (1.10 ether) exceeds the configured cap (1.00 ether)", true, "geth"},
			{"call failed: InsufficientFunds", true, "Nethermind"},
		}
		for _, test := range tests {
			err = evmclient.NewSendErrorS(test.message)
			assert.Equal(t, err.IsTooExpensive(), test.expect)
			err = newSendErrorWrapped(test.message)
			assert.Equal(t, err.IsTooExpensive(), test.expect)
		}

		assert.False(t, randomError.IsTooExpensive())
		// Nil
		err = evmclient.NewSendError(nil)
		assert.False(t, err.IsTooExpensive())
	})

	t.Run("Optimism Fees errors", func(t *testing.T) {
		err := evmclient.NewSendErrorS("primary websocket (wss://ws-mainnet.optimism.io) call failed: fee too high: 5835750750000000, use less than 467550750000000 * 0.700000")
		assert.True(t, err.IsFeeTooHigh())
		assert.False(t, err.IsFeeTooLow())
		err = newSendErrorWrapped("primary websocket (wss://ws-mainnet.optimism.io) call failed: fee too high: 5835750750000000, use less than 467550750000000 * 0.700000")
		assert.True(t, err.IsFeeTooHigh())
		assert.False(t, err.IsFeeTooLow())

		err = evmclient.NewSendErrorS("fee too low: 30365610000000, use at least tx.gasLimit = 5874374 and tx.gasPrice = 15000000")
		assert.False(t, err.IsFeeTooHigh())
		assert.True(t, err.IsFeeTooLow())
		err = newSendErrorWrapped("fee too low: 30365610000000, use at least tx.gasLimit = 5874374 and tx.gasPrice = 15000000")
		assert.False(t, err.IsFeeTooHigh())
		assert.True(t, err.IsFeeTooLow())

		assert.False(t, randomError.IsFeeTooHigh())
		assert.False(t, randomError.IsFeeTooLow())
		// Nil
		err = evmclient.NewSendError(nil)
		assert.False(t, err.IsFeeTooHigh())
		assert.False(t, err.IsFeeTooLow())
	})

	t.Run("moonriver errors", func(t *testing.T) {
		err := evmclient.NewSendErrorS("primary http (http://***REDACTED***:9933) call failed: submit transaction to pool failed: Pool(Stale)")
		assert.True(t, err.IsNonceTooLowError())
		assert.False(t, err.IsTransactionAlreadyInMempool())
		assert.False(t, err.Fatal())
		err = evmclient.NewSendErrorS("primary http (http://***REDACTED***:9933) call failed: submit transaction to pool failed: Pool(AlreadyImported)")
		assert.True(t, err.IsTransactionAlreadyInMempool())
		assert.False(t, err.IsNonceTooLowError())
		assert.False(t, err.Fatal())
	})
}

func Test_Eth_Errors_Fatal(t *testing.T) {
	t.Parallel()

	tests := []errorCase{
		{"some old bollocks", false, "none"},

		{"insufficient funds for transfer", false, "Geth"},
		{"exceeds block gas limit", true, "Geth"},
		{"invalid sender", true, "Geth"},
		{"negative value", true, "Geth"},
		{"oversized data", true, "Geth"},
		{"gas uint64 overflow", true, "Geth"},
		{"intrinsic gas too low", true, "Geth"},
		{"nonce too high", true, "Geth"},

		{"Insufficient funds. The account you tried to send transaction from does not have enough funds. Required 100 and got: 50.", false, "Parity"},
		{"Supplied gas is beyond limit.", true, "Parity"},
		{"Sender is banned in local queue.", true, "Parity"},
		{"Recipient is banned in local queue.", true, "Parity"},
		{"Code is banned in local queue.", true, "Parity"},
		{"Transaction is not permitted.", true, "Parity"},
		{"Transaction is too big, see chain specification for the limit.", true, "Parity"},
		{"Transaction gas is too low. There is not enough gas to cover minimal cost of the transaction (minimal: 100 got: 50) Try increasing supplied gas.", true, "Parity"},
		{"Transaction cost exceeds current gas limit. Limit: 50, got: 100. Try decreasing supplied gas.", true, "Parity"},
		{"Invalid signature: some old bollocks", true, "Parity"},
		{"Invalid RLP data: some old bollocks", true, "Parity"},

		{"invalid message format", true, "Arbitrum"},
		{"forbidden sender address", true, "Arbitrum"},
		{"tx dropped due to L2 congestion", false, "Arbitrum"},
		{"execution reverted: error code", true, "Arbitrum"},

		{"call failed: SenderIsContract", true, "Nethermind"},
		{"call failed: Invalid", true, "Nethermind"},
		{"call failed: Int256Overflow", true, "Nethermind"},
		{"call failed: FailedToResolveSender", true, "Nethermind"},
		{"call failed: GasLimitExceeded", true, "Nethermind"},

		{"invalid shard", true, "Harmony"},
		{"`to` address of transaction in blacklist", true, "Harmony"},
		{"`from` address of transaction in blacklist", true, "Harmony"},
		{"staking message does not match directive message", true, "Harmony"},
	}

	for _, test := range tests {
		t.Run(test.message, func(t *testing.T) {
			err := evmclient.NewSendError(errors.New(test.message))
			assert.Equal(t, test.expect, err.Fatal())
		})
	}
}

func Test_ExtractRevertReasonFromRPCError(t *testing.T) {
	message := "important revert reason"
	messageHex := utils.RemoveHexPrefix(hexutil.Encode([]byte(message)))
	sigHash := "12345678"
	var jsonErr error = &evmclient.JsonError{
		Code:    1,
		Data:    fmt.Sprintf("0x%s%s", sigHash, messageHex),
		Message: "something different",
	}

	t.Run("it extracts revert reasons when present", func(tt *testing.T) {
		revertReason, err := evmclient.ExtractRevertReasonFromRPCError(jsonErr)
		require.NoError(t, err)
		require.Equal(t, message, revertReason)
	})

	t.Run("it unwraps wrapped errors", func(tt *testing.T) {
		wrappedErr := errors.Wrap(jsonErr, "wrapped message")
		revertReason, err := evmclient.ExtractRevertReasonFromRPCError(wrappedErr)
		require.NoError(t, err)
		require.Equal(t, message, revertReason)
	})

	t.Run("it unwraps multi-wrapped errors", func(tt *testing.T) {
		wrappedErr := errors.Wrap(jsonErr, "wrapped message")
		wrappedErr = errors.Wrap(wrappedErr, "wrapped again!!")
		revertReason, err := evmclient.ExtractRevertReasonFromRPCError(wrappedErr)
		require.NoError(t, err)
		require.Equal(t, message, revertReason)
	})

	t.Run("it gracefully errors when no data present", func(tt *testing.T) {
		var jsonErr error = &evmclient.JsonError{
			Code:    1,
			Message: "something different",
		}
		_, err := evmclient.ExtractRevertReasonFromRPCError(jsonErr)
		require.Error(t, err)
	})

	t.Run("gracefully errors when given a normal error", func(tt *testing.T) {
		_, err := evmclient.ExtractRevertReasonFromRPCError(errors.New("normal error"))
		require.Error(tt, err)
	})

	t.Run("gracefully errors when given no error", func(tt *testing.T) {
		_, err := evmclient.ExtractRevertReasonFromRPCError(nil)
		require.Error(tt, err)
	})
}
