// Package gethwrappers provides tools for wrapping solidity contracts with
// golang packages, using abigen.
package gethwrappers

// Chainlink Functions (OCR2DR)

// Version 0 (Testnet Beta)
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.6/functions/v0_0_0/Functions.abi ../../../contracts/solc/v0.8.6/functions/v0_0_0/Functions.bin OCR2DR ocr2dr
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.6/functions/v0_0_0/FunctionsClient.abi ../../../contracts/solc/v0.8.6/functions/v0_0_0/FunctionsClient.bin OCR2DRClient ocr2dr_client
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.6/functions/v0_0_0/FunctionsClientExample.abi ../../../contracts/solc/v0.8.6/functions/v0_0_0/FunctionsClientExample.bin OCR2DRClientExample ocr2dr_client_example
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.6/functions/v0_0_0/FunctionsOracleWithInit.abi ../../../contracts/solc/v0.8.6/functions/v0_0_0/FunctionsOracleWithInit.bin OCR2DROracle ocr2dr_oracle
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.6/functions/v0_0_0/FunctionsBillingRegistryWithInit.abi ../../../contracts/solc/v0.8.6/functions/v0_0_0/FunctionsBillingRegistryWithInit.bin OCR2DRRegistry ocr2dr_registry

// Version 1 (Mainnet Preview)
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.19/functions/v1_0_0/FunctionsRequest.abi ../../../contracts/solc/v0.8.19/functions/v1_0_0/FunctionsRequest.bin Functions functions
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.19/functions/v1_0_0/FunctionsClient.abi ../../../contracts/solc/v0.8.19/functions/v1_0_0/FunctionsClient.bin FunctionsClient functions_client
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.19/functions/v1_0_0/FunctionsClientExample.abi ../../../contracts/solc/v0.8.19/functions/v1_0_0/FunctionsClientExample.bin FunctionsClientExample functions_client_example
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.19/functions/v1_0_0/FunctionsLoadTestClient.abi ../../../contracts/solc/v0.8.19/functions/v1_0_0/FunctionsLoadTestClient.bin FunctionsLoadTestClient functions_load_test_client
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.19/functions/v1_0_0/FunctionsCoordinator.abi ../../../contracts/solc/v0.8.19/functions/v1_0_0/FunctionsCoordinator.bin FunctionsCoordinator functions_coordinator
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.19/functions/v1_0_0/FunctionsRouter.abi ../../../contracts/solc/v0.8.19/functions/v1_0_0/FunctionsRouter.bin FunctionsRouter functions_router
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.19/functions/v1_0_0/TermsOfServiceAllowList.abi ../../../contracts/solc/v0.8.19/functions/v1_0_0/TermsOfServiceAllowList.bin TermsOfServiceAllowList functions_allow_list
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.19/functions/v1_0_0/FunctionsV1EventsMock.abi ../../../contracts/solc/v0.8.19/functions/v1_0_0/FunctionsV1EventsMock.bin FunctionsV1EventsMock functions_v1_events_mock
