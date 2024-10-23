// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.0;

interface ITokenAdminRegistry {
  /// @notice Returns the pool for the given token.
  function getPool(
    address token
  ) external view returns (address);

  /// @notice Proposes an administrator for the given token as pending administrator.
  /// @param localToken The token to register the administrator for.
  /// @param administrator The administrator to register.
  function proposeAdministrator(address localToken, address administrator) external;

  /// @notice Accepts the administrator role for a token.
  /// @param localToken The token to accept the administrator role for.
  /// @dev This function can only be called by the pending administrator.
  function acceptAdminRole(
    address localToken
  ) external;

  /// @notice Sets the pool for a token. Setting the pool to address(0) effectively delists the token
  /// from CCIP. Setting the pool to any other address enables the token on CCIP.
  /// @param localToken The token to set the pool for.
  /// @param pool The pool to set for the token.
  function setPool(address localToken, address pool) external;

  /// @notice Transfers the administrator role for a token to a new address with a 2-step process.
  /// @param localToken The token to transfer the administrator role for.
  /// @param newAdmin The address to transfer the administrator role to. Can be address(0) to cancel
  /// a pending transfer.
  /// @dev The new admin must call `acceptAdminRole` to accept the role.
  function transferAdminRole(address localToken, address newAdmin) external;
}
