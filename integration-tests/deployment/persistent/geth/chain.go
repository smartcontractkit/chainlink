package geth

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	chainselectors "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	ctf_config "github.com/smartcontractkit/chainlink-testing-framework/config"
	ctf_test_env "github.com/smartcontractkit/chainlink-testing-framework/docker/test_env"
	"github.com/smartcontractkit/chainlink/integration-tests/deployment"
	persistent_types "github.com/smartcontractkit/chainlink/integration-tests/deployment/persistent/types"
	"time"
)

type GethChainBuilder struct{}

func (s *GethChainBuilder) Build(evmNetwork blockchain.EVMNetwork, rpcProvider ctf_test_env.RpcProvider) (deployment.Chain, error) {
	chain := deployment.Chain{}
	client, err := ethclient.Dial(evmNetwork.URLs[0])
	if err != nil {
		return chain, err
	}

	sel, err := chainselectors.SelectorFromChainId(uint64(evmNetwork.ChainID))
	if err != nil {
		return deployment.Chain{}, err
	}

	return deployment.Chain{
		Selector:           sel,
		Client:             client,
		DeployerKey:        &bind.TransactOpts{},
		DeployerKeys:       []*bind.TransactOpts{{}},
		EVMNetworkWithRPCs: deployment.NewEVMNetworkWithRPCs(evmNetwork, rpcProvider),
		Confirm: func(txHash common.Hash) (uint64, error) {
			ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Minute)
			receipt, err := client.TransactionReceipt(ctx, txHash)
			cancelFn()
			if err != nil {
				return 0, err
			}
			if receipt.Status != types.ReceiptStatusSuccessful {
				return 0, fmt.Errorf("transaction %s failed with status %d", txHash.Hex(), receipt.Status)
			}
			return receipt.BlockNumber.Uint64(), nil
		},
		RetrySubmit: deployment.NoOpRetrySubmit,
	}, nil
}

type NewEVMChainWithGeth struct {
	GethChainBuilder
	config ctf_config.PrivateEthereumNetworkConfig
}

func (n *NewEVMChainWithGeth) Chain() (deployment.Chain, error) {
	chain := deployment.Chain{}
	if n.config.GetEthereumVersion() == nil {
		return chain, fmt.Errorf("ethereum version is required")
	}

	if n.config.GetExecutionLayer() == nil {
		return chain, fmt.Errorf("execution layer is required")
	}

	ethBuilder := ctf_test_env.NewEthereumNetworkBuilder()
	network, err := ethBuilder.
		WithEthereumVersion(*n.config.GetEthereumVersion()).
		WithExecutionLayer(*n.config.GetExecutionLayer()).
		WithEthereumChainConfig(n.config.GetChainConfig()).
		WithDockerNetworks(n.config.GetDockerNetworkNames()).
		WithCustomDockerImages(n.config.GetCustomDockerImages()).
		Build()

	if err != nil {
		return chain, err
	}

	evmNetwork, rpcProvider, err := network.Start()
	if err != nil {
		return chain, err
	}

	evmNetwork.Name = fmt.Sprintf("%s-%d", *n.config.GetExecutionLayer(), evmNetwork.ChainID)

	return n.Build(evmNetwork, rpcProvider)
}

func CreateNewEVMChainWithGeth(config ctf_config.PrivateEthereumNetworkConfig) persistent_types.NewEVMChainConfig {
	return &NewEVMChainWithGeth{
		config: config,
	}
}

func CreateExistingEVMChainConfigWithGeth(evmNetwork blockchain.EVMNetwork) persistent_types.ExistingEVMChainConfig {
	return &ExistingEVMChainConfigWithGeth{
		evmNetwork: evmNetwork,
	}
}

type ExistingEVMChainConfigWithGeth struct {
	GethChainBuilder
	evmNetwork blockchain.EVMNetwork
}

func (e *ExistingEVMChainConfigWithGeth) EVMNetwork() blockchain.EVMNetwork {
	return e.evmNetwork
}

func (e *ExistingEVMChainConfigWithGeth) Chain() (deployment.Chain, error) {
	chain := deployment.Chain{}
	rpcProvider := ctf_test_env.NewRPCProvider(e.evmNetwork.HTTPURLs, e.evmNetwork.URLs, e.evmNetwork.HTTPURLs, e.evmNetwork.URLs)

	chain, err := e.Build(e.evmNetwork, rpcProvider)
	if err != nil {
		return chain, err
	}

	return chain, nil
}
