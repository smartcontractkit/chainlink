// Package gethwrappers provides tools for wrapping solidity contracts with
// golang packages, using abigen.
package gethwrappers

// Make sure solidity compiler artifacts are up-to-date. Only output stdout on failure.
//go:generate ./generation/compile_contracts.sh

//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.19/Operator/Operator.abi ../../contracts/solc/v0.8.19/Operator/Operator.bin Operator operator_wrapper
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.19/OperatorFactory/OperatorFactory.abi ../../contracts/solc/v0.8.19/OperatorFactory/OperatorFactory.bin OperatorFactory operator_factory
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.19/AuthorizedForwarder/AuthorizedForwarder.abi ../../contracts/solc/v0.8.19/AuthorizedForwarder/AuthorizedForwarder.bin AuthorizedForwarder authorized_forwarder
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.19/AuthorizedReceiver/AuthorizedReceiver.abi ../../contracts/solc/v0.8.19/AuthorizedReceiver/AuthorizedReceiver.bin AuthorizedReceiver authorized_receiver
//go:generate go run ./generation/generate/wrap.go OffchainAggregator/OffchainAggregator.abi - OffchainAggregator offchain_aggregator_wrapper

// Automation
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/CronUpkeepFactory/CronUpkeepFactory.abi - CronUpkeepFactory cron_upkeep_factory_wrapper
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/CronUpkeepFactory/CronUpkeep.abi - CronUpkeep cron_upkeep_wrapper
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/KeeperRegistrar1_2/KeeperRegistrar.abi ../../contracts/solc/v0.8.6/KeeperRegistrar1_2/KeeperRegistrar.bin KeeperRegistrar keeper_registrar_wrapper1_2
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/KeeperRegistrar1_2Mock/KeeperRegistrar1_2Mock.abi ../../contracts/solc/v0.8.6/KeeperRegistrar1_2Mock/KeeperRegistrar1_2Mock.bin KeeperRegistrarMock keeper_registrar_wrapper1_2_mock
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/KeeperRegistry1_2/KeeperRegistry1_2.abi ../../contracts/solc/v0.8.6/KeeperRegistry1_2/KeeperRegistry1_2.bin KeeperRegistry keeper_registry_wrapper1_2
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/KeeperRegistry1_2/TypeAndVersionInterface.abi ../../contracts/solc/v0.8.6/KeeperRegistry1_2/TypeAndVersionInterface.bin TypeAndVersionInterface type_and_version_interface_wrapper
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/KeeperRegistryCheckUpkeepGasUsageWrapper1_2/KeeperRegistryCheckUpkeepGasUsageWrapper1_2.abi ../../contracts/solc/v0.8.6/KeeperRegistryCheckUpkeepGasUsageWrapper1_2/KeeperRegistryCheckUpkeepGasUsageWrapper1_2.bin KeeperRegistryCheckUpkeepGasUsageWrapper gas_wrapper
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/KeeperRegistryCheckUpkeepGasUsageWrapper1_2Mock/KeeperRegistryCheckUpkeepGasUsageWrapper1_2Mock.abi ../../contracts/solc/v0.8.6/KeeperRegistryCheckUpkeepGasUsageWrapper1_2Mock/KeeperRegistryCheckUpkeepGasUsageWrapper1_2Mock.bin KeeperRegistryCheckUpkeepGasUsageWrapperMock gas_wrapper_mock
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/KeeperRegistry1_3/KeeperRegistry1_3.abi ../../contracts/solc/v0.8.6/KeeperRegistry1_3/KeeperRegistry1_3.bin KeeperRegistry keeper_registry_wrapper1_3
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/KeeperRegistryLogic1_3/KeeperRegistryLogic1_3.abi ../../contracts/solc/v0.8.6/KeeperRegistryLogic1_3/KeeperRegistryLogic1_3.bin KeeperRegistryLogic keeper_registry_logic1_3
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/KeeperRegistrar2_0/KeeperRegistrar2_0.abi ../../contracts/solc/v0.8.6/KeeperRegistrar2_0/KeeperRegistrar2_0.bin KeeperRegistrar keeper_registrar_wrapper2_0
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/KeeperRegistry2_0/KeeperRegistry2_0.abi ../../contracts/solc/v0.8.6/KeeperRegistry2_0/KeeperRegistry2_0.bin KeeperRegistry keeper_registry_wrapper2_0
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/KeeperRegistryLogic2_0/KeeperRegistryLogic2_0.abi ../../contracts/solc/v0.8.6/KeeperRegistryLogic2_0/KeeperRegistryLogic2_0.bin KeeperRegistryLogic keeper_registry_logic2_0
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/UpkeepTranscoder/UpkeepTranscoder.abi ../../contracts/solc/v0.8.6/UpkeepTranscoder/UpkeepTranscoder.bin UpkeepTranscoder upkeep_transcoder
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.16/VerifiableLoadUpkeep/VerifiableLoadUpkeep.abi ../../contracts/solc/v0.8.16/VerifiableLoadUpkeep/VerifiableLoadUpkeep.bin VerifiableLoadUpkeep verifiable_load_upkeep_wrapper
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.16/VerifiableLoadStreamsLookupUpkeep/VerifiableLoadStreamsLookupUpkeep.abi ../../contracts/solc/v0.8.16/VerifiableLoadStreamsLookupUpkeep/VerifiableLoadStreamsLookupUpkeep.bin VerifiableLoadStreamsLookupUpkeep verifiable_load_streams_lookup_upkeep_wrapper
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.16/VerifiableLoadLogTriggerUpkeep/VerifiableLoadLogTriggerUpkeep.abi ../../contracts/solc/v0.8.16/VerifiableLoadLogTriggerUpkeep/VerifiableLoadLogTriggerUpkeep.bin VerifiableLoadLogTriggerUpkeep verifiable_load_log_trigger_upkeep_wrapper
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.16/StreamsLookupUpkeep/StreamsLookupUpkeep.abi ../../contracts/solc/v0.8.16/StreamsLookupUpkeep/StreamsLookupUpkeep.bin StreamsLookupUpkeep streams_lookup_upkeep_wrapper
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.16/StreamsLookupCompatibleInterface/StreamsLookupCompatibleInterface.abi ../../contracts/solc/v0.8.16/StreamsLookupCompatibleInterface/StreamsLookupCompatibleInterface.bin StreamsLookupCompatibleInterface streams_lookup_compatible_interface
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.16/AutomationConsumerBenchmark/AutomationConsumerBenchmark.abi ../../contracts/solc/v0.8.16/AutomationConsumerBenchmark/AutomationConsumerBenchmark.bin AutomationConsumerBenchmark automation_consumer_benchmark
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.16/AutomationRegistrar2_1/AutomationRegistrar2_1.abi ../../contracts/solc/v0.8.16/AutomationRegistrar2_1/AutomationRegistrar2_1.bin AutomationRegistrar automation_registrar_wrapper2_1
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.16/KeeperRegistry2_1/KeeperRegistry2_1.abi ../../contracts/solc/v0.8.16/KeeperRegistry2_1/KeeperRegistry2_1.bin KeeperRegistry keeper_registry_wrapper_2_1
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.16/KeeperRegistryLogicA2_1/KeeperRegistryLogicA2_1.abi ../../contracts/solc/v0.8.16/KeeperRegistryLogicA2_1/KeeperRegistryLogicA2_1.bin KeeperRegistryLogicA keeper_registry_logic_a_wrapper_2_1
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.16/KeeperRegistryLogicB2_1/KeeperRegistryLogicB2_1.abi ../../contracts/solc/v0.8.16/KeeperRegistryLogicB2_1/KeeperRegistryLogicB2_1.bin KeeperRegistryLogicB keeper_registry_logic_b_wrapper_2_1
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.16/IKeeperRegistryMaster/IKeeperRegistryMaster.abi ../../contracts/solc/v0.8.16/IKeeperRegistryMaster/IKeeperRegistryMaster.bin IKeeperRegistryMaster i_keeper_registry_master_wrapper_2_1
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.16/AutomationUtils2_1/AutomationUtils2_1.abi ../../contracts/solc/v0.8.16/AutomationUtils2_1/AutomationUtils2_1.bin AutomationUtils automation_utils_2_1
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.19/AutomationRegistry2_2/AutomationRegistry2_2.abi ../../contracts/solc/v0.8.19/AutomationRegistry2_2/AutomationRegistry2_2.bin AutomationRegistry automation_registry_wrapper_2_2
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.19/AutomationRegistryLogicA2_2/AutomationRegistryLogicA2_2.abi ../../contracts/solc/v0.8.19/AutomationRegistryLogicA2_2/AutomationRegistryLogicA2_2.bin AutomationRegistryLogicA automation_registry_logic_a_wrapper_2_2
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.19/AutomationRegistryLogicB2_2/AutomationRegistryLogicB2_2.abi ../../contracts/solc/v0.8.19/AutomationRegistryLogicB2_2/AutomationRegistryLogicB2_2.bin AutomationRegistryLogicB automation_registry_logic_b_wrapper_2_2
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.19/IAutomationRegistryMaster/IAutomationRegistryMaster.abi ../../contracts/solc/v0.8.19/IAutomationRegistryMaster/IAutomationRegistryMaster.bin IAutomationRegistryMaster i_automation_registry_master_wrapper_2_2
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.19/AutomationCompatibleUtils/AutomationCompatibleUtils.abi ../../contracts/solc/v0.8.19/AutomationCompatibleUtils/AutomationCompatibleUtils.bin AutomationCompatibleUtils automation_compatible_utils
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.19/AutomationUtils2_2/AutomationUtils2_2.abi ../../contracts/solc/v0.8.19/AutomationUtils2_2/AutomationUtils2_2.bin AutomationUtils automation_utils_2_2
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.19/AutomationRegistrar2_3/AutomationRegistrar2_3.abi ../../contracts/solc/v0.8.19/AutomationRegistrar2_3/AutomationRegistrar2_3.bin AutomationRegistrar automation_registrar_wrapper2_3
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.19/AutomationRegistry2_3/AutomationRegistry2_3.abi ../../contracts/solc/v0.8.19/AutomationRegistry2_3/AutomationRegistry2_3.bin AutomationRegistry automation_registry_wrapper_2_3
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.19/AutomationRegistryLogicA2_3/AutomationRegistryLogicA2_3.abi ../../contracts/solc/v0.8.19/AutomationRegistryLogicA2_3/AutomationRegistryLogicA2_3.bin AutomationRegistryLogicA automation_registry_logic_a_wrapper_2_3
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.19/AutomationRegistryLogicB2_3/AutomationRegistryLogicB2_3.abi ../../contracts/solc/v0.8.19/AutomationRegistryLogicB2_3/AutomationRegistryLogicB2_3.bin AutomationRegistryLogicB automation_registry_logic_b_wrapper_2_3
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.19/IAutomationRegistryMaster2_3/IAutomationRegistryMaster2_3.abi ../../contracts/solc/v0.8.19/IAutomationRegistryMaster2_3/IAutomationRegistryMaster2_3.bin IAutomationRegistryMaster2_3 i_automation_registry_master_wrapper_2_3
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.19/AutomationUtils2_3/AutomationUtils2_3.abi ../../contracts/solc/v0.8.19/AutomationUtils2_3/AutomationUtils2_3.bin AutomationUtils automation_utils_2_3
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.19/ArbitrumModule/ArbitrumModule.abi ../../contracts/solc/v0.8.19/ArbitrumModule/ArbitrumModule.bin ArbitrumModule arbitrum_module
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.19/ChainModuleBase/ChainModuleBase.abi ../../contracts/solc/v0.8.19/ChainModuleBase/ChainModuleBase.bin ChainModuleBase chain_module_base
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.19/OptimismModule/OptimismModule.abi ../../contracts/solc/v0.8.19/OptimismModule/OptimismModule.bin OptimismModule optimism_module
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.19/ScrollModule/ScrollModule.abi ../../contracts/solc/v0.8.19/ScrollModule/ScrollModule.bin ScrollModule scroll_module
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.19/IChainModule/IChainModule.abi ../../contracts/solc/v0.8.19/IChainModule/IChainModule.bin IChainModule i_chain_module
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.19/IAutomationV21PlusCommon/IAutomationV21PlusCommon.abi ../../contracts/solc/v0.8.19/IAutomationV21PlusCommon/IAutomationV21PlusCommon.bin IAutomationV21PlusCommon i_automation_v21_plus_common

