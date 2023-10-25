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
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.6/TestAPIConsumer.abi ../../contracts/solc/v0.6/TestAPIConsumer.bin TestAPIConsumer test_api_consumer_wrapper
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.6/MockETHLINKAggregator.abi ../../contracts/solc/v0.6/MockETHLINKAggregator.bin MockETHLINKAggregator mock_ethlink_aggregator_wrapper
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.6/MockGASAggregator.abi ../../contracts/solc/v0.6/MockGASAggregator.bin MockGASAggregator mock_gas_aggregator_wrapper

//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.7/Consumer.abi ../../contracts/solc/v0.7/Consumer.bin Consumer consumer_wrapper
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.7/MultiWordConsumer.abi ../../contracts/solc/v0.7/MultiWordConsumer.bin MultiWordConsumer multiwordconsumer_wrapper
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.19/Operator.abi ../../contracts/solc/v0.8.19/Operator.bin Operator operator_wrapper
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.19/OperatorFactory.abi ../../contracts/solc/v0.8.19/OperatorFactory.bin OperatorFactory operator_factory
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.19/AuthorizedForwarder.abi ../../contracts/solc/v0.8.19/AuthorizedForwarder.bin AuthorizedForwarder authorized_forwarder
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.19/AuthorizedReceiver.abi ../../contracts/solc/v0.8.19/AuthorizedReceiver.bin AuthorizedReceiver authorized_receiver
//go:generate go run ./generation/generate/wrap.go OffchainAggregator/OffchainAggregator.abi - OffchainAggregator offchain_aggregator_wrapper

