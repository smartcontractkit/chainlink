A test runner harness for the /chainlink repo.

<goals>
- Provide a single, easy command to setup and run tests in /chainlink repo, eliminating `make` command chaining.
- Enable automatically re-running tests and analyzing results to catch and diagnose flakes and slow tests
- Provide an AI skill for the process in `.agents/skills/fix-chainlink-tests/SKILL.md` (under `tools/test/`)
</goals>

<rules>
- From /chainlink root, document `make new_test`, `make new_gotestsum`, and `make new_test_diagnose`. When working only inside this module, `go run . …` is fine.
- Each output should account for a pretty, human-readable terminal experience, and a minimal version meant for AI ingestion.
- Harness-owned terminal messages go through `internal/output` (`--ai-output` vs human, inline progress policy); child test processes still use raw stdout/stderr passthrough where appropriate.
</rules>

<modes>
<mode name="go test" subcommand="run"> 
Run tests using vanilla `go test` command and arguments
</mode>
<mode name="gotestsum" subcommand="gotestsum"> 
Run tests using gotestsum for those that prefer its output and tools
</mode>
<mode name="diagnose" subcommand="diagnose"> 
Opinionated flow to re-run tests and identify flakes, races, timeouts, and test runtimes.
</mode>
</modes>

<permissions>
This tool and its tests frequently write to the repository root (e.g., `diagnose-*` directories) and the `deployment/` cache directory.
Agents (Gemini CLI, Cursor, Claude Code, Codex) MUST have read/write access to the entire repository root, not just this `tools/test/` directory.
- **Gemini CLI:** Ensure you select "Trust folder" to enable the expanded sandbox in `.gemini/settings.json`.
- **Other Agents:** If you encounter "operation not permitted" or scoping errors, advise the user to open the repository root as the workspace.
</permissions>

<commands>
Run these commands to validate any changes you make
```sh
golangci-lint run ./... --fix
go test ./...
```

DO NOT use other commands like `goimports`, `gofmt`, or `go vet` for formatting and lint checks.
</commands>
