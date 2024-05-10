// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

interface IAny2EVMMultiOffRamp {
  /// @notice Returns the the current nonce for a receiver.
  /// @param sourceChainSelector The source chain to retrieve the nonce for
  /// @param sender The sender address
  /// @return nonce The nonce value belonging to the sender address.
  function getSenderNonce(uint64 sourceChainSelector, address sender) external view returns (uint64 nonce);
}
