A test runner harness for the /chainlink repo.

<goals>
- Provide a single, easy command to setup and run tests in /chainlink repo, eliminating `make` command chaining.
- Enable automatically re-running tests and analyzing results to catch and diagnose flakes and slow tests.
- Provide a dual-layered AI skill set:
  1. `backlog-flaky-test-pipeline`: A high-level orchestration skill for JIRA/Trunk.io backlog clearing.
  2. `debug-test-failure`: A diagnostic skill for focused, iterative test fixing (under `tools/test/`).
</goals>

<rules>
- From /chainlink root, document `make new_test`, `make new_gotestsum`, and `make new_test_diagnose`. When working only inside this module, `go run . …` is fine.
- Each output should account for a pretty, human-readable terminal experience, and a minimal version meant for AI ingestion.
- Harness-owned terminal messages go through `internal/output` (`--ai-output` vs human, inline progress policy); child test processes still use raw stdout/stderr passthrough where appropriate.
</rules>

<modes>
<mode name="go test" subcommand="run">
Run tests using vanilla `go test` command and arguments.
</mode>
<mode name="gotestsum" subcommand="gotestsum">
Run tests using gotestsum for those that prefer its output and tools.
</mode>
<mode name="diagnose" subcommand="diagnose">
Opinionated flow to re-run tests and identify flakes, races, timeouts, and test runtimes.
</mode>
</modes>

<orchestration>
### Skill Hierarchy
The harness supports two primary AI workflows:

1. **`backlog-flaky-test-pipeline`** (Global): Use this for "Macro" tasks. It manages JIRA ticket lifecycle, Trunk.io evidence gathering, and multi-agent debate. It often calls the harness's `diagnose` mode as a sub-task.
2. **`debug-test-failure`** (Local): Use this for "Micro" tasks. It focuses on the immediate code logic, local race detection, and iterative fixes for a specific test file.
</orchestration>

<commands>
Run these commands to validate any changes you make:
```sh
golangci-lint run ./... --fix
go test ./...