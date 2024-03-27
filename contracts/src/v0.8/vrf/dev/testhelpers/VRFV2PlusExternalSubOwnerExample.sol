// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import {LinkTokenInterface} from "../../../shared/interfaces/LinkTokenInterface.sol";
import {VRFConsumerBaseV2Plus} from "../VRFConsumerBaseV2Plus.sol";
import {VRFV2PlusClient} from "../libraries/VRFV2PlusClient.sol";

/// @notice This contract is used for testing only and should not be used for production.
contract VRFV2PlusExternalSubOwnerExample is VRFConsumerBaseV2Plus {
  // solhint-disable-next-line chainlink-solidity/prefix-storage-variables-with-s-underscore
  LinkTokenInterface internal LINKTOKEN;

  uint256[] public s_randomWords;
  uint256 public s_requestId;
  address internal s_owner;

  constructor(address vrfCoordinator, address link) VRFConsumerBaseV2Plus(vrfCoordinator) {
    LINKTOKEN = LinkTokenInterface(link);
    s_owner = msg.sender;
  }

  // solhint-disable-next-line chainlink-solidity/prefix-internal-functions-with-underscore
  function fulfillRandomWords(uint256 requestId, uint256[] memory randomWords) internal override {
    // solhint-disable-next-line gas-custom-errors
    require(requestId == s_requestId, "request ID is incorrect");
    s_randomWords = randomWords;
  }

  function requestRandomWords(
    uint256 subId,
    uint32 callbackGasLimit,
    uint16 requestConfirmations,
    uint32 numWords,
    bytes32 keyHash,
    bool nativePayment
  ) external onlyOwner {
    VRFV2PlusClient.RandomWordsRequest memory req = VRFV2PlusClient.RandomWordsRequest({
      keyHash: keyHash,
      subId: subId,
      requestConfirmations: requestConfirmations,
      callbackGasLimit: callbackGasLimit,
      numWords: numWords,
      extraArgs: VRFV2PlusClient._argsToBytes(VRFV2PlusClient.ExtraArgsV1({nativePayment: nativePayment}))
    });
    // Will revert if subscription is not funded.
    s_requestId = s_vrfCoordinator.requestRandomWords(req);
  }
}
