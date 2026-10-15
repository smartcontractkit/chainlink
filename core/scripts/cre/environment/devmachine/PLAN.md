# Plan: `devmachine` — remote CRE machine provisioning

## Overview

A new cobra subcommand on the Local CRE CLI that takes a developer from nothing to a
ready EC2 DevMachine (Docker installed, images loaded, SSH alias configured) and prints
the 3 manual commands to run the environment.

## Location & wiring

- New package: `core/scripts/cre/environment/devmachine` (module
  `github.com/smartcontractkit/chainlink/core/scripts`).
- Register in `core/scripts/cre/environment/main.go` `init()`:
  `root.RootCmd.AddCommand(devmachine.Cmd)`.
- Invoked as `go run . devmachine up|down|status|handoff`.
- New deps in `core/scripts/go.mod`: `charmbracelet/bubbletea` (+ `bubbles`, `lipgloss`).
  Existing: `spf13/cobra`, `gopkg.in/yaml.v3`.

Suggested files: `cmd.go`, `config.go`, `preflight.go`, `machine.go`, `bootstrap.go`,
`artifacts.go`, `transfer.go`, `load.go`, `handoff.go`, `sshconfig.go`, `exec.go`
(shell-out + dry-run), `ui/` (bubbletea), `devmachine_test.go`.

## Config & cache

- Config: `~/.config/cre/devmachine.yaml` (written after each phase; flags override).
- Cache: `~/.cache/cre/devmachine/{agent,images}/`.
- Env overrides: `CRE_CHAINLINK_REPO`, `CRE_JOB_DISTRIBUTOR_REPO`, `CRE_REMOTE_*`.

```yaml
aws:    { profile: dev-ex-dev/administrator-access, region: us-west-2 }
cluster:{ namespace: cre-staggered-rollout-workflow, privileged: true }
machine:{ name: hybrid-local-cre-bartek, instanceId: i-0d7698..., privateIp: 10.193.2.51, createdByTool: true }
ssh:    { alias: hybrid-local-cre-bartek, key: ~/.ssh/id_ed25519 }
repos:  { chainlink: /.../chainlink, jobDistributor: /.../job-distributor }
images: { chainlink: chainlink-amd:latest, jobDistributor: job-distributor-amd:0.22.1, jdLoadTag: job-distributor:0.22.1 }
agent:  { remotePath: /home/ec2-user/remote-agent, addr: 0.0.0.0:18080, port: 18080 }
env:    { workdir: /.../chainlink/core/scripts/cre/environment, config: configs/workflow-gateway-don-mixed-placement.toml }
```

## CLI surface

```
devmachine up       # P0–P5, then handoff (P6)
devmachine down     # delete DevMachine if tool-created; remove ssh block; --purge-cache
devmachine status   # config + live checks (instance state, docker, agent port)
devmachine handoff  # reprint the 3 manual commands
```

Global: `--config`, `--dry-run`, `--from <phase>`, `--only <phase>`, `--yes`, `-v`.

`up` flags:
- `--name` (default `hybrid-local-cre-<user>`)
- `--namespace` (default `cre-staggered-rollout-workflow`)
- `--team`, `--cost-center`
- `--instance-type` (default **`m8a.2xlarge`**)
- `--volume-size` (default `200`)
- `--az` (default `a`)
- `--instance-id`, `--region`, `--aws-profile` (non-privileged path)
- `--chainlink-repo`, `--job-distributor-repo`
- `--ssh-key` (default `~/.ssh/id_ed25519`)
- `--chainlink-image` (default `chainlink-amd:latest`)
- `--jd-image` (default `job-distributor-amd:0.22.1`), `--jd-tag` (default `v0.22.1`)

## Phases

### P0 — Preflight

- Verify binaries: `aws kubectl docker buildx go rsync ssh` (with install hints).
- `aws sts get-caller-identity` (hint `aws sso login --profile` if expired).
- `kubectl config current-context`; `kubectl auth can-i create devmachines -n <ns>` →
  sets `privileged`.
- Auto-detect repos: `CRE_CHAINLINK_REPO`/`CRE_JOB_DISTRIBUTOR_REPO` env → common paths
  (`~/chainlink`, `~/Desktop/repos/*`, `~/src/*`) → prompt.
- Merge flags/config/prompts; persist config.

### P1 — Machine

Privileged:

```yaml
apiVersion: aws.platform.smartcontract.com/v1alpha1
kind: DevMachine
metadata:
  name: <name>
  namespace: <ns>
  labels:
    chain.link/env: dev
    chain.link/product: cre
    chain.link/team: <team>
    chain.link/cost-center: <cost-center>
spec:
  availabilityZone: <az>
  ami: amazon-linux-2023
  instanceType: <instanceType>   # default m8a.2xlarge
  volumeSize: <volumeSize>
```

