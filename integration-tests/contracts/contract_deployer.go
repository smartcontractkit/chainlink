package contracts

import (
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog/log"

	"github.com/smartcontractkit/libocr/gethwrappers/offchainaggregator"
	"github.com/smartcontractkit/libocr/gethwrappers2/ocr2aggregator"
	ocrConfigHelper "github.com/smartcontractkit/libocr/offchainreporting/confighelper"

	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/automation_consumer_benchmark"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/flags_wrapper"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/flux_aggregator_wrapper"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/functions_billing_registry_events_mock"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/functions_oracle_events_mock"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/keeper_registrar_wrapper1_2"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/keeper_registrar_wrapper2_0"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/keeper_registry_logic1_3"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/keeper_registry_logic2_0"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/keeper_registry_wrapper1_1"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/keeper_registry_wrapper1_2"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/keeper_registry_wrapper1_3"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/keeper_registry_wrapper2_0"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/link_token_interface"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/mock_aggregator_proxy"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/mock_ethlink_aggregator_wrapper"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/mock_gas_aggregator_wrapper"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/operator_factory"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/oracle_wrapper"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/test_api_consumer_wrapper"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/upkeep_transcoder"

	eth_contracts "github.com/smartcontractkit/chainlink/integration-tests/contracts/ethereum"
)

// ContractDeployer is an interface for abstracting the contract deployment methods across network implementations
type ContractDeployer interface {
	DeployAPIConsumer(linkAddr string) (APIConsumer, error)
	DeployOracle(linkAddr string) (Oracle, error)
	DeployFlags(rac string) (Flags, error)
	DeployFluxAggregatorContract(linkAddr string, fluxOptions FluxAggregatorOptions) (FluxAggregator, error)
	DeployLinkTokenContract() (LinkToken, error)
	LoadLinkToken(address common.Address) (LinkToken, error)
	DeployOffChainAggregator(linkAddr string, offchainOptions OffchainOptions) (OffchainAggregator, error)
	DeployVRFContract() (VRF, error)
	DeployMockETHLINKFeed(answer *big.Int) (MockETHLINKFeed, error)
	LoadETHLINKFeed(address common.Address) (MockETHLINKFeed, error)
	DeployMockGasFeed(answer *big.Int) (MockGasFeed, error)
	LoadGasFeed(address common.Address) (MockGasFeed, error)
	DeployKeeperRegistrar(registryVersion eth_contracts.KeeperRegistryVersion, linkAddr string, registrarSettings KeeperRegistrarSettings) (KeeperRegistrar, error)
	LoadKeeperRegistrar(address common.Address, registryVersion eth_contracts.KeeperRegistryVersion) (KeeperRegistrar, error)
	DeployUpkeepTranscoder() (UpkeepTranscoder, error)
	DeployKeeperRegistry(opts *KeeperRegistryOpts) (KeeperRegistry, error)
	LoadKeeperRegistry(address common.Address, registryVersion eth_contracts.KeeperRegistryVersion) (KeeperRegistry, error)
	DeployKeeperConsumer(updateInterval *big.Int) (KeeperConsumer, error)
	DeployKeeperConsumerPerformance(
		testBlockRange,
		averageCadence,
		checkGasToBurn,
		performGasToBurn *big.Int,
	) (KeeperConsumerPerformance, error)
	DeployKeeperConsumerBenchmark() (AutomationConsumerBenchmark, error)
	LoadKeeperConsumerBenchmark(address common.Address) (AutomationConsumerBenchmark, error)
	DeployKeeperPerformDataChecker(expectedData []byte) (KeeperPerformDataChecker, error)
	DeployUpkeepCounter(testRange *big.Int, interval *big.Int) (UpkeepCounter, error)
	DeployUpkeepPerformCounterRestrictive(testRange *big.Int, averageEligibilityCadence *big.Int) (UpkeepPerformCounterRestrictive, error)
	DeployVRFConsumer(linkAddr string, coordinatorAddr string) (VRFConsumer, error)
	DeployVRFConsumerV2(linkAddr string, coordinatorAddr string) (VRFConsumerV2, error)
	DeployVRFv2Consumer(coordinatorAddr string) (VRFv2Consumer, error)
	DeployVRFv2LoadTestConsumer(coordinatorAddr string) (VRFv2LoadTestConsumer, error)
	DeployVRFCoordinator(linkAddr string, bhsAddr string) (VRFCoordinator, error)
	DeployVRFCoordinatorV2(linkAddr string, bhsAddr string, linkEthFeedAddr string) (VRFCoordinatorV2, error)
	DeployDKG() (DKG, error)
	DeployOCR2VRFCoordinator(beaconPeriodBlocksCount *big.Int, linkAddr string) (VRFCoordinatorV3, error)
	DeployVRFBeacon(vrfCoordinatorAddress string, linkAddress string, dkgAddress string, keyId string) (VRFBeacon, error)
	DeployVRFBeaconConsumer(vrfCoordinatorAddress string, beaconPeriodBlockCount *big.Int) (VRFBeaconConsumer, error)
	DeployBlockhashStore() (BlockHashStore, error)
	DeployOperatorFactory(linkAddr string) (OperatorFactory, error)
	DeployStaking(params eth_contracts.StakingPoolConstructorParams) (Staking, error)
	DeployBatchBlockhashStore(blockhashStoreAddr string) (BatchBlockhashStore, error)
	DeployFunctionsOracleEventsMock() (FunctionsOracleEventsMock, error)
	DeployFunctionsBillingRegistryEventsMock() (FunctionsBillingRegistryEventsMock, error)
	DeployMockAggregatorProxy(aggregatorAddr string) (MockAggregatorProxy, error)
	DeployOffchainAggregatorV2(linkAddr string, offchainOptions OffchainOptions) (OffchainAggregatorV2, error)
}

