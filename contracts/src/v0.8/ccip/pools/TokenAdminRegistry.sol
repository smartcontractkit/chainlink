// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.19;

import {ITokenAdminRegistry} from "../interfaces/ITokenAdminRegistry.sol";

import {OwnerIsCreator} from "../../shared/access/OwnerIsCreator.sol";

import {EnumerableSet} from "../../vendor/openzeppelin-solidity/v4.8.3/contracts/utils/structs/EnumerableSet.sol";

// This contract has minimal functionality and minimal test coverage. It will be
// improved upon in future tickets.
contract TokenAdminRegistry is ITokenAdminRegistry, OwnerIsCreator {
  using EnumerableSet for EnumerableSet.AddressSet;

  error OnlyRegistryModule(address sender);
  error OnlyAdministrator(address sender, address token);
  error AlreadyRegistered(address token, address currentAdministrator);
  error UnsupportedToken(address token);

  event AdministratorRegistered(address indexed token, address indexed administrator);
  event PoolSet(address indexed token, address indexed pool);

  struct TokenConfig {
    address administrator; // ────────────────╮ the current administrator of the token
    bool isPermissionedAdmin; //              │ if true, this administrator has been configured by the CCIP owner
    //                                        │ and it could have elevated permissions.
    bool allowPermissionlessReRegistration; //│ if true, the token can be re-registered without the administrator's signature
    bool isRegistered; // ────────────────────╯ if true, the token is registered in the registry
    address tokenPool; // the token pool for this token. Can be address(0) if not deployed or not configured.
  }

  mapping(address token => TokenConfig) internal s_tokenConfig;
  EnumerableSet.AddressSet internal s_tokens;

  EnumerableSet.AddressSet internal s_RegistryModules;

  /// @notice Returns all pools for the given tokens.
  /// @dev Will return address(0) for tokens that do not have a pool.
  function getPools(address[] calldata tokens) external view returns (address[] memory) {
    address[] memory pools = new address[](tokens.length);
    for (uint256 i = 0; i < tokens.length; ++i) {
      pools[i] = s_tokenConfig[tokens[i]].tokenPool;
    }
    return pools;
  }

  /// @inheritdoc ITokenAdminRegistry
  function getPool(address token) external view returns (address) {
    address pool = s_tokenConfig[token].tokenPool;
    if (pool == address(0)) {
      revert UnsupportedToken(token);
    }

    return pool;
  }

  function setPool(address token, address pool) external {
    TokenConfig storage config = s_tokenConfig[token];
    if (config.administrator != msg.sender) {
      revert OnlyAdministrator(msg.sender, token);
    }

    config.tokenPool = pool;

    emit PoolSet(token, pool);
  }

  function getAllConfiguredTokens() external view returns (address[] memory) {
    return s_tokens.values();
  }

  // ================================================================
  // │                    Administrator config                      │
  // ================================================================

  /// @notice Public getter to check for permissions of an administrator
  function isAdministrator(address localToken, address administrator) public view returns (bool) {
    return s_tokenConfig[localToken].administrator == administrator;
  }

  /// @notice Resisters a new local administrator for a token.
  function registerAdministrator(address localToken, address administrator) external {
    // Only allow permissioned registry modules to register administrators
    if (!s_RegistryModules.contains(msg.sender)) {
      revert OnlyRegistryModule(msg.sender);
    }
    TokenConfig storage config = s_tokenConfig[localToken];

    if (config.isRegistered && !config.allowPermissionlessReRegistration) {
      revert AlreadyRegistered(localToken, config.administrator);
    }

    // If the token is not registered yet, or if re-registration is permitted, register the new administrator
    config.administrator = administrator;
    config.isRegistered = true;

    s_tokens.add(localToken);

    emit AdministratorRegistered(localToken, administrator);
  }

  /// @notice Registers a local administrator for a token. This will overwrite any potential current administrator
  /// and set the permissionedAdmin to true.
  function registerAdministratorPermissioned(address localToken, address administrator) external onlyOwner {
    TokenConfig storage config = s_tokenConfig[localToken];

    config.administrator = administrator;
    config.isRegistered = true;
    config.isPermissionedAdmin = true;

    s_tokens.add(localToken);

    emit AdministratorRegistered(localToken, administrator);
  }

  // ================================================================
  // │                      Registry Modules                        │
  // ================================================================

  function addRegistryModule(address module) external onlyOwner {
    s_RegistryModules.add(module);
  }
}
