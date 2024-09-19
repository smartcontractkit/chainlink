package deployment

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/avast/retry-go/v4"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/pkg/errors"
)

const (
	RPC_DEFAULT_RETRY_ATTEMPTS = 10
	RPC_DEFAULT_RETRY_DELAY    = 1000 * time.Millisecond
)

type RetryConfig struct {
	Attempts uint
	Delay    time.Duration
}

func defaultRetryConfig() RetryConfig {
	return RetryConfig{
		Attempts: RPC_DEFAULT_RETRY_ATTEMPTS,
		Delay:    RPC_DEFAULT_RETRY_DELAY,
	}
}

type RPC struct {
	WSURL string
	// TODO: http fallback needed for some networks?
}

// MultiClient should comply with the OnchainClient interface
var _ OnchainClient = &MultiClient{}

type MultiClient struct {
	*ethclient.Client
	Backups     []*ethclient.Client
	RetryConfig RetryConfig
}

func NewMultiClient(rpcs []RPC, opts ...func(client *MultiClient)) (*MultiClient, error) {
	if len(rpcs) == 0 {
		return nil, fmt.Errorf("No RPCs provided, need at least one")
	}
	var mc MultiClient
	clients := make([]*ethclient.Client, 0, len(rpcs))
	for _, rpc := range rpcs {
		client, err := ethclient.Dial(rpc.WSURL)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to dial %s", rpc.WSURL)
		}
		clients = append(clients, client)
	}
	mc.Client = clients[0]
	mc.Backups = clients[1:]
	mc.RetryConfig = defaultRetryConfig()

	for _, opt := range opts {
		opt(&mc)
	}
	return &mc, nil
}

func (mc *MultiClient) TransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error) {
	var receipt *types.Receipt
	err := mc.retryWithBackups(func(client *ethclient.Client) error {
		var err error
		receipt, err = client.TransactionReceipt(ctx, txHash)
		return err
	})
	return receipt, err
}

func (mc *MultiClient) SendTransaction(ctx context.Context, tx *types.Transaction) error {
	return mc.retryWithBackups(func(client *ethclient.Client) error {
		return client.SendTransaction(ctx, tx)
	})
}

func (mc *MultiClient) CodeAt(ctx context.Context, account common.Address, blockNumber *big.Int) ([]byte, error) {
	var code []byte
	err := mc.retryWithBackups(func(client *ethclient.Client) error {
		var err error
		code, err = client.CodeAt(ctx, account, blockNumber)
		return err
	})
	return code, err
}

func (mc *MultiClient) NonceAt(ctx context.Context, account common.Address) (uint64, error) {
	var count uint64
	err := mc.retryWithBackups(func(client *ethclient.Client) error {
		var err error
		count, err = client.NonceAt(ctx, account, nil)
		return err
	})
	return count, err
}

func (mc *MultiClient) retryWithBackups(op func(*ethclient.Client) error) error {
	var err error
	for _, client := range append([]*ethclient.Client{mc.Client}, mc.Backups...) {
		err2 := retry.Do(func() error {
			err = op(client)
			if err != nil {
				// TODO: logger?
				fmt.Printf("Error %v with client %v\n", err, client)
				return err
			}
			return nil
		}, retry.Attempts(mc.RetryConfig.Attempts), retry.Delay(mc.RetryConfig.Delay))
		if err2 == nil {
			return nil
		}
		fmt.Printf("Client %v failed, trying next client\n", client)
	}
	return errors.Wrapf(err, "All backup clients %v failed", mc.Backups)
}
