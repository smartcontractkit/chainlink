---
"chainlink": patch
---

#added

Add configurability to mercury transmitter

```toml
[Mercury.Transmitter]
MaxTransmitQueueSize = 10_000 # Default
TransmitTimeout = "5s" # Default
```
