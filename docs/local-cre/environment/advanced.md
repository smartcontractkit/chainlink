---
id: local-cre-environment-advanced
title: Advanced Topics
sidebar_label: Advanced
sidebar_position: 3
---

# Advanced Topics

This page groups the advanced Local CRE topics in one place.

## Billing

Local CRE can bootstrap the billing service with `--with-billing`, and the smoke package also includes billing coverage. Use this when validating workflow billing integrations locally rather than only running the base DON stack.

## Adding a New Standard Capability

New Local CRE extension work should be expressed as a feature, not through the older installable-capability path.

At a high level, adding a feature means:

1. choose the capability flag the feature is responsible for
2. add a feature implementation under `system-tests/lib/cre/features/...`
3. implement the feature lifecycle hooks such as `PreEnvStartup` and `PostEnvStartup`
4. register any required on-chain capability configuration from the feature hook
5. add any node configuration changes needed for that feature
6. include the feature in the default feature set
7. expose the corresponding flag through one or more topology configs

In the current implementation:

- `InstallableCapability` is deprecated for new work
- `Feature` is the primary interface for new capabilities and related setup
- features are applied before environment startup and can register capability config on-chain

This is contributor-facing material. If you are only choosing or running topologies, use [Topologies and Capabilities](topologies.md) instead.

## Hot Swapping

Local CRE supports swapping:

- the Chainlink node image
- a capability binary

Use hot swapping when you want to refresh part of the running environment without doing a full `env stop` and `env start`.

### Swapping Chainlink nodes

To recreate the Chainlink node containers with updated images or rebuilt local code:

```bash
go run . env swap nodes
```

This command:

- reloads the saved Local CRE state
- removes the current node containers for each DON
- recreates those node sets
- rebuilds the Docker image if the topology is configured to build from source and the local source has changed

Useful flags:

- `--force` to force container removal
- `--wait-time` to control how long Local CRE waits before retrying after removal/startup issues

Use `env swap nodes` when:

- you changed Chainlink core code
- you changed the Docker build inputs for node images
- you want the running DONs to pick up a newly built node image

### Swapping a capability binary

The capability swap flow is still the fastest local loop when you need to rebuild a plugin without recreating the whole environment.

Typical command:

```bash
go run . env swap capability -n cron -b /path/to/cron
```

`env swap capability` supports:

- `--name`
- `--binary`
- `--force`

This command:

- finds DONs that use the named capability flag
- cancels matching job proposals
- copies the new binary into the relevant containers
- re-approves the proposals so the jobs restart with the updated binary

Use `env swap capability` when:

- you rebuilt a swappable capability locally
- you only need to refresh one capability instead of the entire node image

Only capabilities supported by the swappable capability provider can be hot-reloaded this way.

## DX Tracing and Telemetry

Use this part of the stack when you need to answer questions like:

- did the workflow engine initialize successfully
- did a workflow emit the expected user logs or base messages
- are Beholder messages reaching Kafka and Loki
- are node and connector logs reaching Grafana/Loki
- do I need more than plain container logs to debug a failing run

### What each option gives you

- `--with-beholder`
  Starts the Beholder stack used by the CRE tests for workflow-related messages and heartbeat validation.
- `--with-observability`
  Starts the observability stack.
- `--with-dashboards`
  Starts observability and provisions the Grafana dashboards used for inspection.
- `go run . obs up`
  Manages the observability stack directly from the Local CRE CLI.

### When to use which

Use plain container logs first when:

- startup fails early
- the problem is isolated to one container
- you do not need workflow-event correlation

Use Beholder when:

- the test expects specific workflow log messages
- you need to validate workflow heartbeats or emitted events
- you are debugging cron, HTTP action, DON time, or other scenarios that already rely on Beholder listeners in the test helpers

Use observability and dashboards when:

- you need aggregated logs in Grafana/Loki
- you want to inspect log streaming end to end
- you need a broader view across nodes, connectors, and supporting services

### DX tracing

In Local CRE, DX refers to usage tracking for the Local CRE tooling itself, not workflow/node telemetry.

The CLI records events such as:

- environment setup result
- environment startup result and startup time
- workflow deployment
- Beholder startup
- billing startup
- capability swaps
- node swaps

The tracker configuration in the Local CRE code uses:

- GitHub variable name: `API_TOKEN_LOCAL_CRE`
- product name: `local_cre`

