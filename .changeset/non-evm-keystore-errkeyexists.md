---
"chainlink": patch
---

#bugfix wrap duplicate-key Import errors with ErrKeyExists for the Solana, Stellar, Starknet, Sui, TON and Tron keystores, so nodes that import those keys from config can restart against a persisted keystore