// NewContractDeployer returns an instance of a contract deployer based on the client type
func NewContractDeployer(bcClient blockchain.EVMClient) (ContractDeployer, error) {
	switch clientImpl := bcClient.Get().(type) {
	case *blockchain.EthereumClient:
		return NewEthereumContractDeployer(clientImpl), nil
	case *blockchain.KlaytnClient:
		return &KlaytnContractDeployer{NewEthereumContractDeployer(clientImpl)}, nil
	case *blockchain.MetisClient:
		return &MetisContractDeployer{NewEthereumContractDeployer(clientImpl)}, nil
	case *blockchain.ArbitrumClient:
		return &MetisContractDeployer{NewEthereumContractDeployer(clientImpl)}, nil
	case *blockchain.OptimismClient:
		return &OptimismContractDeployer{NewEthereumContractDeployer(clientImpl)}, nil
	case *blockchain.RSKClient:
		return &RSKContractDeployer{NewEthereumContractDeployer(clientImpl)}, nil
	case *blockchain.PolygonClient:
		return &PolygonContractDeployer{NewEthereumContractDeployer(clientImpl)}, nil
	case *blockchain.CeloClient:
		return &CeloContractDeployer{NewEthereumContractDeployer(clientImpl)}, nil
	case *blockchain.QuorumClient:
		return &QuorumContractDeployer{NewEthereumContractDeployer(clientImpl)}, nil
	case *blockchain.ScrollClient:
		return &ScrollContractDeployer{NewEthereumContractDeployer(clientImpl)}, nil
	}
	return nil, errors.New("unknown blockchain client implementation for contract deployer, register blockchain client in NewContractDeployer")
}

// EthereumContractDeployer provides the implementations for deploying ETH (EVM) based contracts
type EthereumContractDeployer struct {
	client blockchain.EVMClient
}

// KlaytnContractDeployer wraps ethereum contract deployments for Klaytn
type KlaytnContractDeployer struct {
	*EthereumContractDeployer
}

// MetisContractDeployer wraps ethereum contract deployments for Metis
type MetisContractDeployer struct {
	*EthereumContractDeployer
}

// ArbitrumContractDeployer wraps for Arbitrum
type ArbitrumContractDeployer struct {
	*EthereumContractDeployer
}

// OptimismContractDeployer wraps for Optimism
type OptimismContractDeployer struct {
	*EthereumContractDeployer
}

// RSKContractDeployer wraps for RSK
type RSKContractDeployer struct {
	*EthereumContractDeployer
}

type PolygonContractDeployer struct {
	*EthereumContractDeployer
}

type CeloContractDeployer struct {
	*EthereumContractDeployer
}

type QuorumContractDeployer struct {
	*EthereumContractDeployer
}

type ScrollContractDeployer struct {
	*EthereumContractDeployer
}

// NewEthereumContractDeployer returns an instantiated instance of the ETH contract deployer
func NewEthereumContractDeployer(ethClient blockchain.EVMClient) *EthereumContractDeployer {
	return &EthereumContractDeployer{
		client: ethClient,
	}
}

// DefaultFluxAggregatorOptions produces some basic defaults for a flux aggregator contract
func DefaultFluxAggregatorOptions() FluxAggregatorOptions {
	return FluxAggregatorOptions{
		PaymentAmount: big.NewInt(1),
		Timeout:       uint32(30),
		MinSubValue:   big.NewInt(0),
		MaxSubValue:   big.NewInt(1000000000000),
		Decimals:      uint8(0),
		Description:   "Test Flux Aggregator",
	}
}

// DeployFlags deploys flags contract
func (e *EthereumContractDeployer) DeployFlags(
	rac string,
) (Flags, error) {
	address, _, instance, err := e.client.DeployContract("Flags", func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		racAddr := common.HexToAddress(rac)
		return flags_wrapper.DeployFlags(auth, backend, racAddr)
	})
	if err != nil {
		return nil, err
	}
	return &EthereumFlags{
		client:  e.client,
		flags:   instance.(*flags_wrapper.Flags),
		address: address,
	}, nil
}

