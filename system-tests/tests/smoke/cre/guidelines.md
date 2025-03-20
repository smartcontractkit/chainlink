# Test Modification and Execution Guide

## Table of Contents
1. [How to Run the Test](#how-to-run-the-test)
2. [Adding a New Capability](#adding-a-new-capability)
   - [Copying the Binary to the Container](#copying-the-binary-to-the-container)
   - [Adding support for the new capability in the testing code](#adding-support-for-the-new-capability-in-the-testing-code)
     - [Defining a CapabilityFlag for the Capability](#defining-a-capabilityflag-for-the-capability)
     - [Defining Additional Node Configuration](#defining-additional-node-configuration)
     - [Defining a Job Spec for the New Capability](#defining-a-job-spec-for-the-new-capability)
     - [Registering the Capability in the Capabilities Registry Contract](#registering-the-capability-in-the-capabilities-registry-contract)
3. [Using a New Workflow](#using-a-new-workflow)
   - [Test Uploads the Binary](#test-uploads-the-binary)
   - [Workflow Configuration](#workflow-configuration)
   - [Workflow Secrets](#workflow-secrets)
   - [Manual Upload of the Binary](#manual-upload-of-the-binary)
4. [Deployer Address or Deployment Sequence Changes](#deployer-address-or-deployment-sequence-changes)
5. [Multiple DONs](#multiple-dons)
   - [DON Type](#don-type)
   - [Capabilities](#capabilities)
   - [HTTP Port Range Start](#http-port-range-start)
   - [Database (DB) Port](#database-db-port)
   - [Number of Nodes](#number-of-nodes)
6. [Price Data Source](#price-data-source)
   - [Live Source](#live-source)
   - [Mocked Data Source](#mocked-data-source)
7. [Using a Specific Docker Image for Chainlink Node](#using-a-specific-docker-image-for-chainlink-node)
8. [Troubleshooting](#troubleshooting)
   - [Chainlink Node migrations fail](#chainlink-node-migrations-fail)
   - [Chainlink image not found in local Docker registry](#chainlink-image-not-found-in-local-docker-registry)
9. [Docker vs Kubernetes (k8s)](#docker-vs-kubernetes-k8s)
10. [CRIB Requirements](#crib-requirements)
11. [Setting Docker Images for CRIB Execution](#setting-docker-images-for-crib-execution)
12. [Running Tests in Local Kubernetes (`kind`)](#running-tests-in-local-kubernetes-kind)
13. [CRIB Deployment Flow](#crib-deployment-flow)
14. [Switching from kind to AWS provider](#switching-from-kind-to-aws-provider)
15. [CRIB Limitations & Considerations](#crib-limitations--considerations)

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

Test also expects you to have the Job Distributor image available locally. By default, `environment-*.toml`'s expect image tagged as `job-distributor:0.9.0`. The easiest way to get it, is to clone the Job Distributor repository and build it locally with:
```bash
docker build -t job-distributor:0.9.0 -f e2e/Dockerfile.e2e .
```

Alternatively, if you have access to the Docker image repository where it's stored you can modify `environment-*.toml`'s with the name of the image stored there.

In the CI test code modifies the config during runtime to production JD image hardcoded in the [.github/e2e-tests.yml](.github/e2e-tests.yml) file as `E2E_JD_VERSION` env var.

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
- Each nodeset **must** use the same Docker image.
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
6. **Start Job Distributor**:
   - Set environment variable: `JOB_DISTRIBUTOR_IMAGE_TAG`.
   - Deploy with `devspace run deploy-jd`.
   - Read JD URLs from `jd-url.json`.
7. **Create Jobs & Configure CRE Contracts** (same as Docker).

---

## Switching from kind to AWS provider
Since `kind` provider uses `/ets/hosts` for routing you **must remove** all entries added previously if you are using the same namespace you used in `kind`. Otherwise traffic will be incorrectly redirected to localhost instead of AWS.

It is thus advised to change namespace names, when switching providers.

---

## CRIB Limitations & Considerations

### Gateway DON
- Must always be on a **dedicated node**.
- Identified using `DON_TYPE=gateway`.
- No bootstrap node required, but multiple worker nodes are allowed.

### Capabilities Binaries
- Unlike Docker, k8s does not support copying capability binaries into a running container.
- Use a pre-built Docker image containing the necessary capabilities.

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

The test configuration is defined in a TOML file (e.g. `environment-*.toml`), which specifies properties for Chainlink nodes. The `capabilities` property of the `node_specs.node` determines which binaries are copied to the container:

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
- **Number of nodes**

### DON Type

Three types of DONs are supported:
- `workflow`
- `capabilities`
- `gateway`

There should only be **one** `workflow` and `gateway` DON, but multiple `capabilities` DONs can be defined.

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

For a working example of a multi-DON setup, refer to the [`environment-capabilities-don.toml`](environment-capabilities-don.toml) file.

### Number of nodes
When defining number of nodes you need not only to modify the `nodes` key, but also too add **nodespecs** for each node. In other words, the number of `nodespecs` needs to be equal to `nodes` count.

This is not enough:
```toml
[[nodesets]]
  nodes = 5
  override_mode = "each"

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

If there are `5` nodes you need to repeat this nodespec `5` times:
```toml
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

Also, `override_mode = "all"` is currently not supported.

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