// Package gethwrappers provides tools for wrapping solidity contracts with
// golang packages, using abigen.
package gethwrappers

// Automation
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.6/MercuryRegistry.abi ../../../contracts/solc/v0.8.6/MercuryRegistry.bin MercuryRegistry mercury_registry_wrapper
