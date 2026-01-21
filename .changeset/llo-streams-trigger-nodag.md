---
"chainlink": minor
---

#added NoDAG CRE LLO Streams Trigger

Adds support for LLO (Low-Latency Oracle) reports to be consumed by CRE workflows via the `streams-trigger@2.0.0` capability. This enables workflows to receive real-time price data from Streams DON oracle networks.

Key changes:
- Modified CRE Transmitter to emit trigger events in both V1 (values.Map) and V2 (proto) formats for backward compatibility
- Added support for multiple LLO report formats (Format 5: protobuf, Format 7: ABI-encoded)
- Added cryptographic signature verification for LLO reports in the SignedReportRemoteAggregator
- Added legacy signature verification for Format 7 reports using v0.3 OCR report context
- Added E2E test demonstrating full data flow: Mock EA → Stream Jobs → LLO Plugin → CRE Transmitter → Workflow
