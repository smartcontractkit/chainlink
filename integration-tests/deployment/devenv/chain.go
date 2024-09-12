package devenv

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/sethvargo/go-retry"
	chainselectors "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/chainlink/integration-tests/deployment"
)

type ChainConfig struct {
	ChainId         uint64
	ChainName       string
	ChainType       string
	WsRpcs          []string
	HttpRpcs        []string
	PrivateWsRpcs   []string // applicable for private chains only so that nodes within same cluster can connect internally
	PrivateHttpRpcs []string // applicable for private chains only so that nodes within same cluster can connect internally
	DeployerKey     *bind.TransactOpts
}

type RegistryConfig struct {
	EVMChainID uint64
	Contract   common.Address
}

func NewChains(logger logger.Logger, configs []ChainConfig) (map[uint64]deployment.Chain, error) {
	chains := make(map[uint64]deployment.Chain)
	for _, chainCfg := range configs {
		selector, err := chainselectors.SelectorFromChainId(chainCfg.ChainId)
		if err != nil {
			return nil, fmt.Errorf("failed to get selector from chain id %d: %w", chainCfg.ChainId, err)
		}
		// TODO : better client handling
		var ec *ethclient.Client
		for _, rpc := range chainCfg.WsRpcs {
			ec, err = ethclient.Dial(rpc)
			if err != nil {
				logger.Warnf("failed to dial ws rpc %s", rpc)
				continue
			}
			logger.Infof("connected to ws rpc %s", rpc)
			break
		}
		if ec == nil {
			return nil, fmt.Errorf("failed to connect to chain %s", chainCfg.ChainName)
		}
		chains[selector] = deployment.Chain{
			Selector:       selector,
			Client:         ec,
			DeployerKey:    chainCfg.DeployerKey,
			LatestBlockNum: ec.BlockNumber,
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
							t, _, err := ec.TransactionByHash(context.Background(), tx)
							if err != nil {
								return fmt.Errorf("tx %s reverted, failed to get transaction: %w", tx, err)
							}
							errReason, err := deployment.GetErrorReasonFromTx(ec, chainCfg.DeployerKey.From, *t, receipt)
							if err == nil && errReason != "" {
								return fmt.Errorf("tx %s reverted,error reason: %s", tx.Hex(), errReason)
							}
							return fmt.Errorf("tx %s reverted, could not decode error reason", tx.Hex())
						}
						return nil
					})
				return blockNumber, err
			},
		}
	}
	return chains, nil
}

// TODO : Remove this when seth is integrated.
func FundAddress(ctx context.Context, from *bind.TransactOpts, to common.Address, amount *big.Int, c deployment.Chain) error {
	nonce, err := c.Client.PendingNonceAt(ctx, from.From)
	if err != nil {
		return fmt.Errorf("failed to get nonce: %w", err)
	}
	gp, err := c.Client.SuggestGasPrice(ctx)
	if err != nil {
		return fmt.Errorf("failed to suggest gas price: %w", err)
	}
	rawTx := types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		GasPrice: gp,
		Gas:      21000,
		To:       &to,
		Value:    amount,
	})
	signedTx, err := from.Signer(from.From, rawTx)
	if err != nil {
		return fmt.Errorf("failed to sign tx: %w", err)
	}
	err = c.Client.SendTransaction(ctx, signedTx)
	if err != nil {
		return fmt.Errorf("failed to send tx: %w", err)
	}
	_, err = c.Confirm(signedTx.Hash())
	return err
}
