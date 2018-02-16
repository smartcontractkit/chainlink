# Chainlink [![Travis-CI](https://travis-ci.com/smartcontractkit/chainlink.svg?token=55tBxbZKUxXXNcTx1P6u&branch=master)](https://travis-ci.com/smartcontractkit/chainlink) [![Maintainability](https://api.codeclimate.com/v1/badges/f45294bf9c01d8055451/maintainability)](https://codeclimate.com/repos/5a864b469b27633ee4002e4b/maintainability) [![CodeCov](https://codecov.io/gh/smartcontractkit/chainlink/branch/master/graph/badge.svg?token=1LacYNWKro)](https://codecov.io/gh/smartcontractkit/chainlink) [![Join the chat at https://gitter.im/smartcontractkit-chainlink/Lobby](https://badges.gitter.im/smartcontractkit-chainlink/Lobby.svg)](https://gitter.im/smartcontractkit-chainlink/Lobby)

## Ethereum Node Requirements

- Parity 1.9+ due to a [fix with pubsub](https://github.com/paritytech/parity/issues/6590).
- Geth 1.7+

## Developer Setup

### Install [Go 1.9+](https://golang.org/dl/)

Create the Go workspace (`~/go` given as an example)

```bash
$ mkdir ~/go && cd ~/go
```

Set environment variables

```bash
$ export GOPATH=$(pwd)
$ export PATH=$PATH:$GOPATH/bin
```

### Create Project Directories

```bash
$ cd $GOPATH
$ mkdir -p src/github.com/smartcontractkit
$ cd src/github.com/smartcontractkit
```

### Clone the repo

```bash
$ git clone https://github.com/smartcontractkit/chainlink.git
$ cd chainlink
```

### Get and run dep

Linux

```bash
$ go get -u github.com/golang/dep/cmd/dep
$ dep ensure
```

Mac

```bash
$ brew install dep
$ dep ensure
```

### Build the project

```bash
$ go build -o chainlink
```

Run the binary

```bash
$ ./chainlink
```

### Testing

```bash
$ cd $GOPATH/src/github.com/smartcontractkit/chainlink
$ go test ./...
```

### Direnv

We use [direnv](https://github.com/direnv/direnv/) to set up PATH and aliases 
for a friendlier developer experience. Here is an example `.envrc` that we use:

```bash
$ cat .envrc
PATH_add tmp
PATH_add solidity/node_modules/.bin
PATH_add internal/bin
```

Direnv can be installed by running

```bash
$ go get -u github.com/direnv/direnv
```

Environment variables that can be set in .envrc, along with default values that get used if no corresponding enviornment variable is found:

    LOG_LEVEL                Default: info
    ROOT                     Default: ~/.chainlink
    PORT                     Default: 6688
    USERNAME                 Default: chainlink
    PASSWORD                 Default: twochains
    ETH_URL                  Default: ws://localhost:8546
    ETH_CHAIN_ID             Default: 0
    POLLING_SCHEDULE         Default: */15 * * * * *
    CLIENT_NODE_URL          Default: http://localhost:6688
    ETH_MIN_CONFIRMATIONS    Default: 12
    ETH_GAS_BUMP_THRESHOLD   Default: 12
    ETH_GAS_BUMP_WEI         Default: 5,000,000,000
    ETH_GAS_PRICE_DEFAULT    Default: 20,000,000,000

### Solidity Development setup

Before proceeding, make sure you have installed [yarn](https://yarnpkg.com/lang/en/docs/install)

```bash
$ cd solidity
$ yarn install
$ truffle test
```

### External Adapters

External adapters are added to the Chainlink node first by adding a bridge type. Bridge types define the name and URL of your external adapter. When a task type is received that is not one of the core adapters, the node will search for a bridge type with that name, creating a bridge to your external adapter.

bridge_type.json

```JSON
{ "name": "randomNumber", "url": "http://localhost:3000/randomNumber" }
```

job.json

```JSON
{
  "initiators": [{
    "type": "runLog"
  }],
  "tasks": [ {"type": "randomNumber"} ]
}
```

#### Adding an External Adapter

POST to `/v2/bridge_types`:

```shell
curl -u chainlink:twochains -X POST -H 'Content-Type: application/json' -d '{"name":"randomNumber","url":"http://localhost:3000/randomNumber"}' http://localhost:6688/v2/bridge_types
```

`"name"` should be unique to the local node, and `"url"` should be the URL of your external adapter, whether local or on a separate machine.

Output should return the JSON given:

```shell
{"name":"randomnumber","url":"http://localhost:3000/randomNumber"}
```

And the node will log the following:

```shell
{"level":"info","ts":1518531822.179224,"caller":"web/router.go:50","msg":"Web request","method":"POST","status":200,"path":"/v2/bridge_types","query":"","body":"{\"name\":\"randomNumber\",\"url\":\"http://localhost:3000/randomNumber\"}","clientIP":"127.0.0.1","comment":"","servedAt":"2018/02/13 - 14:23:42","latency":"1.623398ms"}
```
