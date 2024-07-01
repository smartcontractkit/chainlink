// SPDX-License-Identifier: MIT
pragma solidity 0.8.19;

contract MockFeed {
  function latestRoundData()
    external
    view
    returns (uint80 roundId, int256 answer, uint256 startedAt, uint256 updatedAt, uint80 answeredInRound)
  {
    roundId = 0;
    answer = 4120560980556678;
    startedAt = 0;
    updatedAt = block.timestamp;
    answeredInRound = 0;
  }
}
