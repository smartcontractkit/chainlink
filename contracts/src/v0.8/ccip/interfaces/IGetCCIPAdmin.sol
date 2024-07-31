// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

interface IGetCCIPAdmin {
  /// @notice Returns the admin of the token.
  /// @dev This method is named to never conflict with existing methods.
  function getCCIPAdmin() external view returns (address);
}
