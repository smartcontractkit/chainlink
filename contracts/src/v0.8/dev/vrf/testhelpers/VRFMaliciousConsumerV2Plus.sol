// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "../../../shared/interfaces/LinkTokenInterface.sol";
import "../../interfaces/IVRFCoordinatorV2Plus.sol";
import "../VRFConsumerBaseV2Plus.sol";

contract VRFMaliciousConsumerV2Plus is VRFConsumerBaseV2Plus {
  uint256[] public s_randomWords;
  uint256 public s_requestId;
  IVRFCoordinatorV2Plus COORDINATOR;
  LinkTokenInterface LINKTOKEN;
  uint256 public s_gasAvailable;
  uint256 s_subId;
  bytes32 s_keyHash;

  constructor(address vrfCoordinator, address link) VRFConsumerBaseV2Plus(vrfCoordinator) {
    COORDINATOR = IVRFCoordinatorV2Plus(vrfCoordinator);
    LINKTOKEN = LinkTokenInterface(link);
  }

  function fulfillRandomWords(uint256 requestId, uint256[] memory randomWords) internal override {
    s_gasAvailable = gasleft();
    s_randomWords = randomWords;
    s_requestId = requestId;
    VRFV2PlusClient.RandomWordsRequest memory req = VRFV2PlusClient.RandomWordsRequest({
      keyHash: s_keyHash,
      subId: s_subId,
      requestConfirmations: 1,
      callbackGasLimit: 200000,
      numWords: 1,
      extraArgs: "" // empty extraArgs defaults to link payment
    });
    // Should revert
    COORDINATOR.requestRandomWords(req);
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
    VRFV2PlusClient.RandomWordsRequest memory req = VRFV2PlusClient.RandomWordsRequest({
      keyHash: keyHash,
      subId: s_subId,
      requestConfirmations: 1,
      callbackGasLimit: 500000,
      numWords: 1,
      extraArgs: "" // empty extraArgs defaults to link payment
    });
    return COORDINATOR.requestRandomWords(req);
  }
}
