// SPDX-License-Identifier: MIT
pragma solidity 0.8.16;

import "./KeeperRegistryBase2_1.sol";
import "./interfaces/ILogAutomation.sol";

/**
 * @notice this file exposes structs that are otherwise internal to the automation registry
 * doing this allows those structs to be encoded and decoded with type safety in offchain code
 * and tests because generated wrappers are made available
 */

contract AutomationUtils2_1 {
  function _onChainConfig(KeeperRegistryBase2_1.OnchainConfig memory) external {}

  function _report(KeeperRegistryBase2_1.Report memory) external {}

  function _logTriggerConfig(KeeperRegistryBase2_1.LogTriggerConfig memory) external {}

  function _conditionalTriggerConfig(KeeperRegistryBase2_1.ConditionalTriggerConfig memory) external {}

  function _logTrigger(KeeperRegistryBase2_1.LogTrigger memory) external {}

  function _conditionalTrigger(KeeperRegistryBase2_1.ConditionalTrigger memory) external {}

  function _log(Log memory) external {}
}
