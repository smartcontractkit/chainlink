// SPDX-License-Identifier: MIT
pragma solidity 0.8.6;

import "@openzeppelin/contracts/utils/Address.sol";
import "./ConfirmedOwner.sol";

/**
 * @title PermissionedForwardProxy
 * @notice This proxy is used to forward calls from sender to target. It maintains
 * a permission list to check which sender is allowed to call which target
 */
contract PermissionedForwardProxy is ConfirmedOwner {
  using Address for address;

  error PermissionNotSet();

  event PermissionSet(address indexed sender, address target);
  event PermissionRemoved(address indexed sender);

  mapping(address => address) private s_forwardPermissionList;

  constructor() ConfirmedOwner(msg.sender) {}

  /**
   * @notice Verifies if msg.sender has permission to forward to target address and then forwards the handler
   * @param target address of the contract to forward the handler to
   * @param handler bytes to be passed to target in call data
   */
  function forward(address target, bytes calldata handler) external {
    if (s_forwardPermissionList[msg.sender] != target) {
      revert PermissionNotSet();
    }
    target.functionCall(handler);
  }

  /**
   * @notice Adds permission for sender to forward calls to target via this proxy.
   * Note that it allows to overwrite an existing permission
   * @param sender The address who will use this proxy to forward calls
   * @param target The address where sender will be allowed to forward calls
   */
  function setPermission(address sender, address target) external onlyOwner {
    s_forwardPermissionList[sender] = target;

    emit PermissionSet(sender, target);
  }

  /**
   * @notice Removes permission for sender to forward calls via this proxy
   * @param sender The address who will use this proxy to forward calls
   */
  function removePermission(address sender) external onlyOwner {
    delete s_forwardPermissionList[sender];

    emit PermissionRemoved(sender);
  }

  /**
   * @notice Returns the target address that the sender can use this proxy for
   * @param sender The address to fetch the permissioned target for
   */
  function getPermission(address sender) external view returns (address) {
    return s_forwardPermissionList[sender];
  }
}
