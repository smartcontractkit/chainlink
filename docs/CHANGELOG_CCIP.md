# Changelog CCIP

All notable changes to the CCIP project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## 1.5.0 - Unreleased


## 1.4.0 - 2024-02-16

### Changed

- Token pools are no longer configured with specific on- or offRamps but rather chains through chain selectors. 
  - Token pools validate callers through the router, which should now be passed in via the constructor
  - Many events around adding/removing/updating allowed ramps now have an updated event around adding/removing/updating chains.
  - Rate limits are set for inbound and outbound explicitly, this is not different from the previous method where onRamps and offRamps are set separately.
    - Rate limits apply per remote chain, not per lane. This means that having multiple lanes enabled doesn't increase risk like it used to.
  - Enabling a chain on a token pool will allow every on- and offRamp that is configured on the router for that remote chain.
    - This also means that token pool ramp permissions do not need to be updated on a lane upgrade


- `OffRamp` offchain config format changed:
  - Renamed the `MaxGasPrice` field to `DestMaxGasPrice`.
  - Removed obsolete `SourceFinalityDepth` and `DestFinalityDepth` fields.

  This change is not breaking: the config decoder still accepts old field names.

- `CommitStore` offchain config format changed:
  - Renamed the `MaxGasPrice` field to `SourceMaxGasPrice`.
  - Removed obsolete `SourceFinalityDepth` and `DestFinalityDepth` fields.

  This change is not breaking: the config decoder still accepts old field names.

- Minor changes to the Prometheus metrics emitted by plugins
  - `ccip_unexpired_report_skipped`, `ccip_execution_observation_reports_count`, `ccip_execution_observation_build_duration`, `ccip_execution_build_single_batch`, `ccip_execution_reports_iteration_build_batch`
    removed, because they didn't introduce any additional value compared to the existing OCR2 metrics.
  - Some metrics added to track details of the processing
    - `ccip_unexpired_commit_roots` number of unexpired Commit Roots processed by Exec during the OCR2 iteration
    - `ccip_number_of_messages_processed` number of messages processed by the Exec when building the OCR2 reports and observations
    - `ccip_sequence_number_counter` latest sequence number that was used for generating Commit Report
    
- New `DynamicPriceGetter` implementation of `PriceGetter`.
  - allows for dynamic price fetching from an aggregator contract as well as using static configuration of token prices.
  - current pipeline implementation is still supported.
  - only one of the two PriceGetter implementations can be used at a time (specified in the job spec).

## 1.2.0 - 2023-11-20

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
  - new fields `destDataAvailabilityOverheadGas`, `destGasPerDataAvailabilityByte`, `destDataAvailabilityMultiplierBps` added to **DynamicConfig**.
    - `destDataAvailabilityOverheadGas` is the extra data availability gas charged on top of message data.
    - `destGasPerDataAvailabilityByte` is the amount of gas to charge per byte of data that needs data availability.
    - `destDataAvailabilityMultiplierBps` is the multiplier for data availability gas. It is in multiples of bps, or 0.0001. It can represent calldata compression factor on Rollups.
- OnRamp token transfer fee calculation updated.
  - `minTokenTransferFeeUSD` and `maxTokenTransferFeeUSD` are removed from FeeTokenConfig.
  - `minFeeUSDCents` and `maxFeeUSDCents` are added to TokenTransferFeeConfig, they will be applied at a per-token level.
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

### Renamed

- OffRamps
  - `maxTokensLength` -> `maxNumberOfTokensPerMsg`
  - `maxDataSize` -> `maxDataBytes`
  - `maxPoolGas` -> `maxPoolReleaseOrMintGas`
- OnRamp
  - `maxTokensLength` -> `maxNumberOfTokensPerMsg`
  - `maxDataSize` -> `maxDataBytes`
  - `maxGasLimit` -> `maxPerMsgGasLimit`
  - `gasMultiplier` -> `gasMultiplierWeiPerEth`
  - `premiumMultiplier` -> `premiumMultiplierWeiPerEth`
  - All fees that ended with USD denominated in cents are now suffixed with `Cents`
  - `ratio` -> `deciBps`

### Removed
- All onramp allowlist functionality is removed:
  - events AllowListAdd(address sender), AllowListRemove(address sender), AllowListEnabledSet(bool enabled)
  - applyAllowListUpdates / getAllowList
  - setAllowListEnabled / getAllowListEnabled


## 1.1.0 - 2023-08-23

### Changed
- Changed OnRamp fee calculation logic and corresponding configuration fields.
  - `destGasOverhead` and `destGasPerPayloadByte` are moved from **FeeTokenConfig** to **DynamicConfig**. These values are same on a given lane regardless of fee token.
  - `networkFeeAmountUSD` is renamed to `networkFeeUSD`. It is now multiples of 0.01 USD, as opposed to 1 wei before.
  - `minFee`, `maxFee` are moved from **TokenTransferFeeConfig** to `minTokenTransferFeeUSD`, `maxTokenTransferFeeUSD` in **FeeTokenConfig**.
  - New field `destGasOverhead` added to **TokenTransferFeeConfig**.
    - `destGasOverhead` is the amount of destination token transfer gas, to be billed as part of exec gas fee.
  