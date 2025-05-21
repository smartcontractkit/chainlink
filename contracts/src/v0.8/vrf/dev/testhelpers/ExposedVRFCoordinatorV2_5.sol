// SPDX-License-Identifier: MIT
pragma solidity ^0.8.4;

import {VRFCoordinatorV2_5} from "../VRFCoordinatorV2_5.sol";
import {VRFTypes} from "../../VRFTypes.sol";
import {EnumerableSet} from "../../../vendor/openzeppelin-solidity/v4.7.3/contracts/utils/structs/EnumerableSet.sol";

contract ExposedVRFCoordinatorV2_5 is VRFCoordinatorV2_5 {
  using EnumerableSet for EnumerableSet.UintSet;

  constructor(address blockhashStore) VRFCoordinatorV2_5(blockhashStore) {}

  function computeRequestIdExternal(
    bytes32 keyHash,
    address sender,
    uint256 subId,
    uint64 nonce
  ) external pure returns (uint256, uint256) {
    return _computeRequestId(keyHash, sender, subId, nonce);
  }

  function isTargetRegisteredExternal(address target) external view returns (bool) {
    return _isTargetRegistered(target);
  }

  function getRandomnessFromProofExternal(
    Proof calldata proof,
    VRFTypes.RequestCommitmentV2Plus calldata rc
  ) external view returns (Output memory) {
    return _getRandomnessFromProof(proof, rc);
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

  function setWithdrawableTokensTestingOnlyXXX(uint96 newBalance) external {
    s_withdrawableTokens = newBalance;
  }

  function getWithdrawableTokensTestingOnlyXXX() external view returns (uint96) {
    return s_withdrawableTokens;
  }

  function setWithdrawableNativeTestingOnlyXXX(uint96 newBalance) external {
    s_withdrawableNative = newBalance;
  }

  function getWithdrawableNativeTestingOnlyXXX() external view returns (uint96) {
    return s_withdrawableNative;
  }

  function calculatePaymentAmount(
    uint256 startGas,
    uint256 weiPerUnitGas,
    bool nativePayment,
    bool onlyPremium
  ) external returns (uint96, bool) {
    return _calculatePaymentAmount(startGas, weiPerUnitGas, nativePayment, onlyPremium);
  }
}
