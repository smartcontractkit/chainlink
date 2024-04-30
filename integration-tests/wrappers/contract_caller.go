package wrappers

import (
	"context"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/avast/retry-go"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/rs/zerolog"
	"github.com/smartcontractkit/seth"

	evmClient "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"

	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/logging"
)

// WrappedContractBackend is a wrapper around the go-ethereum ContractBackend interface. It's a thin wrapper
// around the go-ethereum/ethclient.Client, which replaces only CallContract and PendingCallContract calls with
// methods that send data both in "input" and "data" field for backwards compatibility with older clients. Other methods
// are passed through to the underlying client.
type WrappedContractBackend struct {
	evmClient  blockchain.EVMClient
	sethClient *seth.Client
	l          zerolog.Logger
}

// MustNewWrappedContractBackend creates a new WrappedContractBackend with the given clients
func MustNewWrappedContractBackend(evmClient blockchain.EVMClient, sethClient *seth.Client) *WrappedContractBackend {
	if evmClient == nil && sethClient == nil {
		panic("Must provide at least one client")
	}

	return &WrappedContractBackend{
		evmClient:  evmClient,
		sethClient: sethClient,
		l:          logging.GetTestLogger(nil),
	}
}

// MustNewWrappedContractBackendWithLogger creates a new WrappedContractBackend with the given clients and logger instance
func MustNewWrappedContractBackendWithLogger(evmClient blockchain.EVMClient, sethClient *seth.Client, l zerolog.Logger) *WrappedContractBackend {
	if evmClient == nil && sethClient == nil {
		panic("Must provide at least one client")
	}

	return &WrappedContractBackend{
		evmClient:  evmClient,
		sethClient: sethClient,
		l:          l,
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
	var result []byte
	err := retry.Do(func() error {
		var hex hexutil.Bytes
		client := w.getGethClient()
		err := client.Client().CallContext(ctx, &hex, "eth_call", evmClient.ToBackwardCompatibleCallArg(msg), evmClient.ToBackwardCompatibleBlockNumArg(blockNumber))
		if err != nil {
			return err
		}
		result = hex
		return nil
	},
		retry.RetryIf(func(err error) bool {
			if err.Error() == rpc.ErrClientQuit.Error() ||
				err.Error() == rpc.ErrBadResult.Error() ||
				strings.Contains(err.Error(), "connection") ||
				strings.Contains(err.Error(), "EOF") {
				return true
			}

			w.l.Error().Err(err).Msg("Error in CallContract. Not retrying.")

			return false
		}),
		retry.Attempts(uint(10)),
		retry.Delay(1*time.Second),
		retry.OnRetry(func(n uint, err error) {
			w.l.Info().
				Str("Attempt", fmt.Sprintf("%d/%d", n+1, 10)).
				Str("Error", err.Error()).
				Msg("Retrying CallContract")
		}),
	)

	return result, err

	// var hex hexutil.Bytes
	// client := w.getGethClient()
	// err := client.Client().CallContext(ctx, &hex, "eth_call", evmClient.ToBackwardCompatibleCallArg(msg), evmClient.ToBackwardCompatibleBlockNumArg(blockNumber))
	// if err != nil {
	// 	return nil, err
	// }
	// return hex, nil
}

func (w *WrappedContractBackend) PendingCallContract(ctx context.Context, msg ethereum.CallMsg) ([]byte, error) {
	var hex hexutil.Bytes
	client := w.getGethClient()
	err := client.Client().CallContext(ctx, &hex, "eth_call", evmClient.ToBackwardCompatibleCallArg(msg), "pending")
	if err != nil {
		return nil, err
	}
	return hex, nil
}
