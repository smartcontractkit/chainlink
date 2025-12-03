---
"chainlink": patch
---

#bugfix Fixed SendRequestEVM to parse CCIP events directly from transaction receipts instead of using blockchain filtering, making it compatible with both OnRamp v1.5 (CCIPSendRequested) and v1.6 (CCIPMessageSent) events
