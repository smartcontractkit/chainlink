# reconciler

A CLI tool that reconciles Chainlink node configurations for nodes deployed to Kubernetes via Griddle.

> This tool configures already-deployed CRE nodes (contracts, capabilities,
> jobs, TOML); it sits on top of Griddle rather than being part of it. Griddle itself (`griddle.yaml`, chart
> deployment) is unaffected by this rename.

> **`dev` environment only.** This tool exists to quickly stand up a CRE DON in `dev` without going through the
> full GitOps flow. Do not use it against higher environments (`stage`, `prod`, etc.) — those must continue to
> go through GitOps.

## What it does

After Griddle deploys Chainlink nodes and databases to Kubernetes, the nodes are running but unconfigured — no contracts deployed, no capabilities registered, no jobs created. This tool fills that gap:

1. Deploys Keystone smart contracts (Capabilities Registry, Workflow Registry, Forwarders) to the registry chain.
2. Registers worker nodes, capabilities, and DONs on-chain in the Capabilities Registry.
3. Configures the Workflow Registry (allowed DON IDs, workflow owners).
4. Injects contract addresses + capability config into node TOML (one patch, one re-roll).
5. Creates and approves jobs via the Job Distributor (capability jobs + gateway jobs).

In the future it will be **reconciliation-based**: run it repeatedly, and it only applies the delta between desired and actual state. Run it twice with no changes, and the second run is a no-op. Currently it requires you to manually rewing the phase in the state file to apply what comes after it.

## Quick start

```bash
# From the chainlink repo root:
cd core/scripts/cre/reconciler

# Build
go build -o bin/reconciler ./cmd

# Create a desired-state TOML (see example below)
cat > cre/desired.toml << 'EOF'
[infra]
  type = "griddle"
  chart_values = "deploy/config/my-repo"
  namespace = "my-repo-nodeset"

[jd]
  grpc = "grpc-job-distributor.main.stage.cldev.sh:443"
  domain = "cre"
  environment = "dev"

[[chains]]
  chain_id = 1337
  ws_url = "wss://anvil-1337.my-repo-nodeset.svc.cluster.local"
  http_url = "https://anvil-1337.my-repo-nodeset.svc.cluster.local"
  registry = true

[[dons]]
  name = "workflow"
  don_types = ["workflow"]
  capabilities = ["cron", "evm-1337"]

[capability_configs.cron]
  binary_name = "cron"
[capability_configs.evm]
  binary_name = "evm"
EOF

# Run reconciliation
./bin/reconciler apply --desired cre/desired.toml --state cre/state.toml

# If the tool writes a TOML patch, apply it via Griddle (.deploy dev),
# wait for the nodes to re-roll, then re-run:
./bin/reconciler apply --desired cre/desired.toml --state cre/state.toml
```

## Web UI

The tool includes a built-in web UI for visual DON composition and status monitoring. This is the recommended
way to create and edit your `desired.toml` — build the DON layout visually, then use the TOML Preview / Save
features to write it out, rather than hand-writing the TOML from scratch.

```bash
./bin/reconciler serve --desired cre/desired.toml --state cre/state.toml --chart-dir deploy/config/my-repo
```

Then open http://localhost:8089 in your browser.

### Features

- **DON Builder** — drag & drop nodes from a node pool into DON columns. Toggle capabilities per DON from a catalog; chain-scoped capabilities can only be attached to a chain declared in the Chains tab. Add/remove DONs, edit names, set DON types.
- **Chains** — always preloaded on startup from the chart nodes' existing typed EVM config (there's no manual "add chain" action; to add a chain, add its network config to the chart and refresh). If the same chain ID is discovered with more than one RPC URL (e.g. a gateway node pointed at a different one than workers), every variant is shown so you can delete the wrong one. Pick exactly one registry chain.
- **Status** — view the current reconcile phase, deployed contract addresses, on-chain DON IDs, and node status.
- **Config** — edit JD connection settings and infrastructure paths.
- **TOML Preview** — see the generated desired-state TOML before saving.
- **Save** — writes the desired-state TOML directly from the UI.

