// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import {AggregatorV2V3Interface} from "../interfaces/AggregatorV2V3Interface.sol";

contract FeedConsumer {
  AggregatorV2V3Interface public immutable AGGREGATOR;

  constructor(address feedAddress) {
    AGGREGATOR = AggregatorV2V3Interface(feedAddress);
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

  function getRoundData(
    uint80 _roundId
  )
    external
    view
    returns (uint80 roundId, int256 answer, uint256 startedAt, uint256 updatedAt, uint80 answeredInRound)
  {
    return AGGREGATOR.getRoundData(_roundId);
  }

  function latestRoundData()
    external
    view
    returns (uint80 roundId, int256 answer, uint256 startedAt, uint256 updatedAt, uint80 answeredInRound)
  {
    return AGGREGATOR.latestRoundData();
  }
}
