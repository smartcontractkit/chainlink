# Aptos Write Capability EVM Parity Design

Date: 2026-03-16

Author: Gregory Cawthorne

Status: The rebased `feature/aptos-local-cre-minimal-write` stack now delivers working Aptos local CRE read/write parity, including forwarder deploy/config, deterministic one-at-a-time scheduling, smarter retry behavior, typed receiver execution status, and transaction fee propagation. The remaining gaps are mostly structural parity gaps with EVM, not missing local CRE functionality.

## Purpose

This document is the current design and status record for the Aptos local CRE write stack.

It answers four questions:

1. what branch stack is the source of truth
2. what functionality is already implemented
3. what is parity-equivalent to EVM today and what is only partial
4. what CI and follow-up work still remains

## Source Of Truth

The source of truth is the upstream branch named `feature/aptos-local-cre-minimal-write` in all six repos below.

Current upstream heads (all commits signed where noted):

1. `chainlink-protos` branch `feature/aptos-local-cre-minimal-write` at `5ec1a23` (signed)
2. `chainlink-common` branch `feature/aptos-local-cre-minimal-write` at `435e3f3` (signed)
3. `cre-sdk-go` branch `feature/aptos-local-cre-minimal-write` at `11ce4cc` (signed)
4. `chainlink-aptos` branch `feature/aptos-local-cre-minimal-write` at `6b1d669` (signed; includes go.sum tidy fix for check-tidy CI)
5. `capabilities` branch `feature/aptos-local-cre-minimal-write` at `b710d1e` (signed)
6. `chainlink` branch `feature/aptos-local-cre-minimal-write` at `f8cf49b` (signed)

Open PRs:

