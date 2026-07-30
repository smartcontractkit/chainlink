# Griddle Reconciler Phase Plan

This file is the execution plan derived from [TODO.md](/Users/bartektofel/Desktop/repos/chainlink/core/scripts/cre/reconciler/TODO.md). Follow the phases in order.

## Ground Rules

- Backward compatibility does not matter. Prefer deleting old paths over preserving them.
- Do not add transition code, dual schema support, or compatibility shims unless this file explicitly says so.
- Keep each phase narrowly scoped. Do not mix state/schema refactors, discovery changes, and UI work in one patch unless the phase says to do that.
- After every phase:
  - run targeted package tests for touched packages
  - run `go test ./core/scripts/cre/reconciler/...`
  - update this file or `TODO.md` if the phase changes scope or reveals a blocker

## Finalized Decisions

### D1. `registry_based_launch_allowlist`

- Treat `registry_based_launch_allowlist` as a slice of capability names, not addresses.
- Examples:
  - `cron-trigger@1.0.0`
  - `evm-1337`
- UI validation for E6 should validate:
  - non-empty string
  - no duplicates
  - optional simple format check: printable token with no surrounding whitespace
- Do not apply address validation here.

### D2. Deployer key configuration

- Standardize on environment variables only.
- Remove `--deployer-key` entirely.
- Use:
  - `PRIVATE_KEY_<CHAIN_ID>` for per-chain EVM deployer keys
  - `PRIVATE_KEY` as the global fallback
  - Anvil default only when neither env var is set
- Workflow owner must be derived from the effective key for the registry chain.
- Do not mix CLI flags and environment variables.

### D3. Breaking changes

- Breaking changes are acceptable.
- Remove obsolete schema/CLI fields directly instead of supporting old and new forms together.
- Specifically:
  - remove `infra.chart_values`
  - remove `infra.namespace`
  - remove `--deployer-key`
  - choose the final gateway schema directly instead of supporting both old and new forms

### D4. Non-EVM capability scope

- First cut supports:
  - `cron`
  - `consensus`
  - `don-time`
  - `http-action`
  - `http-trigger`
  - `vault`
  - `evm`
  - `solana`
  - `aptos`
- Defer `stellar`.
- Do not expose `write-solana` as a standalone capability. Treat it as Solana-specific config or internal behavior.

### D5. Label validation

- Chart remains the source of truth for:
  - DON membership
  - node role
  - node namespace
- Preflight must validate:
  - chart-side `chainlink-node.registerNodes.labels.don-name`
  - JD-side node labels `p2p_id`, `environment`, and `type`
- `don-<Name>` is not required for the first cut.

### D6. Gateway schema

- Final schema for gateway assignment:

```toml
[[gateway_nodes]]
node = "node-gw-0"
dons = ["workflow-a", "capabilities-b"]
```

- Do not support the old single-`don` shape.
- Scope of support:
  - one gateway nodeset may contain multiple gateway nodes
  - multiple nodesets may each contain gateway nodes
  - any nodeset treated as a gateway nodeset is assumed to contain only gateway nodes
  - mixed worker+gateway nodesets are out of scope

### D7. UI apply security model

- `/api/apply` must remain localhost-only.
- Only one apply may run at a time.
- At the TOML breakpoint, the UI must stop and wait for manual operator confirmation.
- The UI must not run git commands.

### D8. CLD artifacts

- Defer CLD artifacts implementation.
- Do not spend implementation time on G3 in this workstream.
- Remove or rewrite any README text that promises CLD artifacts.

### D9. Parallelization boundary

- Keep `runPreEnvStartup` sequential in the first pass.
- Keep `runPostEnvStartup` sequential in the first pass.
- Parallelize only clearly independent per-node/per-read loops first.

## Phase Order

1. Phase 0: decision lock and planning scaffolding
2. Phase 1: P0 OCR2 discovery correctness
3. Phase 2: resumable on-chain steps
4. Phase 3: per-phase hashing and forced reruns
5. Phase 4: non-EVM discovery and preflight validation
6. Phase 5: capability catalog and backend capability completion
7. Phase 6: UI quick wins and config cleanup
8. Phase 7: gateway model, deployer key unification, safe concurrency
9. Phase 8: UI apply flow and diff/status
10. Phase 9: docs, artifacts, renames, cleanup

---

## Phase 0: Decision Lock And Scaffolding

### Goal

Record the decisions above in the backlog and remove ambiguity before implementation starts.

### Files to touch

- [TODO.md](/Users/bartektofel/Desktop/repos/chainlink/core/scripts/cre/reconciler/TODO.md)
- [phase_plan.md](/Users/bartektofel/Desktop/repos/chainlink/core/scripts/cre/reconciler/phase_plan.md)

### Tasks

1. Update the open-decision sections in `TODO.md` to match the finalized decisions in this file.
2. Mark explicitly deferred items:
   - `stellar`
   - standalone `write-solana`
   - CLD artifacts
3. Update H1 in `TODO.md`:
   - remove any mention of `--deployer-key`
   - replace it with env-only resolution
