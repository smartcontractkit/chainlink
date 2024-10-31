// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.4;

import {Ownable2StepMsgSender} from "./Ownable2StepMsgSender.sol";
import {EnumerableSet} from "../../vendor/openzeppelin-solidity/v4.8.3/contracts/utils/structs/EnumerableSet.sol";

/// @title The AuthorizedCallers contract
/// @notice A contract that manages multiple authorized callers. Enables restricting access to certain functions to a set of addresses.
contract AuthorizedCallers is Ownable2StepMsgSender {
  using EnumerableSet for EnumerableSet.AddressSet;

  event AuthorizedCallerAdded(address caller);
  event AuthorizedCallerRemoved(address caller);

  error UnauthorizedCaller(address caller);
  error ZeroAddressNotAllowed();

  /// @notice Update args for changing the authorized callers
  struct AuthorizedCallerArgs {
    address[] addedCallers;
    address[] removedCallers;
  }

  /// @dev Set of authorized callers
  EnumerableSet.AddressSet internal s_authorizedCallers;

  /// @param authorizedCallers the authorized callers to set
  constructor(address[] memory authorizedCallers) {
    _applyAuthorizedCallerUpdates(
      AuthorizedCallerArgs({addedCallers: authorizedCallers, removedCallers: new address[](0)})
    );
  }

  /// @return authorizedCallers Returns all authorized callers
  function getAllAuthorizedCallers() external view returns (address[] memory) {
    return s_authorizedCallers.values();
  }

  /// @notice Updates the list of authorized callers
  /// @param authorizedCallerArgs Callers to add and remove. Removals are performed first.
  function applyAuthorizedCallerUpdates(AuthorizedCallerArgs memory authorizedCallerArgs) external onlyOwner {
    _applyAuthorizedCallerUpdates(authorizedCallerArgs);
  }

  /// @notice Updates the list of authorized callers
  /// @param authorizedCallerArgs Callers to add and remove. Removals are performed first.
  function _applyAuthorizedCallerUpdates(AuthorizedCallerArgs memory authorizedCallerArgs) internal {
    address[] memory removedCallers = authorizedCallerArgs.removedCallers;
    for (uint256 i = 0; i < removedCallers.length; ++i) {
      address caller = removedCallers[i];

      if (s_authorizedCallers.remove(caller)) {
        emit AuthorizedCallerRemoved(caller);
      }
    }

    address[] memory addedCallers = authorizedCallerArgs.addedCallers;
    for (uint256 i = 0; i < addedCallers.length; ++i) {
      address caller = addedCallers[i];

      if (caller == address(0)) {
        revert ZeroAddressNotAllowed();
      }

      s_authorizedCallers.add(caller);
      emit AuthorizedCallerAdded(caller);
    }
  }

  /// @notice Checks the sender and reverts if it is anyone other than a listed authorized caller.
  function _validateCaller() internal view {
    if (!s_authorizedCallers.contains(msg.sender)) {
      revert UnauthorizedCaller(msg.sender);
    }
  }

  /// @notice Checks the sender and reverts if it is anyone other than a listed authorized caller.
  modifier onlyAuthorizedCallers() {
    _validateCaller();
    _;
  }
}
