// Package gethwrappers provides tools for wrapping solidity contracts with
// golang packages, using abigen.
package gethwrappers

// Chainlink LLO
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.19/Verifier/Verifier.abi ../../../contracts/solc/v0.8.19/Verifier/Verifier.bin Verifier verifier
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.19/VerifierProxy/VerifierProxy.abi ../../../contracts/solc/v0.8.19/VerifierProxy/VerifierProxy.bin VerifierProxy verifier_proxy
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.19/ErroredVerifier/ErroredVerifier.abi ../../../contracts/solc/v0.8.19/ErroredVerifier/ErroredVerifier.bin ErroredVerifier errored_verifier
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.19/ExposedVerifier/ExposedVerifier.abi ../../../contracts/solc/v0.8.19/ExposedVerifier/ExposedVerifier.bin ExposedVerifier exposed_verifier
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.19/RewardManager/RewardManager.abi ../../../contracts/solc/v0.8.19/RewardManager/RewardManager.bin RewardManager reward_manager
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.19/FeeManager/FeeManager.abi ../../../contracts/solc/v0.8.19/FeeManager/FeeManager.bin FeeManager fee_manager
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.19/ChannelConfigStore/ChannelConfigStore.abi ../../../contracts/solc/v0.8.19/ChannelConfigStore/ChannelConfigStore.bin ChannelConfigStore channel_config_store
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.19/ChannelVerifier/ChannelVerifier.abi ../../../contracts/solc/v0.8.19/ChannelVerifier/ChannelVerifier.bin ChannelVerifier channel_verifier
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.19/ExposedChannelVerifier/ExposedChannelVerifier.abi ../../../contracts/solc/v0.8.19/ExposedChannelVerifier/ExposedChannelVerifier.bin ExposedChannelVerifier exposed_channel_verifier
