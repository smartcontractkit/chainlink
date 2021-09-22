// Package gethwrappers provides tools for wrapping solidity contracts with
// golang packages, using abigen.
package gethwrappers

// Make sure solidity compiler artifacts are up to date. Only output stdout on failure.
//go:generate ./generation/compile_contracts.sh

//go:generate go run ./generation/generate/wrap.go ../../../contracts/solc/v0.6/FluxAggregator.abi ../../../contracts/solc/v0.6/FluxAggregator.bin FluxAggregator flux_aggregator_wrapper
//go:generate go run ./generation/generate/wrap.go ../../../contracts/solc/v0.6/VRFTestHelper.abi ../../../contracts/solc/v0.6/VRFTestHelper.bin VRFTestHelper solidity_vrf_verifier_wrapper
//go:generate go run ./generation/generate/wrap.go ../../../contracts/solc/v0.6/VRFCoordinator.abi ../../../contracts/solc/v0.6/VRFCoordinator.bin VRFCoordinator solidity_vrf_coordinator_interface
//go:generate go run ./generation/generate/wrap.go ../../../contracts/solc/v0.6/VRFConsumer.abi ../../../contracts/solc/v0.6/VRFConsumer.bin VRFConsumer solidity_vrf_consumer_interface
//go:generate go run ./generation/generate/wrap.go ../../../contracts/solc/v0.6/VRFRequestIDBaseTestHelper.abi ../../../contracts/solc/v0.6/VRFRequestIDBaseTestHelper.bin VRFRequestIDBaseTestHelper solidity_vrf_request_id
//go:generate go run ./generation/generate/wrap.go ../../../contracts/solc/v0.6/Flags.abi ../../../contracts/solc/v0.6/Flags.bin Flags flags_wrapper
//go:generate go run ./generation/generate/wrap.go ../../../contracts/solc/v0.6/Oracle.abi ../../../contracts/solc/v0.6/Oracle.bin Oracle oracle_wrapper
//go:generate go run ./generation/generate/wrap.go ../../../contracts/solc/v0.7/Consumer.abi ../../../contracts/solc/v0.7/Consumer.bin Consumer consumer_wrapper
//go:generate go run ./generation/generate/wrap.go ../../../contracts/solc/v0.7/MultiWordConsumer.abi ../../../contracts/solc/v0.7/MultiWordConsumer.bin MultiWordConsumer multiwordconsumer_wrapper
//go:generate go run ./generation/generate/wrap.go ../../../contracts/solc/v0.7/Operator.abi ../../../contracts/solc/v0.7/Operator.bin Operator operator_wrapper
//go:generate go run ./generation/generate/wrap.go OffchainAggregator/OffchainAggregator.abi - OffchainAggregator offchain_aggregator_wrapper

//go:generate go run ./generation/generate/wrap.go ../../../contracts/solc/v0.7/KeeperRegistry.abi ../../../contracts/solc/v0.7/KeeperRegistry.bin KeeperRegistry keeper_registry_wrapper

// v0.8 VRFConsumer
//go:generate go run ./generation/generate/wrap.go ../../../contracts/solc/v0.8/VRFConsumer.abi ../../../contracts/solc/v0.8/VRFConsumer.bin VRFConsumer solidity_vrf_consumer_interface_v08
//go:generate go run ./generation/generate/wrap.go ../../../contracts/solc/v0.8/VRFRequestIDBaseTestHelper.abi ../../../contracts/solc/v0.8/VRFRequestIDBaseTestHelper.bin VRFRequestIDBaseTestHelper solidity_vrf_request_id_v08

//go:generate mockery --recursive --name FluxAggregatorInterface --output ../mocks/ --case=underscore --structname FluxAggregator --filename flux_aggregator.go
//go:generate mockery --recursive --name FlagsInterface --output ../mocks/ --case=underscore --structname Flags --filename flags.go

//go:generate go run ./generation/generate_link/wrap_link.go

// VRF V2
//go:generate go run ./generation/generate/wrap.go ../../../contracts/solc/v0.8/VRFCoordinatorV2.abi ../../../contracts/solc/v0.8/VRFCoordinatorV2.bin VRFCoordinatorV2 vrf_coordinator_v2
//go:generate go run ./generation/generate/wrap.go ../../../contracts/solc/v0.8/VRFConsumerV2.abi ../../../contracts/solc/v0.8/VRFConsumerV2.bin VRFConsumerV2 vrf_consumer_v2
//go:generate go run ./generation/generate/wrap.go ../../../contracts/solc/v0.8/VRFMaliciousConsumerV2.abi ../../../contracts/solc/v0.8/VRFMaliciousConsumerV2.bin VRFMaliciousConsumerV2 vrf_malicious_consumer_v2
//go:generate go run ./generation/generate/wrap.go ../../../contracts/solc/v0.8/VRFTestHelper.abi ../../../contracts/solc/v0.8/VRFTestHelper.bin VRFV08TestHelper solidity_vrf_v08_verifier_wrapper

// To run these commands, you must either install docker, or the correct version
// of abigen. The latter can be installed with these commands, at least on linux:
//
//   git clone https://github.com/ethereum/go-ethereum
//   cd go-ethereum/cmd/abigen
//   git checkout v<version-needed>
//   go install
//
// Here, <version-needed> is the version of go-ethereum specified in chainlink's
// go.mod. This will install abigen in "$GOPATH/bin", which you should add to
// your $PATH.
//
// To reduce explicit dependencies, and in case the system does not have the
// correct version of abigen installed , the above commands spin up docker
// containers. In my hands, total running time including compilation is about
// 13s. If you're modifying solidity code and testing against go code a lot, it
// might be worthwhile to generate the the wrappers using a static container
// with abigen and solc, which will complete much faster. E.g.
//
//   abigen -sol ../../../contracts/src/v0.6/VRFAll.sol -pkg vrf -out solidity_interfaces.go
//
// where VRFAll.sol simply contains `import "contract_path";` instructions for
// all the contracts you wish to target. This runs in about 0.25 seconds in my
// hands.
//
// If you're on linux, you can copy the correct version of solc out of the
// appropriate docker container. At least, the following works on ubuntu:
//
//   $ docker run --name solc ethereum/solc:0.6.2
//   $ sudo docker cp solc:/usr/bin/solc /usr/bin
//   $ docker rm solc
//
// If you need to point abigen at your solc executable, you can specify the path
// with the abigen --solc <path-to-executable> option.
