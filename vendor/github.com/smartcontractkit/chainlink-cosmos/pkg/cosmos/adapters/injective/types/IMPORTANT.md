⚠️ Important
============

When syncing types from chain-side OCR module, copy only the relevant files.
Inclusion of extra generated files such as proto gw, keeper helpers, etc, cofuses future users of this integration.

### List of relevant files

Revised on _Oct 7th 2021_

```
codec.go
errors.go
genesis.pb.go
msgs.go
ocr.pb.go
params.go
proposal.go
query.pb.go
tx.pb.go
types.go
```

### OCR Cosmos module

Lives in https://github.com/InjectiveLabs/injective-core/tree/f/ocr/injective-chain/modules/ocr for now.
