# @chainlink/contracts CHANGELOG.md

## Unreleased

...

## 0.5.1 - ..........

- Rename `KeeperBase` -> `AutomationBase` and add alias for backwards compatibility
- Rename `KeeperCompatible` -> `AutomationCompatible` and add alias for backwards compatibility
- Rename `KeeperCompatibleInterface` -> `AutomationCompatibleInterface` and add alias for backwards compatibility
- Rename `KeeperRegistryInterface1_2` -> `AutomationRegistryInterface1_2` and add alias for backwards compatibility

## 0.5.0 - 2022-09-26

### Changed

- Fix EIP-150 Bug in VRFV2Wrapper.sol (b9d8261eaa05838b9b609ea02005ecca3b6adca3)
- Added a new UpkeepFormat version `V2` in `UpkeepFormat`
- Renamed `KeeperRegistry` to `KeeperRegistry1_2` and `KeeperRegistryInterface` to `KeeperRegistryInterface1_2`
- Updated `UpkeepTranscoder` to only do a pass-through for upkeep bytes

## 0.4.2 - 2022-07-20

### Changed

- Downgrade 0.8.13 contracts to 0.8.6 due to [this solc bug](https://medium.com/certora/overly-optimistic-optimizer-certora-bug-disclosure-2101e3f7994d).
- Reintroduce v0.6 `EACAggregatorProxy` after removing it in [this commit](https://github.com/smartcontractkit/chainlink/commit/558f42f5122779cb2e05dc8c2b84d1ae78cc0d71)
- Ignore status update in `ArbitrumSequencerUptimeFeed` if incoming update has stale timestamp
- Revert to using current Arbitrum seq status flag in `ArbitrumSequencerUptimeFeed`
- Moved `VRFV2Wrapper`, `VRFV2WrapperConsumerBase` and `interfaces/VRFV2WrapperInterface` out of `dev` folder.

## 0.4.1 - 2022-05-09

### Changed

- VRFv2 contract pragma versions changed from `^0.8.0` to `^0.8.4`.

## 0.4.0 - 2022-02-07

### Added

- `ArbitrumSequencerUptimeFeedInterface` and `ArbitrumSequencerUptimeFeed` added in v0.8.

### Changed

- Changed `ArbitrumValidator#validate` target to `ArbitrumSequencerUptimeFeed` instead of
  Flags contract.
- Moved `VRFConsumerBaseV2` out of dev

## 0.3.1 - 2022-01-05

### Changed:

- Fixed install issue with npm.

## 0.3.0 - 2021-12-09

### Added

- Prettier Solidity formatting applied to v0.7 and above.
- ERC677ReceiverInterface added in v0.8.
- `KeeperBase.sol` and `KeeperCompatible.sol` in Solidity v0.6 and v0.8

### Changed:

- Operator Contract and Chainlink Client are officially supported. This enables
  multiword requests/response are available through the ChainlinkClient by using
  the newly enabled `buildOperatorRequest` along with `sendOperatorRequest` or
  `sendOperatorRequestTo`.
- `ChainlinkClient` functions `requestOracleData` and `requestOracleDataFrom` have been changed to `sendChainlinkRequest` and
  `sendChainlinkRequestTo` respectively.
- Updated function comments in `v0.6/interfaces/KeeperCompatibleInterface.sol` and `v0.8/interfaces/KeeperCompatibleInterface.sol` to match the latest in v0.7.
- Add `DelegateForwarderInterface` interface and `CrossDomainDelegateForwarder` base contract which implements a new `forwardDelegate()` function to forward delegatecalls from L1 to L2.

## 0.2.2 - 2021-09-21

### Added:

- v0.8 Access Controlled contracts (`SimpleWriteAccessController` and `SimpleReadAccessController`).
- v0.8 Flags contracts (`Flags`).
- v0.8 Contracts for the V2 VRF. `VRFCoordinatorV2.sol`, `VRF.sol`,
  `VRFConsumerBaseV2.sol`, `VRFCoordinatorV2Interface.sol`. Along
  with related test contract `VRFConsumerV2.sol` and example contracts
  `VRFSingleConsumerExample.sol` and `VRFConsumerExternalSubOwnerExampl.sol`.
- v0.6 `MockV3Aggregator` in src/v0.6/tests/.
- v0.7 Added keeper-related smart contracts from the keeper repo. Added tests for `KeeperRegistry` and `UpkeepRegistrationRequests` in `test/v0.7/`.

### Changed:

- Move `Operator` and associated contracts (`AuthorizedForwarder`, `AuthorizedReceiver`, `LinkTokenReceiver`, `OperatorFactory`) from `./src/v0.7/dev/` to `./src/v0.7/`.
- Updated `Denominations` in `./src/<version>` to include additional fiat currencies.
- Updated `./src/v0.8/vender/BufferChainlink.sol` with latest unchecked math version.

## 0.2.1 - 2021-07-13

### Changed:

- Bump hardhat from 2.3.3 to 2.4.1
- Move Solidity version 0.8.x contracts `ChainlinkClient.sol`, `Chainlink.sol`, `VRFConsumerBase.sol` and `VRFRequestIDBase.sol` from `./src/v0.8/dev/` to `./src/v0.8/`.
- Updated `FeedRegistryInterface` to use `base` and `quote` parameter names.
- Move `Denominations` from `./src/<version>/dev/` to `./src/<version>`

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
