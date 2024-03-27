// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import {ConfirmedOwner} from "./ConfirmedOwner.sol";
import {AccessControllerInterface} from "../interfaces/AccessControllerInterface.sol";

/// @title SimpleWriteAccessController
/// @notice Gives access to accounts explicitly added to an access list by the controller's owner.
/// @dev does not make any special permissions for externally, see  SimpleReadAccessController for that.
contract SimpleWriteAccessController is AccessControllerInterface, ConfirmedOwner {
  bool public checkEnabled;
  mapping(address => bool) internal s_accessList;

  event AddedAccess(address user);
  event RemovedAccess(address user);
  event CheckAccessEnabled();
  event CheckAccessDisabled();

  constructor() ConfirmedOwner(msg.sender) {
    checkEnabled = true;
  }

  /// @notice Returns the access of an address
  /// @param _user The address to query
  function hasAccess(address _user, bytes memory) public view virtual override returns (bool) {
    return s_accessList[_user] || !checkEnabled;
  }

  /// @notice Adds an address to the access list
  /// @param _user The address to add
  function addAccess(address _user) external onlyOwner {
    if (!s_accessList[_user]) {
      s_accessList[_user] = true;

      emit AddedAccess(_user);
    }
  }

  /// @notice Removes an address from the access list
  /// @param _user The address to remove
  function removeAccess(address _user) external onlyOwner {
    if (s_accessList[_user]) {
      s_accessList[_user] = false;

      emit RemovedAccess(_user);
    }
  }

  /// @notice makes the access check enforced
  function enableAccessCheck() external onlyOwner {
    if (!checkEnabled) {
      checkEnabled = true;

      emit CheckAccessEnabled();
    }
  }

  /// @notice makes the access check unenforced
  function disableAccessCheck() external onlyOwner {
    if (checkEnabled) {
      checkEnabled = false;

      emit CheckAccessDisabled();
    }
  }

  /// @dev reverts if the caller does not have access
  modifier checkAccess() {
    // solhint-disable-next-line gas-custom-errors
    require(hasAccess(msg.sender, msg.data), "No access");
    _;
  }
}
