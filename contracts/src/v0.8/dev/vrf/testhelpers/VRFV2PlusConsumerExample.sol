// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "../../../shared/interfaces/LinkTokenInterface.sol";
import "../../interfaces/IVRFCoordinatorV2Plus.sol";
import "../VRFConsumerBaseV2Plus.sol";
import "../../../shared/access/ConfirmedOwner.sol";

/// @notice This contract is used for testing only and should not be used for production.
contract VRFV2PlusConsumerExample is ConfirmedOwner, VRFConsumerBaseV2Plus {
  LinkTokenInterface public s_linkToken;
  uint256 public s_recentRequestId;
  IVRFCoordinatorV2Plus public s_vrfCoordinatorApiV1;
  uint256 public s_subId;

  struct Response {
    bool fulfilled;
    address requester;
    uint256 requestId;
    uint256[] randomWords;
  }
  mapping(uint256 /* request id */ => Response /* response */) public s_requests;

  constructor(address vrfCoordinator, address link) VRFConsumerBaseV2Plus(vrfCoordinator) {
    s_vrfCoordinatorApiV1 = IVRFCoordinatorV2Plus(vrfCoordinator);
    s_linkToken = LinkTokenInterface(link);
  }

  function getRandomness(uint256 requestId, uint256 idx) public view returns (uint256 randomWord) {
    Response memory resp = s_requests[requestId];
    require(resp.requestId != 0, "request ID is incorrect");
    return resp.randomWords[idx];
  }

  function subscribe() internal returns (uint256) {
    if (s_subId == 0) {
      s_subId = s_vrfCoordinatorApiV1.createSubscription();
      s_vrfCoordinatorApiV1.addConsumer(s_subId, address(this));
    }
    return s_subId;
  }

  function createSubscriptionAndFundNative() external payable {
    subscribe();
    s_vrfCoordinatorApiV1.fundSubscriptionWithNative{value: msg.value}(s_subId);
  }

  function createSubscriptionAndFund(uint96 amount) external {
    subscribe();
    // Approve the link transfer.
    s_linkToken.transferAndCall(address(s_vrfCoordinator), amount, abi.encode(s_subId));
  }

  function topUpSubscription(uint96 amount) external {
    require(s_subId != 0, "sub not set");
    s_linkToken.transferAndCall(address(s_vrfCoordinator), amount, abi.encode(s_subId));
  }

  function topUpSubscriptionNative() external payable {
    require(s_subId != 0, "sub not set");
    s_vrfCoordinatorApiV1.fundSubscriptionWithNative{value: msg.value}(s_subId);
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
    VRFV2PlusClient.RandomWordsRequest memory req = VRFV2PlusClient.RandomWordsRequest({
      keyHash: keyHash,
      subId: s_subId,
      requestConfirmations: requestConfirmations,
      callbackGasLimit: callbackGasLimit,
      numWords: numWords,
      extraArgs: VRFV2PlusClient._argsToBytes(VRFV2PlusClient.ExtraArgsV1({nativePayment: nativePayment}))
    });
    uint256 requestId = s_vrfCoordinator.requestRandomWords(req);
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
      s_vrfCoordinatorApiV1.addConsumer(s_subId, consumers[i]);
    }
  }

  function setSubId(uint256 subId) external {
    s_subId = subId;
  }
}
