// SPDX-License-Identifier: MIT
pragma solidity 0.8.6;

import "../dev/interfaces/IAggregatorProxy.sol";

contract MockAggregatorProxy is IAggregatorProxy {
  address internal s_aggregator;

  constructor(address aggregator) {
    s_aggregator = aggregator;
  }

  function updateAggregator(address aggregator) external {
    s_aggregator = aggregator;
  }

  function aggregator() external view override returns (address) {
    return s_aggregator;
  }
}
