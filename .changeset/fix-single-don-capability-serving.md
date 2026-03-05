---
"chainlink": patch
---

Fix capability serving in single-DON topologies where the same DON is both workflow and capability DON (e.g. local CRE). Previously capabilities failed with "empty workflowDONs provided" because only remote workflow DONs were passed to serveCapabilities. #bugfix
