// SPDX-License-Identifier: MIT
pragma solidity 0.8.6;

import {AutomationRegistryBaseInterface as IRegistry} from "./interfaces/AutomationRegistryInterface2_1.sol";
import "../../../interfaces/TypeAndVersionInterface.sol";

uint256 constant PERFORM_GAS_CUSHION = 5_000;

/**
 * @title AutomationForwarder is a relayer that sits between the registry and the customer's target contract
 * @dev The purpose of the forwarder is to give customers a consistent address to authorize against,
 * which stays consistent between migrations. The Forwarder also exposes the registry address, so that users who
 * want to programatically interact with the registry (ie top up funds) can do so.
 */
contract AutomationForwarder is TypeAndVersionInterface {
  address private s_registry;
  address private immutable i_target;
  string public constant override typeAndVersion = "AutomationForwarder 1.0.0";

  error NotAuthorized();

  constructor(address target) {
    s_registry = msg.sender;
    i_target = target;
  }

  /**
   * @notice forward is called by the registry and forwards the call to the target
   * @param gasAmount is the amount of gas to use in the call
   * @param data is the 4 bytes function selector + arbitrary function data
   * @return success indicating whether the target call succeeded or failed
   */
  function forward(uint256 gasAmount, bytes memory data) external returns (bool success) {
    if (msg.sender != s_registry) revert NotAuthorized();
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
  function updateRegistry(address newRegistry) external {
    if (msg.sender != s_registry) revert NotAuthorized();
    s_registry = newRegistry;
  }

  // TODO - here we should return a "user interface" that only contains the functions on the registry that a
  // user might want to interract with
  /**
   * @notice gets the registry address
   */
  function getRegistry() external view returns (address) {
    return s_registry;
  }

  /**
   * @notice gets the target contract address
   */
  function getTarget() external view returns (address) {
    return i_target;
  }
}
