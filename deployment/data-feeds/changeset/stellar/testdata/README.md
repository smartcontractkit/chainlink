# Stellar Data Feeds e2e test fixtures

`data_feeds_cache.wasm` and `data_feeds_proxy.wasm` are the compiled Soroban
contracts under test by `e2e_test.go` (`TestStellarDataFeedsE2E`). They are the
byte-for-byte artifacts of `stellar contract build` on chainlink-stellar's
`contracts/data-feeds` workspace at PR #161/#162 content (Bound with explicit
discriminants; literal BytesN ABI).

sha256:

- `data_feeds_cache.wasm` = `e96c7dc85c4cf37acf8acf6bbefd85d5e888d9a9a5342e39151d6f0323c70cbc`
- `data_feeds_proxy.wasm` = `e3f4a4ae333d26218bdb7aca0870623d06bdc2a65ee40440c797cfedd3e456f1`

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
