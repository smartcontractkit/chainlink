// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

import {EnumerableSet} from "../../vendor/openzeppelin-solidity/v5.0.2/contracts/utils/structs/EnumerableSet.sol";

/// @title DONAccessControl
/// @notice Abstract contract for managing DON access control
/// @dev Provides granular permission management for DON IDs and authorized addresses
abstract contract DONAccessControl {
  using EnumerableSet for EnumerableSet.UintSet;
  using EnumerableSet for EnumerableSet.AddressSet;

  // Struct
  struct DONPermission {
    uint32 donID;
    address authorizedAddress;
  }

  // Mappings
  // Set of all allowed DON IDs
  EnumerableSet.UintSet internal s_allowedDONs;
  // Mapping from keccak256(donID, address) to permission bool
  mapping(bytes32 => bool) internal s_DONPermissions;
  // Mapping from DON ID to set of authorized addresses. Needed to list all permissions.
  mapping(uint32 => EnumerableSet.AddressSet) internal s_DONAuthorizedAddresses;

  // Events
  event AllowedDONUpdatedV1(uint32 indexed donID, bool allowed);
  event DONPermissionUpdatedV1(uint32 indexed donID, address indexed authorizedAddress, bool allowed);

  // Errors
  error AddressNotAuthorized(uint32 donID, address caller);
  error DONNotAllowed(uint32 donID);

  /// @notice Updates the allowed status for a single DON
  /// @param donID The ID of the DON to update
  /// @param allowed The new allowed status
  function _updateAllowedDON(uint32 donID, bool allowed) internal {
    if (allowed) {
      s_allowedDONs.add(donID);
    } else {
      s_allowedDONs.remove(donID);
    }

    emit AllowedDONUpdatedV1(donID, allowed);
  }

  /// @notice Updates permission for a single address and DON combination
  /// @param donID The ID of the DON
  /// @param authorizedAddress The address to update permissions for
  /// @param allowed The new permission status
  function _updateDONPermission(uint32 donID, address authorizedAddress, bool allowed) internal {
    bytes32 accessKey = _computeAccessKey(donID, authorizedAddress);
    s_DONPermissions[accessKey] = allowed;

    if (allowed) {
      s_DONAuthorizedAddresses[donID].add(authorizedAddress);
    } else {
      s_DONAuthorizedAddresses[donID].remove(authorizedAddress);
    }

    emit DONPermissionUpdatedV1(donID, authorizedAddress, allowed);
  }

  /// @notice Computes a unique key for storing DON address permissions
  /// @dev Combines donID and address using keccak256
  /// @param donID The ID of the DON
  /// @param authorizedAddress The address to compute the key for
  /// @return bytes32 The computed unique key
  // Helper function to compute a unique key from donID and address
  function _computeAccessKey(uint32 donID, address authorizedAddress) internal pure returns (bytes32) {
    return keccak256(abi.encodePacked(donID, authorizedAddress));
  }

  /// @notice Checks if an address has access to a specific DON
  /// @param donID The ID of the DON
  /// @param addr The address to check
  /// @return bool True if the address has access, false otherwise
  function _hasAccess(uint32 donID, address addr) internal view returns (bool) {
    bytes32 accessKey = _computeAccessKey(donID, addr);
    return s_DONPermissions[accessKey];
  }

  /// @notice Validates access permissions for a given DON and caller
  /// @dev Reverts with DONNotAllowed if the DON is not allowed or AddressNotAuthorized if the caller lacks permission
  /// @param donID The ID of the DON to check
  /// @param caller The address attempting to access the DON
  function _validateDONPermission(uint32 donID, address caller) internal view {
    if (!s_allowedDONs.contains(donID)) {
      // First, ensure the DON is in the allowed list. This is separate from the permission check below because a DON
      // can be removed from the allowed list without removing the permissioned addresses associated with the DON.
      revert DONNotAllowed(donID);
    }

    if (!_hasAccess(donID, caller)) {
      revert AddressNotAuthorized(donID, caller); // Then, ensure the specific address is authorized for the DON
    }
  }
}
