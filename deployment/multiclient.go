package deployment

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/avast/retry-go/v4"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	chainsel "github.com/smartcontractkit/chain-selectors"
)

const (
	// Default retry configuration for RPC calls
	RPCDefaultRetryAttempts = 10
	RPCDefaultRetryDelay    = 1000 * time.Millisecond

	// Default retry configuration for dialing RPC endpoints
	RPCDefaultDialRetryAttempts = 10
	RPCDefaultDialRetryDelay    = 1000 * time.Millisecond
)

type RetryConfig struct {
	Attempts uint
	Delay    time.Duration
}

func defaultRetryConfig() RetryConfig {
	return RetryConfig{
		Attempts: RPCDefaultRetryAttempts,
		Delay:    RPCDefaultRetryDelay,
	}
}

// MultiClient should comply with the OnchainClient interface
var _ OnchainClient = &MultiClient{}

type MultiClient struct {
	*ethclient.Client
	Backups     []*ethclient.Client
	RetryConfig RetryConfig
	lggr        logger.Logger
	chainName   string
}

func NewMultiClient(lggr logger.Logger, rpcsCfg RPCConfig, opts ...func(client *MultiClient)) (*MultiClient, error) {
	if len(rpcsCfg.RPCs) == 0 {
		return nil, errors.New("no RPCs provided, need at least one")
	}
	// Set the chain name
	chain, exists := chainsel.ChainBySelector(rpcsCfg.ChainSelector)
	if !exists {
		return nil, fmt.Errorf("chain with selector %d not found", rpcsCfg.ChainSelector)
	}
	mc := MultiClient{lggr: lggr, chainName: chain.Name}

	clients := make([]*ethclient.Client, 0, len(rpcsCfg.RPCs))
	for i, rpc := range rpcsCfg.RPCs {
		client, err := mc.dialWithRetry(rpc, lggr)
		if err != nil {
			lggr.Warnf("failed to dial client %d for RPC '%s' trying with the next one: %v", i, rpc.Name, err)
			continue
		}
		clients = append(clients, client)
	}

	if len(clients) == 0 {
		return nil, errors.New("no valid RPC clients created")
	}

	mc.Client = clients[0]
	mc.Backups = clients[1:]
	mc.RetryConfig = defaultRetryConfig()

	for _, opt := range opts {
		opt(&mc)
	}
	return &mc, nil
}

func (mc *MultiClient) SendTransaction(ctx context.Context, tx *types.Transaction) error {
	return mc.retryWithBackups("SendTransaction", func(client *ethclient.Client) error {
		return client.SendTransaction(ctx, tx)
	})
}

func (mc *MultiClient) CallContract(ctx context.Context, msg ethereum.CallMsg, blockNumber *big.Int) ([]byte, error) {
	var result []byte
	err := mc.retryWithBackups("CallContract", func(client *ethclient.Client) error {
		var err error
		result, err = client.CallContract(ctx, msg, blockNumber)
		return err
	})
	return result, err
}

func (mc *MultiClient) CallContractAtHash(ctx context.Context, msg ethereum.CallMsg, blockHash common.Hash) ([]byte, error) {
	var result []byte
	err := mc.retryWithBackups("CallContractAtHash", func(client *ethclient.Client) error {
		var err error
		result, err = client.CallContractAtHash(ctx, msg, blockHash)
		return err
	})
	return result, err
}

func (mc *MultiClient) CodeAt(ctx context.Context, account common.Address, blockNumber *big.Int) ([]byte, error) {
	var code []byte
	err := mc.retryWithBackups("CodeAt", func(client *ethclient.Client) error {
		var err error
		code, err = client.CodeAt(ctx, account, blockNumber)
		return err
	})
	return code, err
}

func (mc *MultiClient) CodeAtHash(ctx context.Context, account common.Address, blockHash common.Hash) ([]byte, error) {
	var code []byte
	err := mc.retryWithBackups("CodeAtHash", func(client *ethclient.Client) error {
		var err error
		code, err = client.CodeAtHash(ctx, account, blockHash)
		return err
	})
	return code, err
}

