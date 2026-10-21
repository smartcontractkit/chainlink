---
"chainlink": patch
---

#internal Rebuild on Ubuntu 24.04 with ca-certificates 20260816ubuntu0.24.04.1, which removes the compromised eMudhra intermediate CAs from the image trust store. Ubuntu 20.04 is out of standard support and does not receive the fix. Build stages move from Debian bullseye to bookworm because the bullseye security archive no longer serves packages. No code changes.
