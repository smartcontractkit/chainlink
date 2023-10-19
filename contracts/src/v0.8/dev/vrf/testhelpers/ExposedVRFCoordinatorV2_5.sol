// SPDX-License-Identifier: MIT
pragma solidity ^0.8.4;

import "../../../vrf/VRF.sol";
import {VRFCoordinatorV2_5} from "../VRFCoordinatorV2_5.sol";
import "../../../vendor/openzeppelin-solidity/v4.7.3/contracts/utils/structs/EnumerableSet.sol";

contract ExposedVRFCoordinatorV2_5 is VRFCoordinatorV2_5 {
  using EnumerableSet for EnumerableSet.UintSet;

  constructor(address blockhashStore) VRFCoordinatorV2_5(blockhashStore) {}

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

  function getActiveSubscriptionIdsLength() external view returns (uint256) {
    return s_subIds.length();
  }

  function getSubscriptionConfig(uint256 subId) external view returns (SubscriptionConfig memory) {
    return s_subscriptionConfigs[subId];
  }

  function getSubscriptionStruct(uint256 subId) external view returns (Subscription memory) {
    return s_subscriptions[subId];
  }

  function setTotalBalanceTestingOnlyXXX(uint96 newBalance) external {
    s_totalBalance = newBalance;
  }

  function setTotalNativeBalanceTestingOnlyXXX(uint96 newBalance) external {
    s_totalNativeBalance = newBalance;
  }

  function setWithdrawableTokensTestingOnlyXXX(address oracle, uint96 newBalance) external {
    s_withdrawableTokens[oracle] = newBalance;
  }

  function getWithdrawableTokensTestingOnlyXXX(address oracle) external view returns (uint96) {
    return s_withdrawableTokens[oracle];
  }

  function setWithdrawableNativeTestingOnlyXXX(address oracle, uint96 newBalance) external {
    s_withdrawableNative[oracle] = newBalance;
  }

  function getWithdrawableNativeTestingOnlyXXX(address oracle) external view returns (uint96) {
    return s_withdrawableNative[oracle];
  }
}
