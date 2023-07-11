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
  uint64 public s_subId;
  uint256 public s_recentRequestId;

  struct Response {
    bool fulfilled;
    address requester;
    uint256 requestId;
    uint256[] randomWords;
  }
  mapping(uint256 /* request id */ => Response /* response */) public s_requests;

  constructor(address vrfCoordinator, address link) ConfirmedOwner(msg.sender) VRFConsumerBaseV2Plus(vrfCoordinator) {
    s_vrfCoordinator = IVRFCoordinatorV2Plus(vrfCoordinator);
    s_linkToken = LinkTokenInterface(link);
  }

  function getRandomness(uint256 requestId, uint256 idx) public view returns (uint256 randomWord) {
    Response memory resp = s_requests[requestId];
    require(resp.requestId != 0, "request ID is incorrect");
    return resp.randomWords[idx];
  }

  function createSubscriptionAndFund(uint96 amount) external {
    if (s_subId == 0) {
      s_subId = s_vrfCoordinator.createSubscription();
      s_vrfCoordinator.addConsumer(s_subId, address(this));
      _setSubOwner(address(this));
    }
    // Approve the link transfer.
    s_linkToken.transferAndCall(address(s_vrfCoordinator), amount, abi.encode(s_subId));
  }

  function setSubOwner(address subOwner) external {
    _setSubOwner(subOwner);
  }

  function topUpSubscription(uint96 amount) external {
    s_linkToken.transferAndCall(address(s_vrfCoordinator), amount, abi.encode(s_subId));
  }

  function fulfillRandomWords(uint256 requestId, uint256[] memory randomWords) internal override {
    require(requestId == s_recentRequestId, "request ID is incorrect");
    s_requests[requestId].randomWords = randomWords;
    s_requests[requestId].fulfilled = true;
  }

  function requestRandomWords(
    uint32 callbackGasLimit,
    uint16 requestConfirmations,
    uint32 numWords,
    bytes32 keyHash,
    bool nativePayment
  ) external {
    uint256 requestId = s_vrfCoordinator.requestRandomWords(
      keyHash,
      s_subId,
      requestConfirmations,
      callbackGasLimit,
      numWords,
      nativePayment
    );
    Response memory resp = Response({
      requestId: requestId,
      randomWords: new uint256[](0),
      fulfilled: false,
      requester: msg.sender
    });
    s_requests[requestId] = resp;
    s_recentRequestId = requestId;
  }

  function updateSubscription(address[] memory consumers) external {
    require(s_subId != 0, "subID not set");
    for (uint256 i = 0; i < consumers.length; i++) {
      s_vrfCoordinator.addConsumer(s_subId, consumers[i]);
    }
  }
}
