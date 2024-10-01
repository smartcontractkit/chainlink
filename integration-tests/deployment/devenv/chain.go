package devenv

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	chainselectors "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/chainlink/integration-tests/deployment"
)

const (
	KMS_DEPLOYER_KEY_ENV        = "KMS_DEPLOYER_KEY_ID"
	KMS_DEPLOYER_KEY_REGION_ENV = "KMS_DEPLOYER_KEY_REGION"
	AWS_PROFILE_NAME_ENV        = "AWS_PROFILE"
	EVMChainType                = "EVM"
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

// SetDeployerKey sets the deployer key for the chain. If private key is not provided, it fetches the deployer key from KMS.
func (c *ChainConfig) SetDeployerKey(pvtKeyStr *string) error {
	if pvtKeyStr != nil && *pvtKeyStr != "" {
		pvtKey, err := crypto.HexToECDSA(*pvtKeyStr)
		if err != nil {
			return fmt.Errorf("failed to convert private key to ECDSA: %w", err)
		}
		deployer, err := bind.NewKeyedTransactorWithChainID(pvtKey, new(big.Int).SetUint64(c.ChainID))
		if err != nil {
			return fmt.Errorf("failed to create transactor: %w", err)
		}
		c.DeployerKey = deployer
		return nil
	}
	kmsConfig, err := deployment.KMSConfigFromEnvVars()
	if err != nil {
		return fmt.Errorf("failed to get kms config from env vars: %w", err)
	}
	kmsClient, err := deployment.NewKMSClient(kmsConfig)
	if err != nil {
		return fmt.Errorf("failed to create KMS client: %w", err)
	}
	evmKMSClient := deployment.NewEVMKMSClient(kmsClient, kmsConfig.KmsDeployerKeyId)
	c.DeployerKey, err = evmKMSClient.GetKMSTransactOpts(context.Background(), new(big.Int).SetUint64(c.ChainID))
	if err != nil {
		return fmt.Errorf("failed to get transactor from KMS client: %w", err)
	}
	return nil
}

func NewChains(lggr logger.Logger, configs []ChainConfig) (map[uint64]deployment.Chain, error) {
	chains := make(map[uint64]deployment.Chain)
	for _, chainCfg := range configs {
		selector, err := chainselectors.SelectorFromChainId(chainCfg.ChainID)
		if err != nil {
			return nil, fmt.Errorf("failed to get selector from chain id %d: %w", chainCfg.ChainID, err)
		}
		var rpcs []deployment.RPC
		for _, rpc := range chainCfg.WSRPCs {
			rpcs = append(rpcs, deployment.RPC{WSURL: rpc})
		}
		mc, err := deployment.NewMultiClient(lggr, rpcs, deployment.WithRetryConfig(5, 1*time.Second))
		if err != nil {
			return nil, fmt.Errorf("failed to create multi client: %w for chain id %d", err, chainCfg.ChainID)
		}

		chainID := chainCfg.ChainID
		chains[selector] = deployment.Chain{
			Selector:    selector,
			Client:      mc,
			DeployerKey: chainCfg.DeployerKey,
			Confirm: func(tx *types.Transaction) (uint64, error) {
				var blockNumber uint64
				if tx == nil {
					return 0, fmt.Errorf("tx was nil, nothing to confirm")
				}
				receipt, err := mc.Confirm(context.Background(), tx, 3*time.Minute)
				if err != nil {
					return 0, err
				}
				if receipt == nil {
					return 0, fmt.Errorf("receipt was nil")
				}
				blockNumber = receipt.BlockNumber.Uint64()
				if receipt.Status == 0 {
					errReason, err := deployment.GetErrorReasonFromTx(mc, chainCfg.DeployerKey.From, *tx, receipt)
					if err == nil && errReason != "" {
						return blockNumber, fmt.Errorf("tx %s reverted,error reason: %s", tx.Hash().Hex(), errReason)
					}
					return blockNumber, fmt.Errorf("tx %s reverted, could not decode error reason", tx.Hash().Hex())
				}
				return blockNumber, err
			},
		}
	}
	return chains, nil
}
