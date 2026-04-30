---
"chainlink": patch
---

#nops Drop deprecated parity-ethereum / openethereum references from
user-facing docs and config comments. Both projects have been
unmaintained for years; mentioning them as supported execution clients
is misleading. Updated:

- README's "officially supported execution clients" list (drops
  Parity/Openethereum, leaves Geth and Besu).
- The `MaxInFlight` comment in `chains-evm.toml`, which used parity as a
  comparison point for the conservative default.
- The `Gnosis_Mainnet.toml` consensus comment, which framed AuRa
  finality in terms of parity.

Resolves #4018.
