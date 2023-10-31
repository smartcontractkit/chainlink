// Package gethwrappers provides tools for wrapping solidity contracts with
// golang packages, using abigen.
package gethwrappers

// Chainlink Functions (OCR2DR)

// Version 1 (Mainnet Preview)
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.19/functions/v1_X/FunctionsRequest.abi ../../../contracts/solc/v0.8.19/functions/v1_X/FunctionsRequest.bin Functions functions
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.19/functions/v1_X/FunctionsClient.abi ../../../contracts/solc/v0.8.19/functions/v1_X/FunctionsClient.bin FunctionsClient functions_client
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.19/functions/v1_X/FunctionsClientExample.abi ../../../contracts/solc/v0.8.19/functions/v1_X/FunctionsClientExample.bin FunctionsClientExample functions_client_example
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.19/functions/v1_X/FunctionsLoadTestClient.abi ../../../contracts/solc/v0.8.19/functions/v1_X/FunctionsLoadTestClient.bin FunctionsLoadTestClient functions_load_test_client
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.19/functions/v1_X/FunctionsCoordinator.abi ../../../contracts/solc/v0.8.19/functions/v1_X/FunctionsCoordinator.bin FunctionsCoordinator functions_coordinator
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.19/functions/v1_X/FunctionsRouter.abi ../../../contracts/solc/v0.8.19/functions/v1_X/FunctionsRouter.bin FunctionsRouter functions_router
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.19/functions/v1_X/TermsOfServiceAllowList.abi ../../../contracts/solc/v0.8.19/functions/v1_X/TermsOfServiceAllowList.bin TermsOfServiceAllowList functions_allow_list
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.19/functions/v1_X/FunctionsV1EventsMock.abi ../../../contracts/solc/v0.8.19/functions/v1_X/FunctionsV1EventsMock.bin FunctionsV1EventsMock functions_v1_events_mock
