// SPDX-License-Identifier: MIT
pragma solidity ^0.8.6;

interface IOffchainAggregator {
  function requestNewRound() external returns (uint80);
}
