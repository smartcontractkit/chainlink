# Local CRE — Mixed-Env Topology

Mixed-env is a CI test mode for Local CRE that runs each DON with **two code versions at once** — half the nodes built from your PR, half from `develop` — and **fails the run if the two halves ever disagree**. Its job is to catch **non-determinism and cross-version incompatibilities** a PR might introduce, *before* it merges.

## Why it exists

A DON only works if its nodes compute the **same** result from the same inputs. If a change makes nodes produce different consensus reports, or different cross-DON request/response payloads, the DON can no longer agree. Normally you'd only discover this *after* merge — or during a rolling upgrade, when old and new nodes run side by side.

Mixed-env deliberately puts old (`develop`) and new (PR) nodes in the **same** DON, so any such disagreement shows up as a failed CI check instead of a production surprise.

## The topology it runs (2-2)

It's the standard capabilities topology, except every multi-node DON is split **2-2** across two images:

| DON | Nodes | Image split |
|---|---|---|
| `workflow` | 4 | 2 × PR · 2 × develop |
| `chain-cap` (capabilities) | 4 | 2 × PR · 2 × develop |
| `vault` | 4 | 2 × PR · 2 × develop |
| `bootstrap` / `gateway` | 1 | develop |

```
            workflow DON                        chain-cap / vault DON
   ┌──────┬──────┬──────┬──────┐        ┌──────┬──────┬──────┬──────┐
   │  PR  │  PR  │ dev  │ dev  │        │  PR  │  PR  │ dev  │ dev  │
   └──────┴──────┴──────┴──────┘        └──────┴──────┴──────┴──────┘
```

