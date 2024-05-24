// Package gethwrappers provides tools for wrapping solidity contracts with
// golang packages, using abigen.
package gethwrappers

// Log tester
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.19/LogEmitter/LogEmitter.abi ../../contracts/solc/v0.8.19/LogEmitter/LogEmitter.bin LogEmitter log_emitter
//go:generate go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.19/VRFLogEmitter/VRFLogEmitter.abi ../../contracts/solc/v0.8.19/VRFLogEmitter/VRFLogEmitter.bin VRFLogEmitter vrf_log_emitter
