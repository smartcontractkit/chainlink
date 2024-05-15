// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import {IEVM2AnyOnRampClient} from "./IEVM2AnyOnRampClient.sol";

interface IEVM2AnyMultiOnRamp is IEVM2AnyOnRampClient {
  /// @notice Gets the next sequence number to be used in the onRamp
  /// @param destChainSelector The destination chain selector
  /// @return the next sequence number to be used
  function getExpectedNextSequenceNumber(uint64 destChainSelector) external view returns (uint64);

  /// @notice Returns the current nonce for a sender
  /// @param destChainSelector The destination chain selector
  /// @param sender The sender address
  /// @return The sender's nonce
  function getSenderNonce(uint64 destChainSelector, address sender) external view returns (uint64);
}
