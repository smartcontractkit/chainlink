// SPDX-License-Identifier: MIT

pragma solidity 0.8.6;

import "../automation/upkeeps/EthBalanceMonitor.sol";
import "../dev/automation/upkeeps/EthBalanceMonitorExtended.sol";

contract EthBalanceMonitorExposed is EthBalanceMonitor {
  constructor(
    address keeperRegistryAddress,
    uint256 minWaitPeriod
  ) EthBalanceMonitor(keeperRegistryAddress, minWaitPeriod) {}

  function setLastTopUpXXXTestOnly(address target, uint56 lastTopUpTimestamp) external {
    s_targets[target].lastTopUpTimestamp = lastTopUpTimestamp;
  }
}

contract EthBalanceMonitorExtendedExposed is EthBalanceMonitorExtended {
  constructor(
    address keeperRegistryAddress,
    uint256 minWaitPeriod
  ) EthBalanceMonitorExtended(keeperRegistryAddress, minWaitPeriod) {}

  function setLastTopUpXXXTestOnly(address target, uint56 lastTopUpTimestamp) external {
    s_targets[target].lastTopUpTimestamp = lastTopUpTimestamp;
  }
}
