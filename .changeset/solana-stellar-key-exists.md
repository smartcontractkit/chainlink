---
"chainlink": patch
---

#bugfix Solana and Stellar keystore `Import` now wrap the duplicate-key error with `keystore.ErrKeyExists`, so a node that re-imports its `ImportedSolKeys` / `ImportedStellarKeys` on restart skips the already-present key instead of exiting.
