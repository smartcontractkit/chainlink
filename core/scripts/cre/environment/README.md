# Local CRE environment

## Contact Us
Slack: #topic-local-dev-environments

## Table of content

1. [Using the CLI](#using-the-cli)
   - [Prerequisites](#prerequisites-for-docker)
   - [Start Environment](#start-environment)
    - [Using Existing Docker plugins image](#using-existing-docker-plugins-image)
    - [Beholder](#beholder)
    - [Storage](#storage)
   - [Stop Environment](#stop-environment)
   - [Restart Environment](#restarting-the-environment)
   - [DX Tracing](#dx-tracing)
2. [Job Distributor Image](#job-distributor-image)
3. [Example Workflows](#example-workflows)
4. [Troubleshooting](#troubleshooting)

# Using the CLI

The CLI manages CRE test environments. It is located in `core/scripts/cre/environment`. It doesn't come as a compiled binary, so every command has to be executed as `go run . <command> [subcommand]`.

## Prerequisites (for Docker) ###
1. **Docker installed and running**
    - with usage of default Docker socket **enabled**
    - with Apple Virtualization framework **enabled**
    - with VirtioFS **enabled**
    - with use of containerd for pulling and storing images **disabled**
2. **AWS SSO access to SDLC**
  - REQUIRED: `sdlc` profile (with `PowerUserAccess` role)
>  [See more for configuring AWS in CLL](https://smartcontract-it.atlassian.net/wiki/spaces/INFRA/pages/1045495923/Configure+the+AWS+CLI)


## Prerequisites For CRIB ###
1. telepresence installed: `brew install telepresenceio/telepresence/telepresence-oss`
2. Telepresence will update the /etc/resolver configs and will require to enter sudo password the first time you run it

# QUICKSTART
```
# e.g. AWS_ECR=<PROD_ACCOUNT_ID>.dkr.ecr.<REGION>.amazonaws.com
AWS_ECR=<PROD_AWS_URL> go run . env start --auto-setup
```
> You can find `PROD_ACCOUNT_ID` and `REGION` in the `[profile prod]` section of the [AWS CLI configuration guide](https://smartcontract-it.atlassian.net/wiki/spaces/INFRA/pages/1045495923/Configure+the+AWS+CLI#Configure).

If you are missing requirements, you may need to fix the errors and re-run.

Refer to [this document](https://docs.google.com/document/d/1HtVLv2ipx2jvU15WYOijQ-R-5BIZrTdAaumlquQVZ48/edit?tab=t.0#heading=h.wqgcsrk9ncjs) for troubleshooting and FAQ. Use `#topic-local-dev-environments` for help.

## Start Environment
```bash
# while in core/scripts/cre/environment
go run . env start [--auto-setup]

# to start environment with an example workflow web API-based workflow
go run . env start --with-example

 # to start environment with an example workflow cron-based workflow (this requires the `cron` capability binary to be setup in the `extra_capabilities` section of the TOML config)
go run . env start --with-example --example-workflow-trigger cron

# to start environment using image with all supported capabilities
go run . env start --with-plugins-docker-image <SDLC_ACCOUNT_ID>dkr.ecr.<SDLC_ACCOUNT_REGION>.amazonaws.com/chainlink:nightly-<YYYMMDD>-plugins

# to start environment with local Beholder
go run . env start --with-beholder
```

> Important! **Nightly** Chainlink images are retained only for one day and built at 03:00 UTC. That means that in most cases you should use today's image, not yesterday's.

Optional parameters:
- `-a`: Check if all dependencies are present and if not install them (defaults to `false`)
- `-t`: Topology (`simplified` or `full`)
- `-w`: Wait on error before removing up Docker containers (e.g. to inspect Docker logs, e.g. `-w 5m`)
- `-e`: Extra ports for which external access by the DON should be allowed (e.g. when making API calls or downloading WASM workflows)
- `-x`: Registers an example PoR workflow using CRE CLI and verifies it executed successfuly
- `-s`: Time to wait for example workflow to execute successfuly (defaults to `5m`)
- `-p`: Docker `plugins` image to use (must contain all of the following capabilities: `ocr3`, `cron`, `readcontract` and `logevent`)
- `-y`: Trigger for example workflow to deploy (web-trigger or cron). Default: `web-trigger`. **Important!** `cron` trigger requires user to either provide the capbility binary path in TOML config or Docker image that has it baked in
- `-c`: List of configuration files for `.proto` files that will be registered in Beholder (only if `--with-beholder/-b` flag is used). Defaults to [./proto-configs/default.toml](./proto-configs/default.toml)

### Using existing Docker Plugins image

If you don't want to build Chainlink image from your local branch (default behaviour) or you don't want to go through the hassle of downloading capabilities binaries in order to enable them on your environment you should use the `--with-plugins-docker-image` flag. It is recommended to use a nightly `core plugins` image that's build by [Docker Build action](https://github.com/smartcontractkit/chainlink/actions/workflows/docker-build.yml) as it contains all supported capability binaries.

### Beholder

When environment is started with `--with-beholder` or with `-b` flag after the DON is ready  we will boot up `Chip Ingress` and `Red Panda`, create a `cre` topic and download and install workflow-related protobufs from the [chainlink-protos](https://github.com/smartcontractkit/chainlink-protos/tree/main/workflows) repository.

Once up and running you will be able to access [CRE topic view](http://localhost:8080/topics/cre) to see workflow-emitted events. These include both standard events emitted by the Workflow Engine and custom events emitted from your workflow.

#### Filtering out heartbeats
Heartbeat messages spam the topic, so it's highly recommended that you add a JavaScript filter that will exclude them using the following code: `return value.msg !== 'heartbeat';`.

If environment is aready running you can start just the Beholder stack (and register protos) with:
```bash
go run . env beholder start
```

> This assumes you have `chip-ingress:qa-latest` Docker image on your local machine. Without it Beholder won't be able to start. If you do not, close the [Atlas](https://github.com/smartcontractkit/atlas) repository, and then in `atlas/chip-ingress` run `docker build -t chip-ingress:qa-latest .`

### Storage

By default, workflow artifacts are loaded from the container's filesystem. The Chainlink nodes can only load workflow files from the local filesystem if `WorkflowFetcher` uses the `file://` prefix. Right now, it cannot read workflow files from both the local filesystem and external sources (like S3 or web servers) at the same time.

The environment supports two storage backends for workflow uploads:
- Gist (requires deprecated CRE CLI, remote)
- S3 MinIO (built-in, local)

Configuration details for the CRE CLI are generated automatically into the `cre.yaml` file
(path is printed after starting the environment).

For more details on the URL resolution process and how workflow artifacts are handled, see the [URL Resolution Process](../../../../system-tests/tests/smoke/cre/guidelines.md#url-resolution-process) section in `system-tests/tests/smoke/cre/guidelines.md`.

## Stop Environment
```bash
# while in core/scripts/cre/environment
go run main.go env stop

# or... if you have the CTF binary
ctf d rm
```
---

## Restarting the environment

If you are using Blockscout and you restart the environment **you need to restart the block explorer** if you want to see current block history. If you don't you will see stale state of the previous environment. To restart execute:
```bash
ctf bs r
```
---

## Workflow Commands

The environment provides workflow management commands defined in `core/scripts/cre/environment/environment/workflow.go`:

### `workflow deploy`
Compiles and uploads a workflow to the environment by copying it to workflow nodes and registering with the workflow registry. It checks if a workflow with same name already exists and deletes it, if it does.

**Usage:**
```bash
go run . workflow deploy [flags]
```

**Key flags:**
- `-w, --workflow-file-path`: Path to the workflow file (default: `./examples/workflows/v2/cron/main.go`)
- `-c, --config-file-path`: Path to the config file (optional)
- `-s, --secrets-file-path`: Path to the secrets file (optional)
- `-t, --container-target-dir`: Path to target directory in Docker container (default: `/home/chainlink/workflows`)
- `-o, --container-name-pattern`: Pattern to match container name (default: `workflow-node`)
- `-n, --workflow-name`: Workflow name (default: `exampleworkflow`)
- `-r, --rpc-url`: RPC URL (default: `http://localhost:8545`)
- `-i, --chain-id`: Chain ID (default: `1337`)
- `-a, --workflow-registry-address`: Workflow registry address (default: `0x9fE46736679d2D9a65F0992F2272dE9f3c7fa6e0`)
- `-b, --capabilities-registry-address`: Capabilities registry address (default: `0xe7f1725E7734CE288F8367e1Bb143E90bb3F0512`)
- `-d, --workflow-owner-address`: Workflow owner address (default: `0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266`)
- `-e, --don-id`: DON ID (default: `1`)

**Example:**
```bash
go run . workflow deploy -w ./my-workflow.go -n myworkflow -c ./config.yaml
```

### `workflow delete`
Deletes a specific workflow from the workflow registry contract (but doesn't remove it from Docker containers).

**Usage:**
```bash
go run . workflow delete [flags]
```

**Key flags:**
- `-n, --name`: Workflow name to delete (default: `exampleworkflow`)
- `-r, --rpc-url`: RPC URL (default: `http://localhost:8545`)
- `-i, --chain-id`: Chain ID (default: `1337`)
- `-a, --workflow-registry-address`: Workflow registry address (default: `0x9fE46736679d2D9a65F0992F2272dE9f3c7fa6e0`)
- `-d, --workflow-owner-address`: Workflow owner address (default: `0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266`)

**Example:**
```bash
go run . workflow delete -n myworkflow
```

### `workflow delete-all`
Deletes all workflows from the workflow registry contract.

**Usage:**
```bash
go run . workflow delete-all [flags]
```

**Key flags:**
- `-r, --rpc-url`: RPC URL (default: `http://localhost:8545`)
- `-i, --chain-id`: Chain ID (default: `1337`)
- `-a, --workflow-registry-address`: Workflow registry address (default: `0x9fE46736679d2D9a65F0992F2272dE9f3c7fa6e0`)
- `-d, --workflow-owner-address`: Workflow owner address (default: `0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266`)

**Example:**
```bash
go run . workflow delete-all
```

### `workflow deploy-and-verify-example`
Deploys and verifies the example workflow.

**Usage:**
```bash
go run . workflow deploy-and-verify-example
```

This command uses default values and is useful for testing the workflow deployment process.

---

## Further use
To manage workflows you will need the CRE CLI. You can either:
- download it from [smartcontract/dev-platform](https://github.com/smartcontractkit/dev-platform/releases/tag/v0.2.0) or
- using GH CLI:
  ```bash
  gh release download v0.2.0 --repo smartcontractkit/dev-platform --pattern '*darwin_arm64*'
  ```

Remember that the CRE CLI version needs to match your CPU architecture and operating system.

---

### Advanced Usage:
1. **Choose the Right Topology**
   - For a single DON with all capabilities: `configs/workflow-don.toml` (default)
   - For a single DON with all capabilities, but with a separate gateway node: `configs/workflow-gateway-don.toml`
   - For a full topology (workflow DON + capabilities DON + gateway DON): `configs/workflow-gateway-capabilities-don.toml`
2. **Download or Build Capability Binaries**
   - Some capabilities like `cron`, `log-event-trigger`, or `read-contract` are not embedded in all Chainlink images.
   - If your use case requires them, you should build them manually by:
      - Cloning [smartcontractkit/capabilities](https://github.com/smartcontractkit/capabilities) repository (Make sure they are built for `linux/amd64`!)
      - Building each capability manually by running `GOOS="linux" GOARCH="amd64" CGO_ENABLED=0 go build -o evm` inside capability's folder or building all of them at once with `./nx run-many -t build` in root of `capabilities` folder

     Once that is done reference them in your TOML like:
       ```toml
       [extra_capabilities]
       cron_capability_binary_path = "./cron" # remember to adjust binary name and path
       # log even trigger and read-contract binaries go here
       # they are all commented out by default
       ```
     Do make sure that the path to the binary is either relative to the `environment` folder or absolute. Then the binary will be copied to the Docker image.
   - If the capability is already baked into your CL image (check the Dockerfile), comment out the TOML path line to skip copying. (they will be commented out by default)
3.  **Decide whether to build or reuse Chainlink Docker Image**
     - By default, the config builds the Docker image from your local branch. To use an existing image change to:
     ```toml
     [nodesets.node_specs.node]
     image = "<your-Docker-image>:<your-tag>"
     ```
      - Make these changes for **all** nodes in the nodeset in the TOML config.
      - If you decide to reuse a Chainlink Docker Image using the `--with-plugins-docker-image` flag, please notice that this will not copy any capability binaries to the image.
        You will need to make sure that all the capabilities you need are baked in the image you are using.

4. **Decide whether to use Docker or k8s**
    - Read [Docker vs Kubernetes in guidelines.md](../../../../system-tests/tests/smoke/cre/guidelines.md) to learn how to switch between Docker and Kubernetes
5. **Start Observability Stack (Docker-only)**
      ```bash
      # to start Loki, Grafana and Prometheus run:
      ctf obs up

     # to start Blockscout block explorer run:
      ctf bs u
      ```
    - To download the `ctf` binary follow the steps described [here](https://smartcontractkit.github.io/chainlink-testing-framework/framework/getting_started.html)

Optional environment variables used by the CLI:
- `CTF_CONFIGS`: TOML config paths. Defaults to [./configs/workflow-don.toml](./configs/workflow-don.toml)
- `PRIVATE_KEY`: Plaintext private key that will be used for all deployments (needs to be funded). Defaults to `ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80`
- `TESTCONTAINERS_RYUK_DISABLED`: Set to "true" to disable cleanup. Defaults to `false`

When starting the environment in AWS-managed Kubernetes make sure to source `.env` environment from the `crib/deployments/cre` folder specific for AWS. Remember, that it must include ingress domain settings.

### Testing Billing
Spin up the billing service and necessary migrations in the `billing-platform-service` repo.
The directions there should be sufficient. If they are not, give a shout to @cre-business.

So far, all we have working is the workflow run by:
```go
go run . env start --with-example -w 1m
```

I recommend increasing your docker resources to near max memory, as this is going to slow your local
machine down anyways, and it could mean the difference between a 5 minute and 2 minute iteration cycle.

Add the following TOML config to `core/scripts/cre/environment/configs/workflow-don.toml`:
```toml
[Billing]
URL = 'host.docker.internal:2223'
TLSEnabled = false
```

Outside of the local-CRE, the workflow registry chain ID and address is pulled from
[Capabilities.WorkflowRegistry]. We don't yet have that figured out in the local-CRE
(it will spin up relayers), so instead you want to mimic the following diffs:

```go
@@ -535,7 +537,8 @@ func NewApplication(ctx context.Context, opts ApplicationOpts) (Application, err
                creServices.workflowRateLimiter,
                creServices.workflowLimits,
                workflows.WithBillingClient(billingClient),
-               workflows.WithWorkflowRegistry(cfg.Capabilities().WorkflowRegistry().Address(), cfg.Capabilities().WorkflowReg
istry().ChainID()),
+               // tmp to avoid booting up relayers
+               workflows.WithWorkflowRegistry("0xA15BB66138824a1c7167f5E85b957d04Dd34E468", "11155111"),
```

and

```
@@ -957,7 +960,7 @@ func newCREServices(
                                        workflowLimits,
                                        artifactsStore,
                                        syncer.WithBillingClient(billingClient),
-                                       syncer.WithWorkflowRegistry(capCfg.WorkflowRegistry().Address(), capCfg.WorkflowRegist
ry().ChainID()),
+                                       syncer.WithWorkflowRegistry("0xA15BB66138824a1c7167f5E85b957d04Dd34E468", "11155111"),
                                )
```

The happy-path:
* workflow runs successfully
* no `switch to metering mode` error logs generated by the workflow run
* SubmitWorkflowReceipt in billing service does not have an err message

---

## DX Tracing

To track environment usage and quality metrics (success/failure rate, startup time) local CRE environment is integrated with DX. If you have `gh cli` configured and authenticated on your local machine it will be used to automatically setup DX integration in the background. If you don't, tracing data will be stored locally in `~/.local/share/dx/` and uploaded once either `gh cli` is available or valid `~/.local/share/dx/config.json` file appears.

> Minimum required version of the `GH CLI` is `v2.50.0`

To opt out from tracing use the following environment variable:
```bash
DISABLE_DX_TRACKING=true
```

### Manually creating config file

Valid config file has the following content:
```json
{
  "dx_api_token":"xxx",
  "github_username":"your-gh-username"
}
```

DX API token can be found in 1 Password in the engineering vault as `DX - Local CRE Environment`.

Other environment variables:
* `DX_LOG_LEVEL` -- log level of a rudimentary logger
* `DX_TEST_MODE` -- executes in test mode, which means that data sent to DX won't be included in any reports
* `DX_FORCE_OFFLINE_MODE` -- doesn't send any events, instead saves them on the disk

---

# Job Distributor Image

Tests require a local Job Distributor image. By default, configs expect version `job-distributor:0.12.7`.

To build locally:
```bash
git clone https://github.com/smartcontractkit/job-distributor
cd job-distributor
git checkout v0.12.7
docker build -t job-distributor:0.12.7 -f e2e/Dockerfile.e2e .
```

If you pull the image from the PRO ECR remember to either update the image name in [TOML config](./configs/) for your chosed topology or to tag that image as `job-distributor:0.12.7`.

## Example Workflows

The environment includes several example workflows located in `core/scripts/cre/environment/examples/workflows/`:

### Available Workflows

#### V2 Workflows
- **`v2/cron/`**: Simple cron-based workflow that executes on a schedule
- **`v2/node-mode/`**: Node mode workflow example
- **`v2/http/`**: HTTP-based workflow example

#### V1 Workflows
- **`v1/proof-of-reserve/cron-based/`**: Cron-based proof-of-reserve workflow
- **`v1/proof-of-reserve/web-trigger-based/`**: Web API trigger-based proof-of-reserve workflow

### Deployable Example Workflows

The following workflows can be deployed using the `workflow deploy-and-verify-example` command:

#### Proof-of-Reserve Workflows
Both proof-of-reserve workflows execute a proof-of-reserve-like scenario with the following steps:
- Call external HTTP API and fetch value of test asset
- Reach consensus on that value
- Write that value in the consumer contract on chain

**Usage:**
```bash
go run . workflow deploy-and-verify-example [flags]
```

**Key flags:**
- `-y, --example-workflow-trigger`: Trigger type (`web-trigger` or `cron`, default: `web-trigger`)
- `-u, --example-workflow-timeout`: Time to wait for workflow execution (default: `5m`)
- `-g, --gateway-url`: Gateway URL for web API trigger (default: `http://localhost:5002`)
- `-d, --don-id`: DON ID for web API trigger (default: `vault`)
- `-w, --workflow-registry-address`: Workflow registry address (default: `0x9fE46736679d2D9a65F0992F2272dE9f3c7fa6e0`)
- `-r, --rpc-url`: RPC URL (default: `http://localhost:8545`)

**Examples:**
```bash
# Deploy cron-based proof-of-reserve workflow
go run . workflow deploy-and-verify-example -y cron

# Deploy web-trigger-based proof-of-reserve workflow with custom timeout
go run . workflow deploy-and-verify-example -y web-trigger -u 10m
```

#### Cron-based Workflow
- **Trigger**: Every 30 seconds on a schedule
- **Behavior**: Keeps executing until paused or deleted
- **Requirements**: External `cron` capability binary (must be manually compiled or downloaded and configured in TOML)
- **Source**: [`examples/workflows/v1/proof-of-reserve/cron-based/main.go`](./examples/workflows/v1/proof-of-reserve/cron-based/main.go)

#### Web API Trigger-based Workflow
- **Trigger**: Only when a precisely crafted and cryptographically signed request is made to the gateway node
- **Behavior**: Triggers workflow **once** and only if:
  - Sender is whitelisted in the workflow
  - Topic is whitelisted in the workflow
- **Source**: [`examples/workflows/v1/proof-of-reserve/web-trigger-based/main.go`](./examples/workflows/v1/proof-of-reserve/web-trigger-based/main.go)

**Note**: You might see multiple attempts to trigger and verify the workflow when running the example. This is expected and could happen because:
- Topic hasn't been registered yet (nodes haven't downloaded the workflow yet)
- Consensus wasn't reached in time

### Manual Workflow Deployment

For other workflows (v2/cron, v2/node-mode, v2/http), you can deploy them manually using the `workflow deploy` command:

```bash
# Deploy v2 cron workflow
go run . workflow deploy -w ./examples/workflows/v2/cron/main.go -n cron-workflow

# Deploy v2 http workflow
go run . workflow deploy -w ./examples/workflows/v2/http/main.go -n http-workflow

# Deploy v2 node-mode workflow
go run . workflow deploy -w ./examples/workflows/v2/node-mode/main.go -n node-mode-workflow
```

## Troubleshooting

### Docker fails to download public images
Make sure you are logged in to Docker. Run: `docker login`

```
Error: failed to setup test environment: failed to create blockchains: failed to deploy blockchain: create container: Error response from daemon: Head "https://ghcr.io/v2/foundry-rs/foundry/manifests/stable": denied: denied
```
your ghcr token is stale. do
```
docker logout ghcr.io
docker pull ghcr.io/foundry-rs/foundry
```
and try starting the environment

### GH CLI is not installed
Either download from [cli.github.com](https://cli.github.com) or install with Homebrew with:
```bash
brew install gh
```

Once installed, configure it by running:
```bash
gh auth login
```

For GH CLI to be used by the environment to download the CRE CLI you must have access to [smartcontract/dev-platform](https://github.com/smartcontractkit/dev-platform) repository.