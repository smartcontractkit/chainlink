package client_test

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"

	evmclient "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
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
			{"nonce too low: address 0x336394A3219e71D9d9bd18201d34E95C1Bb7122C, tx: 8089 state: 8090", true, "Arbitrum"},
			{"Nonce too low", true, "Besu"},
			{"nonce too low", true, "Erigon"},
			{"nonce too low", true, "Klaytn"},
			{"Transaction nonce is too low. Try incrementing the nonce.", true, "Parity"},
			{"transaction rejected: nonce too low", true, "Arbitrum"},
			{"invalid transaction nonce", true, "Arbitrum"},
			{"invalid transaction: nonce too low", true, "Optimism"},
			{"call failed: nonce too low: address 0x0499BEA33347cb62D79A9C0b1EDA01d8d329894c current nonce (5833) > tx nonce (5511)", true, "Avalanche"},
			{"call failed: OldNonce", true, "Nethermind"},
			{"call failed: OldNonce, Current nonce: 22, nonce of rejected tx: 17", true, "Nethermind"},
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

	t.Run("IsNonceTooHigh", func(t *testing.T) {

		tests := []errorCase{
			{"call failed: NonceGap", true, "Nethermind"},
			{"call failed: NonceGap, Future nonce. Expected nonce: 10", true, "Nethermind"},
			{"nonce too high: address 0x336394A3219e71D9d9bd18201d34E95C1Bb7122C, tx: 8089 state: 8090", true, "Arbitrum"},
			{"nonce too high", true, "Geth"},
			{"nonce too high", true, "Erigon"},
		}

		for _, test := range tests {
			err = evmclient.NewSendErrorS(test.message)
			assert.Equal(t, err.IsNonceTooHighError(), test.expect)
			err = newSendErrorWrapped(test.message)
			assert.Equal(t, err.IsNonceTooHighError(), test.expect)
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
			{"Replacement transaction underpriced", true, "Besu"},
			{"replacement transaction underpriced", true, "Erigon"},
			{"replacement transaction underpriced", true, "Klaytn"},
			{"there is another tx which has the same nonce in the tx pool", true, "Klaytn"},
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
			{"Known transaction", true, "Besu"},
			{"already known", true, "Erigon"},
			{"block already known", true, "Erigon"},
			{"Transaction with the same hash was already imported.", true, "Parity"},
			{"call failed: AlreadyKnown", true, "Nethermind"},
			{"call failed: OwnNonceAlreadyUsed", true, "Nethermind"},
			{"known transaction", true, "Klaytn"},
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
			{"Gas price below configured minimum gas price", true, "Besu"},
			{"transaction underpriced", true, "Erigon"},
			{"There are too many transactions in the queue. Your transaction was dropped due to limit. Try increasing the fee.", false, "Parity"},
			{"Transaction gas price is too low. It does not satisfy your node's minimal gas price (minimal: 100 got: 50). Try increasing the gas price.", true, "Parity"},
			{"gas price too low", true, "Arbitrum"},
			{"FeeTooLow", true, "Nethermind"},
			{"FeeTooLow, MaxFeePerGas too low. MaxFeePerGas: 50, BaseFee: 100, MaxPriorityFeePerGas:200, Block number: 5", true, "Nethermind"},
			{"FeeTooLow, EffectivePriorityFeePerGas too low 10 < 20, BaseFee: 30", true, "Nethermind"},
			{"FeeTooLow, FeePerGas needs to be higher than 100 to be added to the TxPool. Affordable FeePerGas of rejected tx: 50.", true, "Nethermind"},
			{"FeeTooLowToCompete", true, "Nethermind"},
			{"transaction underpriced", true, "Klaytn"},
			{"intrinsic gas too low", true, "Klaytn"},
		}

		for _, test := range tests {
			err = evmclient.NewSendErrorS(test.message)
			assert.Equal(t, err.IsTerminallyUnderpriced(), test.expect, "expected %q to match %s for client %s", err, "IsTerminallyUnderpriced", test.network)
			err = newSendErrorWrapped(test.message)
			assert.Equal(t, err.IsTerminallyUnderpriced(), test.expect, "expected %q to match %s for client %s", err, "IsTerminallyUnderpriced", test.network)
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
			{"Upfront cost exceeds account balance", true, "Besu"},
			{"insufficient funds for transfer", true, "Erigon"},
			{"insufficient funds for gas * price + value", true, "Erigon"},
			{"insufficient balance for transfer", true, "Erigon"},
			{"Insufficient balance for transaction. Balance=100.25, Cost=200.50", true, "Parity"},
			{"Insufficient funds. The account you tried to send transaction from does not have enough funds. Required 200.50 and got: 100.25.", true, "Parity"},
			{"transaction rejected: insufficient funds for gas * price + value", true, "Arbitrum"},
			{"not enough funds for gas", true, "Arbitrum"},
			{"insufficient funds for gas * price + value: address 0xb68D832c1241bc50db1CF09e96c0F4201D5539C9 have 9934612900000000 want 9936662900000000", true, "Arbitrum"},
			{"invalid transaction: insufficient funds for gas * price + value", true, "Optimism"},
			{"call failed: InsufficientFunds", true, "Nethermind"},
			{"call failed: InsufficientFunds, Account balance: 4740799397601480913, cumulative cost: 22019342038993800000", true, "Nethermind"},
			{"insufficient funds", true, "Klaytn"},
		}
		for _, test := range tests {
			err = evmclient.NewSendErrorS(test.message)
			assert.Equal(t, err.IsInsufficientEth(), test.expect)
			err = newSendErrorWrapped(test.message)
			assert.Equal(t, err.IsInsufficientEth(), test.expect)
		}
	})

	t.Run("IsTxFeeExceedsCap", func(t *testing.T) {
		tests := []errorCase{
			{"tx fee (1.10 ether) exceeds the configured cap (1.00 ether)", true, "geth"},
			{"tx fee (1.10 FTM) exceeds the configured cap (1.00 FTM)", true, "geth"},
			{"tx fee (1.10 foocoin) exceeds the configured cap (1.00 foocoin)", true, "geth"},
			{"Transaction fee cap exceeded", true, "Besu"},
			{"tx fee (1.10 ether) exceeds the configured cap (1.00 ether)", true, "Erigon"},
			{"invalid gas fee cap", true, "Klaytn"},
			{"max fee per gas higher than max priority fee per gas", true, "Klaytn"},
		}
		for _, test := range tests {
			err = evmclient.NewSendErrorS(test.message)
			assert.Equal(t, err.IsTxFeeExceedsCap(), test.expect)
			err = newSendErrorWrapped(test.message)
			assert.Equal(t, err.IsTxFeeExceedsCap(), test.expect)
		}

		assert.False(t, randomError.IsTxFeeExceedsCap())
		// Nil
		err = evmclient.NewSendError(nil)
		assert.False(t, err.IsTxFeeExceedsCap())
	})

	t.Run("L2 Fees errors", func(t *testing.T) {
		err := evmclient.NewSendErrorS("primary websocket (wss://ws-mainnet.optimism.io) call failed: fee too high: 5835750750000000, use less than 467550750000000 * 0.700000")
		assert.True(t, err.IsL2FeeTooHigh())
		assert.False(t, err.L2FeeTooLow())
		err = newSendErrorWrapped("primary websocket (wss://ws-mainnet.optimism.io) call failed: fee too high: 5835750750000000, use less than 467550750000000 * 0.700000")
		assert.True(t, err.IsL2FeeTooHigh())
		assert.False(t, err.L2FeeTooLow())

		err = evmclient.NewSendErrorS("fee too low: 30365610000000, use at least tx.gasLimit = 5874374 and tx.gasPrice = 15000000")
		assert.False(t, err.IsL2FeeTooHigh())
		assert.True(t, err.L2FeeTooLow())
		err = newSendErrorWrapped("fee too low: 30365610000000, use at least tx.gasLimit = 5874374 and tx.gasPrice = 15000000")
		assert.False(t, err.IsL2FeeTooHigh())
		assert.True(t, err.L2FeeTooLow())

		err = evmclient.NewSendErrorS("queue full")
		assert.True(t, err.IsL2Full())
		err = evmclient.NewSendErrorS("sequencer pending tx pool full, please try again")
		assert.True(t, err.IsL2Full())

		assert.False(t, randomError.IsL2FeeTooHigh())
		assert.False(t, randomError.L2FeeTooLow())
		// Nil
		err = evmclient.NewSendError(nil)
		assert.False(t, err.IsL2FeeTooHigh())
		assert.False(t, err.L2FeeTooLow())
	})

	t.Run("Metis gas price errors", func(t *testing.T) {
		err := evmclient.NewSendErrorS("primary websocket (wss://ws-mainnet.optimism.io) call failed: gas price too low: 18000000000 wei, use at least tx.gasPrice = 19500000000 wei")
		assert.True(t, err.L2FeeTooLow())
		err = newSendErrorWrapped("primary websocket (wss://ws-mainnet.optimism.io) call failed: gas price too low: 18000000000 wei, use at least tx.gasPrice = 19500000000 wei")
		assert.True(t, err.L2FeeTooLow())

		assert.False(t, randomError.L2FeeTooLow())
		// Nil
		err = evmclient.NewSendError(nil)
		assert.False(t, err.L2FeeTooLow())
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

		{"Intrinsic gas exceeds gas limit", true, "Besu"},
		{"Transaction gas limit exceeds block gas limit", true, "Besu"},
		{"Invalid signature", true, "Besu"},

		{"insufficient funds for transfer", false, "Erigon"},
		{"exceeds block gas limit", true, "Erigon"},
		{"invalid sender", true, "Erigon"},
		{"negative value", true, "Erigon"},
		{"oversized data", true, "Erigon"},
		{"gas uint64 overflow", true, "Erigon"},
		{"intrinsic gas too low", true, "Erigon"},

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
		{"execution reverted: stale report", true, "Arbitrum"},
		{"execution reverted", true, "Arbitrum"},

		{"call failed: SenderIsContract", true, "Nethermind"},
		{"call failed: Invalid", true, "Nethermind"},
		{"call failed: Invalid, transaction Hash is null", true, "Nethermind"},
		{"call failed: Int256Overflow", true, "Nethermind"},
		{"call failed: FailedToResolveSender", true, "Nethermind"},
		{"call failed: GasLimitExceeded", true, "Nethermind"},
		{"call failed: GasLimitExceeded, Gas limit: 100, gas limit of rejected tx: 150", true, "Nethermind"},

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
