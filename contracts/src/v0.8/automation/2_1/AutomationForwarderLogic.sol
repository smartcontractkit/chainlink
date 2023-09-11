// SPDX-License-Identifier: MIT
pragma solidity 0.8.16;

import {IAutomationRegistryConsumer} from "./interfaces/IAutomationRegistryConsumer.sol";
import {ITypeAndVersion} from "../../../shared/interfaces/ITypeAndVersion.sol";

contract AutomationForwarderLogic is ITypeAndVersion {
  IAutomationRegistryConsumer private s_registry;

  string public constant typeAndVersion = "AutomationForwarder 1.0.0";

  /**
   * @notice updateRegistry is called by the registry during migrations
   * @param newRegistry is the registry that this forwarder is being migrated to
   */
  function updateRegistry(address newRegistry) external {
    if (msg.sender != address(s_registry)) revert();
    s_registry = IAutomationRegistryConsumer(newRegistry);
  }

  function getRegistry() external view returns (IAutomationRegistryConsumer) {
    return s_registry;
  }
}
