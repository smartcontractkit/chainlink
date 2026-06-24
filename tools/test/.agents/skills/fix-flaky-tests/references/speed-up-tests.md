Optimize slow tests. Exclude non-deterministic flakes/panics.

1. Replace `t.Sleep()` and coarse polling with `require.Eventually` and tight intervals.
2. Look to add `t.Parallel()` where safe and possible.

<lessons-learned>
1. Block mining bottleneck. Lower `cltest.Mine` frequency. Simulation resolve faster. Mine faster than `DeltaRound` or hit "cannot access old blocks" error.
2. EVM polling slow. Lower `LogPollInterval` to 100ms. Nodes find events faster.
3. P2P discovery slow. Lower `DeltaDial` and `DeltaReconcile` to 100ms. Nodes sync faster.
4. OCR timeouts not bottleneck. Lowering `DeltaRound` or `MaxDurationObservation` cause timeout flake. Nodes need time for consensus.
5. Libocr bounds strict. Lowering `ContractConfigTrackerPollInterval` under 1s fail job validation instantly. Do not touch core OCR configs.
</lessons-learned>
