// SPDX-License-Identifier: BUSL-1.1
// solhint-disable one-contract-per-file
pragma solidity ^0.8.0;

import {IBridgeAdapter} from "../../interfaces/IBridge.sol";
import {ILiquidityContainer} from "../../interfaces/ILiquidityContainer.sol";

import {IERC20} from "../../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/IERC20.sol";
import {SafeERC20} from "../../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/utils/SafeERC20.sol";

/// @notice Mock L1 Bridge adapter
/// @dev Sends the L1 tokens from the msg sender to address(this)
contract MockL1BridgeAdapter is IBridgeAdapter, ILiquidityContainer {
  using SafeERC20 for IERC20;

  error InsufficientLiquidity();

  IERC20 internal immutable i_token;

  constructor(IERC20 token) {
    i_token = token;
  }

  /// @notice Simply transferFrom msg.sender the tokens that are to be bridged.
  function sendERC20(
    address localToken,
    address /* remoteToken */,
    address /* receiver */,
    uint256 amount,
    bytes calldata /* bridgeSpecificPayload */
  ) external payable override returns (bytes memory) {
    IERC20(localToken).transferFrom(msg.sender, address(this), amount);
    return "";
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

  // No-op
  function finalizeWithdrawERC20(
    address /* remoteSender */,
    address /* localReceiver */,
    bytes calldata /* bridgeSpecificData */
  ) external {}
}

/// @notice Mock L2 Bridge adapter
/// @dev Sends the L2 tokens from the msg sender to address(this)
contract MockL2BridgeAdapter is IBridgeAdapter {
  /// @notice Simply transferFrom msg.sender the tokens that are to be bridged.
  function sendERC20(
    address localToken,
    address /* remoteToken */,
    address /* recipient */,
    uint256 amount,
    bytes calldata /* bridgeSpecificPayload */
  ) external payable override returns (bytes memory) {
    IERC20(localToken).transferFrom(msg.sender, address(this), amount);
    return "";
  }

  function getBridgeFeeInNative() external pure returns (uint256) {
    return 0;
  }

  // No-op
  function finalizeWithdrawERC20(
    address /* remoteSender */,
    address /* localReceiver */,
    bytes calldata /* bridgeSpecificData */
  ) external override {}
}
