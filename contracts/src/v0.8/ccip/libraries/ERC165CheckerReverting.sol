// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import {IERC165} from "../../vendor/openzeppelin-solidity/v4.8.3/contracts/utils/introspection/IERC165.sol";

/// @notice Library used to query support of an interface declared via {IERC165}.
/// @dev These functions return the actual result of the query: they do not `revert` if an interface is not supported.
library ERC165CheckerReverting {
  error InsufficientGasForStaticcall();

  // As per the EIP-165 spec, no interface should ever match 0xffffffff
  bytes4 private constant INTERFACE_ID_INVALID = 0xffffffff;

  /// @dev 30k gas is required to make the staticcall. Under the 63/64 rule this means that 30,477 gas must be available
  /// to ensure that at least 30k is forwarded. Checking for at least 31,000 ensures that after additional
  /// operations are performed there is still >= 30,477 gas remaining.
  /// 30,000 = ((30,477 * 63) / 64)
  uint256 private constant MINIMUM_GAS_REQUIREMENT = 31_000;

  /// @notice Returns true if `account` supports the {IERC165} interface.
  /// @dev Any contract that implements ERC165 must explicitly indicate support of InterfaceId_ERC165 and explicitly
  /// indicate non-support of InterfaceId_Invalid as per the standard.
  /// @param account the address to be queried for support for ERC165
  /// @return true if the contract at account indicates support for ERC165, false otherwise
  function _supportsERC165Reverting(
    address account
  ) internal view returns (bool) {
    return _supportsERC165InterfaceUncheckedReverting(account, type(IERC165).interfaceId)
      && !_supportsERC165InterfaceUncheckedReverting(account, INTERFACE_ID_INVALID);
  }

  /// @notice Returns true if `account` supports a defined interface
  /// @dev The function must support both the interfaceId and interfaces specified by ERC165 generally as per the standard
  /// @param account the contract to be queried for support
  /// @param interfaceId the interface being checked for support
  /// @return true if the contract at account indicates support of the interface with, false otherwise.
  function _supportsInterfaceReverting(address account, bytes4 interfaceId) internal view returns (bool) {
    // query support of both ERC165 as per the spec and support of _interfaceId
    return _supportsERC165Reverting(account) && _supportsERC165InterfaceUncheckedReverting(account, interfaceId);
  }

  /// @notice Query if a contract implements an interface, does not check ERC165 support
  /// @param account The address of the contract to query for support of an interface
  /// @param interfaceId The interface identifier, as specified in ERC-165
  /// @return true if the contract at account indicates support of the interface with
  /// identifier interfaceId, false otherwise
  /// @dev Assumes that account contains a contract that supports ERC165, otherwise
  /// the behavior of this method is undefined. This precondition can be checked
  /// @dev Function will only revert if the minimum gas requirement is not met before the staticcall is performed.
  function _supportsERC165InterfaceUncheckedReverting(address account, bytes4 interfaceId) internal view returns (bool) {
    bytes memory encodedParams = abi.encodeWithSelector(IERC165.supportsInterface.selector, interfaceId);

    bool success;
    uint256 returnSize;
    uint256 returnValue;

    bytes4 notEnoughGasSelector = InsufficientGasForStaticcall.selector;

    assembly {
      // The EVM does not return a specific error code if a revert is due to OOG. This check ensures that
      // the message will not throw an OOG error by requiring that the amount of gas for the following
      // staticcall exists before invoking it.
      if lt(gas(), MINIMUM_GAS_REQUIREMENT) {
        mstore(0x0, notEnoughGasSelector)
        revert(0x0, 0x4)
      }

      success := staticcall(30000, account, add(encodedParams, 0x20), mload(encodedParams), 0x00, 0x20)
      returnSize := returndatasize()
      returnValue := mload(0x00)
    }
    return success && returnSize >= 0x20 && returnValue > 0;
  }
}
