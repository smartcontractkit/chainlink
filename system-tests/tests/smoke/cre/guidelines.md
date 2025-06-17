# Test Modification and Execution Guide

## Table of Contents

1. [Running the Test](#1-running-the-test)
   - [Chainlink Node Image](#chainlink-node-image)
   - [Environment Variables](#environment-variables)
   - [Job Distributor Image](#job-distributor-image)
   - [CRE CLI Binary](#cre-cli-binary)
   - [cron Capability Binary](#cron-capability-binary)
   - [PoR Workflow Source Code](#por-workflow-source-code)
   - [Test Timeout](#test-timeout)
   - [Visual Studio Code Debug Configuration](#visual-studio-code-debug-configuration)
2. [Using the CLI](#2-using-the-cli)
3. [Docker vs Kubernetes (k8s)](#3-docker-vs-kubernetes-k8s)
4. [CRIB Requirements](#4-crib-requirements)
5. [Setting Docker Images for CRIB Execution](#5-setting-docker-images-for-crib-execution)
6. [Running Tests in Local Kubernetes (kind)](#6-running-tests-in-local-kubernetes-kind)
7. [CRIB Deployment Flow](#7-crib-deployment-flow)
8. [Switching from kind to AWS Provider](#8-switching-from-kind-to-aws-provider)
9. [CRIB Limitations & Considerations](#9-crib-limitations--considerations)
10. [Adding a New Capability](#10-adding-a-new-capability)
11. [Using a New Workflow](#11-using-a-new-workflow)
12. [Deployer Address or Deployment Sequence Changes](#12-deployer-address-or-deployment-sequence-changes)
13. [Adding a New Test to the CI](#13-adding-a-new-test-to-the-ci)
14. [Multiple DONs](#14-multiple-dons)
15. [Price Data Source](#15-price-data-source)
16. [Using a Specific Docker Image for Chainlink Node](#16-using-a-specific-docker-image-for-chainlink-node)
17. [Using Existing EVM & P2P Keys](#17-using-existing-evm--p2p-keys)
18. [Troubleshooting](#18-troubleshooting)

---

## Test Modification and Execution Guide

This guide explains how to set up and run system tests for Chainlink workflows using the CRE (Composable Runtime Environment) framework. It includes support for Docker and Kubernetes (via CRIB), multiple capabilities, and integration with Chainlink nodes and job distributor services.

---

## 1. Running the Test

Before starting, you’ll need to configure your environment correctly.

### Chainlink Node Image

The TOML config defines how Chainlink node images are used:

- **Default behavior**: Builds the Docker image from your current branch.
  ```toml
  [nodesets.node_specs.node]
    docker_ctx = "../../../.."
    docker_file = "plugins/chainlink.Dockerfile"
  ```

- **Using a pre-built image**: Replace the config with:
  ```toml
  [nodesets.node_specs.node]
    image = "my-docker-image:my-tag"
  ```
  Apply this to every `nodesets.node_specs.node` section.

**Minimum required version**: Commit [e13e5675](https://github.com/smartcontractkit/chainlink/commit/e13e5675d3852b04e18dad9881e958066a2bf87a) (Feb 25, 2025)

---

### Environment Variables

Set these before running your test:

- `CTF_CONFIGS`: Required. Comma-separated list of TOML config files.
- `PRIVATE_KEY`: Required. Plaintext private key for contract deployment and node funding.
- `GIST_WRITE_TOKEN`: Optional. Needed only when uploading a new workflow. Must have `gist:read:write` scope on GitHub.

---

### Job Distributor Image

Tests require a local Job Distributor image. By default, configs expect version `job-distributor:0.9.0`.

To build locally:
```bash
git clone https://github.com/smartcontractkit/job-distributor
cd job-distributor
git checkout v0.9.0
docker build -t job-distributor:0.9.0 -f e2e/Dockerfile.e2e .
```

Or pull from your internal registry and update the image name in `environment-*.toml`.

---

### CRE CLI Binary

You can build it locally from [smartcontractkit/dev-platform](https://github.com/smartcontractkit/dev-platform) repo, download from [v0.0.2 release page](https://github.com/smartcontractkit/dev-platform/releases/tag/v0.2.0) or download using GH CLI: `gh release download v0.2.0 --repo smartcontractkit/dev-platform --pattern '*darwin_arm64*'`.

**Required version**: ` v0.2.0`

---

### `cron` Capability Binary

This binary is needed for tests using the cron capability.

**Option 1**: Use a CL node image that already includes the binary. If so, comment this in TOML:
```toml
[workflow_config.dependencies]
  # cron_capability_binary_path = "./cron"
```

**Option 2**: Provide a path to a locally built binary:
```toml
[workflow_config.dependencies]
  cron_capability_binary_path = "./some-folder/cron"
```

You can build it from [smartcontractkit/capabilities](https://github.com/smartcontractkit/capabilities) repo, download from [v1.0.2-alpha release page](https://github.com/smartcontractkit/capabilities/releases/tag/v1.0.2-alpha) or use GH CLI to download them, e.g. `gh release download v1.0.2-alpha --repo smartcontractkit/capabilities --pattern 'amd64_cron'`

**Note**: Binary must be compiled for **Linux** and **amd64**.

---

### PoR Workflow Source Code

By default, tests compile the workflow at runtime. Clone the repo:
```bash
git clone https://github.com/smartcontractkit/proof-of-reserves-workflow-e2e-test
```

Then, update the TOML:
```toml
[workflow_config]
  workflow_folder_location = "/path/to/proof-of-reserves-workflow-e2e-test"
```

---

### Test Timeout

- If building the image: Set Go test timeout to **20 minutes**.
- If using pre-built images: Execution takes **4–7 minutes**.

---

### Visual Studio Code Debug Configuration

Example `launch.json` entry:
```json
{
  "name": "Launch Capability Test",
  "type": "go",
  "request": "launch",
  "mode": "test",
  "program": "${workspaceFolder}/system-tests/tests/smoke/cre",
  "env": {
    "CTF_CONFIGS": "environment-one-don.toml",
    "PRIVATE_KEY": "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80",
    "GIST_WRITE_TOKEN": "xxxx"
  },
  "args": [
    "-test.run",
    "TestCRE_OCR3_PoR_Workflow_SingleDon_MockedPrice"
  ]
}
```

**CI behavior differs**: In CI, workflows and binaries are uploaded ahead of time, and images are injected via:
- `E2E_JD_VERSION`
- `E2E_TEST_CHAINLINK_IMAGE`
- `E2E_TEST_CHAINLINK_VERSION`

---

## 2. Using the CLI

Local CRE environment and documentation were migrated to [core/scripts/cre/environment/docs.md](../../../../core/scripts/cre/environment/docs.md).

---

### Further use
To manage workflows you will need the CRE CLI. You can either:
- download it from [smartcontract/dev-platform](https://github.com/smartcontractkit/dev-platform/releases/tag/v0.2.0) or
- using GH CLI: `gh release download v0.2.0 --repo smartcontractkit/dev-platform --pattern '*darwin_arm64*'`

Remember that the CRE CLI version needs to match your CPU architecture and operating system.

---

## 3. Docker vs Kubernetes (k8s)

The environment type is set in your TOML config:
```toml
[infra]
  type = "crib"  # Options: "docker" or "crib"
```

To run tests in Kubernetes, you must use the `crib` option. CRIB supports both:
- **Local cluster (`kind`)**
- **AWS cloud provider**

Example TOML for CRIB:
```toml
[infra.crib]
  namespace = "crib-local"
  folder_location = "/absolute/path/to/crib/deployments/cre"
  provider = "kind"  # or "aws"
```

---

## 4. CRIB Requirements

Before using CRIB, ensure the following:

1. **Read the CRIB Setup Guide**
   Follow the official [CRIB deployment instructions](https://smartcontract-it.atlassian.net/wiki/spaces/INFRA/pages/660145339/General+CRIB+-+Deploy+Access+Instructions).

2. **AWS Role (for AWS provider)**
   Required only for AWS. Local `kind` setup does not require role access.

3. **Pull Local Registry Image** (for `kind` only):
   ```bash
   docker pull registry:2
   ```

4. **Clone CRIB Repository**
   ```bash
   git clone https://github.com/smartcontractkit/crib
   cd crib
   pwd  # to get absolute path for config
   ```

5. **Update `folder_location` in TOML**:
   ```toml
   folder_location = "/your/absolute/path/to/crib/deployments/cre"
   ```

6. **Add Cost Attribution (for AWS)**:
   ```toml
   [infra.crib.team_input]
     team = "your-team"
     product = "product-name"
     cost_center = "crib"
     component = "crib"
   ```

7. **Connect VPN** (for AWS provider only)

---

## 5. Setting Docker Images for CRIB Execution

CRIB does **not** support building Docker images from source during test runtime.

❌ Not allowed:
```toml
[nodesets.node_specs.node]
  docker_ctx = "../../../.."
  docker_file = "plugins/chainlink.Dockerfile"
```

✅ Required:
```toml
[nodesets.node_specs.node]
  image = "localhost:5001/chainlink:112b9323-plugins-cron"
```

`localhost:5001` is the repository name for local Kind registry. In order to push your image there you need:
- **tag a Docker image with prefix**, e.g. `docker tag chainlink-tmp:latest localhost:5001/chainlink:latest`
- **pushh it to the local registry**, e.g. `docker push localhost:5001/chainlink:latest`

**Also, it is crucial that you use an image, where default user is `chainlink`**. That's because Helm charts used for k8s deployment will start that container as that user and if your image was created using a different default user (e.g. `root`), then Chainlink application won't even start due to incorrect filesystem permissions. If you are building the image locally, use the following command:
```bash
# run in root of chainlink repository
docker build -t localhost:5001/chainlink:<your-tag>> --build-arg CHAINLINK_USER=chainlink -f plugins/chainlink.Dockerfile
```

### Notes:
- All nodes in a single nodeset **must** use the same image.
- You must specify an image tag explicitly (e.g., `:v1.2.3`).

### Job Distributor (JD) Image in CRIB

CRIB extracts only the image **tag** from your TOML:
```toml
[jd]
  image = "jd-test-1:my-awesome-tag"
```

If you leave off the tag, CRIB will fail:
```toml
[jd]
  image = "jd-test-1"  # ❌ This will fail
```

---

## 6. Running Tests in Local Kubernetes (`kind`)

### Docker Registry Setup

Pull the required local registry image:
```bash
docker pull registry:2
```

### Hostname Routing with `/etc/hosts`

CRIB dynamically creates hostname entries, but modifying `/etc/hosts` requires root. Tests fail if routing isn't set up.

#### Solutions:
1. **Manually add host entries**.
2. **Run `devspace` manually before starting tests** to allow interactive root access.

#### Example Manual Entries

**Geth Chain**:
```bash
127.0.0.1 crib-local-geth-1337-http.main.stage.cldev.sh
127.0.0.1 crib-local-geth-1337-ws.main.stage.cldev.sh
```

**Job Distributor**:
```bash
127.0.0.1 crib-local-job-distributor-grpc.main.stage.cldev.sh
```

**Chainlink Nodes (1 bootstrap, 3 workers)**:
```bash
127.0.0.1 crib-local-workflow-bt-0.main.stage.cldev.sh
127.0.0.1 crib-local-workflow-0.main.stage.cldev.sh
127.0.0.1 crib-local-workflow-1.main.stage.cldev.sh
127.0.0.1 crib-local-workflow-2.main.stage.cldev.sh
```

### Automating Host Setup with `devspace`

From within `cre/deployment` and a `nix develop` shell:

**Deploy Geth Chain**:
```bash
CHAIN_ID=<id> devspace run deploy-custom-geth-chain
```

**Deploy JD**:
```bash
devspace run deploy-jd
```

**Deploy DON**:
```bash
DON_TYPE=<name> DON_NODE_COUNT=<n> DON_BOOT_NODE_COUNT=<b> devspace run deploy-don
```

Ensure DON type matches the `name` field in your TOML config.

---

## 7. CRIB Deployment Flow

1. **Start a `nix develop` shell**, set:
   - `PROVIDER`, `DEVSPACE_NAMESPACE`, `CONFIG_OVERRIDES_DIR`

2. **Deploy chains**:
   ```bash
   CHAIN_ID=<id>
   devspace run deploy-custom-geth-chain
   ```
   Read endpoints from `chain-<CHAIN_ID>-urls.json`

3. **Deploy Keystone contracts**

4. **Generate CL node configs and secrets** (`./crib-configs`)

5. **Start each DON**:
   - Set `DEVSPACE_IMAGE`, `DEVSPACE_IMAGE_TAG`, `DON_BOOT_NODE_COUNT`, `DON_NODE_COUNT`, `DON_TYPE`
   - Run: `devspace run deploy-don`
   - Get URLs from `don-<DON_TYPE>-urls.json`
   - Copy binaries (if needed): `devspace run copy-to-pods`

6. **Deploy JD**:
   - Set `JOB_DISTRIBUTOR_IMAGE_TAG`
   - Run: `devspace run deploy-jd`
   - Get JD URLs from `jd-url.json`

7. **Create jobs and configure CRE contracts**
   - Same as Docker-based flow

...

---

## 8. Switching from `kind` to AWS Provider

When switching from a local `kind` setup to AWS:

- **Remove `/etc/hosts` entries** added by CRIB for the `kind` namespace. If reused, these entries will redirect traffic incorrectly to `localhost`.
- **Use a new namespace** in your TOML config to avoid DNS conflicts.

Recommended: Always switch namespaces when changing providers.

---

## 9. CRIB Limitations & Considerations

### Gateway DON
- Must run on a **dedicated node**
- Use `DON_TYPE = gateway`
- No bootstrap node required, but worker nodes are supported

### Mocked Price Provider
- Not supported in CRIB
- CRIB can only use **live endpoints**

### Environment Variables
- Some are set by Go code, others by `.env` in `deployments/cre`
- Avoid overlapping values to prevent inconsistent behavior

### DNS Propagation (AWS only)
- DNS may take time to propagate
- If tests fail early, retry after a few minutes

### Ingress Check (local `kind` only)
- May fail even when the environment is healthy
- If this happens, re-run the test

### Connectivity Troubleshooting
Check pod status:
```bash
kubectl get pods
```
Inspect logs:
```bash
kubectl logs <POD_NAME>
```

---

## 10. Adding a New Capability

To add a new capability (e.g., writing to the Aptos chain), follow these detailed steps:

### 1. Define the Capability Flag

Create a unique flag in `flags.go` that represents your capability. This is used throughout the test framework:
```go
const (
  WriteAptosCapability CapabilityFlag = "write-aptos" // <--- NEW ENTRY
)
```

### 2. Copy the Binary to the Container/Pod

Use the TOML config or inject the binary programmatically. The latter is recommended so you can reuse the binary path later in your job spec factory function:

```go
customBinariesPaths := map[string]string{}
containerPath, pathErr := capabilities.DefaultContainerDirectory(in.Infra.InfraType)
require.NoError(t, pathErr, "failed to get container directory")

var aptosBinaryPathInTheContainer string
if in.WorkflowConfig.DependenciesConfig.AptosCapabilityBinaryPath != "" {
  aptosBinaryPathInTheContainer = filepath.Join(containerPath, filepath.Base(in.WorkflowConfig.DependenciesConfig.AptosCapabilityBinaryPath))
  customBinariesPaths[keystonetypes.AptosWriteCapability] = in.WorkflowConfig.DependenciesConfig.AptosCapabilityBinaryPath
} else {
  aptosBinaryPathInTheContainer = filepath.Join(containerPath, "aptos")
}
```

You can now pass `customBinariesPaths` and the constructed path to the `SetupInput`.

> Note: Bootstrap nodes do not run capabilities, so binaries are only copied to worker nodes.

### 3. Define Additional Node Configuration (Optional)

Some capabilities require node-specific TOML config. Here’s a sample:
```go
if hasFlag(flags, WriteAptosCapability) {
  workerNodeConfig += fmt.Sprintf(`
    [[Aptos]]
    ChainID = '%s'
    Enabled = true

    [[Aptos.Nodes]]
    Name = 'aptos'
    URL = '%s'

    [Aptos.TransactionManager]
    BroadcastChanSize = 100
    ConfirmPollSecs = 2
    DefaultMaxGasAmount = 200000
    MaxSimulateAttempts = 5
    MaxSubmitRetryAttempts = 5
    MaxTxRetryAttempts = 3
    PruneIntervalSecs = 14400
    PruneTxExpirationSecs = 7200
    SubmitDelayDuration = 3
    TxExpirationSecs = 30

    [Aptos.Workflow]
    ForwarderAddress = '%s'

    [Aptos.WriteTargetCap]
    ConfirmerPollPeriod = '300ms'
    ConfirmerTimeout = '30s'
  `, chainID, rpcURL, forwarderAddress)
}
```

### 4. Define the Job Spec

Use a factory function to generate the job spec dynamically:
```go
func AptosJobSpecFactoryFn(binaryPath string) keystonetypes.JobSpecFactoryFn {
  return func(ctx context.Context, node *types.Node, env *devenv.Environment) (*jobv1.ProposeJobRequest, error) {
    jobSpec := fmt.Sprintf(`
      type = "standardcapabilities"
      schemaVersion = 1
      externalJobID = "%s"
      name = "aptos-write-capability"
      command = "%s"
      config = ""
    `, uuid.NewString(), binaryPath)

    nodeID, _ := node.FindLabelValue(node, libnode.NodeIDKey)
    return &jobv1.ProposeJobRequest{NodeId: nodeID, Spec: jobSpec}, nil
  }
}
```

### 5. Register Capability in Capabilities Registry

Use a factory to register capabilities dynamically:
```go
func AptosCapabilityFactoryFn() keystonetypes.DONCapabilityWithConfigFactoryFn {
  return func(flags []string) []keystone_changeset.DONCapabilityWithConfig {
    if hasFlag(flags, WriteAptosCapability) {
      return []keystone_changeset.DONCapabilityWithConfig{
        {
          Capability: kcr.CapabilitiesRegistryCapability{
            LabelledName:   "write_aptos-testnet",
            Version:        "1.0.0",
            CapabilityType: 3,
            ResponseType:   1,
          },
          Config: &capabilitiespb.CapabilityConfig{},
        },
      }
    }
    return nil
  }
}
```

### 6. Update DON Topology

Update DON assignment with the new capability:
```go
dons := []*cretypes.CapabilitiesAwareNodeSet{
  {
    Input:              in.NodeSets[0],
    Capabilities:       []string{cretypes.OCR3Capability, cretypes.WriteAptosCapability},
    DONTypes:           []string{cretypes.WorkflowDON},
    BootstrapNodeIndex: 0,
  },
}
```

Ensure that the corresponding TOML defines the nodeset correctly, including node count and mode.

### 7. Pass all inputs to `setupInput`

Now, pass your custom factories to the `setupInput` and the universal setup function:
```go
universalSetupInput := creenv.SetupInput{
  CapabilitiesAwareNodeSets:            dons,
  CapabilitiesContractFactoryFunctions: []keystonetypes.DONCapabilityWithConfigFactoryFn{
    creconsensus.ConsensusJobSpecFactoryFn(chainIDUint64),
    AptosCapabilityFactoryFn,
  },
  BlockchainsInput:                     *in.BlockchainA,
  JdInput:                              *in.JD,
  InfraInput:                           *in.Infra,
  CustomBinariesPaths:                  customBinariesPaths,
  JobSpecFactoryFunctions: []keystonetypes.JobSpecFactoryFn{
    creconsensus.ConsensusJobSpecFactoryFn(chainIDUint64),
    AptosJobSpecFactoryFn(aptosBinaryPathInTheContainer),
  },
}

universalSetupOutput, setupErr := creenv.SetupTestEnvironment(t.Context(), testLogger, cldlogger.NewSingleFileLogger(t), universalSetupInput)
if setupErr != nil {
  panic(setupErr)
}
```

### 8. Run the Test

Once all pieces are configured, run the test as normal. Ensure that the logs show the capability was registered and the job executed successfully.

> Reminder: Capabilities and DON types are defined in Go. Infrastructure (images, ports) lives in TOML.

---

## 11. Using a New Workflow

You can test new workflows in two ways:

### Option 1: Test Uploads the Binary

The test itself compiles the workflow to a WASM binary and uploads it to Gist using the CRE CLI. This is the preferred method for local development.

Update your TOML config:
```toml
[workflow_config]
  use_cre_cli = true
  should_compile_new_workflow = true
  workflow_folder_location = "path-to-folder-with-main.go"
```

If your workflow accepts configuration parameters, define them in a Go struct and serialize to a temporary config file:
```go
configFile, err := os.CreateTemp("", "config.json")
require.NoError(t, err, "failed to create workflow config file")

type PoRWorkflowConfig struct {
  FeedID          string `json:"feedId"`
  URL             string `json:"url"`
  ConsumerAddress string `json:"consumerAddress"`
}

workflowConfig := PoRWorkflowConfig{
  FeedID:          feedID,
  URL:             "https://api.real-time-reserves.verinumus.io/v1/chainlink/proof-of-reserves/TrueUSD",
  ConsumerAddress: feedsConsumerAddress.Hex(),
}
```

> ⚠️ Currently, we treat config changes as a new workflow version and require re-upload.

If the workflow is not configurable, pass an empty byte slice `[]byte{}` during registration.

Workflows can be registered using:
```go
input := keystonetypes.ManageWorkflowWithCRECLIInput{
  NewWorkflow: workflow,
  ShouldCompileNewWorkflow: true,
}
err := creworkflow.RegisterWithCRECLI(input)
```

---

### Option 2: Manual Upload of the Binary

You can manually compile the workflow binary to WASM and upload it to an accessible location (e.g. Gist, S3, IPFS). Then reference it in your TOML config:
```toml
[workflow_config]
  use_cre_cli = true
  should_compile_new_workflow = false

  [workflow_config.compiled_config]
    binary_url = "https://gist.githubusercontent.com/user/binary.wasm.br.b64"
    config_url = "https://gist.githubusercontent.com/user/config.json"  # Optional
    secrets_url = "https://gist.githubusercontent.com/user/secrets.json"  # Optional
```

All URLs must be accessible from within the DON Gateway node.

---

### Workflow Secrets

To use secrets:
1. Create a `secrets.config.yaml` mapping secrets to env vars
2. Encrypt them using the CRE CLI and upload
3. Reference them in your workflow

> Secrets are encrypted using the CSA public key from each DON. They are not shareable across DONs unless they use the same keys.

Use `_ENV_VAR_ALL` suffix for shared secrets across nodes:
```go
newWorkflow := keystonetypes.NewWorkflow{
  SecretsFilePath:  secretsFilePath,
  Secrets: map[string]string{
    "FIRST_SECRET_ENV_VAR_ALL": "first secret value",
    "SECOND_SECRET_ENV_VAR_ALL": "second secret value",
  },
}
```

Register it:
```go
input := keystonetypes.ManageWorkflowWithCRECLIInput{
  NewWorkflow: newWorkflow,
  ShouldCompileNewWorkflow: true,
}

registerErr := creworkflow.RegisterWithCRECLI(input)
```

---

### YAML Workflows (Data Feeds DSL)

No compilation required. Define YAML workflow inline and propose it like any job:

```toml
type = "workflow"
schemaVersion = 1
name = "df-workflow"
externalJobID = "df-workflow-id"
workflow = """
name: df-workflow
owner: '0xabc...'
triggers:
 - id: streams-trigger@1.0.0
   config:
     maxFrequencyMs: 5000
     feedIds:
       - '0xfeed...'
consensus:
 - id: offchain_reporting@1.0.0
   ref: ccip_feeds
   inputs:
     observations:
       - $(trigger.outputs)
   config:
     report_id: '0001'
     key_id: 'evm'
     aggregation_method: data_feeds
     encoder: EVM
     encoder_config:
       abi: (bytes32 FeedID, uint224 Price, uint32 Timestamp)[] Reports
targets:
 - id: write_geth@1.0.0
   inputs:
     signed_report: $(ccip_feeds.outputs)
   config:
     address: '0xcontract...'
     deltaStage: 10s
     schedule: oneAtATime
"""
```

Then propose the job using JD, either directly:
```go
offChainClient.ProposeJob(ctx, &jobv1.ProposeJobRequest{NodeId: nodeID, Spec: workflowSpec})
```

Or using the `CreateJobs` helper:
```go
createJobsInput := keystonetypes.CreateJobsInput{
  CldEnv: env,
  DonTopology: donTopology,
  DonToJobSpecs: donToJobSpecs,
}
createJobsErr := libdon.CreateJobs(testLogger, createJobsInput)
```

---

## 12. Deployer Address or Deployment Sequence Changes

CI assumes a stable deployer address and deployment sequence for consistency. This ensures contracts like the Data Feeds Cache remain at fixed addresses.

### When to Trigger Re-upload
If either the deployer private key or deployment sequence changes, you must recompile and re-upload the workflow:
```toml
[workflow_config]
  use_cre_cli = true
  should_compile_new_workflow = true
  workflow_folder_location = "path-to-your-workflow"
```

After uploading, update the `config_url` in your CI-specific TOML config.

> Note: In most CI flows, existing pre-uploaded workflows are reused to avoid unnecessary binary changes.

---

## 13. Adding a New Test to the CI

CI currently does not support uploading to Gist during test execution. Workflow binaries and configurations must be pre-uploaded.

### Step 1: CI-Specific Configuration
Prepare a TOML file with:
```toml
[workflow_config]
  workflow_name = "my-workflow"
  feed_id = "018e..."
  use_cre_cli = true
  should_compile_new_workflow = false

  [workflow_config.compiled_config]
    binary_url = "https://gist.githubusercontent.com/user/.../binary.wasm.br.b64"
    config_url = "https://gist.githubusercontent.com/user/.../config.json"
```

> Use URLs from a local test run that compiles and uploads the workflow.

### Step 2: Add Entry in `.github/e2e-tests.yaml`
Example entry:
```yaml
- id: system-tests/smoke/cre/por_test.go:TestCRE_OCR3_PoR_Workflow_MultiDon
  path: system-tests/tests/smoke/cre/por_test.go
  test_env_type: docker
  runs_on: ubuntu-latest
  triggers:
    - PR CRE E2E Core Tests
    - Merge Queue CRE E2E Core Tests
  test_cmd: >
    cd tests && pushd smoke/cre/cmd > /dev/null && \
    go run main.go download all --output-dir ../ \
      --gh-token-env-var-name GITHUB_API_TOKEN \
      --cre-cli-version v0.1.5 --capabilities-name cron \
      --capabilities-version v1.0.2-alpha 1>&2 && \
    popd > /dev/null && \
    go test github.com/smartcontractkit/chainlink/system-tests/tests/smoke/cre \
      -v -run "^(TestCRE_OCR3_PoR_Workflow_MultiDon)$" -timeout 30m -count=1 -test.parallel=1 -json; \
    exit_code=$?; ../../tools/ci/wait-for-containers-to-stop.sh 30; exit $exit_code;
  pyroscope_env: ci-smoke-capabilities-evm-simulated
  test_env_vars:
    E2E_TEST_CHAINLINK_VERSION: "{{ env.DEFAULT_CHAINLINK_PLUGINS_VERSION }}"
    E2E_JD_VERSION: 0.9.0
    GITHUB_READ_TOKEN: "{{ env.GITHUB_API_TOKEN }}"
    CI: "true"
    CTF_CONFIGS: "my-ci-config.toml"
    PRIVATE_KEY: "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
  test_go_project_path: system-tests
```

#### Key Fields to Adjust
- `id`: `<relative-path>:<TestFunctionName>`
- `path`: Relative to repo root
- `test_cmd`: Include test runner, download helpers, and cleanup
- `test_env_vars`: Include all required secrets and flags

> Talk to DevEx if you need dynamic or sensitive env vars set securely.

---


... (previous content preserved)

---

## 14. Multiple DONs

You can configure multiple DONs (Decentralized Oracle Networks) by modifying your TOML config and Go code accordingly.

### Supported Capabilities
- `ocr3`
- `cron`
- `custom-compute`
- `write-evm`
- `read contract`
- `log-event-trigger` (under development)
- `web-api-trigger` (under development)
- `web-api-target` (under development)

### DON Types
- `workflow`
- `capabilities`
- `gateway`

You can only have **one** `workflow` and **one** `gateway` DON. You can define multiple `capabilities` DONs.

### Go Code (DON Assignment Example)
```go
[]*cretypes.CapabilitiesAwareNodeSet{
  {
    Input:              in.NodeSets[0],
    Capabilities:       []string{cretypes.OCR3Capability, cretypes.CustomComputeCapability},
    DONTypes:           []string{cretypes.WorkflowDON},
    BootstrapNodeIndex: 0,
  },
  {
    Input:              in.NodeSets[1],
    Capabilities:       []string{cretypes.WriteEVMCapability},
    DONTypes:           []string{cretypes.CapabilitiesDON},
    BootstrapNodeIndex: 0,
  },
  {
    Input:              in.NodeSets[2],
    Capabilities:       []string{},
    DONTypes:           []string{cretypes.GatewayDON},
    BootstrapNodeIndex: -1,
    GatewayNodeIndex:   0,
  },
}
```

> ⚠️ The framework does not validate capability compatibility with DON type. Some capabilities **must** run in the `workflow` DON (e.g. `ocr3`, `cron`, `custom-compute`).

### TOML Configuration
Use `override_mode = "all"` for uniform nodeset config:
```toml
[[nodesets]]
  nodes = 5
  override_mode = "all"
  name = "workflow"

  [[nodesets.node_specs]]
    [nodesets.node_specs.node]
      image = "localhost:5001/chainlink:112b9323-plugins-cron"
      user_config_overrides = """
      [Feature]
        LogPoller = true
      [OCR2]
        Enabled = true
        DatabaseTimeout = '1s'
      [P2P.V2]
        Enabled = true
        ListenAddresses = ['0.0.0.0:5001']
      """
```

For per-node customization (e.g. to import secrets), use `override_mode = "each"`.

---

## 15. Price Data Source

Supports both **live** and **mocked** price feeds.

### Live Source
```toml
[price_provider]
  feed_id = "018bfe8840700040000000000000000000000000000000000000000000000000"
  url = "api.real-time-reserves.verinumus.io/v1/chainlink/proof-of-reserves/TrueUSD"
```

> ⚠️ The PoR workflow only supports this response structure.

### Mocked Data Source
For local testing without live API dependency:
```toml
[price_provider]
  feed_id = "018bfe8840700040000000000000000000000000000000000000000000000000"
  [price_provider.fake]
    port = 8171
```

The mock server:
- Runs on host port 8171
- Returns a new price only when the last value has been confirmed in the contract

---

## 16. Using a Specific Docker Image for Chainlink Node

Default behavior builds an image from your current branch:
```toml
[[nodeset.node_specs]]
  [nodeset.node_specs.node]
    docker_ctx = "../../.."
    docker_file = "plugins/chainlink.Dockerfile"
```

To use a prebuilt image:
```toml
[[nodeset.node_specs]]
  [nodeset.node_specs.node]
    image = "image-you-want-to-use"
```

Apply this to every node in your config.

---

## 17. Using Existing EVM & P2P Keys

When using public chains with limited funding, use pre-funded, encrypted keys:

TOML format:
```toml
[[nodesets.node_specs]]
  [nodesets.node_specs.node]
    test_secrets_overrides = """
    [EVM]
    [[EVM.Keys]]
    JSON = '{...}'
    Password = ''
    ID = 1337
    [P2PKey]
    JSON = '{...}'
    Password = ''
    """
```

> Requires `override_mode = "each"` and the same keys across all chains

These limitations come from the current CRE SDK logic and not Chainlink itself.

---

## 18. Troubleshooting

### Chainlink Node Migrations Fail
Remove old Postgres volumes:
```bash
ctf d rm
```
Or remove volumes manually if `ctf` CLI is unavailable.

### Docker Image Not Found
If Docker build succeeds but the image isn’t found at runtime:
- Restart your machine
- Retry the test

---