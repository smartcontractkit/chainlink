// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import {IAggregatorV2V3} from "../interfaces/IAggregatorV2V3.sol";

contract FeedConsumer {
  IAggregatorV2V3 public immutable AGGREGATOR;

  constructor(address feedAddress) {
    AGGREGATOR = IAggregatorV2V3(feedAddress);
  }

  function latestAnswer() external view returns (int256 answer) {
    return AGGREGATOR.latestAnswer();
  }

  function latestTimestamp() external view returns (uint256) {
    return AGGREGATOR.latestTimestamp();
  }

  function latestRound() external view returns (uint256) {
    return AGGREGATOR.latestRound();
  }

  function getAnswer(uint256 roundId) external view returns (int256) {
    return AGGREGATOR.getAnswer(roundId);
  }

  function getTimestamp(uint256 roundId) external view returns (uint256) {
    return AGGREGATOR.getTimestamp(roundId);
  }

  function decimals() external view returns (uint8) {
    return AGGREGATOR.decimals();
  }

  function description() external view returns (string memory) {
    return AGGREGATOR.description();
  }

  function version() external view returns (uint256) {
    return AGGREGATOR.version();
  }

  function getRoundData(uint80 _roundId)
    external
    view
    returns (
      uint80 roundId,
      int256 answer,
      uint256 startedAt,
      uint256 updatedAt,
      uint80 answeredInRound
    )
  {
    return AGGREGATOR.getRoundData(_roundId);
  }

  function latestRoundData()
    external
    view
    returns (
      uint80 roundId,
      int256 answer,
      uint256 startedAt,
      uint256 updatedAt,
      uint80 answeredInRound
    )
  {
    return AGGREGATOR.latestRoundData();
  }
}
