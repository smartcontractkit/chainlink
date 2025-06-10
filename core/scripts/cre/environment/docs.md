# Local CRE environment

## Table of content

1. [Using the CLI](#using-the-cli)
   - [Prerequisites](#prerequisites-for-docker)
   - [Start Environment](#start-environment)
   - [Stop Environment](#stop-environment)
   - [Restart Environment](#restarting-the-environment)

2. [Job Distributor Image](#job-distributor-image)
3. [Example Workflows](#example-workflows)
4. [Troubleshooting](#troubleshooting)

# Using the CLI

The CLI manages CRE test environments. It is located in `core/scripts/cre/environment`.

## Prerequisites (for Docker) ###
1. **Docker installed and running**
    - with usage of default Docker socket **enabled**
    - with Apple Virtualization framework **enabled**
    - with VirtioFS **enabled**
    - with use of containerd for pulling and storing images **disabled**
2. **Job Distributor Docker image available**
    - Either build it locally as described in [this section](#job-distributor-image)
    - or download it from the PROD ECR (if you have access) and tag as `job-distributor:0.9.0`

If you want to run an example workflow you also need to:

1. **Download CRE CLI v0.2.0**
    - download it either from [smartcontract/dev-platform](https://github.com/smartcontractkit/dev-platform/releases/tag/v0.2.0)
    - or using GH CLI by running: `gh release download v0.2.0 --repo smartcontractkit/dev-platform --pattern 'cre_v0.2.0_darwin_arm64.tar.gz'`
    - once you have the archive downloaded run:
      ```bash
      tar -xf cre_v0.2.0_darwin_arm64.tar.gz
      rm cre_v0.2.0_darwin_arm64.tar.gz
      # do not worry about potential 'No such xattr: com.apple.quarantine' error
      xattr -d com.apple.quarantine cre_v0.2.0_darwin_arm64
      export "PATH=$(pwd):$PATH"
      ```

Optionally:
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

## Start Environment
```bash
# while in core/scripts/cre/environment
go run main.go env start

# to start environment with an example workflow web API-based workflow
go run main.go env start --with-example

 # to start environment with an example workflow cron-based workflow (this requires the `cron` capability binary to be setup in the `extra_capabilities` section of the TOML config)
go run main.go env start --with-example --example-workflow-trigger cron

# to start environment using image with all supported capabilities
go run main.go env start --with-plugins-docker-image <SDLC_ACCOUNT_ID>dkr.ecr.<SDLC_ACCOUNT_REGION>.amazonaws.com/chainlink:nightly-<YYYMMDD>-plugins
```

> Important! **Nightly** Chainlink images are retained only for one day and built at 03:00 UTC. That means that in most cases you should use today's image, not yesterday's.

Optional parameters:
- `-t`: Topology (`simplified` or `full`)
- `-w`: Wait on error before removing up Docker containers (e.g. to inspect Docker logs, e.g. `-w 5m`)
- `-e`: Extra ports for which external access by the DON should be allowed (e.g. when making API calls or downloading WASM workflows)
- `-x`: Registers an example PoR workflow using CRE CLI and verifies it executed successfuly
- `-s`: Time to wait for example workflow to execute successfuly (defaults to `5m`)
- `-p`: Docker `plugins` image to use (must contain all of the following capabilities: `ocr3`, `cron`, `readcontract` and `logevent`)
- `-y`: Trigger for example workflow to deploy (web-trigger or cron). Default: `web-trigger`. **Important!** `cron` trigger requires user to either provide the capbility binary path in TOML config or Docker image that has it baked in.

### Using existing Docker Plugins image

If you don't want to build Chainlink image from your local branch (default behaviour) or you don't want to go through the hassle of downloading capabilities binaries in order to enable them on your environment you should use the `--with-plugins-docker-image` flag. It is recommended to use a nightly `core plugins` image that's build by [Docker Build action](https://github.com/smartcontractkit/chainlink/actions/workflows/docker-build.yml) as it contains all supported capability binaries.

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

## Job Distributor Image

Tests require a local Job Distributor image. By default, configs expect version `job-distributor:0.9.0`.

To build locally:
```bash
git clone https://github.com/smartcontractkit/job-distributor
cd job-distributor
git checkout v0.9.0
docker build -t job-distributor:0.9.0 -f e2e/Dockerfile.e2e .
```

If you pull the image from the PRO ECR remember to either update the image name in [TOML config](./configs/) for your chosed topology or to tag that image as `job-distributor:0.9.0`.

## Example workflows

Two example workflows are available. Both execute a proof-of-reserve-like scenario with following steps:
- call external HTTP API and fetch value of test asset
- reach consensus on that value
- write that value in the consumer contract on chain

The only difference between is the trigger.

### cron-based workflow
This workflow is triggered every 30s, on a schedule. It will keep executing until it is paused or deleted. It requires an external `cron` capability binary, which you have to either manually compile or download **and** a manual TOML config change to indicate its location.

Source code can be found in [proof-of-reserves-workflow-e2e-test](https://github.com/smartcontractkit/proof-of-reserves-workflow-e2e-test/main/dx-891-web-api-trigger-workflow/cron-based/main.go) repository.

### web API trigger-based workflow
This workflow is triggered only, when a precisely crafed and cryptographically signed request is made to the gateway node. It will only trigger the workflow **once** and only if:
* sender is whitelisted in the workflow
* topic is whitelisted in the workflow

Source code can be found in [proof-of-reserves-workflow-e2e-test](https://github.com/smartcontractkit/proof-of-reserves-workflow-e2e-test/blob/dx-891-web-api-trigger-workflow/web-api-trigger-based/main.go) repository.

You might see multiple attempts to trigger and verify that workflow, when running the example. This is expected and could be happening, because:
- topic hasn't been registered yet (nodes haven't downloaded the workflow yet)
- consensus wasn't reached in time

## Troubleshooting

### Docker fails to download public images
Make sure you are logged in to Docker. Run: `docker login`

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