A test runner harness for the /chainlink repo.

<goals>
- Provide a single, easy command to setup and run tests in /chainlink repo, eliminating `make` command chaining.
- Enable automatically re-running tests and analyzing results to catch and diagnose flakes and slow tests
</goals>

<rules>
- From /chainlink root, only document `go -C ./tools/test run . …` (never `go run ./tools/test` from the parent module).
</rules>

<modes>
<mode name="go test" subcommand="test"> 
Run tests using vanilla `go test` command and arguments
</mode>
<mode name="gotestsum" subcommand="gotestsum"> 
Run tests using gotestsum for those that prefer its output and tools
</mode>
<mode name="survey" subcommand="survey"> 
Opinionated flow to re-run tests and identify flakes, races, timeouts, and test runtimes.
</mode>
</modes>

<commands>
Run these commands to validate any changes you make
```sh
golangci-lint run ./... --fix
go test ./...
```
</commands>