4. Update E6 in `TODO.md`:
   - describe the allowlist entries as capability names, not addresses
5. Update E1 in `TODO.md`:
   - use the final `dons = [...]` schema only
6. Update E8 in `TODO.md`:
   - removal of `infra.chart_values` and `infra.namespace` is a direct breaking change

### Validation

- No code behavior changes in this phase.
- Run:
  - `go test ./core/scripts/cre/reconciler/internal/domain/...`

### Exit criteria

- `TODO.md` and `phase_plan.md` no longer contain contradictory decisions.

---

## Phase 1: P0 OCR2 Discovery Correctness

### Goal

Fix the bug where bootstrap and gateway nodes are treated like OCR signer nodes.

### Primary backlog items

- B3

### Files to inspect first

- [internal/discovery/discovery.go](/Users/bartektofel/Desktop/repos/chainlink/core/scripts/cre/reconciler/internal/discovery/discovery.go)
- [internal/onchain/topology.go](/Users/bartektofel/Desktop/repos/chainlink/core/scripts/cre/reconciler/internal/onchain/topology.go)
- [internal/onchain/hydrate.go](/Users/bartektofel/Desktop/repos/chainlink/core/scripts/cre/reconciler/internal/onchain/hydrate.go)
- [deps.go](/Users/bartektofel/Desktop/repos/chainlink/core/scripts/cre/reconciler/deps.go)
- [reconcile.go](/Users/bartektofel/Desktop/repos/chainlink/core/scripts/cre/reconciler/reconcile.go)

### Implementation steps

1. Add role-aware behavior to discovery:
   - worker nodes read OCR2 bundle IDs
   - bootstrap nodes skip OCR2 bundle reads
   - gateway nodes skip OCR2 bundle reads
2. Update the discovery adapter interfaces if needed, but do not widen scope beyond OCR2 gating in this phase.
3. Update topology hydration:
   - only require OCR2 bundle IDs for worker/plugin nodes
   - do not error on missing OCR2 bundles for bootstrap/gateway nodes
4. Keep all EVM address discovery behavior unchanged in this phase.

### Tests to add or update

- [internal/discovery/discovery_test.go](/Users/bartektofel/Desktop/repos/chainlink/core/scripts/cre/reconciler/internal/discovery/discovery_test.go)
- [internal/onchain/hydrate_test.go](/Users/bartektofel/Desktop/repos/chainlink/core/scripts/cre/reconciler/internal/onchain/hydrate_test.go)
- any topology test that currently assumes OCR2 is required for all nodes

### Validation

- Run:
  - `go test ./core/scripts/cre/reconciler/internal/discovery`
  - `go test ./core/scripts/cre/reconciler/internal/onchain`
  - `go test ./core/scripts/cre/reconciler/...`

### Exit criteria

- Bootstrap/gateway discovery performs no OCR2 read.
- Topology build succeeds without OCR2 bundle IDs for bootstrap/gateway nodes.
- Worker behavior is unchanged.

---

## Phase 2: Resumable On-Chain Steps

### Goal

Split the monolithic on-chain apply flow into resumable, independently persisted steps.

### Primary backlog items

- A3
- A2

### Files to inspect first

- [internal/onchain/deployer.go](/Users/bartektofel/Desktop/repos/chainlink/core/scripts/cre/reconciler/internal/onchain/deployer.go)
- [internal/onchain/deploy.go](/Users/bartektofel/Desktop/repos/chainlink/core/scripts/cre/reconciler/internal/onchain/deploy.go)
- [internal/onchain/env.go](/Users/bartektofel/Desktop/repos/chainlink/core/scripts/cre/reconciler/internal/onchain/env.go)
- [reconcile.go](/Users/bartektofel/Desktop/repos/chainlink/core/scripts/cre/reconciler/reconcile.go)
- [internal/domain/state.go](/Users/bartektofel/Desktop/repos/chainlink/core/scripts/cre/reconciler/internal/domain/state.go)

### Implementation steps

1. ~~Identify the existing P1-P8 sequence in `Deployer.Apply`.~~ Done pre-Phase-2 — the `internal/onchain`
   subpackage split already happened and `Apply` already calls 9 distinct steps.
2. ~~Extract each major step into its own method~~ — already done (build CLDF env, build topology, build CRE env,
   deploy contracts, run PreEnvStartup, prepare JD chain configs, configure CapReg, resolve DON IDs, configure
   WorkflowReg all already exist as separate `Deployer` methods). Fixed the stale log header
   (`"=== P1-P6: ... ==="`, wrong step count) — replaced with a plain, non-numbered header.
3. Persist-after-success review: `resolveDONIDs` (`state.DONIDs`) and `configureWorkflowReg` (`state.WorkflowReg`)
   already write their own state artifact only after success, and `Apply` already persists after each. **No new
   checkpoint/phase marker was added** for `PreEnvStartup` or `configureCapReg` — see the decided-scope note below.