func (mc *MultiClient) NonceAt(ctx context.Context, account common.Address, block *big.Int) (uint64, error) {
	var count uint64
	err := mc.retryWithBackups("NonceAt", func(client *ethclient.Client) error {
		var err error
		count, err = client.NonceAt(ctx, account, block)
		return err
	})
	return count, err
}

func (mc *MultiClient) NonceAtHash(ctx context.Context, account common.Address, blockHash common.Hash) (uint64, error) {
	var count uint64
	err := mc.retryWithBackups("NonceAtHash", func(client *ethclient.Client) error {
		var err error
		count, err = client.NonceAtHash(ctx, account, blockHash)
		return err
	})
	return count, err
}

func (mc *MultiClient) WaitMined(ctx context.Context, tx *types.Transaction) (*types.Receipt, error) {
	mc.lggr.Debugf("Waiting for tx %s to be mined for chain %s", tx.Hash().Hex(), mc.chainName)
	// no retries here because we want to wait for the tx to be mined
	resultCh := make(chan *types.Receipt)
	doneCh := make(chan struct{})

	waitMined := func(client *ethclient.Client, tx *types.Transaction) {
		mc.lggr.Debugf("Waiting for tx %s to be mined with chain %s", tx.Hash().Hex(), mc.chainName)
		receipt, err := bind.WaitMined(ctx, client, tx)
		if err != nil {
			mc.lggr.Warnf("WaitMined error %v with chain %s", err, mc.chainName)
			return
		}
		select {
		case resultCh <- receipt:
		case <-doneCh:
			return
		}
	}

	for _, client := range append([]*ethclient.Client{mc.Client}, mc.Backups...) {
		go waitMined(client, tx)
	}
	var receipt *types.Receipt
	select {
	case receipt = <-resultCh:
		close(doneCh)
		mc.lggr.Debugf("Tx %s mined with chain %s", tx.Hash().Hex(), mc.chainName)
		return receipt, nil
	case <-ctx.Done():
		mc.lggr.Warnf("WaitMined context done %v", ctx.Err())
		close(doneCh)
		return nil, ctx.Err()
	}
}

func (mc *MultiClient) retryWithBackups(opName string, op func(*ethclient.Client) error) error {
	var err error
	for i, client := range append([]*ethclient.Client{mc.Client}, mc.Backups...) {
		err2 := retry.Do(func() error {
			mc.lggr.Debugf("Trying op %s with chain %s client index %d", opName, mc.chainName, i)
			err = op(client)
			if err != nil {
				mc.lggr.Warnf("retryable error '%s' for op %s with chain %s client index %d", err.Error(), opName, mc.chainName, i)
				return err
			}
			return nil
		}, retry.Attempts(mc.RetryConfig.Attempts), retry.Delay(mc.RetryConfig.Delay))
		if err2 == nil {
			return nil
		}
		mc.lggr.Infof("Client at index %d failed, trying next client chain %s", i, mc.chainName)
	}
	return errors.Join(err, fmt.Errorf("All backup clients %v failed for chain %s", mc.Backups, mc.chainName))
}

func (mc *MultiClient) dialWithRetry(rpc RPC, lggr logger.Logger) (*ethclient.Client, error) {
	endpoint, err := rpc.ToEndpoint()
	if err != nil {
		return nil, err
	}

	var client *ethclient.Client
	err = retry.Do(func() error {
		var err2 error
		mc.lggr.Debugf("dialing endpoint '%s' for RPC %s for chain %s", endpoint, rpc.Name, mc.chainName)
		client, err2 = ethclient.Dial(endpoint)
		if err2 != nil {
			lggr.Warnf("retryable error for RPC %s:%s for chain %s  %v", rpc.Name, endpoint, mc.chainName, err)
			return err2
		}
		return nil
	}, retry.Attempts(RPCDefaultDialRetryAttempts), retry.Delay(RPCDefaultDialRetryDelay))

	if err != nil {
		return nil, errors.Join(err, fmt.Errorf("failed to dial endpoint '%s' for RPC %s for chain %s after retries", endpoint, rpc.Name, mc.chainName))
	}
	return client, nil
}
