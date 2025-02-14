# Chainlink CCIP Smart Contracts

## Installation

```sh
# via pnpm
$ pnpm add @chainlink/contracts-ccip
# via npm
$ npm install @chainlink/contracts-ccip --save
```

### Directory Structure

```sh
@chainlink/contracts-ccip
├── src # Solidity contracts
│   └── v0.8
└── abi # ABI json output
    └── v0.8
```

### Usage

The solidity smart contracts themselves can be imported via the `src` directory of `@chainlink/contracts-ccip`:

```solidity
import '@chainlink/contracts-ccip/src/v0.8/ccip/applications/CCIPReceiver.sol';
```

### Changesets

We use [changesets](https://github.com/changesets/changesets) to manage versioning the contracts.

Every PR that modifies any configuration or code, should most likely accompanied by a changeset file.

To install `changesets`:

1. Install `pnpm` if it is not already installed - [docs](https://pnpm.io/installation).
2. Run `pnpm install`.

Either after or before you create a commit, run the `pnpm changeset:ccip` command in the `contracts` directory to create an accompanying changeset entry which will reflect on the CHANGELOG for the next release.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/) and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## License

The CCIP repo is licensed under the [BUSL-1.1](./src/v0.8/ccip/LICENSE.md) license, however, there are a few exceptions

- `src/v0.8/ccip/applications/*` is licensed under the [MIT](./src/v0.8/ccip/LICENSE-MIT.md) license
- `src/v0.8/ccip/interfaces/*` is licensed under the [MIT](./src/v0.8/ccip/LICENSE-MIT.md) license
- `src/v0.8/ccip/libraries/{Client.sol, Internal.sol}` is licensed under the [MIT](./src/v0.8/ccip/LICENSE-MIT.md) license
