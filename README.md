# ChainLink [![Maintainability](https://api.codeclimate.com/v1/badges/f45294bf9c01d8055451/maintainability)](https://codeclimate.com/repos/5a864b469b27633ee4002e4b/maintainability) [![Travis-CI](https://travis-ci.com/smartcontractkit/chainlink.svg?token=55tBxbZKUxXXNcTx1P6u&branch=master)](https://travis-ci.com/smartcontractkit/chainlink) [![Join the chat at https://gitter.im/smartcontractkit-chainlink/Lobby](https://badges.gitter.im/smartcontractkit-chainlink/Lobby.svg)](https://gitter.im/smartcontractkit-chainlink/Lobby) [![GoDoc](https://godoc.org/github.com/smartcontractkit/chainlink?status.svg)](https://godoc.org/github.com/smartcontractkit/chainlink)

ChainLink is middleware to simplify communication with blockchains. Here you'll find the ChainLink Golang node, currently in alpha. This initial implementation is intended for use and review by developers, and will go on to form the basis for ChainLink's [decentralized oracle network](https://link.smartcontract.com/whitepaper). Further development of the ChainLink Node and ChainLink Network will happen here, if you are interested in contributing please see our [contribution guidelines](https://github.com/smartcontractkit/chainlink/blob/master/CONTRIBUTING.md). The current node supports:
- easy connectivity of on-chain contracts to any off-chain computation or API
- multiple methods for scheduling both on-chain and off-chain computation for a user's smart contract
- automatic gas price bumping to prevent stuck transactions, assuring your data is delivered in a timely manner
- push notification of smart contract state changes to off-chain systems, by tracking Ethereum logs
- translation of various off-chain data types into EVM consumable types and transactions
- easy to implement smart contract libraries for connecting smart contracts directly to their preferred oracles
- easy to install node, which runs natively across operating systems, blazingly fast, and with a low memory footprint

Examples of how to utilize and integrate Chainlinks can be found in our [Hello Chainlink](https://github.com/smartcontractkit/hello_chainlink) repository. 

## Install

1. [Install Go 1.9+](https://golang.org/doc/install#install), and add your GOPATH's [bin directory to your PATH](https://golang.org/doc/code.html#GOPATH)
2. Install ChainLink: `$ go get -u github.com/smartcontractkit/chainlink`
3. Run the node: `$ chainlink help`

### Ethereum Node Requirements

In order to run the ChainLink node you must have access to a running Ethereum node with an open websocket connection.
Any Ethereum based network will work once you've [configured](https://github.com/smartcontractkit/chainlink#configure) the chain ID.
Ethereum node versions currently tested and supported:

- Parity 1.9+ (due to a [fix with pubsub](https://github.com/paritytech/parity/issues/6590).)
- Geth 1.7+

## Run

To start your ChainLink node, simply run:
```bash
$ chainlink node
```
By default this will start on port 6688, where it exposes a [REST API](https://github.com/smartcontractkit/chainlink/wiki/REST-API).

Once your node has started, you can view your current jobs with:
```bash
$ chainlink jobs
````
View details of a specific job with:
```bash
$ chainlink show $JOB_ID
```

To find out more about the ChainLink CLI, you can always run `chainlink help`.

Check out the [wiki](https://github.com/smartcontractkit/chainlink/wiki)'s pages on [Adapters](https://github.com/smartcontractkit/chainlink/wiki/Adapters) and [Initiators](https://github.com/smartcontractkit/chainlink/wiki/Initiators) to learn more about how to create Jobs and Runs.

## Configure

You can configure your node's behavior by setting environment variables which can be, along with default values that get used if no corresponding enviornment variable is found:

    LOG_LEVEL                Default: info
    ROOT                     Default: ~/.chainlink
    PORT                     Default: 6688
    USERNAME                 Default: chainlink
    PASSWORD                 Default: twochains
    ETH_URL                  Default: ws://localhost:8546
    ETH_CHAIN_ID             Default: 0
    ETH_GAS_BUMP_THRESHOLD   Default: 12
    ETH_MIN_CONFIRMATIONS    Default: 12
    ETH_GAS_BUMP_WEI         Default: 5000000000  (5 gwei)
    ETH_GAS_PRICE_DEFAULT    Default: 20000000000 (20 gwei)

When running the CLI to talk to a ChainLink node on another machine, you can change the following environment variables:

    CLIENT_NODE_URL          Default: http://localhost:6688
    USERNAME                 Default: chainlink
    PASSWORD                 Default: twochains

## External Adapters

External adapters are what make ChainLink easily extensible, providing simple integration of custom computations and specialized APIs.
A ChainLink node communicates with external adapters via a simple REST API.

For more information on creating and using external adapters, please see our [external adapters page](https://github.com/smartcontractkit/chainlink/wiki/External-Adapters).


## Development Setup


- [Install Go 1.9+](https://golang.org/doc/install#install)
- Set up a Go workspace(`~/go` given as an example directory) and add go binaries to your path:
```bash
$ mkdir ~/go && cd ~/go
$ export GOPATH=$(pwd)
$ export PATH=$PATH:$GOPATH/bin
```

- [Install `dep`](https://github.com/golang/dep#installation):
```bash
$ go get -u github.com/golang/dep/cmd/dep
```

- Clone the repo:
```bash
$ go get github.com/smartcontractkit/chainlink
$ cd $GOPATH/src/github.com/smartcontractkit/chainlink
```

- Run:
```bash
$ go run main.go
```

### Build your current version

```bash
$ cd $GOPATH/src/github.com/smartcontractkit/chainlink
$ go build
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
$ cd $GOPATH/src/github.com/smartcontractkit/chainlink/solidity
$ yarn install
```
3. Run tests:
```bash
$ truffle test
```
### Development Tips

For more tips on how to build and test ChainLink, see our [development tips page](https://github.com/smartcontractkit/chainlink/wiki/Development-Tips).

## Contributing

ChainLink's source code is [licensed under the MIT License](https://github.com/smartcontractkit/chainlink/blob/master/LICENSE), and contributions are welcome.

Please check out our [contributing guidelines](https://github.com/smartcontractkit/chainlink/blob/master/CONTRIBUTING.md) for more details.

Thank you!
