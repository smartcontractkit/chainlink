// SPDX-License-Identifier: MIT
pragma solidity ^0.8.6;

import "../interfaces/AggregatorV3Interface.sol";

contract MockLinkNativeAggregator is AggregatorV3Interface {
  function decimals() external pure override returns (uint8) {
    return 10;
  }

  function description() external pure override returns (string memory) {
    return "Mock Feed";
  }

  function version() external pure override returns (uint256) {
    return 1;
  }

  function getRoundData(
    uint80 _roundId
  )
    external
    view
    override
    returns (uint80 roundId, int256 answer, uint256 startedAt, uint256 updatedAt, uint80 answeredInRound)
  {
    return (_roundId, 0, 0, block.timestamp, 0);
  }

  function latestRoundData()
    external
    view
    override
    returns (uint80 roundId, int256 answer, uint256 startedAt, uint256 updatedAt, uint80 answeredInRound)
  {
    return (0, 0, 0, block.timestamp, 0);
  }
}
