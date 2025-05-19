// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import {VRFCoordinatorV2Interface} from "../interfaces/VRFCoordinatorV2Interface.sol";
import {VRFConsumerBaseV2} from "../VRFConsumerBaseV2.sol";
import {ConfirmedOwner} from "../../shared/access/ConfirmedOwner.sol";
import {ChainSpecificUtil} from "../ChainSpecificUtil_v0_8_6.sol";
import {LinkTokenInterface} from "../../shared/interfaces/LinkTokenInterface.sol";

contract VRFV2OwnerTestConsumer is VRFConsumerBaseV2, ConfirmedOwner {
  VRFCoordinatorV2Interface public COORDINATOR;
  LinkTokenInterface public LINKTOKEN;
  uint64 public subId;
  uint256 public s_responseCount;
  uint256 public s_requestCount;
  uint256 public s_averageFulfillmentInMillions = 0; // in millions for better precision
  uint256 public s_slowestFulfillment = 0;
  uint256 public s_fastestFulfillment = 999;
  uint256 public s_lastRequestId;
  mapping(uint256 => uint256) internal requestHeights; // requestIds to block number when rand request was made

  event SubscriptionCreatedFundedAndConsumerAdded(uint64 subId, address consumer, uint256 amount);

  struct RequestStatus {
    bool fulfilled;
    uint256[] randomWords;
    uint requestTimestamp;
    uint fulfilmentTimestamp;
    uint256 requestBlockNumber;
    uint256 fulfilmentBlockNumber;
  }

  mapping(uint256 => RequestStatus) /* requestId */ /* requestStatus */ public s_requests;

  constructor(address _vrfCoordinator, address _link) VRFConsumerBaseV2(_vrfCoordinator) ConfirmedOwner(msg.sender) {
    COORDINATOR = VRFCoordinatorV2Interface(_vrfCoordinator);
    LINKTOKEN = LinkTokenInterface(_link);
  }

  function fulfillRandomWords(uint256 _requestId, uint256[] memory _randomWords) internal override {
    uint256 fulfilmentBlockNumber = ChainSpecificUtil._getBlockNumber();
    uint256 requestDelay = fulfilmentBlockNumber - requestHeights[_requestId];
    uint256 requestDelayInMillions = requestDelay * 1_000_000;

    if (requestDelay > s_slowestFulfillment) {
      s_slowestFulfillment = requestDelay;
    }
    s_fastestFulfillment = requestDelay < s_fastestFulfillment ? requestDelay : s_fastestFulfillment;
    s_averageFulfillmentInMillions = s_responseCount > 0
      ? (s_averageFulfillmentInMillions * s_responseCount + requestDelayInMillions) / (s_responseCount + 1)
      : requestDelayInMillions;

    s_requests[_requestId].fulfilled = true;
    s_requests[_requestId].randomWords = _randomWords;
    s_requests[_requestId].fulfilmentTimestamp = block.timestamp;
    s_requests[_requestId].fulfilmentBlockNumber = fulfilmentBlockNumber;

    s_responseCount++;
  }

  function requestRandomWords(
    uint16 _requestConfirmations,
    bytes32 _keyHash,
    uint32 _callbackGasLimit,
    uint32 _numWords,
    uint16 _requestCount,
    uint256 _subTopUpAmount
  ) external onlyOwner {
    // create a subscription, address(this) will be the owner
    subId = COORDINATOR.createSubscription();
    // add address(this) as a consumer on the subscription
    COORDINATOR.addConsumer(subId, address(this));
    topUpSubscription(_subTopUpAmount);
    emit SubscriptionCreatedFundedAndConsumerAdded(subId, address(this), _subTopUpAmount);

    for (uint16 i = 0; i < _requestCount; i++) {
      uint256 requestId = COORDINATOR.requestRandomWords(
        _keyHash,
        subId,
        _requestConfirmations,
        _callbackGasLimit,
        _numWords
      );
      s_lastRequestId = requestId;
      uint256 requestBlockNumber = ChainSpecificUtil._getBlockNumber();
      s_requests[requestId] = RequestStatus({
        randomWords: new uint256[](0),
        fulfilled: false,
        requestTimestamp: block.timestamp,
        fulfilmentTimestamp: 0,
        requestBlockNumber: requestBlockNumber,
        fulfilmentBlockNumber: 0
      });
      s_requestCount++;
      requestHeights[requestId] = requestBlockNumber;
    }

    COORDINATOR.removeConsumer(subId, address(this));
    COORDINATOR.cancelSubscription(subId, msg.sender);
  }

  function reset() external {
    s_averageFulfillmentInMillions = 0; // in millions for better precision
    s_slowestFulfillment = 0;
    s_fastestFulfillment = 999;
    s_requestCount = 0;
    s_responseCount = 0;
  }

  function getRequestStatus(
    uint256 _requestId
  )
    external
    view
    returns (
      bool fulfilled,
      uint256[] memory randomWords,
      uint requestTimestamp,
      uint fulfilmentTimestamp,
      uint256 requestBlockNumber,
      uint256 fulfilmentBlockNumber
    )
  {
    RequestStatus memory request = s_requests[_requestId];
    return (
      request.fulfilled,
      request.randomWords,
      request.requestTimestamp,
      request.fulfilmentTimestamp,
      request.requestBlockNumber,
      request.fulfilmentBlockNumber
    );
  }

  function topUpSubscription(uint256 amount) public onlyOwner {
    LINKTOKEN.transferAndCall(address(COORDINATOR), amount, abi.encode(subId));
  }
}
