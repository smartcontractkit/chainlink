// SPDX-License-Identifier: MIT
pragma solidity ^0.8.6;

import {VRFV2WrapperConsumerBase} from "../VRFV2WrapperConsumerBase.sol";
import {ConfirmedOwner} from "../../shared/access/ConfirmedOwner.sol";
import {ChainSpecificUtil} from "../../ChainSpecificUtil_v0_8_6.sol";
import {VRFV2WrapperInterface} from "../interfaces/VRFV2WrapperInterface.sol";

contract VRFV2WrapperLoadTestConsumer is VRFV2WrapperConsumerBase, ConfirmedOwner {
  VRFV2WrapperInterface public immutable i_vrfV2Wrapper;
  uint256 public s_responseCount;
  uint256 public s_requestCount;
  uint256 public s_averageFulfillmentInMillions = 0; // in millions for better precision
  uint256 public s_slowestFulfillment = 0;
  uint256 public s_fastestFulfillment = 999;
  uint256 public s_lastRequestId;
  // solhint-disable-next-line chainlink-solidity/prefix-storage-variables-with-s-underscore
  mapping(uint256 => uint256) internal requestHeights; // requestIds to block number when rand request was made
  mapping(uint256 => RequestStatus) /* requestId */ /* requestStatus */ public s_requests;

  event WrappedRequestFulfilled(uint256 requestId, uint256[] randomWords, uint256 payment);
  event WrapperRequestMade(uint256 indexed requestId, uint256 paid);

  struct RequestStatus {
    uint256 paid;
    bool fulfilled;
    uint256[] randomWords;
    uint256 requestTimestamp;
    uint256 fulfilmentTimestamp;
    uint256 requestBlockNumber;
    uint256 fulfilmentBlockNumber;
  }

  constructor(
    address _link,
    address _vrfV2Wrapper
  ) ConfirmedOwner(msg.sender) VRFV2WrapperConsumerBase(_link, _vrfV2Wrapper) {
    i_vrfV2Wrapper = VRFV2WrapperInterface(_vrfV2Wrapper);
  }

  function makeRequests(
    uint32 _callbackGasLimit,
    uint16 _requestConfirmations,
    uint32 _numWords,
    uint16 _requestCount
  ) external onlyOwner {
    for (uint16 i = 0; i < _requestCount; i++) {
      uint256 requestId = requestRandomness(_callbackGasLimit, _requestConfirmations, _numWords);
      s_lastRequestId = requestId;
      uint256 requestBlockNumber = ChainSpecificUtil._getBlockNumber();
      uint256 paid = VRF_V2_WRAPPER.calculateRequestPrice(_callbackGasLimit);
      s_requests[requestId] = RequestStatus({
        paid: paid,
        fulfilled: false,
        randomWords: new uint256[](0),
        requestTimestamp: block.timestamp,
        fulfilmentTimestamp: 0,
        requestBlockNumber: requestBlockNumber,
        fulfilmentBlockNumber: 0
      });
      s_requestCount++;
      requestHeights[requestId] = requestBlockNumber;
      emit WrapperRequestMade(requestId, paid);
    }
  }

  function fulfillRandomWords(uint256 _requestId, uint256[] memory _randomWords) internal override {
    // solhint-disable-next-line gas-custom-errors
    require(s_requests[_requestId].paid > 0, "request not found");
    uint256 fulfilmentBlockNumber = ChainSpecificUtil._getBlockNumber();
    uint256 requestDelay = fulfilmentBlockNumber - requestHeights[_requestId];
    uint256 requestDelayInMillions = requestDelay * 1_000_000;

    if (requestDelay > s_slowestFulfillment) {
      s_slowestFulfillment = requestDelay;
    }
    if (requestDelay < s_fastestFulfillment) {
      s_fastestFulfillment = requestDelay;
    }
    s_averageFulfillmentInMillions = s_responseCount > 0
      ? (s_averageFulfillmentInMillions * s_responseCount + requestDelayInMillions) / (s_responseCount + 1)
      : requestDelayInMillions;

    s_responseCount++;
    s_requests[_requestId].fulfilled = true;
    s_requests[_requestId].randomWords = _randomWords;
    s_requests[_requestId].fulfilmentTimestamp = block.timestamp;
    s_requests[_requestId].fulfilmentBlockNumber = fulfilmentBlockNumber;

    emit WrappedRequestFulfilled(_requestId, _randomWords, s_requests[_requestId].paid);
  }

  function getRequestStatus(
    uint256 _requestId
  )
    external
    view
    returns (
      uint256 paid,
      bool fulfilled,
      uint256[] memory randomWords,
      uint256 requestTimestamp,
      uint256 fulfilmentTimestamp,
      uint256 requestBlockNumber,
      uint256 fulfilmentBlockNumber
    )
  {
    // solhint-disable-next-line gas-custom-errors
    require(s_requests[_requestId].paid > 0, "request not found");
    RequestStatus memory request = s_requests[_requestId];
    return (
      request.paid,
      request.fulfilled,
      request.randomWords,
      request.requestTimestamp,
      request.fulfilmentTimestamp,
      request.requestBlockNumber,
      request.fulfilmentBlockNumber
    );
  }

  /// @notice withdrawLink withdraws the amount specified in amount to the owner
  /// @param amount the amount to withdraw, in juels
  function withdrawLink(uint256 amount) external onlyOwner {
    LINK.transfer(owner(), amount);
  }

  function reset() external {
    s_averageFulfillmentInMillions = 0;
    s_slowestFulfillment = 0;
    s_fastestFulfillment = 999;
    s_requestCount = 0;
    s_responseCount = 0;
  }

  receive() external payable {}
}