// DeployFluxAggregatorContract deploys the Flux Aggregator Contract on an EVM chain
func (e *EthereumContractDeployer) DeployFluxAggregatorContract(
	linkAddr string,
	fluxOptions FluxAggregatorOptions,
) (FluxAggregator, error) {
	address, _, instance, err := e.client.DeployContract("Flux Aggregator", func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		la := common.HexToAddress(linkAddr)
		return flux_aggregator_wrapper.DeployFluxAggregator(auth,
			backend,
			la,
			fluxOptions.PaymentAmount,
			fluxOptions.Timeout,
			fluxOptions.Validator,
			fluxOptions.MinSubValue,
			fluxOptions.MaxSubValue,
			fluxOptions.Decimals,
			fluxOptions.Description)
	})
	if err != nil {
		return nil, err
	}
	return &EthereumFluxAggregator{
		client:         e.client,
		fluxAggregator: instance.(*flux_aggregator_wrapper.FluxAggregator),
		address:        address,
	}, nil
}

func (e *EthereumContractDeployer) DeployStaking(params eth_contracts.StakingPoolConstructorParams) (Staking, error) {
	stakingAddress, _, instance, err := e.client.DeployContract("Staking", func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		return eth_contracts.DeployStaking(auth, backend, params)
	})
	if err != nil {
		return nil, err
	}
	return &EthereumStaking{
		client:  e.client,
		staking: instance.(*eth_contracts.Staking),
		address: stakingAddress,
	}, nil
}

func (e *EthereumContractDeployer) DeployFunctionsOracleEventsMock() (FunctionsOracleEventsMock, error) {
	address, _, instance, err := e.client.DeployContract("FunctionsOracleEventsMock", func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		return functions_oracle_events_mock.DeployFunctionsOracleEventsMock(auth, backend)
	})
	if err != nil {
		return nil, err
	}
	return &EthereumFunctionsOracleEventsMock{
		client:     e.client,
		eventsMock: instance.(*functions_oracle_events_mock.FunctionsOracleEventsMock),
		address:    address,
	}, nil
}

func (e *EthereumContractDeployer) DeployFunctionsBillingRegistryEventsMock() (FunctionsBillingRegistryEventsMock, error) {
	address, _, instance, err := e.client.DeployContract("FunctionsBillingRegistryEventsMock", func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		return functions_billing_registry_events_mock.DeployFunctionsBillingRegistryEventsMock(auth, backend)
	})
	if err != nil {
		return nil, err
	}
	return &EthereumFunctionsBillingRegistryEventsMock{
		client:     e.client,
		eventsMock: instance.(*functions_billing_registry_events_mock.FunctionsBillingRegistryEventsMock),
		address:    address,
	}, nil
}

// DeployLinkTokenContract deploys a Link Token contract to an EVM chain
func (e *EthereumContractDeployer) DeployLinkTokenContract() (LinkToken, error) {
	linkTokenAddress, _, instance, err := e.client.DeployContract("LINK Token", func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		return link_token_interface.DeployLinkToken(auth, backend)
	})
	if err != nil {
		return nil, err
	}

	return &EthereumLinkToken{
		client:   e.client,
		instance: instance.(*link_token_interface.LinkToken),
		address:  *linkTokenAddress,
	}, err
}

// LoadLinkToken returns deployed on given address EthereumLinkToken
func (e *EthereumContractDeployer) LoadLinkToken(address common.Address) (LinkToken, error) {
	instance, err := e.client.LoadContract("LinkToken", address, func(
		address common.Address,
		backend bind.ContractBackend,
	) (interface{}, error) {
		return link_token_interface.NewLinkToken(address, backend)
	})
	if err != nil {
		return nil, err
	}
	return &EthereumLinkToken{
		address:  address,
		client:   e.client,
		instance: instance.(*link_token_interface.LinkToken),
	}, err
}

func (e *EthereumContractDeployer) NewLinkTokenContract(address common.Address) (LinkToken, error) {
	return e.LoadLinkToken(address)
}

// DefaultOffChainAggregatorOptions returns some base defaults for deploying an OCR contract
func DefaultOffChainAggregatorOptions() OffchainOptions {
	return OffchainOptions{
		MaximumGasPrice:         uint32(3000),
		ReasonableGasPrice:      uint32(10),
		MicroLinkPerEth:         uint32(500),
		LinkGweiPerObservation:  uint32(500),
		LinkGweiPerTransmission: uint32(500),
		MinimumAnswer:           big.NewInt(1),
		MaximumAnswer:           big.NewInt(50000000000000000),
		Decimals:                8,
		Description:             "Test OCR",
	}
}

// DefaultOffChainAggregatorConfig returns some base defaults for configuring an OCR contract
func DefaultOffChainAggregatorConfig(numberNodes int) OffChainAggregatorConfig {
	if numberNodes <= 4 {
		log.Err(fmt.Errorf("insufficient number of nodes (%d) supplied for OCR, need at least 5", numberNodes)).
			Int("Number Chainlink Nodes", numberNodes).
			Msg("You likely need more chainlink nodes to properly configure OCR, try 5 or more.")
	}
	s := []int{1}
	// First node's stage already inputted as a 1 in line above, so numberNodes-1.
	for i := 0; i < numberNodes-1; i++ {
		s = append(s, 2)
	}
	return OffChainAggregatorConfig{
		AlphaPPB:         1,
		DeltaC:           time.Minute * 60,
		DeltaGrace:       time.Second * 12,
		DeltaProgress:    time.Second * 35,
		DeltaStage:       time.Second * 60,
		DeltaResend:      time.Second * 17,
		DeltaRound:       time.Second * 30,
		RMax:             6,
		S:                s,
		N:                numberNodes,
		F:                1,
		OracleIdentities: []ocrConfigHelper.OracleIdentityExtra{},
	}
}

