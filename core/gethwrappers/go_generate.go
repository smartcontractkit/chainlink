// Package gethwrappers provides tools for wrapping solidity contracts with
// golang packages, using abigen.
package gethwrappers

// Make sure solidity compiler artifacts are up to date. Only output stdout on failure.
//go:generate ./generation/compile_contracts.sh

//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.6/FluxAggregator.abi ../../contracts/solc/v0.6/FluxAggregator.bin FluxAggregator flux_aggregator_wrapper
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.6/VRF.abi ../../contracts/solc/v0.6/VRF.bin VRF solidity_vrf_wrapper
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.6/VRFTestHelper.abi ../../contracts/solc/v0.6/VRFTestHelper.bin VRFTestHelper solidity_vrf_verifier_wrapper
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.6/VRFCoordinator.abi ../../contracts/solc/v0.6/VRFCoordinator.bin VRFCoordinator solidity_vrf_coordinator_interface
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.6/VRFConsumer.abi ../../contracts/solc/v0.6/VRFConsumer.bin VRFConsumer solidity_vrf_consumer_interface
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.6/VRFRequestIDBaseTestHelper.abi ../../contracts/solc/v0.6/VRFRequestIDBaseTestHelper.bin VRFRequestIDBaseTestHelper solidity_vrf_request_id
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.6/Flags.abi ../../contracts/solc/v0.6/Flags.bin Flags flags_wrapper
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.6/Oracle.abi ../../contracts/solc/v0.6/Oracle.bin Oracle oracle_wrapper
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.6/BlockhashStore.abi ../../contracts/solc/v0.6/BlockhashStore.bin BlockhashStore blockhash_store
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.6/TestAPIConsumer.abi ../../contracts/solc/v0.6/TestAPIConsumer.bin TestAPIConsumer test_api_consumer_wrapper
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.6/MockETHLINKAggregator.abi ../../contracts/solc/v0.6/MockETHLINKAggregator.bin MockETHLINKAggregator mock_ethlink_aggregator_wrapper
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.6/MockGASAggregator.abi ../../contracts/solc/v0.6/MockGASAggregator.bin MockGASAggregator mock_gas_aggregator_wrapper

//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.7/Consumer.abi ../../contracts/solc/v0.7/Consumer.bin Consumer consumer_wrapper
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.7/MultiWordConsumer.abi ../../contracts/solc/v0.7/MultiWordConsumer.bin MultiWordConsumer multiwordconsumer_wrapper
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.7/Operator.abi ../../contracts/solc/v0.7/Operator.bin Operator operator_wrapper
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.7/OperatorFactory.abi ../../contracts/solc/v0.7/OperatorFactory.bin OperatorFactory operator_factory
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.7/AuthorizedForwarder.abi ../../contracts/solc/v0.7/AuthorizedForwarder.bin AuthorizedForwarder authorized_forwarder
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.7/AuthorizedReceiver.abi ../../contracts/solc/v0.7/AuthorizedReceiver.bin AuthorizedReceiver authorized_receiver
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/BatchBlockhashStore.abi ../../contracts/solc/v0.8.6/BatchBlockhashStore.bin BatchBlockhashStore batch_blockhash_store
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/BatchVRFCoordinatorV2.abi ../../contracts/solc/v0.8.6/BatchVRFCoordinatorV2.bin BatchVRFCoordinatorV2 batch_vrf_coordinator_v2
//go:generate go run ./generation/generate/wrap.go OffchainAggregator/OffchainAggregator.abi - OffchainAggregator offchain_aggregator_wrapper

