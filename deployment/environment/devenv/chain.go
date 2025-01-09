package devenv

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	chainselectors "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/chainlink/deployment"
)

const (
	EVMChainType = "EVM"
)

// ChainConfig holds the configuration for a with a deployer key which can be used to send transactions to the chain.
type ChainConfig struct {
	ChainID     uint64               // chain id as per EIP-155, mainly applicable for EVM chains
	ChainName   string               // name of the chain populated from chainselector repo
	ChainType   string               // should denote the chain family. Acceptable values are EVM, COSMOS, SOLANA, STARKNET, APTOS etc
	WSRPCs      []string             // websocket rpcs to connect to the chain
	HTTPRPCs    []string             // http rpcs to connect to the chain
	DeployerKey *bind.TransactOpts   // key to deploy and configure contracts on the chain
	Users       []*bind.TransactOpts // map of addresses to their transact opts to interact with the chain as users
}

func (c *ChainConfig) SetUsers(pvtkeys []string) error {
	if pvtkeys == nil {
		// if no private keys are provided, set deployer key as the user
		if c.DeployerKey != nil {
			c.Users = []*bind.TransactOpts{c.DeployerKey}
			return nil
		} else {
			return errors.New("no private keys provided for users, deployer key is also not set")
		}
	}
	for _, pvtKeyStr := range pvtkeys {
		pvtKey, err := crypto.HexToECDSA(pvtKeyStr)
		if err != nil {
			return fmt.Errorf("failed to convert private key to ECDSA: %w", err)
		}
		user, err := bind.NewKeyedTransactorWithChainID(pvtKey, new(big.Int).SetUint64(c.ChainID))
		if err != nil {
			return fmt.Errorf("failed to create transactor: %w", err)
		}
		c.Users = append(c.Users, user)
	}
	return nil
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
		fmt.Printf("Deployer Address: %s for chain id %d\n", deployer.From.Hex(), c.ChainID)
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
		chainInfo, err := deployment.ChainInfo(selector)
		if err != nil {
			return nil, fmt.Errorf("failed to get chain info for chain %s: %w", chainCfg.ChainName, err)
		}
		chains[selector] = deployment.Chain{
			Selector:    selector,
			Client:      ec,
			DeployerKey: chainCfg.DeployerKey,
			Confirm: func(tx *types.Transaction) (uint64, error) {
				var blockNumber uint64
				if tx == nil {
					return 0, fmt.Errorf("tx was nil, nothing to confirm chain %s", chainInfo.ChainName)
				}
				ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
				defer cancel()
				receipt, err := bind.WaitMined(ctx, ec, tx)
				if err != nil {
					return blockNumber, fmt.Errorf("failed to get confirmed receipt for chain %s: %w", chainInfo.ChainName, err)
				}
				if receipt == nil {
					return blockNumber, fmt.Errorf("receipt was nil for tx %s chain %s", tx.Hash().Hex(), chainInfo.ChainName)
				}
				blockNumber = receipt.BlockNumber.Uint64()
				if receipt.Status == 0 {
					errReason, err := deployment.GetErrorReasonFromTx(ec, chainCfg.DeployerKey.From, tx, receipt)
					if err == nil && errReason != "" {
						return blockNumber, fmt.Errorf("tx %s reverted,error reason: %s chain %s", tx.Hash().Hex(), errReason, chainInfo.ChainName)
					}
					return blockNumber, fmt.Errorf("tx %s reverted, could not decode error reason chain %s", tx.Hash().Hex(), chainInfo.ChainName)
				}
				return blockNumber, nil
			},
		}
	}
	return chains, nil
}
