// SPDX-License-Identifier: MIT
// Copied from https://github.com/ethereum-optimism/optimism/blob/f707883038d527cbf1e9f8ea513fe33255deadbc/packages/contracts-bedrock/src/universal/CrossDomainMessenger.sol#L153
pragma solidity ^0.8.0;

interface IOptimismCrossDomainMessenger {
  /// @notice Emitted whenever a message is sent to the other chain.
  /// @param target       Address of the recipient of the message.
  /// @param sender       Address of the sender of the message.
  /// @param message      Message to trigger the recipient address with.
  /// @param messageNonce Unique nonce attached to the message.
  /// @param gasLimit     Minimum gas limit that the message can be executed with.
  event SentMessage(address indexed target, address sender, bytes message, uint256 messageNonce, uint256 gasLimit);

  /// @notice Relays a message that was sent by the other CrossDomainMessenger contract. Can only
  ///         be executed via cross-chain call from the other messenger OR if the message was
  ///         already received once and is currently being replayed.
  /// @param _nonce       Nonce of the message being relayed.
  /// @param _sender      Address of the user who sent the message.
  /// @param _target      Address that the message is targeted at.
  /// @param _value       ETH value to send with the message.
  /// @param _minGasLimit Minimum amount of gas that the message can be executed with.
  /// @param _message     Message to send to the target.
  function relayMessage(
    uint256 _nonce,
    address _sender,
    address _target,
    uint256 _value,
    uint256 _minGasLimit,
    bytes calldata _message
  ) external payable;
}
