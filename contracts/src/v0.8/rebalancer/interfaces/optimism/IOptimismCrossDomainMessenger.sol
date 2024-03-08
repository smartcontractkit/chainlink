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
}