4. **Decided scope (differs from the original wording here — see TODO.md A3):** `DeployV2RegistryContractsSequence`
   deploys CapabilitiesRegistry + WorkflowRegistry together against a fresh in-memory datastore with no existence
   check of its own — it cannot deploy just one of the two without bypassing the sequence and duplicating its
   datastore-merge logic. Instead of a true independent per-contract guard, `deployContracts` now uses a
   `contractsFullyDeployed(state)` predicate that requires **both** addresses present to skip; if either is missing,
   it redeploys **both** via the existing sequence (accepting that a contract that already existed gets redeployed/
   orphaned in that case).
5. Treat `PreEnvStartup` as atomic — confirmed unchanged: no per-forwarder or per-DON partial-completion logic was
   added. `PreEnvStartup` and `configureCapReg` have no per-run completion marker by design — they're full-replace/
   idempotent operations, and skip-if-input-unchanged is explicitly A1's job (Phase 3), not this phase's. Always
   re-running them on a resumed `Apply` is correct, not a gap.
6. Boundary-state resume, given step 4's fix:
   - only CapReg exists → redeploys both (see step 4), then the rest of `Apply` proceeds normally.
   - CapReg and WorkflowReg exist → deploy step skips; `PreEnvStartup`/`configureCapReg`/`resolveDONIDs`/
     `configureWorkflowReg` safely re-run (idempotent, no data loss).
   - contracts exist but not configured → same as above, resumes correctly with no changes needed beyond step 4.
7. Public CLI behavior unchanged — confirmed, no flags touched.

### Tests to add or update

- `internal/onchain/deploy_test.go` (new): `contractsFullyDeployed` (neither/CapReg-only/WorkflowReg-only/both) and
  `deployContracts` skip-path test.
- No end-to-end `Apply` test harness was built — none existed before this phase, and building one (faking CLDF env +
  JD client) is disproportionate to the narrowed scope above; deferred until it's actually needed.

### Validation

- Run:
  - `go test ./core/scripts/cre/reconciler/internal/onchain` — 36 passed
  - `go test ./core/scripts/cre/reconciler/...` — 147 passed, 9 packages

### Exit criteria

- A failed run can resume from the failed step — true for deploy (per the decided redeploy-both behavior) and for
  every step after it (idempotent re-run, no data loss).
- Already-completed steps (contract deploy when both present) are not repeated.
- Per-contract deployment *check* is independent (both addresses checked separately); the *remedy* when not fully
  deployed is a combined redeploy, not an independent per-contract deploy — see TODO.md A3.

---

## Phase 3: Per-Phase Hashing And Force Reruns

### Goal

Add input-hash memoization for opaque phases and a force-rerun escape hatch.

### Decided (differs from the steps below — see TODO.md A1 for full rationale)

- **PreEnvStartup + Configure CapReg share one hash key**, not two — PreEnvStartup's output contains a protobuf type
  that isn't cleanly TOML-serializable, so they always run or skip together.
- **No `--force` flag.** The repair path is deleting the relevant `phase_hashes` entry (or the whole table) in the
  state file by hand — more precise than a blanket force flag, since the cascade rule (below) propagates it forward
  automatically.
- **Cascade invalidation, Docker-layer-cache style**, instead of embedding contract addresses/DON IDs into every hash:
  `deploy contracts → [PreEnvStartup+CapReg] → resolve DON IDs → Configure WorkflowReg → Jobs`. A fresh contract
  deploy (vs. skip), or any hash-gated phase actually running, forces every phase after it to run too, regardless of
  its own hash — otherwise a freshly (re)deployed-but-unconfigured contract could be silently skipped.
- **`TOMLPatchApplied` (one-time ratchet) is retired.** `injectTOML` runs every apply; the new direct-compare
  (`infra.ReadLayerValue`) makes it a no-op when nothing differs, rather than a bool gate that only ever fires once.

### Primary backlog items

- A1

### Files to inspect first

- [internal/domain/state.go](/Users/bartektofel/Desktop/repos/chainlink/core/scripts/cre/reconciler/internal/domain/state.go)
- [reconcile.go](/Users/bartektofel/Desktop/repos/chainlink/core/scripts/cre/reconciler/reconcile.go)
- [cmd.go](/Users/bartektofel/Desktop/repos/chainlink/core/scripts/cre/reconciler/cmd.go)
- [internal/onchain/deployer.go](/Users/bartektofel/Desktop/repos/chainlink/core/scripts/cre/reconciler/internal/onchain/deployer.go)
- [internal/ui/server.go](/Users/bartektofel/Desktop/repos/chainlink/core/scripts/cre/reconciler/internal/ui/server.go) later, for status/diff consumers

### Implementation steps (as implemented)

1. Added `PhaseHashes map[string]string` to `StateFile`, keyed `pre_env_startup_capreg` / `configure_workflowreg` /
   `jobs` (see decisions above for why PreEnvStartup+CapReg share one key).
2. Added `domain.CanonicalHash` — `sha256(json.Marshal(v))`. `encoding/json` already sorts map keys deterministically;
   callers (in `internal/onchain/hash.go`) explicitly sort slices (DON lists, p2p_id member lists) before hashing, so
   no custom canonicalization walker was needed.
