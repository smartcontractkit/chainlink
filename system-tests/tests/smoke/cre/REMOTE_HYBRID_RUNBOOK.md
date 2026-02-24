# CRE Remote Hybrid Runbook

This runbook covers the EC2-based remote mode for CRE where components can run either locally or remotely.

## Scope

- Remote backend is EC2 + Docker (no Kubernetes path).
- Remote control plane is the CRE agent.
- Access mode is direct-only.

## Core Environment Variables

- `CRE_EC2_INSTANCE_ID=<instance-id>` (used by direct mode auto IP lookup)
- `CRE_EC2_AGENT_PORT=<port>` (defaults to `8080`)
- `CRE_EC2_AGENT_URL=<url>` (optional explicit override)
- `CRE_EC2_HOST_IP=<private-ip>` (optional in direct mode; if missing, resolved from AWS CLI using instance ID)
- `CRE_AWS_PROFILE=<profile>` (optional AWS auth profile)

## Direct Mode Defaults and IP Resolution

- Host IP resolution is:
  1. `CRE_EC2_HOST_IP` if set.
  2. Otherwise, resolve from AWS CLI using `CRE_EC2_INSTANCE_ID`:
     - `aws ec2 describe-instances --instance-ids <id> --query ...`
     - prefers private IP; falls back to public IP if needed.
- Region defaults to `us-west-2` unless AWS env region overrides are present.
- If no explicit host IP and no instance ID are available, startup fails with a clear error.

## AWS Credentials Resolution (CLI)

For direct-mode auto IP lookup, AWS CLI auth selection follows:

1. Static env credentials (`AWS_ACCESS_KEY_ID` + `AWS_SECRET_ACCESS_KEY`)
2. Web identity (`AWS_WEB_IDENTITY_TOKEN_FILE` + `AWS_ROLE_ARN`)
3. `CRE_AWS_PROFILE`
4. `AWS_PROFILE`
5. `AWS_DEFAULT_PROFILE`
6. AWS CLI default credential chain/profile

## Agent Startup

- In `direct` mode, bind agent to all interfaces (for example `0.0.0.0:18080`).

## Placement Rules

- Use `placement = "local" | "remote"` in CRE component config (NodeSets, JD, Blockchains).
- Same placement (`local->local`, `remote->remote`) uses **internal** URLs.
- Cross placement (`local->remote`, `remote->local`) uses **external** URLs.
- Remote NodeSets targeting local gateway are allowed when bridge/tunnel plumbing for gateway ingress is present.

## P2P Peering Rules (SharedPeering)

- `P2P.V2.ListenAddresses` is the **bind** interface used by the node process (CRE sets `0.0.0.0:5001`).
- `P2P.V2.AnnounceAddresses` is the **routable** endpoint set peers learn via discovery.
- In mixed placement:
  - local node announce set includes internal node host and a bridged host (`host.docker.internal:5001`) for remote callers.
  - remote node announce set includes internal node host and direct EC2 host IP address (`<ec2-private-ip>:5001`) in direct mode.
- If announce addresses are not routable from the caller placement, DON2DON discovery can succeed but stream establishment will fail.

## Mixed Bootstrap Reachability

- When remote DONs and a local bootstrap node are both present, CRE starts persistent relay plumbing for bootstrap peering on `5001`.
- Before DON startup, CRE performs a fail-fast sanity check that the remote relay listener for bootstrap peering is reachable.
- If startup fails on bootstrap reachability:
  - ensure relay supervisor was started,
  - ensure EC2 agent is reachable and has relay open for `5001`,
  - verify direct mode host IP resolution (`CRE_EC2_HOST_IP` or `CRE_EC2_INSTANCE_ID` + AWS CLI auth).

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
- Direct mode cannot resolve EC2 IP: ensure `CRE_EC2_INSTANCE_ID` is set and AWS CLI credentials are valid, or set `CRE_EC2_HOST_IP` explicitly.
- `invalid jd placement`: use `placement=local` or `placement=remote` (only supported values).
- Remote nodes hitting local-only fixtures: ensure fixture relay helper is active.
- Mixed remote->local gateway from NodeSets is supported when bridge plumbing is present.
- DON2DON flakiness in mixed mode: check generated node TOML includes `P2P.V2.AnnounceAddresses` that are routable from the opposite placement.
