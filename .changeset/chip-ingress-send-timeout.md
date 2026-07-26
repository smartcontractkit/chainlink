---
"chainlink": patch
---

#changed Chip-ingress batch emitter defaults tuned from staging/prod capacity analysis:
`ChipIngressBufferSize` 1000 → 10000, `ChipIngressMaxBatchSize` 500 → 1000,
`ChipIngressSendInterval` 100ms → 500ms, `ChipIngressSendTimeout` 3s → 10s,
`ChipIngressDrainTimeout` 10s → 30s.
