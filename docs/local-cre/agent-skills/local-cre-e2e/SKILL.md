---
name: local-cre-e2e
description: Configure and run local CRE environments and CRE end-to-end tests in the chainlink repo. Use this when starting local CRE on the default topology, running smoke or regression CRE tests, or creating a custom topology to override flags, limits, capability config, or user config overrides.
---

# Local CRE E2E

Use this skill when working in the `chainlink` repo and you need to:
- start, stop, or restart local CRE
- run CRE smoke or regression e2e tests on the default topology
- run a test against a specific topology
- create a custom topology to override limits, flags, capability config, or user config overrides

This skill is for local CRE system-test workflows, not for generic unit tests.

## Assumptions

- Repo root is the `chainlink` checkout.
- Local CRE commands are run from `core/scripts/cre/environment`.
- CRE e2e test commands are run from `system-tests/tests`.
- Only one local CRE environment should be treated as active at a time unless the harness is explicitly known to support isolation.

## Default Workflow

Use the default topology when the user asks to run the standard local CRE suite or to verify a change without special flags.

1. Stop any existing local CRE environment:

```bash
cd core/scripts/cre/environment
go run . env stop -a
```

2. If the environment has not been prepared yet, set it up once:

```bash
cd core/scripts/cre/environment
go run . env setup
```

3. Start local CRE on the default topology:

```bash
cd core/scripts/cre/environment
go run . env start
```

4. Optionally bring up observability helpers:

```bash
go run . obs up
```

Use `--with-chip-ingress-stack` on `env start` when the test depends on the real stack, or when you want Red Panda Console to debug workflow events. (`--with-beholder` is deprecated.)

## Running E2E Tests On The Default Topology

For the normal CRE smoke suite:

```bash
cd system-tests/tests
go test ./smoke/cre -timeout 20m -run '^Test_CRE_'
```

For only the V2 smoke suite:

```bash
cd system-tests/tests
go test ./smoke/cre -timeout 15m -run '^Test_CRE_V2'
```

For regression tests:

```bash
cd system-tests/tests
go test ./regression/cre -timeout 20m -run '^Test_CRE_'
```

Rule of thumb:
- `smoke` is for happy-path and sanity coverage
- `regression` is for edge cases and negative cases

## Running A Specific Test Or Bucket

Use a narrow regex when debugging a single scenario or bucket:

```bash
cd system-tests/tests
go test ./smoke/cre -timeout 20m -run '^Test_CRE_V2_Suite_Bucket_B$' -count=1 -v
```

Examples:

```bash
cd system-tests/tests
go test ./smoke/cre -timeout 20m -run 'Test_CRE_V2_Suite_Bucket_B/.*/Vault' -count=1 -v
```

```bash
cd system-tests/tests
go test ./regression/cre -timeout 20m -run '^Test_CRE_V2_Consensus_Regression$' -count=1 -v
```

Prefer `-count=1` when re-running flaky or stateful CRE scenarios.

## Using A Specific Topology

Use a non-default topology when the test requires a specific DON layout, chain, or feature configuration.

1. Stop the existing environment:

```bash
cd core/scripts/cre/environment
go run . env stop -a
```

2. Start local CRE with `CTF_CONFIGS` pointing at the topology file:

```bash
cd core/scripts/cre/environment
CTF_CONFIGS=./configs/workflow-gateway-capabilities-don.toml go run . env start
```

3. Run the target test:

```bash
cd system-tests/tests
TOPOLOGY_NAME=workflow-gateway-capabilities \
go test ./smoke/cre -timeout 20m -run '^Test_CRE_V2_Suite_Bucket_B$' -count=1 -v
```

`TOPOLOGY_NAME` is optional but useful because many CRE suite tests include it in subtest names.

## Creating A Custom Topology

Create a custom topology when the user wants to override:
- limits
- feature flags
- capability config
- DON composition
- additional sources
- `user_config_overrides`

Workflow:

1. Pick the closest existing topology from `core/scripts/cre/environment/configs/`.
2. Copy it to a new file in the same directory.
3. Change only the fields needed for the scenario.
4. Start local CRE with `CTF_CONFIGS=<new topology>`.
5. Run only the relevant tests first.

Example:

```bash
cd core/scripts/cre/environment/configs
cp workflow-gateway-capabilities-don.toml workflow-gateway-capabilities-don-my-override.toml
```

Then start it:

```bash
cd ../
CTF_CONFIGS=./configs/workflow-gateway-capabilities-don-my-override.toml go run . env start
```

Then run the intended tests:

```bash
cd ../../../system-tests/tests
TOPOLOGY_NAME=workflow-gateway-capabilities-my-override \
go test ./smoke/cre -timeout 20m -run '^Test_CRE_V2_Suite_Bucket_B$' -count=1 -v
```

## Override Guidelines

When making a custom topology:
- keep the diff small and purpose-specific
- prefer copying the nearest topology instead of building a new one from scratch
- use a descriptive filename that states what changed
- do not change unrelated images, chains, or capabilities unless the test needs it
- if the topology is only for a one-off local check, keep it local and avoid adding it to CI

Typical override points:
- `nodesets.capability_configs`
- `nodesets.user_config_overrides`
- CRE feature flags
- additional mock or support-service endpoints

## Restart And Cleanup

When changing topology or low-level config, prefer a full stop/start instead of assuming the running environment will converge.

Clean restart:

```bash
cd core/scripts/cre/environment
go run . env stop -a
CTF_CONFIGS=./configs/<topology>.toml go run . env start
```

When done:

```bash
cd core/scripts/cre/environment
go run . env stop -a
```

## Troubleshooting

- If tests unexpectedly use the wrong topology, stop local CRE and restart with the intended `CTF_CONFIGS`.
- If the test suite appears to reuse stale state, rerun with `-count=1`.
- If a test depends on logs, traces, or dashboards, bring up `go run . obs up`.
- If a topology-specific failure looks unrelated to the test, first confirm the environment actually started with the intended topology.

## References

For longer repo-specific guidance, see:
- `docs/local-cre/index.md`
- `docs/local-cre/system-tests/index.md`
- `docs/local-cre/system-tests/running-tests.md`
