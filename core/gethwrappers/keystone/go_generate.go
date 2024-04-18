// Package gethwrappers provides tools for wrapping solidity contracts with
// golang packages, using abigen.
package gethwrappers

// Keystone

//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.19/keystone/KeystoneForwarder.sol/KeystoneForwarder.abi.json ../../../contracts/solc/v0.8.19/keystone/KeystoneForwarder.sol/KeystoneForwarder.bin KeystoneForwarder forwarder
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.19/keystone/OCR3Capability.sol/OCR3Capability.abi.json ../../../contracts/solc/v0.8.19/keystone/OCR3Capability.sol/OCR3Capability.bin OCR3Capability ocr3_capability
