// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import {ChainSpecificUtil} from "../../../ChainSpecificUtil.sol";
import {VRFConsumerBaseV2Plus} from "../VRFConsumerBaseV2Plus.sol";
import {VRFV2PlusClient} from "../libraries/VRFV2PlusClient.sol";

/**
 * @title The VRFLoadTestExternalSubOwner contract.
 * @notice Allows making many VRF V2 randomness requests in a single transaction for load testing.
 */
contract VRFV2PlusLoadTestWithMetrics is VRFConsumerBaseV2Plus {
  uint256 public s_responseCount;
  uint256 public s_requestCount;
  uint256 public s_averageResponseTimeInBlocksMillions = 0; // in millions for better precision
  uint256 public s_slowestResponseTimeInBlocks = 0;
  uint256 public s_fastestResponseTimeInBlocks = 999;
  uint256 public s_slowestResponseTimeInSeconds = 0;
  uint256 public s_fastestResponseTimeInSeconds = 999;
  uint256 public s_averageResponseTimeInSecondsMillions = 0;

  uint256 public s_lastRequestId;

  uint32[] public s_requestBlockTimes;

  struct RequestStatus {
    bool fulfilled;
    uint256[] randomWords;
    uint256 requestTimestamp;
    uint256 fulfilmentTimestamp;
    uint256 requestBlockNumber;
    uint256 fulfilmentBlockNumber;
  }

  mapping(uint256 => RequestStatus) /* requestId */ /* requestStatus */ public s_requests;

  constructor(address _vrfCoordinator) VRFConsumerBaseV2Plus(_vrfCoordinator) {}

  // solhint-disable-next-line chainlink-solidity/prefix-internal-functions-with-underscore
  function fulfillRandomWords(uint256 _requestId, uint256[] calldata _randomWords) internal override {
    s_requests[_requestId].fulfilled = true;
    s_requests[_requestId].randomWords = _randomWords;
    s_requests[_requestId].fulfilmentTimestamp = block.timestamp;
    s_requests[_requestId].fulfilmentBlockNumber = ChainSpecificUtil._getBlockNumber();

    uint256 responseTimeInBlocks = s_requests[_requestId].fulfilmentBlockNumber -
      s_requests[_requestId].requestBlockNumber;
    uint256 responseTimeInSeconds = s_requests[_requestId].fulfilmentTimestamp -
      s_requests[_requestId].requestTimestamp;

    (
      s_slowestResponseTimeInBlocks,
      s_fastestResponseTimeInBlocks,
      s_averageResponseTimeInBlocksMillions
    ) = _calculateMetrics(
      responseTimeInBlocks,
      s_fastestResponseTimeInBlocks,
      s_slowestResponseTimeInBlocks,
      s_averageResponseTimeInBlocksMillions,
      s_responseCount
    );
    (
      s_slowestResponseTimeInSeconds,
      s_fastestResponseTimeInSeconds,
      s_averageResponseTimeInSecondsMillions
    ) = _calculateMetrics(
      responseTimeInSeconds,
      s_fastestResponseTimeInSeconds,
      s_slowestResponseTimeInSeconds,
      s_averageResponseTimeInSecondsMillions,
      s_responseCount
    );

    s_responseCount++;

    s_requestBlockTimes.push(uint32(responseTimeInBlocks));
  }

  function requestRandomWords(
    uint256 _subId,
    uint16 _requestConfirmations,
    bytes32 _keyHash,
    uint32 _callbackGasLimit,
    bool _nativePayment,
    uint32 _numWords,
    uint16 _requestCount
  ) external onlyOwner {
    for (uint16 i = 0; i < _requestCount; i++) {
      VRFV2PlusClient.RandomWordsRequest memory req = VRFV2PlusClient.RandomWordsRequest({
        keyHash: _keyHash,
        subId: _subId,
        requestConfirmations: _requestConfirmations,
        callbackGasLimit: _callbackGasLimit,
        numWords: _numWords,
        extraArgs: VRFV2PlusClient._argsToBytes(VRFV2PlusClient.ExtraArgsV1({nativePayment: _nativePayment}))
      });
      // Will revert if subscription is not funded.
      uint256 requestId = s_vrfCoordinator.requestRandomWords(req);

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
    }
  }

  function reset() external {
    s_averageResponseTimeInBlocksMillions = 0; // in millions for better precision
    s_slowestResponseTimeInBlocks = 0;
    s_fastestResponseTimeInBlocks = 999;
    s_averageResponseTimeInSecondsMillions = 0; // in millions for better precision
    s_slowestResponseTimeInSeconds = 0;
    s_fastestResponseTimeInSeconds = 999;
    s_requestCount = 0;
    s_responseCount = 0;
    delete s_requestBlockTimes;
  }

  function getRequestStatus(
    uint256 _requestId
  )
    external
    view
    returns (
      bool fulfilled,
      uint256[] memory randomWords,
      uint256 requestTimestamp,
      uint256 fulfilmentTimestamp,
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

  function _calculateMetrics(
    uint256 _responseTime,
    uint256 _fastestResponseTime,
    uint256 _slowestResponseTime,
    uint256 _averageInMillions,
    uint256 _responseCount
  ) internal pure returns (uint256 slowest, uint256 fastest, uint256 average) {
    uint256 _requestDelayInMillions = _responseTime * 1_000_000;
    if (_responseTime > _slowestResponseTime) {
      _slowestResponseTime = _responseTime;
    }
    _fastestResponseTime = _responseTime < _fastestResponseTime ? _responseTime : _fastestResponseTime;
    uint256 averageInMillions = _responseCount > 0
      ? (_averageInMillions * _responseCount + _requestDelayInMillions) / (_responseCount + 1)
      : _requestDelayInMillions;

    return (_slowestResponseTime, _fastestResponseTime, averageInMillions);
  }

  function getRequestBlockTimes(uint256 offset, uint256 quantity) external view returns (uint32[] memory) {
    uint256 end = offset + quantity;
    if (end > s_requestBlockTimes.length) {
      end = s_requestBlockTimes.length;
    }

    uint32[] memory blockTimes = new uint32[](end - offset);
    for (uint256 i = offset; i < end; i++) {
      blockTimes[i - offset] = s_requestBlockTimes[i];
    }

    return blockTimes;
  }
}
