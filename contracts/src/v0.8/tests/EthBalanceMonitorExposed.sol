// SPDX-License-Identifier: MIT

pragma solidity 0.8.6;

import "../automation/upkeeps/EthBalanceMonitor.sol";

contract EthBalanceMonitorExposed is EthBalanceMonitor {
  constructor(
    address keeperRegistryAddress,
    uint256 minWaitPeriod
  ) EthBalanceMonitor(keeperRegistryAddress, minWaitPeriod) {}

  function setLastTopUpXXXTestOnly(address target, uint56 lastTopUpTimestamp) external {
    s_targets[target].lastTopUpTimestamp = lastTopUpTimestamp;
  }
}