// Automation
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.7/KeeperRegistry1_1.abi ../../contracts/solc/v0.7/KeeperRegistry1_1.bin KeeperRegistry keeper_registry_wrapper1_1
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.7/KeeperRegistry1_1Mock.abi ../../contracts/solc/v0.7/KeeperRegistry1_1Mock.bin KeeperRegistryMock keeper_registry_wrapper1_1_mock
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.7/UpkeepPerformCounterRestrictive.abi ../../contracts/solc/v0.7/UpkeepPerformCounterRestrictive.bin UpkeepPerformCounterRestrictive upkeep_perform_counter_restrictive_wrapper
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.7/UpkeepCounter.abi ../../contracts/solc/v0.7/UpkeepCounter.bin UpkeepCounter upkeep_counter_wrapper
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/CronUpkeepFactory.abi - CronUpkeepFactory cron_upkeep_factory_wrapper
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/CronUpkeep.abi - CronUpkeep cron_upkeep_wrapper
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/KeeperRegistrar.abi ../../contracts/solc/v0.8.6/KeeperRegistrar.bin KeeperRegistrar keeper_registrar_wrapper1_2
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/KeeperRegistrar1_2Mock.abi ../../contracts/solc/v0.8.6/KeeperRegistrar1_2Mock.bin KeeperRegistrarMock keeper_registrar_wrapper1_2_mock
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/KeeperRegistry1_2.abi ../../contracts/solc/v0.8.6/KeeperRegistry1_2.bin KeeperRegistry keeper_registry_wrapper1_2
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/TypeAndVersionInterface.abi ../../contracts/solc/v0.8.6/TypeAndVersionInterface.bin TypeAndVersionInterface type_and_version_interface_wrapper
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/KeeperRegistryCheckUpkeepGasUsageWrapper1_2.abi ../../contracts/solc/v0.8.6/KeeperRegistryCheckUpkeepGasUsageWrapper1_2.bin KeeperRegistryCheckUpkeepGasUsageWrapper gas_wrapper
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/KeeperRegistryCheckUpkeepGasUsageWrapper1_2Mock.abi ../../contracts/solc/v0.8.6/KeeperRegistryCheckUpkeepGasUsageWrapper1_2Mock.bin KeeperRegistryCheckUpkeepGasUsageWrapperMock gas_wrapper_mock
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/KeeperRegistry1_3.abi ../../contracts/solc/v0.8.6/KeeperRegistry1_3.bin KeeperRegistry keeper_registry_wrapper1_3
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/KeeperRegistryLogic1_3.abi ../../contracts/solc/v0.8.6/KeeperRegistryLogic1_3.bin KeeperRegistryLogic keeper_registry_logic1_3
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/KeeperRegistrar2_0.abi ../../contracts/solc/v0.8.6/KeeperRegistrar2_0.bin KeeperRegistrar keeper_registrar_wrapper2_0
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/KeeperRegistry2_0.abi ../../contracts/solc/v0.8.6/KeeperRegistry2_0.bin KeeperRegistry keeper_registry_wrapper2_0
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/KeeperRegistryLogic2_0.abi ../../contracts/solc/v0.8.6/KeeperRegistryLogic2_0.bin KeeperRegistryLogic keeper_registry_logic2_0
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/UpkeepTranscoder.abi ../../contracts/solc/v0.8.6/UpkeepTranscoder.bin UpkeepTranscoder upkeep_transcoder
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.16/VerifiableLoadUpkeep.abi ../../contracts/solc/v0.8.16/VerifiableLoadUpkeep.bin VerifiableLoadUpkeep verifiable_load_upkeep_wrapper
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.16/VerifiableLoadStreamsLookupUpkeep.abi ../../contracts/solc/v0.8.16/VerifiableLoadStreamsLookupUpkeep.bin VerifiableLoadStreamsLookupUpkeep verifiable_load_streams_lookup_upkeep_wrapper
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.16/VerifiableLoadLogTriggerUpkeep.abi ../../contracts/solc/v0.8.16/VerifiableLoadLogTriggerUpkeep.bin VerifiableLoadLogTriggerUpkeep verifiable_load_log_trigger_upkeep_wrapper
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.16/StreamsLookupUpkeep.abi ../../contracts/solc/v0.8.16/StreamsLookupUpkeep.bin StreamsLookupUpkeep streams_lookup_upkeep_wrapper
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.16/StreamsLookupCompatibleInterface.abi ../../contracts/solc/v0.8.16/StreamsLookupCompatibleInterface.bin StreamsLookupCompatibleInterface streams_lookup_compatible_interface
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.16/AutomationConsumerBenchmark.abi ../../contracts/solc/v0.8.16/AutomationConsumerBenchmark.bin AutomationConsumerBenchmark automation_consumer_benchmark
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.16/AutomationRegistrar2_1.abi ../../contracts/solc/v0.8.16/AutomationRegistrar2_1.bin AutomationRegistrar automation_registrar_wrapper2_1
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.16/KeeperRegistry2_1.abi ../../contracts/solc/v0.8.16/KeeperRegistry2_1.bin KeeperRegistry keeper_registry_wrapper_2_1
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.16/KeeperRegistryLogicA2_1.abi ../../contracts/solc/v0.8.16/KeeperRegistryLogicA2_1.bin KeeperRegistryLogicA keeper_registry_logic_a_wrapper_2_1
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.16/KeeperRegistryLogicB2_1.abi ../../contracts/solc/v0.8.16/KeeperRegistryLogicB2_1.bin KeeperRegistryLogicB keeper_registry_logic_b_wrapper_2_1
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.16/IKeeperRegistryMaster.abi ../../contracts/solc/v0.8.16/IKeeperRegistryMaster.bin IKeeperRegistryMaster i_keeper_registry_master_wrapper_2_1
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.16/ILogAutomation.abi ../../contracts/solc/v0.8.16/ILogAutomation.bin ILogAutomation i_log_automation
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.16/AutomationUtils2_1.abi ../../contracts/solc/v0.8.16/AutomationUtils2_1.bin AutomationUtils automation_utils_2_1
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.16/AutomationForwarderLogic.abi ../../contracts/solc/v0.8.16/AutomationForwarderLogic.bin AutomationForwarderLogic automation_forwarder_logic
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/LogUpkeepCounter.abi ../../contracts/solc/v0.8.6/LogUpkeepCounter.bin LogUpkeepCounter log_upkeep_counter_wrapper
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.16/LogTriggeredStreamsLookup.abi ../../contracts/solc/v0.8.16/LogTriggeredStreamsLookup.bin LogTriggeredStreamsLookup log_triggered_streams_lookup_wrapper
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.16/DummyProtocol.abi ../../contracts/solc/v0.8.16/DummyProtocol.bin DummyProtocol dummy_protocol_wrapper

//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.16/KeeperConsumer.abi ../../contracts/solc/v0.8.16/KeeperConsumer.bin KeeperConsumer keeper_consumer_wrapper
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.16/KeeperConsumerPerformance.abi ../../contracts/solc/v0.8.16/KeeperConsumerPerformance.bin KeeperConsumerPerformance keeper_consumer_performance_wrapper
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.16/PerformDataChecker.abi ../../contracts/solc/v0.8.16/PerformDataChecker.bin PerformDataChecker perform_data_checker_wrapper
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.16/UpkeepCounter.abi ../../contracts/solc/v0.8.16/UpkeepCounter.bin UpkeepCounter upkeep_counter_wrapper
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.16/UpkeepPerformCounterRestrictive.abi ../../contracts/solc/v0.8.16/UpkeepPerformCounterRestrictive.bin UpkeepPerformCounterRestrictive upkeep_perform_counter_restrictive_wrapper

