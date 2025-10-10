package devenv

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/aptos-labs/aptos-go-sdk"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gagliardetto/solana-go"
	solRpc "github.com/gagliardetto/solana-go/rpc"
	"golang.org/x/sync/errgroup"

	aptosCrypto "github.com/aptos-labs/aptos-go-sdk/crypto"
	chainselectors "github.com/smartcontractkit/chain-selectors"
	"github.com/zksync-sdk/zksync2-go/accounts"
	"github.com/zksync-sdk/zksync2-go/clients"

	solCommonUtil "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/common"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	cldf_aptos "github.com/smartcontractkit/chainlink-deployments-framework/chain/aptos"
	cldf_evm "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
	cldf_evm_client "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm/provider/rpcclient"
	cldf_solana "github.com/smartcontractkit/chainlink-deployments-framework/chain/solana"
	cldf_tron "github.com/smartcontractkit/chainlink-deployments-framework/chain/tron"
	tronprovider "github.com/smartcontractkit/chainlink-deployments-framework/chain/tron/provider"
	cldf_chain_utils "github.com/smartcontractkit/chainlink-deployments-framework/chain/utils"
	"github.com/smartcontractkit/chainlink/deployment"
)

const (
	EVMChainType   = "EVM"
	SolChainType   = "SOLANA"
	AptosChainType = "APTOS"
	TronChainType  = "TRON"
)

type CribRPCs struct {
	Internal string
	External string
}

