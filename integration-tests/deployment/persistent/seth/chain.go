package seth

import (
	"context"
	"fmt"

	"github.com/avast/retry-go/v4"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pelletier/go-toml/v2"
	"github.com/pkg/errors"

	"path/filepath"
	"strings"

	chainselectors "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/blockchain"
	ctf_config "github.com/smartcontractkit/chainlink-testing-framework/lib/config"
	ctf_test_env "github.com/smartcontractkit/chainlink-testing-framework/lib/docker/test_env"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/osutil"
	seth_utils "github.com/smartcontractkit/chainlink-testing-framework/lib/utils/seth"
	"github.com/smartcontractkit/chainlink-testing-framework/seth"
	"github.com/smartcontractkit/chainlink/integration-tests/deployment"
	persistent_types "github.com/smartcontractkit/chainlink/integration-tests/deployment/persistent/types"
)

func CreateNewEVMChainWithSeth(config ctf_config.EthereumNetworkConfig, sethConfig seth.Config, hooks ctf_test_env.EthereumNetworkHooks) (persistent_types.NewEVMChainProducer, error) {
	contractsRootFolder, err := findGethWrappersFolderRoot(5)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to find contracts root folder")
	}

	return &NewEVMChainWithSeth{
		config:               config,
		sethConfig:           sethConfig,
		contractsRootFolder:  contractsRootFolder,
		EthereumNetworkHooks: hooks,
	}, nil
}

type NewEVMChainWithSeth struct {
	ChainBuilder
	config              ctf_config.EthereumNetworkConfig
	sethConfig          seth.Config
	contractsRootFolder string
	ctf_test_env.EthereumNetworkHooks
}

func (n *NewEVMChainWithSeth) Hooks() ctf_test_env.EthereumNetworkHooks {
	return n.EthereumNetworkHooks
}

func (n *NewEVMChainWithSeth) Chain() (deployment.Chain, persistent_types.RpcProvider, error) {
	chain := deployment.Chain{}

	ethBuilder := ctf_test_env.NewEthereumNetworkBuilder()
	network, err := ethBuilder.
		WithExistingConfig(n.config).
		WithHooks(n.EthereumNetworkHooks).
		Build()

	if err != nil {
		return chain, nil, err
	}

	evmNetwork, rpcProvider, err := network.Start()
	if err != nil {
		return chain, nil, err
	}

	evmNetwork.Name = fmt.Sprintf("%s-%d", *network.ExecutionLayer, evmNetwork.ChainID)

	finalSethConfig, err := prepareFinalSethConfig(n.sethConfig, evmNetwork)
	if err != nil {
		return chain, nil, errors.Wrapf(err, "failed to prepare final seth config")
	}

	sethClient, err := seth.NewClientBuilderWithConfig(finalSethConfig).
		// we want to set it dynamically, because the path depends on the location of the file in the project
		WithGethWrappersFolders([]string{fmt.Sprintf("%s/ccip", n.contractsRootFolder)}).
		WithRpcUrl(evmNetwork.URLs[0]).
		WithPrivateKeys(evmNetwork.PrivateKeys).
		Build()

	if err != nil {
		return chain, nil, errors.Wrapf(err, "failed to create seth client")
	}

	return n.Build(sethClient, evmNetwork, rpcProvider)
}

type NewEVMChainConfigWithSeth struct {
	ctf_config.EthereumNetworkConfig
	sethConfig seth.Config
}

func (n *NewEVMChainConfigWithSeth) SethConfig() seth.Config {
	return n.sethConfig
}

func (n *NewEVMChainConfigWithSeth) DockerNetworks() []string {
	var dockerNetworks []string
	for _, network := range n.DockerNetworkNames {
		contains := false
		for _, dockerNetwork := range dockerNetworks {
			if strings.EqualFold(dockerNetwork, network) {
				contains = true
				break
			}
		}
		if !contains {
			dockerNetworks = append(dockerNetworks, network)
		}
	}
	return dockerNetworks
}

type ExistingEVMChainConfigWithSeth struct {
	ChainBuilder
	evmNetwork          blockchain.EVMNetwork
	sethConfig          seth.Config
	contractsRootFolder string
}

func (e *ExistingEVMChainConfigWithSeth) EVMNetwork() blockchain.EVMNetwork {
	return e.evmNetwork
}

func CreateExistingEVMChainWithSeth(evmNetwork blockchain.EVMNetwork, sethConfig seth.Config) (persistent_types.ExistingEVMChainProducer, error) {
	contractsRootFolder, err := findGethWrappersFolderRoot(5)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to find contracts root folder")
	}

	return &ExistingEVMChainConfigWithSeth{
		evmNetwork:          evmNetwork,
		sethConfig:          sethConfig,
		contractsRootFolder: contractsRootFolder,
	}, nil
}