// v0.8.6 VRFConsumer
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/VRFCoordinatorMock.abi ../../contracts/solc/v0.8.6/VRFCoordinatorMock.bin VRFCoordinatorMock vrf_coordinator_mock
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/VRFConsumer.abi ../../contracts/solc/v0.8.6/VRFConsumer.bin VRFConsumer solidity_vrf_consumer_interface_v08
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/VRFRequestIDBaseTestHelper.abi ../../contracts/solc/v0.8.6/VRFRequestIDBaseTestHelper.bin VRFRequestIDBaseTestHelper solidity_vrf_request_id_v08
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/VRFOwnerlessConsumerExample.abi ../../contracts/solc/v0.8.6/VRFOwnerlessConsumerExample.bin VRFOwnerlessConsumerExample vrf_ownerless_consumer_example
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/VRFLoadTestOwnerlessConsumer.abi ../../contracts/solc/v0.8.6/VRFLoadTestOwnerlessConsumer.bin VRFLoadTestOwnerlessConsumer vrf_load_test_ownerless_consumer
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/VRFLoadTestExternalSubOwner.abi ../../contracts/solc/v0.8.6/VRFLoadTestExternalSubOwner.bin VRFLoadTestExternalSubOwner vrf_load_test_external_sub_owner
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/VRFV2LoadTestWithMetrics.abi ../../contracts/solc/v0.8.6/VRFV2LoadTestWithMetrics.bin VRFV2LoadTestWithMetrics vrf_load_test_with_metrics
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/VRFV2OwnerTestConsumer.abi ../../contracts/solc/v0.8.6/VRFV2OwnerTestConsumer.bin VRFV2OwnerTestConsumer vrf_owner_test_consumer
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/VRFv2Consumer.abi ../../contracts/solc/v0.8.6/VRFv2Consumer.bin VRFv2Consumer vrf_v2_consumer_wrapper

//go:generate go run ./generation/generate_link/wrap_link.go

// VRF V2
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/BlockhashStore.abi ../../contracts/solc/v0.8.6/BlockhashStore.bin BlockhashStore blockhash_store
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/BatchBlockhashStore.abi ../../contracts/solc/v0.8.6/BatchBlockhashStore.bin BatchBlockhashStore batch_blockhash_store
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/BatchVRFCoordinatorV2.abi ../../contracts/solc/v0.8.6/BatchVRFCoordinatorV2.bin BatchVRFCoordinatorV2 batch_vrf_coordinator_v2
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
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/ChainSpecificUtilHelper.abi ../../contracts/solc/v0.8.6/ChainSpecificUtilHelper.bin ChainSpecificUtilHelper chain_specific_util_helper

// VRF V2 Wrapper
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/VRFV2Wrapper.abi ../../contracts/solc/v0.8.6/VRFV2Wrapper.bin VRFV2Wrapper vrfv2_wrapper
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/VRFV2WrapperInterface.abi ../../contracts/solc/v0.8.6/VRFV2WrapperInterface.bin VRFV2WrapperInterface vrfv2_wrapper_interface
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/VRFV2WrapperConsumerExample.abi ../../contracts/solc/v0.8.6/VRFV2WrapperConsumerExample.bin VRFV2WrapperConsumerExample vrfv2_wrapper_consumer_example

// Keepers X VRF v2
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/KeepersVRFConsumer.abi ../../contracts/solc/v0.8.6/KeepersVRFConsumer.bin KeepersVRFConsumer keepers_vrf_consumer

