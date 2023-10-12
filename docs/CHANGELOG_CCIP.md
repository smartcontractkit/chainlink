# Changelog CCIP

All notable changes to the CCIP project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

<!-- unreleased -->
## [dev]

## 1.2.0 - Unreleased

### Added

- USDC Support
  - Added `USDC` token pool
  - Jobspec changes to support USDC information
- Added TypeAndVersion to all token pools & price registry

### Changed
- PriceUpdate now accepts an array of gas price update
  - Removed `destChainSelector` and `usdPerUnitGas` from PriceUpdates
  - Added `GasPriceUpdate[] gasPriceUpdates` to PriceUpdates. Each `GasPriceUpdate` struct contains `destChainSelector` and `usdPerUnitGas`.
- OnRamp fee calculation logic now includes L1 security fee if sending to L2.
  - New field `destBytesOverhead` added to **TokenTransferFeeConfig**.
    - `destBytesOverhead` is the size of additional bytes being passed to destination for token transfers. For example, USDC transfers require additional attestation data.
  - new fields `destDataAvailabilityOverheadGas`, `destGasPerDataAvailabilityByte`, `destDataAvailabilityMultiplier` added to **DynamicConfig**.
    - `destDataAvailabilityOverheadGas` is the extra data availability gas charged on top of message data.
    - `destGasPerDataAvailabilityByte` is the amount of gas to charge per byte of data that needs data availability.
    - `destDataAvailabilityMultiplier` is the multiplier for data availability gas. It is in multiples of 1e-4, or 0.0001. It can represent calldata compression factor on Optimistic Rollups.
- OnRamp token transfer fee calculation updated.
  - `minTokenTransferFeeUSD` and `maxTokenTransferFeeUSD` are removed from FeeTokenConfig.
  - `minFeeUSD` and `maxFeeUSD` are added to TokenTransferFeeConfig, they will be applied at a per-token level.
  - token transfer premium is calculated as the sum of each individual token transfer fee.
- MessageId hashing logic updated.
  - the new `sourceTokenData` field is added to the hash.
  - fixed-size fields are hashed in nested hash function.
- CommitStore OffchainConfig fields updated.
  - New fields `GasPriceHeartBeat`, `DAGasPriceDeviationPPB`, `ExecGasPriceDeviationPPB`, `TokenPriceHeartBeat`, `TokenPriceDeviationPPB` added
    - `GasPriceHeartBeat` specifies an update heartbeat threshold for gas prices
    - `DAGasPriceDeviationPPB` specifies deviation PPB threshold for dava availability (DA) gas price. On chains without DA component, this should be 0.
    - `ExecGasPriceDeviationPPB` specifies deviation PPB threshold for native EVM execution gas price.
    - `TokenPriceHeartBeat` specifies an update heartbeat threshold for token prices
    - `TokenPriceDeviationPPB` specifies deviation PPB threshold for token price.
  - Old Fields `FeeUpdateHeartBeat`, `FeeUpdateDeviationPPB` removed. They are replaced by the fields above.
- OffRamp caps gas passed on to TokenPool when calling `releaseOrMint`.
  - A new `maxPoolGas` field is added to OffRamp **DynamicConfig** to store this gas limit.
- OnRamp will revert with `SourceTokenDataTooLarge` if TokenPool returns too much data.
  - The revert threshold is `destBytesOverhead` in **TokenTransferFeeConfig**.


### Removed
- All onramp allowlist functionality is removed:
  - events AllowListAdd(address sender), AllowListRemove(address sender), AllowListEnabledSet(bool enabled)
  - applyAllowListUpdates / getAllowList
  - setAllowListEnabled / getAllowListEnabled


## 1.1.0 - Unreleased

### Changed
- Changed OnRamp fee calculation logic and corresponding configuration fields.
  - `destGasOverhead` and `destGasPerPayloadByte` are moved from **FeeTokenConfig** to **DynamicConfig**. These values are same on a given lane regardless of fee token.
  - `networkFeeAmountUSD` is renamed to `networkFeeUSD`. It is now multiples of 0.01 USD, as opposed to 1 wei before.
  - `minFee`, `maxFee` are moved from **TokenTransferFeeConfig** to `minTokenTransferFeeUSD`, `maxTokenTransferFeeUSD` in **FeeTokenConfig**.
  - New field `destGasOverhead` added to **TokenTransferFeeConfig**.
    - `destGasOverhead` is the amount of destination token transfer gas, to be billed as part of exec gas fee.
  