// SPDX-License-Identifier: MIT
pragma solidity 0.8.16;

import "./KeeperRegistryBase2_1.sol";
import "./interfaces/ILogAutomation.sol";

/**
 * @notice this file exposes structs that are otherwise internal to the automation registry
 * doing this allows those structs to be encoded and decoded with type safety in offchain code
 * and tests because generated wrappers are made available
 */

/**
 * @notice structure of trigger for log triggers
 */
struct LogTriggerConfig {
  address contractAddress;
  uint8 filterSelector; // denotes which topics apply to filter ex 000, 101, 111...only last 3 bits apply
  bytes32 topic0;
  bytes32 topic1;
  bytes32 topic2;
  bytes32 topic3;
}

contract AutomationUtils2_1 {
  /**
   * @dev this can be removed as OnchainConfig is now exposed directly from the registry
   */
  function _onChainConfig(KeeperRegistryBase2_1.OnchainConfig memory) external {}

  function _report(KeeperRegistryBase2_1.Report memory) external {}

  function _logTriggerConfig(LogTriggerConfig memory) external {}

  function _logTrigger(KeeperRegistryBase2_1.LogTrigger memory) external {}

  function _conditionalTrigger(KeeperRegistryBase2_1.ConditionalTrigger memory) external {}

  function _log(Log memory) external {}
}
