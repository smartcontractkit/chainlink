---
"chainlink": minor
---

Remove VRF v2 support including the v2 coordinator, VRFOwner force-fulfillment and custom reverted-txns pipeline, the legacy `vrfv2` pipeline task, and the v2 blockhash store / block header feeder coordinators. Migrate to VRF v2 Plus. #removed #breaking_change
