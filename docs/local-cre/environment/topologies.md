---
id: local-cre-environment-topologies
title: Topologies and Capabilities
sidebar_label: Topologies
sidebar_position: 2
---

# Topologies and Capabilities

Topologies control the DON layout, chains, infra type, and capability placement for Local CRE.

## Default Behavior

The CLI shell defaults `CTF_CONFIGS` to `configs/workflow-gateway-don.toml`. During environment startup, Local CRE also ensures the default capabilities configuration is prepended.

For smoke tests, the default local flow uses `configs/workflow-gateway-capabilities-don.toml` unless you override it.

## Discovering Topologies

Use the topology commands:

```bash
go run . topology list
go run . topology show --config configs/workflow-gateway-capabilities-don.toml
go run . topology generate
```

Implementation-backed defaults:

- `topology show --config` defaults to `configs/workflow-gateway-don.toml`
- `topology show --output-dir` defaults to `state`
- `topology generate --output-dir` defaults to `docs/topologies`
- `topology generate --index-path` defaults to `docs/TOPOLOGIES.md`

## Generated Topology Docs

Generated topology docs already live in the environment package:

- `core/scripts/cre/environment/docs/TOPOLOGIES.md`
- `core/scripts/cre/environment/docs/topologies/*.md`

Use them for:

- topology class
- DON count
- capability placement
- per-DON node counts and chain assignments

In particular, the generated matrix for `workflow-gateway-capabilities-don.toml` is the most useful reference for the default local smoke-test topology.

## Multiple DONs

Use a multi-DON topology when the workflow stack needs responsibilities split across separate DONs instead of running everything in one place.

In practice, the common DON roles are:

- `workflow` for workflow execution
- `capabilities` for capabilities exposed to other DONs
- `bootstrap` for DON bootstrapping
- `gateway` for connector and gateway traffic

When you inspect or change a multi-DON topology, verify these questions:

- which DON should execute the workflow
- which capabilities must stay local to that DON
- which capabilities must be remotely exposed from a separate DON
- which chains and ports those DONs need to reach

Keep these mental rules:

- workflow-only topologies are simpler for quick local iteration
- capability-enabled topologies are the standard path for realistic smoke coverage
- sharded and specialized topologies should be chosen only when the test or feature needs them

## Enabling Existing Capabilities

Enabling a capability is not one step. All of the following need to line up:

1. the topology TOML must place that capability on the correct DON
2. the node image used by that DON must actually contain the plugin or binary
3. if the workflow needs outbound access through the gateway to services outside the default setup, the environment must allow the required gateway egress ports
4. the generated topology docs should confirm that the final placement matches what you intended

As a rule of thumb:

- use a simpler topology when the workflow only needs local capabilities such as cron or consensus
- use the capability-enabled topology when the workflow needs remotely exposed capabilities such as EVM, read-contract, vault, or web API targets

After changing capability placement, re-run:

```bash
go run . topology show --config <your-topology>.toml
go run . topology generate
```

Then check the generated capability matrix before starting the environment.

## Adding or Modifying Topologies

When introducing a new topology:

1. add or update the TOML config under `configs/`
2. regenerate the topology docs with `go run . topology generate`
3. use `go run . topology show --config <file>` to sanity-check the result
4. update test guidance if the new topology is intended for smoke coverage

## Related Pages

- [Environment](index.md)
- [Advanced Topics](advanced.md)
- [Running Tests](../system-tests/running-tests.md)
