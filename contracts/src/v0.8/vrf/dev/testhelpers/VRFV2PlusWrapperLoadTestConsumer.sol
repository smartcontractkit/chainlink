// SPDX-License-Identifier: MIT
pragma solidity ^0.8.6;

import {VRFV2PlusWrapperConsumerBase} from "../VRFV2PlusWrapperConsumerBase.sol";
import {ConfirmedOwner} from "../../../shared/access/ConfirmedOwner.sol";
import {ChainSpecificUtil} from "../../../ChainSpecificUtil.sol";
import {VRFV2PlusClient} from "../libraries/VRFV2PlusClient.sol";

contract VRFV2PlusWrapperLoadTestConsumer is VRFV2PlusWrapperConsumerBase, ConfirmedOwner {
  uint256 public s_responseCount;
  uint256 public s_requestCount;
  uint256 public s_averageFulfillmentInMillions = 0; // in millions for better precision
  uint256 public s_slowestFulfillment = 0;
  uint256 public s_fastestFulfillment = 999;
  uint256 public s_lastRequestId;
  // solhint-disable-next-line chainlink-solidity/prefix-storage-variables-with-s-underscore
  mapping(uint256 => uint256) internal requestHeights; // requestIds to block number when rand request was made

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
    bool native;
  }

  mapping(uint256 => RequestStatus) /* requestId */ /* requestStatus */ public s_requests;

  constructor(
    address _link,
    address _vrfV2PlusWrapper
  ) ConfirmedOwner(msg.sender) VRFV2PlusWrapperConsumerBase(_link, _vrfV2PlusWrapper) {}

  function makeRequests(
    uint32 _callbackGasLimit,
    uint16 _requestConfirmations,
    uint32 _numWords,
    uint16 _requestCount
  ) external onlyOwner {
    for (uint16 i = 0; i < _requestCount; i++) {
      bytes memory extraArgs = VRFV2PlusClient._argsToBytes(VRFV2PlusClient.ExtraArgsV1({nativePayment: false}));
      (uint256 requestId, uint256 paid) = requestRandomness(
        _callbackGasLimit,
        _requestConfirmations,
        _numWords,
        extraArgs
      );
      s_lastRequestId = requestId;

      uint256 requestBlockNumber = ChainSpecificUtil._getBlockNumber();
      s_requests[requestId] = RequestStatus({
        paid: paid,
        fulfilled: false,
        randomWords: new uint256[](0),
        requestTimestamp: block.timestamp,
        requestBlockNumber: requestBlockNumber,
        fulfilmentTimestamp: 0,
        fulfilmentBlockNumber: 0,
        native: false
      });
      s_requestCount++;
      requestHeights[requestId] = requestBlockNumber;
      emit WrapperRequestMade(requestId, paid);
    }
  }

  function makeRequestsNative(
    uint32 _callbackGasLimit,
    uint16 _requestConfirmations,
    uint32 _numWords,
    uint16 _requestCount
  ) external onlyOwner {
    for (uint16 i = 0; i < _requestCount; i++) {
      bytes memory extraArgs = VRFV2PlusClient._argsToBytes(VRFV2PlusClient.ExtraArgsV1({nativePayment: true}));
      (uint256 requestId, uint256 paid) = requestRandomnessPayInNative(
        _callbackGasLimit,
        _requestConfirmations,
        _numWords,
        extraArgs
      );
      s_lastRequestId = requestId;

      uint256 requestBlockNumber = ChainSpecificUtil._getBlockNumber();
      s_requests[requestId] = RequestStatus({
        paid: paid,
        fulfilled: false,
        randomWords: new uint256[](0),
        requestTimestamp: block.timestamp,
        requestBlockNumber: requestBlockNumber,
        fulfilmentTimestamp: 0,
        fulfilmentBlockNumber: 0,
        native: true
      });
      s_requestCount++;
      requestHeights[requestId] = requestBlockNumber;
      emit WrapperRequestMade(requestId, paid);
    }
  }

  // solhint-disable-next-line chainlink-solidity/prefix-internal-functions-with-underscore
  function fulfillRandomWords(uint256 _requestId, uint256[] memory _randomWords) internal override {
    // solhint-disable-next-line custom-errors
    require(s_requests[_requestId].paid > 0, "request not found");
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
    // solhint-disable-next-line custom-errors
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

  function reset() external {
    s_averageFulfillmentInMillions = 0; // in millions for better precision
    s_slowestFulfillment = 0;
    s_fastestFulfillment = 999;
    s_requestCount = 0;
    s_responseCount = 0;
  }

  /// @notice withdrawLink withdraws the amount specified in amount to the owner
  /// @param amount the amount to withdraw, in juels
  function withdrawLink(uint256 amount) external onlyOwner {
    s_linkToken.transfer(owner(), amount);
  }

  /// @notice withdrawNative withdraws the amount specified in amount to the owner
  /// @param amount the amount to withdraw, in wei
  function withdrawNative(uint256 amount) external onlyOwner {
    (bool success, ) = payable(owner()).call{value: amount}("");
    // solhint-disable-next-line custom-errors
    require(success, "withdrawNative failed");
  }

  receive() external payable {}
}
