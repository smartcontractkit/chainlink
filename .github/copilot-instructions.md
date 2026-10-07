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

In priority order:

1. **Correctness**: wrong conditions, off-by-one, nil and zero-value handling, broken invariants, unhandled state transitions, behaviour that contradicts the PR description. `*big.Int` aliasing (mutating a shared pointer), integer overflow on amounts, unit mismatches (wei/gwei, seconds/ms).
2. **Concurrency**: data races on maps, slices, and struct fields; goroutines with no exit path on context cancel or channel close; deadlocks; locks held across I/O or RPC; sends that can block forever; `WaitGroup` or `Once` misuse.
3. **Lifecycle**: `context.Context` replaced with `context.Background()` mid-chain; missing `cancel()`; `Close` that leaves goroutines started by `Start` running; rows, bodies, files, tickers, and subscriptions released on the happy path only.
4. **Errors**: an error both logged and returned (handle it once, at the layer that can act); wrapping with `%v` where callers use `errors.Is`/`errors.As`; panics reachable from external input; retries that never stop.
5. **Security**: untrusted input reaching SQL, file paths, shell, or unbounded allocation; secrets or keys in logs or errors; bypassable signature, auth, or chain-ID checks; reorg and nonce handling.
6. **Performance** on hot paths: a DB query or RPC per item (N+1), unbounded growth of maps, channels, or caches, allocation inside loops over unbounded input, repeated recomputation of the same result.
7. **Tests**: new behaviour with no test; tests that pass by coincidence (`time.Sleep` synchronisation, ordering luck, shared global state); error paths and edge cases left uncovered; mocks where an in-memory fake would run real logic.
8. **Design**: breaking changes to exported APIs other modules import; one decision leaking across several packages; pass-through layers that add no abstraction.

## Comment format

Each comment names the defect, the line, and a **failure scenario**: concrete input or state leading to a wrong outcome, plus a fix. A concern with no failure scenario belongs in the summary as a question, or nowhere.

## Owned by CI

Formatting, imports, naming style, comment wording, unused code, loop-variable capture, modern-idiom rewrites, preallocation, `Sprintf` vs concatenation, testify assertion choice, and deprecated APIs are all enforced by `golangci-lint`. Leave them to CI. When a PR is purely lint fixes or mechanical refactors, review it for accidental behaviour changes only.