//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.16/ILogAutomation/ILogAutomation.abi ../../contracts/solc/v0.8.16/ILogAutomation/ILogAutomation.bin ILogAutomation i_log_automation
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.16/AutomationForwarderLogic/AutomationForwarderLogic.abi ../../contracts/solc/v0.8.16/AutomationForwarderLogic/AutomationForwarderLogic.bin AutomationForwarderLogic automation_forwarder_logic
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/LogUpkeepCounter/LogUpkeepCounter.abi ../../contracts/solc/v0.8.6/LogUpkeepCounter/LogUpkeepCounter.bin LogUpkeepCounter log_upkeep_counter_wrapper
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/SimpleLogUpkeepCounter/SimpleLogUpkeepCounter.abi ../../contracts/solc/v0.8.6/SimpleLogUpkeepCounter/SimpleLogUpkeepCounter.bin SimpleLogUpkeepCounter simple_log_upkeep_counter_wrapper
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.16/LogTriggeredStreamsLookup/LogTriggeredStreamsLookup.abi ../../contracts/solc/v0.8.16/LogTriggeredStreamsLookup/LogTriggeredStreamsLookup.bin LogTriggeredStreamsLookup log_triggered_streams_lookup_wrapper
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.16/DummyProtocol/DummyProtocol.abi ../../contracts/solc/v0.8.16/DummyProtocol/DummyProtocol.bin DummyProtocol dummy_protocol_wrapper

