# Chainlink

[![Join the chat at https://gitter.im/smartcontractkit-chainlink/Lobby](https://badges.gitter.im/smartcontractkit-chainlink/Lobby.svg)](https://gitter.im/smartcontractkit-chainlink/Lobby)
[![CircleCI](https://circleci.com/gh/smartcontractkit/chainlink.svg?style=shield)](https://circleci.com/gh/smartcontractkit/chainlink)
[![Go Report Card](https://goreportcard.com/badge/github.com/smartcontractkit/chainlink)](https://goreportcard.com/report/github.com/smartcontractkit/chainlink)
[![GoDoc](https://godoc.org/github.com/smartcontractkit/chainlink?status.svg)](https://godoc.org/github.com/smartcontractkit/chainlink)
[![Maintainability](https://api.codeclimate.com/v1/badges/273722bb9f6f22d799bd/maintainability)](https://codeclimate.com/github/smartcontractkit/chainlink/maintainability)
[![Test Coverage](https://api.codeclimate.com/v1/badges/273722bb9f6f22d799bd/test_coverage)](https://codeclimate.com/github/smartcontractkit/chainlink/test_coverage)

Chainlink is middleware to simplify communication with blockchains.
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

Examples of how to utilize and integrate Chainlinks can be found in the [examples](./examples) directory.

## Install

1. [Install Go 1.12+](https://golang.org/doc/install#install), and add your GOPATH's [bin directory to your PATH](https://golang.org/doc/code.html#GOPATH)
2. Install [NodeJS](https://nodejs.org/en/download/package-manager/) & [Yarn](https://yarnpkg.com/lang/en/docs/install/)
3. Download Chainlink: `git clone github.com/smartcontractkit/chainlink && cd chainlink`
4. Build and install Chainlink: `make install`
5. Run the node: `chainlink help`

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
$ chainlink local node
```
By default this will start on port 6688, where it exposes a [REST API](https://github.com/smartcontractkit/chainlink/wiki/REST-API).

Once your node has started, you can view your current jobs with:
```bash
$ chainlink jobspecs
````
View details of a specific job with:
```bash
$ chainlink show $JOB_ID
```

To find out more about the Chainlink CLI, you can always run `chainlink help`.

Check out the [wiki](https://github.com/smartcontractkit/chainlink/wiki)'s pages on [Adapters](https://github.com/smartcontractkit/chainlink/wiki/Adapters) and [Initiators](https://github.com/smartcontractkit/chainlink/wiki/Initiators) to learn more about how to create Jobs and Runs.

## Configure

You can configure your node's behavior by setting environment variables which can be, along with default values that get used if no corresponding environment variable is found. The latest information on configuration variables are available in [the docs](https://docs.chain.link/docs/configuration-variables).

## External Adapters

External adapters are what make Chainlink easily extensible, providing simple integration of custom computations and specialized APIs.
A Chainlink node communicates with external adapters via a simple REST API.

For more information on creating and using external adapters, please see our [external adapters page](https://github.com/smartcontractkit/chainlink/wiki/External-Adapters).


## Development Setup

For the latest information on setting up a development environment, see the [guide here](https://github.com/smartcontractkit/chainlink/wiki/Development-Setup-Guide).

### Build your current version

```bash
$ go build -o chainlink
```

- Run the binary:
```bash
$ ./chainlink
```

### Test

```bash
$ cd $GOPATH/src/github.com/smartcontractkit/chainlink
$ go test ./...
```

### Solidity Development

1. [Install Yarn](https://yarnpkg.com/lang/en/docs/install)
2. Install the dependencies:
```bash
$ cd $GOPATH/src/github.com/smartcontractkit/chainlink/evm
$ yarn install
```
3. Run tests:
```bash
$ yarn run test-sol
```
### Development Tips

For more tips on how to build and test Chainlink, see our [development tips page](https://github.com/smartcontractkit/chainlink/wiki/Development-Tips).

## Contributing

Chainlink's source code is [licensed under the MIT License](https://github.com/smartcontractkit/chainlink/blob/master/LICENSE), and contributions are welcome.

Please check out our [contributing guidelines](./docs/CONTRIBUTING.md) for more details.

Thank you!