// Automation
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.7/KeeperRegistry1_1.abi ../../contracts/solc/v0.7/KeeperRegistry1_1.bin KeeperRegistry keeper_registry_wrapper1_1
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.7/UpkeepPerformCounterRestrictive.abi ../../contracts/solc/v0.7/UpkeepPerformCounterRestrictive.bin UpkeepPerformCounterRestrictive upkeep_perform_counter_restrictive_wrapper
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.7/UpkeepCounter.abi ../../contracts/solc/v0.7/UpkeepCounter.bin UpkeepCounter upkeep_counter_wrapper
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/CronUpkeepFactory.abi - CronUpkeepFactory cron_upkeep_factory_wrapper
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/CronUpkeep.abi - CronUpkeep cron_upkeep_wrapper
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/KeeperRegistrar.abi ../../contracts/solc/v0.8.6/KeeperRegistrar.bin KeeperRegistrar keeper_registrar_wrapper1_2
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/KeeperRegistry1_2.abi ../../contracts/solc/v0.8.6/KeeperRegistry1_2.bin KeeperRegistry keeper_registry_wrapper1_2
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/TypeAndVersionInterface.abi ../../contracts/solc/v0.8.6/TypeAndVersionInterface.bin TypeAndVersionInterface type_and_version_interface_wrapper
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/KeeperRegistryCheckUpkeepGasUsageWrapper1_2.abi ../../contracts/solc/v0.8.6/KeeperRegistryCheckUpkeepGasUsageWrapper1_2.bin KeeperRegistryCheckUpkeepGasUsageWrapper gas_wrapper
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/KeeperRegistry1_3.abi ../../contracts/solc/v0.8.6/KeeperRegistry1_3.bin KeeperRegistry keeper_registry_wrapper1_3
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/KeeperRegistryLogic1_3.abi ../../contracts/solc/v0.8.6/KeeperRegistryLogic1_3.bin KeeperRegistryLogic keeper_registry_logic1_3
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/KeeperRegistrar2_0.abi ../../contracts/solc/v0.8.6/KeeperRegistrar2_0.bin KeeperRegistrar keeper_registrar_wrapper2_0
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/KeeperRegistry2_0.abi ../../contracts/solc/v0.8.6/KeeperRegistry2_0.bin KeeperRegistry keeper_registry_wrapper2_0
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/KeeperRegistryLogic2_0.abi ../../contracts/solc/v0.8.6/KeeperRegistryLogic2_0.bin KeeperRegistryLogic keeper_registry_logic2_0
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/UpkeepTranscoder.abi ../../contracts/solc/v0.8.6/UpkeepTranscoder.bin UpkeepTranscoder upkeep_transcoder
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/VerifiableLoadUpkeep.abi ../../contracts/solc/v0.8.6/VerifiableLoadUpkeep.bin VerifiableLoadUpkeep verifiable_load_upkeep_wrapper
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/VerifiableLoadMercuryUpkeep.abi ../../contracts/solc/v0.8.6/VerifiableLoadMercuryUpkeep.bin VerifiableLoadMercuryUpkeep verifiable_load_mercury_upkeep_wrapper
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.15/MercuryUpkeep.abi ../../contracts/solc/v0.8.15/MercuryUpkeep.bin MercuryUpkeep mercury_upkeep_wrapper
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.15/FeedLookupCompatibleInterface.abi ../../contracts/solc/v0.8.15/FeedLookupCompatibleInterface.bin FeedLookupCompatibleInterface feed_lookup_compatible_interface
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.16/AutomationConsumerBenchmark.abi ../../contracts/solc/v0.8.16/AutomationConsumerBenchmark.bin AutomationConsumerBenchmark automation_consumer_benchmark
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/KeeperRegistry2_1.abi ../../contracts/solc/v0.8.6/KeeperRegistry2_1.bin KeeperRegistry keeper_registry_wrapper_2_1
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/KeeperRegistryLogicA2_1.abi ../../contracts/solc/v0.8.6/KeeperRegistryLogicA2_1.bin KeeperRegistryLogicA keeper_registry_logic_a_wrapper_2_1
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/KeeperRegistryLogicB2_1.abi ../../contracts/solc/v0.8.6/KeeperRegistryLogicB2_1.bin KeeperRegistryLogicB keeper_registry_logic_b_wrapper_2_1
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/IKeeperRegistryMaster.abi ../../contracts/solc/v0.8.6/IKeeperRegistryMaster.bin IKeeperRegistryMaster i_keeper_registry_master_wrapper_2_1
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/LogUpkeepCounter.abi ../../contracts/solc/v0.8.6/LogUpkeepCounter.bin LogUpkeepCounter log_upkeep_counter_wrapper