//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.16/KeeperConsumer/KeeperConsumer.abi ../../contracts/solc/v0.8.16/KeeperConsumer/KeeperConsumer.bin KeeperConsumer keeper_consumer_wrapper
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.16/KeeperConsumerPerformance/KeeperConsumerPerformance.abi ../../contracts/solc/v0.8.16/KeeperConsumerPerformance/KeeperConsumerPerformance.bin KeeperConsumerPerformance keeper_consumer_performance_wrapper
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.16/PerformDataChecker/PerformDataChecker.abi ../../contracts/solc/v0.8.16/PerformDataChecker/PerformDataChecker.bin PerformDataChecker perform_data_checker_wrapper
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.16/UpkeepCounter/UpkeepCounter.abi ../../contracts/solc/v0.8.16/UpkeepCounter/UpkeepCounter.bin UpkeepCounter upkeep_counter_wrapper
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.16/UpkeepPerformCounterRestrictive/UpkeepPerformCounterRestrictive.abi ../../contracts/solc/v0.8.16/UpkeepPerformCounterRestrictive/UpkeepPerformCounterRestrictive.bin UpkeepPerformCounterRestrictive upkeep_perform_counter_restrictive_wrapper

// v0.8.6 VRFConsumer
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/VRFCoordinatorMock/VRFCoordinatorMock.abi ../../contracts/solc/v0.8.6/VRFCoordinatorMock/VRFCoordinatorMock.bin VRFCoordinatorMock vrf_coordinator_mock
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/VRFConsumer/VRFConsumer.abi ../../contracts/solc/v0.8.6/VRFConsumer/VRFConsumer.bin VRFConsumer solidity_vrf_consumer_interface_v08
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/VRFRequestIDBaseTestHelper/VRFRequestIDBaseTestHelper.abi ../../contracts/solc/v0.8.6/VRFRequestIDBaseTestHelper/VRFRequestIDBaseTestHelper.bin VRFRequestIDBaseTestHelper solidity_vrf_request_id_v08
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/VRFOwnerlessConsumerExample/VRFOwnerlessConsumerExample.abi ../../contracts/solc/v0.8.6/VRFOwnerlessConsumerExample/VRFOwnerlessConsumerExample.bin VRFOwnerlessConsumerExample vrf_ownerless_consumer_example
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/VRFLoadTestOwnerlessConsumer/VRFLoadTestOwnerlessConsumer.abi ../../contracts/solc/v0.8.6/VRFLoadTestOwnerlessConsumer/VRFLoadTestOwnerlessConsumer.bin VRFLoadTestOwnerlessConsumer vrf_load_test_ownerless_consumer
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/VRFLoadTestExternalSubOwner/VRFLoadTestExternalSubOwner.abi ../../contracts/solc/v0.8.6/VRFLoadTestExternalSubOwner/VRFLoadTestExternalSubOwner.bin VRFLoadTestExternalSubOwner vrf_load_test_external_sub_owner
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/VRFV2LoadTestWithMetrics/VRFV2LoadTestWithMetrics.abi ../../contracts/solc/v0.8.6/VRFV2LoadTestWithMetrics/VRFV2LoadTestWithMetrics.bin VRFV2LoadTestWithMetrics vrf_load_test_with_metrics
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/VRFV2OwnerTestConsumer/VRFV2OwnerTestConsumer.abi ../../contracts/solc/v0.8.6/VRFV2OwnerTestConsumer/VRFV2OwnerTestConsumer.bin VRFV2OwnerTestConsumer vrf_owner_test_consumer
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/VRFv2Consumer/VRFv2Consumer.abi ../../contracts/solc/v0.8.6/VRFv2Consumer/VRFv2Consumer.bin VRFv2Consumer vrf_v2_consumer_wrapper