// DeployOffChainAggregator deploys the offchain aggregation contract to the EVM chain
func (e *EthereumContractDeployer) DeployOffChainAggregator(
	linkAddr string,
	offchainOptions OffchainOptions,
) (OffchainAggregator, error) {
	address, _, instance, err := e.client.DeployContract("OffChain Aggregator", func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		la := common.HexToAddress(linkAddr)
		return offchainaggregator.DeployOffchainAggregator(auth,
			backend,
			offchainOptions.MaximumGasPrice,
			offchainOptions.ReasonableGasPrice,
			offchainOptions.MicroLinkPerEth,
			offchainOptions.LinkGweiPerObservation,
			offchainOptions.LinkGweiPerTransmission,
			la,
			offchainOptions.MinimumAnswer,
			offchainOptions.MaximumAnswer,
			offchainOptions.BillingAccessController,
			offchainOptions.RequesterAccessController,
			offchainOptions.Decimals,
			offchainOptions.Description)
	})
	if err != nil {
		return nil, err
	}
	return &EthereumOffchainAggregator{
		client:  e.client,
		ocr:     instance.(*offchainaggregator.OffchainAggregator),
		address: address,
	}, err
}

// DeployAPIConsumer deploys api consumer for oracle
func (e *EthereumContractDeployer) DeployAPIConsumer(linkAddr string) (APIConsumer, error) {
	addr, _, instance, err := e.client.DeployContract("TestAPIConsumer", func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		return test_api_consumer_wrapper.DeployTestAPIConsumer(auth, backend, common.HexToAddress(linkAddr))
	})
	if err != nil {
		return nil, err
	}
	return &EthereumAPIConsumer{
		address:  addr,
		client:   e.client,
		consumer: instance.(*test_api_consumer_wrapper.TestAPIConsumer),
	}, err
}

// DeployOracle deploys oracle for consumer test
func (e *EthereumContractDeployer) DeployOracle(linkAddr string) (Oracle, error) {
	addr, _, instance, err := e.client.DeployContract("Oracle", func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		return oracle_wrapper.DeployOracle(auth, backend, common.HexToAddress(linkAddr))
	})
	if err != nil {
		return nil, err
	}
	return &EthereumOracle{
		address: addr,
		client:  e.client,
		oracle:  instance.(*oracle_wrapper.Oracle),
	}, err
}

func (e *EthereumContractDeployer) DeployMockETHLINKFeed(answer *big.Int) (MockETHLINKFeed, error) {
	address, _, instance, err := e.client.DeployContract("MockETHLINKAggregator", func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		return mock_ethlink_aggregator_wrapper.DeployMockETHLINKAggregator(auth, backend, answer)
	})
	if err != nil {
		return nil, err
	}
	return &EthereumMockETHLINKFeed{
		client:  e.client,
		feed:    instance.(*mock_ethlink_aggregator_wrapper.MockETHLINKAggregator),
		address: address,
	}, err
}

// LoadETHLINKFeed returns deployed on given address EthereumMockETHLINKFeed
func (e *EthereumContractDeployer) LoadETHLINKFeed(address common.Address) (MockETHLINKFeed, error) {
	instance, err := e.client.LoadContract("MockETHLINKFeed", address, func(
		address common.Address,
		backend bind.ContractBackend,
	) (interface{}, error) {
		return mock_ethlink_aggregator_wrapper.NewMockETHLINKAggregator(address, backend)
	})
	if err != nil {
		return nil, err
	}
	return &EthereumMockETHLINKFeed{
		address: &address,
		client:  e.client,
		feed:    instance.(*mock_ethlink_aggregator_wrapper.MockETHLINKAggregator),
	}, err
}

func (e *EthereumContractDeployer) DeployMockGasFeed(answer *big.Int) (MockGasFeed, error) {
	address, _, instance, err := e.client.DeployContract("MockGasFeed", func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		return mock_gas_aggregator_wrapper.DeployMockGASAggregator(auth, backend, answer)
	})
	if err != nil {
		return nil, err
	}
	return &EthereumMockGASFeed{
		client:  e.client,
		feed:    instance.(*mock_gas_aggregator_wrapper.MockGASAggregator),
		address: address,
	}, err
}

// LoadGasFeed returns deployed on given address EthereumMockGASFeed
func (e *EthereumContractDeployer) LoadGasFeed(address common.Address) (MockGasFeed, error) {
	instance, err := e.client.LoadContract("MockETHLINKFeed", address, func(
		address common.Address,
		backend bind.ContractBackend,
	) (interface{}, error) {
		return mock_gas_aggregator_wrapper.NewMockGASAggregator(address, backend)
	})
	if err != nil {
		return nil, err
	}
	return &EthereumMockGASFeed{
		address: &address,
		client:  e.client,
		feed:    instance.(*mock_gas_aggregator_wrapper.MockGASAggregator),
	}, err
}

