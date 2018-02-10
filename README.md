# Chainlink Go Implementation

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
```

Direnv can be installed by running

```bash
$ go get -u github.com/direnv/direnv
```

Environment variables that can be set in .envrc, along with default values that get used if no corresponding enviornment variable is found:

    LOG_LEVEL                Default: info
    ROOT                     Default: ~/.chainlink
    USERNAME                 Default: chainlink
    PASSWORD                 Default: twochains
    ETH_URL                  Default: http://localhost:8545
    ETH_CHAIN_ID             Default: 0
    POLLING_SCHEDULE         Default: */15 * * * * *
    CLIENT_NODE_URL          Default: http://localhost:8080
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

### Adding External Adapters

Post to `/v2/jobs`:

```shell
curl -u chainlink:twochains -X POST -H 'Content-Type: application/json' -d '{"name":"randomNumber","url":"https://example.com/randomNumber"}' http://localhost:8080/v2/jobs
```

`"name"` should be unique to the local node, and `"url"` should be the URL of your external adapter, whether local or on a separate machine.

Output should return a unique id:

```shell
{"id":"65c60bd0267941369e2e92e73c5319d6"}
```