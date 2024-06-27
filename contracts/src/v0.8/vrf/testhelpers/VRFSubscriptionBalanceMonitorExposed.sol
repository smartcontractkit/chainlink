// SPDX-License-Identifier: MIT

pragma solidity 0.8.6;

import {VRFSubscriptionBalanceMonitor} from "../dev/VRFSubscriptionBalanceMonitor.sol";

contract VRFSubscriptionBalanceMonitorExposed is VRFSubscriptionBalanceMonitor {
  constructor(
    address linkTokenAddress,
    address coordinatorAddress,
    address keeperRegistryAddress,
    uint256 minWaitPeriodSeconds
  ) VRFSubscriptionBalanceMonitor(linkTokenAddress, coordinatorAddress, keeperRegistryAddress, minWaitPeriodSeconds) {}

  function setLastTopUpXXXTestOnly(uint64 target, uint56 lastTopUpTimestamp) external {
    s_targets[target].lastTopUpTimestamp = lastTopUpTimestamp;
  }
}
