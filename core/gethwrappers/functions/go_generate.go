// Package gethwrappers provides tools for wrapping solidity contracts with
// golang packages, using abigen.
package gethwrappers

// Chainlink Functions (OCR2DR)

// Version 1 (Mainnet Preview)
//go:generate go run ../generation/wrap.go functions FunctionsClient functions_client
//go:generate go run ../generation/wrap.go functions FunctionsClientExample functions_client_example
//go:generate go run ../generation/wrap.go functions FunctionsLoadTestClient functions_load_test_client
//go:generate go run ../generation/wrap.go functions FunctionsCoordinator functions_coordinator
//go:generate go run ../generation/wrap.go functions FunctionsRouter functions_router
//go:generate go run ../generation/wrap.go functions TermsOfServiceAllowList functions_allow_list
//go:generate go run ../generation/wrap.go functions FunctionsV1EventsMock functions_v1_events_mock
