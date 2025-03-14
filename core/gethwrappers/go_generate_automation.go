package gethwrappers

//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/automation/KeeperRegistrar1_2/KeeperRegistrar1_2.sol/KeeperRegistrar.abi.json ../../contracts/solc/automation/KeeperRegistrar1_2/KeeperRegistrar1_2.sol/KeeperRegistrar.bin KeeperRegistrar keeper_registrar_wrapper1_2
//go:generate go run ./generation/generate_automation/wrap.go KeeperRegistrar1_2Mock KeeperRegistrarMock keeper_registrar_wrapper1_2_mock
//go:generate go run ./generation/generate_automation/wrap.go KeeperRegistry1_2 KeeperRegistry keeper_registry_wrapper1_2
//go:generate go run ./generation/generate_automation/wrap.go KeeperRegistryCheckUpkeepGasUsageWrapper1_2Mock KeeperRegistryCheckUpkeepGasUsageWrapperMock gas_wrapper_mock
//go:generate go run ./generation/generate_automation/wrap.go KeeperRegistry1_3 KeeperRegistry keeper_registry_wrapper1_3
//go:generate go run ./generation/generate_automation/wrap.go KeeperRegistryLogic1_3 KeeperRegistryLogic keeper_registry_logic1_3
//go:generate go run ./generation/generate_automation/wrap.go KeeperRegistrar2_0 KeeperRegistrar keeper_registrar_wrapper2_0
//go:generate go run ./generation/generate_automation/wrap.go KeeperRegistry2_0 KeeperRegistry keeper_registry_wrapper2_0
//go:generate go run ./generation/generate_automation/wrap.go KeeperRegistryLogic2_0 KeeperRegistryLogic keeper_registry_logic2_0
//go:generate go run ./generation/wrap.go automation UpkeepTranscoder upkeep_transcoder
//go:generate go run ./generation/wrap.go automation VerifiableLoadUpkeep verifiable_load_upkeep_wrapper
//go:generate go run ./generation/wrap.go automation VerifiableLoadStreamsLookupUpkeep verifiable_load_streams_lookup_upkeep_wrapper
//go:generate go run ./generation/wrap.go automation StreamsLookupUpkeep streams_lookup_upkeep_wrapper
//go:generate go run ./generation/wrap.go automation StreamsLookupCompatibleInterface streams_lookup_compatible_interface
//go:generate go run ./generation/wrap.go automation AutomationConsumerBenchmark automation_consumer_benchmark
//go:generate go run ./generation/generate_automation/wrap.go AutomationRegistrar2_1 AutomationRegistrar automation_registrar_wrapper2_1
//go:generate go run ./generation/generate_automation/wrap.go KeeperRegistry2_1 KeeperRegistry keeper_registry_wrapper_2_1
//go:generate go run ./generation/generate_automation/wrap.go KeeperRegistryLogicA2_1 KeeperRegistryLogicA keeper_registry_logic_a_wrapper_2_1
//go:generate go run ./generation/generate_automation/wrap.go KeeperRegistryLogicB2_1 KeeperRegistryLogicB keeper_registry_logic_b_wrapper_2_1
//go:generate go run ./generation/wrap.go automation IKeeperRegistryMaster i_keeper_registry_master_wrapper_2_1
//go:generate go run ./generation/generate_automation/wrap.go AutomationRegistry2_2 AutomationRegistry automation_registry_wrapper_2_2
//go:generate go run ./generation/generate_automation/wrap.go AutomationRegistryLogicA2_2 AutomationRegistryLogicA automation_registry_logic_a_wrapper_2_2
//go:generate go run ./generation/generate_automation/wrap.go AutomationRegistryLogicB2_2 AutomationRegistryLogicB automation_registry_logic_b_wrapper_2_2
//go:generate go run ./generation/wrap.go automation IAutomationRegistryMaster i_automation_registry_master_wrapper_2_2
//go:generate go run ./generation/wrap.go automation AutomationCompatibleUtils automation_compatible_utils
//go:generate go run ./generation/generate_automation/wrap.go AutomationRegistrar2_3 AutomationRegistrar automation_registrar_wrapper2_3
//go:generate go run ./generation/generate_automation/wrap.go AutomationRegistry2_3 AutomationRegistry automation_registry_wrapper_2_3
//go:generate go run ./generation/generate_automation/wrap.go AutomationRegistryLogicA2_3 AutomationRegistryLogicA automation_registry_logic_a_wrapper_2_3
//go:generate go run ./generation/generate_automation/wrap.go AutomationRegistryLogicB2_3 AutomationRegistryLogicB automation_registry_logic_b_wrapper_2_3
//go:generate go run ./generation/generate_automation/wrap.go AutomationRegistryLogicC2_3 AutomationRegistryLogicC automation_registry_logic_c_wrapper_2_3
//go:generate go run ./generation/wrap.go automation IAutomationRegistryMaster2_3 i_automation_registry_master_wrapper_2_3
//go:generate go run ./generation/wrap.go automation ArbitrumModule arbitrum_module
//go:generate go run ./generation/wrap.go automation ChainModuleBase chain_module_base
//go:generate go run ./generation/wrap.go automation ScrollModule scroll_module
//go:generate go run ./generation/wrap.go automation IChainModule i_chain_module
//go:generate go run ./generation/wrap.go automation IAutomationV21PlusCommon i_automation_v21_plus_common
//go:generate go run ./generation/wrap.go automation MockETHUSDAggregator mock_ethusd_aggregator_wrapper

//go:generate go run ./generation/wrap.go automation ILogAutomation i_log_automation
//go:generate go run ./generation/wrap.go automation AutomationForwarderLogic automation_forwarder_logic
//go:generate go run ./generation/wrap.go automation LogUpkeepCounter log_upkeep_counter_wrapper
//go:generate go run ./generation/wrap.go automation SimpleLogUpkeepCounter simple_log_upkeep_counter_wrapper
//go:generate go run ./generation/wrap.go automation LogTriggeredStreamsLookup log_triggered_streams_lookup_wrapper
//go:generate go run ./generation/wrap.go automation DummyProtocol dummy_protocol_wrapper

//go:generate go run ./generation/wrap.go automation KeeperConsumerPerformance keeper_consumer_performance_wrapper
//go:generate go run ./generation/wrap.go automation PerformDataChecker perform_data_checker_wrapper
//go:generate go run ./generation/wrap.go automation UpkeepCounter upkeep_counter_wrapper
//go:generate go run ./generation/wrap.go automation UpkeepPerformCounterRestrictive upkeep_perform_counter_restrictive_wrapper
