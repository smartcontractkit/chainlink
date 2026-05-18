// Package gethwrappers provides tools for wrapping solidity contracts with
// golang packages, using abigen.
package gethwrappers

// Keystone

//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.24/KeystoneForwarder/KeystoneForwarder.abi ../../../contracts/solc/v0.8.24/KeystoneForwarder/KeystoneForwarder.bin KeystoneForwarder forwarder
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.24/OCR3Capability/OCR3Capability.abi ../../../contracts/solc/v0.8.24/OCR3Capability/OCR3Capability.bin OCR3Capability ocr3_capability
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.24/CapabilitiesRegistry/CapabilitiesRegistry.abi ../../../contracts/solc/v0.8.24/CapabilitiesRegistry/CapabilitiesRegistry.bin CapabilitiesRegistry capabilities_registry
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.24/KeystoneFeedsConsumer/KeystoneFeedsConsumer.abi ../../../contracts/solc/v0.8.24/KeystoneFeedsConsumer/KeystoneFeedsConsumer.bin KeystoneFeedsConsumer feeds_consumer
