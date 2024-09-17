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
)

const (
	RPC_RETRY_ATTEMPTS = 10
	RPC_RETRY_DELAY    = 1000 * time.Millisecond
)

// MultiClient should comply with the coreenv.OnchainClient interface
var _ OnchainClient = &MultiClient{}

type MultiClient struct {
	*ethclient.Client
	backup []*ethclient.Client
}

type RPC struct {
	RPCName string `toml:"rpc_name"`
	HTTPURL string `toml:"http_url"`
	WSURL   string `toml:"ws_url"`
}

func NewMultiClient(rpcs []RPC) *MultiClient {
	if len(rpcs) == 0 {
		panic("No RPCs provided")
	}
	clients := make([]*ethclient.Client, 0, len(rpcs))
	for _, rpc := range rpcs {
		client, err := ethclient.Dial(rpc.HTTPURL)
		if err != nil {
			panic(err)
		}
		clients = append(clients, client)
	}
	return &MultiClient{
		Client: clients[0],
		backup: clients[1:],
	}
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
	for _, client := range append([]*ethclient.Client{mc.Client}, mc.backup...) {
		err2 := retry.Do(func() error {
			err = op(client)
			if err != nil {
				fmt.Printf("  [MultiClient RPC] Retrying with new client, error: %v\n", err)
				return err
			}
			return nil
		}, retry.Attempts(RPC_RETRY_ATTEMPTS), retry.Delay(RPC_RETRY_DELAY))
		if err2 == nil {
			return nil
		}
	}
	return err
}
