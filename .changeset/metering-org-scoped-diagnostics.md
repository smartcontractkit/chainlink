---
"chainlink": patch
---

Add org-scoped metering diagnostics: a single structured Info-level log event per finished workflow execution (emitted from Reports.End, outside the receipt retry loop), gated to a hard-coded incident target organization only. Organization identity is resolved from the billing service's GetWorkflowExecutionRates response with the engine org label as fallback; unknown identity never enables logging. The event captures the full conversion evidence (rate card, per-node SpendValue/SpendValueInGasUnits with the conversion path taken and native-vs-legacy scale delta, aggregated native spend and CRE values, metering mode and reason, total credits, receipt errors) without touching billing arithmetic. Temporary incident tooling: revert when the investigation closes.
