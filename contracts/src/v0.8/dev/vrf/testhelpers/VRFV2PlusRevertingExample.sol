// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "../../../interfaces/LinkTokenInterface.sol";
import "../../interfaces/IVRFCoordinatorV2Plus.sol";
import "../VRFConsumerBaseV2Plus.sol";

// VRFV2RevertingExample will always revert. Used for testing only, useless in prod.
contract VRFV2PlusRevertingExample is VRFConsumerBaseV2Plus {
  uint256[] public s_randomWords;
  uint256 public s_requestId;
  IVRFCoordinatorV2Plus COORDINATOR;
  LinkTokenInterface LINKTOKEN;
  uint64 public s_subId;
  uint256 public s_gasAvailable;

  constructor(address vrfCoordinator, address link) VRFConsumerBaseV2Plus(vrfCoordinator) {
    COORDINATOR = IVRFCoordinatorV2Plus(vrfCoordinator);
    LINKTOKEN = LinkTokenInterface(link);
  }

  function fulfillRandomWords(uint256, uint256[] memory) internal override {
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
    require(s_subId != 0, "sub not set");
    // Approve the link transfer.
    LINKTOKEN.transferAndCall(address(COORDINATOR), amount, abi.encode(s_subId));
  }

  function updateSubscription(address[] memory consumers) external {
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
    s_requestId = COORDINATOR.requestRandomWords(keyHash, subId, minReqConfs, callbackGasLimit, numWords, false);
    return s_requestId;
  }
}
