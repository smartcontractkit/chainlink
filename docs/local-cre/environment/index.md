---
id: local-cre-environment-index
title: Environment Operations
sidebar_label: Environment
sidebar_position: 0
---

# Environment Operations

The Local CRE CLI lives in `core/scripts/cre/environment`. The default invocation style is:

```bash
cd core/scripts/cre/environment
go run . <command> [subcommand]
```

You can also install the binary:

```bash
make install
```

That produces `local_cre`, including the interactive shell (`local_cre sh`).

## Environment Lifecycle

Start:

```bash
go run . env start [--auto-setup]
```

Current `env start` flags from source code include:

- `--auto-setup`
- `--wait-on-error-timeout`
- `--cleanup-on-error`
- `--extra-allowed-gateway-ports`
- `--with-example`
- `--example-workflow-timeout`
- `--with-chip-ingress-stack` (deprecated alias: `--with-beholder`)
- `--with-dashboards`
- `--with-observability`
- `--with-billing`
- `--setup-config`
- `--grpc-port`
- `--local-node`
- `--local-capabilities`
- `--capabilities-path`
- `--local-build-platform`
- `--local-node-image`

## Building the Node and Capabilities from Local Source

By default the node image and CRE capability plugins come from the topology's
`docker_ctx` build (in-Docker) and the pinned upstream capability versions. When
you are actively developing the node or a capability, the following flags build
them from your **local working tree** instead — including local `replace`
directives (e.g. a local `chainlink-common`) that an in-Docker build cannot see.

- `--local-node` — cross-compile the Chainlink node from the local checkout on
  the host, bake it into a minimal image, and use it for every node (overrides
  the topology's `docker_ctx`/`image`).
- `--local-capabilities <names|all>` — build these CRE capabilities from local
  source and inject them into the nodes (via `capabilities` / `capabilities_container_dir`),
  overriding the versions baked into the node image. Unlisted capabilities keep
  the image's version. Accepts `all` or a comma-separated list, e.g.
  `cron,consensus,evm,http-action,http-trigger`.
- `--capabilities-path <dir>` — location of the local capabilities repo
  (default `$CRE_CAPABILITIES_PATH`, else `~/go/src/github.com/smartcontractkit/capabilities`).
- `--local-build-platform <os/arch>` — target for the local builds (default `linux/<host arch>`).
- `--local-node-image <tag>` — image tag for the locally built node (default `cre-node:local`).

Examples:

```bash
# Node + all capabilities from local source (with observability):
go run . env start --with-observability --local-node --local-capabilities all

# Full (pinned) node image, but override only cron + evm with local builds:
go run . env start --local-capabilities cron,evm
```

Prerequisites (one-time, not performed by these flags): a C cross-compiler for
the node's CGO build when the host OS/arch differs from the target (override with
`CRE_LOCAL_CC` / `CRE_LOCAL_CXX` / `CRE_LOCAL_LIBDIR`), and the managed images
present (`go run . env setup`). On Linux hosts building `linux/<same arch>` no
cross-compiler is needed.

Stop:

```bash
go run . env stop
go run . env stop --all
```

Restart:

```bash
go run . env restart
go run . env restart --with-chip-ingress-stack
```

Purge environment state:

```bash
go run . env state purge
```

Use purge when the saved state or cached environment artifacts look inconsistent.

## Setup and Images

By default Local CRE builds the Chainlink image from the local branch. To use a pre-built image instead, set `image` in each node definition in the topology TOML and omit `docker_ctx` and `docker_file`.

The deprecated `-p/--with-plugins-docker-image` flag still exists, but contributors should use TOML-based image selection instead.

`env setup` ensures required managed images are present, including Job Distributor, Chip Router, Chip Ingress, and Chip Config.

When pulling managed images from ECR, configure both registries:

- `MAIN_AWS_ECR` for core managed CRE images
- `SDLC_AWS_ECR` for Chip Router images

Chip Router image resolution during startup is:

1. `CTF_CHIP_ROUTER_IMAGE` (if set)
2. `chip_router.image` from the active topology TOML

If the resolved router image is missing locally, startup follows the same build-or-pull fallback path as the Chip Ingress stack images.

## Chip Ingress stack and Observability

Use `--with-chip-ingress-stack` when you need the ChIP ingress stack and Red Panda. (`--with-beholder` is deprecated.) Use `--with-observability` or `--with-dashboards` when you need the Grafana-based observability stack.

Important related flags:

- `--grpc-port` for ChIP ingress
- `--with-dashboards` to provision the dashboards on top of observability

When dashboards are enabled, the CLI waits for Grafana at `http://localhost:3000`.

### Chip Router Topology

Chip Router is the ingress owner on `50051`. Nodes emit workflow telemetry to the router, and the router fans out to downstream subscribers.

Current local ports:

- `50050`: Chip Router admin API
- `50051`: Chip Router ingress gRPC
- `50052`: chip-config
- `50053`: real ChIP / Chip Ingress stack gRPC

In tests, sink-backed scenarios register a test sink with Chip Router. Stack-backed scenarios register real ChIP ingress with Chip Router.

To override the router image without changing committed TOMLs:

```bash
export CTF_CHIP_ROUTER_IMAGE=chip-router:<commit-sha>
```

This override takes precedence over `chip_router.image`.

## Storage and State

Local CRE persists state to the repo-local state file that the system tests later reuse. This is why the smoke-test helpers can detect an existing environment and avoid recreating it.

If you need to reset the environment completely, use state purge and then re-run setup/start.

## Debugging

For day-to-day debugging, the main patterns remain:

- inspect core node logs and container state
- rebuild or swap capabilities
- use observability and the Chip Ingress stack when tracing workflow activity

Hot-swapping guidance and workflow-specific commands are covered in:

- [Workflow Operations](workflows.md)
- [Advanced Topics](advanced.md)

## Telemetry and Tracing

The Local CRE stack supports:

- OTel-based observability
- Chip Router fanout with Chip Ingress stack integration
- DX tracing

If you need the full tracing stack for debugging or demos, enable observability during startup and follow the environment-specific tracing configuration described in the advanced page.

## Troubleshooting

The most common failures in this area are:

- Chainlink node migrations fail
- Docker image not found
- Docker cannot download required public images
- `gh` is missing or unauthenticated

When startup problems happen:

1. rerun `go run . env setup`
2. confirm image access and authentication
3. purge state if the saved state is stale
4. restart with observability or the Chip Ingress stack if you need more signals

For topology-specific issues, continue with [Topologies and Capabilities](topologies.md).
