// Package gethwrappers provides tools for wrapping solidity contracts with
// golang packages, using abigen.
package gethwrappers

// Chainlink Functions (OCR2DR)
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/Functions.abi ../../contracts/solc/v0.8.6/Functions.bin OCR2DR ocr2dr
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/FunctionsClient.abi ../../contracts/solc/v0.8.6/FunctionsClient.bin OCR2DRClient ocr2dr_client
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/FunctionsClientExample.abi ../../contracts/solc/v0.8.6/FunctionsClientExample.bin OCR2DRClientExample ocr2dr_client_example
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/FunctionsOracleWithInit.abi ../../contracts/solc/v0.8.6/FunctionsOracleWithInit.bin OCR2DROracle ocr2dr_oracle
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/FunctionsBillingRegistryWithInit.abi ../../contracts/solc/v0.8.6/FunctionsBillingRegistryWithInit.bin OCR2DRRegistry ocr2dr_registry