This is separate from observability, Loki, Grafana, Beholder logs, or any workflow-level tracing.

If you are debugging Local CRE usage instrumentation, look at the tracking hooks in the environment commands. If you are debugging workflow execution, logs, or message flow, use the observability and Beholder paths described above instead.

### Practical debugging flow

1. start the environment normally
2. add `--with-beholder` if the scenario depends on workflow messages
3. add `--with-observability` or `--with-dashboards` when you need Grafana/Loki
4. rerun the workflow or smoke test
5. inspect Beholder output, Grafana, and Loki before changing topology or code

## Using Specific Images and Existing Keys

Use this path when you need reproducible node images or stable node identity instead of whatever Local CRE generates from the working tree.

### Using a specific Chainlink image

By default, most Docker-based topologies build the Chainlink node image from the local checkout:

```toml
[nodesets.node_specs.node]
  docker_ctx = "../../../.."
  docker_file = "core/chainlink.Dockerfile"
```

To pin a prebuilt image instead, replace the build settings with an explicit image:

```toml
[nodesets.node_specs.node]
  image = "chainlink-tmp:latest"
```

Apply that change to every node spec in the nodeset that should use the pinned image.

Use explicit images when:

- you want repeatable runs across machines
- you want to test a nightly or CI-built image
- you are comparing behavior between two image versions
- you are running in Kubernetes, where runtime builds are not the normal path

The example override file at `core/scripts/cre/environment/configs/examples/workflow-don-overrides.toml` shows this pattern in practice.

### Reusing existing EVM and P2P keys

Local CRE normally generates fresh node keys. The lower-level CRE types also support importing an existing node-secrets payload instead of generating new keys.

That path is useful when you need:

- stable peer IDs across restarts
- stable EVM addresses for repeatable tests
- a previously generated node identity restored into a new environment

The implementation detail to be aware of is:

- `NodeKeyInput.ImportedSecrets` bypasses key generation and imports existing encrypted node secrets

This is a contributor or integrator workflow rather than a normal Local CRE quickstart path. If you need stable keys, treat that as a deliberate topology/configuration change and validate the resulting peer IDs and on-chain addresses after startup.

## External and Public Blockchains

The default Local CRE flow uses local Anvil chains. Switch to external or public blockchains only when the workflow or capability truly depends on a non-local RPC.

### What actually changes

When you stop using only local Anvil chains, you need all of the following to line up:

- the `[[blockchains]]` entries must use the correct chain type and chain ID
- the topology must support that chain on the relevant DONs
- the node image and enabled features must support the chain family you are targeting
- the chain RPC endpoints must be reachable and healthy before Local CRE startup

### Example: use a preconfigured external EVM RPC

For a non-local chain, the practical pattern is to provide a blockchain entry whose output URLs are already known, instead of asking Local CRE to spin up a local chain container.

Example:

```toml
[[blockchains]]
  chain_id = "11155111"
  type = "anvil"

  [blockchains.out]
    type = "anvil"
    use_cache = true

    [[blockchains.out.nodes]]
      ws_url = "wss://0xrpc.io/sep"
      http_url = "https://0xrpc.io/sep"
      internal_ws_url = "wss://0xrpc.io/sep"
      internal_http_url = "https://0xrpc.io/sep"
```

Then make sure the DON that needs that chain supports it. For example:

```toml
[[nodesets]]
  name = "bootstrap-gateway"
  supported_evm_chains = [1337, 11155111]
```

And if a workflow DON needs chain-specific capabilities on that chain, its capability list must include the matching flag, for example:

```toml
capabilities = ["read-contract-11155111", "write-evm-11155111"]
```

### When to use this pattern

Use it when:

- a workflow must read from or write to a public testnet or mainnet
- a capability depends on chain-specific behavior that Anvil does not reproduce
- you are validating against an already-running external chain endpoint

### Practical checks before startup

Before you run `env start`, verify:

1. the RPC URLs respond and match the expected chain ID
2. the relevant DONs support that chain
3. the capability flags include the chain-qualified capability names you need
4. the node image contains any required plugins for that chain family

Compared with the default Anvil flow, this setup is much more sensitive to RPC health, endpoint latency, and mismatches between topology config and the actual remote chain.

## Kubernetes Deployment

Kubernetes is the alternative infra mode to the default Docker setup.

Switch the topology to Kubernetes by setting:

