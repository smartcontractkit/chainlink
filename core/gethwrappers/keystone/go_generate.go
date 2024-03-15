// Package gethwrappers provides tools for wrapping solidity contracts with
// golang packages, using abigen.
package gethwrappers

// Keystone

//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.19/KeystoneForwarder/KeystoneForwarder.abi ../../../contracts/solc/v0.8.19/KeystoneForwarder/KeystoneForwarder.bin KeystoneForwarder forwarder