func (e *ExistingEVMChainConfigWithSeth) Chain() (deployment.Chain, persistent_types.RpcProvider, error) {
	chain := deployment.Chain{}

	finalSethConfig, err := prepareFinalSethConfig(e.sethConfig, e.evmNetwork)
	if err != nil {
		return chain, nil, errors.Wrapf(err, "failed to prepare final seth config")
	}

	sethClient, err := seth.NewClientBuilderWithConfig(finalSethConfig).
		// we want to set it dynamically, because the path depends on the location of the file in the project
		WithGethWrappersFolders([]string{fmt.Sprintf("%s/ccip", e.contractsRootFolder)}).
		Build()
	if err != nil {
		return chain, nil, errors.Wrapf(err, "failed to create seth client")
	}

	return e.Build(sethClient, e.evmNetwork, ctf_test_env.NewRPCProvider(e.evmNetwork.HTTPURLs, e.evmNetwork.URLs, e.evmNetwork.HTTPURLs, e.evmNetwork.URLs))
}

type ChainBuilder struct{}

func (s *ChainBuilder) Build(sethClient *seth.Client, evmNetwork blockchain.EVMNetwork, rpcProvider ctf_test_env.RpcProvider) (deployment.Chain, persistent_types.RpcProvider, error) {
	shouldRetryOnErrFn := func(err error) bool {
		// some retry logic here
		return true
	}

	prepareReplacementTransactionFn := func(sethClient *seth.Client, tx *types.Transaction) (*types.Transaction, error) {
		// TODO some replacement tx creation logic could go here
		// TODO for example: adjusting base fee aggressively if it's too low for transaction to even be included in the block
		return tx, nil
	}

	sel, err := chainselectors.SelectorFromChainId(uint64(evmNetwork.ChainID))
	if err != nil {
		return deployment.Chain{}, nil, err
	}

	return deployment.Chain{
		Selector: sel,
		Client:   sethClient.Client,
		Keys: func() []*bind.TransactOpts {
			var keys []*bind.TransactOpts
			// use all private keys set for network, in case we want to use them for concurrent transactions
			for i := range sethClient.Cfg.Network.PrivateKeys {
				// we set the nonce to nil, because we want go-ethereum to use pending nonce it gets from the node
				opts := sethClient.NewTXKeyOpts(i, seth.WithNonce(nil))
				keys = append(keys, opts)
			}

			return keys
		}(),
		Confirm: func(txHash common.Hash) (uint64, error) {
			ctx, cancelFn := context.WithTimeout(context.Background(), sethClient.Cfg.Network.TxnTimeout.Duration())
			tx, _, err := sethClient.Client.TransactionByHash(ctx, txHash)
			cancelFn()
			if err != nil {
				return 0, err
			}
			decoded, revertErr := sethClient.DecodeTx(tx)
			if revertErr != nil {
				return 0, revertErr
			}
			if decoded.Receipt == nil {
				return 0, fmt.Errorf("no receipt found for transaction %s even though it wasn't reverted. This should not happen", tx.Hash().String())
			}
			return decoded.Receipt.BlockNumber.Uint64(), nil
		},
		RetrySubmit: func(tx *types.Transaction, err error) (*types.Transaction, error) {
			if err == nil {
				return tx, nil
			}

			retryErr := retry.Do(
				func() error {
					ctx, cancel := context.WithTimeout(context.Background(), sethClient.Cfg.Network.TxnTimeout.Duration())
					defer cancel()

					return sethClient.Client.SendTransaction(ctx, tx)
				}, retry.OnRetry(func(i uint, retryErr error) {
					replacementTx, replacementErr := prepareReplacementTransactionFn(sethClient, tx)
					if replacementErr != nil {
						return
					}
					tx = replacementTx
				}),
				retry.DelayType(retry.FixedDelay),
				retry.Attempts(10),
				retry.RetryIf(shouldRetryOnErrFn),
			)

			return tx, sethClient.DecodeSendErr(retryErr)
		},
		DefaultKey: func() *bind.TransactOpts {
			// this will use the first private key from the seth client
			// if you want to use N private key you can use sethClient.NewTXKeyOpts(N)
			// we set the nonce to nil, because we want go-ethereum to use pending nonce it gets from the node
			return sethClient.NewTXOpts(seth.WithNonce(nil))
		},
	}, persistent_types.NewEVMNetworkWithRPCs(evmNetwork, rpcProvider), nil
}

// findGethWrappersFolderRoot finds the root folder of the geth wrappers. It looks for a file named ".geth_wrappers_root" or ".repo_root" in the current directory and its `folderLimit` parents.
func findGethWrappersFolderRoot(folderLimit int) (string, error) {
	contractsRootFile, err := osutil.FindFile(".geth_wrappers_root", ".repo_root", folderLimit)
	if err != nil {
		return "", fmt.Errorf("failed to find contracts root folder: %w", err)
	}
	return filepath.Dir(contractsRootFile), nil
}

func prepareFinalSethConfig(sethConfig seth.Config, evmNetwork blockchain.EVMNetwork) (*seth.Config, error) {
	// copy it so we don't end up modifying the original config during configuration merge
	marshalled, err := toml.Marshal(sethConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal seth config: %w", err)
	}
	var sethConfigCopy seth.Config
	err = toml.Unmarshal(marshalled, &sethConfigCopy)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal seth config: %w", err)
	}

	sethConfig, err = seth_utils.MergeSethAndEvmNetworkConfigs(evmNetwork, sethConfigCopy)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to merge seth and evm network configs")
	}

	return &sethConfig, nil
}
