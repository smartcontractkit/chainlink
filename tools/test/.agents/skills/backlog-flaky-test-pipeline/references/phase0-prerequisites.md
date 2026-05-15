---
phase: phase0
model_tier: lightweight
---

<phase id="phase0">

<purpose>
Check each required tool before proceeding. Stop with setup instructions if any hard requirement fails. On success, write `phase_outputs.phase0` with `nav_tool`, `lsp_available`, and `golangci_lint_available`.
</purpose>

<pre-step id="mode-pre-detect">
Scan the raw invocation args for the literal token `--local`. If found:
- Set `invocation_mode_hint = "local"`.
- Skip the `trunk-mcp` and `atlassian-mcp` checks below entirely (JIRA and Trunk are not needed).
- Still run the `golangci-lint` and `code-navigation` checks — local mode requires both.

Phase 1 performs the authoritative mode parse; this pre-detect only prevents hard-stopping on missing MCP servers when they are irrelevant.
</pre-step>

<checks>

<check id="trunk-mcp" required="true">
Call `mcp__trunk__search-test` with the actual repo name (e.g. `smartcontractkit/chainlink`) and `testNameSearch="probe"`.

**Fail immediately if `repoName` is the literal string `owner/repo`** — Trunk returns empty results silently for unknown repos.

If the call fails → stop:
```
Trunk MCP is not configured. See: https://docs.trunk.io/flaky-tests/use-mcp-server/configuration
```
</check>

<check id="atlassian-mcp" required="true">
Call `mcp__atlassian__atlassianUserInfo`. Cache `accountId` for later phases.

If it fails → stop:
```
Atlassian MCP is not configured. See: https://support.atlassian.com/atlassian-rovo-mcp-server/docs/getting-started-with-the-atlassian-remote-mcp-server/
(use v2, not v1)
```
</check>

<check id="golangci-lint" required="true">
Run `which golangci-lint`.

If found → set `golangci_lint_available = true` and continue.

If not found → stop:
- If `asdf` is available: `asdf install golangci-lint`
- Otherwise: read the version from `.tool-versions` (or project docs) and run `go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@<version>`
</check>

<check id="diagnose-tool" required="true">
Run `go -C tools/test run . diagnose -h` (exit code 0 expected; help text on stdout).

If it fails → stop:
```
The chainlink `diagnose` tool is not available at tools/test/. Phase 4 relies on it to
verify fixes (10x iteration runner with AI-readable output). Verify you are running
this skill from the chainlink repo root and that tools/test/ contains the diagnose
command.
```
</check>

<check id="code-navigation" required="true">
At least one code navigation tool must work.

**PRE-STEP (mandatory — never skip):** Call `ToolSearch` with `query="select:LSP"` to load the LSP tool schema before attempting step a or b.

**a. LSP for Go** (preferred — always attempt first):
Attempt an `LSP` hover call on any known `.go` source file.

- Works → set `nav_tool = "lsp"`, `lsp_available = true`. Skip step b.
- Fails → run `which gopls`:
  - **gopls missing**: print the following, then proceed to step b (do not stop):
    ```
    gopls is not installed. To enable LSP for Go:

    1. Install gopls:
         go install golang.org/x/tools/gopls@latest
       If using asdf:  asdf reshim golang

    2. Then configure Claude Code with one of:
         (a) Run /lsp in this session to connect the language server.
         (b) Install the official gopls plugin: https://claude.com/plugins/gopls-lsp
    ```
  - **gopls present but LSP still fails**: check whether `~/.claude/plugins/cache/claude-plugins-official/gopls-lsp/` exists.
    - Path **does not exist** (plugin not installed):
      ```
      gopls is installed but LSP is not configured. Options:
        (a) Install the gopls plugin (recommended): https://claude.com/plugins/gopls-lsp
        (b) Run /lsp in this session to connect the language server manually.
      ```
    - Path **exists** (plugin installed but failed):
      ```
      The gopls plugin is installed but LSP failed to connect.
      Try restarting Claude Code, or run /lsp to reconnect the language server.
      ```
    Proceed to step b in either sub-case.

**b. code-review-graph** (fallback):
Call `mcp__code-review-graph__get_minimal_context_tool` with `task="probe"`.

- Works → set `nav_tool = "crg"`.
- Fails → hard stop:
  ```
  Neither LSP nor code-review-graph is available. Without a code navigation tool,
  analyzing tests requires scanning the entire codebase and will consume an excessive
  number of tokens. Please install one of the following before continuing:
    - LSP: go install golang.org/x/tools/gopls@latest  (then /lsp)
    - code-review-graph: see project CLAUDE.md for setup instructions
  ```
</check>

<check id="lsp-supplementary" required="false">
Only run if `nav_tool = "crg"`. Attempt an LSP hover call.

- Works → set `lsp_available = true`.
- Fails → set `lsp_available = false`.

Does not block progress.
</check>

</checks>

<on_complete>
Write to `phase_outputs.phase0`: `{ nav_tool, lsp_available, golangci_lint_available }`.

Announce: "Prerequisites verified. All required tools are available."

Read [phase1-input.md](phase1-input.md) and follow its instructions.
</on_complete>

</phase>
