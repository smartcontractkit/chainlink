// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import {LinkTokenInterface} from "../../../shared/interfaces/LinkTokenInterface.sol";
import {VRFConsumerBaseV2Plus} from "../VRFConsumerBaseV2Plus.sol";
import {VRFV2PlusClient} from "../libraries/VRFV2PlusClient.sol";

contract VRFMaliciousConsumerV2Plus is VRFConsumerBaseV2Plus {
  uint256[] public s_randomWords;
  uint256 public s_requestId;
  // solhint-disable-next-line chainlink-solidity/prefix-storage-variables-with-s-underscore
  LinkTokenInterface internal LINKTOKEN;
  uint256 public s_gasAvailable;
  uint256 internal s_subId;
  bytes32 internal s_keyHash;

  constructor(address vrfCoordinator, address link) VRFConsumerBaseV2Plus(vrfCoordinator) {
    LINKTOKEN = LinkTokenInterface(link);
  }

  // solhint-disable-next-line chainlink-solidity/prefix-internal-functions-with-underscore
  function fulfillRandomWords(uint256 requestId, uint256[] calldata randomWords) internal override {
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
    s_vrfCoordinator.requestRandomWords(req);
  }

  function createSubscriptionAndFund(uint96 amount) external {
    if (s_subId == 0) {
      s_subId = s_vrfCoordinator.createSubscription();
      s_vrfCoordinator.addConsumer(s_subId, address(this));
    }
    // Approve the link transfer.
    LINKTOKEN.transferAndCall(address(s_vrfCoordinator), amount, abi.encode(s_subId));
  }

  function updateSubscription(address[] memory consumers) external {
    // solhint-disable-next-line gas-custom-errors
    require(s_subId != 0, "subID not set");
    for (uint256 i = 0; i < consumers.length; i++) {
      s_vrfCoordinator.addConsumer(s_subId, consumers[i]);
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
    return s_vrfCoordinator.requestRandomWords(req);
  }
}
