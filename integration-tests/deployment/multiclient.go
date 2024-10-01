package deployment

import (
	"context"
	"fmt"
	"math/big"
	"reflect"
	"runtime"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/pkg/errors"
	"github.com/sethvargo/go-retry"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
)

const (
	RPC_DEFAULT_RETRY_ATTEMPTS = uint64(10)
	RPC_DEFAULT_RETRY_DELAY    = 1000 * time.Millisecond
)

type RetryConfig struct {
	Attempts uint64
	Delay    time.Duration
}

func defaultRetryConfig() RetryConfig {
	return RetryConfig{
		Attempts: RPC_DEFAULT_RETRY_ATTEMPTS,
		Delay:    RPC_DEFAULT_RETRY_DELAY,
	}
}

type RPC struct {
	Name  string
	WSURL string
	// TODO: http fallback needed for some networks?
}

// MultiClient should comply with the OnchainClient interface
var _ OnchainClient = &MultiClient{}

type MultiClient struct {
	logger logger.Logger
	*ethclient.Client
	Backups     []*ethclient.Client
	RetryConfig RetryConfig
}

func WithRetryConfig(attempts uint64, delay time.Duration) func(client *MultiClient) {
	return func(client *MultiClient) {
		client.RetryConfig = RetryConfig{
			Attempts: attempts,
			Delay:    delay,
		}
	}
}

func NewMultiClient(lggr logger.Logger, rpcs []RPC, opts ...func(client *MultiClient)) (*MultiClient, error) {
	if len(rpcs) == 0 {
		return nil, fmt.Errorf("no RPCs provided, need at least one")
	}
	mc := &MultiClient{
		logger: lggr,
	}
	clients := make([]*ethclient.Client, 0, len(rpcs))
	for i, rpc := range rpcs {
		client, err := ethclient.Dial(rpc.WSURL)
		if err != nil {
			name := rpc.Name
			if name == "" {
				name = fmt.Sprintf("at index %d", i+1)
			}
			lggr.Warnf("failed to dial rpc %s ending in %s, moving to next one", name, rpc.WSURL[len(rpc.WSURL)-4:])
			continue
		}
		clients = append(clients, client)
	}
	if len(clients) == 0 {
		return nil, fmt.Errorf("all RPCs failed, try again with different URLs")
	}
	mc.Client = clients[0]
	if len(clients) > 1 {
		mc.Backups = clients[1:]
	} else {
		lggr.Warn("Only one RPC provided, no backups available")
	}
	mc.RetryConfig = defaultRetryConfig()
	mc.logger = lggr

	for _, opt := range opts {
		opt(mc)
	}
	return mc, nil
}

func (mc *MultiClient) Confirm(parent context.Context, tx *types.Transaction, timeout time.Duration) (*types.Receipt, error) {
	if tx == nil {
		return nil, fmt.Errorf("tx was nil, nothing to confirm")
	}
	var receipt *types.Receipt
	ctx, cancel := context.WithTimeout(parent, timeout)
	defer cancel()
	success := make(chan bool, 1)
	for i, client := range append([]*ethclient.Client{mc.Client}, mc.Backups...) {
		i := i
		client := client
		go func() {
			var err error
			receipt, err = bind.WaitMined(ctx, client, tx)
			if err != nil {
				mc.logger.Warnf("error %+v with client %d for op Confirm", err, i+1)
			}
			if receipt != nil {
				success <- true
			}
		}()
	}
	for {
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("confirmation timed out, all clients failed")
		case <-success:
			return receipt, nil
		}
	}
}

func (mc *MultiClient) HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error) {
	var header *types.Header
	err := mc.retryWithBackups(ctx, func(client *ethclient.Client) error {
		var err error
		header, err = client.HeaderByNumber(ctx, number)
		return err
	})
	return header, err
}

func (mc *MultiClient) CallContract(ctx context.Context, call ethereum.CallMsg, blockNumber *big.Int) ([]byte, error) {
	var result []byte
	err := mc.retryWithBackups(ctx, func(client *ethclient.Client) error {
		var err error
		result, err = client.CallContract(ctx, call, blockNumber)
		return err
	})
	return result, err
}

func (mc *MultiClient) TransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error) {
	var receipt *types.Receipt
	// TransactionReceipt might return ethereum.NotFound error if the transaction is not yet mined
	err := mc.retryWithBackups(ctx, func(client *ethclient.Client) error {
		var err error
		receipt, err = client.TransactionReceipt(ctx, txHash)
		return err
	}, ethereum.NotFound)
	return receipt, err
}

func (mc *MultiClient) SendTransaction(ctx context.Context, tx *types.Transaction) error {
	return mc.retryWithBackups(ctx, func(client *ethclient.Client) error {
		return client.SendTransaction(ctx, tx)
	})
}

func (mc *MultiClient) CodeAt(ctx context.Context, account common.Address, blockNumber *big.Int) ([]byte, error) {
	var code []byte
	err := mc.retryWithBackups(ctx, func(client *ethclient.Client) error {
		var err error
		code, err = client.CodeAt(ctx, account, blockNumber)
		return err
	})
	return code, err
}

func (mc *MultiClient) NonceAt(ctx context.Context, account common.Address) (uint64, error) {
	var count uint64
	err := mc.retryWithBackups(ctx, func(client *ethclient.Client) error {
		var err error
		count, err = client.NonceAt(ctx, account, nil)
		return err
	})
	return count, err
}

// retryWithBackups will retry the operation with the main client and then with the backup clients
// If the operation returns an error that is in the acceptedErrors list, it will not retry the operation or try the next client
func (mc *MultiClient) retryWithBackups(ctx context.Context, op func(*ethclient.Client) error, acceptedErrors ...error) error {
	var err2 error
	funcName := runtime.FuncForPC(reflect.ValueOf(op).Pointer()).Name()
	for i, client := range append([]*ethclient.Client{mc.Client}, mc.Backups...) {
		err2 = retry.Do(ctx, retry.WithMaxRetries(mc.RetryConfig.Attempts, retry.NewConstant(mc.RetryConfig.Delay)),
			func(ctx context.Context) error {
				err := op(client)
				if err != nil {
					// If the error is one of the accepted errors return err without retrying
					for _, acceptedError := range acceptedErrors {
						if errors.Is(err, acceptedError) {
							return err
						}
					}
					mc.logger.Warnf("error %+v with client %d for op %s", err, i+1, funcName)
					return retry.RetryableError(err)
				}
				return nil
			})
		if err2 == nil {
			return nil
		}
		// Check if the error is one of the accepted errors
		// In case of an accepted error, we don't need to try the next client
		for _, acceptedError := range acceptedErrors {
			if errors.Is(err2, acceptedError) {
				return err2
			}
		}
		mc.logger.Infof("Client %d failed, trying next client\n", i+1)
	}
	return errors.Wrapf(err2, "All backup clients %v failed", mc.Backups)
}
