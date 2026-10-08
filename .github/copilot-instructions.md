# PR Review Instructions

Chainlink node: a Go monorepo of many modules whose code signs transactions, moves funds, and writes on-chain and persisted state.

Every PR runs CI: `go build`, tests, and `golangci-lint` (`.golangci.yml`). Your review covers what CI cannot see: **behaviour**. Spend each comment on a defect that compiles cleanly, passes lint, and still breaks at runtime.

Each comment states a **failure scenario**: the concrete input or state that leads to a wrong outcome, plus a fix. A concern with no failure scenario stays out of the review.

The repo tracks the latest Go release, which is newer than your training data. Treat every syntax form and stdlib API in the diff as valid; the CI compiler decides what builds.

## Where to look hardest

Scale scrutiny to **merge danger**. Read line by line any change with:

- **Wide blast radius**: shared utilities, transaction signing and submission, funds, on-chain or persisted state.
- **Costly revert**: schema migrations, data rewrites, config or protocol changes other services depend on, deploy ordering.
- **High complexity**: concurrency, branching state machines, cross-module contracts, logic without tests.

## What to find

In priority order:

1. **Correctness**: Does the code do what the PR claims, for every input and state it can receive, including edge, empty, and failure cases? Do invariants still hold after every path?
2. **Concurrency**: Is all shared state safely synchronised? Can every goroutine and blocking operation finish or be cancelled?
3. **Lifecycle**: Is everything acquired or started (contexts, goroutines, connections, handles) released or stopped on every path, including errors and shutdown?
4. **Errors**: Is each error handled once, at the layer that can act on it, with enough context for callers to inspect and operators to diagnose?
5. **Security**: Is untrusted input validated at the trust boundary? Do secrets stay out of logs, errors, and responses? Can any check be bypassed?
6. **Performance**: Does cost stay bounded as input grows? Is the hot path free of work that could be batched, cached, or skipped?
7. **Tests**: Is new behaviour tested through its public interface, including error paths? Are tests deterministic, passing by design rather than by timing or ordering luck?
8. **Design**: Does it break contracts other packages or modules rely on? Does it leave structure as good or better than it found it?
   - **Deep modules**: Does each new interface hide more than it costs to learn? Flag pass-through methods and layers that add no abstraction.
   - **Information hiding**: Does each design decision (format, algorithm, storage) live in one place, or must several packages change together?
   - **Pull complexity downward**: Does the callee absorb hard cases, or must every caller handle them? Could the error case be defined out of existence?
   - **DRY**: Is each piece of knowledge represented once? Code that merely looks alike but encodes different knowledge stays apart.
   - **Orthogonality**: Does a change here force changes in unrelated modules? Do callers reach through chains like `a.B().C().D()`?
   - **Good-enough**: Does the code build what the task needs, with generality waiting for a real second use case?

## Owned by CI

Formatting, imports, naming style, comment wording, unused code, loop-variable capture, modern-idiom rewrites, preallocation, `Sprintf` vs concatenation, testify assertion choice, and deprecated APIs are all enforced by `golangci-lint`. Leave them to CI. When a PR is purely lint fixes or mechanical refactors, review it for accidental behaviour changes only.
