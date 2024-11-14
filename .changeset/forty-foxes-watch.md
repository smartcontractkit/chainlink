---
"chainlink": minor
---

Updated Solana TXM's in-memory storage to track statuses across the Solana transaction lifecycle. Added a method to translate Solana transaction statuses into states expected by the ChainWriter interface. Made the duration transactions are retained in storage after finality or error configurable using `TxRetentionTimeout`. #added
