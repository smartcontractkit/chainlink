// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

/// @notice A contract to handle access control of subscription management dependent on signing a Terms of Service
interface ITermsOfServiceAllowList {
  /// @notice Return the message data for the proof given to accept the Terms of Service
  /// @param acceptor - The wallet address that has accepted the Terms of Service on the UI
  /// @param recipient - The recipient address that the acceptor is taking responsibility for
  /// @return Hash of the message data
  function getMessage(address acceptor, address recipient) external pure returns (bytes32);

  /// @notice Check if the address is blocked for usage
  /// @param sender The transaction sender's address
  /// @return True or false
  function isBlockedSender(address sender) external returns (bool);

  /// @notice Get a list of all allowed senders
  /// @dev WARNING: This operation will copy the entire storage to memory, which can be quite expensive. This is designed
  /// to mostly be used by view accessors that are queried without any gas fees. Developers should keep in mind that
  /// this function has an unbounded cost, and using it as part of a state-changing function may render the function
  /// uncallable if the set grows to a point where copying to memory consumes too much gas to fit in a block.
  /// @return addresses - all allowed addresses
  function getAllAllowedSenders() external view returns (address[] memory);

  /// @notice Allows access to the sender based on acceptance of the Terms of Service
  /// @param acceptor - The wallet address that has accepted the Terms of Service on the UI
  /// @param recipient - The recipient address that the acceptor is taking responsibility for
  /// @param r - ECDSA signature r data produced by the Chainlink Functions Subscription UI
  /// @param s - ECDSA signature s produced by the Chainlink Functions Subscription UI
  /// @param v - ECDSA signature v produced by the Chainlink Functions Subscription UI
  function acceptTermsOfService(address acceptor, address recipient, bytes32 r, bytes32 s, uint8 v) external;

  /// @notice Removes a sender's access if already authorized, and disallows re-accepting the Terms of Service
  /// @param sender - Address of the sender to block
  function blockSender(address sender) external;

  /// @notice Re-allows a previously blocked sender to accept the Terms of Service
  /// @param sender - Address of the sender to unblock
  function unblockSender(address sender) external;
}
