// Package gethwrappers provides tools for wrapping solidity contracts with
// golang packages, using abigen.
package gethwrappers

// Make sure solidity compiler artifacts are up to date. Only output stdout on failure.
//go:generate ./generation/compile_contracts.sh

//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/KeeperRegistryLogicA2_1.abi ../../contracts/solc/v0.8.6/KeeperRegistryLogicA2_1.bin KeeperRegistryLogicA keeper_registry_logic_a_wrapper_2_1
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/KeeperRegistryLogicB2_1.abi ../../contracts/solc/v0.8.6/KeeperRegistryLogicB2_1.bin KeeperRegistryLogicB keeper_registry_logic_b_wrapper_2_1
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/IKeeperRegistryMaster.abi ../../contracts/solc/v0.8.6/IKeeperRegistryMaster.bin IKeeperRegistryMaster i_keeper_registry_master_wrapper_2_1
