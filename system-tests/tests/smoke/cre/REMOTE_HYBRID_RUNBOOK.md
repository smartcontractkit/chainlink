# CRE Remote Hybrid Runbook

This runbook covers the EC2-based remote mode for CRE where components can run either locally or remotely.

## Scope

- Remote backend is EC2 + Docker (no Kubernetes path).
- Remote control plane is the CRE agent.
- Access modes:
  - `ssm`: control and endpoint reachability via SSM tunnels.
  - `direct`: endpoint reachability via EC2 host IP, with SSM optional for agent only.

## Core Environment Variables

- `CRE_AGENT_MODE=ec2`
- `CRE_EC2_INSTANCE_ID=<instance-id>` (required for SSM mode)
- `CRE_EC2_AGENT_PORT=<port>` (defaults to `8080`)
- `CRE_EC2_AGENT_URL=<url>` (optional explicit override)
- `CRE_REMOTE_ACCESS_MODE=ssm|direct`
- `CRE_EC2_HOST_IP=<private-ip>` (required when `CRE_REMOTE_ACCESS_MODE=direct`)
- `CRE_AWS_PROFILE=<profile>` (optional SSM auth profile)

## Agent Startup

- In `ssm` mode, bind agent to loopback (for example `127.0.0.1:18080`).
- In `direct` mode, bind agent to all interfaces (for example `0.0.0.0:18080`).

## Placement Rules

- Same placement (`local->local`, `remote->remote`) uses **internal** URLs.
- Cross placement (`local->remote`, `remote->local`) uses **external** URLs.
- Remote NodeSets targeting local gateway are allowed when bridge/tunnel plumbing for gateway ingress is present.

## Bridge and Fixture Relay

- Remote components cannot directly call local in-process fixtures.
- Use fixture relay for local fixtures (CHiP testsink, fake HTTP, billing/PoR mocks).
- Relay is opened per fixture port and uses fixed remote port parity.

## Recommended Test Order

1. All remote.
2. All local.
3. Mixed (for example JD local + NodeSet remote).

## Fast Triage Checklist

- Agent unreachable: verify bind address/port vs chosen access mode.
- `invalid jd target: local`: use `target=local` (supported; `docker` is alias).
- Remote nodes hitting local-only fixtures: ensure fixture relay helper is active.
- Mixed remote->local gateway case: expected failure for now.