3. Hashing scope: PreEnvStartup+CapReg combined (excludes bootstrap-only/gateway DONs, mirroring
   `capRegWorkerNodes`'s `flags.HasNoOtherFlags` check), Configure WorkflowReg (sorted workflow DON names + derived
   workflow-owner address), Jobs (all DON types + gateway service configs).
4. Kept actual-state checks unchanged: contract deployment (`contractsFullyDeployed`, Phase 2), JD chain configs
   (library-level idempotency, unchanged).
5. TOML patching uses direct compare (`infra.ReadLayerValue`), not hashing — and `TOMLPatchApplied`'s one-time ratchet
   was retired so this comparison actually gets a chance to run on every apply, not just the first one.
6. No `--force` flag (Decision 3) — manual `phase_hashes` deletion is the repair path, and the cascade rule (Decision
   2) means deleting one entry already propagates forward.
7. N/A — no CLI changes in this phase.
8. Each hash-gated phase stores its hash only after success (`state.SetPhaseHash`, called post-success in `Apply`/
   `SyncJobs`).
9. Documented in TODO.md A1: not drift-safe for hashed phases; manual `phase_hashes` deletion (or the whole table) is
   the repair path.

### Tests added

- [internal/domain/state_test.go](/Users/bartektofel/Desktop/repos/chainlink/core/scripts/cre/reconciler/internal/domain/state_test.go) — `CanonicalHash` determinism, `PhaseHashes` get/set, `PhaseNeedsRun` cascade semantics.
- [internal/onchain/hash_test.go](/Users/bartektofel/Desktop/repos/chainlink/core/scripts/cre/reconciler/internal/onchain/hash_test.go) (new) — hash determinism/sensitivity per phase, bootstrap/gateway DON exclusion.
- [internal/infra/chartpatch_test.go](/Users/bartektofel/Desktop/repos/chainlink/core/scripts/cre/reconciler/internal/infra/chartpatch_test.go) — `ReadLayerValue` (missing file, absent layer, reflects patched content).
- [reconcile_test.go](/Users/bartektofel/Desktop/repos/chainlink/core/scripts/cre/reconciler/reconcile_test.go) — `filterUnchangedTOMLPatches` (drops identical, keeps changed, all-unchanged yields empty).
- [deploy_test.go](/Users/bartektofel/Desktop/repos/chainlink/core/scripts/cre/reconciler/internal/onchain/deploy_test.go) — updated for `deployContracts`'s new `(bool, error)` return.

### Validation

- Run:
  - `go test ./core/scripts/cre/reconciler/internal/domain` — passing
  - `go test ./core/scripts/cre/reconciler/internal/onchain` — passing
  - `go test ./core/scripts/cre/reconciler/...` — 161 passed, 9 packages (up from 147 baseline, +14 new tests)
  - `go build ./...` in `system-tests/lib` — passing (new `crecontracts.BindCapabilityRegistry` helper)

### Exit criteria

- A no-change apply performs no on-chain or JD writes, and `injectTOML` patches nothing.
- A changed capability config reruns the combined PreEnvStartup+CapReg unit, which cascades into WorkflowReg + Jobs.
- Deleting a `phase_hashes` entry (the manual repair path, replacing `--force`) reruns it and everything after it in
  the pipeline.

---

## Phase 4: Non-EVM Discovery And Preflight

### Goal

Teach discovery/preflight about Solana and Aptos, and fail earlier on configuration issues.

### Primary backlog items

- B1
- B2
- B4

### Files to inspect first

- [deps.go](/Users/bartektofel/Desktop/repos/chainlink/core/scripts/cre/reconciler/deps.go)
- [nodeapi.go](/Users/bartektofel/Desktop/repos/chainlink/core/scripts/cre/reconciler/nodeapi.go)
- [internal/discovery/discovery.go](/Users/bartektofel/Desktop/repos/chainlink/core/scripts/cre/reconciler/internal/discovery/discovery.go)
- [internal/domain/state.go](/Users/bartektofel/Desktop/repos/chainlink/core/scripts/cre/reconciler/internal/domain/state.go)
- [internal/domain/chartvalues.go](/Users/bartektofel/Desktop/repos/chainlink/core/scripts/cre/reconciler/internal/domain/chartvalues.go)
- [internal/onchain/topology.go](/Users/bartektofel/Desktop/repos/chainlink/core/scripts/cre/reconciler/internal/onchain/topology.go)
- [internal/nodeconfig/chains.go](/Users/bartektofel/Desktop/repos/chainlink/core/scripts/cre/reconciler/internal/nodeconfig/chains.go)

### Decided (differs from the steps below — see TODO.md B1/B2/B4 for full rationale)

- **`[[chains]]` unified with a `family` field** (required, no default — breaking change), not split EVM/non-EVM
  chain types, per explicit direction. Non-EVM entries need no RPC URLs and can't be the registry chain.
