// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {ITypeAndVersion} from "../../shared/interfaces/ITypeAndVersion.sol";
import {IGetCCIPAdmin} from "../interfaces/IGetCCIPAdmin.sol";
import {IOwner} from "../interfaces/IOwner.sol";
import {ITokenAdminRegistry} from "../interfaces/ITokenAdminRegistry.sol";

import {AccessControl} from "../../vendor/openzeppelin-solidity/v5.0.2/contracts/access/AccessControl.sol";

contract RegistryModuleOwnerCustom is ITypeAndVersion {
  error CanOnlySelfRegister(address admin, address token);
  error RequiredRoleNotFound(address msgSender, bytes32 role, address token);
  error AddressZero();

  event AdministratorRegistered(address indexed token, address indexed administrator);

  string public constant override typeAndVersion = "RegistryModuleOwnerCustom 1.6.0";

  // The TokenAdminRegistry contract
  ITokenAdminRegistry internal immutable i_tokenAdminRegistry;

  constructor(
    address tokenAdminRegistry
  ) {
    if (tokenAdminRegistry == address(0)) {
      revert AddressZero();
    }
    i_tokenAdminRegistry = ITokenAdminRegistry(tokenAdminRegistry);
  }

  /// @notice Registers the admin of the token using the `getCCIPAdmin` method.
  /// @param token The token to register the admin for.
  /// @dev The caller must be the admin returned by the `getCCIPAdmin` method.
  function registerAdminViaGetCCIPAdmin(
    address token
  ) external {
    _registerAdmin(token, IGetCCIPAdmin(token).getCCIPAdmin());
  }

  /// @notice Registers the admin of the token using the `owner` method.
  /// @param token The token to register the admin for.
  /// @dev The caller must be the admin returned by the `owner` method.
  function registerAdminViaOwner(
    address token
  ) external {
    _registerAdmin(token, IOwner(token).owner());
  }

  /// @notice Registers the admin of the token using OZ's AccessControl DEFAULT_ADMIN_ROLE.
  /// @param token The token to register the admin for.
  /// @dev The caller must have the DEFAULT_ADMIN_ROLE as defined by the contract itself.
  function registerAccessControlDefaultAdmin(
    address token
  ) external {
    bytes32 defaultAdminRole = AccessControl(token).DEFAULT_ADMIN_ROLE();
    if (!AccessControl(token).hasRole(defaultAdminRole, msg.sender)) {
      revert RequiredRoleNotFound(msg.sender, defaultAdminRole, token);
    }

    _registerAdmin(token, msg.sender);
  }

  /// @notice Registers the admin of the token to msg.sender given that the
  /// admin is equal to msg.sender.
  /// @param token The token to register the admin for.
  /// @param admin The caller must be the admin.
  function _registerAdmin(address token, address admin) internal {
    if (admin != msg.sender) {
      revert CanOnlySelfRegister(admin, token);
    }

    i_tokenAdminRegistry.proposeAdministrator(token, admin);

    emit AdministratorRegistered(token, admin);
  }
}
