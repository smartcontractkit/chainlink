// SPDX-License-Identifier: MIT

pragma solidity ^0.8.0;

import {IERC165} from "../../vendor/openzeppelin-solidity/v4.8.3/contracts/utils/introspection/IERC165.sol";

/**
 * @dev Library used to query support of an interface declared via {IERC165}.
 *
 * Note that these functions return the actual result of the query: they do not
 * `revert` if an interface is not supported. It is up to the caller to decide
 * what to do in these cases.
 *
 * Note this is exactly the same as the OZ version, with the exception that external calls will revert
 * if < 31_000 gas is available, to prevent message delivery issues in CCIP.
 */
library ERC165CheckerReverting {
  // As per the EIP-165 spec, no interface should ever match 0xffffffff
  bytes4 private constant _INTERFACE_ID_INVALID = 0xffffffff;

  // bytes4(keccak256("NotEnoughGasForSupportsInterfaceCall()"))
  bytes4 private constant NOT_ENOUGH_GAS_SIG = 0x161c3bf7;

  /**
   * @dev Returns true if `account` supports the {IERC165} interface.
   */
  function _supportsERC165Reverting(
    address account
  ) internal view returns (bool) {
    // Any contract that implements ERC165 must explicitly indicate support of
    // InterfaceId_ERC165 and explicitly indicate non-support of InterfaceId_Invalid
    return _supportsERC165InterfaceUncheckedReverting(account, type(IERC165).interfaceId)
      && !_supportsERC165InterfaceUncheckedReverting(account, _INTERFACE_ID_INVALID);
  }

  /**
   * @dev Returns true if `account` supports the interface defined by
   * `interfaceId`. Support for {IERC165} itself is queried automatically.
   *
   * See {IERC165-supportsInterface}.
   */
  function _supportsInterfaceReverting(address account, bytes4 interfaceId) internal view returns (bool) {
    // query support of both ERC165 as per the spec and support of _interfaceId
    return _supportsERC165Reverting(account) && _supportsERC165InterfaceUncheckedReverting(account, interfaceId);
  }

  /**
   * @dev Returns a boolean array where each value corresponds to the
   * interfaces passed in and whether they're supported or not. This allows
   * you to batch check interfaces for a contract where your expectation
   * is that some interfaces may not be supported.
   *
   * See {IERC165-supportsInterface}.
   *
   * _Available since v3.4._
   */
  function _getSupportedInterfacesReverting(
    address account,
    bytes4[] memory interfaceIds
  ) internal view returns (bool[] memory) {
    // an array of booleans corresponding to interfaceIds and whether they're supported or not
    bool[] memory interfaceIdsSupported = new bool[](interfaceIds.length);

    // query support of ERC165 itself
    if (_supportsERC165Reverting(account)) {
      // query support of each interface in interfaceIds
      for (uint256 i = 0; i < interfaceIds.length; i++) {
        interfaceIdsSupported[i] = _supportsERC165InterfaceUncheckedReverting(account, interfaceIds[i]);
      }
    }

    return interfaceIdsSupported;
  }

  /**
   * @dev Returns true if `account` supports all the interfaces defined in
   * `interfaceIds`. Support for {IERC165} itself is queried automatically.
   *
   * Batch-querying can lead to gas savings by skipping repeated checks for
   * {IERC165} support.
   *
   * See {IERC165-supportsInterface}.
   */
  function _supportsAllInterfacesReverting(address account, bytes4[] memory interfaceIds) internal view returns (bool) {
    // query support of ERC165 itself
    if (!_supportsERC165Reverting(account)) {
      return false;
    }

    // query support of each interface in interfaceIds
    for (uint256 i = 0; i < interfaceIds.length; i++) {
      if (!_supportsERC165InterfaceUncheckedReverting(account, interfaceIds[i])) {
        return false;
      }
    }

    // all interfaces supported
    return true;
  }

  /**
   * @notice Query if a contract implements an interface, does not check ERC165 support
   * @param account The address of the contract to query for support of an interface
   * @param interfaceId The interface identifier, as specified in ERC-165
   * @return true if the contract at account indicates support of the interface with
   * identifier interfaceId, false otherwise
   * @dev Assumes that account contains a contract that supports ERC165, otherwise
   * the behavior of this method is undefined. This precondition can be checked
   * with {supportsERC165}.
   *
   * Some precompiled contracts will falsely indicate support for a given interface, so caution
   * should be exercised when using this function.
   *
   * Interface identification is specified in ERC-165.
   */
  function _supportsERC165InterfaceUncheckedReverting(address account, bytes4 interfaceId) internal view returns (bool) {
    // prepare call
    bytes memory encodedParams = abi.encodeWithSelector(IERC165.supportsInterface.selector, interfaceId);

    // perform static call
    bool success;
    uint256 returnSize;
    uint256 returnValue;

    assembly {
      // Enforce that there's enough gas avilable so that the call will not fail due to OOG error. Without this supportsInterface() may return false when it should return true.

      // 32,000 gas was chosen to ensure enough gas to invoke the staticcall after the check
      // without breaking the 63/64 rule. Under EIP-150 there must be at least ~30,476 gas
      // remaining before the staticcall.
      if lt(gas(), 32000) {
        mstore(0x0, NOT_ENOUGH_GAS_SIG)
        revert(0x0, 0x4)
      }

      success := staticcall(30000, account, add(encodedParams, 0x20), mload(encodedParams), 0x00, 0x20)
      returnSize := returndatasize()
      returnValue := mload(0x00)
    }
    return success && returnSize >= 0x20 && returnValue > 0;
  }
}
