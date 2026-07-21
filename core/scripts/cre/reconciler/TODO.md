# griddle reconciler — groomed backlog

Refined from raw notes into scoped work items, each with current status (verified against the code as of the
phase-0–10 + Phase-11-plan state), the problem, concrete scope, acceptance criteria, open decisions, and
dependencies. This is a planning input, not an implementation plan.

Status legend: **DONE** (in tree) · **PARTIAL** (some of it landed) · **TODO** (not started) · **DECISION** (needs a call before planning).

Priority buckets (suggested): **P0** correctness/security now · **P1** reconciliation model · **P2** feature/UX · **P3** polish.

---

## A. Reconciliation & phase model

### A1 — Per-phase input-hash memoization (skip unchanged phases)  ·  TODO · P1
**Not "true reconciliation".** CLDF changesets are executed, not simulated — there's no way to preview the jobspec /
contract calldata a phase *would* apply, so effect-level desired-vs-actual diffing is infeasible for the opaque steps.
Instead: **hash each phase's desired inputs, store the hash after the phase succeeds, and skip the phase next run when
its input hash is unchanged.** This is input memoization, not drift detection ("current" = the inputs at last successful
apply, not live on-chain state).

**Why it works without simulation.** The expensive phases are already **full-replace**, so re-running one is inherently
correct for add/remove — hashing only decides *whether* to re-run:
- **Jobs** = `DeleteAllForDons` (delete everything) → `runPostEnvStartup` over the *current* feature set. Feature
  removed ⇒ its jobs are deleted and not recreated. **No targeted per-feature deletion needed** (drops the previously-
  planned "job deletion by feature").
- **`ConfigureCapabilitiesRegistry` is a setter, not an appender** (confirmed) — re-running with the full desired set
  handles capability/node add *and* remove. So "if anything changed, re-run the whole phase" is correct; no per-
  capability/per-node diff engine.

`state.DesiredHash` (currently declared-but-unused in `internal/domain/state.go`) is replaced by a per-phase hash map.

**Phase skip mechanism (hybrid — prefer actual-state reads where cheap, hash only the opaque phases):**
| Phase | Mechanism | Hash inputs |
|---|---|---|
| Deploy contracts | **actual state** (per-contract address presence, see A3) — correctness + drift-safe | — |
| PreEnvStartup | hash | per-DON caps + `capability_configs` (global+DON) + allowlist + registry chain + **member set by sorted p2p_id** |
| JD chain configs | **actual state** (`ListNodeChainConfigs`, already done) | — |
| Configure CapReg | hash (opaque effect) | **excluding bootstrap-only & gateway DONs** (mirror `flags.HasNoOtherFlags(flags, {GatewayDON, BootstrapDON})`): DON names/types, caps+configs+allowlist+exposesRemoteCaps, **membership by sorted p2p_id**, discovered CSA/EVM/OCR2 per member |
| Resolve DON IDs | always run (read-only, cheap) | — |
| Configure WorkflowReg | hash | workflow DON IDs + workflow owners |
| TOML injection (breakpoint) | **direct compare** (actual is readable — see below) | — |
| Jobs | hash (backed by full delete+recreate) | topology + feature set + gateway service configs |

**TOML phase is special — compare, don't hash.** The current `30-cre` layer is already in the chart YAML, so we can
read it and compare directly to the freshly-generated TOML per node. If identical for all nodes ⇒ **no patch, no
breakpoint** (also catches manual edits to the layer). Only patch + breakpoint the nodes whose generated TOML differs.
Compare the layer's TOML string value (generation is deterministic), not the YAML serialization.

**Membership is keyed on p2p_id, not chart/instance name** (everywhere a "which nodes" set is hashed: PreEnvStartup,
CapReg). p2p_id is the stable on-chain identity; a node rename/redeploy that keeps its p2p_id must NOT trigger a
re-config, and a re-keyed node MUST.

**Scope.**
- Add a per-phase hash map to `StateFile` (e.g. `PhaseHashes map[string]string`), replacing the unused `DesiredHash`.
- A canonical, deterministic hash helper over each phase's input bundle (sorted keys / stable field order — Go map
  iteration is random; see caveat). Sort p2p_id member sets.
- Gate each opaque phase on `hash(inputs) == state.PhaseHashes[phase]`; on mismatch run, then store the new hash
  **only after success**.
- `--force` flag to ignore all hashes and re-run every phase (drift-repair escape hatch).

**Acceptance.** A no-change `apply` does zero on-chain/JD writes, no TOML patch, no breakpoint, and exits clean;
changing one DON's capability config re-runs (only) CapReg + jobs for the changed set; removing a feature re-runs jobs
(delete+recreate) so its jobs disappear; a node rename that preserves p2p_id is a no-op; `--force` re-runs everything.

