# AGENTS.md - Development Rules for `tools/ci`

This module centralizes CI operations. All changes must adhere strictly to these rules:

1. **Logic in `internal/`**:
   - All domain logic lives in `internal/<package>`, taking explicit structs and returning `(T, error)`. Filesystem reads are allowed; avoid reading env vars, parsing CLI flags, or writing directly to stdout/stderr.
   - Never call `os.Getenv`, read flags, or write to standard streams directly inside `internal/` (except `internal/ghaction`).
2. **`cmd/` is a thin shell**:
   - Parse flags, resolve environment fallbacks, call `internal/`, and render outputs.
3. **Dual output**:
   - All data commands must support `--json` output alongside human-readable output.
   - When `$GITHUB_OUTPUT` is set, write output parameters using `internal/ghaction`. When unset, write `k=v` or structured output to stdout for local execution.
4. **No `${{ }}` interpolation**:
   - Workflows must pass inputs via `env:` or flags, never string-interpolated into `run:` scripts.
5. **Strict TDD**:
   - Write failing tests first before adding commands or domain logic.
   - Golden-file tests for JSON/matrix generation.
   - Fake/HTTP mock servers for GitHub API interactions (`httptest.NewServer`).
