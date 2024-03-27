// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import {LinkTokenInterface} from "../../shared/interfaces/LinkTokenInterface.sol";
import {VRFCoordinatorV2Interface} from "../interfaces/VRFCoordinatorV2Interface.sol";
import {VRFConsumerBaseV2} from "../VRFConsumerBaseV2.sol";

// VRFV2RevertingExample will always revert. Used for testing only, useless in prod.
contract VRFV2RevertingExample is VRFConsumerBaseV2 {
  uint256[] public s_randomWords;
  uint256 public s_requestId;
  VRFCoordinatorV2Interface internal COORDINATOR;
  LinkTokenInterface internal LINKTOKEN;
  uint64 public s_subId;
  uint256 public s_gasAvailable;

  constructor(address vrfCoordinator, address link) VRFConsumerBaseV2(vrfCoordinator) {
    COORDINATOR = VRFCoordinatorV2Interface(vrfCoordinator);
    LINKTOKEN = LinkTokenInterface(link);
  }

  function fulfillRandomWords(uint256, uint256[] memory) internal pure override {
    // solhint-disable-next-line gas-custom-errors, reason-string
    revert();
  }

  function createSubscriptionAndFund(uint96 amount) external {
    if (s_subId == 0) {
      s_subId = COORDINATOR.createSubscription();
      COORDINATOR.addConsumer(s_subId, address(this));
    }
    // Approve the link transfer.
    LINKTOKEN.transferAndCall(address(COORDINATOR), amount, abi.encode(s_subId));
  }

  function topUpSubscription(uint96 amount) external {
    // solhint-disable-next-line gas-custom-errors
    require(s_subId != 0, "sub not set");
    // Approve the link transfer.
    LINKTOKEN.transferAndCall(address(COORDINATOR), amount, abi.encode(s_subId));
  }

  function updateSubscription(address[] memory consumers) external {
    // solhint-disable-next-line gas-custom-errors
    require(s_subId != 0, "subID not set");
    for (uint256 i = 0; i < consumers.length; i++) {
      COORDINATOR.addConsumer(s_subId, consumers[i]);
    }
  }

  function requestRandomness(
    bytes32 keyHash,
    uint64 subId,
    uint16 minReqConfs,
    uint32 callbackGasLimit,
    uint32 numWords
  ) external returns (uint256) {
    s_requestId = COORDINATOR.requestRandomWords(keyHash, subId, minReqConfs, callbackGasLimit, numWords);
    return s_requestId;
  }
}
