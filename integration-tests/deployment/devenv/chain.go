package devenv

import (
	"context"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	chainselectors "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/chainlink/integration-tests/deployment"
)

const (
	EVMChainType = "EVM"
)

// ChainConfig holds the configuration for a with a deployer key which can be used to send transactions to the chain.
type ChainConfig struct {
	ChainID     uint64             // chain id as per EIP-155, mainly applicable for EVM chains
	ChainName   string             // name of the chain populated from chainselector repo
	ChainType   string             // should denote the chain family. Acceptable values are EVM, COSMOS, SOLANA, STARKNET, APTOS etc
	WSRPCs      []string           // websocket rpcs to connect to the chain
	HTTPRPCs    []string           // http rpcs to connect to the chain
	DeployerKey *bind.TransactOpts // key to send transactions to the chain
}

func NewChains(logger logger.Logger, configs []ChainConfig) (map[uint64]deployment.Chain, error) {
	chains := make(map[uint64]deployment.Chain)
	for _, chainCfg := range configs {
		selector, err := chainselectors.SelectorFromChainId(chainCfg.ChainID)
		if err != nil {
			return nil, fmt.Errorf("failed to get selector from chain id %d: %w", chainCfg.ChainID, err)
		}
		// TODO : better client handling
		var ec *ethclient.Client
		for _, rpc := range chainCfg.WSRPCs {
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
			Selector:    selector,
			Client:      ec,
			DeployerKey: chainCfg.DeployerKey,
			Confirm: func(tx *types.Transaction) (uint64, error) {
				var blockNumber uint64
				if tx == nil {
					return 0, fmt.Errorf("tx was nil, nothing to confirm")
				}
				ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
				defer cancel()
				chainId, err := ec.ChainID(ctx)
				if err != nil {
					return blockNumber, fmt.Errorf("failed to get chain id: %w", err)
				}
				receipt, err := bind.WaitMined(ctx, ec, tx)
				if err != nil {
					return blockNumber, fmt.Errorf("failed to get confirmed receipt for chain %d: %w", chainId, err)
				}
				if receipt == nil {
					return blockNumber, fmt.Errorf("receipt was nil for tx %s", tx.Hash().Hex())
				}
				blockNumber = receipt.BlockNumber.Uint64()
				if receipt.Status == 0 {
					errReason, err := deployment.GetErrorReasonFromTx(ec, chainCfg.DeployerKey.From, *tx, receipt)
					if err == nil && errReason != "" {
						return blockNumber, fmt.Errorf("tx %s reverted,error reason: %s", tx.Hash().Hex(), errReason)
					}
					return blockNumber, fmt.Errorf("tx %s reverted, could not decode error reason", tx.Hash().Hex())
				}
				return blockNumber, nil
			},
		}
	}
	return chains, nil
}
