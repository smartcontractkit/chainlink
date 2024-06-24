// SPDX-License-Identifier: MIT
// Copied from https://github.com/ethereum-optimism/optimism/blob/f707883038d527cbf1e9f8ea513fe33255deadbc/packages/contracts-bedrock/src/universal/StandardBridge.sol#L88
pragma solidity ^0.8.0;

interface IOptimismStandardBridge {
  /// @notice Emitted when an ERC20 bridge is finalized on this chain.
  /// @param localToken  Address of the ERC20 on this chain.
  /// @param remoteToken Address of the ERC20 on the remote chain.
  /// @param from        Address of the sender.
  /// @param to          Address of the receiver.
  /// @param amount      Amount of the ERC20 sent.
  /// @param extraData   Extra data sent with the transaction.
  event ERC20BridgeFinalized(
    address indexed localToken,
    address indexed remoteToken,
    address indexed from,
    address to,
    uint256 amount,
    bytes extraData
  );

  /// @notice Finalizes an ERC20 bridge on this chain. Can only be triggered by the other
  ///         StandardBridge contract on the remote chain.
  /// @param _localToken  Address of the ERC20 on this chain.
  /// @param _remoteToken Address of the corresponding token on the remote chain.
  /// @param _from        Address of the sender.
  /// @param _to          Address of the receiver.
  /// @param _amount      Amount of the ERC20 being bridged.
  /// @param _extraData   Extra data to be sent with the transaction. Note that the recipient will
  ///                     not be triggered with this data, but it will be emitted and can be used
  ///                     to identify the transaction.
  function finalizeBridgeERC20(
    address _localToken,
    address _remoteToken,
    address _from,
    address _to,
    uint256 _amount,
    bytes calldata _extraData
  ) external;
}
