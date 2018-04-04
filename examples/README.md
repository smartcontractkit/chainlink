## Hello Chainlink

This is a repository of introductory [Chainlink](https://github.com/smartcontractkit/chainlink) integrations.
Please click on each folder for more information.

## Requirements

- Go 1.9+
- Node JS
- Docker

## Run Chainlink Development Environment

Each example requires the development environment, set up as follows:

1. Check out repo [smartcontractkit/chainlink](https://github.com/smartcontractkit/chainlink) and navigate to folder.
2. Run `./internal/bin/devnet` (or configure direnv `.envrc` to add `internal/bin` to your `PATH`)
3. Run truffle migrations:
  1. `cd solidity`
  2. `yarn install`
  3. `./node_modules/.bin/truffle migrate --network devnet`
4. Run `./internal/bin/cldev` in top level repo folder

Go to our [development wiki](https://github.com/smartcontractkit/chainlink/wiki/Development-Tips) to read more.
