<br/>
<p align="center">
<a href="https://chain.link" target="_blank">
<img src="https://raw.githubusercontent.com/smartcontractkit/chainlink/develop/docs/logo-chainlink-blue.svg" width="225" alt="Chainlink logo">
</a>
</p>
<br/>

[![GitHub tag (latest SemVer)](https://img.shields.io/github/v/tag/smartcontractkit/chainlink?style=flat-square)](https://hub.docker.com/r/smartcontract/chainlink/tags)
[![GitHub license](https://img.shields.io/github/license/smartcontractkit/chainlink?style=flat-square)](https://github.com/smartcontractkit/chainlink/blob/master/LICENSE)
[![GitHub workflow changelog](https://img.shields.io/github/workflow/status/smartcontractkit/chainlink/Changelog?style=flat-square&label=github-actions)](https://github.com/smartcontractkit/chainlink/actions?query=workflow%3AChangelog)
[![CircleCI build](https://img.shields.io/circleci/build/github/smartcontractkit/chainlink/master?style=flat-square&label=circleci&logo=circleci)](https://circleci.com/gh/smartcontractkit/chainlink/tree/master)
[![GitHub contributors](https://img.shields.io/github/contributors-anon/smartcontractkit/chainlink?style=flat-square)](https://github.com/smartcontractkit/chainlink/graphs/contributors)
[![GitHub commit activity](https://img.shields.io/github/commit-activity/y/smartcontractkit/chainlink?style=flat-square)](https://github.com/smartcontractkit/chainlink/commits/master)
[![Official documentation](https://img.shields.io/static/v1?label=docs&message=latest&color=blue)](https://docs.chain.link/)

[Chainlink](https://chain.link/) expands the capabilities of smart contracts by enabling access to real-world data and off-chain computation while maintaining the security and reliability guarantees inherent to blockchain technology.

This repo contains the Chainlink core node, operator UI and contracts. The core node is the bundled binary available to be run by node operators participating in a [decentralized oracle network](https://link.smartcontract.com/whitepaper).
All major release versions have pre-built docker images available for download from the [Chainlink dockerhub](https://hub.docker.com/r/smartcontract/chainlink/tags).
If you are interested in contributing please see our [contribution guidelines](./docs/CONTRIBUTING.md).
If you are here to report a bug or request a feature, please [check currently open Issues](https://github.com/smartcontractkit/chainlink/issues).
For more information about how to get started with Chainlink, check our [official documentation](https://docs.chain.link/).
Resources for Solidity developers can be found in the [Chainlink Hardhat Box](https://github.com/smartcontractkit/hardhat-starter-kit).

## Community

Chainlink has an active and ever growing community. [Discord](https://discordapp.com/invite/aSK4zew)
is the primary communication channel used for day to day communication,
answering development questions, and aggregating Chainlink related content. Take
a look at the [community docs](./docs/COMMUNITY.md) for more information
regarding Chainlink social accounts, news, and networking.

## Build Chainlink

1. [Install Go 1.18](https://golang.org/doc/install), and add your GOPATH's [bin directory to your PATH](https://golang.org/doc/code.html#GOPATH)
   - Example Path for macOS `export PATH=$GOPATH/bin:$PATH` & `export GOPATH=/Users/$USER/go`
2. Install [NodeJS](https://nodejs.org/en/download/package-manager/) & [Yarn](https://yarnpkg.com/lang/en/docs/install/). See the current version in `package.json` at the root of this repo under the `engines.node` key.
   - It might be easier long term to use [nvm](https://nodejs.org/en/download/package-manager/#nvm) to switch between node versions for different projects. For example, assuming $NODE_VERSION was set to a valid version of NodeJS, you could run: `nvm install $NODE_VERSION && nvm use $NODE_VERSION`
3. Install [Postgres (>= 11.x)](https://wiki.postgresql.org/wiki/Detailed_installation_guides).
   - You should [configure Postgres](https://www.postgresql.org/docs/12/ssl-tcp.html) to use SSL connection (or for testing you can set `?sslmode=disable` in your Postgres query string).
4. Ensure you have Python 3 installed (this is required by [solc-select](https://github.com/crytic/solc-select) which is needed to compile solidity contracts)
5. Download Chainlink: `git clone https://github.com/smartcontractkit/chainlink && cd chainlink`
6. Build and install Chainlink: `make install`
   - If you got any errors regarding locked yarn package, try running `yarn install` before this step
   - If `yarn install` throws a network connection error, try increasing the network timeout by running `yarn install --network-timeout 150000` before this step
7. Run the node: `chainlink help`

For the latest information on setting up a development environment, see the [Development Setup Guide](https://github.com/smartcontractkit/chainlink/wiki/Development-Setup-Guide).

### Mac M1/ARM64 [EXPERIMENTAL]

Chainlink can be experimentally compiled with ARM64 as the target arch. You may run into errors with cosmwasm:

```
# github.com/CosmWasm/wasmvm/api
ld: warning: ignoring file ../../../.asdf/installs/golang/1.18/packages/pkg/mod/github.com/!cosm!wasm/wasmvm@v0.16.3/api/libwasmvm.dylib, building for macOS-arm64 but attempting to link with file built for macOS-x86_64
Undefined symbols for architecture arm64:# github.com/CosmWasm/wasmvm/api
ld: warning: ignoring file ../../../.asdf/installs/golang/1.18/packages/pkg/mod/github.com/!cosm!wasm/wasmvm@v0.16.3/api/libwasmvm.dylib, building for macOS-arm64 but attempting to link with file built for macOS-x86_64
Undefined symbols for architecture arm64:
```

In this case, try the following steps:

1. `git clone git@github.com:mandrean/terra-core.git`
2. `cd terra-core; git checkout feat/multiarch`
3. `make install; cd ..`
4. `go work init /path/to/chainlink`
5. `go work use /path/to/terra-core`

### Ethereum Node Requirements

In order to run the Chainlink node you must have access to a running Ethereum node with an open websocket connection.
Any Ethereum based network will work once you've [configured](https://github.com/smartcontractkit/chainlink#configure) the chain ID.
Ethereum node versions currently tested and supported:

- [Parity/Openethereum](https://github.com/openethereum/openethereum)
- [Geth](https://github.com/ethereum/go-ethereum/releases)
- [Nethermind](https://github.com/NethermindEth/nethermind)

We cannot recommend specific version numbers for ethereum nodes since the software is being continually updated, but you should usually try to run the latest version available.

## Running a local Chainlink node

**NOTE**: By default, chainlink will run in TLS mode. For local development you can disable this by setting the following env vars:

```
CHAINLINK_DEV=true
CHAINLINK_TLS_PORT=0
SECURE_COOKIES=false
```

Alternatively, you can generate self signed certificates using `tools/bin/self-signed-certs` or [manually](https://github.com/smartcontractkit/chainlink/wiki/Creating-Self-Signed-Certificates).

To start your Chainlink node, simply run:

```bash
chainlink node start
```

By default this will start on port 6688. You should be able to access the UI at [http://localhost:6688/](http://localhost:6688/).

Chainlink provides a remote CLI client as well as a UI. Once your node has started, you can open a new terminal window to use the CLI. You will need to log in to authorize the client first:

```bash
chainlink admin login
```

(You can also set `ADMIN_CREDENTIALS_FILE=/path/to/credentials/file` in future if you like, to avoid having to login again).

Now you can view your current jobs with:

```bash
chainlink jobs list
```

To find out more about the Chainlink CLI, you can always run `chainlink help`.

Check out the [doc](https://docs.chain.link/) pages on [Jobs](https://docs.chain.link/docs/jobs/) to learn more about how to create Jobs.

### Configuration

Node configuration is managed by a combination of environment variables and direct setting via API/UI/CLI.

Check the [official documentation](https://docs.chain.link/docs/configuration-variables) for more information on how to configure your node.

### External Adapters

External adapters are what make Chainlink easily extensible, providing simple integration of custom computations and specialized APIs. A Chainlink node communicates with external adapters via a simple REST API.

For more information on creating and using external adapters, please see our [external adapters page](https://docs.chain.link/docs/external-adapters).

## Development

### Running tests

1. [Install Yarn](https://yarnpkg.com/lang/en/docs/install)

2. Install [gencodec](https://github.com/fjl/gencodec) and [jq](https://stedolan.github.io/jq/download/) to be able to run `go generate ./...` and `make abigen`

3. Install mockery

`make mockery`

Using the `make` command will install the correct version.

4. Build contracts:

```bash
yarn
yarn setup:contracts
```

4. Generate and compile static assets:

```bash
go generate ./...
```

5. Prepare your development environment:

```bash
export DATABASE_URL=postgresql://127.0.0.1:5432/chainlink_test?sslmode=disable
```

Note: Other environment variables should not be set for all tests to pass

6.  Drop/Create test database and run migrations:

```
go run ./core/main.go local db preparetest
```

If you do end up modifying the migrations for the database, you will need to rerun

7. Run tests:

```bash
go test ./...
```

#### Notes

- The `parallel` flag can be used to limit CPU usage, for running tests in the background (`-parallel=4`) - the default is `GOMAXPROCS`
- The `p` flag can be used to limit the number of _packages_ tested concurrently, if they are interferring with one another (`-p=1`)
- The `-short` flag skips tests which depend on the database, for quickly spot checking simpler tests in around one minute (you may still need a phony env var to pass some validation: `DATABASE_URL=_test`)

#### Race Detector

As of Go 1.1, the runtime includes a data race detector, enabled with the `-race` flag. This is used in CI via the 
`tools/bin/go_core_race_tests` script. If the action detects a race, the artifact on the summary page will include 
`race.*` files with detailed stack traces. 

> _**It will not issue false positives, so take its warnings seriously.**_

For local, targeted race detection, you can run:
```bash
GORACE="log_path=$PWD/race" go test -race ./core/path/to/pkg -count 10
GORACE="log_path=$PWD/race" go test -race ./core/path/to/pkg -count 100 -run TestFooBar/sub_test 
```

https://go.dev/doc/articles/race_detector

#### Fuzz tests

As of Go 1.18, fuzz tests `func FuzzXXX(*testing.F)` are included as part of the normal test suite, so existing cases are executed with `go test`.

Additionally, you can run active fuzzing to search for new cases:
```bash
go test ./pkg/path -run=XXX -fuzz=FuzzTestName
```

https://go.dev/doc/fuzz/

### Solidity

Inside the `contracts/` directory:
1. Install dependencies:

```bash
yarn
```

2. Run tests:

```bash
yarn test
```

### Code Generation

Go generate is used to generate mocks in this project. Mocks are generated with [mockery](https://github.com/vektra/mockery) and live in core/internal/mocks.

### Nix Flake

A [flake](https://nixos.wiki/wiki/Flakes) is provided for use with the [Nix
package manager](https://nixos.org/). It defines a declarative, reproducible
development environment.

To use it:

1. [Nix has to be installed with flake support](https://nixos.wiki/wiki/Flakes#Installing_flakes).
2. Run `nix develop`. You will be put in shell containing all the dependencies.
   Alternatively, a `direnv` integration exists to automatically change the
   environment when `cd`-ing into the folder.
3. Create a local postgres database:

```
cd $PGDATA/
initdb
pg_ctl -l $PGDATA/postgres.log -o "--unix_socket_directories='$PWD'" start
createdb chainlink_test -h localhost
createuser --superuser --no-password chainlink -h localhost
```

4. Start postgres, `pg_ctl -l $PGDATA/postgres.log -o "--unix_socket_directories='$PWD'" start`

Now you can run tests or compile code as usual.

### Tips

For more tips on how to build and test Chainlink, see our [development tips page](https://github.com/smartcontractkit/chainlink/wiki/Development-Tips).

### Contributing

Chainlink's source code is [licensed under the MIT License](./LICENSE), and contributions are welcome.

Please check out our [contributing guidelines](./docs/CONTRIBUTING.md) for more details.

Thank you!

## License

[MIT](https://choosealicense.com/licenses/mit/)
