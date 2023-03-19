// Package gethwrappers provides tools for wrapping solidity contracts with
// golang packages, using abigen.
package gethwrappers

// Transmission
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.15/Greeter.abi ../../../contracts/solc/v0.8.15/Greeter.bin Greeter greeter
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.15/SCA.abi ../../../contracts/solc/v0.8.15/SCA.bin SCA sca
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.15/SmartContractAccountFactory.abi ../../../contracts/solc/v0.8.15/SmartContractAccountFactory.bin SmartContractAccountFactory smart_contract_account_factory
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.15/EntryPoint.abi ../../../contracts/solc/v0.8.15/EntryPoint.bin EntryPoint entry_point
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.15/SmartContractAccountHelper.abi ../../../contracts/solc/v0.8.15/SmartContractAccountHelper.bin SmartContractAccountHelper smart_contract_account_helper
