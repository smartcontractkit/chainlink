---
"chainlink": patch
---

#internal Dockerfile: drop apt version pins so the image takes Ubuntu's current ca-certificates (20260816ubuntu0.24.04.1, eMudhra CAs removed). The pinned 20260601 bundle and the pinned curl and jq builds are no longer in the archive.
