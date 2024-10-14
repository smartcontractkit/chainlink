// Package gethwrappers_ccip provides tools for wrapping solidity contracts with
// golang packages, using abigen.
package ccip

//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.24/CommitStore/CommitStore.abi ../../../contracts/solc/v0.8.24/CommitStore/CommitStore.bin CommitStore commit_store
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.24/CommitStoreHelper/CommitStoreHelper.abi ../../../contracts/solc/v0.8.24/CommitStoreHelper/CommitStoreHelper.bin CommitStoreHelper commit_store_helper
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.24/RMN/RMN.abi ../../../contracts/solc/v0.8.24/RMN/RMN.bin RMNContract rmn_contract
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.24/ARMProxy/ARMProxy.abi ../../../contracts/solc/v0.8.24/ARMProxy/ARMProxy.bin RMNProxyContract rmn_proxy_contract
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.24/TokenAdminRegistry/TokenAdminRegistry.abi ../../../contracts/solc/v0.8.24/TokenAdminRegistry/TokenAdminRegistry.bin TokenAdminRegistry token_admin_registry
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.24/RegistryModuleOwnerCustom/RegistryModuleOwnerCustom.abi ../../../contracts/solc/v0.8.24/RegistryModuleOwnerCustom/RegistryModuleOwnerCustom.bin RegistryModuleOwnerCustom registry_module_owner_custom
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.24/EVM2EVMOnRamp/EVM2EVMOnRamp.abi ../../../contracts/solc/v0.8.24/EVM2EVMOnRamp/EVM2EVMOnRamp.bin EVM2EVMOnRamp evm_2_evm_onramp
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.24/OnRamp/OnRamp.abi ../../../contracts/solc/v0.8.24/OnRamp/OnRamp.bin OnRamp onramp
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.24/EVM2EVMOffRamp/EVM2EVMOffRamp.abi ../../../contracts/solc/v0.8.24/EVM2EVMOffRamp/EVM2EVMOffRamp.bin EVM2EVMOffRamp evm_2_evm_offramp
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.24/OffRamp/OffRamp.abi ../../../contracts/solc/v0.8.24/OffRamp/OffRamp.bin OffRamp offramp
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.24/RMNRemote/RMNRemote.abi ../../../contracts/solc/v0.8.24/RMNRemote/RMNRemote.bin RMNRemote rmn_remote
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.24/RMNHome/RMNHome.abi ../../../contracts/solc/v0.8.24/RMNHome/RMNHome.bin RMNHome rmn_home
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.24/MultiAggregateRateLimiter/MultiAggregateRateLimiter.abi ../../../contracts/solc/v0.8.24/MultiAggregateRateLimiter/MultiAggregateRateLimiter.bin MultiAggregateRateLimiter multi_aggregate_rate_limiter
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.24/Router/Router.abi ../../../contracts/solc/v0.8.24/Router/Router.bin Router router
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.24/FeeQuoter/FeeQuoter.abi ../../../contracts/solc/v0.8.24/FeeQuoter/FeeQuoter.bin FeeQuoter fee_quoter
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.24/CCIPHome/CCIPHome.abi ../../../contracts/solc/v0.8.24/CCIPHome/CCIPHome.bin CCIPHome ccip_home
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.24/NonceManager/NonceManager.abi ../../../contracts/solc/v0.8.24/NonceManager/NonceManager.bin NonceManager nonce_manager

