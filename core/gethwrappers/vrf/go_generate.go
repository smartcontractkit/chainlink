// Package gethwrappers provides tools for wrapping solidity contracts with
// golang packages, using abigen.
package gethwrappers

// Transmission
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.6/VRFV2PlusConsumerExample.abi ../../../contracts/solc/v0.8.6/VRFV2PlusConsumerExample.bin VRFV2PlusConsumerExample vrfv2plus_consumer_example