// ChainConfig holds the configuration for a with a deployer key which can be used to send transactions to the chain.
type ChainConfig struct {
	ChainID             string                              // chain id as per EIP-155
	ChainName           string                              // name of the chain populated from chainselector repo
	ChainType           string                              // should denote the chain family. Acceptable values are EVM, COSMOS, SOLANA, STARKNET, APTOS etc
	PreferredURLScheme  cldf_evm_client.URLSchemePreference // preferred url scheme for the chain
	WSRPCs              []CribRPCs                          // websocket rpcs to connect to the chain
	HTTPRPCs            []CribRPCs                          // http rpcs to connect to the chain
	DeployerKey         *bind.TransactOpts                  // key to deploy and configure contracts on the chain
	IsZkSyncVM          bool
	ClientZkSyncVM      *clients.Client
	DeployerKeyZkSyncVM *accounts.Wallet
	SolDeployerKey      solana.PrivateKey
	SolArtifactDir      string                                 // directory of pre-built solana artifacts, if any
	Users               []*bind.TransactOpts                   // map of addresses to their transact opts to interact with the chain as users
	MultiClientOpts     []func(c *cldf_evm_client.MultiClient) // options to configure the multi client
	AptosDeployerKey    aptos.Account
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

func (c *ChainConfig) ToRPCs() []cldf_evm_client.RPC {
	var rpcs []cldf_evm_client.RPC
	// assuming that the length of WSRPCs and HTTPRPCs is always the same
	for i, rpc := range c.WSRPCs {
		rpcs = append(rpcs, cldf_evm_client.RPC{
			Name:               fmt.Sprintf("%s-%d", c.ChainName, i),
			WSURL:              rpc.External,
			HTTPURL:            c.HTTPRPCs[i].External, // copying the corresponding HTTP RPC
			PreferredURLScheme: c.PreferredURLScheme,
		})
	}
	return rpcs
}

func NewChains(logger logger.Logger, configs []ChainConfig) (cldf_chain.BlockChains, error) {
	var evmSyncMap sync.Map
	var solSyncMap sync.Map
	var aptosSyncMap sync.Map
	var tronSyncMap sync.Map

	g := new(errgroup.Group)
	for _, chainCfg := range configs {
		// capture loop variable
		g.Go(func() error {
			family := chainCfg.ChainType
			if chainCfg.ChainType == TronChainType {
				family = EVMChainType
			}
			chainDetails, err := chainselectors.GetChainDetailsByChainIDAndFamily(chainCfg.ChainID, strings.ToLower(family))
			if err != nil {
				return fmt.Errorf("failed to get selector from chain id %s: %w", chainCfg.ChainID, err)
			}

			rpcConf := cldf_evm_client.RPCConfig{
				ChainSelector: chainDetails.ChainSelector,
				RPCs:          chainCfg.ToRPCs(),
			}

			switch chainCfg.ChainType {
			case EVMChainType:
				ec, err := cldf_evm_client.NewMultiClient(logger, rpcConf, chainCfg.MultiClientOpts...)
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

				chain := cldf_evm.Chain{
					Selector:    chainDetails.ChainSelector,
					Client:      ec,
					DeployerKey: chainCfg.DeployerKey,
					Confirm:     confirmFn,
				}

				if chainCfg.IsZkSyncVM {
					chain.IsZkSyncVM = true
					chain.ClientZkSyncVM = chainCfg.ClientZkSyncVM
					chain.DeployerKeyZkSyncVM = chainCfg.DeployerKeyZkSyncVM
				}

				evmSyncMap.Store(chainDetails.ChainSelector, chain)
				return nil

			case SolChainType:
				solArtifactPath := chainCfg.SolArtifactDir
				if solArtifactPath == "" {
					logger.Info("Creating tmp directory for generated solana programs and keypairs")
					solArtifactPath, err = os.MkdirTemp("", "solana-artifacts")
					logger.Infof("Solana programs tmp dir at %s", solArtifactPath)
					if err != nil {
						return err
					}
				}

				sc := solRpc.New(chainCfg.HTTPRPCs[0].External)
				solSyncMap.Store(chainDetails.ChainSelector, cldf_solana.Chain{
					Selector:    chainDetails.ChainSelector,
					Client:      sc,
					DeployerKey: &chainCfg.SolDeployerKey,
					KeypairPath: solArtifactPath + "/deploy-keypair.json",
					URL:         chainCfg.HTTPRPCs[0].External,
					WSURL:       chainCfg.WSRPCs[0].External,
					Confirm: func(instructions []solana.Instruction, opts ...solCommonUtil.TxModifier) error {
						_, err := solCommonUtil.SendAndConfirm(
							context.Background(), sc, instructions, chainCfg.SolDeployerKey, solRpc.CommitmentConfirmed, opts...,
						)
						return err
					},
					ProgramsPath: solArtifactPath,
				})
				return nil

			case AptosChainType:
				cID, err := strconv.ParseUint(chainCfg.ChainID, 10, 8)
				if err != nil {
					return err
				}

				ac, err := aptos.NewClient(aptos.NetworkConfig{
					Name:    chainCfg.ChainName,
					NodeUrl: chainCfg.HTTPRPCs[0].External,
					ChainId: uint8(cID),
				})

				if err != nil {
					return err
				}

				aptosSyncMap.Store(chainDetails.ChainSelector, cldf_aptos.Chain{
					Selector:       chainDetails.ChainSelector,
					Client:         ac,
					DeployerSigner: &chainCfg.AptosDeployerKey,
					URL:            chainCfg.HTTPRPCs[0].External,
					Confirm: func(txHash string, opts ...any) error {
						tx, err := ac.WaitForTransaction(txHash, opts...)
						if err != nil {
							return err
						}

						if !tx.Success {
							return fmt.Errorf("transaction failed: %s", tx.VmStatus)
						}

						return nil
					},
				})
				return nil
			case TronChainType:
				signerGen, err := tronprovider.SignerGenCTFDefault()
				if err != nil {
					return fmt.Errorf("failed to create signer generator: %w", err)
				}

				fullNodeURL := strings.Replace(chainCfg.HTTPRPCs[0].External, "/jsonrpc", "/wallet", 1)
				solidityNodeURL := strings.Replace(chainCfg.HTTPRPCs[0].External, "/jsonrpc", "/walletsolidity", 1)

				tronRPCProvider := tronprovider.NewRPCChainProvider(chainDetails.ChainSelector, tronprovider.RPCChainProviderConfig{
					FullNodeURL:       fullNodeURL,
					SolidityNodeURL:   solidityNodeURL,
					DeployerSignerGen: signerGen,
				})
				ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
				defer cancel()
				tronChain, err := tronRPCProvider.Initialize(ctx)
				if err != nil {
					return fmt.Errorf("failed to initialize tron chain: %w", err)
				}

				tronChain, ok := tronChain.(cldf_tron.Chain)
				if !ok {
					return fmt.Errorf("expected cldf_tron.Chain, got %T", tronChain)
				}

				tronSyncMap.Store(chainDetails.ChainSelector, tronChain)
				return nil
			default:
				return fmt.Errorf("chain type %s is not supported", chainCfg.ChainType)
			}
		})
	}

	if err := g.Wait(); err != nil {
		return cldf_chain.BlockChains{}, err
	}

	var blockChains []cldf_chain.BlockChain

	evmSyncMap.Range(func(sel, value any) bool {
		blockChains = append(blockChains, value.(cldf_evm.Chain))
		return true
	})

	solSyncMap.Range(func(sel, value any) bool {
		blockChains = append(blockChains, value.(cldf_solana.Chain))
		return true
	})

	aptosSyncMap.Range(func(sel, value any) bool {
		blockChains = append(blockChains, value.(cldf_aptos.Chain))
		return true
	})

	tronSyncMap.Range(func(sel, value any) bool {
		blockChains = append(blockChains, value.(cldf_tron.Chain))
		return true
	})

	return cldf_chain.NewBlockChainsFromSlice(blockChains), nil
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

func (c *ChainConfig) SetAptosDeployerKey(keyString *string) error {
	if keyString == nil || *keyString == "" {
		return errors.New("no Aptos private key provided")
	}

	keyStr := strings.TrimPrefix(*keyString, "0x")

	deployerKey := &aptosCrypto.Ed25519PrivateKey{}
	err := deployerKey.FromHex(keyStr)
	if err != nil {
		return fmt.Errorf("invalid Aptos private key: %w", err)
	}

	aptosAccount, err := aptos.NewAccountFromSigner(deployerKey)
	if err != nil {
		return fmt.Errorf("failed to create Aptos account from private key: %w", err)
	}

	c.AptosDeployerKey = *aptosAccount
	return nil
}
