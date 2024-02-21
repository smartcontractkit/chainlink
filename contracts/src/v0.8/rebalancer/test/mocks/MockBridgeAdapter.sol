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
  error NonceAlreadyUsed(uint256 nonce);

  IERC20 internal immutable i_token;
  uint256 internal s_nonce = 1;
  mapping(uint256 => bool) internal s_nonceUsed;

  constructor(IERC20 token) {
    i_token = token;
  }

  /// @notice Simply transferFrom msg.sender the tokens that are to be bridged to address(this).
  function sendERC20(
    address localToken,
    address /* remoteToken */,
    address /* remoteReceiver */,
    uint256 amount,
    bytes calldata /* bridgeSpecificPayload */
  ) external payable override returns (bytes memory) {
    IERC20(localToken).transferFrom(msg.sender, address(this), amount);
    bytes memory encodedNonce = abi.encode(s_nonce++);
    return encodedNonce;
  }

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

  /// @dev Test setup is trusted, so just transfer the tokens to the localReceiver,
  /// @dev which should be the local rebalancer.
  /// @dev Infer the amount from the bridgeSpecificPayload
  /// @dev Note that this means that this bridge adapter will need to have some tokens,
  /// @dev however this is ok in a test environment since we will have infinite tokens.
  function finalizeWithdrawERC20(
    address /* remoteSender */,
    address localReceiver,
    bytes calldata bridgeSpecificPayload
  ) external {
    (uint256 amount, uint256 nonce) = abi.decode(bridgeSpecificPayload, (uint256, uint256));
    if (s_nonceUsed[nonce]) revert NonceAlreadyUsed(nonce);
    s_nonceUsed[nonce] = true;
    i_token.safeTransfer(localReceiver, amount);
  }
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
