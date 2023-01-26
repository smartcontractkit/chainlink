// SPDX-License-Identifier: MIT
pragma solidity 0.8.6;

contract MockAggregatorProxy {
  address internal s_aggregator;

  constructor(address aggregator) {
    s_aggregator = aggregator;
  }

  function updateAggregator(address aggregator) external {
    s_aggregator = aggregator;
  }

  function aggregator() external view returns (address) {
    return s_aggregator;
  }
}
