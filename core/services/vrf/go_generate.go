package vrf

// See comments in solc.sh for details about this command

//go:generate abigen --sol ../../../evm/v0.5/contracts/VRFAll.sol --solc ../../../tools/bin/solc.sh --pkg vrf --out solidity_verifier_interface.go
