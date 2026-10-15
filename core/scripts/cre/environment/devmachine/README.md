# devmachine

`devmachine` provisions a remote EC2 **DevMachine** for hybrid CRE runs and
prepares it end to end: creates (or adopts) the machine, installs Docker, builds
and loads the required images, and starts the remote agent as a systemd service.
It then prints the commands to start the environment locally.

It is a subcommand of the Local CRE CLI:

```
go run . devmachine <up|down|status|handoff> [flags]
```

## What it does (phases)

| Phase | Description |
|-------|-------------|
| **P0** Preflight | Verify tools (`aws kubectl docker buildx go rsync ssh`), AWS creds, kubectl access + `can-i create devmachines`, auto-detect repos, verify the SSH key is loaded in `ssh-agent` |
| **P1** DevMachine | Create the `DevMachine` CRD and wait until `Ready`, or adopt an existing instance via `--instance-id` |
| **P2** Bootstrap | Over SSM: install Docker + rsync, enable sshd, install your public key (`restorecon`), add users to the `docker` group |
| **P3** Artifacts | Build the remote agent (linux/amd64) and the amd64 `chainlink`/`job-distributor` images; cache tarballs |
| **P4** Transfer | Write a managed `~/.ssh/config` block and `rsync` the agent + image tarballs to the machine |
| **P5** Load | `docker load` the tarballs and retag job-distributor to the name the env expects |
| **P6** Start agent | Install + start `remote-agent` as a systemd service (survives logout, restarts on failure) |
| **P7** Handoff | Print how to check the agent and start the environment locally |

## Prerequisites

- **Tools:** `aws` (v2), `kubectl`, `docker` + `buildx`, `go`, `rsync`, `ssh`.
- **AWS:** a valid SSO session for the chosen profile (`aws sso login --profile ...`).
- **Kubernetes:** a kubectl context that can reach the `DevMachine` CRD. If you
  cannot create DevMachines, ask someone to create one and use `--instance-id`.
- **SSH key:** passphrase handling is **not** implemented — the key must be usable
  non-interactively. P0 checks `ssh-agent` and, when interactive, runs
  `ssh-add <key>` for you; under `--yes` it fails with instructions.
- **Repos:** `chainlink` (for the agent + CL image) and `job-distributor`. These
  are auto-detected from the current directory / common paths, or via
  `CRE_CHAINLINK_REPO` / `CRE_JOB_DISTRIBUTOR_REPO`, or `--chainlink-repo` /
  `--job-distributor-repo`.

## Usage

Create a machine and prepare it (interactive TUI):

```
go run . devmachine up \
  --aws-profile dev-ex-dev/administrator-access \
  --team <team> --cost-center <cost-center>
```

Adopt an existing machine (no create permission):

```
go run . devmachine up --instance-id i-0123456789abcdef0 \
  --aws-profile dev-ex-dev/administrator-access
```

Non-interactive (CI / scripts):

```
go run . devmachine up --yes --aws-profile <p> --team <t> --cost-center <c>
```

Other commands:

```
go run . devmachine status          # config + live checks (EC2, ssh, agent)
go run . devmachine handoff         # reprint the next-steps commands
go run . devmachine down            # delete the machine (only if tool-created) + clean ssh block
go run . devmachine down --purge-cache
```

## Flags (`up`)

| Flag | Default | Notes |
|------|---------|-------|
| `--name` | `hybrid-local-cre-<user>` | DevMachine name (DNS-1123 sanitized) |
| `--namespace` | `cre-staggered-rollout-workflow` | |
| `--team`, `--cost-center` | — | required for create (labels) |
| `--instance-type` | `m8a.2xlarge` | |
| `--volume-size` | `200` | GiB |
| `--az` | `a` | availability zone suffix |
| `--instance-id` | — | adopt an existing instance instead of creating |
| `--region` | `us-west-2` | |
| `--aws-profile` | — | |
| `--chainlink-repo`, `--job-distributor-repo` | auto-detected | |
| `--ssh-key` | `~/.ssh/id_ed25519` | must be in `ssh-agent` |
| `--chainlink-image` | `chainlink-amd:latest` | |
| `--jd-image` | `job-distributor-amd:0.22.1` | |
| `--jd-tag` | `v0.22.1` | git tag built for job-distributor |
| `--env-workdir`, `--env-config` | chainlink env dir / `configs/workflow-gateway-don-mixed-placement.toml` | used in handoff |
| `--dry-run` | | print commands, execute nothing, write no config |
| `--yes` | | non-interactive; fail instead of prompting |
| `--from`, `--only` | | resume / run a single phase (`P0`..`P7`) |
| `-v`, `--verbose` | | log every command |

## Config & cache

- Config/state: `~/.config/cre/devmachine.yaml` (written after each phase; flags
  override). Use `--config <path>` to keep a separate state.
- Artifact cache: `~/.cache/cre/devmachine/` (agent binary + image tarballs,
  keyed by git SHA / image tag so re-runs skip rebuilds).
- Env overrides: `CRE_CHAINLINK_REPO`, `CRE_JOB_DISTRIBUTOR_REPO`.

## After provisioning

The agent runs as the `remote-agent` systemd service. Start the environment on
your laptop:

```
cd <chainlink>/core/scripts/cre/environment
export CRE_REMOTE_AGENT_EC2_INSTANCE_ID=<instance-id>
export AWS_PROFILE=<profile>
CTF_CONFIGS=configs/workflow-gateway-don-mixed-placement.toml go run . env start
```

The environment resolves the agent host from the instance id via AWS CLI (private
IP preferred). Export `CRE_REMOTE_HOST_IP=<private-ip>` to skip that lookup.

## Troubleshooting

- **`unknown command "devmachine"`** — run from the branch/dir that contains this
  package (`core/scripts/cre/environment`).
- **kubectl errors about "No version is set"** — your `kubectl` shim has no
  version selected; run it with a working kubectl on `PATH`.
- **`permission to create DevMachines: no`** — expected under `--dry-run` (the
  check can't run) or if your role lacks access; use `--instance-id`.
- **SSH hangs / host-key prompt** — handled via `StrictHostKeyChecking=accept-new`;
  if it still hangs, the key likely isn't loaded in `ssh-agent`.
- **Docker permission denied on the machine** — re-run `up`; P2 adds users to the
  `docker` group (existing sessions must reconnect).
- **Agent not reachable** — `go run . devmachine status`, then on the machine
  `systemctl status remote-agent`.

See [`PLAN.md`](./PLAN.md) for the design and phase-by-phase details.
