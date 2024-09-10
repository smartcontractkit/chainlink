# Go bindings generation EVM smart contracts and CR / CW configuration

This is a POC to improve the UX of ChainReader and ChainWriter by providing a mechanism to create go bindings for smart contracts that use CR / CW as the underlying mechanism for blockchain connectivity.

This project generate a command line tool that can be plugged within the project build system, using go:generate, to generate client code for on-chain connectivity against solidty smart contracts, on top of ChainReader and ChainWriter abstractions.

Given a set of ABI this tool generate:
- ChainReaderConfig needed to create a ChainReader instance
- ChainWriterConfig needed to create a ChainWriter instance
- One go file per smart contract ABI with a struct that represent the smart contract client and functions to read and write to the smart contract

## Features supported

With probably plenty of bugs to be fixed, the POC supports:

- Strongly typed smart contract state reading. Events not yet supported.
- Strongly typed smart contracts method invocation that requires transactions. Payable methods allow to send ETH, Non-payable don't.
- Simple and complex types inputs and outputs.
- Basic chain reader and chain writer configuration but works for simple use cases.

Solidity primitive types to Go primitive types are not fully mapped. If you encounter an issue regarding type mapping you can update the code in internal/gen/evm/solity_types_mapping.go

## Troubleshooting and issue reporting

This is a POC so only used this tool at your own risk and do not use it yet for anything targeting production soon. We are still evaluating the POC and defining if we will productize this tool or not.
The easiest way to get help would be to push a branch with your projects and explain the issue with DM through slack to Pablo La Greca to get a fix. 

## installation

To install the command line tool:

```sh
git clone https://github.com/smartcontractkit/chainlink.git
cd chainlink
git checkout -b go-bindings-poc origin/go-bindings-poc
cd evm-chain-bindings
make install
```

## Usage

### Example project

There a example project located in examples/basic showing how to setup a project to add in the build process the tool execution using go:generate and Makefile

### Command line tool

```sh
❯ evm-chain-bindings -h  
Usage of evm-bindings:
       evm-bindings[flags] -contracts T [directory]
  -clean
        output folder for the generated code
  -contracts string
      comma-separated list of directories containing EVM smart contracts source and ABI files; must be set (default "contracts")
  -output string
      output folder for the generated code (default "generated/evm/bindings")
  -silent-if-no-contracts
      do not fails if there are not contracts to be processed
  -verbose
      generates debugging output
```


## Test cases

There's also a test case in [here](https://github.com/smartcontractkit/chainlink/blob/go-bindings-poc/core/services/relay/evm/go_bindings_test.go) that showcase a complete setup of this go bindings using ChainReader and ChainWriter.
