package txmgr

import (
	"context"
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/smartcontractkit/chainlink/v2/common/config"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
)

const ErrOutOfCounters = "not enough keccak counters to continue the execution"

// ZK Chain can return an overflow error based on the number of keccak hashes in the call
// This method allows a caller to determine if a tx would fail due to overflow error by simulating the transaction
// Used as an entry point for custom simulation across different chains
func SimulateTransaction(ctx context.Context, c client.Client, chainType config.ChainType, msg ethereum.CallMsg) error {
	switch chainType {
	case config.ChainZkEvm:
		return simulateTransactionZkEvm(ctx, c, msg)
	default:
		return simulateTransactionDefault(ctx, c, msg)
	}
}

// eth_estimateGas returns out-of-counters (OOC) error if the transaction would result in an overflow
func simulateTransactionDefault(ctx context.Context, c client.Client, msg ethereum.CallMsg) error {
	var result uint64
	errCall := c.CallContext(ctx, &result, "eth_estimateGas", toCallArg(msg), "pending")
	jsonErr, err := client.ExtractRPCError(errCall)
	if err != nil {
		return fmt.Errorf("failed to simulate tx: %w", err)
	}
	// Only return error if Zk OOC error is identified
	if jsonErr.Message == ErrOutOfCounters {
		return errors.New(ErrOutOfCounters)
	}
	return nil
}

// zkEVM implemented a custom zkevm_estimateCounters method to detect if a transaction would result in an out-of-counters (OOC) error
func simulateTransactionZkEvm(ctx context.Context, c client.Client, msg ethereum.CallMsg) error {
	var result struct {
		countersUsed struct {
			gasUsed              int
			usedKeccakHashes     int
			usedPoseidonHashes   int
			usedPoseidonPaddings int
			usedMemAligns        int
			usedArithmetics      int
			usedBinaries         int
			usedSteps            int
			usedSHA256Hashes     int
		}
		countersLimit struct {
			maxGasUsed          int
			maxKeccakHashes     int
			maxPoseidonHashes   int
			maxPoseidonPaddings int
			maxMemAligns        int
			maxArithmetics      int
			maxBinaries         int
			maxSteps            int
			maxSHA256Hashes     int
		}
		oocError string
	}
	err := c.CallContext(ctx, &result, "zkevm_estimateCounters", toCallArg(msg), "pending")
	if err != nil {
		return fmt.Errorf("failed to simulate tx: %w", err)
	}
	if len(result.oocError) > 0 {
		return errors.New(result.oocError)
	}
	return nil
}

func toCallArg(msg ethereum.CallMsg) interface{} {
	arg := map[string]interface{}{
		"from": msg.From,
		"to":   msg.To,
	}
	if len(msg.Data) > 0 {
		arg["input"] = hexutil.Bytes(msg.Data)
	}
	if msg.Value != nil {
		arg["value"] = (*hexutil.Big)(msg.Value)
	}
	if msg.Gas != 0 {
		arg["gas"] = hexutil.Uint64(msg.Gas)
	}
	if msg.GasPrice != nil {
		arg["gasPrice"] = (*hexutil.Big)(msg.GasPrice)
	}
	if msg.GasFeeCap != nil {
		arg["maxFeePerGas"] = (*hexutil.Big)(msg.GasFeeCap)
	}
	if msg.GasTipCap != nil {
		arg["maxPriorityFeePerGas"] = (*hexutil.Big)(msg.GasTipCap)
	}
	return arg
}