//go:generate go run ./generation/generate_link/wrap_link.go

// VRF V2
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/BatchVRFCoordinatorV2/BatchVRFCoordinatorV2.abi ../../contracts/solc/v0.8.6/BatchVRFCoordinatorV2/BatchVRFCoordinatorV2.bin BatchVRFCoordinatorV2 batch_vrf_coordinator_v2
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/VRFOwner/VRFOwner.abi ../../contracts/solc/v0.8.6/VRFOwner/VRFOwner.bin VRFOwner vrf_owner
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/VRFCoordinatorV2/VRFCoordinatorV2.abi ../../contracts/solc/v0.8.6/VRFCoordinatorV2/VRFCoordinatorV2.bin VRFCoordinatorV2 vrf_coordinator_v2
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/VRFConsumerV2/VRFConsumerV2.abi ../../contracts/solc/v0.8.6/VRFConsumerV2/VRFConsumerV2.bin VRFConsumerV2 vrf_consumer_v2
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/VRFMaliciousConsumerV2/VRFMaliciousConsumerV2.abi ../../contracts/solc/v0.8.6/VRFMaliciousConsumerV2/VRFMaliciousConsumerV2.bin VRFMaliciousConsumerV2 vrf_malicious_consumer_v2

//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/VRFTestHelper/VRFTestHelper.abi ../../contracts/solc/v0.8.6/VRFTestHelper/VRFTestHelper.bin VRFV08TestHelper solidity_vrf_v08_verifier_wrapper
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/VRFSingleConsumerExample/VRFSingleConsumerExample.abi ../../contracts/solc/v0.8.6/VRFSingleConsumerExample/VRFSingleConsumerExample.bin VRFSingleConsumerExample vrf_single_consumer_example

