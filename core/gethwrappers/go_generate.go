// Package gethwrappers provides tools for wrapping solidity contracts with
// golang packages, using abigen.
package gethwrappers

// Make sure solidity compiler artifacts are up to date. Only output stdout on failure.
//go:generate ./generation/compile_contracts.sh

//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.6/FluxAggregator.abi ../../contracts/solc/v0.6/FluxAggregator.bin FluxAggregator flux_aggregator_wrapper
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.6/VRFTestHelper.abi ../../contracts/solc/v0.6/VRFTestHelper.bin VRFTestHelper solidity_vrf_verifier_wrapper
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.6/VRFCoordinator.abi ../../contracts/solc/v0.6/VRFCoordinator.bin VRFCoordinator solidity_vrf_coordinator_interface
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.6/VRFConsumer.abi ../../contracts/solc/v0.6/VRFConsumer.bin VRFConsumer solidity_vrf_consumer_interface
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.6/VRFRequestIDBaseTestHelper.abi ../../contracts/solc/v0.6/VRFRequestIDBaseTestHelper.bin VRFRequestIDBaseTestHelper solidity_vrf_request_id
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.6/Flags.abi ../../contracts/solc/v0.6/Flags.bin Flags flags_wrapper
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.6/Oracle.abi ../../contracts/solc/v0.6/Oracle.bin Oracle oracle_wrapper
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.6/BlockhashStore.abi ../../contracts/solc/v0.6/BlockhashStore.bin BlockhashStore blockhash_store
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.7/Consumer.abi ../../contracts/solc/v0.7/Consumer.bin Consumer consumer_wrapper
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.7/MultiWordConsumer.abi ../../contracts/solc/v0.7/MultiWordConsumer.bin MultiWordConsumer multiwordconsumer_wrapper
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.7/Operator.abi ../../contracts/solc/v0.7/Operator.bin Operator operator_wrapper
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.7/OperatorFactory.abi ../../contracts/solc/v0.7/OperatorFactory.bin OperatorFactory operator_factory
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.7/AuthorizedForwarder.abi ../../contracts/solc/v0.7/AuthorizedForwarder.bin AuthorizedForwarder authorized_forwarder
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.7/AuthorizedReceiver.abi ../../contracts/solc/v0.7/AuthorizedReceiver.bin AuthorizedReceiver authorized_receiver
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/BatchBlockhashStore.abi ../../contracts/solc/v0.8.6/BatchBlockhashStore.bin BatchBlockhashStore batch_blockhash_store
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/BatchVRFCoordinatorV2.abi ../../contracts/solc/v0.8.6/BatchVRFCoordinatorV2.bin BatchVRFCoordinatorV2 batch_vrf_coordinator_v2
//go:generate go run ./generation/generate/wrap.go OffchainAggregator/OffchainAggregator.abi - OffchainAggregator offchain_aggregator_wrapper

//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.7/KeeperRegistry1_1.abi ../../contracts/solc/v0.7/KeeperRegistry1_1.bin KeeperRegistry keeper_registry_wrapper1_1
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.7/UpkeepPerformCounterRestrictive.abi ../../contracts/solc/v0.7/UpkeepPerformCounterRestrictive.bin UpkeepPerformCounterRestrictive upkeep_perform_counter_restrictive_wrapper
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.7/UpkeepCounter.abi ../../contracts/solc/v0.7/UpkeepCounter.bin UpkeepCounter upkeep_counter_wrapper
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/CronUpkeepFactory.abi - CronUpkeepFactory cron_upkeep_factory_wrapper
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/CronUpkeep.abi - CronUpkeep cron_upkeep_wrapper
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/KeeperRegistry1_2.abi ../../contracts/solc/v0.8.6/KeeperRegistry1_2.bin KeeperRegistry keeper_registry_wrapper1_2
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/TypeAndVersionInterface.abi ../../contracts/solc/v0.8.6/TypeAndVersionInterface.bin TypeAndVersionInterface type_and_version_interface_wrapper
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/KeeperRegistryCheckUpkeepGasUsageWrapper1_2.abi ../../contracts/solc/v0.8.6/KeeperRegistryCheckUpkeepGasUsageWrapper1_2.bin KeeperRegistryCheckUpkeepGasUsageWrapper gas_wrapper
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/KeeperRegistry1_3.abi ../../contracts/solc/v0.8.6/KeeperRegistry1_3.bin KeeperRegistry keeper_registry_wrapper1_3
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/KeeperRegistryLogic1_3.abi ../../contracts/solc/v0.8.6/KeeperRegistryLogic1_3.bin KeeperRegistryLogic keeper_registry_logic1_3
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/KeeperRegistry2_0.abi ../../contracts/solc/v0.8.6/KeeperRegistry2_0.bin KeeperRegistry keeper_registry_wrapper2_0
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/KeeperRegistryLogic2_0.abi ../../contracts/solc/v0.8.6/KeeperRegistryLogic2_0.bin KeeperRegistryLogic keeper_registry_logic2_0