// v0.8.6 VRFConsumer
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/VRFCoordinatorMock.abi ../../contracts/solc/v0.8.6/VRFCoordinatorMock.bin VRFCoordinatorMock vrf_coordinator_mock
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/VRFConsumer.abi ../../contracts/solc/v0.8.6/VRFConsumer.bin VRFConsumer solidity_vrf_consumer_interface_v08
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/VRFRequestIDBaseTestHelper.abi ../../contracts/solc/v0.8.6/VRFRequestIDBaseTestHelper.bin VRFRequestIDBaseTestHelper solidity_vrf_request_id_v08
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/VRFOwnerlessConsumerExample.abi ../../contracts/solc/v0.8.6/VRFOwnerlessConsumerExample.bin VRFOwnerlessConsumerExample vrf_ownerless_consumer_example
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/VRFLoadTestOwnerlessConsumer.abi ../../contracts/solc/v0.8.6/VRFLoadTestOwnerlessConsumer.bin VRFLoadTestOwnerlessConsumer vrf_load_test_ownerless_consumer
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/VRFLoadTestExternalSubOwner.abi ../../contracts/solc/v0.8.6/VRFLoadTestExternalSubOwner.bin VRFLoadTestExternalSubOwner vrf_load_test_external_sub_owner
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/VRFV2LoadTestWithMetrics.abi ../../contracts/solc/v0.8.6/VRFV2LoadTestWithMetrics.bin VRFV2LoadTestWithMetrics vrf_load_test_with_metrics

//go:generate go run ./generation/generate_link/wrap_link.go

// VRF V2
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/VRFOwner.abi ../../contracts/solc/v0.8.6/VRFOwner.bin VRFOwner vrf_owner
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/VRFCoordinatorV2.abi ../../contracts/solc/v0.8.6/VRFCoordinatorV2.bin VRFCoordinatorV2 vrf_coordinator_v2
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/VRFConsumerV2.abi ../../contracts/solc/v0.8.6/VRFConsumerV2.bin VRFConsumerV2 vrf_consumer_v2
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/VRFMaliciousConsumerV2.abi ../../contracts/solc/v0.8.6/VRFMaliciousConsumerV2.bin VRFMaliciousConsumerV2 vrf_malicious_consumer_v2
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/VRFTestHelper.abi ../../contracts/solc/v0.8.6/VRFTestHelper.bin VRFV08TestHelper solidity_vrf_v08_verifier_wrapper
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/VRFSingleConsumerExample.abi ../../contracts/solc/v0.8.6/VRFSingleConsumerExample.bin VRFSingleConsumerExample vrf_single_consumer_example
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/VRFExternalSubOwnerExample.abi ../../contracts/solc/v0.8.6/VRFExternalSubOwnerExample.bin VRFExternalSubOwnerExample vrf_external_sub_owner_example
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/VRFV2RevertingExample.abi ../../contracts/solc/v0.8.6/VRFV2RevertingExample.bin VRFV2RevertingExample vrfv2_reverting_example
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/VRFConsumerV2UpgradeableExample.abi ../../contracts/solc/v0.8.6/VRFConsumerV2UpgradeableExample.bin VRFConsumerV2UpgradeableExample vrf_consumer_v2_upgradeable_example
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/VRFV2TransparentUpgradeableProxy.abi ../../contracts/solc/v0.8.6/VRFV2TransparentUpgradeableProxy.bin VRFV2TransparentUpgradeableProxy vrfv2_transparent_upgradeable_proxy
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/VRFV2ProxyAdmin.abi ../../contracts/solc/v0.8.6/VRFV2ProxyAdmin.bin VRFV2ProxyAdmin vrfv2_proxy_admin

