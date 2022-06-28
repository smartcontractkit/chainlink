# Chainlink Smart Contracts

## Installation

```sh
# via Yarn
$ yarn add @chainlink/contracts
# via npm
$ npm install @chainlink/contracts --save
```

### Directory Structure

```sh
@chainlink/contracts
├── src # Solidity contracts
│   ├── v0.4
│   ├── v0.5
│   ├── v0.6
│   ├── v0.7
│   └── v0.8
└── abi # ABI json output
    ├── v0.4
    ├── v0.5
    ├── v0.6
    ├── v0.7
    └── v0.8
```

### Usage

The solidity smart contracts themselves can be imported via the `src` directory of `@chainlink/contracts`:

```solidity
import '@chainlink/contracts/src/v0.8/KeeperCompatibleInterface.sol';

```

## Local Development

Note: Contracts in `dev/` directories are under active development and are likely unaudited. Please refrain from using these in production applications.

```bash
# Clone Chainlink repository
$ git clone https://github.com/smartcontractkit/chainlink.git
# Continuing via Yarn
$ cd contracts/
$ yarn
$ yarn test
```

## Contributing
Check out our [Github repo](https://github.com/smartcontractkit/chainlink).

Contributions are welcome! Please refer to the `/docs` folder in our [repo](https://github.com/smartcontractkit/chainlink/tree/develop/docs) where you will find `CONTRIBUTING.md` which has detailed information on how you can contribute.

Thank you!

## License

[MIT](https://choosealicense.com/licenses/mit/)
