// SPDX-License-Identifier: MIT
pragma solidity ^0.8.4;

import "../../../vrf/VRF.sol";
import {VRFCoordinatorV2Plus} from "../VRFCoordinatorV2Plus.sol";
import "../../../vendor/openzeppelin-solidity/v4.7.3/contracts/utils/structs/EnumerableSet.sol";

contract ExposedVRFCoordinatorV2Plus is VRFCoordinatorV2Plus {
  using EnumerableSet for EnumerableSet.UintSet;

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

  function setTotalEthBalanceTestingOnlyXXX(uint96 newBalance) external {
    s_totalEthBalance = newBalance;
  }

  function setWithdrawableTokensTestingOnlyXXX(address oracle, uint96 newBalance) external {
    s_withdrawableTokens[oracle] = newBalance;
  }

  function getWithdrawableTokensTestingOnlyXXX(address oracle) external view returns (uint96) {
    return s_withdrawableTokens[oracle];
  }

  function setWithdrawableEthTestingOnlyXXX(address oracle, uint96 newBalance) external {
    s_withdrawableEth[oracle] = newBalance;
  }

  function getWithdrawableEthTestingOnlyXXX(address oracle) external view returns (uint96) {
    return s_withdrawableEth[oracle];
  }
}