**Caveats (the things that bite).**
- **Not drift-safe** for the hashed phases — out-of-band changes to a contract/JD won't be detected (unchanged input ⇒
  skip). `--force` and the actual-state phases (deploy, chain configs) are the mitigations.
- **Canonical hashing is the real work** — non-deterministic serialization ⇒ spurious re-runs (annoying) or spurious
  skips (dangerous). This is the main implementation risk, not the concept.
- **Store-after-success only**; discovery must run before hashing (it feeds the opaque-phase hashes) — both already hold.

**Depends on / relates to:** A2 (per-checkpoint hashing), A3 (deploy stays actual-state), B-items (p2p_id + discovered
keys feed the hashes). This supersedes the old "diff-based reconciliation" idea.

### A2 — Granular, resumable phases  ·  TODO · P1
**Problem.** `deployAndConfigureOnChain` (onchain.go) runs P1–P8 as one function. A failure anywhere forces a full
re-run of all steps. There is coarse skip (see A3) but no per-step checkpoint/resume.

**Scope.** Break the on-chain flow into individually-gated, individually-persisted steps (deploy CapReg, deploy
WorkflowReg, deploy forwarders, PreEnvStartup, JD chain configs, configure CapReg, resolve DON IDs, configure
WorkflowReg), each with its own "already done?" guard driven by state/datastore so a mid-flow failure resumes at the
failed step. Sync internal log labels (P1–P6 vs orchestration P1–P8 are desynced) with the real order.

**Acceptance.** Killing the process after step N and re-running resumes at step N+1 without repeating N; each step logs
a label that matches its actual position.

Each granular step is the natural unit for A1's store-after-success hashing (hash + persist per step, not one mega-hash
for the whole on-chain flow), so build A2 with A1 in mind.

**Note.** The `internal/onchain` subpackage extraction was deliberately **deferred** (dead-stub version was deleted).
This item can be done in place in `onchain.go`; revisit the subpackage split only if it earns its keep.

**Depends on / relates to:** A1 (per-step hashing), A3.

### A3 — Auto-skip contract deployment when addresses present  ·  PARTIAL · P1
**Status.** Already partly done: `deployContracts` skips the deploy sequence and hydrates the datastore from state when
`state.HasAddress("CapabilitiesRegistry")` is true. **Gap:** it's coarse (keys off CapReg only) and not per-contract.

**Scope.** Make the skip per-contract via datastore `MightGetAddressFromDataStore` guards (CapReg, WorkflowReg,
KeystoneForwarder independently), so a partial prior deploy resumes cleanly. Fold into A2's per-step guards.

**Acceptance.** With only CapReg in state, re-run deploys WorkflowReg/Forwarder but not CapReg; with all present,
deploys nothing.

**Depends on / relates to:** A2 (this is one of its steps), A1.

---

## B. Discovery & node validation

### B1 — Read Aptos / Solana / Stellar public keys + their OCR2 key bundles  ·  TODO · P2
**Problem.** Discovery only reads EVM addresses (`MustReadETHKeys`) and OCR2 bundles keyed by chain family (currently
only `evm` is used). The `NodeClient` interface (`deps.go`, `internal/discovery/discovery.go`, impl `nodeapi.go`) has
no methods for Aptos/Solana/Stellar keys, and `NodeRuntimeInfo` (`internal/domain/state.go`) has no fields for them.
The reference config lib already supports Solana/Aptos chains (`appendSolanaChain`/`appendAptosChain`); griddle is
EVM-only.

**Scope.** Extend `NodeClient` with `ReadSolanaKeys`/`ReadAptosKeys`/`ReadStellarKeys` (map/address forms) using the
matching `clclient` read methods; extend `NodeRuntimeInfo` with per-family address maps; store per-family OCR2 bundle
IDs (already keyed by family — verify Aptos/Solana families flow through). Wire into topology hydration
(`hydrateDiscoveredEVMAddresses` needs Solana/Aptos siblings).

**Acceptance.** For a DON with a `solana-*`/`aptos-*` capability, discovery captures that chain's node public key +
OCR2 bundle, and CapReg configuration succeeds for non-EVM chains.

**Open decisions.** Confirm Stellar is actually supported by the node API / `clclient` today (may not be — scope it
out if not). Which families are in scope for the first cut.

**Depends on / relates to:** B3 (role gating shares the discovery loop), A1 (keys feed hashes).

