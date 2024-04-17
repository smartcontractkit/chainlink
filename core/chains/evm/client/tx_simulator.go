package client

import (
	"context"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink/v2/common/config"
)

type simulatorClient interface {
	CallContext(ctx context.Context, result interface{}, method string, args ...interface{}) error
}

// ZK chains can return an out-of-counters error
// This method allows a caller to determine if a tx would fail due to OOC error by simulating the transaction
// Used as an entry point in case custom simulation is required across different chains
func SimulateTransaction(ctx context.Context, client simulatorClient, lggr logger.SugaredLogger, chainType config.ChainType, msg ethereum.CallMsg) *SendError {
	err := simulateTransactionDefault(ctx, client, msg)
	return NewSendError(err)
}

// eth_estimateGas returns out-of-counters (OOC) error if the transaction would result in an overflow
func simulateTransactionDefault(ctx context.Context, client simulatorClient, msg ethereum.CallMsg) error {
	var result hexutil.Big
	return client.CallContext(ctx, &result, "eth_estimateGas", toCallArg(msg), "pending")
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