//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/VRFExternalSubOwnerExample/VRFExternalSubOwnerExample.abi ../../contracts/solc/v0.8.6/VRFExternalSubOwnerExample/VRFExternalSubOwnerExample.bin VRFExternalSubOwnerExample vrf_external_sub_owner_example

//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/VRFV2RevertingExample/VRFV2RevertingExample.abi ../../contracts/solc/v0.8.6/VRFV2RevertingExample/VRFV2RevertingExample.bin VRFV2RevertingExample vrfv2_reverting_example

//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/VRFConsumerV2UpgradeableExample/VRFConsumerV2UpgradeableExample.abi ../../contracts/solc/v0.8.6/VRFConsumerV2UpgradeableExample/VRFConsumerV2UpgradeableExample.bin VRFConsumerV2UpgradeableExample vrf_consumer_v2_upgradeable_example

//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/VRFV2TransparentUpgradeableProxy/VRFV2TransparentUpgradeableProxy.abi ../../contracts/solc/v0.8.6/VRFV2TransparentUpgradeableProxy/VRFV2TransparentUpgradeableProxy.bin VRFV2TransparentUpgradeableProxy vrfv2_transparent_upgradeable_proxy
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/VRFV2ProxyAdmin/VRFV2ProxyAdmin.abi ../../contracts/solc/v0.8.6/VRFV2ProxyAdmin/VRFV2ProxyAdmin.bin VRFV2ProxyAdmin vrfv2_proxy_admin
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/ChainSpecificUtilHelper/ChainSpecificUtilHelper.abi ../../contracts/solc/v0.8.6/ChainSpecificUtilHelper/ChainSpecificUtilHelper.bin ChainSpecificUtilHelper chain_specific_util_helper
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/VRFCoordinatorTestV2/VRFCoordinatorTestV2.abi ../../contracts/solc/v0.8.6/VRFCoordinatorTestV2/VRFCoordinatorTestV2.bin VRFCoordinatorTestV2 vrf_coordinator_test_v2
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/VRFMockETHLINKAggregator/VRFMockETHLINKAggregator.abi ../../contracts/solc/v0.8.6/VRFMockETHLINKAggregator/VRFMockETHLINKAggregator.bin VRFMockETHLINKAggregator vrf_mock_ethlink_aggregator

