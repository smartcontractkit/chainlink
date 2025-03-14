package gethwrappers

//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/automation/KeeperRegistrar1_2/KeeperRegistrar1_2.sol/KeeperRegistrar.abi.json ../../contracts/solc/automation/KeeperRegistrar1_2/KeeperRegistrar1_2.sol/KeeperRegistrar.bin KeeperRegistrar keeper_registrar_wrapper1_2
//go:generate go run ./generation/wrap_automation.go KeeperRegistrar1_2Mock KeeperRegistrar keeper_registrar_wrapper1_2_mock
//go:generate go run ./generation/wrap_automation.go KeeperRegistry1_2 KeeperRegistry keeper_registry_wrapper1_2
//go:generate go run ./generation/wrap_automation.go KeeperRegistryCheckUpkeepGasUsageWrapper1_2Mock KeeperRegistryCheckUpkeepGasUsageWrapper gas_wrapper_mock
//go:generate go run ./generation/wrap_automation.go KeeperRegistry1_3 KeeperRegistry keeper_registry_wrapper1_3
//go:generate go run ./generation/wrap_automation.go KeeperRegistryLogic1_3 KeeperRegistryLogic keeper_registry_logic1_3
//go:generate go run ./generation/wrap_automation.go KeeperRegistrar2_0 KeeperRegistrar keeper_registrar_wrapper2_0
//go:generate go run ./generation/wrap_automation.go KeeperRegistry2_0 KeeperRegistry keeper_registry_wrapper2_0
//go:generate go run ./generation/wrap_automation.go KeeperRegistryLogic2_0 KeeperRegistryLogic keeper_registry_logic2_0
//go:generate go run ./generation/wrap_automation.go UpkeepTranscoder upkeep_transcoder
//go:generate go run ./generation/wrap_automation.go VerifiableLoadUpkeep verifiable_load_upkeep_wrapper
//go:generate go run ./generation/wrap_automation.go VerifiableLoadStreamsLookupUpkeep verifiable_load_streams_lookup_upkeep_wrapper
//go:generate go run ./generation/wrap_automation.go StreamsLookupUpkeep streams_lookup_upkeep_wrapper
//go:generate go run ./generation/wrap_automation.go StreamsLookupCompatibleInterface streams_lookup_compatible_interface
//go:generate go run ./generation/wrap_automation.go AutomationConsumerBenchmark automation_consumer_benchmark
//go:generate go run ./generation/wrap_automation.go AutomationRegistrar2_1 AutomationRegistrar automation_registrar_wrapper2_1
//go:generate go run ./generation/wrap_automation.go KeeperRegistry2_1 KeeperRegistry keeper_registry_wrapper_2_1
//go:generate go run ./generation/wrap_automation.go KeeperRegistryLogicA2_1 KeeperRegistryLogicA keeper_registry_logic_a_wrapper_2_1
//go:generate go run ./generation/wrap_automation.go KeeperRegistryLogicB2_1 KeeperRegistryLogicB keeper_registry_logic_b_wrapper_2_1
//go:generate go run ./generation/wrap_automation.go IKeeperRegistryMaster i_keeper_registry_master_wrapper_2_1
//go:generate go run ./generation/wrap_automation.go AutomationRegistry2_2 AutomationRegistry automation_registry_wrapper_2_2
//go:generate go run ./generation/wrap_automation.go AutomationRegistryLogicA2_2 AutomationRegistryLogicA automation_registry_logic_a_wrapper_2_2
//go:generate go run ./generation/wrap_automation.go AutomationRegistryLogicB2_2 AutomationRegistryLogicB automation_registry_logic_b_wrapper_2_2
//go:generate go run ./generation/wrap_automation.go IAutomationRegistryMaster i_automation_registry_master_wrapper_2_2
//go:generate go run ./generation/wrap_automation.go AutomationCompatibleUtils automation_compatible_utils
//go:generate go run ./generation/wrap_automation.go AutomationRegistrar2_3 AutomationRegistrar automation_registrar_wrapper2_3
//go:generate go run ./generation/wrap_automation.go AutomationRegistry2_3 AutomationRegistry automation_registry_wrapper_2_3
//go:generate go run ./generation/wrap_automation.go AutomationRegistryLogicA2_3 AutomationRegistryLogicA automation_registry_logic_a_wrapper_2_3
//go:generate go run ./generation/wrap_automation.go AutomationRegistryLogicB2_3 AutomationRegistryLogicB automation_registry_logic_b_wrapper_2_3
//go:generate go run ./generation/wrap_automation.go AutomationRegistryLogicC2_3 AutomationRegistryLogicC automation_registry_logic_c_wrapper_2_3
//go:generate go run ./generation/wrap_automation.go IAutomationRegistryMaster2_3 i_automation_registry_master_wrapper_2_3
//go:generate go run ./generation/wrap_automation.go ArbitrumModule arbitrum_module
//go:generate go run ./generation/wrap_automation.go ChainModuleBase chain_module_base
//go:generate go run ./generation/wrap_automation.go ScrollModule scroll_module
//go:generate go run ./generation/wrap_automation.go IChainModule i_chain_module
//go:generate go run ./generation/wrap_automation.go IAutomationV21PlusCommon i_automation_v21_plus_common
//go:generate go run ./generation/wrap_automation.go MockETHUSDAggregator mock_ethusd_aggregator_wrapper

//go:generate go run ./generation/wrap_automation.go ILogAutomation i_log_automation
//go:generate go run ./generation/wrap_automation.go AutomationForwarderLogic automation_forwarder_logic
//go:generate go run ./generation/wrap_automation.go LogUpkeepCounter log_upkeep_counter_wrapper
//go:generate go run ./generation/wrap_automation.go SimpleLogUpkeepCounter simple_log_upkeep_counter_wrapper
//go:generate go run ./generation/wrap_automation.go LogTriggeredStreamsLookup log_triggered_streams_lookup_wrapper
//go:generate go run ./generation/wrap_automation.go DummyProtocol dummy_protocol_wrapper

//go:generate go run ./generation/wrap_automation.go KeeperConsumerPerformance keeper_consumer_performance_wrapper
//go:generate go run ./generation/wrap_automation.go PerformDataChecker perform_data_checker_wrapper
//go:generate go run ./generation/wrap_automation.go UpkeepCounter upkeep_counter_wrapper
//go:generate go run ./generation/wrap_automation.go  UpkeepPerformCounterRestrictive upkeep_perform_counter_restrictive_wrapper
