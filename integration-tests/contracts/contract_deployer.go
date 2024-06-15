package contracts

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/smartcontractkit/chainlink-testing-framework/networks"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/smartcontractkit/libocr/gethwrappers/offchainaggregator"
	"github.com/smartcontractkit/libocr/gethwrappers2/ocr2aggregator"
	ocrConfigHelper "github.com/smartcontractkit/libocr/offchainreporting/confighelper"

	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"

	eth_contracts "github.com/smartcontractkit/chainlink/integration-tests/contracts/ethereum"
	"github.com/smartcontractkit/chainlink/integration-tests/wrappers"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/functions/generated/functions_load_test_client"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/functions/generated/functions_v1_events_mock"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/arbitrum_module"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/automation_consumer_benchmark"
	automationForwarderLogic "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/automation_forwarder_logic"
	registrar21 "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/automation_registrar_wrapper2_1"
	registrylogica22 "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/automation_registry_logic_a_wrapper_2_2"
	registrylogicb22 "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/automation_registry_logic_b_wrapper_2_2"
	registry22 "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/automation_registry_wrapper_2_2"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/chain_module_base"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/flags_wrapper"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/flux_aggregator_wrapper"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/functions_billing_registry_events_mock"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/functions_oracle_events_mock"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/gas_wrapper"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/gas_wrapper_mock"
	iregistry22 "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/i_automation_registry_master_wrapper_2_2"
	ac "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/i_automation_v21_plus_common"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/i_chain_module"
	iregistry21 "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/i_keeper_registry_master_wrapper_2_1"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/keeper_consumer_performance_wrapper"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/keeper_consumer_wrapper"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/keeper_registrar_wrapper1_2"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/keeper_registrar_wrapper1_2_mock"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/keeper_registrar_wrapper2_0"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/keeper_registry_logic1_3"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/keeper_registry_logic2_0"
	registrylogica21 "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/keeper_registry_logic_a_wrapper_2_1"
	registrylogicb21 "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/keeper_registry_logic_b_wrapper_2_1"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/keeper_registry_wrapper1_1"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/keeper_registry_wrapper1_1_mock"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/keeper_registry_wrapper1_2"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/keeper_registry_wrapper1_3"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/keeper_registry_wrapper2_0"
	registry21 "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/keeper_registry_wrapper_2_1"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/link_token_interface"
	le "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/log_emitter"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/log_triggered_streams_lookup_wrapper"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/log_upkeep_counter_wrapper"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/mock_aggregator_proxy"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/mock_ethlink_aggregator_wrapper"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/mock_gas_aggregator_wrapper"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/operator_factory"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/optimism_module"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/oracle_wrapper"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/perform_data_checker_wrapper"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/scroll_module"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/simple_log_upkeep_counter_wrapper"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/streams_lookup_upkeep_wrapper"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/test_api_consumer_wrapper"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/upkeep_counter_wrapper"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/upkeep_perform_counter_restrictive_wrapper"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/upkeep_transcoder"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/llo-feeds/generated/fee_manager"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/llo-feeds/generated/reward_manager"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/llo-feeds/generated/verifier"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/llo-feeds/generated/verifier_proxy"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/shared/generated/werc20_mock"
)

// ContractDeployer is an interface for abstracting the contract deployment methods across network implementations
type ContractDeployer interface {
	DeployAPIConsumer(linkAddr string) (APIConsumer, error)
	DeployOracle(linkAddr string) (Oracle, error)
	DeployFlags(rac string) (Flags, error)
	DeployFluxAggregatorContract(linkAddr string, fluxOptions FluxAggregatorOptions) (FluxAggregator, error)
	DeployLinkTokenContract() (LinkToken, error)
	DeployWERC20Mock() (WERC20Mock, error)
	LoadLinkToken(address common.Address) (LinkToken, error)
	DeployOffChainAggregator(linkAddr string, offchainOptions OffchainOptions) (OffchainAggregator, error)
	LoadOffChainAggregator(address *common.Address) (OffchainAggregator, error)
	DeployMockETHLINKFeed(answer *big.Int) (MockETHLINKFeed, error)
	LoadETHLINKFeed(address common.Address) (MockETHLINKFeed, error)
	DeployMockGasFeed(answer *big.Int) (MockGasFeed, error)
	LoadGasFeed(address common.Address) (MockGasFeed, error)
	DeployKeeperRegistrar(registryVersion eth_contracts.KeeperRegistryVersion, linkAddr string, registrarSettings KeeperRegistrarSettings) (KeeperRegistrar, error)
	LoadKeeperRegistrar(address common.Address, registryVersion eth_contracts.KeeperRegistryVersion) (KeeperRegistrar, error)
	DeployUpkeepTranscoder() (UpkeepTranscoder, error)
	LoadUpkeepTranscoder(address common.Address) (UpkeepTranscoder, error)
	DeployKeeperRegistry(opts *KeeperRegistryOpts) (KeeperRegistry, error)
	LoadKeeperRegistry(address common.Address, registryVersion eth_contracts.KeeperRegistryVersion) (KeeperRegistry, error)
	DeployKeeperConsumer(updateInterval *big.Int) (KeeperConsumer, error)
	DeployAutomationLogTriggerConsumer(testInterval *big.Int) (KeeperConsumer, error)
	DeployAutomationSimpleLogTriggerConsumer(isStreamsLookup bool) (KeeperConsumer, error)
	DeployAutomationStreamsLookupUpkeepConsumer(testRange *big.Int, interval *big.Int, useArbBlock bool, staging bool, verify bool) (KeeperConsumer, error)
	DeployAutomationLogTriggeredStreamsLookupUpkeepConsumer() (KeeperConsumer, error)
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
	DeployOperatorFactory(linkAddr string) (OperatorFactory, error)
	DeployStaking(params eth_contracts.StakingPoolConstructorParams) (Staking, error)
	DeployFunctionsLoadTestClient(router string) (FunctionsLoadTestClient, error)
	DeployFunctionsOracleEventsMock() (FunctionsOracleEventsMock, error)
	DeployFunctionsBillingRegistryEventsMock() (FunctionsBillingRegistryEventsMock, error)
	DeployStakingEventsMock() (StakingEventsMock, error)
	DeployFunctionsV1EventsMock() (FunctionsV1EventsMock, error)
	DeployOffchainAggregatorEventsMock() (OffchainAggregatorEventsMock, error)
	DeployMockAggregatorProxy(aggregatorAddr string) (MockAggregatorProxy, error)
	DeployOffchainAggregatorV2(linkAddr string, offchainOptions OffchainOptions) (OffchainAggregatorV2, error)
	LoadOffChainAggregatorV2(address *common.Address) (OffchainAggregatorV2, error)
	DeployKeeperRegistryCheckUpkeepGasUsageWrapper(keeperRegistryAddr string) (KeeperRegistryCheckUpkeepGasUsageWrapper, error)
	DeployKeeperRegistry11Mock() (KeeperRegistry11Mock, error)
	DeployKeeperRegistrar12Mock() (KeeperRegistrar12Mock, error)
	DeployKeeperGasWrapperMock() (KeeperGasWrapperMock, error)
	DeployMercuryVerifierContract(verifierProxyAddr common.Address) (MercuryVerifier, error)
	DeployMercuryVerifierProxyContract(accessControllerAddr common.Address) (MercuryVerifierProxy, error)
	DeployMercuryFeeManager(linkAddress common.Address, nativeAddress common.Address, proxyAddress common.Address, rewardManagerAddress common.Address) (MercuryFeeManager, error)
	DeployMercuryRewardManager(linkAddress common.Address) (MercuryRewardManager, error)
	DeployLogEmitterContract() (LogEmitter, error)
	DeployMultiCallContract() (common.Address, error)
}

