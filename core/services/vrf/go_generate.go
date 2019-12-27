package vrf

// Make sure solidity compiler artifacts are up to date. Only output stdout on failure.
//go:generate sh -c "out=\"$(yarn workspace chainlinkv0.5 compile)\" || echo \"$out\""

//go:generate ./generation/generate.sh ../../../evm/v0.5/dist/artifacts/VRFTestHelper.json solidity_verifier_wrapper

// To reduce explicit dependencies, the above commands spin up docker
// containers. In my hands, total running time including compilation is about
// 8s. If you're modifying solidity code and testing against go code a lot, it
// might be worthwhile to generate the the wrappers using a static container
// with abigen and solc, which will complete much faster. E.g.
//
//   abigen -sol ../../../evm/v0.5/contracts/VRFAll.sol -pkg vrf -out solidity_interfaces.go
//
// where VRFAll.sol simply contains `import "contract_path";` instructions for
// all the contracts you wish to target. This runs in about 0.25 seconds in my
// hands.
//
// Here is a Dockerfile which can be used for that purpose.
//
//   # Build abigen docker image with the necessary solidity compilers
//   ARG SOLIDITY_VERSION
//   ARG GETH_VERSION
//   FROM ethereum/solc:${SOLIDITY_VERSION}
//   FROM ethereum/client-go:alltools-${GETH_VERSION}
//   RUN apk add bash
//   COPY --from=0 /usr/bin/solc /usr/bin/solc
//   USER 1000 # Or whatever your host user ID is.
//
// Build it with something like
//
//   docker build --build-arg SOLIDITY_VERSION=0.5.0 --build-arg GETH_VERSION=v1.9.9 -t image_name .
//
// Then run it with something like
//
//   docker run -it -v /path/to/solidity:/solidity /wrapper/target/path:/target image_name bash
//
// and run commands like this in there:
//
//   abigen -sol /solidity/VRFAll.sol -pkg vrf -out /target/tst.go
//
