// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import {AggregatorV2V3Interface} from "../../../shared/interfaces/AggregatorV2V3Interface.sol";

contract MockAggregatorV2V3 is AggregatorV2V3Interface {
  function latestAnswer() external pure returns (int256) {
    return 0;
  }

  function latestTimestamp() external pure returns (uint256) {
    return 0;
  }

  function latestRound() external pure returns (uint256) {
    return 0;
  }

  function getAnswer(uint256) external pure returns (int256) {
    return 0;
  }

  function getTimestamp(uint256 roundId) external pure returns (uint256) {
    return roundId;
  }

  function decimals() external pure returns (uint8) {
    return 0;
  }

  function description() external pure returns (string memory) {
    return "";
  }

  function version() external pure returns (uint256) {
    return 0;
  }

  function getRoundData(
    uint80
  )
    external
    pure
    returns (uint80 roundId, int256 answer, uint256 startedAt, uint256 updatedAt, uint80 answeredInRound)
  {
    return (0, 0, 0, 0, 0);
  }

  function latestRoundData()
    external
    pure
    returns (uint80 roundId, int256 answer, uint256 startedAt, uint256 updatedAt, uint80 answeredInRound)
  {
    return (73786976294838220258, 96800000000, 163826896, 1638268960, 73786976294838220258);
  }
}
