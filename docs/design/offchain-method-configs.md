# Offchain `method_configs` (Don2Don) — design note (CRE-4017)

Status: implemented (this branch). Scope: `method_configs` only; `ocr3_configs` stays on-chain
(config digest + signer/transmitter alignment are bound to the on-chain registry).

## Current source (before this change)

`method_configs` (per-method Don2Don remote config: trigger vs executable, timeouts,
aggregation, request hasher) is authored on-chain in the `CapabilitiesRegistry` contract
(`CapabilityConfig.method_configs`, proto field 7) and reaches nodes via the registry syncer:

```
CapabilitiesRegistry contract
  → registrysyncer.Sync (12s tick) → RegistryMetadata.IDsToDONs[id].CapabilityConfigurations[capID].Config (proto bytes)
    → launcher.onNewRegistry
      → addRemoteCapabilities / serveCapabilities: c.Unmarshal() → CapabilityMethodConfig
        → addRemoteCapabilityV2 (client shims: trigger subscribers, executable clients)
        → serveCapabilityV2 (server shims: trigger publishers, executable servers)
      → computeWantedShimKeys (pruning of stale shims)
```

The consumer is the **launcher** (`core/capabilities/launcher.go`), *not*
`LocalCapabilityManager.startCapability` — method configs configure the don2don shims
(remote trigger/executable clients and servers), not the local capability process.

## Proposed seam (implemented)

Mirror the `spec_config` cutover: **offchain-wins-with-fallback**, gated by the same
per-node `[Capabilities.Local] UseOffchainRegistry` flag.

- **Overlay point**: in `launcher.onNewRegistry`, after `c.Unmarshal()` yields the on-chain
  `CapabilityMethodConfig` map, overlay the offchain `method_configs` for that
  (DON name, capability ID) when the gate is on. Offchain method entries **replace** the
  on-chain entry for the same method (a method config is a typed struct, not a key/value
  map, so per-key merging is not meaningful); methods the offchain payload omits keep
  their on-chain config. A capability absent offchain keeps its full on-chain method set.
- **Keying**: the offchain registry is keyed by on-chain DON **name** (same as `spec_config`).
  The launcher resolves each registry DON's name via `localcapmgr.OffchainDONNames`
  (exported from the existing unexported helper; same ambiguity rules: unnamed or
  name-colliding DONs keep their on-chain config).
- **Delivery/reactivity**: the registry syncer re-drives `OnNewRegistry` on every 12s tick
  regardless of on-chain change, so an offchain-only `method_configs` update propagates
  within one tick without new plumbing. (The launcher does not Subscribe to GlobalConfig;
  the 12s re-drive is the reactive path. This is acceptable because shim `SetConfig` is
  idempotent and cheap when unchanged.)
- **Ingestion validation**: `globalconfig.Validate` now also checks every offchain
  `method_configs` entry converts cleanly (pb → `capabilities.CapabilityMethodConfig`),
  so the launcher only ever sees valid entries. Malformed entries are rejected at the job,
  not silently dropped at the launcher.
- **Conversion**: a single reusable converter `globalconfig.MethodConfigsFromProto`
  (pb `CapabilityMethodConfig` → common `capabilities.CapabilityMethodConfig`), mirroring
  the conversion inline in `registry.CapabilityConfiguration.Unmarshal` in chainlink-common.
  It lives in `globalconfig` (leaf package) so both the validator and the launcher use the
  exact same conversion the on-chain path uses.
- **Pruning**: `computeWantedShimKeys` is computed from the **post-overlay** method maps
  (the same maps the create/update loops use), so a method moved offchain-only is pruned
  and re-created consistently.

## Blast radius

- `core/capabilities/globalconfig` — new exported converter + extended `validate` (rejects
  malformed `method_configs` at ingestion). No behavior change for payloads without
  `method_configs`.
- `core/capabilities/localcapmgr` — `offchainDONNames` exported as `OffchainDONNames`
  (pure rename); `offchainChangesConfig` extended to also compare `method_configs`
  (feeds the existing `config_mismatch` divergence kind).
- `core/capabilities` (launcher) — new `SetOffchainRegistry(gc, useOffchainRegistry)`
  setter + overlay in `addRemoteCapabilities`/`serveCapabilities`/`computeWantedShimKeys`.
  Gate off or no payload → byte-for-byte today's behavior.
- `core/services/cre/cre.go` — wires `opts.OffchainCapabilitiesRegistry` +
  `localCfg.UseOffchainRegistry()` into the launcher (one line, next to the existing
  LocalCapabilityManager wiring).
- Config docs (`docs/CONFIG.md`, `core/config/docs/core.toml`, TOML comment) — the
  `UseOffchainRegistry` comment now mentions `method_configs` alongside `spec_config`.
- `.changeset` — minor entry.

Not in scope: `ocr3_configs` (stays on-chain), `oracle_factory_configs` (CRE-1775 was resolved
upstream by enabling the oracle factory on-chain for capabilities with an OCR3 config — see
`cre.go` `newServicesFn`; the offchain registry does not carry it), CLD authoring of
`method_configs` payloads (follows the CLD reusable operation).
