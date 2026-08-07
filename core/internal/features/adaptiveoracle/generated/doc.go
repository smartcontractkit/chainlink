// Package adaptiveoracle contains abigen-generated Go bindings for the v0.3-adaptive-oracle
// contracts (DualAggregator, AdaptiveOracle, AdaptiveRateLogic, CappedAdaptiveRateLogic) and the
// ReferenceRateAdapterMock test helper, all sourced from https://github.com/smartcontractkit/svr-auction-don
// (contracts/src/v0.3-adaptive-oracle). These are copied in manually for this PoC smoke test,
// following the same pattern used for the DualAggregator bindings in
// core/internal/features/svr/dual_aggregator.go. If/when the svr-auction-don contracts are merged
// into this repo (or published as an importable module), these generated files should be replaced
// with a proper import instead of a manual copy.
//
// Bytecode note: svr-auction-don's foundry.toml uses optimizer_runs = 1_000_000, which pushes
// DualAggregator's runtime bytecode to 24,683 bytes -- over the EIP-170 24,576 byte deployment
// limit enforced by real backends (go-ethereum's simulated.Backend included; forge's local anvil
// does not enforce this, which is why the forge test suite didn't catch it). These bindings were
// instead generated from a one-off build with FOUNDRY_OPTIMIZER_RUNS=200 (DualAggregator shrinks to
// 19,400 bytes), without changing the repo's checked-in foundry.toml. If svr-auction-don's default
// optimizer_runs setting changes, or DualAggregator grows further, these bindings need regenerating
// the same way -- see the Stage 3 write-up for the exact commands.
package adaptiveoracle