// Pools
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.24/BurnMintTokenPool/BurnMintTokenPool.abi ../../../contracts/solc/v0.8.24/BurnMintTokenPool/BurnMintTokenPool.bin BurnMintTokenPool burn_mint_token_pool
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.24/BurnFromMintTokenPool/BurnFromMintTokenPool.abi ../../../contracts/solc/v0.8.24/BurnFromMintTokenPool/BurnFromMintTokenPool.bin BurnFromMintTokenPool burn_from_mint_token_pool
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.24/BurnWithFromMintTokenPool/BurnWithFromMintTokenPool.abi ../../../contracts/solc/v0.8.24/BurnWithFromMintTokenPool/BurnWithFromMintTokenPool.bin BurnWithFromMintTokenPool burn_with_from_mint_token_pool
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.24/BurnMintTokenPoolAndProxy/BurnMintTokenPoolAndProxy.abi ../../../contracts/solc/v0.8.24/BurnMintTokenPoolAndProxy/BurnMintTokenPoolAndProxy.bin BurnMintTokenPoolAndProxy burn_mint_token_pool_and_proxy
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.24/BurnWithFromMintTokenPoolAndProxy/BurnWithFromMintTokenPoolAndProxy.abi ../../../contracts/solc/v0.8.24/BurnWithFromMintTokenPoolAndProxy/BurnWithFromMintTokenPoolAndProxy.bin BurnWithFromMintTokenPoolAndProxy burn_with_from_mint_token_pool_and_proxy
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.24/LockReleaseTokenPool/LockReleaseTokenPool.abi ../../../contracts/solc/v0.8.24/LockReleaseTokenPool/LockReleaseTokenPool.bin LockReleaseTokenPool lock_release_token_pool
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.24/LockReleaseTokenPoolAndProxy/LockReleaseTokenPoolAndProxy.abi ../../../contracts/solc/v0.8.24/LockReleaseTokenPoolAndProxy/LockReleaseTokenPoolAndProxy.bin LockReleaseTokenPoolAndProxy lock_release_token_pool_and_proxy
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.24/TokenPool/TokenPool.abi ../../../contracts/solc/v0.8.24/TokenPool/TokenPool.bin TokenPool token_pool
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.24/USDCTokenPool/USDCTokenPool.abi ../../../contracts/solc/v0.8.24/USDCTokenPool/USDCTokenPool.bin USDCTokenPool usdc_token_pool
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.24/BurnWithFromMintRebasingTokenPool/BurnWithFromMintRebasingTokenPool.abi ../../../contracts/solc/v0.8.24/BurnWithFromMintRebasingTokenPool/BurnWithFromMintRebasingTokenPool.bin BurnWithFromMintRebasingTokenPool burn_with_from_mint_rebasing_token_pool

// Helpers
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.24/MockV3Aggregator/MockV3Aggregator.abi ../../../contracts/solc/v0.8.24/MockV3Aggregator/MockV3Aggregator.bin MockV3Aggregator mock_v3_aggregator_contract
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.24/MaybeRevertMessageReceiver/MaybeRevertMessageReceiver.abi ../../../contracts/solc/v0.8.24/MaybeRevertMessageReceiver/MaybeRevertMessageReceiver.bin MaybeRevertMessageReceiver maybe_revert_message_receiver
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.24/PingPongDemo/PingPongDemo.abi ../../../contracts/solc/v0.8.24/PingPongDemo/PingPongDemo.bin PingPongDemo ping_pong_demo
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.24/SelfFundedPingPong/SelfFundedPingPong.abi ../../../contracts/solc/v0.8.24/SelfFundedPingPong/SelfFundedPingPong.bin SelfFundedPingPong self_funded_ping_pong
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.24/MessageHasher/MessageHasher.abi ../../../contracts/solc/v0.8.24/MessageHasher/MessageHasher.bin MessageHasher message_hasher
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.24/MultiOCR3Helper/MultiOCR3Helper.abi ../../../contracts/solc/v0.8.24/MultiOCR3Helper/MultiOCR3Helper.bin MultiOCR3Helper multi_ocr3_helper
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.24/CCIPReaderTester/CCIPReaderTester.abi ../../../contracts/solc/v0.8.24/CCIPReaderTester/CCIPReaderTester.bin CCIPReaderTester ccip_reader_tester
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.24/USDCReaderTester/USDCReaderTester.abi ../../../contracts/solc/v0.8.24/USDCReaderTester/USDCReaderTester.bin USDCReaderTester usdc_reader_tester
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.24/ReportCodec/ReportCodec.abi ../../../contracts/solc/v0.8.24/ReportCodec/ReportCodec.bin ReportCodec report_codec
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.24/EtherSenderReceiver/EtherSenderReceiver.abi ../../../contracts/solc/v0.8.24/EtherSenderReceiver/EtherSenderReceiver.bin EtherSenderReceiver ether_sender_receiver
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.24/WETH9/WETH9.abi ../../../contracts/solc/v0.8.24/WETH9/WETH9.bin WETH9 weth9
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.24/MockE2EUSDCTokenMessenger/MockE2EUSDCTokenMessenger.abi ../../../contracts/solc/v0.8.24/MockE2EUSDCTokenMessenger/MockE2EUSDCTokenMessenger.bin MockE2EUSDCTokenMessenger mock_usdc_token_messenger
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.24/MockE2EUSDCTransmitter/MockE2EUSDCTransmitter.abi ../../../contracts/solc/v0.8.24/MockE2EUSDCTransmitter/MockE2EUSDCTransmitter.bin MockE2EUSDCTransmitter mock_usdc_token_transmitter

// EncodingUtils
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.24/ICCIPEncodingUtils/ICCIPEncodingUtils.abi ../../../contracts/solc/v0.8.24/ICCIPEncodingUtils/ICCIPEncodingUtils.bin EncodingUtils ccip_encoding_utils

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
