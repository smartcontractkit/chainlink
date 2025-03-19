// Package gethwrappers provides tools for wrapping solidity contracts with
// golang packages, using abigen.
package gethwrappers

// Make sure solidity compiler artifacts are up-to-date. Only output stdout on failure.
//go:generate ../../contracts/scripts/native_solc_compile_all

//go:generate go run ./generation/generate/wrap.go OffchainAggregator/OffchainAggregator.abi - OffchainAggregator offchain_aggregator_wrapper
//go:generate go run ./generation/generate_link/wrap_link.go

//go:generate go generate go_generate_automation.go
//go:generate go generate go_generate_vrf.go
//go:generate go generate ./functions
//go:generate go generate ./keystone
//go:generate go generate ./llo-feeds
//go:generate go generate ./operatorforwarder
//go:generate go generate ./shared
//go:generate go generate ./ccip
//go:generate go generate ./workflow
