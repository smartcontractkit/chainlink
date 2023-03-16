package contracts

import (
	"errors"

	"github.com/smartcontractkit/chainlink-testing-framework/contracts/ethereum"
	int_ethereum "github.com/smartcontractkit/chainlink/integration-tests/contracts/ethereum"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink/core/gethwrappers/generated/authorized_forwarder"
	"github.com/smartcontractkit/chainlink/core/gethwrappers/generated/operator_wrapper"
)

// ContractLoader is an interface for abstracting the contract loading methods across network implementations
type ContractLoader interface {
	LoadOperatorContract(address common.Address) (Operator, error)
	LoadAuthorizedForwarder(address common.Address) (AuthorizedForwarder, error)
	LoadKeeperConsumerBenchmark(address common.Address) (KeeperConsumerBenchmark, error)
	LoadUpkeepResetter(address common.Address) (UpkeepResetter, error)
}

// NewContractLoader returns an instance of a contract Loader based on the client type
func NewContractLoader(bcClient blockchain.EVMClient) (ContractLoader, error) {
	switch clientImpl := bcClient.Get().(type) {
	case *blockchain.EthereumClient:
		return NewEthereumContractLoader(clientImpl), nil
	case *blockchain.KlaytnClient:
		return &KlaytnContractLoader{NewEthereumContractLoader(clientImpl)}, nil
	case *blockchain.MetisClient:
		return &MetisContractLoader{NewEthereumContractLoader(clientImpl)}, nil
	case *blockchain.ArbitrumClient:
		return &ArbitrumContractLoader{NewEthereumContractLoader(clientImpl)}, nil
	case *blockchain.PolygonClient:
		return &PolygonContractLoader{NewEthereumContractLoader(clientImpl)}, nil
	case *blockchain.OptimismClient:
		return &OptimismContractLoader{NewEthereumContractLoader(clientImpl)}, nil
	}
	return nil, errors.New("unknown blockchain client implementation for contract Loader, register blockchain client in NewContractLoader")
}

// EthereumContractLoader provides the implementations for deploying ETH (EVM) based contracts
type EthereumContractLoader struct {
	client blockchain.EVMClient
}

// KlaytnContractLoader wraps ethereum contract deployments for Klaytn
type KlaytnContractLoader struct {
	*EthereumContractLoader
}

// MetisContractLoader wraps ethereum contract deployments for Metis
type MetisContractLoader struct {
	*EthereumContractLoader
}

// ArbitrumContractLoader wraps for Arbitrum
type ArbitrumContractLoader struct {
	*EthereumContractLoader
}

// PolygonContractLoader wraps for Polygon
type PolygonContractLoader struct {
	*EthereumContractLoader
}

// OptimismContractLoader wraps for Optimism
type OptimismContractLoader struct {
	*EthereumContractLoader
}

// NewEthereumContractLoader returns an instantiated instance of the ETH contract Loader
func NewEthereumContractLoader(ethClient blockchain.EVMClient) *EthereumContractLoader {
	return &EthereumContractLoader{
		client: ethClient,
	}
}

// LoadOperatorContract returns deployed on given address Operator contract instance
func (e *EthereumContractLoader) LoadOperatorContract(address common.Address) (Operator, error) {
	instance, err := e.client.LoadContract("Operator", address, func(
		address common.Address,
		backend bind.ContractBackend,
	) (interface{}, error) {
		return operator_wrapper.NewOperator(address, backend)
	})
	if err != nil {
		return nil, err
	}
	return &EthereumOperator{
		address:  address,
		client:   e.client,
		operator: instance.(*operator_wrapper.Operator),
	}, err
}

// LoadAuthorizedForwarder returns deployed on given address AuthorizedForwarder contract instance
func (e *EthereumContractLoader) LoadAuthorizedForwarder(address common.Address) (AuthorizedForwarder, error) {
	instance, err := e.client.LoadContract("AuthorizedForwarder", address, func(
		address common.Address,
		backend bind.ContractBackend,
	) (interface{}, error) {
		return authorized_forwarder.NewAuthorizedForwarder(address, backend)
	})
	if err != nil {
		return nil, err
	}
	return &EthereumAuthorizedForwarder{
		address:             address,
		client:              e.client,
		authorizedForwarder: instance.(*authorized_forwarder.AuthorizedForwarder),
	}, err
}

// LoadKeeperConsumerBenchmark returns deployed on given address Keeper Consumer Contract
func (e *EthereumContractLoader) LoadKeeperConsumerBenchmark(address common.Address) (KeeperConsumerBenchmark, error) {
	instance, err := e.client.LoadContract("KeeperConsumerBenchmark", address, func(
		address common.Address,
		backend bind.ContractBackend,
	) (interface{}, error) {
		return ethereum.NewKeeperConsumerBenchmark(address, backend)
	})
	if err != nil {
		return nil, err
	}
	return &EthereumKeeperConsumerBenchmark{
		address:  &address,
		client:   e.client,
		consumer: instance.(*ethereum.KeeperConsumerBenchmark),
	}, err
}

// LoadUpkeepResetter returns deployed on given address Upkeep Resetter
func (e *EthereumContractLoader) LoadUpkeepResetter(address common.Address) (UpkeepResetter, error) {
	instance, err := e.client.LoadContract("KeeperConsumerBenchmark", address, func(
		address common.Address,
		backend bind.ContractBackend,
	) (interface{}, error) {
		return int_ethereum.NewUpkeepResetter(address, backend)
	})
	if err != nil {
		return nil, err
	}
	return &EthereumUpkeepResetter{
		address:  &address,
		client:   e.client,
		consumer: instance.(*int_ethereum.UpkeepResetter),
	}, err
}
