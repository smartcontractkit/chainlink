A test runner harness for the /chainlink repo.

<goals>
- Provide a single, easy command to setup and run tests in /chainlink repo, eliminating `make` command chaining.
- Enable automatically re-running tests and analyzing results to catch and diagnose flakes and slow tests.
- Provide an AI skill: `fix-flaky-tests`: A diagnostic skill for focused, iterative test fixing (under `tools/test/`) capable of analyzing GitHub Actions logs and managing JIRA tickets' lifecycle.
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

<commands>
Run these commands to validate any changes you make:
```sh
golangci-lint run ./... --fix
go test ./...