---
"@chainlink/contracts": minor
---

- Misc VRF V2+ contract changes
  - Reuse struct RequestCommitmentV2Plus from VRFTypes
  - Fix interface name IVRFCoordinatorV2PlusFulfill in BatchVRFCoordinatorV2Plus to avoid confusion with IVRFCoordinatorV2Plus.sol
  - Remove unused errors
  - Rename variables for readability
  - Fix comments
  - Minor gas optimisation (++i)