// v0.8.6 VRFConsumer
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/VRFConsumer.abi ../../contracts/solc/v0.8.6/VRFConsumer.bin VRFConsumer solidity_vrf_consumer_interface_v08
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/VRFRequestIDBaseTestHelper.abi ../../contracts/solc/v0.8.6/VRFRequestIDBaseTestHelper.bin VRFRequestIDBaseTestHelper solidity_vrf_request_id_v08
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/VRFOwnerlessConsumerExample.abi ../../contracts/solc/v0.8.6/VRFOwnerlessConsumerExample.bin VRFOwnerlessConsumerExample vrf_ownerless_consumer_example
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/VRFLoadTestOwnerlessConsumer.abi ../../contracts/solc/v0.8.6/VRFLoadTestOwnerlessConsumer.bin VRFLoadTestOwnerlessConsumer vrf_load_test_ownerless_consumer
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/VRFLoadTestExternalSubOwner.abi ../../contracts/solc/v0.8.6/VRFLoadTestExternalSubOwner.bin VRFLoadTestExternalSubOwner vrf_load_test_external_sub_owner

//go:generate go run ./generation/generate_link/wrap_link.go

// VRF V2
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
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/NoCancelVRFCoordinatorV2.abi ../../contracts/solc/v0.8.6/NoCancelVRFCoordinatorV2.bin NoCancelVRFCoordinatorV2 nocancel_vrf_coordinator_v2

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

// Log tester
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/LogEmitter.abi ../../contracts/solc/v0.8.6/LogEmitter.bin LogEmitter log_emitter

// Chainlink Functions (OCR2DR)
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/Functions.abi ../../contracts/solc/v0.8.6/Functions.bin OCR2DR ocr2dr
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/FunctionsClient.abi ../../contracts/solc/v0.8.6/FunctionsClient.bin OCR2DRClient ocr2dr_client
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/FunctionsClientExample.abi ../../contracts/solc/v0.8.6/FunctionsClientExample.bin OCR2DRClientExample ocr2dr_client_example
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/FunctionsOracleWithInit.abi ../../contracts/solc/v0.8.6/FunctionsOracleWithInit.bin OCR2DROracle ocr2dr_oracle
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/FunctionsBillingRegistryWithInit.abi ../../contracts/solc/v0.8.6/FunctionsBillingRegistryWithInit.bin OCR2DRRegistry ocr2dr_registry

// To run these commands, you must either install docker, or the correct version
// of abigen. The latter can be installed with these commands, at least on linux:
//
//   git clone https://github.com/ethereum/go-ethereum
//   cd go-ethereum/cmd/abigen
//   git checkout v<version-needed>
//   go install
//
// Here, <version-needed> is the version of go-ethereum specified in chainlink's
// go.mod. This will install abigen in "$GOPATH/bin", which you should add to
// your $PATH.
//
// To reduce explicit dependencies, and in case the system does not have the
// correct version of abigen installed , the above commands spin up docker
// containers. In my hands, total running time including compilation is about
// 13s. If you're modifying solidity code and testing against go code a lot, it
// might be worthwhile to generate the the wrappers using a static container
// with abigen and solc, which will complete much faster. E.g.
//
//   abigen -sol ../../contracts/src/v0.6/VRFAll.sol -pkg vrf -out solidity_interfaces.go
//
// where VRFAll.sol simply contains `import "contract_path";` instructions for
// all the contracts you wish to target. This runs in about 0.25 seconds in my
// hands.
//
// If you're on linux, you can copy the correct version of solc out of the
// appropriate docker container. At least, the following works on ubuntu:
//
//   $ docker run --name solc ethereum/solc:0.6.2
//   $ sudo docker cp solc:/usr/bin/solc /usr/bin
//   $ docker rm solc
//
// If you need to point abigen at your solc executable, you can specify the path
// with the abigen --solc <path-to-executable> option.