func (e *EthereumContractDeployer) DeployUpkeepTranscoder() (UpkeepTranscoder, error) {
	address, _, instance, err := e.client.DeployContract("UpkeepTranscoder", func(
		opts *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		return upkeep_transcoder.DeployUpkeepTranscoder(opts, backend)
	})

	if err != nil {
		return nil, err
	}

	return &EthereumUpkeepTranscoder{
		client:     e.client,
		transcoder: instance.(*upkeep_transcoder.UpkeepTranscoder),
		address:    address,
	}, err
}

func (e *EthereumContractDeployer) DeployKeeperRegistrar(registryVersion eth_contracts.KeeperRegistryVersion, linkAddr string,
	registrarSettings KeeperRegistrarSettings) (KeeperRegistrar, error) {

	if registryVersion == eth_contracts.RegistryVersion_2_0 {
		// deploy registrar 2.0
		address, _, instance, err := e.client.DeployContract("KeeperRegistrar", func(
			opts *bind.TransactOpts,
			backend bind.ContractBackend,
		) (common.Address, *types.Transaction, interface{}, error) {
			return keeper_registrar_wrapper2_0.DeployKeeperRegistrar(opts, backend, common.HexToAddress(linkAddr), registrarSettings.AutoApproveConfigType,
				registrarSettings.AutoApproveMaxAllowed, common.HexToAddress(registrarSettings.RegistryAddr), registrarSettings.MinLinkJuels)
		})

		if err != nil {
			return nil, err
		}

		return &EthereumKeeperRegistrar{
			client:      e.client,
			registrar20: instance.(*keeper_registrar_wrapper2_0.KeeperRegistrar),
			address:     address,
		}, err
	}
	// non OCR registrar
	address, _, instance, err := e.client.DeployContract("KeeperRegistrar", func(
		opts *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		return keeper_registrar_wrapper1_2.DeployKeeperRegistrar(opts, backend, common.HexToAddress(linkAddr), registrarSettings.AutoApproveConfigType,
			registrarSettings.AutoApproveMaxAllowed, common.HexToAddress(registrarSettings.RegistryAddr), registrarSettings.MinLinkJuels)
	})

	if err != nil {
		return nil, err
	}

	return &EthereumKeeperRegistrar{
		client:    e.client,
		registrar: instance.(*keeper_registrar_wrapper1_2.KeeperRegistrar),
		address:   address,
	}, err
}

// LoadKeeperRegistrar returns deployed on given address EthereumKeeperRegistrar
func (e *EthereumContractDeployer) LoadKeeperRegistrar(address common.Address, registryVersion eth_contracts.KeeperRegistryVersion) (KeeperRegistrar, error) {
	if registryVersion == eth_contracts.RegistryVersion_1_1 || registryVersion == eth_contracts.RegistryVersion_1_2 ||
		registryVersion == eth_contracts.RegistryVersion_1_3 {
		instance, err := e.client.LoadContract("KeeperRegistrar", address, func(
			address common.Address,
			backend bind.ContractBackend,
		) (interface{}, error) {
			return keeper_registrar_wrapper1_2.NewKeeperRegistrar(address, backend)
		})
		if err != nil {
			return nil, err
		}
		return &EthereumKeeperRegistrar{
			address:   &address,
			client:    e.client,
			registrar: instance.(*keeper_registrar_wrapper1_2.KeeperRegistrar),
		}, err
	}
	instance, err := e.client.LoadContract("KeeperRegistrar", address, func(
		address common.Address,
		backend bind.ContractBackend,
	) (interface{}, error) {
		return keeper_registrar_wrapper2_0.NewKeeperRegistrar(address, backend)
	})
	if err != nil {
		return nil, err
	}
	return &EthereumKeeperRegistrar{
		address:     &address,
		client:      e.client,
		registrar20: instance.(*keeper_registrar_wrapper2_0.KeeperRegistrar),
	}, err
}

