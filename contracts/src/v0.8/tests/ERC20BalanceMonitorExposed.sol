// SPDX-License-Identifier: MIT

pragma solidity 0.8.6;

import "../automation/upkeeps/ERC20BalanceMonitor.sol";

contract ERC20BalanceMonitorExposed is ERC20BalanceMonitor {
  constructor(
    address erc20TokenAddress,
    address keeperRegistryAddress,
    uint256 minWaitPeriod
  ) ERC20BalanceMonitor(erc20TokenAddress, keeperRegistryAddress, minWaitPeriod) {}

  function setLastTopUpXXXTestOnly(address target, uint56 lastTopUpTimestamp) external {
    s_targets[target].lastTopUpTimestamp = lastTopUpTimestamp;
  }
}
