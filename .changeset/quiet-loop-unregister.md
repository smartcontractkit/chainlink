---
"chainlink": patch
---

#bugfix Standard capabilities now unregister their LOOP on close, so restarting a capability after a config change no longer fails with "plugin already registered"
