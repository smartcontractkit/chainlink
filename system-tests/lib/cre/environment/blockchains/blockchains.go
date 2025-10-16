package blockchains

import (
	"context"
	"fmt"

	pkgerrors "github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain"
	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"

	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
)

type Blockchain interface {
	ChainSelector() uint64
	ChainID() uint64
	ChainFamily() string
	Is(chainFamily string) bool

	CtfOutput() *blockchain.Output
	Fund(ctx context.Context, address string, amount uint64) error
	ToCldfChain() (cldf_chain.BlockChain, error)
}

type Deployer interface {
	Deploy(input *blockchain.Input) (Blockchain, error)
}

type DeployedBlockchains struct {
	Outputs         []Blockchain
	CldfBlockChains chain.BlockChains
}

func (s *DeployedBlockchains) RegistryChain() Blockchain {
	return s.Outputs[0]
}

func Start(
	commonLogger logger.Logger,
	inputs []*blockchain.Input,
	deployers map[blockchain.ChainFamily]Deployer,
) (*DeployedBlockchains, error) {
	outputs := make([]Blockchain, 0, len(inputs))

	for _, input := range inputs {
		chainFamily, chErr := blockchain.TypeToFamily(input.Type)
		if chErr != nil {
			return nil, chErr
		}

		deployer, ok := deployers[chainFamily]
		if !ok {
			return nil, fmt.Errorf("no deployer found for blockchain type %s", input.Type)
		}

		deployedBlockchain, deployErr := deployer.Deploy(input)
		if deployErr != nil {
			return nil, pkgerrors.Wrapf(deployErr, "failed to deploy blockchain of type %s", input.Type)
		}

		outputs = append(outputs, deployedBlockchain)
	}

	cldfBlockchains := make([]cldf_chain.BlockChain, 0, len(outputs))
	for _, db := range outputs {
		// chain, chainErr := CldfChainConfigFromBlockchain(db)
		chain, chainErr := db.ToCldfChain()
		if chainErr != nil {
			return nil, pkgerrors.Wrap(chainErr, "failed to create cldf chain from blockchain")
		}
		cldfBlockchains = append(cldfBlockchains, chain)
	}

	// cldfBlockchains, err := NewCldfChains(commonLogger, chainsConfigs)
	// if err != nil {
	// 	return nil, pkgerrors.Wrap(err, "failed to create chains")
	// }

	return &DeployedBlockchains{
		Outputs:         outputs,
		CldfBlockChains: cldf_chain.NewBlockChainsFromSlice(cldfBlockchains),
	}, nil
}

// type RPCs struct {
// 	Internal string // URL to be used by services running in the same namespace
// 	External string // URL to be used when connecting from outside the namespace
// }

// // CldfChainConfig holds the configuration for a with a deployer key which can be used to send transactions to the chain.
// type CldfChainConfig struct {
// 	ChainID            string                              // chain id as per EIP-155
// 	ChainName          string                              // name of the chain populated from chainselector repo
// 	ChainType          string                              // should denote the chain family. Acceptable values are EVM, COSMOS, SOLANA, STARKNET, APTOS etc
// 	PreferredURLScheme cldf_evm_client.URLSchemePreference // preferred url scheme for the chain
// 	WSRPCs             []RPCs                              // websocket rpcs to connect to the chain
// 	HTTPRPCs           []RPCs                              // http rpcs to connect to the chain
// 	DeployerKey        *bind.TransactOpts                  // key to deploy and configure contracts on the chain
// 	SolDeployerKey     solanago.PrivateKey
// 	SolArtifactDir     string                                 // directory of pre-built solana artifacts, if any
// 	Users              []*bind.TransactOpts                   // map of addresses to their transact opts to interact with the chain as users
// 	MultiClientOpts    []func(c *cldf_evm_client.MultiClient) // options to configure the multi client
// 	AptosDeployerKey   aptos.Account
// }

