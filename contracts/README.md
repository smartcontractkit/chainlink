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
import "@chainlink/contracts/src/v0.8/KeeperCompatibleInterface.sol";
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

## Keeper contracts deployment

[Design doc](https://www.notion.so/chainlink/Keeper-V2-94415970f1ef4b46ba0f6aebee1cd477)

The following networks are supported by keeper:

- `kovan`: Ethereum testnet with chain ID 42
- `mumbai`: Polygon testnet with chain ID 80001
- `bsctestnet`: BSC Testnet with chain ID 97

### Deploy

Set environment variables `<NETWORK>_PRIVATE_KEY` and `<NETWORK>_RPC_URL` according to the network smart contracts will be deployed to.
Environment variables per network:

- `kovan`:
  - `KOVAN_PRIVATE_KEY`
  - `KOVAN_RPC_URL`
- `mumbai`:
  - `MUMBAI_PRIVATE_KEY`
  - `MUMBAI_RPC_URL`
- `bsctestnet`:
  - `BSCTESTNET_PRIVATE_KEY`
  - `BSCTESTNET_RPC_URL`

Then run:

```bash
$ yarn keeper:deploy:<network-name>
```

`<network-name>` is the value from the supported networks list above.

### Verify on Etherscan

Set environment variables `ETHERSCAN_API_KEY`.

Then run:

```bash
$ yarn keeper:verify:<network-name>
```

`<network-name>` is the value from the supported networks list above.

## Contributing

Contributions are welcome! Please refer to
[Chainlink's contributing guidelines](../docs/CONTRIBUTING.md) for detailed
contribution information.

Thank you!

## License

[MIT](https://choosealicense.com/licenses/mit/)
