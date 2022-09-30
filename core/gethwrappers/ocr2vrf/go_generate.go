// Package gethwrappers provides tools for wrapping solidity contracts with
// golang packages, using abigen.
package gethwrappers

// OCR2VRF
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.15/DKG.abi ../../../contracts/solc/v0.8.15/DKG.bin DKG dkg
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.15/VRFCoordinator.abi ../../../contracts/solc/v0.8.15/VRFCoordinator.bin VRFCoordinator vrf_coordinator
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.15/VRFBeacon.abi ../../../contracts/solc/v0.8.15/VRFBeacon.bin VRFBeacon vrf_beacon
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.15/VRFCoordinatorProxy.abi ../../../contracts/solc/v0.8.15/VRFCoordinatorProxy.bin VRFCoordinatorProxy vrf_coordinator_proxy
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.15/VRFBeaconProxy.abi ../../../contracts/solc/v0.8.15/VRFBeaconProxy.bin VRFBeaconProxy vrf_beacon_proxy
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.15/VRFProxyAdmin.abi ../../../contracts/solc/v0.8.15/VRFProxyAdmin.bin VRFProxyAdmin vrf_proxy_admin
//go:generate go run ../generation/generate/wrap.go ../../../contracts/solc/v0.8.15/BeaconVRFConsumer.abi ../../../contracts/solc/v0.8.15/BeaconVRFConsumer.bin BeaconVRFConsumer vrf_beacon_consumer
