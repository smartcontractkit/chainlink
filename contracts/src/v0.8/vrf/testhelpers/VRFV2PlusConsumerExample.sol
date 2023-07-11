// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "../../interfaces/LinkTokenInterface.sol";
import "../../interfaces/IVRFCoordinatorV2Plus.sol";
import "../VRFConsumerBaseV2Plus.sol";
import "../../ConfirmedOwner.sol";

/// @notice Example VRF V2Plus consumer which passes costs to the end user.
contract VRFV2PlusConsumerExample is ConfirmedOwner, VRFConsumerBaseV2Plus {
  IVRFCoordinatorV2Plus public s_vrfCoordinator;
  LinkTokenInterface public s_linkToken;

  struct Response {
    bool fulfilled;
    address requester;
    uint256 requestId;
    uint256[] randomWords;
  }
  mapping(uint256 /* request id */ => Response /* response */) public s_requests;

  constructor(address vrfCoordinator, address link, address subOwner) ConfirmedOwner(msg.sender) VRFConsumerBaseV2Plus(vrfCoordinator, subOwner) {
    s_vrfCoordinator = IVRFCoordinatorV2Plus(vrfCoordinator);
    s_linkToken = LinkTokenInterface(link);
  }

  function fulfillRandomWords(uint256 requestId, uint256[] memory randomWords) internal override {
    Response memory resp = s_requests[requestId];
    require(resp.requestId != 0, "request ID is incorrect");
    s_requests[requestId].randomWords = randomWords;
    s_requests[requestId].fulfilled = true;
  }

  function requestRandomWords(
    uint64 subId,
    uint32 callbackGasLimit,
    uint16 requestConfirmations,
    uint32 numWords,
    bytes32 keyHash,
    bool nativePayment
  ) external {
    uint256 requestId = s_vrfCoordinator.requestRandomWords(keyHash, subId, requestConfirmations, callbackGasLimit, numWords, nativePayment);
    Response memory resp = Response({
      requestId: requestId,
      randomWords: new uint256[](0),
      fulfilled: false,
      requester: msg.sender
    });
    s_requests[requestId] = resp;
  }
}
