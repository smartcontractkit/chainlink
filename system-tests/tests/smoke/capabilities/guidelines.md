# Test Modification and Execution Guide

## Table of Contents
1. [How to Run the Test](#how-to-run-the-test)
2. [Adding a New Capability](#adding-a-new-capability)
   - [Copying the Binary to the Container](#copying-the-binary-to-the-container)
   - [Adding support for the new capability in the testing code](#adding-support-for-the-new-capability-in-the-testing-code)
     - [Defining new bitmask flag representing the capability](#defining-new-bitmask-flag-representing-the-capability)
     - [Defining additional node configuration](#defining-additional-node-configuration)
     - [Defining job spec for the new capability](#defining-job-spec-for-the-new-capability)
     - [Registering the capability in the Capabilities Registry contract](#registering-the-capability-in-the-capabilities-registry-contract)
3. [Using a New Workflow](#using-a-new-workflow)
   - [Test Uploads the Binary](#test-uploads-the-binary)
   - [Workflow Configuration](#workflow-configuration)
   - [Workflow Secrets](#workflow-secrets)
   - [Manual Upload of the Binary](#manual-upload-of-the-binary)
4. [Deployer Address or Deployment Sequence Changes](#deployer-address-or-deployment-sequence-changes)
5. [Multiple DONs](#multiple-dons)
   - [DON type](#don-type)
   - [Capabilities](#capabilities)
   - [HTTP port range start](#http-port-range-start)
   - [DB port](#db-port)
6. [Price Data Source](#price-data-source)
   - [Live Source](#live-source)
   - [Mocked Data Source](#mocked-data-source)
7. [Using a Specific Docker Image for Chainlink Node](#using-a-specific-docker-image-for-chainlink-node)
8. [Troubleshooting](#troubleshooting)
   - [Chainlink Node migrations fail](#chainlink-node-migrations-fail)
   - [Chainlink image not found in local Docker registry](#chainlink-image-not-found-in-local-docker-registry)

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

Test also expects you to have the Job Distributor image available locally. By default, `environment.toml` expects image tagged as `job-distributor:latest`. The easiest way to get it, is to clone the Job Distributor repository and build it locally with:
```bash
docker build -t job-distributor:latest -f e2e/Dockerfile.e2e .
```

Alternatively, if you have access to the Docker image repository where it's stored you can modify `environment.toml` with the name of the image stored there.

In the CI test code modifies the config during runtime to production JD image hardcoded in the [.github/e2e-tests.yml](.github/e2e-tests.yml) file as `E2E_JD_VERSION` env var.

---

## Adding a New Capability

To add a new capability to the test, follow these steps:

1. Copy the capability binary to the Chainlink node’s Docker container (must be in `linux/amd64` format).
   - You can skip this step if the capability is already included in the Chainlink image you are using or if it's built-in.
2. Add support for the new capability in the testing code:
   - Define a new `CapabilityFlag` representing the capability.
   - (Optional) Define additional node configuration if required.
   - Define the job spec for the new capability.
   - Register the capability in the Capabilities Registry contract.
3. Update the TOML configuration to assign the new capability to one of the DONs.

Once these steps are complete, you can run a workflow that requires the new capability.

Let's assume we want to add a capability that represents writing to Aptos chain.

### Copying the Binary to the Container

The test configuration is defined in a TOML file (e.g. `environment.toml`), which specifies properties for Chainlink nodes. The `capabilities` property of the `node_specs.node` determines which binaries are copied to the container:

```toml
  [[nodeset.node_specs]]

    [nodeset.node_specs.node]
      capabilities = ["./aptos_linux_amd64"]
```

This instructs the framework to copy `./aptos_linux_amd64` to the container’s `/home/capabilities/` directory, making it available as `/home/capabilities/aptos_linux_amd64`.

> **Note:** Copying the binary to the bootstrap node is unnecessary since it does not handle capability-related tasks.

### Adding Support for the New Capability in the Testing Code

#### Defining a CapabilityFlag for the Capability

The testing code uses string flags to map DON capabilities to node configuration, job creation, and the Capabilities Registry contract. This means that adding a new capability requires defining a unique flag. Let's name our capability flag as `WriteAptosCapability`.

First, define the new flag:

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

This ensures the TOML configuration correctly maps each DON to its capabilities.

Optionally, add the new flag to the default capabilities used in a single DON setup:

```go
var (
	// Add new capabilities here as well, if single DON should have them by default
	SingleDonFlags = []string{"workflow", "capabilities", "ocr3", "cron", "custom-compute", "write-evm", "write-aptos"}
                                                                                                        // <------------ New entry
)
```

Now that the flag is defined, let's configure the nodes and jobs.

#### Defining Additional Node Configuration

This step is optional, as not every capability requires additional node configuration. However, writing to the Aptos chain does. Depending on the capability, adjustments might be needed for the bootstrap node, the workflow nodes, or all nodes.

The following code snippet adds the required settings:

```go
if hasFlag(flags, WriteAptosCapability) {
  writeAptosConfig := fmt.Sprintf(`
    # Required for initializing the capability
    [Aptos.Workflow]
    # Configuration parameters
    param_1 = "%s"
    `,
    "some value",
  )
  workerNodeConfig += writeAptosConfig
}
```

This is a placeholder snippet—you should replace it with the actual configuration required for the capability. Ensure it is added before restarting the nodes.

#### Defining a Job Spec for the New Capability

Unlike node configuration, defining a new job spec is always required for a new capability. Jobs should only be added to worker nodes.

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

> **Note:** If the new capability requires a different job type, you may need to update the Chainlink Node code. If it works with `standardcapabilities`, no changes are necessary.

#### Registering the Capability in the Capabilities Registry Contract

The final step is adding support for registration of the capability with the Capabilities Registry contract:

```go
if hasFlag(donTopology.Flags, WriteAptosCapability) {
  capabilities = append(capabilities, keystone_changeset.DONCapabilityWithConfig{
    Capability: kcr.CapabilitiesRegistryCapability{
      LabelledName:   "write_aptos-testnet",          // <------- Ensure correct name
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

> **Note:** Since this test code is constantly evolving, no specific line numbers or function names are provided.

## Using a New Workflow

To test a new workflow, you have two options:

1. Compile the workflow to a WASM binary and upload it to Gist **inside the test**.
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
> **Note:** Currently, we do not allow to update the configuration alone. Each configuration change is treated as workflow change and thus requires following the **upload mode**.

---

### Workflow Secrets
Currently, workflow secrets are **not supported**.

---

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

You can choose to use one or multiple DONs. Configuring multiple DONs requires only TOML modifications, assuming they use capabilities already supported in the testing code.

Currently, the supported capabilities are:
- `cron`
- `ocr3`
- `custom-compute`
- `write-evm`

To enable multi-DON support, update the configuration file by:
- Defining a new nodeset.
- Explicitly assigning capabilities to each nodeset.
- Copying the required capabilities to the containers (if they are not built into the image already).

Here’s an example configuration for a nodeset that only supports writing to an EVM chain:

```toml
[[nodesets]]
  don_type = "capabilities"
  name = "capabilities"
  capabilities = ["write-evm"]
  nodes = 5
  override_mode = "each"
  http_port_range_start = 10200

  [nodesets.db]
    image = "postgres:12.0"
    port = 13100

  [[nodesets.node_specs]]

    [nodesets.node_specs.node]
      image = "some-CL-image"
      # Rest of the node configuration follows

  # Additional nodes configuration follows
```

### Key Considerations
When configuring multiple DONs, keep the following in mind:
- **DON Type**
- **Capabilities List**
- **HTTP Port Range Start**
- **Database (DB) Port**

### DON Type

Two types of DONs are supported:
- `workflow`
- `capabilities`

There should only be **one** `workflow` DON, but multiple `capabilities` DONs can be defined.

### Capabilities

- In a **single DON setup**, you can omit the capabilities list, as all known capabilities will be assigned and configured by default (as long as they are included in the `SingleDonFlags` constant).
- In a **multi-DON setup**, you must explicitly define the capabilities for each DON.

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

### HTTP Port Range Start

Each node exposes a port to the host. To prevent port conflicts, assign a distinct range to each nodeset. A good practice is to separate port ranges by **50 or 100** between nodesets.

### Database (DB) Port

Similar to HTTP ports, ensure each nodeset has a unique database port.

For a working example of a multi-DON setup, refer to the [`environment-multi-don.toml`](environment-multi-don.toml) file.

---

## Price Data Source

The test supports both **live** and **mocked** data sources, configurable via TOML.

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
    prices = [182.9, 162.71, 172.02]
```

This configuration launches a mock server on **port 8171** on the host machine. It will return prices in the sequence `[182.9, 162.71, 172.02]`. A new price is returned **only after the previous one has been observed in the consumer contract**. The test completes once all prices have been matched.

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
If you have the `ctf` CLI you can use following command: `ctf d rm`.

### Chainlink image not found in local Docker registry

If you are building the Chainlink image using the Dockerfile, image is successfuly built and yet nodes do not start, because image cannot be found in the local machine, simply restart your computer and try again.