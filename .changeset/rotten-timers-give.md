---
"chainlink": minor
---

Support multiple chains evm clients for TXM gas estimator to fetch L1 gas oracle
Introduced a new config flag to `[EVM.GasEstimator.DAOracle]` called `L1ChainID`, represents the L1 layer chain ID, with default value "0", marking the DA client is disabled.  
#added
