package devenv

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gagliardetto/solana-go"
	solRpc "github.com/gagliardetto/solana-go/rpc"
	chainselectors "github.com/smartcontractkit/chain-selectors"
	"golang.org/x/sync/errgroup"

	cldf_chain_utils "github.com/smartcontractkit/chainlink-deployments-framework/chain/utils"

	solCommonUtil "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/common"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	cldf_evm "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
	cldf_solana "github.com/smartcontractkit/chainlink-deployments-framework/chain/solana"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment"
)

const (
	EVMChainType = "EVM"
	SolChainType = "SOLANA"
)

type CribRPCs struct {
	Internal string
	External string
}

// ChainConfig holds the configuration for a with a deployer key which can be used to send transactions to the chain.
type ChainConfig struct {
	ChainID            string                   // chain id as per EIP-155
	ChainName          string                   // name of the chain populated from chainselector repo
	ChainType          string                   // should denote the chain family. Acceptable values are EVM, COSMOS, SOLANA, STARKNET, APTOS etc
	PreferredURLScheme cldf.URLSchemePreference // preferred url scheme for the chain
	WSRPCs             []CribRPCs               // websocket rpcs to connect to the chain
	HTTPRPCs           []CribRPCs               // http rpcs to connect to the chain
	DeployerKey        *bind.TransactOpts       // key to deploy and configure contracts on the chain
	SolDeployerKey     solana.PrivateKey
	Users              []*bind.TransactOpts        // map of addresses to their transact opts to interact with the chain as users
	MultiClientOpts    []func(c *cldf.MultiClient) // options to configure the multi client
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
		chainID, success := new(big.Int).SetString(c.ChainID, 10)
		if !success {
			return fmt.Errorf("invalid chainID %s", c.ChainID)
		}
		user, err := bind.NewKeyedTransactorWithChainID(pvtKey, chainID)
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
		chainID, success := new(big.Int).SetString(c.ChainID, 10)
		if !success {
			return fmt.Errorf("invalid chainID %s", c.ChainID)
		}

		deployer, err := bind.NewKeyedTransactorWithChainID(pvtKey, chainID)
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
	chainID, success := new(big.Int).SetString(c.ChainID, 10)
	if !success {
		return fmt.Errorf("invalid chainID %s", c.ChainID)
	}
	c.DeployerKey, err = evmKMSClient.GetKMSTransactOpts(context.Background(), chainID)
	if err != nil {
		return fmt.Errorf("failed to get transactor from KMS client: %w", err)
	}
	return nil
}

func (c *ChainConfig) ToRPCs() []cldf.RPC {
	var rpcs []cldf.RPC
	// assuming that the length of WSRPCs and HTTPRPCs is always the same
	for i, rpc := range c.WSRPCs {
		rpcs = append(rpcs, cldf.RPC{
			Name:               fmt.Sprintf("%s-%d", c.ChainName, i),
			WSURL:              rpc.External,
			HTTPURL:            c.HTTPRPCs[i].External, // copying the corresponding HTTP RPC
			PreferredURLScheme: c.PreferredURLScheme,
		})
	}
	return rpcs
}

