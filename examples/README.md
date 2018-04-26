## Examples

This is a collection of introductory [Chainlink](https://github.com/smartcontractkit/chainlink) integrations.
Please click on each folder for more information.

## Requirements

- Go 1.10+
- Node JS
- Docker (for [DevNet](https://github.com/smartcontractkit/devnet))

## Run Chainlink Development Environment

Each example requires the development environment, set up as follows:

1. Run `docker pull smartcontract/devnet` to download the DevNet image.
2. Run `../internal/bin/devnet` (or configure direnv `.envrc` to add `chainlink/internal/bin` to your `PATH`)
3. Run truffle migrations for the Chainlink contracts:
  1. `cd ../solidity`
  2. `yarn install`
  3. `./node_modules/.bin/truffle migrate`
4. Run `./internal/bin/cldev` in top level repo folder

Go to our [development wiki](https://github.com/smartcontractkit/chainlink/wiki/Development-Tips) to read more.
