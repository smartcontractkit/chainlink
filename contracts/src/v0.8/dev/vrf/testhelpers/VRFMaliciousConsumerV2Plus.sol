// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "../../../interfaces/LinkTokenInterface.sol";
import "../../interfaces/IVRFCoordinatorV2Plus.sol";
import "../VRFConsumerBaseV2Plus.sol";

contract VRFMaliciousConsumerV2Plus is VRFConsumerBaseV2Plus {
  uint256[] public s_randomWords;
  uint256 public s_requestId;
  IVRFCoordinatorV2Plus COORDINATOR;
  LinkTokenInterface LINKTOKEN;
  uint64 public s_subId;
  uint256 public s_gasAvailable;
  bytes32 s_keyHash;

  constructor(address vrfCoordinator, address link) VRFConsumerBaseV2Plus(vrfCoordinator) {
    COORDINATOR = IVRFCoordinatorV2Plus(vrfCoordinator);
    LINKTOKEN = LinkTokenInterface(link);
  }

  function fulfillRandomWords(uint256 requestId, uint256[] memory randomWords) internal override {
    s_gasAvailable = gasleft();
    s_randomWords = randomWords;
    s_requestId = requestId;
    // Should revert
    COORDINATOR.requestRandomWords(s_keyHash, s_subId, 1, 200000, 1, false);
  }

  function createSubscriptionAndFund(uint96 amount) external {
    if (s_subId == 0) {
      s_subId = COORDINATOR.createSubscription();
      COORDINATOR.addConsumer(s_subId, address(this));
    }
    // Approve the link transfer.
    LINKTOKEN.transferAndCall(address(COORDINATOR), amount, abi.encode(s_subId));
  }

  function updateSubscription(address[] memory consumers) external {
    require(s_subId != 0, "subID not set");
    for (uint256 i = 0; i < consumers.length; i++) {
      COORDINATOR.addConsumer(s_subId, consumers[i]);
    }
  }

  function requestRandomness(bytes32 keyHash) external returns (uint256) {
    s_keyHash = keyHash;
    return COORDINATOR.requestRandomWords(keyHash, s_subId, 1, 500000, 1, false);
  }
}