// VRF V2 Wrapper
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/VRFV2Wrapper/VRFV2Wrapper.abi ../../contracts/solc/v0.8.6/VRFV2Wrapper/VRFV2Wrapper.bin VRFV2Wrapper vrfv2_wrapper
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/VRFV2WrapperInterface/VRFV2WrapperInterface.abi ../../contracts/solc/v0.8.6/VRFV2WrapperInterface/VRFV2WrapperInterface.bin VRFV2WrapperInterface vrfv2_wrapper_interface
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/VRFV2WrapperConsumerExample/VRFV2WrapperConsumerExample.abi ../../contracts/solc/v0.8.6/VRFV2WrapperConsumerExample/VRFV2WrapperConsumerExample.bin VRFV2WrapperConsumerExample vrfv2_wrapper_consumer_example
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/VRFV2WrapperLoadTestConsumer/VRFV2WrapperLoadTestConsumer.abi ../../contracts/solc/v0.8.6/VRFV2WrapperLoadTestConsumer/VRFV2WrapperLoadTestConsumer.bin VRFV2WrapperLoadTestConsumer vrfv2_wrapper_load_test_consumer

// Keepers X VRF v2
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/KeepersVRFConsumer/KeepersVRFConsumer.abi ../../contracts/solc/v0.8.6/KeepersVRFConsumer/KeepersVRFConsumer.bin KeepersVRFConsumer keepers_vrf_consumer

