# Runbook: Adaptive Oracle E2E Smoke Test

Validates that the v0.3-adaptive-oracle contract design (from
[smartcontractkit/svr-auction-don](https://github.com/smartcontractkit/svr-auction-don),
`contracts/src/v0.3-adaptive-oracle`) plays nice with the existing off-chain OCR2 median plugin,
by driving a real 4-oracle DON against a real `DualAggregator`/`AdaptiveOracle` stack on a
simulated backend.

Three tests:
- `adaptive_oracle_test.go` (`TestIntegration_AdaptiveOracle`) -- basic hook wiring: a real report
  lands, the adaptive rate/min-clamp/reset behaviors work as designed. Uses `AdaptiveRateLogic`
  (geometric-convergence placeholder).
- `adaptive_oracle_deviation_test.go` (`TestIntegration_AdaptiveOracle_DeviationConvergence`) --
  proves the off-chain deviation-triggering design decision: `DualAggregator::latestTransmissionDetails()`
  returns the *adaptive* answer (not the original market answer), so the off-chain median plugin's
  deviation check compares each node's fresh market observation against the adaptive rate. With the
  market rate held fixed, the DON keeps resubmitting on its own -- purely because the adaptive rate
  hasn't caught up yet -- and must stop once it's within the configured deviation threshold. Also
  uses `AdaptiveRateLogic`.
- `adaptive_oracle_capped_test.go` (`TestIntegration_AdaptiveOracle_CappedLogic`) -- same
  integration question, but with `AdaptiveOracle` wired to `CappedAdaptiveRateLogic` instead.
  Unlike `AdaptiveRateLogic`, this logic has no memory of its own: it always serves
  min(marketRate, referenceRate) directly, so a single landed report is enough to bring the
  adaptive rate fully in line with the market (below the reference rate) or fully capped (above
  it) -- no multi-round convergence, and no further reports once settled.

## 1. One-time setup

```bash
# Install abigen (only needed once, or if not already installed)
go install github.com/ethereum/go-ethereum/cmd/abigen@latest

# jq is used for extracting ABI/bytecode from forge artifacts
which jq || brew install jq
```

## 2. Compile the contracts (svr-auction-don repo)

```bash
cd svr-auction-don/contracts

# DualAggregator exceeds the EIP-170 24,576-byte deployment limit at the repo's default
# optimizer_runs=1_000_000. Real backends (go-ethereum's simulated.Backend included) enforce
# this; forge's local anvil doesn't, so `forge test` won't catch it. Build at a lower
# optimizer_runs just for generating bindings -- this does NOT touch the checked-in foundry.toml.
FOUNDRY_OPTIMIZER_RUNS=200 forge build
```

## 3. Extract ABI/bytecode and generate Go bindings

```bash
WORK=/tmp/adaptiveoracle-abi   # or your preferred scratch dir
mkdir -p "$WORK"

for pair in \
  "v0.3-adaptive-oracle/DualAggregator.sol/DualAggregator.json:DualAggregator" \
  "AdaptiveOracle.sol/AdaptiveOracle.json:AdaptiveOracle" \
  "AdaptiveRateLogic.sol/AdaptiveRateLogic.json:AdaptiveRateLogic" \
  "CappedAdaptiveRateLogic.sol/CappedAdaptiveRateLogic.json:CappedAdaptiveRateLogic" \
  "ReferenceRateAdapterMock.sol/ReferenceRateAdapterMock.json:ReferenceRateAdapterMock"; do
  path="${pair%%:*}"
  name="${pair##*:}"
  jq '.abi' "foundry-artifacts/$path" > "$WORK/$name.abi"
  jq -r '.bytecode.object' "foundry-artifacts/$path" > "$WORK/$name.bin"
  abigen --abi "$WORK/$name.abi" --bin "$WORK/$name.bin" \
    --pkg adaptiveoracle --type "$name" --out "$WORK/$name.go"
done
```

## 4. Copy bindings into the chainlink repo

```bash
cp "$WORK"/DualAggregator.go "$WORK"/AdaptiveOracle.go "$WORK"/AdaptiveRateLogic.go "$WORK"/CappedAdaptiveRateLogic.go "$WORK"/ReferenceRateAdapterMock.go \
  <chainlink-repo>/core/internal/features/adaptiveoracle/generated/
```

## 5. Restore the contracts repo's normal build artifacts (optional but tidy)

```bash
cd svr-auction-don/contracts
forge build   # back to the repo's default optimizer_runs=1_000_000
```

`foundry-artifacts/` is gitignored, so this step only matters if you're about to keep working in
that repo and don't want the low-optimizer build sitting in your local artifacts.

## 6. Verify the bindings compile

```bash
cd <chainlink-repo>
go build ./core/internal/features/adaptiveoracle/...
go vet ./core/internal/features/adaptiveoracle/...
```

## 7. Start Postgres

The test framework (`pgtestdb`) needs a real Postgres to provision ephemeral per-node test
databases.

```bash
docker run -d --name adaptiveoracle-test-pg -p 5432:5432 \
  -e POSTGRES_USER=chainlink_dev -e POSTGRES_PASSWORD=insecurepassword -e POSTGRES_DB=chainlink_test \
  postgres:16

# wait for it to be ready
docker exec adaptiveoracle-test-pg pg_isready -U chainlink_dev
```

If the container already exists from a previous run: `docker start adaptiveoracle-test-pg`
instead of `docker run`.

## 8. Run the test

```bash
cd <chainlink-repo>
export CL_DATABASE_URL="postgresql://chainlink_dev:insecurepassword@localhost:5432/chainlink_test?sslmode=disable"

# Run all three tests together (all run in parallel via t.Parallel(); ~90s total):
go test -count=1 -v -timeout 10m ./core/internal/features/adaptiveoracle/...

# Or individually:
go test -run '^TestIntegration_AdaptiveOracle$' -count=1 -v -timeout 8m ./core/internal/features/adaptiveoracle/...
go test -run '^TestIntegration_AdaptiveOracle_DeviationConvergence$' -count=1 -v -timeout 8m ./core/internal/features/adaptiveoracle/...
go test -run '^TestIntegration_AdaptiveOracle_CappedLogic$' -count=1 -v -timeout 8m ./core/internal/features/adaptiveoracle/...
```

- `-count=1` disables Go's test cache (forces a real re-run instead of reporting a stale cached
  pass).
