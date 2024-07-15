---
"chainlink": patch
---

#bugfix Bump BSC PriceMin to 3 gwei to match BSC node's required gas price. This value can be pushed back down to 1 gwei to enable cheaper transactions if the GasPrice field under the Eth.Miner header in the BSC node's config is also pushed down to 1000000000
