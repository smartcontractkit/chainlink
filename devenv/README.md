# Chainlink Developer Environment

A self-contained Docker-based environment for running Chainlink system-level tests. Spins up an EVM chain (Anvil/Geth), a fake external adapter server, and Chainlink nodes, then deploys product-specific contracts and jobs.

For detailed architecture documentation, see [design.md](design.md).

## Quickstart

### Prerequisites

- **Docker** -- must be running
- **Go** -- version specified in `go.mod`
- **Just** -- task runner ([install guide](https://github.com/casey/just?tab=readme-ov-file#cross-platform))

```bash
brew install just # macOS; see link above for other platforms
```

### One-Time Setup

From the `devenv/` directory:

```bash
just build-fakes # Build the mock external adapter Docker image
just cli         # Install the `cl` CLI tool
```

### Running a Test

Tests use a two-terminal workflow. Pick a product from the reference table below and run:

**Terminal 1** -- start the environment (from `devenv/`):

```bash
cl u env.toml,products/<product>/basic.toml
```

**Terminal 2** -- run the test (from `devenv/tests/<product>/`):

```bash
go test -v -run <TestName>
```

## Observability

```bash
cl obs up     # Loki + Prometheus + Grafana
cl obs up -f  # Full stack (adds Pyroscope, cadvisor, postgres exporter)
cl obs down   # Tear down
```

### Dashboards

- [All Logs](http://localhost:3000/explore?schemaVersion=1&panes=%7B%22axv%22:%7B%22datasource%22:%22P8E80F9AEF21F6940%22,%22queries%22:%5B%7B%22refId%22:%22A%22,%22expr%22:%22%7Bjob%3D%5C%22ctf%5C%22%7D%22,%22queryType%22:%22range%22,%22datasource%22:%7B%22type%22:%22loki%22,%22uid%22:%22P8E80F9AEF21F6940%22%7D,%22editorMode%22:%22code%22,%22direction%22:%22backward%22%7D%5D,%22range%22:%7B%22from%22:%22now-30m%22,%22to%22:%22now%22%7D,%22compact%22:false%7D%7D&orgId=1)
- [CL Node Errors](http://localhost:3000/d/a7de535b-3e0f-4066-bed7-d505b6ec9ef1/cl-node-errors?orgId=1&refresh=5s&from=now-15m&to=now&timezone=browser)
- [Load Testing](http://localhost:3000/d/WASPLoadTests/wasp-load-test?orgId=1&from=now-30m&to=now&timezone=browser&var-go_test_name=$__all&var-gen_name=$__all&var-branch=$__all&var-commit=$__all&var-call_group=$__all&refresh=5s)

### Product Reference

Each row maps to a CI matrix entry in [devenv-nightly.yml](../.github/workflows/devenv-nightly.yml).

| Product        | envcmd (from `devenv/`)                                                        | testcmd (from `devenv/tests/<dir>/`)                               | tests_dir       |
| -------------- | ------------------------------------------------------------------------------ | ------------------------------------------------------------------ | --------------- |
| Cron           | `cl u env.toml,products/cron/basic.toml`                                       | `go test -v -run TestSmoke`                                        | `cron`          |
| Direct Request | `cl u env.toml,products/directrequest/basic.toml`                              | `go test -v -run TestSmoke`                                        | `directrequest` |
| Flux           | `cl u env.toml,products/flux/basic.toml`                                       | `go test -v -run TestSmoke`                                        | `flux`          |
| VRF            | `cl u env.toml,products/vrf/basic.toml`                                        | `go test -v -timeout 10m -run TestVRFBasic\|TestVRFJobReplacement` | `vrf`           |
| Automation 2.0 | `cl u env.toml,products/automation/basic.toml`                                 | `go test -v -timeout 30m -run TestRegistry_2_0`                    | `automation`    |
| Automation 2.1 | `cl u env.toml,products/automation/basic.toml`                                 | `go test -v -timeout 30m -run TestRegistry_2_1`                    | `automation`    |
| Keepers        | `cl u env.toml,products/keepers/basic.toml`                                    | `go test -v -timeout 1h -run TestKeeperBasic`                      | `keepers`       |
| OCR2 Smoke     | `cl u env.toml,products/ocr2/basic.toml`                                       | `go test -v -run TestSmoke`                                        | `ocr2`          |
| OCR2 Soak      | `cl u env.toml,products/ocr2/basic.toml,products/ocr2/soak.toml; cl obs up -f` | `go test -v -timeout 4h -run TestOCR2Soak/clean`                   | `ocr2`          |

## Interactive Shell

Instead of running CLI commands directly, you can use the interactive shell with autocomplete:

```bash
cl sh
```

Then inside the shell:

```sh
up                                  # start with default OCR2 config
up env.toml,products/vrf/basic.toml # start with a specific product
obs up -f                           # start full observability stack
test ocr2 TestSmoke                 # run a test
down                                # tear down everything
```

## Run with Custom CL Image

Use `env-cl-rebuild.toml` to build a Chainlink image from your local repository:

```bash
cl u env.toml,products/ocr2/basic.toml,env-cl-rebuild.toml
```

Or override the image via environment variable:

```bash
CHAINLINK_IMAGE=my-registry/chainlink:dev cl u env.toml,products/ocr2/basic.toml
```

## Contributing

### Updating Fakes

Fakes are mock external adapters that return controlled feed values. See [design.md](design.md#fakes-mock-external-adapters) for details.

```bash
just build-fakes               # Build locally
just push-fakes <aws_registry> # Push to ECR
```

### Adding Products

Implement the [Product interface](interface.go) and add a switch clause in [environment.go](environment.go). See [design.md](design.md#the-product-interface) for the full lifecycle.
