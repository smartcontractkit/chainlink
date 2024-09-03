package wrappers

import (
	"context"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/avast/retry-go/v4"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/seth"
	evmClient "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
)

// WrappedContractBackend is a wrapper around the go-ethereum ContractBackend interface. It's a thin wrapper
// around the go-ethereum/ethclient.Client, which replaces only CallContract and PendingCallContract calls with
// methods that send data both in "input" and "data" field for backwards compatibility with older clients. Other methods
// are passed through to the underlying client.
type WrappedContractBackend struct {
	evmClient   blockchain.EVMClient
	sethClient  *seth.Client
	logger      zerolog.Logger
	maxAttempts uint
	retryDelay  time.Duration
	withRetries bool
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

// MustNewRetryingWrappedContractBackend creates a new WrappedContractBackend, which retries read-only operations every 'retryDelay' until
// 'maxAttempts' are reached. It works only with Seth, because EVMClient has some retrying capability already included.
func MustNewRetryingWrappedContractBackend(sethClient *seth.Client, logger zerolog.Logger, maxAttempts uint, retryDelay time.Duration) *WrappedContractBackend {
	if sethClient == nil {
		panic("Must provide at Seth client reference")
	}

	return &WrappedContractBackend{
		sethClient:  sethClient,
		logger:      logger,
		maxAttempts: maxAttempts,
		retryDelay:  retryDelay,
		withRetries: true,
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

	var fn = func() ([]byte, error) {
		client := w.getGethClient()
		return client.CodeAt(ctx, contract, blockNumber)
	}

	ethHeadBanger := newEthHeadBangerFromWrapper[[]byte](w)
	return ethHeadBanger.retry("CodeAt", fn)
}

func (w *WrappedContractBackend) PendingCodeAt(ctx context.Context, contract common.Address) ([]byte, error) {
	if ctxErr := w.getErrorFromContext(ctx); ctxErr != nil {
		return nil, errors.Wrapf(ctxErr, "the context you passed had an error set. Won't call PendingCodeAt")
	}

	var fn = func() ([]byte, error) {
		client := w.getGethClient()
		return client.PendingCodeAt(ctx, contract)
	}

	ethHeadBanger := newEthHeadBangerFromWrapper[[]byte](w)
	return ethHeadBanger.retry("PendingCodeAt", fn)
}

func (w *WrappedContractBackend) CodeAtHash(ctx context.Context, contract common.Address, blockHash common.Hash) ([]byte, error) {
	if ctxErr := w.getErrorFromContext(ctx); ctxErr != nil {
		return nil, errors.Wrapf(ctxErr, "the context you passed had an error set. Won't call CodeAtHash")
	}

	var fn = func() ([]byte, error) {
		client := w.getGethClient()
		return client.CodeAtHash(ctx, contract, blockHash)
	}

	ethHeadBanger := newEthHeadBangerFromWrapper[[]byte](w)
	return ethHeadBanger.retry("CodeAtHash", fn)
}

func (w *WrappedContractBackend) CallContractAtHash(ctx context.Context, call ethereum.CallMsg, blockHash common.Hash) ([]byte, error) {
	if ctxErr := w.getErrorFromContext(ctx); ctxErr != nil {
		return nil, errors.Wrapf(ctxErr, "the context you passed had an error set. Won't call CallContractAtHash")
	}

	var fn = func() ([]byte, error) {
		client := w.getGethClient()
		return client.CallContractAtHash(ctx, call, blockHash)
	}

	ethHeadBanger := newEthHeadBangerFromWrapper[[]byte](w)
	return ethHeadBanger.retry("CallContractAtHash", fn)
}

func (w *WrappedContractBackend) HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error) {
	if ctxErr := w.getErrorFromContext(ctx); ctxErr != nil {
		return nil, errors.Wrapf(ctxErr, "the context you passed had an error set. Won't call HeaderByNumber")
	}

	var fn = func() (*types.Header, error) {
		client := w.getGethClient()
		return client.HeaderByNumber(ctx, number)
	}

	ethHeadBanger := newEthHeadBangerFromWrapper[*types.Header](w)
	return ethHeadBanger.retry("HeaderByNumber", fn)
}

func (w *WrappedContractBackend) PendingNonceAt(ctx context.Context, account common.Address) (uint64, error) {
	if ctxErr := w.getErrorFromContext(ctx); ctxErr != nil {
		return 0, errors.Wrapf(ctxErr, "the context you passed had an error set. Won't call PendingNonceAt")
	}

	var fn = func() (uint64, error) {
		client := w.getGethClient()
		return client.PendingNonceAt(ctx, account)
	}

	ethHeadBanger := newEthHeadBangerFromWrapper[uint64](w)
	return ethHeadBanger.retry("PendingNonceAt", fn)
}

func (w *WrappedContractBackend) SuggestGasPrice(ctx context.Context) (*big.Int, error) {
	if ctxErr := w.getErrorFromContext(ctx); ctxErr != nil {
		return nil, errors.Wrapf(ctxErr, "the context you passed had an error set. Won't call SuggestGasPrice")
	}

	var fn = func() (*big.Int, error) {
		client := w.getGethClient()
		return client.SuggestGasPrice(ctx)
	}

	ethHeadBanger := newEthHeadBangerFromWrapper[*big.Int](w)
	return ethHeadBanger.retry("SuggestGasPrice", fn)
}

func (w *WrappedContractBackend) SuggestGasTipCap(ctx context.Context) (*big.Int, error) {
	if ctxErr := w.getErrorFromContext(ctx); ctxErr != nil {
		return nil, errors.Wrapf(ctxErr, "the context you passed had an error set. Won't call SuggestGasTipCap")
	}

	var fn = func() (*big.Int, error) {
		client := w.getGethClient()
		return client.SuggestGasTipCap(ctx)
	}

	ethHeadBanger := newEthHeadBangerFromWrapper[*big.Int](w)
	return ethHeadBanger.retry("SuggestGasTipCap", fn)
}

func (w *WrappedContractBackend) EstimateGas(ctx context.Context, call ethereum.CallMsg) (gas uint64, err error) {
	if ctxErr := w.getErrorFromContext(ctx); ctxErr != nil {
		return 0, errors.Wrapf(ctxErr, "the context you passed had an error set. Won't call EstimateGas")
	}

	var fn = func() (uint64, error) {
		client := w.getGethClient()
		return client.EstimateGas(ctx, call)
	}

	ethHeadBanger := newEthHeadBangerFromWrapper[uint64](w)
	return ethHeadBanger.retry("EstimateGas", fn)
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

	var fn = func() ([]types.Log, error) {
		client := w.getGethClient()
		return client.FilterLogs(ctx, query)
	}

	ethHeadBanger := newEthHeadBangerFromWrapper[[]types.Log](w)
	return ethHeadBanger.retry("FilterLogs", fn)
}

func (w *WrappedContractBackend) SubscribeFilterLogs(ctx context.Context, query ethereum.FilterQuery, ch chan<- types.Log) (ethereum.Subscription, error) {
	if ctxErr := w.getErrorFromContext(ctx); ctxErr != nil {
		return nil, errors.Wrapf(ctxErr, "the context you passed had an error set. Won't call SubscribeFilterLogs")
	}

	var fn = func() (ethereum.Subscription, error) {
		client := w.getGethClient()
		return client.SubscribeFilterLogs(ctx, query, ch)
	}

	ethHeadBanger := newEthHeadBangerFromWrapper[ethereum.Subscription](w)
	return ethHeadBanger.retry("SubscribeFilterLogs", fn)
}

func (w *WrappedContractBackend) CallContract(ctx context.Context, msg ethereum.CallMsg, blockNumber *big.Int) ([]byte, error) {
	if ctxErr := w.getErrorFromContext(ctx); ctxErr != nil {
		return nil, errors.Wrapf(ctxErr, "the context you passed had an error set. Won't call CallContract")
	}

	var fn = func() ([]byte, error) {
		var hex hexutil.Bytes
		client := w.getGethClient()
		err := client.Client().CallContext(ctx, &hex, "eth_call", evmClient.ToBackwardCompatibleCallArg(msg), evmClient.ToBackwardCompatibleBlockNumArg(blockNumber))
		if err != nil {
			return nil, err
		}
		return hex, nil
	}

	ethHeadBanger := newEthHeadBangerFromWrapper[[]byte](w)
	return ethHeadBanger.retry("CallContract", fn)
}

func (w *WrappedContractBackend) PendingCallContract(ctx context.Context, msg ethereum.CallMsg) ([]byte, error) {
	if ctxErr := w.getErrorFromContext(ctx); ctxErr != nil {
		return nil, errors.Wrapf(ctxErr, "the context you passed had an error set. Won't call PendingCallContract")
	}

	var fn = func() ([]byte, error) {
		var hex hexutil.Bytes
		client := w.getGethClient()
		err := client.Client().CallContext(ctx, &hex, "eth_call", evmClient.ToBackwardCompatibleCallArg(msg), "pending")
		if err != nil {
			return nil, err
		}
		return hex, nil
	}

	ethHeadBanger := newEthHeadBangerFromWrapper[[]byte](w)
	return ethHeadBanger.retry("PendingCallContract", fn)
}

func (w *WrappedContractBackend) getErrorFromContext(ctx context.Context) error {
	if ctxErr := ctx.Value(seth.ContextErrorKey{}); ctxErr != nil {
		if v, ok := ctxErr.(error); ok {
			return v
		}
		return errors.Wrapf(errors.New("unknown error type"), "error in context: %v", ctxErr)
	}

	return nil
}

// ethHeadBanger is just a fancy name for a struct that retries a function a number of times with a delay between each attempt
type ethHeadBanger[ReturnType any] struct {
	logger      zerolog.Logger
	maxAttempts uint
	retryDelay  time.Duration
}

func newEthHeadBangerFromWrapper[ResultType any](wrapper *WrappedContractBackend) ethHeadBanger[ResultType] {
	return ethHeadBanger[ResultType]{
		logger:      wrapper.logger,
		maxAttempts: wrapper.maxAttempts,
		retryDelay:  wrapper.retryDelay,
	}
}

func (e ethHeadBanger[ReturnType]) retry(functionName string, fnToRetry func() (ReturnType, error)) (ReturnType, error) {
	var result ReturnType
	err := retry.Do(func() error {
		var err error
		result, err = fnToRetry()

		return err
	},
		retry.RetryIf(func(err error) bool {
			if err.Error() == rpc.ErrClientQuit.Error() ||
				err.Error() == rpc.ErrBadResult.Error() ||
				strings.Contains(err.Error(), "connection") ||
				strings.Contains(err.Error(), "EOF") {
				return true
			}

			e.logger.Error().Err(err).Msgf("Error in %s. Not retrying.", functionName)

			return false
		}),
		retry.Attempts(e.maxAttempts),
		retry.Delay(e.retryDelay),
		retry.OnRetry(func(n uint, err error) {
			e.logger.Info().
				Str("Attempt", fmt.Sprintf("%d/%d", n+1, 10)).
				Str("Error", err.Error()).
				Msgf("Retrying %s", functionName)
		}),
	)

	return result, err
}
