# Devenv Architecture

## Overview

`devenv` is a self-contained Go module (`github.com/smartcontractkit/chainlink/devenv`) that provides a Docker-based development and testing environment for Chainlink products. It orchestrates local blockchain networks, Chainlink nodes, mock external adapters, and product-specific contract deployments.

Key design principles:

- **Dependency isolation** -- devenv does NOT import `github.com/smartcontractkit/chainlink/v2` or any of its child packages. This keeps the test environment decoupled from the core node codebase.
- **TOML-driven configuration** -- all infrastructure and product settings are declared in composable TOML files that merge left-to-right.
- **Two-phase testing** -- environment setup (CLI) and test execution (`go test`) are separate processes, connected by a shared `env-out.toml` output file.
- **Product abstraction** -- each Chainlink product implements a common `Product` interface, making it straightforward to add new products.

The module depends on the [Chainlink Testing Framework (CTF)](https://github.com/smartcontractkit/chainlink-testing-framework) for Docker orchestration, CL node HTTP clients, and observability tooling.

## High-Level Architecture

```mermaid
flowchart TD
    subgraph config [TOML Configuration]
        envToml["env.toml\n(infra: chain, nodes, fakes)"]
        productToml["products/&lt;name&gt;/basic.toml\n(product settings)"]
        overrideToml["env-geth.toml etc.\n(optional overrides)"]
    end

    subgraph cli [CLI]
        clUp["cl up &lt;configs&gt;"]
    end

    subgraph envSetup ["NewEnvironment()"]
        loadConfig["Load and merge\nTOML configs"]
        startInfra["Start infrastructure\n(Anvil, Fake Server)"]
        genNodeConfig["Products generate\nCL node config"]
        startNodes["Start CL node set\n(shared DB)"]
        storeInfra["Store infra output"]
        deployProducts["ConfigureJobsAndContracts()\nper product instance"]
        storeProducts["Store product output"]
    end

    subgraph output [Output]
        envOut["env-out.toml\n(addresses, URLs, job IDs)"]
    end

    subgraph testPhase [Test Execution]
        goTest["go test -v -run TestName"]
        loadOutput["Load env-out.toml"]
        assertions["Assert on-chain state\nvia gethwrappers +\nCL node API via clclient"]
    end

    envToml --> clUp
    productToml --> clUp
    overrideToml --> clUp
    clUp --> loadConfig
    loadConfig --> startInfra
    startInfra --> genNodeConfig
    genNodeConfig --> startNodes
    startNodes --> storeInfra
    storeInfra --> deployProducts
    deployProducts --> storeProducts
    storeProducts --> envOut
    envOut --> loadOutput
    goTest --> loadOutput
    loadOutput --> assertions
```

## Configuration System

The configuration system uses composable TOML files merged via the `CTF_CONFIGS` environment variable.

### Merge Semantics

When `cl up env.toml,products/ocr2/basic.toml` runs, it sets `CTF_CONFIGS=env.toml,products/ocr2/basic.toml`. The `Load[T]()` function reads each file left-to-right, decoding into the same struct. Later files override earlier keys while preserving keys they do not mention.

```mermaid
flowchart LR
    A["env.toml\n(blockchains, fake_server,\nnodesets)"] -->|merge| C["Merged Config"]
    B["products/ocr2/basic.toml\n(products, ocr2 settings,\nnodeset overrides)"] -->|merge| C
    C --> D["NewEnvironment()\nuses merged config"]
    D --> E["env-out.toml\n(infra + product outputs)"]
```

### Config Layers

| Layer               | File                         | Purpose                                                                    |
| ------------------- | ---------------------------- | -------------------------------------------------------------------------- |
| Base infrastructure | `env.toml`                   | Chain type/ID, fake server image, node count and images                    |
| Product config      | `products/<name>/basic.toml` | Product name, instances, product-specific settings, node count override    |
| Chain override      | `env-geth.toml`              | Switch from Anvil to Geth                                                  |
| Image override      | `env-cl-rebuild.toml`        | Build CL image from local Dockerfile                                       |
| Runtime output      | `env-out.toml`               | Generated after `cl up` -- contains deployed addresses, node URLs, job IDs |

### Root Config Struct

The root configuration type (`Cfg` in `environment.go`) defines the top-level TOML schema:

```go
type Cfg struct {
    Products    []*ProductInfo      `toml:"products"`
    Blockchains []*blockchain.Input `toml:"blockchains"`
    FakeServer  *fake.Input         `toml:"fake_server"`
    NodeSets    []*ns.Input         `toml:"nodesets"`
    JD          *jd.Input           `toml:"jd"`
}
```

Each product configurator has its own struct that gets decoded from the same TOML files (e.g., `[[ocr2]]` sections are decoded by the OCR2 `Configurator`).

## The Product Interface

Every product in devenv implements this interface from `interface.go`:

```go
type Product interface {
    Load() error
    Store(path string, instanceIdx int) error
    GenerateNodesSecrets(ctx, fs, bc, ns) (string, error)
    GenerateNodesConfig(ctx, fs, bc, ns) (string, error)
    ConfigureJobsAndContracts(ctx, instanceIdx, fs, bc, ns) error
}
```

### Product Lifecycle

```mermaid
sequenceDiagram
    participant CLI as cl up
    participant Env as NewEnvironment
    participant Product as Product Configurator
    participant Infra as Docker Infrastructure
    participant Chain as Blockchain

    CLI->>Env: Start with merged TOML config
    Env->>Infra: Create blockchain network (Anvil/Geth)
    Env->>Infra: Create fake data provider

    loop For each product in config
        Env->>Product: Load()
        Product-->>Env: Product config loaded from TOML
        Env->>Product: GenerateNodesConfig()
        Product-->>Env: CL node TOML overrides
        Env->>Product: GenerateNodesSecrets()
        Product-->>Env: CL node secrets overrides
    end

    Note over Env,Infra: Merge all product config overrides into node specs
    Env->>Infra: Start CL node set (shared DB)
    Env->>Env: Store infrastructure output

    loop For each product, for each instance
        Env->>Product: ConfigureJobsAndContracts()
        Product->>Chain: Deploy contracts (LINK, product contracts)
        Product->>Infra: Fund CL nodes, create keys
        Product->>Infra: Create jobs on CL nodes
        Product->>Chain: Register on-chain config
        Env->>Product: Store()
        Product-->>Env: Write output (addresses, job IDs)
    end
```

### Registered Products

| Name           | TOML key         | Config dir                | Nodes | Contracts deployed                                 |
| -------------- | ---------------- | ------------------------- | ----- | -------------------------------------------------- |
| Cron           | `cron`           | `products/cron/`          | 1     | None (bridge + cron job only)                      |
| Direct Request | `direct_request` | `products/directrequest/` | 1     | LINK, Oracle, TestAPIConsumer                      |
| Flux Monitor   | `flux`           | `products/flux/`          | 5     | LINK, FluxAggregator                               |
| OCR2           | `ocr2`           | `products/ocr2/`          | 5     | LINK, OCR2Aggregator                               |
| Automation     | `automation`     | `products/automation/`    | 5     | LINK, Registry (2.0-2.3), Registrar, Upkeeps       |
| Keepers        | `keepers`        | `products/keepers/`       | 5     | LINK, KeeperRegistry (1.1-1.3), Registrar, Upkeeps |
| VRF            | `vrf`            | `products/vrf/`           | 1     | LINK, BlockHashStore, VRFCoordinator, VRFConsumer  |

### Adding a New Product

1. Create `products/<name>/configuration.go` implementing the `Product` interface
2. Create `products/<name>/basic.toml` with default config
3. Add a `case "<name>"` in `newProduct()` in `environment.go`
4. Create `tests/<name>/smoke_test.go` that reads `env-out.toml` and asserts behavior
5. Add a matrix entry in `.github/workflows/devenv-nightly.yml`

## The `cl` CLI

The CLI (`cmd/cl/`) is a Cobra-based tool that drives environment lifecycle.

### Commands

| Command                     | Alias  | Description                                                                          |
| --------------------------- | ------ | ------------------------------------------------------------------------------------ |
| `cl up [configs]`           | `cl u` | Spin up environment from TOML configs (default: `env.toml,products/ocr2/basic.toml`) |
| `cl down`                   | `cl d` | Tear down all Docker containers                                                      |
| `cl restart [configs]`      | `cl r` | Tear down then recreate                                                              |
| `cl test <folder> <filter>` |        | Run `go test` in `tests/<folder>` with `-run <filter>`                               |
| `cl obs up [-f]`            |        | Start observability stack (Loki/Prometheus/Grafana; `-f` for full)                   |
| `cl obs down`               |        | Stop observability stack                                                             |
| `cl bs up`                  |        | Start Blockscout block explorer                                                      |
| `cl bs down`                |        | Stop Blockscout                                                                      |
| `cl shell` / `cl sh`        |        | Interactive shell with autocomplete                                                  |

### How `cl up` Works

1. Sets `CTF_CONFIGS` env var from the argument (or defaults to `env.toml,products/ocr2/basic.toml`)
2. Sets `TESTCONTAINERS_RYUK_DISABLED=true` to prevent container cleanup on CLI exit
3. Calls `devenv.NewEnvironment(ctx)` with a 7-minute timeout
4. `NewEnvironment` loads config, starts infra, runs product configurators, writes `env-out.toml`

### Interactive Shell

`cl sh` starts a REPL with tab-completion for commands and config file paths. It executes commands by invoking the same Cobra command tree in-process.

## Fakes (Mock External Adapters)

Fakes are a lightweight HTTP service that replaces real external adapters and data feeds during testing. The fake server runs as a Docker container on port 9111.

### Why Fakes Exist

Chainlink nodes need external data sources (external adapters, price feeds, Mercury endpoints). Instead of depending on real services, fakes provide deterministic, controllable responses that make tests reliable and fast.

### Routes by Product

| Product        | Route                           | Behavior                                                    |
| -------------- | ------------------------------- | ----------------------------------------------------------- |
| Cron           | `POST /cron_response`           | Returns `{"data": {"result": 200}}`                         |
| Direct Request | `POST /direct_request_response` | Returns `{"data": {"result": 200}}`                         |
| OCR2           | `POST /ea`                      | Returns current EA value (default 200)                      |
| OCR2           | `POST /juelsPerFeeCoinSource`   | Returns JUELS/LINK ratio                                    |
| OCR2           | `POST /trigger_deviation`       | Changes the EA return value (query param `?result=<value>`) |
| Automation     | `POST /api/v1/reports/bulk`     | Returns mock Mercury/DataStreams reports                    |
| Automation     | `GET /client`                   | Returns mock Mercury client config                          |

### Building and Using Fakes

```bash
just build-fakes                     # Build image locally as chainlink-fakes:latest
just push-fakes <aws_registry>       # Build for linux/amd64 and push to ECR
```

In CI, the `FAKE_SERVER_IMAGE` environment variable overrides the image used in `env.toml`.

## Test Architecture

Tests follow a two-phase pattern where environment setup and test execution are independent processes.

```mermaid
sequenceDiagram
    participant Dev as Developer
    participant CLI as cl up (Terminal 1)
    participant Docker as Docker Containers
    participant Test as go test (Terminal 2)
    participant EnvOut as env-out.toml

    Dev->>CLI: cl u env.toml,products/vrf/basic.toml
    CLI->>Docker: Start Anvil, Fake Server, CL Nodes
    CLI->>Docker: Deploy contracts, create jobs
    CLI->>EnvOut: Write deployed state

    Dev->>Test: go test -v -run TestVRFBasic
    Test->>EnvOut: Load config + product output
    Test->>Docker: Interact with contracts (gethwrappers)
    Test->>Docker: Query CL node API (clclient)
    Test->>Test: Assert expected behavior
```

### Test File Pattern

Every smoke test follows the same structure:

1. **Load output** -- read `../../env-out.toml` to get infrastructure and product config

```go
in, err := de.LoadOutput[de.Cfg](outputFile)
productCfg, err := products.LoadOutput[<product>.Configurator](outputFile)
```

2. **Setup cleanup** -- save container logs on test completion

```go
t.Cleanup(func() {
    framework.SaveContainerLogs(...)
})
```

3. **Create clients** -- ETH client for on-chain interaction, CL client for node API

```go
c, auth, _, err := products.ETHClient(ctx, wsURL, feeCap, tipCap)
cls, err := clclient.New(in.NodeSets[0].Out.CLNodes)
```

4. **Interact with contracts** -- use gethwrappers directly (never through `chainlink/v2` wrappers)

```go
consumer, err := solidity_vrf_consumer_interface.NewVRFConsumer(addr, c)
```

5. **Assert with polling** -- use `require.EventuallyWithT` to poll until expected state

```go
require.EventuallyWithT(t, func(ct *assert.CollectT) {
    // check on-chain state or job runs
}, 2*time.Minute, 2*time.Second)
```

### Dependency Rule

Tests in `devenv/tests/` must NOT import:
- `github.com/smartcontractkit/chainlink/v2`
- `github.com/smartcontractkit/chainlink/integration-tests`
- `github.com/smartcontractkit/chainlink/deployment`

Allowed imports are:
- `github.com/smartcontractkit/chainlink/devenv` (this module)
- `github.com/smartcontractkit/chainlink-testing-framework/framework` (CTF)
- `github.com/smartcontractkit/chainlink-evm/gethwrappers` (contract bindings)
- `github.com/smartcontractkit/libocr` (OCR bindings)
- Standard library and third-party libraries (testify, go-ethereum, etc.)

## CI Integration

System tests run nightly via [`.github/workflows/devenv-nightly.yml`](../.github/workflows/devenv-nightly.yml).

### Workflow Structure

The workflow uses a GitHub Actions matrix strategy where each entry defines:

| Field               | Purpose                                                          |
| ------------------- | ---------------------------------------------------------------- |
| `display_name`      | Human-readable test name                                         |
| `envcmd`            | Command to set up the environment (runs from `devenv/`)          |
| `testcmd`           | Command to run the tests (runs from `devenv/tests/<tests_dir>/`) |
| `runner`            | GitHub Actions runner label                                      |
| `tests_dir`         | Subdirectory under `devenv/tests/`                               |
| `logs_archive_name` | Name for the uploaded log artifact                               |

### Execution Flow

```mermaid
flowchart TD
    A[Checkout code] --> B[Setup Docker Buildx]
    B --> C[Install Just]
    C --> D[AWS OIDC + ECR Login]
    D --> E[Setup Go + download deps]
    E --> F["Set CHAINLINK_IMAGE\n(nightly build)"]
    F --> G["Install cl CLI\n(go install cmd/cl)"]
    G --> H["eval envcmd\n(starts environment)"]
    H --> I["eval testcmd\n(runs go test)"]
    I --> J["Upload logs artifact\n(always)"]
```

### Adding a Test to CI

Add a new entry to the `matrix.include` array:

```yaml
- display_name: "Test <Product> Smoke"
  testcmd: "go test -v -timeout 10m -run <TestFunction>"
  envcmd: "cl u env.toml,products/<product>/basic.toml"
  runner: "ubuntu-latest"
  tests_dir: "<product>"
  logs_archive_name: "<product>"
```
