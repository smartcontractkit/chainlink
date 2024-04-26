---
"chainlink": patch
---

Add configurability to mercury transmitter

```toml
[Mercury.Transmitter]
MaxTransmitQueueSize = 10_000 # Default
TransmitTimeout = "5s" # Default
```
