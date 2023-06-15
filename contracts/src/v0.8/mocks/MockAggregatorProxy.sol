// SPDX-License-Identifier: MIT
pragma solidity 0.8.6;

import '../HeartbeatRequester.sol';

contract MockAggregatorProxy is IAggregatorProxy, IOffchainAggregator{
  address internal s_aggregator;

  constructor(address aggregator) {
    s_aggregator = aggregator;
  }

  function updateAggregator(address aggregator) external {
    s_aggregator = aggregator;
  }

  function aggregator() external override view returns (address) {
    return s_aggregator;
  }

  function requestNewRound() external override returns (uint80){
    // do we need the actual logic of requestNewRound?
    // or is dummy requestNewRound ok?
    return 1;
  }
}
