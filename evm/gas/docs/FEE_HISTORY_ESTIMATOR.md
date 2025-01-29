# Fee History Estimator

## Overview
`Fee History` estimator is an EVM-based gas estimator that utilizes RPC calls to make gas price estimations. The estimator heavily relies on two RPC calls: `eth_gasPrice` and `eth_feeHistory`. It is built as a service and caches the calculated results in order to minimize overhead. While bumping, it prioritizes using the latest result from most recent blocks if it exceeds the bumped gas price. `Fee History` estimator supports both Legacy and Dynamic(EIP-1559) transactions. It can also handle chains that don't have a mempool and process transactions on an FCFS basis.

## Configs
- `BumpPercent`: is the percentage by which to bump gas on a transaction. This is used when the estimator's bumping API gets called.
- `CacheTimeout`: is the time to wait to refresh the cached values. A small jitter is applied so the timeout won't be exactly the same each time.
- `EIP1559`: enables EIP-1559 mode and deactivates Legacy estimations. This means the estimator will refresh prices and make estimations only for Dynamic transactions.

The rest of the configs are only applicable when `EIP1559` is enabled

- `BlockHistorySize`: controls the number of past blocks to include during gas calculations. If set to 0, the estimator will skip any priority fee calculation and calculate the underlying base fee. This config should be set to 0 for chains that don't have a mempool. 
- `RewardPercentile`: specifies which fee percentile to pick from for each processed past block. 

### Validations
During startup, the estimator will perform two config checks:
- `BumpPercent` is equal to or higher than *MinimumBumpPercentage*. *MinimumBumpPercentage* is fixed at 10% and it's the minimum percentage allowed by Geth when bumping a transaction, to prevent spam attacks. Replacing a transaction with a price less than 10% from the previous one will result in an error on the RPC side. Even for chains that don't enforce that rule, a 10% bump seems reasonable.
- `RewardPercentile` is equal or lower than *ConnectivityPercentile*, when `EIP1559` is enabled. *ConnectivityPercentile* is fixed at the 85th percentile and it's the maximum percentile we're willing to bump a transaction's price. This is used as a sanity check method in order to avoid excessive gas bumping when an RPC is not responding.

## As a Service
`Fee History` estimator is built as a service. This means it will periodically poll the RPC for new prices, perform off-chain calculations, and cache the result for future use. For simplicity, only one type of gas estimation can be enabled at a time, Legacy or Dynamic. The poll interval is controlled by `CacheTimeout`. This value should be close to the block time. For slower chains, like Ethereum, you can set it to 12s, the same as the block time. For faster chains, you can skip a block or two, as prices will be refreshed more frequently. Ideally, 1s should be the absolute minimum. 


## Legacy Gas Price Estimations
### Fetching
Periodically, `Fee History` estimator will call `eth_gasPrice` RPC method to fetch the gas price reported by the RPC. The parameters of this call can not be controlled by the user, meaning the result can sometimes be stale, especially during sudden gas spikes. It is advisable to use EIP-1559 if the chain supports it.

### Bumping
During bumping, `Fee History` will refresh the cached value by making a call to the RPC. The bumped value of the original price will be compared with the market price and the highest value will be returned.

## Dynamic Price Estimations
### Fetching
`Fee History` estimator periodically calls `eth_feeHistory` method to get the most up-to-date information from the RPC (more information about the call can be found [here](https://ethereum.github.io/execution-apis/api-documentation/)). It fetches three things:
- Base Fee of the next block
- The Yth priority fee percentiles of the past X blocks, where Y is controlled by `RewardPercentile` and X by `BlockHistorySize`.
- The 85th priority fee percentiles of the past X blocks.

The above values are used to construct and cache the following:
- **MaxPriorityFeePerGas**: the average of Yth priority fee percentiles, excluding zero values.
- **MaxFeePerGas**: `baseFee * BaseFeeBufferPercentage + MaxPriorityFeePerGas`. *BaseFeeBufferPercentage* is used as a safety to catch any fluctuations in the Base Fee during the next blocks.
- **PriorityFeeThreshold**: the max out of every 85th priority fee percentile. This value is used to stop the estimator from bumping a price above that threshold and represents the maximum allowed value.

*Note*: for chains that don't have a mempool (activated with `BlockHistorySize=0`) **MaxPriorityFeePerGas** and **PriorityFeeThreshold** are set to 0 since there is no concept of gas bumping.

### Bumping
For bumping, `Fee History` estimator bumps both maxPriorityFeePerGas and maxFeePerGas of the original transaction attempt. This is required by Geth, along with the 10% minimum bumping threshold. The bumped price is compared to the cached market prices stored in the estimator and the highest of the two is picked. Finally, the resulting maxPriorityFeePerGas gets compared to the cached PriorityFeeThreshold value. If the bumped value is higher, this indicates a potential connection issue with the RPC, and bumping is skipped, returning an error.

*Note*: for chains that don't have a mempool (activated with `BlockHistorySize=0`) bumping works differently. Instead, we force-fetch the most up-to-date Base Fee value and embed it in the MaxFeePerGas. MaxPriorityFeePerGas remains 0.

### Metrics
The following prometheus metrics are exposed:
- **gas_price_updater**: latest Gas Price stored
- **base_fee_updater**: Base Fee of the next block
- **max_priority_fee_per_gas_updater**: latest MaxPriorityFeePerGas stored
- **max_fee_per_gas_updater**: latest MaxFeePerGas stored