// VRF V2 Wrapper
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/VRFV2Wrapper.abi ../../contracts/solc/v0.8.6/VRFV2Wrapper.bin VRFV2Wrapper vrfv2_wrapper
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/VRFV2WrapperInterface.abi ../../contracts/solc/v0.8.6/VRFV2WrapperInterface.bin VRFV2WrapperInterface vrfv2_wrapper_interface
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/VRFV2WrapperConsumerExample.abi ../../contracts/solc/v0.8.6/VRFV2WrapperConsumerExample.bin VRFV2WrapperConsumerExample vrfv2_wrapper_consumer_example

// Keepers X VRF v2
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/KeepersVRFConsumer.abi ../../contracts/solc/v0.8.6/KeepersVRFConsumer.bin KeepersVRFConsumer keepers_vrf_consumer

// Aggregators
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/AggregatorV2V3Interface.abi ../../contracts/solc/v0.8.6/AggregatorV2V3Interface.bin AggregatorV2V3Interface aggregator_v2v3_interface
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/AggregatorV3Interface.abi ../../contracts/solc/v0.8.6/AggregatorV3Interface.bin AggregatorV3Interface aggregator_v3_interface
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/DerivedPriceFeed.abi ../../contracts/solc/v0.8.6/DerivedPriceFeed.bin DerivedPriceFeed derived_price_feed_wrapper
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/MockAggregatorProxy.abi ../../contracts/solc/v0.8.6/MockAggregatorProxy.bin MockAggregatorProxy mock_aggregator_proxy

// Log tester
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/LogEmitter.abi ../../contracts/solc/v0.8.6/LogEmitter.bin LogEmitter log_emitter

// Chainlink Functions (OCR2DR)
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/Functions.abi ../../contracts/solc/v0.8.6/Functions.bin OCR2DR ocr2dr
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/FunctionsClient.abi ../../contracts/solc/v0.8.6/FunctionsClient.bin OCR2DRClient ocr2dr_client
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/FunctionsClientExample.abi ../../contracts/solc/v0.8.6/FunctionsClientExample.bin OCR2DRClientExample ocr2dr_client_example
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/FunctionsOracleWithInit.abi ../../contracts/solc/v0.8.6/FunctionsOracleWithInit.bin OCR2DROracle ocr2dr_oracle
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/FunctionsBillingRegistryWithInit.abi ../../contracts/solc/v0.8.6/FunctionsBillingRegistryWithInit.bin OCR2DRRegistry ocr2dr_registry

// Mercury
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.16/Verifier.abi ../../contracts/solc/v0.8.16/Verifier.bin MercuryVerifier mercury_verifier
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.16/VerifierProxy.abi ../../contracts/solc/v0.8.16/VerifierProxy.bin MercuryVerifierProxy mercury_verifier_proxy
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.16/ExposedVerifier.abi ../../contracts/solc/v0.8.16/ExposedVerifier.bin MercuryExposedVerifier mercury_exposed_verifier

// Mocks that contain only events and functions to emit them
//go:generate go run ./generation/generate_events_mock/create_events_mock_contract.go ../../contracts/solc/v0.8.6/FunctionsOracle.abi ../../contracts/src/v0.8/mocks/FunctionsOracleEventsMock.sol FunctionsOracleEventsMock
//go:generate go run ./generation/generate_events_mock/create_events_mock_contract.go ../../contracts/solc/v0.8.6/FunctionsBillingRegistry.abi ../../contracts/src/v0.8/mocks/FunctionsBillingRegistryEventsMock.sol FunctionsBillingRegistryEventsMock
//go:generate ./generation/compile_event_mock_contract.sh
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/FunctionsOracleEventsMock.abi ../../contracts/solc/v0.8.6/FunctionsOracleEventsMock.bin FunctionsOracleEventsMock functions_oracle_events_mock
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/FunctionsBillingRegistryEventsMock.abi ../../contracts/solc/v0.8.6/FunctionsBillingRegistryEventsMock.bin FunctionsBillingRegistryEventsMock functions_billing_registry_events_mock

// Transmission
//go:generate go generate ./transmission
