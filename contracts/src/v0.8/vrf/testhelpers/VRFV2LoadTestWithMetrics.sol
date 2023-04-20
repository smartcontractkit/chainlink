// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "../../interfaces/LinkTokenInterface.sol";
import "../../interfaces/VRFCoordinatorV2Interface.sol";
import "../VRFConsumerBaseV2.sol";
import "../../ConfirmedOwner.sol";

/**
 * @title The VRFLoadTestExternalSubOwner contract.
 * @notice Allows making many VRF V2 randomness requests in a single transaction for load testing.
 */
contract VRFV2LoadTestWithMetrics is VRFConsumerBaseV2, ConfirmedOwner {
  VRFCoordinatorV2Interface public immutable COORDINATOR;
  LinkTokenInterface public immutable LINK;

  uint256 public s_responseCount;
  uint256 public s_requestCount;
  uint256 public s_averageFulfillmentInMillions = 0; // in millions for better precision
  uint256 public s_slowestFulfillment = 0;
  uint256 public s_fastestFulfillment = 999;
  mapping(uint256 => uint256) requestHeights; // requestIds to block number when rand request was made

  constructor(address _vrfCoordinator, address _link) VRFConsumerBaseV2(_vrfCoordinator) ConfirmedOwner(msg.sender) {
    COORDINATOR = VRFCoordinatorV2Interface(_vrfCoordinator);
    LINK = LinkTokenInterface(_link);
  }

  function fulfillRandomWords(uint256 reqId, uint256[] memory) internal override {
    uint256 requestDelay = block.number - requestHeights[reqId];
    uint256 requestDelayInMillions = requestDelay * 1_000_000;

    if (requestDelay > s_slowestFulfillment) {
      s_slowestFulfillment = requestDelay;
    }
    s_fastestFulfillment = requestDelay < s_fastestFulfillment ? requestDelay : s_fastestFulfillment;
    s_averageFulfillmentInMillions = s_responseCount > 0
      ? (s_averageFulfillmentInMillions * s_responseCount + requestDelayInMillions) / (s_responseCount + 1)
      : requestDelayInMillions;

    s_responseCount++;
  }

  function requestRandomWords(
    uint64 _subId,
    uint16 _requestConfirmations,
    bytes32 _keyHash,
    uint16 _requestCount
  ) external onlyOwner {
    for (uint16 i = 0; i < _requestCount; i++) {
      uint256 reqId = COORDINATOR.requestRandomWords(_keyHash, _subId, _requestConfirmations, 300_000, 1);
      s_requestCount++;
      requestHeights[reqId] = block.number;
    }
  }

  function reset() external {
    s_averageFulfillmentInMillions = 0; // in millions for better precision
    s_slowestFulfillment = 0;
    s_fastestFulfillment = 999;
    s_requestCount = 0;
    s_responseCount = 0;
  }
}
