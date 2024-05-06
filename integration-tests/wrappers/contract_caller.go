package wrappers

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/seth"

	evmClient "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"

	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
)

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
	if ctxErr := w.getErrorFromContext(ctx); ctxErr != nil {
		return nil, errors.Wrapf(ctxErr, "the context you passed had an error set. Won't call CodeAt")
	}

	client := w.getGethClient()
	return client.CodeAt(ctx, contract, blockNumber)
}

func (w *WrappedContractBackend) PendingCodeAt(ctx context.Context, contract common.Address) ([]byte, error) {
	if ctxErr := w.getErrorFromContext(ctx); ctxErr != nil {
		return nil, errors.Wrapf(ctxErr, "the context you passed had an error set. Won't call PendingCodeAt")
	}

	client := w.getGethClient()
	return client.PendingCodeAt(ctx, contract)
}

func (w *WrappedContractBackend) CodeAtHash(ctx context.Context, contract common.Address, blockHash common.Hash) ([]byte, error) {
	if ctxErr := w.getErrorFromContext(ctx); ctxErr != nil {
		return nil, errors.Wrapf(ctxErr, "the context you passed had an error set. Won't call CodeAtHash")
	}

	client := w.getGethClient()
	return client.CodeAtHash(ctx, contract, blockHash)
}

func (w *WrappedContractBackend) CallContractAtHash(ctx context.Context, call ethereum.CallMsg, blockHash common.Hash) ([]byte, error) {
	if ctxErr := w.getErrorFromContext(ctx); ctxErr != nil {
		return nil, errors.Wrapf(ctxErr, "the context you passed had an error set. Won't call CallContractAtHash")
	}

	client := w.getGethClient()
	return client.CallContractAtHash(ctx, call, blockHash)
}

func (w *WrappedContractBackend) HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error) {
	if ctxErr := w.getErrorFromContext(ctx); ctxErr != nil {
		return nil, errors.Wrapf(ctxErr, "the context you passed had an error set. Won't call HeaderByNumber")
	}

	client := w.getGethClient()
	return client.HeaderByNumber(ctx, number)
}

func (w *WrappedContractBackend) PendingNonceAt(ctx context.Context, account common.Address) (uint64, error) {
	if ctxErr := w.getErrorFromContext(ctx); ctxErr != nil {
		return 0, errors.Wrapf(ctxErr, "the context you passed had an error set. Won't call PendingNonceAt")
	}

	client := w.getGethClient()
	return client.PendingNonceAt(ctx, account)
}

func (w *WrappedContractBackend) SuggestGasPrice(ctx context.Context) (*big.Int, error) {
	if ctxErr := w.getErrorFromContext(ctx); ctxErr != nil {
		return nil, errors.Wrapf(ctxErr, "the context you passed had an error set. Won't call SuggestGasPrice")
	}

	client := w.getGethClient()
	return client.SuggestGasPrice(ctx)
}

func (w *WrappedContractBackend) SuggestGasTipCap(ctx context.Context) (*big.Int, error) {
	if ctxErr := w.getErrorFromContext(ctx); ctxErr != nil {
		return nil, errors.Wrapf(ctxErr, "the context you passed had an error set. Won't call SuggestGasTipCap")
	}

	client := w.getGethClient()
	return client.SuggestGasTipCap(ctx)
}

func (w *WrappedContractBackend) EstimateGas(ctx context.Context, call ethereum.CallMsg) (gas uint64, err error) {
	if ctxErr := w.getErrorFromContext(ctx); ctxErr != nil {
		return 0, errors.Wrapf(ctxErr, "the context you passed had an error set. Won't call EstimateGas")
	}

	client := w.getGethClient()
	return client.EstimateGas(ctx, call)
}

func (w *WrappedContractBackend) SendTransaction(ctx context.Context, tx *types.Transaction) error {
	if ctxErr := w.getErrorFromContext(ctx); ctxErr != nil {
		return errors.Wrapf(ctxErr, "the context you passed had an error set. Won't call SendTransaction")
	}

	client := w.getGethClient()
	return client.SendTransaction(ctx, tx)
}

func (w *WrappedContractBackend) FilterLogs(ctx context.Context, query ethereum.FilterQuery) ([]types.Log, error) {
	if ctxErr := w.getErrorFromContext(ctx); ctxErr != nil {
		return nil, errors.Wrapf(ctxErr, "the context you passed had an error set. Won't call FilterLogs")
	}

	client := w.getGethClient()
	return client.FilterLogs(ctx, query)
}

func (w *WrappedContractBackend) SubscribeFilterLogs(ctx context.Context, query ethereum.FilterQuery, ch chan<- types.Log) (ethereum.Subscription, error) {
	if ctxErr := w.getErrorFromContext(ctx); ctxErr != nil {
		return nil, errors.Wrapf(ctxErr, "the context you passed had an error set. Won't call SubscribeFilterLogs")
	}

	client := w.getGethClient()
	return client.SubscribeFilterLogs(ctx, query, ch)
}

func (w *WrappedContractBackend) CallContract(ctx context.Context, msg ethereum.CallMsg, blockNumber *big.Int) ([]byte, error) {
	if ctxErr := w.getErrorFromContext(ctx); ctxErr != nil {
		return nil, errors.Wrapf(ctxErr, "the context you passed had an error set. Won't call CallContract")
	}

	var hex hexutil.Bytes
	client := w.getGethClient()
	err := client.Client().CallContext(ctx, &hex, "eth_call", evmClient.ToBackwardCompatibleCallArg(msg), evmClient.ToBackwardCompatibleBlockNumArg(blockNumber))
	if err != nil {
		return nil, err
	}
	return hex, nil
}

func (w *WrappedContractBackend) PendingCallContract(ctx context.Context, msg ethereum.CallMsg) ([]byte, error) {
	if ctxErr := w.getErrorFromContext(ctx); ctxErr != nil {
		return nil, errors.Wrapf(ctxErr, "the context you passed had an error set. Won't call PendingCallContract")
	}

	var hex hexutil.Bytes
	client := w.getGethClient()
	err := client.Client().CallContext(ctx, &hex, "eth_call", evmClient.ToBackwardCompatibleCallArg(msg), "pending")
	if err != nil {
		return nil, err
	}
	return hex, nil
}

func (w *WrappedContractBackend) getErrorFromContext(ctx context.Context) error {
	if ctxErr := ctx.Value("error"); ctxErr != nil {
		if v, ok := ctxErr.(error); ok {
			return v
		}
		if v, ok := ctxErr.(seth.ContextErrorValue); ok {
			return v.Error
		}
		return errors.Wrapf(errors.New("unknown error type"), "error in context: %v", ctxErr)
	}

	return nil
}
