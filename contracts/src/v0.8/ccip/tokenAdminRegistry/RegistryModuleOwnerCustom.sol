// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.24;

import {ITypeAndVersion} from "../../shared/interfaces/ITypeAndVersion.sol";
import {IGetCCIPAdmin} from "../interfaces/IGetCCIPAdmin.sol";
import {IOwner} from "../interfaces/IOwner.sol";
import {ITokenAdminRegistry} from "../interfaces/ITokenAdminRegistry.sol";

import {OwnerIsCreator} from "../../shared/access/OwnerIsCreator.sol";

contract RegistryModuleOwnerCustom is ITypeAndVersion, OwnerIsCreator {
  error CanOnlySelfRegister(address admin, address token);

  event AdministratorRegistered(address indexed token, address indexed administrator);

  string public constant override typeAndVersion = "RegistryModuleOwnerCustom 1.5.0-dev";

  // The TokenAdminRegistry contract
  ITokenAdminRegistry internal s_tokenAdminRegistry;

  constructor(address tokenAdminRegistry) {
    s_tokenAdminRegistry = ITokenAdminRegistry(tokenAdminRegistry);
  }

  /// @notice Registers the admin of the token using the `getCCIPAdmin` method.
  /// @param token The token to register the admin for.
  /// @dev The caller must be the admin returned by the `getCCIPAdmin` method.
  function registerAdminViaGetCCIPAdmin(address token) external {
    _registerAdmin(token, IGetCCIPAdmin(token).getCCIPAdmin());
  }

  /// @notice Registers the admin of the token using the `owner` method.
  /// @param token The token to register the admin for.
  /// @dev The caller must be the admin returned by the `owner` method.
  function registerAdminViaOwner(address token) external {
    _registerAdmin(token, IOwner(token).owner());
  }

  /// @notice Registers the admin of the token to msg.sender given that the
  /// admin is equal to msg.sender.
  /// @param token The token to register the admin for.
  /// @param admin The caller must be the admin.
  function _registerAdmin(address token, address admin) internal {
    if (admin != msg.sender) {
      revert CanOnlySelfRegister(admin, token);
    }

    s_tokenAdminRegistry.registerAdministrator(token, admin);

    emit AdministratorRegistered(token, admin);
  }
}
