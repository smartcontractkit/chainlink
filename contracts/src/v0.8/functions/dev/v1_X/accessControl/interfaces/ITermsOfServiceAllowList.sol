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

  /// @notice Get details about the total number of allowed senders
  /// @return count - total number of allowed senders in the system
  function getAllowedSendersCount() external view returns (uint64);

  /// @notice Retrieve a list of allowed senders using an inclusive range
  /// @dev WARNING: getAllowedSendersInRange uses EnumerableSet .length() and .at() methods to iterate over the list
  /// without the need for an extra mapping. These method can not guarantee the ordering when new elements are added.
  /// Evaluate if eventual consistency will satisfy your usecase before using it.
  /// @param allowedSenderIdxStart - index of the allowed sender to start the range at
  /// @param allowedSenderIdxEnd - index of the allowed sender to end the range at
  /// @return allowedSenders - allowed addresses in the range provided
  function getAllowedSendersInRange(
    uint64 allowedSenderIdxStart,
    uint64 allowedSenderIdxEnd
  ) external view returns (address[] memory allowedSenders);

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

  /// @notice Get details about the total number of blocked senders
  /// @return count - total number of blocked senders in the system
  function getBlockedSendersCount() external view returns (uint64);

  /// @notice Retrieve a list of blocked senders using an inclusive range
  /// @dev WARNING: getBlockedSendersInRange uses EnumerableSet .length() and .at() methods to iterate over the list
  /// without the need for an extra mapping. These method can not guarantee the ordering when new elements are added.
  /// Evaluate if eventual consistency will satisfy your usecase before using it.
  /// @param blockedSenderIdxStart - index of the blocked sender to start the range at
  /// @param blockedSenderIdxEnd - index of the blocked sender to end the range at
  /// @return blockedSenders - blocked addresses in the range provided
  function getBlockedSendersInRange(
    uint64 blockedSenderIdxStart,
    uint64 blockedSenderIdxEnd
  ) external view returns (address[] memory blockedSenders);

  /// @notice Enables migrating any previously allowed senders to the new contract
  /// @param previousSendersToAdd - List of addresses to migrate. These address must be allowed on the previous ToS contract and not blocked
  function migratePreviouslyAllowedSenders(address[] memory previousSendersToAdd) external;
}

// ================================================================
// |                     Configuration state                      |
// ================================================================
struct TermsOfServiceAllowListConfig {
  bool enabled; // ═════════════╗ When enabled, access will be checked against s_allowedSenders. When disabled, all access will be allowed.
  address signerPublicKey; // ══╝ The key pair that needs to sign the acceptance data
}
