// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.0;

/// @dev IBridgeAdapter provides a common interface to interact with the native bridge.
interface IBridgeAdapter {
  error BridgeAddressCannotBeZero();
  error MsgValueDoesNotMatchAmount(uint256 msgValue, uint256 amount);
  error InsufficientEthValue(uint256 wanted, uint256 got);
  error MsgShouldNotContainValue(uint256 value);

  /// @notice Send the specified amount of the local token cross-chain to the remote chain.
  /// The tokens on the remote chain will then be sourced from the remoteToken address.
  /// The amount to be sent must be approved by the caller beforehand on the localToken contract.
  /// The caller must provide the bridging fee in native currency, i.e msg.value.
  /// @param localToken The address of the local ERC-20 token.
  /// @param remoteToken The address of the remote ERC-20 token.
  /// @param recipient The address of the recipient on the remote chain.
  /// @param amount The amount of the local token to send.
  /// @param bridgeSpecificPayload The payload of the cross-chain transfer. Bridge-specific.
  function sendERC20(
    address localToken,
    address remoteToken,
    address recipient,
    uint256 amount,
    bytes calldata bridgeSpecificPayload
  ) external payable returns (bytes memory);

  /// @notice Get the bridging fee in native currency. This fee must be provided upon sending tokens via
  /// the sendERC20 function.
  /// @return The bridging fee in native currency.
  function getBridgeFeeInNative() external view returns (uint256);

  /// @notice Finalize the withdrawal of a cross-chain transfer.
  /// Not all implementations will finalize a transfer in a single call to this function.
  /// Optimism, for example, requires a two-step process to finalize a transfer. The first
  /// step requires proving the withdrawal that occurred on L2 on L1. The second step is then
  /// the finalization, whereby funds become available to the recipient. So, in that particular
  /// scenario, `false` is returned from `finalizeWithdrawERC20` when the first step is completed,
  /// and `true` is returned when the second step is completed.
  /// @param remoteSender The address of the sender on the remote chain.
  /// @param localReceiver The address of the receiver on the local chain.
  /// @param bridgeSpecificPayload The payload of the cross-chain transfer, bridge-specific, i.e a proof of some kind.
  /// @return true iff the funds are available, false otherwise.
  function finalizeWithdrawERC20(
    address remoteSender,
    address localReceiver,
    bytes calldata bridgeSpecificPayload
  ) external returns (bool);
}
