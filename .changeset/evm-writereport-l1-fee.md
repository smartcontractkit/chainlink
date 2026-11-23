---
"chainlink": patch
---

#bugfix Include the L1 data fee in the EVM WriteReport transaction fee and gas metering when `FeatureEVMWriteReportL1FeeActivePeriod` is active. The relayer set dropped it before the fee was calculated.