// NewContractDeployer returns an instance of a contract deployer based on the client type
func NewContractDeployer(bcClient blockchain.EVMClient, logger zerolog.Logger) (ContractDeployer, error) {
	switch clientImpl := bcClient.Get().(type) {
	case *blockchain.EthereumClient:
		return NewEthereumContractDeployer(clientImpl, logger), nil
	case *blockchain.KlaytnClient:
		return &KlaytnContractDeployer{NewEthereumContractDeployer(clientImpl, logger)}, nil
	case *blockchain.MetisClient:
		return &MetisContractDeployer{NewEthereumContractDeployer(clientImpl, logger)}, nil
	case *blockchain.ArbitrumClient:
		return &MetisContractDeployer{NewEthereumContractDeployer(clientImpl, logger)}, nil
	case *blockchain.OptimismClient:
		return &OptimismContractDeployer{NewEthereumContractDeployer(clientImpl, logger)}, nil
	case *blockchain.RSKClient:
		return &RSKContractDeployer{NewEthereumContractDeployer(clientImpl, logger)}, nil
	case *blockchain.PolygonClient:
		return &PolygonContractDeployer{NewEthereumContractDeployer(clientImpl, logger)}, nil
	case *blockchain.CeloClient:
		return &CeloContractDeployer{NewEthereumContractDeployer(clientImpl, logger)}, nil
	case *blockchain.QuorumClient:
		return &QuorumContractDeployer{NewEthereumContractDeployer(clientImpl, logger)}, nil
	case *blockchain.BSCClient:
		return &BSCContractDeployer{NewEthereumContractDeployer(clientImpl, logger)}, nil
	case *blockchain.ScrollClient:
		return &ScrollContractDeployer{NewEthereumContractDeployer(clientImpl, logger)}, nil
	case *blockchain.PolygonZkEvmClient:
		return &PolygonZkEvmContractDeployer{NewEthereumContractDeployer(clientImpl, logger)}, nil
	case *blockchain.LineaClient:
		return &LineaContractDeployer{NewEthereumContractDeployer(clientImpl, logger)}, nil
	case *blockchain.FantomClient:
		return &FantomContractDeployer{NewEthereumContractDeployer(clientImpl, logger)}, nil
	case *blockchain.KromaClient:
		return &KromaContractDeployer{NewEthereumContractDeployer(clientImpl, logger)}, nil
	case *blockchain.WeMixClient:
		return &WeMixContractDeployer{NewEthereumContractDeployer(clientImpl, logger)}, nil
	case *blockchain.GnosisClient:
		return &GnosisContractDeployer{NewEthereumContractDeployer(clientImpl, logger)}, nil
	}
	return nil, errors.New("unknown blockchain client implementation for contract deployer, register blockchain client in NewContractDeployer")
}

