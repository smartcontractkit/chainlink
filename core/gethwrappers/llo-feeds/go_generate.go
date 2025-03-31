// Package gethwrappers provides tools for wrapping solidity contracts with
// golang packages, using abigen.
package gethwrappers

// Chainlink LLO V0.3.0
//go:generate go run ../generation/wrap.go llo-feeds/v0.3.0 Verifier verifier
//go:generate go run ../generation/wrap.go llo-feeds/v0.3.0 VerifierProxy verifier_proxy
//go:generate go run ../generation/wrap.go llo-feeds/v0.3.0 RewardManager reward_manager
//go:generate go run ../generation/wrap.go llo-feeds/v0.3.0 FeeManager fee_manager
//go:generate go run ../generation/wrap.go llo-feeds/v0.3.0 ErroredVerifier errored_verifier
//go:generate go run ../generation/wrap.go llo-feeds/v0.3.0 ExposedVerifier exposed_verifier

// Chainlink LLO V0.4.0
//go:generate go run ../generation/wrap.go llo-feeds/v0.4.0 DestinationVerifier destination_verifier
//go:generate go run ../generation/wrap.go llo-feeds/v0.4.0 DestinationVerifierProxy destination_verifier_proxy
//go:generate go run ../generation/wrap.go llo-feeds/v0.4.0 DestinationFeeManager destination_fee_manager
//go:generate go run ../generation/wrap.go llo-feeds/v0.4.0 DestinationRewardManager destination_reward_manager

// Chainlink LLO V0.5.0
//go:generate go run ../generation/wrap.go llo-feeds/v0.5.0 Verifier verifier_v0_5_0
//go:generate go run ../generation/wrap.go llo-feeds/v0.5.0 VerifierProxy verifier_proxy_v0_5_0
//go:generate go run ../generation/wrap.go llo-feeds/v0.5.0 FeeManager fee_manager_v0_5_0
//go:generate go run ../generation/wrap.go llo-feeds/v0.5.0 RewardManager reward_manager_v0_5_0

//go:generate go run ../generation/wrap.go llo-feeds/v0.5.0 ChannelConfigStore channel_config_store
//go:generate go run ../generation/wrap.go llo-feeds/v0.5.0 Configurator configurator
//go:generate go run ../generation/wrap.go llo-feeds/v0.5.0 ExposedConfigurator exposed_configurator
//go:generate go run ../generation/wrap.go llo-feeds/v0.5.0 MockFeeManager mock_fee_manager_v0_5_0
