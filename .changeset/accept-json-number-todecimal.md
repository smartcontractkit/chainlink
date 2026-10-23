---
"chainlink": patch
---

#bugfix ETHABIEncode now accepts `json.Number` when converting integer ABI values, so `uint256[]` (and other integer types) from JSON no longer fail with "type json.Number cannot be converted to decimal.Decimal".
