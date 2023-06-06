// Package gethwrappers provides tools for wrapping solidity contracts with
// golang packages, using abigen.
package gethwrappers

// Chainlink Functions (OCR2DR)
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.16/Verifier.abi ../../contracts/solc/v0.8.16/Verifier.bin LLOVerifier llo_feeds
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.16/VerifierProxy.abi ../../contracts/solc/v0.8.16/VerifierProxy.bin LLOVerifierProxy llo_feeds
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.16/ErroredVerifier.abi ../../contracts/solc/v0.8.16/ErroredVerifier.bin LLOErroredVerifier llo_feeds_test
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.16/ExposedVerifier.abi ../../contracts/solc/v0.8.16/ExposedVerifier.bin LLOExposedVerifier llo_feeds_test