// VRF V2Plus
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/IVRFCoordinatorV2PlusInternal.abi ../../contracts/solc/v0.8.6/IVRFCoordinatorV2PlusInternal.bin IVRFCoordinatorV2PlusInternal vrf_coordinator_v2plus_interface
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/BatchVRFCoordinatorV2Plus.abi ../../contracts/solc/v0.8.6/BatchVRFCoordinatorV2Plus.bin BatchVRFCoordinatorV2Plus batch_vrf_coordinator_v2plus
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/TrustedBlockhashStore.abi ../../contracts/solc/v0.8.6/TrustedBlockhashStore.bin TrustedBlockhashStore trusted_blockhash_store
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/VRFV2PlusConsumerExample.abi ../../contracts/solc/v0.8.6/VRFV2PlusConsumerExample.bin VRFV2PlusConsumerExample vrfv2plus_consumer_example
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/VRFCoordinatorV2_5.abi ../../contracts/solc/v0.8.6/VRFCoordinatorV2_5.bin VRFCoordinatorV2_5 vrf_coordinator_v2_5
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/VRFV2PlusWrapper.abi ../../contracts/solc/v0.8.6/VRFV2PlusWrapper.bin VRFV2PlusWrapper vrfv2plus_wrapper
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/VRFV2PlusWrapperConsumerExample.abi ../../contracts/solc/v0.8.6/VRFV2PlusWrapperConsumerExample.bin VRFV2PlusWrapperConsumerExample vrfv2plus_wrapper_consumer_example
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/VRFMaliciousConsumerV2Plus.abi ../../contracts/solc/v0.8.6/VRFMaliciousConsumerV2Plus.bin VRFMaliciousConsumerV2Plus vrf_malicious_consumer_v2_plus
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/VRFV2PlusSingleConsumerExample.abi ../../contracts/solc/v0.8.6/VRFV2PlusSingleConsumerExample.bin VRFV2PlusSingleConsumerExample vrf_v2plus_single_consumer
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/VRFV2PlusExternalSubOwnerExample.abi ../../contracts/solc/v0.8.6/VRFV2PlusExternalSubOwnerExample.bin VRFV2PlusExternalSubOwnerExample vrf_v2plus_sub_owner
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/VRFV2PlusRevertingExample.abi ../../contracts/solc/v0.8.6/VRFV2PlusRevertingExample.bin VRFV2PlusRevertingExample vrfv2plus_reverting_example
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/VRFConsumerV2PlusUpgradeableExample.abi ../../contracts/solc/v0.8.6/VRFConsumerV2PlusUpgradeableExample.bin VRFConsumerV2PlusUpgradeableExample vrf_consumer_v2_plus_upgradeable_example
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/VRFV2PlusClient.abi ../../contracts/solc/v0.8.6/VRFV2PlusClient.bin VRFV2PlusClient vrfv2plus_client
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/VRFCoordinatorV2Plus_V2Example.abi ../../contracts/solc/v0.8.6/VRFCoordinatorV2Plus_V2Example.bin VRFCoordinatorV2Plus_V2Example vrf_coordinator_v2_plus_v2_example
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/VRFV2PlusMaliciousMigrator.abi ../../contracts/solc/v0.8.6/VRFV2PlusMaliciousMigrator.bin VRFV2PlusMaliciousMigrator vrfv2plus_malicious_migrator
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/VRFV2PlusLoadTestWithMetrics.abi ../../contracts/solc/v0.8.6/VRFV2PlusLoadTestWithMetrics.bin VRFV2PlusLoadTestWithMetrics vrf_v2plus_load_test_with_metrics
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/VRFCoordinatorV2PlusUpgradedVersion.abi ../../contracts/solc/v0.8.6/VRFCoordinatorV2PlusUpgradedVersion.bin VRFCoordinatorV2PlusUpgradedVersion vrf_v2plus_upgraded_version
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/VRFV2PlusWrapperLoadTestConsumer.abi ../../contracts/solc/v0.8.6/VRFV2PlusWrapperLoadTestConsumer.bin VRFV2PlusWrapperLoadTestConsumer vrfv2plus_wrapper_load_test_consumer

// Aggregators
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/AggregatorV2V3Interface.abi ../../contracts/solc/v0.8.6/AggregatorV2V3Interface.bin AggregatorV2V3Interface aggregator_v2v3_interface
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/AggregatorV3Interface.abi ../../contracts/solc/v0.8.6/AggregatorV3Interface.bin AggregatorV3Interface aggregator_v3_interface
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/MockAggregatorProxy.abi ../../contracts/solc/v0.8.6/MockAggregatorProxy.bin MockAggregatorProxy mock_aggregator_proxy

// Log tester
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.19/LogEmitter.abi ../../contracts/solc/v0.8.19/LogEmitter.bin LogEmitter log_emitter

// Chainlink Functions
//go:generate go generate ./functions

// Mercury
//go:generate go generate ./llo-feeds

// Shared
//go:generate go generate ./shared

// Mocks that contain only events and functions to emit them
// These contracts are used in testing Atlas flows. The contracts contain no logic, only events, structures, and functions to emit them.
// The flow is as follows:
// 1. Compile all non events mock contracts.
// 2. Generate events mock .sol files based on ABI of compiled contracts.
// 3. Compile events mock contracts. ./generation/compile_event_mock_contract.sh calls contracts/scripts/native_solc_compile_all_events_mock to compile events mock contracts.
// 4. Generate wrappers for events mock contracts.
//go:generate ./generation/compile_event_mock_contract.sh

// Transmission
//go:generate go generate ./transmission
