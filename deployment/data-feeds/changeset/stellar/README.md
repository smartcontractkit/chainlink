### Stellar Data Feeds Changesets

CLDF `ChangeSetV2` suite for the Soroban `DataFeedsCache` and `DataFeedsProxy`
contracts. Style follows the Solana suite (`ChangeSetV2` + operations);
functional coverage mirrors the EVM Data Feeds suite (deploy / configure /
admin / remove / ownership / proxy — `jd_*`, `migrate_feeds`, `import_address`
are EVM/JD-specific and deliberately have no Stellar equivalent here).

#### Package layout

```
data-feeds/changeset/stellar/
├── deploy_cache.go, deploy_proxy.go   # instantiate contracts, record AddressRefs
├── configure_cache.go                 # SetFeedConfigs — descriptions + permissions
├── remove_feeds.go, freeze_feeds.go   # remove feed configs, freeze/unfreeze feeds
├── feed_admin.go                      # grant/revoke feed-admin
├── ownership.go                       # transfer / accept / renounce (cache + proxy)
├── upgrade.go, recover_tokens.go      # new WASM / token recovery (cache + proxy)
├── set_proxy_cache.go                 # repoint the proxy at a cache
├── deps.go, state.go                  # shared resolve-ref/build-deps + client loaders
├── validation.go, config.go           # input validation, encoding helpers
├── operation/operation.go             # one CLDF operation per on-chain call
└── testdata/                          # e2e WASM fixtures — see testdata/README.md
```

Each changeset's `Apply` resolves a `datastore.AddressRef` (chain selector +
`ContractType` + version + qualifier) to a `stellarApplyDeps` bundle — the
contract ID plus the `stellardeps.StellarDeps` (Deploy + Invoker) needed to
call it — then executes one or more `operation.*` calls from
`operation/operation.go` through `env.OperationsBundle`. Operations are thin:
each wraps exactly one generated-binding call
(`chainlink-stellar/bindings/contracts/data_feeds_{cache,proxy}`) or one
`stellardeps.Deploy` call. `deps.go`'s `verifyContractRef` /
`resolveContractDeps` hold the chain-exists + version-parse + ref-exists +
build-deps skeleton shared by every changeset in this package.

#### Permissions model — no standalone "set forwarder"

There is no `SetForwarder`-style changeset here. A feed's writer allowlist is
part of its `FeedConfig` and is set via `SetFeedConfigs`
(`configure_cache.go`): each `FeedPermission` carries an `AllowedSender` —
the CRE forwarder contract address — plus `AllowedWorkflowOwner` and
`AllowedWorkflowName`, checked on-chain by `on_report`. This mirrors EVM's
`WorkflowMetadata{AllowedSender, AllowedWorkflowOwner, AllowedWorkflowName}`
gate. To point feeds at a new forwarder, re-run `SetFeedConfigs` with updated
`Permissions` — there's no separate on-chain "forwarder" slot to update.

#### Cross-contract version coupling

`DeployProxy` and `SetProxyCache` both resolve their cache `AddressRef` using
the *proxy's own* `req.Version`, not a separate cache-version field. This
suite's convention is that cross-contract datastore lookups share the acting
changeset's `Version` — so a cache recorded under a different version needs a
matching-version record before a proxy at that version can find it. An
optional `CacheVersion` field on these requests (to let the two diverge) is a
recorded follow-up, not implemented here.

#### Regenerating testdata WASM + bindings

The e2e test (`e2e_test.go`) deploys real WASM to a local devnet; the fixtures
in `testdata/` and the generated Go bindings in
`chainlink-stellar/bindings/contracts/data_feeds_{cache,proxy}` are a
**matched pair** — a contract ABI change requires regenerating both together,
or the e2e test will fail to compile or decode. See `testdata/README.md` for
checksums and the exact regeneration steps; in short:

```sh
# 1. rebuild WASM from the contract source repo
cd <data-feeds-stellar-integration>/contracts/data-feeds
stellar contract build            # stellar CLI 27.0.0
cp target/wasm32v1-none/release/data_feeds_{cache,proxy}.wasm \
   <chainlink>/deployment/data-feeds/changeset/stellar/testdata/
# (older stellar CLIs output to target/wasm32-unknown-unknown/release/ instead)

# 2. regenerate bindings from the new WASM, in chainlink-stellar
./scripts/gen_bindings.sh
```

#### Running the e2e test

`TestStellarDataFeedsE2E` boots a local CTF Soroban devnet, deploys both
contracts, and invokes all 40 exported functions (26 cache + 14 proxy): 38
succeed on-chain, and `recover_tokens` on each contract fails with the
documented `InvalidAction` host error (self-transfer re-entry — see the test's
`recover_tokens` comments). Reads on a frozen feed are additionally asserted
to fail with the cache's `FeedFrozen` error. It needs Docker and does not run
in CI (`skipInCI` skips when `CI=true`).

```sh
go test ./data-feeds/changeset/stellar/ -run TestStellarDataFeedsE2E -v -timeout 25m
```

The default path boots `stellar/quickstart:future` via CLDF's
`NewCTFChainProvider`, hardcoded by the CTF framework to `linux/amd64` — on a
native arm64 host without amd64 emulation this fails (`exec format error`).
Override with `STELLAR_E2E_RPC_URL` / `STELLAR_E2E_FRIENDBOT_URL` to attach to
a quickstart container started natively instead:

```sh
docker run -d -p 8000:8000 stellar/quickstart:latest \
  --local --enable-soroban-rpc --protocol-version 26
STELLAR_E2E_RPC_URL=http://127.0.0.1:8000/rpc \
STELLAR_E2E_FRIENDBOT_URL=http://127.0.0.1:8000/friendbot \
go test ./data-feeds/changeset/stellar/ -run TestStellarDataFeedsE2E -v -timeout 25m
```

See the doc comment on `newE2EChain` for why protocol 26 (not quickstart's
default 25) is required.

#### Out of scope

- **Durable-pipeline registration** in `chainlink-deployments` — follow-up
  work, coordinate with PLEX-2923.
- **MCMS** (multi-chain multisig) integration — PLEX-owned.
- **Forwarder / timelock deployment** — PLEX-owned; this package only wires
  an already-deployed forwarder's address into `FeedPermission.AllowedSender`.
