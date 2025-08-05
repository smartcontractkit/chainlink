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

Use `#topic-local-dev-environments` for help.

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
go run . env start-beholder
```

> This assumes you have `chip-ingress:qa-latest` Docker image on your local machine. Without it Beholder won't be able to start. If you do not, close the [Atlas](https://github.com/smartcontractkit/atlas) repository, and then in `atlas/chip-ingress` run `docker build -t chip-ingress:qa-latest .`

### Storage

The environment supports two storage backends for workflow uploads:
- Gist (remote)
- S3 MinIO (built-in, local)

Configuration details are generated automatically into the `cre.yaml` file
(path is printed after starting the environment).

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
   - For a single DON with all capabilities: `configs/single-don.toml` (default)
   - For a full topology (workflow DON + capabilities DON + gateway DON): `configs/workflow-capabilities-don.toml`
2. **Download or Build Capability Binaries**
   - Some capabilities like `cron`, `log-event-trigger`, or `read-contract` are not embedded in all Chainlink images.
   - If your use case requires them, you can either:
      - Download binaries from [smartcontractkit/capabilities](https://github.com/smartcontractkit/capabilities/releases/tag/v1.0.2-alpha) release page or
      - Use GH CLI to download them, e.g. `gh release download v1.0.2-alpha --repo smartcontractkit/capabilities --pattern 'amd64_cron' && mv amd64_cron cron`
      Make sure they are built for `linux/amd64`!

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
- `CTF_CONFIGS`: TOML config paths. Defaults to [./configs/single-don.toml](./configs/single-don.toml)
- `PRIVATE_KEY`: Plaintext private key that will be used for all deployments (needs to be funded). Defaults to `ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80`
- `TESTCONTAINERS_RYUK_DISABLED`: Set to "true" to disable cleanup. Defaults to `false`

When starting the environment in AWS-managed Kubernetes make sure to source `.env` environment from the `crib/deployments/cre` folder specific for AWS. Remember, that it must include ingress domain settings.

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

## Example workflows

Two example workflows are available. Both execute a proof-of-reserve-like scenario with following steps:
- call external HTTP API and fetch value of test asset
- reach consensus on that value
- write that value in the consumer contract on chain

The only difference between is the trigger.

### cron-based workflow
This workflow is triggered every 30s, on a schedule. It will keep executing until it is paused or deleted. It requires an external `cron` capability binary, which you have to either manually compile or download **and** a manual TOML config change to indicate its location.

Source code can be found in [proof-of-reserves-workflow-e2e-test](https://github.com/smartcontractkit/proof-of-reserves-workflow-e2e-test/blob/main/cron-based/main.go) repository.

### web API trigger-based workflow
This workflow is triggered only, when a precisely crafed and cryptographically signed request is made to the gateway node. It will only trigger the workflow **once** and only if:
* sender is whitelisted in the workflow
* topic is whitelisted in the workflow

Source code can be found in [proof-of-reserves-workflow-e2e-test](https://github.com/smartcontractkit/proof-of-reserves-workflow-e2e-test/blob/main/web-api-trigger-based/main.go) repository.

You might see multiple attempts to trigger and verify that workflow, when running the example. This is expected and could be happening, because:
- topic hasn't been registered yet (nodes haven't downloaded the workflow yet)
- consensus wasn't reached in time

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