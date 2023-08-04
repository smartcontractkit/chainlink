// SPDX-License-Identifier: MIT
pragma solidity ^0.8.4;

import "../../../vrf/VRF.sol";
import {VRFCoordinatorV2Plus} from "../VRFCoordinatorV2Plus.sol";

contract ExposedVRFCoordinatorV2Plus is VRFCoordinatorV2Plus {
  constructor(address blockhashStore) VRFCoordinatorV2Plus(blockhashStore) {}

  function computeRequestIdExternal(
    bytes32 keyHash,
    address sender,
    uint256 subId,
    uint64 nonce
  ) external pure returns (uint256, uint256) {
    return computeRequestId(keyHash, sender, subId, nonce);
  }

  function isTargetRegisteredExternal(address target) external view returns (bool) {
    return isTargetRegistered(target);
  }

  function getRandomnessFromProofExternal(
    Proof calldata proof,
    RequestCommitment calldata rc
  ) external view returns (Output memory) {
    return getRandomnessFromProof(proof, rc);
  }
}
