// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.0;

interface ITokenAdminRegistry {
  /// @notice Returns the pool for the given token.
  function getPool(address token) external view returns (address);

  /// @notice Registers an administrator for the given token.
  /// @param localToken The token to register the administrator for.
  /// @param administrator The administrator to register.
  function registerAdministrator(address localToken, address administrator) external;
}