The web UI is self-contained in the binary (embedded files, no external frontend build step). It uses [SortableJS](https://sortablejs.github.io/Sortable/) for drag & drop and [Tailwind CSS](https://tailwindcss.com/) via CDN.

## Desired-state TOML

The desired state is declared in a TOML file in your consumer repo. The easiest way to produce this file is via
the [web UI](#web-ui) above — use its DON Builder and TOML Preview / Save features rather than hand-writing it.
The format below is documented for reference and for making manual edits once the file exists:

```toml
[infra]
  type = "griddle"
  chart_values = "deploy/config/my-repo"    # path to Griddle chart values dir
  namespace = "my-repo-nodeset"              # K8s namespace
  kubeconfig = "~/.kube/config"              # optional, defaults to KUBECONFIG env

[jd]
  grpc = "grpc-job-distributor.main.stage.cldev.sh:443"
  domain = "cre"
  environment = "dev"

# Chains — every EVM chain a capability or the on-chain registry needs, with
# exactly one entry marked as the registry chain (where CapabilitiesRegistry /
# WorkflowRegistry are deployed and where nodes register). This is the only
# source of chain data the reconciler uses — nothing is inferred from the
# chart's anvil.instances, and the tool never writes [[EVM]] node config: the
# chart must already declare the network (and thus generate a node key) for
# any chain a capability references, or `apply` fails fast with an actionable
# error instead of silently skipping the affected work.
[[chains]]
  chain_id = 1337
  ws_url = "wss://anvil-1337.my-repo-nodeset.svc.cluster.local"
  http_url = "https://anvil-1337.my-repo-nodeset.svc.cluster.local"
  registry = true

[[chains]]
  chain_id = 11155111
  ws_url = "wss://sepolia.infura.io/ws/v3/..."
  http_url = "https://sepolia.infura.io/v3/..."

# DONs — each [[dons]] entry's name must equal a Griddle chart node-set's
# don-name label (chainlink-node.registerNodes.labels.don-name); membership
# is always that node-set's nodes, resolved from the chart at run time.
[[dons]]
  name = "workflow"
  don_types = ["workflow"]
  capabilities = ["cron", "http-action", "http-trigger", "consensus", "evm-1337"]
  # bootstrap_node defaults to the node with nodeType: boot in the chart

[[dons]]
  name = "capabilities"
  don_types = ["capabilities"]
  exposes_remote_capabilities = true
  capabilities = ["vault", "evm-11155111"]

# Optional: assign gateway nodes to specific DONs
# [[gateway_nodes]]
#   node = "node-gw-0"
#   don = "workflow"

# Capability configs — same format as local CRE
[capability_configs.evm]
  binary_name = "evm"
[capability_configs.evm.values]
  LogTriggerPollInterval = 1500000000
  ReceiverGasMinimum = 500

[capability_configs.cron]
  binary_name = "cron"
```

### JD access token

Export `GRIDDLE_JD_ACCESS_TOKEN` before running `apply`:

```bash
export GRIDDLE_JD_ACCESS_TOKEN="..."
./bin/reconciler apply --desired cre/desired.toml --state cre/state.toml
```

The token is a live bearer credential and is never stored in `desired.toml`, `state.toml`, or the web UI — it is read
only from this environment variable.

`apply` checks that `GRIDDLE_JD_ACCESS_TOKEN` is set (whenever `jd.grpc` is configured) **before doing any other
work** — including Kubernetes discovery — and fails immediately with a clear error if it's missing, rather than
failing later with a confusing "JD client is required" error after partial progress. `serve` does not require the
token to start, so you can start it first and export the token afterward — but the web UI's "Check JD Connectivity"
button reads the token from the server's environment at click-time and reports immediately if it's missing, instead
of attempting (and failing) a real connection to JD.

### Node roles

Node roles are read from the Griddle chart values (`nodeType` field per instance):
- `standard` (default) — oracle/plugin nodes that form DONs and run capabilities.
- `boot` — bootstrap nodes, one per DON, used as P2P bootstrappers.
- `gateway` — gateway nodes that route capability calls (vault, http-action, http-trigger). Not registered on-chain as signatory nodes.

If `nodeType` is not set in the chart, bootstrap is inferred by naming convention (`*-bt-*` or `*bootstrap*`).

DON membership is derived from the chart's `chainlink-node.registerNodes.labels.don-name` label. A `[[dons]]` entry's `name` must equal a chart node-set's `don-name`; its member nodes are exactly that node-set. A `nodes` key is no longer accepted — desired.toml never lists individual node names.

### Gateway nodes

Gateway nodes are detected from `nodeType: gateway` in the chart — they are discovered automatically, not listed under any DON. The reconciler:

1. Generates gateway node TOML (OCR2, P2P, ExternalRegistry, stub GatewayConnector).
2. Wires worker nodes' `[Capabilities.GatewayConnector]` with `Gateways[]` pointing at the gateway node's WebSocket endpoint.
3. Creates a gateway job (type `"gateway"`) via JD with service handler configs and sharded DON topology.

All TOML changes (gateway + worker) are in the single breakpoint patch.

## Reconcile flow

```
DISCOVER     Read node info from JD (P2P IDs, OCR keys)
             Read node endpoints from K8s (API URLs, credentials)
             Read chain info + node roles from chart values

P1  DEPLOY    Deploy Keystone contracts
P4  CONFIG    Register worker nodes + capabilities + DONs on-chain
P5  RESOLVE   Read actual DON IDs from the contract
P6  CONFIG    Configure Workflow Registry
P9  VAULT     Update vault capability config (if gateway DON)

T1  INJECT    [BREAKPOINT] Write all node TOML into chart values → EXIT
              User applies via Griddle (.deploy dev), re-runs tool.

R1  VERIFY    Confirm nodes re-rolled with new TOML
J1  JOBS      Propose + approve all JD jobs (capability + gateway)
V1  DONE      Verify convergence
```

The breakpoint at T1 is the only place the tool stops. It writes a `30-cre` TOML layer into your chart values YAML, then exits with code 42. After applying the patch via Griddle and re-running, the tool resumes from where it left off.

## State file

The state file (`cre/state.toml`) is committed to the consumer repo. It records:
- Deployed contract addresses
- Resolved DON IDs
- JD node IDs
- Per-node runtime info (PeerID, API URL, CSA key)
- Current reconcile phase
- Hash of the desired-state TOML (to detect changes)

This will enablee idempotent runs — the tool checks actual state before each phase and skips if already done.

## Audit artifacts

In the future per-changeset artifacts will written in CLD-compatible format to `cre/artifacts/durable_pipelines/`. The directory layout and JSON schema match what the `cld` runtime produces.

## Roadmap

See [TODO.md](./TODO.md) for the groomed backlog of planned work (reconciliation model, feature/UX gaps, open
decisions) beyond what's described here.

## Relationship to other tools

- **Griddle** (`infra-griddle-app`, `griddle-stacks`): deploys nodes (and optionally Anvil chains) + databases to K8s. Does not know about contracts, capabilities, or jobs. The reconciler never reads Griddle's `anvil.instances` — chains are declared explicitly in `desired.toml` (see `[[chains]]` above).
- **Local CRE** (`core/scripts/cre/environment`): same concept but for Docker-based environments. This tool reuses the local CRE's changeset library but not its orchestration code.
- **JD** (Job Distributor): shared service. The reconciler talks to it via gRPC to list nodes, check existing jobs, and propose new ones.
- **`fund-nodes`**: a separate script (not part of this tool) that tops up node EVM accounts. The reconciler assumes nodes are funded and only warns if it detects a zero balance.

## Building

```bash
cd core/scripts/cre/reconciler
go build -o bin/reconciler ./cmd
```

## Testing

```bash
go test ./cre/reconciler/... -v
```
