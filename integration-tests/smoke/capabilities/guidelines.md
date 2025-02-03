# Test Modification and Execution Guide

## Table of Contents
1. [How to Run the Test](#how-to-run-the-test)
2. [Adding a New Capability](#adding-a-new-capability)
   - [Copying the Binary to the Container](#copying-the-binary-to-the-container)
   - [Registering the Capability in the Registry Contract](#registering-the-capability-in-the-registry-contract)
   - [Creating a New Capability Job](#creating-a-new-capability-job)
3. [Using a New Workflow](#using-a-new-workflow)
   - [Test Uploads the Binary](#test-uploads-the-binary)
   - [Workflow Configuration](#workflow-configuration)
   - [Workflow Secrets](#workflow-secrets)
   - [Manual Upload of the Binary](#manual-upload-of-the-binary)
4. [Deployer Address or Deployment Sequence Changes](#deployer-address-or-deployment-sequence-changes)
5. [Multiple DONs](#multiple-dons)
6. [Using a Specific Docker Image for Chainlink Node](#using-a-specific-docker-image-for-chainlink-node)

---

## How to Run the Test

The test requires several environment variables. Below is a launch configuration that can be used with the VCS:

```json
{
  "name": "Launch Capability Test",
  "type": "go",
  "request": "launch",
  "mode": "test",
  "program": "${workspaceFolder}/integration-tests/smoke/capabilities",
  "env": {
    "CTF_CONFIGS": "environment.toml",
    "PRIVATE_KEY": "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80",
    "GITHUB_GIST_API_TOKEN": "your-gist-read:write-fpat-token",
    "GITHUB_CAP_API_TOKEN": "your-capabilities-repo-content-read-fpat-token"
  },
  "args": [
    "-test.run",
    "TestKeystoneWithOCR3Workflow"
  ]
}
```

- **`GITHUB_READ_TOKEN`**: Required for downloading the `cron` capability binary and CRE CLI (if enabled). Requires `content:read` permission for `smartcontractkit/capabilities` and `smartcontractkit/dev-platform` repositories. Use a fine-grained personal access token (PAT) tied to the **organization’s GitHub account**.
- **`GIST_WRITE_TOKEN`**: Required only for compiling and uploading a new workflow. It needs `gist:read:write` permissions and should be a fine-grained PAT **tied to your personal GitHub account**.

Test also expects you to have the Job Distributor image available locally. By default, `environment.toml` expects image tagged as `jd-test-1:latest`. The easiest way to get it, is to clone the Job Distributor repository and build it locally with:
```bash
docker build -t jd-test-1 -f e2e/Dockerfile.e2e
```

Alternatively, if you have access to the Docker image repository where it's stored you can modify `environment.toml` with the name of the image stored there.

---

## Adding a New Capability

To use a new capability in the test, follow these three steps:

1. Copy the capability binary to the Chainlink node’s Docker container (must be in `linux/amd64` format).
2. Register the capability in the capability registry contract.
3. Create a new job on each worker node to launch the capability.

Only after completing these steps can you run a workflow that requires the new capability.

### Copying the Binary to the Container

The test configuration is defined in a TOML file (`environment.toml`), which specifies properties for the Chainlink nodes. The `capabilities` property determines which binaries are copied to the container:

```toml
  [[nodeset.node_specs]]

    [nodeset.node_specs.node]
      capabilities = ["./amd64_cron"]
```

This instructs the framework to copy `./amd64_cron` to the container’s `/home/capabilities/` directory, making it available as `/home/capabilities/amd64_cron`.

> **Note:** While copying the binary to the bootstrap node causes no harm, it is unnecessary since the bootstrap node does not handle capability-related tasks.

### Registering the Capability in the Registry Contract

Next, register the new capability in the **capabilities registry contract** to make it accessible to the workflow DON.

Example definition for a cron capability:

```go
{
  LabelledName:   "cron-trigger",
  Version:        "1.0.0",
  CapabilityType: uint8(0), // trigger
}
```

Some capabilities may also require `ResponseType` or `ConfigurationContract`. Check with the capability author for the necessary values and ensure the correct capability type is set.

Since this test's code is a living creature we do not indicate any line numbers or functions names. It's advised that you search for the example used above or usages of `CapabilitiesRegistryCapability` type.

### Creating a New Capability Job

Now that worker nodes have access to the binary and it is registered in the contract, create a capability job for each worker node.

Example cron capability job:

```toml
type = "standardcapabilities"
schemaVersion = 1
name = "cron-capabilities"
forwardingAllowed = false
command = "/home/capabilities/amd64_cron"
config = ""
```

Since this is **not** a built-in capability, you must specify the binary’s path inside the Docker container.

Also, remember that the exact content of the job is hard to know beforehand and you should work with the capability creator to figure out what exactly is required.

---

## Using a New Workflow

To test a new workflow, you have two options:

1. Compile the workflow to a WASM binary and upload it to Gist inside the test.
2. Manually upload the binary and specify the workflow URL in the test configuration.

### Test Uploads the Binary

For the test to compile and upload the binary, modify your TOML configuration:

```toml
[workflow_config]
  use_cre_cli = true
  should_compile_new_workflow = true
  workflow_folder_location = "path-to-folder-with-main.go-of-your-workflow"
```

### Workflow Configuration

If your workflow requires configuration, modify the test to create and pass the configuration data to CRE CLI:

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

### Workflow Secrets
Currently, workflow secrets are **not supported**.

### Manual Upload of the Binary

If you compiled and uploaded the binary yourself, set the following in your configuration:

```toml
[workflow_config]
  use_cre_cli = true
  should_compile_new_workflow = false

  [workflow_config.compiled_config]
    binary_url = "<binary-url>"
    config_url = "<config-url>"
```

Both URLs must be accessible by the bootstrap node.

---

## Deployer Address or Deployment Sequence Changes

By default, the test reuses an existing workflow and configuration. The feed consumer address remains the same **as long as the deployer address (`f39fd6e51aad88f6f4ce6ab8827279cfffb92266`) and contract deployment sequence do not change**.

If the deployer private key or deployment sequence changes, run the test in **upload mode**:

```toml
[workflow_config]
  use_cre_cli = true
  should_compile_new_workflow = true
  workflow_folder_location = "path-to-folder-with-main.go-of-your-workflow"
```

---

## Multiple DONs

Multiple DONs are **not supported**. A single DON serves as both the **workflow** and **capabilities DON**.

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

