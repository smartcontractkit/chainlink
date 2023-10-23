// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import {IEVM2AnyOnRampClient} from "./IEVM2AnyOnRampClient.sol";

import {Internal} from "../libraries/Internal.sol";

interface IEVM2AnyOnRamp is IEVM2AnyOnRampClient {
  /// @notice Gets the next sequence number to be used in the onRamp
  /// @return the next sequence number to be used
  function getExpectedNextSequenceNumber() external view returns (uint64);

  /// @notice Get the next nonce for a given sender
  /// @param sender The sender to get the nonce for
  /// @return nonce The next nonce for the sender
  function getSenderNonce(address sender) external view returns (uint64 nonce);

  /// @notice Adds and removed token pools.
  /// @param removes The tokens and pools to be removed
  /// @param adds The tokens and pools to be added.
  function applyPoolUpdates(Internal.PoolUpdate[] memory removes, Internal.PoolUpdate[] memory adds) external;
}