// VRF V2Plus
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.19/IVRFCoordinatorV2PlusInternal/IVRFCoordinatorV2PlusInternal.abi ../../contracts/solc/v0.8.19/IVRFCoordinatorV2PlusInternal/IVRFCoordinatorV2PlusInternal.bin IVRFCoordinatorV2PlusInternal vrf_coordinator_v2plus_interface
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.19/BatchVRFCoordinatorV2Plus/BatchVRFCoordinatorV2Plus.abi ../../contracts/solc/v0.8.19/BatchVRFCoordinatorV2Plus/BatchVRFCoordinatorV2Plus.bin BatchVRFCoordinatorV2Plus batch_vrf_coordinator_v2plus
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.19/TrustedBlockhashStore/TrustedBlockhashStore.abi ../../contracts/solc/v0.8.19/TrustedBlockhashStore/TrustedBlockhashStore.bin TrustedBlockhashStore trusted_blockhash_store
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.19/VRFV2PlusConsumerExample/VRFV2PlusConsumerExample.abi ../../contracts/solc/v0.8.19/VRFV2PlusConsumerExample/VRFV2PlusConsumerExample.bin VRFV2PlusConsumerExample vrfv2plus_consumer_example
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.19/VRFCoordinatorV2_5/VRFCoordinatorV2_5.abi ../../contracts/solc/v0.8.19/VRFCoordinatorV2_5/VRFCoordinatorV2_5.bin VRFCoordinatorV2_5 vrf_coordinator_v2_5
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.19/VRFV2PlusWrapper/VRFV2PlusWrapper.abi ../../contracts/solc/v0.8.19/VRFV2PlusWrapper/VRFV2PlusWrapper.bin VRFV2PlusWrapper vrfv2plus_wrapper
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.19/VRFV2PlusWrapperConsumerExample/VRFV2PlusWrapperConsumerExample.abi ../../contracts/solc/v0.8.19/VRFV2PlusWrapperConsumerExample/VRFV2PlusWrapperConsumerExample.bin VRFV2PlusWrapperConsumerExample vrfv2plus_wrapper_consumer_example
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.19/VRFMaliciousConsumerV2Plus/VRFMaliciousConsumerV2Plus.abi ../../contracts/solc/v0.8.19/VRFMaliciousConsumerV2Plus/VRFMaliciousConsumerV2Plus.bin VRFMaliciousConsumerV2Plus vrf_malicious_consumer_v2_plus
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.19/VRFV2PlusSingleConsumerExample/VRFV2PlusSingleConsumerExample.abi ../../contracts/solc/v0.8.19/VRFV2PlusSingleConsumerExample/VRFV2PlusSingleConsumerExample.bin VRFV2PlusSingleConsumerExample vrf_v2plus_single_consumer
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.19/VRFV2PlusExternalSubOwnerExample/VRFV2PlusExternalSubOwnerExample.abi ../../contracts/solc/v0.8.19/VRFV2PlusExternalSubOwnerExample/VRFV2PlusExternalSubOwnerExample.bin VRFV2PlusExternalSubOwnerExample vrf_v2plus_sub_owner
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.19/VRFV2PlusRevertingExample/VRFV2PlusRevertingExample.abi ../../contracts/solc/v0.8.19/VRFV2PlusRevertingExample/VRFV2PlusRevertingExample.bin VRFV2PlusRevertingExample vrfv2plus_reverting_example
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.19/VRFConsumerV2PlusUpgradeableExample/VRFConsumerV2PlusUpgradeableExample.abi ../../contracts/solc/v0.8.19/VRFConsumerV2PlusUpgradeableExample/VRFConsumerV2PlusUpgradeableExample.bin VRFConsumerV2PlusUpgradeableExample vrf_consumer_v2_plus_upgradeable_example
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.19/VRFV2PlusClient/VRFV2PlusClient.abi ../../contracts/solc/v0.8.19/VRFV2PlusClient/VRFV2PlusClient.bin VRFV2PlusClient vrfv2plus_client
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.19/VRFCoordinatorV2Plus_V2Example/VRFCoordinatorV2Plus_V2Example.abi ../../contracts/solc/v0.8.19/VRFCoordinatorV2Plus_V2Example/VRFCoordinatorV2Plus_V2Example.bin VRFCoordinatorV2Plus_V2Example vrf_coordinator_v2_plus_v2_example
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.19/VRFV2PlusMaliciousMigrator/VRFV2PlusMaliciousMigrator.abi ../../contracts/solc/v0.8.19/VRFV2PlusMaliciousMigrator/VRFV2PlusMaliciousMigrator.bin VRFV2PlusMaliciousMigrator vrfv2plus_malicious_migrator
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.19/VRFV2PlusLoadTestWithMetrics/VRFV2PlusLoadTestWithMetrics.abi ../../contracts/solc/v0.8.19/VRFV2PlusLoadTestWithMetrics/VRFV2PlusLoadTestWithMetrics.bin VRFV2PlusLoadTestWithMetrics vrf_v2plus_load_test_with_metrics
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.19/VRFCoordinatorV2PlusUpgradedVersion/VRFCoordinatorV2PlusUpgradedVersion.abi ../../contracts/solc/v0.8.19/VRFCoordinatorV2PlusUpgradedVersion/VRFCoordinatorV2PlusUpgradedVersion.bin VRFCoordinatorV2PlusUpgradedVersion vrf_v2plus_upgraded_version
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.19/VRFV2PlusWrapperLoadTestConsumer/VRFV2PlusWrapperLoadTestConsumer.abi ../../contracts/solc/v0.8.19/VRFV2PlusWrapperLoadTestConsumer/VRFV2PlusWrapperLoadTestConsumer.bin VRFV2PlusWrapperLoadTestConsumer vrfv2plus_wrapper_load_test_consumer
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.19/BlockhashStore/BlockhashStore.abi ../../contracts/solc/v0.8.19/BlockhashStore/BlockhashStore.bin BlockhashStore blockhash_store
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.19/BatchBlockhashStore/BatchBlockhashStore.abi ../../contracts/solc/v0.8.19/BatchBlockhashStore/BatchBlockhashStore.bin BatchBlockhashStore batch_blockhash_store

// Aggregators
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/AggregatorV2V3Interface/AggregatorV2V3Interface.abi ../../contracts/solc/v0.8.6/AggregatorV2V3Interface/AggregatorV2V3Interface.bin AggregatorV2V3Interface aggregator_v2v3_interface
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/AggregatorV2V3Interface/AggregatorV3Interface.abi ../../contracts/solc/v0.8.6/AggregatorV2V3Interface/AggregatorV3Interface.bin AggregatorV3Interface aggregator_v3_interface
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/MockAggregatorProxy/MockAggregatorProxy.abi ../../contracts/solc/v0.8.6/MockAggregatorProxy/MockAggregatorProxy.bin MockAggregatorProxy mock_aggregator_proxy

// Log tester

// ChainReader test contract
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.19/ChainReaderTester/ChainReaderTester.abi ../../contracts/solc/v0.8.19/ChainReaderTester/ChainReaderTester.bin ChainReaderTester chain_reader_tester

// Chainlink Functions
//go:generate go generate ./functions

// Chainlink Keystone
//go:generate go generate ./keystone

// Mercury
//go:generate go generate ./llo-feeds

// Operator Forwarder
//go:generate go generate ./operatorforwarder

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
