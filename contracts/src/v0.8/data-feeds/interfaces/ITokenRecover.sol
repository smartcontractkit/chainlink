// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {IERC20} from "./../../vendor/openzeppelin-solidity/v5.0.2/contracts/interfaces/IERC20.sol";

/// @notice ITokenRecover
/// Implements the recoverTokens method, enabling the recovery of ERC-20 or native tokens accidentally sent to a
/// contract outside of normal operations.
interface ITokenRecover {
  /// @notice Transfer any ERC-20 or native tokens accidentally sent to this contract.
  /// @param token Token to transfer
  /// @param to Address to send payment to
  /// @param amount Amount of token to transfer
  function recoverTokens(IERC20 token, address to, uint256 amount) external;
}
