# PR Review Instructions

Every PR runs CI: `go build`, tests, and `golangci-lint` (`.golangci.yml`). Your review covers what CI cannot see: **behaviour**. Spend each comment on a defect that compiles cleanly, passes lint, and still breaks at runtime.

The repo tracks the latest Go release (`go` directive in the nearest `go.mod`), which is newer than your training data. Treat every syntax form and stdlib API in the diff as valid; the CI compiler decides what builds.

## Merge danger

Open the review summary with a table rating the PR on three axes. Scale: 🟢 trivial/small · 🟡 minor/medium · 🟠 large · 🔴 major/extreme. Each rationale names a concrete reason (callers, data, deploy order, on-chain state), never a restated rating.

| Axis | Measures |
| --- | --- |
| **Blast radius** | What breaks if this is wrong: callers, services, nodes, funds, on-chain or persisted state. |
| **Revert** | Cost of undoing it: schema migrations, data rewrites, config or protocol changes other services depend on, deploy ordering. |
| **Complexity** | Effort to verify correctness: concurrency, branching state machines, cross-module contracts, logic without tests. |

Then list the code blocks that need **scrupulous human review**, worst first.

## What to find

Run each check against every changed file, in priority order:

1. **Correctness**: Does the code do what the PR claims, for every input and state it can receive, including edge, empty, and failure cases? Do invariants still hold after every path?
2. **Concurrency**: Is all shared state safely synchronised? Can every goroutine and blocking operation finish or be cancelled?
3. **Lifecycle**: Is everything acquired or started (contexts, goroutines, connections, handles) released or stopped on every path, including errors and shutdown?
4. **Errors**: Is each error handled once, at the layer that can act on it, with enough context for callers to inspect and operators to diagnose?
5. **Security**: Is untrusted input validated at the trust boundary? Do secrets stay out of logs, errors, and responses? Can any check be bypassed?
6. **Performance**: Does cost stay bounded as input grows? Is the hot path free of work that could be batched, cached, or skipped?
7. **Tests**: Is new behaviour tested through its public interface, including error paths? Are tests deterministic, passing by design rather than by timing or ordering luck?
8. **Design**: Does the change keep each decision in one place and each interface small? Does it break contracts other packages or modules rely on?

## Comment format

Each comment names the defect, the line, and a **failure scenario**: concrete input or state leading to a wrong outcome, plus a fix. A concern with no failure scenario belongs in the summary as a question, or nowhere.

## Owned by CI

Formatting, imports, naming style, comment wording, unused code, loop-variable capture, modern-idiom rewrites, preallocation, `Sprintf` vs concatenation, testify assertion choice, and deprecated APIs are all enforced by `golangci-lint`. Leave them to CI. When a PR is purely lint fixes or mechanical refactors, review it for accidental behaviour changes only.