func (e *EthereumContractDeployer) DeployKeeperRegistry(
	opts *KeeperRegistryOpts,
) (KeeperRegistry, error) {
	var mode uint8
	switch e.client.GetChainID() {
	//Arbitrum payment model
	case big.NewInt(421613):
		mode = uint8(1)
	//Optimism payment model
	case big.NewInt(420):
		mode = uint8(2)
	default:
		mode = uint8(0)
	}
	registryGasOverhead := big.NewInt(80000)
	switch opts.RegistryVersion {
	case eth_contracts.RegistryVersion_1_0, eth_contracts.RegistryVersion_1_1:
		address, _, instance, err := e.client.DeployContract("KeeperRegistry1_1", func(
			auth *bind.TransactOpts,
			backend bind.ContractBackend,
		) (common.Address, *types.Transaction, interface{}, error) {
			return keeper_registry_wrapper1_1.DeployKeeperRegistry(
				auth,
				backend,
				common.HexToAddress(opts.LinkAddr),
				common.HexToAddress(opts.ETHFeedAddr),
				common.HexToAddress(opts.GasFeedAddr),
				opts.Settings.PaymentPremiumPPB,
				opts.Settings.FlatFeeMicroLINK,
				opts.Settings.BlockCountPerTurn,
				opts.Settings.CheckGasLimit,
				opts.Settings.StalenessSeconds,
				opts.Settings.GasCeilingMultiplier,
				opts.Settings.FallbackGasPrice,
				opts.Settings.FallbackLinkPrice,
			)
		})
		if err != nil {
			return nil, err
		}
		return &EthereumKeeperRegistry{
			client:      e.client,
			version:     eth_contracts.RegistryVersion_1_1,
			registry1_1: instance.(*keeper_registry_wrapper1_1.KeeperRegistry),
			registry1_2: nil,
			registry1_3: nil,
			address:     address,
		}, err
	case eth_contracts.RegistryVersion_1_2:
		address, _, instance, err := e.client.DeployContract("KeeperRegistry1_2", func(
			auth *bind.TransactOpts,
			backend bind.ContractBackend,
		) (common.Address, *types.Transaction, interface{}, error) {
			return keeper_registry_wrapper1_2.DeployKeeperRegistry(
				auth,
				backend,
				common.HexToAddress(opts.LinkAddr),
				common.HexToAddress(opts.ETHFeedAddr),
				common.HexToAddress(opts.GasFeedAddr),
				keeper_registry_wrapper1_2.Config{
					PaymentPremiumPPB:    opts.Settings.PaymentPremiumPPB,
					FlatFeeMicroLink:     opts.Settings.FlatFeeMicroLINK,
					BlockCountPerTurn:    opts.Settings.BlockCountPerTurn,
					CheckGasLimit:        opts.Settings.CheckGasLimit,
					StalenessSeconds:     opts.Settings.StalenessSeconds,
					GasCeilingMultiplier: opts.Settings.GasCeilingMultiplier,
					MinUpkeepSpend:       opts.Settings.MinUpkeepSpend,
					MaxPerformGas:        opts.Settings.MaxPerformGas,
					FallbackGasPrice:     opts.Settings.FallbackGasPrice,
					FallbackLinkPrice:    opts.Settings.FallbackLinkPrice,
					Transcoder:           common.HexToAddress(opts.TranscoderAddr),
					Registrar:            common.HexToAddress(opts.RegistrarAddr),
				},
			)
		})
		if err != nil {
			return nil, err
		}
		return &EthereumKeeperRegistry{
			client:      e.client,
			version:     eth_contracts.RegistryVersion_1_2,
			registry1_1: nil,
			registry1_2: instance.(*keeper_registry_wrapper1_2.KeeperRegistry),
			registry1_3: nil,
			address:     address,
		}, err
	case eth_contracts.RegistryVersion_1_3:
		logicAddress, _, _, err := e.client.DeployContract("KeeperRegistryLogic1_3", func(
			auth *bind.TransactOpts,
			backend bind.ContractBackend,
		) (common.Address, *types.Transaction, interface{}, error) {
			return keeper_registry_logic1_3.DeployKeeperRegistryLogic(
				auth,
				backend,
				mode,                // Default payment model
				registryGasOverhead, // Registry gas overhead
				common.HexToAddress(opts.LinkAddr),
				common.HexToAddress(opts.ETHFeedAddr),
				common.HexToAddress(opts.GasFeedAddr),
			)
		})
		if err != nil {
			return nil, err
		}
		err = e.client.WaitForEvents()
		if err != nil {
			return nil, err
		}

		address, _, instance, err := e.client.DeployContract("KeeperRegistry1_3", func(
			auth *bind.TransactOpts,
			backend bind.ContractBackend,
		) (common.Address, *types.Transaction, interface{}, error) {
			return keeper_registry_wrapper1_3.DeployKeeperRegistry(
				auth,
				backend,
				*logicAddress,
				keeper_registry_wrapper1_3.Config{
					PaymentPremiumPPB:    opts.Settings.PaymentPremiumPPB,
					FlatFeeMicroLink:     opts.Settings.FlatFeeMicroLINK,
					BlockCountPerTurn:    opts.Settings.BlockCountPerTurn,
					CheckGasLimit:        opts.Settings.CheckGasLimit,
					StalenessSeconds:     opts.Settings.StalenessSeconds,
					GasCeilingMultiplier: opts.Settings.GasCeilingMultiplier,
					MinUpkeepSpend:       opts.Settings.MinUpkeepSpend,
					MaxPerformGas:        opts.Settings.MaxPerformGas,
					FallbackGasPrice:     opts.Settings.FallbackGasPrice,
					FallbackLinkPrice:    opts.Settings.FallbackLinkPrice,
					Transcoder:           common.HexToAddress(opts.TranscoderAddr),
					Registrar:            common.HexToAddress(opts.RegistrarAddr),
				},
			)
		})
		if err != nil {
			return nil, err
		}
		return &EthereumKeeperRegistry{
			client:      e.client,
			version:     eth_contracts.RegistryVersion_1_3,
			registry1_1: nil,
			registry1_2: nil,
			registry1_3: instance.(*keeper_registry_wrapper1_3.KeeperRegistry),
			address:     address,
		}, err
	case eth_contracts.RegistryVersion_2_0:
		logicAddress, _, _, err := e.client.DeployContract("KeeperRegistryLogic2_0", func(
			auth *bind.TransactOpts,
			backend bind.ContractBackend,
		) (common.Address, *types.Transaction, interface{}, error) {
			return keeper_registry_logic2_0.DeployKeeperRegistryLogic(
				auth,
				backend,
				mode, // Default payment model
				common.HexToAddress(opts.LinkAddr),
				common.HexToAddress(opts.ETHFeedAddr),
				common.HexToAddress(opts.GasFeedAddr),
			)
		})
		if err != nil {
			return nil, err
		}
		err = e.client.WaitForEvents()
		if err != nil {
			return nil, err
		}

		address, _, instance, err := e.client.DeployContract("KeeperRegistry2_0", func(
			auth *bind.TransactOpts,
			backend bind.ContractBackend,
		) (common.Address, *types.Transaction, interface{}, error) {

			return keeper_registry_wrapper2_0.DeployKeeperRegistry(
				auth,
				backend,
				*logicAddress,
			)
		})
		if err != nil {
			return nil, err
		}
		return &EthereumKeeperRegistry{
			client:      e.client,
			version:     eth_contracts.RegistryVersion_2_0,
			registry2_0: instance.(*keeper_registry_wrapper2_0.KeeperRegistry),
			address:     address,
		}, err

	default:
		return nil, fmt.Errorf("keeper registry version %d is not supported", opts.RegistryVersion)
	}
}

