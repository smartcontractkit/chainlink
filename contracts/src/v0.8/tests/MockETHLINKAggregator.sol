// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "../shared/interfaces/AggregatorV3Interface.sol";

contract MockETHLINKAggregator is AggregatorV3Interface {
  int256 public answer;

  constructor(int256 _answer) public {
    answer = _answer;
  }

  function decimals() external view override returns (uint8) {
    return 18;
  }

  function description() external view override returns (string memory) {
    return "MockETHLINKAggregator";
  }

  function version() external view override returns (uint256) {
    return 1;
  }

  function getRoundData(
    uint80 _roundId
  )
    external
    view
    override
    returns (uint80 roundId, int256 ans, uint256 startedAt, uint256 updatedAt, uint80 answeredInRound)
  {
    return (1, answer, block.timestamp, block.timestamp, 1);
  }

  function latestRoundData()
    external
    view
    override
    returns (uint80 roundId, int256 ans, uint256 startedAt, uint256 updatedAt, uint80 answeredInRound)
  {
    return (1, answer, block.timestamp, block.timestamp, 1);
  }
}
