package devenv

import (
	"context"
	"fmt"
	"math/big"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/sethvargo/go-retry"
	chainselectors "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink/integration-tests/deployment"
)

const (
	KMS_DEPLOYER_KEY_ENV        = "KMS_DEPLOYER_KEY_ID"
	KMS_DEPLOYER_KEY_REGION_ENV = "KMS_DEPLOYER_KEY_REGION"
	AWS_PROFILE_NAME_ENV        = "AWS_PROFILE"
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

func (c *ChainConfig) SetDeployerKey(pvtKeyStr *string) error {
	if pvtKeyStr != nil {
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
	kmsDeployerKeyId, exists := os.LookupEnv("KMS_DEPLOYER_KEY_ID")
	if !exists {
		return fmt.Errorf("KMS_DEPLOYER_KEY_ID is required")
	}
	kmsDeployerKeyRegion, exists := os.LookupEnv("KMS_DEPLOYER_KEY_REGION")
	if !exists {
		return fmt.Errorf("KMS_DEPLOYER_KEY_REGION is required")
	}
	awsProfileName, exists := os.LookupEnv("AWS_PROFILE")
	if !exists {
		return fmt.Errorf("AWS_PROFILE is required")
	}
	kmsClient, err := deployment.NewKMSClient(deployment.KMS{
		KmsDeployerKeyId:     kmsDeployerKeyId,
		KmsDeployerKeyRegion: kmsDeployerKeyRegion,
		AwsProfileName:       awsProfileName,
	})
	if err != nil {
		return fmt.Errorf("failed to create KMS client: %w", err)
	}
	evmKMSClient := deployment.NewEVMKMSClient(kmsClient, kmsDeployerKeyId)
	c.DeployerKey, err = evmKMSClient.GetKMSTransactOpts(context.Background(), new(big.Int).SetUint64(c.ChainID))
	if err != nil {
		return fmt.Errorf("failed to get transactor from KMS client: %w", err)
	}
	return nil
}

func NewChains(configs []ChainConfig) (map[uint64]deployment.Chain, error) {
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
		mc, err := deployment.NewMultiClient(rpcs)
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
				err := retry.Do(context.Background(),
					retry.WithMaxDuration(3*time.Minute, retry.NewFibonacci(1*time.Second)),
					func(ctx context.Context) error {
						receipt, err := bind.WaitMined(ctx, mc, tx)
						if err != nil {
							return retry.RetryableError(fmt.Errorf("failed to get receipt for chain %d: %w", chainID, err))
						}
						if receipt != nil {
							blockNumber = receipt.BlockNumber.Uint64()
						}
						if receipt.Status == 0 {
							errReason, err := deployment.GetErrorReasonFromTx(mc, chainCfg.DeployerKey.From, *tx, receipt)
							if err == nil && errReason != "" {
								return fmt.Errorf("tx %s reverted,error reason: %s", tx.Hash().Hex(), errReason)
							}
							return fmt.Errorf("tx %s reverted, could not decode error reason", tx.Hash().Hex())
						}
						return nil
					})
				return blockNumber, err
			},
		}
	}
	return chains, nil
}
