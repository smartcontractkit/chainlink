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
- `--with-beholder`
- `--with-dashboards`
- `--with-observability`
- `--with-billing`
- `--with-contracts-version`
- `--setup-config`
- `--grpc-port`

Stop:

```bash
go run . env stop
go run . env stop --all
go run . env remote stop
go run . env stop-all
```

Restart:

```bash
go run . env restart
go run . env restart --with-beholder
```

Purge environment state:

```bash
go run . env state purge
```

Use purge when the saved state or cached environment artifacts look inconsistent.

## Mixed and Remote Execution

When at least one component uses `placement = "remote"`, Local CRE runs in mixed mode (local + remote). In this mode:

- startup prints a `Runtime Placement Matrix` summary
- `env stop` stops local resources only
- `env remote stop` stops remote resources only via the remote agent
- `env stop-all` runs remote stop first, then local stop

Agent URL resolution precedence is:

1. `CRE_REMOTE_AGENT_URL`
2. `CRE_REMOTE_HOST_IP` + `CRE_REMOTE_AGENT_PORT`
3. `CRE_REMOTE_AGENT_EC2_INSTANCE_ID` + `CRE_REMOTE_AGENT_PORT` + AWS profile/credentials
4. default `CRE_REMOTE_AGENT_PORT=18080` when omitted

If `env stop` warns that remote components are still running, run `go run . env remote stop`.

Implementation and ownership boundaries are documented in [Remote Execution Architecture](../reference/remote-execution-architecture.md).

## Setup and Images

By default Local CRE builds the Chainlink image from the local branch. To use a pre-built image instead, set `image` in each node definition in the topology TOML and omit `docker_ctx` and `docker_file`.

The deprecated `-p/--with-plugins-docker-image` flag still exists, but contributors should use TOML-based image selection instead.

## Beholder and Observability

Use `--with-beholder` when you need the ChIP ingress stack and Red Panda. Use `--with-observability` or `--with-dashboards` when you need the Grafana-based observability stack.

Important related flags:

- `--grpc-port` for ChIP ingress
- `--with-dashboards` to provision the dashboards on top of observability

When dashboards are enabled, the CLI waits for Grafana at `http://localhost:3000`.

## Storage and State

Local CRE persists state to the repo-local state file that the system tests later reuse. This is why the smoke-test helpers can detect an existing environment and avoid recreating it.

If you need to reset the environment completely, use state purge and then re-run setup/start.

## Debugging

For day-to-day debugging, the main patterns remain:

- inspect core node logs and container state
- rebuild or swap capabilities
- use observability and Beholder when tracing workflow activity

Hot-swapping guidance and workflow-specific commands are covered in:

- [Workflow Operations](workflows.md)
- [Advanced Topics](advanced.md)

## Telemetry and Tracing

The Local CRE stack supports:

- OTel-based observability
- Chip ingress / Beholder integration
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
4. restart with observability or Beholder if you need more signals

For topology-specific issues, continue with [Topologies and Capabilities](topologies.md).
