Optimize slow tests. See below sections for first checks and former lessons learned on this repo. If done with these suggestions, conduct your own exploration through code and logs to determine bottlenecks and hypothesize fixes. If you need extra logs/traces to confirm/deny theories, add them, or ask the user to.

<first-look>
1. Replace `time.Sleep()` and coarse polling with `require.Eventually` and tight intervals.
2. Look to add `t.Parallel()` where safe and possible.
3. Use `testing/synctest` where sensible to improve speed and stability.
</first-look>

<lessons-learned>
1. Block mining bottleneck. Lower `cltest.Mine` frequency. Simulation resolves faster. Mine faster than `DeltaRound` or hit "cannot access old blocks" error.
2. EVM polling slow. Lower `LogPollInterval` to 100ms. Nodes find events faster.

<node-config>
Prefer not to mess with node/OCR configs except as last resort.
1. P2P discovery slow. Lower `DeltaDial` and `DeltaReconcile` to 100ms. Nodes sync faster.
2. OCR timeouts not bottleneck. Lowering `DeltaRound` or `MaxDurationObservation` cause timeout flake. Nodes need time for consensus.
3. Libocr bounds strict. Lowering `ContractConfigTrackerPollInterval` under 1s fail job validation instantly. Do not touch core OCR configs.
</node-config>
</lessons-learned>
