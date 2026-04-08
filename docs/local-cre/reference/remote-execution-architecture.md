---
id: local-cre-reference-remote-execution-architecture
title: Remote Execution Architecture
sidebar_label: Remote Execution Architecture
sidebar_position: 1
---

# Remote Execution Architecture

## Goal

Keep responsibilities co-located so contributors can reason about hybrid local/remote execution without hopping across unrelated packages.

## Ownership Boundaries

- `system-tests/lib/cre/environment`
  - High-level environment orchestration.
  - Decides **what** to start and in which order (blockchains, JD, DONs, linking, funding).
  - Consumes remote execution APIs; does not own transport/protocol details.

- `system-tests/lib/cre/environment/remoteexec/client`
  - Remote control-plane client logic.
  - Owns runtime resolution, agent HTTP/retry behavior, start/stop/deploy envelopes, remote stop summary, and agent log normalization.
  - Exposes reusable helpers for orchestrators and workflow artifact deployment call sites.

- `system-tests/lib/cre/environment/remoteexec/agent`
  - Remote data-plane/agent runtime.
  - Owns server handlers, deployment execution, relay lifecycle, and transport contracts used by the agent API.

## Runtime Flow (Hybrid)

1. CLI loads config and builds topology summary.
2. `environment` resolves whether remote components exist.
3. If remote components are present, `remoteexec/client` resolves runtime and performs remote operations.
4. Local components are started directly by `environment` + CTF components.
5. Stop commands route:
   - `env stop`: local only.
   - `env remote stop`: remote only via `remoteexec/client`.
   - `env stop-all`: remote then local.

## Invariants

- Remote HTTP protocol details remain in `remoteexec/client` and `remoteexec/agent`.
- `environment` should not re-introduce ad-hoc remote transport code.
- Placement (`local` vs `remote`) remains the single selector for execution target behavior.
- Remote placement visualization is shown only when at least one component is remote.

## Maintenance Guidance

- When changing agent payloads or operations, update both `remoteexec/agent` and `remoteexec/client` in the same PR.
- When changing orchestration order or placement rules, prefer tests in `system-tests/lib/cre/environment`.
- Keep runbook commands and env var precedence synchronized with code changes in `core/scripts/cre/environment/environment`.
