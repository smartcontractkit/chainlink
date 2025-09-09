# Test Modification and Execution Guide

## Table of Contents

- [Test Modification and Execution Guide](#test-modification-and-execution-guide)
  - [Table of Contents](#table-of-contents)
  - [Test Modification and Execution Guide](#test-modification-and-execution-guide-1)
  - [1. Running the Test](#1-running-the-test)
    - [Chainlink Node Image](#chainlink-node-image)
    - [Environment Variables](#environment-variables)
    - [`cron` Capability Binary](#cron-capability-binary)
    - [Test Timeout](#test-timeout)
    - [Visual Studio Code Debug Configuration](#visual-studio-code-debug-configuration)
  - [2. Using the CLI](#2-using-the-cli)
  - [3. Docker vs Kubernetes (k8s)](#3-docker-vs-kubernetes-k8s)
  - [4. CRIB Requirements](#4-crib-requirements)
  - [5. Setting Docker Images for CRIB Execution](#5-setting-docker-images-for-crib-execution)
    - [Notes](#notes)
    - [Job Distributor (JD) Image in CRIB](#job-distributor-jd-image-in-crib)
      - [AWS Provider](#aws-provider)
      - [Kind Provider](#kind-provider)
  - [6. Running Tests in Local Kubernetes (`kind`)](#6-running-tests-in-local-kubernetes-kind)
- [Caution: this section is outdated. We no longer use `devspace`](#caution-this-section-is-outdated-we-no-longer-use-devspace)
    - [Docker Registry Setup](#docker-registry-setup)
    - [Hostname Routing with `/etc/hosts`](#hostname-routing-with-etchosts)
      - [Solutions](#solutions)
      - [Example Manual Entries](#example-manual-entries)
    - [Automating Host Setup with `devspace`](#automating-host-setup-with-devspace)
  - [7. CRIB Deployment Flow](#7-crib-deployment-flow)
- [Caution: this section is outdated. We no longer use `devspace`](#caution-this-section-is-outdated-we-no-longer-use-devspace-1)
  - [8. Switching from `kind` to AWS Provider](#8-switching-from-kind-to-aws-provider)
  - [9. CRIB Limitations \& Considerations](#9-crib-limitations--considerations)
- [Caution: this section is partially outdated. We no longer use `devspace`](#caution-this-section-is-partially-outdated-we-no-longer-use-devspace)
    - [Mocked Price Provider](#mocked-price-provider)
    - [Environment Variables](#environment-variables-1)
    - [DNS Propagation (AWS only)](#dns-propagation-aws-only)
    - [Ingress Check (local `kind` only)](#ingress-check-local-kind-only)
    - [Connectivity Troubleshooting](#connectivity-troubleshooting)
  - [11. Using a New Workflow](#11-using-a-new-workflow)
    - [Workflow Compilation Process](#workflow-compilation-process)
      - [Compilation Example](#compilation-example)
      - [Compilation Requirements](#compilation-requirements)
    - [Workflow Configuration](#workflow-configuration)
      - [Creating Configuration Files (Optional)](#creating-configuration-files-optional)
    - [File Distribution to Containers](#file-distribution-to-containers)
      - [Copying Files to Containers](#copying-files-to-containers)
      - [Container Discovery](#container-discovery)
    - [Workflow Registration](#workflow-registration)
      - [Registration Process](#registration-process)
      - [Registration Parameters](#registration-parameters)
      - [URL Resolution Process](#url-resolution-process)
    - [Complete Workflow Setup Example](#complete-workflow-setup-example)
    - [12. Workflow Secrets](#12-workflow-secrets)
      - [How Secrets Work](#how-secrets-work)
      - [Creating Secrets Configuration](#creating-secrets-configuration)
      - [Environment Variable Naming](#environment-variable-naming)
      - [Using Secrets in Workflows](#using-secrets-in-workflows)
      - [Secrets Encryption Process](#secrets-encryption-process)
      - [Encrypted Secrets File Structure](#encrypted-secrets-file-structure)
      - [Security Considerations](#security-considerations)
      - [Complete Example](#complete-example)
    - [13. YAML Workflows (Data Feeds DSL)](#13-yaml-workflows-data-feeds-dsl)
  - [14. Adding a New Test to the CI](#14-adding-a-new-test-to-the-ci)
    - [How Auto-Discovery Works](#how-auto-discovery-works)
    - [Test Architecture Pattern](#test-architecture-pattern)
      - [Supported DON Topologies](#supported-don-topologies)
    - [What You Need to Do](#what-you-need-to-do)
      - [1. Create Your Test Function](#1-create-your-test-function)
      - [2. Follow Test Naming Convention](#2-follow-test-naming-convention)
      - [3. Use Standard Test Structure](#3-use-standard-test-structure)
    - [CI Configuration Details](#ci-configuration-details)
    - [Environment Setup](#environment-setup)
    - [Test Execution](#test-execution)
    - [Important Notes](#important-notes)
    - [Troubleshooting](#troubleshooting)
  - [15. Price Data Source](#15-price-data-source)
    - [PriceProvider Interface](#priceprovider-interface)
    - [Live Price Source (TrueUSDPriceProvider)](#live-price-source-trueusdpriceprovider)
    - [Mocked Price Source (FakePriceProvider)](#mocked-price-source-fakepriceprovider)
    - [Mock Server Implementation](#mock-server-implementation)
    - [Price Validation Logic](#price-validation-logic)
      - [Live Provider Validation](#live-provider-validation)
      - [Mock Provider Validation](#mock-provider-validation)
    - [Configuration](#configuration)
      - [TOML Configuration](#toml-configuration)
      - [Programmatic Configuration](#programmatic-configuration)
    - [Usage in Tests](#usage-in-tests)
    - [Key Benefits](#key-benefits)

---

## Test Modification and Execution Guide

This guide explains how to set up and run system tests for Chainlink workflows using the CRE (Composable Runtime Environment) framework. It includes support for Docker and Kubernetes (via CRIB), multiple capabilities, and integration with Chainlink nodes and job distributor services.

---

For more information about the local CRE check its [README.md](../../../../core/scripts/cre/environment/README.md).

---

## 1. Running the Test

Before starting, you’ll need to configure your environment correctly. To do so execute the automated setup function:

```bash
cd core/scripts/cre/environment && go run . env setup
```

### Chainlink Node Image

The TOML config defines how Chainlink node images are used:

- **Default behavior**: Builds the Docker image from your current branch.

  ```toml
  [nodesets.node_specs.node]
    docker_ctx = "../../../.."
    docker_file = "core/chainlink.Dockerfile"
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

Only if you want to run the tests on non-default topology you need to set following variables before running the test:

- `CTF_CONFIGS` -- either `configs/workflow-gateway-don.toml` or `configs/workflow-gateway-capabilities-don.toml`
- `CRE_TOPOLOGY` -- either `workflow-gateway` or `workflow-gateway-capabilities`
- `CTF_LOG_LEVEL=debug` -- to display test debug-level logs

---

### `cron` Capability Binary

This binary is needed for tests using the cron capability.

**Option 1**: Use a CL node image that already includes the binary. Make sure it's available under `/usr/local/bin/cron` inside the image.

**Option 2**: Build the capability locally and copy it to: `core/scripts/cre/environment/binaries/cron`.

You can build it from [smartcontractkit/capabilities](https://github.com/smartcontractkit/capabilities) repository.

**Note**: Binary must be compiled for **Linux** and **amd64**.

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
  "args": ["-test.run", "Test_CRE_Suite"]
}
```

**CI behavior differs**: In CI, workflows and binaries are uploaded ahead of time, and images are injected via:

- `E2E_JD_VERSION`
- `E2E_TEST_CHAINLINK_IMAGE`
- `E2E_TEST_CHAINLINK_VERSION`

---

## 2. Using the CLI

Local CRE environment and documentation were migrated to [core/scripts/cre/environment/README.md](../../../../core/scripts/cre/environment/README.md).

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

8. **Install `telepresence`**.  CRIB requires telepresence for connecting to the services deployed within a cluster. You can install it with `brew install telepresenceio/telepresence/telepresence-oss` or with [other methods listed here](https://telepresence.io/docs/install/client#install-with-brew)

---

## 5. Setting Docker Images for CRIB Execution

CRIB does **not** support building Docker images from source during test runtime.

❌ Not allowed:

```toml
[nodesets.node_specs.node]
  docker_ctx = "../../../.."
  docker_file = "core/chainlink.Dockerfile"
```

✅ Required:

```toml
[nodesets.node_specs.node]
  image = "localhost:5001/chainlink:112b9323-plugins-cron"
```

`localhost:5001` is the repository name for local Kind registry. In order to push your image there you need:

- **tag a Docker image with prefix**, e.g. `docker tag chainlink-tmp:latest localhost:5001/chainlink:latest`
- **push it to the local registry**, e.g. `docker push localhost:5001/chainlink:latest`

**Also, it is crucial that you use an image, where default user is `chainlink`**. That's because Helm charts used for k8s deployment will start that container as that user and if your image was created using a different default user (e.g. `root`), then Chainlink application won't even start due to incorrect filesystem permissions. If you are building the image locally, use the following command:

```bash
# run in root of chainlink repository
docker build -t localhost:5001/chainlink:<your-tag>> --build-arg CHAINLINK_USER=chainlink -f core/chainlink.Dockerfile
```

### Notes

- All nodes in a single nodeset **must** use the same image.
- You must specify an image tag explicitly (e.g., `:v1.2.3`).

### Job Distributor (JD) Image in CRIB

#### AWS Provider

If you're working with AWS you will need to set the JD image URL in the `overrides.toml`

```toml
[jd]
  image = "<PROD_ECR_REGISTRY_URL>/job-distributor:0.12.7"
```

Replace `<PROD_ECR_REGISTRY_URL>` placeholder with the actual value.

#### Kind Provider

When working with kind provider, it will require pulling and pushing an image to local registry, similar as with CL node explained before.

```shell
docker pull <PROD_ECR_REGISTRY_URL>/job-distributor:0.12.7
docker tag <PROD_ECR_REGISTRY_URL>/job-distributor:0.12.7 localhost:5001/job-distributor:0.12.7
docker push localhost:5001/job-distributor:0.12.7
```

Now, you can set:

```toml
[jd]
  image = "localhost:5001/job-distributor:0.12.7"
```

---

## 6. Running Tests in Local Kubernetes (`kind`)

# Caution: this section is outdated. We no longer use `devspace`

### Docker Registry Setup

Pull the required local registry image:

```bash
docker pull registry:2
```

### Hostname Routing with `/etc/hosts`

CRIB dynamically creates hostname entries, but modifying `/etc/hosts` requires root. Tests fail if routing isn't set up.

#### Solutions

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

# Caution: this section is outdated. We no longer use `devspace`

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

# Caution: this section is partially outdated. We no longer use `devspace`

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

## 11. Using a New Workflow

This section explains how to compile, upload, and register workflows in the CRE testing framework. The process involves compiling Go workflow source code to WebAssembly, copying files to containers, and registering with the blockchain contract.

### Workflow Compilation Process

The workflow compilation process follows these steps:

1. **Source Code Preparation**: Ensure your workflow source code is in Go and follows the CRE workflow structure
2. **Compilation**: Use `creworkflow.CompileWorkflow()` to compile Go code to WebAssembly
3. **Compression**: The compiled WASM is automatically compressed using Brotli and base64 encoded
4. **File Management**: Temporary files are cleaned up automatically

#### Compilation Example

```go
workflowFileLocation := "path/to/your/workflow/main.go"
workflowName := "my-workflow-" + uuid.New().String()[0:4]

// Compile workflow to compressed WASM
compressedWorkflowWasmPath, compileErr := creworkflow.CompileWorkflow(workflowFileLocation, workflowName)
require.NoError(t, compileErr, "failed to compile workflow")

// Cleanup temporary files
t.Cleanup(func() {
    wasmErr := os.Remove(compressedWorkflowWasmPath)
    if wasmErr != nil {
        framework.L.Warn().Msgf("failed to remove workflow wasm file %s: %s", compressedWorkflowWasmPath, wasmErr.Error())
    }
})
```

#### Compilation Requirements

- **Workflow Name**: Must be at least 10 characters long
- **Go Environment**: Requires `go mod tidy` to be run in the workflow directory
- **Target Platform**: Compiles for `GOOS=wasip1` and `GOARCH=wasm`
- **Output Format**: Produces `.wasm.br.b64` files (compressed and base64 encoded)

### Workflow Configuration

Workflows may require configuration files that define their runtime parameters. Configuration is optional and depends on the specific workflow implementation:

#### Creating Configuration Files (Optional)

```go
func createConfigFile(feedsConsumerAddress common.Address, workflowName, feedID, dataURL, writeTargetName string) (string, error) {
    workflowConfig := portypes.WorkflowConfig{
        ComputeConfig: portypes.ComputeConfig{
            FeedID:                feedID,
            URL:                   dataURL,
            DataFeedsCacheAddress: feedsConsumerAddress.Hex(),
            WriteTargetName:       writeTargetName,
        },
    }

    configMarshalled, err := yaml.Marshal(workflowConfig)
    if err != nil {
        return "", errors.Wrap(err, "failed to marshal workflow config")
    }

    outputFile := workflowName + "_config.yaml"
    if err := os.WriteFile(outputFile, configMarshalled, 0644); err != nil {
        return "", errors.Wrap(err, "failed to write output file")
    }

    outputFileAbsPath, err := filepath.Abs(outputFile)
    if err != nil {
        return "", errors.Wrap(err, "failed to get absolute path of the config file")
    }

    return outputFileAbsPath, nil
}
```

### File Distribution to Containers

After compilation, workflow files must be distributed to the appropriate containers:

#### Copying Files to Containers

```go
containerTargetDir := "/home/chainlink/workflows"

// Copy compiled workflow binary
workflowCopyErr := creworkflow.CopyArtifactsToDockerContainers(
    containerTargetDir,
    "workflow-node",
    compressedWorkflowWasmPath, workflowConfigFilePath
)
require.NoError(t, workflowCopyErr, "failed to copy workflow to docker containers")
```

#### Container Discovery

The framework automatically discovers containers by name pattern:

- Uses Docker API to list running containers
- Matches container names against the provided pattern
- Copies files to all matching containers
- Creates target directories if they don't exist

### Workflow Registration

Workflows are registered with the blockchain contract using the `RegisterWithContract` function:

#### Registration Process

```go
workflowID, registerErr := creworkflow.RegisterWithContract(
    t.Context(),
    sethClient,                    // Blockchain client
    workflowRegistryAddress,       // Contract address
    donID,                        // DON identifier
    workflowName,                 // Unique workflow name
    "file://" + compressedWorkflowWasmPath,  // Binary URL
    ptr.Ptr("file://" + workflowConfigFilePath), // Config URL
    nil,                          // Secrets URL (optional)
    &containerTargetDir,          // Container artifacts directory
)
require.NoError(t, registerErr, "failed to register workflow")
```

#### Registration Parameters

- **Context**: Test context for timeout handling
- **Seth Client**: Blockchain client for contract interaction
- **Registry Address**: Workflow Registry contract address
- **DON ID**: Decentralized Oracle Network identifier
- **Workflow Name**: Unique identifier for the workflow
- **Binary URL**: Path to the compiled workflow binary on the host machine (used to read and calculate workflow ID)
- **Config URL**: Path to the workflow configuration file on the host machine (optional, used to read and calculate workflow ID)
- **Secrets URL**: Path to encrypted secrets on the host machine (optional)
- **Artifacts Directory**: Container directory where workflow files are stored (e.g., `/home/chainlink/workflows`)

#### URL Resolution Process

The `RegisterWithContract` function processes URLs as follows:

1. **Host Paths**: Binary URL, Config URL, and Secrets URL are paths on the host machine
2. **File Reading**: The function reads these files to calculate the workflow ID and validate content
3. **Container Path Construction**: If `artifactsDirInContainer` is provided, the function constructs container paths by:
   - Extracting the filename from the host path using `filepath.Base()`
   - Joining it with the artifacts directory: `file://{artifactsDir}/{filename}`
4. **Contract Registration**: The constructed container paths are registered with the blockchain contract

**Important**: The `Artifacts Directory` must match the `CRE.WorkflowFetcher.URL` configuration in your TOML file:

```toml
[CRE.WorkflowFetcher]
URL = "file:///home/chainlink/workflows"
```

This ensures that the Chainlink nodes can locate and load the workflow files from the correct container path.

> The Chainlink node can only load workflow files from the local filesystem if `WorkflowFetcher` uses the `file://` prefix. Right now, it cannot read workflow files from both the local filesystem and external sources (like S3 or web servers) at the same time.

### Complete Workflow Setup Example

Here's a complete example of setting up a workflow:

```go
func setupWorkflow(t *testing.T, workflowSourcePath, workflowName string, config *portypes.WorkflowConfig) {
    // 1. Compile workflow
    compressedWorkflowWasmPath, compileErr := creworkflow.CompileWorkflow(workflowSourcePath, workflowName)
    require.NoError(t, compileErr, "failed to compile workflow")

    // 2. Create configuration file (optional)
    var configFilePath string
    if config != nil {
        configData, err := yaml.Marshal(config)
        require.NoError(t, err, "failed to marshal config")

        configFilePath = workflowName + "_config.yaml"
        err = os.WriteFile(configFilePath, configData, 0644)
        require.NoError(t, err, "failed to write config file")
    }

    // 3. Copy files to containers
    containerTargetDir := "/home/chainlink/workflows"
    err := creworkflow.CopyArtifactsToDockerContainers(compressedWorkflowWasmPath, "workflow-node", containerTargetDir)
    require.NoError(t, err, "failed to copy workflow binary")

    if configFilePath != "" {
        err = creworkflow.CopyArtifactsToDockerContainers(configFilePath, "workflow-node", containerTargetDir)
        require.NoError(t, err, "failed to copy config file")
    }

    // 4. Register with contract
    var configURL *string
    if configFilePath != "" {
        configURL = ptr.Ptr("file://" + configFilePath)
    }

    workflowID, registerErr := creworkflow.RegisterWithContract(
        t.Context(),
        sethClient,
        workflowRegistryAddress,
        donID,
        workflowName,
        "file://" + compressedWorkflowWasmPath,
        configURL,
        nil, // secrets URL (optional)
        &containerTargetDir,
    )
    require.NoError(t, registerErr, "failed to register workflow")

    // 5. Cleanup
    t.Cleanup(func() {
        os.Remove(compressedWorkflowWasmPath)
        if configFilePath != "" {
            os.Remove(configFilePath)
        }
    })
}
```

---

### 12. Workflow Secrets

Workflow secrets provide a secure way to pass sensitive data (like API keys, private keys, or database credentials) to workflows running on Chainlink nodes. The secrets are encrypted using each node's public encryption key and can only be decrypted by the intended recipient nodes.

#### How Secrets Work

1. **Configuration**: Define which environment variables contain your secrets
2. **Encryption**: Secrets are encrypted using each DON node's public encryption key
3. **Distribution**: Encrypted secrets are distributed to the appropriate nodes
4. **Decryption**: Each node decrypts only the secrets intended for it

#### Creating Secrets Configuration

Create a YAML file that maps secret names to environment variables:

```yaml
# secrets.yaml
secretsNames:
  API_KEY_SECRET:
    - API_KEY_ENV_VAR_ALL
  DATABASE_PASSWORD:
    - DB_PASSWORD_ENV_VAR_ALL
  PRIVATE_KEY:
    - PRIVATE_KEY_ENV_VAR_ALL
```

#### Environment Variable Naming

- Use `_ENV_VAR_ALL` suffix for secrets shared across all nodes in the DON
- Use `_ENV_VAR_NODE_{NODE_INDEX}` suffix for node-specific secrets (where `NODE_INDEX` is the sequential position of the node in the DON: 0, 1, 2, etc.)
- Environment variables must be set before running the workflow registration

**Note**: `NODE_INDEX` refers to the node's position in the DON (0-based indexing), not the P2P ID. For example:

- `API_KEY_ENV_VAR_NODE_0` for the first node in the DON
- `API_KEY_ENV_VAR_NODE_1` for the second node in the DON
- `API_KEY_ENV_VAR_NODE_2` for the third node in the DON

#### Using Secrets in Workflows

```go
// 1. Set environment variables
os.Setenv("API_KEY_ENV_VAR_ALL", "your-api-key-here")
os.Setenv("DB_PASSWORD_ENV_VAR_ALL", "your-db-password")

// 2. Prepare encrypted secrets
secretsFilePath := "path/to/secrets.yaml"
encryptedSecretsPath, err := creworkflow.PrepareSecrets(
    sethClient,
    donID,
    capabilitiesRegistryAddress,
    workflowOwnerAddress,
    secretsFilePath,
)
require.NoError(t, err, "failed to prepare secrets")

// 3. Copy encrypted secrets to containers
err = creworkflow.CopyArtifactsToDockerContainers(
    encryptedSecretsPath,
    "workflow-node",
    "/home/chainlink/workflows",
)
require.NoError(t, err, "failed to copy secrets to containers")

// 4. Register workflow with secrets
workflowID, registerErr := creworkflow.RegisterWithContract(
    ctx,
    sethClient,
    workflowRegistryAddress,
    donID,
    workflowName,
    "file://" + compressedWorkflowWasmPath,
    configURL,
    &secretsURL, // Pass the encrypted secrets file path
    &containerTargetDir,
)
require.NoError(t, registerErr, "failed to register workflow")
```

#### Secrets Encryption Process

The `PrepareSecrets` function performs these steps:

1. **Load Configuration**: Parses the secrets YAML file
2. **Read Environment Variables**: Loads secret values from environment variables
3. **Get Node Information**: Retrieves node public keys from the Capabilities Registry contract
4. **Filter DON Nodes**: Identifies nodes that belong to the specific DON
5. **Encrypt Secrets**: Encrypts secrets using each node's public encryption key
6. **Generate Metadata**: Creates metadata including encryption keys and node assignments
7. **Save Encrypted File**: Outputs a JSON file with encrypted secrets and metadata

#### Encrypted Secrets File Structure

The generated encrypted secrets file contains:

```json
{
  "encryptedSecrets": {
    "node_p2p_id_1": "encrypted_secret_for_node_1",
    "node_p2p_id_2": "encrypted_secret_for_node_2"
  },
  "metadata": {
    "workflowOwner": "0x...",
    "capabilitiesRegistry": "0x...",
    "donId": "1",
    "dateEncrypted": "2024-01-01T00:00:00Z",
    "nodePublicEncryptionKeys": {
      "node_p2p_id_1": "public_key_1",
      "node_p2p_id_2": "public_key_2"
    },
    "envVarsAssignedToNodes": {
      "node_p2p_id_1": ["API_KEY_ENV_VAR_ALL"],
      "node_p2p_id_2": ["API_KEY_ENV_VAR_ALL"]
    }
  }
}
```

#### Security Considerations

- **Node-Specific Encryption**: Each node can only decrypt secrets intended for it
- **DON Isolation**: Secrets are encrypted per DON and cannot be shared across different DONs
- **Environment Variables**: Secrets are loaded from environment variables, not hardcoded
- **Temporary Files**: Encrypted secrets files are automatically cleaned up after registration

#### Complete Example

```go
func setupWorkflowWithSecrets(t *testing.T, workflowSourcePath, workflowName, secretsConfigPath string) {
    // Set environment variables with your secrets
    os.Setenv("API_KEY_ENV_VAR_ALL", "your-actual-api-key")
    os.Setenv("DB_PASSWORD_ENV_VAR_ALL", "your-actual-db-password")

    // Compile workflow
    compressedWorkflowWasmPath, err := creworkflow.CompileWorkflow(workflowSourcePath, workflowName)
    require.NoError(t, err, "failed to compile workflow")

    // Prepare encrypted secrets
    encryptedSecretsPath, err := creworkflow.PrepareSecrets(
        sethClient,
        donID,
        capabilitiesRegistryAddress,
        workflowOwnerAddress,
        secretsConfigPath,
    )
    require.NoError(t, err, "failed to prepare secrets")

    // Copy files to containers
    containerTargetDir := "/home/chainlink/workflows"
    err = creworkflow.CopyArtifactsToDockerContainers(compressedWorkflowWasmPath, "workflow-node", containerTargetDir)
    require.NoError(t, err, "failed to copy workflow")

    err = creworkflow.CopyArtifactsToDockerContainers(encryptedSecretsPath, "workflow-node", containerTargetDir)
    require.NoError(t, err, "failed to copy secrets")

    // Register workflow with secrets
    secretsURL := "file://" + encryptedSecretsPath
    workflowID, registerErr := creworkflow.RegisterWithContract(
        t.Context(),
        sethClient,
        workflowRegistryAddress,
        donID,
        workflowName,
        "file://" + compressedWorkflowWasmPath,
        nil, // config URL (optional)
        &secretsURL,
        &containerTargetDir,
    )
    require.NoError(t, registerErr, "failed to register workflow")

    // Cleanup
    t.Cleanup(func() {
        os.Remove(compressedWorkflowWasmPath)
        os.Remove(encryptedSecretsPath)
    })
}
```

---

### 13. YAML Workflows (Data Feeds DSL)

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

## 14. Adding a New Test to the CI

The CRE system tests use **auto-discovery** to automatically find and run all tests in the `system-tests/tests/smoke/cre` directory. This means you don't need to manually register your test in any CI configuration files.

### How Auto-Discovery Works

The CI workflow (`.github/workflows/cre-system-tests.yaml`) automatically:

1. **Discovers Tests**: Uses `go test -list .` to find all test functions in the package
2. **Creates Test Matrix**: Generates a matrix with each test and supported DON topologies
3. **Runs Tests**: Executes each test with different configurations automatically

### Test Architecture Pattern

The CRE system tests follow a **separated architecture pattern** where:

- **Environment Creation**: The environment (DONs, contracts, nodes) is created once per topology
- **Test Execution**: Multiple tests run on the same environment instance
- **Shared State**: Tests can leverage the same deployed contracts and node infrastructure

This pattern allows for efficient resource usage and enables running the same test logic across different DON topologies without recreating the entire environment for each test.

#### Supported DON Topologies

Each test is automatically run with these three topologies:

- **workflow**: Uses `configs/workflow-don.toml,configs/ci-config.toml`
- **workflow-gateway**: Uses `configs/workflow-gateway-don.toml,configs/ci-config.toml`
- **workflow-gateway-capabilities**: Uses `configs/workflow-gateway-capabilities-don.toml,configs/ci-config.toml`

### What You Need to Do

#### 1. Create Your Test Function

Simply add your test function to any `.go` file in the `system-tests/tests/smoke/cre` directory:

```go
func Test_CRE_My_New_Workflow(t *testing.T) {
    // Your test implementation
    // The CI will automatically discover and run this test
}
```

#### 2. Follow Test Naming Convention

Use the `Test_CRE_` prefix for your test functions to ensure they're properly identified:

```go
func Test_CRE_MyWorkflow(t *testing.T)     // ✅ Good
func Test_CRE_AnotherWorkflow(t *testing.T) // ✅ Good
func TestMyWorkflow(t *testing.T)          // ❌ Will be discovered but not recommended
```

#### 3. Use Standard Test Structure

Follow the existing test patterns in the directory. Note that the environment is created separately and shared across tests:

```go
func Test_CRE_My_New_Workflow(t *testing.T) {
    // 1. Set configuration if needed
    confErr := setConfigurationIfMissing("path/to/config.toml", "topology")
    require.NoError(t, confErr, "failed to set configuration")

    // 2. Load existing environment (created by CI)
    in, err := framework.Load[environment.Config](nil)
    require.NoError(t, err, "couldn't load environment state")

    // 3. Your test logic using the shared environment
    // The environment (DONs, contracts, nodes) is already set up and ready to use
    // ...
}
```

**Important**: Your test should be designed to work with any of the supported DON topologies. The same test logic should ideally be compatible with:

- `workflow` topology
- `workflow-gateway` topology
- `workflow-gateway-capabilities` topology

This ensures maximum test coverage and validates your workflow across different deployment configurations.

### CI Configuration Details

The auto-discovery process works as follows:

```yaml
# From .github/workflows/cre-system-tests.yaml
- name: Define test matrix
  run: |
    tests=$(go test github.com/smartcontractkit/chainlink/system-tests/tests/smoke/cre -list . | grep -v "ok" | grep -v "^$" | jq -R -s 'split("\n")[:-1] | map([{"test_name": ., "topology": "workflow", "configs":"configs/workflow-don.toml,configs/ci-config.toml"}, {"test_name": ., "topology": "workflow-gateway", "configs":"configs/workflow-gateway-don.toml,configs/ci-config.toml"}, {"test_name": ., "topology": "workflow-gateway-capabilities", "configs":"configs/workflow-gateway-capabilities-don.toml,configs/ci-config.toml"}]) | flatten')
```

### Environment Setup

The CI automatically sets up the test environment:

- **Dependencies**: Downloads required capability binaries
- **Local CRE**: Starts the CRE environment with the specified topology (once per topology)
- **Configuration**: Uses the appropriate TOML configuration files
- **Artifacts**: Handles test logs and artifacts automatically
- **Shared Infrastructure**: All tests within the same topology share the same environment instance

This approach ensures that:

- Environment creation overhead is minimized
- Tests can leverage shared contracts and node infrastructure
- The same test logic can be validated across different DON configurations

### Test Execution

Each test runs with:

```bash
go test github.com/smartcontractkit/chainlink/system-tests/tests/smoke/cre \
  -v -run "^(${TEST_NAME})$" \
  -timeout ${TEST_TIMEOUT} \
  -count=1 \
  -test.parallel=1 \
  -json
```

### Important Notes

- **No Manual Registration**: You don't need to add your test to any CI configuration files
- **Automatic Matrix**: Each test runs with all three DON topologies automatically
- **Standard Configurations**: Uses the existing TOML configuration files
- **Dependency Management**: Capabilities and dependencies are handled automatically
- **Logging**: Test logs are automatically captured and uploaded on failure

### Troubleshooting

If your test isn't being discovered:

1. **Check Function Name**: Ensure it starts with `Test_CRE_`
2. **Check Location**: Ensure it's in the `system-tests/tests/smoke/cre` directory
3. **Check Syntax**: Ensure the test function compiles without errors
4. **Check Dependencies**: Ensure all required dependencies are available

> **Note**: The auto-discovery system eliminates the need for manual CI configuration, making it much easier to add new tests to the CI pipeline.

---

## 15. Price Data Source

The CRE system supports both **live** and **mocked** price feeds through a unified `PriceProvider` interface. This allows for flexible testing scenarios while maintaining consistent behavior across different data sources.

### PriceProvider Interface

The system uses a common interface that abstracts price data source logic:

```go
type PriceProvider interface {
    URL() string
    NextPrice(feedID string, price *big.Int, elapsed time.Duration) bool
    ExpectedPrices(feedID string) []*big.Int
    ActualPrices(feedID string) []*big.Int
    AuthKey() string
}
```

### Live Price Source (TrueUSDPriceProvider)

For integration testing with real data:

```go
// Create a live price provider
priceProvider := NewTrueUSDPriceProvider(testLogger, feedIDs)

// The provider uses the live TrueUSD API
// URL: https://api.real-time-reserves.verinumus.io/v1/chainlink/proof-of-reserves/TrueUSD
```

**Characteristics:**

- Uses real-time data from the TrueUSD API
- No authentication required
- Validates that prices are non-zero
- Tracks actual prices received by the workflow
- Limited validation capabilities (can only check for non-zero values)

### Mocked Price Source (FakePriceProvider)

For local testing and controlled scenarios:

```go
// Create a fake price provider
fakeInput := &fake.Input{Port: 8171}
authKey := "your-auth-key"
priceProvider, err := NewFakePriceProvider(testLogger, fakeInput, authKey, feedIDs)
require.NoError(t, err, "failed to create fake price provider")
```

**Characteristics:**

- Generates random prices for testing
- Provides controlled price sequences
- Validates exact price matches
- Supports authentication
- Tracks both expected and actual prices

### Mock Server Implementation

The fake price provider sets up a mock HTTP server that:

1. **Generates Random Prices**: Creates random price values between 1.00 and 200.00
2. **Supports Authentication**: Validates Authorization headers
3. **Responds to Feed Queries**: Handles `GET` requests with `feedID` parameter
4. **Returns Structured Data**: Provides JSON responses in the expected format:

```json
{
  "accountName": "TrueUSD",
  "totalTrust": 123.45,
  "ripcord": false,
  "updatedAt": "2024-01-01T00:00:00Z"
}
```

### Price Validation Logic

Both providers implement smart validation:

#### Live Provider Validation

```go
func (l *TrueUSDPriceProvider) NextPrice(feedID string, price *big.Int, elapsed time.Duration) bool {
    // Wait for non-zero price
    if price == nil || price.Cmp(big.NewInt(0)) == 0 {
        return true // Continue waiting
    }
    // Price found, stop waiting
    return false
}
```

#### Mock Provider Validation

```go
func (f *FakePriceProvider) NextPrice(feedID string, price *big.Int, elapsed time.Duration) bool {
    // Check if this is a new price we haven't seen
    if !f.priceAlreadyFound(feedID, price) {
        // Record the new price
        f.actualPrices[feedID] = append(f.actualPrices[feedID], price)

        // Move to next expected price
        f.priceIndex[feedID] = ptr.Ptr(len(f.actualPrices[feedID]))

        // Continue if more prices expected
        return len(f.actualPrices[feedID]) < len(f.expectedPrices[feedID])
    }
    return true // Continue waiting
}
```

### Configuration

#### TOML Configuration

The price provider is **not** configured directly in TOML. Instead, the TOML only configures the fake server port:

```toml
[fake]
  port = 8171
```

#### Programmatic Configuration

Price providers are created programmatically in your Go test code:

```go
// For live price provider (no TOML configuration needed)
priceProvider := NewTrueUSDPriceProvider(testLogger, feedIDs)

// For fake price provider (uses port from TOML [fake] section)
fakeInput := &fake.Input{Port: 8171} // or use in.Fake from loaded config
authKey := "your-auth-key"
priceProvider, err := NewFakePriceProvider(testLogger, fakeInput, authKey, feedIDs)
require.NoError(t, err)
```

### Usage in Tests

```go
func Test_CRE_Price_Feed(t *testing.T) {
    feedIDs := []string{
        "018e16c39e000320000000000000000000000000000000000000000000000000",
        "018e16c38e000320000000000000000000000000000000000000000000000000",
    }

    // Choose your provider
    var priceProvider PriceProvider

    if useLiveProvider {
        priceProvider = NewTrueUSDPriceProvider(testLogger, feedIDs)
    } else {
        fakeInput := &fake.Input{Port: 8171}
        priceProvider, err = NewFakePriceProvider(testLogger, fakeInput, "auth-key", feedIDs)
        require.NoError(t, err)
    }

    // Use the provider in your workflow configuration
    workflowConfig := &portypes.WorkflowConfig{
        ComputeConfig: portypes.ComputeConfig{
            FeedID: feedIDs[0],
            URL:    priceProvider.URL(),
            // ... other config
        },
    }

    // Validate price updates
    assert.Eventually(t, func() bool {
        price := getLatestPrice(feedID)
        return !priceProvider.NextPrice(feedID, price, time.Since(startTime))
    }, timeout, interval)
}
```

### Key Benefits

1. **Unified Interface**: Same API for both live and mocked sources
2. **Flexible Testing**: Easy switching between real and fake data
3. **Controlled Validation**: Mock provider enables precise price validation
4. **Authentication Support**: Mock server supports auth for realistic testing
5. **Price Tracking**: Both providers track actual prices received by workflows

---
