<br/>
<p align="center">
<a href="https://chain.link/" target="_blank" color="#0d2990">
  <img src="https://raw.githubusercontent.com/smartcontractkit/explorer/develop/styleguide/static/images/logo-contracts-blue.svg" width="225" alt="Chainlink logo">
</a>
</p>
<br/>

[![npm](https://img.shields.io/npm/v/@chainlink/contracts?style=flat-square)](https://www.npmjs.com/package/@chainlink/contracts)
[![GitHub license](https://img.shields.io/github/license/smartcontractkit/chainlink?style=flat-square)](https://github.com/smartcontractkit/chainlink/blob/master/LICENSE)
[![npm downloads](https://img.shields.io/npm/dt/@chainlink/contracts?style=flat-square)](https://npmjs.com/package/@chainlink/contracts)

[Chainlink's](https://chain.link/) solidity contracts and abstractions.

## Package Installation

```sh
# via Yarn
$ yarn add @chainlink/contracts

# via npm
$ npm install @chainlink/contracts --save
```

### Usage

If you wanted to consume the v0.6.x version of our `ChainlinkClient` smart contract, do the following:

```solidity
import "@chainlink/contracts/contracts/v0.6/ChainlinkClient.sol";
```

### Directory Structure

The contracts are compiled, packaged and published from a separate workflow. When installed via NPM, the files are structured like this:

```sh
contracts
├── abi # abi output from src/
│   ├── v0.4
│   ├── v0.5
│   └── v0.6
|   └── v0.7
├── ethers # ethers contract abstractions codegenned from abis
│   ├── v0.4
│   ├── v0.5
│   └── v0.6
|   └── v0.7
├── contracts # the contracts themselves, in .sol form
│   ├── v0.4
│   ├── v0.5
│   └── v0.6
|   └── v0.7
└── truffle  # truffle contract abstractions codegenned from abis
    ├── v0.4
    ├── v0.5
    └── v0.6
    └── v0.7
```

## Local Development

Note: Contracts in `src/v0.7/dev` are under active development and not yet stable.
Please use them for testing and development only.

```bash
# Clone Chainlink repository
$ git clone https://github.com/smartcontractkit/chainlink.git

# Continuing via Yarn
$ yarn install
$ yarn setup:contracts

# Continuing via npm
$ npm install
$ npm run setup:contracts
```

## Testing

After completing the above [Development](#Development) commands, run tests with:

```sh
# From this directory, `evm-contracts` via Yarn
$ yarn test

# via npm
$ npm run test

# From project root
$ yarn wsrun @chainlink/contracts test
```

## Contributing

Contributions are welcome! Please refer to
[Chainlink's contributing guidelines](./docs/CONTRIBUTING.md) for detailed
contribution information.

Thank you!

## License

[MIT](https://choosealicense.com/licenses/mit/)
