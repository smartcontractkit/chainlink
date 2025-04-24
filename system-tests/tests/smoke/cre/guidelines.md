# Test Modification and Execution Guide

## Table of Contents
1. [How to Run the Test](#how-to-run-the-test)
   - [Chainlink Node Image](#chainlink-node-image)
   - [Environment Variables](#environment-variables)
   - [Job Distributor Image](#job-distributor-image)
   - [CRE CLI Binary](#cre-cli-binary)
   - [cron Binary](#cron-binary)
   - [PoR Workflow Source Code](#por-workflow-source-code)
   - [Test Timeout](#test-timeout)
   - [Visual Studio Code configuration](#visual-studio-code-configuration)
2. [Docker vs Kubernetes (k8s)](#docker-vs-kubernetes-k8s)
3. [CRIB Requirements](#crib-requirements)
4. [Setting Docker Images for CRIB Execution](#setting-docker-images-for-crib-execution)
5. [Running Tests in Local Kubernetes (`kind`)](#running-tests-in-local-kubernetes-kind)
6. [CRIB Deployment Flow](#crib-deployment-flow)
7. [Switching from kind to AWS provider](#switching-from-kind-to-aws-provider)
8. [CRIB Limitations & Considerations](#crib-limitations--considerations)
9. [Adding a New Capability](#adding-a-new-capability)
   - [Copying the Binary to the Container](#copying-the-binary-to-the-container)
   - [Defining Additional Node Configuration](#defining-additional-node-configuration)
   - [Defining a Job Spec for the New Capability](#defining-a-job-spec-for-the-new-capability)
   - [Registering the Capability in the Capabilities Registry contract](#registering-the-capability-in-the-capabilities-registry-contract)
   - [Updating the DON topology](#updating-the-don-topology-to-assign-the-new-capability-to-one-of-the-dons)
10. [Using a New Workflow](#using-a-new-workflow)
    - [Test Uploads the Binary](#test-uploads-the-binary)
    - [Workflow Configuration](#workflow-configuration)
    - [Workflow Secrets](#workflow-secrets)
    - [Manual Upload of the Binary](#manual-upload-of-the-binary)
    - [YAML workflows](#yaml-workflows)
11. [Deployer Address or Deployment Sequence Changes](#deployer-address-or-deployment-sequence-changes)
12. [Adding a new test to the CI](#adding-a-new-test-to-the-ci)
13. [Multiple DONs](#multiple-dons)
    - [DON Type](#don-type)
    - [Capabilities](#capabilities)
    - [HTTP Port Range Start](#http-port-range-start)
    - [Database (DB) Port](#database-db-port)
    - [Number of nodes](#number-of-nodes)
14. [Price Data Source](#price-data-source)
    - [Live Source](#live-source)
    - [Mocked Data Source](#mocked-data-source)
    - [Blockchain Configuration](#blockchain-configuration)
15. [Using a Specific Docker Image for Chainlink Node](#using-a-specific-docker-image-for-chainlink-node)
16. [Troubleshooting](#troubleshooting)
    - [Chainlink Node migrations fail](#chainlink-node-migrations-fail)
    - [Chainlink image not found in local Docker registry](#chainlink-image-not-found-in-local-docker-registry)
17. [CLI Usage](#cli-usage)
    - [Before you start](#before-you-start)
    - [Environment Variables](#environment-variables)
    - [Cleanup](#cleanup)
18. [Using existing EVM & P2P keys](#using-existing-evm--p2p-keys)

---

## How to Run the Test

Before running the test, several prerequisites must be met.

### Chainlink Node Image

The TOML configuration allows you to choose whether to:

- Use an existing Docker image
- Build a Docker image from the currently checked-out branch

By default, all TOML test configurations build the image from the current branch. This is expressed in the config as follows:

```toml
[nodesets.node_specs.node]
  docker_ctx = "../../../.."
  docker_file = "plugins/chainlink.Dockerfile"
```

If you prefer to use an existing image, update the config to:

```toml
[nodesets.node_specs.node]
  image = "my-docker-image:my-tag"
```

Make this change for each `nodesets.node_specs.node` entry in the config.

**Supported version**: ≥ [e13e5675d3852b04e18dad9881e958066a2bf87a](https://github.com/smartcontractkit/chainlink/commit/e13e5675d3852b04e18dad9881e958066a2bf87a) (merged on 2025-02-25)

---

### Environment Variables

Required environment variables:

- **`CTF_CONFIGS`** – Always required. A comma-separated list of TOML config files to use.
- **`PRIVATE_KEY`** – A plaintext private key used for all contract deployments, configuration, and Chainlink node funding. This key must be sufficiently funded.
- **`GIST_WRITE_TOKEN`** – Required only when compiling and uploading a new workflow. This must be a fine-grained personal access token with `gist:read:write` permissions tied to your personal GitHub account.

---

### Job Distributor Image

The test requires the Job Distributor image to be available locally. By default, `environment-*.toml` files expect an image tagged as `job-distributor:0.9.0`.

To build it locally checkout [smartcontractkit/job-distributor](https://github.com/smartcontractkit/job-distributor) repository and then run:

```bash
git checkout v0.9.0
docker build -t job-distributor:0.9.0 -f e2e/Dockerfile.e2e .
```

If you have access to the production ECR, you can also pull the image from there. Alternatively, update the `environment-*.toml` files with the full name of the image from your registry.

**Supported version**: v0.9.0

---

### CRE CLI Binary

Download the CRE CLI binary compiled for your host machine's architecture from the [smartcontractkit/dev-platform](https://github.com/smartcontractkit/dev-platform) repository. Alternatively, build it on your local machine for operating system and architecture matching yours.

**Supported version**: v0.1.5

---

### `cron` Binary

You must ensure the test environment has access to the `cron` capability binary. You can either:

1. **Use a CL node image that already includes the binary**, or
2. **Make the binary available on your host machine** so that the test can copy it into the running container.

If you choose the first option, comment out the relevant line in your TOML config that specifies the binary path, like this:
```toml
[workflow_config.dependencies]
  # cron_capability_binary_path = "./cron"
```

If you choose the second option, update the relevant line in your TOML config to match the binary path:

```toml
[workflow_config.dependencies]
  cron_capability_binary_path = "./some-folder/cron"
```

To obtain the binary, you can:
- Clone the [smartcontractkit/capabilities](https://github.com/smartcontractkit/capabilities) repo and build it locally
- Download it from the release assets

Ensure the binary is built for **Linux** and **amd64** architecture.

**Supported version**: v0.1.2-alpha

---

### PoR Workflow Source Code

By default, the TOML configs instruct the code to compile the workflow each time the test runs. Therefore, you must clone the [smartcontractkit/proof-of-reserves-workflow-e2e-test](https://github.com/smartcontractkit/proof-of-reserves-workflow-e2e-test) repository.

Then, update the TOML config to point to the workflow's location (relative or absolute paths are supported):

```toml
[workflow_config]
workflow_folder_location = "/Users/my-user/repositories/proof-of-reserves-workflow-e2e-test"
```

### Test timeout
If building Docker image set Go test timeout to 20 minutes. Test should build the image and execute in 12-15 minutes on most machines. When using existing images execution time varies between 4 and 7 minutes.

### Visual Studio Code configuration
Below is a launch configuration that can be used with the VCS:

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

In CI the flow is a bit different, because we generate one-time access tokens to these two repositories and download required assets before running the test.
Also, the test code modifies the config during runtime to production JD image hardcoded in the [.github/e2e-tests.yml](.github/e2e-tests.yml) file as `E2E_JD_VERSION` env var and injects Chainlink node image via `E2E_TEST_CHAINLINK_IMAGE` and `E2E_TEST_CHAINLINK_VERSION` vars.

## Docker vs Kubernetes (k8s)

The following TOML configuration determines whether a test is executed in Docker or Kubernetes:

```toml
[infra]
  # Choose either "docker" or "crib"
  type = "crib"
```

The only way to execute tests in Kubernetes (k8s) is through CRIB, which supports both a local cluster (`kind`) and AWS. When executing in CRIB, you must provide the following configuration:

```toml
[infra.crib]
  namespace = "crib-local"
  folder_location = "$(pwd of crib repository)/deployments/cre"
  # Choose either "aws" or "kind"
  provider = "kind"
```

---

## CRIB Requirements

Before running tests in CRIB, follow these steps:

1. **Read the CRIB Instructions** – Follow the [CRIB deployment guide](https://smartcontract-it.atlassian.net/wiki/spaces/INFRA/pages/660145339/General+CRIB+-+Deploy+Access+Instructions).
2. **Obtain AWS Role** – If you plan to run tests on AWS, acquire the necessary AWS role for CRIB. Running on a local `kind` cluster does not require any roles.
3. **Manually Download Docker Registry Image** – If using the `kind` provider, download the required Docker registry image:
   ```bash
   docker pull registry:2
   ```
4. **Clone the CRIB Repository** – Clone the [CRIB repository](https://github.com/smartcontractkit/crib) and determine its absolute path using `pwd`.
5. **Update the TOML Configuration** – Set the `folder_location` parameter to the absolute path of the `deployments/cre` folder within the CRIB repository.
   ```toml
   [infra.crib]
   folder_location = "/Users/me/repositories/crib/deployments/cre"
   ```
6. **Adjust Namespace and Provider** – If using AWS, you **must** provide cost attribution details:
   ```toml
   [infra.crib.team_input]
   team = "your team"
   product = "name of the product you are working on"
   cost_center = "crib"
   component = "crib"
   ```
7. **Start VPN** - If using AWS.
---

## Setting Docker Images for CRIB Execution

CRIB does **not** support dynamically built Docker images from local `Dockerfile`s during test execution. Using the following TOML configuration will result in an error:

```toml
[nodesets.node_specs.node]
  docker_ctx = "../../../.."
  docker_file = "plugins/chainlink.Dockerfile"
```

Instead, you **must** use the `image` key:

```toml
[nodesets.node_specs.node]
  image = "localhost:5001/chainlink:112b9323-plugins-cron"
```

### Image Restrictions
- All nodes within a single nodeset **must** use the same Docker image (but each nodeset can use a different image).
- The image tag **must** be explicit (omitting tag, so that implicitly `latest` is used is **not** supported).

#### Job Distribution (JD) Image
Currently, CRIB reads only the **image tag** from the TOML configuration. The following setting:

```toml
[jd]
  image = "jd-test-1:my-awesome-tag"
```

Will result in CRIB using an image from the main AWS ECR repository with the tag `my-awesome-tag`.

If an image tag is omitted, an error will occur:

```toml
[jd]
  image = "jd-test-1"  # This will fail
```

---

## Running Tests in Local Kubernetes (`kind`)

### Docker Registry Setup
Ensure you have pulled the `registry:2` Docker image:
```bash
docker pull registry:2
```

### Hostname Routing
All routing to the `kind` cluster is done via `/etc/hosts`. CRIB automatically adds new host entries for detected ingresses, but since `/etc/hosts` is protected, root privileges are required. However, running tests does **not** allow interactive password input, leading to failures when new hostnames must be added.

#### Workarounds
1. **Manually add `/etc/hosts` entries** (tedious but straightforward).
2. **Run `devspace` manually** for each chain/DON before starting tests. This allows CRIB to add entries while you enter the root password interactively.

### Manually Adding `/etc/hosts` Entries
For each component, manually add the following entries:

#### Geth Chain
```bash
127.0.0.1 <NAMESPACE>-geth-<CHAIN_ID>-http.main.stage.cldev.sh
127.0.0.1 <NAMESPACE>-geth-<CHAIN_ID>-ws.main.stage.cldev.sh
```
Example:
```bash
127.0.0.1 crib-local-geth-1337-http.main.stage.cldev.sh
127.0.0.1 crib-local-geth-1337-ws.main.stage.cldev.sh
```

#### Job Distributor
```bash
127.0.0.1 <NAMESPACE>-job-distributor-grpc.main.stage.cldev.sh
```
Example:
```bash
127.0.0.1 crib-local-job-distributor-grpc.main.stage.cldev.sh
```

#### Chainlink Nodes
For bootstrap nodes:
```bash
127.0.0.1 <NAMESPACE>-<DON_TYPE>-bt-<INDEX>.main.stage.cldev.sh
```
For worker nodes:
```bash
127.0.0.1 <NAMESPACE>-<DON_TYPE>-<INDEX>.main.stage.cldev.sh
```

Example (1 bootstrap + 3 worker nodes in `workflow` DON):
```bash
127.0.0.1 crib-local-workflow-bt-0.main.stage.cldev.sh
127.0.0.1 crib-local-workflow-0.main.stage.cldev.sh
127.0.0.1 crib-local-workflow-1.main.stage.cldev.sh
127.0.0.1 crib-local-workflow-2.main.stage.cldev.sh
```

### Automating Hostname Setup with `devspace`
Run the following commands **inside the `cre/deployment` subfolder** and a shell where `nix develop` was executed:

#### Geth Chain
```bash
CHAIN_ID=<CHAIN_ID> devspace run deploy-custom-geth-chain
```
#### Job Distributor
```bash
devspace run deploy-jd
```
#### Chainlink Nodes
```bash
DON_TYPE=<type of don> DON_NODE_COUNT=<number of worker nodes> DON_BOOT_NODE_COUNT=<number of bootstrap nodes> devspace run deploy-don
```

Ensure `DON_TYPE` matches the `name` field in your TOML config:

```toml
[[nodesets]]
  nodes = 5
  name = "workflow"
```

---

## CRIB Deployment Flow

1. **Initialize a `nix develop` shell** and set environment variables.
    - Set environment variables: `PROVIDER`, `DEVSPACE_NAMESPACE`, `CONFIG_OVERRIDES_DIR`
2. **Start Blockchains**:
   - Set `CHAIN_ID` from TOML.
   - Deploy with `devspace run deploy-custom-geth-chain`.
   - Read endpoints from `chain-<CHAIN_ID>-urls.json`.
3. **Deploy Keystone Contracts**.
4. **Generate CL Node Configs & Secrets** (stored in `./crib-configs`).
5. **Start Each DON**:
   - Set environment variables: `DEVSPACE_IMAGE`, `DEVSPACE_IMAGE_TAG`, `DON_BOOT_NODE_COUNT`, `DON_NODE_COUNT` and `DON_TYPE`.
   - Deploy with `devspace run deploy-don`.
   - Read DON URLs from `don-<DON_TYPE>-urls.json`.
   - Copy capabilities binaries to pods with `devspace run copy-to-pods` (if needed).
6. **Start Job Distributor**:
   - Set environment variable: `JOB_DISTRIBUTOR_IMAGE_TAG`.
   - Deploy with `devspace run deploy-jd`.
   - Read JD URLs from `jd-url.json`.
7. **Create Jobs & Configure CRE Contracts** (same as Docker).

---

## Switching from kind to AWS provider
Since `kind` provider uses `/ets/hosts` for routing you **must remove all entries added previously if you are using the same namespace you used in `kind`**. Otherwise traffic will be incorrectly redirected to localhost instead of AWS.

It is thus advised to change namespace names, when switching providers.

---

## CRIB Limitations & Considerations

### Gateway DON
- Must always be on a **dedicated node**.
- Identified using `DON_TYPE=gateway`.
- No bootstrap node required, but multiple worker nodes are allowed.

### Mocked Price Provider
- CRIB does **not** support the mocked data source used in PoR smoke tests, as it runs outside a container.
- Only tests using **live endpoints** can be executed in CRIB.

### Environment variables
- Some are set by the Go code, others are taken from `./deployments/cre/.env` and applied when `nix develop` is run. Make sure that variables set from Go code are not present in `.env` file as it might lead to inconsistent behaviour.

### DNS propagation
- When running in AWS DNS propagation of Ingress domains might be painfully slow. Be patient and try again if failures occur.
- Ingress check on `kind` cluster is also sometimes faulty and fails, even though all systems are operational. Currently, the only remedy is re-running.

### Connection Issues
If you encounter connection problems:
Check pod health:
```bash
kubectl get pods
```

Ensure all pods show "Running" status.
View pod logs:
```bash
kubectl logs <POD_NAME>
```
---

## Adding a New Capability

To add a new capability to the test, follow these steps:

1. Define a new `CapabilityFlag` representing the capability.
2. Add support for the new capability in the testing code:
   - Add code that copies the capability binary to the Chainlink node's Docker container (must be in `linux/amd64` format).
    - You can skip this step if the capability is already included in the Chainlink image you are using or if it's built-in.
   - (Optional) Define additional node configuration if required.
   - Define the job spec for the new capability.
   - Register the capability in the Capabilities Registry contract.
2. Update the DON topology to assign the new capability to one of the DONs.
3. Pass newly defined factory functions to the composable setup method.

Once these steps are complete, you can run a workflow that requires the new capability.

Let's assume we want to add a capability that represents writing to Aptos chain.

#### Defining a CapabilityFlag for the Capability

The testing code uses string flags to map DON capabilities to node configuration, job creation, and the Capabilities Registry contract. This means that adding a new capability requires defining a unique flag. Let's name our capability flag as `WriteAptosCapability`.

First, define the new flag in [flags.go](../../../lib/cre/types/flags.go):

```go
const (
	OCR3Capability          CapabilityFlag = "ocr3"
	CronCapability          CapabilityFlag = "cron"
	CustomComputeCapability CapabilityFlag = "custom-compute"
	WriteEVMCapability      CapabilityFlag = "write-evm"
    WriteAptosCapability    CapabilityFlag = "write-aptos"               // <------------ New entry

	// Add more capabilities as needed
)
```

Now that the flag is defined, let's handle copying of the capability to the container/pod and then move to configuring the nodes and jobs.

### Adding Support for the New Capability in the Testing Code

### Copying the Binary to the Container

Even though it is possible to instruct the code to copy a capability binary to the container/pod using TOML config, it is recommended to handle copying programmatically, because that way binary name and path will be available later on to the job spec-creating code. It also helps to handle different expected capabilities folders in Docker and Kubernetes.

PoR test uses the following code to copy `cron` binary to the container/pod:
```go
customBinariesPaths := map[string]string{}
containerPath, pathErr := capabilities.DefaultContainerDirectory(in.Infra.InfraType)
require.NoError(t, pathErr, "failed to get default container directory")
var cronBinaryPathInTheContainer string
if in.WorkflowConfig.DependenciesConfig.CronCapabilityBinaryPath != "" {
  // where cron binary is located in the container
  cronBinaryPathInTheContainer = filepath.Join(containerPath, filepath.Base(in.WorkflowConfig.DependenciesConfig.CronCapabilityBinaryPath))
  // where cron binary is located on the host
  customBinariesPaths[keystonetypes.CronCapability] = in.WorkflowConfig.DependenciesConfig.CronCapabilityBinaryPath
} else {
  // assume that if cron binary is already in the image it is in the default location and has default name
  cronBinaryPathInTheContainer = filepath.Join(containerPath, "cron")
}

universalSetupInput := creenv.SetupInput{
  // other fields omitted for breviety
  CustomBinariesPaths:                  customBinariesPaths,
  JobSpecFactoryFunctions: []keystonetypes.JobSpecFactoryFn{
    crecron.CronJobSpecFactoryFn(cronBinaryPathInTheContainer),
  },
}

universalSetupOutput, setupErr := creenv.SetupTestEnvironment(testcontext.Get(t), testLogger, cldlogger.NewSingleFileLogger(t), universalSetupInput)
if setupErr != nil {
  panic(err)
}
```

Let's unpack it:
- first, we define `customBinariesPaths` map, where key is the capability flag and value the path to the file on the host machine
- if no such path is set in test config, we assume that binary is already present in the image and that it is named `cron`
- if path is present, we add it to `customBinariesPaths` and construct a `cronBinaryPathInTheContainer` that is later on passed to function that prepares job spec for cron container

Aptos equivalent could look in a following way:
```go
customBinariesPaths := map[string]string{}
containerPath, pathErr := capabilities.DefaultContainerDirectory(in.Infra.InfraType)
require.NoError(t, pathErr, "failed to get default container directory")
var aptosBinaryPathInTheContainer string
if in.WorkflowConfig.DependenciesConfig.AptosCapabilityBinaryPath != "" {
  // where cron binary is located in the container
  aptosBinaryPathInTheContainer = filepath.Join(containerPath, filepath.Base(in.WorkflowConfig.DependenciesConfig.AptosCapabilityBinaryPath))
  // where cron binary is located on the host
  customBinariesPaths[keystonetypes.AptosWriteCapability] = in.WorkflowConfig.DependenciesConfig.AptosCapabilityBinaryPath
} else {
  // assume that if aptos binary is already in the image it is in the default location and has default name
  aptosBinaryPathInTheContainer = filepath.Join(containerPath, "aptos")
}

universalSetupInput := creenv.SetupInput{
  // other fields omitted for breviety
  CustomBinariesPaths:                  customBinariesPaths,
  JobSpecFactoryFunctions: []keystonetypes.JobSpecFactoryFn{
    creaptos.AptosJobSpecFactoryFn(aptosBinaryPathInTheContainer),
  },
}

universalSetupOutput, setupErr := creenv.SetupTestEnvironment(testcontext.Get(t), testLogger, cldlogger.NewSingleFileLogger(t), universalSetupInput)
if setupErr != nil {
  panic(err)
}
```

This works for both Kubernetes and Docker.

> **Note:** Copying the binary to the bootstrap node is unnecessary since it does not handle capability-related tasks. Our code takes care of that automatically.

#### Defining Additional Node Configuration

This step is optional, as not every capability requires additional node configuration. However, writing to the Aptos chain does. Depending on the capability, adjustments might be needed for the bootstrap node, the workflow nodes, or all nodes.

The following code snippet adds the required settings:

```go
if hasFlag(flags, WriteAptosCapability) {
  writeAptosConfig := fmt.Sprintf(`
    # Required for initializing the capability
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
    `,
    "2",                       // chainID
    "http://some-aptos-rpc",   // RPC
    "0x516e771e1b4a903afe74c27d057c65849ecc1383782f6642d7ff21425f4f9c91" // forwarder address
  )
  workerNodeConfig += writeAptosConfig
}
```

This is a placeholder snippet—you should replace it with the actual configuration required for the capability. Ensure it is added before restarting the nodes. Currently, the config-generating code lives in [por.go](../../../lib/cre/don/config/por/por.go), which is very PoR-test centric and might reqiure refactoring in the future to be more universal and composable.

#### Defining a Job Spec for the New Capability

Unlike node configuration, defining a new job spec is almost always required for a new capability. Jobs should only be added to worker nodes.

Assume the Aptos capability job does not require special configuration (this may or may not be true):

```go
if hasFlag(flags, WriteAptosCapability) {
  aptosJobSpec := fmt.Sprintf(`
    type = "standardcapabilities"
    schemaVersion = 1
    externalJobID = "%s"
    name = "aptos-write-capability"
    command = "/home/capabilities/%s"             # <-------- location of the capability binary within the container
    config = ""
  `,
    uuid.NewString(),
    "aptos_linux_amd64")

  aptosJobRequest := &jobv1.ProposeJobRequest{
    NodeId: node.NodeID,
    Spec:   aptosJobSpec,
  }

  _, aptosErr := ctfEnv.Offchain.ProposeJob(context.Background(), aptosJobRequest)
  if aptosErr != nil {
    errCh <- errors.Wrapf(aptosErr, "failed to propose Aptos write job for node %s", node.NodeID)
    return
  }
}
```

This code must be integrated into the section responsible for proposing and approving new jobs using the Job Distributor (JD).

Following the principles of composable architecture it is recommended that this is expressed as function that retruns `types.JobSpecFactoryFn`, which can later be passed to [SetupTestEnvironment](system-tests/lib/cre/environment/environment.go) (the composable setup function) as one of the job spec factory functions in a following way:
```go

import (
  crecompute "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/compute"
	creconsensus "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/consensus"
	crecron "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/cron"
	cregateway "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/gateway"
)

universalSetupInput := creenv.SetupInput{
  // other fields omitted for breviety
  JobSpecFactoryFunctions: []keystonetypes.JobSpecFactoryFn{
    creconsensus.ConsensusJobSpecFactoryFn(chainIDUint64),
    crecron.CronJobSpecFactoryFn(cronBinaryPathInTheContainer),
    cregateway.GatewayJobSpecFactoryFn(chainIDUint64, extraAllowedPorts, []string{}, []string{"0.0.0.0/0"}),
    crecompute.ComputeJobSpecFactoryFn,
    // your new factory function goes here
  },
}

universalSetupOutput, setupErr := creenv.SetupTestEnvironment(testcontext.Get(t), testLogger, cldlogger.NewSingleFileLogger(t), universalSetupInput)
if err != nil {
  panic(err)
}
```

Some examples of how such factories might be implement can be found here:
- [compute.go](system-tests/lib/cre/don/jobs/compute/compute.go)
- [chainreader.go](system-tests/lib/cre/don/jobs/chainreader/chainreader.go)

> **Note:** If the new capability requires a different job type, you may need to update the Chainlink Node code. If it works with `standardcapabilities`, no changes are necessary.

#### Registering the Capability in the Capabilities Registry Contract

The final step is adding support for registration of the capability with the Capabilities Registry contract:

```go
if hasFlag(donTopology.Flags, WriteAptosCapability) {
  capabilities = append(capabilities, keystone_changeset.DONCapabilityWithConfig{
    Capability: kcr.CapabilitiesRegistryCapability{
      LabelledName:   "write_aptos-testnet",          // <------- Ensure correct name (it might be dynamic and depend on things like chainID)
      Version:        "1.0.0",                        // <------- Ensure correct version
      CapabilityType: 3, // TARGET
      ResponseType:   1, // OBSERVATION_IDENTICAL
    },
    Config: &capabilitiespb.CapabilityConfig{},
  })
}
```

Ensure that the **name and version** match:
- The values used by the capability itself.
- The values used in the workflow definition.

If they do not match, the test will likely fail in a way that is difficult to diagnose.

Some capabilities may also require a `ConfigurationContract`. Check with the capability author for the necessary values and ensure the correct capability type is set.

Just like in case of job specs it is recommended that you wrap that code in a factory function with the `func(donFlags []string) []keystone_changeset.DONCapabilityWithConfig` type and then pass it to the universal setup function in the following way:
```go

import (
  libc "github.com/smartcontractkit/chainlink/system-tests/lib/conversions"
	libcontracts "github.com/smartcontractkit/chainlink/system-tests/lib/cre/contracts"
)

universalSetupInput := creenv.SetupInput{
  // other fields omitted for breviety
  CapabilitiesContractFactoryFunctions: []keystonetypes.DONCapabilityWithConfigFactoryFn{
    libcontracts.DefaultCapabilityFactoryFn,
    libcontracts.ChainWriterCapabilityFactory(libc.MustSafeUint64(int64(chainIDInt))),
    // your new factory function goes here
    }
  }

universalSetupOutput, setupErr := creenv.SetupTestEnvironment(testcontext.Get(t), testLogger, cldlogger.NewSingleFileLogger(t), universalSetupInput)
if err != nil {
  panic(err)
}
```

### Updating the DON topology to assign the new capability to one of the DONs.

As mentioned earlier, infrastructure (number of DONs, number of nodes in each DON, images to use) is defined in TOML config, but capabilities and DON types are defined in Go code. Let's assume we want have 1 workflow DON and 1 gateway DON. Workflow DON will expose a couple of capabilities, including our newly added one.
```go
dons := []*cretypes.CapabilitiesAwareNodeSet{
			{
				Input:              in.NodeSets[0],                   // this comes from TOML
				Capabilities:       []string{cretypes.OCR3Capability, cretypes.CustomComputeCapability, cretypes.WriteAptosCapability}, // <----- added here
				DONTypes:           []string{cretypes.WorkflowDON},
				BootstrapNodeIndex: 0,
			},
			{
				Input:              in.NodeSets[2],
				Capabilities:       []string{},
				DONTypes:           []string{cretypes.GatewayDON}, // <----- it's crucial to set the correct DON type
				BootstrapNodeIndex: -1,                            // <----- it's crucial to indicate there's no bootstrap node
				GatewayNodeIndex:   0,
			},
		}
```
> Note: Remember to define 2 nodesets with appropriate number of nodes in the TOML config!

Now that we have the desired topology, we can pass it to the composable setup function:
```go
universalSetupInput := creenv.SetupInput{
  CapabilitiesAwareNodeSets: dons,
  // other fields omitted for breviety
}

universalSetupOutput, setupErr := creenv.SetupTestEnvironment(testcontext.Get(t), testLogger, cldlogger.NewSingleFileLogger(t), universalSetupInput)
if setupErr != nil {
  panic(setupErr)
}
  ```

## Using a New Workflow

To test a new workflow, you have two options:

1. Compile the workflow to a WASM binary and upload it to Gist **inside the test**.
2. Manually upload the binary and specify the workflow URL in the test configuration.

### Test Uploads the Binary

For the test to compile and upload the binary, modify your TOML configuration:

```toml
[workflow_config]
  use_cre_cli = true                  # must be 'true', can compile new workflow only the CRE CLI
  should_compile_new_workflow = true  # must be 'true'
  workflow_folder_location = "path-to-folder-with-main.go-of-your-workflow"
```

### Workflow Configuration

If your workflow requires configuration, modify the test to create and pass the configuration data to CRE CLI. For example:

```go
configFile, err := os.CreateTemp("", "config.json")
require.NoError(t, err, "failed to create workflow config file")

workflowConfig := PoRWorkflowConfig{
  FeedID:          feedID,
  URL:             "https://api.real-time-reserves.verinumus.io/v1/chainlink/proof-of-reserves/TrueUSD",
  ConsumerAddress: feedsConsumerAddress.Hex(),
}
```

> **Note:** If the workflow is **not configurable**, do not pass configuration data. Instead, pass an empty `[]byte` when compiling or registering it.

> **Note:** Currently, we do not allow to update the configuration alone. Each configuration change is treated as workflow change and thus requires following the **upload mode**.

---

### Workflow Secrets

Broadly speaking using secrets in your workflow requires:
1. creation of `secrets.config.yaml` file, which contains a mapping of secrets to environment variables that hold them
2. encrypting of these secrets using `CRE CLI` and uploading them
3. referencing them in your workflow

> Note: Secrets are encrypted using using CL nodes CSA public keys and so it is impossible to share secrets between different DONs unless they are using the same CSA and P2P keys.

There is a helper function in [secrets.go](system-tests/lib/crecli/secrets.go) that helps in creating `secrets.config.yaml`. Once you have that file ready, you will need to instruct the code to set environment variables secrets reference and then encrypt them.

Depending on whether environment variable should be the same for all nodes, or whether each node should get a different secret value, environment variables get different prefixes. Here we will concern ourselves only with the first case. And that means adding `_ENV_VAR_ALL` suffix to the environment variable containing the secret.

Let's assume that you have added two secrets to that `secrets.config.yaml` file:
- `SECRET_A`
- `SECRET_B`
And that both should be the same for all nodes and that first one should map to env var called `FIRST_SECRET_ENV_VAR_ALL` and second `SECOND_SECRET_ENV_VAR_ALL`.

When calling our wrapper function that compiles and registers a new workflows has secret support you'd craft the following struct representing new workflow:
```go
newWorkflow := keystonetypes.NewWorkflow{
  // other fields skipped for breviety
  SecretsFilePath:  secretsFilePath,
  Secrets: map[string]string{
    "FIRST_SECRET_ENV_VAR_ALL": "first secret value",
    "SECOND_SECRET_ENV_VAR_ALL": "second secret value",
  },
}

input := keystonetypes.RegisterWorkflowWithCRECLIInput{
  // other fields skipped for breviety
  NewWorkflow: newWorkflow,
  ShouldCompileNewWorkflow: true,
}

registerErr := creworkflow.RegisterWithCRECLI(registerWorkflowInput)
if registerErr != nil {
  return errors.Wrap(registerErr, "failed to register workflow with CRE CLI")
}
```

This way you will have secrets encoded with CSA keys of the workflow DON and uploaded to Gist.

---

### Manual Upload of the Binary

If you compiled and uploaded the binary yourself, set the following in your configuration:

```toml
[workflow_config]
  use_cre_cli = true                  # keep set to 'true', we will still use it to register the workflow
  should_compile_new_workflow = false # set to 'false'

  [workflow_config.compiled_config]
    binary_url = "<binary-url>"
    config_url = "<config-url>"       # optional
    secrets_url = "<secrets-url>"     # optional
```

All URLs must be accessible by the gateway node(s).

### YAML workflows

When using workflows expressed in YAML you do not need to compile, upload and register them anywhere, manually or using `CRE CLI`. So all the above-mentioned parts of the test can be left out. All you need to do is:
1. Define the workflow as `string`
2. Use JD to propose is it in the same way that jobs are proposed

#### Workfow definition
Here's an example of a workflow:
```toml
type = "workflow"
schemaVersion = 1
name = "my-df-workflow"
externalJobID = "my-df-workflow-f712hdf"
workflow = """
name: "my-df-workflow"
owner: '0xe7f1725E7734CE288F8367e1Bb143E90bb3F0512'
triggers:
 - id: streams-trigger@1.0.0
   config:
     maxFrequencyMs: 5000
     feedIds:
       - '0x018bfe88407000400000000000000000'
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
     aggregation_config:
       allowedPartialStaleness: '0.5'
       feeds:
        '0x018bfe88407000400000000000000000':
          deviation: '0.01'
          heartbeat: 600
     encoder: EVM
     encoder_config:
       abi: (bytes32 FeedID, uint224 Price, uint32 Timestamp)[] Reports
targets:
  - id: write_geth-testnet@1.0.0
    inputs:
      signed_report: $(ccip_feeds.outputs)
    config:
      address: '0x24309990d635A6C5FF711503BfCb942dd25F96A0'
      deltaStage: 10s
      schedule: oneAtATime
```

It will be triggered by streams update, take the trigger data as input for consensus and once it's reached upload the report on EVM chain.

Things to keep in mind:
- job type must be `workflow`
- actual workflow code is passed in the `workflow` field

#### Proposing the workflow using JD

You start by creating a job proposal:
```go
// assume it contains above-mentioned workflow spec
var workflowSpec string

job := &jobv1.ProposeJobRequest{
		NodeId: nodeID, // nodeID with which it registered itself in JD
    Spec: workflowSpec}
```

Then you either call JD directly:
```go
timeout := time.Second * 60
ctx, cancel := context.WithTimeout(context.Background(), timeout)
defer cancel()
_, err := offChainClient.ProposeJob(ctx, jobReq)
if err != nil {
  return errors.Wrapf(err, "failed to propose job %s for node %s", jobDesc.Flag, jobReq.NodeId)
}
```

Or if you are using our wrappers:
```go
import (
  libnode "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/node"
  keystonetypes "github.com/smartcontractkit/chainlink/system-tests/lib/cre/types"
)

// assuming you have devenv.Environment and keystonetypes.DonTopology references already
var env *devenv.Environment
var donTopology *keystonetypes.DonTopology
donToJobSpecs := make(keystonetypes.DonsToJobSpecs)

createJobsInput := keystonetypes.CreateJobsInput{
  CldEnv:        env,
  DonTopology:   donTopology,
  DonToJobSpecs: donToJobSpecs,
}

// assuming there's only 1 DON in topology
// we want to create that workflow for all worker nodes, bootstrap doesn't need it
workflowNodeSet, err := node.FindManyWithLabel(donTopology[0].NodesMetadata, &keystonetypes.Label{Key: libnode.NodeTypeKey, Value: keystonetypes.WorkerNode}, libnode.EqualLabels)
if err != nil {
  // there should be no DON without worker nodes, even gateway DON is composed of a single worker node
  return nil, errors.Wrap(err, "failed to find worker nodes")
}

donToJobSpecs[donTopology.WorkflowDonID] = make(keystonetypes.DonJobs)

jobDesc := keystonetypes.JobDescription{Flag: keystonetypes.OCR3Capability, NodeType: keystonetypes.WorkerNode}

for idx, workerNode := range workflowNodeSet {
    nodeID, nodeIDErr := workerNode.FindLabelValue(workerNode, libnode.NodeIDKey)
		if nodeIDErr != nil {
			return nil, errors.Wrap(nodeIDErr, "failed to get node id from labels")
		}

  donToJobSpecs[donTopology.WorkflowDonID][jobDesc] = append(donToJobSpecs[donTopology.WorkflowDonID][jobDesc], &jobv1.ProposeJobRequest{
		NodeId: nodeID, // nodeID with which it registered itself in JD
    Spec: workflowSpec} // spec we previously created
  )
}

createJobsErr := libdon.CreateJobs(testLogger, createJobsInput)
if createJobsErr != nil {
  panic(createJobsErr)
}
```

---

## Deployer Address or Deployment Sequence Changes

By default, tests running in the CI reuse existing workflow binaries and configurations. This is important, because contract addresses, including Data Feeds Cache, which acts as the consumer in many tests, remains the same **as long as the deployer address (`f39fd6e51aad88f6f4ce6ab8827279cfffb92266`) and contract deployment sequence do not change**.

If the deployer private key or deployment sequence changes, run the test in **upload mode**:

```toml
[workflow_config]
  use_cre_cli = true
  should_compile_new_workflow = true
  workflow_folder_location = "path-to-folder-with-main.go-of-your-workflow"
```

And then update the config URL in the TOML configs for CI.

---

## Adding a new test to the CI

Due to limitations of CI execution (lack of Gist write token), there are some differences between local and CI execution. They will be removed once `CRE CLI` supports new storage targets.

There are two steps to follow, when adding a new test to the CI:
1. CI-specific config creation
2. Adding tests to [.github/e2e-tests.yaml](../../../../.github/e2e-tests.yml)

### CI-specific configuration
At the very least the CI TOML configuration should contain the following:
```toml
[workflow_config]
  workflow_name = "abcdefgasd"
  feed_id = "018e16c39e000320000000000000000000000000000000000000000000000000"

  use_cre_cli = true
  should_compile_new_workflow = false

  [workflow_config.compiled_config]
    binary_url = "https://gist.githubusercontent.com/Tofel/1a2f6f7d9424bcb176a7cee31112bc11/raw/0e3859564f92c141cb856c65ef7a3db711588c08/binary.wasm.br.b64"
    # if fake is enabled AND we do not compile a new workflow, this config needs to use URL pointing to IP, on which Docker host is available in Linux systems
    # since that's the OS of our CI runners.
    config_url = "https://gist.githubusercontent.com/Tofel/3905a87f22f105da5c0d7196ce7032c4/raw/63c982400f682a95580bf7b5b422aaf8ef4ba511/two_dons_config.json_03_04_2025"
```

The important parts are:
- `should_compile_new_workflow = false`
- `binary_url`
- `config_url` (if used)
- `secrets_url` (if used)

You have to replace the URLs with ones that are created during local test execution, since we cannot use the `CRE CLI` in CI to upload anything to Gist.

### Adding tests to e2e-tests.yaml

This file contains all known e2e tests and is used to run tests matching defined triggers. Here's an example of such entry, which you need to adjust to match your test:
```yaml
  - id: system-tests/smoke/cre/por_test.go:<YOUR TEST NAME>
    path: system-tests/tests/smoke/cre/por_test.go
    test_env_type: docker
    runs_on: ubuntu-latest
    triggers:
      - PR CRE E2E Core Tests
      - Merge Queue CRE E2E Core Tests
      - Nightly E2E Tests
      - Push CRE E2E Core Tests
    test_cmd: cd tests && pushd smoke/cre/cmd > /dev/null && go run main.go download all --output-dir ../ --gh-token-env-var-name GITHUB_API_TOKEN --cre-cli-version v0.1.5 --capabilities-name cron --capabilities-version v1.0.2-alpha 1>&2 && popd > /dev/null && { go test github.com/smartcontractkit/chainlink/system-tests/tests/smoke/cre -v -run "^(YOUR TEST FUNCION NAME)$" -timeout 30m -count=1 -test.parallel=1 -json; exit_code=$?; ../../tools/ci/wait-for-containers-to-stop.sh 30; exit $exit_code; } # Sleep to allow testcontainers to stop
    pyroscope_env: ci-smoke-capabilities-evm-simulated
    test_env_vars:
      E2E_TEST_CHAINLINK_VERSION: "{{ env.DEFAULT_CHAINLINK_PLUGINS_VERSION }}" # This is the chainlink version that has the plugins
      E2E_JD_VERSION: 0.9.0 # there is no latest tag for this repo, so we need to specify the version
      GITHUB_READ_TOKEN: "{{ env.GITHUB_API_TOKEN }}" # GATI-provided token that can read from capabilities and dev-platform repos
      CI: "true"
      CTF_CONFIGS: "<YOUR CONFIG FILE NAME>"
      # Anvil developer key, not a secret
      PRIVATE_KEY: "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
    test_go_project_path: system-tests
```

At the very least you should adjust:
- `id`
- `path`
- `test_cmd`

Any extra environment variables that test needs to access should be added under `test_env_vars`. If these environment variables are dynamic, reach out to DevEx team for further guidance.

Now, lets dive deeper.

#### id

That's just a test identifier, it can be any string, but we prefer to use a following pattern: `<relative-path-to-test-file>:<test function name>`. E.g. `system-tests/smoke/cre/por_test.go:TestCRE_OCR3_PoR_Workflow_CapabilitiesDons_LivePrice`

#### path

Relative path to the file containg test function.

#### test_cmd

Here things get tricky, because current command does e things:
- downloads `CRE CLI` and capabilities binaries
- runs the test
- waits for all Docker containers to exit (feature required by Flakeguard)

At the very least be sure to:
- adjust test folder
- adjust test function name
- pass names of any extra capabilities binaries that need to be available to your test (make sure the version is correct)

Since the command executes in `system-tests/tests/smoke/cre/cmd` folder, once downloaded `CRE CLI` and binaries will be copied to parent folder (as indicated by `--output-dir ../`) and placed in `system-tests/tests/smoke/cre`, which allows them to be used by this TOML config:
```toml
  [workflow_config.dependencies]
  cron_capability_binary_path = "./cron"
  cre_cli_binary_path = "./cre_v0.1.5_darwin_arm64"
```

## Multiple DONs

You can choose to use one or multiple DONs. Configuring multiple DONs requires only TOML modifications, assuming they use capabilities already supported in the testing code.

Currently, the supported capabilities are:
- `cron`
- `ocr3`
- `custom-compute`
- `write-evm`
- `read contract`
- `log-event-trigger` (under development)
- `web-api-trigger` (under development)
- `web-api-target` (under development)

To enable multi-DON support, update the configuration file by:
- Defining a new nodeset.
- Copying the required capabilities to the containers (if they are not built into the image already).

### Key Considerations
When configuring multiple DONs, keep the following in mind:
- **DON name**
- **HTTP Port Range Start**
- **Database (DB) Port**
- **Number of nodes**

### Capabilities

TOML config is primarly concerned with the infrastructure: Docker vs k8s, images to use, etc. But it knows nothing about capabilities of each node set or DON types. These need to be defined in the Go code.

Three types of DONs are supported:
- `workflow`
- `capabilities`
- `gateway`

There should only be **one** `workflow` and `gateway` DON, but multiple `capabilities` DONs can be defined.

For example:
```go
[]*cretypes.CapabilitiesAwareNodeSet{
			{
				Input:              in.NodeSets[0],                   // this comes from TOML
				Capabilities:       []string{cretypes.OCR3Capability, cretypes.CustomComputeCapability, cretypes.WebAPITriggerCapability},
				DONTypes:           []string{cretypes.WorkflowDON},
				BootstrapNodeIndex: 0,
			},
			{
				Input:              in.NodeSets[1],
				Capabilities:       []string{cretypes.WriteEVMCapability, cretypes.WebAPITargetCapability},
				DONTypes:           []string{cretypes.CapabilitiesDON}, // <----- it's crucial to set the correct DON type
				BootstrapNodeIndex: 0,
			},
			{
				Input:              in.NodeSets[2],
				Capabilities:       []string{},
				DONTypes:           []string{cretypes.GatewayDON}, // <----- it's crucial to set the correct DON type
				BootstrapNodeIndex: -1,                            // <----- it's crucial to indicate there's no bootstrap node
				GatewayNodeIndex:   0,
			},
		}
```

Currently, the framework does not enforce validation on whether capabilities are assigned to the correct DON types. However, some capabilities **must** run on the `workflow` DON. These include:
* `ocr3`
* `cron`
* `custom-compute`
and possibly some other ones.

The following capabilities are supported:
- `ocr3`
- `cron`
- `custom-compute`
- `write-evm`
- `read contract` (no test uses it)
- `log-event-trigger` (no test uses it)
- `web-api-trigger` (no test uses it)
- `web-api-target` (no test uses it)

### HTTP Port Range Start

Each node exposes a port to the host. To prevent port conflicts, assign a distinct range to each nodeset. A good practice is to separate port ranges by **50 or 100** between nodesets.

### Database (DB) Port

Similar to HTTP ports, ensure each nodeset has a unique database port.

For a working example of a multi-DON setup, refer to the [`environment-capabilities-don.toml`](environment-capabilities-don.toml) file.

### Number of nodes

In the simples case, when defining a new nodeset, in which every node has the same configuration you can take advantage of the powerful `override_mode = all`, which makes adding a new nodeset as trivial as:
```toml
[[nodesets]]
  nodes = 5                   # set number of nodes you need
  override_mode = "all"       # this is crucial!
  name = "my-new-nodeset"

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

In case each node requires a slightly different config (e.g. you want the bootstrap node to expost some port, or you want to import existing keys) you need to use `override_mode = "each"` and **provide nodespecs for each node**. For example:
```toml
[[nodesets]]
  nodes = 2
  override_mode = "each"
  name = "my-small-nodeset"

  [[nodesets.node_specs]]
    [nodesets.node_specs.node]
      image = "localhost:5001/chainlink:112b9323-plugins-cron"
      custom_ports = ["5002:5002"]          # custom port exposed only on first node
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

  [[nodesets.node_specs]]
    [nodesets.node_specs.node]
      image = "localhost:5001/chainlink:112b9323-plugins-cron"
      # no port exposed here
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

---

## Price Data Source

The PoR test supports both **live** and **mocked** data sources, configurable via TOML.

### Live Source

The PoR workflow is designed to work with the following API:
[http://api.real-time-reserves.verinumus.io](http://api.real-time-reserves.verinumus.io)

Only this response structure is supported. If you want to use a different data source, you must modify both the workflow code and its configuration.

To configure a live data source, use the following TOML settings:

```toml
[price_provider]
  # Without the 0x prefix!
  feed_id = "018bfe8840700040000000000000000000000000000000000000000000000000"
  url = "api.real-time-reserves.verinumus.io/v1/chainlink/proof-of-reserves/TrueUSD"
```

### Blockchain Configuration

Tests are working with `Anvil` by default. All the configurations are using `0s` blocks when deploying contracts but `DeployDataFeedsCache` requires blocks to be mined so `custom_anvil_miner` controls the speed of blocks after the initial deployment is complete.

This allows us to test changes faster.
```
[blockchain_a]
  type = "anvil"
  chain_id = "1337"

[custom_anvil_miner]
  block_speed_seconds = 5
```

If you need to switch to a slow chain you can do it like this, `-b` controls block production speed.
```
[blockchain_a]
  chain_id = "1337"
  docker_cmd_params = ["-b", "5"]
  type = "anvil"
```

### Mocked Data Source

A mocked data source has been introduced to:
- Avoid dependency on a third-party endpoint.
- Enable verification of price values returned by the mock against those stored in the consumer contract.

To configure a mocked data source, use the following TOML settings:

```toml
[price_provider]
  # Without the 0x prefix!
  feed_id = "018bfe8840700040000000000000000000000000000000000000000000000000"

  [price_provider.fake]
    port = 8171
```

This configuration launches a mock server on **port 8171** on the host machine. Prices that mock will returned are set in the Go code of the fake data provider. A new price is returned **only after the previous one has been observed in the consumer contract**. The test completes once all prices have been matched.

---

## Using a Specific Docker Image for Chainlink Node

By default, the test builds a Docker image from the current branch:

```toml
[[nodeset.node_specs]]
  [nodeset.node_specs.node]
  docker_ctx = "../../.."
  docker_file = "plugins/chainlink.Dockerfile"
```

To use an existing image, change it to:

```toml
[[nodeset.node_specs]]
  [nodeset.node_specs.node]
  image = "image-you-want-to-use"
```

Apply this change to **all node entries** in the test configuration.

## Troubleshooting

### Chainlink Node migrations fail

If you see Chainlink Node migrations fail it might, because the Postgres volume has some old data on it. Do remove it and run the test again.
If you have the `ctf` CLI you can use following command: `ctf d rm`. If you do not have the `ctf` binary, remove all volumes associated with CL nodes manually.

### Chainlink image not found in local Docker registry

If you are building the Chainlink image using the Dockerfile, image is successfuly built and yet nodes do not start, because image cannot be found in the local machine, simply restart your computer and try again.

## CLI Usage

The CRE CLI provides commands to manage the local environment. The main commands are:

### Start Environment

```bash
# assuming you are in system-tests/smoke/cre/
cd cmd && go run main.go env start [flags]
```

Flags:
- `-t, --topology string` - Topology to use for the environment (simplified or full) (default "simplified")
- `-w, --wait-on-error-timeout string` - Wait on error timeout duration (e.g. 10s, 1m, 1h)
- `-e, --extra-allowed-ports intSlice` - Extra allowed ports (e.g. 8080,8081)

Example:
```bash
# assuming you are in system-tests/smoke/cre/
cd cmd && go run main.go env start -t simplified -w 5m -e 8080,8081
```

Simplified topology will lanuch a single DON with all capabilities. Full topology will start 3 DONs:
- workflow DON (5 nodes)
- capabilities DON (2 nodes)
- gateway DON (1 node)

Wait on error timeout flag is useful if you want to wait until containers are removed during a failed startup. For example if containers failed to start it allows you to inspect the failure reason.

Extra allowed ports are useful if your gateway needs to access servces running on ports different than `80` and `443`.

Once the environment has started `cre.settings.yaml` file will be created in current directory, with correct contract addresses and RPC URL, so that you can use it together with `CRE CLI` to compile, upload and register workflows.

### Stop Environment

```bash
# assuming you are in system-tests/smoke/cre/
cd cmd && go run main.go env stop
```

This command stops the local CRE environment. If the environment is not running, it will simply fall through.

### Before you start

1. Decide, which topology you need (as you might need to modify corresponding TOML file):
  - `configs/single-don.toml` for simplified topology
  - `configs/workflow-capabilities-don.toml` for full topology
2. Decide, whether your workflows require any capabilities that are not bundled in, such as:
- `cron`
- `read contract`
- `log-event-trigger`

And if it does download them to your machine (**linux/amd64** versions) from [smartcontractkit/capabilities](https://github.com/smartcontractkit/capabilities) repository and point the TOML config to their location:
```toml
[extra_capabilities]
# uncomment as needed and adjust paths to enable these capabilities and have them copied to containers/pods and configured
# cron_capability_binary_path = "../cron"
# log_event_trigger_binary_path = "../logtrigger"
# read_contract_capability_binary_path = "../readcontract"
```

3. Decide, whether you want to build a Docker image based on currently checked out branch or use existing one and modify TOML config [accordingly]((#using-a-specific-docker-image-for-chainlink-node))
> Note: If using existing image, be sure to use `-plugins` version, which contains OCR3 capability

4. Pull or build Job Distributor image as described [here](#job-distributor-image).
5. Optionally download [CTF binary](https://smartcontractkit.github.io/chainlink-testing-framework/framework/getting_started.html) and start observability stack with `ctf obs up`

### Environment Variables

The CLI uses the following environment variables:

- `CTF_CONFIGS` - Path to the TOML configuration file. If not set, defaults to:
  - `configs/single-don.toml` for simplified topology
  - `configs/workflow-capabilities-don.toml` for full topology
- `PRIVATE_KEY` - Private key used for contract deployments and node funding. If not set, defaults to a test key.
- `TESTCONTAINERS_RYUK_DISABLED` - Set to "true" to disable Ryuk container cleanup

### Cleanup

If the environment encounters an unexpected error during startup, you may need to manually clean up resources. Use the following command:

```bash
ctf d rm
```

This will remove all containers with the 'ctf' label and their associated volumes.

# Using existing EVM & P2P keys

It is a good practice, when nodes are connected to public chains on which we have limited access to funding. If nodes use exiting EVM keys we can fund them once and restart/redeploy nodes without losing access to these funds. In both cases we support only encrypted JSON keys. They need to be added to TOML config in a following manner:
```toml
  [[nodesets.node_specs]]
    [nodesets.node_specs.node]
      # other fields go here...
      test_secrets_overrides = """
      [EVM]
      [[EVM.Keys]]
      JSON = '{"address":"4e132a27812dfc644e2c23bfdbc961d9fde6dfca","crypto":{"cipher":"aes-128-ctr","ciphertext":"7ea055f7ef9f643354da1aed91ffc72675b602ed6e4842078ff1ae0a9c2de9e5","cipherparams":{"iv":"83b4e57106f933fc23f70df24d1ccc08"},"kdf":"scrypt","kdfparams":{"dklen":32,"n":262144,"p":1,"r":8,"salt":"dfeab65812944899a7fd9973023d8b5db9985b0b358c08f5ab2445c59bbb7457"},"mac":"59fc362b8c7f5e48369af5e63d8ea8351db3f80a3558fc906ef162b4172fa153"},"id":"00000000-0000-0000-0000-000000000000","version":3}'
      Password = ''
      ID = 1337
      [P2PKey]
      JSON = '{"keyType":"P2P","publicKey":"c317793efae82840441811dd47ea3ddc4ebbc6d315b60021dcc6b3b6fa78eb5a","peerID":"p2p_12D3KooWNwvRtyW3MJQEBK5YfvH8NicvYFMdkByDNaU6DAbkqmsf","crypto":{"cipher":"aes-128-ctr","ciphertext":"4f1fa051e495d9ef2a691989df42f8188a181583a7a58f9daec659686f0c4afba60bccda959bcffa9e538e4e1b0bf3aa3bc2cf3a79e2f9d31a4a9163057f26a50447db97","cipherparams":{"iv":"913448f73b824da66400784a00e64d7a"},"kdf":"scrypt","kdfparams":{"dklen":32,"n":262144,"p":1,"r":8,"salt":"ad42b3dbac2b1da586dec40608eb0341eaaceb5af32f5e93c3012b4db4da6602"},"mac":"b9727ac0aeca91fd9cadd950a0cb4eff82b9ddb4afc6f8346af02fd767f5419b"}}'
      Password = ''
      """
  ```

  That functionality is available only for `nodesets` using `override_mode` equal to `each`, because we need to supply keys for all nodes from the nodeset. Furthermore, if we have more than one DON, then either no DON can use existing keys or all DONs must. This is a limitation of current implementation, which we might remove in the future, if there's such a need. There are some more limitations to bear in mind:
  - you need to import **both** P2P keys and EVM keys
  - when importing EVM keys for multiple chains, the same keys must be used for each chain

  These limitations are related to our CRE SDK and some shortcuts we took in the past (we don't need very complex import support for tests), and not of the CL node itself.