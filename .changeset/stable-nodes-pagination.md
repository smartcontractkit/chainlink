---
"chainlink": patch
---

Sort aggregated NodeStatuses by (ChainID, Name) so the /nodes dashboard paginates deterministically across reloads when more than 100 nodes are configured.

#bugfix
