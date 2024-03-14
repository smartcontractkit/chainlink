---
"chainlink": minor
---

HeadTracker now respects the `FinalityTagEnabled` config option. If the flag is enabled, HeadTracker backfills blocks up to the latest finalized block provided by the corresponding RPC call. To address potential misconfigurations, `HistoryDepth` is now calculated from the latest finalized block instead of the head. NOTE: Consumers (e.g. TXM and LogPoller) do not fully utilize Finality Tag yet.
