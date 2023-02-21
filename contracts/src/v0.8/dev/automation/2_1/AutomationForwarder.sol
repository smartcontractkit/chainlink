// SPDX-License-Identifier: MIT
pragma solidity 0.8.6;

import {AutomationRegistryBaseInterface as IRegistry} from "../../interfaces/automation/2_1/AutomationRegistryInterface2_1.sol";
import "../../../interfaces/TypeAndVersionInterface.sol";

uint256 constant PERFORM_GAS_CUSHION = 5_000;

/**
 * @title AutomationForwarder is a relayer that sits between the registry and the customer's target contract
 * @dev The purpose of the forwarder is to give customers a consistent address to authorize against,
 * which stays consistent between migrations. The Forwarder also exposes the registry address, so that users who
 * want to programatically interact with the registry (ie top up funds) can do so.
 */
contract AutomationForwarder is TypeAndVersionInterface {
  IRegistry s_registry;
  address immutable i_target;
  string public constant override typeAndVersion = "AutomationForwarder 1.0.0";

  error NotAuthorized();

  constructor(IRegistry registry, address target) {
    s_registry = registry;
    i_target = target;
  }

  /**
   * @notice forward is called by the registry and forwards the call to the target
   * @param gasAmount is the amount of gas to use in the call
   * @param data is the 4 bytes function selector + arbitrary function data
   * @dev the forward function reverts with the same revert data as the target call
   */

  function forward(uint256 gasAmount, bytes memory data) external returns (bool success) {
    if (msg.sender != address(s_registry)) revert NotAuthorized();
    address target = i_target;
    assembly {
      let g := gas()
      // Compute g -= PERFORM_GAS_CUSHION and check for underflow
      if lt(g, PERFORM_GAS_CUSHION) {
        revert(0, 0)
      }
      g := sub(g, PERFORM_GAS_CUSHION)
      // if g - g//64 <= gasAmount, revert
      // (we subtract g//64 because of EIP-150)
      if iszero(gt(sub(g, div(g, 64)), gasAmount)) {
        revert(0, 0)
      }
      // solidity calls check that a contract actually exists at the destination, so we do the same
      if iszero(extcodesize(target)) {
        revert(0, 0)
      }
      // call with exact gas
      success := call(gasAmount, target, 0, add(data, 0x20), mload(data), 0, 0)
    }
  }

  /**
   * @notice updateRegistry is called by the registry during migrations
   * @param newRegistry is the registry that this forwarder is being migrated to
   */

  function updateRegistry(IRegistry newRegistry) external {
    if (msg.sender != address(s_registry)) revert NotAuthorized();
    s_registry = newRegistry;
  }

  /**
   * @notice gets the registry address
   */
  function getRegistry() external view returns (IRegistry) {
    return s_registry;
  }

  /**
   * @notice gets the target contract address
   */
  function getTarget() external view returns (address) {
    return i_target;
  }
}

/**
 * @title AutomationForwarderFactory is factory contract that deploys new AutomationForwarders
 * @dev while this functionality *could* live inside the Registry, the consious desision was made to
 * create a separate factory in the interest of saving space inside the registry
 */
contract AutomationForwarderFactory is TypeAndVersionInterface {
  event NewForwarderDeployed(AutomationForwarder forwarder);

  string public constant override typeAndVersion = "AutomationForwarderFactory 1.0.0";

  /**
   * @notice deploy deploys a new AutomationForwarder
   * @return AutomationForwarder the newly deployed contract instance
   */
  function deploy(address target) external returns (AutomationForwarder) {
    AutomationForwarder forwarder = new AutomationForwarder(IRegistry(msg.sender), target);
    emit NewForwarderDeployed(forwarder);
    return forwarder;
  }
}
