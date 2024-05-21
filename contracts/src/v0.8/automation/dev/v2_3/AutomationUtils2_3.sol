// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.19;

import {AutomationRegistryBase2_3} from "./AutomationRegistryBase2_3.sol";

/**
 * @notice this file exposes structs that are otherwise internal to the automation registry
 * doing this allows those structs to be encoded and decoded with type safety in offchain code
 * and tests because generated wrappers are made available
 */

contract AutomationUtils2_3 {
  /**
   * @dev this uses the v2.3 Report, which uses linkUSD instead of linkNative (as in v2.2 and prior). This should be used only in typescript tests.
   */
  function _report(AutomationRegistryBase2_3.Report memory) external {} // 0xe65d6546
}
