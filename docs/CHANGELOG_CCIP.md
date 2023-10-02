# Changelog CCIP

All notable changes to the CCIP project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

<!-- unreleased -->
## [dev]

## 1.2.0 - Unreleased

### Added

### Changed
- OnRamp fee calculation logic now includes L1 security fee if sending to L2.
  - New field `destBytesOverhead` added to **TokenTransferFeeConfig**.
    - `destBytesOverhead` is the size of additional bytes being passed to destination for token transfers. For example, USDC transfers require additional attestation data.
  - new fields `destDataAvailabilityOverheadGas`, `destGasPerDataAvailabilityByte`, `destDataAvailabilityMultiplier` added to **DynamicConfig**.
    - `destDataAvailabilityOverheadGas` is the extra data availability gas charged on top of message data.
    - `destGasPerDataAvailabilityByte` is the amount of gas to charge per byte of data that needs data availability.
    - `destDataAvailabilityMultiplier` is the multiplier for data availability gas. It is in multiples of 1e-4, or 0.0001. It can represent calldata compression factor on Optimistic Rollups.
- MessageId hashing logic updated.
  - the new `sourceTokenData` field is added to the hash.
  - fixed-size fields are hashed in nested hash function.
- CommitStore OffchainConfig fields updated.
  - New fields `GasPriceHeartBeat`, `DAGasPriceDeviationPPB`, `ExecGasPriceDeviationPPB`, `TokenPriceHeartBeat`, `TokenPriceDeviationPPB` added
  - Old Fields `FeeUpdateHeartBeat`, `FeeUpdateDeviationPPB` removed.

### Removed



## 1.1.0 - Unreleased

### Changed
- Changed OnRamp fee calculation logic and corresponding configuration fields.
  - `destGasOverhead` and `destGasPerPayloadByte` are moved from **FeeTokenConfig** to **DynamicConfig**. These values are same on a given lane regardless of fee token.
  - `networkFeeAmountUSD` is renamed to `networkFeeUSD`. It is now multiples of 0.01 USD, as opposed to 1 wei before.
  - `minFee`, `maxFee` are moved from **TokenTransferFeeConfig** to `minTokenTransferFeeUSD`, `maxTokenTransferFeeUSD` in **FeeTokenConfig**.
  - New field `destGasOverhead` added to **TokenTransferFeeConfig**.
    - `destGasOverhead` is the amount of destination token transfer gas, to be billed as part of exec gas fee.
  