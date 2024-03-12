package wrappers

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/smartcontractkit/seth"

	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
)

// It needs to be kept up do date with core/chains/evm/client/rpc_client.go !!!!!

// WrappedContractBackend is a wrapper around the go-ethereum ContractBackend interface. It's a thin wrapper
// around the go-ethereum/ethclient.Client, which replaces only CallContract and PendingCallContract calls with
// methods that send data both in "input" and "data" field for backwards compatibility with older clients. Other methods
// are passed through to the underlying client.
type WrappedContractBackend struct {
	evmClient  blockchain.EVMClient
	sethClient *seth.Client
}

// MustNewWrappedContractBackend creates a new WrappedContractBackend with the given clients
func MustNewWrappedContractBackend(evmClient blockchain.EVMClient, sethClient *seth.Client) *WrappedContractBackend {
	if evmClient == nil && sethClient == nil {
		panic("Must provide at least one client")
	}

	return &WrappedContractBackend{
		evmClient:  evmClient,
		sethClient: sethClient,
	}
}

func (w *WrappedContractBackend) getGethClient() *ethclient.Client {
	if w.sethClient != nil {
		return w.sethClient.Client
	}

	if w.evmClient != nil {
		return w.evmClient.GetEthClient()
	}

	panic("No client found")
}

func (w *WrappedContractBackend) CodeAt(ctx context.Context, contract common.Address, blockNumber *big.Int) ([]byte, error) {
	client := w.getGethClient()
	return client.CodeAt(ctx, contract, blockNumber)
}

func (w *WrappedContractBackend) PendingCodeAt(ctx context.Context, contract common.Address) ([]byte, error) {
	client := w.getGethClient()
	return client.PendingCodeAt(ctx, contract)
}

func (w *WrappedContractBackend) CodeAtHash(ctx context.Context, contract common.Address, blockHash common.Hash) ([]byte, error) {
	client := w.getGethClient()
	return client.CodeAtHash(ctx, contract, blockHash)
}

func (w *WrappedContractBackend) CallContractAtHash(ctx context.Context, call ethereum.CallMsg, blockHash common.Hash) ([]byte, error) {
	client := w.getGethClient()
	return client.CallContractAtHash(ctx, call, blockHash)
}

func (w *WrappedContractBackend) HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error) {
	client := w.getGethClient()
	return client.HeaderByNumber(ctx, number)
}

func (w *WrappedContractBackend) PendingNonceAt(ctx context.Context, account common.Address) (uint64, error) {
	client := w.getGethClient()
	return client.PendingNonceAt(ctx, account)
}

func (w *WrappedContractBackend) SuggestGasPrice(ctx context.Context) (*big.Int, error) {
	client := w.getGethClient()
	return client.SuggestGasPrice(ctx)
}

func (w *WrappedContractBackend) SuggestGasTipCap(ctx context.Context) (*big.Int, error) {
	client := w.getGethClient()
	return client.SuggestGasTipCap(ctx)
}

func (w *WrappedContractBackend) EstimateGas(ctx context.Context, call ethereum.CallMsg) (gas uint64, err error) {
	client := w.getGethClient()
	return client.EstimateGas(ctx, call)
}

func (w *WrappedContractBackend) SendTransaction(ctx context.Context, tx *types.Transaction) error {
	client := w.getGethClient()
	return client.SendTransaction(ctx, tx)
}

func (w *WrappedContractBackend) FilterLogs(ctx context.Context, query ethereum.FilterQuery) ([]types.Log, error) {
	client := w.getGethClient()
	return client.FilterLogs(ctx, query)
}

func (w *WrappedContractBackend) SubscribeFilterLogs(ctx context.Context, query ethereum.FilterQuery, ch chan<- types.Log) (ethereum.Subscription, error) {
	client := w.getGethClient()
	return client.SubscribeFilterLogs(ctx, query, ch)
}

func (w *WrappedContractBackend) CallContract(ctx context.Context, msg ethereum.CallMsg, blockNumber *big.Int) ([]byte, error) {
	var hex hexutil.Bytes
	client := w.getGethClient()
	err := client.Client().CallContext(ctx, &hex, "eth_call", toCallArg(msg), toBlockNumArg(blockNumber))
	if err != nil {
		return nil, err
	}
	return hex, nil
}

func (w *WrappedContractBackend) PendingCallContract(ctx context.Context, msg ethereum.CallMsg) ([]byte, error) {
	var hex hexutil.Bytes
	client := w.getGethClient()
	err := client.Client().CallContext(ctx, &hex, "eth_call", toCallArg(msg), "pending")
	if err != nil {
		return nil, err
	}
	return hex, nil
}

// COPIED FROM go-ethereum/ethclient/gethclient - must be kept up to date!
func toBlockNumArg(number *big.Int) string {
	if number == nil {
		return "latest"
	}
	if number.Sign() >= 0 {
		return hexutil.EncodeBig(number)
	}
	// It's negative.
	if number.IsInt64() {
		return rpc.BlockNumber(number.Int64()).String()
	}
	// It's negative and large, which is invalid.
	return fmt.Sprintf("<invalid %d>", number)
}

// COPIED FROM go-ethereum/ethclient/gethclient - must be kept up to date!
// Modified to include legacy 'data' as well as 'input' in order to support non-compliant servers.
func toCallArg(msg ethereum.CallMsg) interface{} {
	arg := map[string]interface{}{
		"from": msg.From,
		"to":   msg.To,
	}
	if len(msg.Data) > 0 {
		arg["input"] = hexutil.Bytes(msg.Data)
		arg["data"] = hexutil.Bytes(msg.Data) // duplicate legacy field for compatibility (required by Geth < v1.11.0)
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
