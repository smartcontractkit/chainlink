// Package gethwrappers provides tools for wrapping solidity contracts with
// golang packages, using abigen.
package gethwrappers

// OCR2VRF
//go:generate go run ../generation/generate/wrap.go ../../../../contracts/solc/v0.8.15/DKG.abi ../../../../contracts/solc/v0.8.15/DKG.bin DKG dkg
//go:generate go run ../generation/generate/wrap.go ../../../../contracts/solc/v0.8.15/VRFBeaconCoordinator.abi ../../../../contracts/solc/v0.8.15/VRFBeaconCoordinator.bin VRFBeaconCoordinator vrf_beacon_coordinator
//go:generate go run ../generation/generate/wrap.go ../../../../contracts/solc/v0.8.15/BeaconVRFConsumer.abi ../../../../contracts/solc/v0.8.15/BeaconVRFConsumer.bin BeaconVRFConsumer vrf_beacon_consumer