- `kubectl apply` the manifest.
- Poll `kubectl -n <ns> get devmachine <name> -o json` until
  `conditions[type=Ready].status=="True"` **and** `.status.instanceId != ""`
  (timeout ~10m, live status).
- Read `.status.instanceId` + `.status.privateIp`.

Non-privileged: use `--instance-id/--region/--aws-profile`; verify via
`aws ec2 describe-instances`; read private IP; `createdByTool=false`.

### P2 — Bootstrap

One `aws ssm send-command` (AWS-RunShellScript), poll + stream; idempotent:

```sh
command -v docker >/dev/null || dnf install -y docker
command -v rsync  >/dev/null || dnf install -y rsync
systemctl enable --now docker sshd
install -d -m700 -o ec2-user -g ec2-user /home/ec2-user/.ssh
grep -qF '<pubkey>' /home/ec2-user/.ssh/authorized_keys 2>/dev/null || echo '<pubkey>' >> /home/ec2-user/.ssh/authorized_keys
chown -R ec2-user:ec2-user /home/ec2-user/.ssh
chmod 700 /home/ec2-user/.ssh && chmod 600 /home/ec2-user/.ssh/authorized_keys
restorecon -R /home/ec2-user/.ssh 2>/dev/null || true
usermod -aG docker ec2-user || true
usermod -aG docker ssm-user || true
```

### P3 — Artifacts (cached)

- Agent (from `<chainlink>/system-tests/lib`):
  `GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o <cache>/remote-agent ./cre/environment/remoteexec/agent/cmd/local-agent`;
  key by chainlink git SHA.
- CL image (from chainlink repo):
  `docker buildx build --platform linux/amd64 --build-arg COMMIT_SHA=$(git rev-parse HEAD) --build-arg VERSION_TAG=$(git describe --always) --build-arg CL_AUTO_DOCKER_TAG=develop -f core/chainlink.Dockerfile . -t chainlink-amd:latest --load`.
- JD image (from JD repo): checkout `v0.22.1`, build server binary,
  `docker buildx build --platform linux/amd64 -t job-distributor-amd:0.22.1 --load .`.
- `docker save` both → cache tars keyed by tag+digest; skip if cached.

### P4 — Transfer

- Ensure managed `~/.ssh/config` block (marker comments, overwrite stale):

  ```
  Host <alias>
    HostName <instanceId>
    User ec2-user
    IdentityFile <key>
    ProxyCommand aws ssm start-session --target %h --region <region> --profile <profile> --document-name AWS-StartSSHSession --parameters portNumber=%p
  ```

- `rsync -avz --progress <cache>/remote-agent <cache>/*.tar <alias>:~/`

### P5 — Load

Via ssh:

```
docker load -i job-distributor-0.22.1-amd.tar
docker tag job-distributor-amd:0.22.1 job-distributor:0.22.1
docker load -i chainlink-amd-latest.tar
docker images
```

### P6 — Handoff (print only)

```
ssh hybrid-local-cre-bartek
# on machine: ~/remote-agent -addr 0.0.0.0:18080

# on laptop, from <chainlink>/core/scripts/cre/environment:
export CRE_REMOTE_AGENT_EC2_INSTANCE_ID=i-0d769858467ae9042
export AWS_PROFILE=dev-ex-dev/administrator-access
CTF_CONFIGS=configs/workflow-gateway-don-mixed-placement.toml go run . env start
```

Optionally export `CRE_REMOTE_HOST_IP=<privateIp>` to skip the env's AWS lookup.

## TUI (Bubble Tea)

- Phase list with spinners; live build/rsync tail.
- On error: show the exact failing command with retry/skip/abort.
- State checkpointed for `--from`.
- `--yes`/non-tty → plain logs.

## Naming & teardown

- Default name `hybrid-local-cre-<user>`, user sanitized to DNS-1123
  (`bartek.tofel` → `bartek-tofel`).
- `down`: `kubectl delete devmachine <name> -n <ns>` only if `createdByTool`; remove
  managed ssh block; optional `--purge-cache`.

## Testing

- Unit: name sanitization, config merge, manifest rendering, Ready/instanceId jsonpath
  parsing, ssh-block idempotency, `--dry-run` golden.
- Manual e2e: run against a real DevMachine, then `devmachine status` + `env start`.

## Milestones

1. Package scaffold + cobra wiring + config + P0.
2. P1 create/poll + accept-existing.
3. P2 SSM bootstrap.
4. P3 build+cache (agent → images).
5. P4 transfer + managed ssh block.
6. P5 load/retag.
7. P6 handoff + `status`.
8. `down` + `--dry-run`/`--from`.
9. Bubble Tea TUI.

## Risks

- CL image build time + ~1.5 GB transfer (cache mitigates; S3 later needs a bucket the
  instance role can read — current role only has `AmazonSSMManagedInstanceCore`).
- Region/namespace/team defaults vary per team → flags.
- `down` must no-op on non-created machines.