- **Aptos/Solana key reads are capability-gated**, not unconditional like EVM — only nodes whose DON declares the
  matching capability get a read attempted, since most nodes have no such key configured.
- **The chain-config-per-node check was NOT moved earlier.** It stays exactly where the EVM version already ran
  (inside `buildTopology`, before contracts get configured) and was extended to Solana/Aptos there — a brand-new
  node-set legitimately doesn't have chart chain config yet, so running this very early would break first-time
  applies.
- **The JD-label preflight DOES run early** — unconditionally near the top of `Run()`, on every invocation (even
  once on-chain work is complete), via a new standalone `buildOffchainClient` factory decoupled from full CLDF
  environment setup.
- **Solana support was added to the vendored `clclient` dependency** (chainlink-testing-framework), via a temporary
  `replace` directive in `core/scripts/go.mod`, since it had no Solana key-read method at all (only Aptos).

### Implementation steps (as implemented)

1. `NodeClient`/`Client` gained `ReadAptosKeys`/`ReadSolanaKeys` — single native address each, not chain-ID-keyed.
2. `NodeRuntimeInfo` gained `AptosAddress`/`SolanaAddress string`.
3. `discovery.Run`/`discoverOne` take a per-node needed-families map (`reconcile.go`'s `nonEVMFamiliesByNode`,
   computed from `desired.DONs` + `don.WorkerNodes(cv)`) and only read a family's key when needed.
4. New preflight: `onchain.ValidateNodeLabels` (JD labels) runs early/unconditionally in `Run()`. The chain-config
   check stays in `buildTopology` (see decided-scope note above) — no separate early step for it.
5. Preflight validates:
   - chart-side `don-name` — unchanged, already ran at chart-load time.
   - JD-side labels `p2p_id`, `environment` (== `desired.JD.Environment`), `type` (== `bootstrap`/`gateway`/`plugin`
     — workers are `"plugin"` in JD, not `"standard"`).
   - every worker node in every DON has the required chain config for each chain-scoped capability, EVM and now
     Solana/Aptos too (`validateDiscoveredNonEVMAddresses` in `hydrate.go`), at the same call site as before.
6. `stellar` stays out of scope — confirmed nowhere touched.
7. Existing apply-time validation (`validateDiscoveredEVMAddresses`) is unchanged, still the backstop.

### Tests added

- `internal/discovery/discovery_test.go` — capability-gated Aptos/Solana reads.
- `internal/domain/desired_test.go` — `Chain.Family` validation (missing/unsupported family, non-EVM registry
  rejected, solana/aptos capability cross-reference), `DON.NonEVMFamilies()`.
- `internal/onchain/hydrate_test.go` — `validateDiscoveredNonEVMAddresses`.
- `internal/onchain/preflight_test.go` (new) — `expectedJDTypeLabel`, `ValidateNodeLabels` (skips when JD not
  configured, errors on missing discovered CSA key — without needing a real JD connection).

### Validation

- Run:
  - `go test ./core/scripts/cre/reconciler/internal/discovery` — passing
  - `go test ./core/scripts/cre/reconciler/internal/domain` — passing
  - `go test ./core/scripts/cre/reconciler/internal/onchain` — passing
  - `go test ./core/scripts/cre/reconciler/...` — 175 passed, 9 packages (up from 162 baseline, +13 new tests)
  - `go build ./...` in `chainlink-testing-framework/framework` — passing (new Solana key-read support)

### Exit criteria

- Solana and Aptos keys are available in runtime state, for nodes whose DON needs them.
- Misconfigured chain config still fails before contracts get configured (unchanged timing, now covers non-EVM too).
- Missing/wrong JD labels produce precise errors, checked on every `Run()` invocation.

---

## Phase 5: Capability Catalog And Backend Completion

### Goal

Align the UI catalog, defaults, discovery, and backend wiring around the actual supported capability set.

### Primary backlog items

- C2

### Files to inspect first

- [internal/ui/server.go](/Users/bartektofel/Desktop/repos/chainlink/core/scripts/cre/reconciler/internal/ui/server.go)
- [internal/ui/web/capability_defaults.toml](/Users/bartektofel/Desktop/repos/chainlink/core/scripts/cre/reconciler/internal/ui/web/capability_defaults.toml)
- [internal/domain/desired.go](/Users/bartektofel/Desktop/repos/chainlink/core/scripts/cre/reconciler/internal/domain/desired.go)
- [internal/onchain/env.go](/Users/bartektofel/Desktop/repos/chainlink/core/scripts/cre/reconciler/internal/onchain/env.go)
- [internal/onchain/topology.go](/Users/bartektofel/Desktop/repos/chainlink/core/scripts/cre/reconciler/internal/onchain/topology.go)

### Decided (differs from the steps below — see TODO.md C2 for full rationale)