- **PR image** = the Chainlink node image built from your branch.
- **baseline image** = the exact commit of the branch your PR **merges into** — the base of the CI
  auto-merge. For a PR targeting `develop` that's the `develop` commit you branched from; for a PR
  targeting a release branch (e.g. `release/2.57.1`) it's that release commit. Because it's precisely
  the branch your PR image derives from, the two halves differ **only** by your PR's changes — so the
  base branch's own churn never shows up as a false positive.
  - If that image isn't published, the fallback depends on the base branch:
    - **`develop` base** → the latest cached `develop` **nightly** build.
    - **A stacked PR** (base is another feature branch that ultimately targets `develop`) → also the
      **develop nightly**. This compares the *whole stack* (parent branch + your PR) against `develop`,
      so a **parent** branch that changes runtime behavior can surface a marker here — that's a real
      potential non-determinism (it'll be caught when the parent merges too); use the
      [`skip-mixed-env` label](#required-check--emergency-bypass) if it's expected.
    - **A release branch** (e.g. `release/2.57.1`) → **skipped** (a warning is logged). There's no correct
      stand-in: the develop nightly would be the wrong branch and would flood the run with the whole
      release↔develop divergence. Release-branch commits aren't force-built the way `develop` is, so
      their branch-tip image is often absent, and when it is the release PR simply doesn't get this check.
- Chains, capabilities, ports, and everything else are identical to the normal topology — only the per-node images change.

## What it checks

While the normal smoke tests run, mixed-env watches **every node's logs** for the tell-tale signs that the PR and develop nodes disagreed. If **any** of these appears even once, the test is marked **failed** with the reason **"Non-Determinism introduced"**:

| Layer | What a hit means | Log line it watches for |
|---|---|---|
| Consensus (OCR) | nodes produced different reports for the same round | `This is commonly caused by non-determinism` |
| Cross-DON request | same request, but different payloads across nodes | `received messages with the same id and different payloads` |
| Cross-DON response | nodes returned different responses for one request | `received multiple unique responses for the same request` / `response quorum unreachable` |

These lines are **silent** when every node runs the same code, so a single occurrence is a reliable signal that the PR changed behavior in a way `develop` does not agree with.

## What this will catch

- Changes that make a node compute a **different consensus report** than `develop` — serialization/encoding changes, added or reordered fields, map-ordering bugs, rounding, or any non-deterministic logic.
- Changes that alter a **cross-DON request or response payload**.
- **Backward-incompatible** changes that would break a mixed old/new deployment (e.g. a rolling upgrade).

## What it will NOT catch (good to know)

- It only triggers on tests that actually exercise **consensus** or a **cross-DON capability call**. Pure gateway/HTTP-only paths won't set it off.
- It flags **any** divergence from `develop` — including **intentional** report/payload changes. Those are genuine incompatibilities: if the change is deliberate and will be rolled out safely, bypass the check with the [`skip-mixed-env` label](#required-check--emergency-bypass) rather than "fixing" it.
- It compares against the `develop` commit your PR is based on — not the very latest `develop`. So a change that only conflicts with `develop` commits landed **after** you branched is caught by the nightly full-matrix sweep (which pins the latest `develop`), not by the per-PR run. Keep your branch reasonably current for the tightest signal.

## When it runs

- **CI:** automatically on CRE-affecting PRs, in its **own** workflow `.github/workflows/cre-mixed-env-tests.yaml` (kept separate from `cre-system-tests.yaml` so that file stays simple). It runs the OCR3/DON2DON-heavy tests — `Test_CRE_V2_Suite_Bucket_A`, `Test_CRE_V2_Suite_Bucket_B`, and the `Test_CRE_V2_EVM_Read_*` suite — under mixed-env. The PR image and the develop image are both already built (per-PR and nightly), so no extra image builds are added. A dedicated **Check for non-determinism** step scans the live node containers after the suite and fails the job on any marker — kept separate from the auto-quarantined test step so the failure can't be swallowed.

## Required check & emergency bypass

Mixed-env is a **required** check: it feeds the **ETH Smoke Tests** merge gate (`check-e2e-test-results` in `.github/workflows/integration-tests.yml`), so a non-determinism failure **blocks merge**. A run that is legitimately skipped (a non-CRE PR, or a release-branch PR with no baseline image — see above) counts as a non-blocking warning, not a failure.

To bypass it in an emergency, in order of preference:

1. **`skip-mixed-env` PR label (self-service).** Apply the label to the PR and re-run CI. The mixed-env job is then skipped, which the gate treats as a warning, so merge is unblocked — no admin needed. Like the E2E regression tests, mixed-env does not run in the merge queue, so the bypass carries through to the queue. Use this for a *deliberate*, safely-rolled-out change that the check correctly flags, or to unblock while an infra issue is sorted out.
2. **Admin / ruleset bypass.** Repo admins (and ruleset bypass actors) can merge past the failing required check directly.
3. **De-wire the gate.** For a durable disable, remove `run-core-cre-mixed-env-tests` from the `check-e2e-test-results` job's `needs` and its `check_result` line.

> The `CRE_NONDETERMINISM_CHECK` env var only toggles the in-test scan when running locally; it is not a CI merge lever. Use the label in CI.

## Running it locally

Mixed-env compares two pre-built images, so render the topology with your two image refs and start the environment **without** `CTF_CHAINLINK_IMAGE` (a non-empty value forces every node onto one image):

```bash
cd core/scripts/cre/environment

# 1. Render the topology with two images (your build + a develop build).
CRE_PR_IMAGE=<your image> CRE_BASELINE_IMAGE=<develop image> \
  ./configs/render-mixed-env.sh

# 2. Start against the rendered config — do NOT set CTF_CHAINLINK_IMAGE.
CTF_CONFIGS=configs/mixed-env-don.toml go run . env start
```

### Running it locally with local images

By default the rendered topology sets `pull_image = true`, which makes the environment pull each node image from a registry. If both your images are local builds, that pull fails — a bare ref like `chainlink:develop` resolves against Docker Hub, where no such public repo exists (`pull access denied for chainlink`). Pass `--local` to render with `pull_image = false` so the environment uses the images already in your local Docker daemon.

If you've run Local CRE before, your current branch is already built as `chainlink-tmp:latest`. To get a local develop image, build one tagged `develop`:

```bash
DOCKER_TAG=develop make docker
```

Confirm both images exist:

```bash
docker image ls | grep chainlink
# chainlink-tmp:latest
# chainlink:develop
```

Then render with `--local` and start as usual:

```bash
CRE_PR_IMAGE=chainlink-tmp:latest CRE_BASELINE_IMAGE=chainlink:develop \
  ./configs/render-mixed-env.sh --local

CTF_CONFIGS=configs/mixed-env-don.toml go run . env start
```

Then run the smoke suite with the check active:

```bash
cd system-tests/tests
TOPOLOGY_NAME=mixed-env go test ./smoke/cre -run '^Test_CRE_V2_Suite_Bucket_A$' -timeout 20m
```

The check turns on automatically when `TOPOLOGY_NAME` contains `mixed-env` (or set `CRE_NONDETERMINISM_CHECK=true`). After the tests finish, every node's logs are scanned; a single marker fails the run with `Non-Determinism introduced` and lists the offending container(s).
