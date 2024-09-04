package devenv

import (
	"context"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/sethvargo/go-retry"
	chainselectors "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink/integration-tests/deployment"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

type ChainConfig struct {
	ChainId uint64
	// TODO : use a slice of rpc urls for failing over to the available rpcs
	WsRpc       string
	HttpRpc     string
	DeployerKey *bind.TransactOpts
}

type RegistryConfig struct {
	EVMChainID uint64
	Contract   common.Address
}

func NewChainConfig(chainId uint64, wsRpc, httpRpc string, deployerKey *bind.TransactOpts) ChainConfig {
	return ChainConfig{
		ChainId:     chainId,
		WsRpc:       wsRpc,
		HttpRpc:     httpRpc,
		DeployerKey: deployerKey,
	}
}

func NewChains(logger logger.Logger, configs []ChainConfig) (map[uint64]deployment.Chain, error) {
	chains := make(map[uint64]deployment.Chain)
	for _, chainCfg := range configs {
		selector, err := chainselectors.SelectorFromChainId(chainCfg.ChainId)
		if err != nil {
			return nil, fmt.Errorf("failed to get selector from chain id %d: %w", chainCfg.ChainId, err)
		}
		// TODO : better client handling
		ec, err := ethclient.Dial(chainCfg.WsRpc)
		if err != nil {
			return nil, fmt.Errorf("failed to dial ws rpc %s: %w", chainCfg.WsRpc, err)
		}
		chains[selector] = deployment.Chain{
			Selector:    selector,
			Client:      ec,
			DeployerKey: chainCfg.DeployerKey,
			Confirm: func(tx common.Hash) (uint64, error) {
				var blockNumber uint64
				err := retry.Do(context.Background(),
					retry.WithMaxDuration(3*time.Minute, retry.NewFibonacci(1*time.Second)),
					func(ctx context.Context) error {
						receipt, err := ec.TransactionReceipt(ctx, tx)
						if err != nil {
							return retry.RetryableError(fmt.Errorf("failed to get receipt: %w", err))
						}
						if receipt == nil {
							blockNumber = receipt.BlockNumber.Uint64()
						}
						if receipt.Status == 0 {
							return fmt.Errorf("tx %s reverted", tx.Hex())
						}
						return nil
					})
				return blockNumber, err
			},
		}
	}
	return chains, nil
}
