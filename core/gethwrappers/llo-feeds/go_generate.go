// Package gethwrappers provides tools for wrapping solidity contracts with
// golang packages, using abigen.
package gethwrappers

// Chainlink LLO
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.16/Verifier/Verifier.abi ../../../contracts/solc/v0.8.16/Verifier/Verifier.bin Verifier verifier
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.16/VerifierProxy/VerifierProxy.abi ../../../contracts/solc/v0.8.16/VerifierProxy/VerifierProxy.bin VerifierProxy verifier_proxy
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.16/ErroredVerifier/ErroredVerifier.abi ../../../contracts/solc/v0.8.16/ErroredVerifier/ErroredVerifier.bin ErroredVerifier errored_verifier
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.16/ExposedVerifier/ExposedVerifier.abi ../../../contracts/solc/v0.8.16/ExposedVerifier/ExposedVerifier.bin ExposedVerifier exposed_verifier
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.16/RewardManager/RewardManager.abi ../../../contracts/solc/v0.8.16/RewardManager/RewardManager.bin RewardManager reward_manager
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.16/FeeManager/FeeManager.abi ../../../contracts/solc/v0.8.16/FeeManager/FeeManager.bin FeeManager fee_manager
