---
"chainlink": patch
---

#internal Dockerfile: drop apt version pins in the final stage so the image takes Ubuntu's current ca-certificates (20260816ubuntu0.24.04.1, eMudhra CAs removed). The pinned 20260601 bundle and curl 8.5.0-2ubuntu10.15 are no longer in the archive and broke the v2.67.3-rc.1 build.
