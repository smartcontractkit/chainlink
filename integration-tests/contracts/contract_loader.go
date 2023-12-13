package contracts

import (
	"errors"

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrf_coordinator_v2_5"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrf_v2plus_load_test_with_metrics"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog"

	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/functions/generated/functions_coordinator"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/functions/generated/functions_load_test_client"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/functions/generated/functions_router"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/authorized_forwarder"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/link_token_interface"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/operator_wrapper"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrf_coordinator_v2"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrf_load_test_with_metrics"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/llo-feeds/generated/fee_manager"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/llo-feeds/generated/reward_manager"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/llo-feeds/generated/verifier"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/llo-feeds/generated/verifier_proxy"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/shared/generated/werc20_mock"
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

	// Mercury
	LoadMercuryVerifier(addr common.Address) (MercuryVerifier, error)
	LoadMercuryVerifierProxy(addr common.Address) (MercuryVerifierProxy, error)
	LoadMercuryFeeManager(addr common.Address) (MercuryFeeManager, error)
	LoadMercuryRewardManager(addr common.Address) (MercuryRewardManager, error)

	LoadWERC20Mock(addr common.Address) (WERC20Mock, error)

	// VRF
	LoadVRFCoordinatorV2(addr string) (VRFCoordinatorV2, error)
	LoadVRFv2LoadTestConsumer(addr string) (VRFv2LoadTestConsumer, error)
	LoadVRFCoordinatorV2_5(addr string) (VRFCoordinatorV2_5, error)
	LoadVRFv2PlusLoadTestConsumer(addr string) (VRFv2PlusLoadTestConsumer, error)
}

// NewContractLoader returns an instance of a contract Loader based on the client type
func NewContractLoader(bcClient blockchain.EVMClient, logger zerolog.Logger) (ContractLoader, error) {
	switch clientImpl := bcClient.Get().(type) {
	case *blockchain.EthereumClient:
		return NewEthereumContractLoader(clientImpl, logger), nil
	case *blockchain.KlaytnClient:
		return &KlaytnContractLoader{NewEthereumContractLoader(clientImpl, logger)}, nil
	case *blockchain.MetisClient:
		return &MetisContractLoader{NewEthereumContractLoader(clientImpl, logger)}, nil
	case *blockchain.ArbitrumClient:
		return &ArbitrumContractLoader{NewEthereumContractLoader(clientImpl, logger)}, nil
	case *blockchain.PolygonClient:
		return &PolygonContractLoader{NewEthereumContractLoader(clientImpl, logger)}, nil
	case *blockchain.OptimismClient:
		return &OptimismContractLoader{NewEthereumContractLoader(clientImpl, logger)}, nil
	case *blockchain.PolygonZkEvmClient:
		return &PolygonZkEvmContractLoader{NewEthereumContractLoader(clientImpl, logger)}, nil
	case *blockchain.WeMixClient:
		return &WeMixContractLoader{NewEthereumContractLoader(clientImpl, logger)}, nil
	case *blockchain.LineaClient:
		return &LineaContractLoader{NewEthereumContractLoader(clientImpl, logger)}, nil
	case *blockchain.CeloClient:
		return &CeloContractLoader{NewEthereumContractLoader(clientImpl, logger)}, nil
	case *blockchain.ScrollClient:
		return &ScrollContractLoader{NewEthereumContractLoader(clientImpl, logger)}, nil
	case *blockchain.FantomClient:
		return &FantomContractLoader{NewEthereumContractLoader(clientImpl, logger)}, nil
	case *blockchain.BSCClient:
		return &BSCContractLoader{NewEthereumContractLoader(clientImpl, logger)}, nil
	}
	return nil, errors.New("unknown blockchain client implementation for contract Loader, register blockchain client in NewContractLoader")
}