- `TestIntegration_AdaptiveOracle` takes ~70s; `TestIntegration_AdaptiveOracle_DeviationConvergence`
  takes ~90s (it deliberately waits an extra 15s after convergence to confirm the DON has actually
  stopped reporting); `TestIntegration_AdaptiveOracle_CappedLogic` takes ~85s (same 15s
  stopped-reporting check). All three run concurrently via `t.Parallel()`, so running the whole
  package takes about as long as the slowest one, not the sum. `-timeout 10m` gives headroom.
- Since all three tests spin up their own Postgres-backed node stacks concurrently, the default
  Postgres `max_connections=100` can be exhausted, surfacing as `FATAL: remaining connection slots
  are reserved for roles with the SUPERUSER attribute`. If you see that, raise the limit and
  restart the container: `docker exec adaptiveoracle-test-pg psql -U chainlink_dev -d
  chainlink_test -c "ALTER SYSTEM SET max_connections = 300;" && docker restart
  adaptiveoracle-test-pg`.
- Look for `--- PASS:` / `--- FAIL:` at the very end -- with `-v` there's a lot of node log noise
  in between. Grep for `"adaptive rate moved to"` in the output to see the convergence sequence the
  deviation test observed (expect something like 950000000 -> 925000000 -> 912500000 -> 906250000
  for the default constants).

## 9. Cleanup (optional)

```bash
docker stop adaptiveoracle-test-pg && docker rm adaptiveoracle-test-pg
```

## When do you need steps 2-6 again?

Only when `DualAggregator`, `AdaptiveOracle`, `AdaptiveRateLogic`, `CappedAdaptiveRateLogic`, or
`ReferenceRateAdapterMock` change on the contracts side. If you're just iterating on the Go
test/helper code (`adaptive_oracle_test.go`, `adaptive_oracle_helper.go`), skip straight to step
7/8.
