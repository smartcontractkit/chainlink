// SPDX-License-Identifier: MIT
pragma solidity 0.8.6;

import {IAggregatorProxy} from '../HeartbeatRequester.sol';

contract MockAggregatorProxy is IAggregatorProxy {
  address internal s_aggregator;

  constructor(address _aggregator) {
    s_aggregator = _aggregator;
  }

  function updateAggregator(address _aggregator) external {
    s_aggregator = _aggregator;
  }

  function aggregator() external override view returns (address) {
    return s_aggregator;
  }
}
