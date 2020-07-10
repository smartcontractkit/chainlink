package gethwrappers

// Make sure solidity compiler artifacts are up to date. Only output stdout on failure.
//go:generate ./generation/compile_contracts.sh

//go:generate go run ./generation/generate/wrap.go ../../../evm-contracts/abi/v0.6/FluxAggregator.json flux_aggregator_wrapper
//go:generate go run ./generation/generate/wrap.go ../../../evm-contracts/abi/v0.6/VRFTestHelper.json solidity_vrf_verifier_wrapper
//go:generate go run ./generation/generate/wrap.go ../../../evm-contracts/abi/v0.6/VRFCoordinator.json solidity_vrf_coordinator_interface
//go:generate go run ./generation/generate/wrap.go ../../../evm-contracts/abi/v0.6/VRFConsumer.json solidity_vrf_consumer_interface
//go:generate go run ./generation/generate/wrap.go ../../../evm-contracts/abi/v0.6/VRFRequestIDBaseTestHelper.json solidity_vrf_request_id
//go:generate go run ./generation/generate_link/wrap_link.go

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
//   abigen -sol ../../../evm-contracts/src/v0.6/VRFAll.sol -pkg vrf -out solidity_interfaces.go
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
