# Changelog CCIP

All notable changes to the CCIP project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

<!-- unreleased -->
## [dev]

...
## 1.0.0 - Unreleased

### Added

### Changed
- Changed OnRamp fee calculation logic and corresponding configuration fields.
  - `destGasOverhead` and `destGasPerPayloadByte` are moved from **FeeTokenConfig** to **DynamicConfig**. These values are same on a given lane regardless of fee token.
  - `networkFeeAmountUSD` is renamed to `networkFeeUSD`. It is now multiples of 0.01 USD, as opposed to 1 wei before.
  - `minFee`, `maxFee` are moved from **TokenTransferFeeConfig** to `minTokenTransferFeeUSD`, `maxTokenTransferFeeUSD` in **FeeTokenConfig**.
  - A new field called `destGas` is added to **TokenTransferFeeConfig**. It is used to calculate destination token transfer gas cost, to be billed as part of exec gas fee.

### Removed


## 0.7.0 - 30-06-2023

### Added
- TokenPool's now have getOnRamps() and getOffRamps() functions to get the list of on and off ramps for a token pool.
- ARMProxy contract (see reasoning here https://github.com/smartcontractkit/chainlink-ccip/pull/990). It's a new
contract that router, ramps and token pools use to access the ARM contract. Constructor/config changes:
  - Router constructor accepts ARM proxy 
  - ARMProxy constructor accepts ARM contract 
  - TokenPool constructor accepts ARM proxy
  - EVM2EVMOnRamp no longer accepts ARM contract in setDynamicConfig, instead ARMProxy is a new field on StaticConfig 
  - EVM2EVMOffRamp no longer accepts ARM contract in the onChainConfig passed in setOCR2Config, instead ARMProxy is a new field on StaticConfig 
  - ComitStore no longer accepts ARM contract in the onChainConfig passed in setOCR2Config, instead ARMProxy is a new field on StaticConfig 


### Changed
- We now support setting the router to zero as a cheaper way to pause the onramp vs an explicit pause function. Explicit pausing is removed.
- RampUpdates for token pool upgrades now contain a rate limit to be able to configure  
per lane token pool rate limits.
  - Before
  ```
    struct RampUpdate {
      address ramp;
      bool allowed;
    }
    ```
  - After
  ```
    struct RampUpdate {
      address ramp;
      bool allowed;
      RateLimiter.Config rateLimiterConfig;
    }
  ```
- FeeTokenConfigArgs & FeeTokenConfig have two new fields and one renamed field
  - added
    - destGasPerPayloadByte - gas cost per payload byte on destination chain
    - enabled - whether the fee token is enabled
  - renamed
    - multiplier -> gasMultiplier
- EVM2EVMOnRamp.StaticConfig has a new field
  - added
    - maxNopFeesJuels - max NOP fees in juels that the onramp can accure
- Manual execution additionally accepts a gas limit array
  - the array must be same length as messages in the report
  - if an array entry is 0, corresponding message's gas limit is used during execution
  - if an array entry is not 0, it overrides the message's gas limit
  - gas limit override cannot be lower than the original limit defined in the message
- RateLimiter errors are changed to reflect type
  - `ConsumingMoreThanCapacity` error is now `AggregateValueMaxCapacityExceeded` and `TokenMaxCapacityExceeded`
  - `RateLimitReached` error is now `AggregateValueRateLimitReached` and `TokenRateLimitReached`
- CommitOffchainConfig 
  - DestFinalityDepth property added

### Removed
- Token pool constructor no longer takes a `rateLimiterConfig` parameter.


## 0.6.0 - 2023-06-19

### Added

### Changed
- CCIP receiver revert data is now bubbled up from Router to OffRamp `ExecutionStateChanged` event
- `ExecutionStateChanged` event now contains revert data
- ARM emits `VotedToCurse` events even if it's someone's second time voting to curse in a row, or if the contract is already cursed.
- Temporary fix for finality tag support
  - CCIP uses real finality based on the finality tag for chains that support it. For other chains we use a block number based approach for finality.
  - Expect finality times >10 minutes on most chains, which some significantly longer.  

### Removed

- sequence numbers from execution reports because they were not used 

## 0.5.0 - 2023-06-07

### Added

- Commitstore: added isUnpausedAndARMHealthy as a single health check RPC call

### Changed

- TokenPools have changed to require an `allowList` parameter in constructor
  - If `allowList` is non-empty
    - only addresses in allowList can be `originalSender` when invoking lockOrBurn
    - addresses can be added to or removed from `allowList` by owner
    - allowList cannot be disabled later now
  - If `allowList` is empty
    - pool is constructed with allowList disabled
    - allowList cannot be enabled later on
- Added source chain selector and offramp address to `MessageExecuted` event emitted by Router
- AFN renamed to ARM
- Fixed onRamp allowList constructor args
- Disallow the linktoken to be a NOP
- Rework of the BurnMintERC677

## 0.4.0 - 2023-05-24

### Added

- BurnMintERC677 is the new default token that should be deployed whenever there is a need for a burn/mint token
  - Supports ERC677
  - OZ AccessControlEnumerable
  - OZ ERC20Burnable
  - OZ ERC20
  - Compatible with IBurnMintERC20 (CCIP interface)

### Changed

- IBurnMintERC20 interface has changed to follow OZ Burnable tokens
  - New interface
    - function burn(uint256 amount)
    - function burnFrom(address account, uint256 amount)
    - mint(address account, uint256 amount)
  - Old interface
    - function mint(address account, uint256 amount)
- Reduced rate limiting gas usage, this changes the config params to uint128
- Upgrade OZ dependencies to v4.8.0
- Bumped Solidity optimizations from 15k to 30k
- Config changes
  - CommitOffchainConfig
    - SourceIncomingConfirmations renamed to SourceFinalityDepth
    - DestIncomingConfirmations removed
  - ExecOffchainConfig
    - SourceIncomingConfirmations -> SourceFinalityDepth
    - DestIncomingConfirmations -> DestFinalityDepth
    - NEW DestOptimisticConfirmations (required, cannot be 0. Can be DestFinalityDepth)

### Removed

- wrapped token pools
  - Pools should be deployed as burn/mint together with a newly introduced token: BurnMintERC677.
  - This allows us to upgrade the pools without deploying a new token
  - The pool should be allowed to burn and/or mint by calling `grantMintAndBurnRoles(address pool)`


## 0.3.0 - 2023-05-09

### Added
- Added token bps fee to each individual token transfer
  - Fee structure is as follows:
    - bps fee, accurate to 0.1 bps
    - minFee, in US cents, the minimum fee to charge for 1 transfer
    - maxFee, in US cents, the maximum fee to charge for 1 transfer
  - Fee is in the range of [minFee, maxFee] 
  - The fee is configurable per token per lane per direction
  - Edge cases:
    - each token transfer is charged independently, we do not aggregate same-token transfers
    - transfers with 0 token amount is charged the minimum fee
    - all fee fields can be 0
  - The fee is charged in `feeTokens` and added to message execution fee; we do not take breadcrumbs of token transfers

### Changed

- Solidity version bumped to 0.8.19
- AggregateRateLimiter values are now in US dollar amounts with 18 decimals. Previously, it was 36 decimals.
- _setNops calls payNops 
- OnRamp and OffRamp contracts emit PoolAdded event from constructor 
- `EVM2EVMOnRamp.applyAllowListUpdates(address[] calldata removes, address[] calldata adds)` signature changed. Arguments order was `adds`, `removes`

## [0.2.0] - 2023-04-30
## [0.1.0] - 2023-03-14