- **Chain-selector mismatch discovered during investigation:** Solana's blockchain provider (`system-tests/lib/cre/environment/blockchains/solana`) resolves a chain selector from a base58 genesis-hash-style string via `chainselectors.SolanaChainIdToChainSelector()` — incompatible with the opaque numeric `chain_id` Phase 4 established for `solana-<N>` capability scoping. Resolved by adding a generic `domain.Chain.ChainSpecificData map[string]any` field (`toml:"chain_specific"`) rather than a one-off `GenesisHash` struct field, so the same mechanism extends to future non-EVM families. Solana reads `chain_specific.genesis_hash` via a `Chain.SolanaGenesisHash()` helper; `chain_id` keeps its Phase-4 opaque-number semantics unchanged.
- **Aptos has no such mismatch:** its `Deploy()` parses `chain_id` as a plain uint64, so `chain.ChainID` is reused directly — no extra field needed.
- **write-solana fold:** per explicit user direction, the standalone `[capability_configs.write-solana]` table was deleted outright rather than merged into `solana`'s values (no Go code read it by name, so nothing depended on the merge).
- **Step 3 ("extend backend environment construction") turned out to be load-bearing, not optional/cosmetic:** the CRE feature hooks (`system-tests/lib/cre/features/aptos`, `.../solana/v2`) hard-require a matching `creEnv.Blockchains` entry — aptos errors, solana nil-pointer-panics — the moment a DON has that capability enabled and gets applied. So "aptos selectable in the UI" could not ship safely without also completing this wiring in the same phase.
- Solana requires the `SOLANA_PRIVATE_KEY` env var (no default on the Kubernetes/external-provider path, unlike EVM's auto-defaulted `PRIVATE_KEY`) and a writable `ContractsDir`, which the reconciler generates internally via `os.MkdirTemp` (operational detail, not user-declared desired state).

### Implementation steps

1. Update the UI catalog to the finalized support set:
   - add `aptos` ✅
   - keep `solana` ✅
   - do not add `write-solana` ✅ (was already absent)
2. Clean `capability_defaults.toml`:
   - keep defaults for real supported capabilities ✅
   - delete the standalone `write-solana` table (not merged — see Decided above) ✅
3. Extend backend environment construction so the reconciler can actually supply the chain families required by Aptos/Solana feature code ✅ — `buildCreEnvironment` (env.go) now switches on `chain.Family` and calls the matching CRE blockchain deployer (solana/aptos), reusing their existing Kubernetes-path `Deploy()` — no new provider-construction code.
4. Verify chain-scoped capability handling in `desired.go` still matches the final catalog ✅ — `ChainScopedCapabilities`/`stripChainSuffix`/`validChainFamilies` were already consistent (aptos present everywhere) from Phase 4; `validateChains()` gained per-family branches (solana: `ws_url`+`http_url`+`chain_specific.genesis_hash` required; aptos: `http_url` required, `ws_url` rejected).

### Tests to add or update

- UI server capability catalog tests ✅ (`aptos` present + chainScoped, `write-solana`/`mock` absent)
- desired-state tests for chain-scoped capability handling ✅ (new validation-error cases + positive solana/aptos chain fixtures)
- env tests for the family-switch wiring ✅ (`SOLANA_PRIVATE_KEY`-required error path; full EVM/aptos happy-path wiring isn't unit-tested since `evm.Deploy()` dials a real RPC endpoint even on the Kubernetes path — consistent with why no prior test exercised `buildCreEnvironment` end-to-end)

### Validation

- Ran:
  - `go test ./core/scripts/cre/reconciler/internal/ui`
  - `go test ./core/scripts/cre/reconciler/internal/domain`
  - `go test ./core/scripts/cre/reconciler/internal/onchain`
  - `go test ./core/scripts/cre/reconciler/...` → 181 passed, 9 packages (up from 175)
  - `go vet ./core/scripts/cre/reconciler/...` → no issues
  - `gofmt -l` → clean

### Exit criteria

- The visible capability catalog matches the supported set. ✅
- `aptos` is selectable end-to-end. ✅
- standalone `write-solana` does not appear in the UI. ✅

---

## Phase 6: UI Quick Wins And Config Cleanup

### Goal

Ship the low-risk UI improvements and remove stale desired-state config fields.

### Primary backlog items

- E6
- E7
- E5
- E8

### Files to inspect first

- [internal/ui/web/app.js](/Users/bartektofel/Desktop/repos/chainlink/core/scripts/cre/reconciler/internal/ui/web/app.js)
- [internal/ui/web/index.html](/Users/bartektofel/Desktop/repos/chainlink/core/scripts/cre/reconciler/internal/ui/web/index.html)
- [internal/ui/server.go](/Users/bartektofel/Desktop/repos/chainlink/core/scripts/cre/reconciler/internal/ui/server.go)
- [internal/domain/desired.go](/Users/bartektofel/Desktop/repos/chainlink/core/scripts/cre/reconciler/internal/domain/desired.go)
- [cmd.go](/Users/bartektofel/Desktop/repos/chainlink/core/scripts/cre/reconciler/cmd.go)
- [reconcile.go](/Users/bartektofel/Desktop/repos/chainlink/core/scripts/cre/reconciler/reconcile.go)

### Decided (differs slightly from the steps below — see TODO.md E5-E8 for full detail)

- **Config tab removed entirely, not left as an empty placeholder** (explicit user decision): once the Infrastructure
  card (E8) and Job Distributor card (E5) both leave it, nothing remains, so the tab and its nav button were deleted
  rather than kept empty for future settings.
- **E7 fix went one step further than the stated scope:** `removeDON` had a related gap — it only clamped
  `selectedDON` when it fell past the shrunk array's end, but never re-targeted it when a DON *before* or *at* the
  selected index was deleted, silently pointing the selector at the wrong DON. Fixed alongside the default-selection
  change since both use the same new `defaultDONIndex()` helper.
- **E6's UI pattern:** reused the existing chip+`prompt()` add/remove convention already used for chain-scoped
  capabilities, rather than introducing a new textarea-based input — keeps the DON editor visually consistent.

### Implementation steps

1. Add the E6 allowlist input to the DON editor UI ✅ — chip+`prompt()` editor in `donColumnHTML`, backed by
   `don.registryBasedLaunchAllowlist`; non-empty/trimmed/no-duplicates validation; entries stay untyped capability-name
   strings, never reinterpreted as addresses.
2. Change E7 default DON selection ✅ — first workflow DON, fallback to index 0; `removeDON` re-validation fixed too.
3. Move JD connectivity/settings to a JD tab ✅.
4. Remove `infra.chart_values` and `infra.namespace` from the desired-state schema ✅.
5. Remove the corresponding UI inputs ✅ (the whole Config tab, per the Decided note above).
6. Make `apply` load the chart from `--chart-dir`, same as `serve` ✅ — `diff` got the same flag too, since it had the
   identical `ds.Infra.ChartValues` pattern.
7. Derive namespace from chart values everywhere ✅ — `runServe`'s `Infra.Namespace` fallback removed; `checkJD`'s
   namespace source in app.js switched to the already-chart-derived `state.namespace`.

### Tests to add or update

- desired-state parsing tests ✅ (stripped `chart_values`/`namespace` from every `[infra]` fixture; removed the
  obsolete "missing chart_values" validation-error case)
- UI server request/response tests ✅ (`TestResponseToDesiredState`, `TestAPI_Desired_SaveAndLoad` updated)
- UI tests for allowlist round-trip ✅ (`TestAPI_Desired_SaveAndLoad` now round-trips `registryBasedLaunchAllowlist`)

### Validation

- Ran:
  - `go test ./core/scripts/cre/reconciler/internal/ui`
  - `go test ./core/scripts/cre/reconciler/internal/domain`
  - `go test ./core/scripts/cre/reconciler/...` → 180 passed, 9 packages (181 minus the one obsolete test removed)
  - `go vet ./core/scripts/cre/reconciler/...` → no issues
  - `gofmt -l` → clean

### Exit criteria

- UI no longer exposes chart path or namespace fields. ✅
- Desired-state schema no longer contains those fields. ✅
- Allowlist can be edited in the UI as capability names. ✅

---

## Phase 7: Gateway Model, Deployer Key Unification, Safe Concurrency

### Goal

Complete the new gateway model, unify deployer key handling, and parallelize safe loops.

### Primary backlog items

- E1
- E2
- H1
- D1

### Files to inspect first

- [internal/domain/desired.go](/Users/bartektofel/Desktop/repos/chainlink/core/scripts/cre/reconciler/internal/domain/desired.go)
- [internal/ui/server.go](/Users/bartektofel/Desktop/repos/chainlink/core/scripts/cre/reconciler/internal/ui/server.go)
- [internal/ui/web/app.js](/Users/bartektofel/Desktop/repos/chainlink/core/scripts/cre/reconciler/internal/ui/web/app.js)
- [internal/onchain/topology.go](/Users/bartektofel/Desktop/repos/chainlink/core/scripts/cre/reconciler/internal/onchain/topology.go)
- [internal/onchain/gateway.go](/Users/bartektofel/Desktop/repos/chainlink/core/scripts/cre/reconciler/internal/onchain/gateway.go)
- [cmd.go](/Users/bartektofel/Desktop/repos/chainlink/core/scripts/cre/reconciler/cmd.go)
- [reconcile.go](/Users/bartektofel/Desktop/repos/chainlink/core/scripts/cre/reconciler/reconcile.go)
- [internal/onchain/env.go](/Users/bartektofel/Desktop/repos/chainlink/core/scripts/cre/reconciler/internal/onchain/env.go)
- [internal/onchain/workflowreg.go](/Users/bartektofel/Desktop/repos/chainlink/core/scripts/cre/reconciler/internal/onchain/workflowreg.go)

### Implementation steps

1. Replace the single-DON gateway assignment model with `dons = [...]`.
2. Support both gateway deployment shapes:
   - one gateway nodeset containing multiple gateway nodes
   - multiple nodesets, each containing gateway nodes
3. Assume a gateway nodeset contains only gateway nodes.
4. Update topology and connector generation so one gateway node can serve multiple DONs, including `capabilities` DONs.
5. Update UI gateway assignment to multi-select.
6. Block save/preview when any gateway has zero assigned DONs.
7. Remove `--deployer-key` from the CLI.
8. Implement env-only effective key resolution:
   - `PRIVATE_KEY_<CHAIN_ID>`
   - fallback `PRIVATE_KEY`
   - fallback Anvil default
9. Use the effective registry-chain key for workflow-owner derivation.
10. Parallelize only safe loops:
   - `buildNodeSpecs`
   - `buildDonsForJobs`
   - other clearly independent per-node reads
11. Keep `runPreEnvStartup` and `runPostEnvStartup` sequential.

### Tests to add or update

- desired-state gateway schema tests
- topology/gateway tests
- CLI flag tests for removal of `--deployer-key`
- deployer key resolution tests
- race-safe concurrency tests where practical

### Validation

- Run:
  - `go test ./core/scripts/cre/reconciler/internal/domain`
  - `go test ./core/scripts/cre/reconciler/internal/onchain`
  - `go test -race ./core/scripts/cre/reconciler/...`

### Exit criteria

- Gateway assignments use the final multi-DON schema only.
- Deployer key behavior is env-only and consistent.
- Added concurrency passes `-race`.

---

## Phase 8: UI Apply Flow And Diff/Status

### Goal

Expose apply and reconciliation insight in the UI without weakening safety.

### Primary backlog items

- E3
- E4

### Files to inspect first

- [internal/ui/server.go](/Users/bartektofel/Desktop/repos/chainlink/core/scripts/cre/reconciler/internal/ui/server.go)
- [internal/ui/web/app.js](/Users/bartektofel/Desktop/repos/chainlink/core/scripts/cre/reconciler/internal/ui/web/app.js)
- [internal/ui/web/index.html](/Users/bartektofel/Desktop/repos/chainlink/core/scripts/cre/reconciler/internal/ui/web/index.html)
- [internal/domain/state.go](/Users/bartektofel/Desktop/repos/chainlink/core/scripts/cre/reconciler/internal/domain/state.go)
- [reconcile.go](/Users/bartektofel/Desktop/repos/chainlink/core/scripts/cre/reconciler/reconcile.go)

### Implementation steps

1. Add `/api/apply` with:
   - localhost-only access
   - single-flight locking
   - streaming logs
2. Build an execution view in the UI:
   - current phase
   - live log tail
   - final success/failure state
3. At the TOML breakpoint:
   - stop execution
   - show manual instructions
   - resume only after operator confirmation
4. Build a diff/status screen using Phase 3 hashes and persisted state.
5. Do not build a second reconciliation model for the UI.

### Tests to add or update

- UI handler tests for apply and streaming behavior
- state/status response tests
- any small frontend logic tests if a harness exists

### Validation

- Run:
  - `go test ./core/scripts/cre/reconciler/internal/ui`
  - `go test ./core/scripts/cre/reconciler/...`

### Exit criteria

- UI can run apply with live logs.
- UI pauses at the TOML breakpoint and resumes only on confirmation.
- Diff/status accurately reflects hashed skip/rerun behavior.

---

## Phase 9: Docs, Artifacts, Renames, Cleanup

### Goal

Finish documentation and polish after the behavior is stable.

### Primary backlog items

- G1
- G2
- housekeeping

### Files to inspect first

- [README.md](/Users/bartektofel/Desktop/repos/chainlink/core/scripts/cre/reconciler/README.md)
- [cmd.go](/Users/bartektofel/Desktop/repos/chainlink/core/scripts/cre/reconciler/cmd.go)
- [internal/domain/desired.go](/Users/bartektofel/Desktop/repos/chainlink/core/scripts/cre/reconciler/internal/domain/desired.go)
- [internal/ui/web/index.html](/Users/bartektofel/Desktop/repos/chainlink/core/scripts/cre/reconciler/internal/ui/web/index.html)
- [internal/ui/web/app.js](/Users/bartektofel/Desktop/repos/chainlink/core/scripts/cre/reconciler/internal/ui/web/app.js)

### Implementation steps

1. Update README for:
   - current phase model
   - env-only deployer key behavior
   - `--chart-dir`
   - chart-owned chain assumptions
   - schema-breaking changes
2. Remove or rewrite README text that promises CLD artifacts.
3. Rename `boot` to `bootstrap` everywhere relevant.
4. Remove the dead root `web/` directory if it still exists.

### Tests to add or update

- any tests affected by `boot` -> `bootstrap`
### Validation

- Run:
  - `go test ./core/scripts/cre/reconciler/...`

### Exit criteria

- Docs match actual behavior.
- No stale deployer-key docs remain.
- Any remaining aspirational artifact docs are removed.

---

## Stop Conditions

Stop and ask for direction only if one of these happens:

- Aptos or Solana node API reads are impossible with the available client surface and require a dependency upgrade
- the system-tests feature code requires a broader non-EVM environment build refactor than Phase 5 assumes
- a supposedly independent concurrency target fails `-race` and requires a deeper architectural change

If a stop condition is not hit, continue executing phases without asking for more planning.
