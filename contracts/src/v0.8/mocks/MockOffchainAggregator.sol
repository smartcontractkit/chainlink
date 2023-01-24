// SPDX-License-Identifier: MIT
pragma solidity 0.8.6;

import "../dev/interfaces/IOffchainAggregator.sol";

contract MockOffchainAggregator is IOffchainAggregator {
  event RoundIdUpdated(uint80 roundId);

  uint80 public roundId;

  function requestNewRound() external override returns (uint80) {
    roundId++;
    emit RoundIdUpdated(roundId);
    return roundId;
  }
}
