package gethwrappers

// v0.8.6 VRFConsumer
//go:generate go run ./generation/wrap.go vrf VRFCoordinatorMock vrf_coordinator_mock
//go:generate go run ./generation/wrap.go vrf VRFCoordinatorMock vrf_coordinator_mock
//go:generate go run ./generation/wrap.go vrf VRFConsumer solidity_vrf_consumer_interface_v08
//go:generate go run ./generation/wrap.go vrf VRFRequestIDBaseTestHelper solidity_vrf_request_id_v08
//go:generate go run ./generation/wrap.go vrf VRFOwnerlessConsumerExample vrf_ownerless_consumer_example
//go:generate go run ./generation/wrap.go vrf VRFLoadTestOwnerlessConsumer vrf_load_test_ownerless_consumer
//go:generate go run ./generation/wrap.go vrf VRFLoadTestExternalSubOwner vrf_load_test_external_sub_owner
//go:generate go run ./generation/wrap.go vrf VRFV2LoadTestWithMetrics vrf_load_test_with_metrics
//go:generate go run ./generation/wrap.go vrf VRFV2OwnerTestConsumer vrf_owner_test_consumer
//go:generate go run ./generation/wrap.go vrf VRFv2Consumer vrf_v2_consumer_wrapper
//go:generate go run ./generation/wrap.go vrf Counter counter

// VRF V2
//go:generate go run ./generation/wrap.go vrf BatchVRFCoordinatorV2 batch_vrf_coordinator_v2
//go:generate go run ./generation/wrap.go vrf VRFOwner vrf_owner
//go:generate go run ./generation/wrap.go vrf VRFCoordinatorV2 vrf_coordinator_v2
//go:generate go run ./generation/wrap.go vrf VRFConsumerV2 vrf_consumer_v2
//go:generate go run ./generation/wrap.go vrf VRFMaliciousConsumerV2 vrf_malicious_consumer_v2
//go:generate go run ./generation/wrap.go vrf VRFTestHelper solidity_vrf_v08_verifier_wrapper
//go:generate go run ./generation/wrap.go vrf VRFSingleConsumerExample vrf_single_consumer_example
//go:generate go run ./generation/wrap.go vrf VRFExternalSubOwnerExample vrf_external_sub_owner_example
//go:generate go run ./generation/wrap.go vrf VRFV2RevertingExample vrfv2_reverting_example
//go:generate go run ./generation/wrap.go vrf VRFConsumerV2UpgradeableExample vrf_consumer_v2_upgradeable_example
//go:generate go run ./generation/wrap.go vrf VRFV2TransparentUpgradeableProxy vrfv2_transparent_upgradeable_proxy
//go:generate go run ./generation/wrap.go vrf VRFV2ProxyAdmin vrfv2_proxy_admin
//go:generate go run ./generation/wrap.go vrf ChainSpecificUtilHelper chain_specific_util_helper
//go:generate go run ./generation/wrap.go vrf VRFCoordinatorTestV2 vrf_coordinator_test_v2
//go:generate go run ./generation/wrap.go vrf VRFMockETHLINKAggregator vrf_mock_ethlink_aggregator

// VRF V2 Wrapper
//go:generate go run ./generation/wrap.go vrf VRFV2Wrapper vrfv2_wrapper
//go:generate go run ./generation/wrap.go vrf VRFV2WrapperInterface vrfv2_wrapper_interface
//go:generate go run ./generation/wrap.go vrf VRFV2WrapperConsumerExample vrfv2_wrapper_consumer_example
//go:generate go run ./generation/wrap.go vrf VRFV2WrapperLoadTestConsumer vrfv2_wrapper_load_test_consumer

// Keepers X VRF v2
//go:generate go run ./generation/wrap.go vrf KeepersVRFConsumer keepers_vrf_consumer

// VRF V2Plus
//go:generate go run ./generation/wrap.go vrf IVRFCoordinatorV2PlusInternal vrf_coordinator_v2plus_interface
//go:generate go run ./generation/wrap.go vrf BatchVRFCoordinatorV2Plus batch_vrf_coordinator_v2plus
//go:generate go run ./generation/wrap.go vrf TrustedBlockhashStore trusted_blockhash_store
//go:generate go run ./generation/wrap.go vrf VRFV2PlusConsumerExample vrfv2plus_consumer_example
//go:generate go run ./generation/wrap.go vrf VRFCoordinatorV2_5 vrf_coordinator_v2_5
//go:generate go run ./generation/wrap.go vrf VRFV2PlusWrapper vrfv2plus_wrapper
//go:generate go run ./generation/wrap.go vrf VRFV2PlusWrapperConsumerExample vrfv2plus_wrapper_consumer_example
//go:generate go run ./generation/wrap.go vrf VRFMaliciousConsumerV2Plus vrf_malicious_consumer_v2_plus
//go:generate go run ./generation/wrap.go vrf VRFV2PlusSingleConsumerExample vrf_v2plus_single_consumer
//go:generate go run ./generation/wrap.go vrf VRFV2PlusExternalSubOwnerExample vrf_v2plus_sub_owner
//go:generate go run ./generation/wrap.go vrf VRFV2PlusRevertingExample vrfv2plus_reverting_example
//go:generate go run ./generation/wrap.go vrf VRFConsumerV2PlusUpgradeableExample vrf_consumer_v2_plus_upgradeable_example
//go:generate go run ./generation/wrap.go vrf VRFV2PlusClient vrfv2plus_client
//go:generate go run ./generation/wrap.go vrf VRFCoordinatorV2Plus_V2Example vrf_coordinator_v2_plus_v2_example
//go:generate go run ./generation/wrap.go vrf VRFV2PlusMaliciousMigrator vrfv2plus_malicious_migrator
//go:generate go run ./generation/wrap.go vrf VRFV2PlusLoadTestWithMetrics vrf_v2plus_load_test_with_metrics
//go:generate go run ./generation/wrap.go vrf VRFCoordinatorV2PlusUpgradedVersion vrf_v2plus_upgraded_version
//go:generate go run ./generation/wrap.go vrf VRFV2PlusWrapperLoadTestConsumer vrfv2plus_wrapper_load_test_consumer
//go:generate go run ./generation/wrap.go vrf BlockhashStore blockhash_store
//go:generate go run ./generation/wrap.go vrf BatchBlockhashStore batch_blockhash_store
//go:generate go run ./generation/wrap.go vrf VRFCoordinatorV2_5_Arbitrum vrf_coordinator_v2_5_arbitrum
//go:generate go run ./generation/wrap.go vrf VRFCoordinatorV2_5_Optimism vrf_coordinator_v2_5_optimism
//go:generate go run ./generation/wrap.go vrf VRFV2PlusWrapper_Arbitrum vrfv2plus_wrapper_arbitrum
//go:generate go run ./generation/wrap.go vrf VRFV2PlusWrapper_Optimism vrfv2plus_wrapper_optimism
//go:generate go run ./generation/wrap.go vrf VRFCoordinatorTestV2_5 vrf_coordinator_test_v2_5