// LoadKeeperRegistry returns deployed on given address EthereumKeeperRegistry
func (e *EthereumContractDeployer) LoadKeeperRegistry(address common.Address, registryVersion eth_contracts.KeeperRegistryVersion) (KeeperRegistry, error) {
	switch registryVersion {
	case eth_contracts.RegistryVersion_1_1:
		instance, err := e.client.LoadContract("KeeperRegistry", address, func(
			address common.Address,
			backend bind.ContractBackend,
		) (interface{}, error) {
			return keeper_registry_wrapper1_1.NewKeeperRegistry(address, backend)
		})
		if err != nil {
			return nil, err
		}
		return &EthereumKeeperRegistry{
			address:     &address,
			client:      e.client,
			registry1_1: instance.(*keeper_registry_wrapper1_1.KeeperRegistry),
		}, err
	case eth_contracts.RegistryVersion_1_2:
		instance, err := e.client.LoadContract("KeeperRegistry", address, func(
			address common.Address,
			backend bind.ContractBackend,
		) (interface{}, error) {
			return keeper_registry_wrapper1_2.NewKeeperRegistry(address, backend)
		})
		if err != nil {
			return nil, err
		}
		return &EthereumKeeperRegistry{
			address:     &address,
			client:      e.client,
			registry1_2: instance.(*keeper_registry_wrapper1_2.KeeperRegistry),
		}, err
	case eth_contracts.RegistryVersion_1_3:
		instance, err := e.client.LoadContract("KeeperRegistry", address, func(
			address common.Address,
			backend bind.ContractBackend,
		) (interface{}, error) {
			return keeper_registry_wrapper1_3.NewKeeperRegistry(address, backend)
		})
		if err != nil {
			return nil, err
		}
		return &EthereumKeeperRegistry{
			address:     &address,
			client:      e.client,
			registry1_3: instance.(*keeper_registry_wrapper1_3.KeeperRegistry),
		}, err
	case eth_contracts.RegistryVersion_2_0:
		instance, err := e.client.LoadContract("KeeperRegistry", address, func(
			address common.Address,
			backend bind.ContractBackend,
		) (interface{}, error) {
			return keeper_registry_wrapper2_0.NewKeeperRegistry(address, backend)
		})
		if err != nil {
			return nil, err
		}
		return &EthereumKeeperRegistry{
			address:     &address,
			client:      e.client,
			registry2_0: instance.(*keeper_registry_wrapper2_0.KeeperRegistry),
		}, err
	default:
		return nil, fmt.Errorf("keeper registry version %d is not supported", registryVersion)
	}
}

func (e *EthereumContractDeployer) DeployKeeperConsumer(updateInterval *big.Int) (KeeperConsumer, error) {
	address, _, instance, err := e.client.DeployContract("KeeperConsumer", func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		return eth_contracts.DeployKeeperConsumer(auth, backend, updateInterval)
	})
	if err != nil {
		return nil, err
	}
	return &EthereumKeeperConsumer{
		client:   e.client,
		consumer: instance.(*eth_contracts.KeeperConsumer),
		address:  address,
	}, err
}

func (e *EthereumContractDeployer) DeployUpkeepCounter(testRange *big.Int, interval *big.Int) (UpkeepCounter, error) {
	address, _, instance, err := e.client.DeployContract("UpkeepCounter", func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		return eth_contracts.DeployUpkeepCounter(auth, backend, testRange, interval)
	})
	if err != nil {
		return nil, err
	}
	return &EthereumUpkeepCounter{
		client:   e.client,
		consumer: instance.(*eth_contracts.UpkeepCounter),
		address:  address,
	}, err
}

func (e *EthereumContractDeployer) DeployUpkeepPerformCounterRestrictive(testRange *big.Int, averageEligibilityCadence *big.Int) (UpkeepPerformCounterRestrictive, error) {
	address, _, instance, err := e.client.DeployContract("UpkeepPerformCounterRestrictive", func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		return eth_contracts.DeployUpkeepPerformCounterRestrictive(auth, backend, testRange, averageEligibilityCadence)
	})
	if err != nil {
		return nil, err
	}
	return &EthereumUpkeepPerformCounterRestrictive{
		client:   e.client,
		consumer: instance.(*eth_contracts.UpkeepPerformCounterRestrictive),
		address:  address,
	}, err
}

