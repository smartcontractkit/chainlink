// Package gethwrappers provides tools for wrapping solidity contracts with
// golang packages, using abigen.
package gethwrappers

// OCR2VRF
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.15/DKG.abi ../../../contracts/solc/v0.8.15/DKG.bin DKG dkg
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.15/VRFCoordinator.abi ../../../contracts/solc/v0.8.15/VRFCoordinator.bin VRFCoordinator vrf_coordinator
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.15/VRFBeacon.abi ../../../contracts/solc/v0.8.15/VRFBeacon.bin VRFBeacon vrf_beacon
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.15/BeaconVRFConsumer.abi ../../../contracts/solc/v0.8.15/BeaconVRFConsumer.bin BeaconVRFConsumer vrf_beacon_consumer
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.15/LoadTestBeaconVRFConsumer.abi ../../../contracts/solc/v0.8.15/LoadTestBeaconVRFConsumer.bin LoadTestBeaconVRFConsumer load_test_beacon_consumer
