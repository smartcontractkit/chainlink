---
"chainlink": major
---

Remove Chainlink Automation from the node.

The `ocr2automation` OCR2 plugin type is no longer supported: job specs with
`pluginType = "ocr2automation"` will now fail validation. The
`OCR2.CaptureAutomationCustomTelemetry` config field has been removed, and the
`automation-custom`, `ocr2-automation` and `ocr3-automation` telemetry types no
longer exist. The `evm_upkeep_states` table is dropped by migration 0305. The
`chaincli` tool and the devenv automation and log poller E2E suites are removed.

#removed #breaking_change
