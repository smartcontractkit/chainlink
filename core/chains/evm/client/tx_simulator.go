package client

import (
	"context"
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink/v2/common/config"

	commonclient "github.com/smartcontractkit/chainlink/v2/common/client"
)

const ErrOutOfCounters = "not enough counters to continue the execution"

type simulatorClient interface {
	CallContext(ctx context.Context, result interface{}, method string, args ...interface{}) error
}

// ZK chains can return an out-of-counters error
// This method allows a caller to determine if a tx would fail due to OOC error by simulating the transaction
// Used as an entry point for custom simulation across different chains
func SimulateTransaction(ctx context.Context, client simulatorClient, lggr logger.SugaredLogger, chainType config.ChainType, msg ethereum.CallMsg) error {
	var err error
	switch chainType {
	case config.ChainZkEvm:
		err = simulateTransactionZkEvm(ctx, client, lggr, msg)
	default:
		err = simulateTransactionDefault(ctx, client, msg)
	}
	// ClassifySendError will not have the proper fields for logging within the method due to the empty Transaction passed
	code := ClassifySendError(err, lggr, &types.Transaction{}, msg.From, chainType.IsL2())
	// Only return error if ZK OOC error is identified
	if code == commonclient.OutOfCounters {
		return errors.New(ErrOutOfCounters)
	}
	return nil
}

// eth_estimateGas returns out-of-counters (OOC) error if the transaction would result in an overflow
func simulateTransactionDefault(ctx context.Context, client simulatorClient, msg ethereum.CallMsg) error {
	var result hexutil.Big
	errCall := client.CallContext(ctx, &result, "eth_estimateGas", toCallArg(msg), "pending")
	jsonErr, _ := ExtractRPCError(errCall)
	if jsonErr != nil && len(jsonErr.Message) > 0 {
		return errors.New(jsonErr.Message)
	}
	return nil
}

type zkEvmEstimateCountResponse struct {
	CountersUsed struct {
		GasUsed              string
		UsedKeccakHashes     string
		UsedPoseidonHashes   string
		UsedPoseidonPaddings string
		UsedMemAligns        string
		UsedArithmetics      string
		UsedBinaries         string
		UsedSteps            string
		UsedSHA256Hashes     string
	}
	CountersLimit struct {
		MaxGasUsed          string
		MaxKeccakHashes     string
		MaxPoseidonHashes   string
		MaxPoseidonPaddings string
		MaxMemAligns        string
		MaxArithmetics      string
		MaxBinaries         string
		MaxSteps            string
		MaxSHA256Hashes     string
	}
	OocError string
}

// zkEVM implemented a custom zkevm_estimateCounters method to detect if a transaction would result in an out-of-counters (OOC) error
func simulateTransactionZkEvm(ctx context.Context, client simulatorClient, lggr logger.SugaredLogger, msg ethereum.CallMsg) error {
	var result zkEvmEstimateCountResponse
	err := client.CallContext(ctx, &result, "zkevm_estimateCounters", toCallArg(msg), "pending")
	if err != nil {
		return fmt.Errorf("failed to simulate tx: %w", err)
	}
	if detectZkEvmCounterOverflow(result) && len(result.OocError) > 0 {
		lggr.Debugw("zkevm_estimateCounters returned", "result", result)
		return errors.New(result.OocError)
	}
	return nil
}

// Helper method for zkEvm to determine if response indicates an overflow
func detectZkEvmCounterOverflow(result zkEvmEstimateCountResponse) bool {
	if result.CountersUsed.UsedKeccakHashes > result.CountersLimit.MaxKeccakHashes ||
		result.CountersUsed.UsedPoseidonHashes > result.CountersLimit.MaxPoseidonHashes ||
		result.CountersUsed.UsedPoseidonPaddings > result.CountersLimit.MaxPoseidonPaddings ||
		result.CountersUsed.UsedMemAligns > result.CountersLimit.MaxMemAligns ||
		result.CountersUsed.UsedArithmetics > result.CountersLimit.MaxArithmetics ||
		result.CountersUsed.UsedBinaries > result.CountersLimit.MaxBinaries ||
		result.CountersUsed.UsedSteps > result.CountersLimit.MaxSteps ||
		result.CountersUsed.UsedSHA256Hashes > result.CountersLimit.MaxSHA256Hashes {
		return true
	}
	return false
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
