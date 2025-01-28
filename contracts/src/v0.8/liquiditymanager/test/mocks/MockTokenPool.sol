// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.0;

import {ILiquidityContainer} from "../../interfaces/ILiquidityContainer.sol";

import {IERC20} from "../../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/IERC20.sol";
import {SafeERC20} from "../../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/utils/SafeERC20.sol";

contract MockLockReleaseTokenPool is ILiquidityContainer {
  using SafeERC20 for IERC20;

  IERC20 public immutable i_localToken;

  constructor(IERC20 localToken) {
    i_localToken = localToken;
  }

  /// @notice Provide additional liquidity to the container.
  /// @dev Should emit LiquidityAdded
  function provideLiquidity(uint256 amount) external {
    i_localToken.safeTransferFrom(msg.sender, address(this), amount);
    emit LiquidityAdded(msg.sender, amount);
  }

  /// @notice Withdraws liquidity from the container to the msg sender
  /// @dev Should emit LiquidityRemoved
  function withdrawLiquidity(uint256 amount) external {
    i_localToken.safeTransfer(msg.sender, amount);
    emit LiquidityRemoved(msg.sender, amount);
  }
}
