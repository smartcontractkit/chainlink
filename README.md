# Chainlink Go Implementation

## Developer Setup

### Install [Go](https://golang.org/dl/)

Create the Go workspace (`~/go` given as an example)

```bash
$ mkdir ~/go && cd ~/go
```

Set environment variables

```bash
$ export GOPATH=`pwd`
$ export PATH=$PATH:$(go env GOPATH)/bin
```

### Create Project Directories

```bash
$ cd $(go env GOPATH)
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

### Direnv

We use [direnv](https://github.com/direnv/direnv/) to set up PATH and aliases 
for a friendlier developer experience. Here is an example `.envrc` that we use:

```bash
$ PATH_add tmp
$ PATH_add solidity/node_modules/.bin
```

Direnv can be installed by running

```bash
$ go get -u github.com/direnv/direnv
```

Environment variables, along with their default values that can be used:

    ROOT                     Default: ~/.chainlink
    USERNAME                 Default: chainlink
    PASSWORD                 Default: twochains
    ETH_URL                  Default: http://localhost:8545
    ETH_CHAIN_ID             Default: 0
    POLLING_SCHEDULE         Default: */15 * * * * *
    CLIENT_NODE_URL          Default: http://localhost:8080
    ETH_MIN_CONFIRMATIONS    Default: 12
    ETH_GAS_BUMP_THRESHOLD   Default: 12
    ETH_GAS_BUMP_GWEI        Default: 5000000000
    ETH_GAS_PRICE_DEFAULT    Default: 20000000000

### Solidity Development setup

```bash
$ cd solidity
$ yarn install
$ truffle test
```