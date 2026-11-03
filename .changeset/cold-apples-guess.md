---
"chainlink": patch
---

Bump `go-stellar-sdk` to v0.8.0 (Protocol 30) so the Stellar integration decodes the CAP-0088 ledger header extension (millisecond close times) that every ledger carries after the upgrade. Without this bump, ledgers produced after the upgrade fail to decode. #bugfix
