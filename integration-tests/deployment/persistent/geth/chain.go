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

type ChainBuilder struct{}

func (s *ChainBuilder) Build(evmNetwork blockchain.EVMNetwork, rpcProvider ctf_test_env.RpcProvider) (deployment.Chain, persistent_types.RpcProvider, error) {
	chain := deployment.Chain{}
	client, err := ethclient.Dial(evmNetwork.URLs[0])
	if err != nil {
		return chain, nil, err
	}

	sel, err := chainselectors.SelectorFromChainId(uint64(evmNetwork.ChainID))
	if err != nil {
		return chain, nil, err
	}

	return deployment.Chain{
		Selector: sel,
		Client:   client,
		Keys:     []*bind.TransactOpts{{}},
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
		DefaultKey:  func() *bind.TransactOpts { return &bind.TransactOpts{} },
	}, persistent_types.NewEVMNetworkWithRPCs(evmNetwork, rpcProvider), nil
}

type NewEVMChainWithGeth struct {
	ChainBuilder
	config ctf_config.EthereumNetworkConfig
	ctf_test_env.EthereumNetworkHooks
}

func (n *NewEVMChainWithGeth) Chain() (deployment.Chain, persistent_types.RpcProvider, error) {
	chain := deployment.Chain{}

	ethBuilder := ctf_test_env.NewEthereumNetworkBuilder()
	network, err := ethBuilder.
		WithExistingConfig(n.config).
		Build()

	if err != nil {
		return chain, nil, err
	}

	evmNetwork, rpcProvider, err := network.Start()
	if err != nil {
		return chain, nil, err
	}

	evmNetwork.Name = fmt.Sprintf("%s-%d", *network.ExecutionLayer, evmNetwork.ChainID)

	return n.Build(evmNetwork, rpcProvider)
}

func CreateNewEVMChainWithGeth(config ctf_config.EthereumNetworkConfig, hooks ctf_test_env.EthereumNetworkHooks) persistent_types.NewEVMChainProducer {
	return &NewEVMChainWithGeth{
		config:               config,
		EthereumNetworkHooks: hooks,
	}
}

func (n *NewEVMChainWithGeth) Hooks() ctf_test_env.EthereumNetworkHooks {
	return n.EthereumNetworkHooks
}

func CreateExistingEVMChainConfigWithGeth(evmNetwork blockchain.EVMNetwork) persistent_types.ExistingEVMChainProducer {
	return &ExistingEVMChainConfigWithGeth{
		evmNetwork: evmNetwork,
	}
}

type ExistingEVMChainConfigWithGeth struct {
	ChainBuilder
	evmNetwork blockchain.EVMNetwork
}

func (e *ExistingEVMChainConfigWithGeth) EVMNetwork() blockchain.EVMNetwork {
	return e.evmNetwork
}

func (e *ExistingEVMChainConfigWithGeth) Chain() (deployment.Chain, persistent_types.RpcProvider, error) {
	return e.Build(e.evmNetwork, ctf_test_env.NewRPCProvider(e.evmNetwork.HTTPURLs, e.evmNetwork.URLs, e.evmNetwork.HTTPURLs, e.evmNetwork.URLs))
}
