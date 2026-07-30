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

1. Identify the existing P1-P8 sequence in `Deployer.Apply`.
2. Extract each major step into its own method:
   - build CLDF env
   - build topology
   - build CRE env
   - deploy contracts
   - run PreEnvStartup
   - prepare JD chain configs
   - configure CapReg
   - resolve DON IDs
   - configure WorkflowReg
3. After each successful state-mutating step:
   - persist state
   - update the phase/checkpoint marker
4. Replace the current coarse deploy skip with per-contract address guards.
5. Treat `PreEnvStartup` as atomic:
   - do not try to resume midway through `PreEnvStartup`
   - do not introduce per-forwarder resume logic
   - do not introduce per-DON partial completion inside `PreEnvStartup`
6. Ensure the realistic contract/config boundary states resume cleanly:
   - only CapReg exists
   - CapReg and WorkflowReg exist
   - contracts exist, but CapReg and WorkflowReg have not been configured yet
7. Keep the public CLI behavior unchanged in this phase.

### Tests to add or update

- create or expand deployer tests for partial resume
- update any state tests that assume a single on-chain phase checkpoint

### Validation

- Run:
  - `go test ./core/scripts/cre/reconciler/internal/onchain`
  - `go test ./core/scripts/cre/reconciler/...`

### Exit criteria

- A failed run can resume from the failed step.
- Already-completed steps are not repeated.
- Per-contract deployment skip works independently.

---

## Phase 3: Per-Phase Hashing And Force Reruns

### Goal

Add input-hash memoization for opaque phases and a force-rerun escape hatch.

### Primary backlog items

- A1

### Files to inspect first

- [internal/domain/state.go](/Users/bartektofel/Desktop/repos/chainlink/core/scripts/cre/reconciler/internal/domain/state.go)
- [reconcile.go](/Users/bartektofel/Desktop/repos/chainlink/core/scripts/cre/reconciler/reconcile.go)
- [cmd.go](/Users/bartektofel/Desktop/repos/chainlink/core/scripts/cre/reconciler/cmd.go)
- [internal/onchain/deployer.go](/Users/bartektofel/Desktop/repos/chainlink/core/scripts/cre/reconciler/internal/onchain/deployer.go)
- [internal/ui/server.go](/Users/bartektofel/Desktop/repos/chainlink/core/scripts/cre/reconciler/internal/ui/server.go) later, for status/diff consumers

### Implementation steps

1. Add `PhaseHashes map[string]string` to `StateFile`.
2. Add a canonical hashing helper:
   - stable field order
   - stable map serialization
   - sorted DONs where needed
   - sorted members by `p2p_id`
   - deterministic handling of capability config maps
3. Define exactly which phases use hashing:
   - PreEnvStartup
   - Configure CapReg
   - Configure WorkflowReg
   - Jobs
4. Keep actual-state checks for:
   - contract deployment
   - JD chain configs
5. TOML patching must use direct compare, not hashing.
6. Define `--force` explicitly before implementing it:
   - it bypasses hash-based skip decisions
   - it bypasses TOML no-change skip decisions
   - it reruns configurable phases even if their inputs are unchanged
   - it does not force unsafe contract redeployment when contracts already exist
   - deploy steps still respect actual-state safety checks
7. Add `--force` to the CLI and thread it into reconciliation control flow using the semantics above.
8. Store a phase hash only after the phase succeeds.
9. Do not try to make hashed phases drift-safe. Document that `--force` is the repair path for hashed/configurable phases.

### Tests to add or update

- [internal/domain/state_test.go](/Users/bartektofel/Desktop/repos/chainlink/core/scripts/cre/reconciler/internal/domain/state_test.go)
- new unit tests around canonical hash construction
- reconcile/deployer tests for skip behavior

### Validation

- Run:
  - `go test ./core/scripts/cre/reconciler/internal/domain`
  - `go test ./core/scripts/cre/reconciler/internal/onchain`
  - `go test ./core/scripts/cre/reconciler/...`

### Exit criteria

- A no-change apply performs no on-chain or JD writes.
- A changed capability config reruns only the affected hashed phases.
- `--force` reruns all phases.

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

### Implementation steps

1. Extend `NodeClient` with:
   - `ReadSolanaKeys`
   - `ReadAptosKeys`
2. Extend `NodeRuntimeInfo` with non-EVM key storage.
3. Update discovery to read and persist these keys.
4. Add a dedicated preflight step before on-chain apply.
5. Preflight must validate:
   - chart-side `don-name`
   - JD-side labels `p2p_id`, `environment`, `type`
   - every node in every DON has the required chain config for each chain-scoped capability
6. Keep `stellar` out of scope in this phase.
7. Keep the existing apply-time validation as a backstop.

### Tests to add or update

- discovery tests for Aptos and Solana runtime data
- chart/preflight tests for missing chain config
- tests for missing JD labels

### Validation

- Run:
  - `go test ./core/scripts/cre/reconciler/internal/discovery`
  - `go test ./core/scripts/cre/reconciler/internal/domain`
  - `go test ./core/scripts/cre/reconciler/internal/onchain`
  - `go test ./core/scripts/cre/reconciler/...`

### Exit criteria

- Solana and Aptos keys are available in runtime state.
- Misconfigured nodes fail in preflight before on-chain work starts.
- Missing JD labels produce precise errors.

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

### Implementation steps

1. Update the UI catalog to the finalized support set:
   - add `aptos`
   - keep `solana`
   - do not add `write-solana`
2. Clean `capability_defaults.toml`:
   - keep defaults for real supported capabilities
   - fold any `write-solana` values into the Solana capability handling if they are still needed
3. Extend backend environment construction as needed so the reconciler can actually supply the chain families required by Aptos/Solana feature code.
4. Verify chain-scoped capability handling in `desired.go` still matches the final catalog.

### Tests to add or update

- UI server capability catalog tests
- desired-state tests for chain-scoped capability handling
- any env/topology tests needed once Aptos/Solana are materially wired in

### Validation

- Run:
  - `go test ./core/scripts/cre/reconciler/internal/ui`
  - `go test ./core/scripts/cre/reconciler/internal/domain`
  - `go test ./core/scripts/cre/reconciler/internal/onchain`
  - `go test ./core/scripts/cre/reconciler/...`

### Exit criteria

- The visible capability catalog matches the supported set.
- `aptos` is selectable end-to-end.
- standalone `write-solana` does not appear in the UI.

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

### Implementation steps

1. Add the E6 allowlist input to the DON editor UI:
   - add a free-text UI field backed by a slice of strings
   - serialize user-entered values directly to `registryBasedLaunchAllowlist`
   - preserve order
   - reject empty entries
   - do not reinterpret the entries as addresses or typed objects
2. Change E7 default DON selection:
   - first workflow DON
   - fallback to index 0 if none
3. Move JD connectivity/settings to a JD tab.
4. Remove `infra.chart_values` and `infra.namespace` from the desired-state schema.
5. Remove the corresponding UI inputs.
6. Make `apply` load the chart from `--chart-dir`, same as `serve`.
7. Derive namespace from chart values everywhere.

### Tests to add or update

- desired-state parsing tests
- UI server request/response tests
- UI tests for allowlist round-trip if practical

### Validation

- Run:
  - `go test ./core/scripts/cre/reconciler/internal/ui`
  - `go test ./core/scripts/cre/reconciler/internal/domain`
  - `go test ./core/scripts/cre/reconciler/...`

### Exit criteria

- UI no longer exposes chart path or namespace fields.
- Desired-state schema no longer contains those fields.
- Allowlist can be edited in the UI as capability names.

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
