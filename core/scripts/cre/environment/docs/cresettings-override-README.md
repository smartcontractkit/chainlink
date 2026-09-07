# CRE Settings Override (Local CRE e2e tests)

A test helper for **overriding CRE settings inside a Local CRE e2e test at runtime** —
without tearing down and restarting the topology — with **automatic revert** when the
test ends.

- Helper: [`system-tests/tests/test-helpers/cresettings_override.go`](../../../../system-tests/tests/test-helpers/cresettings_override.go)
- Example tests: [`system-tests/tests/smoke/cre/cresettings_override_test.go`](../../../../system-tests/tests/smoke/cre/cresettings_override_test.go)

## TL;DR

```go
func Test_CRE_MyThing(t *testing.T) {
    testEnv := t_helpers.SetupTestEnvironmentWithPerTestKeys(t, t_helpers.GetDefaultTestConfig(t))

    // Scope the override with an option — here to THIS test's workflow, so it can't affect
    // any other test. It auto-reverts when the test ends. myWorkflowID is hex, no 0x.
    t_helpers.ApplyCRESettings(t, testEnv,
        t_helpers.Workflow(myWorkflowID, "[PerWorkflow.HTTPAction]\nCallLimit = '1'"),
    )

    // ... trigger the workflow and assert its behaviour under the override ...
}
```

That's the whole happy path: one call, and cleanup is automatic.