```toml
[infra]
  type = "kubernetes"
```

Use Kubernetes when:

- you need cluster-like execution rather than a single local Docker network
- your test environment depends on prebuilt images instead of local source builds
- you are validating behavior closer to remote or CI-style execution

Unlike Docker mode, Kubernetes mode assumes the nodes are already running in the cluster and Local CRE connects to them by generating the expected service URLs.

### Prerequisites

Before using Kubernetes mode, make sure you have:

1. a Kubernetes cluster with the Chainlink nodes already deployed
2. Helm charts or deployment overlays that support config and secret overrides
3. external ingress configured if you need external access
4. local `kubectl` access to the cluster
5. all DON nodes deployed in a single namespace

### Kubernetes configuration

The Kubernetes fields live under `infra.kubernetes`:

```toml
[infra]
  type = "kubernetes"

  [infra.kubernetes]
    namespace = "my-namespace"
    external_domain = "example.com"
    external_port = 80
    label_selector = "app=chainlink"
    node_api_user = "admin@chain.link"
    node_api_password = "secure-password-here"
```

What these fields are for:

- `namespace`
  The namespace where the DON nodes are already running.
- `external_domain`
  The domain used to derive externally reachable service URLs.
- `external_port`
  The ingress port, usually `80`.
- `label_selector`
  The selector used to discover the relevant Chainlink pods.
- `node_api_user` and `node_api_password`
  The credentials Local CRE uses to talk to the nodes.

### Main differences from Docker

In Docker mode, many topologies use:

```toml
docker_ctx = "../../../.."
docker_file = "core/chainlink.Dockerfile"
```

In Kubernetes mode, prefer explicit images instead:

```toml
[nodesets.node_specs.node]
  image = "chainlink:your-tag"

[jd]
  image = "job-distributor:your-tag"
```

Kubernetes is therefore the wrong choice for fast local code iteration and the right choice for image-based validation.

### Config and secret overrides

Kubernetes mode is designed to work with deployments that accept node-specific config overlays.

The expected model is:

1. Local CRE generates node-specific TOML config from the topology
2. that config is pushed into the cluster as ConfigMaps and Secrets
3. the running nodes mount those overlays through the deployment or Helm chart
4. updated configs are picked up without rebuilding the images

The original Local CRE guidance called out the expected objects:

- a ConfigMap such as `<node-name>-config-override`
- a Secret such as `<node-name>-secrets-override`

If your chart or deployment setup does not support that overlay pattern, Kubernetes mode will not behave like the standard Local CRE flow.

### Example configuration

Representative Kubernetes-connected topology:

```toml
[[blockchains]]
  chain_id = "1337"
  type = "anvil"

  [blockchains.out]
    use_cache = true
    type = "anvil"
    family = "evm"
    chain_id = "1337"

    [[blockchains.out.nodes]]
      ws_url = "wss://anvil-service-rpc.example.com"
      http_url = "https://anvil-service-rpc.example.com"
      internal_ws_url = "ws://anvil-service:8545"
      internal_http_url = "http://anvil-service:8545"

[infra]
  type = "kubernetes"

  [infra.kubernetes]
    namespace = "my-namespace"
    external_domain = "example.com"
    external_port = 80
    label_selector = "app=chainlink"
    node_api_user = "admin@chain.link"
    node_api_password = "secure-password-here"

[jd]
  csa_encryption_key = "d1093c0060d50a3c89c189b2e485da5a3ce57f3dcb38ab7e2c0d5f0bb2314a44"
  image = "job-distributor:your-tag"
```

### Service URL generation

Local CRE derives service URLs in Kubernetes from naming conventions. Using `my-namespace` as an example:

- bootstrap node internal URL: `http://workflow-bt-0:6688`
- bootstrap node external URL: `https://my-namespace-workflow-bt-0.example.com`
- plugin node internal URL: `http://workflow-0:6688`
- plugin node external URL: `https://my-namespace-workflow-0.example.com`

This is why the namespace, external domain, and label selector matter so much in Kubernetes mode.

### Practical guidance

Before using Kubernetes:

1. replace local build settings with explicit image tags
2. set the Kubernetes-specific infra fields required by your environment, such as namespace
3. confirm the referenced images already contain the plugins and binaries your topology needs
4. provide blockchain output URLs for chains that are already running remotely
5. reuse the same topology validation flow with `topology show` and generated topology docs