// func (c *CldfChainConfig) ToRPCs() []cldf_evm_client.RPC {
// 	rpcs := []cldf_evm_client.RPC{}
// 	// assuming that the length of WSRPCs and HTTPRPCs is always the same
// 	for i, rpc := range c.WSRPCs {
// 		rpcs = append(rpcs, cldf_evm_client.RPC{
// 			Name:               fmt.Sprintf("%s-%d", c.ChainName, i),
// 			WSURL:              rpc.External,
// 			HTTPURL:            c.HTTPRPCs[i].External, // copying the corresponding HTTP RPC
// 			PreferredURLScheme: c.PreferredURLScheme,
// 		})
// 	}
// 	return rpcs
// }

// // copied from deployment/environment/devenv/chain.go to avoid dependency on deployment module
// func NewCldfChains(logger logger.Logger, blockchains []Blockchain) (cldf_chain.BlockChains, error) {
// 	var evmSyncMap sync.Map
// 	var solSyncMap sync.Map
// 	var aptosSyncMap sync.Map
// 	var tronSyncMap sync.Map

// 	g := new(errgroup.Group)
// 	for _, chainCfg := range blockchains {
// 		g.Go(func() error {
// 			family := strings.ToLower(chainCfg.ChainFamily())
// 			// tron's devnet chainID maps to many chain selectors, one for tron one for EVM
// 			// we want to force mapping to EVM family here to avoid selector mismatches later
// 			if chainCfg.Is(chainselectors.FamilyTron) {
// 				family = chainselectors.FamilyEVM
// 			}
// 			chainDetails, err := chainselectors.GetChainDetailsByChainIDAndFamily(strconv.FormatUint(chainCfg.ChainID(), 10), family)
// 			if err != nil {
// 				return fmt.Errorf("failed to get selector from chain id %s: %w", chainCfg.ChainID(), err)
// 			}

// 			rpcConf := cldf_evm_client.RPCConfig{
// 				ChainSelector: chainDetails.ChainSelector,
// 				RPCs:          chainCfg.ToRPCs(),
// 			}

// 			switch strings.ToLower(chainCfg.ChainFamily()) {
// 			case chainselectors.FamilyEVM:
// 				ec, evmErr := cldf_evm_client.NewMultiClient(logger, rpcConf)
// 				if evmErr != nil {
// 					return fmt.Errorf("failed to create multi client: %w", evmErr)
// 				}

// 				chainInfo, infoErr := cldf_chain_utils.ChainInfo(chainDetails.ChainSelector)
// 				if infoErr != nil {
// 					return fmt.Errorf("failed to get chain info for chain %s-%d: %w", chainCfg.ChainFamily(), chainCfg.ChainID(), infoErr)
// 				}

// 				confirmFn := func(tx *types.Transaction) (uint64, error) {
// 					var blockNumber uint64
// 					if tx == nil {
// 						return 0, fmt.Errorf("tx was nil, nothing to confirm chain %s", chainInfo.ChainName)
// 					}
// 					ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
// 					defer cancel()
// 					receipt, rErr := bind.WaitMined(ctx, ec, tx)
// 					if rErr != nil {
// 						return blockNumber, fmt.Errorf("failed to get confirmed receipt for chain %s: %w", chainInfo.ChainName, rErr)
// 					}
// 					if receipt == nil {
// 						return blockNumber, fmt.Errorf("receipt was nil for tx %s chain %s", tx.Hash().Hex(), chainInfo.ChainName)
// 					}
// 					blockNumber = receipt.BlockNumber.Uint64()
// 					if receipt.Status == 0 {
// 						errReason, rErr := deployment.GetErrorReasonFromTx(ec, chainCfg, tx, receipt)
// 						if rErr == nil && errReason != "" {
// 							return blockNumber, fmt.Errorf("tx %s reverted,error reason: %s chain %s", tx.Hash().Hex(), errReason, chainInfo.ChainName)
// 						}
// 						return blockNumber, fmt.Errorf("tx %s reverted, could not decode error reason chain %s", tx.Hash().Hex(), chainInfo.ChainName)
// 					}
// 					return blockNumber, nil
// 				}

// 				chain := cldf_evm.Chain{
// 					Selector:    chainDetails.ChainSelector,
// 					Client:      ec,
// 					DeployerKey: chainCfg.DeployerKey,
// 					Confirm:     confirmFn,
// 				}