func NewChains(logger logger.Logger, configs []ChainConfig) (map[uint64]cldf_evm.Chain, map[uint64]cldf_solana.Chain, error) {
	evmChains := make(map[uint64]cldf_evm.Chain)
	solChains := make(map[uint64]cldf_solana.Chain)
	var evmSyncMap sync.Map
	var solSyncMap sync.Map

	g := new(errgroup.Group)
	for _, chainCfg := range configs {
		chainCfg := chainCfg // capture loop variable
		g.Go(func() error {
			chainDetails, err := chainselectors.GetChainDetailsByChainIDAndFamily(chainCfg.ChainID, strings.ToLower(chainCfg.ChainType))
			if err != nil {
				return fmt.Errorf("failed to get selector from chain id %s: %w", chainCfg.ChainID, err)
			}

			rpcConf := cldf.RPCConfig{
				ChainSelector: chainDetails.ChainSelector,
				RPCs:          chainCfg.ToRPCs(),
			}

			switch chainCfg.ChainType {
			case EVMChainType:
				ec, err := cldf.NewMultiClient(logger, rpcConf, chainCfg.MultiClientOpts...)
				if err != nil {
					return fmt.Errorf("failed to create multi client: %w", err)
				}

				chainInfo, err := cldf_chain_utils.ChainInfo(chainDetails.ChainSelector)
				if err != nil {
					return fmt.Errorf("failed to get chain info for chain %s: %w", chainCfg.ChainName, err)
				}

				confirmFn := func(tx *types.Transaction) (uint64, error) {
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
				}

				evmSyncMap.Store(chainDetails.ChainSelector, cldf_evm.Chain{
					Selector:    chainDetails.ChainSelector,
					Client:      ec,
					DeployerKey: chainCfg.DeployerKey,
					Confirm:     confirmFn,
				})
				return nil

			case SolChainType:
				logger.Info("Creating solana programs tmp dir")
				programsPath, err := os.MkdirTemp("", "solana-programs")
				logger.Infof("Solana programs tmp dir at %s", programsPath)
				if err != nil {
					return err
				}

				keyPairDir, err := os.MkdirTemp("", "solana-keypair")
				logger.Infof("Solana keypair dir at %s", keyPairDir)
				if err != nil {
					return err
				}

				keyPairPath, err := generateSolanaKeypair(chainCfg.SolDeployerKey, keyPairDir)
				if err != nil {
					return err
				}

				sc := solRpc.New(chainCfg.HTTPRPCs[0].External)
				solSyncMap.Store(chainDetails.ChainSelector, cldf_solana.Chain{
					Selector:    chainDetails.ChainSelector,
					Client:      sc,
					DeployerKey: &chainCfg.SolDeployerKey,
					KeypairPath: keyPairPath,
					URL:         chainCfg.HTTPRPCs[0].External,
					WSURL:       chainCfg.WSRPCs[0].External,
					Confirm: func(instructions []solana.Instruction, opts ...solCommonUtil.TxModifier) error {
						_, err := solCommonUtil.SendAndConfirm(
							context.Background(), sc, instructions, chainCfg.SolDeployerKey, solRpc.CommitmentConfirmed, opts...,
						)
						return err
					},
					ProgramsPath: programsPath,
				})
				return nil

			default:
				return fmt.Errorf("chain type %s is not supported", chainCfg.ChainType)
			}
		})
	}

	if err := g.Wait(); err != nil {
		return nil, nil, err
	}

	evmSyncMap.Range(func(sel, value interface{}) bool {
		evmChains[sel.(uint64)] = value.(cldf_evm.Chain)
		return true
	})

	solSyncMap.Range(func(sel, value interface{}) bool {
		solChains[sel.(uint64)] = value.(cldf_solana.Chain)
		return true
	})

	return evmChains, solChains, nil
}

func (c *ChainConfig) SetSolDeployerKey(keyString *string) error {
	if keyString == nil || *keyString == "" {
		return errors.New("no Solana private key provided")
	}

	solKey, err := solana.PrivateKeyFromBase58(*keyString)
	if err != nil {
		return fmt.Errorf("invalid Solana private key: %w", err)
	}

	c.SolDeployerKey = solKey
	return nil
}

func generateSolanaKeypair(privateKey solana.PrivateKey, dir string) (string, error) {
	privateKeyBytes := []byte(privateKey)

	intArray := make([]int, len(privateKeyBytes))
	for i, b := range privateKeyBytes {
		intArray[i] = int(b)
	}

	keypairJSON, err := json.Marshal(intArray)
	if err != nil {
		return "", fmt.Errorf("failed to marshal keypair: %w", err)
	}

	keypairPath := filepath.Join(dir, "solana-keypair.json")
	if err := os.WriteFile(keypairPath, keypairJSON, 0600); err != nil {
		return "", fmt.Errorf("failed to write keypair to file: %w", err)
	}

	return keypairPath, nil
}
