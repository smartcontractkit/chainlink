// Package gethwrappers provides tools for wrapping solidity contracts with
// golang packages, using abigen.
package gethwrappers

// Chainlink LLO
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.16/Verifier.abi ../../../contracts/solc/v0.8.16/Verifier.bin Verifier verifier
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.16/VerifierProxy.abi ../../../contracts/solc/v0.8.16/VerifierProxy.bin VerifierProxy verifier_proxy
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.16/ErroredVerifier.abi ../../../contracts/solc/v0.8.16/ErroredVerifier.bin ErroredVerifier errored_verifier
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.16/ExposedVerifier.abi ../../../contracts/solc/v0.8.16/ExposedVerifier.bin ExposedVerifier exposed_verifier
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.16/RewardManager.abi ../../../contracts/solc/v0.8.16/RewardManager.bin RewardManager reward_manager
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.16/FeeManager.abi ../../../contracts/solc/v0.8.16/FeeManager.bin FeeManager fee_manager
