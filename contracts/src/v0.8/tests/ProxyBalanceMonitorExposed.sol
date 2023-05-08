// SPDX-License-Identifier: MIT

pragma solidity 0.8.6;

import "../dev/automation/upkeeps/ProxyBalanceMonitor.sol";

contract ProxyBalanceMonitorExposed is ProxyBalanceMonitor {
  constructor(address linkTokenAddress, uint256 minWaitPeriod) ProxyBalanceMonitor(linkTokenAddress, minWaitPeriod) {}

  function setLastTopUpXXXTestOnly(address target, uint56 lastTopUpTimestamp) external {
    s_targets[target].lastTopUpTimestamp = lastTopUpTimestamp;
  }
}
