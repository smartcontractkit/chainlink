// SPDX-License-Identifier: BUSL-1.1
// solhint-disable one-contract-per-file
pragma solidity ^0.8.0;

import {IBridgeAdapter} from "../../interfaces/IBridge.sol";
import {ILiquidityContainer} from "../../interfaces/ILiquidityContainer.sol";
import {IWrappedNative} from "../../../ccip/interfaces/IWrappedNative.sol";

import {IERC20} from "../../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/IERC20.sol";
import {SafeERC20} from "../../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/utils/SafeERC20.sol";

/// @notice Mock multiple-stage finalization bridge adapter implementation.
/// @dev Funds are only made available after both the prove and finalization steps are completed.
/// Sends the L1 tokens from the msg sender to address(this).
contract MockL1BridgeAdapter is IBridgeAdapter, ILiquidityContainer {
  using SafeERC20 for IERC20;

  error InsufficientLiquidity();
  error NonceAlreadyUsed(uint256 nonce);
  error InvalidFinalizationAction();
  error NonceNotProven(uint256 nonce);
  error NativeSendFailed();

  /// @notice Payload to "prove" the withdrawal.
  /// @dev This is just a mock setup, there's no real proving. This is so that
  /// we can test the multi-step finalization code path.
  /// @param nonce the nonce emitted on the remote chain.
  struct ProvePayload {
    uint256 nonce;
  }

  /// @notice Payload to "finalize" the withdrawal.
  /// @dev This is just a mock setup, there's no real finalization. This is so that
  /// we can test the multi-step finalization code path.
  /// @param nonce the nonce emitted on the remote chain.
  struct FinalizePayload {
    uint256 nonce;
    uint256 amount;
  }

  /// @notice The finalization action to take.
  /// @dev This emulates Optimism's two-step withdrawal process.
  enum FinalizationAction {
    ProveWithdrawal,
    FinalizeWithdrawal
  }

  /// @notice The payload to use for the bridgeSpecificPayload in the finalizeWithdrawERC20 function.
  struct Payload {
    FinalizationAction action;
    bytes data;
  }

  IERC20 internal immutable i_token;
  uint256 internal s_nonce = 1;
  mapping(uint256 => bool) internal s_nonceProven;
  mapping(uint256 => bool) internal s_nonceFinalized;

  /// @dev For test cases where we want to send pure native upon finalizeWithdrawERC20 being called.
  /// This is to emulate the behavior of bridges that do not bridge wrapped native.
  bool internal immutable i_holdNative;

  constructor(IERC20 token, bool holdNative) {
    i_token = token;
    i_holdNative = holdNative;
  }

  /// @dev The receive function is needed for IWrappedNative.withdraw() to work.
  receive() external payable {}

  /// @notice Simply transferFrom msg.sender the tokens that are to be bridged to address(this).
  function sendERC20(
    address localToken,
    address /* remoteToken */,
    address /* remoteReceiver */,
    uint256 amount,
    bytes calldata /* bridgeSpecificPayload */
  ) external payable override returns (bytes memory) {
    IERC20(localToken).transferFrom(msg.sender, address(this), amount);

    // If the flag to hold native is set we assume that i_token points to a WETH contract
    // and withdraw native.
    // This way we can transfer the raw native back to the sender upon finalization.
    if (i_holdNative) {
      IWrappedNative(address(i_token)).withdraw(amount);
    }

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

  /// @dev for easy encoding offchain
  function encodeProvePayload(ProvePayload memory payload) external pure {}

  function encodeFinalizePayload(FinalizePayload memory payload) external pure {}

  function encodePayload(Payload memory payload) external pure {}

  /// @dev Test setup is trusted, so just transfer the tokens to the localReceiver,
  /// which should be the local rebalancer. Infer the amount from the bridgeSpecificPayload.
  /// Note that this means that this bridge adapter will need to have some tokens,
  /// however this is ok in a test environment since we will have infinite tokens.
  /// @param localReceiver the address to transfer the tokens to.
  /// @param bridgeSpecificPayload the payload to use for the finalization or proving.
  /// @return true if the transfer was successful, revert otherwise.
  function finalizeWithdrawERC20(
    address /* remoteSender */,
    address localReceiver,
    bytes calldata bridgeSpecificPayload
  ) external override returns (bool) {
    Payload memory payload = abi.decode(bridgeSpecificPayload, (Payload));
    if (payload.action == FinalizationAction.ProveWithdrawal) {
      return _proveWithdrawal(payload);
    } else if (payload.action == FinalizationAction.FinalizeWithdrawal) {
      return _finalizeWithdrawal(payload, localReceiver);
    }
    revert InvalidFinalizationAction();
  }

  function _proveWithdrawal(Payload memory payload) internal returns (bool) {
    ProvePayload memory provePayload = abi.decode(payload.data, (ProvePayload));
    if (s_nonceProven[provePayload.nonce]) revert NonceAlreadyUsed(provePayload.nonce);
    s_nonceProven[provePayload.nonce] = true;
    return false;
  }

  function _finalizeWithdrawal(Payload memory payload, address localReceiver) internal returns (bool) {
    FinalizePayload memory finalizePayload = abi.decode(payload.data, (FinalizePayload));
    if (!s_nonceProven[finalizePayload.nonce]) revert NonceNotProven(finalizePayload.nonce);
    if (s_nonceFinalized[finalizePayload.nonce]) revert NonceAlreadyUsed(finalizePayload.nonce);
    s_nonceFinalized[finalizePayload.nonce] = true;
    // re-entrancy prevented by nonce checks above.
    _transferTokens(finalizePayload.amount, localReceiver);
    return true;
  }

  function _transferTokens(uint256 amount, address localReceiver) internal {
    if (i_holdNative) {
      (bool success, ) = payable(localReceiver).call{value: amount}("");
      if (!success) {
        revert NativeSendFailed();
      }
    } else {
      i_token.safeTransfer(localReceiver, amount);
    }
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
  ) external pure override returns (bool) {
    return true;
  }
}
