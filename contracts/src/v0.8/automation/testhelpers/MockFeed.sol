// SPDX-License-Identifier: MIT
pragma solidity 0.8.19;

contract MockFeed {
  function latestRoundData()
    external
    view
    returns (uint80 roundId, int256 answer, uint256 startedAt, uint256 updatedAt, uint80 answeredInRound)
  {
    roundId = 0;
    answer = 0;
    startedAt = 0;
    updatedAt = 0;
    answeredInRound = 0;
  }
}
