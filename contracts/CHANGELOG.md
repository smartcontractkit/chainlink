# @chainlink/contracts CHANGELOG.md

## Unreleased

### Changed:

- Bump hardhat from 2.3.3 to 2.4.1
- Move Solidity version 0.8.x contracts `ChainlinkClient.sol` and `Chainlink.sol` from `./src/v0.8/dev/` to `./src/v0.8/`.

## 0.2.0 - 2021-07-01

### Added:

- `@chainlink/contracts` package changelog.
- `KeeperCompatibleInterface` contracts.
- Feeds Registry contracts: `FeedRegistryInterface` and `Denominations`.
- v0.8 Consumable contracts (`ChainlinkClient`, `VRFConsumerBase` and aggregator interfaces).
- Multi-word response handling in v0.7 and v0.8 `ChainlinkClient` contracts.

### Changed:

- Added missing licensees to `KeeperComptibleInterface`'s
- Upgrade solidity v8 compiler version from 0.8.4 to 0.8.6
- Tests converted to Hardhat.
- Ethers upgraded from v4 to v5.
- Contract artifacts in `abi/` are now raw abi .json files, and do not include bytecode or other supplimentary data.

### Removed:

- Removed dependencies: `@chainlink/belt`, `@chainlink/test-helpers` and `@truffle`.
- Ethers and Truffle contract artifacts are no longer published.
