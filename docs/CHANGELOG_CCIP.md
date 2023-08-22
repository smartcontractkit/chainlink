# Changelog CCIP

All notable changes to the CCIP project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

<!-- unreleased -->
## [dev]

...
## 1.1.0 - Unreleased

### Added

### Changed
- Changed OnRamp fee calculation logic and corresponding configuration fields.
  - `destGasOverhead` and `destGasPerPayloadByte` are moved from **FeeTokenConfig** to **DynamicConfig**. These values are same on a given lane regardless of fee token.
  - `networkFeeAmountUSD` is renamed to `networkFeeUSD`. It is now multiples of 0.01 USD, as opposed to 1 wei before.
  - `minFee`, `maxFee` are moved from **TokenTransferFeeConfig** to `minTokenTransferFeeUSD`, `maxTokenTransferFeeUSD` in **FeeTokenConfig**.
  - A new field called `destGasOverhead` is added to **TokenTransferFeeConfig**. It is used to calculate destination token transfer gas cost, to be billed as part of exec gas fee.

### Removed