// EthereumContractLoader provides the implementations for deploying ETH (EVM) based contracts
type EthereumContractLoader struct {
	client blockchain.EVMClient
	l      zerolog.Logger
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
type PolygonZkEvmContractLoader struct {
	*EthereumContractLoader
}

// PolygonZKEVMContractLoader wraps for Polygon zkEVM
type PolygonZKEVMContractLoader struct {
	*EthereumContractLoader
}

// WeMixContractLoader wraps for WeMix
type WeMixContractLoader struct {
	*EthereumContractLoader
}

// LineaContractLoader wraps for Linea
type LineaContractLoader struct {
	*EthereumContractLoader
}

// CeloContractLoader wraps for Celo
type CeloContractLoader struct {
	*EthereumContractLoader
}

// ScrollContractLoader wraps for Scroll
type ScrollContractLoader struct {
	*EthereumContractLoader
}

// FantomContractLoader wraps for Fantom
type FantomContractLoader struct {
	*EthereumContractLoader
}

// BSCContractLoader wraps for BSC
type BSCContractLoader struct {
	*EthereumContractLoader
}

// NewEthereumContractLoader returns an instantiated instance of the ETH contract Loader
func NewEthereumContractLoader(ethClient blockchain.EVMClient, logger zerolog.Logger) *EthereumContractLoader {
	return &EthereumContractLoader{
		client: ethClient,
		l:      logger,
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
		l:        e.l,
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
		l:        e.l,
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
		l:        e.l,
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

// LoadMercuryVerifier returns Verifier contract deployed on given address
func (e *EthereumContractLoader) LoadMercuryVerifier(addr common.Address) (MercuryVerifier, error) {
	instance, err := e.client.LoadContract("Mercury Verifier", addr, func(
		address common.Address,
		backend bind.ContractBackend,
	) (interface{}, error) {
		return verifier.NewVerifier(address, backend)
	})
	if err != nil {
		return nil, err
	}
	return &EthereumMercuryVerifier{
		client:   e.client,
		instance: instance.(*verifier.Verifier),
		address:  addr,
	}, err
}

// LoadMercuryVerifierProxy returns VerifierProxy contract deployed on given address
func (e *EthereumContractLoader) LoadMercuryVerifierProxy(addr common.Address) (MercuryVerifierProxy, error) {
	instance, err := e.client.LoadContract("Mercury Verifier Proxy", addr, func(
		address common.Address,
		backend bind.ContractBackend,
	) (interface{}, error) {
		return verifier_proxy.NewVerifierProxy(address, backend)
	})
	if err != nil {
		return nil, err
	}
	return &EthereumMercuryVerifierProxy{
		client:   e.client,
		instance: instance.(*verifier_proxy.VerifierProxy),
		address:  addr,
	}, err
}

func (e *EthereumContractLoader) LoadMercuryFeeManager(addr common.Address) (MercuryFeeManager, error) {
	instance, err := e.client.LoadContract("Mercury Fee Manager", addr, func(
		address common.Address,
		backend bind.ContractBackend,
	) (interface{}, error) {
		return fee_manager.NewFeeManager(address, backend)
	})
	if err != nil {
		return nil, err
	}
	return &EthereumMercuryFeeManager{
		client:   e.client,
		instance: instance.(*fee_manager.FeeManager),
		address:  addr,
	}, err
}

func (e *EthereumContractLoader) LoadMercuryRewardManager(addr common.Address) (MercuryRewardManager, error) {
	instance, err := e.client.LoadContract("Mercury Reward Manager", addr, func(
		address common.Address,
		backend bind.ContractBackend,
	) (interface{}, error) {
		return reward_manager.NewRewardManager(address, backend)
	})
	if err != nil {
		return nil, err
	}
	return &EthereumMercuryRewardManager{
		client:   e.client,
		instance: instance.(*reward_manager.RewardManager),
		address:  addr,
	}, err
}

func (e *EthereumContractLoader) LoadWERC20Mock(addr common.Address) (WERC20Mock, error) {
	instance, err := e.client.LoadContract("WERC20 Mock", addr, func(
		address common.Address,
		backend bind.ContractBackend,
	) (interface{}, error) {
		return werc20_mock.NewWERC20Mock(address, backend)
	})
	if err != nil {
		return nil, err
	}
	return &EthereumWERC20Mock{
		client:   e.client,
		instance: instance.(*werc20_mock.WERC20Mock),
		address:  addr,
	}, err
}

func (e *EthereumContractLoader) LoadVRFCoordinatorV2_5(addr string) (VRFCoordinatorV2_5, error) {
	address := common.HexToAddress(addr)
	instance, err := e.client.LoadContract("VRFCoordinatorV2_5", address, func(
		address common.Address,
		backend bind.ContractBackend,
	) (interface{}, error) {
		return vrf_coordinator_v2_5.NewVRFCoordinatorV25(address, backend)
	})
	if err != nil {
		return nil, err
	}
	return &EthereumVRFCoordinatorV2_5{
		address:     &address,
		client:      e.client,
		coordinator: instance.(*vrf_coordinator_v2_5.VRFCoordinatorV25),
	}, err
}

func (e *EthereumContractLoader) LoadVRFv2PlusLoadTestConsumer(addr string) (VRFv2PlusLoadTestConsumer, error) {
	address := common.HexToAddress(addr)
	instance, err := e.client.LoadContract("VRFV2PlusLoadTestWithMetrics", address, func(
		address common.Address,
		backend bind.ContractBackend,
	) (interface{}, error) {
		return vrf_v2plus_load_test_with_metrics.NewVRFV2PlusLoadTestWithMetrics(address, backend)
	})
	if err != nil {
		return nil, err
	}
	return &EthereumVRFv2PlusLoadTestConsumer{
		client:   e.client,
		consumer: instance.(*vrf_v2plus_load_test_with_metrics.VRFV2PlusLoadTestWithMetrics),
		address:  &address,
	}, err
}

func (e *EthereumContractLoader) LoadVRFCoordinatorV2(addr string) (VRFCoordinatorV2, error) {
	address := common.HexToAddress(addr)
	instance, err := e.client.LoadContract("VRFCoordinatorV2", address, func(
		address common.Address,
		backend bind.ContractBackend,
	) (interface{}, error) {
		return vrf_coordinator_v2.NewVRFCoordinatorV2(address, backend)
	})
	if err != nil {
		return nil, err
	}
	return &EthereumVRFCoordinatorV2{
		address:     &address,
		client:      e.client,
		coordinator: instance.(*vrf_coordinator_v2.VRFCoordinatorV2),
	}, err
}

func (e *EthereumContractLoader) LoadVRFv2LoadTestConsumer(addr string) (VRFv2LoadTestConsumer, error) {
	address := common.HexToAddress(addr)
	instance, err := e.client.LoadContract("VRFV2LoadTestWithMetrics", address, func(
		address common.Address,
		backend bind.ContractBackend,
	) (interface{}, error) {
		return vrf_load_test_with_metrics.NewVRFV2LoadTestWithMetrics(address, backend)
	})
	if err != nil {
		return nil, err
	}
	return &EthereumVRFv2LoadTestConsumer{
		client:   e.client,
		consumer: instance.(*vrf_load_test_with_metrics.VRFV2LoadTestWithMetrics),
		address:  &address,
	}, err
}
