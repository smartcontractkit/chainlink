---
"chainlink": patch
---

#bugfix Registry-launched OCR instances now take their default bootstrappers from `[P2P.V2].DefaultBootstrappers` instead of `[Capabilities.Peering.V2].DefaultBootstrappers`. The confidential relay now signs responses with the key for `[P2P].PeerID` and no longer reads `[Capabilities.Peering].PeerID`.
