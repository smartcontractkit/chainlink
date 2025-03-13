---
"chainlink": major
---

Upgrade LLO protocol to support sub-seconds reports. #breaking_change

CAUTION: This release needs a careful rollout. It _must_ be tested on a staging DON first by setting half of the nodes to this version, while leaving half at the previous version, to verify inter-version interoperability before node-by-node rollout in production.

NOTE: Protocol version 0 does NOT support gapless handover on sub-second reports. You must upgrade to version 1 for that.

Rollout plan is here: https://smartcontract-it.atlassian.net/browse/MERC-6852
