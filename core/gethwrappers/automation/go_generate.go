// Package gethwrappers provides tools for wrapping solidity contracts with
// golang packages, using abigen.
package gethwrappers

// Automation
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.16/ComposerCompatibleInterfaceV1.abi ../../../contracts/solc/v0.8.16/ComposerCompatibleInterfaceV1.bin ComposerCompatibleInterfaceV1 composer_compatible_interface
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.16/MercuryRegistryComposer.abi ../../../contracts/solc/v0.8.16/MercuryRegistryComposer.bin MercuryRegistryComposer mercury_registry_composer
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.16/ComposerCrossChainSend.abi ../../../contracts/solc/v0.8.16/ComposerCrossChainSend.bin ComposerCrossChainSend composer_cross_chain
