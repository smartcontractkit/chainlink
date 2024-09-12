---
"chainlink": minor
---

#internal Updated the TXM confirmation logic to use the mined transaction count to identify re-org'd or confirmed transactions.

- Confirmer uses the mined transaction count to determine if transactions have been re-org'd or confirmed.
- Confirmer no longer sets transaction states to `confirmed_missing_receipt`. This state is maintained in queries for backwards compatibility.
- Finalizer now responsible for fetching and storing receipts for confirmed transactions.
- Finalizer now responsible for resuming pending task runs.
- Finalizer now responsible for marking old transactions without receipts broadcasted before the finalized head as fatal.
