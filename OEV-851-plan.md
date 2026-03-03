# OEV-851: Ingest Mevshare params to Beholder

**Ticket goal:** Send Flashbots/MEVShare params (and optionally Atlas params) in the Beholder `TxMessage` the core node emits whenever a secondary tx gets broadcasted.

---

## Status

- [x] Exploration complete
- [x] Proto changes (`chainlink-protos`) — field added, Go code regenerated, CHANGELOG + version bumped to 1.2.0
- [x] chainlink-evm changes (`chainlink-evm`) — `EmitTxMessage` updated, tests added and passing; local `replace` directive pointing at local `chainlink-protos` for dev
- [ ] chainlink-protos PR merged + tag cut
- [ ] chainlink-evm `go.mod` updated to published tag (remove replace directive), chainlink-evm PR merged + tag cut
- [ ] chainlink (core) bump of chainlink-evm dep

---

## Architecture overview

The system involves **three repositories** that need coordinated changes:

```
chainlink-protos/svr   →   chainlink-evm   →   chainlink (core)
  (TxMessage proto)       (EmitTxMessage)       (no code change needed;
                                                  just a dep bump)
```

### Flow for a secondary (dual-broadcast) transaction

1. OCR2 plugin calls `ocr2FeedsDualTransmission.CreateSecondaryEthTransaction()` in  
   `chainlink-evm/pkg/transmitter/ocr/dual_transmitter.go`
2. That sets `TxMeta.DualBroadcast = true` and `TxMeta.DualBroadcastParams = urlParams()`  
   where `urlParams()` URL-encodes the `secondaryMeta` map (Flashbots/Atlas params)
3. The TXMv2 picks up the transaction; `sendTransactionWithError` in  
   `chainlink-evm/pkg/txm/txm.go` calls `client.SendTransaction()`
4. The dual-broadcast client (`FlashbotsClient` or `MetaClient`) reads  
   `meta.DualBroadcastParams` from the transaction meta and uses those params
5. After a **successful** broadcast, `txm.go` calls  
   `t.Metrics.EmitTxMessage(ctx, attempt.Hash, fromAddress, tx)`
6. `EmitTxMessage` in `chainlink-evm/pkg/txm/metrics.go` builds `svrv1.TxMessage`  
   and emits it to Beholder

### Current `TxMessage` proto (field numbers 1-7)

```proto
message TxMessage {
  string hash = 1;
  string from_address = 2;
  string to_address = 3;
  string nonce = 4;
  int64 created_at = 5;
  string chain_id = 6;
  string feed_address = 7;
}
```

### `TxMeta.DualBroadcastParams` format

URL-encoded query string with keys depending on the backend:

**Flashbots / MEVShare keys** (parsed by `parseURLParams()` in `flashbots_client.go`):
- `auctionTimeout` (int)
- `builder` (repeated)
- `hint` (repeated)
- `refund` (string, format `address:percent`)

**Atlas / Meta keys** (parsed by `VerifyResponse()` in `meta_client.go`):
- `destination` (address — the metacall contract)
- `dapp` (address — the dApp control contract)

---

## What we want to add

When a secondary transaction is broadcast (i.e. `meta.DualBroadcast == true`),
include the dual-broadcast params in the Beholder message so that the observability
pipeline (Atlas / Chip Ingress) can see which Flashbots/MEVShare/Atlas parameters
were used.

### Proposed new proto fields

The cleanest approach is to add a **single string field** that carries the raw
`DualBroadcastParams` URL-encoded string, plus a **boolean** to indicate this was
a dual-broadcast. This avoids parsing complexity in the proto layer and keeps the
schema stable across both Flashbots and Atlas backends.

Alternatively, we could decompose into structured sub-fields. This is richer but
requires the proto to know about backend-specific keys, which couples it to the
clients. The raw-string approach is simpler and sufficient for "visibility".

**Decision: use a raw string for now** (mirrors the `FastLaneAtlasError` pattern
which also stores the Atlas URL as a plain string).

Proposed additions (field numbers continue from 7):

```proto
message TxMessage {
  // existing fields 1-7 ...
  bool   dual_broadcast        = 8;   // true iff this is a secondary/dual-broadcast tx
  string dual_broadcast_params = 9;   // raw URL-encoded params (Flashbots or Atlas)
}
```

The `beholder_data_schema` version should be bumped from `2` to `3`.

---

## Repositories and files to change

### 1. `chainlink-protos` (repo: `/Users/ggerritsen/dev/cll/chainlink-protos`)

**File:** `svr/v1/beholder_tx_message.proto`
- Add fields `dual_broadcast` (8) and `dual_broadcast_params` (9)

**File:** `svr/v1/beholder_tx_message.pb.go`
- Regenerate with `buf generate` (or protoc)

**File:** `svr/CHANGELOG.md`
- Add minor version entry

