// Package gethwrappers provides tools for wrapping solidity contracts with
// golang packages, using abigen.
package gethwrappers

// Chainlink LLO V0.3.0
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.19/llo-feeds/v0.3.0/Verifier/Verifier.abi ../../../contracts/solc/v0.8.19/llo-feeds/v0.3.0/Verifier/Verifier.bin Verifier verifier
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.19/llo-feeds/v0.3.0/VerifierProxy/VerifierProxy.abi ../../../contracts/solc/v0.8.19/llo-feeds/v0.3.0/VerifierProxy/VerifierProxy.bin VerifierProxy verifier_proxy
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.19/llo-feeds/v0.3.0/RewardManager/RewardManager.abi ../../../contracts/solc/v0.8.19/llo-feeds/v0.3.0/RewardManager/RewardManager.bin RewardManager reward_manager
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.19/llo-feeds/v0.3.0/FeeManager/FeeManager.abi ../../../contracts/solc/v0.8.19/llo-feeds/v0.3.0/FeeManager/FeeManager.bin FeeManager fee_manager
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.19/llo-feeds/v0.3.0/test/mocks/ErroredVerifier/ErroredVerifier.abi ../../../contracts/solc/v0.8.19/llo-feeds/v0.3.0/test/mocks/ErroredVerifier/ErroredVerifier.bin ErroredVerifier errored_verifier
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.19/llo-feeds/v0.3.0/test/mocks/ExposedVerifier/ExposedVerifier.abi ../../../contracts/solc/v0.8.19/llo-feeds/v0.3.0/test/mocks/ExposedVerifier/ExposedVerifier.bin ExposedVerifier exposed_verifier

// Chainlink LLO V0.4.0
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.19/llo-feeds/v0.4.0/DestinationVerifier/DestinationVerifier.abi ../../../contracts/solc/v0.8.19/llo-feeds/v0.4.0/DestinationVerifier/DestinationVerifier.bin DestinationVerifier destination_verifier
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.19/llo-feeds/v0.4.0/DestinationVerifierProxy/DestinationVerifierProxy.abi ../../../contracts/solc/v0.8.19/llo-feeds/v0.4.0/DestinationVerifierProxy/DestinationVerifierProxy.bin DestinationVerifierProxy destination_verifier_proxy
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.19/llo-feeds/v0.4.0/DestinationFeeManager/DestinationFeeManager.abi ../../../contracts/solc/v0.8.19/llo-feeds/v0.4.0/DestinationFeeManager/DestinationFeeManager.bin DestinationFeeManager destination_fee_manager
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.19/llo-feeds/v0.4.0/DestinationRewardManager/DestinationRewardManager.abi ../../../contracts/solc/v0.8.19/llo-feeds/v0.4.0/DestinationRewardManager/DestinationRewardManager.bin DestinationRewardManager destination_reward_manager

// Chainlink LLO V0.5.0
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.19/llo-feeds/v0.5.0/Verifier/Verifier.abi ../../../contracts/solc/v0.8.19/llo-feeds/v0.5.0/Verifier/Verifier.bin Verifier verifier_v0_5_0
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.19/llo-feeds/v0.5.0/VerifierProxy/VerifierProxy.abi ../../../contracts/solc/v0.8.19/llo-feeds/v0.5.0/VerifierProxy/VerifierProxy.bin VerifierProxy verifier_proxy_v0_5_0
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.19/llo-feeds/v0.5.0/FeeManager/FeeManager.abi ../../../contracts/solc/v0.8.19/llo-feeds/v0.5.0/FeeManager/FeeManager.bin FeeManager fee_manager_v0_5_0
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.19/llo-feeds/v0.5.0/RewardManager/RewardManager.abi ../../../contracts/solc/v0.8.19/llo-feeds/v0.5.0/RewardManager/RewardManager.bin RewardManager reward_manager_v0_5_0

//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.19/llo-feeds/v0.5.0/configuration/ChannelConfigStore/ChannelConfigStore.abi ../../../contracts/solc/v0.8.19/llo-feeds/v0.5.0/configuration/ChannelConfigStore/ChannelConfigStore.bin ChannelConfigStore channel_config_store
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.19/llo-feeds/v0.5.0/configuration/Configurator/Configurator.abi ../../../contracts/solc/v0.8.19/llo-feeds/v0.5.0/configuration/Configurator/Configurator.bin Configurator configurator
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.19/llo-feeds/v0.5.0/configuration/test/mocks/ExposedConfigurator/ExposedConfigurator.abi ../../../contracts/solc/v0.8.19/llo-feeds/v0.5.0/configuration/test/mocks/ExposedConfigurator/ExposedConfigurator.bin ExposedConfigurator exposed_configurator
