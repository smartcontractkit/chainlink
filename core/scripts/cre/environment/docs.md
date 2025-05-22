# Local CRE environment

## Table of content

1. [Using the CLI](#using-the-cli)
   - [Prerequisites](#prerequisites-for-docker)
   - [Start Environment](#start-environment)
   - [Stop Environment](#stop-environment)
   - [Restart Environment](#restarting-the-environment)

2. [Job Distributor Image](#job-distributor-image)

# Using the CLI

The CLI manages CRE test environments. It lives in `core/scripts/cre/environment`.

## Prerequisites (for Docker) ###
1. **Docker installed and running**
    - with usage of default Docker socket **enabled**
    - with Apple Virtualization framework **enabled**
    - with VirtioFS **enabled**
    - with use of containerd for pulling and storing images **disabled**
2. **Logged in to Docker**
    - Run `docker login`
3. **Job Distributor Docker image available**
    - [This section](#job-distributor-image) explains how to build it locally
4. **Download CRE CLI v0.2.0**
    - download it from [smartcontract/dev-platform](https://github.com/smartcontractkit/dev-platform/releases/tag/v0.2.0) or
    - using GH CLI: `gh release download v0.2.0 --repo smartcontractkit/dev-platform --pattern '*darwin_arm64*'`
    - remove it from MacOs' quarantine `xattr -d com.apple.quarantine cre_v0.2.0_darwin_arm64`

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
       cron_capability_binary_path = "../cron" # remember to adjust binary name
       ```
   - If the capability is already baked into your CL image (check the Dockerfile), comment out the TOML path line to skip copying.
3. **Ensure Binaries Are in the Right Location**
    - Default config of the CLI command will look for `cron`, and other capability binaries in `core/scripts/cre/environment`
4.  **Decide whether to build or reuse Chainlink Docker Image**
   - To build from your local branch:
     ```toml
     [nodesets.node_specs.node]
     docker_ctx = "../../../.."
     docker_file = "plugins/chainlink.Dockerfile"
     ```
   - To reuse a prebuilt image:
     ```toml
     [nodesets.node_specs.node]
     image = "<your-Docker-image>:<your-tag>"
     ```
  Make these changes for **all** nodes in the nodeset.

5. **Decide whether to use Docker or k8s**
    - Read sections 3 to 9 starting [here](#2-docker-vs-kubernetes-k8s) to learn how to switch between Docker and Kubernetes
6. **Start Observability Stack (Docker-only)**
   - If you want Grafana/Prometheus support, run:
     ```bash
     ctf obs up
     ```
   - If you want Blockscout block explorer, run:
    ```bash
    ctf bs u
    ```
    (warning, that stack is pretty heavy)
    - To download the `ctf` binary follow the steps described [here](https://smartcontractkit.github.io/chainlink-testing-framework/framework/getting_started.html)
7. **Download and configure GH CLI (if CRE CLI  is missing or not in your PATH)**
  - Either download from [cli.github.com](https://cli.github.com/) or install with Homebrew with `brew install gh`
  - Configure with `gh auth login`

Optional environment variables used by the CLI:
- `CTF_CONFIGS`: TOML config path
- `PRIVATE_KEY`: Default test key if not set
- `TESTCONTAINERS_RYUK_DISABLED`: Set to "true" to disable cleanup

When starting the environment in AWS-managed Kubernetes make sure to source `.env` environment from the `crib/deployments/cre` folder specific for AWS. Remember, that it must include ingress domain settings.

---

## Start Environment
```bash
# while in core/scripts/cre/environment
go run main.go env start

# or to start environment with an example workflow
go run main.go env start --with-example
```

Optional parameters:
- `-t`: Topology (`simplified` or `full`)
- `-w`: Wait on error before cleanup (e.g. to inspect Docker logs, e.g. `-w 5m`)
- `-e`: Extra ports for which external access by the DON should be allowed (e.g. when making API calls)
- `-x`: Registers an example PoR workflow using CRE CLI and verifies it
- `-s`: Time to wait for example workflow to execute successfuly (defaults to `5m`)
- `-p`: Docker `Plugins` image to use (must contain all of the following capabilities: `ocr3`, `cron`, `readcontract` and `logevent`)


### Using existing Docker Plugins image

If you don't want to build Chainlink image from your local branch (default behaviour) or you don't want to go through the hassle of downloading capabilities binaries in order to enable them on your environment you should use the `--with-plugins-docker-image` flag. It is recommended to use a nightly `core plugins` image that's build by [Docker Build action](https://github.com/smartcontractkit/chainlink/actions/workflows/docker-build.yml) as it contains all supported capability binaries.

## Stop Environment
```bash
# while in core/scripts/cre/environment
go run main.go env stop
```

Or... if you have the CTF binary:
```
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
- using GH CLI: `gh release download v0.2.0 --repo smartcontractkit/dev-platform --pattern '*darwin_arm64*'`

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

Or pull from your internal registry and update the image name in `environment-*.toml`.