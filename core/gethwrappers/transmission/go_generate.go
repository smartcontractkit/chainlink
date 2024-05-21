// Package gethwrappers provides tools for wrapping solidity contracts with
// golang packages, using abigen.
package gethwrappers

// Transmission
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.19/Greeter/Greeter.abi ../../../contracts/solc/v0.8.19/Greeter/Greeter.bin Greeter greeter_wrapper
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.19/SmartContractAccountFactory/SmartContractAccountFactory.abi ../../../contracts/solc/v0.8.19/SmartContractAccountFactory/SmartContractAccountFactory.bin SmartContractAccountFactory smart_contract_account_factory
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.19/EntryPoint/EntryPoint.abi ../../../contracts/solc/v0.8.19/EntryPoint/EntryPoint.bin EntryPoint entry_point
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.19/SmartContractAccountHelper/SmartContractAccountHelper.abi ../../../contracts/solc/v0.8.19/SmartContractAccountHelper/SmartContractAccountHelper.bin SmartContractAccountHelper smart_contract_account_helper
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.19/SCA/SCA.abi ../../../contracts/solc/v0.8.19/SCA/SCA.bin SCA sca_wrapper
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.19/Paymaster/Paymaster.abi ../../../contracts/solc/v0.8.19/Paymaster/Paymaster.bin Paymaster paymaster_wrapper