Prefer scoping to your own **workflow** (or **owner**) as above — it isolates the change
to your workflow. Use `Global` only when you truly mean an environment-wide change (see
[Choosing a scope](#choosing-a-scope--isolation)).

## When to use it

Use it when a test needs a settings value different from the topology default — a rate
limit, a size bound, a feature gate, a concurrency limit — and you don't want a whole
separate topology just to bake that value into `CL_CRE_SETTINGS`.

Don't reach for it to change things that aren't CRE settings (node TOML, capability
config, chain config); those aren't in the settings schema and this won't touch them.

## How it mirrors what we do in CLD

This is the in-test analogue of the production settings rollout in `chainlink-deployments`
(the `cre-limit-change` flow). Same mechanism, smaller blast radius:

| Production rollout (CLD) | This helper (in a test) |
|---|---|
| Edit `settings/*.toml`, regenerate the compiled `settings.toml` | Pass the same settings as scope options — `Global`/`Org`/`Owner`/`Workflow` (or `FromFS` for a file tree), combined by the same `CombineCRESettingsFiles` |
| Durable pipeline emits `job_propose_arbitrary` → `jobName: CRESettings`, `template: cre-settings`, `donName: all-nodes` | Calls the **same** `ProposeJobSpec{Template: CRESettings}` changeset, targeting every node of every DON |
| Sign + execute the proposal | Auto-approves the proposal on each node |
| Nodes apply it live via `loop.AtomicSettings.Store` (no restart) | Same — no restart |
| Change persists until the next rollout | **Auto-reverts** to the pre-test baseline on cleanup |

Same job type, same all-nodes delivery, same live-apply path — just scoped to one test
and reverted afterwards. If you understand the CLD settings rollout, you understand this.

## The API

```go
// Apply one or more scope options to every DON, register auto-revert, and return a handle.
func ApplyCRESettings(t *testing.T, env *TestEnvironment, opts ...Option) *CRESettingsHandle

// Scope options — each takes that scope's settings TOML (no scope prefix):
func Global(toml string) Option        // [global]        — every org/owner/workflow
func Org(id, toml string) Option       // [org.<id>]      — id: org id
func Owner(id, toml string) Option     // [owner.<id>]    — id: workflow-owner hex (no 0x)
func Workflow(id, toml string) Option  // [workflow.<id>] — id: workflow id hex (no 0x)
func FromFS(fsys fs.FS) Option         // a whole prod-style file tree (os.DirFS / embed.FS)

func (h *CRESettingsHandle) Reset(t *testing.T)             // revert now (also happens automatically)
func (h *CRESettingsHandle) AppliedTOML(donName string) string  // the document that was applied
func (h *CRESettingsHandle) BaselineTOML(donName string) string // the document it reverts to
```

- Each option's string is just that scope's settings — **no scope prefix** (the scope and id
  come from the option), identical to a prod `settings/*.toml`. Group related settings under a
  table, e.g. `[PerWorkflow.HTTPAction]` then `CallLimit = '9'`.
- **Setting values are quoted strings** (`CallLimit = '9'`, `Enabled = 'true'`).
- **Pass every scope in ONE call.** Each `ApplyCRESettings` call rebuilds from the boot
  baseline and *replaces* any prior override rather than stacking on it, so a second apply
  while one is still active is **refused with an error** (it would silently drop the first).
  To swap settings within a test, call `Reset` first, then apply again.
- Under the hood the options build a prod-style file tree combined by the deployment
  `CombineCRESettingsFiles` — the same code that builds the compiled `settings.toml`.
- The canonical list of settings and their formats lives in
  `chainlink-common/pkg/settings/cresettings/defaults.toml`.

## Choosing a scope & isolation

Settings resolve **most-specific-first**. For a given workflow a node looks up:

```
workflow.<id>  →  owner.<id>  →  org.<id>  →  global  →  compiled default
```

**Prefer the narrowest scope.** Because the environment is shared across the suite, the
scope you pick is also your **blast radius**:

- **Workflow / Owner** — the override only resolves for *your* workflow (or owner). Even if
  something went wrong, no other test's workflow is affected. **This is the default you
  should reach for.**
- **Org** — affects every workflow in that org.
- **Global** — affects **every workflow in the environment**. It's the deliberate
  environment-wide escape hatch; use it only when you actually mean that.

> **Where do the IDs come from?** From your deployment step: the org id your workflow runs
> under, and the workflow-owner / workflow-id hex (**without** the `0x` prefix) you get when
> you register/deploy the workflow. `Global` needs no id.

### Examples

Each option's string is plain settings TOML — no scope prefix — exactly like a
`settings/*.toml` file (leading indentation is fine; TOML ignores it). For an on-disk
fixture instead of inline strings, use `t_helpers.FromFS(os.DirFS("testdata/..."))`.

**Workflow — one workflow only (preferred)**

```go
t_helpers.ApplyCRESettings(t, testEnv, t_helpers.Workflow(workflowID, `
[PerWorkflow.HTTPAction]
CallLimit = '9'
[PerWorkflow.HTTPTrigger]
RateLimit = 'every5s:2'`)) // workflowID has NO 0x prefix
```

**Owner — one workflow owner only**

```go
t_helpers.ApplyCRESettings(t, testEnv, // ownerHex has NO 0x prefix
    t_helpers.Owner(ownerHex, "[PerOwner]\nWorkflowExecutionConcurrencyLimit = '5'"),
)
```

**Org — one org only**

```go
t_helpers.ApplyCRESettings(t, testEnv, t_helpers.Org(orgID, `
[PerOrg]
BaseTriggerRetransmitEnabled = 'false'
WorkflowExecutionConcurrencyLimit = '42'`))
```

**Global — the whole environment (escape hatch)**

```go
t_helpers.ApplyCRESettings(t, testEnv, t_helpers.Global(`
[PerWorkflow.HTTPAction]
CallLimit = '7'
[PerWorkflow]
ExecutionConcurrencyLimit = '3'`))
```

**Several scopes at once — compose the options in one call**

```go
t_helpers.ApplyCRESettings(t, testEnv,
    t_helpers.Global("[PerWorkflow.HTTPAction]\nCallLimit = '11'"),
    t_helpers.Org(orgID, "[PerOrg]\nBaseTriggerRetransmitEnabled = 'false'"),
    t_helpers.Workflow(workflowID, "[PerWorkflow]\nExecutionConcurrencyLimit = '1'"),
)
```

All scopes merge into a single document (layered on each DON's baseline) and are delivered
together. To see exactly what got produced, log `h.AppliedTOML("workflow")`.

## When does the change take effect? (live vs. restart)

Every setting is read through the same live `AtomicSettings`, but *when* a consumer reads
it differs — this is the single most important thing to get right:

| Class | Effect of a mid-test override | Examples |
|---|---|---|
| **Immediate** (read per operation) | Next call sees it | gates (`RemoteExecutableWorkflowDONBindingEnabled`, `ExecutionTimestampsEnabled`, `ChainAllowed`), bounds/limits (`HTTPAction.CallLimit`, `ExecutionConcurrencyLimit`, `ChainRead.CallLimit`), timeouts |
| **~5s lag** (poller-resized) | Applies within a few seconds | rate limits (`HTTPTrigger.RateLimit`, gateway rates), queue caps |
| **Registration-time** | Only affects workflows registered *after* the override | trigger subscription/registration limits, WASM size checks, workflow-admission limits |

**Rule of thumb:**
- Per-execution limits/gates → `ApplyCRESettings`, then trigger the workflow.
- Registration-time settings → **`ApplyCRESettings` *before* you deploy/register the
  workflow.** Applying them afterwards leaves the already-registered workflow on the old
  value, and it will look like nothing happened.

## It applies to every DON (you don't choose)

Overrides are delivered to **every node of every DON**, mirroring CLD's all-nodes rollout.
This is deliberate: a setting may be enforced on the workflow nodes, the capabilities
nodes, *or* the gateway nodes, so a partial rollout could silently fail to take effect. You
say *what* to change; the helper makes sure it lands everywhere it could be read.
(Internally it still merges each DON's own boot baseline, since DONs can boot with
different `CL_CRE_SETTINGS` — but that's handled for you.)

## Run these serially — the override guard

Because overrides mutate settings on the shared environment, **only one override may be
active at a time**, and settings-override tests must run **serially**. Don't call
`t.Parallel()` in them, and don't add them to the `CRE_TEST_PARALLEL_ENABLED` set. (The CRE
suite already runs environment-using scenarios serially by default, so this is the norm,
not a special case.)

The helper enforces it. If a second override starts while another is still active, it
**fails fast** with an actionable message:

```
a CRE settings override from "Test_CRE_Other" is still active on the shared environment.
Settings-override tests mutate shared state and must run serially: remove t.Parallel()
from this test (and do not add it to the CRE_TEST_PARALLEL_ENABLED set). ...
```

**If you see this:** make the failing test serial — remove its `t.Parallel()` and keep it
out of the parallel set. To change settings again *within one test*, call `h.Reset(t)`
before the next `ApplyCRESettings`; a second apply while an override is still active is
refused (with a message telling you to `Reset` first), since it would silently replace the
first override.

The guard coordinates override-vs-override. What keeps an *unrelated* concurrent test from
seeing your change is running serially (the default) **plus** scoping narrowly — see
[Choosing a scope & isolation](#choosing-a-scope--isolation).

## How cleanup works

- `ApplyCRESettings` captures each DON's **boot baseline** (its `CL_CRE_SETTINGS`) up front
  and registers a `t.Cleanup` that re-applies it. **You don't have to do anything** — the
  settings revert when the test ends, whether it passes or fails.
- Deleting the settings job does **not** revert (a node keeps the last values it stored), so
  cleanup works by **re-applying the baseline**, not by removing the job. The helper does
  this for you.
- Overrides are **not additive**: each delivery fully replaces the getter, so the helper
  always sends *baseline ⊕ your overrides*. Two consequences worth knowing:
  - Omitting a key is fine — it just falls back to its compiled default.
  - The DON's boot settings are always preserved, because they're part of the baseline.

### Reverting early — `Handle.Reset(t)`

Call `Reset` to revert **inside the test body** — e.g. to assert behaviour before *and*
after, or to A/B two settings in one test:

```go
h := t_helpers.ApplyCRESettings(t, testEnv,
    t_helpers.Global("[PerWorkflow.HTTPAction]\nCallLimit = '1'"),
)
// ... assert the tight limit is enforced ...

h.Reset(t) // back to baseline now
// ... assert normal behaviour is restored ...
```

`Reset` is idempotent and turns the automatic cleanup into a no-op, so it's always safe to
call (and safe to *not* call — cleanup still runs).

## Advanced usage

- **Inspect the exact documents:** `h.AppliedTOML(donName)` / `h.BaselineTOML(donName)`
  return the TOML that was applied / will be restored (`donName` is e.g. `"workflow"`).
  Handy for `t.Logf` and debugging.
- **A/B within one test:** apply → assert → `Reset` → apply a different value → assert.
  Because each apply layers on the *baseline* (not the previous override), values never
  accumulate.
- **Combine scopes** freely in one call; they merge into one document.

## Mistakes & failure modes

The helper is fail-loud for anything it can detect, and it validates your overrides against
the schema **before** delivering anything.

| Mistake | What happens |
|---|---|
| Unknown / misspelled setting key | **Fails the test** at `ApplyCRESettings` (`unknown fields …`). Nothing is delivered. |
| Value in the wrong format for the setting (`"banana"` for an int limit, a malformed rate, …) | **Fails the test** at `ApplyCRESettings` (`invalid toml settings …`) — every value is parsed against its setting type. Nothing is delivered. |
| Non-string value | Impossible — the API only accepts `string` values. |
| `0x`-prefixed org / owner / workflow id | **Fails the test** at `ApplyCRESettings` — ids must be given without the `0x` prefix. |
| Environment not running / a node unreachable / proposal rejected | **Fails the test** at `ApplyCRESettings` — propose/approve errors are surfaced. |
| Revert fails during cleanup | **Surfaced**, not swallowed — `t.Error` from the automatic cleanup, or a hard failure from an explicit `Reset`. |
| Best-effort convergence log can't read container logs | **Ignored** — it's only a visibility aid; the authoritative success signal is that every node approved the job. |
| Flipping a **registration-time** setting *after* the workflow is registered | **Silently no-op** for that workflow (see the live-vs-restart table). Apply it *before* deploying the workflow. |
| Overriding a scope nothing matches (e.g. an org id no running workflow uses) | Delivered fine, simply never consulted. Not an error. |
| Two override tests running **concurrently** on the shared environment | **Fails fast** at `ApplyCRESettings` — the second override is refused while the first is active, with an actionable message. Keep override tests serial; do **not** call `t.Parallel()`. To re-apply within one test, `Reset` first — a second apply while one is active is refused. |
| An *unrelated* concurrent test caught by a `Global` override's blast radius | **Not caught by the guard** — it only coordinates override-vs-override. Prevented by running serially (the CRE default) and by scoping narrowly (Workflow/Owner). |

**Short version:** authoring mistakes (bad keys / values / ids) and overlapping override
tests **fail loudly**. The only quiet risks left are timing ones — flipping a
registration-time setting too late, or a `Global` override touching an unrelated
concurrent test — and both are avoided by scoping narrowly and running serially (the
default).

## Requirements & running

These are e2e tests and need a running Local CRE environment (see
`docs/local-cre/system-tests/running-tests.md`). With the environment up:

```bash
go test ./system-tests/tests/smoke/cre -run '^Test_CRE_CRESettings_' -timeout 20m -v
```

## Under the hood (pointers)

You don't need any of this to use the helper, but if you want to trace it:

- Live apply on the node: `core/services/cresettings/delegate.go` → `loop.AtomicSettings.Store` (`chainlink-common/pkg/loop/settings.go`).
- Delivery changeset: `deployment/cre/jobs` (`ProposeJobSpec{Template: CRESettings}`); validation: `deployment/cre/jobs/settings.go` (`VerifyCRESettings`).
- Settings schema, scopes and resolution: `chainlink-common/pkg/settings/cresettings` and `chainlink-common/pkg/settings`.
