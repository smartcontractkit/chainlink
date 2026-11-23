---
"chainlink": patch
---

#changed Consensus-queue launches only from the Capabilities Registry through `[Capabilities.Local] RegistryBasedLaunchAllowlist`, and reads its OCR limits from the plugin config in the registry's OCR offchain config. Consensus-queue job specs are rejected.