func (e *EthereumContractDeployer) DeployKeeperConsumerPerformance(
	testBlockRange,
	averageCadence,
	checkGasToBurn,
	performGasToBurn *big.Int,
) (KeeperConsumerPerformance, error) {
	address, _, instance, err := e.client.DeployContract("KeeperConsumerPerformance", func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		return eth_contracts.DeployKeeperConsumerPerformance(
			auth,
			backend,
			testBlockRange,
			averageCadence,
			checkGasToBurn,
			performGasToBurn,
		)
	})
	if err != nil {
		return nil, err
	}
	return &EthereumKeeperConsumerPerformance{
		client:   e.client,
		consumer: instance.(*eth_contracts.KeeperConsumerPerformance),
		address:  address,
	}, err
}

func (e *EthereumContractDeployer) DeployKeeperConsumerBenchmark() (AutomationConsumerBenchmark, error) {
	address, _, instance, err := e.client.DeployContract("AutomationConsumerBenchmark", func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		return automation_consumer_benchmark.DeployAutomationConsumerBenchmark(
			auth,
			backend,
		)
	})
	if err != nil {
		return nil, err
	}
	return &EthereumAutomationConsumerBenchmark{
		client:   e.client,
		consumer: instance.(*automation_consumer_benchmark.AutomationConsumerBenchmark),
		address:  address,
	}, err
}

// LoadKeeperConsumerBenchmark returns deployed on given address EthereumAutomationConsumerBenchmark
func (e *EthereumContractDeployer) LoadKeeperConsumerBenchmark(address common.Address) (AutomationConsumerBenchmark, error) {
	instance, err := e.client.LoadContract("AutomationConsumerBenchmark", address, func(
		address common.Address,
		backend bind.ContractBackend,
	) (interface{}, error) {
		return automation_consumer_benchmark.NewAutomationConsumerBenchmark(address, backend)
	})
	if err != nil {
		return nil, err
	}
	return &EthereumAutomationConsumerBenchmark{
		address:  &address,
		client:   e.client,
		consumer: instance.(*automation_consumer_benchmark.AutomationConsumerBenchmark),
	}, err
}

func (e *EthereumContractDeployer) DeployKeeperPerformDataChecker(expectedData []byte) (KeeperPerformDataChecker, error) {
	address, _, instance, err := e.client.DeployContract("PerformDataChecker", func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		return eth_contracts.DeployPerformDataChecker(
			auth,
			backend,
			expectedData,
		)
	})
	if err != nil {
		return nil, err
	}
	return &EthereumKeeperPerformDataCheckerConsumer{
		client:             e.client,
		performDataChecker: instance.(*eth_contracts.PerformDataChecker),
		address:            address,
	}, err
}

// DeployOperatorFactory deploys operator factory contract
func (e *EthereumContractDeployer) DeployOperatorFactory(linkAddr string) (OperatorFactory, error) {
	addr, _, instance, err := e.client.DeployContract("OperatorFactory", func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		return operator_factory.DeployOperatorFactory(auth, backend, common.HexToAddress(linkAddr))
	})
	if err != nil {
		return nil, err
	}
	return &EthereumOperatorFactory{
		address:         addr,
		client:          e.client,
		operatorFactory: instance.(*operator_factory.OperatorFactory),
	}, err
}

// DeployMockAggregatorProxy deploys a mock aggregator proxy contract
func (e *EthereumContractDeployer) DeployMockAggregatorProxy(aggregatorAddr string) (MockAggregatorProxy, error) {
	addr, _, instance, err := e.client.DeployContract("MockAggregatorProxy", func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		return mock_aggregator_proxy.DeployMockAggregatorProxy(auth, backend, common.HexToAddress(aggregatorAddr))
	})
	if err != nil {
		return nil, err
	}
	return &EthereumMockAggregatorProxy{
		address:             addr,
		client:              e.client,
		mockAggregatorProxy: instance.(*mock_aggregator_proxy.MockAggregatorProxy),
	}, err
}

// DeployOffChainAggregator deploys the offchain aggregation contract to the EVM chain
func (e *EthereumContractDeployer) DeployOffchainAggregatorV2(
	linkAddr string,
	offchainOptions OffchainOptions,
) (OffchainAggregatorV2, error) {
	address, _, instance, err := e.client.DeployContract("OffChain Aggregator v2", func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		la := common.HexToAddress(linkAddr)
		return ocr2aggregator.DeployOCR2Aggregator(
			auth,
			backend,
			la,
			offchainOptions.MinimumAnswer,
			offchainOptions.MaximumAnswer,
			offchainOptions.BillingAccessController,
			offchainOptions.RequesterAccessController,
			offchainOptions.Decimals,
			offchainOptions.Description,
		)
	})
	if err != nil {
		return nil, err
	}
	return &EthereumOffchainAggregatorV2{
		client:   e.client,
		contract: instance.(*ocr2aggregator.OCR2Aggregator),
		address:  address,
	}, err
}