// 				evmSyncMap.Store(chainDetails.ChainSelector, chain)
// 				return nil

// 			case chainselectors.FamilySolana:
// 				solArtifactPath := chainCfg.SolArtifactDir
// 				if solArtifactPath == "" {
// 					logger.Info("Creating tmp directory for generated solana programs and keypairs")
// 					solArtifactPath, err = os.MkdirTemp("", "solana-artifacts")
// 					logger.Infof("Solana programs tmp dir at %s", solArtifactPath)
// 					if err != nil {
// 						return err
// 					}
// 				}

// 				sc := solRpc.New(chainCfg.HTTPRPCs[0].External)
// 				solSyncMap.Store(chainDetails.ChainSelector, cldf_solana.Chain{
// 					Selector:    chainDetails.ChainSelector,
// 					Client:      sc,
// 					DeployerKey: &chainCfg.SolDeployerKey,
// 					KeypairPath: solArtifactPath + "/deploy-keypair.json",
// 					URL:         chainCfg.HTTPRPCs[0].External,
// 					WSURL:       chainCfg.WSRPCs[0].External,
// 					Confirm: func(instructions []solanago.Instruction, opts ...solCommonUtil.TxModifier) error {
// 						_, err := solCommonUtil.SendAndConfirm(
// 							context.Background(), sc, instructions, chainCfg.SolDeployerKey, solRpc.CommitmentConfirmed, opts...,
// 						)
// 						return err
// 					},
// 					ProgramsPath: solArtifactPath,
// 				})
// 				return nil
// 			case chainselectors.FamilyTron:
// 				signerGen, err := tronprovider.SignerGenCTFDefault()
// 				if err != nil {
// 					return fmt.Errorf("failed to create signer generator: %w", err)
// 				}

// 				fullNodeURL := strings.Replace(chainCfg.HTTPRPCs[0].External, "/jsonrpc", "/wallet", 1)
// 				solidityNodeURL := strings.Replace(chainCfg.HTTPRPCs[0].External, "/jsonrpc", "/walletsolidity", 1)

// 				tronRPCProvider := tronprovider.NewRPCChainProvider(chainDetails.ChainSelector, tronprovider.RPCChainProviderConfig{
// 					FullNodeURL:       fullNodeURL,
// 					SolidityNodeURL:   solidityNodeURL,
// 					DeployerSignerGen: signerGen,
// 				})
// 				ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
// 				defer cancel()
// 				tronChain, err := tronRPCProvider.Initialize(ctx)
// 				if err != nil {
// 					return fmt.Errorf("failed to initialize tron chain: %w", err)
// 				}

// 				tronChain, ok := tronChain.(cldf_tron.Chain)
// 				if !ok {
// 					return fmt.Errorf("expected cldf_tron.Chain, got %T", tronChain)
// 				}

// 				tronSyncMap.Store(chainDetails.ChainSelector, tronChain)
// 				return nil
// 			default:
// 				return fmt.Errorf("chain type %s is not supported", chainCfg.ChainType)
// 			}
// 		})
// 	}

// 	if err := g.Wait(); err != nil {
// 		return cldf_chain.BlockChains{}, err
// 	}

// 	var blockChains []cldf_chain.BlockChain

// 	evmSyncMap.Range(func(sel, value any) bool {
// 		blockChains = append(blockChains, value.(cldf_evm.Chain))
// 		return true
// 	})

// 	solSyncMap.Range(func(sel, value any) bool {
// 		blockChains = append(blockChains, value.(cldf_solana.Chain))
// 		return true
// 	})

// 	aptosSyncMap.Range(func(sel, value any) bool {
// 		blockChains = append(blockChains, value.(cldf_aptos.Chain))
// 		return true
// 	})

// 	tronSyncMap.Range(func(sel, value any) bool {
// 		blockChains = append(blockChains, value.(cldf_tron.Chain))
// 		return true
// 	})

// 	return cldf_chain.NewBlockChainsFromSlice(blockChains), nil
// }
