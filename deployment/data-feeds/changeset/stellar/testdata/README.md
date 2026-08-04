# Stellar Data Feeds e2e test fixtures

`data_feeds_cache.wasm` and `data_feeds_proxy.wasm` are the compiled Soroban
contracts under test by `e2e_test.go` (`TestStellarDataFeedsE2E`). They are the
byte-for-byte artifacts of `stellar contract build` on chainlink-stellar's
`contracts/data-feeds` workspace at commit `971e21ae` (PR #161/#162 heads;
adds `is_configured`/`is_frozen`/`set_feed_frozen` and optional
`decimals`/`description` returns).

sha256:

- `data_feeds_cache.wasm` = `c494264eaaa2cf5e5a745c0022ac791ffa07d1beb24a1d2ac67fcb6d7f28ef0f`
- `data_feeds_proxy.wasm` = `f670dc316877fac6f6c594133f411af25254d4a305038eeade28d34b9382fe6a`

## Regeneration

These are committed binaries (the e2e test needs a real WASM to deploy to a
local CTF Soroban devnet). To regenerate after a contract change, rebuild from
the contract source repo and copy the artifacts back:

```sh
cd <chainlink-stellar>/contracts/data-feeds
stellar contract build            # stellar CLI 27.0.0
cp target/wasm32v1-none/release/data_feeds_cache.wasm  <this-dir>/
cp target/wasm32v1-none/release/data_feeds_proxy.wasm  <this-dir>/
```

(Older stellar CLI versions output to `target/wasm32-unknown-unknown/release/`
instead — check which target triple your `stellar contract build` uses.)

If the contract ABI changes, regenerate the Go bindings
(`chainlink-stellar/bindings/contracts/data_feeds_{cache,proxy}`) as well, or the
e2e test will fail to compile / decode.
