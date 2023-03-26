// SPDX-License-Identifier: MIT
pragma solidity 0.8.6;

contract MockOffchainAggregator {
  event RoundIdUpdated(uint80 roundId);

  uint80 public roundId;

  function requestNewRound() external returns (uint80) {
    roundId++;
    emit RoundIdUpdated(roundId);
    return roundId;
  }
}
