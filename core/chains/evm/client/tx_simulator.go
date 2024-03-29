package client

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rpc"

	"github.com/smartcontractkit/chainlink/v2/common/config"
)

type TxSimulationRequest struct {
	From  common.Address
	To    *common.Address
	Data  []byte
	Error *SendError
}

type simulatorClient interface {
	CallContext(ctx context.Context, result interface{}, method string, args ...interface{}) error
	BatchCallContext(ctx context.Context, b []rpc.BatchElem) error
}

// ZK chains can return an out-of-counters error
// This method allows a caller to determine if a tx would fail due to OOC error by simulating the transaction
// Used as an entry point in case custom simulation is required across different chains
func SimulateTransaction(ctx context.Context, client simulatorClient, chainType config.ChainType, req TxSimulationRequest) *SendError {
	reqs := []TxSimulationRequest{req}
	err := simulateTransactionDefault(ctx, client, reqs)
	if err != nil {
		return NewSendError(err)
	}
	return reqs[0].Error
}

func BatchSimulateTransaction(ctx context.Context, client simulatorClient, chainType config.ChainType, reqs []TxSimulationRequest) error {
	return simulateTransactionDefault(ctx, client, reqs)
}

// eth_estimateGas returns out-of-counters (OOC) error if the transaction would result in an overflow
func simulateTransactionDefault(ctx context.Context, client simulatorClient, reqs []TxSimulationRequest) error {
	rpcBatchCalls := make([]rpc.BatchElem, len(reqs))
	for i, req := range reqs {
		toAddress := common.Address{}
		if req.To != nil {
			toAddress = *req.To
		}
		rpcBatchCalls[i] = rpc.BatchElem{
			Method: "eth_estimateGas",
			Args: []any{
				map[string]interface{}{
					"from": req.From,
					"to":   toAddress,
					"data": hexutil.Bytes(req.Data),
				},
				"pending",
			},
			Result: new(hexutil.Big),
		}
	}
	err := client.BatchCallContext(ctx, rpcBatchCalls)
	if err != nil {
		return err
	}

	for _, elem := range rpcBatchCalls {
		params := elem.Args[0].(map[string]interface{})
		for i := 0; i < len(reqs); i++ {
			req := &reqs[i]
			// Match the request to rpc response to set the error in the proper object
			if params["from"] == req.From &&
				(req.To == nil && params["to"] == common.Address{} || req.To != nil && params["to"] == *req.To) &&
				params["data"].(hexutil.Bytes).String() == hexutil.Bytes(req.Data).String() {
				// Wrap RPC error in SendError
				req.Error = NewSendError(elem.Error)
			}
		}
	}
	return nil
}
