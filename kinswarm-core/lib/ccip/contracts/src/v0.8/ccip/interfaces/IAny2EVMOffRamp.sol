// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

interface IAny2EVMOffRamp {
  /// @notice Returns the the current nonce for a receiver.
  /// @param sender The sender address
  /// @return nonce The nonce value belonging to the sender address.
  function getSenderNonce(
    address sender
  ) external view returns (uint64 nonce);
}
