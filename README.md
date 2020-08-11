<br/>
<p align="center">
<a href="https://chain.link/" target="_blank" color="#0d2990">
  <img src="./styleguide/static/logo-blue.svg" width="225" alt="Chainlink logo">
</a>
</p>
<br/>

[![Join the chat at https://discordapp.com/invite/aSK4zew](https://img.shields.io/discord/592041321326182401.svg?logoColor=white)](https://discordapp.com/invite/aSK4zew)
[![CircleCI](https://circleci.com/gh/smartcontractkit/chainlink.svg?style=shield)](https://circleci.com/gh/smartcontractkit/chainlink)
[![Go Report Card](https://goreportcard.com/badge/github.com/smartcontractkit/chainlink)](https://goreportcard.com/report/github.com/smartcontractkit/chainlink)
[![GoDoc](https://godoc.org/github.com/smartcontractkit/chainlink?status.svg)](https://godoc.org/github.com/smartcontractkit/chainlink)
[![Maintainability](https://api.codeclimate.com/v1/badges/273722bb9f6f22d799bd/maintainability)](https://codeclimate.com/github/smartcontractkit/chainlink/maintainability)
[![Test Coverage](https://api.codeclimate.com/v1/badges/273722bb9f6f22d799bd/test_coverage)](https://codeclimate.com/github/smartcontractkit/chainlink/test_coverage)



[Chainlink](https://chain.link/) is middleware to simplify communication with blockchains.
Here you'll find the Chainlink Golang node, currently in alpha.
This initial implementation is intended for use and review by developers,
and will go on to form the basis for Chainlink's [decentralized oracle network](https://link.smartcontract.com/whitepaper).
Further development of the Chainlink Node and Chainlink Network will happen here,
if you are interested in contributing please see our [contribution guidelines](./docs/CONTRIBUTING.md).
The current node supports:

- easy connectivity of on-chain contracts to any off-chain computation or API
- multiple methods for scheduling both on-chain and off-chain computation for a user's smart contract
- automatic gas price bumping to prevent stuck transactions, assuring your data is delivered in a timely manner
- push notification of smart contract state changes to off-chain systems, by tracking Ethereum logs
- translation of various off-chain data types into EVM consumable types and transactions
- easy to implement smart contract libraries for connecting smart contracts directly to their preferred oracles
- easy to install node, which runs natively across operating systems, blazingly fast, and with a low memory footprint

Examples of how to utilize and integrate Chainlinks can be found in the [Chainlink Truffle Box](https://github.com/smartcontractkit/box).

## Install

1. [Install Go 1.14](https://golang.org/doc/install#install), and add your GOPATH's [bin directory to your PATH](https://golang.org/doc/code.html#GOPATH)
2. Install [NodeJS](https://nodejs.org/en/download/package-manager/) & [Yarn](https://yarnpkg.com/lang/en/docs/install/)
3. Install [Postgres (>= 9.6)](https://wiki.postgresql.org/wiki/Detailed_installation_guides).
4. Download Chainlink: `git clone https://github.com/smartcontractkit/chainlink && cd chainlink`
5. Build and install Chainlink: `make install`
6. Run the node: `chainlink help`

### Ethereum Node Requirements

In order to run the Chainlink node you must have access to a running Ethereum node with an open websocket connection.
Any Ethereum based network will work once you've [configured](https://github.com/smartcontractkit/chainlink#configure) the chain ID.
Ethereum node versions currently tested and supported:

- [Parity 1.11+](https://github.com/paritytech/parity-ethereum/releases) (due to a [fix with pubsub](https://github.com/paritytech/parity/issues/6590).)
- [Geth 1.8+](https://github.com/ethereum/go-ethereum/releases)

## Run

**NOTE**: By default, chainlink will run in TLS mode. For local development you can either disable this by setting CHAINLINK_DEV to true, or generate self signed certificates using `tools/bin/self-signed-certs` or [manually](https://github.com/smartcontractkit/chainlink/wiki/Creating-Self-Signed-Certificates).

To start your Chainlink node, simply run:

```bash
chainlink node start
```

By default this will start on port 6688, where it exposes a [REST API](https://github.com/smartcontractkit/chainlink/wiki/REST-API).

Once your node has started, you can view your current jobs with:

```bash
chainlink jobs list
```

View details of a specific job with:

```bash
chainlink jobs show "$JOB_ID"
```

To find out more about the Chainlink CLI, you can always run `chainlink help`.

Check out the [docs'](https://docs.chain.link/) pages on [Adapters](https://docs.chain.link/docs/adapters) and [Initiators](https://docs.chain.link/docs/initiators) to learn more about how to create Jobs and Runs.

## Configure

You can configure your node's behavior by setting environment variables which can be, along with default values that get used if no corresponding environment variable is found. The latest information on configuration variables are available in [the docs](https://docs.chain.link/docs/configuration-variables).

## Project Directory

This project contains several sub-projects, some with their own documentation.

- [core](./core) - the core Chainlink node
- [@chainlink/belt](./belt) - tools for performing commands on Chainlink smart contracts
- [@chainlink/contracts](./evm-contracts) - smart contracts
- [@chainlink/test-helpers](./evm-test-helpers) - smart contract-related resources
- [explorer](./explorer) - [Chainlink Explorer](https://explorer.chain.link/)
- [integration/forks](./integration/forks) - integration test for [ommers](https://ethereum.stackexchange.com/a/46/19503) and [re-orgs](https://en.bitcoin.it/wiki/Chain_Reorganization)
- [sgx](./core/sgx) - experimental, optional module that can be loaded into Chainlink to do processing within an [SGX](https://software.intel.com/en-us/sgx) enclave
- [styleguide](./styleguide) - Chainlink style guide
- [tools](./tools) - Chainlink tools

## External Adapters

External adapters are what make Chainlink easily extensible, providing simple integration of custom computations and specialized APIs.
A Chainlink node communicates with external adapters via a simple REST API.

For more information on creating and using external adapters, please see our [external adapters page](https://docs.chain.link/docs/external-adapters).

## Development Setup

For the latest information on setting up a development environment, see the [guide here](https://github.com/smartcontractkit/chainlink/wiki/Development-Setup-Guide).

### Build your current version

```bash
go build -o chainlink ./core/
```

- Run the binary:

```bash
./chainlink
```

### Test Core

1. [Install Yarn](https://yarnpkg.com/lang/en/docs/install)

2. Install [gencodec](https://github.com/fjl/gencodec), [mockery version 1.0.0](https://github.com/vektra/mockery/releases/tag/v1.0.0), and [jq](https://stedolan.github.io/jq/download/) to be able to run `go generate ./...` and `make abigen`

3. Build contracts:

```bash
yarn
yarn setup:contracts
```

4. Generate and compile static assets:

```bash
go generate ./...
go run ./packr/main.go ./core/eth/
```

5. Prepare your development environment:

```bash
export DATABASE_URL=postgresql://127.0.0.1:5432/chainlink_test?sslmode=disable
export CHAINLINK_DEV=true # I prefer to use direnv and skip this
```

6.  Drop/Create test database and run migrations:
```
go run ./core/main.go local db preparetest
```

If you do end up modifying the migrations for the database, you will need to rerun

7. Run tests:

```bash
go test -parallel=1 ./...
```


### Solidity Development

1. [Install Yarn](https://yarnpkg.com/lang/en/docs/install)
2. Install the dependencies:

```bash
cd evm
yarn install
```

3. Run tests:

```bash
yarn run test-sol
```

### Use of Go Generate

Go generate is used to generate mocks in this project. Mocks are generate with [mockery](https://github.com/vektra/mockery) and live in core/internal/mocks.

### Development Tips

For more tips on how to build and test Chainlink, see our [development tips page](https://github.com/smartcontractkit/chainlink/wiki/Development-Tips).

## Contributing

Chainlink's source code is [licensed under the MIT License](./LICENSE), and contributions are welcome.

Please check out our [contributing guidelines](./docs/CONTRIBUTING.md) for more details.

Thank you!

## License

[MIT](https://choosealicense.com/licenses/mit/)

