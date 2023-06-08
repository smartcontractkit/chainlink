// Package gethwrappers provides tools for wrapping solidity contracts with
// golang packages, using abigen.
package gethwrappers

// OCR2VRF - remove the _disabled tag to run these locally.
//go:generate_disabled go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.19/DKG.abi ../../../contracts/solc/v0.8.19/DKG.bin DKG dkg
//go:generate_disabled go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.19/VRFCoordinator.abi ../../../contracts/solc/v0.8.19/VRFCoordinator.bin VRFCoordinator vrf_coordinator
//go:generate_disabled go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.19/VRFBeacon.abi ../../../contracts/solc/v0.8.19/VRFBeacon.bin VRFBeacon vrf_beacon
//go:generate_disabled go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.19/BeaconVRFConsumer.abi ../../../contracts/solc/v0.8.19/BeaconVRFConsumer.bin BeaconVRFConsumer vrf_beacon_consumer
//go:generate_disabled go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.19/LoadTestBeaconVRFConsumer.abi ../../../contracts/solc/v0.8.19/LoadTestBeaconVRFConsumer.bin LoadTestBeaconVRFConsumer load_test_beacon_consumer