**File:** `svr/iron-flask.yaml`
- Schema registration stays the same (entity is the same proto file)

After merging, create a new pre-release tag so `chainlink-evm` can depend on it.

### 2. `chainlink-evm` (repo: `/Users/ggerritsen/dev/cll/chainlink-evm`)

**File:** `pkg/txm/metrics.go` — `EmitTxMessage()`
- Read `meta.DualBroadcast` and `meta.DualBroadcastParams`
- Populate new proto fields when meta is present
- Bump schema version string from `"/beholder-tx-message/versions/2"` to `"/beholder-tx-message/versions/3"`

**File:** `pkg/txm/metrics_test.go`
- Add test case for dual-broadcast transaction: verify `DualBroadcast=true` and
  `DualBroadcastParams` populated correctly
- Add test case for non-dual-broadcast: verify `DualBroadcast=false` / field absent

**Dependency:** bump `chainlink-protos/svr` to new version in `go.mod`

After merging, create a new pre-release tag so `chainlink` (core) can depend on it.

### 3. `chainlink` (core) — this repo

**No code changes needed.** Only a `go.mod` bump of `chainlink-evm` to pick up the
new `chainlink-evm` tag with the updated `EmitTxMessage`.

---

## Step-by-step execution plan

### Step 1: Update the proto (`chainlink-protos`)

1. Edit `svr/v1/beholder_tx_message.proto` — add fields 8 and 9
2. Regenerate Go code (`buf generate` in the `svr/` directory)
3. Update `svr/CHANGELOG.md` with a minor version entry
4. Open PR, get review, merge
5. Cut a new pre-release tag (e.g. `svr/v1.1.2-...`)

### Step 2: Update `EmitTxMessage` (`chainlink-evm`)

1. Bump `chainlink-protos/svr` in `go.mod` to the new tag
2. Edit `pkg/txm/metrics.go`:
   - After resolving `meta`, check `meta.DualBroadcast` and `meta.DualBroadcastParams`
   - Set `message.DualBroadcast` and `message.DualBroadcastParams` accordingly
   - Change schema version to `"/beholder-tx-message/versions/3"`
3. Update `pkg/txm/metrics_test.go` with new test cases
4. Open PR, get review, merge
5. Cut a new pre-release tag

### Step 3: Bump dep in `chainlink` (core)

1. `go get github.com/smartcontractkit/chainlink-evm@<new-tag>`
2. `go mod tidy`
3. Open PR

---

## Decisions (confirmed)

**Q1: Field structure** → Raw URL-encoded string for `dual_broadcast_params`. Simple, opaque, works for both Flashbots and Atlas without coupling the proto to backend-specific keys.

**Q2: Atlas params** → Use `DualBroadcastParams` as-is. It already contains `destination` and `dapp` for Atlas, and the Flashbots keys. No additional Atlas-specific fields needed.

**Q3: Schema version bump** → No bump. Proto field additions are backward-compatible; keep `beholder_data_schema` at `"/beholder-tx-message/versions/2"`.

**Q4: Primary txs** → Don't set `dual_broadcast_params` when `meta.DualBroadcast` is nil/false (proto optional field — omitempty means not set ≠ empty string). Field simply absent for primary txs.

---

## Key file locations

| Repo | File | Purpose |
|------|------|---------|
| chainlink-protos | `svr/v1/beholder_tx_message.proto` | Proto definition |
| chainlink-protos | `svr/v1/beholder_tx_message.pb.go` | Generated Go (auto) |
| chainlink-evm | `pkg/txm/metrics.go` | `EmitTxMessage` implementation |
| chainlink-evm | `pkg/txm/metrics_test.go` | Tests for `EmitTxMessage` |
| chainlink-evm | `pkg/txm/types/transaction.go` | `TxMeta` struct (read-only reference) |
| chainlink-evm | `pkg/transmitter/ocr/dual_transmitter.go` | Where `DualBroadcastParams` is set |
| chainlink-evm | `pkg/txm/clientwrappers/dualbroadcast/flashbots_client.go` | Flashbots param parsing |
| chainlink-evm | `pkg/txm/clientwrappers/dualbroadcast/meta_client.go` | Atlas param parsing |
| chainlink (core) | `go.mod` | Dep on chainlink-evm (bump needed) |

---

## Notes on proto regeneration

In `chainlink-protos/svr/`:
```
buf generate
```
This uses `buf.gen.jd.yaml` and `buf.yaml`. The generated `.pb.go` file must be
committed alongside the `.proto` file.

---

## Notes on `beholder_data_schema` versioning

The string `"/beholder-tx-message/versions/2"` is an Atlas Chip Ingress schema registry
path. Bumping it tells Atlas that the message format changed. Since proto additions of
optional fields are backward-compatible for decoders, the need to bump depends on whether
the Atlas ingestion pipeline requires a new schema registration. This is worth confirming
with the platform/Atlas team.