// EthereumContractDeployer provides the implementations for deploying ETH (EVM) based contracts
type EthereumContractDeployer struct {
	client blockchain.EVMClient
	l      zerolog.Logger
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

type BSCContractDeployer struct {
	*EthereumContractDeployer
}

type ScrollContractDeployer struct {
	*EthereumContractDeployer
}

type PolygonZkEvmContractDeployer struct {
	*EthereumContractDeployer
}

type LineaContractDeployer struct {
	*EthereumContractDeployer
}

type FantomContractDeployer struct {
	*EthereumContractDeployer
}

type KromaContractDeployer struct {
	*EthereumContractDeployer
}

type WeMixContractDeployer struct {
	*EthereumContractDeployer
}

type GnosisContractDeployer struct {
	*EthereumContractDeployer
}

// NewEthereumContractDeployer returns an instantiated instance of the ETH contract deployer
func NewEthereumContractDeployer(ethClient blockchain.EVMClient, logger zerolog.Logger) *EthereumContractDeployer {
	return &EthereumContractDeployer{
		client: ethClient,
		l:      logger,
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
	return &LegacyEthereumFluxAggregator{
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

func (e *EthereumContractDeployer) DeployFunctionsLoadTestClient(router string) (FunctionsLoadTestClient, error) {
	address, _, instance, err := e.client.DeployContract("FunctionsLoadTestClient", func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		return functions_load_test_client.DeployFunctionsLoadTestClient(auth, backend, common.HexToAddress(router))
	})
	if err != nil {
		return nil, err
	}
	return &LegacyEthereumFunctionsLoadTestClient{
		client:   e.client,
		instance: instance.(*functions_load_test_client.FunctionsLoadTestClient),
		address:  *address,
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

func (e *EthereumContractDeployer) DeployStakingEventsMock() (StakingEventsMock, error) {
	address, _, instance, err := e.client.DeployContract("StakingEventsMock", func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		return eth_contracts.DeployStakingEventsMock(auth, backend)
	})
	if err != nil {
		return nil, err
	}
	return &LegacyEthereumStakingEventsMock{
		client:     e.client,
		eventsMock: instance.(*eth_contracts.StakingEventsMock),
		address:    address,
	}, nil
}

func (e *EthereumContractDeployer) DeployFunctionsV1EventsMock() (FunctionsV1EventsMock, error) {
	address, _, instance, err := e.client.DeployContract("FunctionsV1EventsMock", func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		return functions_v1_events_mock.DeployFunctionsV1EventsMock(auth, backend)
	})
	if err != nil {
		return nil, err
	}
	return &LegacyEthereumFunctionsV1EventsMock{
		client:     e.client,
		eventsMock: instance.(*functions_v1_events_mock.FunctionsV1EventsMock),
		address:    address,
	}, nil
}

func (e *EthereumContractDeployer) DeployKeeperRegistry11Mock() (KeeperRegistry11Mock, error) {
	address, _, instance, err := e.client.DeployContract("KeeperRegistry11Mock", func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		return keeper_registry_wrapper1_1_mock.DeployKeeperRegistryMock(auth, backend)
	})
	if err != nil {
		return nil, err
	}
	return &LegacyEthereumKeeperRegistry11Mock{
		client:       e.client,
		registryMock: instance.(*keeper_registry_wrapper1_1_mock.KeeperRegistryMock),
		address:      address,
	}, nil
}

func (e *EthereumContractDeployer) DeployKeeperRegistrar12Mock() (KeeperRegistrar12Mock, error) {
	address, _, instance, err := e.client.DeployContract("KeeperRegistrar12Mock", func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		return keeper_registrar_wrapper1_2_mock.DeployKeeperRegistrarMock(auth, backend)
	})
	if err != nil {
		return nil, err
	}
	return &LegacyEthereumKeeperRegistrar12Mock{
		client:        e.client,
		registrarMock: instance.(*keeper_registrar_wrapper1_2_mock.KeeperRegistrarMock),
		address:       address,
	}, nil
}

func (e *EthereumContractDeployer) DeployKeeperGasWrapperMock() (KeeperGasWrapperMock, error) {
	address, _, instance, err := e.client.DeployContract("KeeperGasWrapperMock", func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		return gas_wrapper_mock.DeployKeeperRegistryCheckUpkeepGasUsageWrapperMock(auth, backend)
	})
	if err != nil {
		return nil, err
	}
	return &LegacyEthereumKeeperGasWrapperMock{
		client:         e.client,
		gasWrapperMock: instance.(*gas_wrapper_mock.KeeperRegistryCheckUpkeepGasUsageWrapperMock),
		address:        address,
	}, nil
}

func (e *EthereumContractDeployer) DeployOffchainAggregatorEventsMock() (OffchainAggregatorEventsMock, error) {
	address, _, instance, err := e.client.DeployContract("OffchainAggregatorEventsMock", func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		return eth_contracts.DeployOffchainAggregatorEventsMock(auth, backend)
	})
	if err != nil {
		return nil, err
	}
	return &EthereumOffchainAggregatorEventsMock{
		client:     e.client,
		eventsMock: instance.(*eth_contracts.OffchainAggregatorEventsMock),
		address:    address,
	}, nil
}

