// SPDX-License-Identifier: BUSL-1.1
// solhint-disable one-contract-per-file
pragma solidity ^0.8.0;

import {IBridgeAdapter, IL1BridgeAdapter} from "../../interfaces/IBridge.sol";
import {ILiquidityContainer} from "../../interfaces/ILiquidityContainer.sol";

import {IERC20} from "../../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/IERC20.sol";
import {SafeERC20} from "../../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/utils/SafeERC20.sol";

/// @notice Mock L1 Bridge adapter
/// @dev Sends the L1 tokens from the msg sender to address(this)
contract MockL1BridgeAdapter is IL1BridgeAdapter, ILiquidityContainer {
  using SafeERC20 for IERC20;

  error InsufficientLiquidity();

  IERC20 internal immutable i_token;

  constructor(IERC20 token) {
    i_token = token;
  }

  function sendERC20(address l1Token, address, address, uint256 amount) external payable {
    IERC20(l1Token).transferFrom(msg.sender, address(this), amount);
  }

  /// @notice Mock function to finalize a withdrawal from L2
  /// @dev Does nothing as the indented action cannot be inferred from the inputs
  function finalizeWithdrawERC20FromL2(
    address l2Sender,
    address l1Receiver,
    bytes calldata bridgeSpecificPayload
  ) external {}

  function getBridgeFeeInNative() external pure returns (uint256) {
    return 0;
  }

  function provideLiquidity(uint256 amount) external {
    i_token.safeTransferFrom(msg.sender, address(this), amount);
    emit LiquidityAdded(msg.sender, amount);
  }

  function withdrawLiquidity(uint256 amount) external {
    if (i_token.balanceOf(address(this)) < amount) revert InsufficientLiquidity();
    i_token.safeTransfer(msg.sender, amount);
    emit LiquidityRemoved(msg.sender, amount);
  }
}

/// @notice Mock L2 Bridge adapter
/// @dev Sends the L2 tokens from the msg sender to address(this)
contract MockL2BridgeAdapter is IBridgeAdapter {
  function sendERC20(address, address l2token, address, uint256 amount) external payable {
    IERC20(l2token).transferFrom(msg.sender, address(this), amount);
  }

  function getBridgeFeeInNative() external pure returns (uint256) {
    return 0;
  }
}