### B2 — Validate node labels during discovery/preflight  ·  TODO · P2
**Problem.** Only the `don-name` chart label is validated (`donNameLabel` in chartvalues.go). Nothing checks that JD-
/chart-registered nodes carry the other expected labels: `p2p_id`, `don-<Name>`, `environment`, `type`.

**Scope.** Add a preflight check (in discovery or a dedicated `preflight` step) that reads each node's labels (from JD
`ListNodes` labels and/or chart `registerNodes.labels`) and asserts the expected set is present and consistent with
the resolved role/DON/env; fail with a precise per-node message listing missing/mismatched labels.

**Acceptance.** A node missing `p2p_id` or with a wrong `environment` label is reported by name with the exact missing
label before any on-chain work.

**Open decisions.** Exact label keys/values and their source of truth (JD vs chart). Enumerate the required set.

### B4 — Validate every chain-capability has a matching `[[EVM]]` per member node  ·  PARTIAL · P1
**Status.** The apply-time half is already covered by **C1**: `validateDiscoveredEVMAddresses` (onchain.go:1080) sources
its required-chain set from the model (registry chain + each DON's capability chains) and fails fast if a member node
discovered no EVM address for a required chain — which transitively proves the chart declared no `[[EVM]]` for it.
**Gap:** that check runs late (after discovery, at apply) and is address-shaped, not chart-shaped. There's no early,
authoring-time preflight that reads the chart's per-node `[[EVM]]` blocks directly.

**Problem.** If a DON carries `evm-1337` but one of its member nodes is missing the `[[EVM]]` chain config for 1337 in
the chart, the failure only surfaces mid-apply as a "no EVM address for chain N" error — after JD/on-chain work has
begun for other DONs.

**Scope.** Add a preflight (in discovery/`preflight`, or `/api/preview` server-side) that, for each DON and each of its
chain-scoped capabilities (`evm-<id>`, later `solana-*`/`aptos-*`), asserts **every member node** of that DON has the
matching chain config declared in the chart. Report per-node, per-chain misses before any JD/on-chain write. Keep the
existing discovery-time check as the backstop.

**Acceptance.** A DON with `evm-1337` where one member node lacks `[[EVM]]` for 1337 is reported by node name + chain
ID at preflight, before discovery/on-chain work runs.

**Depends on / relates to:** C1 (DONE; this makes it earlier + chart-shaped), B2 (label validation shares the preflight),
B1 (extend to non-EVM chains once their discovery lands).

### B3 — Don't read OCR2 keys from bootstrap and gateway nodes  ·  TODO · P0 (bug)
**Problem.** `discoverOne` (`internal/discovery/discovery.go`) reads OCR2 bundle IDs for **every** node with no role
gating, and topology hydration (`hydrateDiscoveredOCR2BundleIDs`) errors on "no discovered OCR2 key bundles". Bootstrap
and gateway nodes are not OCR signers and may not have (or need) OCR2 bundles — reading/requiring them is wrong and can
spuriously fail.

**Scope.** Gate the OCR2 read on role: only worker/plugin nodes. Skip for `RoleBootstrap` and `RoleGateway`
(`node.NodeType`). Ensure `hydrateDiscoveredOCR2BundleIDs` is only called for / only requires worker nodes.

**Acceptance.** Discovery on a bootstrap/gateway node performs no OCR2 read and topology build doesn't require OCR2
bundles for them; workers unchanged.

**Depends on / relates to:** B1 (same loop), A1.

---

## C. Node config generation

### C1 — Never emit `[[EVM]]`; require the chain to pre-exist (validate)  ·  DONE · P0
**Status.** Done. Two parts landed together:
1. EVM emission was already **dropped** — `nodeconfig.Generate` emits no `[[EVM]]` (Phase 2, golden-tested).
2. The precondition is now **enforced explicitly**: chains are a first-class part of the desired-state schema
   (`[[chains]]` in `internal/domain/desired.go` — chain ID + `ws_url`/`http_url` + exactly one `registry = true`).
   Nothing is derived from the chart's `anvil.instances` anymore (`ChartAnvilInfo`/`Anvils`/`getAnvilInstances`/
   `GetAnvilByChainID` were removed entirely). `validateDiscoveredEVMAddresses` (onchain.go) sources its required-chain
   set from the model (registry chain + each DON's capability chains) and fails fast with an actionable message
   ("node has no EVM address for chain N — the chart's node config defines no `[[EVM]]` for this chain... add it to
   the chart, roll the nodes, and re-run") before any on-chain work runs. The UI's chain-scoped capability picker
   (`app.js` `toggleCapability`) only allows attaching a chain the user has already declared, so the `evm-1337`-typo
   class of bug (declaring a capability chain that isn't a real, addressed chain) is caught at authoring time, not at
   `PreEnvStartup` time.

**Depends on / relates to:** B1/B3 (discovery), and the documented "chart owns EVM chains" assumption.

### C2 — Add currently-unsupported capabilities  ·  TODO · P2
**Problem.** The set of capabilities a user can actually pick and deploy is narrower than what the tree half-supports —
there's drift between the UI catalog, the default-config file, and end-to-end wiring:
- **UI catalog** (`handleCapabilities`, server.go:327-334) exposes only 8: `cron`, `evm`, `http-action`,
  `http-trigger`, `vault`, `consensus`, `don-time`, `solana`.
- **`capability_defaults.toml` ships default configs for capabilities that are *not* in the catalog** — `aptos`,
  `write-solana`, `mock` — so they have configs but can't be selected in the UI.
- **`aptos` is registered as chain-scoped** (`chainScopedCapabilities`, desired.go:457-459: `evm`/`solana`/`aptos`) yet
  is absent from the catalog, so it can't be attached to a DON at all today.

**Scope.** Close the gap end-to-end for each capability we intend to support:
1. Add it to the UI catalog (name/label/description/`chainScoped`) so it's selectable.
2. Ensure a default config exists (or is intentionally omitted) in `capability_defaults.toml`.
3. Confirm the backend actually wires it into CapReg config + jobs (the reference feature hooks), and — for chain-scoped
   caps — that discovery reads the matching keys (this is **gated on B1** for Solana/Aptos/Stellar).
Start from the concrete known gaps (`aptos`, `write-solana`, `mock`), then enumerate the remaining CRE capabilities the
reconciler should support.

**Acceptance.** Each newly-supported capability is selectable in the UI, has (or intentionally lacks) a default config,
and produces correct CapReg + job output on `apply`; a chain-scoped addition (e.g. `aptos-<id>`) works alongside its
discovery keys.

**Open decisions.** The authoritative list of capabilities the reconciler must support (enumerate against the current
CRE capability set — don't guess). Which are in scope for the first cut vs. deferred. Whether `mock` is a real target or
test-only. Per-capability: chain-scoped or not, gateway-routed or not.

**Depends on / relates to:** B1 (non-EVM discovery — hard prerequisite for Solana/Aptos/Stellar caps), C1 (chain
precondition applies to each new chain-scoped cap), E6 (allowlist), the UI capability catalog + `capability_defaults.toml`.

---

## D. Parallelization

### D1 — Parallelize remaining sequential loops  ·  PARTIAL · P2
**Status.** Job cleanup is already parallel + rate-limited (`jobs.DeleteAllForDons`, Phase 4); discovery is parallel
(Phase 3); JD chain configs are parallel (Phase 6). The cancel→list→delete chain *within a node* is necessarily
sequential (cancel must precede delete) — that's correct, not a gap.

**Remaining candidates (verify each is safe, then bound with errgroup):**
- `runPreEnvStartup` — per-feature × per-DON loop is sequential (onchain.go). Features may have shared state
  (`donsCapabilities`, OCR3 maps) — needs a mutex or per-DON result merge before parallelizing.
- `buildDonsForJobs` — builds `clnode.Output` + enriches CapReg node per node sequentially (postenv/reconcile).
- `buildNodeSpecs` / topology build — `GetNodeSecretsToml` per node sequential (onchain.go `buildTopology`).
- `runPostEnvStartup` — per-feature × per-DON, similar caveats to PreEnvStartup.

**Scope.** For each, confirm no shared-state races, then wrap in a bounded `errgroup` (mirror Phase 3/4/6). Rate-limit
anything hitting JD.

**Acceptance.** Named loops run concurrently under `-race` with no data races; wall-clock improves on multi-node DONs.

**Open decisions.** PreEnvStartup/PostEnvStartup parallelization depends on Feature-hook thread-safety — may need to
stay sequential or aggregate results. Confirm before planning.

---

## E. Gateway & UI

### E1 — One gateway serving multiple DONs, incl. `capabilities`-type DONs  ·  PARTIAL · P2
**Status.** DON type `capabilities` already exists as an option in the UI DON-type selector (`app.js`/`index.html`).
**Gaps:**
1. **Backend:** a gateway currently maps to a single DON (`GatewayNodeAssignment{Node, DON}`, `GatewayDONFor` returns
   one; `GatewayConnectorConfig.DonID` single). Needs to support a gateway serving **multiple** DONs (a set of
   DON IDs per gateway connector).
2. **UI:** `renderGatewayAssignments` offers only **workflow** DONs as assignment targets. Must also allow
   **capabilities**-type DONs, and allow assigning one gateway to **multiple** DONs (multi-select rather than single).

**Scope.** Backend: change gateway→DON to a one-to-many mapping in `desired.toml` schema (`[[gateway_nodes]]` gains a
`dons = [...]` list, or repeatable) + topology/connector build to emit multiple `Gateways[].DonID` / a DON list.
UI: multi-select of eligible DONs (workflow + capabilities) per gateway; render membership read-only as today.
Generated `[Capabilities.GatewayConnector]` must reflect all served DONs.

**Acceptance.** A single gateway node assigned to both a workflow DON and a capabilities DON produces a gateway job +
connector config serving both; UI lets you pick multiple DONs of either type.

**Open decisions.** `desired.toml` shape for many-DON gateway (list field vs repeated entries) — **breaking schema
change, document it.** Whether capabilities DONs need gateway access in all cases or only for gateway-routed caps.

**Depends on / relates to:** gateway job creation, `NeedsGateway`/`NeedsGatewayAccess` logic.

### E2 — Block save/preview when a gateway DON is left unassigned  ·  TODO · P2
**Status.** The gateway→DON assignment **is** used downstream — `renderGatewayAssignments` (app.js) writes
`state.gatewayNodes`, which persists to `desired.toml` `[[gateway_nodes]]` (`GatewayNodeAssignment`, desired.go:171) and
feeds topology/connector build (`GatewayDONFor` desired.go:444, `buildGatewayNodeSets` onchain.go:903). **Gap:** the UI
happily saves/previews with a gateway DON that has no DON selected (empty `don`), producing a gateway that serves
nothing.

**Scope.** On `save()` and `previewTOML()` (and/or server-side in the `/api/desired` + `/api/preview-toml` handlers),
reject when any gateway DON has no assignment; surface a precise "gateway <name> is not assigned to a DON" error and
don't write TOML. Since a gateway may serve **multiple** DONs (see E1), validate "≥1 assignment", not "exactly one".

**Acceptance.** Saving/previewing with an unassigned gateway shows a blocking error naming the gateway; assigning it to
≥1 DON clears the error.

**Depends on / relates to:** E1 (multi-DON gateway — the many-DON shape this validation must tolerate).

### E3 — Apply the desired state from the UI, with live execution logs  ·  TODO · P2
**Problem.** `apply` is CLI-only (`cmd.go`); the UI can only author/preview/save `desired.toml`. Operators must drop to
a terminal to run the reconciler and can't see progress in the UI.

**Scope.** Add an `/api/apply` endpoint that runs the reconciler and **streams execution logs to the UI** (SSE or
websocket — this is a hard requirement, not optional). Drive a new "Execution" view that shows phase progress + log
tail. The existing `--wait-at-breakpoint` mechanism becomes a UI state: at the TOML-injection breakpoint the view
transitions to a **"commit, push & deploy"** step that waits for the operator to confirm they've committed/pushed the
patched chart and rolled the nodes, then resumes the remaining phases on confirm.

**Acceptance.** Clicking Apply in the UI streams live logs; at the breakpoint the UI prompts for commit/push/deploy and
only continues on explicit confirmation; the run completes (or fails) with the full log visible in the UI.

**Open decisions.** Security/authz of triggering `apply` from a web server (it does on-chain writes + JD mutation) —
localhost-only? confirm gate? How the commit/push/deploy step is represented (manual confirm vs. the UI shelling out to
git). Concurrency: refuse a second apply while one is running.

**Depends on / relates to:** A2 (resumable phases make the breakpoint + resume clean), E4 (status/diff view shares the
execution surface), the existing `waitAtBreakpoint`/`confirm` flags.

### E4 — Diff and status screens in the UI  ·  PARTIAL · P2
**Status.** A **Status** tab exists (`loadState` → `/api/state`): shows phase, deployed contract addresses, on-chain DON
IDs, and node list. **Gap:** no **diff** screen (desired vs. last-applied), and status is a flat dump.

**Scope.** Diff view: render desired-state-in-editor vs. the last successfully-applied inputs, per phase — "what would
this apply change?". This is the natural consumer of A1's per-phase input hashes (a phase whose hash matches is shown
as "no change / will skip"; a mismatch shows the changed fields). Enrich Status with per-phase state and last-run
outcome.

**Acceptance.** The diff screen shows, per phase, whether the next apply would re-run or skip it, and highlights which
inputs changed; status shows per-phase last-applied state.

**Depends on / relates to:** A1 (per-phase input-hash memoization — supplies the "changed vs. unchanged" signal), E3
(apply surface).

### E5 — Move JD connectivity + settings to a dedicated JD tab  ·  TODO · P2
**Problem.** JD settings (gRPC endpoint, domain, environment, TLS) and the "Check JD Connectivity" flow live inside the
**Config** tab (index.html:131-169; `checkJD` in app.js). They're conceptually separate from infra config and crowd the
tab (which E8 shrinks further).

**Scope.** Add a new **JD** tab in the nav; move the Job Distributor card (endpoint, domain, environment, TLS toggle,
access-token note, Check button, node-validation results) there. Keep `state.jd` wiring; `syncConfigInputs`/`checkJD`
just read from the new tab's inputs. Route the tab in `switchTab`.

**Acceptance.** JD endpoint/domain/env/TLS and the connectivity+node-validation check all live on their own JD tab; the
Config tab no longer shows JD fields.

**Depends on / relates to:** E8 (Config-tab cleanup — do together).

### E6 — Per-DON `RegistryBasedLaunchAllowlist` input  ·  TODO · P2
**Status.** Backend is fully wired: `don.RegistryBasedAllowlist` flows into CapReg config (onchain.go:865, :989) and
node config (`RegistryBasedLaunchAllowlist`, nodeconfig.go:127); the UI state model already carries
`registryBasedLaunchAllowlist` (app.js) and it round-trips through `/api/desired`. **Gap:** there is **no UI control** to
view/edit it — it's always empty unless hand-written into `desired.toml`.

**Scope.** Add an editable list input per DON in `donColumnHTML` (alongside the "Exposes remote capabilities" checkbox)
that reads/writes `don.registryBasedLaunchAllowlist`. Entries are workflow-owner addresses / IDs (confirm exact form).

**Acceptance.** Adding allowlist entries to a DON in the UI persists to `desired.toml` and appears in the generated
CapReg + node config.

**Open decisions.** Exact entry format/validation (address vs. arbitrary string) and whether it applies to all DON
types or workflow/capabilities only.

### E7 — Capability picker DON selector defaults to first workflow DON  ·  TODO · P3
**Problem.** `renderCapabilityCatalog` (app.js) defaults `state.selectedDON` to `0` — the first DON of *any* type,
which is often a bootstrap/gateway DON that can't carry capabilities.

**Scope.** Default `selectedDON` to the index of the first DON whose `donTypes` includes `workflow` (fall back to `0`
if none). Only on initial default — don't override an explicit user selection.

**Acceptance.** On load, the capability catalog targets the first workflow DON, not a bootstrap/gateway DON.

### E8 — Remove `Chart Values Dir` and `Namespace` from the Config tab  ·  DECISION · P2
**Problem.** The Config tab's `Chart Values Dir` (`Infra.ChartValues`) and `Namespace` (`Infra.Namespace`) fields are
now largely redundant:
- **Namespace** is a fallback only — the real per-node namespace comes from the chart (`GetNodeNamespace`); `serve`
  already prefers `Infra.Namespace` but falls back to the chart (cmd.go:255-261).
- **Chart Values Dir** is still read by the **apply** path (`LoadChartValues(ds.Infra.ChartValues, …)`, cmd.go:205),
  whereas `serve` uses the `--chart-dir` flag. So the field can't just be deleted without redirecting apply.

**Scope.** Drop both inputs from the Config tab; source namespace from the chart everywhere; make `apply` read the
chart dir from a `--chart-dir` flag (align with `serve`) instead of `Infra.ChartValues`; remove `Infra.ChartValues` /
`Infra.Namespace` from the `desired.toml` schema (**breaking change — document in G1**).

**Acceptance.** Config tab no longer shows these fields; `apply` locates the chart via `--chart-dir`; `desired.toml`
no longer needs `[infra] chart_values`/`namespace`.

**Open decisions (the reason this is DECISION, not TODO).** Confirm nothing else reads `Infra.ChartValues`/
`Infra.Namespace`. Decide the apply CLI contract (`--chart-dir` required? default `.`?). Whether to keep `[infra]`
at all once these two fields go (only `type` would remain).

**Depends on / relates to:** E5 (Config tab shrinks — the JD move + this cleanup leave the tab nearly empty; consider
whether Config survives as a tab), G1 (document the breaking schema change).

---

## F. JD chain configs

### F1 — Skip creating a chain config that already exists  ·  DONE (verify) · P3
**Status.** Implemented in Phase 6: `cre.CreateJDChainConfigs` is list-then-skip + duplicate-key tolerant, and the
per-node loop is bounded-parallel + race-guarded (`r.mu` around `state.JDNodeIDs`). **Action:** confirm against a live
JD that a second `apply` creates zero new chain configs (integration check, out of band). No code work expected unless
the live check shows redundant creates; then add the cheap pre-skip (`ListNodeChainConfigs` before building the node).

---

## G. Cross-cutting

### G1 — Documentation  ·  TODO · P2
**Scope (make concrete).**
- README: update flow diagram to the current phase model; document `--deployer-key` and the `GRIDDLE_JD_ACCESS_TOKEN`
  env var (Phase 11); document the "chart owns all `[[EVM]]` chains, must pre-exist" contract (C1); document the
  breaking `desired.toml` changes (no `nodes[]`, no `access_token`, many-DON gateway shape from E1).
- Package doc comments for the new subpackages (`internal/domain|infra|nodeconfig|discovery|jobs|ui`).
- A short "reconciliation model" doc once A1 lands (what triggers what).
**Acceptance.** A new operator can run `apply` and `serve` from the README alone, including required env vars.

### G2 — Naming: `boot` → `bootstrap`, `Cap Configs` → `Capability Configs`  ·  TODO · P3
Two small renames, grouped because they're pure label/value churn — do last to avoid conflicts with in-flight items.

**G2a — Node role `boot` → `bootstrap`.** The role *constant value* is still `RoleBootstrap NodeRole = "boot"`
(desired.go:22). This "boot" string leaks to the UI (`string(n.NodeType)`, server.go:163) where the node badge and
DON-type derivation key on it (`role === 'boot'`, `n.nodeType === 'boot'` — app.js:363,718) and `.badge-boot` CSS.
(The DON-*type* selector option is **already** `bootstrap` — app.js:335 — so only the node-role value/badge lags.)
- Scope: change the constant value to `"bootstrap"`; `parseNodeRole` already accepts both (chartvalues.go:431), keep
  that for back-compat; update app.js badge/`donTypeForNodes` comparisons and rename `.badge-boot` → `.badge-bootstrap`;
  display "bootstrap" on the label. Grep for any other `"boot"` string comparisons before flipping the value.
- Acceptance: a bootstrap node shows a `bootstrap` badge; nothing still compares against `"boot"`.

**G2b — Tab label `Cap Configs` → `Capability Configs`.** The nav button reads "Cap Configs" (index.html:34) while the
card header inside already reads "Capability Configs" (index.html:73). Rename the nav label for consistency. The
`data-tab="configs"` id / `renderCapConfigs` internals can stay as-is (cosmetic label only).
- Acceptance: the tab nav reads "Capability Configs".

### G3 — Produce CLD-compatible audit artifacts (README promise is unimplemented)  ·  TODO · P2
**Problem.** The README advertises audit artifacts (§"Audit artifacts", README.md:208-210): *"per-changeset artifacts
will [be] written in CLD-compatible format to `cre/artifacts/durable_pipelines/` … directory layout and JSON schema
match what the `cld` runtime produces."* This is aspirational — **nothing writes them today** (the only file writes in
the tree are the state file, the chart patch, and `desired.toml`). The README's own "In the future" wording flags the
gap.

**Scope.** After each executed changeset/phase, emit a CLD-compatible artifact into `cre/artifacts/durable_pipelines/`
whose directory layout + JSON schema match the `cld` runtime's output, so the same tooling can consume reconciler runs.
Because CLDF changesets are *executed, not simulated* (see A1), the artifact is a record of what was applied per
changeset. Wire it into the per-step flow (A2) so each step emits its artifact on success.

**Acceptance.** A completed `apply` populates `cre/artifacts/durable_pipelines/` with per-changeset JSON artifacts whose
layout/schema validate against the `cld` runtime's format; the README's "In the future" caveat is removed.

**Open decisions.** The **authoritative CLD schema + directory layout** — pin it against the real `cld` runtime source
(not in this repo per a quick grep; confirm where it lives and version it). Which phases count as a "changeset"
(the opaque CLDF steps vs. every phase). Whether to gate behind a flag (`--artifacts-dir` / on by default). What goes
in each artifact (inputs, addresses, tx hashes, changeset id/hash — reuse A1's per-phase hashes).

**Depends on / relates to:** A1/A2 (per-phase/changeset execution + hashes feed the artifact contents), G1 (update the
README from "in the future" to documenting the real output).

## Housekeeping (discovered during review — not in the original notes)
- **Delete dead root `web/` dir** — `web/{app.js,index.html,styles.css,capability_defaults.toml}` are still git-tracked
  but no longer embedded (live assets moved to `internal/ui/web/`). `git rm -r web/`. · P3
- **Rotate the leaked JD token** — a real `access_token` was committed to a consumer repo's `cre/desired.toml`;
  removing the field (Phase 11) does not un-leak git history. · P0 (ops, not code)

---

## H. Deployer keys & signing

### H1 — Per-chain deployer private keys via `PRIVATE_KEY_<CHAIN_ID>` env vars  ·  TODO · P2
**Problem.** There is a single deployer key for **all** chains: `--deployer-key` (cmd.go:116), defaulting to the Anvil
dev account (`reconcile.go:97`, `blockchain.DefaultAnvilPrivateKey`). It's used for every chain's transactor
(`DeployerTransactorGen: TransactorFromRaw(r.deployerKey)`, onchain.go:247) and for the workflow-owner address on the
registry chain (`deployerAddress(r.deployerKey)`, onchain.go:550). A real multi-chain deployment needs a different
funded key per chain.

**Scope.** Resolve a per-chain key from `PRIVATE_KEY_<CHAIN_ID>` env vars (e.g. `PRIVATE_KEY_1337=0x…`), falling back
to `--deployer-key`, then the Anvil default. Thread a `chainID → key` lookup into the per-chain transactor construction
(onchain.go:247 and the deploy loop at :764-769) instead of the single `r.deployerKey`. Decide what the workflow-owner
address (registry chain, onchain.go:550) resolves to — presumably the registry chain's key.

**Acceptance.** With `PRIVATE_KEY_1337` set, deploys/txs on chain 1337 are signed by that key; chains without a
`PRIVATE_KEY_<id>` fall back to `--deployer-key`/Anvil; the workflow owner is derived from the registry chain's key.

**Open decisions.** Precedence order (env vs. flag) — env should win per-chain, flag as global fallback. Whether to
keep `--deployer-key` at all or make it purely the fallback. Whether keys ever come from a file/secret store rather
than env (out of scope for the first cut). Never log the keys.

**Depends on / relates to:** the on-chain deploy flow (A2/A3), G1 (document the env-var convention alongside
`GRIDDLE_JD_ACCESS_TOKEN`).

---

## Suggested sequencing for the eventual implementation plan
1. **P0 correctness first:** B3 (OCR2 role gating), C1 (EVM precondition validation), Phase 11 (token env-only), token rotation.
2. **Reconciliation core:** A3 → A2 → A1 (they build on each other). A1 is per-phase input-hash memoization backed by
   the existing full-replace phases — no targeted job-deletion or diff engine required.
3. **Discovery breadth:** B1 (non-EVM keys), B2 (label validation), B4 (chart-level chain-capability preflight).
   Then C2 (add unsupported capabilities) once B1 unblocks the non-EVM ones.
4. **UI quick wins (cheap, mostly independent):** E6 (allowlist input), E7 (workflow-DON default), E5 (JD tab),
   E8 (Config-tab cleanup, needs a decision), E2 (block unassigned-gateway save).
5. **Feature/UX:** E1 (multi-DON gateway + capabilities type), H1 (per-chain deployer keys), remaining D1
   parallelization.
6. **Bigger UI:** E3 (apply-from-UI + live logs) → E4 (diff/status, needs A1).
7. **Polish:** G1 (docs), G2 (renames), G3 (CLD audit artifacts — after the A phase/changeset model lands), housekeeping.

## Cross-item dependency map (quick reference)
- A1 (per-phase input-hash) builds on A2 (per-step checkpoints) + A3 (deploy stays actual-state) + B (p2p_id +
  discovered keys feed the hashes). No dependency on targeted job-deletion — the full-replace jobs phase handles removal.
- B1 and B3 share the discovery loop — do together. B4 shares the preflight with B2.
- C1 depends on discovery (B) producing/validating ETH addresses; B4 makes C1's check earlier + chart-shaped.
- C2 (add unsupported caps) is gated on B1 for non-EVM caps; it spans the UI catalog + `capability_defaults.toml` +
  backend wiring, so touch all three per capability.
- E1 is mostly independent (schema + UI + connector build). E2 depends on E1's many-DON shape.
- E3 (apply-from-UI) relates to A2 (resumable breakpoint); E4 (diff) consumes A1's per-phase hashes.
- E5 (JD tab) + E8 (Config-tab cleanup) touch the same tabs — do together; E8 is a breaking `desired.toml` change (G1).
- E6/E7 are self-contained UI changes (backend already supports the allowlist).
- H1 (per-chain keys) touches the on-chain deploy flow (A2/A3); document the env convention in G1.
- G2 (renames) should be last to avoid churn across other in-flight items.
- G3 (CLD audit artifacts) consumes A1/A2 (per-changeset execution + hashes) and closes the README's "in the future"
  promise (G1) — build it after the A phase/changeset model exists.
