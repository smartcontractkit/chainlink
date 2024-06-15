// Package gethwrappers provides tools for wrapping solidity contracts with
// golang packages, using abigen.
package gethwrappers

// Keystone

//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.19/KeystoneForwarder/KeystoneForwarder.abi ../../../contracts/solc/v0.8.19/KeystoneForwarder/KeystoneForwarder.bin KeystoneForwarder forwarder
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.19/OCR3Capability/OCR3Capability.abi ../../../contracts/solc/v0.8.19/OCR3Capability/OCR3Capability.bin OCR3Capability ocr3_capability
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.19/CapabilityRegistry/CapabilityRegistry.abi ../../../contracts/solc/v0.8.19/CapabilityRegistry/CapabilityRegistry.bin CapabilityRegistry keystone_capability_registry
