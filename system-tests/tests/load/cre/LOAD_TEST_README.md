## Prerequisites

Before running the load tests, ensure you have the following set up:

### Required Software
1. **Docker**
    - Local Docker repository configured
    - Chainlink node image built and available in registry
    - Job Distributor (JD) image built and available  in registry

2. **Nix Package Manager**
    - If you’re on macOS, the recommended way to install Nix is using the Determinate Systems graphical installer

3. **Mock Capability**
    - Must be compiled from the capabilities repository, details on how to compile it can be found in the `How to compile the mock capability` below
    - Required for simulating capability interactions

4. **Observability stack**
    - Tests require Grafana and Loki for monitoring and logging
    - To set up the observability stack, use the chainlink-testing-framework command: `ctf obs up`
 

### Cloud Requirements
- **For AWS Testing Only:**
    - Access to the CRIB repository
    - AWS credentials configured
    - Appropriate permissions set up
    - Devspace installed
    - VPN turned on

### How to compile the mock capability
#### 1. Clone the capabilities repository
`git clone https://github.com/smartcontractkit/capabilities`

#### 2. Navigate to the repository
`cd capabilities`

#### 3. Build the mock capability
`nx build mock`

After compilation, the mock binary will be located at `capabilities/bin/amd64/mock`.

Then update your test configuration file:
```
[binaries_config]
mock_capability_binary_path = "/path/to/capabilities/bin/amd64/mock"
```
Use either an absolute path or a relative path pointing to your compiled binary.

### Setting Up Observability stack

You can skip this step if either:
- You already have Grafana and Loki running locally
- You plan to use remote instances of these services

For installation:
1. See the [CTF documentation](https://smartcontractkit.github.io/chainlink-testing-framework/framework/getting_started.html) for up-to-date setup instructions
2. After installing the CTF CLI, run `ctf obs up` to launch local Grafana and Loki instances

### Clone CRIB repo
1. Clone the CRIB repository and checkout branch `dx-50-cre-crib-df-perf-testing`
2. Create a `.env` file in `crib/deployments/cre` (use `.env.example` as a template)

## Deployment Types

After completing the prerequisites, configure your testing TOML file. The load tests support two deployment types (`...-kind.toml` and `...-crib.toml`):

### Kind Deployment
- **Environment**: Runs locally using Docker with Anvil for blockchain simulation
- **Hardware Requirements**: Consider your machine capabilities when setting the number of nodes
- **Best For**: Development and simple tests due to faster deployment time
- **Limitations**: Not suitable for large deployments due to local hardware constraints

### Crib Deployment (AWS)
- **Environment**: Runs in AWS, can be set up to providing a production-like testing environment
- **Configuration**: Use `LOAD_TEST="true"` environment variable to set production-like pod resource limits
- **Deployment Notes**:
    - Deployment is slower than Kind
    - Requires active VPN connection (may fail due to connectivity issues)
    - Deployments persist longer-term
- **Advanced Testing**: Values are saved after deployment to support `TestXXXWithReconnect` methods, allowing you to simulate different load scenarios without redeploying





### 10. Run test
Cd into the current directory and run the command:
```
CTF_CONFIGS=writer-don-load-test-kind.toml \
DASHBOARD_FOLDER=LoadTests \
DASHBOARD_NAME=Wasp \
DATA_SOURCE_NAME=Loki \
GRAFANA_URL=http://localhost:3000 \
LOKI_URL=http://localhost:3030/loki/api/v1/push \
PRIVATE_KEY=ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80 \
CHAINLINK_CLUSTER_VERSION="2.6.0-preview-1896-f6cd441" \
LOAD_TEST="true" \
go test -run TestLoad_Writer_MockCapabilities -timeout 60m
```
The test runs with a 60-minute timeout using the local Kind deployment configuration.

**Environment Variables**

| Variable                  | Purpose                                                                                          |
|---------------------------|--------------------------------------------------------------------------------------------------|
| CTF_CONFIGS               | Specifies the test configuration file with node setup and test parameters                        |
| DASHBOARD_FOLDER          | Grafana folder for storing test metrics                                                          |
| DASHBOARD_NAME            | Name of the dashboard for viewing results                                                        |
| DATA_SOURCE_NAME          | Log data source (Loki)                                                                           |
| GRAFANA_URL               | URL of your local Grafana instance                                                               |
| LOKI_URL                  | Endpoint for sending test logs                                                                   |
| PRIVATE_KEY               | Ethereum key for blockchain transactions (default Anvil key), change only if deploying with Geth |
| CHAINLINK_CLUSTER_VERSION | Helmchart version with ingress mock support                                                      |
| LOAD_TEST                 | When "true", increases resource limits for production-like testing                               |






## Chaos tests

! These tests require CRIB environment with 3 DONS: Asset, Workflow and Writer, and at least one EVM chain !

Chaos tests are added to the load test by default but can be run only on `main.stage`.

There are 3 modes:
- "clean" - no chaos tests
- "reorg" - EVM specific chaos suite - block reorganizations
- "rpc" - simulating realistic RPC latency
- "full" - running all the experiments

Turn your VPN on, add `[chaos.dashboard_uids]` on which the test will automatically add annotations, export 1Password `CRE chaos test vars` under `Eng Shared Vault` and run the test.

## Chaos tests without workload

If you want to apply the chaos test without workload or debug it you can use
```
go test -v -timeout 2h -run TestChaos
```

