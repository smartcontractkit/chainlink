package contracts

import (
	"errors"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/functions/generated/functions_coordinator"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/functions/generated/functions_load_test_client"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/functions/generated/functions_router"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/authorized_forwarder"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/link_token_interface"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/operator_wrapper"
)

// ContractLoader is an interface for abstracting the contract loading methods across network implementations
type ContractLoader interface {
	LoadLINKToken(address string) (LinkToken, error)
	LoadOperatorContract(address common.Address) (Operator, error)
	LoadAuthorizedForwarder(address common.Address) (AuthorizedForwarder, error)

	/* functions 1_0_0 */
	LoadFunctionsCoordinator(addr string) (FunctionsCoordinator, error)
	LoadFunctionsRouter(addr string) (FunctionsRouter, error)
	LoadFunctionsLoadTestClient(addr string) (FunctionsLoadTestClient, error)
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

// LoadLINKToken returns deployed on given address LINK Token contract instance
func (e *EthereumContractLoader) LoadLINKToken(addr string) (LinkToken, error) {
	instance, err := e.client.LoadContract("LINK Token", common.HexToAddress(addr), func(
		address common.Address,
		backend bind.ContractBackend,
	) (interface{}, error) {
		return link_token_interface.NewLinkToken(address, backend)
	})
	if err != nil {
		return nil, err
	}
	return &EthereumLinkToken{
		client:   e.client,
		instance: instance.(*link_token_interface.LinkToken),
		address:  common.HexToAddress(addr),
	}, err
}

// LoadFunctionsCoordinator returns deployed on given address FunctionsCoordinator contract instance
func (e *EthereumContractLoader) LoadFunctionsCoordinator(addr string) (FunctionsCoordinator, error) {
	instance, err := e.client.LoadContract("Functions Coordinator", common.HexToAddress(addr), func(
		address common.Address,
		backend bind.ContractBackend,
	) (interface{}, error) {
		return functions_coordinator.NewFunctionsCoordinator(address, backend)
	})
	if err != nil {
		return nil, err
	}
	return &EthereumFunctionsCoordinator{
		client:   e.client,
		instance: instance.(*functions_coordinator.FunctionsCoordinator),
		address:  common.HexToAddress(addr),
	}, err
}

// LoadFunctionsRouter returns deployed on given address FunctionsRouter contract instance
func (e *EthereumContractLoader) LoadFunctionsRouter(addr string) (FunctionsRouter, error) {
	instance, err := e.client.LoadContract("Functions Router", common.HexToAddress(addr), func(
		address common.Address,
		backend bind.ContractBackend,
	) (interface{}, error) {
		return functions_router.NewFunctionsRouter(address, backend)
	})
	if err != nil {
		return nil, err
	}
	return &EthereumFunctionsRouter{
		client:   e.client,
		instance: instance.(*functions_router.FunctionsRouter),
		address:  common.HexToAddress(addr),
	}, err
}

// LoadFunctionsLoadTestClient returns deployed on given address FunctionsLoadTestClient contract instance
func (e *EthereumContractLoader) LoadFunctionsLoadTestClient(addr string) (FunctionsLoadTestClient, error) {
	instance, err := e.client.LoadContract("FunctionsLoadTestClient", common.HexToAddress(addr), func(
		address common.Address,
		backend bind.ContractBackend,
	) (interface{}, error) {
		return functions_load_test_client.NewFunctionsLoadTestClient(address, backend)
	})
	if err != nil {
		return nil, err
	}
	return &EthereumFunctionsLoadTestClient{
		client:   e.client,
		instance: instance.(*functions_load_test_client.FunctionsLoadTestClient),
		address:  common.HexToAddress(addr),
	}, err
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
