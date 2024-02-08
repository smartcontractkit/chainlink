// Package gethwrappers_ccip provides tools for wrapping solidity contracts with
// golang packages, using abigen.
package ccip

//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.19/CommitStore/CommitStore.abi ../../../contracts/solc/v0.8.19/CommitStore/CommitStore.bin CommitStore commit_store
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.19/CommitStoreHelper/CommitStoreHelper.abi ../../../contracts/solc/v0.8.19/CommitStoreHelper/CommitStoreHelper.bin CommitStoreHelper commit_store_helper
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.19/BurnMintTokenPool/BurnMintTokenPool.abi ../../../contracts/solc/v0.8.19/BurnMintTokenPool/BurnMintTokenPool.bin BurnMintTokenPool burn_mint_token_pool
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.19/BurnFromMintTokenPool/BurnFromMintTokenPool.abi ../../../contracts/solc/v0.8.19/BurnFromMintTokenPool/BurnFromMintTokenPool.bin BurnFromMintTokenPool burn_from_mint_token_pool
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.19/BurnWithFromMintTokenPool/BurnWithFromMintTokenPool.abi ../../../contracts/solc/v0.8.19/BurnWithFromMintTokenPool/BurnWithFromMintTokenPool.bin BurnWithFromMintTokenPool burn_with_from_mint_token_pool
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.19/LockReleaseTokenPool/LockReleaseTokenPool.abi ../../../contracts/solc/v0.8.19/LockReleaseTokenPool/LockReleaseTokenPool.bin LockReleaseTokenPool lock_release_token_pool
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.19/CustomTokenPool/CustomTokenPool.abi ../../../contracts/solc/v0.8.19/CustomTokenPool/CustomTokenPool.bin CustomTokenPool custom_token_pool
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.19/TokenPool/TokenPool.abi ../../../contracts/solc/v0.8.19/TokenPool/TokenPool.bin TokenPool token_pool
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.19/ARM/ARM.abi ../../../contracts/solc/v0.8.19/ARM/ARM.bin ARMContract arm_contract
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.19/ARMProxy/ARMProxy.abi ../../../contracts/solc/v0.8.19/ARMProxy/ARMProxy.bin ARMProxyContract arm_proxy_contract
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.19/MockARM/MockARM.abi ../../../contracts/solc/v0.8.19/MockARM/MockARM.bin MockARMContract mock_arm_contract
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.19/MockE2EUSDCTokenMessenger/MockE2EUSDCTokenMessenger.abi ../../../contracts/solc/v0.8.19/MockE2EUSDCTokenMessenger/MockE2EUSDCTokenMessenger.bin MockE2EUSDCTokenMessenger mock_usdc_token_messenger
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.19/MockE2EUSDCTransmitter/MockE2EUSDCTransmitter.abi ../../../contracts/solc/v0.8.19/MockE2EUSDCTransmitter/MockE2EUSDCTransmitter.bin MockE2EUSDCTransmitter mock_usdc_token_transmitter

//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.19/MockV3Aggregator/MockV3Aggregator.abi ../../../contracts/solc/v0.8.19/MockV3Aggregator/MockV3Aggregator.bin MockV3Aggregator mock_v3_aggregator_contract

//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.19/EVM2EVMOnRamp/EVM2EVMOnRamp.abi ../../../contracts/solc/v0.8.19/EVM2EVMOnRamp/EVM2EVMOnRamp.bin EVM2EVMOnRamp evm_2_evm_onramp
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.19/EVM2EVMOffRamp/EVM2EVMOffRamp.abi ../../../contracts/solc/v0.8.19/EVM2EVMOffRamp/EVM2EVMOffRamp.bin EVM2EVMOffRamp evm_2_evm_offramp
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.19/Router/Router.abi ../../../contracts/solc/v0.8.19/Router/Router.bin Router router
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.19/PriceRegistry/PriceRegistry.abi ../../../contracts/solc/v0.8.19/PriceRegistry/PriceRegistry.bin PriceRegistry price_registry

//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.19/MaybeRevertMessageReceiver/MaybeRevertMessageReceiver.abi ../../../contracts/solc/v0.8.19/MaybeRevertMessageReceiver/MaybeRevertMessageReceiver.bin MaybeRevertMessageReceiver maybe_revert_message_receiver
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.19/PingPongDemo/PingPongDemo.abi ../../../contracts/solc/v0.8.19/PingPongDemo/PingPongDemo.bin PingPongDemo ping_pong_demo
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.19/SelfFundedPingPong/SelfFundedPingPong.abi ../../../contracts/solc/v0.8.19/SelfFundedPingPong/SelfFundedPingPong.bin SelfFundedPingPong self_funded_ping_pong
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.19/WETH9/WETH9.abi ../../../contracts/solc/v0.8.19/WETH9/WETH9.bin WETH9 weth9

// Customer contracts
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.19/USDCTokenPool/USDCTokenPool.abi ../../../contracts/solc/v0.8.19/USDCTokenPool/USDCTokenPool.bin USDCTokenPool usdc_token_pool

// Generate mocks for our contracts
//go:generate mockery --quiet --dir ./generated/evm_2_evm_onramp/ --name EVM2EVMOnRampInterface --output ./mocks/ --outpkg mock_contracts --case=underscore
//go:generate mockery --quiet --dir ./generated/evm_2_evm_offramp/ --name EVM2EVMOffRampInterface --output ./mocks/ --outpkg mock_contracts --case=underscore
//go:generate mockery --quiet --dir ./generated/evm_2_evm_offramp_1_2_0/ --name EVM2EVMOffRampInterface --output ./mocks/v1_2_0/ --outpkg mock_contracts --case=underscore
//go:generate mockery --quiet --dir ./generated/evm_2_evm_offramp_1_0_0/ --name EVM2EVMOffRampInterface --output ./mocks/v1_0_0/ --outpkg mock_contracts --case=underscore
//go:generate mockery --quiet --dir ./generated/commit_store/ --name CommitStoreInterface --output ./mocks/ --outpkg mock_contracts --case=underscore
//go:generate mockery --quiet --dir ./generated/price_registry/ --name PriceRegistryInterface --output ./mocks/ --outpkg mock_contracts --case=underscore
//go:generate mockery --quiet --dir ../generated/link_token_interface/ --name LinkTokenInterface --output ./mocks/ --outpkg mock_contracts --case=underscore

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