1. `chainlink-protos` [#314](https://github.com/smartcontractkit/chainlink-protos/pull/314)
2. `chainlink-common` [#1896](https://github.com/smartcontractkit/chainlink-common/pull/1896)
3. `cre-sdk-go` [#119](https://github.com/smartcontractkit/cre-sdk-go/pull/119)
4. `chainlink-aptos` [#386](https://github.com/smartcontractkit/chainlink-aptos/pull/386)
5. `capabilities` [#472](https://github.com/smartcontractkit/capabilities/pull/472)
6. `chainlink` [#21515](https://github.com/smartcontractkit/chainlink/pull/21515)

## What This Stack Delivers

The current rebased stack delivers the write side of local CRE for Aptos, including:

1. Aptos forwarder deployment during local CRE startup
2. Aptos forwarder `set_config`
3. worker/job wiring to the deployed Aptos forwarder
4. saved-state correctness for Aptos forwarder reuse
5. deterministic one-at-a-time Aptos write scheduling
6. `p2pToTransmitterMap` peer-to-transmitter mapping
7. Aptos read/write/roundtrip/expected-failure local CRE smoke coverage
8. smarter retry and failure classification closer to EVM
9. typed receiver execution status in Aptos `WriteReportReply`
10. transaction fee propagation in Aptos `WriteReportReply`

## Authoritative End-To-End Validation

The authoritative local validation remains:

1. from `/Users/gregorycawthorne/Documents/chainlink/core/scripts/cre/environment`:
   - `CTF_CONFIGS=configs/workflow-gateway-don-aptos.toml go run . env start`
2. from `/Users/gregorycawthorne/Documents/chainlink/system-tests/tests`:
   - `go test -timeout 20m -run '^Test_CRE_V2_Aptos_Suite$' ./smoke/cre/ -count=1 -v`

Validated scenarios:

1. `Aptos_Read`
2. `Aptos_Write`
3. `Aptos_Write_Read_Roundtrip`
4. `Aptos_Write_Expected_Failure`

## Corrections To Older Drafts

The following older assumptions are now wrong and should not be reused:

1. the write config key is `p2pToTransmitterMap`, not `aptosTransmitters`
2. this is a six-repo stack, not a five-repo stack
3. `chainlink-protos` is part of the required slice
4. Aptos write replies now include typed receiver execution status and transaction fee
5. the final correct behavior depended on a `cre-sdk-go` decoder fix, not just proto regeneration
6. Aptos retry behavior is no longer success/failure-only polling

## Repo Responsibilities And Implemented Changes

### 1. `chainlink-protos`

Owns the Aptos write reply protobuf contract.

Implemented:

1. `receiver_contract_execution_status`
2. `transaction_fee`
3. protobuf regeneration using the repo-pinned toolchain

Pinned tool versions used:

1. `protoc 29.3`
2. `protoc-gen-go 1.36.6`
3. `protoc-gen-go-grpc 1.5.1`

### 2. `chainlink-common`

Owns shared Aptos type surfaces and status helpers used by capability and relayer layers.

Implemented:

1. Aptos read-capability baseline from the earlier read PR stack
2. `LedgerVersion` support
3. transmitter-bytes compatibility used by the write path
4. shared Aptos write-status helpers
5. shared Aptos reply-shape support for receiver execution status and fee

### 3. `cre-sdk-go`

Owns the Aptos SDK surface consumed by the workflows.

Implemented:

1. Aptos `ViewRequest` SDK support
2. stable hand-written Aptos client layout that survives `clean-generate`
3. corrected hand-written Aptos `WriteReportReply` decoder
4. refreshed generated Aptos client output so checked-in files match the repo's current `make generate` behavior

Current effective Aptos `WriteReportReply` field order:

1. `tx_status`
2. `receiver_contract_execution_status`
3. `tx_hash`
4. `transaction_fee`
5. `error_message`

Note:

The checked-in Aptos generated client now reports `protoc-gen-go v1.36.11`. That is correct for `cre-sdk-go`. Unlike `chainlink-protos`, this repo's `Makefile` installs `protoc-gen-go` from the current `google.golang.org/protobuf` module version during `make generate`, so the generated header comment is expected to track that module version.

### 4. `chainlink-aptos`

Owns the Aptos relayer, TXM integration, binding compilation path, and write-target confirmation behavior.

Implemented:

1. rebased Aptos relayer/TXM write-path support
2. sequence-safe Aptos write submission
3. Aptos write-status classification used in relayer/write-target confirmation
4. expected terminal failures still submit once so workflows receive a real failed on-chain tx hash
5. transaction-fee propagation through the Aptos relayer service surface
6. Dockerized Aptos CLI fallback for Move compilation when host `aptos` is unavailable

### 5. `capabilities`

Owns the Aptos capability startup contract, scheduling, retry semantics, and final reply population.

Implemented:

1. `p2pToTransmitterMap` as the only Aptos write config shape
2. strict startup validation of peer-to-transmitter mapping
3. deterministic one-at-a-time scheduling
4. prior-transmitter failure inspection before local retry
5. `F+1` stopping behavior
6. VM-status-based retry classification
7. typed receiver execution status in Aptos write replies
8. transaction-fee propagation in Aptos write replies
9. `ABORTED` status for confirmed on-chain aborts instead of collapsing them into generic fatal failure

### 6. `chainlink`

Owns local CRE topology, forwarder deployment/config/state persistence, worker config, smoke-suite wiring, and CI/test bootstrap integration.

Implemented:

1. Aptos local CRE topology and capability wiring
2. Capability Registry spec config emitting `p2pToTransmitterMap` only
3. Aptos `WriteReport` configured for deterministic one-at-a-time scheduling
4. Aptos forwarder deployment, `set_config`, worker wiring, and saved-state persistence
5. full local CRE smoke coverage for read, write, roundtrip, and expected failure
6. CI/test bootstrap workarounds for Aptos tests that are forced through the generic CRE matrix

## Implemented Parity Behavior

### 1. Capability Registry Config Shape

Implemented in `chainlink` and enforced in `capabilities`:

1. Aptos spec config emits `p2pToTransmitterMap` only
2. `aptosTransmitters` is obsolete and not used
3. keys are lowercase hex encodings of 32-byte peer IDs
4. values are normalized Aptos transmitter addresses

Startup validation now requires:

1. map exists and is non-empty
2. every DON member peer ID is present
3. every transmitter parses as a valid Aptos address
4. duplicate transmitters are rejected after normalization

### 2. Deterministic One-At-A-Time Scheduling

Implemented:

1. `chainlink` configures Aptos `WriteReport` as a remote executable with one-at-a-time scheduling and `deltaStage`
2. `capabilities` computes a deterministic queue position per transmission ID
3. each node waits `queuePosition * deltaStage`
4. each node checks whether the report already succeeded before submitting

This matches EVM’s architectural layering:

1. one-at-a-time scheduling is owned by the capability layer
2. TXM is not DON-scheduling-aware
3. write-target confirmation is not the scheduler

### 3. Smarter Retry And Failure Behavior

Implemented:

1. later transmitters inspect prior scheduled failed txs before submitting
2. non-forwarder aborts are treated as terminal receiver or user-code failures
3. known terminal forwarder aborts are treated as non-retryable
4. unknown forwarder aborts remain retryable
5. `F+1` stopping behavior exists in the capability layer
6. own failed txs are classified from `vm_status`
7. write-target confirmation and relayer status handling use similar Aptos write-status classification
8. `already processed` is handled specially rather than treated as a generic failure

This is materially closer to EVM than the earlier Aptos implementation.

### 4. Canonical Success/Failure Hash Handling

Implemented:

1. canonical success hash return when another node already succeeded
2. deterministic failed-hash resolution by scanning prior scheduled transmitters first
3. expected-failure workflows receive a real failed on-chain Aptos tx hash

### 5. Forwarder Deployment And Saved-State Correctness

Implemented:

1. Aptos forwarder deployment during local CRE startup
2. Aptos forwarder `set_config` during startup
3. Aptos forwarder address propagation into worker config
4. Aptos forwarder address persistence into `local_cre.toml`
5. saved-state rebuild reuses the same Aptos forwarder rather than silently redeploying another one

This fixed the earlier split-brain forwarder mismatch.

### 6. Receiver Execution Status

Implemented.

Aptos write replies now carry a typed receiver execution status. That makes the Aptos reply surface much closer to EVM and lets downstream consumers distinguish more cleanly between:

1. receiver or user-code failure
2. forwarder or platform failure
3. successful receiver execution

### 7. Fee Propagation

Implemented.

Aptos write replies now propagate transaction fee information rather than returning only status, hash, and error text.

This is the Aptos equivalent of the EVM metering concept, even though it does not affect local CRE write correctness directly.

### 8. Shared Aptos Write-Status Handling

Implemented partially.

There is now meaningful shared Aptos write-status handling across:

1. `capabilities`
2. `chainlink-aptos`

This is a real improvement over the earlier split implementation.

## What Is Still Not Full EVM Parity

The remaining gaps are mostly structural rather than functional blockers for local CRE.

### 1. No Full Shared Forwarder-State Vocabulary Across All Layers

This remains only partial.

EVM still has a cleaner transmission-state model. Aptos is closer than before, but it still does not have one canonical forwarder/transmission-state vocabulary used consistently across:

1. capability layer
2. write target
3. relayer / TXM boundary

### 2. No EVM-Style Gas-Limit-Aware Retry Decision

This is still not implemented.

Reason:

1. EVM has richer forwarder transmission state that can support gas-limit-aware retry decisions
2. Aptos does not currently expose an equivalent high-value signal in the forwarder state shape being used here

This remains lower priority than the already completed parity work.

## Why Aptos Cannot Just Reuse The EVM Local CRE CI Setup

Short answer: because the shared generic CRE topology is a superset for the EVM tests but not for the Aptos tests.

EVM tests generally use the default local CRE config path, and the generic shared CI topology already contains the required EVM chains. The compatibility check therefore passes and no recreation is needed.

Aptos tests explicitly require `workflow-gateway-don-aptos.toml`, which contains:

1. the Aptos blockchain entry
2. Aptos capability wiring
3. Aptos chain-allowlist configuration
4. Aptos write/read capability flags

The generic shared topology does not contain those. So the Aptos tests must recreate into the Aptos topology unless the workflow matrix itself changes.

## Signing

All PR branch commits should be signed. As of 2026-03-16:

- **chainlink-protos [#314](https://github.com/smartcontractkit/chainlink-protos/pull/314)**: 1 commit signed and pushed (head `5ec1a23`).
- **chainlink-common, capabilities, cre-sdk-go, chainlink-aptos, chainlink**: Signed via `sign-and-push-aptos-prs.sh` (or manually with hardware key). Script now includes `chainlink-protos`; run from chainlink repo with signing key available.

## Current CI Status As Of 2026-03-16

Before upstreaming, CI that can be fixed locally has been addressed (e.g. chainlink-aptos check-tidy, chainlink plugin refs—see below). **Validate go.mod dependencies** is expected to fail on any PR in this stack that pins other PRs’ branch commits (e.g. chainlink-aptos, capabilities, chainlink-common at feature-branch refs). That check will stay red until the full stack is merged to main/develop; no code change is required for it.

### `chainlink#21515`

Build fix (plugin install): `plugins.public.yaml` and `go.mod` were updated to pin chainlink-aptos to the current signed branch head (`v0.0.0-20260316114633-6b1d6694544a`) so the Docker plugin-install step can fetch and build the aptos plugin. The previous ref pointed at a pre–force-push commit that was no longer on the remote.

Current non-green checks:

1. `Validate go.mod dependencies` (expected until stack is upstreamed)
2. `sigscanner-check`

Aptos-specific CI on head `f8cf49b` (or later after this fix):

1. `Test_CRE_V2_Aptos_Read` shard is green
2. `Test_CRE_V2_Aptos_Suite` shard is green

Aptos CRE CI bootstrap is resolved on the current branch head.

### `chainlink-aptos#386`

Current non-green checks:

1. `Validate go.mod dependencies`

This is the expected PR-branch-ref validation failure.

Actionable repo CI is otherwise green, including `run-tests`.

### `capabilities#472`

Current non-green checks:

1. `Analyze (javascript-typescript)`
2. `SonarQube Code Analysis`

Actionable Go/build/test CI is green.

### `chainlink-common#1896`

Current non-green checks:

1. `SonarQube Code Analysis`
2. `Validate go.mod dependencies`
3. `check-tidy`

Interpretation:

1. `Validate go.mod dependencies` is the expected PR-branch-ref failure until the stack is upstreamed.
2. `check-tidy`: run `make gomodtidy` in chainlink-common and commit/push any changes to fix.
3. Sonar is external to the code changes themselves.

### `cre-sdk-go#119`

Current non-green checks:

1. `Validate go.mod dependencies`

Interpretation:

1. `Validate go.mod dependencies` is the expected PR-branch-ref failure against `chainlink-protos` `feature/aptos-local-cre-minimal-write`
2. the previously failing `check-tidy` issue was fixed by:
   - removing stale `go.sum` entries
   - refreshing the generated Aptos client so it matches the repo's current generator version

### `chainlink-protos#314`

- All commits signed (head `5ec1a23`).
- Buf build/format/lint CI is green. No changeset required for this proto-only change per maintainer workflow.

## Obsolete Design Choices

These are intentionally obsolete and should not be revived:

1. `aptosTransmitters`
2. launcher changes in production code to special-case Aptos local CRE
3. local-only `.local` configs as part of the upstreamable solution
4. treating Aptos success/failure as only binary visibility without VM-status-aware classification

## Recommended Current Framing

The accurate summary for review is:

1. Aptos local CRE write parity is implemented and passes locally end-to-end
2. the rebased PR stack now includes six repos, not five
3. Aptos write behavior is materially closer to EVM than the earlier implementation
4. key parity items already implemented include:
   - one-at-a-time scheduling
   - peer-to-transmitter mapping
   - smarter retry and abort classification
   - receiver execution status
   - fee propagation
5. the main remaining product-level parity gap is a more unified shared forwarder/transmission state model
6. the main remaining PR-level cleanup is CI hygiene in `chainlink-common` and `cre-sdk-go`, not Aptos local CRE correctness itself