// DeployLinkTokenContract deploys a Link Token contract to an EVM chain
func (e *EthereumContractDeployer) DeployLinkTokenContract() (LinkToken, error) {
	linkTokenAddress, _, instance, err := e.client.DeployContract("LINK Token", func(
		auth *bind.TransactOpts,
		_ bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		return link_token_interface.DeployLinkToken(auth, wrappers.MustNewWrappedContractBackend(e.client, nil))
	})
	if err != nil {
		return nil, err
	}

	return &LegacyEthereumLinkToken{
		client:   e.client,
		instance: instance.(*link_token_interface.LinkToken),
		address:  *linkTokenAddress,
		l:        e.l,
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
	return &LegacyEthereumLinkToken{
		address:  address,
		client:   e.client,
		instance: instance.(*link_token_interface.LinkToken),
		l:        e.l,
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
	return &LegacyEthereumOffchainAggregator{
		client:  e.client,
		ocr:     instance.(*offchainaggregator.OffchainAggregator),
		address: address,
		l:       e.l,
	}, err
}

// LoadOffChainAggregator loads an already deployed offchain aggregator contract
func (e *EthereumContractDeployer) LoadOffChainAggregator(address *common.Address) (OffchainAggregator, error) {
	instance, err := e.client.LoadContract("OffChainAggregator", *address, func(
		address common.Address,
		backend bind.ContractBackend,
	) (interface{}, error) {
		return offchainaggregator.NewOffchainAggregator(address, backend)
	})
	if err != nil {
		return nil, err
	}
	return &LegacyEthereumOffchainAggregator{
		address: address,
		client:  e.client,
		ocr:     instance.(*offchainaggregator.OffchainAggregator),
		l:       e.l,
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
	return &LegacyEthereumAPIConsumer{
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
	return &LegacyEthereumOracle{
		address: addr,
		client:  e.client,
		oracle:  instance.(*oracle_wrapper.Oracle),
	}, err
}

func (e *EthereumContractDeployer) DeployMockETHLINKFeed(answer *big.Int) (MockETHLINKFeed, error) {
	address, _, instance, err := e.client.DeployContract("MockETHLINKAggregator", func(
		auth *bind.TransactOpts,
		_ bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		return mock_ethlink_aggregator_wrapper.DeployMockETHLINKAggregator(auth, wrappers.MustNewWrappedContractBackend(e.client, nil), answer)
	})
	if err != nil {
		return nil, err
	}
	return &LegacyEthereumMockETHLINKFeed{
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
	return &LegacyEthereumMockETHLINKFeed{
		address: &address,
		client:  e.client,
		feed:    instance.(*mock_ethlink_aggregator_wrapper.MockETHLINKAggregator),
	}, err
}

func (e *EthereumContractDeployer) DeployMockGasFeed(answer *big.Int) (MockGasFeed, error) {
	address, _, instance, err := e.client.DeployContract("MockGasFeed", func(
		auth *bind.TransactOpts,
		_ bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		return mock_gas_aggregator_wrapper.DeployMockGASAggregator(auth, wrappers.MustNewWrappedContractBackend(e.client, nil), answer)
	})
	if err != nil {
		return nil, err
	}
	return &LegacyEthereumMockGASFeed{
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
	return &LegacyEthereumMockGASFeed{
		address: &address,
		client:  e.client,
		feed:    instance.(*mock_gas_aggregator_wrapper.MockGASAggregator),
	}, err
}

func (e *EthereumContractDeployer) DeployUpkeepTranscoder() (UpkeepTranscoder, error) {
	address, _, instance, err := e.client.DeployContract("UpkeepTranscoder", func(
		opts *bind.TransactOpts,
		_ bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		return upkeep_transcoder.DeployUpkeepTranscoder(opts, wrappers.MustNewWrappedContractBackend(e.client, nil))
	})

	if err != nil {
		return nil, err
	}

	return &LegacyEthereumUpkeepTranscoder{
		client:     e.client,
		transcoder: instance.(*upkeep_transcoder.UpkeepTranscoder),
		address:    address,
	}, err
}

func (e *EthereumContractDeployer) LoadUpkeepTranscoder(address common.Address) (UpkeepTranscoder, error) {
	instance, err := e.client.LoadContract("UpkeepTranscoder", address, func(
		address common.Address,
		backend bind.ContractBackend,
	) (interface{}, error) {
		return upkeep_transcoder.NewUpkeepTranscoder(address, backend)
	})

	if err != nil {
		return nil, err
	}

	return &LegacyEthereumUpkeepTranscoder{
		client:     e.client,
		transcoder: instance.(*upkeep_transcoder.UpkeepTranscoder),
		address:    &address,
	}, err
}

func (e *EthereumContractDeployer) DeployKeeperRegistrar(registryVersion eth_contracts.KeeperRegistryVersion, linkAddr string,
	registrarSettings KeeperRegistrarSettings) (KeeperRegistrar, error) {

	if registryVersion == eth_contracts.RegistryVersion_2_0 {
		// deploy registrar 2.0
		address, _, instance, err := e.client.DeployContract("KeeperRegistrar", func(
			opts *bind.TransactOpts,
			_ bind.ContractBackend,
		) (common.Address, *types.Transaction, interface{}, error) {
			return keeper_registrar_wrapper2_0.DeployKeeperRegistrar(opts, wrappers.MustNewWrappedContractBackend(e.client, nil), common.HexToAddress(linkAddr), registrarSettings.AutoApproveConfigType,
				registrarSettings.AutoApproveMaxAllowed, common.HexToAddress(registrarSettings.RegistryAddr), registrarSettings.MinLinkJuels)
		})

		if err != nil {
			return nil, err
		}

		return &LegacyEthereumKeeperRegistrar{
			client:      e.client,
			registrar20: instance.(*keeper_registrar_wrapper2_0.KeeperRegistrar),
			address:     address,
		}, err
	} else if registryVersion == eth_contracts.RegistryVersion_2_1 || registryVersion == eth_contracts.RegistryVersion_2_2 { // both 2.1 and 2.2 registry use registrar 2.1
		// deploy registrar 2.1
		address, _, instance, err := e.client.DeployContract("AutomationRegistrar", func(
			opts *bind.TransactOpts,
			_ bind.ContractBackend,
		) (common.Address, *types.Transaction, interface{}, error) {
			// set default TriggerType to 0(conditional), AutoApproveConfigType to 2(auto approve enabled), AutoApproveMaxAllowed to 1000
			triggerConfigs := []registrar21.AutomationRegistrar21InitialTriggerConfig{
				{TriggerType: 0, AutoApproveType: registrarSettings.AutoApproveConfigType,
					AutoApproveMaxAllowed: uint32(registrarSettings.AutoApproveMaxAllowed)},
				{TriggerType: 1, AutoApproveType: registrarSettings.AutoApproveConfigType,
					AutoApproveMaxAllowed: uint32(registrarSettings.AutoApproveMaxAllowed)},
			}

			return registrar21.DeployAutomationRegistrar(
				opts,
				wrappers.MustNewWrappedContractBackend(e.client, nil),
				common.HexToAddress(linkAddr),
				common.HexToAddress(registrarSettings.RegistryAddr),
				registrarSettings.MinLinkJuels,
				triggerConfigs)
		})

		if err != nil {
			return nil, err
		}

		return &LegacyEthereumKeeperRegistrar{
			client:      e.client,
			registrar21: instance.(*registrar21.AutomationRegistrar),
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

	return &LegacyEthereumKeeperRegistrar{
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
		return &LegacyEthereumKeeperRegistrar{
			address:   &address,
			client:    e.client,
			registrar: instance.(*keeper_registrar_wrapper1_2.KeeperRegistrar),
		}, err
	} else if registryVersion == eth_contracts.RegistryVersion_2_0 {
		instance, err := e.client.LoadContract("KeeperRegistrar", address, func(
			address common.Address,
			backend bind.ContractBackend,
		) (interface{}, error) {
			return keeper_registrar_wrapper2_0.NewKeeperRegistrar(address, backend)
		})
		if err != nil {
			return nil, err
		}
		return &LegacyEthereumKeeperRegistrar{
			address:     &address,
			client:      e.client,
			registrar20: instance.(*keeper_registrar_wrapper2_0.KeeperRegistrar),
		}, err
	}
	instance, err := e.client.LoadContract("AutomationRegistrar", address, func(
		address common.Address,
		backend bind.ContractBackend,
	) (interface{}, error) {
		return registrar21.NewAutomationRegistrar(address, backend)
	})
	if err != nil {
		return nil, err
	}
	return &LegacyEthereumKeeperRegistrar{
		address:     &address,
		client:      e.client,
		registrar21: instance.(*registrar21.AutomationRegistrar),
	}, err
}

func (e *EthereumContractDeployer) DeployKeeperRegistry(
	opts *KeeperRegistryOpts,
) (KeeperRegistry, error) {
	var mode uint8
	switch e.client.GetChainID().Int64() {
	//Arbitrum payment model
	case networks.ArbitrumMainnet.ChainID, networks.ArbitrumSepolia.ChainID:
		mode = uint8(1)
	//Optimism payment model
	case networks.OptimismMainnet.ChainID, networks.OptimismSepolia.ChainID:
		mode = uint8(2)
	//Base
	case networks.BaseMainnet.ChainID, networks.BaseSepolia.ChainID:
		mode = uint8(2)
	default:
		mode = uint8(0)
	}
	registryGasOverhead := big.NewInt(80000)
	switch opts.RegistryVersion {
	case eth_contracts.RegistryVersion_1_0, eth_contracts.RegistryVersion_1_1:
		address, _, instance, err := e.client.DeployContract("KeeperRegistry1_1", func(
			auth *bind.TransactOpts,
			_ bind.ContractBackend,
		) (common.Address, *types.Transaction, interface{}, error) {
			return keeper_registry_wrapper1_1.DeployKeeperRegistry(
				auth,
				wrappers.MustNewWrappedContractBackend(e.client, nil),
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
		return &LegacyEthereumKeeperRegistry{
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
			_ bind.ContractBackend,
		) (common.Address, *types.Transaction, interface{}, error) {
			return keeper_registry_wrapper1_2.DeployKeeperRegistry(
				auth,
				wrappers.MustNewWrappedContractBackend(e.client, nil),
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
		return &LegacyEthereumKeeperRegistry{
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
			_ bind.ContractBackend,
		) (common.Address, *types.Transaction, interface{}, error) {
			return keeper_registry_logic1_3.DeployKeeperRegistryLogic(
				auth,
				wrappers.MustNewWrappedContractBackend(e.client, nil),
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
			_ bind.ContractBackend,
		) (common.Address, *types.Transaction, interface{}, error) {
			return keeper_registry_wrapper1_3.DeployKeeperRegistry(
				auth,
				wrappers.MustNewWrappedContractBackend(e.client, nil),
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
		return &LegacyEthereumKeeperRegistry{
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
			_ bind.ContractBackend,
		) (common.Address, *types.Transaction, interface{}, error) {
			return keeper_registry_logic2_0.DeployKeeperRegistryLogic(
				auth,
				wrappers.MustNewWrappedContractBackend(e.client, nil),
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
			_ bind.ContractBackend,
		) (common.Address, *types.Transaction, interface{}, error) {

			return keeper_registry_wrapper2_0.DeployKeeperRegistry(
				auth,
				wrappers.MustNewWrappedContractBackend(e.client, nil),
				*logicAddress,
			)
		})
		if err != nil {
			return nil, err
		}
		return &LegacyEthereumKeeperRegistry{
			client:      e.client,
			version:     eth_contracts.RegistryVersion_2_0,
			registry2_0: instance.(*keeper_registry_wrapper2_0.KeeperRegistry),
			address:     address,
		}, err

	case eth_contracts.RegistryVersion_2_1:
		automationForwarderLogicAddr, err := deployAutomationForwarderLogic(e.client)
		if err != nil {
			return nil, err
		}

		registryLogicBAddr, _, _, err := e.client.DeployContract("KeeperRegistryLogicB2_1", func(
			auth *bind.TransactOpts,
			_ bind.ContractBackend,
		) (common.Address, *types.Transaction, interface{}, error) {

			return registrylogicb21.DeployKeeperRegistryLogicB(
				auth,
				wrappers.MustNewWrappedContractBackend(e.client, nil),
				mode,
				common.HexToAddress(opts.LinkAddr),
				common.HexToAddress(opts.ETHFeedAddr),
				common.HexToAddress(opts.GasFeedAddr),
				*automationForwarderLogicAddr,
			)
		})
		if err != nil {
			return nil, err
		}

		if err := e.client.WaitForEvents(); err != nil {
			return nil, err
		}

		registryLogicAAddr, _, _, err := e.client.DeployContract("KeeperRegistryLogicA2_1", func(
			auth *bind.TransactOpts,
			_ bind.ContractBackend,
		) (common.Address, *types.Transaction, interface{}, error) {

			return registrylogica21.DeployKeeperRegistryLogicA(
				auth,
				wrappers.MustNewWrappedContractBackend(e.client, nil),
				*registryLogicBAddr,
			)
		})
		if err != nil {
			return nil, err
		}
		if err := e.client.WaitForEvents(); err != nil {
			return nil, err
		}

		address, _, _, err := e.client.DeployContract("KeeperRegistry2_1", func(
			auth *bind.TransactOpts,
			_ bind.ContractBackend,
		) (common.Address, *types.Transaction, interface{}, error) {
			return registry21.DeployKeeperRegistry(
				auth,
				wrappers.MustNewWrappedContractBackend(e.client, nil),
				*registryLogicAAddr,
			)
		})
		if err != nil {
			return nil, err
		}
		if err := e.client.WaitForEvents(); err != nil {
			return nil, err
		}

		registryMaster, err := iregistry21.NewIKeeperRegistryMaster(
			*address,
			wrappers.MustNewWrappedContractBackend(e.client, nil),
		)
		if err != nil {
			return nil, err
		}
		return &LegacyEthereumKeeperRegistry{
			client:      e.client,
			version:     eth_contracts.RegistryVersion_2_1,
			registry2_1: registryMaster,
			address:     address,
		}, err
	case eth_contracts.RegistryVersion_2_2:
		var chainModuleAddr *common.Address
		var err error
		chainId := e.client.GetChainID().Int64()

		if chainId == networks.ScrollMainnet.ChainID || chainId == networks.ScrollSepolia.ChainID {
			chainModuleAddr, _, _, err = e.client.DeployContract("ScrollModule", func(
				auth *bind.TransactOpts,
				_ bind.ContractBackend,
			) (common.Address, *types.Transaction, interface{}, error) {
				return scroll_module.DeployScrollModule(auth, wrappers.MustNewWrappedContractBackend(e.client, nil))
			})
		} else if chainId == networks.ArbitrumMainnet.ChainID || chainId == networks.ArbitrumSepolia.ChainID {
			chainModuleAddr, _, _, err = e.client.DeployContract("ArbitrumModule", func(
				auth *bind.TransactOpts,
				_ bind.ContractBackend,
			) (common.Address, *types.Transaction, interface{}, error) {
				return arbitrum_module.DeployArbitrumModule(auth, wrappers.MustNewWrappedContractBackend(e.client, nil))
			})
		} else if chainId == networks.OptimismMainnet.ChainID || chainId == networks.OptimismSepolia.ChainID {
			chainModuleAddr, _, _, err = e.client.DeployContract("OptimismModule", func(
				auth *bind.TransactOpts,
				_ bind.ContractBackend,
			) (common.Address, *types.Transaction, interface{}, error) {
				return optimism_module.DeployOptimismModule(auth, wrappers.MustNewWrappedContractBackend(e.client, nil))
			})
		} else {
			chainModuleAddr, _, _, err = e.client.DeployContract("ChainModuleBase", func(
				auth *bind.TransactOpts,
				_ bind.ContractBackend,
			) (common.Address, *types.Transaction, interface{}, error) {
				return chain_module_base.DeployChainModuleBase(auth, wrappers.MustNewWrappedContractBackend(e.client, nil))
			})
		}
		if err != nil {
			return nil, err
		}
		if err = e.client.WaitForEvents(); err != nil {
			return nil, err
		}

		automationForwarderLogicAddr, err := deployAutomationForwarderLogic(e.client)
		if err != nil {
			return nil, err
		}

		var allowedReadOnlyAddress common.Address
		if chainId == networks.PolygonZkEvmMainnet.ChainID || chainId == networks.PolygonZkEvmCardona.ChainID {
			allowedReadOnlyAddress = common.HexToAddress("0x1111111111111111111111111111111111111111")
		} else {
			allowedReadOnlyAddress = common.HexToAddress("0x0000000000000000000000000000000000000000")
		}
		registryLogicBAddr, _, _, err := e.client.DeployContract("AutomationRegistryLogicB2_2", func(
			auth *bind.TransactOpts,
			_ bind.ContractBackend,
		) (common.Address, *types.Transaction, interface{}, error) {

			return registrylogicb22.DeployAutomationRegistryLogicB(
				auth,
				wrappers.MustNewWrappedContractBackend(e.client, nil),
				common.HexToAddress(opts.LinkAddr),
				common.HexToAddress(opts.ETHFeedAddr),
				common.HexToAddress(opts.GasFeedAddr),
				*automationForwarderLogicAddr,
				allowedReadOnlyAddress,
			)
		})
		if err != nil {
			return nil, err
		}

		if err := e.client.WaitForEvents(); err != nil {
			return nil, err
		}

		registryLogicAAddr, _, _, err := e.client.DeployContract("AutomationRegistryLogicA2_2", func(
			auth *bind.TransactOpts,
			backend bind.ContractBackend,
		) (common.Address, *types.Transaction, interface{}, error) {

			return registrylogica22.DeployAutomationRegistryLogicA(
				auth,
				wrappers.MustNewWrappedContractBackend(e.client, nil),
				*registryLogicBAddr,
			)
		})
		if err != nil {
			return nil, err
		}
		if err := e.client.WaitForEvents(); err != nil {
			return nil, err
		}

		address, _, _, err := e.client.DeployContract("AutomationRegistry2_2", func(
			auth *bind.TransactOpts,
			backend bind.ContractBackend,
		) (common.Address, *types.Transaction, interface{}, error) {
			return registry22.DeployAutomationRegistry(
				auth,
				wrappers.MustNewWrappedContractBackend(e.client, nil),
				*registryLogicAAddr,
			)
		})
		if err != nil {
			return nil, err
		}
		if err := e.client.WaitForEvents(); err != nil {
			return nil, err
		}

		registryMaster, err := iregistry22.NewIAutomationRegistryMaster(
			*address,
			wrappers.MustNewWrappedContractBackend(e.client, nil),
		)
		if err != nil {
			return nil, err
		}

		chainModule, err := i_chain_module.NewIChainModule(
			*chainModuleAddr,
			wrappers.MustNewWrappedContractBackend(e.client, nil),
		)
		if err != nil {
			return nil, err
		}

		return &LegacyEthereumKeeperRegistry{
			client:      e.client,
			version:     eth_contracts.RegistryVersion_2_2,
			registry2_2: registryMaster,
			chainModule: chainModule,
			address:     address,
		}, err
	default:
		return nil, fmt.Errorf("keeper registry version %d is not supported", opts.RegistryVersion)
	}
}

func deployAutomationForwarderLogic(client blockchain.EVMClient) (*common.Address, error) {
	automationForwarderLogicAddr, _, _, err := client.DeployContract("automationForwarderLogic", func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		return automationForwarderLogic.DeployAutomationForwarderLogic(auth, backend)
	})

	if err != nil {
		return nil, err
	}

	if err := client.WaitForEvents(); err != nil {
		return nil, err
	}
	return automationForwarderLogicAddr, nil
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
		return &LegacyEthereumKeeperRegistry{
			address:     &address,
			client:      e.client,
			registry1_1: instance.(*keeper_registry_wrapper1_1.KeeperRegistry),
			version:     registryVersion,
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
		return &LegacyEthereumKeeperRegistry{
			address:     &address,
			client:      e.client,
			registry1_2: instance.(*keeper_registry_wrapper1_2.KeeperRegistry),
			version:     registryVersion,
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
		return &LegacyEthereumKeeperRegistry{
			address:     &address,
			client:      e.client,
			registry1_3: instance.(*keeper_registry_wrapper1_3.KeeperRegistry),
			version:     registryVersion,
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
		return &LegacyEthereumKeeperRegistry{
			address:     &address,
			client:      e.client,
			registry2_0: instance.(*keeper_registry_wrapper2_0.KeeperRegistry),
			version:     registryVersion,
		}, err
	case eth_contracts.RegistryVersion_2_1:
		instance, err := e.client.LoadContract("KeeperRegistry", address, func(
			address common.Address,
			backend bind.ContractBackend,
		) (interface{}, error) {
			return ac.NewIAutomationV21PlusCommon(address, backend)
		})
		if err != nil {
			return nil, err
		}
		return &LegacyEthereumKeeperRegistry{
			address:     &address,
			client:      e.client,
			registry2_1: instance.(*iregistry21.IKeeperRegistryMaster),
			version:     registryVersion,
		}, err
	case eth_contracts.RegistryVersion_2_2: // why the contract name is not the same as the actual contract name?
		instance, err := e.client.LoadContract("AutomationRegistry", address, func(
			address common.Address,
			backend bind.ContractBackend,
		) (interface{}, error) {
			return iregistry22.NewIAutomationRegistryMaster(address, backend)
		})
		if err != nil {
			return nil, err
		}
		return &LegacyEthereumKeeperRegistry{
			address:     &address,
			client:      e.client,
			registry2_2: instance.(*iregistry22.IAutomationRegistryMaster),
			version:     registryVersion,
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
		return keeper_consumer_wrapper.DeployKeeperConsumer(auth, backend, updateInterval)
	})
	if err != nil {
		return nil, err
	}
	return &EthereumKeeperConsumer{
		client:   e.client,
		consumer: instance.(*keeper_consumer_wrapper.KeeperConsumer),
		address:  address,
	}, err
}

func (e *EthereumContractDeployer) DeployAutomationLogTriggerConsumer(testInterval *big.Int) (KeeperConsumer, error) {
	address, _, instance, err := e.client.DeployContract("LogUpkeepCounter", func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		return log_upkeep_counter_wrapper.DeployLogUpkeepCounter(
			auth, wrappers.MustNewWrappedContractBackend(e.client, nil), testInterval,
		)
	})
	if err != nil {
		return nil, err
	}
	return &LegacyEthereumAutomationLogCounterConsumer{
		client:   e.client,
		consumer: instance.(*log_upkeep_counter_wrapper.LogUpkeepCounter),
		address:  address,
	}, err
}

func (e *EthereumContractDeployer) DeployAutomationSimpleLogTriggerConsumer(isStreamsLookup bool) (KeeperConsumer, error) {
	address, _, instance, err := e.client.DeployContract("SimpleLogUpkeepCounter", func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		return simple_log_upkeep_counter_wrapper.DeploySimpleLogUpkeepCounter(
			auth, backend, isStreamsLookup,
		)
	})
	if err != nil {
		return nil, err
	}
	return &LegacyEthereumAutomationSimpleLogCounterConsumer{
		client:   e.client,
		consumer: instance.(*simple_log_upkeep_counter_wrapper.SimpleLogUpkeepCounter),
		address:  address,
	}, err
}

func (e *EthereumContractDeployer) DeployAutomationStreamsLookupUpkeepConsumer(testRange *big.Int, interval *big.Int, useArbBlock bool, staging bool, verify bool) (KeeperConsumer, error) {
	address, _, instance, err := e.client.DeployContract("StreamsLookupUpkeep", func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		return streams_lookup_upkeep_wrapper.DeployStreamsLookupUpkeep(
			auth, wrappers.MustNewWrappedContractBackend(e.client, nil), testRange, interval, useArbBlock, staging, verify,
		)
	})
	if err != nil {
		return nil, err
	}
	return &LegacyEthereumAutomationStreamsLookupUpkeepConsumer{
		client:   e.client,
		consumer: instance.(*streams_lookup_upkeep_wrapper.StreamsLookupUpkeep),
		address:  address,
	}, err
}

func (e *EthereumContractDeployer) DeployAutomationLogTriggeredStreamsLookupUpkeepConsumer() (KeeperConsumer, error) {
	address, _, instance, err := e.client.DeployContract("LogTriggeredStreamsLookupUpkeep", func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		return log_triggered_streams_lookup_wrapper.DeployLogTriggeredStreamsLookup(
			auth, wrappers.MustNewWrappedContractBackend(e.client, nil), false, false, false,
		)
	})
	if err != nil {
		return nil, err
	}
	return &LegacyEthereumAutomationLogTriggeredStreamsLookupUpkeepConsumer{
		client:   e.client,
		consumer: instance.(*log_triggered_streams_lookup_wrapper.LogTriggeredStreamsLookup),
		address:  address,
	}, err
}

func (e *EthereumContractDeployer) DeployUpkeepCounter(testRange *big.Int, interval *big.Int) (UpkeepCounter, error) {
	address, _, instance, err := e.client.DeployContract("UpkeepCounter", func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		return upkeep_counter_wrapper.DeployUpkeepCounter(auth, wrappers.MustNewWrappedContractBackend(e.client, nil), testRange, interval)
	})
	if err != nil {
		return nil, err
	}
	return &LegacyEthereumUpkeepCounter{
		client:   e.client,
		consumer: instance.(*upkeep_counter_wrapper.UpkeepCounter),
		address:  address,
	}, err
}

func (e *EthereumContractDeployer) DeployUpkeepPerformCounterRestrictive(testRange *big.Int, averageEligibilityCadence *big.Int) (UpkeepPerformCounterRestrictive, error) {
	address, _, instance, err := e.client.DeployContract("UpkeepPerformCounterRestrictive", func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		return upkeep_perform_counter_restrictive_wrapper.DeployUpkeepPerformCounterRestrictive(auth, backend, testRange, averageEligibilityCadence)
	})
	if err != nil {
		return nil, err
	}
	return &LegacyEthereumUpkeepPerformCounterRestrictive{
		client:   e.client,
		consumer: instance.(*upkeep_perform_counter_restrictive_wrapper.UpkeepPerformCounterRestrictive),
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
		return keeper_consumer_performance_wrapper.DeployKeeperConsumerPerformance(
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
	return &LegacyEthereumKeeperConsumerPerformance{
		client:   e.client,
		consumer: instance.(*keeper_consumer_performance_wrapper.KeeperConsumerPerformance),
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
	return &LegacyEthereumAutomationConsumerBenchmark{
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
	return &LegacyEthereumAutomationConsumerBenchmark{
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
		return perform_data_checker_wrapper.DeployPerformDataChecker(
			auth,
			backend,
			expectedData,
		)
	})
	if err != nil {
		return nil, err
	}
	return &LegacyEthereumKeeperPerformDataCheckerConsumer{
		client:             e.client,
		performDataChecker: instance.(*perform_data_checker_wrapper.PerformDataChecker),
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
	return &LegacyEthereumOperatorFactory{
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

func (e *EthereumContractDeployer) DeployKeeperRegistryCheckUpkeepGasUsageWrapper(keeperRegistryAddr string) (KeeperRegistryCheckUpkeepGasUsageWrapper, error) {
	addr, _, instance, err := e.client.DeployContract("KeeperRegistryCheckUpkeepGasUsageWrapper", func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		return gas_wrapper.DeployKeeperRegistryCheckUpkeepGasUsageWrapper(auth, backend, common.HexToAddress(keeperRegistryAddr))
	})
	if err != nil {
		return nil, err
	}
	return &EthereumKeeperRegistryCheckUpkeepGasUsageWrapper{
		address:         addr,
		client:          e.client,
		gasUsageWrapper: instance.(*gas_wrapper.KeeperRegistryCheckUpkeepGasUsageWrapper),
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
	return &LegacyEthereumOffchainAggregatorV2{
		client:   e.client,
		contract: instance.(*ocr2aggregator.OCR2Aggregator),
		address:  address,
		l:        e.l,
	}, err
}

// LoadOffChainAggregatorV2 loads an already deployed offchain aggregator v2 contract
func (e *EthereumContractDeployer) LoadOffChainAggregatorV2(address *common.Address) (OffchainAggregatorV2, error) {
	instance, err := e.client.LoadContract("OffChainAggregatorV2", *address, func(
		address common.Address,
		backend bind.ContractBackend,
	) (interface{}, error) {
		return ocr2aggregator.NewOCR2Aggregator(address, backend)
	})
	if err != nil {
		return nil, err
	}
	return &LegacyEthereumOffchainAggregatorV2{
		client:   e.client,
		contract: instance.(*ocr2aggregator.OCR2Aggregator),
		address:  address,
		l:        e.l,
	}, err
}

func (e *EthereumContractDeployer) DeployMercuryVerifierContract(verifierProxyAddr common.Address) (MercuryVerifier, error) {
	address, _, instance, err := e.client.DeployContract("Mercury Verifier", func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		return verifier.DeployVerifier(auth, backend, verifierProxyAddr)
	})
	if err != nil {
		return nil, err
	}
	return &EthereumMercuryVerifier{
		client:   e.client,
		instance: instance.(*verifier.Verifier),
		address:  *address,
		l:        e.l,
	}, err
}

func (e *EthereumContractDeployer) DeployMercuryVerifierProxyContract(accessControllerAddr common.Address) (MercuryVerifierProxy, error) {
	address, _, instance, err := e.client.DeployContract("Mercury Verifier Proxy", func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		return verifier_proxy.DeployVerifierProxy(auth, backend, accessControllerAddr)
	})
	if err != nil {
		return nil, err
	}
	return &EthereumMercuryVerifierProxy{
		client:   e.client,
		instance: instance.(*verifier_proxy.VerifierProxy),
		address:  *address,
		l:        e.l,
	}, err
}

func (e *EthereumContractDeployer) DeployMercuryFeeManager(linkAddress common.Address, nativeAddress common.Address, proxyAddress common.Address, rewardManagerAddress common.Address) (MercuryFeeManager, error) {
	address, _, instance, err := e.client.DeployContract("Mercury Fee Manager", func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		return fee_manager.DeployFeeManager(auth, backend, linkAddress, nativeAddress, proxyAddress, rewardManagerAddress)
	})
	if err != nil {
		return nil, err
	}
	return &EthereumMercuryFeeManager{
		client:   e.client,
		instance: instance.(*fee_manager.FeeManager),
		address:  *address,
		l:        e.l,
	}, err
}

func (e *EthereumContractDeployer) DeployMercuryRewardManager(linkAddress common.Address) (MercuryRewardManager, error) {
	address, _, instance, err := e.client.DeployContract("Mercury Reward Manager", func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		return reward_manager.DeployRewardManager(auth, backend, linkAddress)
	})
	if err != nil {
		return nil, err
	}
	return &EthereumMercuryRewardManager{
		client:   e.client,
		instance: instance.(*reward_manager.RewardManager),
		address:  *address,
		l:        e.l,
	}, err
}

func (e *EthereumContractDeployer) DeployWERC20Mock() (WERC20Mock, error) {
	address, _, instance, err := e.client.DeployContract("WERC20 Mock", func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		return werc20_mock.DeployWERC20Mock(auth, backend)
	})
	if err != nil {
		return nil, err
	}
	return &EthereumWERC20Mock{
		client:   e.client,
		instance: instance.(*werc20_mock.WERC20Mock),
		address:  *address,
		l:        e.l,
	}, err
}

func (e *EthereumContractDeployer) DeployLogEmitterContract() (LogEmitter, error) {
	address, _, instance, err := e.client.DeployContract("Log Emitter", func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		return le.DeployLogEmitter(auth, backend)
	})
	if err != nil {
		return nil, err
	}
	return &LegacyLogEmitterContract{
		client:   e.client,
		instance: instance.(*le.LogEmitter),
		address:  *address,
		l:        e.l,
	}, err
}

func (e *EthereumContractDeployer) DeployMultiCallContract() (common.Address, error) {
	multiCallABI, err := abi.JSON(strings.NewReader(MultiCallABI))
	if err != nil {
		return common.Address{}, err
	}
	address, tx, _, err := e.client.DeployContract("MultiCall Contract", func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		address, tx, contract, err := bind.DeployContract(auth, multiCallABI, common.FromHex(MultiCallBIN), backend)
		if err != nil {
			return common.Address{}, nil, nil, err
		}
		return address, tx, contract, err
	})
	if err != nil {
		return common.Address{}, err
	}
	r, err := bind.WaitMined(context.Background(), e.client.DeployBackend(), tx)
	if err != nil {
		return common.Address{}, err
	}
	if r.Status != types.ReceiptStatusSuccessful {
		return common.Address{}, fmt.Errorf("deploy multicall failed")
	}
	return *address, nil

}
