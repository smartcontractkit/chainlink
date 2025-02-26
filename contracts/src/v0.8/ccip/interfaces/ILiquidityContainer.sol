// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.0;

/// @notice Interface for a liquidity container, this can be a CCIP token pool.
interface ILiquidityContainer {
  event LiquidityAdded(address indexed provider, uint256 indexed amount);
  event LiquidityRemoved(address indexed provider, uint256 indexed amount);

  /// @notice Provide additional liquidity to the container.
  /// @dev Should emit LiquidityAdded
  function provideLiquidity(
    uint256 amount
  ) external;

  /// @notice Withdraws liquidity from the container to the msg sender
  /// @dev Should emit LiquidityRemoved
  function withdrawLiquidity(
    uint256 amount
  ) external;
}
