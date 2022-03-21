// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "./ConfirmedOwner.sol";
import {getRevertMsg} from "./utils/utils.sol";

/**
 * @title PermissionedForwardProxy
 * @notice This proxy is used to forward calls from sender to target. It maintains
 * a permission list to check which sender is allowed to call which target
 */
contract PermissionedForwardProxy is ConfirmedOwner {
  error CallFailed(string reason);

  mapping(address => address) public forwardPermissionList;

  constructor() ConfirmedOwner(msg.sender) {}

  /**
   * @notice Verifies if msg.sender has permission to forward to target address and then forwards the handler
   * @param target address of the contract to forward the handler to
   * @param handler bytes to be passed to target in call data
   */
  function forward(address target, bytes calldata handler) external {
    require(forwardPermissionList[msg.sender] == target, "Forwarding permission not found");
    (bool success, bytes memory payload) = target.call(handler);
    if (!success) {
      revert CallFailed(getRevertMsg(payload));
    }
  }

  /**
   * @notice Adds permission for sender to forward calls to target via this proxy.
   * Note that it allows to overwrite an existing permission
   * @param sender The address who will use this proxy to forward calls
   * @param target The address where sender will be allowed to forward calls
   */
  function addPermission(address sender, address target) external onlyOwner {
    forwardPermissionList[sender] = target;
  }

  /**
   * @notice Removes permission for sender to forward calls via this proxy
   * @param sender The address who will use this proxy to forward calls
   */
  function removePermission(address sender) external onlyOwner {
    delete forwardPermissionList[sender];
  }
}
